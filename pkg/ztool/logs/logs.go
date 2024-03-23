// 旧版本
package logs

import (
	"fmt"
	"os"
	"time"
)

type Logser struct {
	name  string
	group []string
}

var Level = 2 // 日志等级 0: off, 1: info, 2: warn, 3:debug

// 新建日志对象
func New(n string) *Logser {
	return &Logser{
		name: n,
	}
}

// 新建日志组
func (l *Logser) Group(n string) *Logser {
	return &Logser{
		name:  l.name,
		group: append(l.group, n),
	}
}

// 添加组项
func (l *Logser) AddGroup(n string) *Logser {
	l.group = append(l.group, n)
	return l
}

func (l *Logser) print(t, format string, a ...any) {
	var groups string
	for _, v := range l.group {
		groups += `[` + v + `] `
	}
	fmt.Printf(`[`+l.name+`] `+time.Now().Format(time.DateTime)+` [`+t+`] `+groups+format+"\n", a...)
}

func (l *Logser) Info(format string, a ...any) {
	if Level >= 1 {
		l.print(`INFO`, format, a...)
	}
}

func (l *Logser) Warn(format string, a ...any) {
	if Level >= 2 {
		l.print(`WARN`, format, a...)
	}
}
func (l *Logser) Debug(format string, a ...any) {
	if Level >= 3 {
		l.print(`DEBU`, format, a...)
	}
}

func (l *Logser) Fatal(format string, a ...any) {
	l.print(`FATA`, format, a...)
	os.Exit(1)
}
