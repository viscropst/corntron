package log

import "github.com/rs/zerolog"

func (l Level) toZeroLogLevel() zerolog.Level {
	switch l {
	case DebugLevel:
		return zerolog.DebugLevel
	case ErrorLevel:
		return zerolog.ErrorLevel
	case WarnLevel:
		return zerolog.WarnLevel
	case TraceLevel:
		return zerolog.TraceLevel
	case InfoLevel:
		return zerolog.InfoLevel
	case Disabled:
		return zerolog.Disabled
	case FatalLevel:
		return zerolog.FatalLevel
	case PanicLevel:
		return zerolog.PanicLevel
	case NoLevel:
		fallthrough
	default:
		return zerolog.NoLevel

	}
}

type zeroLogger struct {
	logger zerolog.Logger
	level  Level
}

func (l zeroLogger) With(options ...zerolog.Context) zeroLogger {
	tmp := &l
	tmp.logger = options[0].Logger()
	return *tmp
}

func (l zeroLogger) Write(p []byte) (n int, err error) {
	n = len(p)
	if n > 0 && p[n-1] == '\n' {
		// Trim CR added by stdlog.
		p = p[0 : n-1]
	}
	l.logger.WithLevel(l.level.toZeroLogLevel()).Msg(string(p))
	return
}

func ZeroLogger(level ...Level) zeroLogger {
	zLog := zeroLogger{}
	if len(level) > 0 {
		zLog.level = level[0]
	} else {
		zLog.level = NoLevel
	}
	cw := zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.PartsExclude = []string{zerolog.TimestampFieldName}
		if len(level) == 0 {
			w.PartsExclude = append(w.PartsExclude, zerolog.LevelFieldName)
		}
	})
	zerolog.FormattedLevels[zerolog.ErrorLevel] = "error:"
	zerolog.FormattedLevels[zerolog.WarnLevel] = "WARN:>"
	zerolog.FormattedLevels[zerolog.FatalLevel] = "fatal:"
	zerolog.FormattedLevels[zerolog.DebugLevel] = "debug:"
	zerolog.SetGlobalLevel(CLIOutputLevel.toZeroLogLevel())
	zLog.logger = zerolog.New(cw)
	return zLog
}
