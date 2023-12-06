package engines

import (
	"fmt"

	"tum.de/huan/dpms-matrix/data"
)

type Camera struct {
	*Engine
}

func (c *Camera) checkRequirements(reqs []data.Constraint) bool {
	for _, req := range reqs {
		if req.Name != "camera" {
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

func NewCamera(process *data.Task, name string, matrixInfo data.MatrixInfo) *Camera {
	c := &Camera{
		Engine: NewEngine(process, &data.Task{}, name, matrixInfo),
	}
	c.requirementsCheck = c.checkRequirements
	c.completeTask = c.cameraTask
	return c
}

func (c *Camera) cameraTask() {
	fmt.Println(red + "Camera completed task" + reset)
	c.FindNextMachine(c.currentTask.NextTask)
	c.currentTask = &data.Task{}
}
