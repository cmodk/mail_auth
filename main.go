package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const (
	COST = 14
)

var (
	userdb       = flag.String("userdb", "/user.db", "File with user database")
	http_listen  = flag.String("http-listen", "", "Listen address for http")
	http_port    = flag.Int("http-port", 9000, "Listen port for http")
	add_user     = flag.Bool("adduser", false, "Add new user to database")
	new_password = flag.String("password", "", "Password for new user")
	new_user     = flag.String("username", "", "New username")

	log = logrus.New()

	users   []User
	usermap map[string]User
)

type User struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func main() {
	flag.Parse()

	file, err := os.Open(*userdb)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err := yaml.NewDecoder(file).Decode(&users); err != nil {
		panic(err)
	}

	usermap = make(map[string]User)
	for _, u := range users {
		usermap[u.Username] = u
	}

	if *add_user {
		if err := AddNewUser(); err != nil {
			log.Fatal(err)
		}
		return
	}

	mux := http.ServeMux{}

	mux.HandleFunc("/auth", authHandler)

	s := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", *http_listen, *http_port),
		Handler:        &mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())

}

func invalidLoginOrPassword(w http.ResponseWriter) {
	w.Header().Add("Auth Status", "Invalid login or password")
	fmt.Fprintf(w, "")
}

func success(w http.ResponseWriter) {
	w.Header().Add("Auth Status", "OK")
	fmt.Fprintf(w, "")
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("Auth-User")
	password := r.Header.Get("Auth-Pass")

	user, ok := usermap[username]
	if !ok {
		log.Errorf("User %s not found\n", username)
		invalidLoginOrPassword(w)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		invalidLoginOrPassword(w)
		return
	}

	success(w)

}
