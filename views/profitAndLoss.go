package views

import (
	"encoding/json"
	"strconv"

	"html/template"
	"net/http"

	"github.com/antihax/evedata/evedata"
	"github.com/antihax/evedata/models"
	"github.com/antihax/evedata/templates"
)

func init() {
	evedata.AddRoute("profitandloss", "GET", "/profitAndLoss", profitAndLossPage)
	evedata.AddAuthRoute("profitandloss", "GET", "/U/walletSummary", walletSummaryAPI)
}

func profitAndLossPage(w http.ResponseWriter, r *http.Request) {
	setCache(w, 60*60)
	p := newPage(r, "Profit and Loss Statement")
	templates.Templates = template.Must(template.ParseFiles("templates/profitAndLoss.html", templates.LayoutPath))

	if err := templates.Templates.ExecuteTemplate(w, "base", p); err != nil {
		httpErr(w, err)
		return
	}
}

func walletSummaryAPI(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		rangeI int64
	)

	setCache(w, 5*60)
	s := evedata.SessionFromContext(r.Context())

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int32)
	if !ok {
		httpErrCode(w, nil, http.StatusUnauthorized)
		return
	}

	// Get range in days
	rangeTxt := r.FormValue("range")
	if rangeTxt != "" {
		rangeI, err = strconv.ParseInt(rangeTxt, 10, 64)
		if err != nil {
			httpErrCode(w, err, http.StatusBadRequest)
			return
		}
	} else {
		httpErrCode(w, err, http.StatusBadRequest)
		return
	}

	v, err := models.GetWalletSummary(characterID, rangeI)
	if err != nil {
		httpErr(w, err)
		return
	}

	json.NewEncoder(w).Encode(v)
}
