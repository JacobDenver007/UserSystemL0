// Copyright (C) 2017, Beijing Bochen Technology Co.,Ltd.  All rights reserved.
//
// This file is part of L0
//
// The L0 is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The L0 is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package config

import (
	"github.com/bocheninc/L0/rpc"
	"github.com/spf13/viper"
)

func JrpcConfig(logFilePath, cfgFilePath string) *rpc.Config {
	config := rpc.NewDefaultConfig()
	config.Enabled = viper.GetBool("jrpc.enabled")
	config.Port = getString("jrpc.port", config.Port)
	config.Logcache = getInt("jrpc.logcache", config.Logcache)
	config.LogFilePath = logFilePath
	config.ConfigFilePath = cfgFilePath
	return config
}
