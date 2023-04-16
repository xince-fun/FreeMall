package config

type MysqlConfig struct {
	SegmentEnable bool   `mapstructure:"segment_enable" json:"segment_enable" yaml:"segment_enable"`
	Host          string `mapstructure:"host" json:"host" yaml:"host"`
	Port          int    `mapstructure:"port" json:"port" yaml:"port"`
	User          string `mapstructure:"user" json:"user" yaml:"user"`
	Password      string `mapstructure:"password" json:"password" yaml:"password"`
	DbName        string `mapstructure:"db_name" json:"db_name" yaml:"db_name"`
	TableName     string `mapstructure:"table_name" json:"table_name" yaml:"table_name"`
	OpenConn      int    `mapstructure:"open_conn" json:"open_conn" yaml:"open_conn"`
	Idle          int    `mapstructure:"idle" json:"idle" yaml:"idle"`
	IdleTimeout   int    `mapstructure:"idle_timeout" json:"idle_timeout" yaml:"idle_timeout"`
}

type EtcdConfig struct {
	SnowflakeEnable bool     `mapstructure:"snowflake_enable" json:"snowflake_enable" yaml:"snowflake_enable"`
	DiscoveryEnable bool     `mapstructure:"discovery_enable" json:"discovery_enable" yaml:"discovery_enable"`
	Endpoints       []string `mapstructure:"endpoints" json:"endpoints" yaml:"endpoints"`
	DialTimeout     int      `mapstructure:"dial_timeout" json:"dial_timeout" yaml:"dial_timeout"`
	TimeDeviation   int      `mapstructure:"time_deviation" json:"time_deviation" yaml:"time_deviation"`
}

type ServerConfig struct {
	Name        string      `mapstructure:"name" json:"name" yaml:"name"`
	Host        string      `mapstructure:"host" json:"host" yaml:"host"`
	Port        int         `mapstructure:"port" json:"port" yaml:"port"`
	MysqlConfig MysqlConfig `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	EtcdConfig  EtcdConfig  `mapstructure:"etcd" json:"etcd" yaml:"etcd"`
	LogLevel    string      `mapstructure:"log_level" json:"log_level" yaml:"log_level"`
}
