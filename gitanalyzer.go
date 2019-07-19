package main

import (
	"fmt"
	"os/exec"
	"strings"
)

type Pair struct {
	author1 string
	author2 string
}

type GitRepository struct {
	path        string
	singleStats map[string]*Stats
	pairStats   map[Pair]*Stats
}

func NewGitRepository(path string) *GitRepository {
	return &GitRepository{path, make(map[string]*Stats), make(map[Pair]*Stats)}
}

func (r *GitRepository) Analyze() {
	cmd := exec.Command("git", "log", "--shortstat", "--pretty=format:'COMMIT_START###Author: %an###Date: %aD###Message: %B'")
	cmd.Dir = r.path
	out, _ := cmd.Output()
	result := string(out)

	commits := strings.Split(result, "COMMIT_START###")
	fmt.Println("Analyzing", len(commits), "Commits")
	//ioutil.WriteFile("test.txt", out, 0644)
	for _, commitLog := range commits {
		commit, err := NewCommit(commitLog)
		if err == nil {
			r.AnalyzeCommit(commit)
		} else {
			fmt.Println(err.Error())
		}
	}
}

func (r *GitRepository) OutputPairStats() {
	for pair, stats := range r.pairStats {
		fmt.Println("")
		fmt.Println("Pair: ", pair.author1, "+", pair.author2)
		fmt.Println(stats.String())
	}
}

func (r *GitRepository) OutputSingleStats() {
	for author, stats := range r.singleStats {
		fmt.Println("")
		fmt.Println("Author: ", author)
		fmt.Println("Number Of Commits:", stats.numberOfCommits, "Insertions:", stats.insertions, "Deletions:", stats.deletions)
	}
}

func (r *GitRepository) AnalyzeCommit(commit *Commit) {
	pair := Pair{commit.firstAuthor, commit.secondAuthor}
	commitStats := Stats{1, commit.insertions, commit.deletions}

	pairStats := r.getStatsForPair(pair)
	pairStats.Add(commitStats)
	r.setStatsForPair(pair, pairStats)

	author1Stats := r.getStatsForSingleAuthor(commit.firstAuthor)
	author1Stats.Add(commitStats)
	r.setStatsForSingleAuthor(commit.firstAuthor, author1Stats)

	if commit.firstAuthor != commit.secondAuthor {
		author2Stats := r.getStatsForSingleAuthor(commit.secondAuthor)
		author2Stats.Add(commitStats)
		r.setStatsForSingleAuthor(commit.secondAuthor, author2Stats)
	}
}

func (r *GitRepository) getStatsForSingleAuthor(author string) *Stats {
	val, ok := r.singleStats[author]
	if !ok {
		return &Stats{0, 0, 0}
	}
	return val
}

func (r *GitRepository) setStatsForSingleAuthor(author string, stats *Stats) {
	r.singleStats[author] = stats
}

func (r *GitRepository) getStatsForPair(pair Pair) *Stats {
	val, ok := r.pairStats[pair]
	if !ok {
		return &Stats{0, 0, 0}
	}
	return val
}

func (r *GitRepository) setStatsForPair(pair Pair, stats *Stats) {
	r.pairStats[pair] = stats
}

// func (r *GitRepository) getCommitInfo(commit string) (*CommitInfo, error) {
// 	firstAuthorMatches := r.firstAuthorMatcher.FindStringSubmatch(commit)
// 	if len(firstAuthorMatches) == 0 {
// 		return nil, errors.New("Failed to find first author")
// 	}
// 	secondAuthorMatches := r.secondAuthorMatcher.FindStringSubmatch(commit)
// 	if len(secondAuthorMatches) == 0 {
// 		secondAuthorMatches = firstAuthorMatches
// 	}
// 	dateMatches := r.dateMatcher.FindStringSubmatch(commit)
// 	if len(dateMatches) == 0 {
// 		return nil, errors.New("Failed to find date")
// 	}
// 	insertionMatches := r.insertionMatcher.FindStringSubmatch(commit)
// 	insertions := 0
// 	if len(insertionMatches) > 1 {
// 		insertions, _ = strconv.Atoi(insertionMatches[1])
// 	}
// 	deletionsMatches := r.deletionMatcher.FindStringSubmatch(commit)
// 	deletions := 0
// 	if len(deletionsMatches) > 1 {
// 		deletions, _ = strconv.Atoi(deletionsMatches[1])
// 	}

// 	firstAuthor := firstAuthorMatches[1]
// 	secondAuthor := secondAuthorMatches[1]
// 	date, _ := time.Parse(time.ANSIC, dateMatches[1])

// 	if strings.Compare(firstAuthor, secondAuthor) > 0 {
// 		tmp := firstAuthor
// 		firstAuthor = secondAuthor
// 		secondAuthor = tmp
// 	}

// 	return &CommitInfo{firstAuthor, secondAuthor, date, insertions, deletions}, nil
// }
