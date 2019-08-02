package main

import "fmt"

type A struct {
	a int
}

func f(x []A) {
	//for i := 0; i < len(x); i++ {
	for _, X := range x {
		X.a *= 10
	}
}

func main() {
	ar := make([]A, 0)
	//a := A{a: 100}
	k := 10
	x := make([]A, 1)
	x[0] = A{a: 200}
	for i := 0; i < k; i++ {
		ar = append(ar, x[0])
	}
	x[0].a = 300
	for i := 0; i < k; i++ {
		fmt.Println(ar[i].a)
	}

	f(ar)
	for i := 0; i < k; i++ {
		fmt.Println(ar[i].a)
	}
}
