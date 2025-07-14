package repositories

import (
	"context"
	"go-EdTech/logger"
	"go-EdTech/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type UsersRepository struct {
	db *pgxpool.Pool
}

func NewUsersRepository(conn *pgxpool.Pool) *UsersRepository {
	return &UsersRepository{db: conn}
}

// FWP
func (r *UsersRepository) Create(c context.Context, user models.User) (uuid.UUID, error) {
	logger := logger.GetLogger()

	var id uuid.UUID
	err := r.db.QueryRow(c, "insert into Users(user_name, user_surname, user_type, status, email, password_hash) values($1, $2, $3, $4, $5, $6) returning uuid", user.User_Name, user.User_Surname, user.User_Type, "active", user.Email, user.PasswordHash).Scan(&id)
	if err != nil {
		logger.Error("could not query database", zap.String("db_msg", err.Error()))
		return uuid.Nil, err
	}
	return id, nil
}

// FWP
func (r *UsersRepository) FindOne(c context.Context, strUUID string) (models.User, error) {
	logger := logger.GetLogger()
	parsed_UUID, err := uuid.Parse(strUUID)
	if err != nil {
		logger.Error("could not parse UUID", zap.String("Backend_msg", err.Error()))
		return models.User{}, err
	}

	var user models.User
	row := r.db.QueryRow(c, "select uuid, user_name, user_surname, user_type, email, password_hash from Users where uuid = $1", parsed_UUID)
	err = row.Scan(&user.Id, &user.User_Name, &user.User_Surname, &user.User_Type, &user.Email, &user.PasswordHash)

	if err != nil {
		logger.Error("could not scan query row", zap.String("db_msg", err.Error()))
		return models.User{}, err
	}

	return user, nil
}

func (r *UsersRepository) FindAll(c context.Context) ([]models.User, error) {
	logger := logger.GetLogger()

	rows, err := r.db.Query(c, "select uuid, user_name, user_surname, user_type, email, password_hash from Users")
	if err != nil {
		logger.Error("could not query database", zap.String("db_msg", err.Error()))
		return nil, err
	}

	users := make([]models.User, 0)
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Id, &user.User_Name, &user.User_Surname, &user.User_Type, &user.Email, &user.PasswordHash)
		if err != nil {
			logger.Error("could not scan query row", zap.String("db_msg", err.Error()))
			return nil, err
		}

		users = append(users, user)
	}
	if rows.Err() != nil {
		return nil, err
	}

	return users, nil
}

func (r *UsersRepository) Update(c context.Context, user models.User, strUUID string) error {
	logger := logger.GetLogger()
	parsed_UUID, err := uuid.Parse(strUUID)
	if err != nil {
		logger.Error("could not parse UUID", zap.String("Backend_msg", err.Error()))
		return err
	}

	_, err = r.db.Exec(c, "update users set user_name = $1, user_surname =$2, email = $3 where uuid = $4", user.User_Name, user.User_Surname, user.Email, parsed_UUID)
	if err != nil {
		logger.Error("could not execute in database", zap.String("db_msg", err.Error()))
		return err
	}

	return nil
}

func (r *UsersRepository) ChangePassword(c context.Context, password []byte, strUUID string) error {
	logger := logger.GetLogger()
	parsed_UUID, err := uuid.Parse(strUUID)
	if err != nil {
		logger.Error("could not parse UUID", zap.String("Backend_msg", err.Error()))
		return err
	}

	_, err = r.db.Exec(c, "update users set password_hash = $1 where uuid = $2", password, parsed_UUID)
	if err != nil {
		logger.Error("could not execute in database", zap.String("db_msg", err.Error()))
		return err
	}
	return nil
}

func (r *UsersRepository) Delete(c context.Context, strUUID string) error {
	logger := logger.GetLogger()
	parsed_UUID, err := uuid.Parse(strUUID)
	if err != nil {
		logger.Error("could not parse UUID", zap.String("Backend_msg", err.Error()))
		return err
	}

	_, err = r.db.Exec(c, "delete from users where uuid = $1", parsed_UUID)
	if err != nil {
		logger.Error("could not execute in database", zap.String("db_msg", err.Error()))
		return err
	}
	return nil
}

func (u *UsersRepository) FindByEmail(c context.Context, email string) (models.User, error) {
	logger := logger.GetLogger()

	var user models.User
	row := u.db.QueryRow(c, "select uuid, user_name, user_surname, user_type, email, password_hash from users where email = $1", email)
	err := row.Scan(&user.Id, &user.User_Name, &user.User_Surname, &user.User_Type, &user.Email, &user.PasswordHash)

	if err != nil {
		logger.Error("could not scan query row", zap.String("db_msg", err.Error()))
		return models.User{}, err
	}

	return user, nil
}

func (r *UsersRepository) Deactivate(c context.Context, strUUID string) error {
	logger := logger.GetLogger()
	parsed_UUID, err := uuid.Parse(strUUID)
	if err != nil {
		logger.Error("could not parse UUID", zap.String("Backend_msg", err.Error()))
		return err
	}

	_, err = r.db.Exec(c, "update users set status = 'inactive' where uuid = $1", parsed_UUID)
	if err != nil {
		logger.Error("could not execute in database", zap.String("db_msg", err.Error()))
		return err
	}

	return nil
}

func (r *UsersRepository) Activate(c context.Context, strUUID string) error {
	logger := logger.GetLogger()
	parsed_UUID, err := uuid.Parse(strUUID)
	if err != nil {
		logger.Error("could not parse UUID", zap.String("Backend_msg", err.Error()))
		return err
	}

	_, err = r.db.Exec(c, "update users set status = ' active' where uuid = $1", parsed_UUID)
	if err != nil {
		logger.Error("could not execute in database", zap.String("db_msg", err.Error()))
		return err
	}

	return nil
}
