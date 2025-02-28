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

package v2

type Mappings struct {
	Properties map[string]Property `json:"properties,omitempty"`
}

type Property struct {
	Type           string `json:"type"` // text, keyword, date, numeric, boolean, geo_point
	Analyzer       string `json:"analyzer,omitempty"`
	SearchAnalyzer string `json:"search_analyzer,omitempty"`
	Format         string `json:"format,omitempty"` // date format yyyy-MM-dd HH:mm:ss || yyyy-MM-dd || epoch_millis
	Index          bool   `json:"index"`
	Store          bool   `json:"store"`
	Sortable       bool   `json:"sortable"`
	Aggregatable   bool   `json:"aggregatable"`
	Highlightable  bool   `json:"highlightable"`
}

func NewMappings() *Mappings {
	return &Mappings{
		Properties: make(map[string]Property),
	}
}

func NewProperty(typ string) Property {
	p := Property{
		Type:           typ,
		Analyzer:       "",
		SearchAnalyzer: "",
		Format:         "",
		Index:          true,
		Store:          false,
		Sortable:       false,
		Aggregatable:   false,
		Highlightable:  false,
	}

	switch typ {
	case "keyword":
		p.Aggregatable = true
	case "numeric", "date":
		p.Sortable = true
		p.Aggregatable = true
	}

	return p
}
