package config

import (
	"os"

	cfgutil "github.com/namhoai1109/tabi/core/cfg"
)

type Configuration struct {
	Stage            string   `env:"UP_STAGE"`
	AppName          string   `env:"APP_NAME"`
	Port             int      `env:"PORT" envDefault:"3000"`
	ReadTimeout      int      `env:"READ_TIMEOUT"`
	WriteTimeout     int      `env:"WRITE_TIMEOUT"`
	ReloadConfigTime int      `env:"RELOAD_CONFIG_TIME" envDefault:"300"`
	AllowOrigins     []string `env:"ALLOW_ORIGINS"`
	DbDsn            string   `env:"DB_DSN"`
	DbLog            bool     `env:"DB_LOG"`
	ReqLog           bool     `env:"REQ_LOG"`
	Region           string   `env:"REGION"`
	ConfigStage      string   `env:"CONFIG_STAGE"`

	JwtPartnerSecret    string `env:"JWT_PARTNER_SECRET"`
	JwtPartnerDuration  int    `env:"JWT_PARTNER_DURATION"`
	JwtPartnerAlgorithm string `env:"JWT_PARTNER_ALGORITHM"`
	JwtUserSecret       string `env:"JWT_USER_SECRET"`
	JwtUserDuration     int    `env:"JWT_USER_DURATION"`
	JwtUserAlgorithm    string `env:"JWT_USER_ALGORITHM"`

	EndpointFile string `env:"ENDPOINT_FILE"`
}

func Load() (*Configuration, error) {
	appName := os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
	if configname := os.Getenv("CONFIG_NAME"); configname != "" {
		appName = configname
	}
	stage := os.Getenv("UP_STAGE")

	cfg := new(Configuration)
	cfg.AppName = appName
	cfg.Stage = stage
	if err := cfgutil.LoadWithAPS(cfg, appName, stage); err != nil {
		return nil, err
	}
	return cfg, nil
}
