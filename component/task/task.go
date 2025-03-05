package task

type Config string

type TaskType string

const (
	TaskTypeMongo TaskType = "mongo"
	TaskTypeMySQL TaskType = "mysql"
)

type Task struct {
	TaskType TaskType
	Source   Source
	Target   Target
}
