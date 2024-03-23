// 简化数学运算(beta); 分类: Math_(数学)

package ztool

import "math"

// import (
// 	"math"
// 	"strings"
// )

// type number interface{ float64 | ~int } // 可用于运算的数字

// 取绝对值
// func Math_Abs(x int) uint {
// 	if x < 0 {
// 		return uint(-x)
// 	}
// 	return uint(x)
// }

// func zeroChk[T number](x T) {
// 	if x < 0 {
// 		panic(`zeroChk: v < 0`)
// 	}
// }

// 平方运算
/*
 计算一个数的 n 次方 (使用递归)
*/
// func Math_Pow[T number](x T, n uint) T {
// 	if n == 0 {
// 		return 1
// 	}
// 	return x * Math_Pow(x, n-1)
// }

// 计算一个数的 n 次方 (使用循环)
// func Math_PowX[T number](x T, n int) T {
// 	// zeroChk(n)
// 	n = int(Math_Abs(n))
// 	var y T = 1
// 	for i := 0; i < n; i++ {
// 		y *= x
// 	}
// 	return y
// }

// 四舍五入并保留 n 位小数
// func Math_RoundN(x float64, n uint) float64 {
// 	// if n < 0 {
// 	// 	panic(`Math_RoundN: n < 0`)
// 	// }
// 	f := Math_Pow(10.0, n)
// 	// f := math.Pow10(n)
// 	y := math.Round(x*f) / f
// 	return y
// }

// 移动字符串中的小数点
/*
 注：正反都可以哦 :)
*/
// func Math_RoundStr(s string, n int) string {
// 	i := strings.Index(s, `.`)
// 	if i == -1 {
// 		return s
// 	}
// 	r := s[:i] + s[i+1:]
// 	for i := 0; i < len(s)-n; i++ {
// 		r = r + "0"
// 	}
// 	return r + "." + s[i:i+1]
// }

// 计算位数 (使用对数)
/*
// digits returns the number of digits of a positive integer n

这段代码的功能是定义一个名为digits的函数，它接受一个正整数n作为参数，返回n的位数。函数的实现原理是使用对数函数，如我之前所解释的。

函数首先检查n是否为零，如果是，就返回1，因为零只有一位。然后，函数使用math包中的Log10函数，计算n的以10为底的对数，得到一个浮点数。然后，函数使用math包中的Ceil函数，将这个浮点数向上取整，得到一个整数。这个整数就是n的位数。
*/
func Math_digits(n int) int {
	// check if n is zero
	if n == 0 {
		return 1
	}
	// calculate the logarithm of n with base 10
	log := math.Log10(float64(n))
	// round up the logarithm and convert it to an integer
	return int(math.Ceil(log))
}

// 计算位数 (使用循环)
/*
// digits returns the number of digits of a positive integer n
不使用对数函数的一种方法是使用循环和除法。我们可以不断地将数字除以10，每次除法都会使数字的位数减少一位，直到数字变为0为止。我们可以用一个计数器来记录除法的次数，这个计数器就是数字的位数。
例如，如果我们要计算数字1234的位数，我们可以先将1234除以10，得到123，计数器加一；然后将123除以10，得到12，计数器加一；然后将12除以10，得到1，计数器加一；最后将1除以10，得到0，计数器加一。这时，计数器的值是4，这就是数字1234的位数。
这种方法的原理是，如果一个正整数n有k位，那么它一定满足10k−1≤n<10k
。也就是说，n一定在10的k-1次方和10的k次方之间。因此，如果我们不断地将n除以10，那么它一定会在k次之内变为0。如果我们记录下除法的次数，那么就得到k，也就是n的位数。
*/
func Math_digitx(n int) int {
	// check if n is zero
	if n == 0 {
		return 1
	}
	// initialize a counter
	count := 0
	// loop until n becomes zero
	for n > 0 {
		// divide n by 10
		n = n / 10
		// increment the counter
		count++
	}
	// return the counter
	return count
}
