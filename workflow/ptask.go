package workflow

import (
	"database/sql"

	"github.com/jczornik/glacier_backup/persistent"
	"github.com/jczornik/glacier_backup/persistent/tasks"
)

type ptask struct {
	id   int64
	db   persistent.DBClient
	task task
}

func newPTask(c persistent.DBClient, task task, workflow int64, tx *sql.Tx) (ptask, error) {
	id, err := tasks.Create(tx, workflow, task.Name())
	if err != nil {
		return ptask{}, err
	}

	return ptask{id, c, task}, nil
}

func (pt ptask) Name() string {
	return pt.task.Name()
}

func (pt ptask) Exec() error {
	db, err := pt.db.OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	err = tasks.UpdateStatus(db, pt.id, tasks.RunningStatus)
	if err != nil {
		return err
	}

	res := pt.task.Exec()
	if res != nil {
		err = tasks.UpdateStatus(db, pt.id, tasks.FailedStatus)
		if err != nil {
			return err
		}
		return res
	}

	err = tasks.UpdateStatus(db, pt.id, tasks.FinishedStatus)
	if err != nil {
		return err
	}

	return res
}

func (pt ptask) Rollback() error {
	db, err := pt.db.OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	err = tasks.UpdateStatus(db, pt.id, tasks.RollbackStatus)
	if err != nil {
		return err
	}

	res := pt.task.Rollback()
	if res != nil {
		err = tasks.UpdateStatus(db, pt.id, tasks.FailedRollbackStatus)
		if err != nil {
			return err
		}
		return res
	}

	err = tasks.UpdateStatus(db, pt.id, tasks.RollbackedStatus)
	if err != nil {
		return err
	}

	return res
}
