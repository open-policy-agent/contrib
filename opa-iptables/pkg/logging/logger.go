package logging

import (
	"os"
	"github.com/sirupsen/logrus"
)

// initialize logger
var log = logrus.New()

func init(){
	// setting logger output to stdout
	log.SetOutput(os.Stdout)
}

// Config defines format and level of logging for logger
type Config struct {
	// Format can be text | json | json-pretty. Default format is text.
	Format string
	// Level can be info | debug | error. Default level is info.
	Level string
}

// SetupLogging setting up logger using given configuration
func SetupLogging(config Config) {

	switch config.Format {
	case "text":
		log.SetFormatter(&logrus.TextFormatter{TimestampFormat:"2006-01-02 15:04:05",DisableSorting:true,FullTimestamp:true,DisableLevelTruncation:true})
	case "json-pretty":
		log.SetFormatter(&logrus.JSONFormatter{PrettyPrint:true,TimestampFormat:"2006-01-02 15:04:05"})
	case "json":
		fallthrough
	default:
		log.SetFormatter(&logrus.JSONFormatter{TimestampFormat:"2006-01-02 15:04:05"})
	}

	level := logrus.InfoLevel
	if config.Level != "" {
		var err error
		level,err = logrus.ParseLevel(config.Level)
		if err != nil {
			logrus.Fatalf("Unable to parse log level: %v", err)
		}
	}
	log.SetLevel(level)
}

// GetLogger returns an instance of logger
func GetLogger() *logrus.Logger {
	return log
}