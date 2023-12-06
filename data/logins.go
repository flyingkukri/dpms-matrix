package data

import "maunium.net/go/mautrix/id"

type MatrixInfo struct {
	Server    string
	UserID    id.UserID
	AuthToken string
}

var (
	MotionSensorInfo = MatrixInfo{
		Server:    "",
		UserID:    "",
		AuthToken: "",
	}
	CameraInfo = MatrixInfo{
		Server:    "",
		UserID:    "",
		AuthToken: "",
	}
	DoorInfo = MatrixInfo{
		Server:    "",
		UserID:    "",
		AuthToken: "",
	}
	PhoneInfo = MatrixInfo{
		Server:    "",
		UserID:    "",
		AuthToken: "",
	}
)
