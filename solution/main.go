package main

import (
	"bufio"
	"flag"
	"github.com/ruseinov/quiz/solution/sort"
	"log"
	"os"
)

var filePath string

func main() {
	words, err := readFile()

	wordsByLength := words
	sort.Sort(ByLength(wordsByLength))

}

/*
* Filepath comes from global scope since it's a flag
 */
func readFile() (words, err) {
	file, err := os.Open(filePath)

	if err != nil {
		return
	}

	var words []string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	defer file.Close()

	return
}

func init() {
	flag.StringVar(&filePath, "filePath", "../word.list", "Specify filepath to fetch word list from")
}
