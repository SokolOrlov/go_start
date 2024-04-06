package ports

import (
	"context"
	common "qwe/pkg/models"
)

type IService interface {
	//обработать входящее сообщение
	HandleConsumerMessage(context.Context, common.KafkaMessage) error
}
