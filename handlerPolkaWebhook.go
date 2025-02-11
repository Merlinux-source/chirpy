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
	"net/http"
)

func handlerPolkaWebhook(writer http.ResponseWriter, request *http.Request, config *apiConfig) {
	type InputData struct {
		UserID uuid.UUID `json:"user_id"`
	}
	type expectedInput struct {
		EventName string    `json:"event"`
		Data      InputData `json:"data"`
	}
	var input expectedInput
	err := json.NewDecoder(request.Body).Decode(&input)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if input.EventName != "user.upgraded" {
		writer.WriteHeader(http.StatusNoContent)
		return
	}
	_, err = config.query.GetUserByID(request.Context(), input.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("error occured handlerPolkaWebhook -> fetching user. User does not exist.")
			writer.WriteHeader(http.StatusNotFound)
			return
		}
		fmt.Println("error occured handlerPolkaWebhook -> fetching user.", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = config.query.SetUserUpgradeStatusTrue(request.Context(), input.Data.UserID)
	if err != nil {
		fmt.Println("error occured handlerPolkaWebhook -> updating user.", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusNoContent)
}
