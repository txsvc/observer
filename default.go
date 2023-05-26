package observer

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/txsvc/stdlib/v2/deprecated/provider"
)

type (
	// defaultObserverImpl provides a simple implementation in the absence of any configuration
	defaultObserverImpl struct {
		loggingDisabled bool
	}
)

var (
	// Interface guards.

	// This enforces a compile-time check of the provider implmentation,
	// making sure all the methods defined in the interfaces are implemented.

	_ provider.GenericProvider = (*defaultObserverImpl)(nil)

	_ ErrorReportingProvider = (*defaultObserverImpl)(nil)
	_ LoggingProvider        = (*defaultObserverImpl)(nil)
	_ MetricsProvider        = (*defaultObserverImpl)(nil)

	// the instance, a singleton
	theDefaultProvider *defaultObserverImpl
)

func init() {
	Init()
}

func Init() {
	// force a reset
	theDefaultProvider = nil

	// initialize the observer with a NULL provider that prevents NPEs in case someone forgets to initialize the platform with a real provider
	loggingConfig := provider.WithProvider("observer.default.logger", TypeLogger, NewDefaultProvider)
	errorReportingConfig := provider.WithProvider("observer.default.errorreporting", TypeErrorReporter, NewDefaultProvider)
	metricsConfig := provider.WithProvider("observer.default.metrics", TypeMetrics, NewDefaultProvider)

	NewConfig(loggingConfig, errorReportingConfig, metricsConfig)
}

// a default provider that does nothing but prevents NPEs in case someone forgets to actually initializa the 'real' provider
func NewDefaultProvider() interface{} {
	if theDefaultProvider == nil {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		//zerolog.SetGlobalLevel(zerolog.InfoLevel)

		theDefaultProvider = &defaultObserverImpl{
			loggingDisabled: false,
		}
		theDefaultProvider.EnableLogging()
	}
	return theDefaultProvider
}

// IF GenericProvider

func (np *defaultObserverImpl) Close() error {
	return nil
}

// IF ErrorReportingProvider

func (np *defaultObserverImpl) ReportError(e error) error {
	log.Error().Stack().Err(e).Msg("")
	return e
}

// IF LoggingProvider

func (np *defaultObserverImpl) Log(msg string, keyValuePairs ...string) {
	if np.loggingDisabled {
		return // just do nothing
	}
	np.LogWithLevel(LevelInfo, msg, keyValuePairs...)
}

func (np *defaultObserverImpl) _LogWithLevel(lvl Severity, msg string, keyValuePairs ...string) {
	if np.loggingDisabled {
		return // just do nothing
	}

	var kv *zerolog.Array
	if len(keyValuePairs) > 0 {
		kv = zerolog.Arr()
		for i := range keyValuePairs {
			kv = kv.Str(keyValuePairs[i])
		}
	}

	switch lvl {
	case LevelDebug:
		if kv != nil {
			log.Debug().Array(ValuesLogId, kv).Msg(msg)
		} else {
			log.Debug().Msg(msg)
		}
	case LevelInfo:
		if kv != nil {
			log.Info().Array(ValuesLogId, kv).Msg(msg)
		} else {
			log.Info().Msg(msg)
		}
	case LevelNotice:
		if kv != nil {
			log.Info().Array(ValuesLogId, kv).Msg(msg)
		} else {
			log.Info().Msg(msg)
		}
	case LevelWarn:
		if kv != nil {
			log.Warn().Array(ValuesLogId, kv).Msg(msg)
		} else {
			log.Warn().Msg(msg)
		}
	case LevelError:
		if kv != nil {
			log.Error().Array(ValuesLogId, kv).Msg(msg)
		} else {
			log.Error().Msg(msg)
		}
	case LevelAlert:
		if kv != nil {
			log.Fatal().Array(ValuesLogId, kv).Msg(msg)
		} else {
			log.Fatal().Msg(msg)
		}
	}
}

func (np *defaultObserverImpl) LogWithLevel(lvl Severity, msg string, keyValuePairs ...string) {
	if np.loggingDisabled {
		return // just do nothing
	}

	var kv *zerolog.Array
	if len(keyValuePairs) > 0 {
		kv = zerolog.Arr()
		for i := range keyValuePairs {
			kv = kv.Str(keyValuePairs[i])
		}
	}

	var e *zerolog.Event

	switch lvl {
	case LevelDebug:
		e = log.Debug()
	case LevelInfo:
		e = log.Info()
	case LevelNotice:
		e = log.Trace()
	case LevelWarn:
		e = log.Warn()
	case LevelError:
		e = log.Error()
	case LevelAlert:
		e = log.Error()
	}

	if kv != nil {
		e.Array(ValuesLogId, kv).Msg(msg)
	} else {
		e.Msg(msg)
	}
}

func (np *defaultObserverImpl) EnableLogging() {
	np.loggingDisabled = false
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
}

func (np *defaultObserverImpl) DisableLogging() {
	np.loggingDisabled = true
	zerolog.SetGlobalLevel(zerolog.Disabled)
}

// IF MetricsProvider

func (np *defaultObserverImpl) Meter(ctx context.Context, metric string, vals ...string) {
	kv := zerolog.Arr()
	if len(vals) > 0 {
		for i := range vals {
			kv = kv.Str(vals[i])
		}
	}
	log.Trace().Array(ValuesLogId, kv).Str(MetricsLogId, metric).Msg("")
}
