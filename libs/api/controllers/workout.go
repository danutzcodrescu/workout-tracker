package api_controllers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	api_repositories "workout-tracker/libs/api/repositories"
	api_utils "workout-tracker/libs/api/utils"
)

type LinkActivityWithGroupBody struct {
	GroupId int `json:"groupId"`
}

const dateFormat = "2006-01-02T15:04:05Z"

const file_size_in_mb = 3

// 3MB file upload size
const max_upload_size = 1024 * 1024 * file_size_in_mb

func parseActivityFile(fileBytes []byte) (api_utils.Workout, error) {
	var err error
	var workout api_utils.TrainingCenterDatabase
	xml.Unmarshal(fileBytes, &workout)
	var laps []api_utils.WorkoutLap
	var totalTime = 0
	var distance = 0
	var calories = 0
	for _, lap := range workout.Activities.Activity.Laps {
		workoutLap := api_utils.WorkoutLap{Intensity: lap.Intensity, StartTime: lap.StartTime}
		workoutLap.Calories, err = strconv.Atoi(lap.Calories)
		calories += workoutLap.Calories
		if err != nil {
			log.Println(err, "calories parsing")
			return api_utils.Workout{}, err
		}
		workoutLap.DistanceMeters, err = strconv.Atoi(lap.DistanceMeters)
		distance += workoutLap.DistanceMeters
		if err != nil {
			log.Println(err, "distance rowed")
			return api_utils.Workout{}, err
		}
		workoutLap.TotalTimeSeconds, err = strconv.Atoi(lap.TotalTimeSeconds)
		totalTime += workoutLap.TotalTimeSeconds
		if err != nil {
			log.Println(err, "rowed time")
			return api_utils.Workout{}, err
		}
		for _, effort := range lap.Track.Trackpoints {
			workoutEffort := api_utils.Effort{}
			startTime, err := time.Parse(dateFormat, strings.TrimRight(lap.StartTime, " "))
			if err != nil {
				log.Println(err)
				return api_utils.Workout{}, err
			}
			effortTime, err := time.Parse(dateFormat, effort.Time)
			if err != nil {
				log.Println(err)
				return api_utils.Workout{}, err
			}
			workoutEffort.Cadence, err = strconv.Atoi(effort.Cadence)
			if err != nil {
				log.Println(err, "cadence")
				return api_utils.Workout{}, err
			}
			workoutEffort.DistanceMeters, err = strconv.Atoi(effort.DistanceMeters)
			if err != nil {
				log.Println(err, "distance per effort")
				return api_utils.Workout{}, err
			}
			if effort.Extensions.TPX.Watts == "" {
				workoutEffort.Watts = 0
			} else {
				workoutEffort.Watts, err = strconv.Atoi(effort.Extensions.TPX.Watts)
				if err != nil {
					log.Println(err, "wats", effort.Time)
					return api_utils.Workout{}, err
				}
			}
			workoutEffort.Time = int(effortTime.Sub(startTime).Seconds())
			workoutLap.Efforts = append(workoutLap.Efforts, workoutEffort)
		}
		laps = append(laps, workoutLap)
	}
	return api_utils.Workout{
		Date:     workout.Activities.Activity.ID,
		Laps:     laps,
		Distance: distance,
		Time:     totalTime,
		Calories: calories,
		Pace:     float32(totalTime) / float32(distance) * 500,
	}, nil
}

func UploadActivityController(application *Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, max_upload_size)
		if err := r.ParseMultipartForm(max_upload_size); err != nil {
			clientError(w, err, fmt.Sprintln("The uploaded file is too big. Please choose an file that's less than", file_size_in_mb, "MB in size."))(application)
			return
		}

		file, _, err := r.FormFile("activity")
		if err != nil {
			clientError(w, err, fmt.Sprintln("The form does not contain any file under activity form field"))(application)
			return
		}
		defer file.Close()
		fileBytes, _ := io.ReadAll(file)
		activityWorkout, err := parseActivityFile(fileBytes)
		if err != nil {
			serverError(w, err, "Error parsing workout")(application)
		}
		err = application.Repositories.Activity.Insert(activityWorkout)
		if err != nil {
			serverError(w, err, "Error inserting record")(application)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(activityWorkout)
	}
}

func RetrieveActivities(application *Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()

		limitParam := params.Get("limit")
		if limitParam == "" {
			limitParam = "10"
		}
		limit, err := strconv.Atoi(limitParam)
		if err != nil {
			serverError(w, err, "Error parsing limit")(application)
		}

		offsetParam := params.Get("offset")
		if offsetParam == "" {
			offsetParam = "0"
		}
		offset, err := strconv.Atoi(offsetParam)
		if err != nil {
			serverError(w, err, "Error parsing offset")(application)
		}

		orderBy := params.Get("orderBy")
		if orderBy == "" {
			orderBy = "date"
		}

		orderDirection := params.Get("orderDirection")
		if orderDirection == "" {
			orderDirection = "DESC"
		}

		activities, err := application.Repositories.Activity.Get(api_repositories.GetAllParams{
			Limit:          limit,
			Offset:         offset,
			OrderBy:        orderBy,
			OrderDirection: orderDirection,
		})
		if err != nil {
			serverError(w, err, "Error retrieving activities")(application)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(activities)
	}
}

func RetrieveActivity(application *Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		activityDate := strings.TrimPrefix(r.URL.Path, fmt.Sprintf("/%s/activities/", api_utils.ACTIVITY_API_VERSION))

		activity, err := application.Repositories.Activity.GetByDate(activityDate)
		if err != nil && err.Error() == "sql: no rows in result set" {
			notFoundError(w, err, "Activity not found")(application)
			return
		}

		if err != nil {
			serverError(w, err, "Error retrieving activity")(application)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(activity)
	}
}

func LinkActivityWithGroup(application *Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body = LinkActivityWithGroupBody{}
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			serverError(w, err, "Error parsing group id")(application)
			return
		}
		activityDate := strings.TrimPrefix(r.URL.Path, fmt.Sprintf("/%s/activities/", api_utils.ACTIVITY_API_VERSION))
		activity, err := application.Repositories.Activity.SetGroupId(activityDate, body.GroupId)
		if err != nil {
			serverError(w, err, "Error linking activity with group")(application)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(activity)
	}
}
