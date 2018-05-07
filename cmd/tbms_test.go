package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"testing"

	api "github.com/alikhil/TBMS/internals/api"
	en "github.com/alikhil/TBMS/internals/engine"
	io "github.com/alikhil/TBMS/internals/io"
	tuple "github.com/kmanley/golang-tuple"
)

func BenchmarkLocalIO(b *testing.B) {
	benchmarkCreate(io.LocalIO{}, b)
	benchmarkInsert(io.LocalIO{}, b)
}

func benchmarkCreate(i io.IO, b *testing.B) {
	var re = &en.RealEngine{IO: i}
	re.InitDatabase()
	api.Init(re)
	re.DropDatabase()
}

func benchmarkInsert(i io.IO, b *testing.B) {
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
	re.DropDatabase()
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
