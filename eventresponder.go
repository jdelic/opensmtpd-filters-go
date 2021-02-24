package opensmtpd

import (
	"fmt"
	"strings"
)

type EventResponder interface {
	Proceed()
	HardReject(response string)
	SoftReject(response string)
	Greylist(response string)
	DatalineReply(line string)
	DatalineEnd()
	WriteMultilineHeader(header, value string)
	SafePrintln(msg string)
	Respond(msgType, sessionId, token, format string, params... interface{})
	FlushMessage(session *SMTPSession)
}

type SafePrinter struct{}

type EventResponderImpl struct {
	SafePrinter
	event FilterEvent
}

func (evr *EventResponderImpl) Proceed() {
	evr.Respond("filter-result", evr.event.GetToken(), evr.event.GetSessionId(), "%s", "proceed")
}

func (evr *EventResponderImpl) HardReject(response string) {
	evr.Respond("filter-result", evr.event.GetToken(), evr.event.GetSessionId(),
		"reject|550 %s", response)
}

func (evr *EventResponderImpl) Greylist(response string) {
	evr.Respond("filter-result", evr.event.GetToken(), evr.event.GetSessionId(),
		"reject|421 %s", response)
}

func (evr *EventResponderImpl) SoftReject(response string) {
	evr.Respond("filter-result", evr.event.GetToken(), evr.event.GetSessionId(),
		"reject|451 %s", response)
}

func (evr *EventResponderImpl) FlushMessage(session *SMTPSession) {
	token := evr.event.GetToken()
	for _, line := range session.Message {
		evr.Respond("filter-dataline", token, session.Id, "%s", line)
	}
	evr.DatalineEnd()
}

func (evr *EventResponderImpl) DatalineEnd() {
	evr.Respond("filter-dataline", evr.event.GetToken(), evr.event.GetSessionId(), "%s", ".")
}

func (evr *EventResponderImpl) DatalineReply(line string) {
	prefix := ""
	// Output raw SMTP data - escape leading dots.
	if strings.HasPrefix(line, ".") {
		prefix = "."
	}
	evr.Respond("filter-dataline", evr.event.GetToken(), evr.event.GetSessionId(), "%s%s", prefix, line)
}

func (evr *EventResponderImpl) WriteMultilineHeader(header, value string) {
	token := evr.event.GetToken()
	sessionId := evr.event.GetSessionId()
	for i, line := range strings.Split(value, "\n") {
		if i == 0 {
			evr.Respond("filter-dataline", token, sessionId, "%s: %s", header, line)
		} else {
			evr.Respond("filter-dataline", token, sessionId, "%s", line)
		}
	}
}

func (evr *EventResponderImpl) Respond(msgType, sessionId, token, format string, params... interface{}) {
	var prefix string
	if evr.event.GetProtocolVersion() < "0.5" {
		prefix = msgType + "|" + token + "|" + sessionId
	} else {
		prefix = msgType + "|" + sessionId + "|" + token
	}
	evr.SafePrintln(prefix + "|" + fmt.Sprintf(format, params...))
}

func (sp *SafePrinter) SafePrintln(msg string) {
	stdoutChannel <- msg + "\n"
}

func NewEventResponder(_event FilterEvent) EventResponder {
	resp := EventResponderImpl{
		event: _event,
	}

	return &resp
}

