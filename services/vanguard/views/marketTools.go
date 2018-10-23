package views

import (
	"net/http"
	"strconv"
	"time"

	"github.com/antihax/evedata/services/vanguard"
	"github.com/antihax/evedata/services/vanguard/models"
)

func init() {
	vanguard.AddRoute("marketTools", "GET", "/marketUnderValue",
		func(w http.ResponseWriter, r *http.Request) {
			renderTemplate(w,
				"marketUnderValue.html",
				time.Hour*24*31,
				newPage(r, "EVE Online Undervalued Market Items"))
		})
	vanguard.AddRoute("marketTools", "GET", "/marketStationStocker",
		func(w http.ResponseWriter, r *http.Request) {
			renderTemplate(w,
				"marketStationStocker.html",
				time.Hour*24*31,
				newPage(r, "EVE Online Station Stocker"))
		})
	vanguard.AddRoute("marketRegions", "GET", "/J/marketRegions", marketRegionsAPI)
	vanguard.AddRoute("marketUnderValue", "GET", "/J/marketUnderValue", marketUnderValueAPI)
	vanguard.AddRoute("marketStationStocker", "GET", "/J/marketStationStocker", marketStationStockerAPI)
}

func marketRegionsAPI(w http.ResponseWriter, r *http.Request) {
	v, err := models.GetMarketRegions()
	if err != nil {
		httpErr(w, err)
		return
	}

	renderJSON(w, v, time.Hour*24*31)
}

func marketUnderValueAPI(w http.ResponseWriter, r *http.Request) {
	marketRegionID, err := strconv.ParseInt(r.FormValue("marketRegionID"), 10, 64)
	if err != nil {
		httpErr(w, err)
		return
	}
	sourceRegionID, err := strconv.ParseInt(r.FormValue("sourceRegionID"), 10, 64)
	if err != nil {
		httpErr(w, err)
		return
	}
	destinationRegionID, err := strconv.ParseInt(r.FormValue("destinationRegionID"), 10, 64)
	if err != nil {
		httpErr(w, err)
		return
	}
	discount, err := strconv.ParseFloat(r.FormValue("discount"), 64)
	if err != nil {
		httpErr(w, err)
		return
	}
	discount = discount / 100

	v, err := models.MarketUnderValued(marketRegionID, sourceRegionID, destinationRegionID, discount)
	if err != nil {
		httpErr(w, err)
		return
	}

	renderJSON(w, v, time.Hour)
}

func marketStationStockerAPI(w http.ResponseWriter, r *http.Request) {
	marketRegionID, err := strconv.ParseInt(r.FormValue("marketRegionID"), 10, 64)
	if err != nil {
		httpErr(w, err)
		return
	}

	destinationRegionID, err := strconv.ParseInt(r.FormValue("destinationRegionID"), 10, 64)
	if err != nil {
		httpErr(w, err)
		return
	}
	markup, err := strconv.ParseFloat(r.FormValue("markup"), 64)
	if err != nil {
		httpErr(w, err)
		return
	}
	markup = markup / 100

	v, err := models.MarketStationStocker(marketRegionID, destinationRegionID, markup)
	if err != nil {
		httpErr(w, err)
		return
	}

	renderJSON(w, v, time.Hour)
}
