package trajectory

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

/*
{
    "request": "requestid.path.appid.version.module.instanceid",
    "items": {
        "cpu_usage": "0.0|g",
        "memory_usage": "0.0|g",
        "exec_time": "0|ms",
        "overhead": "0|g",
        "system_GetSystemStats_offset": "6|g",
        "end_memory": "0.0|g",
        "end_cpu": "0.0|g"
    }
}

"requestid.path.appid.version.module.instanceid.cpu_usage:0.0|g"
"requestid.path.appid.version.module.instanceid.exec_time:0|ms"

"requestid.path.appid.version.module.instanceid.exec_time:0|ms"
*/

var statKeyPrefix = "Stat:"
var timestampKey = "timestamp"
var StatKeys = "StatKeys"

type Request struct {
	Id      string            `json:"id"`
	Machine string            `json:"machine"`
	Stats   map[string]string `json:"stats"`
}

type Stat struct {
	Id      string `redis:"id"`
	Machine string
	Value   string `redis:"value"`
	Type    string `redis:"type"`
}

func (s *Stat) Key() string {
	return statKeyPrefix + s.Id
}

func (s *Stat) KeyWithOutStatType() string {
	splitId := strings.Split(s.Id, ".")
	return s.Machine + "$" + strings.Join(splitId[:len(splitId)-1], ".")
}

func ParseStat(message []byte) []Stat {
	// If starts with { assume json else comma delemited string
	if bytes.HasPrefix(message, []byte{'{'}) {
		return ParseJsonStats(message)
	}

	// Grab the first piece of the string up to the first period.
	fmt.Println("Parse string only stat %s", message)
	//return ParseStringStats(message)
	return nil
}

//func ParseStringStats(message []byte) []Stat {
//stats := make([]Stat, 0)

//messageString := string(message[:])
////timestamp := time.Now().UTC().Format("20060102150405")

//for _, mess := range strings.Split(messageString, ",") {
//// TODO: Convert all to regex. For now just hacking ;)
//keyValueSplit := strings.Split(mess, ":")
//valueSplit := strings.Split(keyValueSplit[1], "|")

//stats = append(stats,
//makeStat(keyValueSplit[0], valueSplit[0], valueSplit[1]))
//}

//return stats
//}

func makeStat(id, machine, value, statType string) Stat {
	return Stat{
		id,
		machine,
		value,
		statType,
	}
}

func ParseJsonStats(message []byte) []Stat {
	stats := make([]Stat, 0)

	request := new(Request)
	if err := json.Unmarshal(message, request); err != nil {
		log.Println(err)
		return stats
	}

	for key, value := range request.Stats {
		valueSplit := strings.Split(value, "|")
		stats = append(stats,
			makeStat(request.Id+"."+key, request.Machine, valueSplit[0],
				valueSplit[1]))
	}

	return stats
}

// Make redis hash off the request id.
