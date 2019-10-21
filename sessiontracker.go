package opensmtpd

import (
	"log"
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

func (sf *SessionTrackingMixin) LinkConnect(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
	if len(params) != 4 {
		log.Fatal("invalid input, shouldn't happen")
	}

	s := SMTPSession{}
	s.Id = sessionId
	s.Rdns = params[0]
	s.Src = params[2]

	sh.SetSession(&s)
}

func (sf *SessionTrackingMixin) LinkDisconnect(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
	if len(params) != 0 {
		log.Fatal("invalid input, shouldn't happen")
	}
	delete(sh.GetSessions(), sessionId)
}

func (sf *SessionTrackingMixin) LinkGreeting(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
	if len(params) != 1 {
		log.Fatal("invalid input, shouldn't happen")
	}

	s := sh.GetSession(sessionId)
	s.MtaName = params[0]
	sh.SetSession(s)
}

func (sf *SessionTrackingMixin) LinkIdentify(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
	if len(params) != 2 {
		log.Fatal("invalid input, shouldn't happen")
	}

	s := sh.GetSession(sessionId)
	s.HeloName = params[1]
	sh.SetSession(s)
}

func (sf *SessionTrackingMixin) LinkAuth(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
	if len(params) != 2 {
		log.Fatal("invalid input, shouldn't happen")
	}
	if params[1] != "pass" {
		return
	}
	s := sh.GetSession(sessionId)
	s.UserName = params[0]
	sh.SetSession(s)
}

func (sf *SessionTrackingMixin) TxReset(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
	if len(params) != 1 {
		log.Fatal("invalid input, shouldn't happen")
	}

	s := sh.GetSession(sessionId)
	s.Msgid = ""
	s.MailFrom = ""
	s.RcptTo = nil
	s.Message = nil
	sh.SetSession(s)
}

func (sf *SessionTrackingMixin) TxBegin(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
	if len(params) != 1 {
		log.Fatal("invalid input, shouldn't happen")
	}

	s := sh.GetSession(sessionId)
	s.Msgid = params[0]
	sh.SetSession(s)
}

func (sf *SessionTrackingMixin) TxMail(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
	if len(params) != 3 {
		log.Fatal("invalid input, shouldn't happen")
	}

	if params[2] != "ok" {
		return
	}

	s := sh.GetSession(sessionId)
	s.MailFrom = params[1]
	sh.SetSession(s)
}

func (sf *SessionTrackingMixin) TxRcpt(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
	if len(params) != 3 {
		log.Fatal("invalid input, shouldn't happen")
	}

	if params[2] != "ok" {
		return
	}

	s := sh.GetSession(sessionId)
	s.RcptTo = append(s.RcptTo, params[1])
	sh.SetSession(s)
}

func (sf *SessionTrackingMixin) Dataline(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
	if len(params) < 2 {
		log.Fatal("invalid input, shouldn't happen")
	}
	//token := params[0]
	line := strings.Join(params[1:], "")

	s := sh.GetSession(sessionId)
	if line == "." {
		s.Message = append(s.Message, line)
		if cb, ok := (interface {})(*sf).(MessageReceivedCallback); ok {
			s := sh.GetSession(sessionId)
			cb.MessageComplete(params[0], s)
			sh.SetSession(s)
		} else {
			FlushMessage(params[0], sh.GetSession(sessionId))
		}
		return
	}

	// Input is raw SMTP data - unescape leading dots.
	line = strings.TrimPrefix(line, ".")

	s.Message = append(s.Message, line)
	sh.SetSession(s)
}

func (sf *SessionTrackingMixin) Commit(fw FilterWrapper, verb string, sh SessionHolder, sessionId string, params []string) {
	if len(params) != 2 {
		log.Fatal("invalid input, shouldn't happen")
	}
}

