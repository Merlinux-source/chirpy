/*
 * Copyright 2025 Samuel Kemper
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
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"main/internal/database"
	"net/http"
	"time"
)

func handlerCreateChirp(w http.ResponseWriter, r *http.Request, conf *apiConfig) {
	var err error
	type reqParam struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}
	type result struct {
		Id         uuid.UUID `json:"id"`
		CreatedAt  time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Body       string    `json:"body"`
		UserID     uuid.UUID `json:"user_id"`
	}

	var reqVal reqParam
	var cleanChirp string
	var valid bool
	reqDecoder := json.NewDecoder(r.Body)
	err = reqDecoder.Decode(&reqVal)
	if err != nil {
		fmt.Println("Error occured handlerCreateChirp -> reqVal Decoding", err)
	}
	valid = validateChirp(reqVal.Body)
	if !valid {
		w.WriteHeader(400)
		w.Write([]byte("{\"error\": \"Chirp is too long\"}"))
		return
	}

	cleanChirp = sanitizeChirp(reqVal.Body)
	chirp, err := conf.query.CreateChirp(r.Context(), database.CreateChirpParams{Body: cleanChirp, UserID: reqVal.UserId})
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("some error occured."))
		return
	}
	ret, err := json.Marshal(result{Id: chirp.ID, CreatedAt: chirp.CreatedAt, Updated_at: chirp.UpdatedAt, Body: chirp.Body, UserID: chirp.UserID})
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("some error occured."))
		return
	}
	w.WriteHeader(201)
	w.Write(ret)
	return
}
