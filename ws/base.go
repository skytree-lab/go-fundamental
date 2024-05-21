package ws

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/skytree-lab/go-fundamental/util"
)

type Base struct {
	Host               string
	Ws                 *Conn
	EventStream        chan *Event
	MessageHandler     func(msg []byte)
	OnWSCreated        func()
	heartbeatFailCount int
}

func (h *Base) Close() {
	h.Ws.Shutdown()
}

func (h *Base) CreateWsConn(host, path string) {
	h.EventStream = make(chan *Event, 100)
	go h.eventHandler()

	util.Logger().Info("creating websocket connection...")
	h.Ws = NewWsConn(host, path, h.EventStream)
}

func (h *Base) onWSCreated(event *Event) {
	if event.err == nil {
		util.Logger().Info("creating websocket connection OK.")
		if h.OnWSCreated != nil {
			h.OnWSCreated()
		}
		event.ws.StartToReceiveMessage()
	} else {
		msg := fmt.Sprintf("create websocket failed. err = %v, args = %v", event.err, event.arg)
		util.Logger().Error(msg)
		event.ws.Reconnect(event.err)
	}
}

func (h *Base) onHeartBeatStarted(event *Event) {
	if event.err != nil {
		h.heartbeatFailCount++
		if h.heartbeatFailCount == 5 {
			msg := fmt.Sprintf("send heartbeat failed 5 times. arg = %v, err = %v", event.arg, event.err)
			util.Logger().Error(msg)
			event.ws.Reconnect(event.err)
		}
	}
}

func (h *Base) onSubscribed(event *Event) {
	if event.err != nil {
		msg := fmt.Sprintf("subcribe failed. arg = %v, err = %v", event.arg, event.err)
		util.Logger().Error(msg)
		event.ws.Reconnect(event.err)
	}
}

func (h *Base) onReconnect(event *Event) {
	if event.err != nil {
		msg := fmt.Sprintf("reconnect reason err = %v", event.err)
		util.Logger().Error(msg)
	}
}

func (h *Base) onWSClosed(event *Event) {
	if event.err != nil {
		msg := fmt.Sprintf("close websocket failed. arg = %v, err = %v", event.arg, event.err)
		util.Logger().Info(msg)
	}
}

func (h *Base) onMessageArrive(event *Event) {
	if event.err != nil {
		msg := fmt.Sprintf("read failed. arg = %v, err = %v", event.arg, event.err)
		util.Logger().Error(msg)
	} else {
		switch event.messageType {
		case websocket.PongMessage:
			msg := fmt.Sprintf("Base PONG!. err = %v", event.err)
			util.Logger().Error(msg)
		case websocket.CloseMessage:
			msg := fmt.Sprintf("websocket closed. err = %v", event.err)
			util.Logger().Error(msg)
		case websocket.TextMessage, websocket.BinaryMessage:
			h.MessageHandler(event.message)
		}
	}
}

func (h *Base) eventHandler() {
	for event := range h.EventStream {
		switch event.event {
		case EventCreateConnection:
			h.onWSCreated(event)
		case EventHeartBeat:
			h.onHeartBeatStarted(event)
		case EventMessage:
			h.onMessageArrive(event)
		case EventSocketClose:
			h.onWSClosed(event)
		case EventSubscribe:
			h.onSubscribed(event)
		case EventReconnect:
			h.onReconnect(event)
		}
	}
}
