package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Configuration holds the application configuration
type Configuration struct {
	App   AppConfig
	DB    DBConfig
	Log   LogConfig
	JWT   JWTConfig
	Cache CacheConfig
}

// AppConfig holds application-level configuration
type AppConfig struct {
	Name string
	Env  string
	Port string
}

// DBConfig holds database configuration
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// LogConfig holds logger configuration
type LogConfig struct {
	Level  string
	Output string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret string
	Expire time.Duration
}

// CacheConfig holds cache configuration
type CacheConfig struct {
	MaxSize    int
	GCInterval time.Duration
}

var (
	config      *Configuration
	watcherChan chan struct{}
)

// Load loads configuration from file or environment variables
func Load() (*Configuration, error) {
	// Set default values first
	setDefaults()

	// Bind environment variables
	bindEnvs()

	// Load .env file if exists
	if _, err := os.Stat(".env"); err == nil {
		// Load .env file using godotenv
		err := godotenv.Load()
		if err != nil {
			return nil, fmt.Errorf("failed to load .env file: %w", err)
		}
	}

	// Use automatic env
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Unmarshal config
	cfg := &Configuration{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate config
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	config = cfg
	return config, nil
}

// Get returns the loaded configuration
func Get() *Configuration {
	return config
}

// WatcherChan returns a channel that receives notifications when config changes
func WatcherChan() <-chan struct{} {
	if watcherChan == nil {
		watcherChan = make(chan struct{})

		// Set up viper callback for config changes
		viper.OnConfigChange(func(e fsnotify.Event) {
			// Reload configuration
			cfg := &Configuration{}
			if err := viper.Unmarshal(cfg); err == nil {
				if err := cfg.validate(); err == nil {
					config = cfg
					// Notify watchers
					select {
					case watcherChan <- struct{}{}:
					default:
					}
				}
			}
		})
	}
	return watcherChan
}

func setDefaults() {
	viper.SetDefault("app.name", "go-admin")
	viper.SetDefault("app.env", "local")
	viper.SetDefault("app.port", "8080")

	viper.SetDefault("db.host", "localhost")
	viper.SetDefault("db.port", "3306")
	viper.SetDefault("db.user", "root")
	viper.SetDefault("db.password", "password")
	viper.SetDefault("db.name", "go_admin")

	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.output", "console")

	viper.SetDefault("jwt.secret", "go-admin-secret")
	viper.SetDefault("jwt.expire", "24h")

	viper.SetDefault("cache.maxsize", 10000)
	viper.SetDefault("cache.gcinterval", "10m")
}

func bindEnvs() {
	// App config
	viper.BindEnv("app.name", "APP_NAME")
	viper.BindEnv("app.env", "APP_ENV")
	viper.BindEnv("app.port", "APP_PORT")

	// DB config
	viper.BindEnv("db.host", "DB_HOST")
	viper.BindEnv("db.port", "DB_PORT")
	viper.BindEnv("db.user", "DB_USER")
	viper.BindEnv("db.password", "DB_PASSWORD")
	viper.BindEnv("db.name", "DB_NAME")

	// Log config
	viper.BindEnv("log.level", "LOG_LEVEL")
	viper.BindEnv("log.output", "LOG_OUTPUT")

	// JWT config
	viper.BindEnv("jwt.secret", "JWT_SECRET")
	viper.BindEnv("jwt.expire", "JWT_EXPIRE")

	// Cache config
	viper.BindEnv("cache.maxsize", "CACHE_MAXSIZE")
	viper.BindEnv("cache.gcinterval", "CACHE_GCINTERVAL")
}

func (c *Configuration) validate() error {
	if c.App.Name == "" {
		return fmt.Errorf("app.name is required")
	}

	if c.DB.Host == "" {
		return fmt.Errorf("db.host is required")
	}

	if c.JWT.Secret == "" {
		return fmt.Errorf("jwt.secret is required")
	}

	return nil
}

// IsDevelopment checks if the application is running in development mode
func (c *Configuration) IsDevelopment() bool {
	return c.App.Env == "local" || c.App.Env == "dev" || c.App.Env == "development"
}

// DSN returns the database connection string
func (c *DBConfig) DSN() string {
	// Remove surrounding quotes from password if present
	password := c.Password
	if strings.HasPrefix(password, `"`) && strings.HasSuffix(password, `"`) {
		password = strings.Trim(password, `"`)
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, password, c.Host, c.Port, c.Name)
}
