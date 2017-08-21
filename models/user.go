package models

import "time"

type User struct {
    Id int
    Os string
    Device string
    Locale string
    Voiceover bool
    Bold_text bool
    Reduce_motion bool
    Reduce_transparency bool
    Created_at time.Time
}
