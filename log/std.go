package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync"
)

var _ Logger = (*stdLogger)(nil)

type stdLogger struct {
	log   *log.Logger
	pool  *sync.Pool
	poolm *sync.Pool
}

// NewStdLogger new a logger with writer.
func NewStdLogger(w io.Writer) Logger {
	return &stdLogger{
		log: log.New(w, "", 0),
		pool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
		poolm: &sync.Pool{
			New: func() interface{} {
				return map[string]interface{}{}
			},
		},
	}
}

// Log print the kv pairs log.
func (l *stdLogger) Log(level Level, keyvals ...interface{}) error {
	m := l.poolm.Get().(map[string]interface{})
	defer l.poolm.Put(m)

	if len(keyvals) == 0 {
		return nil
	}
	if (len(keyvals) & 1) == 1 {
		keyvals = append(keyvals, "KEYVALS UNPAIRED")
	}
	m["level"] = level.String()
	for i := 0; i < len(keyvals); i += 2 {
		m[fmt.Sprintf("%s", keyvals[i])] = fmt.Sprintf("%v", keyvals[i+1])

	}
	jsonBytes, err := json.Marshal(m)
	if err != nil {
		l.log.Output(4, err.Error())
		return err
	}

	_ = l.log.Output(4, string(jsonBytes)) //nolint:gomnd

	return nil
}
