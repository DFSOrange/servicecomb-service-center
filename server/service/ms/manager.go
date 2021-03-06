// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ms

import (
	"fmt"
	"github.com/apache/servicecomb-service-center/pkg/log"
)

type dataSourceEngine func(opts Options) (DataSource, error)

var (
	plugins          = make(map[ImplName]dataSourceEngine)
	msDataSourceInst DataSource
)

// load plugins configuration into plugins
func Install(pluginImplName string, engineFunc dataSourceEngine) {
	plugins[ImplName(pluginImplName)] = engineFunc
}

// construct storage plugin instance
// invoked by sc main process
func Init(opts Options) error {
	if opts.PluginImplName == "" {
		return nil
	}

	msDataSourceEngine, ok := plugins[opts.PluginImplName]
	if !ok {
		return fmt.Errorf("microservice datasource implementation not supported [%s]", opts.PluginImplName)
	}
	var err error
	msDataSourceInst, err = msDataSourceEngine(opts)
	if err != nil {
		return err
	}
	log.Info(fmt.Sprintf("microservice datasource mgt [%s@%s] enabled", opts.PluginImplName, opts.Endpoint))
	return nil
}

// usage: ms.MicroService().create()
func MicroService() DataSource {
	return msDataSourceInst
}
