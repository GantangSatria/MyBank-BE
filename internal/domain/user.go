package domain

import "time"

type User struct {
	ID                       uint64
	Name                     string
	Email                    string
	Phone                    string
	Password                 string  
	DateOfBirth              *time.Time
	Occupation               string
	Segment                  string
	IsPersonalizationEnabled bool
	IsActive                 bool
	CreatedAt                time.Time
	UpdatedAt                time.Time
}