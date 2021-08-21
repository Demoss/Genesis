package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"genesis/global"
	"genesis/pkg/resources"
	"github.com/gorilla/schema"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"net/http"

	"regexp"
)

type Server struct {
	config *Config
	logger *logrus.Logger
}

func New(config *Config) *Server {
	return &Server{
		config: config,
		logger: logrus.New(),
	}
}

func (s *Server) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}

	s.logger.Info("Starting server")

	return nil
}

func (s *Server) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}

	s.logger.SetLevel(level)
	return nil
}

type IServer interface {
	AddNewUser(u resources.User) error
	CheckIsUserCreated(user resources.User, users resources.AllUsers) bool
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

func (s Server) CheckIsUserCreated(user resources.User, users resources.AllUsers) bool {
	for i := range users.Users {
		if user.Email == users.Users[i].Email {
			return false
		}
	}
	return true
}

func (s Server) AddNewUser(user resources.User) error {
	var AllUsers resources.AllUsers

	rawDataIn, err := ReadInfoFromFile("users.json")
	err = json.Unmarshal(rawDataIn, &AllUsers)
	if err != nil {
		println(err)
	}

	user = resources.User{user.Email, user.Pass}

	is := s.CheckIsUserCreated(user, AllUsers)

	if is == true {

		AllUsers.Users = append(AllUsers.Users, user)

		WriteInfoToFile("users.json", AllUsers)

		return nil
	} else {
		return errors.New("user is already exist")
	}

}

func valid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

func (s Server) CreateUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("context-type", "application/json")

	keys1, ok1 := r.URL.Query()["email"]
	keys2, ok2 := r.URL.Query()["pass"]
	if !ok1 || len(keys1[0]) == 0 {
		log.Println("Url Email is missing")
		w.Write([]byte("Please check params spelling"))
		return
	}
	if !ok2 || len(keys2[0]) == 0 {
		log.Println("Url Password is missing")
		w.Write([]byte("Please check params spelling"))
		return
	}

	key1 := keys1[0]
	if valid(key1) {
		key2 := keys2[0]

		log.Println("Email " + string(key1) + " and pass " + string(key2))

		var decoder = schema.NewDecoder()

		var user resources.User

		err := decoder.Decode(&user, r.URL.Query())
		if err != nil {
			log.Println(err)
			return
		}

		isCreated := s.AddNewUser(user)

		if isCreated == nil {
			w.Write([]byte("User successfully created"))
		} else {
			log.Println(isCreated)
			w.Write([]byte("User is already exist"))
		}
	} else {
		w.Write([]byte("Incorrect email"))
		log.Println("Incorrect email")
	}
}

func (s Server) AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	var AllUsers resources.AllUsers

	rawDataIn, err := ReadInfoFromFile("users.json")

	err = json.Unmarshal(rawDataIn, &AllUsers)
	if err != nil {
		println(err)
	}
	keys1, ok1 := r.URL.Query()["email"]
	keys2, ok2 := r.URL.Query()["pass"]
	if !ok1 || len(keys1[0]) < 1 {
		log.Println("Url Email 'key' is missing")
	}
	if !ok2 || len(keys2[0]) < 1 {
		log.Println("Url Password 'key' is missing")
	}

	var decoder = schema.NewDecoder()

	var user resources.User

	err = decoder.Decode(&user, r.URL.Query())
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < len(AllUsers.Users); i++ {
		if user.Email != AllUsers.Users[i].Email || user.Pass != AllUsers.Users[i].Pass {
			global.Logged = 0
		} else {
			global.Logged = 1
			break
		}
	}
	if global.Logged == 1 {
		w.Write([]byte("You are logged in"))
	} else {
		w.Write([]byte("User doesn't exist"))
	}
}
