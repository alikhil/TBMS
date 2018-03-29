package io

import (
	"testing"
)

func TestCreateFile(t *testing.T) {
	l := LocalIO{}
	c := LRUCache{l}
	c.createFile("test")
	if !c.fileExists("test") {
		t.Errorf("file does not created")
	}
}
