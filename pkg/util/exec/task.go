package exec

import (
	"os/exec"
	"strings"

	"github.com/google/uuid"

	pb "simpleagent/pkg/proto/pbgo/simpleagent"
	"simpleagent/pkg/util/log"
	"simpleagent/pkg/util/task"
)

var (
	ResTasks []*pb.ExecTaskReply
	restask  *pb.ExecTaskReply
)

type ExecTask struct {
	Name    string
	Uuid    uuid.UUID
	Cmd     string
	Out     string
	Err     string
	Success bool
}

func (execTask *ExecTask) Execute() (string, error) {
	log.Infof("*Execute* [%s], UUID: [%v], CMD: %v \n", execTask.Name, execTask.Uuid, execTask.Cmd)
	out, err := exec.Command("sh", "-c", execTask.Cmd).Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

func (execTask *ExecTask) CallBack(out string, result *task.Result, err error) {
	var reserr interface{}
	if err != nil {
		reserr = err.Error()
	} else {
		reserr = ""
	}

	restask = &pb.ExecTaskReply{
		Name:    execTask.Name,
		Uuid:    execTask.Uuid.String(),
		Cmd:     execTask.Cmd,
		Error:   reserr.(string),
		Result:  out,
		Success: result.IsSuccessful(),
	}
	ResTasks = append(ResTasks, restask)

	if result.IsSuccessful() {
		log.Infof("*task [%s] exec success*, result:\n %v\n", execTask.Name, out)
	} else {
		log.Errorf("*task [%s] exec failed*, error: [%v] \n", execTask.Name, err)
	}
}

func FlushRes() {
	ResTasks = []*pb.ExecTaskReply{}
}