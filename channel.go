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

func (c *channel) sendMessage() {

}
