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

// CreateTeam creates new item
func (server *Server) CreateTeam(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	team := models.Team{}
	err = json.Unmarshal(body, &team)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	team.Prepare()
	err = team.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Uncreatorized"))
		return
	}
	if uid != team.CreatorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	teamCreated, err := team.SaveTeam(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, teamCreated.ID))
	responses.JSON(w, http.StatusCreated, teamCreated)
}

// GetTeams gets all teams
func (server *Server) GetTeams(w http.ResponseWriter, r *http.Request) {
	team := models.Team{}
	teams, err := team.FindAllTeams(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, teams)
}

// GetTeam gets by ID
func (server *Server) GetTeam(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	team := models.Team{}
	teamReceived, err := team.FindTeamByID(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, teamReceived)
}

// UpdateTeam updates by ID
func (server *Server) UpdateTeam(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Check if the team id is valid
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	//CHeck if the auth token is valid and  get the user id from it
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Uncreatorized"))
		return
	}

	// Check if the team exist
	team := models.Team{}
	err = server.DB.Debug().Model(models.Team{}).Where("id = ?", pid).Take(&team).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Team not found"))
		return
	}

	// If a user attempt to update a team not belonging to him
	if uid != team.CreatorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Uncreatorized"))
		return
	}
	// Read the data teamed
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	teamUpdate := models.Team{}
	err = json.Unmarshal(body, &teamUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//Also check if the request user id is equal to the one gotten from token
	if uid != teamUpdate.CreatorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Uncreatorized"))
		return
	}

	teamUpdate.Prepare()
	err = teamUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	teamUpdate.ID = team.ID //this is important to tell the model the team id to update, the other update field are set above

	teamUpdated, err := teamUpdate.UpdateATeam(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, teamUpdated)
}

// DeleteTeam deletes by ID
func (server *Server) DeleteTeam(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Is a valid team id given to us?
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Is this user authenticated?
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Uncreatorized"))
		return
	}

	// Check if the team exist
	team := models.Team{}
	err = server.DB.Debug().Model(models.Team{}).Where("id = ?", pid).Take(&team).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Uncreatorized"))
		return
	}

	// Is the authenticated user, the owner of this team?
	if uid != team.CreatorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Uncreatorized"))
		return
	}
	_, err = team.DeleteATeam(server.DB, pid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")
}
