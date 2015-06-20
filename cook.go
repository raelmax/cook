// Cook generate projects from templates based on github repositories.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// Clone fetch a git repository to current directory and returns a directory name
func Clone(repoName string) string {
	repoURL := "git@github.com:" + repoName + ".git"

	cloneCmd := exec.Command("git", "clone", repoURL)
	_, err := cloneCmd.Output()

	if err != nil {
		panic(err)
	}

	return strings.Split(repoName, "/")[1]
}

// Parse a json config file to ask user to new values
func Parse(repoPath string) map[string]interface{} {
	configName := "cookiecutter.json"
	config, err := ioutil.ReadFile(repoPath + "/" + configName)
	if err != nil {
		panic(err)
	}

	var configJson map[string]interface{}
	json.Unmarshal([]byte(config), &configJson)
	return configJson
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
	config := Parse(repoPath)
	fmt.Print(config["full_name"])
}
