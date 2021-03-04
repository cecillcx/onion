package main

import (
	"controller"
	"fmt"
	"onion"
)

func add(a, b, c int) int {
	fmt.Printf("I'll print:%v add %v\n", a, b)
	return a + b + c
}

type IsSStructTest struct {
	name string
	re   int
}

func returnStruct(a, b, c int, name string) IsSStructTest {
	return IsSStructTest{name, a + b + c}
}

func main() {
	o := onion.NewOnion()
	r := controller.NewRateLimiterCtrl(10)
	errCheckFunc := func(ret []interface{}) error {
		if ret[0].(int) < 10 {
			return fmt.Errorf("is error")
		}
		return nil
	}
	er := controller.NewRetryCtrl(3, errCheckFunc)

	o.ApplyFunc(add).AddCtrl(r).AddCtrl(er)

	o.Invoke(1, 100, 2)
	fmt.Println("-------------")
	fmt.Printf("res: %v\n", o.Invoke(1, 2, 10)[0])

	o.ApplyFunc(returnStruct).AddCtrl(r)
	fmt.Printf("%v", o.Invoke(1, 2, 3, "test1"))
}
