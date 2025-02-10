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
	"crypto/rand"
	"encoding/hex"
	"errors"
)

func MakeRefreshToken() (string, error) {
	var bytes = make([]byte, 32)
	numOfReads, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	if numOfReads != len(bytes) {
		return "", errors.New("read Error, check your systems entropy")
	}
	return hex.EncodeToString(bytes), err
}
