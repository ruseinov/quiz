package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
)

type Result struct {
	Error  error
	Result string
}

var filePath string
var wordsByLength []string
var wordsMap map[string]int
var reply chan Result

func init() {
	reply = make(chan Result, 0)

	flag.StringVar(&filePath, "filePath", "../word.list", "Specify filepath to fetch word list from")
	flag.Parse()
}

func main() {
	var err error
	wordsByLength, wordsMap, err = readFile()

	if err != nil {
		log.Fatal(err)
	}

	sort.Sort(ByLength(wordsByLength))

	go getLongestWord()

	res := <-reply

	if res.Error != nil {
		log.Fatal(res.Error)
	}

	fmt.Printf("Longest word is: %s\n", res.Result)
}

func lookForWords(originalWord string, lowerBound int, upperBound int) (string, error) {
	for upperBound > 0 {
		if lowerBound >= upperBound {
			break
		}

		word := originalWord[lowerBound:upperBound]
		if _, ok := wordsMap[word]; ok {
			if upperBound != len(originalWord) {
				_, err := lookForWords(originalWord, upperBound, len(originalWord))
				if err != nil {
					upperBound--
				}
			} else {
				reply <- Result{Result: originalWord}
				return word, nil
			}
		} else {
			upperBound--
		}
	}
	return "", errors.New("Unable to find a match")
}

func getLongestWord() {
	currentIndex := 0

	for true {
		if currentIndex == len(wordsByLength)-1 {
			reply <- Result{Error: errors.New("Unable to find a match")}
			break
		}
		var words []string

		for i := len(wordsByLength) - currentIndex - 1; i > 0; i-- {
			word := wordsByLength[currentIndex]

			if len(words) > 0 && len(word) < len(words[0]) {
				break
			}

			words = append(words, word)
			currentIndex++

			lookForWords(word, 0, len(word)-1)
		}
	}
}

/*
* Filepath comes from global scope since it's a flag
 */
func readFile() (words []string, wordsMap map[string]int, err error) {
	file, err := os.Open(filePath)

	if err != nil {
		return
	}

	wordsMap = make(map[string]int)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		words = append(words, scanner.Text())
		wordsMap[scanner.Text()] = 1
	}

	defer file.Close()

	return
}
