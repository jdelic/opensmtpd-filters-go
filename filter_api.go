//
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
	id string

	rdns string
	src string
	heloName string
	userName string
	mtaName string

	msgid string
	mailFrom string
	rcptTo []string
	message []string
}

/*
    { "link-connect" },
    { "link-disconnect" },
    { "link-greeting" },
    { "link-identify" },
    { "link-tls" },
    { "link-auth" },

    { "tx-reset" },
    { "tx-begin" },
    { "tx-mail" },
    { "tx-rcpt" },
    { "tx-envelope" },
    { "tx-data" },
    { "tx-commit" },
    { "tx-rollback" },

    { "protocol-client" },
    { "protocol-server" },

    { "filter-report" },
    { "filter-response" },

    { "timeout" },

	var reporters = []string {
		"link-connect":    linkConnect,
		"link-disconnect": linkDisconnect,
		"link-greeting":   linkGreeting,
		"link-identify":   linkIdentify,
		"link-auth":       linkAuth,
		"tx-reset":        txReset,
		"tx-begin":        txBegin,
		"tx-mail":         txMail,
		"tx-rcpt":         txRcpt,
	}


	var filters = map[string]func(string, []string) {
		"data-line": dataLine,
		"commit": dataCommit,
	}
*/


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

type MessageReceivedCallback interface {
	MessageComplete(string, *SMTPSession)
}

type CommitFilter interface {
	Commit(string, []string)
}

type ConfigReceiver interface {
	Config([]string)
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
	s.id = sessionId
	s.rdns = params[0]
	s.src = params[2]
	sf.Sessions[s.id] = s
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
	s.mtaName = params[0]
	sf.Sessions[s.id] = s
}

func (sf *SessionTrackingFilter) LinkIdentify(sessionId string, params []string) {
	if len(params) != 2 {
		log.Fatal("invalid input, shouldn't happen")
	}

	s := sf.Sessions[sessionId]
	s.heloName = params[1]
	sf.Sessions[s.id] = s
}

func (sf *SessionTrackingFilter) LinkAuth(sessionId string, params []string) {
	if len(params) != 2 {
		log.Fatal("invalid input, shouldn't happen")
	}
	if params[1] != "pass" {
		return
	}
	s := sf.Sessions[sessionId]
	s.userName = params[0]
	sf.Sessions[s.id] = s
}

func (sf *SessionTrackingFilter) TxReset(sessionId string, params []string) {
	if len(params) != 1 {
		log.Fatal("invalid input, shouldn't happen")
	}

	s := sf.Sessions[sessionId]
	s.msgid = ""
	s.mailFrom = ""
	s.rcptTo = nil
	s.message = nil
	sf.Sessions[s.id] = s
}

func (sf *SessionTrackingFilter) TxBegin(sessionId string, params []string) {
	if len(params) != 1 {
		log.Fatal("invalid input, shouldn't happen")
	}

	s := sf.Sessions[sessionId]
	s.msgid = params[0]
	sf.Sessions[s.id] = s
}

func (sf *SessionTrackingFilter) TxMail(sessionId string, params []string) {
	if len(params) != 3 {
		log.Fatal("invalid input, shouldn't happen")
	}

	if params[2] != "ok" {
		return
	}

	s := sf.Sessions[sessionId]
	s.mailFrom = params[1]
	sf.Sessions[s.id] = s
}

func (sf *SessionTrackingFilter) TxRcpt(sessionId string, params []string) {
	if len(params) != 3 {
		log.Fatal("invalid input, shouldn't happen")
	}

	if params[2] != "ok" {
		return
	}

	s := sf.Sessions[sessionId]
	s.rcptTo = append(s.rcptTo, params[1])
	sf.Sessions[s.id] = s
}

