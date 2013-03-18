package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

// documentation for csv is at http://golang.org/pkg/encoding/csv/
// TODO; could not find
func main() {
	file, error := os.Open("csvdata/csvtest.csv")
	if error != nil {
		// error is printable
		// element passed are separated by space automatically
		fmt.Println("Error: ", error)
		return
	}
	// automatically call Close() at the end of current method
	defer file.Close()
	//
	reader := csv.NewReader(file)
	// options are available at:
	// http://golang.org/src/pkg/encoding/csv/reader.go?s=3213:3671#194
	reader.Comma = ';'
	lineCount := 1
	for {
		// read just one record, but we could ReadAll() as well
		record, error := reader.Read()
		// EOF is fitted into error
		if error == io.EOF {
			break
		} else if error != nil {
			fmt.Println("Error:", error)
			return
		}
		// record is an array of string so is directly printable
		fmt.Println("Record", lineCount, "is", record, "and has", len(record), "fields")
		// iterate on top of that skipping the 1st record
		for i := 0; i < len(record); i++ {
			fmt.Println(" ", record[i])
		}
		fmt.Println()
		lineCount += 1
	}
}
