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

package api

import (
	"fmt"
)

type Consent struct {
	Consentid  string
	Action     string
	Appid      string
	Ownerid    string
	Consumerid string
	Datatype   string
	Dataaccess string
	Dt_begin   string
	Dt_end     string
}

type IsConsent struct {
	Consent string
}

func (c *Consent) Print() string {
	consentStr := fmt.Sprintf("ConsentID:%s ConsumerID:%s OwnerID:%s Datatype:%s Dataaccess:%s Dt_begin:%s Dt_end:%s", c.Consentid, c.Consumerid, c.Ownerid, c.Datatype, c.Dataaccess, c.Dt_begin, c.Dt_end)
	return consentStr
}

type DecodedConsent struct {
	Txuuid     string
	Appid      string
	Ownerid    string
	Consumerid string
	Datatype   string
	Dataaccess string
	Dt_begin   string
	Dt_end     string
	Ccid       string
}
