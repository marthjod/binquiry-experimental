package logging

import (
	log "github.com/Sirupsen/logrus"
)

// Use timestamped format. If user-supplied log level cannot be set, exit.
func MustSetLoglevel(lvl string) {
	var (
		level log.Level
		err   error
	)

	if level, err = log.ParseLevel(lvl); err != nil {
		log.Fatal(err)
	}
	log.SetLevel(level)
	formatter := &log.TextFormatter{
		FullTimestamp: true,
	}
	log.SetFormatter(formatter)
}
