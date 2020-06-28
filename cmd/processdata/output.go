package main

import (
	"strconv"
	"time"
)

type satLaunch map[string]string

func groupSatsPerLaunch(satcatRecords []satcatRecord) (satLaunch satLaunch) {
	satLaunch = make(map[string]string)
	for _, satcatRecord := range satcatRecords {
		satLaunch[satcatRecord.SatName] = satcatRecord.Launch
	}
	return satLaunch
}

type satLocation struct {
	RaOfAscNode float64
	MeanAnomaly float64
}

type groupedSatLocations map[string]map[string]satLocation

func groupRecords(ommRecords []ommRecord) (groupedSatLocations groupedSatLocations) {
	groupedSatLocations = make(map[string]map[string]satLocation)
	for _, ommRecord := range ommRecords {
		epoch, _ := time.Parse(epochFormat, ommRecord.Epoch)
		epochDate := epoch.Format("2006-01-02")
		meanAnomaly, err := strconv.ParseFloat(ommRecord.MeanAnomaly, 64)
		check(err)
		raOfAscNode, err := strconv.ParseFloat(ommRecord.RaOfAscNode, 64)
		check(err)
		satName := ommRecord.ObjectName

		if _, ok := groupedSatLocations[epochDate]; !ok {
			groupedSatLocations[epochDate] = make(map[string]satLocation)
		}
		if _, ok := groupedSatLocations[epochDate][satName]; !ok {
			groupedSatLocations[epochDate][satName] = satLocation{
				MeanAnomaly: meanAnomaly,
				RaOfAscNode: raOfAscNode,
			}
		}
	}
	return groupedSatLocations
}
