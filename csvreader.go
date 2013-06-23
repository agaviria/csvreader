package main

import (
	"flag"
	"fmt"
	"os"
	"text/template"
)

const usageTmpl = `Usage: csvreader -f=<csvpath> [Options]

Basset Reporter manages BSA/AML profile alerts uniformly.
Use "-help" for more information about flag options.

`

var versionStr = "0.2"

func showUsage() {
	fmt.Fprintf(os.Stderr, usageTmpl)
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

var (
	helpFlag = flag.Bool("help", false, "Display Help Menu")
	outFlag  = flag.Bool("o", false, "Display Output in Terminal")
	// bundleFlag  = flag.Bool("b", false, "Bundle ALL rows found in CSV")
	versionFlag = flag.Bool("v", false, "Application Version")
	acctFlag    = flag.String("acct", "", "Search By Account Number")
	fileFlag    = flag.String("f", "", "CSV Path: /csvdata/csvtest.csv")
	outPathFlag = flag.String("fo", "out", "Directory Path For Generated Markdown files")
)

var tmplAlert *template.Template

func main() {
	// setup alert template
	var err error
	tmplAlert, err = template.ParseFiles("alert.tmpl")
	if err != nil {
		panic(err)
	}

	// parse flags
	flag.Parse()

	// show help if requested
	if *helpFlag {
		showUsage()
		os.Exit(0)
	}

	// show version if requested
	if *versionFlag {
		fmt.Printf("Version: %s\n\n", versionStr)
		showUsage()
		os.Exit(0)
	}

	// do actual work
	read(*fileFlag)
}
