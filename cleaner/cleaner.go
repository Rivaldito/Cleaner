package cleaner

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/charmbracelet/log"
)

// The number of seconds in a day
const secondsInADay int64 = 86400

// Cleaner struct.
//   - OSFile: The root were the cleaner remove the old data.
//   - fileExtension: The extension file of the files that will be "cleaned".
//   - files: the list of the all files in the dir
//   - currentTime: The current time in UNIX format
//   - daysDiffToClean: The numbers of day to remove a file
//   - secondsDiffToClean: The numbers of seconds to remove a file
type clean struct {
	OSFile             string
	fileExtension      string
	files              []fs.FileInfo
	currentTime        int64
	daysDiffToClean    int
	secondsDiffToClean int64
}

// Singleton instance for the struct
var cl *clean

// Ensures concurrency safety
var once sync.Once

// Constructor of the cleaner
//
// Args
//   - OSFile: The root were the cleaner remove the old data.
//   - fileExtension: The extension file of the files that will be "cleaned".
//   - daysDiffToClean: The numbers of day to remove a file
//
// Return
//
//   - clean: An instace of clean, if the instance is created previously the NewCleaner return the same instance
func NewCleaner(osFile string, fileExtension string, daysDiffToClean int) *clean {

	// heck if the instance exist and also prevent race events
	once.Do(func() {
		cl = &clean{
			OSFile:             osFile,
			fileExtension:      fileExtension,
			daysDiffToClean:    daysDiffToClean,
			secondsDiffToClean: int64(daysDiffToClean) * secondsInADay,
		}
	})
	// Singleton instance
	return cl
}

// Read the entire files of a OSFile of the clean structure
//
// The ReadFiles method save all the information of the files in a slice []fs.FileInfo in the clean structure
func (c *clean) ReadFiles() {
	// Open the dierctory
	f, err := os.Open(c.OSFile)
	if err != nil {
		log.Error(err.Error())
	}
	// Read all the files in the directory
	files, err := f.Readdir(0)
	if err != nil {
		log.Error(err.Error())
	}
	// close the directory connection
	defer f.Close()
	c.files = files
}

//	Check if the extension of the file match with the extension set in the clean structure.
//
// Args
//   - fileName: the name of the file, including the exension
//
// Return.
//   - True: if the file extension match
//   - False: if the file extension not match
func (c clean) CheckExtension(fileName string) bool {
	//Check the file extension
	fileExtension := filepath.Ext(fileName)
	if fileExtension == c.fileExtension {
		return true
	} else {
		return false
	}
}

//	Check if the last modification of the file
//
// Return.
//   - True: if the last modification is major or equal to the secondsDiffToClean
//   - False: if the last modification is less to the secondsDiffToClean
func (c clean) DateComparation(unixTimeFile int64) bool {
	//Get the difference between the (currente time) - (last modification of the file)
	diffTime := c.currentTime - unixTimeFile
	if diffTime >= c.secondsDiffToClean {
		return true
	} else {
		return false
	}
}

// Cleaner
//
// The cleaner execute the clean removes all the files, that match with the extension and the last modification date
func (c *clean) Cleaner() {

	var (
		isExtension  bool
		isDateDelete bool
		filePath     string
	)
	//Set the current time in the object
	c.SetUNIXTimeNow()
	//Read the files
	c.ReadFiles()
	//itarate through the all files
	for _, file := range c.files {
		//start conditional evaluation
		isExtension = c.CheckExtension(file.Name())
		//Check if the file has the extension permitted to delete
		if !isExtension {
			continue
		}
		//Check if the "file" is a dir
		if file.IsDir() {
			continue
		}
		//Check if the date is major or equal to delete the file
		isDateDelete = c.DateComparation(file.ModTime().Unix())
		if !isDateDelete {
			continue
		}
		//The path of the file
		filePath = c.OSFile + "/" + file.Name()
		//Remove the all files
		c.DeleteFile(filePath)
	}
	// clean the file slice
	c.files = nil
}

// The CleanerWithContext execute the cleaner with like a cronJob, thats mean that the function receive an hour and a specific minute to excute the Clean process
//
// Args
//   - ctx: A context parent to cancel the function
//   - hour: The hour that the clean process will be excecute it
//   - minutes: The minutes that the clean process will be excecute it
func (c *clean) CleanerWithContext(ctx context.Context, hour int, minutes int) {

	//Create a channel that will send a value on a channel every minute
	ticker := time.NewTicker(time.Minute)

	for {
		select {
		//If the ctx is done, the CleanerWithContext will stop the excecution
		case <-ctx.Done():
			log.Info("Cleaner with context is done")
			return
		// Channel to read the ticker
		case <-ticker.C:
			h, m, _ := time.Now().Clock()
			// Check if the time to execute the Cleaner match wit the current time
			if m == minutes && h == hour {
				c.Cleaner()
			}
		}
	}

}

// DeleteFile
//
// This method remove the file.
//
//	Args:
//		-filePath: the path of the file, the NON relative path
func (c *clean) DeleteFile(filePath string) {
	err := os.Remove(filePath)
	if err != nil {
		log.Error(err.Error())
		return
	} else {
		log.Info("Succefully deleted file: ", filePath)
	}
}

/*
	Define the setter and getter methods
*/

// Get the all files in the OSFile
func (c *clean) GetFiles() []fs.FileInfo {
	c.ReadFiles()
	return c.files
}

func (c clean) PrintFiles() {
	for _, file := range c.files {
		fmt.Println(file.Name())
	}
}

// Get the current time in UNIX format
func (c clean) GetCurrentUNIXTime() int64 {
	return time.Now().Unix()
}

// Set the CurrenteTime attribute to the current time in UNIX format
func (c *clean) SetUNIXTimeNow() {
	c.currentTime = time.Now().Unix()
}

// Get the file extension that the files extension will be removed
func (c clean) GetFileExtensionToClean() string {
	return c.fileExtension
}

// Set the file extension that the files extension will be removed
func (c *clean) SetFileExtensionToClean(fileExtension string) {
	c.fileExtension = fileExtension
}

// Get the days of the difference of the clean that will be removed
func (c clean) GetDaysDiffToClean() int {
	return c.daysDiffToClean
}

// Set the days of last modification of the files that will be removed
func (c *clean) SetDaysDiffToClean(daysDiff int) {
	c.daysDiffToClean = daysDiff
	c.secondsDiffToClean = int64(daysDiff) * secondsInADay
}
