package localcache

import (
	"errors"
	"lx-source/src/caches"
	"lx-source/src/env"
	"lx-source/src/middleware/util"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/ZxwyWebSite/ztool"
)

type Cache struct {
	Path  string // 本地缓存目录 cache
	Bind  string // Api地址，用于生成外链 http://192.168.10.22:1011/
	state bool   // 激活状态
}

// var loger = env.Loger.NewGroup(`Caches`) //caches.Loger.AppGroup(`local`)

func (c *Cache) getLink(q *caches.Query) (furl string) {
	// fmt.Printf("%#v\n", q.Request)
	if env.Config.Cache.Local_Auto {
		// 注：此方式无法确定是否支持HTTPS，暂时默认HTTP，让Nginx重定向
		furl = ztool.Str_FastConcat(util.GetPath(q.Request, `link/`), `file/`, q.Query())
	} else {
		furl = ztool.Str_FastConcat(c.Bind, `/file/`, q.Query())
	}
	return
}

func (c *Cache) Get(q *caches.Query) string {
	// 加一层缓存，减少重复检测文件造成的性能损耗 (缓存已移至Router)
	// if _, ok := env.Cache.Get(q.Query()); !ok {
	if _, e := os.Stat(ztool.Str_FastConcat(c.Path, `/`, q.Query())); e != nil {
		return ``
	}
	// env.Cache.Set(q.Query(), struct{}{}, 3600)
	// }
	return c.getLink(q)
	// fpath := filepath.Join(c.Path, q.Source, q.MusicID, q.Quality)
	// if _, e := os.Stat(fpath); e != nil {
	// 	return ``
	// }
	// return c.getLink(fpath)
}

func (c *Cache) Set(q *caches.Query, l string) string {
	fpath := ztool.Str_FastConcat(c.Path, `/`, q.Query())
	// if env.Config.Main.FFConv && q.Source == `kg` { // ztool.Chk_IsMatch(q.Source, `kg`)
	// 	err := ztool.Fbj_MkdirAll(fpath, 0644)
	// 	if err != nil {
	// 		loger.Error(`DownloadFile_Mkdir: %v`, err)
	// 		return ``
	// 	}
	// 	out, err := ztool.Cmd_aWaitExec(ztool.Str_FastConcat(`ffmpeg -i "`, l, `" -vn`, ` -c:a copy`, ` "`, fpath, `"`))
	// 	if err != nil {
	// 		loger.Error(`DownloadFile_Exec: %v, Output: %v`, err, out)
	// 		return ``
	// 	}
	// 	loger.Debug(`FFMpeg_Out: %v`, out)
	// } else {
	for i := 0; true; i++ {
		err := ztool.Net_DownloadFile(l, fpath, nil)
		if err == nil {
			break
		}
		caches.Loger.Error(`DownloadFile: %v, Retry: %v`, err, i)
		if i == 1 || !strings.Contains(err.Error(), `context deadline exceeded`) {
			if err := os.Remove(fpath); err != nil {
				caches.Loger.Error(`RemoveFile: %s`, err)
			}
			return ``
		}
		time.Sleep(time.Second)
	}
	// }
	// env.Cache.Set(q.Query(), struct{}{}, 3600)
	return c.getLink(q)
	// fpath := filepath.Join(c.Path, q.String)
	// os.MkdirAll(filepath.Dir(fpath), fs.ModePerm)
	// g := c.Loger.NewGroup(`localcache`)
	// ret, err := ztool.Net_HttpReq(http.MethodGet, l, nil, nil, nil)
	// if err != nil {
	// 	g.Warn(`HttpReq: %s`, err)
	// 	return ``
	// }
	// if err := os.WriteFile(fpath, ret, fs.ModePerm); err != nil {
	// 	g.Warn(`WriteFile: %s`, err)
	// 	return ``
	// }
	// return c.getLink(fpath)
}

func (c *Cache) Stat() bool {
	return c.state
}

func (c *Cache) Init() error {
	if c.Bind == `` {
		return errors.New(`请输入Api地址以生成外链`)
	} else {
		ubj, err := url.Parse(c.Bind)
		if err != nil {
			return err
		}
		ubj.Path = strings.TrimSuffix(ubj.Path, `/`)
		c.Bind = ubj.String()
	}
	c.state = true
	return nil
}

// func New(path, addr string, loger *logs.Logger) *Cache {
// 	return &Cache{
// 		Path:  path,
// 		Addr:  addr,
// 		Loger: loger,
// 		emsg:  cache.ErrNotInited,
// 	}
// }
