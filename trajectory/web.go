package trajectory

import (
	"log"
	"net/http"

	"github.com/bmizerany/pat"
)

func StartWeb() {
	server := pat.New()

	server.Get("/", http.FileServer(http.Dir("./web")))
	server.Get("/css/", http.FileServer(http.Dir("./web")))
	server.Get("/fonts/", http.FileServer(http.Dir("./web")))
	server.Get("/img/", http.FileServer(http.Dir("./web")))
	server.Get("/js/", http.FileServer(http.Dir("./web")))
	server.Get("/js/vendor", http.FileServer(http.Dir("./web")))

	http.Handle("/", server)

	err := http.ListenAndServe(":8123", nil)

	if err != nil {
		log.Println(err)
	}
}
