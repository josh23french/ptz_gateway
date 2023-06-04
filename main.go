//  ptz_gateway - a router/translator for various Pan-Tilt-Zoom protocols
// 	Copyright (C) 2020  Joshua French
//
// 	This program is free software: you can redistribute it and/or modify
// 	it under the terms of the GNU Affero General Public License as published
// 	by the Free Software Foundation, either version 3 of the License, or
// 	(at your option) any later version.
//
// 	This program is distributed in the hope that it will be useful,
// 	but WITHOUT ANY WARRANTY; without even the implied warranty of
// 	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// 	GNU Affero General Public License for more details.
//
// 	You should have received a copy of the GNU Affero General Public License
// 	along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"flag"
	"os"
	"sync"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	var (
		config     = &Config{}
		configPath = flag.String("config.file", "", "path to config file")
		wg         = sync.WaitGroup{}
	)

	flag.Parse()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if *configPath == "" {
		path, err := findConfigFile("ptz_gateway")
		if err != nil {
			panic(err)
		}
		configPath = &path
	}

	loadConfig(*configPath, config)

	wg.Add(1)
	wg.Done()

	wg.Wait()
}
