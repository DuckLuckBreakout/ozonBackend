package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/admin"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/errors"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/tools/http_utils"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/tools/validator"
	"github.com/DuckLuckBreakout/ozonBackend/pkg/tools/logger"
)

type AdminHandler struct {
	AdminUCase admin.UseCase
}

func NewHandler(adminUCase admin.UseCase) admin.Handler {
	return &AdminHandler{
		AdminUCase: adminUCase,
	}
}

func (h *AdminHandler) ChangeOrderStatus(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("admin_handler", "ChangeOrderStatus", requireId, err)
		}
	}()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var updateOrder models.UpdateOrder
	err = json.Unmarshal(body, &updateOrder)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrCanNotUnmarshal, http.StatusBadRequest)
		return
	}
	updateOrder.Sanitize()

	err = validator.ValidateStruct(updateOrder)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusBadRequest)
		return
	}

	err = h.AdminUCase.ChangeOrderStatus(&updateOrder)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrProductNotFound, http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponseSuccess(w, http.StatusOK)
}
