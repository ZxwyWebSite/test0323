// 类型转换驱动 (beta)
// 注：接口类型无法进行更多操作，推荐使用Value()获取原值后手动转换
package cache

import (
	"errors"

	"github.com/ZxwyWebSite/ztool"
)

// ItemVal 对象类型转换
type ItemVal struct {
	name  string
	value interface{}
}

var (
	Err_FormatErr = errors.New(`FormatErr`) // 类型转换错误
)

// 获取键名
func (iv *ItemVal) Name() string {
	return iv.name
}

// 获取原值
func (iv *ItemVal) Value() interface{} {
	return iv.value
}

// 到结构体
func (iv *ItemVal) MapTo(to interface{}) error {
	return ztool.Val_MapToStruct(iv.value, to)
}

// 到字节组
func (iv *ItemVal) Bytes() ([]byte, error) {
	b, ok := iv.value.([]byte)
	if !ok {
		return nil, Err_FormatErr
	}
	return b, nil
}
func (iv *ItemVal) MustBytes() []byte {
	out, err := iv.Bytes()
	if err != nil {
		return []byte{}
	}
	return out
}

// 到字符串
func (iv *ItemVal) String() (string, error) {
	s, ok := iv.value.(string)
	if !ok {
		return ``, Err_FormatErr
	}
	return s, nil
}
func (iv *ItemVal) MustString() string {
	out, _ := iv.String()
	return out
}

// 到整数型
func (iv *ItemVal) Int() (int, error) {
	i, ok := iv.value.(int)
	if !ok {
		return 0, Err_FormatErr
	}
	return i, nil
}
func (iv *ItemVal) MustInt() int {
	out, _ := iv.Int()
	return out
}

// func Any[T interface{}](iv *ItemVal) (T, bool) {
// 	val, ok := iv.value.(T)
// 	return val, ok
// }

// func (iv *ItemVal) ()  {
// }
