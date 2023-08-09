package wsasyncroutine

import "encoding/json"

type AsyncRoutine struct {
	messageHandler *MessageHandler
	channel        chan *MessageLayout
	regIndex       int
}

func (a *AsyncRoutine) Send(in MessageData, out MessageData) error {
	a.messageHandler.Send(in)
	a.regIndex = out.Type()
	err := a.messageHandler.Register(out.Type(), func(layout *MessageLayout) error {
		println("a.channel <- layout", string(layout.Message), a.channel)
		a.channel <- layout
		println("1 a.channel <- layout")

		e := a.messageHandler.Unregister(a.regIndex)
		if e != nil {
			return e
		}
		return nil
	})
	if err != nil {
		return err
	}
	println("register ", a.regIndex)
	println("register ", a.messageHandler.handlers[1])
	println("pre")
	outLayout := <-a.channel
	println("post")

	err = json.Unmarshal(outLayout.Message, &out)
	if err != nil {
		return err
	}
	return nil
}
