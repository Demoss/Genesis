package server

import (
	"errors"
	"genesis/pkg/resources"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestServer_CheckIsUserCreated(t *testing.T) {
	test := []struct {
		name      string
		user      resources.User
		users     resources.AllUsers
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

	for _, tc := range test {
		result := CheckIsUserCreated(tc.user, tc.users)
		t.Run(tc.name, func(tt *testing.T) {
			if tc.expectErr != result {
				tt.Errorf("expected %v, returned - %v", tc.expectErr, result)
			}
		})
	}
}



func TestAddNewUser(t *testing.T) {
	tests := []struct{
		name string
		user resources.User
		users resources.AllUsers
		err error
	}{
		{
			name: "all right",
			user: resources.User{Email: "test@gmail.com", Pass: "ss"},
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
			err: nil,
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
			err: errors.New("user is already exist"),
		},
	}
	for _,tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			
			WriteInfoToFile("users.json", tc.users)

			erro := AddNewUser(tc.user)
			if !reflect.DeepEqual(tc.err, erro) {
				tt.Errorf("expected %v, get - %v", tc.err, erro)
			}
		})
		err := os.Remove("users.json")
		if err != nil {
			return
		}
	}

}

func TestCreateUser(t *testing.T) {
	tests := []struct{
		name string
		email string
		pass string
		answer string
	}{
		{
			name: "all right",
			email: "ss@gmail.com",
			pass: "ss",
			answer: "User successfully created",
		},
		{
			name: "user is already exist",
			email: "ss@gmail.com",
			pass: "ss",
			answer: "User is already exist",
		},
		{
			name: "missing email",
			email: "",
			pass: "ss",
			answer: "Please check params spelling",
		},
		{
			name: "missing pass",
			email: "ss@gmail.com",
			pass: "",
			answer: "Please check params spelling",
		},
		{
			name: "not valid email",
			email: ".com",
			pass: "ss",
			answer: "Incorrect email",
		},
	}
	for _,tc := range tests{
		t.Run(tc.name, func(tt *testing.T) {
			url := "http://localhost:8000/user/create?email=" + tc.email + "&pass=" + tc.pass
			req := httptest.NewRequest("GET",url,nil)
			w := httptest.NewRecorder()

			CreateUser(w,req)

			resp := w.Result()
			body,_ := ioutil.ReadAll(resp.Body)
			bodyStr := string(body)
			if bodyStr != tc.answer{
				tt.Errorf("expected %v, get - %v",tc.answer,bodyStr)
			}
		})
	}
	err := os.Remove("users.json")
	if err != nil {
		return
	}
}

func TestAuthenticateUser(t *testing.T) {
	tests := []struct{
		name string
		email string
		pass string
		answer string
	}{

		{
			name: "all right",
			email: "ss@gmail.com",
			pass: "ss",
			answer: "You are logged in",
		},
		{
			name: "fail",
			email: "sss@gmail.com",
			pass: "ss",
			answer: "User doesn't exist",
		},
	}
	for _,tc := range tests{
		t.Run(tc.name, func(tt *testing.T) {
			url := "http://localhost:8000/user/auth?email=" + tc.email + "&pass=" + tc.pass
			req := httptest.NewRequest("GET",url,nil)
			w := httptest.NewRecorder()

			users := resources.AllUsers{Users: []resources.User{
				{
					Email: "ss@gmail.com",
					Pass:  "ss",
				}}}
			WriteInfoToFile("users.json",users)

			AuthenticateUser(w,req)

			resp := w.Result()
			body,_ := ioutil.ReadAll(resp.Body)
			bodyStr := string(body)
			if bodyStr != tc.answer{
				tt.Errorf("expected %v, get - %v",tc.answer,bodyStr)
			}
			err := os.Remove("users.json")
			if err != nil {
				return
			}
		})
	}
}