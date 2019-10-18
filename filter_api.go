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
}

type FilterDef struct {
	Sessions map[string]*SMTPSession
}

type EventHandler = func(FilterDef, string, SessionHolder, string, []string)

type FilterDispatchMap = map[string]map[string]EventHandler

func (fdef *FilterDef) GetCapabilities() FilterDispatchMap {
	capabilities := make(FilterDispatchMap)

	capabilities["report"] = make(map[string]EventHandler)
	reporters := capabilities["report"]
	if _, ok := (interface{})(fdef).(LinkConnectFilter); ok {
		reporters["link-connect"] =
			func(fd FilterDef, verb string, sh SessionHolder, sessionId string, params []string) {
				fd.Filter.(LinkConnectFilter).LinkConnect(verb, sh, sessionId, params)
			}
	}
	if _, ok := fdef.Filter.(LinkDisconnectFilter); ok {
		reporters["link-disconnect"] =
			func(fd FilterDef, verb string, sh SessionHolder, sessionId string, params []string) {
				fd.Filter.(LinkDisconnectFilter).LinkDisconnect(verb, sh, sessionId, params)
			}
	}
	if _, ok := fdef.Filter.(LinkGreetingFilter); ok {
		reporters["link-greeting"] =
			func(fd FilterDef, verb string, sh SessionHolder, sessionId string, params []string) {
				fd.Filter.(LinkGreetingFilter).LinkGreeting(verb, sh, sessionId, params)
			}
	}
	if _, ok := fdef.Filter.(LinkIdentifyFilter); ok {
		reporters["link-identify"] =
			func(fd FilterDef, verb string, sh SessionHolder, sessionId string, params []string) {
				fd.Filter.(LinkIdentifyFilter).LinkIdentity(verb, sh, sessionId, params)
			}
	}
	if _, ok := fdef.Filter.(LinkTLSFilter); ok {
		reporters["link-tls"] =
			func(fd FilterDef, verb string, sh SessionHolder, sessionId string, params []string) {
				fd.Filter.(LinkTLSFilter).LinkTLS(verb, sh, sessionId, params)
			}
	}
	if _, ok := fdef.Filter.(LinkAuthFilter); ok {
		reporters["link-auth"] =
			func(fd FilterDef, verb string, sh SessionHolder, sessionId string, params []string) {
				fd.Filter.(LinkAuthFilter).LinkAuth(verb, sh, sessionId, params)
			}
	}
	if _, ok := fdef.Filter.(TxResetFilter); ok {
		reporters["tx-reset"] =
			func(fd FilterDef, verb string, sh SessionHolder, sessionId string, params []string) {
				fd.Filter.(TxResetFilter).TxReset(verb, sh, sessionId, params)
			}
	}
	if _, ok := fdef.Filter.(TxBeginFilter); ok {
		reporters["tx-begin"] =
			func(fd FilterDef, verb string, sh SessionHolder, sessionId string, params []string) {
				fd.Filter.(TxBeginFilter).TxBegin(verb, sh, sessionId, params)
			}
	}
	if _, ok := fdef.Filter.(TxMailFilter); ok {
		reporters["tx-mail"] =
			func(fd FilterDef, verb string, sh SessionHolder, sessionId string, params []string) {
				fd.Filter.(TxMailFilter).TxMail(verb, sh, sessionId, params)
			}
	}
	if _, ok := fdef.Filter.(TxRcptFilter); ok {
		reporters["tx-rcpt"] =
			func(fd FilterDef, verb string, sh SessionHolder, sessionId string, params []string) {
				fd.Filter.(TxRcptFilter).TxRcpt(verb, sh, sessionId, params)
			}
	}
	if _, ok := fdef.Filter.(TxEnvelopeFilter); ok {
		reporters["tx-envelope"] =
			func(fd FilterDef, verb string, sh SessionHolder, sessionId string, params []string) {
				fd.Filter.(TxEnvelopeFilter).TxEnvelope(verb, sh, sessionId, params)
			}
	}
	if _, ok := fdef.Filter.(TxDataFilter); ok {
		reporters["tx-data"] =
			func(fd FilterDef, verb string, sh SessionHolder, sessionId string, params []string) {
				fd.Filter.(TxDataFilter).TxData(verb, sh, sessionId, params)
			}
	}
	if _, ok := fdef.Filter.(TxCommitFilter); ok {
		reporters["tx-commit"] =
			func(fd FilterDef, verb string, sh SessionHolder, sessionId string, params []string) {
				fd.Filter.(TxCommitFilter).TxCommit(verb, sh, sessionId, params)
			}
	}
	if _, ok := fdef.Filter.(TxRollbackFilter); ok {
		reporters["link-connect"] =
			func(fd FilterDef, verb string, sh SessionHolder, sessionId string, params []string) {
				fd.Filter.(LinkConnectFilter).LinkConnect(verb, sh, sessionId, params)
			}
	}
	if _, ok := fdef.Filter.(ProtocolClientFilter); ok {
		reporters["protocol-client"] =
			func(fd FilterDef, verb string, sh SessionHolder, sessionId string, params []string) {
				fd.Filter.(ProtocolClientFilter).ProtocolClient(verb, sh, sessionId, params)
			}
	}
	if _, ok := fdef.Filter.(ProtocolServerFilter); ok {
		reporters["protocol-server"] =
			func(fd FilterDef, verb string, sh SessionHolder, sessionId string, params []string) {
				fd.Filter.(ProtocolServerFilter).ProtocolServer(verb, sh, sessionId, params)
			}
	}
	if _, ok := fdef.Filter.(TimeoutFilter); ok {
		reporters["timeout"] =
			func(fd FilterDef, verb string, sh SessionHolder, sessionId string, params []string) {
				fd.Filter.(TimeoutFilter).Timeout(verb, sh, sessionId, params)
			}
	}
	capabilities["report"] = reporters

	capabilities["filter"] = make(map[string]EventHandler)
	filters := capabilities["filter"]
	if _, ok := fdef.Filter.(DatalineFilter); ok {
		filters["data-line"] =
			func(fd FilterDef, verb string, sh SessionHolder, sessionId string, params []string) {
				fd.Filter.(DatalineFilter).Dataline(verb, sh, sessionId, params)
			}
	}
	if _, ok := fdef.Filter.(CommitFilter); ok {
		filters["commit"] =
			func(fd FilterDef, verb string, sh SessionHolder, sessionId string, params []string) {
				fd.Filter.(CommitFilter).Commit(verb, sh, sessionId, params)
			}
	}
	if _, ok := fdef.Filter.(FilterResponseFilter); ok {
		filters["filter-response"] =
			func(fd FilterDef, verb string, sh SessionHolder, sessionId string, params []string) {
				fd.Filter.(FilterResponseFilter).FilterResponse(verb, sh, sessionId, params)
			}
	}
	if _, ok := fdef.Filter.(FilterReportFilter); ok {
		filters["filter-report"] =
			func(fd FilterDef, verb string, sh SessionHolder, sessionId string, params []string) {
				fd.Filter.(FilterReportFilter).FilterReport(verb, sh, sessionId, params)
			}
	}
	capabilities["filter"] = filters

	log.Printf("%v", capabilities)

	return capabilities
}

