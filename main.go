package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	dataload "vube/practice/points/load"
	"vube/practice/points/reporter"
	"vube/practice/points/reset"
	"vube/practice/points/saves"
	"vube/practice/points/startup"
	"vube/practice/points/update"
)

const port int = 8999

/*TODO: signal handle and save before exit? save after each update?*/

func main() {
	log.Printf("RUNNING ON PORT %d\n", port)

	r := mux.NewRouter()
	r.HandleFunc("/update", updateHandler).Methods("POST")
	r.HandleFunc("/info", infoHandler).Methods("GET")
	r.HandleFunc("/reset", resetHandler).Methods("GET")
	r.HandleFunc("/save", saveHandler).Methods("GET")
	http.Handle("/", r)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatalf("LISTEN AND SERVE FAILED. err: %v", err)
	}
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	res, err := reporter.NewReport()
	if err != nil {
		log.Printf("report failed")
		fmt.Fprint(w, `{"result":"failed"}`)
		return
	}
	fmt.Fprint(w, res)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	result, err := update.UnmarshalResult(r)
	if err != nil {
		log.Printf("unmarshalResult failed. err: %s", err)
		fmt.Fprint(w, `{"result":"failed"}`)
		return
	}
	res := update.UpdateAndRecordResults(result)
	fmt.Fprint(w, res)
}

func resetHandler(w http.ResponseWriter, r *http.Request) {
	err := reset.ResetOverallData()
	if err != nil {
		log.Printf("failed to reset overall data. err: %s", err)
		fmt.Fprint(w, `{"result":"failed"}`)
		return
	}
	fmt.Fprint(w, `{"result":"ok"}`)
	return
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	err := saves.SaveAsFile()
	if err != nil {
		fmt.Fprint(w, `{"result":"failed"}`)
		return
	}
	fmt.Fprint(w, `{"result":"ok"}`)
	return
}

func init() {
	load := flag.String("load", "", "load a saved file.")
	flag.Parse()
	if *load != "" {
		err := dataload.LoadFile(*load)
		if err != nil {
			os.Exit(404)
		}
	}
	startup.LoadXterVals()
}
