package nextdate

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DateFormat string = "20060102"

func AfterNow(date, now time.Time) bool {
	return date.After(now)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", err
	}
	if !CheckRepeat(repeat) {
		return "", errors.New("ошибка: неверный формат правила повторения")
	}
	if repeat == "y" {
		for {
			date = date.AddDate(1, 0, 0)
			if AfterNow(date, now) {
				break
			}
		}
		dateString := date.Format(DateFormat)
		return dateString, nil
	} else if strings.Contains(repeat, "d") && strings.Contains(repeat, " ") {
		rule := strings.Split(repeat, " ")
		if len(rule) != 2 {
			return "", err
		}
		timeToAdd, err := strconv.Atoi(rule[1])
		if rule[0] != "d" || timeToAdd > 400 || timeToAdd < 0 {
			err = errors.New("ошибка: неверный формат правила повторения")
			return "", err
		}
		if err != nil {
			return "", err
		}
		for {
			date = date.AddDate(0, 0, timeToAdd)
			if AfterNow(date, now) {
				break
			}
		}
		dateString := date.Format(DateFormat)
		return dateString, nil
	}
	return "", err
}

func NexDateHandler(w http.ResponseWriter, r *http.Request) {
	// if r.Method != http.MethodPost {
	// 	http.Error(w, "ошибка: метод не поддерживается", http.StatusMethodNotAllowed)
	// 	return
	// }
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	nowString := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")
	now, err := time.Parse(DateFormat, nowString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if nowString == "" {
		now = time.Now()
	}
	nexDate, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(nexDate))
}

func CheckRepeat(repeat string) bool {
	if len(repeat) == 0 {
		return true
	}
	if repeat == "y" {
		return true
	}
	if strings.Contains(repeat, " ") {
		rule := strings.Split(repeat, " ")
		if len(rule) == 2 {
			ruleInt, err := strconv.Atoi(rule[1])
			if err != nil {
				return false
			}
			if rule[0] == "d" && ruleInt < 400 && ruleInt > 0 {
				return true
			}
		}
	}
	return false
}
