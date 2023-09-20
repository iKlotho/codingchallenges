package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

type WCOutput interface {
	PrintResult()
}

type WCFile struct {
	Name   string
	reader *bufio.Reader
}

func (wcf *WCFile) loadFile() error {
	fd, err := os.Open(wcf.Name)
	reader := bufio.NewReader(fd)
	wcf.reader = reader
	return err
}

func (wcf *WCFile) Size() int {
	return wcf.reader.Size()
}

func (wcf *WCFile) LineCount() int {
	// TODO: optimize?
	lc := 0
	for {
		_, err := wcf.reader.ReadBytes('\n')
		if err != nil {
			break
		}
		lc++
	}
	return lc
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
	lineFlagPtr := flag.Bool("l", false, "Print number of line in each input file")
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
		if *lineFlagPtr {
			out += fmt.Sprintf("%d ", file.LineCount())
		}
		out += file.Name
		fmt.Println(out)
	}
}
