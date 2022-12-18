package main

import (
	"fmt"
	"time"
)

func main() {
	var a = "initial"
	fmt.Println(a)

	// variables
	var (
		aa string
		b  int
		c  bool
	)
	fmt.Println(aa, b, c)

	// for loop
	for j := 7; j <= 9; j++ {
		fmt.Println(j)
	}

	// for range
	x := []int{1, 2, 3}
	for k, v := range x {
		fmt.Println(k, v)
	}

	// if else
	if num := 9; num < 0 {
		fmt.Println(num, "negative")
	} else if num < 10 {
		fmt.Println(num, "has 1 digit")
	} else {
		fmt.Println(num, "has multiple digits")
	}

	// switch case (increase readability)
	tt := time.Now().Weekday()
	switch tt {
	case time.Saturday, time.Sunday:
		{
			fmt.Println("It's the weekend")
		}
	default:
		fmt.Println("It's a weekday")

	}

	// another way of using switch case instead of if/else
	t := time.Now()
	switch {
	case t.Hour() < 12:
		{
			fmt.Println("It's before noon")
		}
	default:
		fmt.Println("It's after noon")

	}

	// switch GOLANG advance （反射）
	// i interface{} --> 表示 type 为 interface 类型的 parameters
	whatAmI := func(i interface{}) {
		switch t := i.(type) {
		case bool:
			{
				fmt.Println("I'm a bool")
			}
		case int:
			fmt.Println("I'm a int")

		default:
			fmt.Println("Don't know type", t)

		}
	}
	whatAmI(true)
	whatAmI(10)
	whatAmI("DavidDong")

	// array
	arry := [5]int{1, 2, 3, 4, 5}
	fmt.Println(arry)

	// two-dimensional array
	// declare
	var twoD [2][3]int
	// give values
	for i := 0; i < 2; i++ {
		for j := 0; j < 3; j++ {
			twoD[i][j] = i + j
		}
	}
	// display
	fmt.Println("2d-array", twoD)

	// slice normal declaration
	var ss []int
	fmt.Println(ss)

	// slice
	s := make([]string, 3)
	fmt.Println(s)
	s[0] = "Hao"
	s[1] = "De"
	s[2] = "Ba"
	fmt.Println(s)

	// slice append
	s = append(s, "app1")
	s = append(s, "app1", "app2")
	fmt.Println("APPEND:", s)

	// slice copy
	cc := make([]string, len(s))
	copy(cc, s) // deep copy which means cc and s are independent
	fmt.Println("Deep Copy:", cc)

	// slice copy
	xx := []int{1, 2, 3}
	y := xx
	xx[0] = 555
	fmt.Println("Light Copy:", xx, y) // same address

	// slicing operation (>=2 and <5)
	l := s[2:5]
	fmt.Println("sl1: ", l)

	l = s[:5]
	fmt.Println("sl2:", l)

	l = s[2:]
	fmt.Println("sl3:", l)

	ttt := []string{"declare", "and", "initialization"}
	fmt.Println("dandi: ", ttt)

	// two dimensional array (each dimension can be different length)
	twoDD := make([][]int, 3)
	for i := 0; i < 3; i++ {

		// dynamic inner length
		innerLen := i + 1
		twoDD[i] = make([]int, innerLen)
		for j := 0; j < innerLen; j++ {
			twoDD[i][j] = i + j
		}
	}
	fmt.Println(twoDD)

	// map
	m := make(map[string]int)

	m["k1"] = 7
	m["k2"] = 13

	n := map[string]int{"key1": 33, "key2": 34}
	fmt.Println(n["key1"])

	// variable as 函数 (utility)
	myFunc := func(a int, b int) (int, int) {
		fmt.Println("call myFunc")
		return a * b, a + b
	}
	res1, res2 := myFunc(3, 5)
	fmt.Println(res1, res2)

	//sum 可变参数
	sum()
	sum(1, 1, 1)
	ok := []int{1, 2, 3}
	sum(ok...)

	//intSeq (闭包保存状态)
	nextInt := intSeq()
	fmt.Println(nextInt())
	fmt.Println(nextInt())
	fmt.Println(nextInt())

	newInts := intSeq()
	fmt.Println(newInts())

	// 递归
	fmt.Println(fact(7))

	// 结构体方法
	r := rect{10, 5}
	fmt.Println("area: ", r.area())
	fmt.Println("perim: ", r.perim())
}

// 函数 (plusPlus --> private; PlusPlus --> public)
func plusPlus(a, b, c int, d, e, f string) (returnValue1 int, returnValue2 string) {
	return a + b + c, d + e + f
}

// 可变参数 (可以使用 interface{} --> 来判断阵对不同类型)
func sum(nums ...int) {
	fmt.Println(nums, " ")
	total := 0
	for _, num := range nums {
		total += num
	}
	fmt.Println(total)
}

// 闭包 intSeq函数 返回一个匿名函数， 该匿名函数入参也为空，返回一个int
// 它通过 intSeq 将匿名函数包起来了
func intSeq() func() int {
	i := 0
	return func() int {
		i++
		return i
	}
}

// 递归
func fact(n int) int {
	if n == 0 {
		return 1
	}
	return n * fact(n-1)
}

// 指针
// 结构体
type rect struct {
	width, height int
}

// 结构体方法
func (r rect) area() int {
	return r.width * r.height
}

func (r *rect) perim() int {
	return 2*r.width + 2*r.height
}
