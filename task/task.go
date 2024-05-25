package task

import (
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
)

type State int

const (
	Pending State = iota
	Scheduled
	Running
	Completed
	Failed
)

type Task struct {
	ID    uuid.UUID
	Name  string
	State State
	// The following fields are Docker specific. This means our app is coupled to Docker.
	Image         string
	CPU           float64
	Memory        int64
	Disk          int64
	ExportedPorts nat.PortSet
	PortBindings  map[string]string
	RestartPolicy string
	// Some basic metrics for the task. Not sure if these are also Docker specific.
	// As in I am not sure if this code is shaped to deal with the Docker API.
	// But I am going to assume it is.
	StartTime  time.Time
	FinishTime time.Time
}

type TaskEvent struct {
	ID        uuid.UUID
	State     State
	Timestamp time.Time
	Task      Task
}

// With these two structs we have defined the skeleton of how a Task is and how a Task can move from one state to the other
