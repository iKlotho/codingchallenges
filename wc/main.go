package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

type WCFile struct {
	Name string
	fd   *os.File
}

func (wcf *WCFile) loadFile() error {
	fd, err := os.Open(wcf.Name)
	wcf.fd = fd
	return err
}

func (wcf *WCFile) Size() int64 {
	fileInfo, err := os.Stat(wcf.Name)
	if err != nil {
		panic(err)
	}
	return fileInfo.Size()
}

func (wcf *WCFile) LineCount() int {
	// TODO: optimize?
	wcf.fd.Seek(0, 0)
	lc := 0
	reader := bufio.NewReader(wcf.fd)
	for {
		_, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}
		lc++
	}
	return lc
}

func (wcf *WCFile) WordCount() int {
	wcf.fd.Seek(0, 0)
	wordCount := 0
	scanner := bufio.NewScanner(wcf.fd)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		wordCount++
	}
	return wordCount
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
	lineFlagPtr := flag.Bool("l", false, "Print number of lines in each input file")
	wordFlagPtr := flag.Bool("w", false, "Print number of words in each input file")
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
		if *lineFlagPtr {
			out += fmt.Sprintf("%d ", file.LineCount())
		}
		if *wordFlagPtr {
			out += fmt.Sprintf("%d ", file.WordCount())
		}
		if *byteFlagPtr {
			out += fmt.Sprintf("%d ", file.Size())
		}
		out += file.Name
		fmt.Println(out)
	}
}
