package provider

func IsValidResourceType(resourceType string) bool {
	switch ResourceType(resourceType) {
	case Database:
		return true
	case Collection:
		return true
	case Record:
		return true
	default:
		return false
	}
}
