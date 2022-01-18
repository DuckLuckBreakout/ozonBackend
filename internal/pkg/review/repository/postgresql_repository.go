package repository

import (
	"database/sql"
	"fmt"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/dto"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/usecase"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/review"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/errors"
	"math"
)

type PostgresqlRepository struct {
	db *sql.DB
}

func NewSessionPostgresqlRepository(db *sql.DB) review.Repository {
	return &PostgresqlRepository{
		db: db,
	}
}

// Select range of reviews
func (r *PostgresqlRepository) SelectRangeReviews(
	rangeReviews *dto.DtoRangeReviews,
	paginator *dto.DtoPaginatorReviews,
) ([]*dto.DtoViewReview, error) {
	rows, err := r.db.Query(
		"SELECT rating, advantages, disadvantages, comment, is_public, "+
			"date_added, user_id "+
			"FROM reviews "+
			"WHERE product_id = $1 "+
			rangeReviews.SortString+
			"LIMIT $2 OFFSET $3",
		rangeReviews.ProductId,
		paginator.Count,
		paginator.Count*(paginator.PageNum-1),
	)
	if err != nil {
		return nil, errors.ErrDBInternalError
	}
	defer rows.Close()

	reviews := make([]*dto.DtoViewReview, 0)
	for rows.Next() {
		userReview := &dto.DtoViewReview{}
		err = rows.Scan(
			&userReview.Rating,
			&userReview.Advantages,
			&userReview.Disadvantages,
			&userReview.Comment,
			&userReview.IsPublic,
			&userReview.DateAdded,
			&userReview.UserId,
		)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, userReview)
	}

	return reviews, nil
}

// Get count of all review pages for this product
func (r *PostgresqlRepository) GetCountPages(countPages *dto.DtoCountPages) (*dto.DtoCounter, error) {
	row := r.db.QueryRow(
		"SELECT count(id) "+
			"FROM reviews "+
			"WHERE product_id = $1",
	)

	var counter int
	if err := row.Scan(&counter); err != nil {
		return nil, errors.ErrDBInternalError
	}
	counter = int(math.Ceil(float64(counter) / 1))

	return &dto.DtoCounter{Count: counter}, nil
}

// Create sort string for query
func (r *PostgresqlRepository) CreateSortString(sortString *dto.DtoSortString) (string, error) {
	// Select order target
	var orderTarget string
	switch sortString.SortKey {
	case usecase.ReviewDateAddedSort:
		orderTarget = "date_added"
	default:
		return "", errors.ErrIncorrectPaginator
	}

	// Select order direction
	var orderDirection string
	switch sortString.SortDirection {
	case usecase.ReviewPaginatorASC:
		orderDirection = "ASC"
	case usecase.ReviewPaginatorDESC:
		orderDirection = "DESC"
	default:
		return "", errors.ErrIncorrectPaginator
	}

	return fmt.Sprintf("ORDER BY %s %s ", orderTarget, orderDirection), nil
}

// Select all statistics about reviews by product id
func (r *PostgresqlRepository) SelectStatisticsByProductId(productId *dto.DtoProductId) (*dto.DtoReviewStatistics, error) {
	rows, err := r.db.Query(
		"SELECT count(id), rating "+
			"FROM reviews "+
			"WHERE product_id = $1 "+
			"GROUP BY rating "+
			"ORDER BY rating",
	)
	if err != nil {
		return nil, errors.ErrIncorrectPaginator
	}
	defer rows.Close()

	statistics := &dto.DtoReviewStatistics{}
	statistics.Stars = make([]int, 5)
	var countStars int
	var rating int
	for rows.Next() {
		err = rows.Scan(
			&countStars,
			&rating,
		)
		if err != nil {
			return nil, err
		}
		statistics.Stars[rating-1] = countStars
	}

	return statistics, nil
}

// Check rights for review (the user has completed orders)
func (r *PostgresqlRepository) CheckReview(review *dto.DtoCheckReview) bool {
	row := r.db.QueryRow(
		"SELECT (SELECT count(us.id) "+
			"FROM user_orders us "+
			"JOIN ordered_products op ON us.id = op.order_id "+
			"WHERE (us.user_id = $1 AND op.product_id = $2)) - "+
			"(SELECT count(rv.id) FROM reviews rv "+
			"WHERE (rv.user_id = $1 AND rv.product_id = $2))",
		review.UserId,
		review.ProductId,
	)

	var isExist int
	if err := row.Scan(&isExist); err != nil || isExist == 0 {
		return false
	}

	return true
}

// Add new review for product
func (r *PostgresqlRepository) AddReview(userId *dto.DtoUserId, review *dto.DtoReview) (*dto.DtoReviewId, error) {
	row := r.db.QueryRow(
		"INSERT INTO reviews(product_id, user_id, rating, advantages, "+
			"disadvantages, comment, is_public) "+
			"VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		review.ProductId,
		userId.Id,
		review.Rating,
		review.Advantages,
		review.Disadvantages,
		review.Comment,
		review.IsPublic,
	)

	var reviewId uint64
	if err := row.Scan(&reviewId); err != nil {
		return nil, errors.ErrDBInternalError
	}

	return &dto.DtoReviewId{Id: reviewId}, nil
}
