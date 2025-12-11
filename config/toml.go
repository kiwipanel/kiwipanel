package config

var (
	ProductionConfigPath  = "/opt/kiwipanel/config/kiwipanel.toml"
	DevelopmentConfigPath = "kiwipanel/config/kiwipanel.toml"
)

// Root config structure ---------------------------

type Config struct {
	Server   ServerConfig   `toml:"server"`
	Database DatabaseConfig `toml:"database"`
	Security SecurityConfig `toml:"security"`
	Log      LogConfig      `toml:"log"`
	Paths    PathConfig     `toml:"paths"`
}

// server ------------------------------------------

type ServerConfig struct {
	Bind        string `toml:"bind_addr"`    // 0.0.0.0
	Port        int    `toml:"bind_port"`    // 8443
	EnableHTTPS bool   `toml:"enable_https"` // true
	CertFile    string `toml:"cert_file"`
	KeyFile     string `toml:"key_file"`
}

// database ----------------------------------------

type DatabaseConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"` // encrypted
	Name     string `toml:"name"`
}

// security ----------------------------------------

type SecurityConfig struct {
	JWTSecret  string   `toml:"jwt_secret"` // encrypted
	AdminEmail string   `toml:"admin_email"`
	TrustedIPs []string `toml:"trusted_ips"`
}

// logger ------------------------------------------

type LogConfig struct {
	Level string `toml:"level"`
	File  string `toml:"file"`
}

// paths -------------------------------------------

type PathConfig struct {
	DataDir    string `toml:"data_dir"`
	BackupDir  string `toml:"backup_dir"`
	TempDir    string `toml:"temp_dir"`
	BinaryPath string `toml:"binary_path"`
}
