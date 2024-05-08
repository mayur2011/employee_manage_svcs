package mapstore

import (
	"fmt"
	"sort"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"employee_manage_svcs/domain"
)

// Mutex to protect InMemoryStore
var mu sync.Mutex

// Structure defination
type MapStore struct {
	store map[int]domain.Employee
}

// Factory method gives a new instance of MapStore, Also this is for caller packages, not for mapstoer itself
func NewMapStore() *MapStore {
	return &MapStore{store: make(map[int]domain.Employee)}
}

func (ms *MapStore) isRecordExists(id int) bool {
	_, ok := ms.store[id]
	return ok
}

func (ms *MapStore) Create(employee domain.Employee) error {
	mu.Lock()
	defer mu.Unlock()
	if ms.isRecordExists(employee.ID) {
		return fmt.Errorf("employee already exists")
		//return domain.ErrorIDExists
	}
	employee.CreatedAt = time.Now().UTC()
	ms.store[employee.ID] = employee
	log.Println("Employee has been created")
	return nil
}

func (ms *MapStore) GetAll(page, size int) ([]domain.Employee, error) {
	mu.Lock()
	defer mu.Unlock()
	// Create a slice of keys (employee IDs)
	var keys []int
	for k := range ms.store {
		keys = append(keys, k)
	}

	// Sort the keys based on timestamps (created_at)
	sort.Slice(keys, func(i, j int) bool {
		return ms.store[keys[i]].CreatedAt.Before(ms.store[keys[j]].CreatedAt)
	})

	startIndex := page * size
	endIndex := startIndex + size

	// Final result of employees
	var employees []domain.Employee
	var i int
	for _, id := range keys {
		if i >= startIndex && i < endIndex {
			emp := ms.store[id]
			employees = append(employees, emp)
		}
		i++
	}
	return employees, nil
}

func (ms *MapStore) GetById(id int) (domain.Employee, error) {
	mu.Lock()
	defer mu.Unlock()
	if ms.isRecordExists(id) {
		employee := ms.store[id]
		return employee, nil
	}
	return domain.Employee{}, fmt.Errorf("employee does not exist for this id")
}

func (ms *MapStore) Update(id int, employee domain.Employee) error {
	mu.Lock()
	defer mu.Unlock()
	if ms.isRecordExists(id) {
		ms.store[id] = employee
		log.Println("Employee has been updated")
		return nil
	}
	return fmt.Errorf("employee does not exist for this id")
}

func (ms *MapStore) Delete(id int) error {
	mu.Lock()
	defer mu.Unlock()
	if ms.isRecordExists(id) {
		delete(ms.store, id)
		log.Println("Employee has been deleted")
		return nil
	}
	return fmt.Errorf("employee does not exist for this id")
}
