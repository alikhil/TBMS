package main

import (
	// "fmt"
	// "github.com/alikhil/distributed-fs/utils"
	"strings"

	"github.com/kmanley/golang-tuple"

	api "github.com/alikhil/TBMS/internals/api"
	en "github.com/alikhil/TBMS/internals/engine"
	io "github.com/alikhil/TBMS/internals/io"
	"github.com/alikhil/TBMS/internals/logger"
)

func runExample() {
	// masterEndpoint := fmt.Sprintf("%s:%v", utils.GetIPAddress(), 5001)
	// var mapping = en.GetFileToBytesMap()

	// client, ok := utils.GetRemoteClient(masterEndpoint)
	// if !ok {
	// 	panic("failed to connect to remote client")
	// }
	// dfs := utils.DFSClient{Client: client}
	// dfs.InitRecordMappings(mapping)

	// cache := io.LRUCache{}
	// cache.Init(&dfs, mapping, 5)

	// var re = &en.RealEngine{IO: &cache}
	var re = &en.RealEngine{IO: &io.LocalIO{}}
	re.InitDatabase()
	api.Init(re)
	logger.Info.Printf("From example")
	defer re.DropDatabase()

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
		logger.Info.Printf("%v", node)
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
		logger.Info.Printf("%v", node)
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
		logger.Info.Printf("%v", node)
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

	logger.Info.Printf("Print all authors who wrote about services:")
	for _, node := range allComputerScientist {
		logger.Info.Printf("%v", node)
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

	logger.Info.Printf("Print all prof. Mazzara couathors:")
	for _, node := range coauthors {
		logger.Info.Printf("%v", node)
	}

}
