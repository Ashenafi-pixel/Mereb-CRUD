package response

import "mereb-crud/pkg/domain"

type PersonResponse struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Hobbies []string `json:"hobbies"`
}

func NewSuccessResponse(message string, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"status":  "success",
		"message": message,
		"data":    data,
	}
}

func MerebErrorResponse(message string) map[string]interface{} {
	return map[string]interface{}{
		"status":  "error",
		"message": message,
	}
}

func ToPersonResponse(person *domain.Person) *PersonResponse {
	return &PersonResponse{
		ID:      person.ID,
		Name:    person.Name,
		Age:     person.Age,
		Hobbies: person.Hobbies,
	}
}

func ToPersonResponseList(persons []*domain.Person) []*PersonResponse {
	var responses []*PersonResponse
	for _, persons := range persons {
		responses = append(responses, ToPersonResponse(persons))
	}
	return responses
}
