package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/DuckLuckBreakout/web/backend/internal/pkg/models"
	"github.com/DuckLuckBreakout/web/backend/internal/pkg/notification"
	"github.com/DuckLuckBreakout/web/backend/internal/server/errors"
	"github.com/DuckLuckBreakout/web/backend/internal/server/tools/http_utils"
	"github.com/DuckLuckBreakout/web/backend/internal/server/tools/validator"
	"github.com/DuckLuckBreakout/web/backend/pkg/tools/logger"
	"github.com/DuckLuckBreakout/web/backend/pkg/tools/server_push"
)

type NotificationHandler struct {
	NotificationUCase notification.UseCase
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

	var credentials models.NotificationCredentials
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

	err = h.NotificationUCase.SubscribeUser(currentSession.UserData.Id, &credentials)
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

	var userIdentifier models.UserIdentifier
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

	err = h.NotificationUCase.UnsubscribeUser(currentSession.UserData.Id, userIdentifier.Endpoint)
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
