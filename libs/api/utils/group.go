package api_utils

type GroupCreate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Group struct {
	GroupCreate
	ID int `json:"id"`
}

const GROUP_API_VERSION = "v1"
