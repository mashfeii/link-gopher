package models

type (
	State       string
	SessionType string
)

const (
	StateTrackInputURL     State = "track_input_url"
	StateTrackInputTags    State = "track_input_tags"
	StateTrackInputFilters State = "track_input_filters"
	StateTrackConfirm      State = "track_confirm"

	StateListUntrackSelectTags    State = "untrack_select_tags"
	StateListUntrackSelectFilters State = "untrack_select_filters"
	StateUntrackSelectURL         State = "untrack_select_url"
	StateListUntrackConfirm       State = "untrack_confirm"

	StateListConfirm State = "list_confirm"
)
