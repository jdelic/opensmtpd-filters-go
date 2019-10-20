package opensmtpd

type LinkConnectReceiver interface {
	LinkConnect(string, SessionHolder, string, []string)
}

type LinkDisconnectReceiver interface {
	LinkDisconnect(string, SessionHolder, string, []string)
}

type LinkGreetingReceiver interface {
	LinkGreeting(string, SessionHolder, string, []string)
}

type LinkIdentifyReceiver interface {
	LinkIdentity(string, SessionHolder, string, []string)
}

type LinkTLSReceiver interface {
	LinkTLS(string, SessionHolder, string, []string)
}

type LinkAuthReceiver interface {
	LinkAuth(string, SessionHolder, string, []string)
}

type TxResetReceiver interface {
	TxReset(string, SessionHolder, string, []string)
}

type TxBeginReceiver interface {
	TxBegin(string, SessionHolder, string, []string)
}

type TxMailReceiver interface {
	TxMail(string, SessionHolder, string, []string)
}

type TxRcptReceiver interface {
	TxRcpt(string, SessionHolder, string, []string)
}

type TxEnvelopeReceiver interface {
	TxEnvelope(string, SessionHolder, string, []string)
}

type TxDataReceiver interface {
	TxData(string, SessionHolder, string, []string)
}

type TxCommitReceiver interface {
	TxCommit(string, SessionHolder, string, []string)
}

type TxRollbackReceiver interface {
	TxRollback(string, SessionHolder, string, []string)
}

type ProtocolClientReceiver interface {
	ProtocolClient(string, SessionHolder, string, []string)
}

type ProtocolServerReceiver interface {
	ProtocolServer(string, SessionHolder, string, []string)
}

type FilterReportReceiver interface {
	FilterReport(string, SessionHolder, string, []string)
}

type FilterResponseReceiver interface {
	FilterResponse(string, SessionHolder, string, []string)
}

type TimeoutReceiver interface {
	Timeout(string, SessionHolder, string, []string)
}
