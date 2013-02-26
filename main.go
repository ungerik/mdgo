package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func makeGoFile(filename string, info os.FileInfo) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	lines := bytes.Split(data, []byte{'\n'})
	for i, line := range lines {
		if len(line) == 0 {
			continue
		}
		if line[0] == '\t' {
			lines[i] = line[1:]
		} else {
			lines[i] = append([]byte("// "), line...)
		}
	}
	data = bytes.Join(lines, []byte{'\n'})

	if strings.HasSuffix(filename, ".md") {
		filename = filename[:len(filename)-len(".md")]
		if !strings.HasSuffix(filename, ".go") {
			filename += ".go"
		}
	}
	fmt.Println(filename)
	return ioutil.WriteFile(filename, data, info.Mode())
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
