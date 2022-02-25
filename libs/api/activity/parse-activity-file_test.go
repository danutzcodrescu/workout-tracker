package api_activity

import (
	"io/ioutil"
	"reflect"
	"testing"
)

func TestParseActivityFile(t *testing.T) {
	workout, err := ioutil.ReadFile("../../../tools/testFiles/workout.tcx")
	if err != nil {
		t.Errorf("cannot parse workout file file")
	}
	workoutWitoutWatts, err := ioutil.ReadFile("../../../tools/testFiles/workoutWithoutWatts.tcx")
	if err != nil {
		t.Errorf("cannot parse workout without watts file file")
	}
	result, err := ParseActivityFile(workout)
	if err != nil {
		t.Errorf("cannot parse activity")
	}
	expectedWorkout := &Workout{Date: "2022-02-24T13:39:00Z", Laps: []WorkoutLap{
		{
			StartTime:        "2022-02-24T13:39:00Z",
			TotalTimeSeconds: 60,
			DistanceMeters:   244,
			Calories:         16,
			Intensity:        "Active",
			Efforts: []Effort{
				{
					Time:           1,
					DistanceMeters: 4,
					Cadence:        0,
					Watts:          61,
				},
				{
					Time:           2,
					DistanceMeters: 8,
					Cadence:        0,
					Watts:          61,
				},
				{
					Time:           58,
					DistanceMeters: 238,
					Cadence:        22,
					Watts:          201,
				},
			},
		},
		{
			StartTime:        "2022-02-24T13:40:00Z ",
			TotalTimeSeconds: 60,
			DistanceMeters:   172,
			Calories:         9,
			Intensity:        "Resting",
			Efforts: []Effort{
				{
					Time:           2,
					DistanceMeters: 11,
					Cadence:        22,
					Watts:          201,
				},
				{
					Time:           58,
					DistanceMeters: 172,
					Cadence:        22,
					Watts:          63,
				},
			},
		},
	}}
	if reflect.DeepEqual(result, expectedWorkout) {
		t.Errorf("want %+v; got %+v", expectedWorkout, result)
	}
	expectedWorkout.Laps[0].Efforts[0].Watts = 0
	res, err := ParseActivityFile(workoutWitoutWatts)
	if err != nil {
		t.Errorf("cannot parse activity")
	}
	if reflect.DeepEqual(res, expectedWorkout) {
		t.Errorf("want %+v; got %+v", expectedWorkout, result)
	}
}
