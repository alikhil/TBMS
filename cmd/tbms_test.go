package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/alikhil/TBMS/internals/logger"
	"github.com/alikhil/distributed-fs/utils"

	api "github.com/alikhil/TBMS/internals/api"
	en "github.com/alikhil/TBMS/internals/engine"
	io "github.com/alikhil/TBMS/internals/io"
	tuple "github.com/kmanley/golang-tuple"
)

var remoteIP = "10.91.41.109:5001"

func BenchmarkLocalIO(b *testing.B) { benchmarkIO(io.LocalIO{}, b) }

func BenchmarkLocalIOCache(b *testing.B) { benchmarkIO(getCache(), b) }

func BenchmarkDFS(b *testing.B) { benchmarkIO(getDFS(), b) }

func BenchmarkDFSCache(b *testing.B) { benchmarkIO(getDFSWithCache(), b) }

func benchmarkIO(i io.IO, b *testing.B) {
	startTime := time.Now()
	var re = &en.RealEngine{IO: i}
	re.InitDatabase()
	api.Init(re)
	b.Run("insert_data", func(b *testing.B) {
		path := "../test_data/1000/papers.in"
		createObjsFromFile(path, strToPaperParam,
			func(properties []*tuple.Tuple) bool {
				_, ok := api.CreateNode("Paper", properties...)
				return ok
			})
	})
	endTime := time.Now()
	logger.Info.Printf("Insertion of 1000 nodes took %v", endTime.Sub(startTime))

	startTime = time.Now()
	api.SelectNodesWhere(func(*api.Node) bool { return true })
	endTime = time.Now()
	logger.Info.Printf("Retrieval of 1000 nodes took %v", endTime.Sub(startTime))

	// startTime = time.Now()
	re.DropDatabase()
}

func getCache() io.IO {
	var mapa = en.GetFileToBytesMap()

	cache := io.LRUCache{}
	cache.Init(io.LocalIO{}, mapa, 5)
	return &cache
}

func getDFS() io.IO {
	var mapa = en.GetFileToBytesMap()

	client, ok := utils.GetRemoteClient(remoteIP)
	if !ok {
		panic("failed to connect to remote client")
	}
	dfs := utils.DFSClient{Client: client}
	dfs.InitRecordMappings(mapa)
	return &dfs
}

func getDFSWithCache() io.IO {

	var mapa = en.GetFileToBytesMap()

	client, ok := utils.GetRemoteClient(remoteIP)
	if !ok {
		panic("failed to connect to remote client")
	}
	dfs := utils.DFSClient{Client: client}
	dfs.InitRecordMappings(mapa)

	cache := io.LRUCache{}
	cache.Init(&dfs, mapa, 5)

	return &cache
}

func strToPaperParam(text string) []*tuple.Tuple {
	words := strings.Fields(text)
	title := strings.Join(words[1:], " ")
	id, err := strconv.Atoi(words[0])
	if err != nil {
		panic("Can't convert paper ID to int" + words[0])
	}

	params := []*tuple.Tuple{
		tuple.NewTupleFromItems("id", id),
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
