package main

import "github.com/bugsnag/bugsnag-go"

func init() {
	bugsnag.Configure(bugsnag.Configuration{
		APIKey:       "670bbdff860972d79e45819f01347765",
		ReleaseStage: "development",
	})
}
