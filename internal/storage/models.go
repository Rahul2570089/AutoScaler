package storage

import "encoding/json"

const (
	GroupPrefix           = "group:"
	PolicyPrefix          = "policy:"
	ScheduledActionPrefix = "schedule:"
	ActivityPrefix        = "activity:"
)

func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func GroupKey(name string) []byte {
	return []byte(GroupPrefix + name)
}

func PolicyKey(name string) []byte {
	return []byte(PolicyPrefix + name)
}

func ScheduledActionKey(name string) []byte {
	return []byte(ScheduledActionPrefix + name)
}

func ActivityKey(name string) []byte {
	return []byte(ActivityPrefix + name)
}
