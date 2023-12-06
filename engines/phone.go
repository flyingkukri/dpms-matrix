package engines

import (
	"fmt"

	"tum.de/huan/dpms-matrix/data"
)

type Phone struct {
	*Engine
}

func (p *Phone) checkRequirements(reqs []data.Constraint) bool {
	for _, req := range reqs {
		if req.Name != "phone" {
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

func NewPhone(process *data.Task, name string, matrixInfo data.MatrixInfo) *Phone {

	p := &Phone{
		Engine: NewEngine(process, &data.Task{}, name, matrixInfo),
	}
	p.Engine.requirementsCheck = p.checkRequirements
	p.Engine.completeTask = p.phoneTask
	return p
}

func (p *Phone) phoneTask() {
	fmt.Println(red + "The phone has decided to open the door." + reset)
	p.FindNextMachine(p.currentTask.NextTask)
	p.currentTask = &data.Task{}
}
