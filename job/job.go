/**
 * Created by lock
 * Date: 2019-08-09
 * Time: 18:22
 */
package job

type Job struct {
}

func New() *Job {
	return new(Job)
}

func (job *Job) Run() {
	//read from redis queue
	//rpc call connect layer send msg
}
