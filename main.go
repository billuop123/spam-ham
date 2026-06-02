package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	dataread "github.com/billuop123/spam-ham/utils"
)

func main() {
	trainingData := dataread.ReadData()
	var inputString string
	// bow
	fmt.Println("Enter a text:")
	reader := bufio.NewReader(os.Stdin)
	inputString, _ = reader.ReadString('\n')
	inputString = strings.TrimSpace(inputString)
	bow := map[string]int{}
	// spambow
	spamBow := map[string]int{}
	// ham bow
	hamBow := map[string]int{}
	bowCount(trainingData, bow, spamBow, hamBow)
	spamCount, hamCount := spamHamCount(trainingData)
	totalCount := spamCount + hamCount
	spamProb := float64(spamCount) / float64(totalCount)
	hamProb := float64(hamCount) / float64(totalCount)
	testDecider(inputString, spamCount, hamCount, spamBow, hamBow, spamProb, hamProb)
}

func spamHamCount(data [][]string) (int, int) {
	spamCount := 0
	hamCount := 0
	for _, item := range data {
		words := strings.Split(item[0], " ")
		switch item[1] {
		case "spam":
			spamCount += len(words)
		case "ham":
			hamCount += len(words)
		default:
			fmt.Println("class not found")
		}
	}
	return spamCount, hamCount
}

func calculateProb(data string, spamCount, hamCount int, spamBow, hamBow map[string]int, spamProb, hamProb float64) (float64, float64) {
	splittedData := strings.Split(data, " ")
	totalSpamProb := 1.0
	totalHamProb := 1.0
	for _, singleString := range splittedData {
		spamVal := spamBow[singleString]
		hamVal := hamBow[singleString]
		totalSpamProb = totalSpamProb * (float64(spamVal+1) / float64(spamCount+len(spamBow)))
		totalHamProb = totalHamProb * (float64(hamVal+1) / float64(hamCount+len(hamBow)))
	}
	totalSpamProb = spamProb * totalSpamProb
	totalHamProb = hamProb * totalHamProb
	return totalSpamProb, totalHamProb
}

func bowCount(trainingData [][]string, bow, spamBow, hamBow map[string]int) {
	for _, data := range trainingData {
		splitted := strings.SplitSeq(data[0], " ")
		for word := range splitted {
			if word == "" {
				continue
			}
			bowVal, ok := bow[word]
			if ok {
				bow[word] = bowVal + 1
			} else {
				bow[word] = 1
			}
			switch data[1] {
			case "spam":
				val, ok := spamBow[word]
				if ok {
					spamBow[word] = val + 1
				} else {
					spamBow[word] = 1
				}
			case "ham":
				val, ok := hamBow[word]
				if ok {
					hamBow[word] = val + 1
				} else {
					hamBow[word] = 1
				}
			default:
				fmt.Println("class not found")
			}
		}
	}
}

func testDecider(testData string, spamCount, hamCount int, spamBow, hamBow map[string]int, spamProb, hamProb float64) {
	finalSpamProb, finalHamProb := calculateProb(testData, spamCount, hamCount, spamBow, hamBow, spamProb, hamProb)
	if finalSpamProb > finalHamProb {
		fmt.Println("This is a Spam")
	} else {
		fmt.Println("This is ham")
	}
}
