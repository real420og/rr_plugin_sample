package plugin

import (
	"github.com/roadrunner-server/api/v2/plugins/config"
	"github.com/roadrunner-server/errors"
	"go.uber.org/zap"
	"net/http"
)

const name = "custom_plugin"

type Plugin struct {
	log *zap.Logger
	cfg *Config
}

func (p *Plugin) Init(cfg config.Configurer, log *zap.Logger) error {
	if !cfg.Has(name) {
		return errors.E(errors.Disabled)
	}

	p.cfg = &Config{}
	err := cfg.UnmarshalKey(name, p.cfg)
	if err != nil {
		return err
	}

	p.log = new(zap.Logger)
	*p.log = *log

	return nil
}

func (p *Plugin) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rd := NewResponseDecorator(w)

		rd.Header().Set("Hello", "World")

		next.ServeHTTP(rd, r)

		p.log.Info(
			"response info",
			zap.Int("code", rd.Code()),
			zap.Binary("body", rd.OriginalBody()),
			zap.Int("size", rd.Size()),
		)
	})
}

func (p *Plugin) Serve() chan error {
	errCh := make(chan error, 1)

	p.cfg.InitDefaults()

	p.log.Info(p.cfg.Say)

	return errCh
}

func (p *Plugin) Stop() error {
	return nil
}

func (p *Plugin) Name() string {
	return name
}
