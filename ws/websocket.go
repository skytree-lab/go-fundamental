package ws

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/skytree-lab/go-fundamental/util"
)

type Conn struct {
	sync.RWMutex
	*websocket.Conn
	host          string
	path          string
	heartbeatDone chan struct{}
	readMsgDone   chan struct{}
	subs          []interface{}
	eventStream   chan *Event
	writeStream   chan func()
	closed        bool
}

const (
	// Create Connection
	EventCreateConnection = iota
	// HeartBeat
	EventHeartBeat
	// Message
	EventMessage
	// Close Socket
	EventSocketClose
	// Subscribe
	EventSubscribe
	// Reconnect
	EventReconnect
)

type Event struct {
	ws *Conn

	event int
	arg   interface{}
	err   error

	messageType int
	message     []byte
}

func NewWsConn(host, path string, eventStream chan *Event) *Conn {
	Conn := &Conn{
		host:        host,
		path:        path,
		eventStream: eventStream,
		writeStream: make(chan func(), 1000000),
		closed:      true,
	}

	Conn.connect()
	go Conn.writer()
	return Conn
}

func (ws *Conn) connect() {
	var err error

	dialer := &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 45 * time.Second,
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		TLSClientConfig:  &tls.Config{InsecureSkipVerify: true},
	}

	ws.Conn, _, err = dialer.Dial(fmt.Sprintf("wss://%s%s", ws.host, ws.path), nil)
	if err == nil {
		ws.setClosed(false)
	}

	ws.fireEvent(EventCreateConnection, fmt.Sprintf("%s/%s", ws.host, ws.path), err)
}

func (ws *Conn) Reconnect(err error) {
	ws.fireEvent(EventReconnect, "", err)
	ws.CloseWs()
	go func() {
		time.Sleep(time.Second)
		ws.connect()
	}()
}

func (ws *Conn) writer() {
	for f := range ws.writeStream {
		if !ws.isClosed() {
			f()
		} else {
			util.Logger().Info("ws closed, ignore write")
		}
	}
}

func (ws *Conn) Subscribe(subEvent interface{}) {
	ws.subs = append(ws.subs, subEvent)
	ws.Write(func() {
		ws.SetWriteDeadline(time.Now().Add(time.Second))
		err := ws.WriteJSON(subEvent)
		ws.fireEvent(EventSubscribe, subEvent, err)
	})
}

func (ws *Conn) Write(f func()) {
	if ws.writeStream != nil {
		ws.writeStream <- f
	}
}

func (ws *Conn) StartToHeartbeat(f func() string, interval time.Duration) {
	ws.heartbeatDone = make(chan struct{}, 1)
	go func() {
		timer := time.NewTicker(interval)
		defer timer.Stop()
		util.Logger().Info("StartToHeartbeat start")
		for {
			select {
			case <-timer.C:
				ws.Write(func() {
					ping := f()
					err := ws.WriteControl(websocket.PingMessage, []byte(ping), time.Now().Add(time.Second))
					if err != nil {
						msg := fmt.Sprintf("SetWriteDeadline failed, err = %v", err)
						util.Logger().Error(msg)
					}
					ws.fireEvent(EventHeartBeat, ping, err)
				})
			case <-ws.heartbeatDone:
				util.Logger().Info("StartToHeartbeat exit")
				return
			}
		}
	}()
}

func (ws *Conn) isClosed() bool {
	ws.Lock()
	closed := ws.closed
	ws.Unlock()
	return closed
}

func (ws *Conn) setClosed(closed bool) {
	ws.Lock()
	ws.closed = closed
	ws.Unlock()
}

func (ws *Conn) StartToReceiveMessage() {
	ws.readMsgDone = make(chan struct{}, 1)
	go func() {
		util.Logger().Info("receive message loop start")
		for {
			if ws.isClosed() {
				util.Logger().Info("receive message loop exit")
				return
			}

			t, msg, err := ws.ReadMessage()
			if err != nil {
				if !ws.isClosed() {
					msg := fmt.Sprintf("receive message failed, err = %v", err)
					util.Logger().Error(msg)
					ws.Reconnect(err)
				}
				util.Logger().Info("receive message loop exit")
				return
			}
			ws.fireMessage(err, t, msg)
		}
	}()
}

func (ws *Conn) CloseWs() {
	ws.setClosed(true)

	if ws.heartbeatDone != nil {
		close(ws.heartbeatDone)
		ws.heartbeatDone = nil
	}

	if ws.readMsgDone != nil {
		close(ws.readMsgDone)
		ws.readMsgDone = nil
	}

	if ws.Conn != nil {
		err := ws.Close()
		ws.fireEvent(EventSocketClose, "", err)
	}
}

func (ws *Conn) Shutdown() {
	ws.CloseWs()
	close(ws.eventStream)
	ws.eventStream = nil
	close(ws.writeStream)
	ws.writeStream = nil
}

func (ws *Conn) fireEvent(event int, arg interface{}, err error) {
	if ws.eventStream != nil {
		ws.eventStream <- &Event{
			ws:    ws,
			event: event,
			arg:   arg,
			err:   err,
		}
	}
}

func (ws *Conn) fireMessage(err error, messageType int, message []byte) {
	if ws.eventStream != nil {
		ws.eventStream <- &Event{
			ws:          ws,
			event:       EventMessage,
			messageType: messageType,
			message:     message,
		}
	}
}
