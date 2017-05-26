package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"unicode/utf8"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	lineno := 0
	var chars []rune
	pinyinsOf := map[rune][]string{}
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

		chars = append(chars, char)
		pinyinsOf[char] = pinyins
	}

	sort.Sort(runes(chars))
	for _, char := range chars {
		fmt.Printf("%c %s\n", char, strings.Join(pinyinsOf[char], " "))
	}
}

type runes []rune

func (rs runes) Len() int           { return len(rs) }
func (rs runes) Less(i, j int) bool { return rs[i] < rs[j] }
func (rs runes) Swap(i, j int)      { rs[i], rs[j] = rs[j], rs[i] }
