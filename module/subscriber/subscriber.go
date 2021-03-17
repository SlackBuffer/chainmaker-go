/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package subscriber

import (
	"chainmaker.org/chainmaker-go/common/msgbus"
	commonPb "chainmaker.org/chainmaker-go/pb/protogo/common"
	"chainmaker.org/chainmaker-go/subscriber/model"
	feed "github.com/ethereum/go-ethereum/event"
)

// EventSubscriber - new EventSubscriber struct
type EventSubscriber struct {
	blockFeed feed.Feed
}

// OnMessage - deal msgbus.BlockInfo message
func (s *EventSubscriber) OnMessage(msg *msgbus.Message) {
	if blockInfo, ok := msg.Payload.(*commonPb.BlockInfo); ok {
		go s.blockFeed.Send(model.NewBlockEvent{BlockInfo: blockInfo})
	}
}

// OnQuit - deal msgbus OnQuit message
func (s *EventSubscriber) OnQuit() {
	// do nothing
}

// NewSubscriber - new and register msgbus.BlockInfo object
func NewSubscriber(msgBus msgbus.MessageBus) *EventSubscriber {
	subscriber := &EventSubscriber{}
	msgBus.Register(msgbus.BlockInfo, subscriber)
	return subscriber
}

// SubscribeBlockEvent - subscribe block event
func (s *EventSubscriber) SubscribeBlockEvent(ch chan<- model.NewBlockEvent) feed.Subscription {
	return s.blockFeed.Subscribe(ch)
}