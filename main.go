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
	csvFilename := flag.String("csv", "problems.csv", "csv file with 'question,answer' format")
	timeLimit := flag.Int64("limit", 30, "Time limit in seconds")
	flag.Parse()

	//Load csv file
	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open file %s\n", *csvFilename))
	}

	//Read csv file
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to read csv file")
	}

	//Parse csv file to struct
	problems := parseLines(lines)

	//Create timer
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	//Check for correct answers
	correct := 0
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.q)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Println("\nTime has run out!")
			finalScore(correct, len(problems))
			return
		case answer := <-answerCh:
			if answer == p.a {
				correct++
			}
		}
	}

	finalScore(correct, len(problems))
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func finalScore(correct int, questionsLength int) {
	fmt.Printf("You score %d out of %d\n", correct, questionsLength)
}
