package api_repositories

import (
	"encoding/json"
	"fmt"
	api_utils "workout-tracker/libs/api/utils"

	"github.com/jmoiron/sqlx"
)

type ActivityRepository struct {
	DB *sqlx.DB
}

type GetAllParams struct {
	Limit          int
	Offset         int
	OrderBy        string
	OrderDirection string
}

func (r ActivityRepository) Insert(workout api_utils.Workout) error {
	laps_json, _ := json.Marshal(workout.Laps)
	result := r.DB.MustExec("INSERT INTO workouts.activities (date, laps, distance, duration, calories, pace) VALUES ($1, $2, $3, $4, $5, $6)", workout.Date, laps_json, workout.Distance, workout.Time, workout.Calories, workout.Pace)
	_, err := result.RowsAffected()
	return err
}

func (r ActivityRepository) Get(params GetAllParams) ([]api_utils.WorkoutData, error) {
	activities := []api_utils.WorkoutData{}
	query := fmt.Sprintf("SELECT date, calories, distance, to_char(duration * interval '1 second', 'MI:SS.MS') as time, to_char(pace * interval '1 second', 'MI:SS.MS') as pace FROM workouts.activities ORDER BY %s %s LIMIT $1 OFFSET $2", params.OrderBy, params.OrderDirection)
	err := r.DB.Select(&activities, query, params.Limit, params.Offset)
	return activities, err
}

func (r ActivityRepository) GetByDate(date string) (api_utils.WorkoutWithLaps, error) {
	activity := api_utils.WorkoutWithLaps{}
	err := r.DB.Get(&activity, "SELECT date, calories, distance, to_char(duration * interval '1 second', 'MI:SS.MS') as time, to_char(pace * interval '1 second', 'MI:SS.MS') as pace, laps FROM workouts.activities WHERE date=$1", date)
	return activity, err
}

func (r ActivityRepository) SetGroupId(date string, groupId int) (api_utils.WorkoutData, error) {
	activity := api_utils.WorkoutData{}
	err := r.DB.Get(&activity, "UPDATE workouts.activities SET group_id=$1 WHERE date=$2 RETURNING date, calories, distance, to_char(duration * interval '1 second', 'MI:SS.MS') as time, to_char(pace * interval '1 second', 'MI:SS.MS') as pace", groupId, date)
	return activity, err
}
