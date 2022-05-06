package game

import "hlinspect/internal/hlrpc/schema"

type Sync struct {
	inputControlReqChan  chan *schema.CommandInput
	inputControlRespChan chan bool
}

func NewSync() *Sync {
	return &Sync{}
}

func (s *Sync) InputControlReq(handlerSide bool) chan *schema.CommandInput {
	if handlerSide {
		return s.inputControlReqChan
	}

	if s.inputControlReqChan != nil && s.inputControlRespChan == nil {
		s.inputControlRespChan = make(chan bool)
	}
	return s.inputControlReqChan
}

func (s *Sync) InputControlResp(handlerSide bool) chan bool {
	if handlerSide {
		return s.inputControlRespChan
	}

	c := s.inputControlRespChan
	if s.inputControlReqChan == nil && s.inputControlRespChan != nil {
		close(s.inputControlRespChan)
		s.inputControlRespChan = nil
	}
	return c
}

func (s *Sync) StartInputControl() {
	if s.inputControlReqChan != nil {
		return
	}
	s.inputControlReqChan = make(chan *schema.CommandInput)
}

func (s *Sync) StopInputControl() {
	if s.inputControlReqChan == nil {
		return
	}
	select {
	case s.inputControlReqChan <- nil:
		close(s.inputControlRespChan)
		s.inputControlRespChan = nil
	default:
	}
	close(s.inputControlReqChan)
	s.inputControlReqChan = nil
}
