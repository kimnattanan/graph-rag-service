package query

import (
	"fmt"
	"strings"
	"testing"
)

func TestLoginGoStringRedactsPassword(t *testing.T) {
	q := Login{Email: "ada@example.com", Password: "super-secret-password"}
	printed := fmt.Sprintf("%#v", q)
	if strings.Contains(printed, "super-secret-password") {
		t.Fatalf("password leaked: %s", printed)
	}
}
