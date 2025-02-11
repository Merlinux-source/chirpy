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
)

func handlerChangeUser(writer http.ResponseWriter, request *http.Request, config *apiConfig) {
	type expectedInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var input expectedInput

	jwtUserToken, err := auth.GetBearerToken(request.Header)
	if err != nil {
		fmt.Println("error in handlerChangeUser -> jwt token is invalid:", err)
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	userUUID, err := auth.ValidateJWT(jwtUserToken, config.jwt_secret)
	if err != nil {
		fmt.Println("error in handlerChangeUser -> jwt token is invalid:", err)
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = json.NewDecoder(request.Body).Decode(&input)
	if err != nil {
		fmt.Println("error in handlerChangeUser -> json decoding failed:", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	validMail, err := mail.ParseAddress(input.Email)
	if err != nil {
		fmt.Println("error in handlerChangeUser -> mail address is invalid:", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	input.Email = validMail.Address

	passwordHash, err := auth.HashPassword(input.Password)
	if err != nil {
		fmt.Println("error in handlerChangeUser -> hashing password failed:", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := config.query.SetUserEmailAndPassword(request.Context(), database.SetUserEmailAndPasswordParams{ID: userUUID, Email: input.Email, HashedPassword: passwordHash})
	if err != nil {
		fmt.Println("error in handlerChangeUser -> query.SetUserEmailAndPassword() failed:", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(UserReturnObject{Id: user.ID, Email: user.Email, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, IsChirpyRed: user.IsChirpyRed})
	if err != nil {
		fmt.Println("error in handlerChangeUser -> json encoding failed:", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write(response)
}
