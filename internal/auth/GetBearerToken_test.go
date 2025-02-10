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

package auth

import (
	"github.com/google/uuid"
	"net/http"
	"testing"
	"time"
)

func TestGetBearerToken(t *testing.T) {
	type args struct {
		headers http.Header
	}
	type test struct {
		name    string
		args    args
		want    string
		wantErr bool
	}
	tests := []test{
		{"Test no auth header", args{http.Header{}}, "", true},
		{"Test invalid auth header", args{http.Header{"Authorization": []string{"Haha"}}}, "", true},
		{"Test incomplete auth header", args{http.Header{"Authorization": []string{"Bearer "}}}, "", true},
	}

	// Generate a valid JWT for the complete valid header test case
	validJWT, err := MakeJWT(uuid.New(), "123", time.Second*5) // Adjust parameters as needed
	if err != nil {
		t.Fatalf("Failed to create valid JWT: %v", err)
	}
	tests = append(tests, test{"Test complete valid header", args{http.Header{"Authorization": []string{"Bearer " + validJWT}}}, validJWT, false})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBearerToken(tt.args.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetBearerToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}
