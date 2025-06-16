package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Quiz struct {
	Question string
	Answer   string
}

func ReadFile(filePath string) ([]Quiz, error) {
	if filePath == "" {
		return nil, nil
	}
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var quizzes []Quiz
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			continue
		}
		quizzes = append(quizzes, Quiz{
			Question: parts[0],
			Answer:   parts[1],
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return quizzes, nil
}
func getInput(input chan string) {
	for {
		in := bufio.NewReader(os.Stdin)
		result, err := in.ReadString('\n')
		if err != nil {
			log.Fatal("Error reading input:", err)
			continue
		}
		input <- result
	}
}
func startQuiz(quiz []Quiz, duration time.Duration) int {
	input := make(chan string)
	timer := time.After(duration * time.Second)
	go getInput(input)
	fmt.Println("Starting Quiz...")
	totalScore := 0
	for _, q := range quiz {
		score, err := checkAnswer(q.Question, q.Answer, input, timer)
		if err != nil {
			fmt.Println("Error:", err)
			break
		} else if score == 0 {
			fmt.Println("Wrong Answer!")
		} else {
			totalScore += score
		}
	}
	return totalScore
}
func checkAnswer(question string, answer string, input <-chan string, timer <-chan time.Time) (int, error) {
	fmt.Println(question)
	select {
	case <-timer:
		return 0, errors.New("time is up")
	case ans := <-input:
		if strings.Compare(strings.TrimSpace(ans), strings.TrimSpace(answer)) == 0 {
			fmt.Println("Correct Answer!")
			return 1, nil
		} else {
			return 0, nil
		}
	}
}

func init() {
	flag.Int("time", 30, "Time limit for all questions in seconds")
	flag.Int("s", 0, "Shuffle the question order (default: 0)")
	flag.Parse()
	fmt.Println("no. of flags", flag.NFlag())
	if flag.NFlag() == 0 {
		fmt.Println("No flags provided. Using default time limit of 30 seconds.")
	}
	if flag.Lookup("time") == nil {
		fmt.Println("No time flag provided. Using default time limit of 30 seconds.")
	} else {
		fmt.Println("Time limit set to:", flag.Lookup("time").Value)
	}

}

func main() {
	quizzes, err := ReadFile("problems.csv")
	if err != nil {
		log.Fatal("Error reading file:", err)
		return
	}
	if len(quizzes) == 0 {
		log.Fatal("No quizzes found in the file.")
		return
	}
	fmt.Println("Shuffle flag:", flag.Lookup("s").Value)
	if flag.Lookup("s") != nil && flag.Lookup("s").Value.(flag.Getter).Get().(int) == 1 {
		for i := range quizzes {
			j := rand.Intn(i + 1)
			quizzes[i], quizzes[j] = quizzes[j], quizzes[i]

		}
	}
	totalScore := startQuiz(quizzes, time.Duration(flag.Lookup("time").Value.(flag.Getter).Get().(int)))
	fmt.Println("Total Score:", totalScore)
}
