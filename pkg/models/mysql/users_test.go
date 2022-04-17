package mysql

import (
	"reflect"
	"testing"
	"time"

	"github.com/mwettste/snippetbox/pkg/models"
)

func TestUserModelGet(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	tests := []struct {
		name          string
		userID        int
		expectedUser  *models.User
		expectedError error
	}{
		{
			name:   "Valid ID",
			userID: 1,
			expectedUser: &models.User{
				ID:      1,
				Name:    "Alice Jones",
				Email:   "alice@example.com",
				Created: time.Date(2018, 12, 23, 17, 25, 22, 0, time.UTC),
				Active:  true,
			},
			expectedError: nil,
		},
		{
			name:          "Zero ID",
			userID:        0,
			expectedUser:  nil,
			expectedError: models.ErrNoRecord,
		},
		{
			name:          "Non-existent ID",
			userID:        2,
			expectedUser:  nil,
			expectedError: models.ErrNoRecord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			m := UserModel{db}

			user, err := m.Get(tt.userID)

			if err != tt.expectedError {
				t.Errorf("expected %v; got %s", tt.expectedError, err)
			}

			if !reflect.DeepEqual(user, tt.expectedUser) {
				t.Errorf("expected %v; got %v", tt.expectedUser, user)
			}
		})
	}
}
