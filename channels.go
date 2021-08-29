package panda

import (
	"sync"
)

type channels struct {
	allChannels map[string]*channel
}

var lock = &sync.Mutex{}
var channelsInstance *channels

func getChannelsInstance() *channels {
	if channelsInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if channelsInstance == nil {
			channelsInstance = &channels{
				allChannels: make(map[string]*channel),
			}
		}
	}
	return channelsInstance
}

func (c *channels) getChannelByName(chName string) *channel {
	if ch, ok := c.allChannels[chName]; ok {
		return ch
	}
	return nil
}

func (c *channels) addChannel(chName string) {
	if _, ok := c.allChannels[chName]; !ok {
		// channel := panda.NewChannel(chName)
		// c.allChannels[chName] = channel
	}
}
