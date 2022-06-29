package api

import (
	"HTTPChiSqlite/db"
	"HTTPChiSqlite/model"
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log"
	"net/http"
	"strconv"
)

const peopleResourceName = "person"

var peopleRepository db.PeopleRepository

func PeopleRouter(r chi.Router) {
	r.Get("/", listPeople)
	r.Put("/", addPerson)
	r.Route("/{id:[0-9]+}", func(r chi.Router) {
		r.Get("/", getPerson)
		r.Patch("/", updatePerson)
		r.Delete("/", deletePerson)
	})

	peopleRepository = db.PeopleRepositoryImpl{}
}

func listPeople(w http.ResponseWriter, r *http.Request) {
	rows, err := peopleRepository.FindAllPeople()
	if err != nil {
		log.Print(err)
		internalError(w, r, err)
		return
	}

	var rl []render.Renderer
	for _, obj := range rows {
		rl = append(rl, obj)
	}

	if len(rl) == 0 {
		rl = []render.Renderer{}
	}

	render.Status(r, http.StatusOK)
	err = render.RenderList(w, r, rl)
	if err != nil {
		log.Print(err)
	}
}

func getPerson(w http.ResponseWriter, r *http.Request) {
	personId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || personId < 0 {
		log.Print(err)
		badRequest(w, r, err)
		return
	}

	person, err := peopleRepository.FindPerson(personId)
	if err != nil {
		log.Print(err)
		resourceNotFound(w, r, peopleResourceName, person.Id)
		return
	}

	render.Status(r, http.StatusOK)
	err = render.Render(w, r, &person)
	if err != nil {
		log.Print(err)
	}
}

func addPerson(w http.ResponseWriter, r *http.Request) {
	person, ok := bindPerson(w, r)
	if !ok {
		return
	}

	id, err := peopleRepository.SavePerson(r.Context(), person)
	if err != nil {
		log.Print(err)
		internalError(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
	err = render.Render(w, r, &AddPersonResponse{Id: id})
	if err != nil {
		log.Print(err)
	}
}

func updatePerson(w http.ResponseWriter, r *http.Request) {
	personId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || personId < 0 {
		log.Print(err)
		badRequest(w, r, err)
		return
	}

	person, ok := bindPerson(w, r)
	if !ok {
		render.Status(r, http.StatusBadRequest)
		return
	}

	person, err = peopleRepository.UpdatePerson(r.Context(), personId, person)
	if err != nil {
		log.Print(err)

		if err == sql.ErrNoRows {
			resourceNotFound(w, r, peopleResourceName, personId)
			return
		}

		switch v := err.(type) {
		case db.NoRowsAffectedError:
			resourceNotFound(w, r, peopleResourceName, personId)
			return
		case db.ParamsNotValidError:
			badRequest(w, r, v)
			return
		default:
			internalError(w, r, err)
			return
		}
	}

	render.Status(r, http.StatusOK)
	err = render.Render(w, r, &person)
	if err != nil {
		log.Print(err)
	}
}

func deletePerson(w http.ResponseWriter, r *http.Request) {
	personId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || personId < 0 {
		log.Print(err)
		badRequest(w, r, err)
		return
	}

	err = peopleRepository.DeletePerson(r.Context(), personId)
	if err != nil {
		log.Print(err)
		internalError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
}

func bindPerson(w http.ResponseWriter, r *http.Request) (model.Person, bool) {
	data := model.Person{}
	if err := render.Bind(r, &data); err != nil {
		badRequest(w, r, err)
		return data, false
	}

	return data, true
}
