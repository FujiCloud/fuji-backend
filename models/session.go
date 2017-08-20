package models

import (
    "time"
)

type Session struct {
    Id int
    User_id int
    Duration int
    Created_at time.Time
}
