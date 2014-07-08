package rest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func (s *StatServices) addStatToFile(w http.ResponseWriter, r *http.Request) {
	writeToFile(r, "data/stats.json")
}

func (s *TaskServices) addTaskToFile(w http.ResponseWriter, r *http.Request) {
	writeToFile(r, "data/tasks.json")
}

func writeToFile(r *http.Request, fileName string) {
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	fmt.Println(string(body))

	// open output file
	//fo, err := os.Open(fileName)
	fo, err := os.OpenFile(fileName, os.O_RDWR|os.O_APPEND, 0660)

	if err != nil {
		fo, err = os.Create(fileName)

		if err != nil {
			panic(err)
		}
	}

	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	fo.Write(body)
}
