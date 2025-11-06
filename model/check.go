package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"text/template"

	"github.com/coroot/coroot/timeseries"
	"github.com/coroot/coroot/utils"
	"github.com/dustin/go-humanize/english"
	"k8s.io/klog"
)

type CheckId string

type CheckType int

const (
	CheckTypeEventBased CheckType = iota
	CheckTypeItemBased
	CheckTypeValueBased
	CheckTypeManual
)

type CheckUnit string

const (
	CheckUnitPercent          = "percent"
	CheckUnitSecond           = "second"
	CheckUnitByte             = "byte"
	CheckUnitSecondsPerSecond = "seconds/second"
)

func (u CheckUnit) FormatValue(v float32) string {
	switch u {
	case CheckUnitSecond:
		return utils.FormatDuration(timeseries.Duration(v), 1)
	case CheckUnitSecondsPerSecond:
		return utils.FormatDuration(timeseries.Duration(v), 1) + "/second"
	case CheckUnitByte:
		value, unit := utils.FormatBytes(v)
		return value + unit
	case CheckUnitPercent:
		return utils.FormatPercentage(v)
	}
	return utils.FormatFloat(v)
}

type CheckConfig struct {
	Id    CheckId
	Type  CheckType
	Title string

	DefaultThreshold        float32
	Unit                    CheckUnit
	MessageTemplate         string
	ConditionFormatTemplate string
}

