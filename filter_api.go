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


type SMTPSession struct {
	Id string

	Rdns string
	Src string
	HeloName string
	UserName string
	MtaName string

	Msgid string
	MailFrom string
	RcptTo []string
	Message []string
}


type EventHandler = func(sessionId string, params[] string)
type FilterDispatchMap = map[string]EventHandler

type LinkConnectFilter interface {
	LinkConnect(string, []string)
}

type LinkDisconnectFilter interface {
	LinkDisconnect(string, []string)
}

type LinkGreetingFilter interface {
	LinkGreeting(string, []string)
}

type LinkIdentifyFilter interface {
	LinkIdentity(string, []string)
}

type LinkTLSFilter interface {
	LinkTLS(string, []string)
}

type LinkAuthFilter interface {
	LinkAuth(string, []string)
}

type TxResetFilter interface {
	TxReset(string, []string)
}

type TxBeginFilter interface {
	TxBegin(string, []string)
}

type TxMailFilter interface {
	TxMail(string, []string)
}

type TxRcptFilter interface {
	TxRcpt(string, []string)
}

type TxEnvelopeFilter interface {
	TxEnvelope(string, []string)
}

type TxDataFilter interface {
	TxData(string, []string)
}

type TxCommitFilter interface {
	TxCommit(string, []string)
}

type TxRollbackFilter interface {
	TxRollback(string, []string)
}

type ProtocolClientFilter interface {
	ProtocolClient(string, []string)
}

type ProtocolServerFilter interface {
	ProtocolServer(string, []string)
}

type FilterReportFilter interface {
	FilterReport(string, []string)
}

type FilterResponseFilter interface {
	FilterResponse(string, []string)
}

type TimeoutFilter interface {
	Timeout(string, []string)
}

type DatalineFilter interface {
	Dataline(string, []string)
}

type CommitFilter interface {
	Commit(string, []string)
}

type ConfigReceiver interface {
	Config([]string)
}

type MessageReceivedCallback interface {
	MessageComplete(string, *SMTPSession)
}

type TxBeginCallback interface {
	TxBeginCallback(string, *SMTPSession)
}

func setIfSet(themap map[string]EventHandler, handler EventHandler, key string) {
	if handler != nil {
		themap[key] = handler
	}
}

func GetMapping(filter interface{}) (FilterDispatchMap, FilterDispatchMap) {
	reporters := make(map[string]EventHandler)
	if f, ok := filter.(LinkConnectFilter); ok {
		setIfSet(reporters, f.LinkConnect, "link-connect")
	}
	if f, ok := filter.(LinkDisconnectFilter); ok {
		setIfSet(reporters, f.LinkDisconnect, "link-disconnect")
	}
	if f, ok := filter.(LinkGreetingFilter); ok {
		setIfSet(reporters, f.LinkGreeting, "link-greeting")
	}
	if f, ok := filter.(LinkIdentifyFilter); ok {
		setIfSet(reporters, f.LinkIdentity, "link-identify")
	}
	if f, ok := filter.(LinkTLSFilter); ok {
		setIfSet(reporters, f.LinkTLS, "link-tls")
	}
	if f, ok := filter.(LinkAuthFilter); ok {
		setIfSet(reporters, f.LinkAuth, "link-auth")
	}
	if f, ok := filter.(TxResetFilter); ok {
		setIfSet(reporters, f.TxReset, "tx-reset")
	}
	if f, ok := filter.(TxBeginFilter); ok {
		setIfSet(reporters, f.TxBegin, "tx-begin")
	}
	if f, ok := filter.(TxMailFilter); ok {
		setIfSet(reporters, f.TxMail, "tx-mail")
	}
	if f, ok := filter.(TxRcptFilter); ok {
		setIfSet(reporters, f.TxRcpt, "tx-rcpt")
	}
	if f, ok := filter.(TxEnvelopeFilter); ok {
		setIfSet(reporters, f.TxEnvelope, "tx-envelope")
	}
	if f, ok := filter.(TxDataFilter); ok {
		setIfSet(reporters, f.TxData, "tx-data")
	}
	if f, ok := filter.(TxCommitFilter); ok {
		setIfSet(reporters, f.TxCommit, "tx-commit")
	}
	if f, ok := filter.(TxRollbackFilter); ok {
		setIfSet(reporters, f.TxRollback, "tx-rollback")
	}
	if f, ok := filter.(ProtocolClientFilter); ok {
		setIfSet(reporters, f.ProtocolClient, "protocol-client")
	}
	if f, ok := filter.(ProtocolServerFilter); ok {
		setIfSet(reporters, f.ProtocolServer, "protocol-server")
	}
	if f, ok := filter.(TimeoutFilter); ok {
		setIfSet(reporters, f.Timeout, "timeout")
	}

	filters := make(map[string]EventHandler)
	if f, ok := filter.(DatalineFilter); ok {
		setIfSet(filters, f.Dataline, "data-line")
	}
	if f, ok := filter.(CommitFilter); ok {
		setIfSet(filters, f.Commit, "commit")
	}
	if f, ok := filter.(FilterResponseFilter); ok {
		setIfSet(filters, f.FilterResponse, "filter-response")
	}
	if f, ok := filter.(FilterReportFilter); ok {
		setIfSet(filters, f.FilterReport, "filter-report")
	}

	return reporters, filters
}

