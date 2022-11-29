package loggernr

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type LogMessage struct {
	GlobalEventTimestamp string `json:"global_event_timestamp"`
	GlobalEventName      string `json:"global_event_name,omitempty"`
	Level                string `json:"level"`
	Context              string `json:"context,omitempty"`
	Message              string `json:"message"`
	ServiceName          string `json:"service_name"`
	SessionId            string `json:"session_id,omitempty"`
	TraceId              string `json:"trace_id,omitempty"`
}

const (
	EMERGENCY string = "EMERGENCY"
	ERROR     string = "ERROR"
	WARN      string = "WARN"
	INFO      string = "INFO"
	DEBUG     string = "DEBUG"
	TRACE     string = "TRACE"
)

func (m *LogMessage) SendLog(level string, message string, object interface{}) {

	//Chave para desligar ou ligar o modo debug
	debugMode := true
	if os.Getenv("ENVIRONMENT_NAME") == "production" {
		debugMode = false
	}

	m.GlobalEventTimestamp = time.Now().Format(time.RFC3339)
	m.Level = level
	m.Message = message

	objectJson, err := json.Marshal(object)
	if err != nil {
		fmt.Println("Message not marshal: ", err.Error())
	}
	var trace []string

	switch m.Level {
	case EMERGENCY,
		ERROR,
		WARN,
		INFO:
		trace = getCallerTrace(2)

	case DEBUG:
		if debugMode {
			trace = getCallerTrace(2)
		}

	case TRACE:
		if debugMode {
			trace = getCallerTrace(-1)
		}
	}

	m.Context = fmt.Sprintf(`{"trace":"%s","object":"%s"}`, trace, objectJson)

	j, err := json.Marshal(m)
	if err != nil {
		fmt.Println("Message not marshal: ", err.Error())
	}

	fmt.Println(string(j))
}

func getCallerTrace(caller int) []string {
	var trace []string
	if caller == -1 {
		ok := true
		for i := 2; ok; i++ {
			pc, file, line, o := runtime.Caller(i)
			ok = o
			if ok {
				path := strings.Split(file, "/")
				pasta := path[len(path)-2]
				arquivo := path[len(path)-1]
				f := runtime.FuncForPC(pc)

				trace = append(trace, fmt.Sprintf(`{"source":"%s/%s:%s::%s"}`, pasta, arquivo, strconv.Itoa(line), f.Name()))
			}
		}
		return trace
	}

	pc, file, line, ok := runtime.Caller(caller) // pc, file, line, ok
	if !ok {
		fmt.Println("Caller not found")
		return []string{}
	}

	path := strings.Split(file, "/")
	pasta := path[len(path)-2]
	arquivo := path[len(path)-1]
	f := runtime.FuncForPC(pc)
	trace = append(trace, fmt.Sprintf(`{"source":"%s/%s:%s::%s"}`, pasta, arquivo, strconv.Itoa(line), f.Name()))

	return trace
}
