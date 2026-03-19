package tracer

type Config struct {
	ServiceName string
	EndPoint    string
	Insecure    bool
	SampleRatio float64
}

type Option func(*Config)

func WithServiceName(name string) Option {
	return func(c *Config) {
		c.ServiceName = name
	}
}

func WithEndPoint(endpoint string) Option {
	return func(c *Config) {
		c.EndPoint = endpoint
	}
}

func WithInsecure(insecure bool) Option {
	return func(c *Config) {
		c.Insecure = insecure
	}
}

func WithSampleRation(sampleRation float64) Option {
	return func(c *Config) {
		c.SampleRatio = sampleRation
	}
}
