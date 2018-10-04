package api

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/rickb777/date"

	"github.com/avalchev94/to-do-app/database"
)

// TaskRepeater is responsible for asynchronously replicating repetitive task
// when their "repeat date(time)" comes. It also, checks if there are repetitive
// tasks which haven't been replicating(in case of down time or similiar reason)
type TaskRepeater struct {
	time time.Time
}

// NewTaskRepeater creates a TaskRepeater. However, you should call Run consequently.
func NewTaskRepeater() (*TaskRepeater, error) {
	var startTime time.Time

	result := redisConn.Get("last_task_repeat")
	switch result.Err() {
	case nil:
		if err := result.Scan(&startTime); err != nil {
			return nil, fmt.Errorf("NewTaskRepeater: can't scan 'last_task_repeat' value")
		}
	case redis.Nil:
		startTime = time.Now().UTC()
	default:
		return nil, result.Err()
	}

	return &TaskRepeater{startTime}, nil
}

// Run starts the TaskRepeater. It is adjustable on the number of the threads to be used
// for scheduling the repetitive tasks.
// Run could be called on different thread(go Run(..))
func (t *TaskRepeater) Run(thread int) {
	log.Println("TaskRepeater: Runing threads..")

	tasks := make(chan *database.RepetitiveTask, 100)
	scheduleDate := date.NewAt(t.time)

	for i := 0; i < thread; i++ {
		go runThread(i, tasks, &scheduleDate)
	}

	for {
		scheduleDate = date.NewAt(t.time)

		err := database.GetRepetitiveTasksAt(t.time, tasks, db)
		if err != nil && err != sql.ErrNoRows {
			log.Println("TaskRepeater: GetRepetitiveTasksAt failed:", err)
		}

		t.sleep()
	}
}

func (t *TaskRepeater) sleep() {
	if err := redisConn.Set("last_task_repeat", t.time, 0).Err(); err != nil {
		log.Println("TaskRepeater: Couldn't save last_task_repeat:", err)
	}

	t.time = t.time.Add(time.Minute)
	now := time.Now().UTC()
	if t.time.After(now) {
		time.Sleep(t.time.Sub(now))
	}
}

func runThread(id int, tasks <-chan *database.RepetitiveTask, scheduleDate *date.Date) {
	log.Printf("TaskRepater: thread %d started.\n", id)

	for {
		t := <-tasks
		if err := t.Schedule(*scheduleDate, db); err != nil {
			log.Printf("TaskRepeater: thread %d FAILED to process task %d\n", id, t.Task.ID)
		} else {
			log.Printf("TaskRepeater: thread %d SUCCEEDED to process task %d\n", id, t.Task.ID)
		}
	}

}
