package repository

import (
	"database/sql"

	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/dto"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/user"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/errors"
)

type PostgresqlRepository struct {
	db *sql.DB
}

func NewSessionPostgresqlRepository(db *sql.DB) user.Repository {
	return &PostgresqlRepository{
		db: db,
	}
}

// Add new user profile
func (r *PostgresqlRepository) AddProfile(user *dto.DtoProfileUser) (*dto.DtoUserId, error) {
	row := r.db.QueryRow(
		"INSERT INTO data_users(id, email) "+
			"VALUES ($1, $2) RETURNING id",
		user.AuthId,
		user.Email,
	)

	var userId uint64
	if err := row.Scan(&userId); err != nil {
		return nil, errors.ErrDBInternalError
	}

	return &dto.DtoUserId{Id: userId}, nil
}

// Select one profile by id
func (r *PostgresqlRepository) SelectProfileById(userId *dto.DtoUserId) (*dto.DtoProfileUser, error) {
	row := r.db.QueryRow(
		"SELECT id, first_name, last_name, avatar, email "+
			"FROM data_users WHERE id = $1",
		userId.Id,
	)

	userById := dto.DtoProfileUser{}

	firstName := sql.NullString{}
	lastName := sql.NullString{}
	avatarUrl := sql.NullString{}
	err := row.Scan(
		&userById.Id,
		&firstName,
		&lastName,
		&avatarUrl,
		&userById.Email,
	)
	userById.FirstName = firstName.String
	userById.LastName = lastName.String
	userById.Avatar.Url = avatarUrl.String

	switch err {
	case sql.ErrNoRows:
		return nil, errors.ErrUserNotFound
	case nil:
		return &userById, nil
	default:
		return nil, errors.ErrDBInternalError
	}
}

// Update info in user profile
func (r *PostgresqlRepository) UpdateProfile(userId *dto.DtoUserId, user *dto.DtoUpdateUser) error {
	_, err := r.db.Exec(
		"UPDATE data_users SET "+
			"first_name = $1, "+
			"last_name = $2 "+
			"WHERE id = $3",
		user.FirstName,
		user.LastName,
		userId.Id,
	)
	if err != nil {
		return errors.ErrDBInternalError
	}

	return nil
}

// Update user avatar
func (r *PostgresqlRepository) UpdateAvatar(userId *dto.DtoUserId, avatarUrl string) error {
	_, err := r.db.Exec(
		"UPDATE data_users SET "+
			"avatar = $1 "+
			"WHERE id = $2",
		avatarUrl,
		userId.Id,
	)
	if err != nil {
		return errors.ErrDBInternalError
	}

	return nil
}
