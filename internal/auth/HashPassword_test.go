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
	"testing"
)

func TestCheckPasswordHash(t *testing.T) {
	type args struct {
		password string
		hash     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "valid password",
			args:    args{password: "mysecretpassword", hash: hashPassword("mysecretpassword")},
			wantErr: false,
		},
		{
			name:    "invalid password",
			args:    args{password: "wrongpassword", hash: hashPassword("mysecretpassword")},
			wantErr: true,
		},
		{
			name:    "empty password",
			args:    args{password: "", hash: hashPassword("mysecretpassword")},
			wantErr: true,
		},
		{
			name:    "empty hash",
			args:    args{password: "mysecretpassword", hash: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckPasswordHash(tt.args.password, tt.args.hash); (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHashPassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "valid password",
			args:    args{password: "mysecretpassword"},
			wantErr: false,
		},
		{
			name:    "empty password",
			args:    args{password: ""},
			wantErr: false, // Hashing an empty password should not return an error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHash, err := HashPassword(tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.args.password != "" && gotHash == "" {
				t.Errorf("HashPassword() gotHash = %v, want non-empty hash", gotHash)
			}
		})
	}
}

// Helper function to hash a password for testing
func hashPassword(password string) string {
	hash, _ := HashPassword(password)
	return hash
}
