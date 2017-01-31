/*
Copyright Pascal Limeux. 2016 All Rights Reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
		 http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package setting

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"net"
	"strconv"
	"strings"
	"time"
)

type Settings struct {
	LogMode            string
	LogFileName        string
	DataSourceName     string
	ExpireInToken      time.Duration
	HttpHostUrl        string
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	HandlerTimeout     time.Duration
	HLTimeout          time.Duration
	DeployTimeout      time.Duration
	TransactionTimeout time.Duration
	HttpHyperledger    string
	ChainCodePath      string
	ChainCodeName      string
}

func (s *Settings) ToString() string {
	st := "Logger  --> file:" + s.LogFileName + " in " + s.LogMode + " mode \n"
	st = st + "Database--> name:" + s.DataSourceName + "\n"
	st = st + "Server  --> url :" + s.HttpHostUrl
	return st
}

func GetSettings(configPath, configFileName string) (Settings, error) {
	var configuration Settings
	viper.SetConfigName(configFileName)
	viper.AddConfigPath(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err.Error())
		return configuration, errors.New("Config file not found...")
	} else {
		configuration.LogMode = viper.GetString("logger.mode")
		configuration.LogFileName = viper.GetString("logger.logFileName")

		configuration.DataSourceName = viper.GetString("database.dataSourceName")

		configuration.HttpHostUrl, err = getHostUrl()
		if err != nil {
			return configuration, err
		}
		configuration.ReadTimeout = viper.GetDuration("server.readTimeout")
		configuration.WriteTimeout = viper.GetDuration("server.writeTimeout")
		configuration.HandlerTimeout = viper.GetDuration("server.handlerTimeout")
		configuration.HLTimeout = viper.GetDuration("server.hLTimeout")
		configuration.DeployTimeout = viper.GetDuration("server.deployTimeout")

		configuration.ExpireInToken = viper.GetDuration("token.expireInToken")

		configuration.HttpHyperledger = viper.GetString("hyperledger.httpHyperledger")
		configuration.ChainCodePath = viper.GetString("hyperledger.chainCodePath")
		configuration.ChainCodeName = viper.GetString("hyperledger.chainCodeName")

		fmt.Println(configuration.ToString())
		return configuration, nil
	}
}

func getHostUrl() (string, error) {
	ipAddress := viper.GetString("server.httpHostIp")
	ipPort := viper.GetInt("server.httpHostPort")
	var err error

	if ipAddress == "" {
		ipAddress, err = getOutboundIP()
		if err != nil {
			return ipAddress, errors.New(" Error to get local IP address")
		}
	}

	ipAddress = ipAddress + ":" + strconv.Itoa(ipPort)
	return ipAddress, nil
}

func getOutboundIP() (string, error) {
	var localAddr string

	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return localAddr, err
	}
	defer conn.Close()

	localAddr = conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")

	return localAddr[0:idx], nil
}
