package workflow

import (
	"errors"
	"testing"
)

type execfunc = func() error
type rollbackfunc = func() error

type simpleTask struct {
	executed bool
	rollbacked bool

	execf execfunc
	rollbackf rollbackfunc
}

func (t *simpleTask) Exec() error {
	t.executed = true
	return t.execf()
}

func (t *simpleTask) Rollback() error {
	t.rollbacked = true
	return t.rollbackf()
}

func newSimpleTask(e execfunc, r rollbackfunc) *simpleTask {
	return &simpleTask{false, false, e, r}
}

func TestExecAll(t *testing.T) {
	// Given
	v := func() error {
		return nil
	}

	actions := []action{newSimpleTask(v, v), newSimpleTask(v, v)}
	w := Workflow{actions}

	// When
	err := w.Exec()

	// Then
	if err != nil {
		t.Errorf("Workflow returned err: %s", err)
	}

	for no, task := range w.actions {
		st, ok := task.(*simpleTask)
		if !ok {
			t.Error("Expected task to be instance of simpleTask")
		}

		if st.executed == false || st.rollbacked == true {
			t.Errorf("Task number %d should be executed and not rollbacked: executed == %t, rollbacked == %t", no, st.executed, st.rollbacked)
		}
	}
}

func TestFailShouldNotCallNextTaskAndRollback(t *testing.T) {
	// Given
	ff := func() error {
		return errors.New("First fail")
	}

	sf := func() error {
		return nil
	}

	rollback := func() error {
		return nil
	}

	task1 := newSimpleTask(ff, rollback)
	task2 := newSimpleTask(sf, rollback)

	flow := Workflow{[]action{task1, task2}}

	// When
	err := flow.Exec();

	// Then
	if err == nil {
		t.Error("Flow should return an error")
	}

	if err.execError == nil {
		t.Error("Flow error should return exec error")
	}

	if err.rollbackError != nil {
		t.Error("Flow error should not return rallback error")
	}

	if task2.executed || task2.rollbacked {
		t.Error("Second task should not be executer nor rollbacked")
	}

	if !task1.executed || !task1.rollbacked {
		t.Errorf("First task should be executed and rollbacked: executed == %t, rollbacked == %t", task1.executed, task1.rollbacked)
	}
}
