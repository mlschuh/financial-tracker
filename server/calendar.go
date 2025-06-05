package main

// FullCalendarEvent (Structure FullCalendar expects for event objects)
type FullCalendarEvent struct {
	ID            string                 `json:"id"`
	Title         string                 `json:"title"`
	Start         string                 `json:"start"` // FC expects ISO 8601 string or Date object
	End           string                 `json:"end,omitempty"`
	AllDay        bool                   `json:"allDay"`
	URL           string                 `json:"url,omitempty"`
	ClassName     string                 `json:"className,omitempty"`     // For styling
	ExtendedProps map[string]interface{} `json:"extendedProps,omitempty"` // For your custom data
	Rrule         string                 `json:"rrule,omitempty"`         // FC v5+ supports rrule
	Exdate        string                 `json:"exdate,omitempty"`        // For exceptions (comma separated ISO 8601 dates)
}

func getEventsAsFullcalendar() []FullCalendarEvent {
	state := getState()
	eventOccurances := state.EventOccurances
	events := state.Events

	ret := []FullCalendarEvent{}

	// Make a slice of full calendar formatted events
	for _, occurance := range eventOccurances {
		var matchingEvent *Event
		for i := range events {
			if events[i].ID == occurance.EventId {
				matchingEvent = &events[i]
				break
			}
		}

		if matchingEvent == nil {
			break
		}

		ret = append(ret, FullCalendarEvent{
			ID:        occurance.ID,
			Title:     matchingEvent.Name,
			Start:     occurance.At.String(),
			End:       occurance.At.String(),
			AllDay:    false,
			ClassName: matchingEvent.Type,
			Rrule:     matchingEvent.RRule,
		})
	}

	return ret
}
