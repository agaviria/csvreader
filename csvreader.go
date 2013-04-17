package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

const file = "csvdata/csvtest.csv"

type Amounts struct {
	EstimatedAmt float64
	ActualAmt    float64
	Percent      string
	EstimatedTxn string
}

type Account struct {
	Num  string
	Name string
	Risk string
	In   *Amounts
	Out  *Amounts
}

const reportSeparator = "====================================================================\n"

func printIncActivity(a *Account) error {
	p := fmt.Printf
	p(reportSeparator)
	// Next 3 Lines set up date parsing of stdLongMonth (m-1)
	y, m, _ := time.Now().Date()
	t := time.Date(y, m-1, 1, 0, 0, 0, 0, time.UTC)
	t = time.Date(y, m-1, 1, 0, 0, 0, 0, time.UTC)
	p("Account: %+s / %+s\n", a.Num, a.Name)
	p("Risk: %+s / ", a.Risk)
	p("Exception Date: %+v", t.Format("January, 2006\n"))
	p("\nThe account exceeded the incoming profile by $%+v,\n", a.In.Percent)
	p("the same as $%+v over the monthly incoming amount of $%+v.\n", a.In.ActualAmt-a.In.EstimatedAmt, a.In.ActualAmt)
	p("Current profile is established at $%+v with an expectancy of (%+v).\n", a.In.EstimatedAmt, a.In.EstimatedTxn)
	return nil
}

func readAmounts(r []string) (a *Amounts, err error) {
	a = new(Amounts)
	est := r[0]
	a.EstimatedAmt, err = strconv.ParseFloat(est, 64)
	if err != nil {
		return nil, fmt.Errorf("Error converting string: +v", err)
	}
	act := r[1]
	a.ActualAmt, err = strconv.ParseFloat(act, 64)
	if err != nil {
		return nil, fmt.Errorf("Error converting string: +v", err)
	}
	a.Percent = r[2]
	a.EstimatedTxn = r[3]
	return a, nil
}

func accountMonth(record []string) error {
	var err error
	var a Account
	a.Num = record[2]
	a.Name = record[3]
	a.Risk = record[0]
	a.In, err = readAmounts(record[5 : 5+6])
	if err != nil {
		return err
	}
	a.Out, err = readAmounts(record[11 : 11+6])
	if err != nil {
		return err
	}
	err = printIncActivity(&a)
	return err
}

func main() {
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer f.Close()

	rdr := csv.NewReader(f)
	rdr.Comma = ';'

	for {
		record, err := rdr.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		err = accountMonth(record)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf(reportSeparator)
	}
}
