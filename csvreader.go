package main

import (
	"flag"
	"fmt"
	"os"
	"text/template"
)

var versionStr = "0.1.0"

func showUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s -f=<csvpath> <Options>\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\n\n")
}

var (
	helpFlag    = flag.Bool("help", false, "Display Help Menu")
	outFlag     = flag.Bool("o", false, "Display Output in Terminal")
	mdFlag      = flag.Bool("m", false, "Generate Markdown File (.md)")
	bundleFlag  = flag.Bool("b", false, "Bundle ALL rows found in CSV")
	versionFlag = flag.Bool("v", false, "Application Version")
	acctFlag    = flag.String("acct", "", "Search By Account Number")
	fileFlag    = flag.String("f", "", "CSV Path: /csvdata/csvtest.csv")
	outPathFlag = flag.String("fo", "out", "Directory Path For Generated Markdown files")
)

var tmpl *template.Template

func init() {
	var err error
	tmpl, err = template.ParseFiles("alert.tmpl")
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Usage = showUsage
	flag.Parse()
	if *helpFlag || *fileFlag == "" {
		flag.Usage()
		os.Exit(0)
	}
	if *versionFlag || *fileFlag == "" {
		fmt.Printf("Version: %s\n\n", versionStr)
		flag.Usage()
		os.Exit(0)
	}
	read(*fileFlag)
}
