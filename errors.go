package main

type ErrorList struct {
	List []string
}

var ErrorMessages = map[string]string{
	"EncodingError":     "There has been a problem with the encoding.",
	"UrlError":          "Not all URL parameters filled.",
	"addressError":      "Access denied",
	"DBConnectionError": "Can not access database",
}

var Messages = map[string]string{
	"DeploymentSuccess": "Deployment successful!",
}
