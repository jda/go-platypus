package platypus

import (
	"errors"
	"time"
)

type LastRunMethodParameters struct {
	Query    string `xml:"query"`
	Datatype string `xml:"datatype"`
}

func (p Platypus) LastRun() (time.Time, error) {
	lastRun := time.Unix(0, 0)

	params := LastRunMethodParameters{
		Datatype: "XML",
		Query:    "SELECT TOP 1 eventlog_etid, eventlog_result, eventlog_msg, eventlog_date FROM eventlog ORDER BY eventlog_date DESC",
	}

	res, err := p.Exec("SQL", params)
	if err != nil {
		return lastRun, err
	}

	if res.Success == 0 {
		return lastRun, errors.New(res.ResponseText)
	}

	// make sure we got data
	if len(res.Attributes.Block) < 1 {
		return lastRun, errors.New(ERR_INSUFFICIENT_RESPONSE)
	}

	ab := unwrapAttributeBlock(res.Attributes.Block[0])

	rawDate := ab["eventlog_date"]
	lastRun, err = time.ParseInLocation("2006-01-02T15:04:05", rawDate, getLocalTimeLocation())
	if err != nil {
		return lastRun, err
	}

	return lastRun, nil
}
