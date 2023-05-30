package models

type Metadata struct {
	DocumentationURI    string `json:"documentation_uri"`
	BugTrackerURI       string `json:"bug_tracker_uri"`
	MailingListURI      string `json:"mailing_list_uri"`
	ChangelogURI        string `json:"changelog_uri"`
	SourceCodeURI       string `json:"source_code_uri"`
	RubygemsMfaRequired string `json:"rubygems_mfa_required"`
	WikiURI             string `json:"wiki_uri"`
	HomepageURI         string `json:"homepage_uri"`
}
