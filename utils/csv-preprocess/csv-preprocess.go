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
	"flag"
	"fmt"
	"os"
	"encoding/csv"
	"strconv"
	"strings"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("Usage: csv-preprocess <CSV Filename>")
		fmt.Println(args)
		os.Exit(1)
	}
	f, err := os.Open(args[0])
	if err != nil {
		panic(err)
	}
	c := csv.NewReader(f)
	c.FieldsPerRecord = -1
	recs, err := c.ReadAll()
	if err != nil {
		panic(err)
	}
	f.Close()

	for i, rec := range recs {
		if len(rec) != 2 {
			fmt.Println("Which comma is Separator?")
			fmt.Println(rec)
			fmt.Println(join(rec))
			index := 0
			fmt.Scanf("%d\n", &index)
			rec = append([]string{strings.Join(rec[:index], ",")}, rec[index:]...)
			recs[i] = rec
		} else {
			continue
		}
	}

	f, err = os.Create("out.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	w := csv.NewWriter(f)
	w.WriteAll(recs)
	w.Flush()
}

func join(rec []string) string {
	switch len(rec) {
	case 3:
		return rec[0] + " (1) " + rec[1] + " (2) " + rec[2]
	}
	tmp := ""
	for i, v := range rec {
		if i == 0 {
			tmp = v
			continue
		}
		tmp += " (" + strconv.Itoa(i) + ") " + v
	}
	return tmp
}
