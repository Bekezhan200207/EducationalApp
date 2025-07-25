package repositories

import (
	"context"
	"go-EdTech/logger"
	"go-EdTech/models"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type SubjectsRepository struct {
	db *pgxpool.Pool
}

func NewSubjectsRepository(conn *pgxpool.Pool) *SubjectsRepository {
	return &SubjectsRepository{db: conn}
}

func (r *SubjectsRepository) FindById(c context.Context, id int) (models.Subject, error) {
	logger := logger.GetLogger()

	var subject models.Subject
	row := r.db.QueryRow(c, "select subject_id, subject_name from subjects where subject_id = $1", id)
	err := row.Scan(&subject.Id, &subject.Title)
	if err != nil {
		logger.Error("could not scan query row", zap.String("db_msg", err.Error()))
		return models.Subject{}, err
	}

	return subject, nil
}

func (r *SubjectsRepository) FindAll(c context.Context) ([]models.Subject, error) {
	logger := logger.GetLogger()

	rows, err := r.db.Query(c, "select subject_id, subject_name from subjects")
	defer rows.Close()
	if err != nil {
		logger.Error("could not query database", zap.String("db_msg", err.Error()))
		return []models.Subject{}, err
	}

	subjects := make([]models.Subject, 0)

	for rows.Next() {
		var subject models.Subject
		err := rows.Scan(&subject.Id, &subject.Title)
		if err != nil {
			logger.Error("could not scan query row", zap.String("db_msg", err.Error()))
			return nil, err
		}

		subjects = append(subjects, subject)
	}

	return subjects, nil
}

func (r *SubjectsRepository) Create(c context.Context, subject models.Subject) (int, error) {
	logger := logger.GetLogger()

	var id int
	row := r.db.QueryRow(c, "insert into subjects(subject_name) values($1) returning subject_id", subject.Title)
	err := row.Scan(&id)
	if err != nil {
		logger.Error("could not scan query row", zap.String("db_msg", err.Error()))
		return 0, err
	}

	return id, nil

}

func (r *SubjectsRepository) Update(c context.Context, id int, Updsubject models.Subject) error {
	logger := logger.GetLogger()

	_, err := r.db.Exec(c, "update subjects set subject_name = $1 where subject_id = $2", Updsubject.Title, id)
	if err != nil {
		logger.Error("could not execute in database", zap.String("db_msg", err.Error()))
		return err
	}
	return nil
}

func (r *SubjectsRepository) Delete(c context.Context, id int) error {
	logger := logger.GetLogger()

	_, err := r.db.Exec(c, "delete from subjects where subject_id = $1", id)
	if err != nil {
		logger.Error("could not execute in database", zap.String("db_msg", err.Error()))
		return err
	}
	return nil
}
