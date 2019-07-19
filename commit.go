package main


import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Commit struct {
	firstAuthor  string
	secondAuthor string
	date         time.Time
	insertions   int
	deletions    int
}

var firstAuthorRegexp, _ = regexp.Compile("Author: (.*?)###")
var secondAuthorRegexp, _ = regexp.Compile("Signed-off-by: (.*) <.*>")
var dateMatcher, _ = regexp.Compile("Date: (.*?)###")
var insertionMatcher, _ = regexp.Compile("([0-9]*) insertion(s)*\\(\\+\\)")
var deletionMatcher, _ = regexp.Compile("([0-9]*) deletion(s)*\\(-\\)")

func NewCommit(commit string) (*Commit, error) {
	firstAuthorMatches := firstAuthorRegexp.FindStringSubmatch(commit)
	if len(firstAuthorMatches) == 0 {
		return nil, errors.New("Failed to find first author")
	}
	secondAuthorMatches := secondAuthorRegexp.FindStringSubmatch(commit)
	if len(secondAuthorMatches) == 0 {
		secondAuthorMatches = firstAuthorMatches
	}
	dateMatches := dateMatcher.FindStringSubmatch(commit)
	if len(dateMatches) == 0 {
		return nil, errors.New("Failed to find date")
	}
	insertionMatches := insertionMatcher.FindStringSubmatch(commit)
	insertions := 0
	if len(insertionMatches) > 1 {
		insertions, _ = strconv.Atoi(insertionMatches[1])
	}
	deletionsMatches := deletionMatcher.FindStringSubmatch(commit)
	deletions := 0
	if len(deletionsMatches) > 1 {
		deletions, _ = strconv.Atoi(deletionsMatches[1])
	}

	firstAuthor := firstAuthorMatches[1]
	secondAuthor := secondAuthorMatches[1]
	date, _ := time.Parse(time.ANSIC, dateMatches[1])

	if strings.Compare(firstAuthor, secondAuthor) > 0 {
		tmp := firstAuthor
		firstAuthor = secondAuthor
		secondAuthor = tmp
	}

	return &Commit{firstAuthor, secondAuthor, date, insertions, deletions}, nil
}