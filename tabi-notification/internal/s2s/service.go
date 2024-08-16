package s2s

import (
	"tabi-notification/config"

	s2s "github.com/namhoai1109/tabi/core/s2s"
)

func New(cfg *config.Configuration) *S2S {
	return &S2S{
		cfg: cfg,
		s2s: s2s.New(s2s.Config{
			JwtAlgorithm: cfg.JwtPartnerAlgorithm,
			JwtSecret:    cfg.JwtPartnerSecret,
			JwtDuration:  cfg.JwtPartnerDuration,
			Debug:        cfg.DbLog,
			Timeout:      30,
			Duration:     300,
			Region:       cfg.Region,
		}),
	}
}

type S2S struct {
	cfg *config.Configuration
	s2s s2s.Intf
}
