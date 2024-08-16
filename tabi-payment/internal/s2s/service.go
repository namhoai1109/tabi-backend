package s2s

import (
	"tabi-payment/config"

	"github.com/namhoai1109/tabi/core/middleware/jwt"
	"github.com/namhoai1109/tabi/core/s2s"
)

func New(cfg *config.Configuration, jwt *jwt.Service) *S2S {
	return &S2S{
		cfg: cfg,
		jwt: jwt,
		s2s: s2s.New(s2s.Config{
			JwtAlgorithm: cfg.JwtPartnerAlgorithm,
			JwtSecret:    cfg.JwtPartnerSecret,
			JwtDuration:  cfg.JwtPartnerDuration,
			Region:       cfg.Region,
			Timeout:      20,
			Duration:     300,
			Debug:        cfg.ReqLog,
		}),
	}
}

type S2S struct {
	cfg *config.Configuration
	jwt *jwt.Service
	s2s Service
}

type Service interface {
	s2s.Intf
}
