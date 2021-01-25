package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type problems struct {
	q string
	a string
}

func parseProblems(lines [][]string) []problems {
	ps := make([]problems, len(lines))
	for i, line := range lines {
		ps[i] = problems{
			q: line[0],
			a: strings.ToLower(strings.Trim(line[1], " ")),
		}
	}
	return ps
}

func calculateScore(usrAns []string, ps []problems) int {
	var corr int
	for i := 0; i < len(ps); i++ {
		if usrAns[i] == "" {
			continue
		}
		usrAns[i] = strings.Replace(usrAns[i], "\n", "", -1)
		usrAns[i] = strings.ToLower(strings.Trim(usrAns[i], " "))
		if usrAns[i] == ps[i].a {
			corr++
		}
	}
	return corr
}

func main() {
	var t int
	var fn string
	flag.StringVar(&fn, "fn", "problems.csv", "Specify filename for the problems in csv format as question,answer")
	flag.IntVar(&t, "time", 30, "Enter duration of timer for each question")
	shfl := flag.Bool("shuffle", false, "To shuffle or not to shuffle")
	flag.Parse()

	fmt.Printf("File: %s\nTimer: %v\nShuffle: %t\n", fn, t, *shfl)

	f, err := os.Open(fn)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer f.Close()
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)

	ps := parseProblems(lines)
	var usrAns = make([]string, len(ps))

	timer := time.NewTimer(time.Duration(t) * time.Second)
	usrCh := make(chan string)

	for i, problem := range ps {
		fmt.Printf("Problem #%d: %s? :  ", i+1, problem.q)
		go func() {
			ans, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("There's an error in your input")
			}
			usrCh <- ans
		}()

		select {
		case <-timer.C:
			fmt.Println("\nTime's Up")
			fmt.Printf("%d/%d answers were correct\n", calculateScore(usrAns, ps), len(lines))
			os.Exit(0)

		case answer := <-usrCh:
			usrAns[i] = answer
		}
	}

	fmt.Printf("%d/%d answers were correct\n", calculateScore(usrAns, ps), len(lines))
}
