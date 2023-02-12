package enforce

import (
	"fmt"
	"time"

	"github.com/jacobtie/github-classrooms-deadline-enforcer/internal/config"
	"github.com/rs/zerolog/log"
)

const EASTERN_TIME_ZONE = "America/New_York"

const YYYY_MM_DD = "2006-01-02"

func resolveReposToEnforce(cfg *config.Config, enforceCfg *enforcerConfig) ([]*enforcement, error) {
	today, err := getCurrentDate(cfg)
	if err != nil {
		return nil, err
	}
	studentsToUsernames := getStudentsToUsernames(enforceCfg.Students)
	enforcements := make([]*enforcement, 0)
	for _, assignment := range enforceCfg.Assignments {
		extensionsToday, extensionsNotToday := bucketStudentExtensions(assignment.Extensions, today)
		if assignment.Deadline == today {
			// Add all students except those with extensions on other days
			for name, username := range studentsToUsernames {
				if _, ok := extensionsNotToday[name]; ok {
					continue
				}
				enforcements = append(enforcements, getEnforcement(assignment.Name, username))
			}
			continue
		}
		// Add all students who had extensions until today on other assignments
		for name := range extensionsToday {
			username := studentsToUsernames[name]
			enforcements = append(enforcements, getEnforcement(assignment.Name, username))
		}
	}
	return enforcements, nil
}

func getCurrentDate(cfg *config.Config) (string, error) {
	loc, err := time.LoadLocation(EASTERN_TIME_ZONE)
	if err != nil {
		return "", fmt.Errorf("[enforce.getCurrentDate] failed to get the America/New_York timezone: %w", err)
	}
	now := timeNow(cfg, loc)
	return now.Format(YYYY_MM_DD), nil
}

func timeNow(cfg *config.Config, loc *time.Location) time.Time {
	if cfg.Test.IS_TEST {
		t, err := time.Parse(YYYY_MM_DD, cfg.Test.TEST_DATE)
		if err != nil {
			log.Fatal().Err(err).Msg("could not parse test date")
		}
		return t
	}
	return time.Now().In(loc)
}

func getStudentsToUsernames(students []*enforcerConfigStudent) map[string]string {
	studentsToUsernames := make(map[string]string)
	for _, student := range students {
		studentsToUsernames[student.Name] = student.Username
	}
	return studentsToUsernames
}

func bucketStudentExtensions(extensions []*enforcerConfigAssignmentExtension, today string) (map[string]struct{}, map[string]struct{}) {
	// map[string]struct{} is functionally a set
	extensionsToday := make(map[string]struct{}, 0)
	extensionsNotToday := make(map[string]struct{}, 0)
	for _, extension := range extensions {
		if extension.Deadline == today {
			extensionsToday[extension.Name] = struct{}{}
		} else {
			extensionsNotToday[extension.Name] = struct{}{}
		}
	}
	return extensionsToday, extensionsNotToday
}

func getEnforcement(assignmentName, username string) *enforcement {
	repoName := fmt.Sprintf("%s-%s", assignmentName, username)
	return &enforcement{
		repoName: repoName,
		username: username,
	}
}
