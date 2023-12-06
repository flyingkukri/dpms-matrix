package engines

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/mitchellh/mapstructure"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
	"tum.de/huan/dpms-matrix/data"
)

var (
	red   string = "\033[31m"
	reset string = "\033[0m"
	blue  string = "\033[34m"
)

type Engine struct {
	clientSend        *mautrix.Client
	clientReceive     *mautrix.Client
	cancelSync        context.CancelFunc
	syncStopWait      *sync.WaitGroup
	process           *data.Task
	currentTask       *data.Task
	taskToSend        *data.Task
	name              string
	requirementsCheck func([]data.Constraint) bool
	searchingForNext  bool
	completeTask      func()
}

type EngineInterface interface {
	Stop()
}

func NewEngine(process *data.Task, currentTask *data.Task, name string, matrixInfo data.MatrixInfo) *Engine {
	//client, cancelSync, syncStopWait := createMatrixClient("https://matrix.org", "@lukashuan:matrix.org", "syt_bHVrYXNodWFu_cCIsNOYxvmTBHznkzFFa_2pqWPS")
	engine := Engine{
		process:     process,
		currentTask: currentTask,
		name:        name,
	}
	clientSend, clientReceive, cancelSync, syncStopWait := createMatrixClients(matrixInfo.Server, matrixInfo.UserID, matrixInfo.AuthToken, &engine)
	engine.clientSend = clientSend
	engine.clientReceive = clientReceive
	engine.cancelSync = cancelSync
	engine.syncStopWait = syncStopWait
	return &engine
}

func (e *Engine) FindNextMachine(newTask *data.Task) {
	if newTask == nil {
		fmt.Println(blue + "Process is finished" + reset)
		return
	}
	//	constraints := newTask.Constraint
	jsonified := map[string]interface{}{
		"msgtype":      "de.tum.huan/requirements",
		"body":         "This message is used to send the requirements to the engines",
		"requirements": newTask.Requirements,
	}
	e.searchingForNext = true
	e.taskToSend = newTask
	_, err := e.clientSend.SendMessageEvent("!SxURpLPLXhLsLHWqAP:tum.de", event.EventMessage, jsonified)
	if err != nil {
		fmt.Printf("failed to send message %v", err)
		//fmt.Println(response.EventID.String())
	}
	log.Printf("%s sent the requirements\n", e.name)
}

func (e Engine) Stop() {
	e.cancelSync()
	e.syncStopWait.Wait()
}

func getRequirements(evt *event.Event) ([]data.Constraint, error) {
	if reqs, ok := evt.Content.Raw["requirements"]; ok {
		if cmp.Equal(reqs, []any{}) {
			return []data.Constraint{}, nil
		}
		return data.GetConstraints(reqs)
	}
	return nil, errors.New("requirements not found")
}

func findDirectChat(client *mautrix.Client, userID id.UserID) (id.RoomID, error) {
	fmt.Println("Finding direct chat")
	// Get the user's rooms
	rooms, err := client.JoinedRooms()
	if err != nil {
		return "", err
	}
	// Iterate over the rooms
	for _, room := range rooms.JoinedRooms {
		// Get the room info
		statestore := client.StateStore

		if !(statestore.GetMember(room, userID).Membership == event.MembershipJoin) {
			continue
		}
		if !(statestore.GetMember(room, client.UserID).Membership == event.MembershipJoin) {
			continue
		}
		members, err := statestore.GetRoomJoinedOrInvitedMembers(room)
		if err != nil {
			return "", fmt.Errorf("failed to get room members: %w", err)
		}
		if len(members) != 2 {
			continue
		}
		//		if !slices.Contains(members, userID) {
		//			continue
		//		}
		//		if !slices.Contains(members, client.UserID) {
		//			continue
		//		}
		return room, nil
	}
	// otherwise create a new room
	return createNewDirectRoom(client, userID)
}

