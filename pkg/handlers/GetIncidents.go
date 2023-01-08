package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"code.crogge.rs/chris/lcwc_api/pkg/models"
	"github.com/ungerik/go-rss"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func GetAllIncidents(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp, err := rss.Read("https://webcad.lcwc911.us/Pages/Public/LiveIncidentsFeed.aspx", false)
	if err != nil {
		fmt.Println(err)
	}

	channel, err := rss.Regular(resp)
	if err != nil {
		fmt.Println(err)
	}

	incidents := []models.Incident{}

	for _, item := range channel.Item {

		time, err := item.PubDate.ParseWithFormat(time.RFC1123)
		if err != nil {
			fmt.Println(err)
		}
		splitStr := strings.Split(item.Description, ";")
		township := cases.Title(language.AmericanEnglish).String(strings.TrimSpace(splitStr[0]))
		intersection := cases.Title(language.AmericanEnglish).String(strings.TrimSpace(splitStr[1]))
		units := cases.Title(language.AmericanEnglish).String(strings.TrimSpace(splitStr[2]))
		unitArr := strings.Split(units, "<Br>")
		for i := range unitArr {
			unitArr[i] = strings.TrimSpace(unitArr[i])
		}
		inc := models.Incident{Title: item.Title, Township: township, Intersection: intersection, Units: unitArr, PubDateUtc: time}
		inc.Type = getIncidentType(inc)
		incidents = append(incidents, inc)
	}

	json.NewEncoder(w).Encode(incidents)
}

func getIncidentType(i models.Incident) string {
	for _, unit := range i.Units {
		for _, hint := range models.FireUnitHints {
			if strings.Contains(strings.ToLower(unit), strings.ToLower(hint)) {
				return "fire"
			}
		}
		for _, hint := range models.MedicalUnitHints {
			if strings.Contains(strings.ToLower(unit), strings.ToLower(hint)) {
				return "medical"
			}
		}
		for _, hint := range models.TrafficUnitHints {
			if strings.Contains(strings.ToLower(unit), strings.ToLower(hint)) {
				return "traffic"
			}
		}
	}
	for _, hint := range models.FireTitleHints {
		if strings.Contains(strings.ToLower(i.Title), strings.ToLower(hint)) {
			return "fire"
		}
	}
	// don't know based off of data, default to traffic
	return "traffic"
}
