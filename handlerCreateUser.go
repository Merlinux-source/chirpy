package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"time"

	"github.com/google/uuid"
)

func handlerCreateUser(rw http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	var err error
	type requestForm struct {
		Email string `json:"email"`
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
		fmt.Println("error occured handlerCreateUser -> decode request", err)
		rw.WriteHeader(400)
		rw.Write([]byte("{\"error\":\"api usage.\"}"))
		return
	}
	if len(reqVal.Email) < 1 {
		fmt.Println("user submited empty mail")
		return
	}
	mailAddr, err := mail.ParseAddress(reqVal.Email)
	if err != nil {
		fmt.Println("error occured handlerCreateUser -> validate email", err)
		rw.WriteHeader(400)
		rw.Write([]byte("{\"error\":\"invalid email.\"}"))
		return
	}

	user, err := cfg.query.CreateUser(r.Context(), mailAddr.Address)
	if err != nil {
		fmt.Println("error occured handlerCreateUser -> Create user", err)
		rw.WriteHeader(400)
		rw.Write([]byte("{\"error\":\"could not create user.\"}"))
		return
	}
	rw.WriteHeader(201)
	result, err := json.Marshal(response{Id: user.ID, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Email: user.Email})
	if err != nil {
		fmt.Println(err)
		return
	}

	rw.Write(result)
	return
}