func createNewDirectRoom(client *mautrix.Client, userID id.UserID) (id.RoomID, error) {
	resp, err := client.CreateRoom(&mautrix.ReqCreateRoom{
		Visibility: "private",
		Invite:     []id.UserID{userID},
		IsDirect:   true,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create room: %w", err)
	}
	return resp.RoomID, nil
}

func sendConfirmation(to id.UserID, room id.RoomID, client *mautrix.Client) error {
	jsonified := map[string]interface{}{
		"msgtype": "de.tum.huan/confirmreqs",
		"body":    "This message is used to confirm that the requirements are fulfilled",
	}
	_, err := client.SendMessageEvent(room, event.EventMessage, jsonified)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	fmt.Printf("Sent confirmation to %s\n", room)
	return nil
}

func answerRequirements(engine *Engine, evt *event.Event, clientSend *mautrix.Client) {
	log.Printf("%s Received requirements \n", engine.name)
	// Get the requirements
	reqs, err := getRequirements(evt)
	if err != nil {
		panic(err)
	}
	if engine.requirementsCheck == nil {
		panic("requirementsCheck has to be set")
	}
	if engine.requirementsCheck(reqs) {
		log.Printf("%s checked that Requirements fulfilled", engine.name)
		room, err := findDirectChat(engine.clientReceive, evt.Sender)
		if err != nil {
			panic(err)
		}
		sendConfirmation(evt.Sender, room, clientSend)
	}
}

func sendNextTask(engine *Engine, evt *event.Event, clientSend *mautrix.Client) {
	engine.searchingForNext = false
	jsonified := map[string]interface{}{
		"msgtype": "de.tum.huan/task",
		"body":    fmt.Sprintf("This message is used to send the next task called %s", engine.taskToSend.Name),
		"task":    engine.taskToSend,
	}
	_, err := clientSend.SendMessageEvent(evt.RoomID, event.EventMessage, jsonified)
	if err != nil {
		panic("unable to send next task")
	}
	log.Printf("%s Is sending the task data for %s\n", engine.name, engine.taskToSend.Name)
}

func retrieveTask(evt *event.Event) data.Task {
	if task, ok := evt.Content.Raw["task"]; ok {
		if cmp.Equal(task, nil) {
			panic("task is empty")
		}
		taskConverted := data.Task{}
		if err := mapstructure.Decode(task, &taskConverted); err != nil {
			panic(fmt.Errorf("failed to decode task: %w", err))
		} else {
			return taskConverted
		}
	}
	panic("There is no task.")
}
func (e *Engine) initTask(evt *event.Event) {
	taskConverted := retrieveTask(evt)
	e.currentTask = &taskConverted
	e.completeTask()
}

func createMatrixClients(server string, userID id.UserID, authToken string, engine *Engine) (*mautrix.Client, *mautrix.Client, context.CancelFunc, *sync.WaitGroup) {
	lastProcessedTimestamp := time.Now()
	// Create clients
	clientSend, err := mautrix.NewClient(server, userID, authToken)
	clientReceive, err := mautrix.NewClient(server, userID, authToken)
	if err != nil {
		panic(err)
	}
	// Setup the matrix client to receive messages
	//syncer := clientReceive.Syncer.(*mautrix.DefaultSyncer)

	clientReceive.StateStore = mautrix.NewMemoryStateStore()
	clientReceive.Syncer.(mautrix.ExtensibleSyncer).OnEventType(event.EventMessage, func(source mautrix.EventSource, evt *event.Event) {
		if evt.Timestamp < lastProcessedTimestamp.UnixMilli()-1000 {
			//	lastProcessedTimestamp = time.Now()
			return
		}
		log.Printf("%s Received message \"%s\" from %s\n", engine.name, evt.Content.AsMessage().Body, evt.Sender)
		//lastProcessedTimestamp = time.Now()
		if evt.Sender == clientReceive.UserID {
			return
		}
		if evt.Content.Raw["msgtype"] == "de.tum.huan/requirements" {
			answerRequirements(engine, evt, clientSend)
		} else if evt.Content.Raw["msgtype"] == "de.tum.huan/confirmreqs" && engine.searchingForNext {
			sendNextTask(engine, evt, clientSend)
		} else if evt.Content.Raw["msgtype"] == "de.tum.huan/task" {
			engine.initTask(evt)
		}
	})
	clientReceive.StateStore = mautrix.NewMemoryStateStore()
	clientReceive.Syncer.(mautrix.ExtensibleSyncer).OnEventType(event.StateMember,
		clientReceive.StateStoreSyncHandler)
	clientReceive.Syncer.(mautrix.ExtensibleSyncer).OnEventType(event.StateMember,
		func(source mautrix.EventSource, evt *event.Event) {
			if evt.Content.AsMember().Membership.IsInviteOrJoin() {
				clientReceive.JoinRoom(evt.RoomID.String(), "", nil)
			}
		})

	syncCtx, cancelSync := context.WithCancel(context.Background())
	var syncStopWait sync.WaitGroup
	syncStopWait.Add(1)
	go func() {
		err = clientReceive.SyncWithContext(syncCtx)
		defer syncStopWait.Done()
		if err != nil && !errors.Is(err, context.Canceled) {
			panic(err)
		}
	}()

	return clientSend, clientReceive, cancelSync, &syncStopWait
}
