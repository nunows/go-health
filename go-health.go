package main

import (
	"fmt"
	"os"
	"time"
	"strings"
	"io/ioutil"
	"encoding/xml"
	"encoding/csv"
)

type Query struct {
	ExportDate ExportDate
	Record	[]Record
}

type ExportDate struct {
	Value string `xml:"value,attr"`
}

type Record struct {
	Type string `xml:"type,attr"`
	SourceName string `xml:"sourceName,attr"`
	SourceVersion string `xml:"sourceVersion,attr"`
	Device string `xml:"device,attr"`
	Unit string `xml:"unit,attr"`
	CreationDate string `xml:"creationDate,attr"`
	StartDate string `xml:"startDate,attr"` 
	EndDate string `xml:"endDate,attr"`
	Value string `xml:"value,attr"`
}

func types() []string {
	
	types := []string {
		"HKQuantityTypeIdentifierBodyMass", 
		"HKQuantityTypeIdentifierHeartRate", 
		"HKQuantityTypeIdentifierBodyTemperature", 
		"HKQuantityTypeIdentifierStepCount",
		"HKQuantityTypeIdentifierDistanceWalkingRunning",
		"HKQuantityTypeIdentifierActiveEnergyBurned",
		"HKQuantityTypeIdentifierFlightsClimbed",
		"HKCategoryTypeIdentifierSleepAnalysis" }

	return types
}

func parseData(file string) (Query, error) {
	var q Query

	xmlFile, err := os.Open(file)

	if err != nil {
		return q, err
	}
	
	defer xmlFile.Close()

	data, _ := ioutil.ReadAll(xmlFile)
	xml.Unmarshal(data, &q)

	return q, nil
}

func recordToCsv(r Record) []string {
	return []string {r.Type, r.SourceName, r.SourceVersion, r.Device, r.Unit, r.CreationDate, r.StartDate, r.EndDate, r.Value}
}

func createCSV(typeName string, record []Record, fromDate string) error {
	format := "2006-01-02 15:04:05"
	date, _ := time.Parse(format,fromDate)

	file, err := os.Create(typeName + ".csv")

	if err != nil {
        return err
    }

    defer file.Close()

    writer := csv.NewWriter(file)

    for _, r := range record {
    	creationDate, _ := time.Parse(format,strings.Replace(r.CreationDate, " +0000", "", -1))
    	
    	if r.Type == typeName && creationDate.After(date) {
    		line := recordToCsv(r)
			err := writer.Write(line)
			if err != nil {
				return err
        	}
    	}
    }

    defer writer.Flush()

    return nil
}

func main() {

	var fromDate = "2016-01-01 00:00:00"
	var dataFile = "export.xml"

	fmt.Printf("Parsing %s file... \n", dataFile)
	q, err := parseData(dataFile)

	if err != nil {
		panic(err)
	}
	
	for _, t := range types() {
		fmt.Printf("Exporting %s data... ", t)
		err := createCSV(t, q.Record, fromDate)

		if err != nil {
			fmt.Println("Error", err)
		} else {
			fmt.Println("done!")
		}
	}	
}
