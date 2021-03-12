package common

// Message ...
var Message []string

// AddMessage ...
func AddMessage(s ...string) {
	Message = append(Message, s...)
}
