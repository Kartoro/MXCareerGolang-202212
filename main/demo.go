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

}
