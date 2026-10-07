package ports

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/kimnattanan/graph-rag-service/internal/common/auth"
	commonerrors "github.com/kimnattanan/graph-rag-service/internal/common/errors"
	"github.com/kimnattanan/graph-rag-service/internal/common/server/httperr"
	"github.com/kimnattanan/graph-rag-service/internal/user/app"
	"github.com/kimnattanan/graph-rag-service/internal/user/app/command"
	"github.com/kimnattanan/graph-rag-service/internal/user/app/query"
)

type HttpServer struct {
	app app.Application
}

func NewHttpServer(application app.Application) HttpServer {
	return HttpServer{app: application}
}

func (h HttpServer) Register(w http.ResponseWriter, r *http.Request) {
	var body RegisterRequest
	if err := render.Decode(r, &body); err != nil {
		httperr.BadRequest("invalid-request", err, w, r)
		return
	}

	err := h.app.Commands.Register.Handle(r.Context(), command.Register{
		Email:    string(body.Email),
		Username: body.Username,
		Password: body.Password,
	})
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h HttpServer) Login(w http.ResponseWriter, r *http.Request) {
	var body LoginRequest
	if err := render.Decode(r, &body); err != nil {
		httperr.BadRequest("invalid-request", err, w, r)
		return
	}

	result, err := h.app.Queries.Login.Handle(r.Context(), query.Login{
		Email:    string(body.Email),
		Password: body.Password,
	})
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	userResponse, err := userToResponse(result.User)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	render.Respond(w, r, AuthResult{
		AccessToken: result.AccessToken,
		ExpiresIn:   result.ExpiresIn,
		TokenType:   result.TokenType,
		User:        userResponse,
	})
}

func (h HttpServer) Logout(w http.ResponseWriter, r *http.Request) {
	token, err := auth.BearerToken(r)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	if err := h.app.Commands.Logout.Handle(r.Context(), command.Logout{Token: token}); err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h HttpServer) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	current, err := h.authenticate(r)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	view, err := h.app.Queries.GetCurrentUser.Handle(r.Context(), query.GetCurrentUser{UserID: current.ID})
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	response, err := userToResponse(view)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	render.Respond(w, r, response)
}

func (h HttpServer) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	current, err := h.authenticate(r)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	if err := h.app.Commands.DeleteAccount.Handle(r.Context(), command.DeleteAccount{UserID: current.ID}); err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h HttpServer) authenticate(r *http.Request) (auth.User, error) {
	token, err := auth.BearerToken(r)
	if err != nil {
		return auth.User{}, err
	}
	return h.app.Queries.Authenticate.Handle(r.Context(), query.Authenticate{Token: token})
}

func userToResponse(view query.User) (User, error) {
	id, err := uuid.Parse(view.ID)
	if err != nil {
		return User{}, commonerrors.NewIncorrectInputError("invalid user id", "invalid-user-id")
	}

	role := UserRole(view.Role)
	if !role.Valid() {
		return User{}, commonerrors.NewIncorrectInputError("invalid role", "invalid-role")
	}

	permissions := make([]UserPermissions, len(view.Permissions))
	for i, permission := range view.Permissions {
		parsed := UserPermissions(permission)
		if !parsed.Valid() {
			return User{}, commonerrors.NewIncorrectInputError("invalid permission", "invalid-permission")
		}
		permissions[i] = parsed
	}

	return User{
		Id:          id,
		Email:       openapi_types.Email(view.Email),
		Username:    view.Username,
		Role:        role,
		Permissions: permissions,
		CreatedAt:   view.CreatedAt,
		UpdatedAt:   view.UpdatedAt,
	}, nil
}
