package configure

var (
	// agent版本
	AgentVersion string

	// agent开始运行的时间
	StartTime int64

	// 本地环境的类型 windows | linux
	OsType int

	// 主机名
	HostName string

	// 本机IP
	LocalIp string

	// agent运行状态
	AgentStatus bool

	// 是否开启日志
	IsDebug bool

	// 日志的路径
	LogPath string

	// 日志大小
	LogSize int64 = 10000000

	// 主进程的id
	MainProcessPid string

	// 监控程序的id
	DaemonPid string

	// 配置文件路径
	ConfigPath string

	//当前项目 agent-go 的目录
	ProjectPath string

	// 脚本存放路径
	ScriptsPath string

	// 当前正在运行的脚本的数量
	ScriptsNum int

	// 监听更新命令的通道
	ScriptUpdateChannal chan bool

	// 获取到到脚本的信息Map
	ScriptsInfo []map[string]interface{}

	// 获取cpu信息map
	CpuMap []map[interface{}]interface{}

	//// 脚本返回结果队列
	//ResultQueue *Queue

	// 更新项目install脚本路径
	NewProjectShellPath string

	// 当前是否正在更新脚本
	UpdateScriptSignal bool
)