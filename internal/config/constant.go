package config

// Acceptable logger levels
const (
	// TraceLevel = "trace"
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	// PanicLevel = "panic"
	// FatalLevel = "fatal"
)

// Acceptable logger formats
const (
	TextFormat = "text"
	JsonFormat = "json"
)

// Acceptable ssl mods for database
const (
	SSLDisable = "disable"
	SSLRequire = "require"
	// SSLVerifyCA   = "verify-ca"
	// SSLVerifyFull = "verify-full"
)
