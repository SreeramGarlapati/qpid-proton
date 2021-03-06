/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package electron

import (
	"qpid.apache.org/proton"
)

// Session is an AMQP session, it contains Senders and Receivers.
type Session interface {
	Endpoint

	// Sender opens a new sender.
	Sender(...LinkOption) (Sender, error)

	// Receiver opens a new Receiver.
	Receiver(...LinkOption) (Receiver, error)
}

type session struct {
	endpoint
	eSession   proton.Session
	connection *connection
	capacity   uint
}

// SessionOption can be passed when creating a Session
type SessionOption func(*session)

// IncomingCapacity returns a Session Option that sets the size (in bytes) of
// the sessions incoming data buffer..
func IncomingCapacity(cap uint) SessionOption { return func(s *session) { s.capacity = cap } }

// in proton goroutine
func newSession(c *connection, es proton.Session, setting ...SessionOption) *session {
	s := &session{
		connection: c,
		eSession:   es,
		endpoint:   makeEndpoint(es.String()),
	}
	for _, set := range setting {
		set(s)
	}
	c.handler.sessions[s.eSession] = s
	s.eSession.SetIncomingCapacity(s.capacity)
	s.eSession.Open()
	return s
}

func (s *session) Connection() Connection     { return s.connection }
func (s *session) eEndpoint() proton.Endpoint { return s.eSession }
func (s *session) engine() *proton.Engine     { return s.connection.engine }

func (s *session) Close(err error) {
	s.engine().Inject(func() {
		if s.Error() == nil {
			localClose(s.eSession, err)
		}
	})
}

func (s *session) Sender(setting ...LinkOption) (snd Sender, err error) {
	err = s.engine().InjectWait(func() error {
		if s.Error() != nil {
			return s.Error()
		}
		l, err := localLink(s, true, setting...)
		if err == nil {
			snd = newSender(l)
		}
		return err
	})
	return
}

func (s *session) Receiver(setting ...LinkOption) (rcv Receiver, err error) {
	err = s.engine().InjectWait(func() error {
		if s.Error() != nil {
			return s.Error()
		}
		l, err := localLink(s, false, setting...)
		if err == nil {
			rcv = newReceiver(l)
		}
		return err
	})
	return
}

// IncomingSender is sent on the Connection.Incoming() channel when there is an
// incoming request to open a session.
type IncomingSession struct {
	incoming
	h        *handler
	pSession proton.Session
	capacity uint
}

func newIncomingSession(h *handler, ps proton.Session) *IncomingSession {
	return &IncomingSession{incoming: makeIncoming(ps), h: h, pSession: ps}
}

// SetCapacity sets the session buffer capacity of an incoming session in bytes.
func (in *IncomingSession) SetCapacity(bytes uint) { in.capacity = bytes }

// Accept an incoming session endpoint.
func (in *IncomingSession) Accept() Endpoint {
	return in.accept(func() Endpoint {
		return newSession(in.h.connection, in.pSession, IncomingCapacity(in.capacity))
	})
}
