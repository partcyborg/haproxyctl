package haproxyctl

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ConfigOption func(*HAProxyConfig) error

// NewHAProxyConfig creates a new HAProxyConfig object
func NewHAProxyConfig(proxyUrl string, opts ...ConfigOption) (*HAProxyConfig, error) {
	endpoint, err := url.Parse(fmt.Sprintf("%s/", strings.TrimRight(proxyUrl, "/")))
	if err != nil {
		return nil, err
	}
	c := &HAProxyConfig{
		URL: *endpoint,
	}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	c.setupClient()
	return c, nil
}

func WithStatsPath(path string) ConfigOption {
	return func(c *HAProxyConfig) error {
		c.StatsPath = path
		return nil
	}
}

func WithAuthString(auth string) ConfigOption {
	return func(c *HAProxyConfig) error {
		return c.SetCredentialsFromAuthString(auth)
	}
}

func WithAuthInfo(user string, pass string) ConfigOption {
	return func(c *HAProxyConfig) error {
		c.Username = user
		c.Password = pass
		return nil
	}
}

func WithHttpClient(client *http.Client) ConfigOption {
	return func(c *HAProxyConfig) error {
		c.client = client
		return nil
	}
}

// HAProxyConfig holds the basic configuration options for haproxyctl
type HAProxyConfig struct {
	URL       url.URL
	StatsPath string
	Username  string
	Password  string
	client    *http.Client
	setupdone bool
}

func (c *HAProxyConfig) setupClient() {
	if c.setupdone {
		return
	}

	if c.client == nil {
		c.client = &http.Client{}
	}
	c.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	c.setupdone = true
}

// Statistics is a slice of HAProxy Statistics
type Statistics []Statistic

// Statistic contains a set of HAProxy Statistics
type Statistic struct {
	BackendName             string    `csv:"# pxname"`
	FrontendName            string    `csv:"svname"`
	QueueCurrent            uint64    `csv:"qcur"`
	QueueMax                uint64    `csv:"qmax"`
	SessionsCurrent         uint64    `csv:"scur"`
	SessionsMax             uint64    `csv:"smax"`
	SessionLimit            uint64    `csv:"slim"`
	SessionsTotal           uint64    `csv:"stot"`
	BytesIn                 uint64    `csv:"bin"`
	BytesOut                uint64    `csv:"bout"`
	DeniedRequests          uint64    `csv:"dreq"`
	DeniedResponses         uint64    `csv:"dresp"`
	ErrorsRequests          uint64    `csv:"ereq"`
	ErrorsConnections       uint64    `csv:"econ"`
	ErrorsResponses         uint64    `csv:"eresp"`
	WarningsRetries         uint64    `csv:"wretr"`
	WarningsDispatches      uint64    `csv:"wredis"`
	Status                  string    `csv:"status"`
	Weight                  uint64    `csv:"weight"`
	IsActive                uint64    `csv:"act"`
	IsBackup                uint64    `csv:"bck"`
	CheckFailed             uint64    `csv:"chkfail"`
	CheckDowned             uint64    `csv:"chkdown"`
	StatusLastChanged       Duration  `csv:"lastchg"`
	Downtime                Duration  `csv:"downtime"`
	QueueLimit              uint64    `csv:"qlimit"`
	ProcessID               uint64    `csv:"pid"`
	ProxyID                 uint64    `csv:"iid"`
	ServiceID               uint64    `csv:"sid"`
	Throttle                uint64    `csv:"throttle"`
	LBTotal                 uint64    `csv:"lbtot"`
	Tracked                 uint64    `csv:"tracked"`
	Type                    EntryType `csv:"type"`
	Rate                    uint64    `csv:"rate"`
	RateLimit               uint64    `csv:"rate_lim"`
	RateMax                 uint64    `csv:"rate_max"`
	CheckStatus             string    `csv:"check_status"`
	CheckCode               string    `csv:"check_code"`
	CheckDuration           uint64    `csv:"check_duration"`
	HTTPResponse1xx         uint64    `csv:"hrsp_1xx"`
	HTTPResponse2xx         uint64    `csv:"hrsp_2xx"`
	HTTPResponse3xx         uint64    `csv:"hrsp_3xx"`
	HTTPResponse4xx         uint64    `csv:"hrsp_4xx"`
	HTTPResponse5xx         uint64    `csv:"hrsp_5xx"`
	HTTPResponseOther       uint64    `csv:"hrsp_other"`
	CheckFailedDets         uint64    `csv:"hanafail"`
	RequestRate             uint64    `csv:"req_rate"`
	RequestRateMax          uint64    `csv:"req_rate_max"`
	RequestTotal            uint64    `csv:"req_tot"`
	AbortedByClient         uint64    `csv:"cli_abrt"`
	AbortedByServer         uint64    `csv:"srv_abrt"`
	CompressedBytesIn       uint64    `csv:"comp_in"`
	CompressedBytesOut      uint64    `csv:"comp_out"`
	CompressedBytesBypassed uint64    `csv:"comp_byp"`
	CompressedResponses     uint64    `csv:"comp_rsp"`
	LastSession             Duration  `csv:"lastsess"`
	LastCheck               string    `csv:"last_chk"`
	LastAgentCheck          string    `csv:"last_agt"`
	AvgQueueTime            uint64    `csv:"qtime"`
	AvgConnectTime          uint64    `csv:"ctime"`
	AvgResponseTime         uint64    `csv:"rtime"`
	AvgTotalTime            uint64    `csv:"ttime"`
	AgentStatus             uint64    `csv:"agent_status"`
	AgentCode               uint64    `csv:"agent_code"`
	AgentDuration           uint64    `csv:"agent_duration"`
	CheckDesc               string    `csv:"check_desc"`
	AgentDesc               string    `csv:"agent_desc"`
	CheckRise               uint64    `csv:"check_rise"`
	CheckFall               uint64    `csv:"check_fall"`
	CheckHealth             uint64    `csv:"check_health"`
	AgentRise               uint64    `csv:"agent_rise"`
	AgentFall               uint64    `csv:"agent_fall"`
	AgentHealth             uint64    `csv:"agent_health"`
	Address                 string    `csv:"addr"`
	Cookie                  uint64    `csv:"cookie"`
	Mode                    string    `csv:"mode"`
	LBAlgorithm             string    `csv:"algo"`
	ConnRate                uint64    `csv:"conn_rate"`
	ConnRateMax             uint64    `csv:"conn_rate_max"`
	ConnTotal               uint64    `csv:"conn_tot"`
	Intercepted             uint64    `csv:"intercepted"`
	DeniedCon               uint64    `csv:"dcon"`
	DeniedSes               uint64    `csv:"dses"`
	Wrew                    uint64    `csv:"wrew"`
	Connect                 uint64    `csv:"connect"`
	Reuse                   uint64    `csv:"reuse"`
	CacheLookups            uint64    `csv:"cache_lookups"`
	CacheHits               uint64    `csv:"cache_hits"`
	IdleConAvail            uint64    `csv:"srv_icur"`
	IdleConLimit            uint64    `csv:"src_ilim"`
	QtimeMax                uint64    `csv:"qtime_max"`
	CtimeMax                uint64    `csv:"ctime_max"`
	RtimeMax                uint64    `csv:"rtime_max"`
	TtimeMax                uint64    `csv:"ttime_max"`
	InternalErr             uint64    `csv:"eint"`
	IdleConnCur             uint64    `csv:"idle_conn_cur"`
	SafeConnCur             uint64    `csv:"safe_conn_cur"`
	UsedConnCur             uint64    `csv:"used_conn_cur"`
	NeedConnEst             uint64    `csv:"need_conn_est"`
}

