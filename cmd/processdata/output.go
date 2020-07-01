package main

import (
	"strconv"
	"time"
)

type satLaunchDate map[string]string

func groupSatsPerLaunchDate(satcatRecords []satcatRecord) (satLaunchDate satLaunchDate) {
	satLaunchDate = make(map[string]string)
	for _, satcatRecord := range satcatRecords {
		satLaunchDate[satcatRecord.SatName] = satcatRecord.Launch
	}
	return satLaunchDate
}

// One graph is displayed for each day
type graphData []graphForDate

// The satellite's are grouped per launch. Each group has a different color.
type graphForDate struct {
	Date     string
	Launches map[string]*launchGroup
}

// Launch group contrains the labels, x and y cooridinates.
type launchGroup struct {
	SatNames      []string
	MeanAnomalies []float64
	RaOfAscNodes  []float64
}

func createGraphData(ommRecords []ommRecord, satLaunchDate satLaunchDate) (graphData graphData) {
	var satAdded map[string]struct{}
	var currentGraph graphForDate
	for _, ommRecord := range ommRecords {
		epoch, _ := time.Parse(epochFormat, ommRecord.Epoch)
		epochDate := epoch.Format("2006-01-02")
		meanAnomaly, err := strconv.ParseFloat(ommRecord.MeanAnomaly, 64)
		check(err)
		raOfAscNode, err := strconv.ParseFloat(ommRecord.RaOfAscNode, 64)
		check(err)
		satName := ommRecord.ObjectName

		if currentGraph.Date != epochDate {
			currentGraph = graphForDate{
				Date:     epochDate,
				Launches: make(map[string]*launchGroup),
			}
			graphData = append(graphData, currentGraph)
			satAdded = make(map[string]struct{})
		}
		if _, ok := satAdded[satName]; !ok {
			launchDate := satLaunchDate[satName]
			var currentLaunch *launchGroup
			if foundLaunch, ok := currentGraph.Launches[launchDate]; !ok {
				currentLaunch = &launchGroup{}
				currentGraph.Launches[launchDate] = currentLaunch
			} else {
				currentLaunch = foundLaunch
			}
			currentLaunch.MeanAnomalies = append(currentLaunch.MeanAnomalies, meanAnomaly)
			currentLaunch.RaOfAscNodes = append(currentLaunch.RaOfAscNodes, raOfAscNode)
			currentLaunch.SatNames = append(currentLaunch.SatNames, satName)
			satAdded[satName] = struct{}{}
		}
	}
	return graphData
}
