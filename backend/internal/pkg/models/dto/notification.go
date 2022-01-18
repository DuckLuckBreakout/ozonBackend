package dto

type DtoNotificationKeys struct {
	Auth   string `json:"auth"`
	P256dh string `json:"p256dh"`
}

type DtoSubscribes struct {
	Credentials map[string]*DtoNotificationKeys `json:"subscribes" valid:"notnull"`
}
