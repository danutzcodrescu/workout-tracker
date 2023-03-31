package api_controllers

import (
	"encoding/xml"
	"log"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "2006-01-02T15:04:05Z"

type TrainingCenterDatabase struct {
	XMLName    xml.Name   `xml:"TrainingCenterDatabase"`
	Activities Activities `xml:"Activities"`
}

type Activities struct {
	Activity Activity
}

type Activity struct {
	Sport string `xml:"Sport,attr"`
	ID    string `xml:"Id"`
	Laps  []Lap  `xml:"Lap"`
}

type Lap struct {
	StartTime        string `xml:"StartTime,attr"`
	TotalTimeSeconds string `xml:"TotalTimeSeconds"`
	DistanceMeters   string `xml:"DistanceMeters"`
	MaximumSpeed     string `xml:"MaximumSpeed"`
	Calories         string `xml:"Calories"`
	Intensity        string `xml:"Intensity"`
	TriggerMethod    string `xml:"TriggerMethod"`
	Track            Track  `xml:"Track"`
}

type Track struct {
	Trackpoints []TrackPoint `xml:"Trackpoint"`
}

type TrackPoint struct {
	Text           string `xml:",chardata"`
	Time           string `xml:"Time"`
	DistanceMeters string `xml:"DistanceMeters"`
	Cadence        string `xml:"Cadence"`
	SensorState    string `xml:"SensorState"`
	Extensions     struct {
		TPX struct {
			Watts string `xml:"Watts"`
		} `xml:"TPX"`
	} `xml:"Extensions"`
}

type Workout struct {
	Date string       `json:"date"`
	Laps []WorkoutLap `json:"laps"`
}

type WorkoutLap struct {
	StartTime        string `json:"startTime"`
	TotalTimeSeconds int    `json:"totalTimeSeconds"`
	DistanceMeters   int    `json:"distanceMeters"`
	// maximumSpeed     float32
	Calories  int      `json:"calories"`
	Intensity string   `json:"intensity"`
	Efforts   []Effort `json:"efforts"`
}

type Effort struct {
	Time           int `json:"time"`
	DistanceMeters int `json:"distanceMeters"`
	Cadence        int `json:"cadence"`
	Watts          int `json:"watts"`
}

func ParseActivityFile(fileBytes []byte) (Workout, error) {
	var err error
	var workout TrainingCenterDatabase
	xml.Unmarshal(fileBytes, &workout)
	var laps []WorkoutLap
	for _, lap := range workout.Activities.Activity.Laps {
		workoutLap := WorkoutLap{Intensity: lap.Intensity, StartTime: lap.StartTime}
		workoutLap.Calories, err = strconv.Atoi(lap.Calories)
		if err != nil {
			// TODO: add it to app logging
			log.Println(err, "calories parsing")
			return Workout{}, err
		}
		workoutLap.DistanceMeters, err = strconv.Atoi(lap.DistanceMeters)
		if err != nil {
			// TODO: add it to app logging
			log.Println(err, "distance rowed")
			return Workout{}, err
		}
		workoutLap.TotalTimeSeconds, err = strconv.Atoi(lap.TotalTimeSeconds)
		if err != nil {
			// TODO: add it to app logging
			log.Println(err, "rowed time")
			return Workout{}, err
		}
		for _, effort := range lap.Track.Trackpoints {
			workoutEffort := Effort{}
			startTime, err := time.Parse(dateFormat, strings.TrimRight(lap.StartTime, " "))
			if err != nil {
				// TODO: add it to app logging
				log.Println(err)
				return Workout{}, err
			}
			effortTime, err := time.Parse(dateFormat, effort.Time)
			if err != nil {
				// TODO: add it to app logging
				log.Println(err)
				return Workout{}, err
			}
			workoutEffort.Cadence, err = strconv.Atoi(effort.Cadence)
			if err != nil {
				// TODO: add it to app logging
				log.Println(err, "cadence")
				return Workout{}, err
			}
			workoutEffort.DistanceMeters, err = strconv.Atoi(effort.DistanceMeters)
			if err != nil {
				// TODO: add it to app logging
				log.Println(err, "distance per effort")
				return Workout{}, err
			}
			if effort.Extensions.TPX.Watts == "" {
				workoutEffort.Watts = 0
			} else {
				workoutEffort.Watts, err = strconv.Atoi(effort.Extensions.TPX.Watts)
				if err != nil {
					// TODO: add it to app logging
					log.Println(err, "wats", effort.Time)
					return Workout{}, err
				}
			}
			workoutEffort.Time = int(effortTime.Sub(startTime).Seconds())
			workoutLap.Efforts = append(workoutLap.Efforts, workoutEffort)
		}
		laps = append(laps, workoutLap)
	}
	return Workout{
		Date: workout.Activities.Activity.ID,
		Laps: laps,
	}, nil
}
