package core

type ContextParamKey string

const (
	ContextKeyChatID ContextParamKey = "chat_id"
	ContextLinkURL   ContextParamKey = "link_url"
)

const (
	CommandNameStart  = "start"
	CommandNameHelp   = "help"
	CommandNameCancel = "cancel"
	CommandNameTrack  = "track"
)

const (
	CallbackReturnTags  = "return_tags"
	CallbackSkipFilters = "skip_filters"
	CallbackReturnURL   = "return_url"
	CallbackSkipTags    = "skip_tags"
)
