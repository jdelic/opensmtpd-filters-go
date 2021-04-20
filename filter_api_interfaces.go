package opensmtpd


type ConfigReceiver interface {
	Config([]string)
}

type MessageReceivedCallback interface {
	MessageComplete(*FilterEvent, *SMTPSession)
}

type ConnectFilter interface {
	Connect(FilterWrapper, FilterEvent)
}

type HeloFilter interface {
	Helo(FilterWrapper, FilterEvent)
}

type EhloFilter interface {
	Ehlo(FilterWrapper, FilterEvent)
}

type StartTLSFilter interface {
	StartTLS(FilterWrapper, FilterEvent)
}

type AuthFilter interface {
	Auth(FilterWrapper, FilterEvent)
}

type MailFromFilter interface {
	MailFrom(FilterWrapper, FilterEvent)
}

type RcptToFilter interface {
	RcptTo(FilterWrapper, FilterEvent)
}

type DataFilter interface {
	Data(FilterWrapper, FilterEvent)
}

type DatalineFilter interface {
	Dataline(FilterWrapper, FilterEvent)
}

type RsetFilter interface {
	Rset(FilterWrapper, FilterEvent)
}

type QuitFilter interface {
	Quit(FilterWrapper, FilterEvent)
}

type NoopFilter interface {
	Noop(FilterWrapper, FilterEvent)
}

type HelpFilter interface {
	Help(FilterWrapper, FilterEvent)
}

type WizFilter interface {
	Wiz(FilterWrapper, FilterEvent)
}

type CommitFilter interface {
	Commit(FilterWrapper, FilterEvent)
}
