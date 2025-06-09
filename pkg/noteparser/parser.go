package noteparser

import (
	"regexp"
	"strings"
	"time"

	"github.com/gitrus/digikeeper-bot/internal/note"
)

var (
	tagRegex = regexp.MustCompile(`#([A-Za-zА-Яа-я_]+)`)
)

const (
	DateTime = "2006-01-02 15:04:05"
	DateOnly = "2006-01-02"
	TimeOnly = "15:04:05"
	TimeHM   = "15:04"
)

const DateHM = "2006-01-02 15:04"

var dateLayouts = []string{time.RFC3339, DateTime, DateHM, DateOnly, TimeOnly, TimeHM}

type segmentParser func(n *note.Note, in string) (string, error)

func compose(parsers ...segmentParser) segmentParser {
	return func(n *note.Note, in string) (string, error) {
		var err error
		for _, p := range parsers {
			in, err = p(n, in)
			if err != nil {
				return "", err
			}
		}
		return in, nil
	}
}

func parseTags(n *note.Note, in string) (string, error) {
	matches := tagRegex.FindAllStringSubmatch(in, -1)
	for _, m := range matches {
		if len(m) > 1 {
			n.Tags = append(n.Tags, m[1])
		}
	}
	cleaned := tagRegex.ReplaceAllString(in, "")
	return cleaned, nil
}

func parseDate(n *note.Note, in string) (string, error) {
	fields := strings.Fields(in)
	remaining := make([]string, 0, len(fields))
	var (
		fullDT   *time.Time
		dateOnly *time.Time
		timeOnly *time.Time
	)
	for _, f := range fields {
		parsed := false
		for _, layout := range dateLayouts {
			if t, err := time.Parse(layout, f); err == nil {
				parsed = true
				switch layout {
				case DateTime, DateHM:
					temp := t
					fullDT = &temp
				case DateOnly:
					temp := t
					dateOnly = &temp
				case TimeOnly, TimeHM:
					temp := t
					timeOnly = &temp
				}
				break
			}
		}
		if !parsed {
			remaining = append(remaining, f)
		}
	}

	if fullDT != nil {
		n.Payload.EventAt = *fullDT
	} else {
		event := n.Payload.EventAt
		if dateOnly != nil {
			d := *dateOnly
			event = time.Date(d.Year(), d.Month(), d.Day(), event.Hour(), event.Minute(), event.Second(), 0, time.Local)
		}
		if timeOnly != nil {
			tm := *timeOnly
			event = time.Date(event.Year(), event.Month(), event.Day(), tm.Hour(), tm.Minute(), tm.Second(), 0, time.Local)
		}
		n.Payload.EventAt = event
	}
	return strings.Join(remaining, " "), nil
}

func Parse(createdAt time.Time, input string) (note.Note, error) {
	n := note.Note{
		CreatedAt: createdAt,
		Payload: note.Payload{
			EventAt: createdAt,
		},
	}

	p := compose(parseTags, parseDate)
	remaining, err := p(&n, input)
	if err != nil {
		return note.Note{}, err
	}
	n.Payload.Text = strings.TrimSpace(remaining)
	return n, nil
}
