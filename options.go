/*
 * Concert (C) 2016 Minio, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import "github.com/minio/cli"

var commands = []cli.Command{
	genCmd,
	renewCmd,
	serverCmd,
}

var genCmd = cli.Command{
	Name:  "gen",
	Usage: "Generate certificates (valid for 90 days).",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "dir",
			Value: "certs", // Default.
			Usage: "Generated certs destination directory. [DEFAULT: \"certs\"]",
		},
		cli.StringFlag{
			Name:  "sub-domains",
			Usage: "Generate certs for requested sub-domains.",
		},
		cli.StringFlag{
			Name:  "san-domains",
			Usage: "Generate certs for requested san-domains.",
		},
	},
	Action: genMain,
}

var renewCmd = cli.Command{
	Name:  "renew",
	Usage: "Renew certificates (valid for 90 days).",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "dir",
			Value: "certs", // Default.
		},
	},
	Action: renewMain,
}

var serverCmd = cli.Command{
	Name:  "server",
	Usage: "Run in server mode to automatically renew certificate, once in every 45 days.",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "dir",
			Value: "certs", // Default.
			Usage: "Generated certs destination directory. [DEFAULT: \"certs\"]",
		},
		cli.StringFlag{
			Name:  "sub-domains",
			Usage: "Generate certs for requested sub-domains.",
		},
		cli.StringFlag{
			Name:  "san-domains",
			Usage: "Generate certs for requested san-domains.",
		},
	},
	Action: serverMain,
}
