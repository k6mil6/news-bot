package state

import (
	"sync"
)

type Machine struct {
	userState map[int64]State
	mu        sync.Mutex
}

func NewMachine() *Machine {
	return &Machine{
		userState: make(map[int64]State),
	}
}

func (m *Machine) Set(user int64, state State) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.userState[user] = state
}

func (m *Machine) Get(user int64) (State, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	userState, ok := m.userState[user]
	return userState, ok
}
