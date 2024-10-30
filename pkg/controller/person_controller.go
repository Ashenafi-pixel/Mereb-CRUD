package controller

import (
	"net/http"

	"mereb-crud/pkg/controller/request"
	"mereb-crud/pkg/controller/response"
	"mereb-crud/pkg/service"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PersonController struct {
	personService service.IPersonService
}

func NewPersonController(personService service.IPersonService) *PersonController {
	return &PersonController{personService}
}

func (controller *PersonController) RegisterPersonRoutes(e *echo.Echo) {
	peopleGroup := e.Group("/person")

	peopleGroup.GET("", controller.GetAllPersons)
	peopleGroup.GET("/:personId", controller.GetPersonById)
	peopleGroup.POST("", controller.AddNewPerson)
	peopleGroup.PUT("/:personId", controller.UpdatePersonById)
	peopleGroup.DELETE("/:personId", controller.DeletePersonById)

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if c.Echo().Renderer != nil {
		}
		c.JSON(http.StatusNotFound, response.NewErrorResponse("Endpoint not found"))
	}
}

func (controller *PersonController) GetAllPersons(c echo.Context) error {
	persons := controller.personService.GetAllPersons()
	return c.JSONPretty(http.StatusOK, response.ToPersonResponseList(persons), "  ")
}

func (controller *PersonController) GetPersonById(c echo.Context) error {
	param := c.Param("personId")
	person, err := controller.personService.GetPersonById(param)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.NewErrorResponse(err.Error()))
	}

	return c.JSONPretty(http.StatusOK, response.ToPersonResponse(person), "  ")
}

func (controller *PersonController) AddNewPerson(c echo.Context) error {
	addPersonRequest := new(request.AddPersonRequest)

	if err := c.Bind(addPersonRequest); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request: unable to bind data"))
	}

	// Generate a new UUID
	newID := uuid.New().String()
	personModel := addPersonRequest.ToModel()
	personModel.ID = newID

	if err := controller.personService.Add(personModel); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, response.NewErrorResponse(err.Error()))
	}

	// Return the created person
	return c.JSON(http.StatusCreated, response.NewSuccessResponse("Person created successfully", personModel))
}

func (controller *PersonController) UpdatePersonById(c echo.Context) error {
	param := c.Param("personId")
	updatePersonRequest := new(request.UpdatePersonRequest)

	if err := c.Bind(updatePersonRequest); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request: unable to bind the provided data"))
	}

	err := controller.personService.UpdatePersonById(updatePersonRequest.ToModel(), param)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, response.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, response.NewSuccessResponse("Person updated successfully", nil))
}

func (controller *PersonController) DeletePersonById(c echo.Context) error {
	personID := c.Param("personId")
	if err := controller.personService.DeleteById(personID); err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, response.NewSuccessResponse("Person deleted successfully", nil))
}
