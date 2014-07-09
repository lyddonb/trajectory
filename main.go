package main

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	auth "github.com/abbot/go-http-auth"
	"github.com/lyddonb/trajectory/db"
	"github.com/lyddonb/trajectory/pipe"
	"github.com/lyddonb/trajectory/rest"
)

const (
	TASK_PREFIX = "/api/tasks/"
	STAT_PREFIX = "/api/stats/"
)

type LoginConfig struct {
	Username string
	Password string
}

func middleware(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}

	// TODO: Get this to only be setup once.
	authenticator := auth.NewBasicAuthenticator(
		"trajectory.com", GetSecret(getLoginConfig()))
	return auth.JustCheck(authenticator, h)
}

func setupTasks(pool db.DBPool, writeToFile bool) {
	router := rest.SetupTaskRouter(pool, TASK_PREFIX, writeToFile, middleware)

	http.Handle(TASK_PREFIX, router)
}

func setupStats(pool db.DBPool, writeToFile bool) {
	router := rest.SetupStatRouter(pool, STAT_PREFIX, writeToFile, middleware)

	http.Handle(STAT_PREFIX, router)
}

func GetSecret(loginConfigs []LoginConfig) auth.SecretProvider {

	return func(user string, realm string) string {
		for _, loginConfig := range loginConfigs {
			if loginConfig.Username == user {
				h := sha1.New()
				io.WriteString(h, loginConfig.Password)
				hashed := base64.StdEncoding.EncodeToString(h.Sum(nil))
				hashed = fmt.Sprintf("{SHA}%s", hashed)
				return hashed
			}
		}

		return ""
	}
}

func getLoginConfig() []LoginConfig {
	file, err := os.Open("login.json")

	loginConfig := []LoginConfig{}

	if err != nil || file == nil {
		fmt.Println("Error loading login config: ", err)
		return loginConfig
	}

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&loginConfig)

	if err != nil {
		fmt.Println("Error decoding login config: ", err)
	}

	return loginConfig
}

func setupWeb() {
	authenticator := auth.NewBasicAuthenticator(
		"trajectory.com", GetSecret(getLoginConfig()))
	http.HandleFunc(
		"/",
		authenticator.Wrap(func(
			res http.ResponseWriter, req *auth.AuthenticatedRequest) {
			http.FileServer(http.Dir("./web/")).ServeHTTP(res, &req.Request)
		}))
}

func main() {
	// Stand up redis pool.
	pool := db.StartDB("127.0.0.1:6379", "")

	go func() {
		listener := pipe.MakeConnection("tcp", ":1300")
		taskPipeline := pipe.NewTaskPipeline(pool)
		pipe.Listen(listener, taskPipeline)
	}()

	writeToFile := false

	setupStats(pool, writeToFile)
	setupTasks(pool, writeToFile)
	setupWeb()

	http.ListenAndServe(":3000", nil)
}
