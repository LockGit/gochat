/**
 * Created by lock
 * Date: 2019-08-09
 * Time: 15:24
 */
package config

import (
	"github.com/spf13/viper"
	"os"
	"runtime"
	"strings"
	"sync"
)

var once sync.Once
var realPath string
var Conf *Config

type Config struct {
	Connect ConnectConfig
	Job     JobConfig
	Logic   LogicConfig
}

func init() {
	Init()
}

func getFilePath() string {
	_, file, _, _ := runtime.Caller(0)
	return file
}

func Init() {
	once.Do(func() {
		env := GetMode()
		pathArr := strings.Split(getFilePath(), "/")
		realPath = strings.Join(pathArr[0:len(pathArr)-2], "/")
		configFilePath := realPath + "/config/" + env + "/"
		viper.SetConfigType("toml")
		viper.SetConfigName("/connect")
		viper.AddConfigPath(configFilePath)
		err := viper.ReadInConfig()
		if err != nil {
			panic(err)
		}
		viper.SetConfigName("/job")
		err = viper.MergeInConfig()
		if err != nil {
			panic(err)
		}
		viper.SetConfigName("/logic")
		err = viper.MergeInConfig()
		if err != nil {
			panic(err)
		}
		Conf = new(Config)
		viper.Unmarshal(&Conf.Connect)
		viper.Unmarshal(&Conf.Job)
		viper.Unmarshal(&Conf.Logic)
	})
}

func GetMode() string {
	env := os.Getenv("RUN_MODE")
	if env == "" {
		env = "dev"
	}
	return env
}

type ConnectBase struct {
	ServerId int    `mapstructure:"serverId"`
	CertPath string `mapstructure:"certPath"`
	KeyPath  string `mapstructure:"keyPath"`
}

type ConnectBucket struct {
	Num           int    `mapstructure:"num"`
	Channel       int    `mapstructure:"channel"`
	Room          int    `mapstructure:"room"`
	SrvProto      int    `mapstructure:"svrProto"`
	RoutineAmount uint64 `mapstructure:"routineAmount"`
	RoutineSize   int    `mapstructure:"routineSize"`
}

type ConnectWebsocket struct {
	Bind string `mapstructure:"bind"`
}

type ConnectConfig struct {
	ConnectBase      ConnectBase      `mapstructure:"connect-base"`
	ConnectBucket    ConnectBucket    `mapstructure:"connect-bucket"`
	ConnectWebsocket ConnectWebsocket `mapstructure:"connect-websocket"`
}

type JobBase struct {
	RedisAddr     string `mapstructure:"redisAddr"`
	RedisPassword string `mapstructure:"redisPassword"`
	RpcAddress    string `mapstructure:"rpcAddress"`
	PushChan      int    `mapstructure:"pushChan"`
	PushChanSize  int    `mapstructure:"pushChanSize"`
}

type JobConfig struct {
	JobBase JobBase `mapstructure:"job-base"`
}

type LogicConfig struct {
}
