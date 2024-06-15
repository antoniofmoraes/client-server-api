package internal

import (
	"encoding/json"
	"strconv"
)

func (f *Float32) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	float, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return err
	}

	*f = Float32(float)
	return nil
}
