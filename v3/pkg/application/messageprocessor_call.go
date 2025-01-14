package application

import (
	"context"
	"fmt"
	"net/http"
	"encoding/json"
)

const (
	CallBinding = 0
)

func (m *MessageProcessor) callErrorCallback(window Window, message string, callID *string, err error) {
	errorMsg := fmt.Sprintf(message, err)
	m.Error(errorMsg)
	window.CallError(*callID, errorMsg)
}

func (m *MessageProcessor) callCallback(window Window, callID *string, result string, isJSON bool) {
	window.CallResponse(*callID, result)
}

func (m *MessageProcessor) processCallCancelMethod(method int, rw http.ResponseWriter, r *http.Request, window Window, params QueryParams) {
	args, err := params.Args()
	if err != nil {
		m.httpError(rw, "Unable to parse arguments: %s", err.Error())
		return
	}
	callID := args.String("call-id")
	if callID == nil || *callID == "" {
		m.Error("call-id is required")
		return
	}

	m.l.Lock()
	cancel := m.runningCalls[*callID]
	m.l.Unlock()

	if cancel != nil {
		cancel()
	}
	m.ok(rw)
}

func (m *MessageProcessor) processCallMethod(method int, rw http.ResponseWriter, r *http.Request, window Window, params QueryParams) {
	args, err := params.Args()
	if err != nil {
		m.httpError(rw, "Unable to parse arguments: %s", err.Error())
		return
	}
	callID := args.String("call-id")
	if callID == nil || *callID == "" {
		m.Error("call-id is required")
		return
	}

	switch method {
	case CallBinding:
		var options CallOptions
		err := params.ToStruct(&options)
		if err != nil {
			m.callErrorCallback(window, "Error parsing call options: %s", callID, err)
			return
		}
		var boundMethod *BoundMethod
		if options.PackageName != "" {
			boundMethod = globalApplication.bindings.Get(&options)
			if boundMethod == nil {
				m.callErrorCallback(window, "Error getting binding for method: %s", callID, fmt.Errorf("method '%s' not found", options.Name()))
				return
			}
		} else {
			boundMethod = globalApplication.bindings.GetByID(options.MethodID)
		}
		if boundMethod == nil {
			m.callErrorCallback(window, "Error getting binding for method: %s", callID, fmt.Errorf("method ID '%s' not found", options.Name()))
			return
		}

		ctx, cancel := context.WithCancel(context.WithoutCancel(r.Context()))

		ambiguousID := false
		m.l.Lock()
		if m.runningCalls[*callID] != nil {
			ambiguousID = true
		} else {
			m.runningCalls[*callID] = cancel
		}
		m.l.Unlock()

		if ambiguousID {
			cancel()
			m.callErrorCallback(window, "Error calling method: %s, a method call with the same id is already running", callID, err)
			return
		}

		go func() {
			defer func() {
				cancel()

				m.l.Lock()
				delete(m.runningCalls, *callID)
				m.l.Unlock()
			}()

			result, err := boundMethod.Call(ctx, options.Args)
			if err != nil {
				m.callErrorCallback(window, "Error calling method: %s", callID, err)
				return
			}
			var jsonResult = []byte("{}")
			if result != nil {
				// convert result to json
				jsonResult, err = m.serializer(result)
				if err != nil {
					m.callErrorCallback(window, "Error converting result to json: %s", callID, err)
					return
				}
			}
			m.callCallback(window, callID, string(jsonResult), true)

			var jsonArgs struct {
				Args json.RawMessage `json:"args"`
			}
			params.ToStruct(&jsonArgs)
			m.Info("Call Binding:", "method", boundMethod, "args", string(jsonArgs.Args), "result", result)
		}()
		m.ok(rw)
	default:
		m.httpError(rw, "Unknown call method: %d", method)
	}

}
