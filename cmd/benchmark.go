package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/kmanley/golang-tuple"

	api "github.com/alikhil/TBMS/internals/api"
	en "github.com/alikhil/TBMS/internals/engine"
	io "github.com/alikhil/TBMS/internals/io"
	"github.com/alikhil/TBMS/internals/logger"
)

type strToParam func(string) *tuple.Tuple

func strToAuthorParam(text string) []*tuple.Tuple {
	words := strings.Fields(text)
	name := strings.Join(words[1:], " ")
	id, err := strconv.Atoi(words[0])
	if err != nil {
		panic("Can't convert to int" + words[0])
	}

	logger.Info.Printf("Read: " + strconv.Itoa(id) + " -> " + name)
	params := []*tuple.Tuple{
		tuple.NewTupleFromItems("id", id),
		tuple.NewTupleFromItems("name", name)}
	return params

}

func scanParamFromFile(filename string) {

}

func runBenchmark() {
	var re = &en.RealEngine{IO: io.LocalIO{}}
	re.InitDatabase()
	api.Init(re)
	logger.Info.Printf("from benchmark")

	var path = "./authors.in"

	inFile, _ := os.Open(path)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		params := strToAuthorParam(scanner.Text())
		logger.Info.Printf("%v", (*params[1]).Get(1))
	}

	// var author = api.createNode()   // (type = "Author", properties = {id: 1; name: "andrew NG"})
	// var article1 = api.createNode() // (type = "Paper", properties = {id: 10; title: "Bitcons for breakfast"})
	// var article2 = api.createNode() // (type = "Paper", properties = {id: 15; title: "Bitcons for lunch"})

	// api.CreateRelationship(author, article1, "wrote")
	// api.CreateRelationship(author, article2, "wrote")

	// var id = 10
	// var id = 15
	// var paper1 = api.getNodeByTypeParameter() // (type="Paper", parameter={"id" == 10})
	// var paper2 = api.getNodeByID()            // (type="Paper", parameter={"id" == 15})

	// api.CreateRelationship(paper1, paper2, "cites")

	// var nodes = api.search() // ???

}
