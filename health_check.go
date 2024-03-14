package golaze

type HealthCheckConfig struct {
	Port           string
	LivenessHooks  []func() error
	ReadinessHooks []func() error
}

type HealthCheck struct {
	*HealthCheckConfig
}

func NewHealthCheck(config *HealthCheckConfig) *HealthCheck {
	if config.Port == "" {
		config.Port = "8081"
	}

	return &HealthCheck{
		config,
	}
}