var Checks = struct {
	index map[CheckId]*CheckConfig

	SLOAvailability            CheckConfig
	SLOLatency                 CheckConfig
	CPUNode                    CheckConfig
	CPUContainer               CheckConfig
	MemoryOOM                  CheckConfig
	MemoryLeakPercent          CheckConfig
	MemoryPressure             CheckConfig
	StorageSpace               CheckConfig
	StorageIOLoad              CheckConfig
	NetworkRTT                 CheckConfig
	NetworkConnectivity        CheckConfig
	NetworkTCPConnections      CheckConfig
	InstanceAvailability       CheckConfig
	DeploymentStatus           CheckConfig
	InstanceRestarts           CheckConfig
	RedisAvailability          CheckConfig
	RedisLatency               CheckConfig
	MongodbAvailability        CheckConfig
	MongodbReplicationLag      CheckConfig
	MemcachedAvailability      CheckConfig
	PostgresAvailability       CheckConfig
	PostgresLatency            CheckConfig
	PostgresReplicationLag     CheckConfig
	PostgresConnections        CheckConfig
	LogErrors                  CheckConfig
	JvmAvailability            CheckConfig
	JvmSafepointTime           CheckConfig
	DotNetAvailability         CheckConfig
	PythonGILWaitingTime       CheckConfig
	NodejsEventLoopBlockedTime CheckConfig
	DnsLatency                 CheckConfig
	DnsServerErrors            CheckConfig
	DnsNxdomainErrors          CheckConfig
	MysqlAvailability          CheckConfig
	MysqlReplicationStatus     CheckConfig
	MysqlReplicationLag        CheckConfig
	MysqlConnections           CheckConfig
}{
	index: map[CheckId]*CheckConfig{},

	SLOAvailability: CheckConfig{
		Type:                    CheckTypeManual,
		Title:                   "可用性",
		MessageTemplate:         `应用正在提供错误`,
		DefaultThreshold:        99,
		Unit:                    CheckUnitPercent,
		ConditionFormatTemplate: "成功的请求百分比 < <threshold>",
	},
	SLOLatency: CheckConfig{
		Type:                    CheckTypeManual,
		Title:                   "延迟",
		MessageTemplate:         `应用正在缓慢执行`,
		DefaultThreshold:        99,
		Unit:                    CheckUnitPercent,
		ConditionFormatTemplate: "请求百分比 < <threshold> 的请求比 <bucket> 更快",
	},
	CPUNode: CheckConfig{
		Type:                    CheckTypeItemBased,
		Title:                   "节点 CPU 利用率",
		MessageTemplate:         `high CPU utilization of {{.Items "node"}}`,
		DefaultThreshold:        80,
		Unit:                    CheckUnitPercent,
		ConditionFormatTemplate: "节点的 CPU 使用率 > <threshold>",
	},
	CPUContainer: CheckConfig{
		Type:                    CheckTypeItemBased,
		Title:                   "容器 CPU 利用率",
		DefaultThreshold:        80,
		Unit:                    CheckUnitPercent,
		MessageTemplate:         `高 CPU 利用率 {{.Items "container"}}`,
		ConditionFormatTemplate: "容器的 CPU 使用率 > <threshold> 的 CPU 限制",
	},
	MemoryOOM: CheckConfig{
		Type:                    CheckTypeEventBased,
		Title:                   "内存不足",
		DefaultThreshold:        0,
		MessageTemplate:         `应用容器由于内存不足被重启 {{.Count "time"}} 次`,
		ConditionFormatTemplate: "the number of container terminations due to Out of Memory > <threshold>",
	},
	MemoryLeakPercent: CheckConfig{
		Type:                    CheckTypeValueBased,
		Title:                   "内存泄漏",
		DefaultThreshold:        10,
		MessageTemplate:         `内存使用率正在增长 {{.Value}} %% 每小时`,
		ConditionFormatTemplate: "内存使用率正在增长 > <threshold> % 每小时",
	},
	MemoryPressure: CheckConfig{
		Type:                    CheckTypeItemBased,
		Title:                   "内存压力",
		DefaultThreshold:        0.02,
		Unit:                    CheckUnitSecond,
		MessageTemplate:         `高内存阻塞时间 {{.Items "instances"}}`,
		ConditionFormatTemplate: "memory stall time > <threshold> per second",
	},
	StorageIOLoad: CheckConfig{
		Type:                    CheckTypeItemBased,
		Title:                   "磁盘 I/O 负载",
		DefaultThreshold:        5,
		Unit:                    CheckUnitSecondsPerSecond,
		MessageTemplate:         `高 I/O 负载 {{.Items "volume"}}`,
		ConditionFormatTemplate: "磁盘的 I/O 负载 > <threshold>",
	},
	StorageSpace: CheckConfig{
		Type:                    CheckTypeItemBased,
		Title:                   "磁盘空间",
		DefaultThreshold:        80,
		Unit:                    CheckUnitPercent,
		MessageTemplate:         `磁盘空间 {{.Items "volume"}} 即将耗尽`,
		ConditionFormatTemplate: "磁盘空间使用率 > <threshold>",
	},
	NetworkRTT: CheckConfig{
		Type:                    CheckTypeItemBased,
		Title:                   "网络往返时间 (RTT)",
		DefaultThreshold:        0.01,
		Unit:                    CheckUnitSecond,
		MessageTemplate:         `高网络延迟 {{.Items "upstream service"}}`,
		ConditionFormatTemplate: "网络往返时间 > <threshold>",
	},
	NetworkConnectivity: CheckConfig{
		Type:                    CheckTypeItemBased,
		Title:                   "网络连通性",
		DefaultThreshold:        0,
		MessageTemplate:         `没有连通性 {{.Items "upstream service"}}`,
		ConditionFormatTemplate: "the number of unavailable upstream services > <threshold>",
	},
	NetworkTCPConnections: CheckConfig{
		Type:                    CheckTypeItemBased,
		Title:                   "TCP 连接",
		DefaultThreshold:        0,
		MessageTemplate:         `无法连接到 {{.Items "upstream service"}}`,
		ConditionFormatTemplate: "the number of upstream services to which the app failed to connect > <threshold>",
	},
	InstanceAvailability: CheckConfig{
		Type:                    CheckTypeManual,
		Title:                   "实例可用性",
		DefaultThreshold:        75,
		Unit:                    CheckUnitPercent,
		MessageTemplate:         `{{.ItemsWithToBe "instance"}} 不可用`,
		ConditionFormatTemplate: "可用的实例 < <threshold> 的所需",
	},
	InstanceRestarts: CheckConfig{
		Type:                    CheckTypeEventBased,
		Title:                   "重启",
		DefaultThreshold:        0,
		MessageTemplate:         `应用容器已重启 {{.Count "time"}} 次`,
		ConditionFormatTemplate: "容器的重启次数 > <threshold>",
	},
	DeploymentStatus: CheckConfig{
		Type:                    CheckTypeValueBased,
		Title:                   "部署状态",
		DefaultThreshold:        180,
		Unit:                    CheckUnitSecond,
		MessageTemplate:         `部署已进行 {{.Value}} 秒`,
		ConditionFormatTemplate: "a rollout is in progress > <threshold>",
	},
	RedisAvailability: CheckConfig{
		Type:                    CheckTypeItemBased,
		Title:                   "Redis 可用性",
		DefaultThreshold:        0,
		MessageTemplate:         `{{.ItemsWithToBe "redis instance"}} 不可用`,
		ConditionFormatTemplate: "不可用的 Redis 实例 > <threshold>",
	},
	RedisLatency: CheckConfig{
		Type:                    CheckTypeItemBased,
		Title:                   "Redis latency",
		DefaultThreshold:        0.005,
		Unit:                    CheckUnitSecond,
		MessageTemplate:         `{{.ItemsWithToBe "redis instance"}} performing slowly`,
		ConditionFormatTemplate: "the average command execution time of a redis instance > <threshold>",
	},
	// MongodbAvailability: CheckConfig{
	// 	Type:                    CheckTypeItemBased,
	// 	Title:                   "Mongodb availability",
	// 	DefaultThreshold:        0,
	// 	MessageTemplate:         `{{.ItemsWithToBe "mongodb instance"}} unavailable`,
	// 	ConditionFormatTemplate: "the number of unavailable mongodb instances > <threshold>",
	// },
	// MongodbReplicationLag: CheckConfig{
	// 	Type:                    CheckTypeItemBased,
	// 	Title:                   "Mongodb replication lag",
	// 	DefaultThreshold:        30,
	// 	MessageTemplate:         `{{.ItemsWithToBe "mongodb replica"}} far behind the primary`,
	// 	ConditionFormatTemplate: "replication lag > <threshold>",
	// 	Unit:                    CheckUnitSecond,
	// },
	// MemcachedAvailability: CheckConfig{
	// 	Type:                    CheckTypeItemBased,
	// 	Title:                   "Memcached availability",
	// 	DefaultThreshold:        0,
	// 	MessageTemplate:         `{{.ItemsWithToBe "memcached instance"}} unavailable`,
	// 	ConditionFormatTemplate: "the number of unavailable memcached instances > <threshold>",
	// },
	// PostgresAvailability: CheckConfig{
	// 	Type:                    CheckTypeItemBased,
	// 	Title:                   "Postgres availability",
	// 	DefaultThreshold:        0,
	// 	MessageTemplate:         `{{.ItemsWithToBe "postgres instance"}} unavailable`,
	// 	ConditionFormatTemplate: "the number of unavailable postgres instances > <threshold>",
	// },
	// PostgresLatency: CheckConfig{
	// 	Type:                    CheckTypeItemBased,
	// 	Title:                   "Postgres latency",
	// 	DefaultThreshold:        0.1,
	// 	Unit:                    CheckUnitSecond,
	// 	MessageTemplate:         `{{.ItemsWithToBe "postgres instance"}} performing slowly`,
	// 	ConditionFormatTemplate: "the average query execution time of a postgres instance > <threshold>",
	// },
	// PostgresReplicationLag: CheckConfig{
	// 	Type:                    CheckTypeItemBased,
	// 	Title:                   "Postgres replication lag",
	// 	DefaultThreshold:        30,
	// 	MessageTemplate:         `{{.ItemsWithToBe "postgres replica"}} far behind the primary`,
	// 	ConditionFormatTemplate: "replication lag > <threshold>",
	// 	Unit:                    CheckUnitSecond,
	// },
	// PostgresConnections: CheckConfig{
	// 	Type:                    CheckTypeItemBased,
	// 	Title:                   "Postgres connections",
	// 	DefaultThreshold:        90,
	// 	MessageTemplate:         `{{.ItemsWithHave "postgres instance"}} too many connections`,
	// 	ConditionFormatTemplate: "the number of connections > <threshold> of `max_connections`",
	// 	Unit:                    CheckUnitPercent,
	// },
	LogErrors: CheckConfig{
		Type:                    CheckTypeEventBased,
		Title:                   "错误",
		DefaultThreshold:        0,
		MessageTemplate:         `{{.Count "error"}} 发生`,
		ConditionFormatTemplate: "错误和严重级别的消息数量 > <threshold>",
	},
	JvmAvailability: CheckConfig{
		Type:                    CheckTypeItemBased,
		Title:                   "JVM availability",
		DefaultThreshold:        0,
		MessageTemplate:         `{{.ItemsWithToBe "JVM instance"}} 不可用`,
		ConditionFormatTemplate: "不可用的 JVM 实例 > <threshold>",
	},
	JvmSafepointTime: CheckConfig{
		Type:                    CheckTypeItemBased,
		Title:                   "JVM 安全点",
		DefaultThreshold:        0.05,
		MessageTemplate:         `高安全点时间 {{.Items "JVM instance"}}`,
		ConditionFormatTemplate: "应用安全点操作停止的时间 > <threshold>",
		Unit:                    CheckUnitSecond,
	},
	// DotNetAvailability: CheckConfig{
	// 	Type:                    CheckTypeItemBased,
	// 	Title:                   ".NET runtime availability",
	// 	DefaultThreshold:        0,
	// 	MessageTemplate:         `{{.ItemsWithToBe ".NET instance"}} unavailable`,
	// 	ConditionFormatTemplate: "the number of unavailable .NET instances > <threshold>",
	// },
	// PythonGILWaitingTime: CheckConfig{
	// 	Type:                    CheckTypeItemBased,
	// 	Title:                   "Python GIL (Global Interpreter Lock) waiting time",
	// 	DefaultThreshold:        0.05,
	// 	MessageTemplate:         `high GIL waiting times on {{.Items "Python instance"}}`,
	// 	ConditionFormatTemplate: "the time Python threads have been waiting for acquiring the GIL (Global Interpreter Lock) > <threshold>",
	// 	Unit:                    CheckUnitSecond,
	// },
	// NodejsEventLoopBlockedTime: CheckConfig{
	// 	Type:                    CheckTypeItemBased,
	// 	Title:                   "Node.js event loop blocked time",
	// 	DefaultThreshold:        0.7,
	// 	MessageTemplate:         `high Node.js event loop blocked times on {{.Items "Node.js instance"}}`,
	// 	ConditionFormatTemplate: "the time Node.js event loop executes blocking code > <threshold>",
	// 	Unit:                    CheckUnitSecond,
	// },
	DnsLatency: CheckConfig{
		Type:                    CheckTypeValueBased,
		Title:                   "DNS 延迟",
		DefaultThreshold:        0.1,
		Unit:                    CheckUnitSecond,
		MessageTemplate:         `高延迟`,
		ConditionFormatTemplate: "DNS 响应时间的第 95 百分位数 > <threshold>",
	},
	DnsServerErrors: CheckConfig{
		Type:                    CheckTypeEventBased,
		Title:                   "DNS 服务器错误",
		DefaultThreshold:        0,
		MessageTemplate:         `{{.Count "server DNS error"}} 发生`,
		ConditionFormatTemplate: "服务器 DNS 错误数量 (不包括 NXDOMAIN) > <threshold>",
	},
	DnsNxdomainErrors: CheckConfig{
		Type:                    CheckTypeEventBased,
		Title:                   "DNS NXDOMAIN 错误",
		DefaultThreshold:        0,
		MessageTemplate:         `应用收到空 DNS 响应 {{.Count "time"}}`,
		ConditionFormatTemplate: "NXDOMAIN DNS 错误数量 (对于之前有效的请求) > <threshold>",
	},
	MysqlAvailability: CheckConfig{
		Type:                    CheckTypeItemBased,
		Title:                   "Mysql 可用性",
		DefaultThreshold:        0,
		MessageTemplate:         `{{.ItemsWithToBe "mysql instance"}} 不可用`,
		ConditionFormatTemplate: "不可用的 MySQL 实例 > <threshold>",
	},
	MysqlReplicationStatus: CheckConfig{
		Type:                    CheckTypeItemBased,
		Title:                   "Mysql 复制状态",
		DefaultThreshold:        0,
		MessageTemplate:         `{{.ItemsWithHave "mysql replica"}} IO 或 SQL 复制线程问题`,
		ConditionFormatTemplate: "IO 或 SQL 复制线程未运行 ",
		Unit:                    CheckUnitSecond,
	},
	MysqlReplicationLag: CheckConfig{
		Type:                    CheckTypeItemBased,
		Title:                   "Mysql 复制延迟",
		DefaultThreshold:        30,
		MessageTemplate:         `{{.ItemsWithToBe "mysql replica"}} 落后主节点`,
		ConditionFormatTemplate: "复制延迟 > <threshold>",
		Unit:                    CheckUnitSecond,
	},
	MysqlConnections: CheckConfig{
		Type:                    CheckTypeItemBased,
		Title:                   "Mysql 连接",
		DefaultThreshold:        90,
		MessageTemplate:         `{{.ItemsWithHave "mysql instance"}} 太多连接`,
		ConditionFormatTemplate: "连接数量 > <threshold> 的 `max_connections`",
		Unit:                    CheckUnitPercent,
	},
}

