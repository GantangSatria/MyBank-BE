package domain

import "time"

type User struct {
	ID                       uint64
	Name                     *string
	Email                    string
	Phone                    *string
	PIN						 *string
	Password                 string  
	Gender                   *string    // L / P
	DateOfBirth              *time.Time
	Occupation               *string
	MaritalStatus            *string    // single / menikah
	Segment                  *string
	MonthlyIncome            float64
	MonthlyIncomeRange       *string    // "10-20 juta", ">20 juta", etc.
	ABGroup                  string
	IsPersonalizationEnabled bool
	IsActive                 bool
	LastLoginAt              *time.Time
	CreatedAt                time.Time
	UpdatedAt                time.Time
}