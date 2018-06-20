package api

type LBS struct {
	Policy string `json:"policy"`
}

type Server struct {
	Name    string `json:"servername"`
	Version string `json:"version"`
	Lb      LBS    `json:"lb"`
}

type HttpProxyCfg struct {
	Addr string   `json:"listen"`
	Svc  []Server `json:"server"`
}

type TcpProxyCfg struct {
	Addr string   `json:"listen"`
	Svc  []Server `json:"server"`
}

type ProxyCfg struct {
	Http []HttpProxyCfg `json:"https"`
	Tcp  []TcpProxyCfg  `json:"tcps"`
}
