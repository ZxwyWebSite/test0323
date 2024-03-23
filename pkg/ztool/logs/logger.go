package logs

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/ZxwyWebSite/ztool"
	color "github.com/ZxwyWebSite/ztool/mod/v1.16.0-color"
	"github.com/mattn/go-colorable"
	// "github.com/fatih/color"
)

const (
	LevelNone = iota
	LevelInfo
	LevelWarn
	LevelDebu
)

const (
	l_Info = `Info`  // 提示
	l_Warn = `Warn`  // 警告
	l_Erro = `Error` // 错误
	l_Fata = `Fatal` // 致命错误
	l_Debu = `Debug` // 调试
	l_Pani = `Panic` // 惊慌失措

	// l_Main = `Main`  // 主分区
	// l_Grop = `Group` // 组分区
)

var (
	// 全局日志等级 0: off, 1: info, 2: warn, 3:debug (报错不受此限制)
	Levell = LevelWarn
	// 全局互斥锁
	mu = new(sync.Mutex)
	// 默认日志对象 (需先使用DefLoger函数创建)
	Main    = `Logs`
	Default = &Logger{name: &Main, level: &Levell, group: []string{`Main`}, output: &color.Output}
	// 全局打印输出接口
	// Writer = color.Output
	// Color对象池 (*color.Color)
	color_pool = &sync.Pool{New: func() any { return color.New() }}
	loger_pool = &sync.Pool{New: func() any { return createLgr() }}
)

// Logger 日志对象
type Logger struct {
	level *int
	// mu    sync.Mutex
	name   *string
	group  []string
	output *io.Writer
	deferr func()
}

func (l *Logger) SetOutput(w io.Writer) {
	// color.FastSetWritter(w)
	// color.Output = w
	l.output = &w
}
func (l *Logger) SetOutputSync(w io.Writer, sync bool) {
	color.FastSetWritterFunc(func(o io.Writer) io.Writer {
		return ztool.Fbj_MultiWriter(ztool.Fbj_MultiWriterConf{
			IgnoreErr: true, ASync: sync,
		}, o, w)
	})
}

// 设置输出到文件 <文件路径,是否输出到终端> <文件对象,延迟执行函数,错误信息>
func (l *Logger) SetOutFile(path string, stdout bool) (*os.File, func() error, error) {
	f, e := ztool.Fbj_OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if e != nil {
		return nil, nil, e
	}
	w := bufio.NewWriter(f) // bufio.NewWriterSize(f, 1024)
	t := ztool.Str_NumConcat(`=`, 20)
	ztool.Cmd_FastFprint(w, ztool.Str_FastConcat("\n", t, ` [`, time.Now().Format(time.RFC3339), `] `, t, "\n"))
	if stdout {
		*l.output = io.MultiWriter(colorable.NewColorableStdout(), colorable.NewNonColorable(w))
	} else {
		*l.output = colorable.NewNonColorable(w)
	}
	// l.deferr = func() { w.Flush() /*; f.Close()*/ }
	return f, w.Flush, nil
}
func (l *Logger) GetOutput() io.Writer {
	return *l.output
}

// 日志颜色
var colors = map[string]func(s string) string{
	l_Info: color.New(color.FgCyan).Add(color.Bold).FastWrapFunc(),
	l_Warn: color.New(color.FgYellow).Add(color.Bold).FastWrapFunc(),
	l_Erro: color.New(color.FgRed).Add(color.Bold).FastWrapFunc(),
	l_Fata: color.New(color.BgRed).Add(color.Bold).FastWrapFunc(),
	l_Debu: color.New(color.FgWhite).Add(color.Bold).FastWrapFunc(),
	l_Pani: color.New(color.BgRed).Add(color.Bold).FastWrapFunc(),
}

// var colors = map[string]func(a ...interface{}) string{
// 	l_Info: color.New(color.FgCyan).Add(color.Bold).SprintFunc(),
// 	l_Warn: color.New(color.FgYellow).Add(color.Bold).SprintFunc(),
// 	l_Erro: color.New(color.FgRed).Add(color.Bold).SprintFunc(),
// 	l_Fata: color.New(color.BgRed).Add(color.Bold).SprintFunc(),
// 	l_Debu: color.New(color.FgWhite).Add(color.Bold).SprintFunc(),
// 	l_Pani: color.New(color.BgRed).Add(color.Bold).SprintFunc(),
// 	// l_Main: color.New(color.FgGreen).Add(color.Bold).SprintFunc(),
// 	// l_Grop: color.New(color.FgMagenta).Add(color.Bold).SprintFunc(),
// }

// 不同级别前缀与时间的间隔，保持宽度一致
// var spaces = map[string]string{
// 	l_Info: ` `,
// 	l_Warn: ` `,
// 	l_Erro: ``,
// 	l_Fata: ``,
// 	l_Debu: ``,
// 	l_Pani: ``,
// }

// 获取前缀间隔的另一种实现方式
func getSpaces(s string) string {
	if s == l_Info || s == l_Warn {
		return ` `
	}
	return ``
}

