package wshandler

import (
	"encoding/json"
	"errors"

	"github.com/gorilla/websocket"
)

func New(conn *websocket.Conn, length int) *MessageHandler {
	return &MessageHandler{
		handlers: make([]*handler, length, length),
		conn:     conn,
	}
}

type MessageHandler struct {
	handlers []*handler
	conn     *websocket.Conn
}
type handler struct {
	callback func(layout *MessageLayout) error
}

func (m *MessageHandler) HandleMessages(data []byte) error {
	println("received ", string(data))

	mLayout := MessageLayout{}
	err := json.Unmarshal(data, &mLayout)
	if err != nil {
		return err
	}
	handler := m.handlers[mLayout.Type]
	println("0")
	if handler == nil {
		return errors.New("handler does exiest")
	}
	println("1")

	err = handler.callback(&mLayout)
	if err != nil {
		println(err.Error())
		println("error")
		return err
	}

	println("handler.callback")
	return nil
}

func (m *MessageHandler) Register(index int, callback func(layout *MessageLayout) error) error {
	if m.handlers[index] != nil {
		return errors.New("already exist")
	}

	m.handlers[index] = &handler{callback: callback}
	return nil
}
func (m *MessageHandler) Unregister(index int) error {
	if m.handlers[index] == nil {
		return errors.New("doesnt exist")
	}

	m.handlers[index] = nil
	return nil
}
func (m MessageHandler) CreateRoutine() *AsyncRoutine {
	return &AsyncRoutine{
		messageHandler: &m,
		channel:        make(chan *MessageLayout),
	}
}

func (m *MessageHandler) Send(in MessageData) error {
	inLayout := MessageLayout{
		Type: in.Type(),
	}
	buff, err := json.Marshal(in)
	if err != nil {
		return err
	}
	inLayout.Message = buff

	buff, err = json.Marshal(inLayout)
	if err != nil {
		return err
	}
	err = m.conn.WriteMessage(websocket.TextMessage, buff)
	if err != nil {
		return err
	}
	return nil
}

type MessageLayout struct {
	Type    int
	Message json.RawMessage
}
type MessageData interface {
	Type() int
}
