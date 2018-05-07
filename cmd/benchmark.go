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




	// script

	//x2
	node1, _ := api.CreateNode("Paper", tuple.NewTupleFromItems("id", 1), tuple.NewTupleFromItems("title", "Bitcons for breakfast"))

	//x1
	rel, _ := api.CreateRelationship(node1, node2, "cites")

	//---------------------
	// select *

	allNodes := api.SelectNodesWhere(func (node *api.Node) bool {
		return true
	})

	// select where
	allPapers := api.SelectNodesWhere(func (node *api.Node) bool {
		return api.Contains((node.GetLabels(), "Paper")
	})

	// select where parameter.title == computer
	allPapersAboutComputers := api.SelectNodesWhere(func (node *api.Node) bool {
		if !api.Contains((node.GetLabels(), "Paper") {
			return false
		}

		title, ok := node.GetProperty("title")

		if ok {
			if strings.Contains("computer", strings.ToLower(title)) {
				return true
			}
		}
		return false
	})

	// select where links to papers from author nodes > 1
	allProductiveAuthors := api.SelectNodesWhere(func (node *api.Node) bool {
		if !api.Contains((node.GetLabels(), "Author") {
			return false
		}

		relationships := node.GetRelationships()

		count := 0

		for _, r := range relationships {
			if r.GetType() == "wrote" {
				count++
			}
		}

		return count > 1

	})

	// select where 
	allComputerScientist := api.SelectNodesWhere(func (node *api.Node) bool {
		if !api.Contains((node.GetLabels(), "Author") {
			return false
		}

		relationships := node.GetRelationships()

		for _, r := range relationships {
			paper := r.GetTo()

			title, ok := paper.GetProperty("title")

			if ok {
				if strings.Contains("computer", strings.ToLower(title)) {
					return true
				}
			}
		}

		return false
	}

	// select all Mazzara coauthors

	MazzaraName := "Manuel Mazzara"
	coauthors := make([]string, 0)

	Mazzara := api.SelectNodesWhere(func (node *api.Node) bool {
		if api.Contains((node.GetLabels(), "Author") {
			name, ok := node.GetProperty("name")
			if ok {
				return name == MazzaraName
			}
		}
		return false
	})[0]

	allMazzaraRels := Mazzara.GetRelationships()
	for _,r := range allMazzaraRels {
		if r.GetType() == "wrote" {
			links := r.GetTo().GetRelationships()
			for _, rr := range links {
				if rr.GetType() == "wrote" {
					name, ok := rr.GetFrom().GetProperty("name")
					if ok && name != MazzaraName {
						coauthors = append(coauthors, name)
					}
				}
			}
		}
	}

}
