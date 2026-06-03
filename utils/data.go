// Package dataread ...
package dataread

import (
	"encoding/csv"
	"os"
	"regexp"
	"strings"
)

func ReadData() [][]string {
	file, err := os.Open("spamhamdata.csv")
	if err != nil {
		panic("The file could not be read")
	}
	defer func() {
		err = file.Close()
		if err != nil {
			panic("cannot close file")
		}
	}()
	reader := csv.NewReader(file)
	reader.Comma = '\t'
	reader.LazyQuotes = true
	records, err := reader.ReadAll()
	if err != nil {
		panic("Error reading files")
	}
	newData := [][]string{}
	for _, record := range records {
		label := strings.TrimSpace(record[0])
		if label != "spam" && label != "ham" {
			continue
		}
		newData = append(newData, []string{cleanText(record[1]), label})
	}
	return newData
}

func cleanText(text string) string {
	text = strings.ToLower(text)
	text = regexp.MustCompile(`[^a-z\s]`).ReplaceAllString(text, "") // remove non-letters
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")     // collapse whitespace
	return strings.TrimSpace(text)
}
