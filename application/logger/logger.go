package logger

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"serverless-go-template/application/elastic"
)

// Level used to filter log message by the Logger.
type Level string

// logging levels.
const (
	PanicLevel Level = "PANIC"
	FatalLevel Level = "FATAL"
	ErrorLevel Level = "ERROR"
	WarnLevel  Level = "WARN"
	InfoLevel  Level = "INFO"
	DebugLevel Level = "DEBUG"
)

type Config struct {
	AppEnv      string
	ElasticVars elastic.ElasticConfig
}

// Logger has the configurations required to send logs. It can be created using the NewLogger function
type Logger struct {
	// LoggerConfig settings and vars to connect elastic
	LoggerConfig Config
}

// NewLogger creates a logger instance
func NewLogger(config Config) *Logger {
	return &Logger{
		LoggerConfig: config,
	}
}

//type Log struct {
//	// Level used to filter log message
//	Level Level
//	// Index index name if send to elastic
//	Index   string
//	Message string
//	// JsonData Body log to send
//	JsonData string
//	// SendConsole boolean to define print log in console
//	SendConsole bool
//	// SendConsole boolean to define send log to elastic
//	SendElastic bool
//	// SendConsole boolean to define send log to Kinesis
//	SendKinesis bool
//
//	file string
//	line int
//}

type Log struct {
	// Level used to filter log message
	Level Level
	// Index index name if send to elastic
	Index   string
	Message string
	// JsonData Body log to send
	JsonData string

	Command         string
	Method          string
	Type            string
	Id              string
	PayloadJson     string
	Response        string
	Extra           string
	ReturnLevelCode int
	ReturnLevelName string
	// SendConsole boolean to define print log in console
	SendConsole bool
	// SendConsole boolean to define send log to elastic
	SendElastic bool
	// SendConsole boolean to define send log to Kinesis
	SendKinesis bool

	file string
	line int
}

// Log print and send message log
// SendConsole log has printed in console
// SendElastic sends logs to ElasticSearch
func (l *Logger) Log(lo Log) {
	_, lo.file, lo.line, _ = runtime.Caller(1)

	if lo.SendConsole {
		if len(lo.PayloadJson) > 0 {
			log.Printf("%s | %s (%d) | %s | %s | %s", lo.Level, lo.file, lo.line, lo.Message, lo.PayloadJson,
				lo.Response)
			// log.Println(reflect.TypeOf(lo.PayloadJson))
			fmt.Println(string(lo.PayloadJson))
		} else {
			log.Printf("%s | %s (%d) | %s", lo.Level, lo.file, lo.line, lo.Message)
		}
	}

	if lo.SendElastic {
		l.sendElastic(&lo)
	}

	if lo.SendKinesis {
		l.sendKinesis(lo.JsonData)
	}
}

// sendElastic connects to elasticHost and send bulk logs
func (l *Logger) sendElastic(lo *Log) {
	if len(l.LoggerConfig.ElasticVars.Hosts) <= 0 {
		log.Println("Hosts undefined to send elastic")
		return
	}

	elasticConnection, err := elastic.NewElasticSearch(l.LoggerConfig.ElasticVars)
	if err != nil {
		log.Println("Error to conect elastic", err)
		return
	}

	elasticConnection.AddLog(elastic.ElasticDocs{
		Channel:   "Middle Earth",
		Extra:     lo.Extra,
		Level:     lo.ReturnLevelCode,
		LevelName: lo.ReturnLevelName,
		Message:   lo.Message,
		Context: elastic.Context{lo.Command, lo.Method, lo.Type, lo.Id,
			lo.PayloadJson, lo.Response},

		//Summary: "["+string(lo.Level)+"] "+lo.Message,
		//Body:    lo.JsonData,
	})

	go elasticConnection.SendLogs(context.Background(), lo.Index, "index", "_doc")
}

func (l *Logger) sendKinesis(jsonData string) {
	//go kinesisLog(jsonData)
}
