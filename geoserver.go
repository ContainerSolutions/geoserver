package main

import (
  "fmt"
  "net/http"
  "strings"
  "log"
  "os"
  "bytes"
  "time"
  "encoding/json"
)

const zoneUrl string = "http://metadata.google.internal/computeMetadata/v1/instance/network-interfaces/0/forwarded-ips"

func getZone() string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", zoneUrl, nil)
	req.Header.Add("Metadata-Flavor", "Google")
	response, err := client.Do(req)
	if err != nil {
		panic("Unable to retrieve zone information: " + err.Error())
	} else {
		defer response.Body.Close()
		buf := new(bytes.Buffer)
		buf.ReadFrom(response.Body)
		return buf.String()
	}
}

var myIp string

const reportingApiUrl string = "https://z4a24pdwp6.execute-api.eu-west-1.amazonaws.com/prod/AccessGeoLocations"
const ipServiceUrl string = "https://api.ipify.org"
//const geoServiceUrl string = "http://ip-api.com/json/"
const geoServiceUrl string = "http://ipinfo.io/"

func sendLocation(jsonStr string) {
  var jsonByte = []byte(jsonStr)
	req, err := http.NewRequest("POST", reportingApiUrl, bytes.NewBuffer(jsonByte))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
}

func getMyIPAddress() string {
	response, err := http.Get(ipServiceUrl)
	if err != nil {
		panic("Unable to retrieve IP: " + err.Error())
	} else {
		defer response.Body.Close()
		buf := new(bytes.Buffer)
		buf.ReadFrom(response.Body)
		myIp := buf.String()
		return myIp
	}
}

func getCoordinates() map[string]string {
	//response, err := http.Get(geoServiceUrl)
client := &http.Client{
}
	req, err := http.NewRequest("GET", geoServiceUrl, nil)
	// ...
	req.Header.Add("User-Agent", "curl/7.49.1")
	response, err := client.Do(req)
	if err != nil {
		panic("Unable to retrieve coordinates: " + err.Error())
	} else {
		defer response.Body.Close()
		buf := new(bytes.Buffer)
		buf.ReadFrom(response.Body)
		var dat map[string]interface{}
		//fmt.Println("response data:", buf.String())
		if err := json.Unmarshal(buf.Bytes(), &dat); err != nil {
        panic(err)
    }
		res := strings.Split(dat["loc"].(string), ",")
    lat := res[0]
    lon := res[1]
	  return map[string]string{"lat": lat, "lon": lon, "ip": dat["ip"].(string)}
	}
}

func reportLocation() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
  coord := getCoordinates()
	hostname, _ := os.Hostname()
	timestamp := int64(time.Now().Unix())
	mapD := map[string]interface{}{
		"IP": coord["ip"],
		"hostname": hostname,
		"date": timestamp,
		"lat": coord["lat"],
		"lng": coord["lon"],
	}
	mapB, _ := json.Marshal(mapD)
	fmt.Println("JSON:", string(mapB))
	sendLocation(string(mapB))
}

func sayhelloName(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()  // parse arguments, you have to call this by yourself
    fmt.Println(r.Form)  // print form information in server side
    fmt.Println("path", r.URL.Path)
    fmt.Println("scheme", r.URL.Scheme)
    fmt.Println(r.Form["url_long"])
    for k, v := range r.Form {
        fmt.Println("key:", k)
        fmt.Println("val:", strings.Join(v, ""))
    }
    fmt.Fprintf(w, getZone()) // send data to client side
}

func main() {
  go reportLocation()
	ticker := time.NewTicker(10 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <- ticker.C:
				go reportLocation()
			case <- quit:
				ticker.Stop()
				return
			}
		}
	}()
	http.HandleFunc("/", sayhelloName) // set router
	err := http.ListenAndServe(":9090", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
