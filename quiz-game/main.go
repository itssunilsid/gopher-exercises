package main

import (
	"fmt"
	"time"
)

func takeInput(sendToChan chan string) {
	var answer string
	fmt.Scanln(&answer)
	sendToChan <- answer
}

func askQuestions(questionsAndAnswers [][]string, stop chan bool, startTimer chan bool, scoreChan chan int, inputChan chan string) {
	counter := 0

	fmt.Println("press enter to start quiz")
	fmt.Scanln()
	startTimer <- true
	for _, qAndA := range questionsAndAnswers {
		fmt.Println(qAndA[0])
		go takeInput(inputChan)

		select {
		case <-stop:
			scoreChan <- counter
			return
		case answer := <-inputChan:
			if answer == qAndA[1] {
				counter++
			}
		}
	}
	scoreChan <- counter
}

func runTimer(quizTimeout int, startTimer chan bool, trigger chan bool) {
	<-startTimer
	<-time.NewTimer(time.Duration(quizTimeout) * time.Second).C
	trigger <- true
}

func main() {
	//TODO add these as command line arguments
	fileName := "problems.csv"
	quizTimeout := 3
	csvContents := readCsvFile(fileName)
	stop := make(chan bool)
	startTimer := make(chan bool)
	scoreChan := make(chan int)
	inputChan := make(chan string)
	go runTimer(quizTimeout, startTimer, stop)
	go askQuestions(csvContents, stop, startTimer, scoreChan, inputChan)
	fmt.Println("score is ", <-scoreChan)
}
