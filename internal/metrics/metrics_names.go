package metrics

type EventType string

const (
	EventTypeCreate   EventType = "create"
	EventTypeRedirect EventType = "redirect"
)

type ResponseType string

const (
	StatusOk            ResponseType = "200"
	StatusBadRequest    ResponseType = "400"
	StatusInternalError ResponseType = "500"
)
