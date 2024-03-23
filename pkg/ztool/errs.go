// 错误处理 (beta); 分类: Err_(错误)

package ztool

import (
	"errors"
	"strconv"
)

// 简单错误批量处理模型
type Err_HandleList struct {
	Num      int
	Err      error
	Callback func(*Err_HandleList)
}

var (
	Err_EsContinue = errors.New(`ErrContinue`) // 忽略错误继续运行
)

func Err_NewDefHandleList() Err_HandleList {
	return Err_HandleList{
		Callback: func(l *Err_HandleList) {
			Cmd_FastPrintln(Str_FastConcat(`任务编号 `, l.NumStr(), ` 执行结束`))
		},
	}
}

func (l *Err_HandleList) sync(e error) {
	l.Err = e
	if l.Callback != nil {
		l.Callback(l)
	}
	l.Num++
}

func (l *Err_HandleList) NumStr() string {
	return strconv.Itoa(l.Num)
}

func (l *Err_HandleList) Do(f func() error) {
	if l.Err != nil {
		return
	}
	l.sync(f())
}
func (l *Err_HandleList) WithValue(f func() (any, error)) any {
	if l.Err != nil {
		return nil
	}
	v, err := f()
	l.sync(err)
	return v
}

// 查错 (不建议此方法调用，无法提前终止函数执行)
func (l *Err_HandleList) Check(e error) {
	if l.Err != nil {
		return
	}
	l.sync(e)
}

// 获取执行结果，
func (l *Err_HandleList) Result() *Err_HandleListResult {
	if l.Err != nil {
		return &Err_HandleListResult{Num: l.NumStr(), Err: l.Err.Error()}
		// return errors.New(FastStrConcat(`Err_HandleListIndex`, strconv.Itoa(l.Num), `: `, l.Err.Error()))
	}
	return nil
}

type Err_HandleListResult struct {
	Num string
	Err string
}

// 生成带任务编号的错误信息
func (e *Err_HandleListResult) Errors() string {
	return Str_FastConcat(`Err_HandleListIndex`, e.Num, `: `, e.Err)
}

// 格式化默认错误信息
func (e *Err_HandleListResult) Format() error {
	return errors.New(Str_FastConcat(`执行到第 `, e.Num, ` 次时发生错误：`, e.Err))
}

// type Err_DefFunc func() error

// 退出延迟执行模型
type Err_DeferList struct {
	tasks []func()
}

// 添加任务
func (l *Err_DeferList) Add(f func()) {
	l.tasks = append(l.tasks, f)
}

// 执行任务
func (l *Err_DeferList) Do() {
	for i, r := 0, len(l.tasks); i < r; i++ {
		l.tasks[i]()
	}
	l = nil
}
