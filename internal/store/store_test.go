package store

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"os"
	"testing"

	"github.com/githubuseradmin/nova-dashboard/internal/models"
)

// testStore connects to the database named by NOVA_TEST_DATABASE_URL and applies
// migrations. The whole test is skipped when that variable is unset, so `go test`
// stays green without a database (CI sets it to a Postgres service).
func testStore(t *testing.T) *Store {
	t.Helper()
	dsn := os.Getenv("NOVA_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("NOVA_TEST_DATABASE_URL not set; skipping DB integration test")
	}
	st, err := New(context.Background(), dsn)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	if err := st.Migrate(context.Background()); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	t.Cleanup(st.Close)
	return st
}

func randEmail() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return "test-" + hex.EncodeToString(b) + "@example.com"
}

func TestUserRoundTrip(t *testing.T) {
	st := testStore(t)
	ctx := context.Background()
	email := randEmail()

	created, err := st.CreateUser(ctx, email, "Test User", models.RoleUser, "hash")
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	t.Cleanup(func() {
		_, _ = st.pool.Exec(context.Background(), `DELETE FROM users WHERE id = $1`, created.ID)
	})

	got, err := st.GetUserByEmail(ctx, email)
	if err != nil {
		t.Fatalf("GetUserByEmail: %v", err)
	}
	if got.ID != created.ID || got.Email != email || got.Role != models.RoleUser {
		t.Fatalf("round-trip mismatch: %+v", got)
	}

	updated, err := st.UpdateUserRole(ctx, created.ID, models.RoleAdmin)
	if err != nil {
		t.Fatalf("UpdateUserRole: %v", err)
	}
	if !updated.IsAdmin() {
		t.Fatal("role was not updated to admin")
	}
}

func TestGetUserByEmailNotFound(t *testing.T) {
	st := testStore(t)
	if _, err := st.GetUserByEmail(context.Background(), randEmail()); err != ErrNotFound {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}
