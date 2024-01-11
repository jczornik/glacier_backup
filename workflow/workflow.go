package workflow

import "fmt"

type action interface {
	Exec() error
	Rollback() error
}

type WorkflowError struct {
	execError error
	rollbackError error
}

func (error WorkflowError) Error() string {
	return fmt.Sprintf("Exec error: %s\nRollback error %s", error.execError, error.rollbackError)
}

type Workflow struct {
	actions []action
}

func (flow Workflow) Exec() *WorkflowError {
	for i, action := range flow.actions {
		if execError := action.Exec(); execError != nil {
			if rollbackError := flow.rollback(i); rollbackError != nil {
				return &WorkflowError{execError, rollbackError}
			}

			return &WorkflowError{execError, nil}
		}
	}

	return nil
}

func (flow Workflow) rollback(idx int) error {
	for i := idx; i >= 0; i-- {
		if err := flow.actions[i].Rollback(); err != nil {
			return err
		}
	}

	return nil
}
