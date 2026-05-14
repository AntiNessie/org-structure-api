package models

import "time"

type Employee struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	DepartmentID uint       `gorm:"not null;index" json:"department_id"`
	FullName     string     `gorm:"not null;size:200" json:"full_name"`
	Position     string     `gorm:"not null;size:200" json:"position"`
	HiredAt      *time.Time `json:"hired_at,omitempty"`
	CreatedAt    time.Time  `gorm:"not null" json:"created_at"`
}
