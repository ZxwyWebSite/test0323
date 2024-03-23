// 高性能字符串操作; 分类: Str_(字符串)

package ztool

import "strings"

// Deprecated: 批量添加前后缀 (性能问题已弃用, 请使用FastJoinPrefixSuffix)
// var JoinPrefixSuffix = Str_FastJoinPrefixSuffix

// 计算字符串切片总长度
/*
 传入字符串切片, 返回其所有元素长度之和
 e.g. this({`ele3`,`ele2`,`ele1`}) => 12
*/
func Str_LenArray(a []string) (o int) {
	for i, r := 0, len(a); i < r; i++ {
		o += len(a[i])
	}
	return
}

// type (
// 	// 可取长度类型
// 	LenAble interface{ string | []any | chan any }
// )
// func LenArray[T LenAble](a []T) int {
// 	var o int
// 	for i, r := 0, len(a); i < r; i++ {
// 		o += len(a[i])
// 	}
// 	return o
// }

// 计算字符串总长度
/*
 注：由 Str_LenArray() 实现
*/
func Str_LenStrings(a ...string) int {
	return Str_LenArray(a)
}

// 字符串快速拼接
/*
 传入多个字符串参数, 返回拼接后的结果
 e.g. this("str1", "str2", "str3") => "str1str2str3"
 注：可用于替换字符串相加
*/
func Str_FastConcat(a ...string) string {
	var b strings.Builder
	b.Grow(Str_LenArray(a))
	for i, r := 0, len(a); i < r; i++ {
		b.WriteString(a[i])
	}
	return b.String()
}

// 字符串批量添加前后缀
/*
 给定一个字符串切片t, 前缀p, 后缀s, 遍历并输出添加前后缀的结果
 e.g. this({"A", "B", "C"}, "[", "] ") => "[A] [B] [C] "
 注：若切片为空, 则返回空字符串
*/
func Str_JoinPreSuffix(t []string, p string, s string) string { // this({`A`, `B`, `C`}, `[`, `] `) => `[A] [B] [C] `
	r := len(t)
	if r == 0 {
		return ``
	}
	var b strings.Builder
	b.Grow(Str_LenArray(t) * (1 + len(p) + len(s))) // sl := LenStrArray(t); sl + (len(p)+len(s))*sl
	for i := 0; i < r; i++ {
		b.WriteString(p)
		b.WriteString(t[i])
		b.WriteString(s)
	}
	return b.String()
}

// 字符串首尾相接
/*
 给定一个字符串str, 输出复制num份后的结果
 e.g. this(`=`, 20) => `====================`
 注：如str为空或num小于1, 将返回空字符串
*/
func Str_NumConcat(str string, num int) string {
	r := len(str)
	if num < 1 || r == 0 {
		return ``
	}
	var b strings.Builder
	b.Grow(r * num)
	for i := 0; i < num; i++ {
		b.WriteString(str)
	}
	return b.String()
}

// 字符串获取标识符前内容
/*
 给定一个字符串str, 输出sub前的内容
 e.g. this(`link?query`, `?`) => `link`
 注：未找到标识符则返回空字符串本身
*/
func Str_Before(str, sub string) string {
	idx := strings.Index(str, sub)
	if idx == -1 {
		return str //``
	}
	return str[:idx]
}
func Str_LastBefore(str, sub string) string {
	idx := strings.LastIndex(str, sub)
	if idx == -1 {
		return str
	}
	return str[:idx]
}

// 依次获取不为空的字符串
/*
 同JavaScript中 '??' '||'
*/
func Str_Select(str ...string) string {
	r := len(str)
	for i := 0; i < r; i++ {
		if str[i] != `` {
			return str[i]
		}
	}
	return str[r-1]
}

// func Str_PadStart(maxLength int, fillString string) string {
// 	length := len(fillString)
// }
