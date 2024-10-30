package request

import "mereb-crud/pkg/service/dto"

type AddPersonRequest struct {
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Hobbies []string `json:"hobbies"`
}

func (request *AddPersonRequest) ToModel() *dto.PersonCreate {
	return &dto.PersonCreate{
		Name:    request.Name,
		Age:     request.Age,
		Hobbies: request.Hobbies,
	}
}

type UpdatePersonRequest struct {
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Hobbies []string `json:"hobbies"`
}

func (request *UpdatePersonRequest) ToModel() *dto.PersonCreate {
	return &dto.PersonCreate{
		Name:    request.Name,
		Age:     request.Age,
		Hobbies: request.Hobbies,
	}
}
