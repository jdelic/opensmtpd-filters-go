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


type FilterWrapper interface {
	GetCapabilities() FilterDispatchMap
	Register()
	Dispatch([]string)
	ProcessConfig(*bufio.Scanner)
	GetFilter() interface{}
}

type Filter interface {
	GetName() string
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

func (fwi *FilterWrapperImpl) GetFilter() interface{} {
	return fwi.Filter
}

type EventHandler = func(FilterWrapper, string, SessionHolder, string, []string)

type FilterDispatchMap = map[string]map[string]EventHandler

func (fwi *FilterWrapperImpl) GetCapabilities() FilterDispatchMap {
	capabilities := make(FilterDispatchMap)

	capabilities["report"] = make(map[string]EventHandler)
	reportReceivers := capabilities["report"]
	if _, ok := fwi.Filter.(LinkConnectReceiver); ok {
		reportReceivers["link-connect"] =
			func(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
				fw.GetFilter().(LinkConnectReceiver).LinkConnect(fw, verb, sh, sessionId, params)
			}
	}
	if _, ok := fwi.Filter.(LinkDisconnectReceiver); ok {
		reportReceivers["link-disconnect"] =
			func(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
				fw.GetFilter().(LinkDisconnectReceiver).LinkDisconnect(fw, verb, sh, sessionId, params)
			}
	}
	if _, ok := fwi.Filter.(LinkGreetingReceiver); ok {
		reportReceivers["link-greeting"] =
			func(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
				fw.GetFilter().(LinkGreetingReceiver).LinkGreeting(fw, verb, sh, sessionId, params)
			}
	}
	if _, ok := fwi.Filter.(LinkIdentifyReceiver); ok {
		reportReceivers["link-identify"] =
			func(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
				fw.GetFilter().(LinkIdentifyReceiver).LinkIdentity(fw, verb, sh, sessionId, params)
			}
	}
	if _, ok := fwi.Filter.(LinkTLSReceiver); ok {
		reportReceivers["link-tls"] =
			func(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
				fw.GetFilter().(LinkTLSReceiver).LinkTLS(fw, verb, sh, sessionId, params)
			}
	}
	if _, ok := fwi.Filter.(LinkAuthReceiver); ok {
		reportReceivers["link-auth"] =
			func(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
				fw.GetFilter().(LinkAuthReceiver).LinkAuth(fw, verb, sh, sessionId, params)
			}
	}
	if _, ok := fwi.Filter.(TxResetReceiver); ok {
		reportReceivers["tx-reset"] =
			func(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
				fw.GetFilter().(TxResetReceiver).TxReset(fw, verb, sh, sessionId, params)
			}
	}
	if _, ok := fwi.Filter.(TxBeginReceiver); ok {
		reportReceivers["tx-begin"] =
			func(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
				fw.GetFilter().(TxBeginReceiver).TxBegin(fw, verb, sh, sessionId, params)
			}
	}
	if _, ok := fwi.Filter.(TxMailReceiver); ok {
		reportReceivers["tx-mail"] =
			func(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
				fw.GetFilter().(TxMailReceiver).TxMail(fw, verb, sh, sessionId, params)
			}
	}
	if _, ok := fwi.Filter.(TxRcptReceiver); ok {
		reportReceivers["tx-rcpt"] =
			func(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
				fw.GetFilter().(TxRcptReceiver).TxRcpt(fw, verb, sh, sessionId, params)
			}
	}
	if _, ok := fwi.Filter.(TxEnvelopeReceiver); ok {
		reportReceivers["tx-envelope"] =
			func(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
				fw.GetFilter().(TxEnvelopeReceiver).TxEnvelope(fw, verb, sh, sessionId, params)
			}
	}
	if _, ok := fwi.Filter.(TxDataReceiver); ok {
		reportReceivers["tx-data"] =
			func(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
				fw.GetFilter().(TxDataReceiver).TxData(fw, verb, sh, sessionId, params)
			}
	}
	if _, ok := fwi.Filter.(TxCommitReceiver); ok {
		reportReceivers["tx-commit"] =
			func(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
				fw.GetFilter().(TxCommitReceiver).TxCommit(fw, verb, sh, sessionId, params)
			}
	}
	if _, ok := fwi.Filter.(TxRollbackReceiver); ok {
		reportReceivers["link-connect"] =
			func(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
				fw.GetFilter().(LinkConnectReceiver).LinkConnect(fw, verb, sh, sessionId, params)
			}
	}
	if _, ok := fwi.Filter.(ProtocolClientReceiver); ok {
		reportReceivers["protocol-client"] =
			func(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
				fw.GetFilter().(ProtocolClientReceiver).ProtocolClient(fw, verb, sh, sessionId, params)
			}
	}
	if _, ok := fwi.Filter.(ProtocolServerReceiver); ok {
		reportReceivers["protocol-server"] =
			func(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
				fw.GetFilter().(ProtocolServerReceiver).ProtocolServer(fw, verb, sh, sessionId, params)
			}
	}
	if _, ok := fwi.Filter.(TimeoutReceiver); ok {
		reportReceivers["timeout"] =
			func(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
				fw.GetFilter().(TimeoutReceiver).Timeout(fw, verb, sh, sessionId, params)
			}
	}
	if _, ok := fwi.Filter.(FilterResponseReceiver); ok {
		reportReceivers["filter-response"] =
			func(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
				fw.GetFilter().(FilterResponseReceiver).FilterResponse(fw, verb, sh, sessionId, params)
			}
	}
	if _, ok := fwi.Filter.(FilterReportReceiver); ok {
		reportReceivers["filter-report"] =
			func(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
				fw.GetFilter().(FilterReportReceiver).FilterReport(fw, verb, sh, sessionId, params)
			}
	}
	capabilities["report"] = reportReceivers

	capabilities["filter"] = make(map[string]EventHandler)
	filters := capabilities["filter"]
	if _, ok := fwi.Filter.(ConnectFilter); ok {
		filters["connect"] =
			func(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
				fw.GetFilter().(ConnectFilter).Connect(fw, verb, sh, sessionId, params)
			}
	}
	if _, ok := fwi.Filter.(HeloFilter); ok {
		filters["helo"] =
			func(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
				fw.GetFilter().(HeloFilter).Helo(fw, verb, sh, sessionId, params)
			}
	}
	if _, ok := fwi.Filter.(EhloFilter); ok {
		filters["ehlo"] =
			func(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
				fw.GetFilter().(EhloFilter).Ehlo(fw, verb, sh, sessionId, params)
			}
	}
	if _, ok := fwi.Filter.(StartTLSFilter); ok {
		filters["starttls"] =
			func(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
				fw.GetFilter().(StartTLSFilter).StartTLS(fw, verb, sh, sessionId, params)
			}
	}
	if _, ok := fwi.Filter.(AuthFilter); ok {
		filters["auth"] =
			func(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
				fw.GetFilter().(AuthFilter).Auth(fw, verb, sh, sessionId, params)
			}
	}
	if _, ok := fwi.Filter.(MailFromFilter); ok {
		filters["mail-from"] =
			func(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
				fw.GetFilter().(MailFromFilter).MailFrom(fw, verb, sh, sessionId, params)
			}
	}
	if _, ok := fwi.Filter.(RcptToFilter); ok {
		filters["rcpt-to"] =
			func(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
				fw.GetFilter().(RcptToFilter).RcptTo(fw, verb, sh, sessionId, params)
			}
	}
	if _, ok := fwi.Filter.(DataFilter); ok {
		filters["data"] =
			func(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
				fw.GetFilter().(DataFilter).Data(fw, verb, sh, sessionId, params)
			}
	}
	if _, ok := fwi.Filter.(DatalineFilter); ok {
		filters["data-line"] =
			func(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
				fw.GetFilter().(DatalineFilter).Dataline(fw, verb, sh, sessionId, params)
			}
	}
	if _, ok := fwi.Filter.(RsetFilter); ok {
		filters["rset"] =
			func(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
				fw.GetFilter().(RsetFilter).Rset(fw, verb, sh, sessionId, params)
			}
	}
	if _, ok := fwi.Filter.(QuitFilter); ok {
		filters["quit"] =
			func(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
				fw.GetFilter().(QuitFilter).Quit(fw, verb, sh, sessionId, params)
			}
	}
	if _, ok := fwi.Filter.(NoopFilter); ok {
		filters["noop"] =
			func(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
				fw.GetFilter().(NoopFilter).Noop(fw, verb, sh, sessionId, params)
			}
	}
	if _, ok := fwi.Filter.(HelpFilter); ok {
		filters["help"] =
			func(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
				fw.GetFilter().(HelpFilter).Help(fw, verb, sh, sessionId, params)
			}
	}
	if _, ok := fwi.Filter.(WizFilter); ok {
		filters["wiz"] =
			func(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
				fw.GetFilter().(WizFilter).Wiz(fw, verb, sh, sessionId, params)
			}
	}
	if _, ok := fwi.Filter.(CommitFilter); ok {
		filters["commit"] =
			func(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
				fw.GetFilter().(CommitFilter).Commit(fw, verb, sh, sessionId, params)
			}
	}

	capabilities["filter"] = filters
	return capabilities
}

func (fwi *FilterWrapperImpl) Register() {
	capabilities := fwi.GetCapabilities()
	for typ := range capabilities {
		for op := range capabilities[typ] {
			SafePrintf("register|%v|smtp-in|%v\n", typ, op)
		}
	}
	SafePrintln("register|ready")
}


func Proceed(token, sessionId string) {
	SafePrintf("filter-result|%s|%s|proceed\n", token, sessionId)
}

func HardReject(token, sessionId, response string) {
	SafePrintf("filter-result|%s|%s|reject|550 %s\n", token, sessionId, response)
}

func Greylist(token, sessionId, response string) {
	SafePrintf("filter-result|%s|%s|reject|421 %s\n", token, sessionId, response)
}

func SoftReject(token, sessionId, response string) {
	SafePrintf("filter-result|%s|%s|reject|451 %s\n", token, sessionId, response)
}

func FlushMessage(token string, session *SMTPSession) {
	for _, line := range session.Message {
		SafePrintf("filter-dataline|%s|%s|%s\n", token, session.Id, line)
	}
	SafePrintf("filter-dataline|%s|%s|.\n", token, session.Id)
}

func DatalineEnd(token, sessionId string) {
	SafePrintf("filter-dataline|%s|%s|.\n", token, sessionId)
}

func DatalineReply(token, sessionId, line string) {
	prefix := ""
	// Output raw SMTP data - escape leading dots.
	if strings.HasPrefix(line, ".") {
		prefix = "."
	}
	SafePrintf("filter-dataline|%s|%s|%s%s\n", token, sessionId, prefix, line)
}

func WriteMultilineHeader(token, sessionId, header, value string) {
	for i, line := range strings.Split(value, "\n") {
		if i == 0 {
			SafePrintf("filter-dataline|%s|%s|%s: %s\n", token, sessionId, header, line)
		} else {
			SafePrintf("filter-dataline|%s|%s|%s\n", token, sessionId, line)
		}
	}
}

func (fwi *FilterWrapperImpl) Dispatch(atoms []string) {
	var sh SessionHolder
	var ok bool
	if sh, ok = fwi.Filter.(SessionHolder); !ok {
		sh = nil
	}

	fcap := fwi.GetCapabilities()
	if handler, ok := fcap[atoms[0]][atoms[4]]; ok {
		handler(fwi, atoms[4], sh, atoms[5], atoms[6:])
	}
}

func (fwi *FilterWrapperImpl) ProcessConfig(scanner *bufio.Scanner) {
	for {
		if !scanner.Scan() {
			os.Exit(0)
		}
		line := scanner.Text()
		if cr, ok := fwi.Filter.(ConfigReceiver); ok {
			cr.Config(strings.Split(line, "|"))
		}

		if line == "config|ready" {
			return
		}
	}
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

func SafePrintf(format string, params... interface{}) {
	stdoutChannel <- fmt.Sprintf(format, params...)
}

func SafePrintln(msg string) {
	stdoutChannel <- msg + "\n"
}

func Run(fw FilterWrapper) {
	// start the stdout writer goroutine so we can write thread safe
	go stdoutWriter(stdoutChannel)

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
