package engines

import (
	"fmt"

	"tum.de/huan/dpms-matrix/data"
)

type Door struct {
	*Engine
}

func (d *Door) checkRequirements(reqs []data.Constraint) bool {
	for _, req := range reqs {
		if req.Name != "door" {
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

func NewDoor(process *data.Task, name string, matrixInfo data.MatrixInfo) *Door {

	d := &Door{
		Engine: NewEngine(process, &data.Task{}, name, matrixInfo),
	}
	d.Engine.requirementsCheck = d.checkRequirements
	d.Engine.completeTask = d.doorTask
	return d
}

func (d *Door) doorTask() {
	fmt.Println(red + "The Door has been opened." + reset)
	d.FindNextMachine(d.currentTask.NextTask)
	d.currentTask = &data.Task{}
}
