package repository

import (
	"context"
	"fmt"
	"mereb-crud/pkg/domain"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/gommon/log"
)

type IPersonRepository interface {
	GetAllPersons() []*domain.Person
	GetPersonById(personId string) (*domain.Person, error)
	AddPerson(person *domain.Person) error
	CheckPersonExistence(personId string) error
	UpdatePersonById(updatedPerson *domain.Person, personId string) error
	DeletePersonById(personId string) error
}

type PersonRepository struct {
	dbPool *pgxpool.Pool
}

func NewPersonRepository(dbPool *pgxpool.Pool) IPersonRepository {
	return &PersonRepository{dbPool}
}

func (repository *PersonRepository) GetAllPersons() []*domain.Person {
	ctx := context.Background()
	personRows, err := repository.dbPool.Query(ctx, "SELECT * FROM persons")
	if err != nil {
		log.Errorf("error while getting all persons: %v", err)
		return nil // Consider returning an error here instead
	}
	defer personRows.Close()
	return extractPersonsFromRows(personRows)
}

func (repository *PersonRepository) GetPersonById(personId string) (*domain.Person, error) {
	ctx := context.Background()
	var person domain.Person
	var hobbiesString string

	if err := repository.dbPool.QueryRow(ctx, "SELECT id, name, age, hobbies FROM persons WHERE id = $1", personId).Scan(&person.ID, &person.Name, &person.Age, &hobbiesString); err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("person not found with id: %s", personId) // Friendly error message
		}
		log.Errorf("error while getting person by ID: %v", err)
		return nil, fmt.Errorf("an error occurred while retrieving the person, please try again later")
	}

	person.Hobbies = strings.Split(hobbiesString, ", ")
	return &person, nil
}

func (repository *PersonRepository) AddPerson(person *domain.Person) error {
	ctx := context.Background()

	person.ID = uuid.New().String()
	hobbiesString := strings.Join(person.Hobbies, ", ")
	insertStatement := "INSERT INTO persons (id, name, age, hobbies) VALUES ($1, $2, $3, $4)"
	if _, err := repository.dbPool.Exec(ctx, insertStatement, person.ID, person.Name, person.Age, hobbiesString); err != nil {
		log.Errorf("error while adding a new person: %v", err)
		return fmt.Errorf("could not add the person, please check the input data")
	}

	return nil
}

func (repository *PersonRepository) CheckPersonExistence(personId string) error {
	ctx := context.Background()

	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM persons WHERE id = $1)"
	if err := repository.dbPool.QueryRow(ctx, query, personId).Scan(&exists); err != nil {
		log.Errorf("error while checking person existence: %v", err)
		return fmt.Errorf("could not check if the person exists, please try again later")
	}

	if !exists {
		return fmt.Errorf("person with id %s does not exist", personId)
	}

	return nil
}

func (repository *PersonRepository) UpdatePersonById(updatedPerson *domain.Person, personId string) error {
	ctx := context.Background()
	hobbiesString := strings.Join(updatedPerson.Hobbies, ", ")
	updateStatement := "UPDATE persons SET name = $1, age = $2, hobbies = $3 WHERE id = $4"
	if _, err := repository.dbPool.Exec(ctx, updateStatement, updatedPerson.Name, updatedPerson.Age, hobbiesString, personId); err != nil {
		log.Errorf("error while updating person: %v", err)
		return fmt.Errorf("could not update the person, please check the input data")
	}

	return nil
}

func (repository *PersonRepository) DeletePersonById(personId string) error {
	ctx := context.Background()

	deleteExec, err := repository.dbPool.Exec(ctx, "DELETE FROM persons WHERE id = $1", personId)
	if err != nil {
		log.Errorf("error while deleting person: %v", err)
		return fmt.Errorf("could not delete the person, please try again later")
	}

	if deleteExec.RowsAffected() == 0 {
		return fmt.Errorf("no person found with id %s", personId) // More informative
	}

	log.Info(fmt.Sprintf("%v rows affected", deleteExec.RowsAffected()))
	return nil
}

func extractPersonsFromRows(personRows pgx.Rows) []*domain.Person {
	var persons []*domain.Person

	for personRows.Next() {
		person := &domain.Person{}
		var hobbiesString string
		if err := personRows.Scan(&person.ID, &person.Name, &person.Age, &hobbiesString); err != nil {
			log.Errorf("error while scanning person row: %v", err)
			continue
		}
		person.Hobbies = strings.Split(hobbiesString, ", ")
		persons = append(persons, person)
	}

	return persons
}
