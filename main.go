package main

import (
	"fmt"

	"carbonaut.cloud/pkg/schema"
)

func main() {
	e := schema.Energy{
		Amount: 2669267,
		Unit:   schema.MICROWATT,
		Name:   "host",
	}
	fmt.Println(e.ConvertToKilowatt())
}
