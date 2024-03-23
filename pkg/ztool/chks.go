// 数据检查; 分类: Chk_(检查)

package ztool

// 检查输入是否为空
func Chk_IsNil(d ...any) (isnil bool) {
	// for _, v := range d {
	// 	switch t := v.(type) {
	// 	case string:
	// 		isnil = t == ``
	// 	case int:
	// 		isnil = t == 0
	// 	case bool:
	// 		isnil = !t
	// 	default:
	// 		isnil = true
	// 	}
	// 	if isnil {
	// 		return
	// 	}
	// }
	for i, r := 0, len(d); i < r; i++ {
		switch t := d[i].(type) {
		case string:
			isnil = t == ``
		case int:
			isnil = t == 0
		case bool:
			isnil = !t
		default:
			isnil = true
		}
		if isnil {
			return
		}
	}
	return
}

// 检查字符串是否为空
/*
 纯字符串对比请使用此函数以提高性能
*/
func Chk_IsNilStr(s ...string) (b bool) {
	for i, r := 0, len(s); i < r; i++ {
		if b = s[i] == ``; b {
			return
		}
	}
	return
}

// 检查是否符合给定值之一
func Chk_IsMatch(s string, h ...string) (isok bool) {
	// for _, v := range h {
	// 	isok = s == v
	// 	if isok {
	// 		return
	// 	}
	// }
	for i, r := 0, len(h); i < r; i++ {
		isok = s == h[i]
		if isok {
			return
		}
	}
	return
}
func Chk_IsMatchInt(s int, h ...int) (isok bool) {
	for i, r := 0, len(h); i < r; i++ {
		isok = s == h[i]
		if isok {
			return
		}
	}
	return
}
func Chk_IsMatchAny[T comparable](s T, h ...T) (isok bool) {
	for i, r := 0, len(h); i < r; i++ {
		isok = s == h[i]
		if isok {
			return
		}
	}
	return
}
func Chk_In[T comparable](s T, h ...[]T) (isok bool) {
	for i, r := 0, len(h); i < r; i++ {
		s2 := h[i]
		for i2, r2 := 0, len(s2); i2 < r2; i++ {
			isok = s == s2[i2]
			if isok {
				return
			}
		}
	}
	return
}

// 检查字符串是否相同
// func Chk_SameStr(a, b string) (ok bool) {
// 	la, lb := len(a), len(b)
// 	if la != lb {
// 		return
// 	}
// 	for i := 0; i < la; i++ {
// 		ok = a[i] == b[i]
// 		if !ok {
// 			return
// 		}
// 	}
// 	return
// }
