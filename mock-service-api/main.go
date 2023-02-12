package main

import (
	_ "embed"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type repoState struct {
	PutUsername    string `json:"putUsername"`
	DeleteUsername string `json:"deleteUsername"`
}

var repoStates map[string]*repoState

//go:embed config.json
var testConfig []byte

var pathRE = regexp.MustCompile("^/repos/([^/]+)/([^/]+)/([^/]+)/([^/]+)$")

const TEST_ORG_NAME = "test-org"
const TEST_CONFIG_REPO_NAME = "test-config-repo"
const TEST_CONFIG_FILENAME = "config.json"

func init() {
	resetState()
}

func main() {
	http.HandleFunc("/github/reset", func(w http.ResponseWriter, r *http.Request) {
		resetState()
		w.WriteHeader(http.StatusOK)
	})
	http.HandleFunc("/github/state", func(w http.ResponseWriter, r *http.Request) {
		b, err := json.Marshal(&repoStates)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	})
	http.HandleFunc("/github/", handleGitHub)
	log.Default().Println("Starting on port 3000")
	http.ListenAndServe(":3000", nil)
}

func handleGitHub(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != "Bearer test-token" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Not Found"))
	}
	subPath := strings.Split(r.URL.Path, "/github")[1]
	pathSegments := pathRE.FindStringSubmatch(subPath)
	if len(pathSegments) != 5 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
		return
	}
	orgName := pathSegments[1]
	repoName := pathSegments[2]
	route := pathSegments[3]
	param := pathSegments[4]
	if orgName != TEST_ORG_NAME {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
		return
	}
	if route == "contents" && r.Method == http.MethodGet {
		handleGetContents(w, r, repoName, param)
		return
	}
	if route == "collaborators" && r.Method == http.MethodPut {
		handlePutCollaborator(w, r, repoName, param)
		return
	}
	if route == "collaborators" && r.Method == http.MethodDelete {
		handleDeleteCollaborator(w, r, repoName, param)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Not Found"))
}

func handleGetContents(w http.ResponseWriter, r *http.Request, repoName, fileName string) {
	if repoName != TEST_CONFIG_REPO_NAME || fileName != TEST_CONFIG_FILENAME {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
		return
	}
	if r.Header.Get("Accept") != "application/vnd.github.raw" {
		// Not a typical response from GitHub but we want the tests to fail
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
		return
	}
	w.Header().Add("Content-Type", "application/vnd.github.raw")
	w.WriteHeader(http.StatusOK)
	w.Write(testConfig)
}

func handlePutCollaborator(w http.ResponseWriter, r *http.Request, repoName, username string) {
	if _, ok := repoStates[repoName]; !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
		return
	}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
		return
	}
	if string(b) != `{"permission":"pull"}` {
		// Not a typical response from GitHub but we want the tests to fail
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
		return
	}
	repoStates[repoName].PutUsername = username
	w.WriteHeader(http.StatusCreated)
	// Not a typical response from GitHub but we do not read the response
	w.Write([]byte("Done"))
}

func handleDeleteCollaborator(w http.ResponseWriter, r *http.Request, repoName, username string) {
	if _, ok := repoStates[repoName]; !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
		return
	}
	repoStates[repoName].DeleteUsername = username
	w.WriteHeader(http.StatusNoContent)
	// Not a typical response from GitHub but we do not read the response
	w.Write([]byte("Done"))
}

func resetState() {
	repoStates = map[string]*repoState{
		"test-assignment-1-number-four":     {},
		"test-assignment-1-number-six":      {},
		"test-assignment-1-number-ten":      {},
		"test-assignment-1-number-twelve":   {},
		"test-assignment-1-number-thirteen": {},
		"test-assignment-2-number-four":     {},
		"test-assignment-2-number-six":      {},
		"test-assignment-2-number-ten":      {},
		"test-assignment-2-number-twelve":   {},
		"test-assignment-2-number-thirteen": {},
		"test-assignment-3-number-four":     {},
		"test-assignment-3-number-six":      {},
		"test-assignment-3-number-ten":      {},
		"test-assignment-3-number-twelve":   {},
		"test-assignment-3-number-thirteen": {},
		"test-assignment-4-number-four":     {},
		"test-assignment-4-number-six":      {},
		"test-assignment-4-number-ten":      {},
		"test-assignment-4-number-twelve":   {},
		"test-assignment-4-number-thirteen": {},
	}
}
