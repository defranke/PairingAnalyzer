package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type CountItem struct {
	Key   string
	Count int
}

func main() {
	repoPath := ""
	if len(os.Args) > 1 {
		repoPath = os.Args[1]
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("")
	fmt.Println("")
	fmt.Println("################################################")
	fmt.Println("🎉🎉🎉 Welcome to the Pairing Analyzer! 🎉🎉🎉")
	fmt.Println("################################################")
	fmt.Println("")
	for repoPath == "" {
		fmt.Print("Path to Git Repository: ")
		pathBytes, _ := reader.ReadString('\n')
		repoPath = string(pathBytes)
	}

	repo := NewGitRepository(strings.Trim(repoPath, "\n"))
	fmt.Println("⏳ Analyzing...")
	repo.Analyze()

	for true {
		fmt.Println("What do you want to do?")
		fmt.Println("[1] 👬 Output Pairing Statistics")
		fmt.Println("[2] 🕺 Output Single Statistics")
		fmt.Println("[3] 👋 Exit")

		fmt.Print("➡️  ")
		cmdBytes, _ := reader.ReadString('\n')
		cmd := strings.Trim(string(cmdBytes), "\n")

		if cmd == "1" {
			repo.OutputPairStats()
		} else if cmd == "2" {
			repo.OutputSingleStats()
		} else if cmd == "3" {
			fmt.Println("Bye 👋")
			break
		} else {
			fmt.Println("❓ Unknown command")
		}
		fmt.Println("\n")
	}

	// repos := make([]string, 0)

	// if len(os.Args) > 1 {
	// 	for i := 1; i < len(os.Args); i++ {
	// 		repos = append(repos, os.Args[i])
	// 	}
	// }

	// globalAuthorCount := make([]*CountItem, 0)
	// globalPairingCount := make([]*CountItem, 0)

	// for _, path := range repos {
	// 	authorCount, pairingCount := analyzeRepository(path)

	// 	printData(path, authorCount, pairingCount)
	// 	print("\n\n\n")

	// 	globalAuthorCount = mergeLists(authorCount, globalAuthorCount)
	// 	globalPairingCount = mergeLists(pairingCount, globalPairingCount)
	// }

	// printData("TOTAL", globalAuthorCount, globalPairingCount)
}

func analyzeRepository(path string) ([]*CountItem, []*CountItem) {
	authorCount := make([]*CountItem, 0)
	pairingCount := make([]*CountItem, 0)

	r, err := git.PlainOpen(path)
	if err != nil {
		println("Failed to open Git: ", err)
		os.Exit(1)
	}

	ref, _ := r.Head()

	authorMatcher, err := regexp.Compile("Signed-off-by: (.*) <.*>")
	if err != nil {
		println("Failed to create Regex: ", err)
		os.Exit(1)
	}

	commitLogs, err := r.Log(&git.LogOptions{From: ref.Hash()})

	err = commitLogs.ForEach(func(c *object.Commit) error {
		firstAuthor := c.Author.Name
		secondAuthorMatches := authorMatcher.FindStringSubmatch(c.Message)
		if len(secondAuthorMatches) > 0 {
			secondAuthor := secondAuthorMatches[1]

			authorCount = updateCount(authorCount, firstAuthor)
			authorCount = updateCount(authorCount, secondAuthor)

			var combined string
			if strings.Compare(firstAuthor, secondAuthor) < 0 {
				combined = fmt.Sprintf("%s + %s", firstAuthor, secondAuthor)
			} else {
				combined = fmt.Sprintf("%s + %s", secondAuthor, firstAuthor)
			}
			pairingCount = updateCount(pairingCount, combined)
		}
		return nil
	})
	return authorCount, pairingCount
}

func updateCount(data []*CountItem, key string) []*CountItem {
	for _, item := range data {
		if item.Key == key {
			item.Count++
			return data
		}
	}
	newItem := CountItem{key, 1}
	return append(data, &newItem)
}

func sortItems(data []*CountItem) {
	sort.Slice(data, func(i int, j int) bool {
		return data[i].Count > data[j].Count
	})
}

func printData(label string, authorCount []*CountItem, pairCount []*CountItem) {
	println("GLOBAL")
	println("============ Commits per Author ============")
	printItems(authorCount)
	print("\n")
	println("============ Commits per Pair ============")
	printItems(pairCount)
}

func printItems(data []*CountItem) {
	sortItems(data)
	for _, item := range data {
		fmt.Printf("%s:%s%d\n", item.Key, strings.Repeat(" ", 50-len(item.Key)), item.Count)
	}
}

func mergeLists(source []*CountItem, target []*CountItem) []*CountItem {
OUTER:
	for _, item := range source {
		for _, globalItem := range target {
			if item.Key == globalItem.Key {
				globalItem.Count += item.Count
				continue OUTER
			}
		}
		target = append(target, item)
	}
	return target
}
