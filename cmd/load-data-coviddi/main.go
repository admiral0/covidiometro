package main

import (
	"fmt"
	"github.com/admiral0/covidiometro/covid"
	"github.com/admiral0/covidiometro/vaccines"
	"os"
)

func main() {
	g, err := covid.New(os.Args[1])
	if err != nil {
		panic(err)
	}

	h, err := g.Head()
	if err != nil {
		panic(err)
	}

	err = g.Load(h)
	if err != nil {
		panic(err)
	}

	fmt.Println(g.Data.Ref.Permalink)
	fmt.Println(g.Data.Ref.Updated)
	fmt.Println("----------------------")

	v, err := vaccines.New(os.Args[1])
	if err != nil {
		panic(err)
	}

	vh, err := v.Head()
	if err != nil {
		panic(err)
	}

	err = v.Load(vh)
	if err != nil {
		panic(err)
	}

	fmt.Println(v.Data.Ref.Permalink)
	fmt.Println(v.Data.Ref.Updated)
}
