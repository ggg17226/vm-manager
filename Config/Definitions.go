package Config

type TomlConfig struct {
	Db     DbConfig     `toml:"db"`
	Server ServerConfig `toml:"server"`
	Host   HostConfig   `toml:"host"`
}

type HostConfig struct {
	NetType string `toml:"netType"`
}

type DbConfig struct {
	Host     string `toml:"host"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	DbName   string `toml:"dbName"`
}

type ServerConfig struct {
	Listen string `toml:"listen"`
}
