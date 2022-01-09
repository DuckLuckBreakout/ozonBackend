package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/session"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/user"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/errors"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/tools/http_utils"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/tools/validator"
	"github.com/DuckLuckBreakout/ozonBackend/pkg/tools/logger"
)

type UserHandler struct {
	UserUCase    user.UseCase
	SessionUCase session.UseCase
}

func NewHandler(userUCase user.UseCase, sessionUCase session.UseCase) user.Handler {
	return &UserHandler{
		UserUCase:    userUCase,
		SessionUCase: sessionUCase,
	}
}

// Handle login user
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("user_handler", "Login", requireId, err)
		}
	}()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var authUser models.LoginUser
	err = json.Unmarshal(body, &authUser)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrCanNotUnmarshal, http.StatusBadRequest)
		return
	}
	authUser.Sanitize()

	err = validator.ValidateStruct(authUser)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusBadRequest)
		return
	}

	userId, err := h.UserUCase.Authorize(&authUser)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusUnauthorized)
		return
	}

	currentSession, err := h.SessionUCase.CreateNewSession(userId)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetCookie(w, models.SessionCookieName, currentSession.Value, models.ExpireSessionCookie*time.Second)
	http_utils.SetJSONResponseSuccess(w, http.StatusOK)
}

// Handle update user profile
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("user_handler", "UpdateProfile", requireId, err)
		}
	}()

	currentSession := http_utils.MustGetSessionFromContext(r.Context())

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var updateUser models.UpdateUser
	err = json.Unmarshal(body, &updateUser)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrCanNotUnmarshal, http.StatusBadRequest)
		return
	}
	updateUser.Sanitize()

	err = validator.ValidateStruct(updateUser)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusBadRequest)
		return
	}

	err = h.UserUCase.UpdateProfile(currentSession.UserData.Id, &updateUser)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponseSuccess(w, http.StatusOK)
}

// Handle update avatar in user profile
func (h *UserHandler) UpdateProfileAvatar(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("user_handler", "UpdateProfileAvatar", requireId, err)
		}
	}()

	currentSession := http_utils.MustGetSessionFromContext(r.Context())

	// Max size - 10 Mb
	err = r.ParseMultipartForm(10 * 1024 * 1024)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrFileNotRead, http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("avatar")
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrFileNotRead, http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileUrl, err := h.UserUCase.SetAvatar(currentSession.UserData.Id, &file, header)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponse(w, models.Avatar{Url: fileUrl}, http.StatusOK)
}

// Handle get user avatar
func (h *UserHandler) GetProfileAvatar(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("user_handler", "GetProfileAvatar", requireId, err)
		}
	}()

	currentSession := http_utils.MustGetSessionFromContext(r.Context())

	fileUrl, err := h.UserUCase.GetAvatar(currentSession.UserData.Id)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrUserNotFound, http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponse(w, models.Avatar{Url: fileUrl}, http.StatusOK)
}

// Handle get profile of current user
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("user_handler", "GetProfile", requireId, err)
		}
	}()

	currentSession := http_utils.MustGetSessionFromContext(r.Context())

	profileUser, err := h.UserUCase.GetUserById(currentSession.UserData.Id)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponse(w, profileUser, http.StatusOK)
}

// Handle signup user
func (h *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("user_handler", "Signup", requireId, err)
		}
	}()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var newUser models.SignupUser
	err = json.Unmarshal(body, &newUser)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrCanNotUnmarshal, http.StatusBadRequest)
		return
	}
	newUser.Sanitize()

	err = validator.ValidateStruct(newUser)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusBadRequest)
		return
	}

	addedUserId, err := h.UserUCase.AddUser(&newUser)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusConflict)
		return
	}

	currentSession, err := h.SessionUCase.CreateNewSession(addedUserId)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetCookie(w, models.SessionCookieName, currentSession.Value, models.ExpireSessionCookie*time.Second)
	http_utils.SetJSONResponseSuccess(w, http.StatusCreated)
}

// Handle logout user
func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("user_handler", "Logout", requireId, err)
		}
	}()

	// Middleware auth add session in context
	currentSession := http_utils.MustGetSessionFromContext(r.Context())

	err = h.SessionUCase.DestroySession(currentSession.Value)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	// Auth middleware control existence of session cookie
	sessionCookie, _ := r.Cookie(models.SessionCookieName)
	http_utils.DestroyCookie(w, sessionCookie)

	http_utils.SetJSONResponseSuccess(w, http.StatusOK)
}
