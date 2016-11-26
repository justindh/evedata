/* 
 * EVE Swagger Interface
 *
 * An OpenAPI for EVE Online
 *
 * OpenAPI spec version: 0.2.6.dev1
 * 
 * Generated by: https://github.com/swagger-api/swagger-codegen.git
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package swagger

import (
	"time"
)

// 200 ok object
type GetCharactersCharacterIdOk struct {

	// ancestry_id integer
	AncestryId int32 `json:"ancestry_id,omitempty"`

	// Creation date of the character
	Birthday time.Time `json:"birthday,omitempty"`

	// bloodline_id integer
	BloodlineId int32 `json:"bloodline_id,omitempty"`

	// The character's corporation ID
	CorporationId int32 `json:"corporation_id,omitempty"`

	// description string
	Description string `json:"description,omitempty"`

	// gender string
	Gender string `json:"gender,omitempty"`

	// The name of the character
	Name string `json:"name,omitempty"`

	// race_id integer
	RaceId int32 `json:"race_id,omitempty"`

	// security_status number
	SecurityStatus float32 `json:"security_status,omitempty"`
}
