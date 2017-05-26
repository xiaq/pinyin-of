package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"unicode/utf8"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	lineno := 0
	for {
		lineno++
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("error when reading line %d: %v", lineno, err)
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
		if size != len(fields[0]) {
			continue
		}

		// Write out character.
		fmt.Print(fields[0])

		// Write out every pinyin.
		for _, f := range fields[2:] {
			if i := strings.Index(f, ":"); i > -1 {
				f = f[:i]
			}
			// Make sure that this only consists of small letter.
			if strings.TrimLeft(f, "abcdefghijklmnopqrstuvwxyz") != "" {
				log.Fatalf("line %d has non-pinyin:", lineno)
			}
			fmt.Print(" " + f)
		}
		fmt.Println()
	}
}
