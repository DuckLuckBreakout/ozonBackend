package handler

import (
	"encoding/json"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/api"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/usecase"
	"io/ioutil"
	"net/http"
	"time"

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

// Login godoc
// @Summary Авторизация.
// @Description Авторизация пользователя.
// @Accept json
// @Produce json
// @Param Review body api.ApiLoginUser true "Данные пользователя."
// @Success 200 {object} usecase.Session "Jwt auth токен успешно получен."
// @Failure 400 {object} errors.Error "Некорректное тело запроса."
// @Failure 500 {object} errors.Error "Непредвиденная ошибка сервера."
// @Router /user/login [POST]
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

	var authUser api.ApiLoginUser
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

	userId, err := h.UserUCase.Authorize(&usecase.LoginUser{
		Email:    authUser.Email,
		Password: authUser.Password,
	})
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusUnauthorized)
		return
	}

	currentSession, err := h.SessionUCase.CreateNewSession(&usecase.UserId{Id: userId.Id})
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetCookie(w, usecase.SessionCookieName, currentSession.Value, usecase.ExpireSessionCookie*time.Second)
	http_utils.SetJSONResponseSuccess(w, http.StatusOK)
}

// UpdateProfile godoc
// @Summary Обновление профиля.
// @Description Обновление пользовательских данных.
// @Accept json
// @Produce json
// @Param UpdateUser body api.ApiUpdateUser true "Обновление пользовательских данных."
// @Success 200 {object} errors.Error "Данные пользователя успешно обновлены."
// @Failure 400 {object} errors.Error "Некорректное тело запроса."
// @Failure 500 {object} errors.Error "Непредвиденная ошибка сервера."
// @Router /user/profile [PUT]
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

	var updateUser api.ApiUpdateUser
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

	err = h.UserUCase.UpdateProfile(
		&usecase.UserId{Id: currentSession.UserData.Id},
		&usecase.UpdateUser{
			FirstName: updateUser.FirstName,
			LastName:  updateUser.LastName,
		},
	)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponseSuccess(w, http.StatusOK)
}

// UpdateProfileAvatar godoc
// @Summary Обновление аватара.
// @Description Обновление пользовательского аватара.
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Param avatar formData file true "Аватарка пользователя."
// @Success 200 {object} errors.Error "Аватар пользователя успешно обновлён."
// @Failure 400 {object} errors.Error "Некорректное тело запроса."
// @Failure 500 {object} errors.Error "Непредвиденная ошибка сервера."
// @Router /user/profile/avatar [PUT]
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

	fileUrl, err := h.UserUCase.SetAvatar(&usecase.UserId{Id: currentSession.UserData.Id}, &file, header)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponse(w, api.ApiAvatar{Url: fileUrl}, http.StatusOK)
}

// GetProfileAvatar godoc
// @Summary Получение аватара.
// @Description Получение пользовательского аватара.
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Success 200 {object} api.ApiAvatar "Аватар пользователя успешно получен."
// @Failure 400 {object} errors.Error "Некорректное тело запроса."
// @Failure 500 {object} errors.Error "Непредвиденная ошибка сервера."
// @Router /user/profile/avatar [GET]
func (h *UserHandler) GetProfileAvatar(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("user_handler", "GetProfileAvatar", requireId, err)
		}
	}()

	currentSession := http_utils.MustGetSessionFromContext(r.Context())

	fileUrl, err := h.UserUCase.GetAvatar(&usecase.UserId{Id: currentSession.UserData.Id})
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrUserNotFound, http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponse(w, api.ApiAvatar{Url: fileUrl}, http.StatusOK)
}

// GetProfile godoc
// @Summary Получение профиля.
// @Description Получение пользовательского профиля.
// @Accept json
// @Produce json
// @Success 200 {object} api.ApiProfileUser "Профиль пользователя успешно получен."
// @Failure 400 {object} errors.Error "Некорректное тело запроса."
// @Failure 500 {object} errors.Error "Непредвиденная ошибка сервера."
// @Router /user/profile [GET]
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("user_handler", "GetProfile", requireId, err)
		}
	}()

	currentSession := http_utils.MustGetSessionFromContext(r.Context())

	profileUser, err := h.UserUCase.GetUserById(&usecase.UserId{Id: currentSession.UserData.Id})
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponse(w, api.ApiProfileUser{
		Id:        profileUser.Id,
		FirstName: profileUser.FirstName,
		LastName:  profileUser.LastName,
		Avatar: api.ApiAvatar{
			Url: profileUser.Avatar.Url,
		},
		AuthId: profileUser.AuthId,
		Email:  profileUser.Email,
	}, http.StatusOK)
}

// Signup godoc
// @Summary Регистрация.
// @Description Регистрация нового пользователя.
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Param Review body api.ApiSignupUser true "Данные пользователя."
// @Success 200 {object} usecase.Session "Jwt auth токен успешно получен."
// @Failure 400 {object} errors.Error "Некорректное тело запроса."
// @Failure 500 {object} errors.Error "Непредвиденная ошибка сервера."
// @Router /user/signup [POST]
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

	var newUser api.ApiSignupUser
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

	addedUserId, err := h.UserUCase.AddUser(&usecase.SignupUser{
		Email:    newUser.Email,
		Password: newUser.Password,
	})
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusConflict)
		return
	}

	currentSession, err := h.SessionUCase.CreateNewSession(&usecase.UserId{Id: addedUserId.Id})
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetCookie(w, usecase.SessionCookieName, currentSession.Value, usecase.ExpireSessionCookie*time.Second)
	http_utils.SetJSONResponseSuccess(w, http.StatusCreated)
}

// Logout godoc
// @Summary Выход.
// @Description Выход из аккаунта пользователя.
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Success 200 {object} errors.Error "Выход из аккаунта успешно выполнен."
// @Failure 400 {object} errors.Error "Некорректное тело запроса."
// @Failure 500 {object} errors.Error "Непредвиденная ошибка сервера."
// @Router /user/logout [DELETE]
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
	sessionCookie, _ := r.Cookie(usecase.SessionCookieName)
	http_utils.DestroyCookie(w, sessionCookie)

	http_utils.SetJSONResponseSuccess(w, http.StatusOK)
}
