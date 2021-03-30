package main

import (
	"fmt"
	"io/fs"
	"os/exec"
	"path/filepath"
)

func main() {
	var testPath = "../dep-updater"
	files := getFiles(testPath)
	fmt.Println(files)
	for _, f := range files {
		err := runPrune(f)
		if err != nil {
			panic(err)
		}
	}
	// for _, f := range files {
	// 	err := runUpdate(f)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
	for _, f := range files {
		outdated, err := getOutdated(f)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(outdated))
	}
}

func getFiles(path string) []string {
	var files = make([]string, 0)
	err := filepath.WalkDir(path, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == "node_modules" {
			return fs.SkipDir
		}
		if info.Name() == "package.json" {
			files = append(files, filepath.Dir(path))
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	uniqFiles := removeDuplicates(files)
	return uniqFiles
}

func runPrune(path string) error {
	cmd := exec.Command("npm", "prune")
	cmd.Dir = path
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// func runUpdate(path string) error {
// 	cmd := exec.Command("npm", "update")
// 	cmd.Dir = path
// 	err := cmd.Run()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func getOutdated(path string) ([]byte, error) {
	args := []string{"outdated", "--json"}
	cmd := exec.Command("npm", args...)
	cmd.Dir = path
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return output, nil
}

func removeDuplicates(fileList []string) []string {
	uniqueMap := map[string]bool{}
	for v := range fileList {
		uniqueMap[fileList[v]] = true
	}
	result := []string{}
	for key := range uniqueMap {
		result = append(result, key)
	}
	return result
}
