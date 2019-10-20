package opensmtpd


type ConnectFilter interface {
	Connect(FilterWrapper, string, SessionHolder, string, []string)
}

type HeloFilter interface {
	Helo(FilterWrapper, string, SessionHolder, string, []string)
}

type EhloFilter interface {
	Ehlo(FilterWrapper, string, SessionHolder, string, []string)
}

type StartTLSFilter interface {
	StartTLS(FilterWrapper, string, SessionHolder, string, []string)
}

type AuthFilter interface {
	Auth(FilterWrapper, string, SessionHolder, string, []string)
}

type MailFromFilter interface {
	MailFrom(FilterWrapper, string, SessionHolder, string, []string)
}

type RcptToFilter interface {
	RcptTo(FilterWrapper, string, SessionHolder, string, []string)
}

type DataFilter interface {
	Data(FilterWrapper, string, SessionHolder, string, []string)
}

type DatalineFilter interface {
	Dataline(FilterWrapper, string, SessionHolder, string, []string)
}

type RsetFilter interface {
	Rset(FilterWrapper, string, SessionHolder, string, []string)
}

type QuitFilter interface {
	Quit(FilterWrapper, string, SessionHolder, string, []string)
}

type NoopFilter interface {
	Noop(FilterWrapper, string, SessionHolder, string, []string)
}

type HelpFilter interface {
	Help(FilterWrapper, string, SessionHolder, string, []string)
}

type WizFilter interface {
	Wiz(FilterWrapper, string, SessionHolder, string, []string)
}

type CommitFilter interface {
	Commit(FilterWrapper, string, SessionHolder, string, []string)
}
