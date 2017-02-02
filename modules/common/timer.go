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

package common

import (
	"github.com/pascallimeux/ocms2/modules/log"
	"time"
)

type Timer struct {
	Start   time.Time
	Elapsed time.Duration
}

func (t *Timer) StartTimer() {
	t.Start = time.Now()
}

func (t *Timer) GetTimer() time.Duration {
	t.Elapsed = time.Since(t.Start)
	return t.Elapsed
}

func (t *Timer) LogElapsed(location, label string) {
	t.Elapsed = time.Since(t.Start)
	log.Info(location, "processing ", label, " in ", t.Elapsed.String())
}
