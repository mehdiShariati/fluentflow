package api

// Scenario is a static catalog (PRD §12.2). Client lists these before starting a session.
type Scenario struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Level       string `json:"level"`
}

var catalog = []Scenario{
	{ID: "coffee_shop", Title: "Ordering coffee", Description: "Practice ordering and small talk at a café.", Level: "A2"},
	{ID: "job_interview", Title: "Job interview", Description: "Answer common interview questions clearly.", Level: "B1"},
	{ID: "airport_immigration", Title: "Airport immigration", Description: "Border control questions and answers.", Level: "A2"},
	{ID: "business_small_talk", Title: "Business small talk", Description: "Networking and polite workplace chat.", Level: "B1"},
	{ID: "apartment_rental", Title: "Apartment rental", Description: "Ask about rent, utilities, and viewing.", Level: "B1"},
	{ID: "exam_speaking", Title: "Exam speaking simulation", Description: "Timed prompts similar to IELTS/TOEFL speaking.", Level: "B2"},
}

func listScenarios() []Scenario {
	out := make([]Scenario, len(catalog))
	copy(out, catalog)
	return out
}

// IsValidScenarioID returns true if id is in the static catalog.
func IsValidScenarioID(id string) bool {
	for _, sc := range catalog {
		if sc.ID == id {
			return true
		}
	}
	return false
}

// ScenarioTitle returns the catalog title for id, or empty if unknown.
func ScenarioTitle(id string) string {
	for _, sc := range catalog {
		if sc.ID == id {
			return sc.Title
		}
	}
	return ""
}
