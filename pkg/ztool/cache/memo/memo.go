package memo

import (
	"encoding/gob"
	"fmt"
	"os"
	"sync"

	// "sync/atomic"
	"time"
	// "unsafe"

	"github.com/ZxwyWebSite/ztool"
	"github.com/ZxwyWebSite/ztool/logs"
)

// MemoStore 内存存储驱动
type MemoStore struct {
	Store *sync.Map
	Loger *logs.Logger
	// writed uint32
	// lastGC time.Time
}

// item 存储的对象
type itemWithTTL struct {
	Expires int64
	Value   interface{}
}

func newItem(value interface{}, expires int) itemWithTTL {
	expires64 := int64(expires) // *(*int64)(unsafe.Pointer(&expires))
	if expires > 0 {
		expires64 += time.Now().Unix() // expires64 = time.Now().Unix() + expires64
	}
	return itemWithTTL{
		Value:   value,
		Expires: expires64,
	}
}

// store.getValue 从itemWithTTL中取值
func (store *MemoStore) getValue(key any) (interface{}, bool) {
	item, ok := store.Store.Load(key)
	if !ok {
		return nil, ok
	}

	var itemObj itemWithTTL
	if itemObj, ok = item.(itemWithTTL); !ok {
		return item, true
	}

	if itemObj.Expires > 0 && itemObj.Expires < time.Now().Unix() {
		// store.Store.Delete(key) // 自动删除过期项
		return nil, false
	}

	return itemObj.Value, ok
}

// GarbageCollect 回收已过期的缓存
func (store *MemoStore) GarbageCollect() {
	now := time.Now().Unix()
	store.Store.Range(func(key, value interface{}) bool {
		if item, ok := value.(itemWithTTL); ok {
			if item.Expires > 0 && item.Expires < now {
				if store.Loger != nil {
					store.Loger.Debug(`清除已过期数据 [%q]`, key.(string))
				}
				store.Store.Delete(key)
			}
		}
		return true
	})
}

// 设置自动GC <检测间隔>
/*
 参考Redis和GolangGC策略
 每s秒检测一次
 距上次运行超过2分钟强制运行一次
 写入超过30次则继续运行
*/
func (store *MemoStore) AutoGC(s int) {
	// store.lastGC = time.Now()
	// doGC := func() {
	// 	store.GarbageCollect()               // 执行GC
	// 	atomic.StoreUint32(&store.writed, 0) // 清空写入计数
	// 	store.lastGC = time.Now()            // 更新时间戳
	// }
	go func() {
		for {
			time.Sleep(time.Second * time.Duration(s))
			store.GarbageCollect()
			// last := store.lastGC.Add(time.Minute)
			// if !last.Before(time.Now()) {
			// 	// 距上次GC未满2分钟
			// 	wd := atomic.LoadUint32(&store.writed)
			// 	if wd < 5 {
			// 		// 写入不超过30次
			// 		continue
			// 	}
			// }
			// doGC()
		}
	}()
}

// NewMemoStore 新建内存存储
func NewMemoStore() *MemoStore {
	return &MemoStore{
		Store: &sync.Map{},
	}
}
func NewMemoStoreConf(loger *logs.Logger, gc int) *MemoStore {
	out := &MemoStore{
		Store: &sync.Map{},
		Loger: loger.NewGroup(`MemoCache`),
	}
	if gc > 0 {
		out.AutoGC(gc)
	}
	return out
}

// Set 存储值
func (store *MemoStore) Set(key string, value interface{}, ttl int) error {
	store.Store.Store(key, newItem(value, ttl))
	// atomic.AddUint32(&store.writed, 1)
	// if store.Loger != nil {
	// 	store.Loger.Debug(`SetKV: %v, %v`, key, value)
	// }
	return nil
}

// Get 取值
func (store *MemoStore) Get(key string) (interface{}, bool) {
	// if store.Loger != nil {
	// 	store.Loger.Debug(`GetKV: %v`, key)
	// }
	return store.getValue(key)
}

// Gets 批量取值
func (store *MemoStore) Gets(keys []string, prefix string) (map[string]interface{}, []string) {
	var res = make(map[string]interface{})
	var notFound = make([]string, 0, len(keys))

	for _, key := range keys {
		if value, ok := store.getValue(prefix + key); ok {
			res[key] = value
		} else {
			notFound = append(notFound, key)
		}
	}

	return res, notFound
}

