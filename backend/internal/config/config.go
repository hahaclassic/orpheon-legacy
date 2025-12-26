package config

import (
	"log"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPConfig struct {
	Host string `env:"HOST"`
	Port string `env:"PORT"`
}

type PostgresConfig struct {
	Host         string        `env:"POSTGRES_HOST"`
	Port         string        `env:"POSTGRES_PORT"`
	User         string        `env:"POSTGRES_USER"`
	Password     string        `env:"POSTGRES_PASSWORD"`
	DB           string        `env:"POSTGRES_DB"`
	SSLMode      string        `env:"POSTGRES_SSL_MODE"`
	StartTimeout time.Duration `env:"POSTGRES_START_TIMEOUT"`
}

type MinIOConfig struct {
	Endpoint           string `env:"MINIO_ENDPOINT"`
	AccessKey          string `env:"MINIO_ROOT_USER"`
	SecretKey          string `env:"MINIO_ROOT_PASSWORD"`
	Secure             bool   `env:"MINIO_SECURE"`
	BucketPlaylist     string `env:"MINIO_BUCKET_PLAYLIST_COVERS"`
	BucketAlbum        string `env:"MINIO_BUCKET_ALBUM_COVERS"`
	BucketArtistAvatar string `env:"MINIO_BUCKET_ARTIST_AVATARS"`
	BucketAudio        string `env:"MINIO_BUCKET_AUDIO_FILES"`
}

type RedisConfig struct {
	Addr     string `env:"REDIS_ADDR"`
	Password string `env:"REDIS_PASSWORD"`
	DB       int    `env:"REDIS_DB"`
}

type PasswordHasherConfig struct {
	Cost int `env:"PASSWORD_HASHER_COST"`
}

type AccessTokenConfig struct {
	TTL       time.Duration `env:"ACCESS_TOKEN_TTL"`
	Jitter    time.Duration `env:"ACCESS_TOKEN_JITTER"`
	SecretKey []byte        `env:"SECRET_KEY"`
}

type RefreshTokenConfig struct {
	TTL    time.Duration `env:"REFRESH_TOKEN_TTL"`
	Jitter time.Duration `env:"REFRESH_TOKEN_JITTER"`
}

type RedisAccessMetaConfig struct {
	TTL    time.Duration `env:"ACCESS_META_TTL"`
	Jitter time.Duration `env:"ACCESS_META_JITTER"`
}

type LocalAccessMetaConfig struct {
	Size int `env:"LOCAL_CACHE_SIZE"`
}

type CookieConfig struct {
	Domain     string        `env:"COOKIE_DOMAIN"`
	Path       string        `env:"COOKIE_PATH"`
	Secure     bool          `env:"COOKIE_SECURE"`
	HttpOnly   bool          `env:"COOKIE_HTTP_ONLY"`
	RefreshTTL time.Duration `env:"COOKIE_REFRESH_TTL"`
	AccessTTL  time.Duration `env:"COOKIE_ACCESS_TTL"`
}

type AudioStorageConfig struct {
	Type     string `env:"AUDIO_STORAGE_TYPE"`
	BasePath string `env:"AUDIO_STORAGE_BASE_PATH"`
}

type LoggerConfig struct {
	Level string `env:"LOG_LEVEL"`
	Path  string `env:"LOG_PATH"`
}

type Config struct {
	HTTP                 HTTPConfig
	Postgres             PostgresConfig
	MinIO                MinIOConfig
	Redis                RedisConfig
	PasswordHasher       PasswordHasherConfig
	AccessToken          AccessTokenConfig
	RefreshToken         RefreshTokenConfig
	RedisAccessMetaCache RedisAccessMetaConfig
	LocalAccessMetaCache LocalAccessMetaConfig
	Cookie               CookieConfig
	AudioStorage         AudioStorageConfig
	Logger               LoggerConfig
}

var (
	cfg  *Config
	once sync.Once
)

func MustLoad(configPath string) *Config {
	once.Do(func() {
		log.Println("Loading config from environment variables...")
		cfg = &Config{}
		if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}
	})
	return cfg
}
