// Package commandProcessor is the first logical gate our tool encounters
package commandProcessor

// InitProcess Start the machine, handle which services to set up.
//func InitProcess() error {
//	// Create a waitgroup, so we can run all services at once.
//	var wg sync.WaitGroup
//	errChan := make(chan error, 1)
//	outChan := make(chan string, 1)
//
//	// if conf.ServicesConfig.Ssh {
//	// TODO
//	// }
//	// if conf.ServicesConfig.Teamviewer {
//	// TODO
//	// }
//	// wait for all the goroutines to complete
//	wg.Wait()
//
//	// Open the log file for writing
//	logFile := "app.log"
//	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
//	if err != nil {
//		log.Fatal(err)
//	}
//	logger := log.New(file, "", log.LstdFlags)
//
//	// check if there were any errors
//	select {
//	case err := <-errChan:
//		logger.Fatalf("errorChan: %v", err)
//		close(errChan)
//	case out := <-outChan:
//		logger.Println(out)
//		log.Printf("Much success! Check the %s file for details!", logFile)
//	default:
//		log.Println("Much success!")
//	}
//	return nil
//}
