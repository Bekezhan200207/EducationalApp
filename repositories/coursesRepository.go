package repositories

import (
	"context"
	"go-EdTech/logger"
	"go-EdTech/models"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Coursesrepository struct {
	db *pgxpool.Pool
}

func NewCoursesRepository(conn *pgxpool.Pool) *Coursesrepository {
	return &Coursesrepository{db: conn}
}

func (r *Coursesrepository) Create(c context.Context, course models.Course) (int, error) {
	logger := logger.GetLogger()

	row := r.db.QueryRow(c, `insert into courses (name, description, is_published) values ($1, $2, $3) returning id`, course.Name, course.Description, course.Is_published)
	err := row.Scan(&course.Id)
	if err != nil {
		logger.Error("could not scan query row", zap.String("db_msg", err.Error()))
		return 0, err
	}
	return course.Id, nil
}

func (r *Coursesrepository) FindById(c context.Context, courseId int) (models.Course, error) {
	logger := logger.GetLogger()

	var course models.Course
	row := r.db.QueryRow(c, `select id, name, description, is_published, created_at, updated_at from courses where id = $1`, courseId)
	if err := row.Scan(&course.Id, &course.Name, &course.Description, &course.Is_published, &course.Created_at, &course.Updated_at); err != nil {
		logger.Error("could not scan query row", zap.String("db_msg", err.Error()))
		return models.Course{}, err
	}
	return course, nil
}

func (r *Coursesrepository) FindAll(c context.Context) ([]models.Course, error) {
	logger := logger.GetLogger()

	rows, err := r.db.Query(c, `select id, name, description, is_published, created_at, updated_at from courses`)
	if err != nil {
		logger.Error("could not query database", zap.String("db_msg", err.Error()))
		return []models.Course{}, err
	}
	defer rows.Close()

	courses := make([]models.Course, 0)

	for rows.Next() {
		var course models.Course
		err := rows.Scan(&course.Id, &course.Name, &course.Description, &course.Is_published, &course.Created_at, &course.Updated_at)
		if err != nil {
			logger.Error("could not scan query row", zap.String("db_msg", err.Error()))
			return nil, err
		}

		courses = append(courses, course)
	}

	return courses, nil
}

func (r *Coursesrepository) Update(c context.Context, id int, Updcourse models.Course) error {
	logger := logger.GetLogger()

	_, err := r.db.Exec(c, "update courses set name = $1, description = $2, is_published = $3 where id = $4", Updcourse.Name, Updcourse.Description, Updcourse.Is_published, id)
	if err != nil {
		logger.Error("could not execute in database", zap.String("db_msg", err.Error()))
		return err
	}
	return nil
}

func (r *Coursesrepository) Delete(c context.Context, id int) error {
	logger := logger.GetLogger()

	_, err := r.db.Exec(c, "delete from courses where id = $1", id)
	if err != nil {
		logger.Error("could not execute in database", zap.String("db_msg", err.Error()))
		return err
	}
	return nil
}
