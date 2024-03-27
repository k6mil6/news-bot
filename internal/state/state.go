package state

type State string

const (
	None                          State = "none"
	WaitingForSourceName          State = "waiting_for_source_name"
	WaitingForSourceURL           State = "waiting_for_source_url"
	WaitingForSourcePriority      State = "waiting_for_source_priority"
	WaitingForSourceIDAndPriority State = "waiting_for_source_id_and_priority"
)
