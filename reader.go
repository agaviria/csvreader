package main

import (
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

func (a *Account) IsPrintable() bool {
	if *bundleFlag {
		return true
	}
	if a.In.IsPrintable() || a.Out.IsPrintable() {
		return true
	}
	return false
}

func (a *Amounts) IsPrintable() bool {
	if *bundleFlag {
		return true
	}
	// Removes false-positive alerts from printable template
	// e.g. -275%
	if strings.Count(a.Percent, "-") != 0 {
		return false
	}
	return true
}

// Returns actual amount of funds in a monthly cycle vs. declared profile
func (a *Amounts) ClientActualAmt() float64 {
	return a.ActualAmt - a.EstimatedAmt
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

func printMthlyActivity(w io.Writer, a *Account) error {
	var context = map[interface{}]interface{}{"date": time.Now(), "account": a}
	err := tmpl.ExecuteTemplate(w, "alert.tmpl", context)
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

func read(file string) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer f.Close()

	rdr := csv.NewReader(f)
	rdr.Comma = ';'
	err = os.MkdirAll(*outPathFlag, 0700)
	if err != nil {
		panic(err)
	}
	for {
		record, err := rdr.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("Error CSV format %v", err)
		}
		ac, err := accountMonth(record)
		if *acctFlag != "" && *acctFlag != ac.Num {
			continue
		}
		if err != nil {
			panic(err)
		}
		if *outFlag && ac.IsPrintable() {
			err = printMthlyActivity(os.Stdout, ac)
		}
		if *mdFlag && ac.IsPrintable() {
			fo, err := os.Create(filepath.Join(*outPathFlag, fmt.Sprintf("%s-%s.md", ac.Name, ac.Num)))
			if err != nil {
				panic(err)
			}
			err = printMthlyActivity(fo, ac)
			fo.Close()
		}
	}
}
