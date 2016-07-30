package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type NameCheapXMLResponse struct {
	Foo string
	Bar string
}

func handleError(message string, err error) {
	log.Printf("message %v", err)
}
func getSavedIP(savedIP string) (IP string, err error) {
	f, err := os.Open(savedIP)
	if os.IsNotExist(err) {
		f, err = os.Create(savedIP)
		if err != nil {
			handleError("didnt find a savedIP file, could not create file", err)
		}

	} else if err != nil {
		handleError("couldn't open file!", err)
	}
	savedIPfromFile, err := ioutil.ReadAll(f)
	if err != nil {
		handleError("could not read file", err)
	}
	savedIP = string(savedIPfromFile)
	return savedIP, err
}
func saveToFile(saveIPFile, IP string) error {
	f, err := os.Open(saveIPFile)
	if err != nil {
		return err
	}
	defer f.Close()
	f.Write([]byte(IP))
	return nil
}

func main() {
	savedIPFile := "savedIPAddress.dat"
	savedIP, err := getSavedIP(savedIPFile)
	if err != nil {
		handleError("could not get saved IP!", err)
		return
	}
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
	if ip != savedIP {
		//Do the update
	}
	saveToFile(savedIPFile, ip)
	return
}
