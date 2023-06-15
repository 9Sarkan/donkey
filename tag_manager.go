package donkey

// donkeyTag this struct responsible to manage donkey tag in each field
// fields activation status and also validation handled by this struct
type donkeyTag struct {
}

// IsActive check if current field is a donkey return true
func (d *donkeyTag) IsActive() bool {
	return false
}

// Name return current donkeys name
func (d *donkeyTag) Name() string {
	return ""
}

// IsRequired check if current fields value is required return true
func (d *donkeyTag) IsRequired() bool {
	return false
}

// IsValid validate current field value
func (d *donkeyTag) IsValid(value any) bool {
	return false
}
