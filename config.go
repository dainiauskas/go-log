package log

import (
	"os/user"
	"strings"
	"sync"
)

// logger configuration
type config struct {
	logPath    string
	pathPrefix string
	logflags   uint32
	maxdays    int // limit log files by days, zero unlimited
	purgeLock  sync.Mutex
	enabled    bool
}

var gConf = config{
	logPath:  "./log/",
	logflags: flagLogFilenameLineNum | flagLogThrough,
	maxdays:  30,
	enabled:  true,
}

func (conf *config) setFlags(flag uint32, on bool) {
	if on {
		conf.logflags = conf.logflags | flag
	} else {
		conf.logflags = conf.logflags & ^flag
	}
}

func (conf *config) logTrace() bool {
	return (conf.logflags & flagLogTrace) != 0
}

func (conf *config) logDebug() bool {
	return (conf.logflags & flagLogDebug) != 0
}

func (conf *config) logThrough() bool {
	return (conf.logflags & flagLogThrough) != 0
}

func (conf *config) logFuncName() bool {
	return (conf.logflags & flagLogFuncName) != 0
}

func (conf *config) logFilenameLineNum() bool {
	return (conf.logflags & flagLogFilenameLineNum) != 0
}

func (conf *config) logToConsole() bool {
	return (conf.logflags & flagLogToConsole) != 0
}

func (conf *config) isEnabled() bool {
	return conf.enabled
}

func (conf *config) setFilenamePrefix(filenamePrefix, symlinkPrefix string) {
	username := "Unknown"
	curUser, err := user.Current()
	if err == nil {
		tmpUsername := strings.Split(curUser.Username, "\\") // for compatible with Windows
		username = tmpUsername[len(tmpUsername)-1]
	}

	conf.pathPrefix = conf.logPath
	if len(filenamePrefix) > 0 {
		filenamePrefix = strings.Replace(filenamePrefix, "%P", gProgname, -1)
		filenamePrefix = strings.Replace(filenamePrefix, "%H", infoHostName, -1)
		filenamePrefix = strings.Replace(filenamePrefix, "%U", username, -1)
		conf.pathPrefix = conf.pathPrefix + filenamePrefix + "."
	}

	if len(symlinkPrefix) > 0 {
		symlinkPrefix = strings.Replace(symlinkPrefix, "%P", gProgname, -1)
		symlinkPrefix = strings.Replace(symlinkPrefix, "%H", infoHostName, -1)
		symlinkPrefix = strings.Replace(symlinkPrefix, "%U", username, -1)
		symlinkPrefix += "."
	}

	isSymlink = map[string]bool{}
	for i := 0; i != logLevelMax; i++ {
		gLoggers[i].level = i
		gSymlinks[i] = symlinkPrefix + gLogLevelNames[i]
		isSymlink[gSymlinks[i]] = true
		gFullSymlinks[i] = conf.logPath + gSymlinks[i]
	}
}

// SetMaxDays - change maxdays parameter
func SetMaxDays(days int) {
	gConf.maxdays = days
}
