package models

type RepositoryDownloadCount struct {
	TotalDownloads int `json:"total"`
}

type VersionDownloadCount struct {
	VersionDownloads int `json:"version_downloads"`
	TotalDownloads   int `json:"total_downloads"`
}
