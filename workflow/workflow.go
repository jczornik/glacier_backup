package workflow

import (
	"fmt"
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
	AddTasks([]task) error
}

type BasicWorkflow struct {
	tasks []task
}

func newBasicWorkflow(tasks []task) Workflow {
	return &BasicWorkflow{tasks}
}

func (flow *BasicWorkflow) Exec() *WorkflowError {
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

func (flow *BasicWorkflow) rollback(idx int) error {
	for i := idx; i >= 0; i-- {
		if err := flow.tasks[i].Rollback(); err != nil {
			return err
		}
	}

	return nil
}

func (flow *BasicWorkflow) AddTasks(tasks []task) error {
	flow.tasks = append(flow.tasks, tasks...)
	return nil
}
