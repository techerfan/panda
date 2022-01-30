package panda

import (
	"sync"

	"github.com/techerfan/panda/logger"
)

type channels struct {
	allChannels map[string]*channel
	logger      logger.Logger
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
	ch := c.addChannel(chName)
	return ch
}

func (c *channels) addChannel(chName string) *channel {
	if ch, ok := c.allChannels[chName]; !ok {
		channel := NewChannel(ch.logger, chName)
		c.allChannels[chName] = channel
		return channel
	} else {
		return ch
	}
}
