package usercase

import (
	"encoding/json"
	"log"
	"net/http"
	"peltoar/v2/model"
	repository "peltoar/v2/repositories"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type EmployeeService struct {
	MongoCollection *mongo.Collection
}

type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

func (svc *EmployeeService) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	res := &Response{}
	defer json.NewEncoder(w).Encode(res)

	var emp model.Employee
	err := json.NewDecoder(r.Body).Decode(&emp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("invalid body", err)
		res.Error = err.Error()
		return
	}

	// assign new employee id
	emp.EmployeeID = uuid.NewString()

	repo := repository.EmployeeRepo{MongoCollection: svc.MongoCollection}

	//insert employee
	insertID, err := repo.InsertEmployee(&emp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("insert error ", err)
		res.Error = err.Error()
		return
	}
	res.Data = emp.EmployeeID
	w.WriteHeader(http.StatusOK)

	log.Println("employee inserted with id ", insertID, emp)
}

func (svc *EmployeeService) GetEmployeeByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	res := &Response{}
	defer json.NewEncoder(w).Encode(res)

	//get employee id
	empID := mux.Vars(r)["id"]
	log.Println("employee id ", empID)

	repo := repository.EmployeeRepo{MongoCollection: svc.MongoCollection}

	emp, err := repo.FindEmployeeByID(empID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("insert error ", err)
		res.Error = err.Error()
		return
	}

	res.Data = emp
	w.WriteHeader(http.StatusOK)
}

func (svc *EmployeeService) GetAllEmployees(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	res := &Response{}
	defer json.NewEncoder(w).Encode(res)

	repo := repository.EmployeeRepo{MongoCollection: svc.MongoCollection}

	emp, err := repo.FindAllEmployees()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("insert error ", err)
		res.Error = err.Error()
		return
	}

	res.Data = emp
	w.WriteHeader(http.StatusOK)
}

func (svc *EmployeeService) UpdateEmployeeByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	res := &Response{}
	defer json.NewEncoder(w).Encode(res)

	//get employee id
	empID := mux.Vars(r)["id"]
	log.Println("employee id", empID)
	if empID == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("invalid employee id")
		res.Error = "invalid employee id"
		return
	}

	var emp model.Employee
	err := json.NewDecoder(r.Body).Decode(&emp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("invalid body", err)
		res.Error = err.Error()
		return
	}

	emp.EmployeeID = empID
	repo := repository.EmployeeRepo{MongoCollection: svc.MongoCollection}
	count, err := repo.UpdateEmployeeByID(empID, &emp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("error: ", err)
		res.Error = err.Error()
		return
	}

	res.Data = count
	w.WriteHeader(http.StatusOK)
}

func (svc *EmployeeService) DeleteByEmployeeID(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	res := &Response{}
	defer json.NewEncoder(w).Encode(res)

	//get employee id
	empID := mux.Vars(r)["id"]
	log.Println("employee id", empID)
	repo := repository.EmployeeRepo{MongoCollection: svc.MongoCollection}
	count, err := repo.DeleteByEmployeeID(empID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("error: ", err)
		res.Error = err.Error()
		return
	}

	res.Data = count
	w.WriteHeader(http.StatusOK)

}

func (svc *EmployeeService) DeleteAllEmployees(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	res := &Response{}
	defer json.NewEncoder(w).Encode(res)

	repo := repository.EmployeeRepo{MongoCollection: svc.MongoCollection}
	count, err := repo.DeleteAllEmployees()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("error: ", err)
		res.Error = err.Error()
		return
	}

	res.Data = count
	w.WriteHeader(http.StatusOK)
}
