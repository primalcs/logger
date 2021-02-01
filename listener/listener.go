package listener

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/rybnov/logger/config"

	"github.com/rybnov/logger/types"

	"github.com/gorilla/mux"
)

func NewListener(port int, config *config.Config) {
	r := mux.NewRouter()
	r.HandleFunc("/logger/setvar", SetVar(config)).Methods(http.MethodGet)
	r.HandleFunc("/logger/health", Health).Methods(http.MethodGet)
	r.Use(mux.CORSMethodMiddleware(r))

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", "0.0.0.0", port),
		WriteTimeout: time.Second * 10,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go log.Println(srv.ListenAndServe())
}

func SetVar(config *config.Config) func(http.ResponseWriter, *http.Request) {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		lvl := request.FormValue("level")
		logLevel, ok := types.LogLevelsInv[lvl]
		if !ok {
			return
		}
		config.SetLogLevel(logLevel)
	}
}

func Health(w http.ResponseWriter, r *http.Request) {
	msg := `{"status":"alive"}`
	w.Write([]byte(msg))
}
