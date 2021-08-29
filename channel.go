package panda

type channel struct {
	name    string
	clients []*client
}

func newChannel(name string) *channel {
	channel := &channel{
		name: name,
	}
	return channel
}

func (c *channel) addClient() {

}

func (c *channel) removeClient() {

}

func (c *channel) sendMessage(msg string) {

}

func (c *channel) destroy() {
	c = nil
}
