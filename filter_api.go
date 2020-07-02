//
// Copyright (c) 2019 Jonas Maurus <@jdelic>
//
// largely based on code originally developed for filter-rspamd
// Copyright (c) 2019 Gilles Chehade <gilles@poolp.org>
//
// Permission to use, copy, modify, and distribute this software for any
// purpose with or without fee is hereby granted, provided that the above
// copyright notice and this permission notice appear in all copies.
//
// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
// WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
// MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
// ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
// WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
// ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
// OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
//

package opensmtpd

import (
	"bufio"

	"fmt"
	"log"

	"os"
	"strings"
)

/*
 * Every filter must implement this interface
 */
type Filter interface {
	GetName() string
}

/*
 * A general type for filter event handlers
 */
type EventHandler = func(FilterWrapper, FilterEvent)

/*
 * Used to store callbacks in a filter and map them to received events
 */
type FilterDispatchMap = map[string]map[string]EventHandler


type EventResponder interface {
	Proceed()
	HardReject(response string)
	SoftReject(response string)
	Greylist(response string)
	DatalineReply(line string)
	DatalineEnd()
	WriteMultilineHeader(header, value string)
	SafePrintf(format string, params... interface{})
	SafePrintln(msg string)
	FlushMessage(session *SMTPSession)
}

type FilterEventData struct {
	atoms[] string
}

type FilterEvent interface {
	GetAtoms() []string
	GetVerb() string
	GetSessionId() string
	GetToken() string
	GetParams() []string
	Responder() EventResponder
}

type EventResponderImpl struct {
	event FilterEvent
}

type FilterEventImpl struct {
	FilterEventData
}

func (freq FilterEventImpl) GetVerb() string {
	return freq.atoms[4]
}

func (freq FilterEventImpl) GetSessionId() string {
	return freq.atoms[5]
}

func (freq FilterEventImpl) GetToken() string {
	return freq.atoms[6]
}

func (freq FilterEventImpl) GetParams() []string {
	if len(freq.atoms) >= 6 {
		return freq.atoms[6:]
	} else {
		return []string{}
	}
}

func (freq FilterEventImpl) GetAtoms() []string {
	return freq.atoms
}

func (freq *FilterEventImpl) Responder() EventResponder {
	return NewEventResponder(freq)
}

func (evr *EventResponderImpl) Proceed() {
	evr.SafePrintf("filter-result|%s|%s|proceed\n", evr.event.GetToken(), evr.event.GetSessionId())
}

func (evr *EventResponderImpl) HardReject(response string) {
	evr.SafePrintf("filter-result|%s|%s|reject|550 %s\n", evr.event.GetToken(), evr.event.GetSessionId(), response)
}

func (evr *EventResponderImpl) Greylist(response string) {
	evr.SafePrintf("filter-result|%s|%s|reject|421 %s\n", evr.event.GetToken(), evr.event.GetSessionId(), response)
}

func (evr *EventResponderImpl) SoftReject(response string) {
	evr.SafePrintf("filter-result|%s|%s|reject|451 %s\n", evr.event.GetToken(), evr.event.GetSessionId(), response)
}

func (evr *EventResponderImpl) FlushMessage(session *SMTPSession) {
	token := evr.event.GetToken()
	for _, line := range session.Message {
		evr.SafePrintf("filter-dataline|%s|%s|%s\n", token, session.Id, line)
	}
	evr.SafePrintf("filter-dataline|%s|%s|.\n", token, session.Id)
}

func (evr *EventResponderImpl) DatalineEnd() {
	evr.SafePrintf("filter-dataline|%s|%s|.\n", evr.event.GetToken(), evr.event.GetSessionId())
}

func (evr *EventResponderImpl) DatalineReply(line string) {
	prefix := ""
	// Output raw SMTP data - escape leading dots.
	if strings.HasPrefix(line, ".") {
		prefix = "."
	}
	evr.SafePrintf("filter-dataline|%s|%s|%s%s\n", evr.event.GetToken(), evr.event.GetSessionId(), prefix, line)
}

func (evr *EventResponderImpl) WriteMultilineHeader(header, value string) {
	token := evr.event.GetToken()
	sessionId := evr.event.GetSessionId()
	for i, line := range strings.Split(value, "\n") {
		if i == 0 {
			evr.SafePrintf("filter-dataline|%s|%s|%s: %s\n", token, sessionId, header, line)
		} else {
			evr.SafePrintf("filter-dataline|%s|%s|%s\n", token, sessionId, line)
		}
	}
}

func (evr *EventResponderImpl) SafePrintf(format string, params... interface{}) {
	stdoutChannel <- fmt.Sprintf(format, params...)
}

func (evr *EventResponderImpl) SafePrintln(msg string) {
	stdoutChannel <- msg + "\n"
}

func NewEventResponder(_event FilterEvent) EventResponder {
	resp := EventResponderImpl{
		event: _event,
	}

	return &resp
}

func NewFilterEvent(_atoms []string) FilterEvent {
	ev := FilterEventImpl{
		FilterEventData{
			atoms: _atoms,
		},
	}

	return &ev
}

func NewFilter(filter Filter) FilterWrapper {
	return &FilterWrapperImpl{
		Filter: filter,
    }
}

var stdoutChannel = make(chan string)

func stdoutWriter(out <-chan string) {
	for str := range out {
		fmt.Print(str)
	}
}

func Run(fw FilterWrapper) {
	// start the stdout writer goroutine so we can write thread safe
	go stdoutWriter(stdoutChannel)

	scanner := bufio.NewScanner(os.Stdin)

	fw.ProcessConfig(scanner)
	fw.Register(NewEventResponder(NewFilterEvent([]string{})))

	for {
		if !scanner.Scan() {
			log.Printf("Scanner closed")
			os.Exit(0)
		}

		atoms := strings.Split(scanner.Text(), "|")
		if len(atoms) < 6 {
			log.Fatal("Less than 6 atoms")
		}

		fw.Dispatch(atoms)
	}
}
