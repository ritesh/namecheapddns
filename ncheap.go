package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

//<?xml version="1.0"?><interface-response><Command>SETDNSHOST</Command><Language>eng</Language><ErrCount>1</ErrCount><errors><Err1>Domain name not found</Err1></errors><ResponseCount>1</ResponseCount>
//<responses><response><ResponseNumber>316153</ResponseNumber><ResponseString>Validation error; not //found; domain name(s)</ResponseString></response></responses><Done>true</Done><debug><![CDATA[]]></debug></interface-response>

type Response struct {
	ResponseNumber int16  `xml:"response>ResponseNumber"`
	ResponseString string `xml:"response>ResponseString"`
}
type Err struct {
	Err1, Err2, Err3, Err4 string
}
type NameCheapXMLResponse struct {
	Command       string
	Language      string
	ErrCount      int
	Errors        []Err
	ResponseCount int
	Responses     []Response
	Done          bool
	Debug         []byte
}

const saveIPFile = "savedIPAddress.dat"

func handleError(message string, err error) {
	log.Printf("message %v", err)
}
func getSavedIP() (IP string, err error) {
	f, err := os.Open(saveIPFile)
	switch err != nil {
	case os.IsNotExist(err):
		func(f *os.File) {
			f, err := os.Create(saveIPFile)
			if err != nil {
				return
			}
		}(f)
	default:
		handleError("could not handle error!", err)
	}
	savedIPfromFile, err := ioutil.ReadAll(f)
	if err != nil {
		handleError("could not read file", err)
	}
	savedIP := string(savedIPfromFile)
	return savedIP, err
}
func saveToFile(IP string) error {
	f, err := os.Open(saveIPFile)
	if err != nil {
		return err
	}
	defer f.Close()
	f.Write([]byte(IP))
	return nil
}

func main() {
	//## http://dynamicdns.park-your-domain.com/update?host=host_name&
	// domain=domain.com&password=domain_password[&ip=your_ip]
	updateURL := "https://dynamicdns.park-your-domain.com/update?"
	params := map[string]string{
		"host":     "@",
		"domain":   "ritesh.xyz",
		"password": "password",
	}
	savedIP, err := getSavedIP()
	if err != nil {
		handleError("could not get saved IP!", err)
	} else {
		params["ip"] = savedIP
	}
	// Get IP directly instead
	currentIP, err := http.Get("https://ifconfig.co")
	if err != nil {
		handleError("could not get IP address!", err)
		return
	}
	a, err := ioutil.ReadAll(currentIP.Body)
	if err != nil {
		handleError("could not read the current IP", err)
		return
	}
	currentIP.Body.Close()
	ip := string(a)
	if ip == params["ip"] {
		//No change here
		return
	}
	ncheapURL, err := url.Parse(updateURL)
	if err != nil {
		handleError("could not parse URL!", err)
	}
	updateURL = fmt.Sprintf(updateURL+"host=%s&domain=%s&password=%s&ip=%v", params["host"], params["domain"], params["password"], params["ip"])
	log.Printf(updateURL)
	resp, err := http.Get(ncheapURL.String())
	if err != nil {
		handleError("could not update!", err)
		return
	}
	a, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		handleError("could not read response body!", err)
		return
	}
	n := &NameCheapXMLResponse{}
	err = xml.Unmarshal(a, n)
	if err != nil {
		handleError("could not unmarshal response", err)
	}
	resp.Body.Close()
	log.Printf("%v", n)
	log.Printf("%v", n.Responses)
	log.Printf("%v", n.Errors)

	saveToFile(ip)
	return
}
