package enforce_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/jacobtie/github-classrooms-deadline-enforcer/internal/config"
	"github.com/jacobtie/github-classrooms-deadline-enforcer/internal/enforce"
	"github.com/stretchr/testify/assert"
)

type repoState struct {
	PutUsername    string `json:"putUsername"`
	DeleteUsername string `json:"deleteUsername"`
}

func runTest(t *testing.T, date string) map[string]*repoState {
	cfg := config.Init()
	cfg.Test.TEST_DATE = date
	http.Get(fmt.Sprintf("%s/reset", cfg.GitHub.BaseURL))
	if err := enforce.Run(context.Background(), cfg); err != nil {
		t.Fatal(err)
	}
	resp, err := http.Get(fmt.Sprintf("%s/state", cfg.GitHub.BaseURL))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	var states map[string]*repoState
	if err := json.Unmarshal(b, &states); err != nil {
		t.Fatal(err)
	}
	assert.Len(t, states, 20)
	return states
}

func TestNoChanges(t *testing.T) {
	// Nothing is updated
	// 10th February 2023
	states := runTest(t, "2023-02-10")
	assert.Equal(t, map[string]*repoState{
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
	}, states)
}

func TestAssignmentWithoutExtensions(t *testing.T) {
	// Assignment is updated for all students
	// 11th February 2023
	states := runTest(t, "2023-02-11")
	assert.Equal(t, map[string]*repoState{
		"test-assignment-1-number-four": {
			PutUsername:    "number-four",
			DeleteUsername: "number-four",
		},
		"test-assignment-1-number-six": {
			PutUsername:    "number-six",
			DeleteUsername: "number-six",
		},
		"test-assignment-1-number-ten": {
			PutUsername:    "number-ten",
			DeleteUsername: "number-ten",
		},
		"test-assignment-1-number-twelve": {
			PutUsername:    "number-twelve",
			DeleteUsername: "number-twelve",
		},
		"test-assignment-1-number-thirteen": {
			PutUsername:    "number-thirteen",
			DeleteUsername: "number-thirteen",
		},
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
	}, states)
}

func TestAssignmentWithExtensions(t *testing.T) {
	// Assignment is updated but not for those with extensions
	// 12th February 2023
	states := runTest(t, "2023-02-12")
	assert.Equal(t, map[string]*repoState{
		"test-assignment-1-number-four":     {},
		"test-assignment-1-number-six":      {},
		"test-assignment-1-number-ten":      {},
		"test-assignment-1-number-twelve":   {},
		"test-assignment-1-number-thirteen": {},
		"test-assignment-2-number-four":     {},
		"test-assignment-2-number-six": {
			PutUsername:    "number-six",
			DeleteUsername: "number-six",
		},
		"test-assignment-2-number-ten":    {},
		"test-assignment-2-number-twelve": {},
		"test-assignment-2-number-thirteen": {
			PutUsername:    "number-thirteen",
			DeleteUsername: "number-thirteen",
		},
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
	}, states)
}

func TestAssignmentAndOtherExtensions(t *testing.T) {
	// Assignment is updated for all students and another assignment's extensions are updated
	// 13th February 2023
	states := runTest(t, "2023-02-13")
	assert.Equal(t, map[string]*repoState{
		"test-assignment-1-number-four":     {},
		"test-assignment-1-number-six":      {},
		"test-assignment-1-number-ten":      {},
		"test-assignment-1-number-twelve":   {},
		"test-assignment-1-number-thirteen": {},
		"test-assignment-2-number-four":     {},
		"test-assignment-2-number-six":      {},
		"test-assignment-2-number-ten": {
			PutUsername:    "number-ten",
			DeleteUsername: "number-ten",
		},
		"test-assignment-2-number-twelve":   {},
		"test-assignment-2-number-thirteen": {},
		"test-assignment-3-number-four": {
			PutUsername:    "number-four",
			DeleteUsername: "number-four",
		},
		"test-assignment-3-number-six": {
			PutUsername:    "number-six",
			DeleteUsername: "number-six",
		},
		"test-assignment-3-number-ten": {
			PutUsername:    "number-ten",
			DeleteUsername: "number-ten",
		},
		"test-assignment-3-number-twelve": {
			PutUsername:    "number-twelve",
			DeleteUsername: "number-twelve",
		},
		"test-assignment-3-number-thirteen": {
			PutUsername:    "number-thirteen",
			DeleteUsername: "number-thirteen",
		},
		"test-assignment-4-number-four":     {},
		"test-assignment-4-number-six":      {},
		"test-assignment-4-number-ten":      {},
		"test-assignment-4-number-twelve":   {},
		"test-assignment-4-number-thirteen": {},
	}, states)
}

func TestNoAssigmentAndOtherExceptions(t *testing.T) {
	// No assignment has all students updated and another assignment's extensions are updated
	// 14th February 2023
	states := runTest(t, "2023-02-14")
	assert.Equal(t, map[string]*repoState{
		"test-assignment-1-number-four":     {},
		"test-assignment-1-number-six":      {},
		"test-assignment-1-number-ten":      {},
		"test-assignment-1-number-twelve":   {},
		"test-assignment-1-number-thirteen": {},
		"test-assignment-2-number-four": {
			PutUsername:    "number-four",
			DeleteUsername: "number-four",
		},
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
	}, states)
}

func TestAssignmentWithExtensionsAndOtherExtensions(t *testing.T) {
	// Assignment is updated for all students except those with extensions and another assignment's extensions are updated
	// 15th February 2023
	states := runTest(t, "2023-02-15")
	assert.Equal(t, map[string]*repoState{
		"test-assignment-1-number-four":     {},
		"test-assignment-1-number-six":      {},
		"test-assignment-1-number-ten":      {},
		"test-assignment-1-number-twelve":   {},
		"test-assignment-1-number-thirteen": {},
		"test-assignment-2-number-four":     {},
		"test-assignment-2-number-six":      {},
		"test-assignment-2-number-ten":      {},
		"test-assignment-2-number-twelve": {
			PutUsername:    "number-twelve",
			DeleteUsername: "number-twelve",
		},
		"test-assignment-2-number-thirteen": {},
		"test-assignment-3-number-four":     {},
		"test-assignment-3-number-six":      {},
		"test-assignment-3-number-ten":      {},
		"test-assignment-3-number-twelve":   {},
		"test-assignment-3-number-thirteen": {},
		"test-assignment-4-number-four":     {},
		"test-assignment-4-number-six": {
			PutUsername:    "number-six",
			DeleteUsername: "number-six",
		},
		"test-assignment-4-number-ten": {
			PutUsername:    "number-ten",
			DeleteUsername: "number-ten",
		},
		"test-assignment-4-number-twelve": {
			PutUsername:    "number-twelve",
			DeleteUsername: "number-twelve",
		},
		"test-assignment-4-number-thirteen": {
			PutUsername:    "number-thirteen",
			DeleteUsername: "number-thirteen",
		},
	}, states)
}
