package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/denisecase/go-hunt-sql/api/auth"
	"github.com/denisecase/go-hunt-sql/api/models"
	"github.com/denisecase/go-hunt-sql/api/responses"
	"github.com/denisecase/go-hunt-sql/api/util/formaterror"
	"github.com/gorilla/mux"
)

// CreateUser creates a new user (same as user/register)
func (server *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("Email")
	pwd := r.FormValue("Password")
	pwd2 := r.FormValue("Password")
	fmt.Println("user_controller: CreateUser (Register)", email, pwd, pwd2)
	redirectTarget := "/user/register"

	if email == "" || pwd == "" || pwd2 == "" {
		responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("Please complete all fields"))
		return
	}

	if len(pwd) < 8 || len(pwd) > 20 {
		responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("Passwords must be 8-20 characters"))
		return
	}

	if pwd != pwd2 {
		responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("Passwords must match"))
		return
	}

	user := models.User{}
	user.Email = email
	user.Password = pwd
	user.Prepare()

	err := user.Validate("register")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	} else {
		_, err = user.SaveUser(server.DB)
		if err != nil {
			formattedError := formaterror.FormatError(err.Error())
			responses.ERROR(w, http.StatusInternalServerError, formattedError)
		} else {
			redirectTarget = "/"
		}
	}
	http.Redirect(w, r, redirectTarget, 302)
}

// GetUsers gets all users
func (server *Server) GetUsers(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	users, err := user.FindAllUsers(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, users)
}

// GetUser gets a user by ID
func (server *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	user := models.User{}
	userGotten, err := user.FindUserByID(server.DB, uint32(uid))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, userGotten)
}

// UpdateUser updates a user by ID
func (server *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if tokenID != uint32(uid) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	user.Prepare()
	err = user.Validate("update")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	updatedUser, err := user.UpdateAUser(server.DB, uint32(uid))
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, updatedUser)
}

// DeleteUser deletes a user by ID
func (server *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := models.User{}
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if tokenID != 0 && tokenID != uint32(uid) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	_, err = user.DeleteAUser(server.DB, uint32(uid))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
	responses.JSON(w, http.StatusNoContent, "")
}
