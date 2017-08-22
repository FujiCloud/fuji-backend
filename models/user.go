package models

import "time"

type User struct {
    Id int `json:"id"`
    Os string `json:"os"`
    Device string `json:"device"`
    Locale string `json:"locale"`
    Voiceover bool `json:"voiceover"`
    Bold_text bool `json:"bold_text"`
    Reduce_motion bool `json:"reduce_motion"`
    Reduce_transparency bool `json:"reduce_transparency"`
    Created_at time.Time `json:"created_at"`
}
