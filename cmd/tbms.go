package main

import "github.com/alikhil/TBMS/internals/logger"

func main() {
	// myFigure := figure.NewFigure("TBMS", "", true)
	// myFigure.Print()
	// fmt.Println("Starting...")

	logger.Info.Printf("Start")

	runBenchmark()

}
