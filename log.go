// Copyright 2019 The SQLFlow Authors. All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package gomaxcompute

import (
	"os"
	"path"

	"github.com/sirupsen/logrus"
)

var log *logrus.Entry

func init() {
	logDir := "./logs"
	if e := os.MkdirAll(logDir, 0744); e != nil {
		log.Panicf("create log directory failed: %v", e)
	}

	f, e := os.Create(path.Join(logDir, "maxcompute.log"))
	if e != nil {
		log.Panicf("open log file failed: %v", e)
	}

	lg := logrus.New()
	lg.SetOutput(f)
	lg.SetLevel(logrus.InfoLevel)
	log = lg.WithFields(logrus.Fields{"driver": "odps"})
}
