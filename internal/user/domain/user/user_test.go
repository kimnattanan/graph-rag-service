package user

import (
	"strings"
	"testing"

	"github.com/kimnattanan/graph-rag-service/internal/common/auth"
)

func TestNewUserStoresRoleUserAndChecksPassword(t *testing.T) {
	account, err := NewUser("user-1", " Ada@Example.com ", " Ada ", "correct-horse")
	if err != nil {
		t.Fatal(err)
	}
	if account.Email() != "ada@example.com" {
		t.Fatalf("email = %s", account.Email())
	}
	if account.Username() != "Ada" {
		t.Fatalf("username = %s", account.Username())
	}
	if account.Role() != RoleUser {
		t.Fatalf("role = %s", account.Role())
	}
	if account.PasswordHash() == "correct-horse" || !strings.HasPrefix(account.PasswordHash(), "$2") {
		t.Fatal("password was not hashed")
	}
	if err := account.CheckPassword("correct-horse"); err != nil {
		t.Fatal(err)
	}
	if err := account.CheckPassword("wrong-password"); err != ErrInvalidCredentials {
		t.Fatalf("err = %v", err)
	}

	permissions := account.Permissions()
	if len(permissions) != 1 || permissions[0] != auth.PermissionConversationAsk {
		t.Fatalf("permissions = %#v", permissions)
	}
}

func TestNewAdminIncludesKnowledgePermission(t *testing.T) {
	account, err := NewAdmin("admin-1", "admin@example.com", "admin", "correct-horse")
	if err != nil {
		t.Fatal(err)
	}
	if account.Role() != RoleAdmin {
		t.Fatalf("role = %s", account.Role())
	}
	if !auth.HasPermission(auth.User{Permissions: account.Permissions()}, auth.PermissionKnowledgeWrite) {
		t.Fatalf("permissions = %#v", account.Permissions())
	}
}

func TestNewUserRejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		username string
		password string
	}{
		{name: "email", email: "not-an-email", username: "ada", password: "correct-horse"},
		{name: "username", email: "ada@example.com", username: " ", password: "correct-horse"},
		{name: "password", email: "ada@example.com", username: "ada", password: "short"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := NewUser("user-1", test.email, test.username, test.password); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}
