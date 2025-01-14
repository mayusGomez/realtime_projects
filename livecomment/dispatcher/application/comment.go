package application

import (
	"livecomments/dispatcher/domain"
	"log"
)

type Comment struct {
	gatewayConfig domain.GatewayConfig
	asyncComm     domain.AsyncCommunication
}

func NewComment(gatewayConfig domain.GatewayConfig, asyncComm domain.AsyncCommunication) *Comment {
	return &Comment{gatewayConfig: gatewayConfig, asyncComm: asyncComm}
}

func (g *Comment) PostComment(connectionId, video, comment string) error {
	queues := g.gatewayConfig.GetQueues(video)
	log.Printf("post comment: %s video: %s queues: %v", comment, video, queues)

	message := domain.CommentMessage{
		ConnectionId: connectionId,
		Video:        video,
		Message:      comment,
	}
	err := g.asyncComm.PostMessage(queues, &message)
	if err != nil {
		return err
	}

	return nil
}
