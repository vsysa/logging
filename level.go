package logging

type Level = int32

const (
	TraceLevel Level = 4
	DebugLevel Level = 8
	InfoLevel  Level = 12
	WarnLevel  Level = 16
	ErrorLevel Level = 20
	FatalLevel Level = 24
)
