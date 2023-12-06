package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"tum.de/huan/dpms-matrix/data"
	"tum.de/huan/dpms-matrix/engines"
)

func main() {
	//testAPI()
	//	testSending()
	engineList := startMachines()
	sensor, ok := engineList[0].(*engines.MotionSensor)
	if !ok {
		fmt.Printf("Error: %s\n", errors.New("not a motion sensor"))
		return
	}
	time.Sleep(time.Second) // so that the machines can start up
	log.Println("Printed previous stuff")
	sensor.Trigger()
	time.Sleep(60 * time.Second) // so that we wait for the process to finish

	for _, engine := range engineList {
		engine.Stop()
	}
}

func startMachines() []engines.EngineInterface {
	process := data.SensorTask

	camera := engines.NewCamera(&data.Task{}, "camera", data.CameraInfo)
	motionSensor := engines.NewMotionSensor(&process, "motionSensor", data.MotionSensorInfo)
	phone := engines.NewPhone(&data.Task{}, "phone", data.PhoneInfo)
	door := engines.NewDoor(&data.Task{}, "door", data.DoorInfo)
	return []engines.EngineInterface{motionSensor,
		camera, phone, door}
}

func testSending() {
	clientMat, err := mautrix.NewClient("https://matrix.org", "@lukashuan:matrix.org", "syt_bHVrYXNodWFu_cCIsNOYxvmTBHznkzFFa_3pqWPS")
	if err != nil {
		panic(err)
	}

	if _, err := clientMat.SendText("!SxURpLPLXhLsLHWqAP:tum.de", "sending test"); err != nil {
		panic(err)
	}
}

func testAPI() {
	// Create clients
	clientTum, err := mautrix.NewClient("https://matrix.tum.de", "@ge59wox:tum.de", "syt_Z2U1OXdveA_ElKmGkzsBORQKtQPjIup_27GtPJ")
	if err != nil {
		panic(err)
	}
	clientMat, err := mautrix.NewClient("https://matrix.org", "@lukashuan:matrix.org", "syt_bHVrYXNodWFu_cCIsNOYxvmTBHznkzFFa_3pqWPS")
	if err != nil {
		panic(err)
	}

	// Setup the matrix client to receive messages
	syncer := clientMat.Syncer.(*mautrix.DefaultSyncer)
	syncer.OnEventType(event.EventMessage, func(source mautrix.EventSource, evt *event.Event) {
		fmt.Println(evt.Content.AsMessage().Body)
		fmt.Println(evt.Sender)
		fmt.Println(time.UnixMilli(evt.Timestamp).GoString())
	})

	syncCtx, cancelSync := context.WithCancel(context.Background())
	var syncStopWait sync.WaitGroup
	syncStopWait.Add(1)
	go func() {
		err = clientMat.SyncWithContext(syncCtx)
		defer syncStopWait.Done()
		if err != nil && !errors.Is(err, context.Canceled) {
			panic(err)
		}
	}()

	if _, err := clientTum.SendText("!SxURpLPLXhLsLHWqAP:tum.de", "Hello, world!"); err != nil {
		panic(err)
	}

	cancelSync()
	syncStopWait.Wait()
}
