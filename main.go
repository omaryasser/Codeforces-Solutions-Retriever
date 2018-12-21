package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Type Your Codeforces Handle.")
	var handle string
	fmt.Scan(&handle)
	fmt.Println("Please Wait ...")
	submissions := getSubmissions(handle)
	acceptedProblems := make([]AcceptedProblem, 0, len(submissions))
	for _, v := range submissions {
		if v.Verdict == "OK" {
			acceptedProblem := AcceptedProblem{v.ContestID, v.ID, v.Problem.Name, v.ProgrammingLanguage, v.Problem.Index}
			acceptedProblems = append(acceptedProblems, acceptedProblem)
		}
	}

	folderName := "./Codeforces Solutions"
	err := os.Mkdir(folderName, 0700)
	if err != nil {
		panic(err)
	}
	for _, acceptedProblem := range acceptedProblems {
		f, err := os.OpenFile(folderName+"/"+acceptedProblem.getFileName(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0700)
		if err != nil {
			panic(err)
		}
		_, err = f.WriteString(acceptedProblem.getLink())
		if err != nil {
			panic(err)
		}
		fmt.Println("Created File:", acceptedProblem.getFileName())
		f.Close()
	}
}

func getSubmissions(handle string) Submissions {
	resp, err := http.Get("http://codeforces.com/api/user.status?handle=" + handle + "&from=1&count=10000")
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("HTTP Request Error! Check Your Connection and Check whether codeforces is running or not. : ", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	arr := status{}
	json.Unmarshal(body, &arr)
	return arr.Result
}

type Submissions []struct {
	ID                  int   `json:"id"`
	ContestID           int   `json:"contestId"`
	CreationTimeSeconds int   `json:"creationTimeSeconds"`
	RelativeTimeSeconds int64 `json:"relativeTimeSeconds"`
	Problem             struct {
		ContestID int      `json:"contestId"`
		Index     string   `json:"index"`
		Name      string   `json:"name"`
		Type      string   `json:"type"`
		Tags      []string `json:"tags"`
	} `json:"problem"`
	Author struct {
		ContestID int `json:"contestId"`
		Members   []struct {
			Handle string `json:"handle"`
		} `json:"members"`
		ParticipantType  string `json:"participantType"`
		Ghost            bool   `json:"ghost"`
		StartTimeSeconds int    `json:"startTimeSeconds"`
	} `json:"author"`
	ProgrammingLanguage string `json:"programmingLanguage"`
	Verdict             string `json:"verdict"`
	Testset             string `json:"testset"`
	PassedTestCount     int    `json:"passedTestCount"`
	TimeConsumedMillis  int    `json:"timeConsumedMillis"`
	MemoryConsumedBytes int    `json:"memoryConsumedBytes"`
}
type status struct {
	Status string      `json:"status"`
	Result Submissions `json:"result"`
}

type AcceptedProblem struct {
	contestId, submissionID int
	name, language, index   string
}

func (l AcceptedProblem) getLink() string {
	var gymOrContest string
	if l.contestId >= 10000 {
		gymOrContest = "gym"
	} else {
		gymOrContest = "contest"
	}
	return "https://codeforces.com/" + gymOrContest + "/" + strconv.Itoa(l.contestId) + "/submission/" + strconv.Itoa(l.submissionID)
}

func (l AcceptedProblem) getFileName() string {
	return strconv.Itoa(l.contestId) + "-" + l.index + "_" + normalizeProblemName(l.name) + "." + normalizeLanguageName(l.language)
}

func normalizeLanguageName(s string) string {
	if strings.Contains(s, "Java") || strings.Contains(s, "java") {
		return "java"
	} else if strings.Contains(s, "Gnu") || strings.Contains(s, "C++") || strings.Contains(s, "c++") {
		return "cpp"
	} else if strings.Contains(s, "Python") || strings.Contains(s, "python"){
		return "python"
	} else if strings.Contains(s, "Go") || strings.Contains(s, "go"){
		return "go"
	}
	return s
}

func normalizeProblemName(name string) string {
	splitted := strings.Split(name, " ")
	var res string
	for _, v := range splitted {
		res += removeSlashes(v)
	}
	return res
}

func removeSlashes(s string) string {
	var res string
	for i := 0; i < len(s); i++ {
		if s[i] != '/' {
			res += string(s[i])
		}
	}
	return res
}
