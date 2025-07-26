package models

import (
	"fmt"
	"time"
)

type Subscription struct {
	ID          int        `json:"id"`
	ServiceName string     `json:"service_name" binding:"required" example:"netflix"`
	Price       int        `json:"price" binding:"required,gt=0" example:"100"`
	UserID      string     `json:"user_id" binding:"required,uuid" example:"f47ac10b-58cc-4372-a567-0e02b2c3d479"`
	StartDate   time.Time  `json:"-" example:"12-2025"`
	EndDate     *time.Time `json:"-" example:"12-2025"`

	RawStartDate string `json:"start_date" binding:"required" example:"12-2025"`
	RawEndDate   string `json:"end_date,omitempty" example:"12-2025"`
}

func (s *Subscription) ParseDates() error {
	startDate, err := time.Parse("01-2006", s.RawStartDate)
	if err != nil {
		return fmt.Errorf("invalid start_date format, expected MM-YYYY")
	}
	s.StartDate = startDate
	if s.RawEndDate != "" {
		endDate, err := time.Parse("01-2006", s.RawEndDate)
		if err != nil {
			return fmt.Errorf("invalid end_date format, expected MM-YYYY")
		}
		s.EndDate = &endDate
	}
	return nil
}

func (s *Subscription) FormatDates() {
	s.RawStartDate = s.StartDate.Format("01-2006")
	if s.EndDate != nil {
		s.RawEndDate = s.EndDate.Format("01-2006")
	}
}

type CUSubscription struct {
	ServiceName string `json:"service_name" binding:"required" example:"netflix"`
	Price       int    `json:"price" binding:"gt=0" example:"100"`
	UserID      string `json:"user_id" binding:"uuid" example:"f47ac10b-58cc-4372-a567-0e02b2c3d479"`
	StartDate   string `json:"start_date" example:"01-2025"`
	EndDate     string `json:"end_date,omitempty" example:"12-2025"`
}
type PeriodRequest struct {
	ServiceName string `json:"service_name,omitempty" example:"netflix"`
	Start       string `json:"start" binding:"required" example:"01-2025"`
	End         string `json:"end" binding:"required" example:"12-2025"`
	UserID      string `json:"user_id,omitempty" binding:"omitempty,uuid" example:"f47ac10b-58cc-4372-a567-0e02b2c3d479"`
}
