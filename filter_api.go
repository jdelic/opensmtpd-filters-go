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

type FilterEventData struct {
	atoms []string
}

type FilterEvent interface {
	GetAtoms() []string
	GetProtocolVersion() string
	GetVerb() string
	GetSessionId() string
	GetToken() string
	GetParams() []string
	Responder() EventResponder
}

type FilterEventImpl struct {
	FilterEventData
}

func (freq FilterEventImpl) GetProtocolVersion() string {
	return freq.atoms[1]
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

type SafePrinter struct{}

func (sp *SafePrinter) SafePrintln(msg string) {
	stdoutChannel <- msg + "\n"
}

func Run(fw FilterWrapper) {
	// start the stdout writer goroutine so we can write thread safe
	go stdoutWriter(stdoutChannel)

	scanner := bufio.NewScanner(os.Stdin)

	fw.ProcessConfig(scanner)
	fw.Register(NewEventResponder(NewFilterEvent([]string{})))

	for {
		if !scanner.Scan() {
			log.Println("Scanner closed")
			os.Exit(0)
		}

		atoms := strings.Split(scanner.Text(), "|")
		if len(atoms) < 6 {
			log.Fatal("Less than 6 atoms")
		}

		fw.Dispatch(atoms)
	}
}
