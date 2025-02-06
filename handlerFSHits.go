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
	"net/http"
)

func handlerFSHits(w http.ResponseWriter, req *http.Request, cfg *apiConfig) {
	var response = fmt.Sprintf("    <p>Chirpy has been visited %d times!</p>", cfg.fileserverHits.Load())
	w.Header().Add("Content-Type", "text/html")
	_, _ = w.Write([]byte("<html>\n  <body>\n    <h1>Welcome, Chirpy Admin</h1>\n" + response + "  </body>\n</html>"))
}
