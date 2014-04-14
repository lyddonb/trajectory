package trajectory

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"

	"github.com/garyburd/redigo/redis"
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

type Request struct {
	Request string            `json:"request"`
	Items   map[string]string `json:"items"`
}

type Stat struct {
	Parent string
	Id     string `redis:"id"`
	Value  string `redis:"value"`
	Type   string `redis:"type"`
}

func (s *Stat) ToRedis() redis.Args {
	return redis.Args{}.Add(s.Id).AddFlat(&s)
}

func (s *Stat) RedisKey() string {
	return statKeyPrefix + s.Parent
}

func ParseStat(message []byte) []Stat {
	// If starts with { assume json else comma delemited string
	if bytes.HasPrefix(message, []byte{'{'}) {
		return ParseJsonStats(message)
	}

	// Grab the first piece of the string up to the first period and assume some
	// sort of parent id.
	return ParseStringStats(message)
}

func ParseStringStats(message []byte) []Stat {
	stats := make([]Stat, 0)

	messageString := string(message[:])

	for _, mess := range strings.Split(messageString, ",") {
		// TODO: Convert all to regex. For now just hacking ;)
		keyValueSplit := strings.Split(mess, ":")
		parent := strings.Split(keyValueSplit[0], ".")[0]
		valueSplit := strings.Split(keyValueSplit[1], "|")

		stat := Stat{
			parent,
			keyValueSplit[0],
			valueSplit[0],
			valueSplit[1],
		}

		stats = append(stats, stat)
	}

	return stats
}

func ParseJsonStats(message []byte) []Stat {
	stats := make([]Stat, 0)

	request := new(Request)
	if err := json.Unmarshal(message, request); err != nil {
		log.Println(err)
		return stats
	}

	// Split the parent out of the request id.
	parent := strings.Split(request.Request, ".")[0]

	for key, value := range request.Items {
		valueSplit := strings.Split(value, "|")
		stat := Stat{
			parent,
			request.Request + "." + key,
			valueSplit[0],
			valueSplit[1],
		}
		stats = append(stats, stat)
	}

	return stats
}

// Make redis hash off the request id.
