package opensmtpd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type FilterWrapper interface {
	GetCapabilities() FilterDispatchMap
	Register(EventResponder)
	Dispatch([]string)
	ProcessConfig(*bufio.Scanner)
	GetFilter() interface{}
}

type FilterWrapperImpl struct {
	Filter interface{}
}

func (fwi *FilterWrapperImpl) GetFilter() interface{} {
	return fwi.Filter
}

func (fwi *FilterWrapperImpl) GetCapabilities() FilterDispatchMap {
	capabilities := make(FilterDispatchMap)

	capabilities["report"] = make(map[string]EventHandler)
	reportReceivers := capabilities["report"]
	if _, ok := fwi.Filter.(LinkConnectReceiver); ok {
		reportReceivers["link-connect"] =
			func(fw FilterWrapper, ev FilterEvent) {
				fw.GetFilter().(LinkConnectReceiver).LinkConnect(fw, ev)
			}
	}
	if _, ok := fwi.Filter.(LinkDisconnectReceiver); ok {
		reportReceivers["link-disconnect"] =
			func(fw FilterWrapper, ev FilterEvent) {
				fw.GetFilter().(LinkDisconnectReceiver).LinkDisconnect(fw, ev)
			}
	}
	if _, ok := fwi.Filter.(LinkGreetingReceiver); ok {
		reportReceivers["link-greeting"] =
			func(fw FilterWrapper, ev FilterEvent) {
				fw.GetFilter().(LinkGreetingReceiver).LinkGreeting(fw, ev)
			}
	}
	if _, ok := fwi.Filter.(LinkIdentityReceiver); ok {
		reportReceivers["link-identify"] =
			func(fw FilterWrapper, ev FilterEvent) {
				fw.GetFilter().(LinkIdentityReceiver).LinkIdentity(fw, ev)
			}
	}
	if _, ok := fwi.Filter.(LinkTLSReceiver); ok {
		reportReceivers["link-tls"] =
			func(fw FilterWrapper, ev FilterEvent) {
				fw.GetFilter().(LinkTLSReceiver).LinkTLS(fw, ev)
			}
	}
	if _, ok := fwi.Filter.(LinkAuthReceiver); ok {
		reportReceivers["link-auth"] =
			func(fw FilterWrapper, ev FilterEvent) {
				fw.GetFilter().(LinkAuthReceiver).LinkAuth(fw, ev)
			}
	}
	if _, ok := fwi.Filter.(TxResetReceiver); ok {
		reportReceivers["tx-reset"] =
			func(fw FilterWrapper, ev FilterEvent) {
				fw.GetFilter().(TxResetReceiver).TxReset(fw, ev)
			}
	}
	if _, ok := fwi.Filter.(TxBeginReceiver); ok {
		reportReceivers["tx-begin"] =
			func(fw FilterWrapper, ev FilterEvent) {
				fw.GetFilter().(TxBeginReceiver).TxBegin(fw, ev)
			}
	}
	if _, ok := fwi.Filter.(TxMailReceiver); ok {
		reportReceivers["tx-mail"] =
			func(fw FilterWrapper, ev FilterEvent) {
				fw.GetFilter().(TxMailReceiver).TxMail(fw, ev)
			}
	}
	if _, ok := fwi.Filter.(TxRcptReceiver); ok {
		reportReceivers["tx-rcpt"] =
			func(fw FilterWrapper, ev FilterEvent) {
				fw.GetFilter().(TxRcptReceiver).TxRcpt(fw, ev)
			}
	}
	if _, ok := fwi.Filter.(TxEnvelopeReceiver); ok {
		reportReceivers["tx-envelope"] =
			func(fw FilterWrapper, ev FilterEvent) {
				fw.GetFilter().(TxEnvelopeReceiver).TxEnvelope(fw, ev)
			}
	}
	if _, ok := fwi.Filter.(TxDataReceiver); ok {
		reportReceivers["tx-data"] =
			func(fw FilterWrapper, ev FilterEvent) {
				fw.GetFilter().(TxDataReceiver).TxData(fw, ev)
			}
	}
	if _, ok := fwi.Filter.(TxCommitReceiver); ok {
		reportReceivers["tx-commit"] =
			func(fw FilterWrapper, ev FilterEvent) {
				fw.GetFilter().(TxCommitReceiver).TxCommit(fw, ev)
			}
	}
	if _, ok := fwi.Filter.(TxRollbackReceiver); ok {
		reportReceivers["link-connect"] =
			func(fw FilterWrapper, ev FilterEvent) {
				fw.GetFilter().(LinkConnectReceiver).LinkConnect(fw, ev)
			}
	}
	if _, ok := fwi.Filter.(ProtocolClientReceiver); ok {
		reportReceivers["protocol-client"] =
			func(fw FilterWrapper, ev FilterEvent) {
				fw.GetFilter().(ProtocolClientReceiver).ProtocolClient(fw, ev)
			}
	}
	if _, ok := fwi.Filter.(ProtocolServerReceiver); ok {
		reportReceivers["protocol-server"] =
			func(fw FilterWrapper, ev FilterEvent) {
				fw.GetFilter().(ProtocolServerReceiver).ProtocolServer(fw, ev)
			}
	}
	if _, ok := fwi.Filter.(TimeoutReceiver); ok {
		reportReceivers["timeout"] =
			func(fw FilterWrapper, ev FilterEvent) {
				fw.GetFilter().(TimeoutReceiver).Timeout(fw, ev)
			}
	}
	if _, ok := fwi.Filter.(FilterResponseReceiver); ok {
		reportReceivers["filter-response"] =
			func(fw FilterWrapper, ev FilterEvent) {
				fw.GetFilter().(FilterResponseReceiver).FilterResponse(fw, ev)
			}
	}
	if _, ok := fwi.Filter.(FilterReportReceiver); ok {
		reportReceivers["filter-report"] =
			func(fw FilterWrapper, ev FilterEvent) {
				fw.GetFilter().(FilterReportReceiver).FilterReport(fw, ev)
			}
	}
	capabilities["report"] = reportReceivers

	capabilities["filter"] = make(map[string]EventHandler)
	filters := capabilities["filter"]
	if _, ok := fwi.Filter.(ConnectFilter); ok {
		filters["connect"] =
			func(fw FilterWrapper, ev FilterEvent) {
				fw.GetFilter().(ConnectFilter).Connect(fw, ev)
			}
	}
	if _, ok := fwi.Filter.(HeloFilter); ok {
		filters["helo"] =
			func(fw FilterWrapper, ev FilterEvent) {
				fw.GetFilter().(HeloFilter).Helo(fw, ev)
			}
	}
	if _, ok := fwi.Filter.(EhloFilter); ok {
		filters["ehlo"] =
			func(fw FilterWrapper, ev FilterEvent) {
				fw.GetFilter().(EhloFilter).Ehlo(fw, ev)
			}
	}
	if _, ok := fwi.Filter.(StartTLSFilter); ok {
		filters["starttls"] =
			func(fw FilterWrapper, ev FilterEvent) {
				fw.GetFilter().(StartTLSFilter).StartTLS(fw, ev)
			}
	}
	if _, ok := fwi.Filter.(AuthFilter); ok {
		filters["auth"] =
			func(fw FilterWrapper, ev FilterEvent) {
				fw.GetFilter().(AuthFilter).Auth(fw, ev)
			}
	}
	if _, ok := fwi.Filter.(MailFromFilter); ok {
		filters["mail-from"] =
			func(fw FilterWrapper, ev FilterEvent) {
				fw.GetFilter().(MailFromFilter).MailFrom(fw, ev)
			}
	}
	if _, ok := fwi.Filter.(RcptToFilter); ok {
		filters["rcpt-to"] =
			func(fw FilterWrapper, ev FilterEvent) {
				fw.GetFilter().(RcptToFilter).RcptTo(fw, ev)
			}
	}
	if _, ok := fwi.Filter.(DataFilter); ok {
		filters["data"] =
			func(fw FilterWrapper, ev FilterEvent) {
				fw.GetFilter().(DataFilter).Data(fw, ev)
			}
	}
	if _, ok := fwi.Filter.(DatalineFilter); ok {
		filters["data-line"] =
			func(fw FilterWrapper, ev FilterEvent) {
				fw.GetFilter().(DatalineFilter).Dataline(fw, ev)
			}
	}
	if _, ok := fwi.Filter.(RsetFilter); ok {
		filters["rset"] =
			func(fw FilterWrapper, ev FilterEvent) {
				fw.GetFilter().(RsetFilter).Rset(fw, ev)
			}
	}
	if _, ok := fwi.Filter.(QuitFilter); ok {
		filters["quit"] =
			func(fw FilterWrapper, ev FilterEvent) {
				fw.GetFilter().(QuitFilter).Quit(fw, ev)
			}
	}
	if _, ok := fwi.Filter.(NoopFilter); ok {
		filters["noop"] =
			func(fw FilterWrapper, ev FilterEvent) {
				fw.GetFilter().(NoopFilter).Noop(fw, ev)
			}
	}
	if _, ok := fwi.Filter.(HelpFilter); ok {
		filters["help"] =
			func(fw FilterWrapper, ev FilterEvent) {
				fw.GetFilter().(HelpFilter).Help(fw, ev)
			}
	}
	if _, ok := fwi.Filter.(WizFilter); ok {
		filters["wiz"] =
			func(fw FilterWrapper, ev FilterEvent) {
				fw.GetFilter().(WizFilter).Wiz(fw, ev)
			}
	}
	if _, ok := fwi.Filter.(CommitFilter); ok {
		filters["commit"] =
			func(fw FilterWrapper, ev FilterEvent) {
				fw.GetFilter().(CommitFilter).Commit(fw, ev)
			}
	}
	capabilities["filter"] = filters
	return capabilities
}

func (fwi *FilterWrapperImpl) Register(out EventResponder) {
	capabilities := fwi.GetCapabilities()
	for typ := range capabilities {
		for op := range capabilities[typ] {
			out.SafePrintln(fmt.Sprintf("register|%v|smtp-in|%v", typ, op))
		}
	}
	out.SafePrintln("register|ready")
}

func (fwi *FilterWrapperImpl) Dispatch(atoms []string) {
	fcap := fwi.GetCapabilities()
	if handler, ok := fcap[atoms[0]][atoms[4]]; ok {
		event := NewFilterEvent(atoms)
		//handler(fwi, verb: atoms[4], sh, sessionId: atoms[5], params: atoms[6:])
		handler(fwi, event)
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

