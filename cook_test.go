package main

import (
	"os"
	"testing"
)

func TestGet(t *testing.T) {
	repoName := "raelmax/cook-basic-template"
	repoPath := Get(repoName)

	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		t.Errorf("Repository %s was not cloned.", repoName)
	}

	os.RemoveAll(repoPath)
}