func init() {
	cs := reflect.ValueOf(&Checks).Elem()
	for i := 0; i < cs.NumField(); i++ {
		if !cs.Type().Field(i).IsExported() {
			continue
		}
		ch := cs.Field(i).Addr().Interface().(*CheckConfig)
		ch.Id = CheckId(cs.Type().Field(i).Name)
		Checks.index[ch.Id] = ch
	}
}

type CheckContext struct {
	items *utils.StringSet
	count int64
	value float32
	unit  CheckUnit
}

func (c CheckContext) Items(singular string) string {
	return english.Plural(c.items.Len(), singular, "")
}

func (c CheckContext) ItemsWithToBe(singular string) string {
	verb := "is"
	if c.items.Len() > 1 {
		verb = "are"
	}
	return c.Items(singular) + " " + verb
}

func (c CheckContext) ItemsWithHave(singular string) string {
	verb := "has"
	if c.items.Len() > 1 {
		verb = "have"
	}
	return c.Items(singular) + " " + verb
}

func (c CheckContext) Count(singular string) string {
	return english.Plural(int(c.count), singular, "")
}

func (c CheckContext) Value() string {
	return c.unit.FormatValue(c.value)
}

type Check struct {
	Id                      CheckId   `json:"id"`
	Title                   string    `json:"title"`
	Status                  Status    `json:"status"`
	Message                 string    `json:"message"`
	Threshold               float32   `json:"threshold"`
	Unit                    CheckUnit `json:"unit"`
	ConditionFormatTemplate string    `json:"condition_format_template"`

	typ             CheckType
	messageTemplate string
	items           *utils.StringSet
	count           int64
	desired         int64
	value           float32
	values          *timeseries.TimeSeries
	fired           bool
}

