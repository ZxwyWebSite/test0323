// 数据整理(beta); 分类: Sort_(归类)

package ztool

// 翻转数组 <数组>
func Sort_Reverse[T any](arr []T) {
	for i, r := 0, len(arr); i < r/2; i++ {
		arr[i], arr[r-1-i] = arr[r-1-i], arr[i] // 交换两个元素的值
	}
}

// 反转数组并返回结果 <数组> <结果>
func Sort_ReverseNew[T any](arr []T) []T {
	length := len(arr)
	buf := make([]T, 0, length)
	for i := length - 1; i >= 0; i-- {
		buf = append(buf, arr[i]) // 将a中的元素从后向前追加到b中
	}
	return buf
}
