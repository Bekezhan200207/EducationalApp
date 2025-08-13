package repositories

import (
	"context"
	"go-EdTech/logger"
	"go-EdTech/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Lessonsrepository struct {
	db *pgxpool.Pool
}

func NewLessonsRepository(conn *pgxpool.Pool) *Lessonsrepository {
	return &Lessonsrepository{db: conn}
}

func (r *Lessonsrepository) FindById(c context.Context, id int) (models.Lesson, error) {
	sql :=
		`
	select 
	l.lesson_id,
	l.lesson_title,
	l.description,
	l.subject_id,
	l.order,
	l.level,
	l.interest,
	l.target_age_min,
	l.target_age_max,
	l.video_data,
	l.video_filename,
	l.video_mime_type,
	l.duration_sec,
	l.is_published,
	l.created_at,
	l.updated_at,
	s.subject_id,
	s.subject_name
	from lessons l
	join subjects s on l.subject_id = s.subject_id,
	extract(epoch from l.created_at) as created_at,
	extract(epoch from l.updated_at) as updated_at
	where l.lesson_id = $1
	`

	logger := logger.GetLogger()

	rows, err := r.db.Query(c, sql, id)
	if err != nil {
		logger.Error("could not query database", zap.String("db_msg", err.Error()))
		return models.Lesson{}, err
	}
	defer rows.Close()

	var lesson *models.Lesson

	for rows.Next() {
		var les models.Lesson
		var sub models.Subject

		err := rows.Scan(
			&les.Id,
			&les.Title,
			&les.Description,
			&les.Subject_id,
			&les.Order,
			&les.Level,
			&les.Interest,
			&les.Target_age_min,
			&les.Target_age_max,
			&les.Video_data,
			&les.Video_filename,
			&les.Video_mime_type,
			&les.Duration_sec,
			&les.Is_published,
			&les.Created_at,
			&les.Updated_at,
			&sub.Id,
			&sub.Title,
		)

		if err != nil {
			logger.Error("could not scan query row", zap.String("db_msg", err.Error()))
			return models.Lesson{}, err
		}

		les.Created_at = time.Unix(les.Created_at.Unix(), 0)
		les.Updated_at = time.Unix(les.Updated_at.Unix(), 0)

		les.Subject_id = sub.Id
		lesson = &les

	}

	err = rows.Err()
	if err != nil {
		logger.Error(err.Error())
		return models.Lesson{}, err
	}

	return *lesson, nil

}

func (r *Lessonsrepository) FindAll(c context.Context) ([]models.Lesson, error) {
	sql :=
		`
	select 
	l.lesson_id,
	l.lesson_title,
	l.description,
	l.subject_id,
	l.order,
	l.level,
	l.interest,
	l.target_age_min,
	l.target_age_max,
	l.video_data,
	l.video_filename,
	l.video_mime_type,
	l.duration_sec,
	l.is_published,
	l.created_at,
	l.updated_at,
	s.subject_id,
	s.subject_name
	from lessons l
	join subjects s on l.subject_id = s.subject_id,
	extract(epoch from l.created_at) as created_at,
	extract(epoch from l.updated_at) as updated_at
	`

	logger := logger.GetLogger()

	rows, err := r.db.Query(c, sql)

	if err != nil {
		logger.Error("could not query database", zap.String("db_msg", err.Error()))
		return nil, err
	}

	lessons := make([]models.Lesson, 0)

	for rows.Next() {
		var les models.Lesson
		var sub models.Subject

		err := rows.Scan(
			&les.Id,
			&les.Title,
			&les.Description,
			&les.Subject_id,
			&les.Order,
			&les.Level,
			&les.Interest,
			&les.Target_age_min,
			&les.Target_age_max,
			&les.Video_data,
			&les.Video_filename,
			&les.Video_mime_type,
			&les.Duration_sec,
			&les.Is_published,
			&les.Created_at,
			&les.Updated_at,
			&sub.Id,
			&sub.Title,
		)
		if err != nil {
			logger.Error("could not scan query row", zap.String("db_msg", err.Error()))
			return []models.Lesson{}, err
		}

		lessons = append(lessons, les)

	}
	err = rows.Err()
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return lessons, nil

}

func (r *Lessonsrepository) Create(c context.Context, lesson models.Lesson) (int, error) {
	var id int

	logger := logger.GetLogger()

	row := r.db.QueryRow(c,
		`
	insert into lessons
	(
	lesson_title, 
	description, 
	subject_id, 
	"order", 
	"level", 
	interest, 
	target_age_min, 
	target_age_max, 
	video_data, 
	video_filename, 
	video_mime_type, 
	duration_sec, 
	is_published
	) 
	values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	returning lesson_id
	`,

		lesson.Title,
		lesson.Description,
		lesson.Subject_id,
		lesson.Order,
		lesson.Level,
		lesson.Interest,
		lesson.Target_age_min,
		lesson.Target_age_max,
		lesson.Video_data,
		lesson.Video_filename,
		lesson.Video_mime_type,
		lesson.Duration_sec,
		lesson.Is_published,
	)

	err := row.Scan(&id)
	if err != nil {
		logger.Error("could not query database", zap.String("db_msg", err.Error()))
		return 0, err
	}

	return id, nil
}

func (r *Lessonsrepository) Update(c context.Context, id int, updLesson models.Lesson) error {
	logger := logger.GetLogger()

	_, err := r.db.Exec(
		c,
		`
	update lessons
	set 
	lesson_title = $1, 
	description = $2, 
	subject_id = $3, 
	"order" = $4, 
	"level" = $5, 
	interest = $6, 
	target_age_min = $7, 
	target_age_max = $8, 
	video_data = $9, 
	video_filename = $10, 
	video_mime_type = $11, 
	duration_sec = $12, 
	is_published = $13,
	updated_at = now()
	where lesson_id = $14
		`,
		updLesson.Title,
		updLesson.Description,
		updLesson.Subject_id,
		updLesson.Order,
		updLesson.Level,
		updLesson.Interest,
		updLesson.Target_age_min,
		updLesson.Target_age_max,
		updLesson.Video_data,
		updLesson.Video_filename,
		updLesson.Video_mime_type,
		updLesson.Duration_sec,
		updLesson.Is_published,
		id,
	)
	if err != nil {
		logger.Error("could not query database", zap.String("db_msg", err.Error()))
		return err
	}

	return nil
}

func (r *Lessonsrepository) Delete(c context.Context, id int) error {
	logger := logger.GetLogger()

	_, err := r.db.Exec(c, "delete from lessons where lesson_id = $1", id)
	if err != nil {
		logger.Error("could not execute in database", zap.String("db_msg", err.Error()))
		return err
	}
	return nil
}
