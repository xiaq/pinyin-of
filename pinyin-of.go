package main

import (
	"bufio"
	"bytes"
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
	maxLine  = flag.Int("max-line", 32, "Maximum byte size of a line in data file")
	dataFile *os.File
	dataSize int64
)

func main() {
	flag.Parse()
	prepareData()

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

func prepareData() {
	if *dataPath == "" {
		flag.Usage()
		log.Fatalf("must specify -data")
	}

	file, err := os.Open(*dataPath)
	if err != nil {
		log.Fatalf("cannot open data file %s: %v", *dataPath, err)
	}

	// File must start with a leading \n.
	var b [1]byte
	_, err = file.Read(b[:])
	if err != nil {
		log.Fatalf("cannot read first byte of data file: %v", err)
	}
	if b[0] != '\n' {
		log.Fatal("data file does not start with a newline")
	}

	size, err := file.Seek(0, 2)
	if err != nil {
		log.Fatalf("cannot seek to end of data file: %v", err)
	}

	// File must end with a trailing \n.
	_, err = file.ReadAt(b[:], size-1)
	if err != nil {
		log.Fatalf("cannot read last byte of data file: %v", err)
	}
	if b[0] != '\n' {
		log.Fatal("data file does not end with a newline")
	}

	dataFile = file
	dataSize = size
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
	// Do a binary search on the file.
	b := make([]byte, *maxLine)
	low, high := int64(1), dataSize-1
	for low < high {
		mid := low + (high-low)/2
		readAt(b, mid)
		// Find the next newline and call it the real mid.
		i := findNewline(b)
		realmid := mid + i

		if realmid > high {
			// file[high] must be a '\n', so we should not have moved past it.
			// This must be a program bug.
			panic("realmid > high, program bug")
		} else if realmid == high {
			// file[mid:high] does not have any newlines. Seek backwards for the
			// last newline and call that the real mid instead.
			off := mid - int64(*maxLine)
			size := int64(*maxLine)
			if off < 0 {
				// If we are too very close to the start of the file, be careful
				// to only look as far as to the start.
				off = 0
				size = mid
			}
			readAt(b[:size], off)
			i := findLastNewline(b[:size])
			realmid = off + i
			if realmid < low-1 {
				// file[low-1] must be a '\n', so we should not have moved past
				// it. This must be a program bug.
				panic("realmid < low-1, program bug")
			}
			// We succesfully found another newline. Go ahead.
		}

		// We found a newline. The entry immediately after it is our middle
		// entry, decode it.
		readAt(b, realmid+1)
		i = findNewline(b)
		line := b[:i]

		char, size := utf8.DecodeRune(line)
		switch {
		case r == char:
			// This is what we are looking for. Decode the rest of the line
			// and return the results.
			rest := line[size:]
			return strings.Split(string(rest), ",")
		case r < char:
			// Continue searching in file[low:realmid].
			high = realmid
		case r > char:
			// Continue searching in file[realmid+1+i:high].
			low = realmid + 1 + i
		}
	}
	// We didn't find anything.
	log.Fatalf("no pinyin found for character %c", r)
	panic("unreachable")
}

func readAt(b []byte, offset int64) {
	_, err := dataFile.ReadAt(b, offset)
	if err != nil && err != io.EOF {
		log.Fatalf("cannot read data file at %d: %v", offset, err)
	}
}

func findNewline(b []byte) int64 {
	i := bytes.IndexByte(b, '\n')
	if i == -1 {
		log.Fatalf("data file has line longer than %d, specify a correct -max-line option\n", *maxLine)
	}
	return int64(i)
}

func findLastNewline(b []byte) int64 {
	i := bytes.LastIndexByte(b, '\n')
	if i == -1 {
		log.Fatalf("data file has line longer than %d, specify a correct -max-line option\n", *maxLine)
	}
	return int64(i)
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