// Duration is a type that we can attach CSV marshalling to for getting time.Duration
type Duration struct {
	time.Duration
}

// UnmarshalCSV converts the seconds timestamp into a golang time.Duration
func (date *Duration) UnmarshalCSV(csv string) (err error) {
	if csv == "" {
		return nil
	}
	timeString := fmt.Sprintf("%vs", csv)
	date.Duration, err = time.ParseDuration(timeString)
	if err != nil {
		return err
	}
	return nil
}

func (date *Duration) MarshalCSV() (string, error) {
	return fmt.Sprintf("%d", date.Nanoseconds()), nil
}

// You could also use the standard Stringer interface
func (date *Duration) String() string {
	return time.Duration(date.Nanoseconds()).String()
}

// EntryType can be a Frontend, Backend, Server or Socket
type EntryType int

const (
	// Frontend indicates this is a front-end
	Frontend EntryType = iota
	// Backend indicates this is a back-end
	Backend
	// Server indicates this is a server
	Server
	// Socket indicates this is a socket
	Socket
)

// Action is a set of actions that we can send to a HAProxy server
type Action string

const (
	ActionSetStateToReady     Action = "ready"
	ActionSetStateToDrain     Action = "drain"
	ActionSetStateToMaint     Action = "maint"
	ActionHealthDisableChecks Action = "dhlth"
	ActionHealthEnableChecks  Action = "ehlth"
	ActionHealthForceUp       Action = "hrunn"
	ActionHealthForceNoLB     Action = "hnolb"
	ActionHealthForceDown     Action = "hdown"
	ActionAgentDisablechecks  Action = "dagent"
	ActionAgentEnablechecks   Action = "eagent"
	ActionAgentForceUp        Action = "arunn"
	ActionAgentForceDown      Action = "adown"
	ActionKillSessions        Action = "shutdown"
)
