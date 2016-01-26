package discovery

import "time"

type (
	Config struct {
		NetworkID     string
		ListenPort    int
		RequestPeriod time.Duration
		Output        chan string
	}
)
