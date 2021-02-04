package listener

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/primalcs/logger/config"
	"github.com/primalcs/logger/types"
)

// NewListener creates an instance for getting http requests for modifying configs
func NewListener(port int, config *config.Config) {
	r := mux.NewRouter()
	r.HandleFunc("/logger/setvar", setVarHandler(config)).Methods(http.MethodGet)
	r.HandleFunc("/logger/health", healthHandler).Methods(http.MethodGet)
	r.Use(mux.CORSMethodMiddleware(r))

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", "0.0.0.0", port),
		WriteTimeout: time.Second * 10,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	log.Println(srv.ListenAndServe())
}

func setVarHandler(config *config.Config) func(http.ResponseWriter, *http.Request) {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		lvl := request.FormValue("level")
		logLevel, ok := types.LogLevelsInv[lvl]
		if !ok {
			return
		}
		config.SetLogLevel(logLevel)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	msg := `{"status":"alive"}`
	w.Write([]byte(msg))
}
