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


type FilterWrapper interface {
	GetCapabilities() FilterDispatchMap
	Register()
	Dispatch([]string)
	ProcessConfig(*bufio.Scanner)
	GetFilter() interface{}
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

type FilterWrapperImpl struct {
	Filter interface{}
}

func (fw *FilterWrapperImpl) GetFilter() interface{} {
	return fw.Filter
}

type EventHandler = func(string, SessionHolder, string, []string)

type FilterDispatchMap = map[string]map[string]EventHandler

func (fw *FilterWrapperImpl) GetCapabilities() FilterDispatchMap {
	capabilities := make(FilterDispatchMap)

	capabilities["report"] = make(map[string]EventHandler)
	reportReceivers := capabilities["report"]
	if _, ok := fw.Filter.(LinkConnectReceiver); ok {
		reportReceivers["link-connect"] =
			func(verb string, sh SessionHolder, sessionId string, params []string) {
				fw.Filter.(LinkConnectReceiver).LinkConnect(verb, sh, sessionId, params)
			}
	}
	if _, ok := fw.Filter.(LinkDisconnectReceiver); ok {
		reportReceivers["link-disconnect"] =
			func(verb string, sh SessionHolder, sessionId string, params []string) {
				fw.Filter.(LinkDisconnectReceiver).LinkDisconnect(verb, sh, sessionId, params)
			}
	}
	if _, ok := fw.Filter.(LinkGreetingReceiver); ok {
		reportReceivers["link-greeting"] =
			func(verb string, sh SessionHolder, sessionId string, params []string) {
				fw.Filter.(LinkGreetingReceiver).LinkGreeting(verb, sh, sessionId, params)
			}
	}
	if _, ok := fw.Filter.(LinkIdentifyReceiver); ok {
		reportReceivers["link-identify"] =
			func(verb string, sh SessionHolder, sessionId string, params []string) {
				fw.Filter.(LinkIdentifyReceiver).LinkIdentity(verb, sh, sessionId, params)
			}
	}
	if _, ok := fw.Filter.(LinkTLSReceiver); ok {
		reportReceivers["link-tls"] =
			func(verb string, sh SessionHolder, sessionId string, params []string) {
				fw.Filter.(LinkTLSReceiver).LinkTLS(verb, sh, sessionId, params)
			}
	}
	if _, ok := fw.Filter.(LinkAuthReceiver); ok {
		reportReceivers["link-auth"] =
			func(verb string, sh SessionHolder, sessionId string, params []string) {
				fw.Filter.(LinkAuthReceiver).LinkAuth(verb, sh, sessionId, params)
			}
	}
	if _, ok := fw.Filter.(TxResetReceiver); ok {
		reportReceivers["tx-reset"] =
			func(verb string, sh SessionHolder, sessionId string, params []string) {
				fw.Filter.(TxResetReceiver).TxReset(verb, sh, sessionId, params)
			}
	}
	if _, ok := fw.Filter.(TxBeginReceiver); ok {
		reportReceivers["tx-begin"] =
			func(verb string, sh SessionHolder, sessionId string, params []string) {
				fw.Filter.(TxBeginReceiver).TxBegin(verb, sh, sessionId, params)
			}
	}
	if _, ok := fw.Filter.(TxMailReceiver); ok {
		reportReceivers["tx-mail"] =
			func(verb string, sh SessionHolder, sessionId string, params []string) {
				fw.Filter.(TxMailReceiver).TxMail(verb, sh, sessionId, params)
			}
	}
	if _, ok := fw.Filter.(TxRcptReceiver); ok {
		reportReceivers["tx-rcpt"] =
			func(verb string, sh SessionHolder, sessionId string, params []string) {
				fw.Filter.(TxRcptReceiver).TxRcpt(verb, sh, sessionId, params)
			}
	}
	if _, ok := fw.Filter.(TxEnvelopeReceiver); ok {
		reportReceivers["tx-envelope"] =
			func(verb string, sh SessionHolder, sessionId string, params []string) {
				fw.Filter.(TxEnvelopeReceiver).TxEnvelope(verb, sh, sessionId, params)
			}
	}
	if _, ok := fw.Filter.(TxDataReceiver); ok {
		reportReceivers["tx-data"] =
			func(verb string, sh SessionHolder, sessionId string, params []string) {
				fw.Filter.(TxDataReceiver).TxData(verb, sh, sessionId, params)
			}
	}
	if _, ok := fw.Filter.(TxCommitReceiver); ok {
		reportReceivers["tx-commit"] =
			func(verb string, sh SessionHolder, sessionId string, params []string) {
				fw.Filter.(TxCommitReceiver).TxCommit(verb, sh, sessionId, params)
			}
	}
	if _, ok := fw.Filter.(TxRollbackReceiver); ok {
		reportReceivers["link-connect"] =
			func(verb string, sh SessionHolder, sessionId string, params []string) {
				fw.Filter.(LinkConnectReceiver).LinkConnect(verb, sh, sessionId, params)
			}
	}
	if _, ok := fw.Filter.(ProtocolClientReceiver); ok {
		reportReceivers["protocol-client"] =
			func(verb string, sh SessionHolder, sessionId string, params []string) {
				fw.Filter.(ProtocolClientReceiver).ProtocolClient(verb, sh, sessionId, params)
			}
	}
	if _, ok := fw.Filter.(ProtocolServerReceiver); ok {
		reportReceivers["protocol-server"] =
			func(verb string, sh SessionHolder, sessionId string, params []string) {
				fw.Filter.(ProtocolServerReceiver).ProtocolServer(verb, sh, sessionId, params)
			}
	}
	if _, ok := fw.Filter.(TimeoutReceiver); ok {
		reportReceivers["timeout"] =
			func(verb string, sh SessionHolder, sessionId string, params []string) {
				fw.Filter.(TimeoutReceiver).Timeout(verb, sh, sessionId, params)
			}
	}
	if _, ok := fw.Filter.(FilterResponseReceiver); ok {
		reportReceivers["filter-response"] =
			func(verb string, sh SessionHolder, sessionId string, params []string) {
				fw.Filter.(FilterResponseReceiver).FilterResponse(verb, sh, sessionId, params)
			}
	}
	if _, ok := fw.Filter.(FilterReportReceiver); ok {
		reportReceivers["filter-report"] =
			func(verb string, sh SessionHolder, sessionId string, params []string) {
				fw.Filter.(FilterReportReceiver).FilterReport(verb, sh, sessionId, params)
			}
	}
	capabilities["report"] = reportReceivers

	capabilities["filter"] = make(map[string]EventHandler)
	filters := capabilities["filter"]
	if _, ok := fw.Filter.(ConnectFilter); ok {
		filters["connect"] =
			func(verb string, sh SessionHolder, sessionId string, params []string) {
				fw.Filter.(ConnectFilter).Connect(verb, sh, sessionId, params)
			}
	}
	if _, ok := fw.Filter.(HeloFilter); ok {
		filters["helo"] =
			func(verb string, sh SessionHolder, sessionId string, params []string) {
				fw.Filter.(HeloFilter).Helo(verb, sh, sessionId, params)
			}
	}
	if _, ok := fw.Filter.(EhloFilter); ok {
		filters["ehlo"] =
			func(verb string, sh SessionHolder, sessionId string, params []string) {
				fw.Filter.(EhloFilter).Ehlo(verb, sh, sessionId, params)
			}
	}
	if _, ok := fw.Filter.(StartTLSFilter); ok {
		filters["starttls"] =
			func(verb string, sh SessionHolder, sessionId string, params []string) {
				fw.Filter.(StartTLSFilter).StartTLS(verb, sh, sessionId, params)
			}
	}
	if _, ok := fw.Filter.(AuthFilter); ok {
		filters["auth"] =
			func(verb string, sh SessionHolder, sessionId string, params []string) {
				fw.Filter.(AuthFilter).Auth(verb, sh, sessionId, params)
			}
	}
	if _, ok := fw.Filter.(MailFromFilter); ok {
		filters["mail-from"] =
			func(verb string, sh SessionHolder, sessionId string, params []string) {
				fw.Filter.(MailFromFilter).MailFrom(verb, sh, sessionId, params)
			}
	}
	if _, ok := fw.Filter.(RcptToFilter); ok {
		filters["rcpt-to"] =
			func(verb string, sh SessionHolder, sessionId string, params []string) {
				fw.Filter.(RcptToFilter).RcptTo(verb, sh, sessionId, params)
			}
	}
	if _, ok := fw.Filter.(DataFilter); ok {
		filters["data"] =
			func(verb string, sh SessionHolder, sessionId string, params []string) {
				fw.Filter.(DataFilter).Data(verb, sh, sessionId, params)
			}
	}
	if _, ok := fw.Filter.(DatalineFilter); ok {
		filters["data-line"] =
			func(verb string, sh SessionHolder, sessionId string, params []string) {
				fw.Filter.(DatalineFilter).Dataline(verb, sh, sessionId, params)
			}
	}
	if _, ok := fw.Filter.(RsetFilter); ok {
		filters["rset"] =
			func(verb string, sh SessionHolder, sessionId string, params []string) {
				fw.Filter.(RsetFilter).Rset(verb, sh, sessionId, params)
			}
	}
	if _, ok := fw.Filter.(QuitFilter); ok {
		filters["quit"] =
			func(verb string, sh SessionHolder, sessionId string, params []string) {
				fw.Filter.(QuitFilter).Quit(verb, sh, sessionId, params)
			}
	}
	if _, ok := fw.Filter.(NoopFilter); ok {
		filters["noop"] =
			func(verb string, sh SessionHolder, sessionId string, params []string) {
				fw.Filter.(NoopFilter).Noop(verb, sh, sessionId, params)
			}
	}
	if _, ok := fw.Filter.(HelpFilter); ok {
		filters["help"] =
			func(verb string, sh SessionHolder, sessionId string, params []string) {
				fw.Filter.(HelpFilter).Help(verb, sh, sessionId, params)
			}
	}
	if _, ok := fw.Filter.(WizFilter); ok {
		filters["wiz"] =
			func(verb string, sh SessionHolder, sessionId string, params []string) {
				fw.Filter.(WizFilter).Wiz(verb, sh, sessionId, params)
			}
	}
	if _, ok := fw.Filter.(CommitFilter); ok {
		filters["commit"] =
			func(verb string, sh SessionHolder, sessionId string, params []string) {
				fw.Filter.(CommitFilter).Commit(verb, sh, sessionId, params)
			}
	}

	capabilities["filter"] = filters
	return capabilities
}

func (fw *FilterWrapperImpl) Register() {
	capabilities := fw.GetCapabilities()
	for typ := range capabilities {
		for op := range capabilities[typ] {
			fmt.Printf("register|%v|smtp-in|%v\n", typ, op)
		}
	}
	fmt.Println("register|ready")
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

func (fw *FilterWrapperImpl) Dispatch(atoms []string) {
	var sh SessionHolder
	var ok bool
	if sh, ok = fw.Filter.(SessionHolder); !ok {
		sh = nil
	}

	fcap := fw.GetCapabilities()
	log.Printf("type: %v op: %v", atoms[0], atoms[4])
	if handler, ok := fcap[atoms[0]][atoms[4]]; ok {
		handler(atoms[4], sh, atoms[5], atoms[6:])
	}
}

func (fw *FilterWrapperImpl) ProcessConfig(scanner *bufio.Scanner) {
	for {
		if !scanner.Scan() {
			os.Exit(0)
		}
		line := scanner.Text()
		if cr, ok := fw.Filter.(ConfigReceiver); ok {
			cr.Config(strings.Split(line, "|"))
		}

		if line == "config|ready" {
			return
		}
	}
}


func NewFilter(filter interface{}) FilterWrapper {
	return &FilterWrapperImpl{
		Filter: filter,
    }
}


func Run(fw FilterWrapper) {
	scanner := bufio.NewScanner(os.Stdin)

	fw.ProcessConfig(scanner)
	fw.Register()

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

		fw.Dispatch(atoms)
	}
}
