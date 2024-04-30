package main

import (
	"fmt"

	"carbonaut.cloud/pkg/providers/scaphandre"
)

func main() {
	p := scaphandre.Provider{}
	e, err := p.CollectEnergy("136.144.49.227:8080/metrics")
	if err != nil {
		panic(err)
	}
	fmt.Println(e)
}
