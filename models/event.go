package models

import "time"

type Event struct {
    Id int `json:"id"`
    Name string `json:"name"`
    Attributes map[string]string `json:"attributes"`
    Created_at time.Time `json:"created_at"`
}
