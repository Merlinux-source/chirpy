/*
Copyright 2025 Merlinux-Source

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"main/internal/database"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	query          *database.Queries
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
func (cfg *apiConfig) middlewareAddCFGContext(next func(http.ResponseWriter, *http.Request, *apiConfig)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r, cfg)
	}
}

func main() {
	var serverMux = http.NewServeMux()
	var config = apiConfig{}
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("FATAL ERROR, cannot connect to the database.")
		os.Exit(2)
	}

	dbQueries := database.New(db)
	config.query = dbQueries

	serverMux.Handle("GET /app/", http.StripPrefix("/app", config.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
	serverMux.HandleFunc("GET /api/healthz", handlerHealth)
	serverMux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	serverMux.HandleFunc("POST /admin/reset", config.middlewareAddCFGContext(handlerFSHitsReset))
	serverMux.HandleFunc("GET /admin/metrics", config.middlewareAddCFGContext(handlerFSHits))
	serverMux.HandleFunc("POST /api/users", config.middlewareAddCFGContext(handlerCreateUser))

	var httpServer = http.Server{Addr: ":8080", Handler: serverMux}
	httpServer.ListenAndServe()
}

func handlerValidateChirp(rw http.ResponseWriter, req *http.Request) {
	type requestParam struct {
		Body string `json:"body"`
	}
	type errReturnParam struct {
		Error string `json:"error"`
	}
	type succReturnParam struct {
		Valid bool `json:"valid"`
	}
	var reqDecoder = json.NewDecoder(req.Body)
	var reqVal requestParam
	var err = reqDecoder.Decode(&reqVal)
	if err != nil {
		fmt.Println("error occured decoind a API request (Validate Chrip)")
		fmt.Println(err)
		rw.WriteHeader(500)
		var marshaledError, err = json.Marshal(errReturnParam{Error: "string"})
		if err != nil {
			fmt.Println("SOMETHING WENT SERIOUSLY WRONG!", err)
			rw.Write([]byte("SOMETHING WENT SERIOUSLY WRONG!"))
			return

		}
		rw.Write(marshaledError)
	}
	if len(reqVal.Body) > 139 {
		rw.WriteHeader(400)
		var marshaledError, err = json.Marshal(errReturnParam{Error: "Chirp is too long"})
		if err != nil {
			fmt.Println("SOMETHING WENT SERIOUSLY WRONG!", err)
		}
		rw.Write(marshaledError)
		return

	}
	res, err := json.Marshal(succReturnParam{Valid: true})
	if err != nil {
		fmt.Println("SOMETHING WENT SERIOUSLY WRONG!", err)
		rw.WriteHeader(500)
		rw.Write([]byte("error"))
		return
	}

	rw.WriteHeader(200)
	rw.Write(res)
}
