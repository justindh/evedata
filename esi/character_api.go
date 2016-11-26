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
	"net/url"
	"encoding/json"
	"fmt"
	"strings"
)

type CharacterApi struct {
	Configuration Configuration
}

func NewCharacterApi() *CharacterApi {
	configuration := NewConfiguration()
	return &CharacterApi{
		Configuration: *configuration,
	}
}

func NewCharacterApiWithBasePath(basePath string) *CharacterApi {
	configuration := NewConfiguration()
	configuration.BasePath = basePath

	return &CharacterApi{
		Configuration: *configuration,
	}
}

/**
 * Get character&#39;s public information
 * Public information about a character  ---  Alternate route: &#x60;/v3/characters/{character_id}/&#x60;  Alternate route: &#x60;/legacy/characters/{character_id}/&#x60;  Alternate route: &#x60;/dev/characters/{character_id}/&#x60;   ---  This route is cached for up to 3600 seconds
 *
 * @param characterId An EVE character ID
 * @param datasource The server name you would like data from
 * @return *GetCharactersCharacterIdOk
 */
func (a CharacterApi) GetCharactersCharacterId(characterId int32, datasource string) (*GetCharactersCharacterIdOk, *APIResponse, error) {

	var httpMethod = "Get"
	// create path and map variables
	path := a.Configuration.BasePath + "/characters/{character_id}/"
	path = strings.Replace(path, "{"+"character_id"+"}", fmt.Sprintf("%v", characterId), -1)


	headerParams := make(map[string]string)
	queryParams := url.Values{}
	formParams := make(map[string]string)
	var postBody interface{}
	var fileName string
	var fileBytes []byte
	// add default headers if any
	for key := range a.Configuration.DefaultHeader {
		headerParams[key] = a.Configuration.DefaultHeader[key]
	}
		queryParams.Add("datasource", a.Configuration.APIClient.ParameterToString(datasource, ""))
	

	// to determine the Content-Type header
	localVarHttpContentTypes := []string{  }

	// set Content-Type header
	localVarHttpContentType := a.Configuration.APIClient.SelectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		headerParams["Content-Type"] = localVarHttpContentType
	}
	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{
		"application/json",
	}

	// set Accept header
	localVarHttpHeaderAccept := a.Configuration.APIClient.SelectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		headerParams["Accept"] = localVarHttpHeaderAccept
	}
	var successPayload = new(GetCharactersCharacterIdOk)
	httpResponse, err := a.Configuration.APIClient.CallAPI(path, httpMethod, postBody, headerParams, queryParams, formParams, fileName, fileBytes)
	if err != nil {
		return successPayload, NewAPIResponse(httpResponse.RawResponse), err
	}
	err = json.Unmarshal(httpResponse.Body(), &successPayload)
	return successPayload, NewAPIResponse(httpResponse.RawResponse), err
}

/**
 * Get corporation history
 * Get a list of all the corporations a character has been a member of  ---  Alternate route: &#x60;/v1/characters/{character_id}/corporationhistory/&#x60;  Alternate route: &#x60;/legacy/characters/{character_id}/corporationhistory/&#x60;  Alternate route: &#x60;/dev/characters/{character_id}/corporationhistory/&#x60;   ---  This route is cached for up to 3600 seconds
 *
 * @param characterId An EVE character ID
 * @param datasource The server name you would like data from
 * @return []GetCharactersCharacterIdCorporationhistory200Ok
 */
