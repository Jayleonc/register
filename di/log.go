package di

import (
	"Jayleonc/gateway/internal/client/log"
	"Jayleonc/gateway/pkg/netx"
)

func InitLogClient(httpClient netx.HTTPClientI) log.ClientI {
	client := log.Client{
		TargetUrl:  "http://gateway-service-url/ser-logs/api/logs",
		AppName:    "your-app-name",
		ApiKey:     "your-api-key",
		HttpClient: httpClient,
	}
	return log.NewLogClient(client)
}

func InitLogSender(logClient log.ClientI) log.SenderI {
	return log.NewLogSender(logClient)
}

func InitLogger(logSender log.SenderI) log.Logger {
	return log.NewLogger(logSender, netx.GetOutboundIP(), "dev")
}
