package workflow

import (
	"fmt"

	"github.com/jczornik/glacier_backup/persistent"
	"github.com/jczornik/glacier_backup/persistent/workflows"
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
	return fmt.Sprintf("Exec error: %s\nRollback error %s", error.execError, error.rollbackError)
}

type Workflow struct {
	tasks []ptask
	db    persistent.DBClient
}

func NewWorkflow(tasks []task, client persistent.DBClient) (Workflow, error) {
	db, err := client.OpenDB()
	if err != nil {
		return Workflow{}, err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return Workflow{}, err
	}

	id, err := workflows.Create(tx)
	if err != nil {
		terr := tx.Rollback()
		if terr != nil {
			return Workflow{}, terr
		}
		return Workflow{}, err
	}

	ptasks := make([]ptask, len(tasks))
	for i, t := range tasks {
		pt, err := newPTask(client, t, id, tx)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return Workflow{}, err
			}
			return Workflow{}, err
		}

		ptasks[i] = pt
	}

	workflow := Workflow{ptasks, client}
	err = tx.Commit()

	return workflow, err
}

func (flow Workflow) Exec() *WorkflowError {
	for i, task := range flow.tasks {
		if execError := task.exec(); execError != nil {
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
		if err := flow.tasks[i].rollback(); err != nil {
			return err
		}
	}

	return nil
}
