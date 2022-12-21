package data

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/zxq97/gotool/constant"
	"github.com/zxq97/gotool/generate"
	"github.com/zxq97/gotool/kafka"
)

func sendKafkaMsg(ctx context.Context, producer *kafka.Producer, topic, key string, req proto.Message, eventtype int32) error {
	trace, ok := ctx.Value(constant.TraceIDKey).(string)
	if !ok {
		trace = generate.UUIDStr()
	}
	bs, err := proto.Marshal(req)
	if err != nil {
		return err
	}
	kfkmsg := &kafka.KafkaMessage{
		TraceId:   trace,
		EventType: eventtype,
		Message:   bs,
	}
	bs, err = proto.Marshal(kfkmsg)
	if err != nil {
		return err
	}
	return producer.SendMessage(topic, []byte(key), bs)
}
