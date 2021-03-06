package wirenet

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/yamux"
)

// Session represents a connection between a client and a server.
// Each session can have from one to N named streams.
// Each named stream is a payload processing (file transfer, video transfer, etc.).
type Session interface {

	// ID returns an unique session id.
	ID() uuid.UUID

	// IsClosed returns a true flag if the session is closed, otherwise returns a false flag.
	IsClosed() bool

	// Close closes gracefully shutdown the all active streams.
	Close() error

	// StreamNames returns a list of open stream names.
	StreamNames() []string

	// OpenStream opens a named stream and returns it.
	// After the named stream is successfully opened, an authentication frame is sent.
	OpenStream(name string) (Stream, error)

	// Identification returns some information specified by the user on the client side using WithIdentification().
	Identification() Identification

	// CloseWire closes gracefully shutdown the server without interrupting any active connections.
	CloseWire() error
}

type session struct {
	id             uuid.UUID
	conn           *yamux.Session
	w              *wire
	streamNames    []string
	closed         bool
	closeCh        chan chan error
	activeStreams  int
	streams        map[uuid.UUID]Stream
	mu             sync.RWMutex
	timeoutDur     time.Duration
	identification Identification
}

func openSession(sid uuid.UUID, id Identification, conn *yamux.Session, w *wire, streamNames []string) {
	sess := &session{
		id:             sid,
		conn:           conn,
		w:              w,
		streamNames:    streamNames,
		closeCh:        make(chan chan error),
		streams:        make(map[uuid.UUID]Stream),
		timeoutDur:     w.sessCloseTimeout,
		identification: id,
	}
	go sess.open()
}

func (s *session) Identification() Identification {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.identification
}

func (s *session) StreamNames() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.streamNames
}

func (s *session) String() string {
	return fmt.Sprintf("wirenet session: %s", s.id)
}

func (s *session) registerStream(stream Stream) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.streams[stream.ID()] = stream
	s.activeStreams++
}

func (s *session) unregisterStream(stream Stream) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.streams, stream.ID())
	if s.activeStreams > 0 {
		s.activeStreams--
	}
}

func (s *session) activeStreamCounter() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.activeStreams
}

func (s *session) shutdown() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			errCh, ok := <-s.closeCh
			if !ok {
				return
			}

			timeout := time.Now().Add(s.timeoutDur)
			for {
				if s.activeStreamCounter() <= 0 {
					break
				}

				if timeout.Unix() <= time.Now().Unix() {
					for _, stream := range s.streams {
						stream.Close()
					}
					continue
				}
				time.Sleep(300 * time.Millisecond)
			}

			cancel()

			var closeErr error
			if !s.conn.IsClosed() {
				closeErr = s.conn.Close()
			}
			errCh <- closeErr
			close(errCh)
			break
		}
	}()
	return ctx
}

func (s *session) validateStreamName(streamName string) (err error) {
	isHubMode := s.w.isHubMode() && !s.w.role.IsClientSide()
	if isHubMode {
		_, err = s.w.findSession(streamName)
		if err == ErrSessionNotFound {
			_, err = s.w.findHandler(streamName)
		}
	} else {
		_, err = s.w.findHandler(streamName)
	}
	return err
}

func (s *session) validateToken(streamName string, token []byte) (err error) {
	if s.w.verifyToken != nil {
		if err := s.w.verifyToken(streamName, s.identification, token); err != nil {
			return err
		}
	}
	return
}

func (s *session) readFrame(conn *yamux.Stream) (frm frame, err error) {
	frm, err = recvFrame(conn, func(f frame) error {
		command := f.Command()
		if err := s.validateToken(command, f.Payload()); err != nil {
			return err
		}
		return s.validateStreamName(command)
	})
	return frm, err
}

func (s *session) errLog(ctx context.Context, err error, op string) {
	opErr := &OpError{
		Op:             op,
		Err:            err,
		SessionID:      s.id,
		Identification: s.identification,
		RemoteAddr:     s.conn.RemoteAddr(),
		LocalAddr:      s.conn.LocalAddr(),
	}
	s.w.errorHandler(ctx, opErr)
}

func (s *session) dispatchStream(ctx context.Context, conn *yamux.Stream) {
	defer func() {
		_ = conn.Close()
	}()

	frm, err := s.readFrame(conn)
	if err != nil {
		s.errLog(ctx, err, "validate stream")
		return
	}

	conn.Shrink()

	streamName := frm.Command()
	isHubMode := s.w.isHubMode() && !s.w.role.IsClientSide()
	if isHubMode {
		err = s.serveHub(ctx, streamName, conn)
		if err == ErrSessionNotFound {
			err = s.serve(ctx, streamName, conn)
		}

	} else {
		err = s.serve(ctx, streamName, conn)
	}
	if err != nil {
		s.errLog(ctx, err, "serve stream")
	}
}

func (s *session) serveHub(_ context.Context, streamName string, conn *yamux.Stream) error {
	sess, err := s.w.findSession(streamName)
	if err != nil {
		return err
	}
	dst, err := sess.OpenStream(streamName)
	if err != nil {
		return err
	}
	defer func() {
		conn.Close()
		dst.Close()
	}()
	dstConn := dst.(*stream).conn
	go func() {
		_ = pipe(dstConn, conn)
	}()
	return pipe(conn, dstConn)
}

func (s *session) serve(ctx context.Context, streamName string, conn *yamux.Stream) error {
	defer func() {
		if err := recover(); err != nil {
			s.errLog(ctx, fmt.Errorf("recover %v", err), "recover stream")
		}
	}()
	handler, err := s.w.findHandler(streamName)
	if err != nil {
		return err
	}
	stream := openStream(s, streamName, conn)
	handler(ctx, stream)
	if !stream.IsClosed() {
		_ = stream.Close()
	}
	return nil
}

func (s *session) open() {
	defer func() {
		_ = s.Close()
	}()

	ctx := s.shutdown()
	s.w.registerSession(s)
	go s.w.openSessHook(s)

	for {
		conn, err := s.conn.AcceptStream()
		if err != nil {
			return
		}

		if s.IsClosed() || s.w.isClosed() {
			_ = conn.Close()
			continue
		}

		go s.dispatchStream(ctx, conn)
	}
}

func (s *session) ID() uuid.UUID {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.id
}

func (s *session) IsClosed() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.closed
}

func (s *session) CloseWire() error {
	return s.w.Close()
}

func (s *session) Close() error {
	if s.IsClosed() {
		return ErrSessionClosed
	}

	s.mu.Lock()
	s.closed = true
	s.mu.Unlock()

	errCh := make(chan error)
	s.closeCh <- errCh
	closeErr := <-errCh

	s.w.unregisterSession(s)
	go s.w.closeSessHook(s)

	return closeErr
}

func (s *session) OpenStream(name string) (Stream, error) {
	if s.IsClosed() {
		return nil, ErrSessionClosed
	}

	conn, err := s.conn.OpenStream()
	if err != nil {
		return nil, err
	}

	frm, err := sendFrame(name, permFrameTyp, s.w.token, conn)
	if err != nil {
		conn.Close()
		return nil, err
	}
	if frm.Command() != name && !frm.IsRecvFrame() {
		conn.Close()
		return nil, io.ErrUnexpectedEOF
	}

	conn.Shrink()
	stream := openStream(s, name, conn)

	return stream, nil
}
