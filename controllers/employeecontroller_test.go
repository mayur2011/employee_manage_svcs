package controllers_test

import (
	"encoding/json"

	"net/http"
	"net/http/httptest"
	"sort"
	"strconv"
	"strings"
	"time"

	"employee_manage_svcs/controllers"
	"employee_manage_svcs/domain"

	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EmployeeController", func() {
	var r *mux.Router
	var w *httptest.ResponseRecorder
	var store *FakeEmployeeStore
	var controller controllers.EmployeeController
	BeforeEach(func() {
		r = mux.NewRouter()
		store = newFakeEmployeeStore()
		controller = controllers.EmployeeController{
			Store: store,
		}
	})

	Describe("Post a employee which does not exist", func() {
		Context("Provide a valid employee data", func() {
			It("Should create a new employee and get HTTP Status: 201", func() {
				r.Handle("/employee", controllers.ResponseHandler(controller.PostEmployee)).Methods("POST")
				employeeJson := `{"ID": 211, "Name":"Andrea Le", "position": "Sr. SW Engg", "salary": 75000.00}`
				req, err := http.NewRequest(
					"POST",
					"/employee",
					strings.NewReader(employeeJson),
				)
				Expect(err).NotTo(HaveOccurred())
				w = httptest.NewRecorder()
				r.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(201))
			})
		})

		Context("Provide a employee data that contains dupliate employee ID", func() {
			It("Should get HTTP Status: 400", func() {
				r.Handle("/employee", controllers.ResponseHandler(controller.PostEmployee)).Methods("POST")
				employeeJson := `{"ID": 111, "Position":"Jr. SW Engg", "Name":"Anni Yi", "salary": 45000.00}`
				req, err := http.NewRequest(
					"POST",
					"/employee",
					strings.NewReader(employeeJson),
				)
				Expect(err).NotTo(HaveOccurred())
				w = httptest.NewRecorder()
				r.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(400))
			})
		})
	})

	Describe("Get a employee for given id", func() {
		Context("Get a specified employee from data store", func() {
			It("Should get a employee record", func() {
				r.Handle("/employees/{id}", controllers.ResponseHandler(controller.GetEmployeeById)).Methods("GET")
				empID := 112
				empIDStr := strconv.Itoa(empID)
				req, err := http.NewRequest("GET", "/employees/"+empIDStr, nil)
				Expect(err).NotTo(HaveOccurred())
				w = httptest.NewRecorder()
				r.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(200))
				//-- unmarshaling the api response for employee ID validation--
				var resp response
				json.Unmarshal(w.Body.Bytes(), &resp)
				tempData := resp.Data.(map[string]interface{})
				resultID := (int)(tempData["id"].(float64))
				Expect(resultID).To(Equal(empID))
			})
		})

		Context("Get 0 record from data store", func() {
			It("Should get a null employee record", func() {
				r.Handle("/employees/{id}", controllers.ResponseHandler(controller.GetEmployeeById)).Methods("GET")
				empID := 100
				empIDStr := strconv.Itoa(empID)
				req, err := http.NewRequest("GET", "/employees/"+empIDStr, nil)
				Expect(err).NotTo(HaveOccurred())
				w := httptest.NewRecorder()
				r.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(404))
				var resp response
				json.Unmarshal(w.Body.Bytes(), &resp)
				Expect(resp.Data).To(BeNil())
			})
		})
	})

	Describe("Update a employee for given id", func() {
		Context("Provide a valid employee data to update", func() {
			It("Should get a HTTP Status: 202", func() {
				r.Handle("/employees/{id}", controllers.ResponseHandler(controller.UpdateEmployee)).Methods("PUT")
				empID := 113
				empIDStr := strconv.Itoa(empID)
				employeeJson := `{"ID":113, "Position": "Sr. SW Engg", "Name": "Kishore Sharma", "salary": 155000.00}`
				req, err := http.NewRequest("PUT", "/employees/"+empIDStr, strings.NewReader(employeeJson))
				Expect(err).NotTo(HaveOccurred())
				w := httptest.NewRecorder()
				r.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(202))
			})
		})
	})

	Describe("Delete a employee for given id", func() {
		Context("Provide a valid employee id to delete", func() {
			It("Should get a HTTP Status: 202", func() {
				r.Handle("/employees/{id}", controllers.ResponseHandler(controller.DeleteEmployee)).Methods("DELETE")
				empID := 113
				empIDStr := strconv.Itoa(empID)
				req, err := http.NewRequest("DELETE", "/employees/"+empIDStr, nil)
				Expect(err).NotTo(HaveOccurred())
				w := httptest.NewRecorder()
				r.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(202))
			})
		})
	})

	Describe("Get list of employees", func() {
		Context("Get all employees from data store", func() {
			It("Should get list of employees", func() {
				r.Handle("/employees", controllers.ResponseHandler(controller.GetAllEmployees)).Methods("GET")
				req, err := http.NewRequest("GET", "/employees", nil)
				Expect(err).NotTo(HaveOccurred())
				w = httptest.NewRecorder()
				r.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(200))
				var employees []domain.Employee
				json.Unmarshal(w.Body.Bytes(), &employees)
				Expect(len(employees)).To(Equal(0))
				var resp response
				json.Unmarshal(w.Body.Bytes(), &resp)
			})
		})
	})

	Describe("Get list of employees paginated", func() {
		Context("Get all employees from data store for page=2 & size=2", func() {
			It("Should get list of employees", func() {
				r.Handle("/employees", controllers.ResponseHandler(controller.GetAllEmployees)).Methods("GET")
				req, err := http.NewRequest("GET", "/employees?page=2&size=2", nil)
				Expect(err).NotTo(HaveOccurred())
				w = httptest.NewRecorder()
				r.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(200))
				var employees []domain.Employee
				json.Unmarshal(w.Body.Bytes(), &employees)
				empID := []int{118, 117}
				Expect(len(employees)).To(Equal(0))
				var resp response
				json.Unmarshal(w.Body.Bytes(), &resp)

				respDataSlice := resp.Data.([]interface{})
				for i, data := range respDataSlice {
					emp := data.(map[string]interface{})
					id := (int)(emp["id"].(float64))
					Expect(empID[i]).To(Equal(id))
				}
			})
		})
	})
})

