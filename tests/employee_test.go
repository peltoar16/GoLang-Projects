package tests

import (
	"context"
	"log"
	"peltoar/v2/model"
	repository "peltoar/v2/repositories"
	"testing"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func mongoClient() *mongo.Client {
	mongoTestClient, err := mongo.Connect(context.Background(),
		options.Client().ApplyURI("mongodb+srv://peltoar-db-user:AustinPelto15!@cluster0.pr0lz.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"))

	if err != nil {
		log.Fatal("error while connecting to mongodb", err)
	}

	log.Println("mongodb successfully connected.")
	err = mongoTestClient.Ping(context.Background(), readpref.Primary())

	if err != nil {
		log.Fatal("ping failed", err)
	}

	log.Println("ping success")
	return mongoTestClient
}

func TestMongoOperations(t *testing.T) {
	mongoTestClient := mongoClient()
	defer mongoTestClient.Disconnect(context.Background())

	// dummy data
	emp1 := uuid.New().String()
	// emp2 := uuid.New().String()

	// connect to colleciton
	col := mongoTestClient.Database("companydb").Collection("employee_test")
	empRepo := repository.EmployeeRepo{MongoCollection: col}

	// Insert Employee 1 data
	t.Run("Insert Employee 1", func(t *testing.T) {
		emp := model.Employee{
			Name:       "Austin Pelto",
			Department: "ComSci",
			EmployeeID: emp1,
		}
		result, err := empRepo.InsertEmployee(&emp)

		if err != nil {
			t.Fatal("insert 1 operation failed", err)
		}

		t.Log("Insert 1 employee successful", result)
	})

	// Get Employee by ID Data
	t.Run("Get Employee 1", func(t *testing.T) {
		result, err := empRepo.FindEmployeeByID(emp1)

		if err != nil {
			t.Fatal("get operation failed", err)
		}

		t.Log("emp 1", result.Name)
	})

	// Get All Employee Data
	t.Run("Get All Employees", func(t *testing.T) {
		results, err := empRepo.FindAllEmployees()

		if err != nil {
			t.Fatal("get operation failed", err)
		}

		t.Log("employees", results)
	})

	// Update Employee 1 data
	t.Run("Update Employee 1 Name", func(t *testing.T) {
		emp := model.Employee{
			Name:       "Tony Stark aka Iron Man",
			Department: "Physics",
			EmployeeID: emp1,
		}

		result, err := empRepo.UpdateEmployeeByID(emp1, &emp)

		if err != nil {
			log.Fatal("update operation failed", err)
		}

		t.Log("update count", result)
	})

	t.Run("Get Employee 1 After update", func(t *testing.T) {
		result, err := empRepo.FindEmployeeByID(emp1)

		if err != nil {
			t.Fatal("get operation failed", err)
		}

		t.Log("emp 1", result.Name)
	})
	// Delete Specific Employee
	t.Run("Delete Employee 1", func(t *testing.T) {
		result, err := empRepo.DeleteByEmployeeID(emp1)
		if err != nil {
			log.Fatal("delete operation failed", err)
		}

		t.Log("delete count", result)
	})

	// Get All Employees Data After Delete
	t.Run("Get All Employees After Delete", func(t *testing.T) {
		results, err := empRepo.FindAllEmployees()

		if err != nil {
			t.Fatal("get operation failed", err)
		}

		t.Log("employees", results)
	})

	// Delete All Employees for Clean Up Purposes
	t.Run("Delete All Employee for Cleanup", func(t *testing.T) {
		result, err := empRepo.DeleteAllEmployees()
		if err != nil {
			log.Fatal("delete operation failed", err)
		}
		t.Log("deleted count", result)
	})
}
