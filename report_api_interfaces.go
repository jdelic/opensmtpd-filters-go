package opensmtpd

type LinkConnectReceiver interface {
	LinkConnect(FilterWrapper, string, SessionHolder, string, []string)
}

type LinkDisconnectReceiver interface {
	LinkDisconnect(FilterWrapper, string, SessionHolder, string, []string)
}

type LinkGreetingReceiver interface {
	LinkGreeting(FilterWrapper, string, SessionHolder, string, []string)
}

type LinkIdentifyReceiver interface {
	LinkIdentity(FilterWrapper, string, SessionHolder, string, []string)
}

type LinkTLSReceiver interface {
	LinkTLS(FilterWrapper, string, SessionHolder, string, []string)
}

type LinkAuthReceiver interface {
	LinkAuth(FilterWrapper, string, SessionHolder, string, []string)
}

type TxResetReceiver interface {
	TxReset(FilterWrapper, string, SessionHolder, string, []string)
}

type TxBeginReceiver interface {
	TxBegin(FilterWrapper, string, SessionHolder, string, []string)
}

type TxMailReceiver interface {
	TxMail(FilterWrapper, string, SessionHolder, string, []string)
}

type TxRcptReceiver interface {
	TxRcpt(FilterWrapper, string, SessionHolder, string, []string)
}

type TxEnvelopeReceiver interface {
	TxEnvelope(FilterWrapper, string, SessionHolder, string, []string)
}

type TxDataReceiver interface {
	TxData(FilterWrapper, string, SessionHolder, string, []string)
}

type TxCommitReceiver interface {
	TxCommit(FilterWrapper, string, SessionHolder, string, []string)
}

type TxRollbackReceiver interface {
	TxRollback(FilterWrapper, string, SessionHolder, string, []string)
}

type ProtocolClientReceiver interface {
	ProtocolClient(FilterWrapper, string, SessionHolder, string, []string)
}

type ProtocolServerReceiver interface {
	ProtocolServer(FilterWrapper, string, SessionHolder, string, []string)
}

type FilterReportReceiver interface {
	FilterReport(FilterWrapper, string, SessionHolder, string, []string)
}

type FilterResponseReceiver interface {
	FilterResponse(FilterWrapper, string, SessionHolder, string, []string)
}

type TimeoutReceiver interface {
	Timeout(FilterWrapper, string, SessionHolder, string, []string)
}
