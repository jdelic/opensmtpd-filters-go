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


type LinkConnectFilter interface {
	LinkConnect(string, SessionHolder, string, []string)
}

type LinkDisconnectFilter interface {
	LinkDisconnect(string, SessionHolder, string, []string)
}

type LinkGreetingFilter interface {
	LinkGreeting(string, SessionHolder, string, []string)
}

type LinkIdentifyFilter interface {
	LinkIdentity(string, SessionHolder, string, []string)
}

type LinkTLSFilter interface {
	LinkTLS(string, SessionHolder, string, []string)
}

type LinkAuthFilter interface {
	LinkAuth(string, SessionHolder, string, []string)
}

type TxResetFilter interface {
	TxReset(string, SessionHolder, string, []string)
}

type TxBeginFilter interface {
	TxBegin(string, SessionHolder, string, []string)
}

type TxMailFilter interface {
	TxMail(string, SessionHolder, string, []string)
}

type TxRcptFilter interface {
	TxRcpt(string, SessionHolder, string, []string)
}

type TxEnvelopeFilter interface {
	TxEnvelope(string, SessionHolder, string, []string)
}

type TxDataFilter interface {
	TxData(string, SessionHolder, string, []string)
}

type TxCommitFilter interface {
	TxCommit(string, SessionHolder, string, []string)
}

type TxRollbackFilter interface {
	TxRollback(string, SessionHolder, string, []string)
}

type ProtocolClientFilter interface {
	ProtocolClient(string, SessionHolder, string, []string)
}

type ProtocolServerFilter interface {
	ProtocolServer(string, SessionHolder, string, []string)
}

type FilterReportFilter interface {
	FilterReport(string, SessionHolder, string, []string)
}

type FilterResponseFilter interface {
	FilterResponse(string, SessionHolder, string, []string)
}

type TimeoutFilter interface {
	Timeout(string, SessionHolder, string, []string)
}

type DatalineFilter interface {
	Dataline(string, SessionHolder, string, []string)
}

