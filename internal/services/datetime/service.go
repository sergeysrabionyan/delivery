package datetime

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const dateLayout = "02.01.2006"

var days = map[string]string{
	"January":   "января",
	"February":  "февраля",
	"March":     "марта",
	"April":     "апреля",
	"May":       "мая",
	"June":      "июня",
	"July":      "июля",
	"August":    "августа",
	"September": "сентября",
	"October":   "октября",
	"November":  "ноября",
	"December":  "декабря",
}

var errorDateParse = errors.New("ошибка преобразования даты")
var datesNotFound = errors.New("нет доступных дат")

type DateInfo struct {
	RawDates string
	dates    dates
	Year     string
}

func New() *DateInfo {
	return &DateInfo{
		dates: make(dates, 0),
	}
}

func (d *DateInfo) Init() error {
	dateList, err := d.parseDates()
	if err != nil {
		return err
	}
	d.dates = append(d.dates, dateList...)
	return nil
}

func (d *DateInfo) Validate() bool {
	return d.RawDates != "" && d.Year != ""
}

type dates []string

func (d *DateInfo) GetDateRange() (int, int, error) {
	if len(d.dates) == 0 {
		return 0, 0, errors.New("нет доступных дат")
	}
	if len(d.dates) == 1 {
		date, err := time.Parse(dateLayout, d.dates[0])
		if err != nil {
			return 0, 0, errorDateParse
		}
		return date.Day(), date.Day(), nil
	}
	start, err := time.Parse(dateLayout, d.dates[0])
	if err != nil {
		return 0, 0, errorDateParse
	}
	end, err := time.Parse(dateLayout, d.dates[len(d.dates)-1])
	if err != nil {
		return 0, 0, errorDateParse
	}
	return start.Day(), end.Day(), nil
}

func (d *DateInfo) CountDates() int {
	return len(d.dates)
}

func (d *DateInfo) GetDates() []string {
	return d.dates
}

func (d *DateInfo) GetMonth() (string, error) {
	if len(d.dates) == 0 {
		return "", datesNotFound
	}
	date, err := time.Parse(dateLayout, d.dates[0])
	if err != nil {
		return "", errorDateParse
	}
	return days[date.Month().String()], nil
}

func (d *DateInfo) GetYear() string {
	return d.Year + " года"
}

func (d *DateInfo) GetDateRangeString() (string, error) {
	start, end, err := d.GetDateRange()
	if err != nil {
		return "", err
	}
	month, err := d.GetMonth()
	if err != nil {
		return "", err
	}
	return "c " + strconv.Itoa(start) + " по " + strconv.Itoa(end) + " " + month, nil
}

func (d *DateInfo) parseDates() ([]string, error) {
	if d.RawDates == "" {
		return nil, errors.New("не заполнены даты")
	}
	reg := regexp.MustCompile(`\[(.*?)\]`)

	dates := make([]string, 0)
	for _, v := range reg.FindAllStringSubmatch(d.RawDates, -1) {
		dates = append(dates, addDatesFromParts(v[1], d.Year)...)
	}
	return dates, nil
}

func addDatesFromParts(datePeriod string, year string) []string {
	if !strings.Contains(datePeriod, "-") {
		return []string{datePeriod + "." + year}
	}
	dateParts := strings.Split(datePeriod, "-")
	return getDateRange(dateParts[0], dateParts[1], year)
}

func getDateRange(start string, end string, year string) []string {
	start = start + "." + year
	end = end + "." + year
	dates := make([]string, 0)
	startDate, err := time.Parse(dateLayout, start)
	if err != nil {
		return nil
	}
	endDate, err := time.Parse(dateLayout, end)
	if err != nil {
		return nil
	}
	countLimit := 0
	for startDate != endDate.AddDate(0, 0, 1) {
		dates = append(dates, startDate.Format(dateLayout))
		startDate = startDate.AddDate(0, 0, 1)
		countLimit++
		if countLimit > 50 {
			return nil
		}
	}
	return dates
}
