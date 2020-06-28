package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
)

var (
	uriBase              = "https://www.space-track.org"
	requestLogin         = "/ajaxauth/login"
	requestCmdAction     = "/basicspacedata/query"
	requestFindStarlinks = "/class/satcat/NORAD_CAT_ID/>40000/SATNAME/STARLINK~~/format/json/orderby/NORAD_CAT_ID%20asc"
	requestOMMStarlink1  = "/class/omm/NORAD_CAT_ID/"
	requestOMMStarlink2  = "/orderby/EPOCH%20asc/format/json"
)

const (
	epochFormat = "2006-01-02T15:04:05"
)

func createClient() *http.Client {
	cookieJar, err := cookiejar.New(nil)
	check(err)

	return &http.Client{
		Jar: cookieJar,
	}
}

func login(client *http.Client) {
	loginData := getLoginData()
	res, err := client.Post(uriBase+requestLogin, "application/json", strings.NewReader(loginData))
	check(err)
	checkStatus(res)
}

type satcatRecord struct {
	SatName    string `json:"SATNAME"`
	NoradCatID string `json:"NORAD_CAT_ID"`
	Launch     string `json:"LAUNCH"`
}

func findStarlinkSats(client *http.Client) (satcatRecords []satcatRecord) {
	res, err := client.Get(uriBase + requestCmdAction + requestFindStarlinks)
	check(err)
	checkStatus(res)

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&satcatRecords)
	check(err)

	return satcatRecords
}

type ommRecord struct {
	NoradCatID  string `json:"NORAD_CAT_ID"`
	ObjectName  string `json:"OBJECT_NAME"`
	Epoch       string `json:"EPOCH"`
	RaOfAscNode string `json:"RA_OF_ASC_NODE"`
	MeanAnomaly string `json:"MEAN_ANOMALY"`
}

func getOomRecords(client *http.Client, satcatRecords []satcatRecord) (ommRecords []ommRecord) {
	satIds := ""
	for _, satRecord := range satcatRecords {
		if len(satIds) > 0 {
			satIds += ","
		}
		satIds += satRecord.NoradCatID
		fmt.Printf("%s %s Launched: %s\n", satRecord.NoradCatID, satRecord.SatName, satRecord.Launch)
	}

	res, err := client.Get(uriBase + requestCmdAction + requestOMMStarlink1 + satIds + requestOMMStarlink2)
	check(err)
	checkStatus(res)

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ommRecords)
	check(err)

	return ommRecords
}

func getLoginData() string {
	data, err := ioutil.ReadFile("logindata.json")
	check(err)
	return string(data)
}

func checkStatus(res *http.Response) {
	if res.StatusCode != 200 {
		fmt.Printf("Invalid request. Status: %s Body: %s", res.Status, readerToString(res.Body))
		os.Exit(-1)
	}
}
