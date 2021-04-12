package model

type GroupInfo struct {
	Gid            string           `json:"gid,omitempty" bson:"_id,omitempty"`
	Name           string           `json:"name,omitempty" bson:"name"`
	Avatar         string           `json:"avatar,omitempty" bson:"avatar"`
	Remark         string           `json:"remark,omitempty" bson:"remark"`
	Store          []int32          `json:"store,omitempty" bson:"store"`
	AddTime        int64            `json:"addTime,omitempty" bson:"addTime"`
	Owner          string           `json:"owner,omitempty" bson:"owner"`
	Members        []SimpleUserInfo `json:"members" bson:"members"`
	RedisDSN       string           `json:"RedisDSN,omitempty" bson:"RedisDSN,omitempty"`
	RedisPassword  string           `json:"RedisPassword,omitempty" bson:"RedisPassword,omitempty"`
	RedisKeyPrefix string           `json:"RedisKeyPrefix,omitempty" bson:"RedisKeyPrefix,omitempty"`
	MysqlDSN       string           `json:"MysqlDSN,omitempty" bson:"MysqlDSN,omitempty"`
	MongodbDSN     string           `json:"MongodbDSN,omitempty" bson:"MongodbDSN,omitempty"`
	GRpcDsn        string           `json:"GRpcDsn,omitempty" bson:"GRpcDsn,omitempty"`
	GRpcAppKey     string           `json:"GRpcAppKey,omitempty" bson:"GRpcAppKey,omitempty"`
	GRpcAppSecret  string           `json:"GRpcAppSecret,omitempty" bson:"GRpcAppSecret,omitempty"`
	IsDev          bool             `json:"IsDev,omitempty" bson:"IsDev"`
	UnionGroupId   string           `json:"UnionGroupId,omitempty" bson:"UnionGroupId,omitempty"`
}
