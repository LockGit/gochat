/**
 * Created by lock
 * Date: 2019-08-09
 * Time: 18:22
 */
package task

type Task struct {
}

func New() *Task {
	return new(Task)
}

func (task *Task) Run() {
	//read from redis queue
	//rpc call connect layer send msg
}
