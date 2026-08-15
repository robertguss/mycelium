// Package lifecycle encodes PHASE-02 commanded state edges.
package lifecycle

// AllowedTargets are legal argv targets for mycelium state this phase.
func AllowedTargets() []string {
	return []string{"exploring", "simmering", "clarified", "archived"}
}

// LegalNext returns allowed commanded next states from from.
func LegalNext(from string) []string {
	switch from {
	case "spark":
		return []string{"exploring", "archived"}
	case "exploring":
		return []string{"simmering", "clarified", "archived"}
	case "simmering":
		return []string{"exploring", "archived"}
	case "clarified":
		return []string{"archived"}
	case "handed-off", "archived":
		return nil
	default:
		return nil
	}
}

// Legal reports whether from → to is a commanded edge this phase.
func Legal(from, to string) bool {
	for _, n := range LegalNext(from) {
		if n == to {
			return true
		}
	}
	return false
}

// IsWake is simmering → exploring.
func IsWake(from, to string) bool {
	return from == "simmering" && to == "exploring"
}

// RevisitRequired is true when target is simmering.
func RevisitRequired(to string) bool {
	return to == "simmering"
}

// RevisitForbidden is true when target is not simmering.
func RevisitForbidden(to string) bool {
	return to != "simmering"
}
