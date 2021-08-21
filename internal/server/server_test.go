package server

import (
	"genesis/pkg/resources"
	"testing"
)

func TestServer_CheckIsUserCreated(t *testing.T) {
	var IsUserCreated = func(user resources.User, users resources.AllUsers) bool {
		for i := range users.Users {
			if user.Email == users.Users[i].Email {
				return false
			}
		}
		return true
	}

	test := []struct {
		name      string
		user      resources.User
		users     resources.AllUsers
		mock      func()
		expectErr bool
	}{
		{
			name: "all right",
			user: resources.User{Email: "ss@gmail.com", Pass: "ss"},
			users: resources.AllUsers{Users: []resources.User{
				{
					Email: "sss@gmail.com",
					Pass:  "ss",
				},
				{
					Email: "sdasd@gmail.com",
					Pass:  "ss",
				},
			}},
			expectErr: true,
		},
		{
			name: "fail",
			user: resources.User{Email: "sss@gmail.com", Pass: "ss"},
			users: resources.AllUsers{Users: []resources.User{
				{
					Email: "sss@gmail.com",
					Pass:  "ss",
				},
				{
					Email: "sdasd@gmail.com",
					Pass:  "ss",
				},
			}},
			expectErr: false,
		},
	}
	s := IsUserCreated
	for _, tc := range test {
		result := IsUserCreated(tc.user, tc.users)
		t.Run(tc.name, func(tt *testing.T) {
			if tc.expectErr != result {
				tt.Errorf("expected %v, returned - %v", tc.expectErr, result)
			}
		})
	}
	IsUserCreated = s
}
