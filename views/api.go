package views

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/danielkraic/knihomol/configuration"
	"github.com/danielkraic/knihomol/controllers"

	"github.com/go-playground/validator/v10"
)

//API json rest api views
type API struct {
	controller *controllers.BooksController
	timeout    time.Duration
	validate   *validator.Validate
}

//NewAPI creates new API
func NewAPI(conf *configuration.Configuration, controller *controllers.BooksController) *API {
	timeoutSec := time.Duration(conf.Timeout) * time.Second

	return &API{
		controller: controller,
		timeout:    timeoutSec,
		validate:   validator.New(),
	}
}

//GetBooks return all books
func (api *API) GetBooks(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()

	books, err := api.controller.GetBooks(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	api.jsonResponse(w, r, books)
}

func (api *API) jsonResponse(w http.ResponseWriter, r *http.Request, data interface{}) {
	response, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("content/type", "application/json")
	_, err = w.Write(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//RefreshBooks refresh status of books
func (api *API) RefreshBooks(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()

	errs := api.controller.Refresh(ctx)
	if errs != nil {
		if len(errs) > 0 {
			http.Error(w, fmt.Sprintf("failed to refresh %d books with error: %s", len(errs), errs[0]), http.StatusInternalServerError)
			return
		}

		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}
}

//AddBook add book to DB
func (api *API) AddBook(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()

	var request manageBookRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = api.validate.Struct(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if strings.HasPrefix(request.BookID, "http") {
		url := request.BookID
		request.BookID = api.controller.GetBookIDFromURL(url)
		if request.BookID == "" {
			http.Error(w, "invalid url", http.StatusInternalServerError)
			return
		}
	}

	err = api.controller.AddBook(ctx, request.BookID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type manageBookRequest struct {
	BookID string `json:"book_id" validate:"required,min=1"`
}

//RemoveBook remove book from DB
func (api *API) RemoveBook(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()

	var request manageBookRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = api.validate.Struct(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = api.controller.RemoveBook(ctx, request.BookID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//RemoveBooks remove books from DB
func (api *API) RemoveBooks(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()

	var request removeBooksRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = api.validate.Struct(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = api.controller.RemoveBooks(ctx, request.BookIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type removeBooksRequest struct {
	BookIDs []string `json:"books_ids" validate:"required,min=1,dive,min=1"`
}
