package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Amounts struct {
	EstimatedAmount float64
	ActualAmount    float64
	Percent         string
	EstimatedTxn    string
}

type Account struct {
	Num  string
	Name string
	Risk string
	In   *Amounts
	Out  *Amounts
}

type alertTemplateData struct {
	Date    time.Time
	Account *Account
}

func (a *Account) IsValidAlert() bool {
	if a.In.IsValidAlert() || a.Out.IsValidAlert() {
		return true
	}
	return false
}

func (a *Amounts) IsValidAlert() bool {
	// Removes false-positive alerts from printable template
	// e.g. -275%
	if strings.Count(a.Percent, "-") != 0 {
		return false
	}
	return true
}

// Returns actual amount of funds in a monthly cycle vs. declared profile
func (a *Amounts) ClientActualAmount() float64 {
	return a.ActualAmount - a.EstimatedAmount
}

func readAmounts(r []string) (a *Amounts, err error) {
	a = new(Amounts)
	est := r[0]
	a.EstimatedAmount, err = strconv.ParseFloat(est, 64)
	if err != nil {
		return nil, fmt.Errorf("Error converting string: %s", err)
	}
	act := r[1]
	a.ActualAmount, err = strconv.ParseFloat(act, 64)
	if err != nil {
		return nil, fmt.Errorf("Error converting string: %s", err)
	}
	a.Percent = r[2]
	if a.Percent == "***" {
		a.Percent = ">1000%"
	}
	a.EstimatedTxn = r[3]
	return a, nil
}

func printMonthlyActivity(w io.Writer, a *Account) error {
	var data = &alertTemplateData{
		Date:    time.Now(),
		Account: a,
	}
	err := tmplAlert.ExecuteTemplate(w, "alert.tmpl", data)
	if err != nil {
		log.Fatalf("execution failed: %s", err)
	}
	return nil
}

// Slice of accounts alerted in a single month retrieved
func accountMonth(record []string) (*Account, error) {
	var err error
	a := new(Account)
	a.In, err = readAmounts(record[5 : 5+6])
	if err != nil {
		return nil, err
	}
	a.Out, err = readAmounts(record[11 : 11+6])
	if err != nil {
		return nil, err
	}
	a.Num = record[2]
	a.Name = record[3]
	a.Risk = record[0]
	// replaces data value with risk representation
	if a.Risk == "10" {
		a.Risk = "Low"
	} else if a.Risk == "20" {
		a.Risk = "Moderate"
	} else {
		a.Risk = "High"
	}
	return a, nil
}

func read(filename string) {
	// open csv file
	csvFile, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening file: %s", err)
	}
	defer csvFile.Close()

	// create buffered reader to skip the first line
	bufferedReader := bufio.NewReader(csvFile)
	_, _, err = bufferedReader.ReadLine()
	if err != nil {
		log.Fatalf("Could not setup buffer for csv reader. %s", err)
	}

	// create new csv.Reader
	csvReader := csv.NewReader(bufferedReader)
	csvReader.Comma = ';'
	err = os.MkdirAll(*outPathFlag, 0700)
	if err != nil {
		log.Fatalf("Error creating output directory %s: %s", *outPathFlag, err)
	}

	// loop over records
	for {
		record, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("Error CSV format %s", err)
		}

		ac, err := accountMonth(record)
		if err != nil {
			log.Fatalf("Error getting account record: %s", err)
		}
		if *acctFlag != "" && *acctFlag != ac.Num {
			continue
		}
		if *outFlag && ac.IsValidAlert() {
			err = printMonthlyActivity(os.Stdout, ac)
		}
		if len(*outPathFlag) > 0 && ac.IsValidAlert() {
			filePath := filepath.Join(*outPathFlag, fmt.Sprintf("%s-%s.md", ac.Name, ac.Num))
			outputFile, err := os.Create(filePath)
			if err != nil {
				log.Fatalf("Error writing to file %s: %s", filePath, err)
			}
			err = printMonthlyActivity(outputFile, ac)
			outputFile.Close()
		}
	}
}
