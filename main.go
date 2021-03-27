package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/thedevsaddam/gojsonq/v2"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

//concurrent fetch from endpoint
func getJson(url string, ch chan [] byte, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	ch <- body
}

//Sort responses from multiple clinic providers based on unique responses
func sortJson(byte []byte) []Clinic{

	var results []map[string]interface{}

	err := json.Unmarshal(byte, &results)

	var clinics []Clinic

	if err != nil{
		log.Println("Invalid JSON passed")
		return clinics
	}

	for _, result := range results {
		//Could use reflect.typeof here.
		availability, ok := result["availability"].(map[string]interface{})
		if !ok{
			availability = result["opening"].(map[string]interface{})
		}
		a := Availability{
			From: availability["from"].(string),
			To:   availability["to"].(string),
		}

		if result["name"] != nil {
			clinics = append(clinics, Clinic{
				Name:     result["name"].(string),
				State:    result["stateName"].(string),
				Availability: a,
			})
		}else{
			clinics = append(clinics, Clinic{
				Name:     result["clinicName"].(string),
				State:    result["stateCode"].(string),
				Availability: a,
			})
		}
	}
	return clinics
}

//Fetch requests from multiple clinic providers
func fetchJson() (string, error){
	var slice [] byte

	urls := []string{"https://storage.googleapis.com/scratchpay-code-challenge/dental-clinics.json", "https://storage.googleapis.com/scratchpay-code-challenge/vet-clinics.json"}
	ch := make(chan [] byte, 2)

	var wg sync.WaitGroup
	for _,url := range urls{
		wg.Add(2)
		go getJson(url, ch, &wg)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for range urls{
		slice = append(slice, <-ch...)
	}

	if slice == nil {
		return "", errors.New("no response")
	}

	//removing ][ to validate json ; more like merge the different json arrays
	slice = bytes.ReplaceAll([]byte (slice), []byte("]["),[]byte(","))
	slice = bytes.ReplaceAll([]byte (slice), []byte("] ["),[]byte(","))
	slice = bytes.ReplaceAll([]byte (slice), []byte("]\n["),[]byte(","))
	sorted:= Clinics{sortJson(slice)}

	b, err := json.Marshal(sorted)
	if err != nil {
		return "", errors.New("something went wrong")
	}
	return string(b), nil
}

func search(w http.ResponseWriter, r *http.Request){
	fmt.Println("Incoming Request")
	var clinicName, state, from, to string

	//If method is GET(Params passed as GET Queries)
	if r.Method == "GET" {
		clinicName = r.URL.Query().Get("name")
		state = r.URL.Query().Get("state")
		from = r.URL.Query().Get("from")
		to = r.URL.Query().Get("to")
	}

	//If method is POST(Json request passed)
	if r.Method == "POST" {
		var request Clinic
		err := json.NewDecoder(r.Body).Decode(&request)

		if err != nil {
			response := map[string]string{"status": "Invalid request inputs passed"}
			_ = json.NewEncoder(w).Encode(response)
			return
		}
		clinicName = request.Name
		state = request.State
		from = request.Availability.From
		to = request.Availability.To
	}

	b, err := fetchJson()

	//if error occurred fetching json
	if err != nil{
		response := map[string]string{"status": "Something went wrong fetching request. Kindly try again or check network if error persists"}
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	query := gojsonq.New().JSONString(b).From("clinics").Where("name", ifEmptyOperator(clinicName), clinicName).Where("state",ifEmptyOperator(state),state).Where("availability.from",ifEmptyOperator(from),from).Where("availability.to",ifEmptyOperator(to),to)
	d:= query.Get()

	var c []Clinic
	mapstructure.Decode(d, &c)
	response, _ := json.Marshal(c)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	return
}

//Function to set operator if empty to != else =
func ifEmptyOperator(query string) string {
	if query == ""{
		return "!="
	}
	return "="
}

func main(){
	http.HandleFunc("/", search)
	fmt.Println("API kicking off now  on port 9000")
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
