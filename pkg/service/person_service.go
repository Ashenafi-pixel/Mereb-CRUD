package service

import (
	"errors"
	"mereb-crud/pkg/domain"
	"mereb-crud/pkg/repository"
	"mereb-crud/pkg/service/dto"
)

type IPersonService interface {
	Add(personCreate *dto.PersonCreate) error
	GetAllPersons() []*domain.Person
	GetPersonById(personId string) (*domain.Person, error)
	UpdatePersonById(updatedPerson *dto.PersonCreate, personId string) error
	DeleteById(personId string) error
}

type PersonService struct {
	personRepository repository.IPersonRepository
}

func NewPersonService(personRepository repository.IPersonRepository) IPersonService {
	return &PersonService{personRepository}
}

func (service *PersonService) Add(personCreate *dto.PersonCreate) error {
	err := validatePersonCreate(personCreate)
	if err != nil {
		return err
	}

	person := personCreateToPerson(personCreate)
	return service.personRepository.AddPerson(person)
}

func (service *PersonService) GetAllPersons() []*domain.Person {
	return service.personRepository.GetAllPersons()
}

func (service *PersonService) GetPersonById(personId string) (*domain.Person, error) {
	return service.personRepository.GetPersonById(personId)
}

func (service *PersonService) UpdatePersonById(updatedPerson *dto.PersonCreate, personId string) error {
	err := service.personRepository.CheckPersonExistence(personId)
	if err != nil {
		return err
	}

	err = validatePersonCreate(updatedPerson)
	if err != nil {
		return err
	}

	person := personCreateToPerson(updatedPerson)
	return service.personRepository.UpdatePersonById(person, personId)
}

func (service *PersonService) DeleteById(personId string) error {
	err := service.personRepository.CheckPersonExistence(personId)
	if err != nil {
		return err
	}

	return service.personRepository.DeletePersonById(personId)
}

func validatePersonCreate(personCreate *dto.PersonCreate) error {
	if personCreate.Name == "" {
		return errors.New("name can't be empty")
	}
	if personCreate.Age < 0 {
		return errors.New("age can't be less than zero")
	}
	if len(personCreate.Hobbies) == 0 { // Check if the hobbies slice is empty
		return errors.New("Hobbies can't be empty")
	}
	return nil
}

func personCreateToPerson(personCreate *dto.PersonCreate) *domain.Person {
	return &domain.Person{
		Name:    personCreate.Name,
		Age:     personCreate.Age,
		Hobbies: personCreate.Hobbies,
	}
}