func (ch *Check) Value() float32 {
	return ch.value
}

func (ch *Check) SetValue(v float32) {
	ch.value = v
}

func (ch *Check) Values() *timeseries.TimeSeries {
	return ch.values
}
func (ch *Check) SetValues(vs *timeseries.TimeSeries) {
	ch.values = vs
}

func (ch *Check) Fire() {
	ch.fired = true
}

func (ch *Check) SetStatus(status Status, format string, a ...any) {
	ch.Status = status
	ch.Message = fmt.Sprintf(format, a...)
}

func (ch *Check) AddItem(format string, a ...any) {
	if len(a) == 0 {
		ch.items.Add(format)
		return
	}
	ch.items.Add(fmt.Sprintf(format, a...))
}

func (ch *Check) Count() int64 {
	return ch.count
}

func (ch *Check) Inc(amount int64) {
	ch.count += amount
}

func (ch *Check) ResetCounter() {
	ch.count = 0
}

func (ch *Check) Desired() int64 {
	return ch.desired
}

func (ch *Check) SetDesired(desired int64) {
	ch.desired = desired
}

func (ch *Check) Items() *utils.StringSet {
	return ch.items
}

func (ch *Check) Calc() {
	switch ch.typ {
	case CheckTypeEventBased:
		if ch.count <= int64(ch.Threshold) {
			return
		}
	case CheckTypeItemBased:
		if ch.items.Len() == 0 {
			return
		}
	case CheckTypeValueBased:
		if ch.value <= ch.Threshold {
			return
		}
	case CheckTypeManual:
		if !ch.fired {
			return
		}
	default:
		return
	}
	t, err := template.New("").Parse(ch.messageTemplate)
	if err != nil {
		ch.SetStatus(UNKNOWN, "invalid template: %s", err)
		return
	}
	buf := &bytes.Buffer{}
	if err := t.Execute(buf, CheckContext{items: ch.items, count: ch.count, value: ch.value, unit: ch.Unit}); err != nil {
		ch.SetStatus(UNKNOWN, "failed to render message: %s", err)
		return
	}
	ch.SetStatus(WARNING, buf.String())
}

