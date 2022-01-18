package dto

type DtoProfileUser struct {
	Id        uint64    `json:"-"`
	FirstName string    `json:"first_name" valid:"utfletter, stringlength(1|30)"`
	LastName  string    `json:"last_name" valid:"utfletter, stringlength(1|30)"`
	Avatar    DtoAvatar `json:"avatar" valid:"notnull, json"`
	AuthId    uint64    `json:"-"`
	Email     string    `json:"email" valid:"email"`
}

type DtoUpdateUser struct {
	FirstName string `json:"first_name" valid:"utfletter, stringlength(1|30)"`
	LastName  string `json:"last_name" valid:"utfletter, stringlength(1|30)"`
}

type DtoAvatar struct {
	Url string `json:"url" valid:"minstringlength(1)"`
}
