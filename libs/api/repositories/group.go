package api_repositories

import (
	api_utils "workout-tracker/libs/api/utils"

	"github.com/jmoiron/sqlx"
)

type GroupRepository struct {
	DB *sqlx.DB
}

func (r GroupRepository) CreateGroup(group api_utils.GroupCreate) (api_utils.Group, error) {
	createdGroup := api_utils.Group{}
	err := r.DB.Get(&createdGroup, "INSERT INTO workouts.groups (name, description) VALUES ($1, $2) RETURNING id, name, description", group.Name, group.Description)
	return createdGroup, err
}

func (r GroupRepository) GetGroups() ([]api_utils.Group, error) {
	groups := []api_utils.Group{}
	err := r.DB.Select(&groups, "SELECT id, name, description FROM workouts.groups")
	return groups, err
}

func (r GroupRepository) GetGroup(id int) (api_utils.Group, error) {
	group := api_utils.Group{}
	err := r.DB.Get(&group, "SELECT id, name, description FROM workouts.groups WHERE id=$1", id)
	return group, err
}

func (r GroupRepository) UpdateGroup(id int, description string) (api_utils.Group, error) {
	updatedGroup := api_utils.Group{}
	err := r.DB.Get(&updatedGroup, "UPDATE workouts.groups SET description=$1 WHERE id=$2 RETURNING id, name, description", description, id)
	return updatedGroup, err
}

func (r GroupRepository) DeleteGroup(id int) error {
	_, err := r.DB.Exec("DELETE FROM workouts.groups WHERE id=$1", id)
	return err
}
