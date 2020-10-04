package publish_event

import (
	"context"
	"time"

	"github.com/asaskevich/govalidator"
)

var (
	DateFormat = "2006-01-02"
)

type Form struct {
	Title       string    `json:"title" valid:"stringlength(1|100)~Title must be present and less then 100 chars"`
	Description string    `json:"description" valid:"stringlength(1|1000)~Description must be present and less then 1000 chars"`
	Cost        float64   `json:"cost" valid:"float"`
	Date        string    `json:"date" valid:"matches(^\d\d\d\d-\d\d-\d\d$)~Date must be of format YYYY-MM-DD"`
	Latitude    float64   `json"latitude" valid:"float"`
	Longitude   float64   `json:"longitude" valid:"float"`
	EventTime   time.Time `json:"-" valid:"-"`
}

func (form *Form) Submit(ctx context.Context) (map[string]string, error) {
	_, err := govalidator.ValidateStruct(form)
	errorMessages := govalidator.ErrorsByField(err)

	if _, ok := errorMessages["date"]; !ok {
		eventTime, err := time.Parse(DateFormat, form.Date)
		if err != nil {
			errorMessages["date"] = "Wrong date format. Example format: 2017-07-21"
		}
		form.EventTime = eventTime
	}

	if _, ok := errorMessages["latitude"]; !ok {
		if form.Latitude < -90 || form.Latitude > 90 {
			errorMessages["latitude"] = "Latitude is invalid"
		}
	}

	if _, ok := errorMessages["longitude"]; !ok {
		if form.Longitude < -180 || form.Longitude > 180 {
			errorMessages["longitude"] = "Longitude is invalid"
		}
	}

	return errorMessages, nil
}

func NewForm() Form {
	return Form{}
}
