package main

import (
	"fmt"
	"strings"

	dataread "github.com/billuop123/spam-ham/utils"
)

func main() {
	totalData := dataread.ReadData()
	limit := len(totalData) - int(0.2*float64(len(totalData)))
	testData := totalData[limit:]
	_ = testData
	trainingData := totalData[:limit]
	bow := map[string]int{}
	// spambow
	spamBow := map[string]int{}
	// ham bow
	hamBow := map[string]int{}
	bowCount(trainingData, bow, spamBow, hamBow)
	spamCount, hamCount, spamMsgs, hamMsgs := spamHamCount(trainingData)
	totalCount := spamMsgs + hamMsgs
	spamProb := float64(spamMsgs) / float64(totalCount)
	hamProb := float64(hamMsgs) / float64(totalCount)
	eval(testData, spamCount, hamCount, spamBow, hamBow, spamProb, hamProb)
}

func spamHamCount(data [][]string) (int, int, int, int) {
	spamCount := 0
	hamCount := 0
	spamMsgs := 0
	hamMsgs := 0
	for _, item := range data {
		words := strings.Split(item[0], " ")
		switch item[1] {
		case "spam":
			spamCount += len(words)
			spamMsgs++
		case "ham":
			hamCount += len(words)
			hamMsgs++
		default:
			fmt.Println("class not found")
		}
	}
	return spamCount, hamCount, spamMsgs, hamMsgs
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

func testDecider(testData string, spamCount, hamCount int, spamBow, hamBow map[string]int, spamProb, hamProb float64) (bool, bool) {
	finalSpamProb, finalHamProb := calculateProb(testData, spamCount, hamCount, spamBow, hamBow, spamProb, hamProb)
	if finalSpamProb > finalHamProb {
		return true, false
	} else {
		return false, true
	}
}

func eval(testData [][]string, spamCount, hamCount int, spamBow, hamBow map[string]int, spamProb, hamProb float64) {
	truePositive := 0.0
	trueNegative := 0.0
	falsePositive := 0.0
	falseNegative := 0.0
	for _, test := range testData {
		val1, _ := testDecider(test[0], spamCount, hamCount, spamBow, hamBow, spamProb, hamProb)
		if val1 {
			if test[1] == "spam" {
				truePositive++
			} else {
				falsePositive++
			}
		} else {
			if test[1] == "ham" {
				trueNegative++
			} else {
				falseNegative++
			}
		}
	}
	accuracy := (truePositive + trueNegative) / (truePositive + trueNegative + falsePositive + falseNegative)
	recall := truePositive / (truePositive + falseNegative)
	precision := truePositive / (truePositive + falsePositive)
	fmt.Println("Confusion Matrix:")
	fmt.Println("----------------------------------------------")
	fmt.Println("|           |  Predicted+         |Predicted- |")
	fmt.Println("|-----------|---------------------------------|")
	fmt.Printf("|   Actual+ |-----%.0f------------%.0f------------|\n", truePositive, falseNegative)
	fmt.Printf("|   Actual- |-----%.0f------------%.0f-----------|\n", falsePositive, trueNegative)
	fmt.Println("|---------------------------------------------|")
	fmt.Printf("accuracy %.2f\n", accuracy)
	fmt.Printf("recall %.2f\n", recall)
	fmt.Printf("precision %.2f\n", precision)
}
