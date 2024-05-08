package router

import (
	"employee_manage_svcs/controllers"
	"employee_manage_svcs/mapstore"

	"github.com/gorilla/mux"
)

// Setting Employee routes
func SetEmployeeRoutes(router *mux.Router) *mux.Router {
	employeeStore := mapstore.NewMapStore()
	employeeController := controllers.EmployeeController{Store: employeeStore}
	router.Handle("/employee", controllers.ResponseHandler(employeeController.PostEmployee)).Methods("POST")
	router.Handle("/employees", controllers.ResponseHandler(employeeController.GetAllEmployees)).Methods("GET")
	router.Handle("/employees/{id}", controllers.ResponseHandler(employeeController.GetEmployeeById)).Methods("GET")
	router.Handle("/employees/{id}", controllers.ResponseHandler(employeeController.UpdateEmployee)).Methods("PUT")
	router.Handle("/employees/{id}", controllers.ResponseHandler(employeeController.DeleteEmployee)).Methods("DELETE")
	return router
}

// InitRoutes registers all employee routes for the application.
func InitRoutes() *mux.Router {
	router := mux.NewRouter()
	router = SetEmployeeRoutes(router)
	return router
}
