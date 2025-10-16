package models

import "time"

type Course struct {
	Id           int
	Name         string
	Description  string
	Is_published bool
	Created_at   time.Time
	Updated_at   time.Time
}
