package main

import (
	"encoding/csv"
	"fmt"
	force "github.com/ForceCLI/force/lib"
	batch "github.com/octoberswimmer/batchforce"
	"log"
	"os"
	"strings"
	"time"
)

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file " + filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for " + filePath, err)
	}

	return records
}

func main() {
	records := readCsvFile("/Users/abhishekgbhat/GolandProjects/batchforce/examples/SheetData.csv")
	var processRecords string
	var a int = 0
	for value := range records {
		a = a + 1
		processRecords = processRecords + fmt.Sprintf(
			"'%v',",
			records[value],
		)
		if (a >= 630) {
			a = 0
			processRecords = strings.ReplaceAll(processRecords, "[", "")
			processRecords = strings.ReplaceAll(processRecords, "]", "")
			processRecords = processRecords[:len(processRecords)-1]
			println(len(processRecords))
			println(processRecords)
			process(processRecords)
			time.Sleep(1 * time.Second)
			processRecords = ""
		}


	}
}


func process(recordsProcess string) {

	query := `select 
				Id, Name, Inactive_Reason__c, Active__c, Inactive_Reason_Detail__c,
				Outlet__r.Outlet_Midtrans_MID__c 
			  from 
				Account 
			  where 
				Product_Outlet_Record_Type__c = 'Product Outlet (GO-PAY)'
				and Outlet__r.Outlet_Midtrans_MID__c in (` + recordsProcess  + `)`
	batch.Run("Account", query, updateRecords)
}

func updateRecords(record force.ForceRecord) (updates []force.ForceRecord) {
	update := force.ForceRecord{}
	update["Id"] = record["Id"].(string)
	update["Active__c"] = false
	update["Inactive_Reason__c"] = "Others"
	update["Inactive_Reason_Detail__c"] = "Inactive (0 tx in 6 months) P1"
	updates = append(updates, update);
	return
}
