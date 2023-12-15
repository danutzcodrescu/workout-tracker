package api_utils

import "encoding/xml"

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
	Date     string       `json:"date"`
	Laps     []WorkoutLap `json:"laps"`
	Distance int
	Time     int
	Calories int
	Pace     float32
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
