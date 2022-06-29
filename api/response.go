package api

import (
	"github.com/go-chi/render"
	"log"
	"net/http"
)

type InvalidRequestResponse struct{}

func (ir *InvalidRequestResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type GenericErrorResponse struct {
	Error string `json:"error"`
}

func (ge *GenericErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type ResourceNotFoundResponse struct {
	Name string `json:"name"`
	Id   any    `json:"id"`
}

func (rnf *ResourceNotFoundResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type AddPersonResponse struct {
	Id int64 `json:"id"`
}

func (ap *AddPersonResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func badRequest(w http.ResponseWriter, r *http.Request, err error) {
	errorResponse(w, r, http.StatusBadRequest, err)
}

func internalError(w http.ResponseWriter, r *http.Request, err error) {
	errorResponse(w, r, http.StatusInternalServerError, err)
}

func errorResponse(w http.ResponseWriter, r *http.Request, status int, err error) {
	render.Status(r, status)
	err = render.Render(w, r, &GenericErrorResponse{Error: err.Error()})
	if err != nil {
		log.Print(err)
	}
}

func resourceNotFound(w http.ResponseWriter, r *http.Request, resName string, id int) {
	render.Status(r, http.StatusNotFound)
	err := render.Render(w, r, &ResourceNotFoundResponse{Name: resName, Id: id})
	if err != nil {
		log.Print(err)
	}
}
