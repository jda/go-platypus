package platypus // import "go.jona.me/platypus"

import (
	"errors"
)

type LoginMethodParameters struct {
	Logintype string `xml:"logintype"`
	Username  string `xml:"username"`
	Password  string `xml:"password"`
	Datatype  string `xml:"datatype"`
}

// Login checks if staff login credentials are valid
func (p Platypus) Login(username string, password string) error {
	params := LoginMethodParameters{
		Logintype: "Staff",
		Datatype:  "XML",
		Username:  username,
		Password:  password,
	}

	res, err := p.Exec("Login", params, nil)
	if err != nil {
		return err
	}

	if res.Success == 0 {
		return errors.New(res.ResponseText)
	}

	return nil
}
