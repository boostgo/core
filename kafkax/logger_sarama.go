package kafkax

import (
	"fmt"
	"strings"

	"github.com/boostgo/core/convert"
	"github.com/boostgo/core/log"

	"github.com/IBM/sarama"
)

type saramaLogger struct{}

// NewSaramaLogger create custom logger for "sarama" library for debugging
func NewSaramaLogger() sarama.StdLogger {
	return &saramaLogger{}
}

func (l *saramaLogger) Print(args ...interface{}) {
	values := make([]string, 0, len(args))
	for _, arg := range args {
		values = append(values, convert.String(arg))
	}

	log.
		Info().
		Str("logger", "SARAMA").
		Strs("args", values).
		Msg("Sarama logger Print")
}

func (l *saramaLogger) Printf(format string, v ...interface{}) {
	message := strings.Builder{}
	_, _ = fmt.Fprintf(&message, format, v...)

	log.
		Info().
		Str("logger", "SARAMA").
		Msg(fmt.Sprintf(format, v...))
}

func (l *saramaLogger) Println(args ...interface{}) {
	values := make([]string, 0, len(args))
	for _, arg := range args {
		values = append(values, convert.String(arg))
	}

	log.
		Info().
		Str("logger", "SARAMA").
		Strs("args", values).
		Msg("Sarama logger Println")
}
