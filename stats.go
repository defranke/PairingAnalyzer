package main

import (
	"fmt"
)

type Stats struct {
	numberOfCommits int
	insertions      int
	deletions       int
}

func (s *Stats) Add(stats Stats) {
	s.numberOfCommits += stats.numberOfCommits
	s.insertions += stats.insertions
	s.deletions += stats.deletions
}

func (s *Stats) String() string {
	return fmt.Sprintf("Number of Commits %d, Insertions %d, Deletions: %d", s.numberOfCommits, s.insertions, s.deletions)
}