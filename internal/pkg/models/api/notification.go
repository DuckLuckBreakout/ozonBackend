package api

import "github.com/DuckLuckBreakout/ozonBackend/internal/server/tools/sanitizer"

type ApiNotificationKeys struct {
	Auth   string `json:"auth"`
	P256dh string `json:"p256dh"`
}

func (nk *ApiNotificationKeys) Sanitize() {
	sanitizer := sanitizer.NewSanitizer()
	nk.Auth = sanitizer.Sanitize(nk.Auth)
	nk.P256dh = sanitizer.Sanitize(nk.P256dh)
}

type ApiUserIdentifier struct {
	Endpoint string `json:"endpoint"`
}

func (ui *ApiUserIdentifier) Sanitize() {
	sanitizer := sanitizer.NewSanitizer()
	ui.Endpoint = sanitizer.Sanitize(ui.Endpoint)
}

type ApiNotificationCredentials struct {
	ApiUserIdentifier
	Keys ApiNotificationKeys `json:"keys"`
}

func (nc *ApiNotificationCredentials) Sanitize() {
	nc.ApiUserIdentifier.Sanitize()
	nc.Keys.Sanitize()
}
