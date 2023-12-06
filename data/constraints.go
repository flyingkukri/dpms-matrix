package data

import (
	"errors"

	"github.com/mitchellh/mapstructure"
)

type Constraint struct {
	Name      string
	Condition string
	Value     string
}

func GetConstraints(reqs any) ([]Constraint, error) {
	constraintList := make([]Constraint, 0)
	if reqsMap, ok := reqs.([]any); ok {
		for _, req := range reqsMap {
			result := Constraint{}
			if err := mapstructure.Decode(req, &result); err != nil {
				return nil, err
			}
			constraintList = append(constraintList, result)
		}
		return constraintList, nil
	} else {
		return nil, errors.New("requirements could not be converted")
	}
}

var (
	SensorConstraints = []Constraint{
		{
			Name:      "motionSensor",
			Condition: "equals",
			Value:     "true",
		},
	}

	CameraConstraints = []Constraint{
		{
			Name:      "camera",
			Condition: "equals",
			Value:     "true",
		},
	}

	PhoneConstraints = []Constraint{
		{
			Name:      "phone",
			Condition: "equals",
			Value:     "true",
		},
	}

	DoorConstraints = []Constraint{
		{
			Name:      "door",
			Condition: "equals",
			Value:     "true",
		},
	}
)
