package api

import (
	"encoding/json"

	"github.com/lyddonb/trajectory/db"
)

type StatAPI struct {
	dal db.StatDAL
}

func NewStatAPI(pool db.DBPool) *StatAPI {
	return &StatAPI{db.NewStatDataAccess(pool)}
}

func (s *StatAPI) SaveRequestStats(request *db.Request, stat db.Stat) (string, error) {
	return s.dal.SaveRequestStats(request, stat)
}

func (s *StatAPI) GetStatForRequest(requestId string) (db.Stat, error) {
	return s.dal.GetStatForRequest(requestId)
}

func (s *StatAPI) MakeRequestStats(statJson map[string]*json.RawMessage) (*db.Request, db.Stat) {
	var url string
	var machine_info string
	var request_id string
	var statList []map[string]string

	stat := make(db.Stat)

	for key, value := range statJson {
		if key == db.URL {
			json.Unmarshal(*value, &url)
		} else if key == db.ID {
			json.Unmarshal(*value, &request_id)
		} else if key == db.MACHINE_INFO {
			json.Unmarshal(*value, &machine_info)
		} else if key == db.STATS {
			json.Unmarshal(*value, &statList)

			// TODO: Do some aggregation based off the types. Might need to
			// combine with existing, etc

			for _, statMap := range statList {
				for k, v := range statMap {
					stat[k] = v
				}
			}
		}
	}

	return &db.Request{request_id, machine_info, url}, stat
}
