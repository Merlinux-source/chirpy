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
	"errors"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type chirpsResponse struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func handlerListChirps(writer http.ResponseWriter, request *http.Request, config *apiConfig) {
	authorIDString := request.URL.Query().Get("author_id")
	authorID, uuidParseErr := uuid.Parse(authorIDString)

	var chirps, err = config.query.GetChrips(request.Context())
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if uuidParseErr == nil {
		chirps, err = config.query.GetChirpsByUserId(request.Context(), authorID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				writer.WriteHeader(http.StatusNotFound)
				return
			}
			writer.WriteHeader(http.StatusInternalServerError)
		}
	}

	writer.WriteHeader(http.StatusOK)
	var response []chirpsResponse
	for _, chirp := range chirps {
		response = append(response, chirpsResponse{Id: chirp.ID, CreatedAt: chirp.CreatedAt, UpdatedAt: chirp.UpdatedAt, Body: chirp.Body, UserID: chirp.UserID})
	}
	bytes, err := json.Marshal(response)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}
	writer.Header().Set("Content-Type", "application/json")
	_, _ = writer.Write(bytes)
	return
}
