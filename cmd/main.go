package main

import (
	"fmt"

	utils "github.com/midea-media-llc/mm-go-utilities"
)

type B struct {
	Code string
}

type A struct {
	Name string
	B    *B
}

func main() {
	a := &A{
		Name: "Name",
		B: &B{
			Code: "Code",
		},
	}

	str := utils.ToSqlScript(a, "Model")
	fmt.Println(str)
}
