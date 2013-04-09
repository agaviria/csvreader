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
		var num string = (record[2])
		var name string = (record[3])
		var eia string = (record[5])
		var cia string = (record[6])
		var inc_percent = (record[7])
		var estIncTxn string = (record[8])
		var inc_diff float64

		for i := 0; i < len(record[i]); i++ {
			estInc, err := strconv.ParseFloat(eia, 64)
			actInc, err := strconv.ParseFloat(cia, 64)
			inc_diff = (actInc - estInc)
			if err == nil {
				fmt.Println("==============================================================================\n")
				fmt.Printf("Account: %+s - %+s exceeded the IncAmt by %+v same as $%+v\n", num, name, inc_percent, inc_diff)
				fmt.Printf("over the monthly incoming amount of $%+v. Currently, the declared\n", actInc)
				fmt.Printf("profile is established at $%+v with an expectancy of (%+v).\n", estInc, estIncTxn)
			} else {
				log.Fatalf("Error converting strings: +v", err)
			}

		}
		fmt.Println()
	}
}
