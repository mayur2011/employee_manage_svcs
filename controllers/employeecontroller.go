package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"employee_manage_svcs/domain"

	"github.com/gorilla/mux"
)

type EmployeeController struct {
	// Dependencies and States
	Store domain.EmployeeStore
}

func (ec EmployeeController) PostEmployee(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var employee domain.Employee
	// decode the incoming request
	err := json.NewDecoder(r.Body).Decode(&employee)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("unable to decode JSON request body")
	}
	err = ec.Store.Create(employee)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	return nil, http.StatusCreated, nil
}

func (ec EmployeeController) GetAllEmployees(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	reqMap := r.URL.Query()

	page := 0
	size := 10
	var err error
	// Accessing the "page" value
	pageNum, pageExists := reqMap["page"]
	if pageExists {
		//log.Printf("Page value: %v\n", pageNum)
		page, err = strconv.Atoi(pageNum[0])
		if err != nil {
			return nil, http.StatusBadRequest, fmt.Errorf("unable to process pageNum value from the request")
		}
	}
	// Accessing the "size" value
	pageSize, sizeExists := reqMap["size"]
	if sizeExists {
		size, err = strconv.Atoi(pageSize[0])
		if err != nil {
			return nil, http.StatusBadRequest, fmt.Errorf("unable to process pageSize value from the request")
		}
	}

	employees, err := ec.Store.GetAll(page, size)
	if err != nil {
		return nil, http.StatusOK, fmt.Errorf("failed to fetch all employees")
	}
	return employees, http.StatusOK, nil
}

func (ec EmployeeController) GetEmployeeById(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	vars := mux.Vars(r)
	id_str := vars["id"]
	id, e := strconv.Atoi(id_str)
	if e != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("unable to process ID value from the request")
	}
	employee, err := ec.Store.GetById(id)
	if err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("no employee info")
	}
	return employee, http.StatusOK, nil
}

func (ec EmployeeController) UpdateEmployee(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var employee domain.Employee
	vars := mux.Vars(r)
	id_str := vars["id"]
	id, e := strconv.Atoi(id_str)
	if e != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("unable to process ID value from the request")
	}
	err := json.NewDecoder(r.Body).Decode(&employee)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("unable to decode JSON request body")
	}
	err = ec.Store.Update(id, employee)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return nil, http.StatusAccepted, err
}

func (ec EmployeeController) DeleteEmployee(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	vars := mux.Vars(r)
	id_str := vars["id"]
	id, e := strconv.Atoi(id_str)
	if e != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("unable to process ID value from the request")
	}
	err := ec.Store.Delete(id)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return nil, http.StatusAccepted, err
}