type CommitFilter interface {
	Commit(string, SessionHolder, string, []string)
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

type Filter interface {
	GetCapabilities() FilterDispatchMap
	Register()
	Dispatch([]string)
	ProcessConfig(*bufio.Scanner)
}

type EventHandler = func(Filter, string, SessionHolder, string, []string)

type FilterDispatchMap = map[string]map[string]EventHandler

func GetCapabilities(fdef Filter) FilterDispatchMap {
	capabilities := make(FilterDispatchMap)

	capabilities["report"] = make(map[string]EventHandler)
	reporters := capabilities["report"]
	if _, ok := fdef.(LinkConnectFilter); ok {
		reporters["link-connect"] =
			func(fd Filter, verb string, sh SessionHolder, sessionId string, params []string) {
				fd.(LinkConnectFilter).LinkConnect(verb, sh, sessionId, params)
			}
	}
	if _, ok := fdef.(LinkDisconnectFilter); ok {
		reporters["link-disconnect"] =
			func(fd Filter, verb string, sh SessionHolder, sessionId string, params []string) {
				fd.(LinkDisconnectFilter).LinkDisconnect(verb, sh, sessionId, params)
			}
	}
	if _, ok := fdef.(LinkGreetingFilter); ok {
		reporters["link-greeting"] =
			func(fd Filter, verb string, sh SessionHolder, sessionId string, params []string) {
				fd.(LinkGreetingFilter).LinkGreeting(verb, sh, sessionId, params)
			}
	}
	if _, ok := fdef.(LinkIdentifyFilter); ok {
		reporters["link-identify"] =
			func(fd Filter, verb string, sh SessionHolder, sessionId string, params []string) {
				fd.(LinkIdentifyFilter).LinkIdentity(verb, sh, sessionId, params)
			}
	}
	if _, ok := fdef.(LinkTLSFilter); ok {
		reporters["link-tls"] =
			func(fd Filter, verb string, sh SessionHolder, sessionId string, params []string) {
				fd.(LinkTLSFilter).LinkTLS(verb, sh, sessionId, params)
			}
	}
	if _, ok := fdef.(LinkAuthFilter); ok {
		reporters["link-auth"] =
			func(fd Filter, verb string, sh SessionHolder, sessionId string, params []string) {
				fd.(LinkAuthFilter).LinkAuth(verb, sh, sessionId, params)
			}
	}
	if _, ok := fdef.(TxResetFilter); ok {
		reporters["tx-reset"] =
			func(fd Filter, verb string, sh SessionHolder, sessionId string, params []string) {
				fd.(TxResetFilter).TxReset(verb, sh, sessionId, params)
			}
	}
	if _, ok := fdef.(TxBeginFilter); ok {
		reporters["tx-begin"] =
			func(fd Filter, verb string, sh SessionHolder, sessionId string, params []string) {
				fd.(TxBeginFilter).TxBegin(verb, sh, sessionId, params)
			}
	}
	if _, ok := fdef.(TxMailFilter); ok {
		reporters["tx-mail"] =
			func(fd Filter, verb string, sh SessionHolder, sessionId string, params []string) {
				fd.(TxMailFilter).TxMail(verb, sh, sessionId, params)
			}
	}
	if _, ok := fdef.(TxRcptFilter); ok {
		reporters["tx-rcpt"] =
			func(fd Filter, verb string, sh SessionHolder, sessionId string, params []string) {
				fd.(TxRcptFilter).TxRcpt(verb, sh, sessionId, params)
			}
	}
	if _, ok := fdef.(TxEnvelopeFilter); ok {
		reporters["tx-envelope"] =
			func(fd Filter, verb string, sh SessionHolder, sessionId string, params []string) {
				fd.(TxEnvelopeFilter).TxEnvelope(verb, sh, sessionId, params)
			}
	}
	if _, ok := fdef.(TxDataFilter); ok {
		reporters["tx-data"] =
			func(fd Filter, verb string, sh SessionHolder, sessionId string, params []string) {
				fd.(TxDataFilter).TxData(verb, sh, sessionId, params)
			}
	}
	if _, ok := fdef.(TxCommitFilter); ok {
		reporters["tx-commit"] =
			func(fd Filter, verb string, sh SessionHolder, sessionId string, params []string) {
				fd.(TxCommitFilter).TxCommit(verb, sh, sessionId, params)
			}
	}
	if _, ok := fdef.(TxRollbackFilter); ok {
		reporters["link-connect"] =
			func(fd Filter, verb string, sh SessionHolder, sessionId string, params []string) {
				fd.(LinkConnectFilter).LinkConnect(verb, sh, sessionId, params)
			}
	}
	if _, ok := fdef.(ProtocolClientFilter); ok {
		reporters["protocol-client"] =
			func(fd Filter, verb string, sh SessionHolder, sessionId string, params []string) {
				fd.(ProtocolClientFilter).ProtocolClient(verb, sh, sessionId, params)
			}
	}
	if _, ok := fdef.(ProtocolServerFilter); ok {
		reporters["protocol-server"] =
			func(fd Filter, verb string, sh SessionHolder, sessionId string, params []string) {
				fd.(ProtocolServerFilter).ProtocolServer(verb, sh, sessionId, params)
			}
	}
	if _, ok := fdef.(TimeoutFilter); ok {
		reporters["timeout"] =
			func(fd Filter, verb string, sh SessionHolder, sessionId string, params []string) {
				fd.(TimeoutFilter).Timeout(verb, sh, sessionId, params)
			}
	}
	capabilities["report"] = reporters

	capabilities["filter"] = make(map[string]EventHandler)
	filters := capabilities["filter"]
	if _, ok := fdef.(DatalineFilter); ok {
		filters["data-line"] =
			func(fd Filter, verb string, sh SessionHolder, sessionId string, params []string) {
				fd.(DatalineFilter).Dataline(verb, sh, sessionId, params)
			}
	}
	if _, ok := fdef.(CommitFilter); ok {
		filters["commit"] =
			func(fd Filter, verb string, sh SessionHolder, sessionId string, params []string) {
				fd.(CommitFilter).Commit(verb, sh, sessionId, params)
			}
	}
	if _, ok := fdef.(FilterResponseFilter); ok {
		filters["filter-response"] =
			func(fd Filter, verb string, sh SessionHolder, sessionId string, params []string) {
				fd.(FilterResponseFilter).FilterResponse(verb, sh, sessionId, params)
			}
	}
	if _, ok := fdef.(FilterReportFilter); ok {
		filters["filter-report"] =
			func(fd Filter, verb string, sh SessionHolder, sessionId string, params []string) {
				fd.(FilterReportFilter).FilterReport(verb, sh, sessionId, params)
			}
	}
	capabilities["filter"] = filters
	return capabilities
}

func Register(fdef Filter) {
	capabilities := fdef.GetCapabilities()
	for typ := range capabilities {
		for op := range capabilities[typ] {
			fmt.Printf("register|%v|smtp-in|%v\n", typ, op)
		}
	}
	fmt.Println("register|ready")
}

type SessionHolder interface {
	GetSessions() map[string]*SMTPSession
	GetSession(string) *SMTPSession
	SetSession(*SMTPSession)
}

type SessionHolderImpl struct {
	Sessions map[string]*SMTPSession
}

func (shi *SessionHolderImpl) GetSessions() map[string]*SMTPSession {
	return shi.Sessions
}

func (shi *SessionHolderImpl) GetSession(sessionId string) *SMTPSession {
	if shi.Sessions == nil {
		return nil
	}
	return shi.Sessions[sessionId]
}

func (shi *SessionHolderImpl) SetSession(session *SMTPSession) {
	if shi.Sessions == nil {
		shi.Sessions = make(map[string]*SMTPSession)
	}
	shi.Sessions[session.Id] = session
}

type SessionTrackingMixin struct {
	SessionHolderImpl
}

func (sf *SessionTrackingMixin) LinkConnect(verb string, sh SessionHolder, sessionId string, params []string) {
	if len(params) != 4 {
		log.Fatal("invalid input, shouldn't happen")
	}

	s := SMTPSession{}
	s.Id = sessionId
	s.Rdns = params[0]
	s.Src = params[2]

	sh.SetSession(&s)
}

func (sf *SessionTrackingMixin) LinkDisconnect(verb string, sh SessionHolder, sessionId string, params []string) {
	if len(params) != 0 {
		log.Fatal("invalid input, shouldn't happen")
	}
	delete(sh.GetSessions(), sessionId)
}

func (sf *SessionTrackingMixin) LinkGreeting(verb string, sh SessionHolder, sessionId string, params []string) {
	if len(params) != 1 {
		log.Fatal("invalid input, shouldn't happen")
	}

	s := sh.GetSession(sessionId)
	s.MtaName = params[0]
	sh.SetSession(s)
}

func (sf *SessionTrackingMixin) LinkIdentify(verb string, sh SessionHolder, sessionId string, params []string) {
	if len(params) != 2 {
		log.Fatal("invalid input, shouldn't happen")
	}

	s := sh.GetSession(sessionId)
	s.HeloName = params[1]
	sh.SetSession(s)
}

func (sf *SessionTrackingMixin) LinkAuth(verb string, sh SessionHolder, sessionId string, params []string) {
	if len(params) != 2 {
		log.Fatal("invalid input, shouldn't happen")
	}
	if params[1] != "pass" {
		return
	}
	s := sh.GetSession(sessionId)
	s.UserName = params[0]
	sh.SetSession(s)
}

func (sf *SessionTrackingMixin) TxReset(verb string, sh SessionHolder, sessionId string, params []string) {
	if len(params) != 1 {
		log.Fatal("invalid input, shouldn't happen")
	}

	s := sh.GetSession(sessionId)
	s.Msgid = ""
	s.MailFrom = ""
	s.RcptTo = nil
	s.Message = nil
	sh.SetSession(s)
}

func (sf *SessionTrackingMixin) TxBegin(verb string, sh SessionHolder, sessionId string, params []string) {
	if len(params) != 1 {
		log.Fatal("invalid input, shouldn't happen")
	}

	s := sh.GetSession(sessionId)
	s.Msgid = params[0]
	sh.SetSession(s)

	if cb, ok := (interface {})(*sf).(TxBeginCallback); ok {
		cb.TxBeginCallback(params[0], s)
	}
}

func (sf *SessionTrackingMixin) TxMail(verb string, sh SessionHolder, sessionId string, params []string) {
	if len(params) != 3 {
		log.Fatal("invalid input, shouldn't happen")
	}

	if params[2] != "ok" {
		return
	}

	s := sh.GetSession(sessionId)
	s.MailFrom = params[1]
	sh.SetSession(s)
}

func (sf *SessionTrackingMixin) TxRcpt(verb string, sh SessionHolder, sessionId string, params []string) {
	if len(params) != 3 {
		log.Fatal("invalid input, shouldn't happen")
	}

	if params[2] != "ok" {
		return
	}

	s := sh.GetSession(sessionId)
	s.RcptTo = append(s.RcptTo, params[1])
	sh.SetSession(s)
}

func (sf *SessionTrackingMixin) Dataline(verb string, sh SessionHolder, sessionId string, params []string) {
	if len(params) < 2 {
		log.Fatal("invalid input, shouldn't happen")
	}
	//token := params[0]
	line := strings.Join(params[1:], "")

	s := sh.GetSession(sessionId)
	if line == "." {
		s.Message = append(s.Message, line)
		if cb, ok := (interface {})(*sf).(MessageReceivedCallback); ok {
			s := sh.GetSession(sessionId)
			cb.MessageComplete(params[0], s)
			sh.SetSession(s)
		} else {
			FlushMessage(params[0], sh.GetSession(sessionId))
		}
		return
	}
	s.Message = append(s.Message, line)
	sh.SetSession(s)
}

func (sf *SessionTrackingMixin) Commit(verb string, sh SessionHolder, sessionId string, params []string) {
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

func FlushMessage(token string, session *SMTPSession) {
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

func Dispatch(fdef Filter, atoms []string) {
	var sh SessionHolder
	var ok bool
	if sh, ok = fdef.(SessionHolder); !ok {
		sh = nil
	}

	fcap := fdef.GetCapabilities()
	log.Printf("type: %v op: %v", atoms[0], atoms[4])
	if handler, ok := fcap[atoms[0]][atoms[4]]; ok {
		handler(fdef, atoms[4], sh, atoms[5], atoms[6:])
	}
}

func ProcessConfig(fdef Filter, scanner *bufio.Scanner) {
	for {
		if !scanner.Scan() {
			os.Exit(0)
		}
		line := scanner.Text()
		if cr, ok := fdef.(ConfigReceiver); ok {
			cr.Config(strings.Split(line, "|"))
		}

		if line == "config|ready" {
			return
		}
	}
}


func Run(fdef Filter) {
	scanner := bufio.NewScanner(os.Stdin)

	fdef.ProcessConfig(scanner)
	fdef.Register()

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

		fdef.Dispatch(atoms)
	}
}
