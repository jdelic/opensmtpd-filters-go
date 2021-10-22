package opensmtpd

import (
	"log"
	"strings"
)

type SMTPSession struct {
	Id string

	Rdns string
	Src string
	SrcIp string
	SrcPort string
	HeloName string
	UserName string
	MtaName string

	Msgid string
	MailFrom string
	RcptTo []string
	Message []string
}

type SessionHolder interface {
	GetSessions() map[string]*SMTPSession
	GetSession(string) *SMTPSession
	SetSession(*SMTPSession)
}

type SessionHolderImpl struct {
	Sessions map[string]*SMTPSession
}

func (shi *SessionHolderImpl) GetSessions() map[string]*SMTPSession {
	return shi.Sessions
}

func (shi *SessionHolderImpl) GetSession(sessionId string) *SMTPSession {
	if shi.Sessions == nil {
		return nil
	}
	return shi.Sessions[sessionId]
}

func (shi *SessionHolderImpl) SetSession(session *SMTPSession) {
	if shi.Sessions == nil {
		shi.Sessions = make(map[string]*SMTPSession)
	}
	shi.Sessions[session.Id] = session
}

type SessionTrackingMixin struct {
	SessionHolderImpl
}

func (sf *SessionTrackingMixin) LinkConnect(fw FilterWrapper, ev FilterEvent) {
	if len(ev.GetParams()) != 4 {
		log.Fatal("invalid input, shouldn't happen")
	}

	params := ev.GetParams()

	s := SMTPSession{}
	s.Id = ev.GetSessionId()
	s.Rdns = params[0]
	s.Src = params[2]

	// parse ipv6 if necessary
	tmp := strings.Split(s.Src, ":")
	s.SrcPort = tmp[len(tmp)-1]

	// remove the port (last section)
	tmp = tmp[0:len(tmp)-1]

	// reassemble ipv6 address with : separator
	srcIp := strings.Join(tmp, ":")
	if strings.HasPrefix(srcIp, "[") {
		// remove the ipv6 wrapper []
		srcIp = srcIp[1:len(srcIp)-1]
	}
	s.SrcIp = srcIp

	sf.SetSession(&s)
}

func (sf *SessionTrackingMixin) LinkDisconnect(fw FilterWrapper, ev FilterEvent) {
	if len(ev.GetParams()) != 0 {
		log.Fatal("invalid input, shouldn't happen")
	}

	delete(sf.GetSessions(), ev.GetSessionId())
}

func (sf *SessionTrackingMixin) LinkGreeting(fw FilterWrapper, ev FilterEvent) {
	params := ev.GetParams()
	if len(params) != 1 {
		log.Fatal("invalid input, shouldn't happen")
	}

	s := sf.GetSession(ev.GetSessionId())
	s.MtaName = params[0]
	sf.SetSession(s)
}

func (sf *SessionTrackingMixin) LinkIdentify(fw FilterWrapper, ev FilterEvent) {
	params := ev.GetParams()
	if len(params) != 2 {
		log.Fatal("invalid input, shouldn't happen")
	}

	if sh, ok := fw.GetFilter().(SessionHolder); ok {
		s := sh.GetSession(ev.GetSessionId())
		s.HeloName = params[1]
		sh.SetSession(s)
	}
}

func (sf *SessionTrackingMixin) LinkAuth(fw FilterWrapper, ev FilterEvent) {
	params := ev.GetParams()
	if len(params) != 2 {
		log.Fatal("invalid input, shouldn't happen")
	}

	// don't store usernames that didn't successfully authenticate
	if params[1] != "pass" {
		return
	}
	s := sf.GetSession(ev.GetSessionId())
	s.UserName = params[0]
	sf.SetSession(s)
}

func (sf *SessionTrackingMixin) TxReset(fw FilterWrapper, ev FilterEvent) {
	params := ev.GetParams()
	if len(params) != 1 {
		log.Fatal("invalid input, shouldn't happen")
	}

	s := sf.GetSession(ev.GetSessionId())
	s.Msgid = ""
	s.MailFrom = ""
	s.RcptTo = nil
	s.Message = nil
	sf.SetSession(s)
}

func (sf *SessionTrackingMixin) TxBegin(fw FilterWrapper, ev FilterEvent) {
	params := ev.GetParams()
	if len(params) != 1 {
		log.Fatal("invalid input, shouldn't happen")
	}

	s := sf.GetSession(ev.GetSessionId())
	s.Msgid = params[0]
	sf.SetSession(s)
}

func (sf *SessionTrackingMixin) TxMail(fw FilterWrapper, ev FilterEvent) {
	params := ev.GetParams()
	if len(params) != 3 {
		log.Fatal("invalid input, shouldn't happen")
	}

	if params[2] != "ok" {
		return
	}

	s := sf.GetSession(ev.GetSessionId())
	s.MailFrom = params[1]
	sf.SetSession(s)
}

func (sf *SessionTrackingMixin) TxRcpt(fw FilterWrapper, ev FilterEvent) {
	params := ev.GetParams()
	if len(params) != 3 {
		log.Fatal("invalid input, shouldn't happen")
	}

	if params[2] != "ok" {
		return
	}

	s := sf.GetSession(ev.GetSessionId())
	s.RcptTo = append(s.RcptTo, params[1])
	sf.SetSession(s)
}

func (sf *SessionTrackingMixin) Dataline(fw FilterWrapper, ev FilterEvent) {
	params := ev.GetParams()
	if len(params) < 2 {
		log.Fatal("invalid input, shouldn't happen")
	}
	//token := params[0]
	line := strings.Join(params[1:], "")

	s := sf.GetSession(ev.GetSessionId())
	if line == "." {
		if cb, ok := fw.GetFilter().(MessageReceivedCallback); ok {
			sf.SetSession(s)
			cb.MessageComplete(&ev, s)
		} else {
			ev.Responder().FlushMessage(s)
		}
		return
	}

	// Input is raw SMTP data - unescape leading dots.
	line = strings.TrimPrefix(line, ".")
	s.Message = append(s.Message, line)
	sf.SetSession(s)
}
