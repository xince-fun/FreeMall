package model

type SnowflakeEtcdHolder struct {
	EtcdAddressNode string
	ListenAddress   string
	Ip              string
	Port            string
	LastUpdateTime  int64
	WorkerId        int
}

type Endpoint struct {
	IP        string `json:"ip"`
	Port      string `json:"port"`
	Timestamp int64  `json:"timestamp"`
}
