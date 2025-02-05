package main

import (
	"fmt"
	"net/http"
)

func handlerFSHits(w http.ResponseWriter, req *http.Request, cfg *apiConfig) {
	var response = fmt.Sprintf("    <p>Chirpy has been visited %d times!</p>", cfg.fileserverHits.Load())
	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte("<html>\n  <body>\n    <h1>Welcome, Chirpy Admin</h1>\n" + response + "  </body>\n</html>"))
}
