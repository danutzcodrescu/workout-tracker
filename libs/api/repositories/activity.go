package api_repositories

import (
	"encoding/json"
	api_utils "workout-tracker/libs/api/utils"

	"github.com/jmoiron/sqlx"
)

type ActivityRepository struct {
	DB *sqlx.DB
}

func (r ActivityRepository) Insert(workout api_utils.Workout) error {
	laps_json, _ := json.Marshal(workout.Laps)
	result := r.DB.MustExec("INSERT INTO workouts.activities (date, laps, distance, duration, calories, pace) VALUES ($1, $2, $3, $4, $5, $6)", workout.Date, laps_json, workout.Distance, workout.Time, workout.Calories, workout.Pace)
	_, err := result.RowsAffected()
	return err
}
