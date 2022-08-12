package model

// DestinationInput missing godoc
type DestinationInput struct {
	Name                    string `json:"Name"`
	Type                    string `json:"Type"`
	URL                     string `json:"URL"`
	Authentication          string `json:"Authentication"`
	XFSystemName            string `json:"XFSystemName"`
	CommunicationScenarioID string `json:"communicationScenarioId"`
	ProductName             string `json:"product.name"`
}
