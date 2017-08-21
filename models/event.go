package models

import "time"

type Event struct {
    Id int
    Name string
    Attributes map[string]string
    Created_at time.Time
}
