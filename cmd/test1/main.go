package main

import (
	"fmt"

	"github.com/Un4ea42Zer0/GOBETA/properties"
)

func main() {
	fmt.Println("Hello")
	props, err := properties.LoadFrom("test.properties")
	if err != nil {
		panic(err)
	}

	value, _ := props.Get("name")
	fmt.Println("Name = ", value)
}
