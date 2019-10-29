package common

type NewReleaseEvent struct {
	Event       string `json:"event"`
	BuildId     int    `json:"buildId"`
	Branch      string `json:"branch"`
	DownloadUrl string `json:"downloadUrl"`
}