// Sets 批量设置值
func (store *MemoStore) Sets(values map[string]interface{}, prefix string) error {
	for key, value := range values {
		store.Store.Store(prefix+key, value)
	}
	return nil
}

// Delete 批量删除值
func (store *MemoStore) Delete(keys []string, prefix string) error {
	for _, key := range keys {
		store.Store.Delete(prefix + key)
	}
	return nil
}

// Persist 将内存存储持久写入缓存
func (store *MemoStore) Persist(path string) error {
	file, err := ztool.Fbj_CreatFile(path)
	if err != nil {
		return err
	}

	defer func() {
		file.Chmod(0644)
		file.Close()
	}()

	persisted := make(map[string]itemWithTTL)
	now := time.Now().Unix()
	store.Store.Range(func(key, value interface{}) bool {
		itemObj, ok := value.(itemWithTTL)
		if !ok {
			itemObj.Value = value
		} else if itemObj.Expires > 0 && itemObj.Expires < now {
			return true
		}
		persisted[key.(string)] = itemObj
		return true
	})

	enc := gob.NewEncoder(file)
	err = enc.Encode(persisted)
	if err != nil {
		return fmt.Errorf(`缓存序列化失败: %s`, err)
	}

	return nil
}

// func (store *MemoStore) Persist(path string) error {
// 	persisted := make(map[string]itemWithTTL)
// 	now := time.Now().Unix()
// 	store.Store.Range(func(key, value interface{}) bool {
// 		itemObj, ok := value.(itemWithTTL)
// 		if !ok {
// 			itemObj.Value = value
// 		} else if itemObj.Expires > 0 && itemObj.Expires < now {
// 			return true
// 		}
// 		persisted[key.(string)] = itemObj
// 		return true
// 	})

// 	buf := new(bytes.Buffer)
// 	enc := gob.NewEncoder(buf)
// 	err := enc.Encode(persisted)
// 	if err != nil {
// 		return fmt.Errorf(`缓存序列化失败: %s`, err)
// 	}

// 	if ztool.Fbj_IsExists(path) {
// 		old, err := os.ReadFile(path)
// 		if err == nil {
// 			if bytes.Equal(old, buf.Bytes()) {
// 				return nil
// 			}
// 		}
// 	}
// 	file, err := ztool.Fbj_CreatFile(path)
// 	if err != nil {
// 		return fmt.Errorf(`缓存文件创建失败: %s`, err)
// 	}
// 	defer func() {
// 		file.Chmod(0644)
// 		file.Close()
// 	}()
// 	_, err = io.Copy(file, buf)

//		return err
//	}
func (store *MemoStore) MustPersist(path string) error {
	err := store.Persist(path)
	if err != nil && store.Loger != nil {
		store.Loger.Error(err.Error())
	}
	return err
}

// Restore 从磁盘文件恢复内存缓存
func (store *MemoStore) Restore(path string) error {
	if !ztool.Fbj_IsExists(path) {
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf(`读取缓存文件失败: %s`, err)
	}
	defer file.Close()

	// defer func() {
	// 	file.Close()
	// 	os.Remove(path)
	// }()

	items := make(map[string]itemWithTTL)
	dec := gob.NewDecoder(file)
	if err := dec.Decode(&items); err != nil {
		return fmt.Errorf(`缓存反序列化失败: %s`, err)
	}

	loaded := 0
	now := time.Now().Unix()
	for k, v := range items {
		if v.Expires > 0 && v.Expires < now && store.Loger != nil {
			store.Loger.Debug(`缓存数据 %q 已过期`, k)
		} else {
			loaded++
			store.Store.Store(k, v)
		}
	}

	if loaded > 0 && store.Loger != nil {
		store.Loger.Info(`已恢复 %d 个持久化缓存项目`, loaded) // 已将 %d 个项目恢复到内存缓存中
	}
	return nil
}
func (store *MemoStore) MustRestore(path string) error {
	err := store.Restore(path)
	if err != nil && store.Loger != nil {
		store.Loger.Error(err.Error())
	}
	return err
}
