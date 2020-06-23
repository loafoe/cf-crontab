package crontab

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"sync"
)

type State struct {
	list []*Task
	cronTab *cron.Cron
	mux sync.Mutex
}

func (e *State)Entries() []*Task {
	return e.list
}

func (e *State)StartCron() {
	e.cronTab.Start()
}

func (e *State)AddEntries(newEntries []Task) {
	e.mux.Lock()
	defer e.mux.Unlock()
	for i := range newEntries {
		_ = newEntries[i].Add(e.cronTab)
		e.list = append(e.list, &newEntries[i])
	}
}

func (e *State) DeleteEntry(id int) error {
	e.mux.Lock()
	defer e.mux.Unlock()
	entryID := cron.EntryID(id)
	for i, t := range e.Entries() {
		if int(t.EntryID) == id {
			fmt.Printf("Removing %d\n", id)
			e.list = append(e.list[:i], e.list[i+1:]...)
			e.cronTab.Remove(entryID)
			return nil
		}
	}
	return fmt.Errorf("entry not found: %d", id)
}

func NewState() *State {
	state := &State{
		cronTab: cron.New(cron.WithSeconds()),
		list: make([]*Task, 0),
	}
	return state
}
