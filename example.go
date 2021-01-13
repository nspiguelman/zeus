package main

import (
	"fmt"
	"time"
)

func question(answerCh chan string, doneCh chan bool, timeout int) {
	timer := time.NewTimer(time.Duration(timeout) * time.Second)
	for i := 1; ; i++ {
		select {
		case <-timer.C:
			fmt.Println("Time is up!")
			doneCh <- true
			return
		case answer := <-answerCh:
			// do something
			fmt.Printf("Answer #%d: %s\n", i, answer)
		}
	}
}

func main() {
	answerCh := make(chan string)
	doneCh := make(chan bool)

	fmt.Println("What does the fox say?")
	go question(answerCh, doneCh, 5)

	for {
		var answer string
		fmt.Scanf("%s\n", &answer)
		select {
		case <-doneCh:
			fmt.Println("Bye!")
			return
		default:
			answerCh <- answer
		}
	}
}

