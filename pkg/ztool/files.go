// 文件对象化操作; 分类: Fbj_(文件)

package ztool

import (
	"io/fs"
	"os"
	"path/filepath"
)

// const (
// 	Fbj_PermR = 4
// 	Fbj_PermW = 2
// 	Fbj_PermX = 1
// )

var (
	// 默认权限
	/*
	  0 | 所有者 | 用户组 | 公共

	 读取 (4): f t t t f f f t
	 写入 (2): f f t t t f t f
	 执行 (1): f f f t t t f t
	 权限 (+): 0 4 6 7 3 1 2 5

	 R: × × × × √ √ √ √
	 W: × × √ √ × × √ √
	 X: × √ × √ × √ × √
	 P: 0 1 2 3 4 5 6 7
	*/
	Fbj_DefPerm fs.FileMode = fs.ModePerm //0644
)

// 判断文件是否存在 <文件路径(相对或绝对皆可)> <存在?>
func Fbj_IsExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// 为文件创建上级目录
func Fbj_MkdirAll(path string, perm fs.FileMode) error {
	if dir := filepath.Dir(path); !Fbj_IsExists(dir) {
		return os.MkdirAll(dir, perm)
	}
	return nil
}
func Fbj_Mkdir(path string) error {
	return os.MkdirAll(path, Fbj_DefPerm)
}

// 打开文件 <路径, 标识, 权限> <*文件, 错误>
/*
 用法同 os.OpenFile
 调用前会检测文件目录是否存在，否则补充创建
*/
func Fbj_OpenFile(path string, flag int, perm fs.FileMode) (*os.File, error) {
	if err := Fbj_MkdirAll(path, perm); err != nil {
		return nil, err
	}
	// if dir := filepath.Dir(path); !Fbj_IsExists(dir) {
	// 	if err := os.MkdirAll(dir, perm); err != nil {
	// 		return nil, err
	// 	}
	// }
	return os.OpenFile(path, flag, perm)
}

// 快速创建文件 (同 os.Create)
func Fbj_CreatFile(path string) (*os.File, error) {
	// basePath := filepath.Dir(path)
	// if !Fbj_IsExists(basePath) {
	// 	err := os.MkdirAll(basePath, 0644)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }
	// return os.Create(path)
	return Fbj_OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, Fbj_DefPerm)
}

// func Fbj_Open(name string) (*os.File, error) {
// 	return os.OpenFile(name, os.O_RDONLY, 0)
// }

// 读取文件内容 (同 os.ReadFile)
// func Fbj_ReadFile(path string) ([]byte, error) {
// 	// Fbj_MkdirAll(path, Fbj_DefPerm)
// 	// os.Open(path)
// 	f, err := Fbj_OpenFile(path, os.O_RDONLY, Fbj_DefPerm)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer f.Close()
// 	var size int
// 	if info, err := f.Stat(); err == nil {
// 		size64 := info.Size()
// 		if int64(int(size64)) == size64 {
// 			size = int(size64)
// 		}
// 	}
// 	size++
// 	if size < 512 {
// 		size = 512
// 	}
// 	data := make([]byte, 0, size)
// 	for {
// 		if len(data) >= cap(data) {
// 			d := append(data[:cap(data)], 0)
// 			data = d[:len(data)]
// 		}
// 		n, err := f.Read(data[len(data):cap(data)])
// 		data = data[:len(data)+n]
// 		if err != nil {
// 			if err == io.EOF {
// 				err = nil
// 			}
// 			return data, err
// 		}
// 	}
// }

// 写入文件内容 (同 os.WriteFile)
func Fbj_WriteFile(name string, data []byte, perm fs.FileMode) error {
	f, err := Fbj_OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err1 := f.Close(); err1 != nil && err == nil {
		err = err1
	}
	return err
}
func Fbj_SaveFile(name string, data []byte) error {
	return Fbj_WriteFile(name, data, Fbj_DefPerm)
}

// 获取文件列表
func Fbj_MkList(path string) ([]string, error) {
	dir, err := os.ReadDir(path)
	var out []string
	if err == nil {
		for i, r := 0, len(dir); i < r; i++ {
			if !dir[i].IsDir() {
				out = append(out, dir[i].Name())
			}
		}
	}
	return out, err
}
