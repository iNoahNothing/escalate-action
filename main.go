package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type IssueData struct {
	Title  string
	Body   string
	Labels []Label
}

type Label struct {
	Name string
}

type Issue struct {
	Url    string
	Id     uint
	Title  string
	Labels []Label
	Body   string
}

type Comment struct {
	Body string
}

type CommentObject struct {
	Action  string
	Issue   Issue
	Comment Comment
}

var commentBlob string
var accessToken string

func parseFlags() {

	flag.StringVar(&commentBlob, "issue", "", "The object of the latest comment.")
	flag.StringVar(&accessToken, "accessToken", "", "Used to authenticate to GitHub api.")
	flag.Parse()

	if commentBlob == "" {
		fmt.Printf("URL is required.\n")
		os.Exit(2)
	}

	if accessToken == "" {
		fmt.Printf("accessToken is required. \n")
	}
}

func createIssue(comment CommentObject) {

	var repo string
	var labels []Label
	for _, tag := range comment.Issue.Labels {
		if strings.Contains(tag.Name, "product/") {
			label := strings.Split(tag.Name, "/")
			repo = label[1]
		} else {
			labels = append(labels, tag)
		}
	}

	issueUrl := strings.Split(comment.Issue.Url, "/")

	owner := issueUrl[4]

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues", owner, repo)

	body := IssueData{Title: comment.Issue.Title, Body: comment.Issue.Body, Labels: labels}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	authHeader := fmt.Sprintf("token %s", accessToken)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", authHeader)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.Status)
		os.Exit(0)
	}

	fmt.Printf("Issue: %s created in repository: %s\n", comment.Issue.Title, repo)

}

func addLabel() {

}

func main() {

	parseFlags()

	var comment CommentObject

	err := json.Unmarshal([]byte(commentBlob), &comment)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	escalate, err := regexp.MatchString("(\n|^)\\/escalate(\n|$)", comment.Comment.Body)

	if escalate {

		createIssue(comment)

	}

}
