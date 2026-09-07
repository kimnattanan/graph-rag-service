package httperr

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/kimnattanan/graph-rag-service/internal/common/errors"
	"github.com/kimnattanan/graph-rag-service/internal/common/logs"
)

func InternalError(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpResponseWithError(err, slug, w, r, "Internal Server Error", http.StatusInternalServerError)
}

func Unauthorized(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpResponseWithError(err, slug, w, r, "Unauthorized", http.StatusUnauthorized)
}

func BadRequest(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpResponseWithError(err, slug, w, r, "Bad Request", http.StatusBadRequest)
}

func RespondWithSlugError(err error, w http.ResponseWriter, r *http.Request) {
	slugError, ok := err.(errors.SlugError)
	if !ok {
		InternalError("internal-server-error", err, w, r)
		return
	}

	switch slugError.ErrorType() {
	case errors.ErrorTypeAuthorization:
		Unauthorized(slugError.Slug(), slugError, w, r)
	case errors.ErrorTypeIncorrectInput:
		BadRequest(slugError.Slug(), slugError, w, r)
	default:
		InternalError(slugError.Slug(), slugError, w, r)
	}
}

func httpResponseWithError(err error, slug string, w http.ResponseWriter, r *http.Request, logMsg string, status int) {
	logs.GetLogEntry(r).WithError(err).WithField("error-slug", slug).Warn(logMsg)
	resp := ErrorResponse{slug, status}

	if err := render.Render(w, r, resp); err != nil {
		panic(err)
	}
}

type ErrorResponse struct {
	Slug       string `json:"slug"`
	httpStatus int
}

func (e ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(e.httpStatus)
	return nil
}
