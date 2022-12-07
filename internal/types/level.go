package types

// Level is error, warn, info, debug
type Level string

func (l Level) String() string {
	return string(l)
}

func (l Level) IsZero() bool {
	return l == ""
}

func (l Level) Less(level Level) bool {
	return l.intVal() < level.intVal()
}

func (l Level) intVal() int {
	switch l {
	case "error":
		return 4
	case "warn":
		return 3
	case "debug":
		return 1
	default:
		return 2
	}
}

func ParseLevel(l string) Level {
	switch l {
	case "error", "ERROR", "fatal", "FATAL", "critical", "CRITICAL":
		return "error"
	case "warning", "WARNING", "warn", "WARN":
		return "warn"
	case "info", "INFO":
		return "info"
	case "debug", "DEBUG", "trace", "TRACE":
		return "debug"
	default:
		return ""
	}
}
