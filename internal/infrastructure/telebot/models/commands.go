package models

type CommandName string

const (
	CommandNameStart   CommandName = "start"
	CommandNameHelp    CommandName = "help"
	CommandNameCancel  CommandName = "cancel"
	CommandNameTrack   CommandName = "track"
	CommandNameList    CommandName = "list"
	CommandNameUntrack CommandName = "untrack"
)

func (c CommandName) String() string {
	return string(c)
}
