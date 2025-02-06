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
	"encoding/json"
	"fmt"
	"main/internal/auth"
	"main/internal/database"
	"net/http"
	"net/mail"
	"time"

	"github.com/google/uuid"
)

func handlerCreateUser(rw http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	var err error
	type requestForm struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	var reqVal requestForm
	rDecoder := json.NewDecoder(r.Body)
	err = rDecoder.Decode(&reqVal)
	if err != nil {
		fmt.Println("error occurred handlerCreateUser -> decode request", err)
		rw.WriteHeader(400)
		_, _ = rw.Write([]byte("{\"error\":\"api usage.\"}"))
		return
	}
	if len(reqVal.Email) < 1 {
		rw.WriteHeader(400)
		fmt.Println("user submitted empty mail")
		return
	}
	if len(reqVal.Password) < 4 {
		rw.WriteHeader(400)
		fmt.Println("user submitted empty password")
		return
	}

	mailAddr, err := mail.ParseAddress(reqVal.Email)
	if err != nil {
		fmt.Println("error occurred handlerCreateUser -> validate email", err)
		rw.WriteHeader(400)
		_, _ = rw.Write([]byte("{\"error\":\"invalid email.\"}"))
		return
	}
	pass, err := auth.HashPassword(reqVal.Password)
	if err != nil {
		rw.WriteHeader(400)
		fmt.Println("error occurred handlerCreateUser -> validate password", err)
		_, _ = rw.Write([]byte("{\"error\":\"invalid password.\"}"))
		return
	}

	user, err := cfg.query.CreateUser(r.Context(), database.CreateUserParams{Email: mailAddr.Address, HashedPassword: pass})
	if err != nil {
		fmt.Println("error occurred handlerCreateUser -> Create user", err)
		rw.WriteHeader(400)
		_, _ = rw.Write([]byte("{\"error\":\"could not create user.\"}"))
		return
	}

	// create user successful.
	rw.WriteHeader(201)
	result, err := json.Marshal(response{Id: user.ID, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Email: user.Email})
	if err != nil {
		fmt.Println(err)
		return
	}

	_, _ = rw.Write(result)
	return
}
