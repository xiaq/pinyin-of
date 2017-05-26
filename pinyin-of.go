package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	dataPath = flag.String("data", "", "Path to pinyin data file")
	dataFile *os.File
)

func main() {
	flag.Parse()
	if *dataPath == "" {
		flag.Usage()
		log.Fatalf("must specify -data")
	}
	var err error
	dataFile, err = os.Open(*dataPath)
	if err != nil {
		log.Fatalf("cannot open pinyin data file %s: %v", *dataPath, err)
	}

	args := flag.Args()
	if len(args) == 0 {
		rd := bufio.NewReader(os.Stdin)
		lineno := 0
		for {
			lineno++
			line, err := rd.ReadString('\n')
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatalf("error when reading line %d: %v", lineno, err)
			}
			convert(line[:len(line)-1], fmt.Sprintf("line %d", lineno))
		}
	} else {
		for i, arg := range args {
			convert(arg, fmt.Sprintf("arg %d", i))
		}
	}
}

var m = map[rune][]string{
	'你': []string{"ni"},
	'的': []string{"de", "di"},
}

func convert(word, what string) {
	for _, r := range word {
		if !unicode.In(r, unicode.Han) {
			log.Fatalf("%s %s contains non-Han characters", what, word)
		}
	}
	var pinyins [][]string
	for _, r := range word {
		pinyins = append(pinyins, find(r))
	}
	printAll(pinyins, 0, "")
}

func find(r rune) []string {
	_, err := dataFile.Seek(0, 0)
	if err != nil {
		log.Fatalf("cannot seek pinyin data file: %v", err)
	}
	rd := bufio.NewReader(dataFile)
	lineno := 0
	for {
		lineno++
		line, err := rd.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("error when reading line %d of pinyin data file: %v",
				lineno, err)
		}

		// We are looking for lines like
		// 的 4886 de:99.9671% di:0.0329%
		// or
		// 梀 4356 su yin
		// Namely, one character followed by a number (we drop this field), and
		// pinyins with optional colon and probability.
		line = line[:len(line)-1]
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		char, size := utf8.DecodeRuneInString(fields[0])
		if size < len(fields[0]) {
			// We have gone past the section of single characters.
			break
		}

		if char == r {
			// Collect all pinyins of this character.
			pinyins := make([]string, len(fields)-2)
			for i, f := range fields[2:] {
				if j := strings.Index(f, ":"); j > -1 {
					f = f[:j]
				}
				// Make sure that this only consists of small letter.
				if strings.TrimLeft(f, "abcdefghijklmnopqrstuvwxyz") != "" {
					log.Fatalf("line %d of pinyin data file has non-pinyin:",
						lineno)
				}
				pinyins[i] = f
			}
			return pinyins
		}
	}
	log.Fatal("No pinyin data for character %c", r)
	panic("unreachable")
}

func printAll(pinyins [][]string, i int, acc string) {
	(&printer{pinyins, 0, false}).print("")
	fmt.Println()
}

type printer struct {
	pinyins [][]string
	depth   int
	printed bool
}

func (p *printer) print(acc string) {
	if p.depth == len(p.pinyins) {
		if p.printed {
			fmt.Print(" ")
		} else {
			p.printed = true
		}
		fmt.Print(acc)
		return
	}

	p.depth++
	for _, choice := range p.pinyins[p.depth-1] {
		p.print(acc + choice)
	}
	p.depth--
}
