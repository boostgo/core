package kafkax

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/boostgo/core/convert"
	"github.com/boostgo/core/reflectx"
	"github.com/boostgo/core/trace"
	"github.com/boostgo/core/validator"

	"github.com/IBM/sarama"
)

// Parse message body to provided export object (which must be ptr) and validate for "validate" tags.
func Parse(message *sarama.ConsumerMessage, export any) error {
	// check export type
	if err := reflectx.CheckExport(export); err != nil {
		return err
	}

	// parse message
	if err := json.Unmarshal(message.Value, export); err != nil {
		return err
	}

	// validate struct
	if err := validator.Get().Struct(export); err != nil {
		return err
	}

	return nil
}

// Header search header in provided message by header name.
func Header(message *sarama.ConsumerMessage, name string) string {
	nameBlob := convert.BytesFromString(name)

	for _, header := range message.Headers {
		if bytes.Equal(header.Key, nameBlob) {
			return convert.StringFromBytes(header.Value)
		}
	}

	return ""
}

// Headers returns all headers from message as map and [param.Param] object
func Headers(message *sarama.ConsumerMessage) map[string]string {
	headers := make(map[string]string, len(message.Headers))
	for _, header := range message.Headers {
		headers[string(header.Key)] = convert.StringFromBytes(header.Value)
	}
	return headers
}

// SetHeaders convert provided headers map to sarama headers slice
func SetHeaders(headers map[string]any) []sarama.RecordHeader {
	messageHeaders := make([]sarama.RecordHeader, len(headers))

	for name, value := range headers {
		messageHeaders = append(messageHeaders, sarama.RecordHeader{
			Key:   convert.BytesFromString(name),
			Value: convert.Bytes(value),
		})
	}

	return messageHeaders
}

func setTrace(ctx context.Context, messages ...*sarama.ProducerMessage) {
	traceID := trace.Get(ctx)
	if traceID == "" {
		return
	}

	for _, message := range messages {
		message.Headers = append(message.Headers, sarama.RecordHeader{
			Key:   convert.Bytes(TraceKey),
			Value: convert.Bytes(traceID),
		})
	}
}
