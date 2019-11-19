package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Must pass the Changelog file name as the first argument")
	}
	file := os.Args[1]
	data, err := os.Open(file)
	if err != nil {
		log.Fatalln(err)
	}
	scanner := bufio.NewScanner(data)
	var startPrint bool
	tagLinePrefix := "## v"
	for scanner.Scan() {
		line := scanner.Text()
		if startPrint && !strings.HasPrefix(line, tagLinePrefix) {
			fmt.Println(scanner.Text())
		}
		if strings.HasPrefix(line, tagLinePrefix) {
			if startPrint {
				// Exit at the next matching tag
				break
			}
			startPrint = true
		}
	}
	if err = scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

}
