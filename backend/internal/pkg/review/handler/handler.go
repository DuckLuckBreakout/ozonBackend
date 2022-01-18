package handler

import (
	"encoding/json"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/api"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/usecase"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/review"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/errors"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/tools/http_utils"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/tools/validator"
	"github.com/DuckLuckBreakout/ozonBackend/pkg/tools/logger"

	"github.com/gorilla/mux"
)

type ReviewHandler struct {
	ReviewUCase review.UseCase
}

func NewHandler(reviewUCase review.UseCase) review.Handler {
	return &ReviewHandler{
		ReviewUCase: reviewUCase,
	}
}

// Get statistics about reviews for product
func (h *ReviewHandler) GetReviewStatistics(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("review_handler", "GetReviewInfo", requireId, err)
		}
	}()

	vars := mux.Vars(r)
	productId, err := strconv.Atoi(vars["id"])
	if err != nil || productId < 1 {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}

	productById, err := h.ReviewUCase.GetStatisticsByProductId(&usecase.ProductId{Id: uint64(productId)})
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrProductNotFound, http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponse(w, api.ApiReviewStatistics{
		Stars: productById.Stars,
	}, http.StatusOK)
}

// Check rights for write new review
func (h *ReviewHandler) CheckReviewRights(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("review_handler", "AddCompletedOrder", requireId, err)
		}
	}()

	currentSession := http_utils.MustGetSessionFromContext(r.Context())

	vars := mux.Vars(r)
	productId, err := strconv.Atoi(vars["id"])
	if err != nil || productId < 1 {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}

	err = h.ReviewUCase.CheckReviewUserRights(
		&usecase.UserId{Id: currentSession.UserData.Id},
		&usecase.ProductId{Id: uint64(productId)},
	)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusConflict)
		return
	}

	http_utils.SetJSONResponseSuccess(w, http.StatusOK)
}

// Add new review for product
func (h *ReviewHandler) AddNewReview(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("review_handler", "AddNewReview", requireId, err)
		}
	}()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var userReview api.ApiReview
	err = json.Unmarshal(body, &userReview)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrCanNotUnmarshal, http.StatusBadRequest)
		return
	}
	userReview.Sanitize()

	err = validator.ValidateStruct(userReview)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusBadRequest)
		return
	}

	currentSession := http_utils.MustGetSessionFromContext(r.Context())

	err = h.ReviewUCase.AddNewReviewForProduct(
		&usecase.UserId{Id: currentSession.UserData.Id},
		&usecase.Review{
			ProductId:     userReview.ProductId,
			Rating:        userReview.Rating,
			Advantages:    userReview.Advantages,
			Disadvantages: userReview.Disadvantages,
			Comment:       userReview.Comment,
			IsPublic:      userReview.IsPublic,
		},
	)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrCanNotAddReview, http.StatusBadRequest)
		return
	}

	http_utils.SetJSONResponseSuccess(w, http.StatusOK)
}

// Get all reviews for product
func (h *ReviewHandler) GetReviewsForProduct(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("review_handler", "GetReviewsForProduct", requireId, err)
		}
	}()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var paginator api.ApiPaginatorReviews
	err = json.Unmarshal(body, &paginator)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrCanNotUnmarshal, http.StatusBadRequest)
		return
	}

	err = validator.ValidateStruct(paginator)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	productId, err := strconv.Atoi(vars["id"])
	if err != nil || productId < 1 {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}

	reviews, err := h.ReviewUCase.GetReviewsByProductId(
		&usecase.ProductId{Id: uint64(productId)},
		&usecase.PaginatorReviews{
			PageNum: paginator.PageNum,
			Count:   paginator.Count,
			SortReviewsOptions: usecase.SortReviewsOptions{
				SortKey:       paginator.SortKey,
				SortDirection: paginator.SortDirection,
			},
		},
	)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrProductNotFound, http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponse(w, &api.ApiRangeReviews{
		ListPreviews:  reviews.ListPreviews,
		MaxCountPages: reviews.MaxCountPages,
	}, http.StatusOK)
}
