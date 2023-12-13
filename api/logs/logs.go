package logs

import (
	"log"
	"os"
	"io"
)

var errorLogFile *os.File

func InitLogger(){
	errorLogFile, _ = os.OpenFile("./logs/error.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)

	iw := io.MultiWriter(os.Stdout, errorLogFile)

	log.SetOutput(iw)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}