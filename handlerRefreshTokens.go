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
	"net/http"
	"time"
)

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
		fmt.Println("FATAL ERROR, cannot get bearer token.", err.Error())
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	if dbToken.ExpiresAt.Before(time.Now()) {
		// token is expired
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	if dbToken.RevokedAt.Valid {
		// token was revoked
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	jwtToken, err := auth.MakeJWT(dbToken.UserID, config.jwt_secret, time.Hour*1)
	if err != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	response, err := json.Marshal(struct {
		Token string `json:"token"`
	}{jwtToken})
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Println("FATAL ERROR, cannot serialize the token.", err.Error())
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write(response)
	return
}
