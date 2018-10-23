package views

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"github.com/antihax/goesi"

	"github.com/antihax/evedata/evedata"
	"github.com/antihax/evedata/models"
	"github.com/antihax/evedata/templates"
)

func init() {
	evedata.AddRoute("account", "GET", "/account", accountPage)

	evedata.AddAuthRoute("account", "GET", "/X/accountInfo", accountInfo)
	evedata.AddAuthRoute("account", "POST", "/X/cursorChar", cursorChar)

	evedata.AddAuthRoute("crestTokens", "GET", "/U/crestTokens", apiGetCRESTTokens)
	evedata.AddAuthRoute("crestTokens", "DELETE", "/U/crestTokens", apiDeleteCRESTToken)

}

func accountPage(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)
	p := newPage(r, "Account Information")
	templates.Templates = template.Must(template.ParseFiles("templates/account.html", templates.LayoutPath))

	p["ScopeGroups"] = models.GetCharacterScopeGroups()

	if err := templates.Templates.ExecuteTemplate(w, "base", p); err != nil {
		httpErr(w, err)
		return
	}
}

func accountInfo(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)

	s := evedata.SessionFromContext(r.Context())
	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int64)
	if !ok {
		httpErrCode(w, http.StatusUnauthorized)
		return
	}

	char, ok := s.Values["character"].(goesi.VerifyResponse)
	if !ok {
		httpErrCode(w, http.StatusForbidden)
		return
	}

	accountInfo, ok := s.Values["accountInfo"].([]byte)
	if !ok {
		if err := updateAccountInfo(s, characterID, char.CharacterName); err != nil {
			httpErr(w, err)
			return
		}

		if err := s.Save(r, w); err != nil {
			httpErr(w, err)
			return
		}
	}

	w.Write(accountInfo)
}

func cursorChar(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)
	s := evedata.SessionFromContext(r.Context())

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int64)
	if !ok {
		httpErrCode(w, http.StatusUnauthorized)
		return
	}

	// Parse the cursorCharacterID
	cursorCharacterID, err := strconv.ParseInt(r.FormValue("cursorCharacterID"), 10, 64)
	if err != nil {
		httpErrCode(w, http.StatusForbidden)
		return
	}

	// Set our new cursor
	err = models.SetCursorCharacter(characterID, cursorCharacterID)
	if err != nil {
		httpErr(w, err)
		return
	}

	char, ok := s.Values["character"].(goesi.VerifyResponse)
	if !ok {
		httpErrCode(w, http.StatusForbidden)
		return
	}

	// Update the account information in redis
	if err = updateAccountInfo(s, characterID, char.CharacterName); err != nil {
		httpErr(w, err)
		return
	}

	if err = s.Save(r, w); err != nil {
		httpErr(w, err)
		return
	}
}

func apiGetCRESTTokens(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)
	s := evedata.SessionFromContext(r.Context())

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int64)
	if !ok {
		httpErrCode(w, http.StatusUnauthorized)
		return
	}

	v, err := models.GetCRESTTokens(characterID)
	if err != nil {
		httpErr(w, err)
		return
	}

	json.NewEncoder(w).Encode(v)

	if err = s.Save(r, w); err != nil {
		httpErr(w, err)
		return
	}
}

func apiDeleteCRESTToken(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)
	s := evedata.SessionFromContext(r.Context())

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int64)
	if !ok {
		httpErrCode(w, http.StatusUnauthorized)
		return
	}

	cid, err := strconv.ParseInt(r.FormValue("tokenCharacterID"), 10, 64)
	if err != nil {
		httpErrCode(w, http.StatusNotFound)
		return
	}

	if err := models.DeleteCRESTToken(characterID, cid); err != nil {
		httpErrCode(w, http.StatusConflict)
		return
	}

	char, ok := s.Values["character"].(goesi.VerifyResponse)
	if !ok {
		httpErrCode(w, http.StatusForbidden)
		return
	}

	if err = updateAccountInfo(s, characterID, char.CharacterName); err != nil {
		httpErr(w, err)
		return
	}

	if err = s.Save(r, w); err != nil {
		httpErr(w, err)
		return
	}
}
