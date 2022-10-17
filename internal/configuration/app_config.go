package configuration

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	json "github.com/json-iterator/go"
	"github.com/spf13/viper"
)

const (
	EnvProduction      = "Production"
	URLShortenerPrefix = "URL_SHORTENER"
)

var ErrUnmarshalConfig = errors.New("viper failed to unmarshal app config")

type Configuration struct {
	/* ---------------------------  HTTP  ----------------------------------- */

	Addr string `mapstructure:"address"`

	/* ---------------------------  API  ----------------------------------- */

	API API `mapstructure:"api"`

	/* ---------------------------  Metrics (Prometheus)  ----------------------- */

	Metrics Metrics `mapstructure:"metrics"`

	/* ---------------------------  Zap Logger  ----------------------------------- */

	// Filter parameter for the leveled logger. Valid values: production, development, none
	ZapLoggerMode string `mapstructure:"zap_logger_mode"`

	/* ---------------------------  Redis  ------------------------------------- */

	Redis Redis `mapstructure:"redis"`

	/* ---------------------------  Base URL  ------------------------------------- */
}

type API struct {
	BaseURL string `mapstructure:"base_url"`
}

type Metrics struct {
	Addr      string `mapstructure:"addr"`
	Namespace string `mapstructure:"namespace"`
	Subsystem string `mapstructure:"subsystem"`
}

type Redis struct {
	Host string `mapstructure:"host"`
}

type ProductionConfigurationLogging struct {
	Level   string    `json:"level"`
	TS      time.Time `json:"ts"`
	Msg     string    `json:"msg"`
	Payload struct {
		Configuration *Configuration `json:"configuration"`
	} `json:"payload"`
}

func NewAppConfiguration(env string, writeConfig bool) (cfg *Configuration, err error) {
	var filename string

	switch env {
	case EnvProduction:
		filename = "url_shortener.settings"
	default:
		filename = "url_shortener.settings.development"
	}

	v := newViper(filename)

	cfg, err = unmarshalConfig(v)
	if err != nil {
		return
	}

	switch env {
	case EnvProduction:
		data, errMarshal := json.ConfigCompatibleWithStandardLibrary.Marshal(ProductionConfigurationLogging{
			Level: "info",
			TS:    time.Now().UTC(),
			Msg:   "resolved_configuration",
			Payload: struct {
				Configuration *Configuration `json:"configuration"`
			}{
				Configuration: cfg,
			},
		})

		if errMarshal != nil {
			return nil, errMarshal
		}

		fmt.Printf("%s\n", data)
	default:
		fmt.Println("Logging the resolved configuration:")
		_, _ = fmt.Println(cfg)
	}

	if writeConfig {
		if err = v.WriteConfig(); err != nil {
			log.Println("viper failed to write app config file:", err)
		}
	}

	return cfg, nil
}

// Set the default config values for the viper object we are using.
// nolint: funlen
func newViper(filename string) *viper.Viper {
	v := viper.New()

	if filename != "" {
		v.SetConfigName(filename)
		v.AddConfigPath("./")
	}

	{
		/* ---------------------------  Transport  -------------------------------- */

		v.SetDefault("address", ":8081")
	}
	{
		/* ---------------------------  Transport  -------------------------------- */

		v.SetDefault("api.base_url", "127.0.0.1:8081")
	}
	{
		/* ---------------------------  Metrics (Prometheus)  --------------------- */
		{
			v.SetDefault("metrics.addr", ":5070")
			v.SetDefault("metrics.namespace", "app")
			v.SetDefault("metrics.subsystem", "url_shortener")
		}
	}
	{
		/* ---------------------------  Zap Logger  ------------------------------- */

		v.SetDefault("zap_logger_mode", "development")
	}
	{
		/* ---------------------------  Redis  ------------------------------------- */
		{
			v.SetDefault("redis.host", ":6379")
		}
	}

	// Set environment variable support:
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetEnvPrefix(URLShortenerPrefix)
	v.AutomaticEnv()

	// ReadInConfig will discover and load the configuration file from disk
	// and key/value stores, searching in one of the defined paths.
	if err := v.ReadInConfig(); err != nil {
		log.Println("viper failed to read app config file:", err)
	}

	return v
}

// unmarshalConfig uses viper to get app configuration.
func unmarshalConfig(v *viper.Viper) (*Configuration, error) {
	var c *Configuration

	if err := v.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnmarshalConfig, err)
	}

	return c, nil
}
