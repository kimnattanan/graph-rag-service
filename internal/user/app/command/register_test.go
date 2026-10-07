package command

import (
	"fmt"
	"strings"
	"testing"
)

func TestRegisterGoStringRedactsPassword(t *testing.T) {
	cmd := Register{Email: "ada@example.com", Username: "ada", Password: "super-secret-password"}
	printed := fmt.Sprintf("%#v", cmd)
	if strings.Contains(printed, "super-secret-password") {
		t.Fatalf("password leaked: %s", printed)
	}
}
