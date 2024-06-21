package config_center

import (
	"github.com/Jayleonc/register/internal/core/configuration"
	"time"
)

type EtcdConfigCenterOptions = configuration.EtcdConfigCenterOptions

func DefaultEtcdConfigCenterOptions() EtcdConfigCenterOptions {
	return configuration.DefaultEtcdConfigCenterOptions()
}

func WithEtcdAddress(address string) ClientOption {
	return func(o *configuration.EtcdConfigCenterOptions) {
		o.EtcdAddress = address
	}
}

func WithCredentials(username, password string) ClientOption {
	return func(o *configuration.EtcdConfigCenterOptions) {
		o.Username = username
		o.Password = password
	}
}

func WithDialTimeout(timeout time.Duration) ClientOption {
	return func(o *configuration.EtcdConfigCenterOptions) {
		o.DialTimeout = timeout
	}
}

func WithLogLevel(level string) ClientOption {
	return func(o *configuration.EtcdConfigCenterOptions) {
		o.LogLevel = level
	}
}

func WithTLS(certFile, keyFile, caFile string) ClientOption {
	return func(o *configuration.EtcdConfigCenterOptions) {
		o.TLSCertFile = certFile
		o.TLSKeyFile = keyFile
		o.TLSCaFile = caFile
	}
}

func WithRetry(attempts int, interval time.Duration) ClientOption {
	return func(o *configuration.EtcdConfigCenterOptions) {
		o.RetryAttempts = attempts
		o.RetryInterval = interval
	}
}
