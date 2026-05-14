package models

import "time"

type Department struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null;size:200" json:"name"`
	ParentID  *uint     `gorm:"index" json:"parent_id,omitempty"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
}
