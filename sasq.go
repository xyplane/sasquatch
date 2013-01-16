package main

import "os"
import "fmt"
//import "flag"
import "time"




func main() {
	fmt.Println("Hello Sasquatch!")
	
	var perr = parseArguments()
	if perr != nil {
		fmt.Println(perr)
		os.Exit(-1)
	}


	// Concurrent Section //

	var readCount = 0
	var readErrExit = false
	var readErrChan = make(chan error)
	var readTimeout = time.NewTimer(30*time.Second)

	for _, tf := range templateFlags {
		go tf.ReadTemplate(readErrChan)
		readCount++
	}

	for _, df := range dataFlags {
		go df.ReadData(readErrChan)
		readCount++
	}

	for readCount > 0 {
		select {
		case err := <-readErrChan:
			readCount--
			if err != nil {
				fmt.Println(err)
				readErrExit = true
			}
		case <-readTimeout.C:
			readCount = 0
			fmt.Println("Timeout!")
			readErrExit = true
		}
	}

	if readErrExit {
		os.Exit(-1)
	}

	// Sequential Section //

	for _, tf := range templateFlags {
		var err = tf.AddTemplate()
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}

	for _, df := range dataFlags {
		var err = df.AddData()
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}

	if tmpl == nil {
		fmt.Println("no templates defined")
		os.Exit(-1)
	}

	var err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	os.Exit(0)
}
