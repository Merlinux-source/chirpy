/*
 * Copyright 2025 Merlinux-source
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"main/internal/auth"
	"main/internal/database"
	"net/http"
	"os"
	"regexp"
	"sync/atomic"
	"time"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	query          *database.Queries
	jwt_secret     string
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
	_ = godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	jwtSecret := os.Getenv("JWT_SECRET")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("FATAL ERROR, cannot connect to the database.")
		os.Exit(2)
	}

	dbQueries := database.New(db)
	config.query = dbQueries
	config.jwt_secret = jwtSecret

	serverMux.Handle("GET /app/", http.StripPrefix("/app", config.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
	serverMux.HandleFunc("GET /api/healthz", handlerHealth)
	serverMux.HandleFunc("POST /api/chirps", config.middlewareAddCFGContext(handlerCreateChirp))
	serverMux.HandleFunc("POST /admin/reset", config.middlewareAddCFGContext(handlerFSHitsReset))
	serverMux.HandleFunc("GET /admin/metrics", config.middlewareAddCFGContext(handlerFSHits))
	serverMux.HandleFunc("POST /api/users", config.middlewareAddCFGContext(handlerCreateUser))
	serverMux.HandleFunc("POST /api/login", config.middlewareAddCFGContext(handlerLoginUser))
	serverMux.HandleFunc("POST /api/refresh", config.middlewareAddCFGContext(handlerRefreshTokens))

	var httpServer = http.Server{Addr: ":8080", Handler: serverMux}
	err = httpServer.ListenAndServe()
	if err != nil {
		fmt.Println("A fatal error occured while starting the server.")
		fmt.Println(err)
		os.Exit(2)
		return
	}
}

func handlerRefreshTokens(writer http.ResponseWriter, request *http.Request, config *apiConfig) {
	userToken, err := auth.GetBearerToken(request.Header)
	if err != nil {
		// token was not correctly supplied by the client.
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	dbToken, err := config.query.GetToken(request.Context(), userToken)
	if err != nil {
		// token was not found in the database
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	if dbToken.ExpiresAt.After(time.Now()) {
		// token is expired
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	response, err := json.Marshal(struct {
		Token string `json:"token"`
	}{dbToken.Token})
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(response)
	return
}

func validateChirp(chirp string) (valid bool) {
	return len(chirp) < 140
}

func sanitizeChirp(chirp string) (cleanChirp string) {
	var badWords = []string{ // this can be used for a regex list as well.
		"kerfuffle",
		"sharbert",
		"fornax",
	}
	cleanChirp = chirp
	for _, word := range badWords {
		regex := regexp.MustCompile("(i?)" + word)
		cleanChirp = regex.ReplaceAllString(cleanChirp, "****")
	}
	return cleanChirp
}
