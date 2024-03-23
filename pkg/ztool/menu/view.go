package menu

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ZxwyWebSite/ztool"
)

type (
	List struct {
		Prefix string // 前缀
		Suffix string // 后缀
		Priter bool   // 美化
	}
)

/*
[00] Item0
[01] Item1
[10] Item10
[99] Item99

*/

// 格式化列表 <列表> <输出, 序号>
func (app *App) ShowList(list []string) string {
	var b strings.Builder
	length := len(list)
	pading := ztool.Math_digitx(length)
	// fmt.Printf("length: %v, padding: %v\n", length, pading)
	for i := 0; i < length; i++ {
		n := i + 1
		b.WriteString(app.Conf.List.Prefix)
		b.WriteString(ztool.Str_NumConcat(`0`, pading-ztool.Math_digitx(n)))
		b.WriteString(strconv.Itoa(n))
		b.WriteString(app.Conf.List.Suffix)
		b.WriteByte(' ')
		b.WriteString(list[i])
		b.WriteByte('\n')
	}
	return b.String()
}

// 格式化映射 <map> <输出>
func (app *App) ShowMaps(maps map[string]any) string {
	var b strings.Builder
	length := len(maps)
	l := make(map[string]int, length)
	var m int
	for k := range maps {
		n := len(k)
		l[k] = n
		if m < n {
			m = n
		}
	}
	i := 0
	for k, v := range maps {
		b.WriteString(app.Conf.List.Prefix)
		b.WriteString(k)
		b.WriteString(ztool.Str_NumConcat(` `, m-l[k]))
		b.WriteString(app.Conf.List.Suffix)
		b.WriteByte(' ')
		b.WriteString(fmt.Sprint(v))
		b.WriteByte('\n')
		i++
	}
	return b.String()
}
