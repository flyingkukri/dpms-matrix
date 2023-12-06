package engines

import (
	"errors"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"tum.de/huan/dpms-matrix/data"
)

type MotionSensor struct {
	*Engine
}

func (m *MotionSensor) checkRequirements(reqs []data.Constraint) bool {
	for _, req := range reqs {
		if req.Name != "motionSensor" {
			return false
		}
		if req.Condition != "equals" {
			return false
		}
		if req.Value != "true" {
			return false
		}
	}
	return true
}

func NewMotionSensor(process *data.Task, name string, matrixInfo data.MatrixInfo) *MotionSensor {

	m := &MotionSensor{
		Engine: NewEngine(process, &data.Task{}, name, matrixInfo),
	}
	m.Engine.requirementsCheck = m.checkRequirements
	m.Engine.completeTask = m.motionTask
	return m
}

func (m MotionSensor) Trigger() error {
	if cmp.Equal(m.process, data.Task{}) {
		return errors.New("no process defined")
	}
	m.currentTask = m.process
	m.completeTask()
	return nil
}

func (m *MotionSensor) motionTask() {
	fmt.Println(red + "Motion Sensor completed task." + reset)
	m.FindNextMachine(m.currentTask.NextTask)
	m.currentTask = &data.Task{}
}
