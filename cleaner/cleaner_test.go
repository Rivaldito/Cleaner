package cleaner_test

import (
	"testing"
	"time"

	"github.com/Rivaldito/Cleaner/cleaner"
)

func TestNewCleaner(t *testing.T) {
	// Call NewCleaner twice with same arguments
	cleaner1 := cleaner.NewCleaner("test_dir", ".txt", 7)
	cleaner2 := cleaner.NewCleaner("test_dir", ".txt", 7)

	// Check if both instances point to the same memory location
	if cleaner1 != cleaner2 {
		t.Errorf("NewCleaner should return the same instance for subsequent calls")
	}

	// Check the values set in the cleaner object
	if cleaner1.OSFile != "test_dir" || cleaner1.GetFileExtensionToClean() != ".txt" || cleaner1.GetDaysDiffToClean() != 7 {
		t.Errorf("NewCleaner did not set the correct values in the cleaner object")
	}
}

func TestCheckExtension(t *testing.T) {
	cleaner := cleaner.NewCleaner("", "", 0)

	// Test with matching extension
	if !cleaner.CheckExtension("test_file.txt") {
		t.Errorf("checkExtension did not identify a matching extension")
	}

	// Test with non-matching extension
	if cleaner.CheckExtension("test_file.jpg") {
		t.Errorf("checkExtension identified a non-matching extension")
	}
}

func TestGettersAndSetters(t *testing.T) {
	cleaner := cleaner.NewCleaner("", "", 0)

	// Test setters and getters for file extension
	cleaner.SetFileExtensionToClean(".log")
	if cleaner.GetFileExtensionToClean() != ".log" {
		t.Errorf("SetFileExtensionToClean or GetFileExtensionToClean malfunctioned")
	}

	// Similar tests can be written for other getters and setters
}

func TestDateComparation(t *testing.T) {
	tests := []struct {
		name           string
		cleanerDays    int
		fileModTime    time.Time
		expectedResult bool
	}{
		{
			name:           "File older than clean threshold",
			cleanerDays:    1,                                   // Clean files older than 1 day
			fileModTime:    time.Now().Add(-time.Hour * 24 * 2), // File modified 2 days ago
			expectedResult: true,
		},
		{
			name:           "Recent file (not deleted)",
			cleanerDays:    1,                               // Clean files older than 1 day
			fileModTime:    time.Now().Add(-time.Hour * 12), // File modified 12 hours ago
			expectedResult: false,
		},
		{
			name:           "File exactly at threshold (not deleted)",
			cleanerDays:    1,                               // Clean files older than 1 day
			fileModTime:    time.Now().Add(-time.Hour * 24), // File modified 1 day ago
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cleaner := cleaner.NewCleaner("", "", tc.cleanerDays)

			if cleaner.DateComparation(tc.fileModTime.Unix()) != tc.expectedResult {
				t.Errorf("dateComparation did not meet expectation for test case: %s", tc.name)
			}
		})
	}
}
