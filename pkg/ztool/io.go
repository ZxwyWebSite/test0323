// IO扩展操作 (beta)

package ztool

import (
	"io"
)

// 异步写入器列表
type Fbj_MultiWriterConf struct {
	IgnoreErr bool // 忽略错误继续写入下一个Writer
	ASync     bool // 异步操作，通常搭配IgnoreErr使用
}
type multiWriter struct {
	writers []io.Writer
	conf    Fbj_MultiWriterConf
	// mu      *sync.Mutex
}

func (t *multiWriter) WaitWrite(p []byte) (n int, err error) {
	for _, w := range t.writers {
		n, err = w.Write(p)
		if err != nil && !t.conf.IgnoreErr {
			return
		} else {
			if n != len(p) {
				err = io.ErrShortWrite
				return
			}
		}
	}
	return len(p), nil
}

func (t *multiWriter) Write(p []byte) (n int, err error) {
	if t.conf.ASync {
		go func() {
			// t.mu.Lock()
			t.WaitWrite(p)
			// t.mu.Unlock()
		}()
		return len(p), nil
	}
	return t.WaitWrite(p)
}

func Fbj_MultiWriter(conf Fbj_MultiWriterConf, writers ...io.Writer) io.Writer {
	allWriters := make([]io.Writer, 0, len(writers))
	for _, w := range writers {
		if mw, ok := w.(*multiWriter); ok {
			allWriters = append(allWriters, mw.writers...)
		} else {
			allWriters = append(allWriters, w)
		}
	}
	return &multiWriter{writers: allWriters, conf: conf}
}
