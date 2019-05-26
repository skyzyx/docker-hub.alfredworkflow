// Copyright (c) 2019 Ryan Parman <https://ryanparman.com>
// Copyright (c) 2019 Contributors <https://github.com/skyzyx/terraform-registry.alfredworkflow/graphs/contributors>
//
// https://www.alfredapp.com/help/workflows/inputs/script-filter/json/

package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	// "github.com/davecgh/go-spew/spew"
	"github.com/parnurzeal/gorequest"
)

var (
	doc         InputDocument
	querystring string
)

type InputDocument struct {
	Count    int64         `json:"count"`
	Next     string        `json:"next"`
	Previous string        `json:"previous"`
	Results  []InputResult `json:"results"`
}

type InputResult struct {
	RepoName         string `json:"repo_name"`
	ShortDescription string `json:"short_description"`
	StarCount        int64  `json:"star_count"`
	PullCount        int64  `json:"pull_count"`
	RepoOwner        string `json:"repo_owner"`
	IsAutomated      bool   `json:"is_automated"`
	IsOfficial       bool   `json:"is_official"`
}

type AlfredDocument struct {
	Items []AlfredItem `json:"items,omitempty"`
}

type AlfredItem struct {
	// Simple objects
	Arg          string `json:"arg,omitempty"`
	Autocomplete string `json:"autocomplete,omitempty"`
	Match        string `json:"match,omitempty"`
	QuicklookUrl string `json:"quicklookurl,omitempty"`
	Subtitle     string `json:"subtitle,omitempty"`
	Title        string `json:"title,omitempty"`
	Type         string `json:"type,omitempty"`
	UID          string `json:"uid,omitempty"`
	Valid        bool   `json:"valid,omitempty"`

	// Complex objects
	Icon AlfredIcon         `json:"icon,omitempty"`
	Mods AlfredModifierKeys `json:"mods,omitempty"`
	Text AlfredText         `json:"text,omitempty"`
}

type AlfredIcon struct {
	Path string `json:"path,omitempty"`
	Type string `json:"type,omitempty"`
}

type AlfredText struct {
	Copy      string `json:"copy,omitempty"`
	LargeType string `json:"largetype,omitempty"`
}

type AlfredModifierKeys struct {
	Alt     AlfredModifierKey `json:"alt,omitempty"`
	Command AlfredModifierKey `json:"cmd,omitempty"`
}

type AlfredModifierKey struct {
	Arg          string `json:"arg,omitempty"`
	Subtitle     string `json:"subtitle,omitempty"`
	QuicklookUrl string `json:"quicklookurl,omitempty"`
	Valid        bool   `json:"valid,omitempty"`
}

// The core function
func main() {
	alfred := new(AlfredDocument)
	querystring = url.QueryEscape(os.Args[1])

	_, body, _ := gorequest.New().Get("https://hub.docker.com/v2/search/repositories?query=" + querystring).End()
	err := json.Unmarshal([]byte(body), &doc)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// spew.Dump(doc)

	// No results
	if doc.Count == 0 {
		alfred.Items = append(alfred.Items, AlfredItem{
			Title: "No results found.",
			Valid: false,
			Type:  "default",
			Icon: AlfredIcon{
				Path: "hub.png",
			},
		})
	}

	// Results
	for _, result := range doc.Results {
		regUrl := fmt.Sprintf(
			"https://hub.docker.com/%s",
			map[bool]string{
				true:  fmt.Sprintf("_/%s", result.RepoName),
				false: fmt.Sprintf("r/%s", result.RepoName),
			}[result.IsOfficial],
		)

		alfred.Items = append(alfred.Items, AlfredItem{
			UID:          result.RepoName,
			Title:        result.RepoName,
			Subtitle:     result.ShortDescription,
			Arg:          regUrl,
			QuicklookUrl: regUrl,
			Valid:        true,
			Type:         "default",
			// Autocomplete string `json:"autocomplete,omitempty"`
			// Match        string `json:"match,omitempty"`
			Icon: AlfredIcon{
				// Type: "fileicon",
				Path: map[bool]string{true: "verified.png", false: "not-verified.png"}[result.IsOfficial],
			},
			Text: AlfredText{
				Copy:      result.RepoName,
				LargeType: result.RepoName,
			},
			Mods: AlfredModifierKeys{
				Alt: AlfredModifierKey{
					Arg: regUrl,
					Subtitle: fmt.Sprintf(
						"%d %s • %d %s • %s",
						result.StarCount,
						map[bool]string{true: "star", false: "stars"}[result.StarCount == 1],
						result.PullCount,
						map[bool]string{true: "pull", false: "pulls"}[result.PullCount == 1],
						map[bool]string{true: "Automated", false: "Not Automated"}[result.IsAutomated],
					),
					QuicklookUrl: regUrl,
					Valid:        true,
				},
			},
		})
	}

	// output, err := json.Marshal(alfred)
	output, err := json.MarshalIndent(alfred, "", "    ")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(string(output))
}
