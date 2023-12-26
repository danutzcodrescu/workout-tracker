package api_controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	api_utils "workout-tracker/libs/api/utils"
)

type UpdateGroupBody struct {
	Description string `json:"description"`
}

func CreateGroup(application *Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var group api_utils.GroupCreate
		err := json.NewDecoder(r.Body).Decode(&group)
		if err != nil {
			serverError(w, err, "Error parsing group")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		createdGroup, err := application.Repositories.Group.CreateGroup(group)
		if err != nil {
			serverError(w, err, "Error creating group")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdGroup)
	}
}

func GetGroups(application *Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groups, err := application.Repositories.Group.GetGroups()
		if err != nil {
			serverError(w, err, "Error retrieving groups")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(groups)
	}
}

func GetGroup(application *Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		group := strings.TrimPrefix(r.URL.Path, fmt.Sprintf("/%s/groups/", api_utils.GROUP_API_VERSION))
		groupId, err := strconv.Atoi(group)
		if err != nil {
			serverError(w, err, "Error parsing group")
		}
		groups, err := application.Repositories.Group.GetGroup(groupId)
		if err != nil {
			serverError(w, err, "Error retrieving groups")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(groups)
	}
}

func UpdateGroup(application *Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var groupData = UpdateGroupBody{}
		err := json.NewDecoder(r.Body).Decode(&groupData)
		if err != nil {
			serverError(w, err, "Error parsing description")
		}
		group := strings.TrimPrefix(r.URL.Path, fmt.Sprintf("/%s/groups/", api_utils.GROUP_API_VERSION))
		groupId, err := strconv.Atoi(group)
		if err != nil {
			serverError(w, err, "Error parsing group")
		}
		groups, err := application.Repositories.Group.UpdateGroup(groupId, groupData.Description)
		if err != nil {
			serverError(w, err, "Error retrieving groups")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(groups)
	}
}
