package riverlet

import "time"

type STATUS int8

const (
	FAILED    STATUS = -1
	AVAILABLE STATUS = 0
	PENDING   STATUS = 1
	RETRY     STATUS = 2
	RUNNING   STATUS = 3
	COMPLETE  STATUS = 4
)

type Instance struct {
}

type Queue struct {
	Seq         int64
	Args        Argument
	Priority    int
	Identity    string
	Attempts    int
	MaxAttempts int
	Status      STATUS
	Created     *time.Time
	Updated     *time.Time
}
