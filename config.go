/*
 *     yinshiGo
 *     Copyright (C) 2017  bobo liu
 *
 *     This program is free software: you can redistribute it and/or modify
 *     it under the terms of the GNU Affero General Public License as published
 *     by the Free Software Foundation, either version 3 of the License, or
 *     (at your option) any later version.
 *
 *     This program is distributed in the hope that it will be useful,
 *     but WITHOUT ANY WARRANTY; without even the implied warranty of
 *     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *     GNU Affero General Public License for more details.
 *
 *     You should have received a copy of the GNU Affero General Public License
 *     along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

var (
	conf *config
)

func init() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println(`
 @@@@@@@@@@@@@@@@@@@@@@@                   @@@@@@@@@@@@@@@@@@@@@@@ 
 @                      @                 @                      @
@@                      @                 @                      @@
@@                      @@@@@@@@@@@@@@@@@@@                      @@
@@                      @                 @                      @@
 @                      @                 @@                     @ 
 #@@@@@@@@@@@@@@@@@@@@@@                   @@@@@@@@@@@@@@@@@@@@@@# 
                                                                   
                                          @                        
                                       @@                         
                         @@        @@@@                     
                           @@@@@@@@`)
	fmt.Println("\n\t稻花香里说丰年，听取人生经验。\n\n")

	confn := flag.String("conf", "./config.json", "Path to config file.")
	help := flag.Bool("h", false, "Show this message.")
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(1)
	}
	conf = &config{}
	conf.load(*confn)
}

type config struct {
	Reload  string    `json:"reload"`
	Addr    string    `json:"address"`
	Servers []*server `json:"servers"`
	FeedURL string    `json:"feed_url"`

	reloadDelay time.Duration
	doReload    bool
}

func (c *config) load(fn string) {
	log.Println("Loading config from", fn)

	data, err := ioutil.ReadFile(fn)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	err2 := json.Unmarshal(data, &c)
	if err2 != nil {
		log.Fatalln("Cannot parse config:", err2)
	}
	if len(c.Servers) == 0 {
		log.Fatalln("No server's config found.")
	}

	c.doReload = c.Reload != "" && c.Reload != "off"
	if c.doReload {
		c.reloadDelay, err = time.ParseDuration(c.Reload)
		if err != nil {
			log.Fatalln("Cannot parse reload delay:", err)
		}
		log.Println("Reload delay:", c.Reload)
	} else {
		log.Println("Reload disabled.")
	}

	log.Println(itoa(len(c.Servers)), "server(s) found.")
	for _, s := range c.Servers {
		if !strings.HasPrefix(s.Path, "/") {
			log.Fatalln("Wrong format of path:", s.Path)
		}
		s.load()
	}
}
