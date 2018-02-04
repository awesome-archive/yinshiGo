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
	"encoding/csv"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
	"sync"
)

var separator = []byte(" --")

type server struct {
	Path string `json:"path"`
	File string `json:"file"`

	records    [][][]byte
	recHaveSrc []bool
	recCount   int

	modTime time.Time
	reload  sync.RWMutex
}

func (s *server) load() {
	s.reload.Lock()
	defer s.reload.Unlock()

	log.Println("Loading records of", s.Path, "from", s.File)

	var f *os.File
	for {
		tf, err := os.Open(s.File)
		if err == nil {
			f = tf
			break
		}
	}
	defer f.Close()

	c := csv.NewReader(f)
	c.FieldsPerRecord = 2
	c.TrimLeadingSpace = true
	rec, err := c.ReadAll()
	if err != nil {
		log.Fatalln("Cannot parse data file:", err)
	}

	s.records = make([][][]byte, len(rec))
	s.recHaveSrc = make([]bool, len(rec))
	for i, r := range rec {
		s.records[i] = [][]byte{[]byte(r[0]), []byte(r[1])}
		s.recHaveSrc[i] = len(r[1]) != 0
	}
	s.recCount = len(rec)

	log.Println(itoa(s.recCount), "records loaded.")

	if stat, err := f.Stat(); conf.doReload && err == nil {
		s.modTime = stat.ModTime()
		go s.reloadWatch()
	}
}

func (s *server) reloadWatch() {
	log.Println("Watching server", s.Path)
	for {
		time.Sleep(conf.reloadDelay)
		stat, err := os.Stat(s.File)
		if err != nil {
			time.Sleep(conf.reloadDelay * 3)
			continue
		}
		if s.modTime != stat.ModTime() {
			log.Println("Reloading", s.Path)
			s.load()
			return
		}
	}
}

func (s server) Get(w http.ResponseWriter, r *http.Request) {
	s.reload.RLock()
	defer s.reload.RUnlock()
	w.Write(s.records[rand.Intn(s.recCount)][0])
}

func (s server) GetLong(w http.ResponseWriter, r *http.Request) {
	n := rand.Intn(s.recCount)
	s.reload.RLock()
	defer s.reload.RUnlock()
	w.Write(s.records[n][0])
	if s.recHaveSrc[n] {
		w.Write(separator)
		w.Write(s.records[n][1])
	}
}