func (a CharacterApi) GetCharactersCharacterIdCorporationhistory(characterId int32, datasource string) ([]GetCharactersCharacterIdCorporationhistory200Ok, *APIResponse, error) {

	var httpMethod = "Get"
	// create path and map variables
	path := a.Configuration.BasePath + "/characters/{character_id}/corporationhistory/"
	path = strings.Replace(path, "{"+"character_id"+"}", fmt.Sprintf("%v", characterId), -1)


	headerParams := make(map[string]string)
	queryParams := url.Values{}
	formParams := make(map[string]string)
	var postBody interface{}
	var fileName string
	var fileBytes []byte
	// add default headers if any
	for key := range a.Configuration.DefaultHeader {
		headerParams[key] = a.Configuration.DefaultHeader[key]
	}
		queryParams.Add("datasource", a.Configuration.APIClient.ParameterToString(datasource, ""))
	

	// to determine the Content-Type header
	localVarHttpContentTypes := []string{  }

	// set Content-Type header
	localVarHttpContentType := a.Configuration.APIClient.SelectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		headerParams["Content-Type"] = localVarHttpContentType
	}
	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{
		"application/json",
	}

	// set Accept header
	localVarHttpHeaderAccept := a.Configuration.APIClient.SelectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		headerParams["Accept"] = localVarHttpHeaderAccept
	}
	var successPayload = new([]GetCharactersCharacterIdCorporationhistory200Ok)
	httpResponse, err := a.Configuration.APIClient.CallAPI(path, httpMethod, postBody, headerParams, queryParams, formParams, fileName, fileBytes)
	if err != nil {
		return *successPayload, NewAPIResponse(httpResponse.RawResponse), err
	}
	err = json.Unmarshal(httpResponse.Body(), &successPayload)
	return *successPayload, NewAPIResponse(httpResponse.RawResponse), err
}

/**
 * Get character portraits
 * Get portrait urls for a character  ---  Alternate route: &#x60;/v2/characters/{character_id}/portrait/&#x60;  Alternate route: &#x60;/dev/characters/{character_id}/portrait/&#x60;   ---  This route is cached for up to 3600 seconds
 *
 * @param characterId An EVE character ID
 * @param datasource The server name you would like data from
 * @return *GetCharactersCharacterIdPortraitOk
 */
func (a CharacterApi) GetCharactersCharacterIdPortrait(characterId int32, datasource string) (*GetCharactersCharacterIdPortraitOk, *APIResponse, error) {

	var httpMethod = "Get"
	// create path and map variables
	path := a.Configuration.BasePath + "/characters/{character_id}/portrait/"
	path = strings.Replace(path, "{"+"character_id"+"}", fmt.Sprintf("%v", characterId), -1)


	headerParams := make(map[string]string)
	queryParams := url.Values{}
	formParams := make(map[string]string)
	var postBody interface{}
	var fileName string
	var fileBytes []byte
	// add default headers if any
	for key := range a.Configuration.DefaultHeader {
		headerParams[key] = a.Configuration.DefaultHeader[key]
	}
		queryParams.Add("datasource", a.Configuration.APIClient.ParameterToString(datasource, ""))
	

	// to determine the Content-Type header
	localVarHttpContentTypes := []string{  }

	// set Content-Type header
	localVarHttpContentType := a.Configuration.APIClient.SelectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		headerParams["Content-Type"] = localVarHttpContentType
	}
	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{
		"application/json",
	}

	// set Accept header
	localVarHttpHeaderAccept := a.Configuration.APIClient.SelectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		headerParams["Accept"] = localVarHttpHeaderAccept
	}
	var successPayload = new(GetCharactersCharacterIdPortraitOk)
	httpResponse, err := a.Configuration.APIClient.CallAPI(path, httpMethod, postBody, headerParams, queryParams, formParams, fileName, fileBytes)
	if err != nil {
		return successPayload, NewAPIResponse(httpResponse.RawResponse), err
	}
	err = json.Unmarshal(httpResponse.Body(), &successPayload)
	return successPayload, NewAPIResponse(httpResponse.RawResponse), err
}

/**
 * Get character names
 * Resolve a set of character IDs to character names  ---  Alternate route: &#x60;/v1/characters/names/&#x60;  Alternate route: &#x60;/legacy/characters/names/&#x60;  Alternate route: &#x60;/dev/characters/names/&#x60;   ---  This route is cached for up to 3600 seconds
 *
 * @param characterIds A comma separated list of character IDs
 * @param datasource The server name you would like data from
 * @return []GetCharactersNames200Ok
 */
