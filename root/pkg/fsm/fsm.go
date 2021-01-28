package fsm

import "fmt"

// 状态对象
type stateHandler interface {
	State() int
	Enter(*FSM)
	Leave(*FSM)
	Handle(*FSM, interface{})
}

// 状态机
type (
	FSM struct {
		curr_state int                  // 当前状态
		states     map[int]stateHandler // 状态处理对象
	}
)

// 创建一个状态机
func New() *FSM {
	fsm := &FSM{}
	fsm.curr_state = -1 // 此状态无效
	fsm.states = make(map[int]stateHandler, 2)
	return fsm
}

// 新增一个状态
func (s *FSM) Add(state stateHandler) error {
	if _, exit := s.states[state.State()]; exit {
		return fmt.Errorf("repeated state:%v", state.State())
	}
	s.states[state.State()] = state
	return nil
}

// 当前状态
func (s *FSM) State() int {
	return s.curr_state
}

// 当前状态对象
func (s *FSM) CurrentStateHandler() stateHandler {
	return s.states[s.curr_state]
}

// 状态迁移
func (s *FSM) Swtich(next_state int) error {
	if _, exit := s.states[next_state]; !exit {
		return fmt.Errorf("not found next_state:%v", next_state)
	}

	if s.State() != -1 {
		s.states[s.State()].Leave(s)
	}
	s.curr_state = next_state
	s.states[s.State()].Enter(s)

	return nil
}

func (s *FSM) Handle(v interface{}) {
	s.CurrentStateHandler().Handle(s, v)
}
