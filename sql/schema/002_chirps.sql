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

-- +goose Up
CREATE TABLE chirps
(
    id         UUID PRIMARY KEY                                      DEFAULT gen_random_uuid(),
    created_at TIMESTAMP                                    NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMP                                    NOT NULL DEFAULT current_timestamp,
    body       TEXT                                         NOT NULL,
    user_id    UUID REFERENCES users (id) ON DELETE CASCADE NOT NULL
);

-- +goose Down
DROP TABLE chirps
