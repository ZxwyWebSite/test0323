// 文件压缩 (beta); 分类：Pak_(打包)

package ztool

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// 压缩配置
type Pak_ZipConfig struct {
	UnPath bool // 不使用路径样式
}

// 封装逻辑
// func zipfile() {
// }

// Zip压缩文件
/*
 传入文件地址, 输出地址, 压缩配置, 尝试压缩并返回发生的错误
 e.g. this(`path/to/example.txt`, ``, nil) >> Zip{"path/to/example.txt"}
 e.g. this(`path/to/example.txt`, ``, { Unpath: true }) >> Zip{"example.txt"}
 注：传入的应是文件而不是目录，可添加Unpath参数忽略文件路径，不填输出地址默认文件名+zip
*/
func Pak_ZipFile(src_file, dst_zip_name string, config Pak_ZipConfig) error {
	if !Fbj_IsExists(src_file) {
		return errors.New(`src file not exist`) // 待压缩文件不存在 no such file or directory
	}
	if dst_zip_name == `` {
		dst_zip_name = Str_FastConcat(src_file, `.zip`)
	}
	os.RemoveAll(dst_zip_name)
	zipfile, err := Fbj_CreatFile(dst_zip_name)
	if err != nil {
		return err
	}
	defer zipfile.Close()
	archive := zip.NewWriter(zipfile)
	defer archive.Close()
	info, _ := os.Lstat(src_file)
	header, _ := zip.FileInfoHeader(info)
	if info.IsDir() {
		return nil // 压缩文件函数不支持压缩目录
	} else {
		header.Method = zip.Deflate
	}
	if config.UnPath {
		header.Name = filepath.Base(src_file)
	}
	writer, err := archive.CreateHeader(header)
	if err != nil {
		return err
	}
	file, _ := os.Open(src_file)
	defer file.Close()
	_, err = io.Copy(writer, file)
	return err
}

// Zip压缩目录
/*
 传入目录地址, 输出文件名, 返回错误信息
 e.g. see Zip_PakFile()
 注：未完善, 慎用
 TODO：封装压缩逻辑
*/
func Pak_ZipDir(src_dir string, zip_file_name string) error {
	var handler Err_HandleList
	// 快捷：输出文件名为空自动补充
	if zip_file_name == `` {
		zip_file_name = Str_FastConcat(src_dir, `.zip`)
	}
	// 预防：旧文件无法覆盖
	os.RemoveAll(zip_file_name)

	// 创建：zip文件
	zipfile, _ := Fbj_CreatFile(zip_file_name)
	defer zipfile.Close()

	// 打开：zip文件
	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	// 遍历路径信息
	handler.Check(filepath.Walk(src_dir, func(path string, info os.FileInfo, e error) error {
		if e != nil {
			return e
		}
		// 如果是源路径，提前进行下一个遍历
		if path == src_dir {
			return nil
		}
		// 获取：文件头信息
		header, _ := zip.FileInfoHeader(info)
		header.Name = strings.TrimPrefix(path, src_dir+`/`)

		// 判断：文件是不是文件夹
		if info.IsDir() {
			header.Name += `/`
		} else {
			// 设置：zip的文件压缩算法
			header.Method = zip.Deflate
		}

		// 创建：压缩包头部信息
		writer, _ := archive.CreateHeader(header)
		if !info.IsDir() {
			file, _ := os.Open(path)
			defer file.Close()
			io.Copy(writer, file)
		}
		return nil
	}))

	// 返回错误
	if res := handler.Result(); res != nil {
		return res.Format()
	}
	return nil
}
