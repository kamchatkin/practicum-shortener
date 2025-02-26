package models

import "time"

// Alias Сокращение
type Alias struct {
	ID        int64
	Alias     string
	Source    string
	Quantity  int64
	CreatedAt time.Time
	UserID    int64
	DeletedAt time.Time
}

// Redirected засчет редиректа в стату
func (a *Alias) Redirected() {
	a.Quantity++
	// @todo реализовать
}

func (a *Alias) Found() bool {
	return a.Alias != ""
}

func (a *Alias) NotFound() bool {
	return a.Alias == ""
}

func (a *Alias) IsDeleted() bool {
	return !a.DeletedAt.IsZero()
}
