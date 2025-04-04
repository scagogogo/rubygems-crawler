package models

import "time"

// PackageInformation
// Example:
// {
//    "name": "rails",
//    "downloads": 436090160,
//    "version": "7.0.5",
//    "version_created_at": "2023-05-24T19:21:28.229Z",
//    "version_downloads": 54428,
//    "platform": "ruby",
//    "authors": "David Heinemeier Hansson",
//    "info": "Ruby on Rails is a full-stack web framework optimized for programmer happiness and sustainable productivity. It encourages beautiful code by favoring convention over configuration.",
//    "licenses": [
//        "MIT"
//    ],
//    "metadata": {
//        "changelog_uri": "https://github.com/rails/rails/releases/tag/v7.0.5",
//        "bug_tracker_uri": "https://github.com/rails/rails/issues",
//        "source_code_uri": "https://github.com/rails/rails/tree/v7.0.5",
//        "mailing_list_uri": "https://discuss.rubyonrails.org/c/rubyonrails-talk",
//        "documentation_uri": "https://api.rubyonrails.org/v7.0.5/",
//        "rubygems_mfa_required": "true"
//    },
//    "yanked": false,
//    "sha": "57ef2baa4a1f5f954bc6e5a019b1fac8486ece36f79c1cf366e6de33210637fe",
//    "project_uri": "https://rubygems.org/gems/rails",
//    "gem_uri": "https://rubygems.org/gems/rails-7.0.5.gem",
//    "homepage_uri": "https://rubyonrails.org",
//    "wiki_uri": null,
//    "documentation_uri": "https://api.rubyonrails.org/v7.0.5/",
//    "mailing_list_uri": "https://discuss.rubyonrails.org/c/rubyonrails-talk",
//    "source_code_uri": "https://github.com/rails/rails/tree/v7.0.5",
//    "bug_tracker_uri": "https://github.com/rails/rails/issues",
//    "changelog_uri": "https://github.com/rails/rails/releases/tag/v7.0.5",
//    "funding_uri": null,
//    "dependencies": {
//        "development": [],
//        "runtime": [
//            {
//                "name": "actioncable",
//                "requirements": "= 7.0.5"
//            },
//            {
//                "name": "actionmailbox",
//                "requirements": "= 7.0.5"
//            },
//            {
//                "name": "actionmailer",
//                "requirements": "= 7.0.5"
//            },
//            {
//                "name": "actionpack",
//                "requirements": "= 7.0.5"
//            },
//            {
//                "name": "actiontext",
//                "requirements": "= 7.0.5"
//            },
//            {
//                "name": "actionview",
//                "requirements": "= 7.0.5"
//            },
//            {
//                "name": "activejob",
//                "requirements": "= 7.0.5"
//            },
//            {
//                "name": "activemodel",
//                "requirements": "= 7.0.5"
//            },
//            {
//                "name": "activerecord",
//                "requirements": "= 7.0.5"
//            },
//            {
//                "name": "activestorage",
//                "requirements": "= 7.0.5"
//            },
//            {
//                "name": "activesupport",
//                "requirements": "= 7.0.5"
//            },
//            {
//                "name": "bundler",
//                "requirements": ">= 1.15.0"
//            },
//            {
//                "name": "railties",
//                "requirements": "= 7.0.5"
//            }
//        ]
//    }
//}
type PackageInformation struct {
	Name             string       `json:"name"`
	Downloads        int          `json:"downloads"`
	Version          string       `json:"version"`
	VersionCreatedAt time.Time    `json:"version_created_at"`
	VersionDownloads int          `json:"version_downloads"`
	Platform         string       `json:"platform"`
	Authors          string       `json:"authors"`
	Info             string       `json:"info"`
	Licenses         []string     `json:"licenses"`
	Metadata         Metadata     `json:"metadata"`
	Yanked           bool         `json:"yanked"`
	Sha              string       `json:"sha"`
	ProjectURI       string       `json:"project_uri"`
	GemURI           string       `json:"gem_uri"`
	HomepageURI      string       `json:"homepage_uri"`
	WikiURI          interface{}  `json:"wiki_uri"`
	DocumentationURI string       `json:"documentation_uri"`
	MailingListURI   string       `json:"mailing_list_uri"`
	SourceCodeURI    string       `json:"source_code_uri"`
	BugTrackerURI    string       `json:"bug_tracker_uri"`
	ChangelogURI     string       `json:"changelog_uri"`
	FundingURI       interface{}  `json:"funding_uri"`
	Dependencies     Dependencies `json:"dependencies"`
}

type Dependencies struct {
	Development []*Dependency `json:"development"`
	Runtime     []*Dependency `json:"runtime"`
}

type Dependency struct {
	Name         string `json:"name"`
	Requirements string `json:"requirements"`
}
