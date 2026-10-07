package auth

import (
	"testing"
	"time"
)

func TestIssueAndParse(t *testing.T) {
	principal := User{
		ID:          "user-1",
		Email:       "ada@example.com",
		Username:    "ada",
		Role:        RoleUser,
		Permissions: []string{PermissionConversationAsk},
	}

	token, expiresIn, err := Issue("secret", principal, "session-1", time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	if expiresIn != int((time.Hour).Seconds()) {
		t.Fatalf("expiresIn = %d", expiresIn)
	}

	got, err := Parse("secret", token)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != principal.ID || got.Email != principal.Email || got.Username != principal.Username {
		t.Fatalf("principal = %#v", got)
	}
	if got.Role != RoleUser || got.SessionID != "session-1" {
		t.Fatalf("claims = %#v", got)
	}
	if !HasPermission(got, PermissionConversationAsk) || HasPermission(got, PermissionKnowledgeWrite) {
		t.Fatalf("permissions = %#v", got.Permissions)
	}
}

func TestParseRejectsWrongSecret(t *testing.T) {
	token, _, err := Issue("secret", User{ID: "user-1", Role: RoleUser}, "session-1", time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Parse("other", token); err == nil {
		t.Fatal("expected error")
	}
}

func TestPermissionsForRole(t *testing.T) {
	admin, err := PermissionsForRole(RoleAdmin)
	if err != nil {
		t.Fatal(err)
	}
	if len(admin) != 2 {
		t.Fatalf("admin permissions = %#v", admin)
	}

	if _, err := PermissionsForRole("owner"); err == nil {
		t.Fatal("expected error")
	}
}
