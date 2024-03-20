package repository

import (
	"context"
	"fmt"
	"peltoar/v2/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type EmployeeRepo struct {
	MongoCollection *mongo.Collection
}

func (r *EmployeeRepo) InsertEmployee(emp *model.Employee) (interface{}, error) {
	result, err := r.MongoCollection.InsertOne(context.Background(), emp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *EmployeeRepo) FindEmployeeByID(empID string) (*model.Employee, error) {
	var emp model.Employee

	err := r.MongoCollection.
		FindOne(context.Background(), bson.D{{Key: "employee_id", Value: empID}}).
		Decode(&emp)

	if err != nil {
		return nil, err
	}

	return &emp, nil
}

func (r *EmployeeRepo) FindAllEmployees() ([]model.Employee, error) {
	results, err := r.MongoCollection.Find(context.Background(), bson.D{})

	if err != nil {
		return nil, err
	}

	var emps []model.Employee
	err = results.All(context.Background(), &emps)
	if err != nil {
		return nil, fmt.Errorf("results decode error %s", err.Error())
	}

	return emps, nil
}

func (r *EmployeeRepo) UpdateEmployeeByID(empID string, updateEmp *model.Employee) (int64, error) {
	results, err := r.MongoCollection.UpdateOne(context.Background(), bson.D{{Key: "employee_id", Value: empID}},
		bson.D{{Key: "$set", Value: updateEmp}})

	if err != nil {
		return 0, err
	}

	return results.ModifiedCount, nil
}

func (r *EmployeeRepo) DeleteByEmployeeID(empID string) (int64, error) {
	result, err := r.MongoCollection.DeleteOne(context.Background(),
		bson.D{{Key: "employee_id", Value: empID}})

	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}

func (r *EmployeeRepo) DeleteAllEmployees() (int64, error) {
	result, err := r.MongoCollection.DeleteMany(context.Background(), bson.D{})

	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}