func Register(filter interface{}) {
	reporters, filters := GetMapping(filter)
	for k := range reporters {
		fmt.Printf("register|report|smtp-in|%s\n", k)
	}
	for k := range filters {
		fmt.Printf("register|filter|smtp-in|%s\n", k)
	}
	fmt.Println("register|ready")
}

type SessionTrackingFilter struct {
	Sessions map[string]SMTPSession
}

func (sf *SessionTrackingFilter) LinkConnect(sessionId string, params []string) {
	if len(params) != 4 {
		log.Fatal("invalid input, shouldn't happen")
	}

	s := SMTPSession{}
	s.Id = sessionId
	s.Rdns = params[0]
	s.Src = params[2]
	sf.Sessions[s.Id] = s
}

func (sf *SessionTrackingFilter) LinkDisconnect(sessionId string, params []string) {
	if len(params) != 0 {
		log.Fatal("invalid input, shouldn't happen")
	}
	delete(sf.Sessions, sessionId)
}

func (sf *SessionTrackingFilter) LinkGreeting(sessionId string, params []string) {
	if len(params) != 1 {
		log.Fatal("invalid input, shouldn't happen")
	}

	s := sf.Sessions[sessionId]
	s.MtaName = params[0]
	sf.Sessions[s.Id] = s
}

func (sf *SessionTrackingFilter) LinkIdentify(sessionId string, params []string) {
	if len(params) != 2 {
		log.Fatal("invalid input, shouldn't happen")
	}

	s := sf.Sessions[sessionId]
	s.HeloName = params[1]
	sf.Sessions[s.Id] = s
}

func (sf *SessionTrackingFilter) LinkAuth(sessionId string, params []string) {
	if len(params) != 2 {
		log.Fatal("invalid input, shouldn't happen")
	}
	if params[1] != "pass" {
		return
	}
	s := sf.Sessions[sessionId]
	s.UserName = params[0]
	sf.Sessions[s.Id] = s
}

func (sf *SessionTrackingFilter) TxReset(sessionId string, params []string) {
	if len(params) != 1 {
		log.Fatal("invalid input, shouldn't happen")
	}

	s := sf.Sessions[sessionId]
	s.Msgid = ""
	s.MailFrom = ""
	s.RcptTo = nil
	s.Message = nil
	sf.Sessions[s.Id] = s
}

func (sf *SessionTrackingFilter) TxBegin(sessionId string, params []string) {
	if len(params) != 1 {
		log.Fatal("invalid input, shouldn't happen")
	}

	s := sf.Sessions[sessionId]
	s.Msgid = params[0]
	sf.Sessions[s.Id] = s

	if cb, ok := (interface {})(*sf).(TxBeginCallback); ok {
		cb.TxBeginCallback(params[0], &s)
	}
}

func (sf *SessionTrackingFilter) TxMail(sessionId string, params []string) {
	if len(params) != 3 {
		log.Fatal("invalid input, shouldn't happen")
	}

	if params[2] != "ok" {
		return
	}

	s := sf.Sessions[sessionId]
	s.MailFrom = params[1]
	sf.Sessions[s.Id] = s
}