type CheckConfigSource string

const (
	CheckConfigSourceKubernetesAnnotations CheckConfigSource = "kubernetes-annotations"
)

type CheckConfigSimple struct {
	Threshold float32 `json:"threshold"`
}

type CheckConfigSLOAvailability struct {
	Custom              bool    `json:"custom"`
	TotalRequestsQuery  string  `json:"total_requests_query"`
	FailedRequestsQuery string  `json:"failed_requests_query"`
	ObjectivePercentage float32 `json:"objective_percentage"`

	Source CheckConfigSource `json:"source,omitempty"`
	Error  string            `json:"error,omitempty"`
}

func (cfg *CheckConfigSLOAvailability) Total() string {
	return fmt.Sprintf(`sum(rate(%s[$RANGE]))`, cfg.TotalRequestsQuery)
}

func (cfg *CheckConfigSLOAvailability) Failed() string {
	return fmt.Sprintf(`sum(rate(%s[$RANGE]))`, cfg.FailedRequestsQuery)
}

type CheckConfigSLOLatency struct {
	Custom              bool    `json:"custom"`
	HistogramQuery      string  `json:"histogram_query"`
	ObjectiveBucket     float32 `json:"objective_bucket"`
	ObjectivePercentage float32 `json:"objective_percentage"`

	Source CheckConfigSource `json:"source,omitempty"`
	Error  string            `json:"error,omitempty"`
}

