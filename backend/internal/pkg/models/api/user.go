package api

import "github.com/DuckLuckBreakout/ozonBackend/internal/server/tools/sanitizer"

type ApiLoginUser struct {
	Email    string `json:"email" valid:"email"`
	Password string `json:"password" valid:"stringlength(6|32)"`
}

func (lu *ApiLoginUser) Sanitize() {
	sanitizer := sanitizer.NewSanitizer()
	lu.Email = sanitizer.Sanitize(lu.Email)
	lu.Password = sanitizer.Sanitize(lu.Password)
}

type ApiUpdateUser struct {
	FirstName string `json:"first_name" valid:"utfletter, stringlength(1|30)"`
	LastName  string `json:"last_name" valid:"utfletter, stringlength(1|30)"`
}

func (uu *ApiUpdateUser) Sanitize() {
	sanitizer := sanitizer.NewSanitizer()
	uu.FirstName = sanitizer.Sanitize(uu.FirstName)
	uu.LastName = sanitizer.Sanitize(uu.LastName)
}

type ApiAvatar struct {
	Url string `json:"url" valid:"minstringlength(1)"`
}

type ApiSignupUser struct {
	Email    string `json:"email" valid:"email"`
	Password string `json:"password" valid:"stringlength(6|32)"`
}

func (su *ApiSignupUser) Sanitize() {
	sanitizer := sanitizer.NewSanitizer()
	su.Email = sanitizer.Sanitize(su.Email)
	su.Password = sanitizer.Sanitize(su.Password)
}

type ApiProfileUser struct {
	Id        uint64    `json:"-"`
	FirstName string    `json:"first_name" valid:"utfletter, stringlength(1|30)"`
	LastName  string    `json:"last_name" valid:"utfletter, stringlength(1|30)"`
	Avatar    ApiAvatar `json:"avatar" valid:"notnull, json"`
	AuthId    uint64    `json:"-"`
	Email     string    `json:"email" valid:"email"`
}
