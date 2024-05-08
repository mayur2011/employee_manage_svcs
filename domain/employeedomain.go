package domain

import (
	"errors"
	"time"
)

// ErrorIDExists is an error value for duplicate customer id
var ErrorIDExists = errors.New("Employee ID exists")

type Employee struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Position  string  `json:"position"`
	Salary    float64 `json:"salary"`
	CreatedAt time.Time
}

type EmployeeStore interface {
	Create(Employee) error
	Update(int, Employee) error
	Delete(int) error
	GetById(int) (Employee, error)
	GetAll(int, int) ([]Employee, error)
}