func (sf *SessionTrackingFilter) Dataline(sessionId string, params []string) {
	if len(params) < 2 {
		log.Fatal("invalid input, shouldn't happen")
	}
	//token := params[0]
	line := strings.Join(params[1:], "|")

	s := sf.Sessions[sessionId]
	if line == "." {
		if cb, ok := (interface {})(*sf).(MessageReceivedCallback); ok {
			s := sf.Sessions[sessionId]
			cb.MessageComplete(params[0], &s)
			sf.Sessions[sessionId] = s
		} else {
			FlushMessage(params[0], sf.Sessions[sessionId])
		}
		return
	}
	s.message = append(s.message, line)
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
	for _, line := range session.message {
		fmt.Printf("filter-dataline|%s|%s|%s\n", token, session.id, line)
	}
	fmt.Printf("filter-dataline|%s|%s|.\n", token, session.id)
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
			os.Exit(0)
		}

		atoms := strings.Split(scanner.Text(), "|")
		if len(atoms) < 6 {
			os.Exit(1)
		}

		reporters, filters := GetMapping(filter)

		switch atoms[0] {
		case "report":
			Dispatch(reporters, atoms)
		case "filter":
			Dispatch(filters, atoms)
		default:
			os.Exit(1)
		}
	}
}

/*
func rspamdQuery(s session, token string) {
	r := strings.NewReader(strings.Join(s.message, "\n"))
	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/checkv2", *rspamdURL), r)
	if err != nil {
		flushMessage(s, token)
		return
	}

	req.Header.Add("Pass", "All")
	if !strings.HasPrefix(s.src, "unix:") {
		if s.src[0] == '[' {
			ip := strings.Split(strings.Split(s.src, "]")[0], "[")[1]
			req.Header.Add("Ip", ip)
		} else {
			ip := strings.Split(s.src, ":")[0]
			req.Header.Add("Ip", ip)
		}
	} else {
		req.Header.Add("Ip", "127.0.0.1")
	}

	req.Header.Add("Hostname", s.rdns)
	req.Header.Add("Helo", s.heloName)
	req.Header.Add("MTA-Name", s.mtaName)
	req.Header.Add("Queue-Id", s.msgid)
	req.Header.Add("From", s.mailFrom)

	if s.userName != "" {
		req.Header.Add("User", s.userName)
	}

	for _, rcptTo := range s.rcptTo {
		req.Header.Add("Rcpt", rcptTo)
	}

	resp, err := client.Do(req)
	if err != nil {
		flushMessage(s, token)
		return
	}
	defer resp.Body.Close()

	rr := &rspamd{}
	if err := json.NewDecoder(resp.Body).Decode(rr); err != nil {
		flushMessage(s, token)
		return
	}

	switch rr.Action {
	case "reject":
		fallthrough
	case "greylist":
		fallthrough
	case "soft reject":
		s.action = rr.Action
		s.response = rr.Messages.SMTP
		sessions[s.id] = s
		flushMessage(s, token)
		return
	}

	if rr.DKIMSig != "" {
		WriteMultilineHeader(s, token, "DKIM-Signature", rr.DKIMSig)
	}

	if rr.Action == "add header" {
		fmt.Printf("filter-dataline|%s|%s|%s: %s\n",
			token, s.id, "X-Spam", "yes")
		fmt.Printf("filter-dataline|%s|%s|%s: %s\n",
			token, s.id, "X-Spam-Score",
			fmt.Sprintf("%v / %v",
				rr.Score, rr.RequiredScore))
	}

	inhdr := true
	for _, line := range s.message {
		if line == "" {
			inhdr = false
		}
		if rr.Action == "rewrite subject" && inhdr && strings.HasPrefix(line, "Subject: ") {
			fmt.Printf("filter-dataline|%s|%s|Subject: %s\n", token, s.id, rr.Subject)
		} else {
			fmt.Printf("filter-dataline|%s|%s|%s\n", token, s.id, line)
		}
	}
	fmt.Printf("filter-dataline|%s|%s|.\n", token, s.id)
	sessions[s.id] = s
}
 */
