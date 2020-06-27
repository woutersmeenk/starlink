package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
	"time"
)

var (
	uriBase              = "https://www.space-track.org"
	requestLogin         = "/ajaxauth/login"
	requestCmdAction     = "/basicspacedata/query"
	requestFindStarlinks = "/class/satcat/NORAD_CAT_ID/>40000/SATNAME/STARLINK~~/format/json/orderby/NORAD_CAT_ID%20asc"
	requestOMMStarlink1  = "/class/omm/NORAD_CAT_ID/"
	requestOMMStarlink2  = "/orderby/EPOCH%20asc/format/json"
)

type satcatRecord struct {
	SatName    string `json:"SATNAME"`
	NoradCatID string `json:"NORAD_CAT_ID"`
	Launch     string `json:"LAUNCH"`
}

type ommRecord struct {
	Epoch       string `json:"EPOCH"`
	RaOfAscNode string `json:"RA_OF_ASC_NODE"`
	MeanAnomaly string `json:"MEAN_ANOMALY"`
}

func main() {
	loginData := getLoginData()

	cookieJar, _ := cookiejar.New(nil)

	client := &http.Client{
		Jar: cookieJar,
	}

	res, err := client.Post(uriBase+requestLogin, "application/json", strings.NewReader(loginData))
	check(err)
	checkStatus(res)

	res, err = client.Get(uriBase + requestCmdAction + requestFindStarlinks)
	check(err)
	checkStatus(res)

	decoder := json.NewDecoder(res.Body)
	var satcatRecords []satcatRecord
	err = decoder.Decode(&satcatRecords)
	check(err)

	fmt.Printf("# Sats: %d\n", len(satcatRecords))

	requests := 1
	for _, satRecord := range satcatRecords {
		fmt.Printf("%s %s Launched: %s\n", satRecord.NoradCatID, satRecord.SatName, satRecord.Launch)

		res, err = client.Get(uriBase + requestCmdAction + requestOMMStarlink1 + satRecord.NoradCatID + requestOMMStarlink2)
		check(err)
		checkStatus(res)
		requests++

		decoder := json.NewDecoder(res.Body)
		var ommRecords []ommRecord
		err = decoder.Decode(&ommRecords)
		check(err)
		fmt.Printf("# Epochs: %d\n", len(ommRecords))

		if requests > 28 {
			fmt.Printf("Sleep for %s\n", time.Minute)
			time.Sleep(time.Minute)
			requests = 0
		}
	}
}

func getLoginData() string {
	data, err := ioutil.ReadFile("logindata.json")
	check(err)
	return string(data)
}

func readerToString(r io.Reader) string {
	buf := new(strings.Builder)
	_, err := io.Copy(buf, r)
	check(err)
	return buf.String()
}

func checkStatus(res *http.Response) {
	if res.StatusCode != 200 {
		fmt.Printf("Invalid request. Status: %s Body: %s", res.Status, readerToString(res.Body))
		os.Exit(-1)
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
