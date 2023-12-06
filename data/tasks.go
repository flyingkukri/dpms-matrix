package data

type Task struct {
	Requirements []Constraint
	Name         string
	NextTask     *Task
}

var (
	DoorTask = Task{
		Requirements: DoorConstraints,
		Name:         "door",
		NextTask:     nil,
	}
	PhoneTask = Task{
		Requirements: PhoneConstraints,
		Name:         "phone",
		NextTask:     &DoorTask,
	}
	CameraTask = Task{
		Requirements: CameraConstraints,
		Name:         "camera",
		NextTask:     &PhoneTask,
	}
	SensorTask = Task{
		Requirements: SensorConstraints,
		Name:         "motionSensor",
		NextTask:     &CameraTask,
	}
)
