package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	files := readDirForNode(path, printFiles)

	printNodes(out, files, path, 0)
	return nil
}

func printNodes(out io.Writer, files []Entity, removePath string, isLastCount int) {
	fileLen := len(files)
	for i, file := range files {
		joinSymbol := "├───"
		if i == fileLen-1 {
			joinSymbol = "└───"
		}

		paths := strings.Split(file.Path, "/")
		splachCount := len(paths)
		fileName := paths[splachCount-1]

		notLast := strings.Repeat("│\t", splachCount-2-isLastCount)

		isLast := strings.Repeat("\t", isLastCount)

		printedName := notLast + isLast + joinSymbol + fileName
		if file.IsDir == false && file.Size != "" {
			printedName += " " + file.Size
		}
		fmt.Fprintln(out, printedName)

		if file.IsDir && len(file.Childs) > 0 {
			if joinSymbol == "└───" {
				isLastCount++
			}
			printNodes(out, file.Childs, removePath, isLastCount)
		}
	}
}

type Entity struct {
	Path   string
	IsDir  bool
	Size   string
	Childs []Entity
}

func readDirForNode(path string, printFiles bool) (nodes []Entity) {
	readFiles, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range readFiles {
		fullName := path + "/" + file.Name()
		node := Entity{Path: fullName, IsDir: file.IsDir(), Childs: []Entity{}}
		if file.IsDir() {
			subDirectoryFiles := readDirForNode(fullName, printFiles)
			node.Childs = append(node.Childs, subDirectoryFiles...)
			nodes = append(nodes, node)
		} else {
			if file.Size() == 0 {
				node.Size = "(empty)"
			} else {
				node.Size = fmt.Sprintf("(%db)", file.Size())
			}
			if printFiles {
				nodes = append(nodes, node)
			}
		}

	}

	return nodes
}
