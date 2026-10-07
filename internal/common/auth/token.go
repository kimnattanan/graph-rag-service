package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	commonerrors "github.com/kimnattanan/graph-rag-service/internal/common/errors"
)

const (
	RoleAdmin = "admin"
	RoleUser  = "user"

	PermissionKnowledgeWrite  = "knowledge:write"
	PermissionConversationAsk = "conversation:ask"

	TokenTypeBearer = "Bearer"
)

// User is the principal carried by an access token.
// Other services authorize from Role and Permissions.
type User struct {
	ID          string
	Email       string
	Username    string
	Role        string
	Permissions []string
	SessionID   string
}

type Claims struct {
	Email       string   `json:"email"`
	Username    string   `json:"username"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
	SessionID   string   `json:"sid"`
	jwt.RegisteredClaims
}

func PermissionsForRole(role string) ([]string, error) {
	switch role {
	case RoleAdmin:
		return []string{PermissionKnowledgeWrite, PermissionConversationAsk}, nil
	case RoleUser:
		return []string{PermissionConversationAsk}, nil
	default:
		return nil, fmt.Errorf("unknown role %q", role)
	}
}

func HasPermission(user User, permission string) bool {
	for _, candidate := range user.Permissions {
		if candidate == permission {
			return true
		}
	}
	return false
}

func Issue(secret string, user User, sessionID string, ttl time.Duration) (string, int, error) {
	if secret == "" {
		return "", 0, fmt.Errorf("empty jwt secret")
	}
	if sessionID == "" {
		return "", 0, fmt.Errorf("empty session id")
	}
	if ttl <= 0 {
		return "", 0, fmt.Errorf("non-positive token ttl")
	}

	permissions := append([]string(nil), user.Permissions...)
	now := time.Now().UTC()
	expiresAt := now.Add(ttl)
	claims := Claims{
		Email:       user.Email,
		Username:    user.Username,
		Role:        user.Role,
		Permissions: permissions,
		SessionID:   sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			ID:        sessionID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		return "", 0, err
	}
	return signed, int(ttl.Seconds()), nil
}

func Parse(secret string, token string) (User, error) {
	if secret == "" || token == "" {
		return User{}, invalidToken()
	}

	claims := &Claims{}
	parsed, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil || !parsed.Valid {
		return User{}, invalidToken()
	}
	if claims.Subject == "" || claims.SessionID == "" || claims.Role == "" {
		return User{}, invalidToken()
	}

	return User{
		ID:          claims.Subject,
		Email:       claims.Email,
		Username:    claims.Username,
		Role:        claims.Role,
		Permissions: append([]string(nil), claims.Permissions...),
		SessionID:   claims.SessionID,
	}, nil
}

func invalidToken() error {
	return commonerrors.NewAuthorizationError("invalid token", "invalid-token")
}