func (cfg *CheckConfigSLOLatency) Histogram() string {
	return fmt.Sprintf("sum by(le)(rate(%s[$RANGE]))", cfg.HistogramQuery)
}

type CheckConfigs map[ApplicationId]map[CheckId]json.RawMessage

func (cc CheckConfigs) getRaw(appId ApplicationId, checkId CheckId) (json.RawMessage, bool) {
	appIdStr := appId.String()

	if appConfigs, ok := cc[appId]; ok {
		if cfg, ok := appConfigs[checkId]; ok {
			return cfg, false
		}
	}

	for configAppId, appConfigs := range cc {
		if configAppId.IsZero() {
			continue
		}
		configAppIdStr := configAppId.String()
		if utils.GlobMatch(appIdStr, configAppIdStr) {
			if cfg, ok := appConfigs[checkId]; ok {
				return cfg, false
			}
		}
	}

	if appConfigs, ok := cc[ApplicationId{}]; ok {
		if cfg, ok := appConfigs[checkId]; ok {
			return cfg, true
		}
	}

	return nil, false
}

func (cc CheckConfigs) GetSimple(checkId CheckId, appId ApplicationId) CheckConfigSimple {
	cfg := CheckConfigSimple{Threshold: Checks.index[checkId].DefaultThreshold}
	raw, _ := cc.getRaw(appId, checkId)
	if raw == nil {
		return cfg
	}
	v, err := unmarshal[CheckConfigSimple](raw)
	if err != nil {
		klog.Warningln("failed to unmarshal check config:", err)
		return cfg
	}
	return v
}

