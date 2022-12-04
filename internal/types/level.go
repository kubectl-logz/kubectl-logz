package types

// Level is error, warn, info, debug
type Level string

func (l Level) String() string {
	return string(l)
}

func (l Level) IsZero() bool {
	return l == ""
}

func ParseLevel(l string) Level {
	switch l {
	case "error", "ERROR", "fatal", "FATAL":
		return "error"
	case "warning", "WARNING", "warn", "WARN":
		return "warn"
	case "info", "INFO":
		return "info"
	case "debug", "DEBUG":
		return "debug"
	default:
		return ""
	}
}
