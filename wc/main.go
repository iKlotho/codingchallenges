package main

import (
	"flag"
	"fmt"
	"os"
)

type WCOutput interface {
	PrintResult()
}

type WCFile struct {
	Name string
	data []byte
}

func (wcf *WCFile) loadFile() error {
	f, err := os.ReadFile(wcf.Name)
	wcf.data = f
	return err
}

func (wcf *WCFile) Size() int {
	return len(wcf.data)
}

func NewFile(name string) (*WCFile, error) {
	wcf := &WCFile{Name: name}
	err := wcf.loadFile()
	if err != nil {
		return nil, err
	}
	return wcf, err
}

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("Please provide a filename\n")
		os.Exit(1)
	}

	byteFlagPtr := flag.Bool("c", false, "Print number of bytes in each input file")
	flag.Parse()

	allFiles := flag.Args()
	files := make([]*WCFile, 0)
	for _, file := range allFiles {
		wcf, err := NewFile(file)
		if err != nil {
			panic(err)
		}
		files = append(files, wcf)
	}

	for _, file := range files {
		out := "  "
		if *byteFlagPtr {
			out += fmt.Sprintf("%d ", file.Size())
		}
		out += file.Name
		fmt.Println(out)
	}
}
