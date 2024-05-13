package model

import (
	"encoding/json"
	"io"
)

type User struct {
	Id           string `json:"id,omitempty"`
	Username     string `json:"username,omitempty"`
	ProfileImage string `json:"profileImage,omitempty"`
}

type Users []*User

func (o *Users) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(o)
}

func (o *User) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(o)
}
