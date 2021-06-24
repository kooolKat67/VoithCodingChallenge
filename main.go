package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/mux"
)

type Section struct {
	Section string `json:"Section"`
	Content string `json:"Content"`
}

var Sections map[string]string

func init() {
	fmt.Println("This will get called on main initialization")
	Sections = make(map[string]string)

	// Make HTTP request

	response, err := http.Get("https://en.wikipedia.org/wiki/Heidenheim_an_der_Brenz")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}
	document.Find("h2").Each(processElement)

}

// This will get called for each HTML element found
func processElement(index int, element *goquery.Selection) {
	sectText := element.NextUntil("h2")
	sect := element.Text()
	sect = strings.ReplaceAll(sect, "[edit]", "")
	Sections[sect] = sectText.Text()

}

func getSection(w http.ResponseWriter, r *http.Request) {
	// once again, we will need to parse the path parameters
	vars := mux.Vars(r)
	if sectName, ok := vars["sectionName"]; ok {
		if section, ok := Sections[sectName]; ok {
			fmt.Println("Endpoint Hit: section")
			json.NewEncoder(w).Encode(section)
		}
	}
}

func handleRequests() {
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/section/{sectionName}", getSection).Methods("GET")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	handleRequests()
}
