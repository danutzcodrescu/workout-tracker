package api_repositories

import (
	"encoding/json"
	api_controllers "workout-tracker/libs/api/controllers"

	"github.com/jmoiron/sqlx"
)

type ActivityRepository struct {
	DB *sqlx.DB
}

func (r ActivityRepository) Insert(workout api_controllers.Workout) error {
	laps_json, _ := json.Marshal(workout.Laps)
	result := r.DB.MustExec("INSERT INTO workouts.activities (date, laps) VALUES ($1, $2)", workout.Date, laps_json)
	_, err := result.RowsAffected()
	return err
}
