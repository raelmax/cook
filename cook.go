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

var regex, _ = regexp.Compile("{{(\\s|)cook?\\w+\\.\\w+(\\s|)}}")

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
		config, err := ioutil.ReadFile(repoPath + string(os.PathSeparator) + configNames[index])

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
			text = strings.TrimSuffix(text, "\n")
			config[k] = text
		}
	}

	return config
}

// getKey receive a string with a placeholder, parse and return a key string
func getKey(placeholder string) string {
	placeholder = regex.FindString(placeholder)
	strParts := strings.Split(placeholder, ".")[1]

	key := strings.Replace(strParts, "}}", "", -1)
	key = strings.TrimSpace(key)
	return key
}

// getPaths receive a repository path and returns a slice with template files/dirs path
func getPaths(repoPath string) []string {
	var paths = make([]string, 0)

	filepath.Walk(repoPath, func(fp string, fi os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err) // can't walk here,
			return nil       // but continue walking elsewhere
		}

		matched := regex.MatchString(fi.Name())

		if matched {
			paths = append(paths, fp)
		}

		return nil
	})

	return paths
}

// ReplacePaths receives a slices of paths and a config map to replace
// variables with config values
func ReplacePaths(paths []string, config map[string]interface{}) string {
	var newFolder string

	// reverse paths list to rename files/dirs without lost references
	for i := len(paths) - 1; i >= 0; i-- {
		parts := strings.Split(paths[i], string(os.PathSeparator))
		replacePart := parts[len(parts)-1]
		originalPart := parts[len(parts)-1]

		key := getKey(replacePart)
		value := config[key].(string)

		originalPart = regex.ReplaceAllString(originalPart, value)
		parts[len(parts)-1] = originalPart
		newPath := strings.Join(parts, string(os.PathSeparator))

		os.Rename(paths[i], newPath)

		if i == 0 {
			folderParts := strings.Split(newPath, string(os.PathSeparator))
			oldFolder := folderParts[0]
			newFolder = folderParts[1]

			os.Rename(newPath, newFolder)
			os.RemoveAll(oldFolder)
		}
	}
	return newFolder
}

//ReplaceContent receives a config map to replace all files with your variables
func ReplaceContent(repoPath string, config map[string]interface{}) {
	filepath.Walk(repoPath, func(fp string, fi os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err) // can't walk here,
			return nil       // but continue walking elsewhere
		}

		if !!fi.IsDir() {
			return nil // not a file.  ignore.
		}

		file, err := ioutil.ReadFile(fp)
		if err != nil {
			fmt.Println(err)
		}

		lines := strings.Split(string(file), "\n")
		for i, line := range lines {
			placeholders := regex.FindAllString(line, -1)
			for _, ph := range placeholders {
				key := getKey(ph)
				value := config[key].(string)
				lines[i] = strings.Replace(lines[i], ph, value, 1)
			}
		}
		output := strings.Join(lines, "\n")
		err = ioutil.WriteFile(fp, []byte(output), 0644)
		if err != nil {
			fmt.Println(err)
		}

		return nil
	})
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
	repoPath = ReplacePaths(paths, config)
	ReplaceContent(repoPath, config)
	fmt.Println("Project genereated: ", repoPath)
}
