package opensmtpd

import (
	"log"
	"strings"
)

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

func (sf *SessionTrackingMixin) LinkConnect(verb string, sh SessionHolder, sessionId string, params []string) {
	if len(params) != 4 {
		log.Fatal("invalid input, shouldn't happen")
	}

	s := SMTPSession{}
	s.Id = sessionId
	s.Rdns = params[0]
	s.Src = params[2]

	sh.SetSession(&s)
}

func (sf *SessionTrackingMixin) LinkDisconnect(verb string, sh SessionHolder, sessionId string, params []string) {
	if len(params) != 0 {
		log.Fatal("invalid input, shouldn't happen")
	}
	delete(sh.GetSessions(), sessionId)
}

func (sf *SessionTrackingMixin) LinkGreeting(verb string, sh SessionHolder, sessionId string, params []string) {
	if len(params) != 1 {
		log.Fatal("invalid input, shouldn't happen")
	}

	s := sh.GetSession(sessionId)
	s.MtaName = params[0]
	sh.SetSession(s)
}

func (sf *SessionTrackingMixin) LinkIdentify(verb string, sh SessionHolder, sessionId string, params []string) {
	if len(params) != 2 {
		log.Fatal("invalid input, shouldn't happen")
	}

	s := sh.GetSession(sessionId)
	s.HeloName = params[1]
	sh.SetSession(s)
}

func (sf *SessionTrackingMixin) LinkAuth(verb string, sh SessionHolder, sessionId string, params []string) {
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

func (sf *SessionTrackingMixin) TxReset(verb string, sh SessionHolder, sessionId string, params []string) {
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

func (sf *SessionTrackingMixin) TxBegin(verb string, sh SessionHolder, sessionId string, params []string) {
	if len(params) != 1 {
		log.Fatal("invalid input, shouldn't happen")
	}

	s := sh.GetSession(sessionId)
	s.Msgid = params[0]
	sh.SetSession(s)
}

func (sf *SessionTrackingMixin) TxMail(verb string, sh SessionHolder, sessionId string, params []string) {
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

func (sf *SessionTrackingMixin) TxRcpt(verb string, sh SessionHolder, sessionId string, params []string) {
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

func (sf *SessionTrackingMixin) Dataline(verb string, sh SessionHolder, sessionId string, params []string) {
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
	s.Message = append(s.Message, line)
	sh.SetSession(s)
}

func (sf *SessionTrackingMixin) Commit(verb string, sh SessionHolder, sessionId string, params []string) {
	if len(params) != 2 {
		log.Fatal("invalid input, shouldn't happen")
	}
}

