package handler

const (
	// ExtraMsgLength defines the max length of extra message
	// ExtraMsgLength = ExtraMsgFlagLength + extra message ID length + extra message body length.
	ExtraMsgLength = 4

	// ExtraMsgFlagLength defines 4 as the extra message flag to avoid repeat with Msg ID.
	ExtraMsgFlagLength = 4

	// ExtraMsgFlag depends on ExtraMsgFlagLength, the max value can be describes by ExtraMsgFlagLength bytes.
	ExtraMsgFlag = 2 ^ (ExtraMsgFlagLength * 8) - 1
)
