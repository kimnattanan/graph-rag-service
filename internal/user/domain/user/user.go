package user

import (
	"net/mail"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"

	"github.com/kimnattanan/graph-rag-service/internal/common/auth"
	commonerrors "github.com/kimnattanan/graph-rag-service/internal/common/errors"
)

const (
	minPasswordLength = 8
	maxPasswordLength = 72
	maxUsernameLength = 64
)

var (
	ErrEmptyUserID        = commonerrors.NewIncorrectInputError("empty user id", "empty-user-id")
	ErrInvalidEmail       = commonerrors.NewIncorrectInputError("invalid email", "invalid-email")
	ErrInvalidUsername    = commonerrors.NewIncorrectInputError("invalid username", "invalid-username")
	ErrInvalidPassword    = commonerrors.NewIncorrectInputError("password must be 8 to 72 characters", "invalid-password")
	ErrInvalidRole        = commonerrors.NewIncorrectInputError("invalid role", "invalid-role")
	ErrEmailAlreadyExists = commonerrors.NewConflictError("email already registered", "email-already-exists")
	ErrInvalidCredentials = commonerrors.NewAuthorizationError("invalid email or password", "invalid-credentials")
	ErrNotFound           = commonerrors.NewNotFoundError("user not found", "user-not-found")
	ErrSessionNotFound    = commonerrors.NewNotFoundError("session not found", "session-not-found")
)

type Role string

const (
	RoleUser  Role = auth.RoleUser
	RoleAdmin Role = auth.RoleAdmin
)

type User struct {
	id           string
	email        string
	username     string
	passwordHash string
	role         Role
	createdAt    time.Time
	updatedAt    time.Time
}

func NewUser(id string, email string, username string, password string) (*User, error) {
	return newUser(id, email, username, password, RoleUser)
}

func NewAdmin(id string, email string, username string, password string) (*User, error) {
	return newUser(id, email, username, password, RoleAdmin)
}

func newUser(id string, email string, username string, password string, role Role) (*User, error) {
	if id == "" {
		return nil, ErrEmptyUserID
	}
	normalizedEmail, err := NormalizeEmail(email)
	if err != nil {
		return nil, err
	}
	normalizedUsername, err := normalizeUsername(username)
	if err != nil {
		return nil, err
	}
	if err := validatePassword(password); err != nil {
		return nil, err
	}
	if _, err := parseRole(string(role)); err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	return &User{
		id:           id,
		email:        normalizedEmail,
		username:     normalizedUsername,
		passwordHash: string(hash),
		role:         role,
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

// UnmarshalUserFromDatabase unmarshals a user from the database.
//
// It should be used only for unmarshalling from the database.
// You can't use UnmarshalUserFromDatabase as a constructor - it may put the domain into an invalid state.
func UnmarshalUserFromDatabase(
	id string,
	email string,
	username string,
	passwordHash string,
	role string,
	createdAt time.Time,
	updatedAt time.Time,
) (*User, error) {
	if id == "" {
		return nil, ErrEmptyUserID
	}
	parsedRole, err := parseRole(role)
	if err != nil {
		return nil, err
	}
	normalizedEmail, err := NormalizeEmail(email)
	if err != nil {
		return nil, err
	}
	normalizedUsername, err := normalizeUsername(username)
	if err != nil {
		return nil, err
	}
	if passwordHash == "" {
		return nil, ErrInvalidPassword
	}

	return &User{
		id:           id,
		email:        normalizedEmail,
		username:     normalizedUsername,
		passwordHash: passwordHash,
		role:         parsedRole,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}, nil
}

func NormalizeEmail(email string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	parsed, err := mail.ParseAddress(email)
	if err != nil || parsed.Address != email {
		return "", ErrInvalidEmail
	}
	return email, nil
}

func normalizeUsername(username string) (string, error) {
	username = strings.TrimSpace(username)
	length := utf8.RuneCountInString(username)
	if length < 1 || length > maxUsernameLength {
		return "", ErrInvalidUsername
	}
	return username, nil
}

func validatePassword(password string) error {
	if utf8.RuneCountInString(password) < minPasswordLength || len(password) > maxPasswordLength {
		return ErrInvalidPassword
	}
	return nil
}

func parseRole(role string) (Role, error) {
	switch Role(role) {
	case RoleUser, RoleAdmin:
		return Role(role), nil
	default:
		return "", ErrInvalidRole
	}
}

func (u *User) CheckPassword(password string) error {
	if bcrypt.CompareHashAndPassword([]byte(u.passwordHash), []byte(password)) != nil {
		return ErrInvalidCredentials
	}
	return nil
}

func (u *User) ID() string {
	return u.id
}

func (u *User) Email() string {
	return u.email
}

func (u *User) Username() string {
	return u.username
}

func (u *User) PasswordHash() string {
	return u.passwordHash
}

func (u *User) Role() Role {
	return u.role
}

func (u *User) Permissions() []string {
	permissions, err := auth.PermissionsForRole(string(u.role))
	if err != nil {
		return nil
	}
	return permissions
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}
