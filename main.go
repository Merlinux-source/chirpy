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
	"embed"
	"fmt"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"main/internal/database"
	"net/http"
	"os"
	"regexp"
	"sync/atomic"
	"time"
)

type UserReturnObject struct {
	Id          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

type apiConfig struct {
	fileserverHits atomic.Int32
	query          *database.Queries
	jwt_secret     string
	polka_key      string
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

//go:embed sql/schema/*.sql
var embedMigrations embed.FS //

func main() {
	var serverMux = http.NewServeMux()
	var config = apiConfig{}
	_ = godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	jwtSecret := os.Getenv("JWT_SECRET")
	polkaAPIKey := os.Getenv("POKA_KEY")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("FATAL ERROR, cannot connect to the database.")
		os.Exit(2)
	}

	_ = goose.SetDialect("postgres")
	goose.SetBaseFS(embedMigrations)
	if err := goose.Up(db, "sql/schema"); err != nil { //
		fmt.Println("automatic database migrations failed. This may be a problem with the database. If you want automatic migrations, please ensure that the DB_URL provided user owns the specified database.")
		fmt.Println(err)
	}
	if err := goose.Version(db, "migrations"); err != nil {
		fmt.Println(err)
	}
	dbQueries := database.New(db)
	config.query = dbQueries
	config.jwt_secret = jwtSecret
	config.polka_key = polkaAPIKey

	serverMux.Handle("GET /app/", http.StripPrefix("/app", config.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
	serverMux.HandleFunc("GET /api/healthz", handlerHealth)
	serverMux.HandleFunc("POST /api/chirps", config.middlewareAddCFGContext(handlerCreateChirp))
	serverMux.HandleFunc("GET /api/chirps/{chirpID}", config.middlewareAddCFGContext(handlerGetChirp))
	serverMux.HandleFunc("GET /api/chirps", config.middlewareAddCFGContext(handlerListChirps))
	serverMux.HandleFunc("DELETE /api/chirps/{chirpID}", config.middlewareAddCFGContext(handlerDeleteChirp))
	serverMux.HandleFunc("POST /admin/reset", config.middlewareAddCFGContext(handlerFSHitsReset))
	serverMux.HandleFunc("GET /admin/metrics", config.middlewareAddCFGContext(handlerFSHits))
	serverMux.HandleFunc("POST /api/users", config.middlewareAddCFGContext(handlerCreateUser))
	serverMux.HandleFunc("PUT /api/users", config.middlewareAddCFGContext(handlerChangeUser))
	serverMux.HandleFunc("POST /api/login", config.middlewareAddCFGContext(handlerLoginUser))
	serverMux.HandleFunc("POST /api/refresh", config.middlewareAddCFGContext(handlerRefreshTokens))
	serverMux.HandleFunc("POST /api/revoke", config.middlewareAddCFGContext(handlerRevokeToken))
	serverMux.HandleFunc("POST /api/polka/webhooks", config.middlewareAddCFGContext(handlerPolkaWebhook))

	var httpServer = http.Server{Addr: ":8080", Handler: serverMux}
	err = httpServer.ListenAndServe()
	if err != nil {
		fmt.Println("A fatal error occured while starting the server.")
		fmt.Println(err)
		os.Exit(2)
		return
	}
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
