package def

const PasswdSalt = "excel2config"

type ErrorCode int

//go:generate stringer -type=ErrorCode -linecomment=true
const (
	ErrCodeStart                ErrorCode = -100
	ErrUserNotExist             ErrorCode = -101 //用户不存在
	ErrUserPasswd               ErrorCode = -102 //密码错误
	ErrEmailFormat              ErrorCode = -103 //邮箱格式错误
	ErrPwdNotConfirmed          ErrorCode = -104 //两次输入密码不一致
	ErrNeedLogin                ErrorCode = -105 //您尚未登录或者登录已过期
	ErrTableHead                ErrorCode = -106 //配置表头格式错误
	ErrPermissionDenied         ErrorCode = -107 //权限不足
	ErrSheetName                ErrorCode = -108 //Sheet名称错误
	ErrTableNotExist            ErrorCode = -109 //表格不存在
	ErrGroupNotExist            ErrorCode = -110 //项目不存在
	ErrGroupStoreEmpty          ErrorCode = -111 //请先配置数据仓库
	ErrGroupExportRedisFailed   ErrorCode = -112 //Redis导表失败，请检查配置格式
	ErrGroupExportMysqlFailed   ErrorCode = -113 //Mysql导表失败，请检查配置格式
	ErrGroupGetConfigFailed     ErrorCode = -114 //读取配置失败，请检查数据仓库配置
	ErrTableFormat              ErrorCode = -115 //配置表格式异常，请先检查
	ErrLoginParam               ErrorCode = -116 //登录参数错误
	ErrLoginFailed              ErrorCode = -117 //登录失败，请联系管理员
	ErrLoginDenied              ErrorCode = -118 //暂无登录权限，请联系管理员
	ErrDingtalkConfig           ErrorCode = -119 //钉钉登录配置错误，请联系管理员
	ErrLdapConfig               ErrorCode = -120 //LDAP配置错误，请联系管理员
	ErrGroupExportDatabusFailed ErrorCode = -121 //Databus更新配置失败，请检查连接配置
	ErrInvalidParam             ErrorCode = -122 //请求参数错误
	ErrCannotCreateExcel        ErrorCode = -123 //正式环境不允许创建Excel
	ErrCodeEnd                  ErrorCode = -1000
)

const (
	DaySeconds   = 86400
	MonthSeconds = DaySeconds * 30

	NeedCompressSheetRows = 200 //需要触发压缩的sheet行数
)

const (
	DefaultIntroductionSheet  = "配置说明"
	DefaultMongoDataBaseName  = "sheet_configs"
	DefaultRedisPubsubChannel = "config_refresh"
)

const (
	DsnTypeRedis   = 1
	DsnTypeMysql   = 2
	DsnTypeMongodb = 3
	DsnTypeRpc     = 4
)

const (
	RoleTypeAdmin     = 1 // 管理员
	RoleTypeDeveloper = 2 // 开发者
)