type FakeEmployeeStore struct {
	employeeStore []domain.Employee
}

type response struct {
	Data  interface{}
	Error string
}

func (es *FakeEmployeeStore) GetById(Id int) (domain.Employee, error) {
	for _, emp := range es.employeeStore {
		if emp.ID == Id {
			return emp, nil
		}
	}
	return domain.Employee{}, domain.ErrorIDExists
}

func (es *FakeEmployeeStore) GetAll(page, size int) ([]domain.Employee, error) {

	// Create a slice of keys (employee IDs)
	var keys []int
	for k := range es.employeeStore {
		keys = append(keys, k)
	}

	// Sort the keys based on timestamps (created_at)
	sort.Slice(keys, func(i, j int) bool {
		return es.employeeStore[keys[i]].CreatedAt.Before(es.employeeStore[keys[j]].CreatedAt)
	})

	startIndex := page * size
	endIndex := startIndex + size

	// Final result of employee
	var employees []domain.Employee
	var i int
	for _, id := range keys {
		if i >= startIndex && i < endIndex {
			emp := es.employeeStore[id]
			employees = append(employees, emp)
		}
		i++
	}
	return employees, nil
}

func (es *FakeEmployeeStore) Create(employee domain.Employee) error {
	for _, u := range es.employeeStore {
		if u.ID == employee.ID {
			return domain.ErrorIDExists
		}
	}
	es.employeeStore = append(es.employeeStore, employee)
	return nil
}

func (es *FakeEmployeeStore) Delete(Id int) error {
	tempStore := es
	es = &FakeEmployeeStore{}
	for _, emp := range tempStore.employeeStore {
		if emp.ID != Id {
			es.employeeStore = append(es.employeeStore, emp)
		}
	}
	return nil
}

func (es *FakeEmployeeStore) Update(Id int, employee domain.Employee) error {
	for n, emp := range es.employeeStore {
		if emp.ID == Id {
			es.employeeStore[n] = employee
			return nil
		}
	}
	return nil
}

func newFakeEmployeeStore() *FakeEmployeeStore {
	store := &FakeEmployeeStore{}
	store.Create(domain.Employee{
		ID:        111,
		Name:      "Andrea Le",
		Position:  "Sr. SW Engg",
		Salary:    75000.00,
		CreatedAt: time.Now().UTC(),
	})
	store.Create(domain.Employee{
		ID:        112,
		Name:      "Manoj Yakkla",
		Position:  "Principal SW Engg",
		Salary:    32000.00,
		CreatedAt: time.Now().UTC(),
	})
	store.Create(domain.Employee{
		ID:        113,
		Name:      "Pintu SS",
		Position:  "Sr. SW Engg",
		Salary:    75000.00,
		CreatedAt: time.Now().UTC(),
	})
	store.Create(domain.Employee{
		ID:        119,
		Name:      "Anuj G",
		Position:  "Principal SW Engg",
		Salary:    185000.00,
		CreatedAt: time.Now().UTC(),
	})
	store.Create(domain.Employee{
		ID:        118,
		Name:      "Rajani Singh",
		Position:  "Principal SW Engg",
		Salary:    132000.00,
		CreatedAt: time.Now().UTC(),
	})
	store.Create(domain.Employee{
		ID:        117,
		Name:      "Ashish Shah",
		Position:  "Sr. SW Engg",
		Salary:    35000.00,
		CreatedAt: time.Now().UTC(),
	})
	return store
}
