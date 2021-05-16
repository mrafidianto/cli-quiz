package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format 'question,answer'")
	timeLimit := flag.Int64("limit", 30, "the time of the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		fmt.Printf("Failed to open file : %s\n", *csvFilename)
		os.Exit(1)
	}
	
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse provided CSV file")
	}

	problems := ParseLines(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correct := 0;

	problemloop:
	for i, problem := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, problem.question)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Println()
			break problemloop
		case answer := <-answerCh: 
			if answer == problem.answer {
				correct++
			}
		}
	}

	fmt.Printf("You scored %d out of %d", correct, len(problems))
}

func ParseLines(lines [][]string) []problem {
	
	ret := make([]problem, len(lines))
	
	for i, line := range lines {
		ret[i] = problem {
			question: line[0],
			answer: strings.TrimSpace(line[1]),
		}
	}

	return ret
}

type problem struct {
	question string
	answer string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}