func (fdef *FilterDef) Register() {
	capabilities := fdef.GetCapabilities()
	for typ := range capabilities {
		log.Printf("enumerating type %v", typ)
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


func (fdef *FilterDef) GetSessions() map[string]*SMTPSession {
	return fdef.Sessions
}

func (fdef *FilterDef) GetSession(sessionId string) *SMTPSession {
	return fdef.Sessions[sessionId]
}

func (fdef *FilterDef) SetSession(session *SMTPSession) {
	fdef.Sessions[session.Id] = session
}

type SessionTrackingMixin struct {}

func (sf *SessionTrackingMixin) LinkConnect(sh SessionHolder, sessionId string, params []string) {
	if len(params) != 4 {
		log.Fatal("invalid input, shouldn't happen")
	}

	log.Printf("LinkConnect sf:%v sh:%v", sf, sh)

	s := SMTPSession{}
	s.Id = sessionId
	s.Rdns = params[0]
	s.Src = params[2]

	sh.SetSession(&s)
}

func (sf *SessionTrackingMixin) LinkDisconnect(sh SessionHolder, sessionId string, params []string) {
	if len(params) != 0 {
		log.Fatal("invalid input, shouldn't happen")
	}
	delete(sh.GetSessions(), sessionId)
}

func (sf *SessionTrackingMixin) LinkGreeting(sh SessionHolder, sessionId string, params []string) {
	if len(params) != 1 {
		log.Fatal("invalid input, shouldn't happen")
	}

	s := sh.GetSession(sessionId)
	s.MtaName = params[0]
	sh.SetSession(s)
}

func (sf *SessionTrackingMixin) LinkIdentify(sh SessionHolder, sessionId string, params []string) {
	if len(params) != 2 {
		log.Fatal("invalid input, shouldn't happen")
	}

	s := sh.GetSession(sessionId)
	s.HeloName = params[1]
	sh.SetSession(s)
}

func (sf *SessionTrackingMixin) LinkAuth(sh SessionHolder, sessionId string, params []string) {
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

func (sf *SessionTrackingMixin) TxReset(sh SessionHolder, sessionId string, params []string) {
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

func (sf *SessionTrackingMixin) TxBegin(sh SessionHolder, sessionId string, params []string) {
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

func (sf *SessionTrackingMixin) TxMail(sh SessionHolder, sessionId string, params []string) {
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

func (sf *SessionTrackingMixin) TxRcpt(sh SessionHolder, sessionId string, params []string) {
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

func (sf *SessionTrackingMixin) Dataline(sh SessionHolder, sessionId string, params []string) {
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

func (sf *SessionTrackingMixin) Commit(sh SessionHolder, sessionId string, params []string) {
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

func (fdef *FilterDef) Dispatch(atoms []string) {
	var sh SessionHolder = fdef

	fcap := fdef.GetCapabilities()
	if handler, err := fcap[atoms[0]][atoms[4]]; !err {
		handler(*fdef, atoms[4], sh, atoms[5], atoms[6:])
	}
}

func (fdef *FilterDef) ProcessConfig(scanner *bufio.Scanner) {
	for {
		if !scanner.Scan() {
			os.Exit(0)
		}
		line := scanner.Text()
		if cr, ok := fdef.Filter.(ConfigReceiver); ok {
			cr.Config(strings.Split(line, "|"))
		}

		if line == "config|ready" {
			return
		}
	}
}


func (fdef *FilterDef) Run() {
	scanner := bufio.NewScanner(os.Stdin)

	if fdef.Sessions == nil {
		fdef.Sessions = make(map[string]*SMTPSession)
	}

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
