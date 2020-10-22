package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"fmt"
	"strings"
	_"os"
	_"strconv"
)

// Table : https://raw.githubusercontent.com/izik1/gbops/master/dmgops.json
type Table struct {
	Unprefixed []struct {
		Name            string `json:"Name"`
		Group           string `json:"Group"`
		TCyclesBranch   int    `json:"TCyclesBranch"`
		TCyclesNoBranch int    `json:"TCyclesNoBranch"`
		Length          int    `json:"Length"`
		Flags           struct {
			Z string `json:"Z"`
			N string `json:"N"`
			H string `json:"H"`
			C string `json:"C"`
		} `json:"Flags"`
		TimingNoBranch []struct {
			Type    string `json:"Type"`
			Comment string `json:"Comment"`
		} `json:"TimingNoBranch,omitempty"`
		TimingBranch []struct {
			Type    string `json:"Type"`
			Comment string `json:"Comment"`
		} `json:"TimingBranch,omitempty"`
	} `json:"Unprefixed"`
	CBPrefixed []struct {
		Name            string `json:"Name"`
		Group           string `json:"Group"`
		TCyclesBranch   int    `json:"TCyclesBranch"`
		TCyclesNoBranch int    `json:"TCyclesNoBranch"`
		Length          int    `json:"Length"`
		Flags           struct {
			Z string `json:"Z"`
			N string `json:"N"`
			H string `json:"H"`
			C string `json:"C"`
		} `json:"Flags"`
		TimingNoBranch []struct {
			Type    string `json:"Type"`
			Comment string `json:"Comment"`
		} `json:"TimingNoBranch"`
	} `json:"CBPrefixed"`
}

const (
	// HEADER : app header
	HEADER = "/*parse_ops - 0.0.1: auto generate C/C++ array from table data*/\n\n"
	// INFILE : table data to parse
	INFILE = "table.json"
	// OUTFILE : file output
	OUTFILE = "cycle_table.h"
)

// todo: pass cmd args
func main() {
	tableData, err := ioutil.ReadFile(INFILE)
	if err != nil {
		log.Fatal(err)
	}

	var table Table
	err = json.Unmarshal(tableData, &table)
	if err != nil {
		log.Fatal(err)
	}

	for i,v := range table.Unprefixed{
		if strings.Contains(v.Name,"AND "){ //&& strings.Contains(v.Name, "u16") {//|| strings.Contains(v.Name, "DEC"){
				fmt.Printf("%08b: %s\n",i,v.Name)
		}
	}
	// ayyyyyy
	println("done!")
}
