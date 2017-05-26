package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	rd := bufio.NewReader(os.Stdin)
	max := 0
	for {
		line, err := rd.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("error when reading line: %v", err)
		}

		size := len(line)
		if max < size {
			max = size
		}
	}
	fmt.Println(max)
}
