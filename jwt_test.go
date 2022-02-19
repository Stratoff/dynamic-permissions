package permissions

import (
	"testing"
)

func TestGetRoleFromPayload(t *testing.T) {
	// TODO: SET MORE CASES
	// Posible tokens
	testCases := map[string]bool{
		// Valid JWT
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjpbImFkbWluIl19.npXlL-uh44DD99rX2ofXYZzcdermJILZ5F9v11LtdiY": true,
	}

	for token, v := range testCases {
		if role, err := GetRoleFromPayload(token); (err == nil) != v {
			t.Fatal("Unexpected behavior", role, err)
		}
	}
}
