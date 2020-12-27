package router

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/denisecase/go-hunt-sql/api/auth"
	"github.com/denisecase/go-hunt-sql/api/models"
	"github.com/gorilla/securecookie"
)

// LoginPageData defines the data needed this page
type LoginPageData struct {
	Email    string
	Password string
	Errors   []string
}

// RegisterPageData defines the data needed this page
type RegisterPageData struct {
	Email     string
	Password  string
	Password2 string
	Errors    []string
}

// ShowLogin handles request to show login page
func (server *Server) ShowLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Show Login")
	data := LoginPageData{
		Email:    "dcase@nwmissouri.edu",
		Password: "password",
	}
	templ.ExecuteTemplate(w, "login", data)
}

// ShowLogout signs user out and returns to home
func (server *Server) ShowLogout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Show Logout")
	clearSession(w)
	http.Redirect(w, r, "/", 302)
}

// ShowRegister handles request to show register page
func (server *Server) ShowRegister(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Show Register")
	data := RegisterPageData{
		Email:     "dcase@nwmissouri.edu",
		Password:  "password",
		Password2: "password",
	}
	templ.ExecuteTemplate(w, "register", data)
}

// PostLogin attempts to sign in user
func (server *Server) PostLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Post Login")
	result := &LoginPageData{
		Email:    r.PostFormValue("Email"),
		Password: r.PostFormValue("Password"),
	}
	if result.Validate() == false {
		fmt.Println("Validation Errors: ", result.Errors)
	//	templ.ExecuteTemplate(w, "login", result) // try again
		http.Redirect(w, r, "/user/login", 302) // success
	}
	http.Redirect(w, r, "/", 302) // success
}

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func setSession(userName string, w http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	}
}

func clearSession(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}

func (server *Server) signIn(email, password string) (string, error) {
	var err error
	user := models.User{}
	err = server.DB.Model(models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return "", err
	}
	err = models.VerifyPassword(user.Password, password)
	if err != nil {
		return "", err
	}
	return auth.CreateToken(user.ID)
}

var rxEmail = regexp.MustCompile(".+@.+\\..+")

// Validate returns true if valid, otherwise adds Errors
func (input *LoginPageData) Validate() bool {
	input.Errors = nil

	match := rxEmail.Match([]byte(input.Email))
	if match == false {
		input.Errors = append(input.Errors, "Please enter a valid email address")
	}

	if strings.TrimSpace(input.Password) == "" {
		input.Errors = append(input.Errors, "Please enter a password")
	}

	return len(input.Errors) == 0
}

// Validate returns true if valid, otherwise adds Errors
func (input *RegisterPageData) Validate() bool {
	input.Errors = nil

	match := rxEmail.Match([]byte(input.Email))
	if match == false {
		input.Errors = append(input.Errors, "Please enter a valid email address")
	}

	if strings.TrimSpace(input.Password) == "" {
		input.Errors = append(input.Errors, "Please enter a password")
	}

	if strings.TrimSpace(input.Password2) == "" {
		input.Errors = append(input.Errors, "Please enter a confirmation password")
	}

	if input.Password != input.Password2 {
		input.Errors = append(input.Errors, "Passwords do not match")
	}

	if len(input.Password) < 8 || len(input.Password) > 20 {
		input.Errors = append(input.Errors, "Passwords must be 8-20 characters long")
	}

	return len(input.Errors) == 0
}

func render(w http.ResponseWriter, filename string, data interface{}) {
	tmpl, err := template.ParseFiles(filename)
	if err != nil {
		log.Println(err)
		http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Println(err)
		http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
	}
}