// Println 打印
func (ll *Logger) println(prefix string, msg string) {
	// TODO Release时去掉
	// color.NoColor = false

	// TODO 使用对象池
	c := color_pool.Get().(*color.Color) // c := color.New()
	defer color_pool.Put(c)

	// TODO 使用全局互斥锁
	mu.Lock()         // ll.mu.Lock()
	defer mu.Unlock() // defer ll.mu.Unlock()

	// TODO 合并Printf参数, 最后仅调用一次fmt以提高性能
	c.FastFprint(*ll.output, ztool.Str_FastConcat(
		// [Main][Type][Time][Group]Logs
		color.New(color.FgGreen).Add(color.Bold).FastWrap(ztool.Str_FastConcat(`[`, *ll.name, `]`)), ` `,
		colors[prefix](ztool.Str_FastConcat(`[`, prefix, `]`)), getSpaces(prefix), ` `,
		color.New(color.FgBlue).Add(color.Bold).FastWrap(time.Now().Format(time.DateTime)), ` `,
		// color.New(color.FgBlue).Add(color.Bold).FastWrap(ztool.Str_FastConcat(`[`, time.Now().Format(time.DateTime), `]`)), ` `,
		color.New(color.FgMagenta).Add(color.Bold).FastWrap(ztool.Str_JoinPreSuffix(ll.group, `[`, `] `)),
		msg, "\n",
	))
	// _, _ = c.Printf(
	// 	"%s %s%s %s %s%s\n", // [Main] [Type] yyyy-mm-dd hh:mm:ss [Groups] Logs
	// 	color.New(color.FgGreen).Add(color.Bold).SprintFunc()(`[`+*ll.name+`]`),
	// 	colors[prefix](`[`+prefix+`]`), getSpaces(prefix),
	// 	color.New(color.FgBlue).Add(color.Bold).SprintFunc()(time.Now().Format(time.DateTime)),
	// 	color.New(color.FgMagenta).Add(color.Bold).SprintFunc()(ztool.Str_JoinPreSuffix(ll.group, `[`, `] `)),
	// 	msg,
	// )
}

// Panic 极端错误
func (ll *Logger) Panic(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	if *ll.level > LevelNone {
		ll.println(l_Pani, msg)
	}
	panic(msg)
}

// Error 错误
func (ll *Logger) Error(format string, v ...interface{}) *Logger {
	if *ll.level > LevelNone {
		ll.println(l_Erro, fmt.Sprintf(format, v...))
	}
	return ll
}

// Warning 警告
func (ll *Logger) Warn(format string, v ...interface{}) *Logger {
	if *ll.level >= LevelWarn {
		ll.println(l_Warn, fmt.Sprintf(format, v...))
	}
	return ll
}

// Info 信息
func (ll *Logger) Info(format string, v ...interface{}) *Logger {
	if *ll.level >= LevelInfo {
		ll.println(l_Info, fmt.Sprintf(format, v...))
	}
	return ll
}

// Debug 校验
func (ll *Logger) Debug(format string, v ...interface{}) *Logger {
	if *ll.level >= LevelDebu {
		ll.println(l_Debu, fmt.Sprintf(format, v...))
	}
	return ll
}

// Fatal 致命错误
func (ll *Logger) Fatal(format string, v ...interface{}) {
	if *ll.level > LevelNone {
		ll.println(l_Fata, fmt.Sprintf(format, v...))
	}
	ll.deferr()
	if runtime.GOOS == `windows` {
		// ztool.Cmd_FastPrint(`按回车键继续...`)
		ztool.Cmd_aSyncExec(`pause`)
	}
	os.Exit(1)
}

// 新建日志对象
func NewLogger(name string) *Logger {
	return &Logger{
		name:   &name,
		level:  &Levell,
		output: &color.Output,
		deferr: func() {},
	}
}
func createLgr() *Logger {
	return &Logger{
		name:   &Main,
		level:  &Levell,
		output: &color.Output,
	}
}

// 创建副本
func (ll *Logger) Clone() *Logger {
	// l := new(Logger)
	l := loger_pool.Get().(*Logger)
	*l = *ll
	return l
}
func (ll *Logger) Free() {
	loger_pool.Put(ll)
	// ll = nil
}

// 重建日志组
/* 同 (*Logger).Clone().SetGroup([]string{item}) */
func (ll *Logger) NewGroup(g string) *Logger {
	l := ll.Clone()
	// l.group = []string{g}
	l.SetGroup(g)
	return l
	// l2 := new(Logger)
	// *l2 = *ll
	// l2.group = []string{g}
	// return l2
}

// 添加组项
func (ll *Logger) AddGroup(item string) *Logger {
	ll.group = append(ll.group, item)
	return ll
}

// 重建组项
/* 同 (*Logger).Clone().AddGroup(item) */
func (ll *Logger) AppGroup(item string) *Logger {
	l := ll.Clone()
	return l.AddGroup(item)
}

// 设置组项
func (ll *Logger) SetGroup(item string) *Logger {
	ll.group = []string{item}
	return ll
}

// 生成默认日志对象
// func DefLogger(name string, level int) *Logger {
// 	// if Default == nil {
// 	Default = &Logger{
// 		level: &level,
// 		name:  &name,
// 	}
// 	// }
// 	return Default
// }
