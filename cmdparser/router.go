package cmdparser

import (
	"strings"

	"github.com/the-forges/example-net/internal"
)

type handleFunc struct {
	callback func(internal.Conn, internal.Command) internal.Result
}

func (h handleFunc) Call(conn internal.Conn, cmd internal.Command) internal.Result {
	return h.callback(conn, cmd)
}

func newHandleFunc(handler func(conn internal.Conn, cmd internal.Command) internal.Result) handleFunc {
	return handleFunc{handler}
}

type router struct {
	Handlers map[string]internal.CommandHandler
}

func (r *router) ensureHandlersExist() {
	if r.Handlers == nil {
		r.Handlers = make(map[string]internal.CommandHandler, 0)
	}
}

func (r *router) addHandler(cmd string, handler internal.CommandHandler) {
	r.ensureHandlersExist() // Prevents nil pointer issues with the handlers map
	r.Handlers[cmd] = handler
}

func (r *router) Handle(cmd string, handler internal.CommandHandler) {
	r.addHandler(cmd, handler)
}

func (r *router) HandleFunc(cmd string, handler func(conn internal.Conn, cmd internal.Command) internal.Result) {
	r.addHandler(cmd, newHandleFunc(handler))
}

func (r router) Parse(req string) internal.Result {
	cmd := newCommand(req, "", "")
	// Verify the command parses correctly
	req = strings.TrimSpace(req)
	if len(req) == 0 {
		return newResult(req, cmd, nil, errorForErrorType(ErrCommandEmpty))
	}
	if !strings.HasPrefix(req, "/") {
		return newResult(req, cmd, nil, errorForErrorType(ErrCommandMalformed))
	}
	// Extract command details
	parts := strings.SplitAfterN(req, " ", 2)
	cmd.action = strings.TrimSpace(parts[0])
	if len(cmd.action) == 1 {
		return newResult(req, cmd, nil, errorForErrorType(ErrCommandMalformed))
	}
	if len(parts) == 2 {
		cmd.body = parts[1]
	}
	// See if a handler exists for the requested command
	r.ensureHandlersExist() // Prevents nil pointer issues with the handlers map
	handler, ok := r.Handlers[cmd.action]
	if !ok {
		return newResult(req, cmd, nil, errorForErrorType(ErrCommandNotFound))
	}
	return newResult(req, cmd, handler, nil)
}

func (r router) ParseAndCall(conn internal.Conn, req string) internal.Result {
	pr := r.Parse(req)
	if pr == nil || !pr.Ok() {
		return pr
	}
	preq, ok := pr.(result)
	if !ok {
		return pr
	}
	return preq.Handler.Call(conn, preq.Cmd)
}

func NewRouter() internal.CommandRouter {
	r := &router{}
	r.ensureHandlersExist()
	return r
}
