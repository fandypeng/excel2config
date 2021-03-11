package model

import (
	"encoding/json"
	"strconv"
)

//this type can be used in json parse function,
//when you want to convert a string value to int
type FlexInt int

func (fi *FlexInt) UnmarshalJSON(b []byte) error {
	if b[0] != '"' {
		return json.Unmarshal(b, (*int)(fi))
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	if s == "" {
		*fi = FlexInt(0)
		return nil
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*fi = FlexInt(i)
	return nil
}

//this type can be used in json parse function,
//when you want to convert a string value to int
type FlexFloat float64

func (fi *FlexFloat) UnmarshalJSON(b []byte) error {
	if b[0] != '"' {
		return json.Unmarshal(b, (*float64)(fi))
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	i, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	*fi = FlexFloat(i)
	return nil
}
