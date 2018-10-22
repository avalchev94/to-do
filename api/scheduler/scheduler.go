package scheduler

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/avalchev94/to-do-app/database"
	"github.com/go-redis/redis"
	"github.com/rickb777/date"
)

// Scheduler is responsible for asynchronously replicating repetitive task
// when their "repeat date(time)" comes. It also, checks if there are repetitive
// tasks which haven't been replicating(in case of down time or similiar reason)
type Scheduler struct {
	time        time.Time
	redisClient *redis.Client
	db          *sql.DB
}

// New creates a Scheduler. However, you should call Run consequently.
func New(redisCli *redis.Client, db *sql.DB) (*Scheduler, error) {
	var startTime time.Time

	result := redisCli.Get("last_task_repeat")
	switch result.Err() {
	case nil:
		if err := result.Scan(&startTime); err != nil {
			return nil, fmt.Errorf("New: can't scan 'last_task_repeat' value")
		}
	case redis.Nil:
		startTime = time.Now().UTC()
	default:
		return nil, result.Err()
	}

	return &Scheduler{startTime, redisCli, db}, nil
}

// Run starts the Scheduler. It is adjustable on the number of the threads to be used
// for scheduling the repetitive tasks.
// Run could be called on different thread(go Run(..))
func (t *Scheduler) Run(thread int) error {
	tasks := make(chan *database.RepetitiveTask, 100)
	scheduled := date.NewAt(t.time)

	for i := 0; i < thread; i++ {
		go runThread(i, tasks, &scheduled, t.db)
	}

	for {
		scheduled = date.NewAt(t.time)

		err := database.GetRepetitiveTasksAt(t.time, tasks, t.db)
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("scheduler: GetRepetitiveTasksAt failed: %v", err)
		}

		err = t.sleep()
		if err != nil {
			return err
		}
	}
}

func (t *Scheduler) sleep() error {
	err := t.redisClient.Set("last_task_repeat", t.time, 0).Err()
	if err != nil {
		return fmt.Errorf("scheduler: couldn't save last_task_repeat: %v", err)
	}

	t.time = t.time.Add(time.Minute)
	now := time.Now().UTC()
	if t.time.After(now) {
		time.Sleep(t.time.Sub(now))
	}
	return nil
}

func runThread(id int, tasks <-chan *database.RepetitiveTask, schedule *date.Date, db *sql.DB) {
	log.Printf("Scheduler: thread %d started.\n", id)

	for {
		t := <-tasks
		if err := t.Schedule(*schedule, db); err != nil {
			log.Printf("Scheduler: thread %d FAILED to process task %d\n", id, t.Task.ID)
		} else {
			log.Printf("Scheduler: thread %d SUCCEEDED to process task %d\n", id, t.Task.ID)
		}
	}

}
