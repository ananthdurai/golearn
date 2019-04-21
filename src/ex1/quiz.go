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

var (
	CsvFileLocation = flag.String("csv_location", "problems.csv", "Input csv file location")
	QuizTimer       = flag.Int("quiz_timer", 300, "Total quiz time in seconds")
)

func main() {
	flag.Parse()

	// Load the csv file from the given location
	f, err := os.Open(*CsvFileLocation)

	if err != nil {
		fmt.Errorf("Unable to open the quiz file")
		panic(err)
	}

	defer f.Close()

	// Read the csv file
	lines, err := csv.NewReader(f).ReadAll()

	if err != nil {
		fmt.Errorf("Unable to read the quiz file")
		panic(err)
	}

	// initialize the score to zero.
	var score int

	// start the quiz timer
	ticker := time.NewTicker(time.Second * time.Duration(*QuizTimer))

	// initialize a channel to signal the quiz finished and print the score
	done := make(chan bool)

	// initialize a reader to read answer from the console
	reader := bufio.NewReader(os.Stdin)

	// initialize a goroutine for the quiz
	go func() {
		for _, line := range lines {
			fmt.Println("What is:", line[0])
			for {
				answer, _ := reader.ReadString('\n')
				answer = strings.Replace(answer, "\n", "", -1)
				if strings.Compare(answer, line[1]) == 0 {
					fmt.Println("You are correct")
					score++ // send the score to the signal
				} else {
					fmt.Println("You are wrong")
				}
				break
			}
		}
		done <- true
	}()

	select {
	case <-done:
		fmt.Println("finished You score is:", score)
	case <-ticker.C:
		fmt.Println("time's up! You score is:", score)
	}

}
