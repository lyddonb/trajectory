package trajectory

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/bmizerany/pat"
	"github.com/garyburd/redigo/redis"
)

func writeJson(w http.ResponseWriter, content interface{}) {
	m, _ := json.Marshal(content)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(m)
}

func ListMachines(w http.ResponseWriter, req *http.Request) {
	redisClient, err := Connect()
	if err != err {
		fmt.Fprintf(os.Stderr, "Could not connection to redis: %s", err.Error())
		// TODO: Return bad response.
		return
	}

	keys, _ := redis.Values(redisClient.Do("ZREVRANGE", Machine, 0, 150))

	machineNames := make([]string, 0)

	if err := redis.ScanSlice(keys, &machineNames); err != nil {
		fmt.Println("Error: ", err)
	}

	req.ParseForm()
	fmt.Println("Filter %s", req.FormValue("filter"))

	filter := req.FormValue("filter")

	// TODO: Check for best "set" type option
	names := make([]string, 0)

	for _, name := range machineNames {
		if filter != "" {
		}

		names = append(names, strings.Split(name, ".")[0])
	}

	writeJson(w, &names)
}

func ListStats(w http.ResponseWriter, req *http.Request) {
	// TODO: Get the machine/parent info to get the requests by.

	redisClient, err := Connect()
	if err != err {
		fmt.Fprintf(os.Stderr, "Could not connection to redis: %s", err.Error())
		// TODO: Return bad response.
		return
	}

	keys, _ := redis.Values(redisClient.Do("ZREVRANGE", StatKeys, 0, 150))

	statKeys := make([]string, 0)
	if err := redis.ScanSlice(keys, &statKeys); err != nil {
		fmt.Println("Error: ", err)
	}

	writeJson(w, &statKeys)
}

//func Stats(w http.ResponseWriter, req *http.Request) {
//redisClient, err := Connect()

//if err != err {
//fmt.Fprintf(os.Stderr, "Could not connection to redis: %s", err.Error())
//}

//keys, _ := redis.Values(redisClient.Do("KEYS", "Stat:*"))

//var keyStrs []string
//redis.ScanSlice(keys, &keyStrs)

//stats := make([]Stat, 0)

//for _, key := range keyStrs {
//var stat Stat
//fmt.Println(key)
//itemDict, _ := redis.Values(redisClient.Do("HGETALL", key))

//if err := redis.ScanStruct(itemDict, &stat); err != nil {
//fmt.Printf("Failed %s\n", err)
//} else {
//stats = append(stats, stat)
//}
//}

//m, _ := json.Marshal(&stats)

//io.WriteString(w, string(m))
//}

func StartRest() {
	server := pat.New()

	server.Get("/api/stats", http.HandlerFunc(ListStats))
	server.Get("/api/hosts", http.HandlerFunc(ListMachines))

	http.Handle("/api/", server)

	err := http.ListenAndServe(":8124", nil)

	if err != nil {
		log.Println(err)
	}
}
