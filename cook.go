// Cook generate projects from templates based on github repositories.
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
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
	var configJson map[string]interface{}
	configNames := [2]string{"cook.json", "cookiecutter.json"}

	for index := range configNames {
		config, err := ioutil.ReadFile(repoPath + "/" + configNames[index])

		if err != nil {
			continue
		}

		json.Unmarshal([]byte(config), &configJson)
	}

	if len(configJson) == 0 {
		panic("This is not a valid repository.")
	}

	return configJson
}

// Ask receive a config map, iterate over and update user project data
func Ask(config map[string]interface{}) map[string]interface{} {
	for k, v := range config {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter ", k, " (default: ", v, "): ")
		text, _ := reader.ReadString('\n')

		if len(strings.TrimSpace(text)) > 0 {
			config[k] = text
		}
	}

	return config
}

// getPaths receive a repository path and returns a slice with template files/dirs path
func getPaths(repoPath string) []string {
	var paths = make([]string, 0)

	filepath.Walk(repoPath, func(fp string, fi os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err) // can't walk here,
			return nil       // but continue walking elsewhere
		}

		regex, _ := regexp.Compile("{{(\\s|)cook?\\w+\\.\\w+(\\s|)}}")
		matched := regex.MatchString(fi.Name())

		if matched {
			paths = append(paths, fp)
		}

		return nil
	})

	return paths
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
	config = Ask(config)

	paths := getPaths(repoPath)
	fmt.Println(paths)
}
