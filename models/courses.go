package models

import "time"

type Course struct {
	Course_id    int
	Course_title string
	Description  string
	Is_published bool
	Created_at   time.Time
	Updated_at   time.Time
}
