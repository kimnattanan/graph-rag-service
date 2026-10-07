package httperr

import (
	stderrors "errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/kimnattanan/graph-rag-service/internal/common/errors"
	"github.com/kimnattanan/graph-rag-service/internal/common/logs"
)

func InternalError(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpResponseWithError(err, slug, w, r, "Internal Server Error", http.StatusInternalServerError)
}

func Unauthorized(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpResponseWithError(err, slug, w, r, userMessage(err, "Unauthorized"), http.StatusUnauthorized)
}

func BadRequest(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpResponseWithError(err, slug, w, r, userMessage(err, "Bad Request"), http.StatusBadRequest)
}

func NotFound(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpResponseWithError(err, slug, w, r, userMessage(err, "Not Found"), http.StatusNotFound)
}

func Conflict(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpResponseWithError(err, slug, w, r, userMessage(err, "Conflict"), http.StatusConflict)
}

func RespondWithSlugError(err error, w http.ResponseWriter, r *http.Request) {
	var slugError errors.SlugError
	if !stderrors.As(err, &slugError) {
		InternalError("internal-server-error", err, w, r)
		return
	}

	switch slugError.ErrorType() {
	case errors.ErrorTypeAuthorization:
		Unauthorized(slugError.Slug(), slugError, w, r)
	case errors.ErrorTypeIncorrectInput:
		BadRequest(slugError.Slug(), slugError, w, r)
	case errors.ErrorTypeNotFound:
		NotFound(slugError.Slug(), slugError, w, r)
	case errors.ErrorTypeConflict:
		Conflict(slugError.Slug(), slugError, w, r)
	default:
		InternalError(slugError.Slug(), slugError, w, r)
	}
}

// userMessage returns the message of errors that were written to be shown to the user
// (the ones implementing errors.SlugError), and a generic fallback for all the others.
func userMessage(err error, fallback string) string {
	var slugError errors.SlugError
	if stderrors.As(err, &slugError) {
		return slugError.Error()
	}

	return fallback
}

func httpResponseWithError(err error, slug string, w http.ResponseWriter, r *http.Request, message string, status int) {
	logs.GetLogEntry(r).WithError(err).WithField("error-slug", slug).Warn(message)
	resp := ErrorResponse{slug, message, status}

	if err := render.Render(w, r, resp); err != nil {
		panic(err)
	}
}

type ErrorResponse struct {
	Slug       string `json:"slug"`
	Message    string `json:"message"`
	httpStatus int
}

func (e ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(e.httpStatus)
	return nil
}