func (a CharacterApi) GetCharactersNames(characterIds []int64, datasource string) ([]GetCharactersNames200Ok, *APIResponse, error) {

	var httpMethod = "Get"
	// create path and map variables
	path := a.Configuration.BasePath + "/characters/names/"


	headerParams := make(map[string]string)
	queryParams := url.Values{}
	formParams := make(map[string]string)
	var postBody interface{}
	var fileName string
	var fileBytes []byte
	// add default headers if any
	for key := range a.Configuration.DefaultHeader {
		headerParams[key] = a.Configuration.DefaultHeader[key]
	}
	var collectionFormat = "csv"
	if collectionFormat == "multi" {
		for _, value := range characterIds {
			queryParams.Add("character_ids", value)
		}
	} else {
		queryParams.Add("character_ids", a.Configuration.APIClient.ParameterToString(characterIds, collectionFormat))
	}
			queryParams.Add("datasource", a.Configuration.APIClient.ParameterToString(datasource, ""))
	

	// to determine the Content-Type header
	localVarHttpContentTypes := []string{  }

	// set Content-Type header
	localVarHttpContentType := a.Configuration.APIClient.SelectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		headerParams["Content-Type"] = localVarHttpContentType
	}
	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{
		"application/json",
	}

	// set Accept header
	localVarHttpHeaderAccept := a.Configuration.APIClient.SelectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		headerParams["Accept"] = localVarHttpHeaderAccept
	}
	var successPayload = new([]GetCharactersNames200Ok)
	httpResponse, err := a.Configuration.APIClient.CallAPI(path, httpMethod, postBody, headerParams, queryParams, formParams, fileName, fileBytes)
	if err != nil {
		return *successPayload, NewAPIResponse(httpResponse.RawResponse), err
	}
	err = json.Unmarshal(httpResponse.Body(), &successPayload)
	return *successPayload, NewAPIResponse(httpResponse.RawResponse), err
}

/**
 * Calculate a CSPA charge cost
 * Takes a source character ID in the url and a set of target character ID&#39;s in the body, returns a CSPA charge cost  ---  Alternate route: &#x60;/v3/characters/{character_id}/cspa/&#x60;  Alternate route: &#x60;/legacy/characters/{character_id}/cspa/&#x60;  Alternate route: &#x60;/dev/characters/{character_id}/cspa/&#x60; 
 *
 * @param characterId An EVE character ID
 * @param characters The target characters to calculate the charge for
 * @param datasource The server name you would like data from
 * @return *PostCharactersCharacterIdCspaCreated
 */
func (a CharacterApi) PostCharactersCharacterIdCspa(characterId int32, characters PostCharactersCharacterIdCspaCharacters, datasource string) (*PostCharactersCharacterIdCspaCreated, *APIResponse, error) {

	var httpMethod = "Post"
	// create path and map variables
	path := a.Configuration.BasePath + "/characters/{character_id}/cspa/"
	path = strings.Replace(path, "{"+"character_id"+"}", fmt.Sprintf("%v", characterId), -1)


	headerParams := make(map[string]string)
	queryParams := url.Values{}
	formParams := make(map[string]string)
	var postBody interface{}
	var fileName string
	var fileBytes []byte
	// authentication '(evesso)' required
	// oauth required
	if a.Configuration.AccessToken != ""{
		headerParams["Authorization"] =  "Bearer " + a.Configuration.AccessToken
	}
	// add default headers if any
	for key := range a.Configuration.DefaultHeader {
		headerParams[key] = a.Configuration.DefaultHeader[key]
	}
		queryParams.Add("datasource", a.Configuration.APIClient.ParameterToString(datasource, ""))
	

	// to determine the Content-Type header
	localVarHttpContentTypes := []string{  }

	// set Content-Type header
	localVarHttpContentType := a.Configuration.APIClient.SelectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		headerParams["Content-Type"] = localVarHttpContentType
	}
	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{
		"application/json",
	}

	// set Accept header
	localVarHttpHeaderAccept := a.Configuration.APIClient.SelectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		headerParams["Accept"] = localVarHttpHeaderAccept
	}
	// body params
	postBody = &characters

	var successPayload = new(PostCharactersCharacterIdCspaCreated)
	httpResponse, err := a.Configuration.APIClient.CallAPI(path, httpMethod, postBody, headerParams, queryParams, formParams, fileName, fileBytes)
	if err != nil {
		return successPayload, NewAPIResponse(httpResponse.RawResponse), err
	}
	err = json.Unmarshal(httpResponse.Body(), &successPayload)
	return successPayload, NewAPIResponse(httpResponse.RawResponse), err
}

