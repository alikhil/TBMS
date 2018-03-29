package main

import (
	"fmt"

	"github.com/common-nighthawk/go-figure"
)

func main() {
	myFigure := figure.NewFigure("TBMS", "", true)
	myFigure.Print()
	fmt.Println("Starting...")
}
