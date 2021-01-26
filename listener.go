package logger

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func NewListener(port int, logger *Logger) {
	r := mux.NewRouter()
	r.HandleFunc("/logger/setvar", SetVar(logger)).Methods(http.MethodGet)
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

func SetVar(logger *Logger) func(http.ResponseWriter, *http.Request) {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		lvl := request.FormValue("level")
		logLevel, ok := LogLevelsInv[lvl]
		if !ok {
			return
		}
		logger.config.SetLogLevel(logLevel)
	}
}

func Health(w http.ResponseWriter, r *http.Request) {
	msg := `{"status":"alive"}`
	w.Write([]byte(msg))
}
