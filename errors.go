package main

type ErrorList struct {
	List []string
}

var ErrorMessages = map[string]string{
	"EncodingError":       "There has been a problem with the encoding. ",
	"UrlError":            "Not all URL parameters filled. ",
	"addressError":        "Access denied. ",
	"hostError":           "Error selecting Docker host. ",
	"DBConnectionError":   "Can not access database. ",
	"HostNameExistsError": "That hostname seems to be in use. ",
	"DeploymentError":     "There was a problem deploying the Docker container. ",
}

var Messages = map[string]string{
	"DeploymentSuccess": "Deployment successful!",
}
