package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"unicode"
)

type WCResult struct {
	lineCount int64
	wordCount int64
	runeCount int64
	charCount int64
}

type WCFile struct {
	fd *os.File
}

func (wcf *WCFile) Calculate() *WCResult {
	wcResult := &WCResult{}
	reader := bufio.NewReader(wcf.fd)
	var prevRune rune
	for {
		r, size, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				if !unicode.IsSpace(prevRune) {
					wcResult.wordCount++
				}
				break
			} else {
				panic(err)
			}
		} else {
			if r == '\n' {
				wcResult.lineCount++
			}
			if unicode.IsSpace(r) && !unicode.IsSpace(prevRune) {
				wcResult.wordCount++
			}
			prevRune = r
			wcResult.runeCount += int64(size)
			wcResult.charCount++
		}
	}
	return wcResult
}

func NewFile(file *os.File) *WCFile {
	wcf := &WCFile{fd: file}
	return wcf
}

func isPipe() bool {
	fi, _ := os.Stdin.Stat()
	return (fi.Mode() & os.ModeCharDevice) == 0
}

func main() {

	byteFlagPtr := flag.Bool("c", false, "Print number of bytes in each input file")
	lineFlagPtr := flag.Bool("l", false, "Print number of lines in each input file")
	wordFlagPtr := flag.Bool("w", false, "Print number of words in each input file")
	charFlagPtr := flag.Bool("m", false, "Print number of characters in each input file")
	flag.Parse()

	// TODO: find a better way
	// If no option provided
	if !*byteFlagPtr && !*lineFlagPtr && !*wordFlagPtr && !*charFlagPtr {
		*byteFlagPtr = true
		*lineFlagPtr = true
		*wordFlagPtr = true
	}

	allFiles := flag.Args()
	files := make([]*WCFile, 0)

	if isPipe() {
		files = append(files, NewFile(os.Stdin))
	} else {
		for _, filename := range allFiles {
			fd, err := os.Open(filename)
			if err != nil {
				panic(err)
			}
			files = append(files, NewFile(fd))
		}
	}

	for _, file := range files {
		wcResult := file.Calculate()
		out := "  "
		if *lineFlagPtr {
			out += fmt.Sprintf("%d ", wcResult.lineCount)
		}
		if *wordFlagPtr {
			out += fmt.Sprintf("%d ", wcResult.wordCount)
		}
		if *byteFlagPtr || *charFlagPtr {
			if *charFlagPtr {
				out += fmt.Sprintf("%d ", wcResult.charCount)
			} else {
				out += fmt.Sprintf("%d ", wcResult.runeCount)

			}
		}
		filename := file.fd.Name()
		if filename != "/dev/stdin" {
			out += filename
		}
		fmt.Println(out)
	}
}
