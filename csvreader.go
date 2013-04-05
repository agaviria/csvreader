package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

// documentation for csv is at http://golang.org/pkg/encoding/csv/
// TODO; refactor
func main() {
	file, err := os.Open("csvdata/csvtest.csv")
	if err != nil {
		log.Fatalf("Error reading all lines: %v", err)
	}
	// automatically call Close() at the end of current method
	defer file.Close()

	// options are available at:
	// http://golang.org/src/pkg/encoding/csv/reader.go?s=3213:3671#194
	for {
		reader := csv.NewReader(file)
		reader.Comma = ';'
		record, err := reader.Read()
		if err != nil {
			log.Print(err)
			os.Exit(-1)
		}

		var eia string = (record[5])
		var cia string = (record[6])
		var inc_percent = (record[7])
		var inc_diff float64
		// for loop needs to be looked at!!!
		for i := 0; i < len(record[i]); i++ {
			estInc, err := strconv.ParseFloat(eia, 64)
			if err == nil {
				fmt.Printf("Estimated Incoming Amount: $%+v\n", estInc)
			}
			actInc, err := strconv.ParseFloat(cia, 64)
			if err == nil {
				fmt.Printf("Actual Customer Activity: $%+v\n", actInc)
			}
			inc_diff = (actInc - estInc)
			fmt.Printf("The account exceeded the incoming amount by %+v same as $%+v", inc_percent, inc_diff)
		}
		fmt.Println()
	}
}
