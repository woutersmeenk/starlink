package main

import (
	"fmt"
)

func main() {
	client := createClient()
	login(client)

	satcatRecords := findStarlinkSats(client)

	fmt.Printf("# Sats: %d\n", len(satcatRecords))

	ommRecords := getOomRecords(client, satcatRecords[5:7])
	fmt.Printf("# Records: %d\n", len(ommRecords))

	groupedRecords := groupRecords(ommRecords)
	fmt.Printf("# Dates: %d\n", len(groupedRecords))
}
