package opensmtpd


type ConnectFilter interface {
	Connect(string, SessionHolder, string, []string)
}

type HeloFilter interface {
	Helo(string, SessionHolder, string, []string)
}

type EhloFilter interface {
	Ehlo(string, SessionHolder, string, []string)
}

type StartTLSFilter interface {
	StartTLS(string, SessionHolder, string, []string)
}

type AuthFilter interface {
	Auth(string, SessionHolder, string, []string)
}

type MailFromFilter interface {
	MailFrom(string, SessionHolder, string, []string)
}

type RcptToFilter interface {
	RcptTo(string, SessionHolder, string, []string)
}

type DataFilter interface {
	Data(string, SessionHolder, string, []string)
}

type DatalineFilter interface {
	Dataline(string, SessionHolder, string, []string)
}

type RsetFilter interface {
	Rset(string, SessionHolder, string, []string)
}

type QuitFilter interface {
	Quit(string, SessionHolder, string, []string)
}

type NoopFilter interface {
	Noop(string, SessionHolder, string, []string)
}

type HelpFilter interface {
	Help(string, SessionHolder, string, []string)
}

type WizFilter interface {
	Wiz(string, SessionHolder, string, []string)
}

type CommitFilter interface {
	Commit(string, SessionHolder, string, []string)
}
