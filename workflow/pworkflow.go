package workflow

import (
	"log"

	"github.com/jczornik/glacier_backup/persistent"
	"github.com/jczornik/glacier_backup/persistent/workflows"
)

type PWorkflow struct {
	id   int64
	workflow Workflow
	db    persistent.DBClient
}

func NewPWorkflow(tasks []task, client persistent.DBClient) (PWorkflow, error) {
	db, err := client.OpenDB()
	if err != nil {
		return PWorkflow{}, err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return PWorkflow{}, err
	}

	id, err := workflows.Create(tx)
	if err != nil {
		terr := tx.Rollback()
		if terr != nil {
			return PWorkflow{}, terr
		}
		return PWorkflow{}, err
	}

	ptasks := make([]task, len(tasks))
	for i, t := range tasks {
		pt, err := newPTask(client, t, id, tx)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return PWorkflow{}, err
			}
			return PWorkflow{}, err
		}

		ptasks[i] = pt
	}

	workflow := BasicWorkflow{ptasks}
	err = tx.Commit()

	return PWorkflow{id, workflow, client}, nil
}

func (flow PWorkflow) Exec() *WorkflowError {
	db, dbErr := flow.db.OpenDB()
	if dbErr != nil {
		return &WorkflowError{dbErr, nil}
	}
	defer db.Close()

	if dbErr = workflows.UpdateStatus(db, flow.id, workflows.RunningStatus); dbErr != nil {
		return &WorkflowError{dbErr, nil}
	}

	err := flow.workflow.Exec()

	if err != nil {
		if err.execError != nil {
			// TODO: Handle db error
			workflows.UpdateStatus(db, flow.id, workflows.RollbackedStatus)
		} else if err.rollbackError != nil {
			// TODO: Handle db error
			workflows.UpdateStatus(db, flow.id, workflows.FailedRollbackStatus)
		} else {
			log.Fatal("This should not happen - error should have either exec or rollback error")
		}
	} else {
		if dbErr = workflows.UpdateStatus(db, flow.id, workflows.FinishedStatus); dbErr != nil {
			return &WorkflowError{dbErr, nil}
		}
	}

	return err
}
