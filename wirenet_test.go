package wirenet

import (
	"context"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
)

func TestRole_IsClientSide(t *testing.T) {
	assert.Equal(t, clientSide, role(1))
}

func TestRole_IsServerSide(t *testing.T) {
	assert.Equal(t, serverSide, role(2))
}

func TestRole_String(t *testing.T) {
	assert.Equal(t, clientSide.String(), "client side")
	assert.Equal(t, serverSide.String(), "server side")
	assert.Equal(t, role(999).String(), "unknown")
}

func TestHub(t *testing.T) {
	w, err := Hub(":8989", nil)
	assert.Nil(t, err)
	assert.True(t, w.(*wire).hubMode)
}

func TestJoin(t *testing.T) {
	w, err := Join(":8989", nil)
	assert.Nil(t, err)
	assert.True(t, w.(*wire).role.IsClientSide())
}

func TestMount(t *testing.T) {
	w, err := Mount("")
	assert.Nil(t, w)
	assert.Equal(t, ErrAddrEmpty, err)

	w, err = Mount(":8989", nil)
	assert.Nil(t, err)
	assert.True(t, w.(*wire).role.IsServerSide())
}

func TestWire_Connect(t *testing.T) {
	addr := genAddr(t)
	conn := make(chan struct{})

	// server side
	server, err := Mount(addr, WithConnectHook(func(closer io.Closer) {
		close(conn)
	}))
	assert.Nil(t, err)
	go func() {
		assert.Nil(t, server.Connect())
	}()
	<-conn

	client, err := Join(addr, WithConnectHook(func(closer io.Closer) {
		time.Sleep(300 * time.Millisecond)
		assert.Nil(t, closer.Close())
		assert.Nil(t, server.Close())
	}))
	assert.Nil(t, err)
	assert.Nil(t, client.Connect())
}

func TestWire_ConnectTLS(t *testing.T) {
	addr := genAddr(t)
	conn := make(chan struct{})

	// server side
	serverTLSConf, err := LoadCertificates("server", "./certs")
	assert.Nil(t, err)
	server, err := Mount(addr,
		WithTLS(serverTLSConf),
		WithConnectHook(func(closer io.Closer) {
			close(conn)
		}))
	assert.Nil(t, err)
	go func() {
		assert.Nil(t, server.Connect())
	}()
	<-conn

	clientTLSConf, err := LoadCertificates("client", "./certs")
	assert.Nil(t, err)
	clientTLSConf.InsecureSkipVerify = true
	client, err := Join(addr,
		WithTLS(clientTLSConf),
		WithConnectHook(func(closer io.Closer) {
			time.Sleep(time.Second)
			assert.Nil(t, closer.Close())
			assert.Nil(t, server.Close())
		}))
	assert.Nil(t, err)
	assert.Nil(t, client.Connect())
}

func TestWire_ReConnect(t *testing.T) {
	addr := genAddr(t)
	retryMin := 1 * time.Second
	retryMax := 2 * time.Second
	retryNum := 2
	var retryCounter int
	client, err := Join(addr,
		WithRetryPolicy(func(min, max time.Duration, attemptNum int) time.Duration {
			retryCounter++
			return DefaultRetryPolicy(min, max, attemptNum)
		}),
		WithRetryMax(retryNum),
		WithRetryWait(retryMin, retryMax),
	)
	assert.Nil(t, err)
	assert.Contains(t, client.Connect().Error(), "connection refused")
	assert.Equal(t, retryNum, retryCounter)
}

func TestWire_Session(t *testing.T) {
	addr := genAddr(t)
	initSrv := make(chan struct{})
	initCli := make(chan struct{})

	// server side
	server, err := Mount(addr, WithConnectHook(func(closer io.Closer) {
		close(initSrv)
	}))
	assert.Nil(t, err)
	go func() {
		assert.Nil(t, server.Connect())
	}()
	<-initSrv

	// client side
	client, err := Join(addr, WithConnectHook(func(closer io.Closer) {
		time.Sleep(time.Second)
		close(initCli)
	}))
	assert.Nil(t, err)
	go func() {
		assert.Nil(t, client.Connect())
	}()
	<-initCli

	s, err := server.Session(uuid.New())
	assert.Error(t, ErrSessionNotFound, err)
	assert.Nil(t, s)

	for _, sess := range server.Sessions() {
		found, err := client.Session(sess.ID())
		assert.Nil(t, err)
		assert.Equal(t, found.ID(), sess.ID())
	}

	assert.Nil(t, client.Close())
	assert.Nil(t, server.Close())
}

func TestWire_ErrorHandler(t *testing.T) {
	addr := genAddr(t)
	initSrv := make(chan struct{})
	initCli := make(chan Session)

	// server side
	server, err := Mount(addr,
		WithErrorHandler(func(_ context.Context, err error) {
			assert.Contains(t, err.Error(), "validate stream")
		}),
		WithConnectHook(func(closer io.Closer) {
			time.Sleep(time.Second)
			close(initSrv)
		}))
	assert.Nil(t, err)
	go func() {
		assert.Nil(t, server.Connect())
	}()
	<-initSrv

	// client side
	client, err := Join(addr, WithSessionOpenHook(func(s Session) {
		initCli <- s
	}))
	assert.Nil(t, err)
	go func() {
		assert.Nil(t, client.Connect())
	}()
	sess := <-initCli

	stream, err := sess.OpenStream("unknown")
	assert.Nil(t, stream)
	assert.Equal(t, ErrStreamHandlerNotFound, err)
}

func genAddr(t *testing.T) string {
	if t == nil {
		t = new(testing.T)
	}
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:0")
	assert.Nil(t, err)
	listener, err := net.ListenTCP("tcp", addr)
	assert.Nil(t, err)
	defer listener.Close()
	port := listener.Addr().(*net.TCPAddr).Port
	return fmt.Sprintf(":%d", port)
}
