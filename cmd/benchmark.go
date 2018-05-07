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

	// Create authors
	RobertoLucchi, _ := api.CreateNode("Author", tuple.NewTupleFromItems("name", "Roberto Lucchi"))
	ClaudioGuidi, _ := api.CreateNode("Author", tuple.NewTupleFromItems("name", "Claudio Guidi"))
	IvanLanese, _ := api.CreateNode("Author", tuple.NewTupleFromItems("name", "Ivan Lanese"))
	ManuelMazzara, _ := api.CreateNode("Author", tuple.NewTupleFromItems("name", "Manuel Mazzara"))

	// Create papers
	paper1, _ := api.CreateNode("Paper", tuple.NewTupleFromItems("title", "A pi-calculus based semantics for WS-BPEL"))
	paper2, _ := api.CreateNode("Paper", tuple.NewTupleFromItems("title", "A formal framework for web services coordination"))
	paper3, _ := api.CreateNode("Paper", tuple.NewTupleFromItems("title", "Towards a unifying theory for web services composition"))
	paper4, _ := api.CreateNode("Paper", tuple.NewTupleFromItems("title", "Timing issues in web services composition"))

	// Create relationships
	api.CreateRelationship(RobertoLucchi, paper1, "wrote")
	api.CreateRelationship(RobertoLucchi, paper3, "wrote")
	api.CreateRelationship(ClaudioGuidi, paper2, "wrote")
	api.CreateRelationship(IvanLanese, paper3, "wrote")
	api.CreateRelationship(ManuelMazzara, paper1, "wrote")
	api.CreateRelationship(ManuelMazzara, paper2, "wrote")
	api.CreateRelationship(ManuelMazzara, paper3, "wrote")
	api.CreateRelationship(ManuelMazzara, paper4, "wrote")

	// script

	//---------------------
	// select *

	allNodes, _ := api.SelectNodesWhere(func(node *api.Node) bool {
		return true
	})

	logger.Info.Printf("Print all nodes:")
	for _, node := range allNodes {
		logger.Info.Printf("Node: %v", node)
	}

	// select where
	allPapers, _ := api.SelectNodesWhere(func(node *api.Node) bool {
		labels := node.GetLabels()
		l := make([]interface{}, 0)
		for _, label := range *labels {
			l = append(l, label)
		}
		return api.Contains(&l, "Paper")
	})

	logger.Info.Printf("Print all papers:")
	for _, node := range allPapers {
		logger.Info.Printf("Node: %v", node)
	}

	// select where parameter.title contains == services
	allPapersAboutServices, _ := api.SelectNodesWhere(func(node *api.Node) bool {
		labels := node.GetLabels()
		l := make([]interface{}, 0)
		for _, label := range *labels {
			l = append(l, label)
		}

		if !api.Contains(&l, "Paper") {
			return false
		}

		title, ok := node.GetProperty("title")

		if ok {
			if strings.Contains(strings.ToLower(title.(string)), "services") {
				return true
			}
		}
		return false
	})

	logger.Info.Printf("Print all papers about services:")
	for _, node := range allPapersAboutServices {
		logger.Info.Printf("Node: %v", node)
	}

	// select where links to papers from author nodes > 1
	allProductiveAuthors, _ := api.SelectNodesWhere(func(node *api.Node) bool {
		labels := node.GetLabels()
		l := make([]interface{}, 0)
		for _, label := range *labels {
			l = append(l, label)
		}

		if !api.Contains(&l, "Author") {
			return false
		}

		relationships := node.GetRelationships()
		count := 0

		for _, r := range *relationships {
			if r.GetType() == "wrote" {
				count++
			}
		}

		return count > 1

	})

	logger.Info.Printf("Print all productive authors (papers count > 1):")
	for _, node := range allProductiveAuthors {
		logger.Info.Printf("Node: %v", node)
	}

	// select where
	allComputerScientist, _ := api.SelectNodesWhere(func(node *api.Node) bool {
		labels := node.GetLabels()
		l := make([]interface{}, 0)
		for _, label := range *labels {
			l = append(l, label)
		}

		if !api.Contains(&l, "Author") {
			return false
		}

		relationships := node.GetRelationships()

		for _, r := range *relationships {
			paper := r.GetTo()

			title, ok := paper.GetProperty("title")
			if ok {
				if strings.Contains(strings.ToLower(title.(string)), "services") {
					return true
				}
			}
		}
		return false
	})

	logger.Info.Printf("Print all authors about services:")
	for _, node := range allComputerScientist {
		logger.Info.Printf("Node: %v", node)
	}

	// select all Mazzara coauthors
	MazzaraName := "Manuel Mazzara"
	coauthors := make([]string, 0)

	m, _ := api.SelectNodesWhere(func(node *api.Node) bool {
		labels := node.GetLabels()
		l := make([]interface{}, 0)
		for _, label := range *labels {
			l = append(l, label)
		}
		if api.Contains(&l, "Author") {
			name, ok := node.GetProperty("name")
			if ok {
				return name == MazzaraName
			}
		}
		return false
	})
	Mazzara := *m[0]

	allMazzaraRels := Mazzara.GetRelationships()
	for _, r := range *allMazzaraRels {
		if r.GetType() == "wrote" {
			links := r.GetTo().GetRelationships()
			for _, rr := range *links {
				if rr.GetType() == "wrote" {
					name, ok := rr.GetFrom().GetProperty("name")
					if ok && name != MazzaraName {
						coauthors = append(coauthors, name.(string))
					}
				}
			}
		}
	}

	logger.Info.Printf("Print all Mazzara couathors:")
	for _, node := range coauthors {
		logger.Info.Printf("Node: %v", node)
	}

}
