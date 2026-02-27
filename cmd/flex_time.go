// SPDX-License-Identifier: MPL-2.0
//
// Copyright (c) 2023 Mark Cornick
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cmd

import (
	"encoding/json"
	"time"
)

// FlexTime represents a time.Time that won't fail to unmarshal if the
// API returns something that isn't a string (like null, false, or an empty object)
// due to API bugs with malformed records.
type FlexTime struct {
	time.Time
}

func (f *FlexTime) UnmarshalJSON(b []byte) error {
	// If it's a JSON string, try to parse it as normal time.Time
	if len(b) >= 2 && b[0] == '"' && b[len(b)-1] == '"' {
		return json.Unmarshal(b, &f.Time)
	}

	// For anything else (null, false, empty object {}, numbers), just leave it
	// as the zero time value. We don't want the whole struct to fail parsing
	// because of one bad timestamp on a bugged record.
	f.Time = time.Time{}
	return nil
}

func (f FlexTime) MarshalJSON() ([]byte, error) {
	if f.Time.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(f.Time)
}
