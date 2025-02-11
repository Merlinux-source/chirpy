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
	"fmt"
	"github.com/google/uuid"
	"main/internal/auth"
	"net/http"
)

func handlerDeleteChirp(writer http.ResponseWriter, request *http.Request, config *apiConfig) {
	// request input validation
	chirpIDString := request.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil { // obvious invalid request, but still don't let them know that we know that they know.
		writer.WriteHeader(http.StatusNotFound)
		return // this would be the perfect time to flag the requesting IP as malicious.
	}

	// user authorization
	userJWT, err := auth.GetBearerToken(request.Header) // user supplied access token
	if err != nil {
		fmt.Println("FATAL ERROR, cannot get JWT token.")
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	userUUID, err := auth.ValidateJWT(userJWT, config.jwt_secret) // access token is valid
	if err != nil {
		fmt.Println("FATAL ERROR, cannot validate JWT token.")
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	chirp, err := config.query.GetChirpById(request.Context(), chirpID) // chrip actually exists.
	if err != nil {
		fmt.Println("FATAL ERROR, cannot get chirp by id.")
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	if chirp.UserID != userUUID { // access token is valid for his chirp
		fmt.Println("FATAL ERROR, chirp id and user id do not match.")
		writer.WriteHeader(http.StatusForbidden)
		return
	}

	err = config.query.DeleteChirp(request.Context(), chirp.ID)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(204)
	return
}
