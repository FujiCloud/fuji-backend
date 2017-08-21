package models

import "time"

type Session struct {
    Id int `json:"id"`
    User_id int `json:"user_id"`
    Duration *int `json:"duration"`
    Created_at time.Time `json:"created_at"`
}
