// Copyright 2022 Ralf Geschke. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package timefix

import (
	"strings"
	"time"
)

// InterterTime defines a time.Time value
type InverterTime struct {
	time.Time
}

// UnmarshalJSON implements the json.Marshaler interface
// At first it tries to use the RFC3339 format, but with the values returned from Kostal Inverter API it goes wrong,
// because the trailing "Z" is missing. This differs from the generated API documentation. So the fix is to add the "Z"
// character.
func (m *InverterTime) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" || string(data) == `""` {
		return nil
	}
	// Fractional seconds are handled implicitly by Parse.
	tt, err := time.Parse(`"`+time.RFC3339+`"`, string(data))
	if err != nil {
		// Append "Z" to data string
		s := strings.TrimRight(string(data), "\"") + "Z\""
		tt2, err := time.Parse(`"`+time.RFC3339+`"`, s)
		*m = InverterTime{tt2}
		return err
	} else {
		*m = InverterTime{tt}
		return err
	}
}
