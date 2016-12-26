package main

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type Query struct {
	ExportDate ExportDate
	Record     []Record
}

type ExportDate struct {
	Value string `xml:"value,attr"`
}

type Record struct {
	Type          string `xml:"type,attr"`
	SourceName    string `xml:"sourceName,attr"`
	SourceVersion string `xml:"sourceVersion,attr"`
	Device        string `xml:"device,attr"`
	Unit          string `xml:"unit,attr"`
	CreationDate  string `xml:"creationDate,attr"`
	StartDate     string `xml:"startDate,attr"`
	EndDate       string `xml:"endDate,attr"`
	Value         string `xml:"value,attr"`
}

func types() map[string]string {

	types := map[string]string{
		"HKQuantityTypeIdentifierBodyMass":               "1",
		"HKQuantityTypeIdentifierHeartRate":              "2",
		"HKQuantityTypeIdentifierBodyTemperature":        "3",
		"HKQuantityTypeIdentifierStepCount":              "4",
		"HKQuantityTypeIdentifierDistanceWalkingRunning": "5",
		"HKQuantityTypeIdentifierActiveEnergyBurned":     "6",
		"HKQuantityTypeIdentifierFlightsClimbed":         "7",
		"HKCategoryTypeIdentifierSleepAnalysis":          "8",
		"HKQuantityTypeIdentifierBodyMassIndex":          "9",
		"HKQuantityTypeIdentifierHeight":                 "10",
		"HKQuantityTypeIdentifierBloodPressureDiastolic": "11",
		"HKQuantityTypeIdentifierBloodPressureSystolic":  "12",
	}

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
	var v string

	if r.Value == "HKCategoryValueSleepAnalysisAsleep" {
		v = "1"
	} else if r.Value == "HKCategoryValueSleepAnalysisInBed" {
		v = "0"
	} else {
		v = r.Value
	}

	return []string{types()[r.Type], r.SourceName, r.SourceVersion, r.Device, r.Unit, r.CreationDate, r.StartDate, r.EndDate, v}
}

func createCSV(all bool, typeName string, record []Record, fromDate string) error {
	format := "2006-01-02 15:04:05"
	date, _ := time.Parse(format, fromDate)

	file, err := os.Create(typeName + ".csv")

	if err != nil {
		return err
	}

	defer file.Close()

	writer := csv.NewWriter(file)

	for _, r := range record {
		creationDate, _ := time.Parse(format, strings.Replace(r.CreationDate, " +0000", "", -1))

		var line []string

		if r.Type == typeName && creationDate.After(date) {
			line = recordToCsv(r)
		} else if all == true && creationDate.After(date) {
			line = recordToCsv(r)
		}

		if line != nil {
			err = writer.Write(line)
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

	// create one file per type
	for t, _ := range types() {
		fmt.Printf("Exporting %s data... ", t)
		err := createCSV(false, t, q.Record, fromDate)

		if err != nil {
			fmt.Println("Error", err)
		} else {
			fmt.Println("done!")
		}
	}

	//create one file with all data
	fmt.Printf("Exporting all data... ")
	errAll := createCSV(true, "all", q.Record, fromDate)

	if errAll != nil {
		fmt.Println("Error", errAll)
	} else {
		fmt.Println("done!")
	}
}
