package main

import (
	"fmt"
	"net/http"
)
func handlerFSHitsReset(w http.ResponseWriter, req *http.Request, cfg *apiConfig) {
	cfg.fileserverHits.Add(-cfg.fileserverHits.Load())
	err := cfg.query.ClearUsers(req.Context())
	if err != nil {
		fmt.Println(err)
	}
}
