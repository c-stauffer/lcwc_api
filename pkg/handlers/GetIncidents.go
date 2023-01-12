package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
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

		var township string
		var intersection string
		var units string
		var unitArr []string

		if len(splitStr) == 3 {
			// simple case where all values have been supplied in rss
			township = strTitle(splitStr[0])
			intersection = strTitle(splitStr[1])
			units = strTitle(splitStr[2])
			unitArr = strings.Split(units, "<Br>")
		} else {
			// missing at least one element, need to use hints to determine what we have rather than just assigning by position
			townshipEle := -1
			intersectionEle := -1
			unitsEle := -1
			unitRegexp, rErr := regexp.Compile(`.+[0-9]+-[0-9+].*`)
			for i := range splitStr {
				if (strContains(splitStr[1], "&") || strContains(splitStr[1], "/")) && intersectionEle == -1 {
					intersectionEle = i
					continue
				}
				if strContains(splitStr[1], "county") && townshipEle == -1 {
					townshipEle = i
					continue
				}

				if unitsEle == -1 {
					if strContains(splitStr[1], "<br>") || strContains(splitStr[i], "pending") {
						unitsEle = i
						continue
					}
					// looks like most units have a name and numeric hyphenated callsign i.e. "AMB 06-2" "MEDIC 02-42"
					match := unitRegexp.MatchString(splitStr[i])
					if rErr != nil {
						log.Println(rErr.Error())
					}
					if match {
						unitsEle = i
						continue
					}
				}
			}
			// we have assigned what we could, now fill in the rest with a guess, empty string the content if needed.
			// the content *should* be township > intersection > units
			if townshipEle > -1 {
				// assign via hint
				township = splitStr[townshipEle]
			} else if intersectionEle != 0 && unitsEle != 0 && len(splitStr) > 0 {
				// no other assignments at this pos so just take what we got
				township = splitStr[0]
			} else {
				// no hint to confirm a match, and something else took this pos, so must be a missing value.
				township = ""
			}

			if intersectionEle > -1 {
				intersection = splitStr[intersectionEle]
			} else if townshipEle != 1 && unitsEle != 1 && len(splitStr) > 1 {
				intersection = splitStr[1]
			} else {
				intersection = ""
			}

			if unitsEle > -1 {
				units = splitStr[unitsEle]
			} else if townshipEle != 2 && intersectionEle != 2 && len(splitStr) > 2 {
				// This doesn't make sense since hitting this means we had the full array all along, but just
				// including to match the other conditionals
				units = splitStr[2]
			} else {
				units = ""
			}

			township = strTitle(township)
			intersection = strTitle(intersection)
			units = strTitle(units)
			unitArr = strings.Split(units, "<Br>")
		}
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
			if strContains(unit, hint) {
				return "fire"
			}
		}
		for _, hint := range models.MedicalUnitHints {
			if strContains(unit, hint) {
				return "medical"
			}
		}
		for _, hint := range models.TrafficUnitHints {
			if strContains(unit, hint) {
				return "traffic"
			}
		}
	}
	for _, hint := range models.FireTitleHints {
		if strContains(i.Title, hint) {
			return "fire"
		}
	}
	// don't know based off of data, default to traffic
	return "traffic"
}

func strContains(searchStr string, subStr string) bool {
	return strings.Contains(strings.ToLower(searchStr), strings.ToLower(subStr))
}

func strTitle(str string) string {
	return cases.Title(language.AmericanEnglish).String(strings.TrimSpace(str))
}
