package main

import (
	"fmt"
	"net/http" // network
	"time"
	"io/ioutil"
	"strings"
	//"reflect"
	"regexp"
	"strconv"
	"github.com/davecgh/go-spew/spew" // DUMPS - https://stackoverflow.com/questions/24512112/how-to-print-struct-variables-in-console
	"flag" // Cli params - https://stackoverflow.com/questions/2707434/how-to-access-command-line-arguments-passed-to-a-go-program
)

// program v GO - sosani info o hrajicich a nasledujicich kapelach

// jsem zmrd, budu se k tomu chovat jako k PHP3 a PHP4 - nema tridy a namespaces, takze podtrzitka resi naming trable
// mezistavy se pak hodi vedle do odpovidajicich structu - proste je ta trida rozlozena na prvocinitele
// "dynamicke" metody budou dostavat "self" jako v Pythonu a v nem prave structy one tridy
// jen nepujde retezit, protoze self instanci fakt nejde vratit

// na stringy
// @see https://pkg.go.dev/strings
// na struktury
// @see https://www.digitalocean.com/community/tutorials/defining-structs-in-go
// na instance
// @see https://stackoverflow.com/questions/7850140/how-do-you-create-a-new-instance-of-a-struct-from-its-type-at-run-time-in-go
// extra packages
// @see https://stackoverflow.com/questions/58114452/how-to-import-package-from-github

// festak ma nekolik dni a ma stage
// stage maji nekolik lineupu - pro kazdy aktivni den 1
// lineup obsahuje hrajici kapely a casy odkdy dokdy

type Band struct {
	Name string
	From time.Time
	To time.Time
}

type Stage struct {
	Name string
	From time.Time
	To time.Time
	Bands []Band
}

type Festival struct {
	Name string
	Stages []Stage
}

// mame co plnit, tak tedka jak

// obecny getter z HTTP; urcite nekde bude lepsi
//@staticMethod
func Request_getContent(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Read body: %v", err)
	}

	return data, nil
}

// class ValnikParser extends AParser
// AParser provides path storage and abstract method parseBody
type Parsers_AParser_data struct {
	Path string
}

func Parsers_ValnikParser_init(self Parsers_AParser_data) Parsers_AParser_data {
	self.Path = "https://valnik.cz/program/"
	return self
}

// @see https://www.digitalocean.com/community/tutorials/how-to-use-dates-and-times-in-go
// @staticMethod
func Parsers_ValnikParser_parseDay(today string) (time.Time, error) {
	res, err := time.Parse("02.1.2006", today)
	return res, err
}

// @staticMethod
func Parsers_ValnikParser_parseTime(band []string, hourPos int, minutePos int, start time.Time) (time.Time, error) {
	var timeCount = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.Local)
	hour, err1 := strconv.ParseInt(band[hourPos], 10, 32)
	minute, err2 := strconv.ParseInt(band[minutePos], 10, 32)
	if err1 != nil {
		return timeCount, err1
	}
	if err2 != nil {
		return timeCount, err2
	}
	// zaciname v 8, driv stejne neni uklizeno a pankaci nejsou venku; takze cokoliv mensiho je realne dalsi den
	if hour < 8 {
		hour += 24
	}
	// ted seskladat casy
	// wtf - obihani objektu prace s casem - chce Duration
	toAdd := time.Duration(hour) * time.Hour + time.Duration(minute) * time.Minute
	timeCount = timeCount.Add(toAdd)
	return timeCount, nil
}

