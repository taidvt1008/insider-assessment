package constants

// Message status constants
const (
	MessageStatusPending = "pending"
	MessageStatusSent    = "sent"
	MessageStatusFailed  = "failed"
)

// MessageStatusValues returns all valid message status values
func MessageStatusValues() []string {
	return []string{
		MessageStatusPending,
		MessageStatusSent,
		MessageStatusFailed,
	}
}

// IsValidMessageStatus checks if the given status is valid
func IsValidMessageStatus(status string) bool {
	for _, validStatus := range MessageStatusValues() {
		if status == validStatus {
			return true
		}
	}
	return false
}
