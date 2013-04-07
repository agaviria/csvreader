package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

const file = "csvdata/csvtest.csv"

func main() {
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("Error reading all lines: %v", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.Comma = ';'

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Print(err)
			os.Exit(-1)
		}
		// need to refactor variables
		var eia string = (record[5])
		var cia string = (record[6])
		var inc_percent = (record[7])
		var inc_diff float64
		
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
