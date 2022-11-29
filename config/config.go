package config

import (
	"os"
	"serverless-go-template/application/elastic"
	"serverless-go-template/application/logger"
	"serverless-go-template/application/loggernr"
	"strconv"

	"github.com/joho/godotenv"
	newrelic "github.com/newrelic/go-agent/v3/newrelic"
)

type AppConfig struct {
	SecretKey                 string
	AppEnv                    string
	Debug                     string
	LogLevel                  string
	RegionName                string
	SqsEndpoint               string
	SqsLocalstack             string
	LocalstackEndpoint        string
	AppQueue                  string
	ApiServer                 string
	ApiServerDescription      string
	ApiRoot                   string
	LocalApiServer            string
	LocalApiServerDescription string
	RedisHost                 string
	RedisPort                 string
	DbHost                    string
	DbUser                    string
	DbPassword                string
	Db                        string
	ElasticHost               string
	ElasticPort               string
	ElasticIndex              string
	ElasticKey                string
	ElasticSecret             string
	AppBucket                 string
	AppLambda                 string
	LoggerConfig              logger.Config
	NewRelicLicense           string
	NewRelicConfig            newrelic.Config
	NewRelicApp               *newrelic.Application
	NewRelicTransaction       newrelic.Transaction
	NewRelicExternalSegment   newrelic.ExternalSegment

	LoggerNR loggernr.LogMessage
}

var GlobalConfig AppConfig

func (cfg *AppConfig) LoadVariables(envPath string) error {

	err := godotenv.Load(envPath)
	if err != nil {
		cfg.LoggerNR.SendLog(loggernr.ERROR, ".env not found", err.Error())
		return err
	} else {
		cfg.LoggerNR.SendLog(loggernr.INFO, ".env found", struct{}{})
	}

	cfg.AppEnv = os.Getenv("ENVIRONMENT_NAME")

	cfg.LoggerNR = loggernr.LogMessage{
		GlobalEventName: "MAIN",
		//TODO adicionar o app name
		ServiceName: "serverless-go-template-" + cfg.AppEnv,
	}

	cfg.SecretKey = os.Getenv("SECRET_KEY")
	cfg.AppEnv = os.Getenv("APP_ENV")
	cfg.Debug = os.Getenv("DEBUG")
	cfg.LogLevel = os.Getenv("LOG_LEVEL")
	cfg.RegionName = os.Getenv("REGION_NAME")
	cfg.SqsEndpoint = os.Getenv("SQS_ENDPOINT")
	cfg.SqsLocalstack = os.Getenv("SQS_LOCALSTACK")
	cfg.LocalstackEndpoint = os.Getenv("LOCALSTACK_ENDPOINT")
	cfg.AppQueue = os.Getenv("APP_QUEUE")
	cfg.ApiServer = os.Getenv("API_SERVER")
	cfg.ApiServerDescription = os.Getenv("API_SERVER_DESCRIPTION")
	cfg.ApiRoot = os.Getenv("API_ROOT")
	cfg.LocalApiServer = os.Getenv("LOCAL_API_SERVER")
	cfg.LocalApiServerDescription = os.Getenv("LOCAL_API_SERVER_DESCRIPTION")
	cfg.RedisHost = os.Getenv("REDIS_HOST")
	cfg.RedisPort = os.Getenv("REDIS_PORT")
	cfg.DbHost = os.Getenv("DB_HOST")
	cfg.DbUser = os.Getenv("DB_USER")
	cfg.DbPassword = os.Getenv("DB_PASSWORD")
	cfg.Db = os.Getenv("DB")
	cfg.ElasticHost = os.Getenv("ELASTIC_HOST")
	cfg.ElasticPort = os.Getenv("ELASTIC_PORT")
	cfg.ElasticIndex = os.Getenv("ELASTIC_INDEX")
	cfg.LoggerConfig = logger.Config{
		AppEnv: cfg.AppEnv,
		ElasticVars: elastic.ElasticConfig{
			Hosts:        []string{cfg.ElasticHost + ":" + cfg.ElasticPort},
			AwsRegion:    cfg.RegionName,
			IndexDefault: cfg.ElasticIndex,
			AwsKey:       cfg.ElasticKey,
			AwsSecret:    cfg.ElasticSecret,
		},
	}
	cfg.NewRelicLicense = os.Getenv("NEW_RELIC_LICENSE")

	return nil
}

func getEnvInt(key string) int {
	val, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		GlobalConfig.LoggerNR.SendLog(loggernr.ERROR, "Invalid key: "+key+" It should be an integer",
			struct{}{})
		os.Exit(1)
	}

	return val
}
