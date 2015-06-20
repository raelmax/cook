// Cook generate projects from templates based on github repositories.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	GITHUB_URL = "git@github.com:"
)

// Clone fetch a git repository to current directory and returns
// a directory name
func Clone(repoName string) string {
	repoUrl := GITHUB_URL + repoName + ".git"

	cloneCmd := exec.Command("git", "clone", repoUrl)
	cloneOut, err := cloneCmd.Output()

	if err != nil {
		panic(err)
	}

	fmt.Println(string(cloneOut))

	return strings.Split(repoName, "/")[1]
}

func main() {
	var repoName, repoPath string

	if len(os.Args) > 1 {
		repoName = os.Args[1]
	} else {
		fmt.Println("You must provide a github <username>/<repository>")
		return
	}

	repoPath = Clone(repoName)
	fmt.Println(repoPath)
}