// @staticMethod
func Parsers_ValnikParser_parseBody(dataLeft string) (Festival, error) {
	// parse it!
	// catch by...
	var titleStart = "<div class=\"et_pb_text_inner\"><h3>"
	var festival Festival
	festival.Name = "Valník"
	// now get days, each as extra stage

	// parse with regexp - jako svato, cos to tam narval?
	rTitles, errX := regexp.Compile(`inner"><h3>([^<]+)`)
	rBands, errY := regexp.Compile(`([\d]{2})[\D]+([\d]{2})[\s]+[\S]+[\s]+([\d]{2})[\D]([\d]{2})\s+([^<]+)`)
	
	if errX != nil {
		return festival, errX
	}
	if errY != nil {
		return festival, errY
	}
	
	splitDays := strings.Split(dataLeft, titleStart)
	days := splitDays[2:]
	
	for _, day := range days {
		// repair data from split
		myDay := titleStart + day
		// stage with title
		match := rTitles.FindStringSubmatch(myDay)
		//if match == nil {
		//	return nil, nil
		//}
		currentDay := strings.Split(match[1], " ")
		var stage Stage
		var err2 error
		var err3 error
		var err4 error
		stage.Name = currentDay[0]
		stage.From, err2 = Parsers_ValnikParser_parseDay(currentDay[1])
		
		// now bands
		match2 := rBands.FindAllStringSubmatch(myDay, -1)
		if err2 != nil {
			return festival, err2
		}
		
		//spew.Dump(stage) // DUMP
		//spew.Dump(match2) // DUMP
		for _, rgmt := range match2 {
			// tak tedka moc nevidim... potrebuju to podle matchu, pak se musi zpracovat skupiny v matchi
			var band Band
			band.Name = rgmt[5]
			band.From, err3 = Parsers_ValnikParser_parseTime(rgmt, 1, 2, stage.From)
			band.To, err4 = Parsers_ValnikParser_parseTime(rgmt, 3, 4, stage.From)
			if err3 != nil {
				return festival, err3
			}
			if err4 != nil {
				return festival, err4
			}
			//return festival, nil // DUMP - intentionally after first one
			stage.Bands = append(stage.Bands, band)
			stage.To = band.To
		}
		
		festival.Stages = append(festival.Stages, stage)
	}
	
	// and now return
	return festival, nil
}

func Parsers_AParser_getFestival(self Parsers_AParser_data) (Festival, error) {
	var festival Festival
	body, err := Request_getContent(self.Path)
	if err != nil {
		return festival, err
	}
	// parse it!
	festival, errx := Parsers_ValnikParser_parseBody(string(body[:]))
	if errx != nil {
		return festival, errx
	}
	return festival, nil
}

// takze je nacteno, ted zpracovani - chceme jen omezene mnozstvi dat, ne vsechno
// defaultni stav je aktualni datum a cas a data proti ni

func Process_Filter_filter(festival Festival, current time.Time, limited int) Festival {
	var filtered Festival
	filtered.Name = festival.Name
	for _, stgs := range festival.Stages {
		//spew.Dump(current) // DUMP
		if stgs.From.Before(current) && stgs.To.After(current) {
			var entries = 0
			var extStg Stage
			extStg.Name = stgs.Name
			extStg.From = stgs.From
			extStg.To = stgs.To
			for _, bnd := range stgs.Bands {
				if current.Before(bnd.To) && current.After(bnd.From) && entries < limited {
					entries += 1
					extStg.Bands = append(extStg.Bands, bnd)
				}
				if current.Before(bnd.From) && entries < limited {
					entries += 1
					extStg.Bands = append(extStg.Bands, bnd)
				}
			}
			filtered.Stages = append(filtered.Stages, extStg)
		}
	}
	return filtered
}

// tim je nacteno a zpracovano a je potreba zacit resit vypis
// do vypisu prijde aktualni a nasledujici kapela v tom, co zbylo po filtru

func Output_Cli_render(festival Festival, current time.Time) {
	// udelat nejaky rozumny dump...
	fmt.Println(festival.Name)
	fmt.Println(current)
	spew.Dump(festival.Stages)
}

// v parametrech je potreba predat festak a volitelne datum a cas
// vypis je pak den -> stage -> kapela stavajici a nasledujici; festak chodi jako pozadavek

func main() {

	var actual = time.Now()
	// param z CLI
	dayPtr := flag.Int("day", actual.Day(), "Day in month to get")
	hourPtr := flag.Int("hour", actual.Hour(), "Hour to get")
	minPtr := flag.Int("min", actual.Minute(), "Minute to get")
	nextPtr := flag.Int("bands", 2, "Next X bands on stage, default is 2")
	flag.Parse()

	var parsedCl Parsers_AParser_data
	parsedCl = Parsers_ValnikParser_init(parsedCl)
	festival, err := Parsers_AParser_getFestival(parsedCl) // natahat z netu
	if err != nil {
		fmt.Println(err)
		return
	}
	
	var current = time.Date(actual.Year(), actual.Month(), *dayPtr, *hourPtr, *minPtr, 0, 0, time.Local)
	filteredFestival := Process_Filter_filter(festival, current, *nextPtr) // filtr
	Output_Cli_render(filteredFestival, current) // vystup
}
