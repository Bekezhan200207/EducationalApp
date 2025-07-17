package models

import "time"

type Lesson struct {
	Id              int
	Title           string
	Description     string
	Subject_id      int //under question
	Order           int
	Level           string
	Interest        string
	Target_age_min  int
	Target_age_max  int
	Video_data      []byte
	Video_filename  string
	Video_mime_type string
	Duration_sec    int
	Is_published    bool
	Created_at      time.Time
	Updated_at      time.Time
}
