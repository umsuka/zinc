/* Copyright 2022 Zinc Labs Inc. and Contributors
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

package meta

import (
	"time"

	"gopkg.in/segmentio/analytics-go.v3"
)

var (
	// SEGMENT_CLIENT := analytics.New("hQYncuWEjDJC23MnU6jHXiye5k7qP2PL")
	SEGMENT_CLIENT analytics.Client
)

func init() {
	cf := analytics.Config{
		Interval:  15 * time.Second,
		BatchSize: 100,
		// Endpoint: "http://localhost:8080/api/v1/segment",
	}

	SEGMENT_CLIENT, _ = analytics.NewWithConfig("hQYncuWEjDJC23MnU6jHXiye5k7qP2PL", cf)
}
