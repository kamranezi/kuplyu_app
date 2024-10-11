package models

import "time"

type Request struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    UserID    uint      `json:"user_id"`
    Category  string    `json:"category"`
    Brand     string    `json:"brand"`
    Model     string    `json:"model"`
    Year      int       `json:"year"`
    PartName  string    `json:"part_name"`
    MaxPrice  float64   `json:"max_price"`
    Distance  int       `json:"distance"`
    CreatedAt time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}
