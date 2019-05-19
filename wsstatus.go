package wsstatus

import (
	"fmt"
	"log"
	"sync"
)

// Handler websocket status handler
type Handler struct {
	mut      sync.Mutex
	payloads map[string][]byte
	status   []string
	cur      int
}

// NewHandler new websocket status handler
func NewHandler() *Handler {
	return &Handler{
		payloads: make(map[string][]byte),
		status:   make([]string, 0),
	}
}

// Registry payload to state
func (h *Handler) Registry(state string, payload []byte) {
	h.mut.Lock()
	defer h.mut.Unlock()

	h.payloads[state] = []byte(payload)
	h.status = append(h.status, state)
}

// Next go to next state
func (h *Handler) Next() {
	h.mut.Lock()
	defer h.mut.Unlock()

	h.cur++
	if h.cur == len(h.status) {
		h.cur = 0
	}
}

// GetStatus get current state
func (h *Handler) GetStatus() string {
	h.mut.Lock()
	defer h.mut.Unlock()

	return h.status[h.cur]
}

// NextTo go to state
func (h *Handler) NextTo(state string) {
	h.mut.Lock()
	defer h.mut.Unlock()

	for k, v := range h.status {
		if state == v {
			h.cur = k
			return
		}
	}

	log.Panic("can not find state,check registry")
}

// GetPayload get payload by status
func (h *Handler) GetPayload(args ...interface{}) []byte {
	h.mut.Lock()
	defer h.mut.Unlock()

	if len(args) == 0 {
		return h.payloads[h.status[h.cur]]
	}

	return []byte(fmt.Sprintf(string(h.payloads[h.status[h.cur]]), args...))
}
