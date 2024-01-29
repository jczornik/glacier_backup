package workflow

import (
	"fmt"

	"github.com/jczornik/glacier_backup/persistent"
)

type task interface {
	Exec() error
	Rollback() error
	Name() string
}

type WorkflowError struct {
	execError     error
	rollbackError error
}

func (error WorkflowError) Error() string {
	return fmt.Sprintf("Exec error: %s\nRollback error %s\n", error.execError, error.rollbackError)
}

type Workflow interface {
	Exec() *WorkflowError
}

type BasicWorkflow struct {
	tasks []task
}

func newBasicWorkflow(tasks []task, client persistent.DBClient) Workflow {
	return BasicWorkflow{tasks}
}

func (flow BasicWorkflow) Exec() *WorkflowError {
	for i, task := range flow.tasks {
		if execError := task.Exec(); execError != nil {
			if rollbackError := flow.rollback(i); rollbackError != nil {
				return &WorkflowError{execError, rollbackError}
			}

			return &WorkflowError{execError, nil}
		}
	}

	return nil
}

func (flow BasicWorkflow) rollback(idx int) error {
	for i := idx; i >= 0; i-- {
		if err := flow.tasks[i].Rollback(); err != nil {
			return err
		}
	}

	return nil
}
