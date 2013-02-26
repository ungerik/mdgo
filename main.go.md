main.go
-------

	package main

	import (
		"bytes"
		"fmt"
		"io/ioutil"
		"os"
		"os/exec"
		"path/filepath"
		"strings"
	)

makeGoFile creates a .go file from a .md file

	func makeGoFile(filename string, info os.FileInfo) error {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}

		lastEmptyLine := -666
		lastCommentedLine := -666
		lines := bytes.Split(data, []byte{'\n'})
		for i := 0; i < len(lines); i++ {
			line := lines[i]
			if len(line) == 0 {
				lastEmptyLine = i
				continue
			}
			if line[0] == '\t' {
				lines[i] = line[1:]
				if i-1 == lastEmptyLine && i-2 == lastCommentedLine {
					lines[i-1] = []byte("// ")
				}
			} else {
				lines[i] = append([]byte("// "), line...)
				lastCommentedLine = i
			}
		}
		data = bytes.Join(lines, []byte{'\n'})

Be robust about file names, remove .md:

		if strings.HasSuffix(filename, ".md") {
			filename = filename[:len(filename)-len(".md")]

Add .go extension if not present:

			if !strings.HasSuffix(filename, ".go") {
				filename += ".go"
			}
		}

Print out the name of the generated file to be able to pipe it to other programs:

		fmt.Println(filename)
		err = ioutil.WriteFile(filename, data, info.Mode())
		if err != nil {
			return err
		}
		return exec.Command("gofmt", "-w", filename).Run()
	}

	func main() {
		args := os.Args[1:]
		if len(args) == 0 {
			args = append(args, ".")
		}
		for _, arg := range args {
			info, err := os.Stat(arg)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			if info.IsDir() {
				err = filepath.Walk(arg, func(path string, info os.FileInfo, err error) error {
					if info.IsDir() || err != nil {
						return err
					}
					if match, err := filepath.Match("*.go.md", path); !match {
						return err
					}
					return makeGoFile(path, info)
				})
			} else {
				err = makeGoFile(arg, info)
			}
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
