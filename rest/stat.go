package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/lyddonb/trajectory/api"
)

type StatServices struct {
	api *api.StatAPI
}

func NewStatServices(statAPI *api.StatAPI) *StatServices {
	return &StatServices{statAPI}
}

func (s *StatServices) addStat(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var statJson map[string]*json.RawMessage

	err := decoder.Decode(&statJson)

	if err != nil {
		SendJsonErrorResponse(w, &statJson, err.Error())
		return
	}

	request, stat := s.api.MakeRequestStats(statJson)

	//_, ok := task[db.REQUEST_ADDRESS]

	//if !ok {
	//task[db.REQUEST_ADDRESS] = r.Host
	//}

	timestamp, e := s.api.SaveRequestStats(request, stat)

	if e != nil {
		SendJsonErrorResponse(w, &statJson, e.Error())
		fmt.Println("Stat failed with %s at %s", e, timestamp)
		return
	}

	SendJsonResponse(w, timestamp, nil)
}

func (s *StatServices) getAllStats(w http.ResponseWriter, r *http.Request) {
	SendJsonResponse(w, "Test", nil)
}

func (s *StatServices) getStatByRequestId(w http.ResponseWriter, r *http.Request) {
	requestId := r.URL.Query().Get(":requestId")

	if requestId == "" {
		SendJsonErrorResponse(w, nil, "No request id passed in.")
	}

	stat, err := s.api.GetStatForRequest(requestId)

	SendJsonResponse(w, stat, err)
}
