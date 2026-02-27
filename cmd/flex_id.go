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
	"fmt"
)

// FlexID handles JSON fields that can be either an int/float or a string,
// which is necessary because the omg.lol API returns numeric IDs for older records
// but alphanumeric hex strings for newer records.
type FlexID string

func (f *FlexID) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch val := v.(type) {
	case string:
		*f = FlexID(val)
	case float64:
		*f = FlexID(fmt.Sprintf("%.0f", val))
	default:
		*f = FlexID(fmt.Sprintf("%v", val))
	}
	return nil
}

func (f FlexID) String() string {
	return string(f)
}
