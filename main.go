// Cook generate projects from templates based on github repositories.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const gitHub = "git@github.com:"

// Clone fetch a git repository to current directory and returns
// a directory name
func Clone(repoName string) string {
	repoURL := gitHub + repoName + ".git"

	cloneCmd := exec.Command("git", "clone", repoURL)
	cloneOut, err := cloneCmd.Output()

	if err != nil {
		panic(err)
	}

	fmt.Println(string(cloneOut))

	return strings.Split(repoName, "/")[1]
}

// Parse a json config file to ask user to new values
func Parse(repoPath string) string {
	return "YAY"
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
