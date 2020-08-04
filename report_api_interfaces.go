package opensmtpd

type LinkConnectReceiver interface {
	LinkConnect(FilterWrapper, FilterEvent)
}

type LinkDisconnectReceiver interface {
	LinkDisconnect(FilterWrapper, FilterEvent)
}

type LinkGreetingReceiver interface {
	LinkGreeting(FilterWrapper, FilterEvent)
}

type LinkIdentifyReceiver interface {
	LinkIdentify(FilterWrapper, FilterEvent)
}

type LinkTLSReceiver interface {
	LinkTLS(FilterWrapper, FilterEvent)
}

type LinkAuthReceiver interface {
	LinkAuth(FilterWrapper, FilterEvent)
}

type TxResetReceiver interface {
	TxReset(FilterWrapper, FilterEvent)
}

type TxBeginCallback interface {
	TxBeginCallback(string, *SMTPSession)
}

type TxBeginReceiver interface {
	TxBegin(FilterWrapper, FilterEvent)
}

type TxMailReceiver interface {
	TxMail(FilterWrapper, FilterEvent)
}

type TxRcptReceiver interface {
	TxRcpt(FilterWrapper, FilterEvent)
}

type TxEnvelopeReceiver interface {
	TxEnvelope(FilterWrapper, FilterEvent)
}

type TxDataReceiver interface {
	TxData(FilterWrapper, FilterEvent)
}

type TxCommitReceiver interface {
	TxCommit(FilterWrapper, FilterEvent)
}

type TxRollbackReceiver interface {
	TxRollback(FilterWrapper, FilterEvent)
}

type ProtocolClientReceiver interface {
	ProtocolClient(FilterWrapper, FilterEvent)
}

type ProtocolServerReceiver interface {
	ProtocolServer(FilterWrapper, FilterEvent)
}

type FilterReportReceiver interface {
	FilterReport(FilterWrapper, FilterEvent)
}

type FilterResponseReceiver interface {
	FilterResponse(FilterWrapper, FilterEvent)
}

type TimeoutReceiver interface {
	Timeout(FilterWrapper, FilterEvent)
}
