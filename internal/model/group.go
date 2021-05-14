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
	RedisDSN       string           `json:"RedisDSN,omitempty" bson:"RedisDSN"`
	RedisPassword  string           `json:"RedisPassword,omitempty" bson:"RedisPassword"`
	RedisKeyPrefix string           `json:"RedisKeyPrefix,omitempty" bson:"RedisKeyPrefix"`
	MysqlDSN       string           `json:"MysqlDSN,omitempty" bson:"MysqlDSN"`
	MongodbDSN     string           `json:"MongodbDSN,omitempty" bson:"MongodbDSN"`
	GRpcDsn        string           `json:"GRpcDsn,omitempty" bson:"GRpcDsn"`
	GRpcAppKey     string           `json:"GRpcAppKey,omitempty" bson:"GRpcAppKey"`
	GRpcAppSecret  string           `json:"GRpcAppSecret,omitempty" bson:"GRpcAppSecret"`
	IsDev          bool             `json:"IsDev,omitempty" bson:"IsDev"`
	UnionGroupId   string           `json:"UnionGroupId,omitempty" bson:"UnionGroupId"`
	AccessToken    string           `json:"AccessToken,omitempty" bson:"AccessToken"`
}
