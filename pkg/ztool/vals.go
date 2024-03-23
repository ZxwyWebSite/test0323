// 数据操作; 分类: Val_(数值)

package ztool

import (
	"bytes"
	"encoding/gob"

	"github.com/ZxwyWebSite/ztool/x/json"
)

// Data映射到结构体 <数据源, 映射到(指针)>
/*
 将 map[string]any 映射到 struct, 常用于接口数据处理
 注：目前只能用json二次转换来解决，实测mapstructor性能较差
*/
func Val_MapToStruct(from, to any) error {
	data, err := json.Marshal(from)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, to)
}

// Gob编码 <数据><结果,错误>
func Val_GobEncode(data any) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Gob解码 <数据,目标><错误>
func Val_GobDecode(data []byte, to any) error {
	buf := bytes.NewReader(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(to)
}