func (sf *SessionTrackingFilter) TxRcpt(sessionId string, params []string) {
	if len(params) != 3 {
		log.Fatal("invalid input, shouldn't happen")
	}

	if params[2] != "ok" {
		return
	}

	s := sf.Sessions[sessionId]
	s.RcptTo = append(s.RcptTo, params[1])
	sf.Sessions[s.Id] = s
}

func (sf *SessionTrackingFilter) Dataline(sessionId string, params []string) {
	if len(params) < 2 {
		log.Fatal("invalid input, shouldn't happen")
	}
	//token := params[0]
	line := strings.Join(params[1:], "")

	s := sf.Sessions[sessionId]
	if line == "." {
		s.Message = append(s.Message, line)
		if cb, ok := (interface {})(*sf).(MessageReceivedCallback); ok {
			s := sf.Sessions[sessionId]
			cb.MessageComplete(params[0], &s)
			sf.Sessions[sessionId] = s
		} else {
			FlushMessage(params[0], sf.Sessions[sessionId])
		}
		return
	}
	s.Message = append(s.Message, line)
	sf.Sessions[sessionId] = s
}

func (sf *SessionTrackingFilter) Commit(sessionId string, params []string) {
	if len(params) != 2 {
		log.Fatal("invalid input, shouldn't happen")
	}
}

func Proceed(token, sessionId string) {
	fmt.Printf("filter-result|%s|%s|proceed\n", token, sessionId)
}

func HardReject(token, sessionId, response string) {
	fmt.Printf("filter-result|%s|%s|reject|550 %s\n", token, sessionId, response)
}

func Greylist(token, sessionId, response string) {
	fmt.Printf("filter-result|%s|%s|reject|421 %s\n", token, sessionId, response)
}

func SoftReject(token, sessionId, response string) {
	fmt.Printf("filter-result|%s|%s|reject|451 %s\n", token, sessionId, response)
}

func FlushMessage(token string, session SMTPSession) {
	for _, line := range session.Message {
		fmt.Printf("filter-dataline|%s|%s|%s\n", token, session.Id, line)
	}
	fmt.Printf("filter-dataline|%s|%s|.\n", token, session.Id)
}

func DatalineReply(token, sessionId, line string) {
	fmt.Printf("filter-dataline|%s|%s|%s\n", token, sessionId, line)
}

func WriteMultilineHeader(token, sessionId, header, value string) {
	for i, line := range strings.Split(value, "\n") {
		if i == 0 {
			fmt.Printf("filter-dataline|%s|%s|%s: %s\n",
				token, sessionId, header, line)
		} else {
			fmt.Printf("filter-dataline|%s|%s|%s\n",
				token, sessionId, line)
		}
	}
}

func Dispatch(mapping FilterDispatchMap, atoms []string) {
	found := false
	for action, handler := range mapping {
		if action == atoms[4] {
			handler(atoms[5], atoms[6:])
			found = true
			break
		}
	}
	if !found {
		log.Printf("Received event for unregistered handler %s", atoms[4])
		os.Exit(1)
	}
}

func ProcessConfig(scanner *bufio.Scanner, filter interface{}) {
	for {
		if !scanner.Scan() {
			os.Exit(0)
		}
		line := scanner.Text()
		if cr, ok := filter.(ConfigReceiver); ok {
			cr.Config(strings.Split(line, "|"))
		}

		if line == "config|ready" {
			return
		}
	}
}

func Run(filter interface{}) {
	scanner := bufio.NewScanner(os.Stdin)
	ProcessConfig(scanner, filter)

	Register(filter)

	for {
		if !scanner.Scan() {
			log.Printf("Scanner closed")
			os.Exit(0)
		}

		atoms := strings.Split(scanner.Text(), "|")
		if len(atoms) < 6 {
			log.Printf("Less than 6 atoms")
			os.Exit(1)
		}

		reporters, filters := GetMapping(filter)

		switch atoms[0] {
		case "report":
			Dispatch(reporters, atoms)
		case "filter":
			Dispatch(filters, atoms)
		default:
			log.Printf("No matching handler type %s", atoms[0])
			os.Exit(1)
		}
	}
}
