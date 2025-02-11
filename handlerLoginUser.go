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
	"fmt"
	"github.com/google/uuid"
	"main/internal/auth"
	"main/internal/database"
	"net/http"
	"net/mail"
	"time"
)

func handlerLoginUser(writer http.ResponseWriter, request *http.Request, config *apiConfig) {
	var err error
	type requestUser struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		Id          uuid.UUID `json:"id"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		Email       string    `json:"email"`
		IsChirpyRed bool      `json:"is_chirpy_red"`
		Token       string    `json:"token"`
		Refresh     string    `json:"refresh_token"`
	}
	var reqVal requestUser
	var reqDecoder *json.Decoder
	var expiresIn = 3600

	reqDecoder = json.NewDecoder(request.Body)
	err = reqDecoder.Decode(&reqVal) // validate post data
	if err != nil {
		fmt.Println("error decoding json payload", err.Error())
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	validEmail, err := mail.ParseAddress(reqVal.Email) // validate email input
	if err != nil {
		fmt.Println("error parsing email", err.Error())
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := config.query.GetUserByEmail(request.Context(), validEmail.Address) // validate user password
	if err != nil {
		fmt.Println("error getting user", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			writer.WriteHeader(http.StatusNotFound)
			return
		}
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = auth.CheckPasswordHash(reqVal.Password, user.HashedPassword)
	if err != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	jwt, err := auth.MakeJWT(user.ID, config.jwt_secret, time.Duration(expiresIn)*time.Second)
	if err != nil {
		fmt.Println("error generating jwt token", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	refreshTokenString, err := auth.MakeRefreshToken()
	if err != nil {
		fmt.Println("error generating refresh token", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = config.query.CreateRefreshToken(request.Context(), database.CreateRefreshTokenParams{refreshTokenString, user.ID})
	if err != nil {
		fmt.Println("error creating refresh token", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	// authorization successful
	bytes, err := json.Marshal(response{Id: user.ID, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Email: user.Email, Token: jwt, Refresh: refreshTokenString, IsChirpyRed: user.IsChirpyRed})
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write(bytes)
}
