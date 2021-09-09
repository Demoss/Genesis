package handlers

import (
	"genesis/internal/server"
	"genesis/pkg"
	"genesis/pkg/resources"
	"net/http/httptest"
	"os"
	"testing"
)

func TestConnector_GetBTC(t *testing.T) {
	test := []struct {
		name     string
		response bool
		url      string
		user     resources.User
		s        server.Server
	}{
		{
			name:     "all right",
			response: true,
			url:      "https://api.cryptonator.com/api/ticker/btc-uah",
			user:     resources.User{Email: "ss@gmail.com", Pass: "ss"},
			s:        server.Server{},
		},
		{
			name:     "fail to connect",
			response: false,
			url:      "https://api.cryptonator.com/api/ticker/btc-uah1",
			user:     resources.User{Email: "ss@gmail.com", Pass: "ss"},
			s:        server.Server{},
		},
	}
	for _, tc := range test {
		t.Run(tc.name, func(tt *testing.T) {
			req := httptest.NewRequest("GET", tc.url, nil)

			w := httptest.NewRecorder()

			urls := "http://localhost:8000/user/auth?email=" + tc.user.Email + "&pass=" + tc.user.Pass
			reqt := httptest.NewRequest("GET", urls, nil)
			wr := httptest.NewRecorder()

			err := server.AddNewUser(tc.user)
			if err != nil {
				return
			}

			tc.s.AuthenticateUser(wr, reqt)

			bit := pkg.NewConnector()

			istrue := bit.GetBTC().Success

			GetBTC(w, req)

			if istrue != tc.response {
				tt.Errorf("exp %v, get %v", tc.response, istrue)
			}
		})
	}
	err := os.Remove("users.json")
	if err != nil {
		return
	}
}
