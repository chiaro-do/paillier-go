package utils

import (
	"encoding/csv"
	"os"
	"log"
	"strconv"
)

func ReadCSV(filePath string) [][]interface{} {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	records = records[1:]
	var result [][]interface{}

	for _, record := range records {
		a, err := strconv.Atoi(record[5])
		l, err := strconv.Atoi(record[6])
		e, err := strconv.Atoi(record[7])
		if err != nil {
			log.Fatalf("Error parsing age: %v\n", err)
		}
		newRecord := []interface{}{record[0], record[1], mapping(a, l, e)}
		result = append(result, newRecord)
	}

	return result
}

func mapping(a int, l int, e int) int {
	return 4 * a + 2 * l + e
}