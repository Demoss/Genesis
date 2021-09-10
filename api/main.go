package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"genesis/global"
	"genesis/pkg"
	"genesis/pkg/api"
	"genesis/pkg/resources"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io/ioutil"
	"log"
	"net"
)

type Server struct {
}

func ReadInfoFromFile(path string) ([]byte, error) {
	rawDataIn, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println(err)
	}
	return rawDataIn, err
}

func WriteInfoToFile(path string, users resources.AllUsers) {
	rawDataOut, err := json.MarshalIndent(&users, "", "  ")
	if err != nil {
		log.Fatal("JSON marshaling failed:", err)
	}

	err = ioutil.WriteFile(path, rawDataOut, 0644)
	if err != nil {
		log.Fatal("Cannot write updated settings file:", err)
	}
}

func CheckIsUserCreated(user resources.User, users resources.AllUsers) bool {
	for i := range users.Users {
		if user.Email == users.Users[i].Email {
			return false
		}
	}
	return true
}

func AddNewUser(user resources.User) error {
	var AllUsers resources.AllUsers

	rawDataIn, err := ReadInfoFromFile("users.json")
	err = json.Unmarshal(rawDataIn, &AllUsers)
	if err != nil {
		println(err)
	}

	is := CheckIsUserCreated(user, AllUsers)

	if is == true {

		AllUsers.Users = append(AllUsers.Users, user)

		WriteInfoToFile("users.json", AllUsers)

		return nil
	} else {
		return errors.New("user is already exist")
	}

}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer()
	api.RegisterCreateUserServer(srv, &Server{})
	reflection.Register(srv)

	if e := srv.Serve(listener); e != nil {
		panic(err)
	}
}

func (s *Server) Create(ctx context.Context, req *api.UserRequest) (*api.UserResponse, error) {
	email, pass := req.GetEmail(), req.GetPass()
	u := resources.User{Email: email, Pass: pass}

	err := AddNewUser(u)
	if err != nil {
		return nil, errors.New("user is already exist")
	}

	return &api.UserResponse{Response: "user successfully created"}, nil
}

func (s *Server) Auth(ctx context.Context, req *api.UserRequest) (*api.UserResponse, error) {
	var AllUsers resources.AllUsers

	rawDataIn, err := ReadInfoFromFile("users.json")

	err = json.Unmarshal(rawDataIn, &AllUsers)
	if err != nil {
		log.Fatal(err)
	}
	email, pass := req.GetEmail(), req.GetPass()
	if email == "" || pass == "" {
		return &api.UserResponse{Response: "Please check params spelling"}, nil
	}
	log.Printf("email from site %v and pass %v", email, pass)
	for i := 0; i < len(AllUsers.Users); i++ {
		if email == AllUsers.Users[i].Email && pass == AllUsers.Users[i].Pass {
			global.Logged = 1
			return &api.UserResponse{Response: "u ar logged in"}, nil

		} else {
			global.Logged = 0
		}
	}
	return &api.UserResponse{Response: "u ar not logged in"}, nil
}

func (s *Server) GetBTC(ctx context.Context, req *api.URL) (*api.UserResponse, error) {
	if req.GetX() != "btc" {
		return &api.UserResponse{Response: "page not found"}, nil
	}
	con := pkg.NewConnector()
	if global.Logged == 1 {
		resp := con.GetBTC()
		return &api.UserResponse{Response: fmt.Sprintf(
			"1 %v = %v %v",
			resp.Ticker.Base,
			resp.Ticker.Price,
			resp.Ticker.Target,
		)}, nil
	}
	return &api.UserResponse{Response: "you are not logged in"}, nil
}