func (cc CheckConfigs) GetSimpleAll(checkId CheckId, appId ApplicationId) []*CheckConfigSimple {
	def := Checks.index[checkId]
	if def == nil {
		klog.Warningln("unknown check:", checkId)
		return nil
	}
	res := []*CheckConfigSimple{{Threshold: Checks.index[checkId].DefaultThreshold}}
	ids := []ApplicationId{ApplicationIdZero}
	if !appId.IsZero() {
		ids = append(ids, appId)
	}
	for _, id := range ids {
		if appConfigs, ok := cc[id]; ok {
			if raw, ok := appConfigs[checkId]; ok {
				if cfg, err := unmarshal[CheckConfigSimple](raw); err != nil {
					klog.Warningln("failed to unmarshal check config:", err)
				} else {
					res = append(res, &cfg)
					continue
				}
			}
		}
		res = append(res, nil)
	}
	return res
}

func (cc CheckConfigs) GetByCheck(id CheckId) map[ApplicationId][]any {
	res := map[ApplicationId][]any{}
	for appId, appConfigs := range cc {
		for checkId, raw := range appConfigs {
			if checkId != id {
				continue
			}
			var cfg any
			var err error
			switch id {
			case Checks.SLOAvailability.Id:
				cfg, err = unmarshal[[]CheckConfigSLOAvailability](raw)
			case Checks.SLOLatency.Id:
				cfg, err = unmarshal[[]CheckConfigSLOLatency](raw)
			default:
				cfg, err = unmarshal[CheckConfigSimple](raw)
			}
			if err != nil {
				klog.Warningln("failed to unmarshal check config:", err)
				continue
			}
			res[appId] = append(res[appId], cfg)
		}
	}
	return res
}

func (cc CheckConfigs) GetAvailability(appId ApplicationId) (CheckConfigSLOAvailability, bool) {
	defaultCfg := CheckConfigSLOAvailability{
		Custom:              false,
		ObjectivePercentage: Checks.SLOAvailability.DefaultThreshold,
	}

	raw, _ := cc.getRaw(appId, Checks.SLOAvailability.Id)
	if raw == nil {
		if appId.Kind == ApplicationKindExternalService {
			defaultCfg.ObjectivePercentage = 0
		}
		return defaultCfg, true
	}
	res, err := unmarshal[[]CheckConfigSLOAvailability](raw)
	if err != nil {
		klog.Warningln("failed to unmarshal check config:", err)
		return defaultCfg, true
	}
	if len(res) == 0 {
		return defaultCfg, true
	}
	return res[0], false
}

func (cc CheckConfigs) GetLatency(appId ApplicationId, category ApplicationCategory) (CheckConfigSLOLatency, bool) {
	objectiveBucket := float32(0.5)
	auxObjectiveBucket := float32(5)
	if category.Auxiliary() {
		objectiveBucket = auxObjectiveBucket
	}
	defaultCfg := CheckConfigSLOLatency{
		Custom:              false,
		ObjectivePercentage: Checks.SLOLatency.DefaultThreshold,
		ObjectiveBucket:     objectiveBucket,
	}
	raw, projectDefault := cc.getRaw(appId, Checks.SLOLatency.Id)
	if raw == nil {
		if appId.Kind == ApplicationKindExternalService {
			defaultCfg.ObjectivePercentage = 0
		}
		return defaultCfg, true
	}
	res, err := unmarshal[[]CheckConfigSLOLatency](raw)
	if err != nil {
		klog.Warningln("failed to unmarshal check config:", err)
		return defaultCfg, true
	}
	if len(res) == 0 {
		return defaultCfg, true
	}
	if projectDefault && category.Auxiliary() {
		v := res[0]
		v.ObjectiveBucket = auxObjectiveBucket
		return v, false
	}
	return res[0], false
}

func unmarshal[T any](raw json.RawMessage) (T, error) {
	var cfg T
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
