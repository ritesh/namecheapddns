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

func main() {
	savedIP := "$HOME/.savedIPAddress"
	f, err := os.Open(savedIP)
	if err != nil {
		log.Printf("could not get a savedIP file %v", err)
	}
	defer f.Close()
	savedIPfromFile, err := ioutil.ReadAll(f)
	if err != nil {
		log.Printf("could not read ip address!")
	}
	savedIP = string(savedIPfromFile)

	currentIP, err := http.Get("http://ifconfig.co")
	if err != nil {
		log.Printf("could not get current external IP from ifconfig! %v", err)
		return
	}
	a, err := ioutil.ReadAll(currentIP.Body)
	if err != nil {
		log.Printf("could not read HTTP response body! %v", err)
		currentIP.Body.Close()
		return
	}
	currentIP.Body.Close()
	ipString := string(a)
	if ipString != savedIP {
		//Do the update
	}
	f.Write([]byte(ipString))
	return
}
