package handler

import (
	"encoding/json"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/tools/sanitizer"
	"io/ioutil"
	"net/http"

	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/notification"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/errors"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/tools/http_utils"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/tools/validator"
	"github.com/DuckLuckBreakout/ozonBackend/pkg/tools/logger"
	"github.com/DuckLuckBreakout/ozonBackend/pkg/tools/server_push"
)

type NotificationHandler struct {
	NotificationUCase notification.UseCase
}

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

func NewHandler(notificationUCase notification.UseCase) notification.Handler {
	return &NotificationHandler{
		NotificationUCase: notificationUCase,
	}
}

func (h *NotificationHandler) SubscribeUser(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("notification_handler", "SubscribeUser", requireId, err)
		}
	}()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var credentials ApiNotificationCredentials
	err = json.Unmarshal(body, &credentials)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrCanNotUnmarshal, http.StatusBadRequest)
		return
	}
	credentials.Sanitize()

	err = validator.ValidateStruct(credentials)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusBadRequest)
		return
	}

	currentSession := http_utils.MustGetSessionFromContext(r.Context())

	err = h.NotificationUCase.SubscribeUser(
		&models.UserId{Id: currentSession.UserData.Id},
		&models.NotificationCredentials{
			UserIdentifier: models.UserIdentifier{
				Endpoint: credentials.Endpoint,
			},
			Keys:           models.NotificationKeys{
				Auth:   credentials.Keys.Auth,
				P256dh: credentials.Keys.P256dh,
			},
		})
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrCanNotAddReview, http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponseSuccess(w, http.StatusOK)
}

func (h *NotificationHandler) UnsubscribeUser(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("notification_handler", "UnsubscribeUser", requireId, err)
		}
	}()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var userIdentifier ApiUserIdentifier
	err = json.Unmarshal(body, &userIdentifier)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrCanNotUnmarshal, http.StatusBadRequest)
		return
	}
	userIdentifier.Sanitize()

	err = validator.ValidateStruct(userIdentifier)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusBadRequest)
		return
	}

	currentSession := http_utils.MustGetSessionFromContext(r.Context())

	err = h.NotificationUCase.UnsubscribeUser(&models.UserId{Id: currentSession.UserData.Id}, userIdentifier.Endpoint)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrCanNotAddReview, http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponseSuccess(w, http.StatusOK)
}

func (h *NotificationHandler) GetNotificationPublicKey(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("notification_handler", "GetNotificationPublicKey", requireId, err)
		}
	}()

	http_utils.SetJSONResponse(w, models.NotificationPublicKey{Key: server_push.VAPIDPublicKey}, http.StatusOK)
}
