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
	Respond(msgType, sessionId, token, format string, params ...interface{})
	FlushMessage(session *SMTPSession)
}

type EventResponderImpl struct {
	SafePrinter
	event FilterEvent
}

func (evr *EventResponderImpl) Proceed() {
	evr.Respond("filter-result", evr.event.GetSessionId(), evr.event.GetToken(), "%s", "proceed")
}

func (evr *EventResponderImpl) HardReject(response string) {
	evr.Respond("filter-result", evr.event.GetSessionId(), evr.event.GetToken(),
		"reject|550 %s", response)
}

func (evr *EventResponderImpl) Greylist(response string) {
	evr.Respond("filter-result", evr.event.GetSessionId(), evr.event.GetToken(),
		"reject|421 %s", response)
}

func (evr *EventResponderImpl) SoftReject(response string) {
	evr.Respond("filter-result", evr.event.GetSessionId(), evr.event.GetToken(),
		"reject|451 %s", response)
}

func (evr *EventResponderImpl) FlushMessage(session *SMTPSession) {
	for _, line := range session.Message {
		evr.DatalineReply(line)
	}
	evr.DatalineEnd()
}

func (evr *EventResponderImpl) DatalineEnd() {
	evr.Respond("filter-dataline", evr.event.GetSessionId(), evr.event.GetToken(), "%s", ".")
}

func (evr *EventResponderImpl) DatalineReply(line string) {
	prefix := ""
	// Output raw SMTP data - escape leading dots.
	if strings.HasPrefix(line, ".") {
		prefix = "."
	}
	evr.Respond("filter-dataline", evr.event.GetSessionId(), evr.event.GetToken(), "%s%s", prefix, line)
}

func (evr *EventResponderImpl) WriteMultilineHeader(header, value string) {
	token := evr.event.GetToken()
	sessionId := evr.event.GetSessionId()
	for i, line := range strings.Split(value, "\n") {
		if i == 0 {
			evr.Respond("filter-dataline", sessionId, token, "%s: %s", header, line)
		} else {
			evr.Respond("filter-dataline", sessionId, token, "%s", line)
		}
	}
}

func (evr *EventResponderImpl) Respond(msgType, sessionId, token, format string, params ...interface{}) {
	var prefix string
	if evr.event.GetProtocolVersion() > "0.5" {
		prefix = msgType + "|" + sessionId + "|" + token
	} else {
		prefix = msgType + "|" + token + "|" + sessionId
	}
	evr.SafePrintln(prefix + "|" + fmt.Sprintf(format, params...))
}

func NewEventResponder(_event FilterEvent) EventResponder {
	resp := EventResponderImpl{
		event: _event,
	}

	return &resp
}
