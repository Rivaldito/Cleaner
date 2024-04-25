# Cleaner 🧹🧼🪣🧽

Cleaner is a package to allow you delete files in a especif directory.

I'll wrote this package solving a common issue in many of my programs. I'll always have a log dir where I save the logs of every day, but with the pass of the time this files keep inside of the directory like garbage 🗑️.

So I got down to work, and decide to write this little package

## Case of use

You can use the `Cleaner` package in two ways:

### #1. Execute the function when it's need it ⏰

Use the cleaner like a normal function.

Clean on demand, I mean, clean only whe the function is called 🥸.

```golang
package main

import (
	"github.com/Rivaldito/Cleaner/cleaner"
)

func main() {

	const (
		osFile   string = "/tmp"
		daysDiff int    = 2
	)

	c := cleaner.NewCleaner(
		osFile,
		cleaner.LOG,
		daysDiff)

	//Clean the garbage 🌪️
	c.Cleaner()

	/*
		Do some other stuff, like mount a HTTP server
	*/

	//Call cancel if you want to stop the cleaner
}
```

### #2. Execute the cleaner like a cronjob ⏱️

In this way to excute we can run the cleaner in background mode.

``` Golang
package main

import (
	"context"

	"github.com/Rivaldito/Cleaner/cleaner"
)

func main() {

	const (
        //The directory we want to clean 🫧
		osFile   string = "/tmp"
        //The days of difference we want to remove the files
		daysDiff int    = 2
	)

	ctx, cancel := context.WithCancel(context.Background())

	c := cleaner.NewCleaner(
		osFile,
        //The extension of the files we want to clean
		cleaner.LOG,
		daysDiff)


    // Time to execute the function

    const (
        hour    int = 14
        minutes int = 22
    )

    // Execute the cleaner in background mode
	go c.CleanerWithContext(ctx, hour, minutes)

	/*
		Do some other stuff, like mount a HTTP server
	*/

	//Call cancel if you want to stop the cleaner
	cancel()

}
```