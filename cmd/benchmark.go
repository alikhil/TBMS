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

type strToParam func(string) []*tuple.Tuple

func strToAuthorParam(text string) []*tuple.Tuple {
	words := strings.Fields(text)
	name := strings.Join(words[1:], " ")
	id, err := strconv.Atoi(words[0])
	if err != nil {
		panic("Can't convert author ID to int" + words[0])
	}

	logger.Info.Printf("Read: " + strconv.Itoa(id) + " -> " + name)
	params := []*tuple.Tuple{
		tuple.NewTupleFromItems("id", int32(id)),
		tuple.NewTupleFromItems("name", name)}
	return params

}

func strToPaperParam(text string) []*tuple.Tuple {
	words := strings.Fields(text)
	title := strings.Join(words[1:], " ")
	id, err := strconv.Atoi(words[0])
	if err != nil {
		panic("Can't convert paper ID to int" + words[0])
	}

	logger.Info.Printf("Read: " + strconv.Itoa(id) + " -> " + title)
	params := []*tuple.Tuple{
		tuple.NewTupleFromItems("id", int32(id)),
		tuple.NewTupleFromItems("title", title)}
	return params
}

func createObjsFromFile(filepath string, fn func(string) []*tuple.Tuple, create func([]*tuple.Tuple) bool) {
	inFile, _ := os.Open(filepath)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		params := fn(scanner.Text())
		ok := create(params)
		if !ok {
			panic("Didn't succeed to apply function")
		}
	}
}

func runBenchmark() {
	var re = &en.RealEngine{IO: io.LocalIO{}}
	re.InitDatabase()
	api.Init(re)
	logger.Info.Printf("from benchmark")

	// var path = "../test_data/1000/authors.in"
	// createObjsFromFile(path, strToAuthorParam,
	// 	func(properties []*tuple.Tuple) bool {
	// 		id, ok := api.CreateNode("Author", properties...)
	// 		logger.Info.Printf("Author added to id: %v", id)
	// 		return ok
	// 	})
	// path = "../test_data/1000/papers.in"
	// createObjsFromFile(path, strToAuthorParam,
	// 	func(properties []*tuple.Tuple) bool {
	// 		id, ok := api.CreateNode("Paper", properties...)
	// 		logger.Info.Printf("Paper added to id: %v", id)
	// 		return ok
	// 	})

	// createObjsFromFile(path, strToPaperParam,
	// 	func(properties []*tuple.Tuple) bool {

	// 	}
	// )

	// var author = api.createNode(("Author", {
	// 	tuple.NewTupleFromItems("id", 1),
	// 	tuple.NewTupleFromItems("name", "Einshtein")})
	node1, _ := api.CreateNode("Paper", tuple.NewTupleFromItems("id", 1), tuple.NewTupleFromItems("title", "Bitcons for breakfast"))
	node2, _ := api.CreateNode("Paper", tuple.NewTupleFromItems("id", 1), tuple.NewTupleFromItems("title", "Bitcons for lunch"))
	// var article2 = api.CreateNode("Paper",
	// 	tuple.NewTupleFromItems("id", 2),
	// 	tuple.NewTupleFromItems("title", "Bitcons for lunch"))

	api.CreateRelationship(node1, node2, "cites")
	// api.CreateRelationship(author, article1, "wrote")
	// api.CreateRelationship(author, article2, "wrote")

	// var id = 10
	// var id = 15
	// var paper1 = api.getNodeByTypeParameter() // (type="Paper", parameter={"id" == 10})
	// var paper2 = api.getNodeByID()            // (type="Paper", parameter={"id" == 15})

	// api.CreateRelationship(paper1, paper2, "cites")

	// var nodes = api.search() // ???

}
