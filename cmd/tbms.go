package main

import (
	"fmt"

	logger "github.com/alikhil/TBMS/internals/logger"
	"github.com/common-nighthawk/go-figure"
)

func main() {
	myFigure := figure.NewFigure("TBMS", "", true)
	myFigure.Print()
	fmt.Println("Starting...")
	logger.Info.Printf("hello from logger")

}
