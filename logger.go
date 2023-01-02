/*
Package logger is a logging facility which provides functions Trace, Info, Warn, Error, Panic and Abort to
write logs with different severity levels. Logs with different severity levels are written to different logfiles.
Sorry for my poor English, I've tried my best.
Features:
	1. Auto rotation: It'll create a new logfile whenever day changes or size of the current logfile exceeds the configured size limit.
	2. Auto purging: It'll delete some oldest logfiles whenever the number of logfiles exceeds the configured limit.
	3. Log-through: Logs with higher severity level will be written to all the logfiles with lower severity level.
	4. Logs are not buffered, they are written to logfiles immediately with os.(*File).Write().
	5. Symlinks `PROG_NAME`.`USER_NAME`.`SEVERITY_LEVEL` will always link to the most current logfiles.
	6. Goroutine-safe.
Basic example:
	// logger.Init must be called first to setup logger
	logger.Init("./log", // specify the directory to save the logfiles
			400, // maximum logfiles allowed under the specified log directory
			20, // number of logfiles to delete when number of logfiles exceeds the configured limit
			100, // maximum size of a logfile in MB
			false) // whether logs with Trace level are written down
	logger.Info("Failed to find player! uid=%d plid=%d cmd=%s xxx=%d", 1234, 678942, "getplayer", 102020101)
	logger.Warn("Failed to parse protocol! uid=%d plid=%d cmd=%s", 1234, 678942, "getplayer")
*/

package log

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// log level
const (
	logLevelTrace = iota
	logLevelInfo
	logLevelWarn
	logLevelError
	logLevelUpdate
	logLevelPanic
	logLevelAbort
	logLevelQuery
	logLevelDebug

	logLevelMax
)

// log flags
const (
	flagLogTrace = 1 << iota
	flagLogThrough
	flagLogFuncName
	flagLogFilenameLineNum
	flagLogToConsole
	flagLogDebug
)

// log info
var (
	infoUserName string
	infoHostName string
)

// const strings
const (
	// Default filename prefix for logfiles
	DefFilenamePrefix = "%P.%H.%U"
	// Default filename prefix for symlinks to logfiles
	DefSymlinkPrefix = "%P.%U"

	logLevelChar = "TIWEUPAQD"
)

// Init must be called first, otherwise this logger will not function properly!
// It returns nil if all goes well, otherwise it returns the corresponding error.
//
//	maxdays: Maximum days to keep logs.
//	logTrace: If set to false, `logger.Trace("xxxx")` will be mute.
func Init(logpath string, maxdays int, logTrace bool) error {
	err := os.MkdirAll(logpath, 0755)
	if err != nil {
		return err
	}

	infoHostName, err = os.Hostname()
	if err != nil {
		return err
	}

	gConf.logPath = logpath + "/"
	gConf.maxdays = maxdays
	gConf.setFlags(flagLogTrace, logTrace)

	SetFilenamePrefix(DefFilenamePrefix, DefSymlinkPrefix)

	return nil
}

// SetLogTrace sets to write trace log file
func SetLogTrace(on bool) {
	gConf.setFlags(flagLogTrace, on)
}

// SetLogDebug sets to write trace log file
func SetLogDebug(on bool) {
	gConf.setFlags(flagLogDebug, on)
}

// SetLogThrough sets whether to write log to all the logfiles with less severe log level.
// By default, logthrough is turn on. You can turn it off for better performance.
func SetLogThrough(on bool) {
	gConf.setFlags(flagLogThrough, on)
}

// SetLogFunctionName sets whether to log down the function name where the log takes place.
// By default, function name is not logged down for better performance.
func SetLogFunctionName(on bool) {
	gConf.setFlags(flagLogFuncName, on)
}

// SetLogFilenameLineNum sets whether to log down the filename and line number where the log takes place.
// By default, filename and line number are logged down. You can turn it off for better performance.
func SetLogFilenameLineNum(on bool) {
	gConf.setFlags(flagLogFilenameLineNum, on)
}

// SetLogToConsole sets whether to output logs to the console.
// By default, logs are not output to the console.
func SetLogToConsole(on bool) {
	gConf.setFlags(flagLogToConsole, on)
}

// SetLogUserName sets user name to write to log.
// By default, empty
func SetLogUserName(name string) {
	infoUserName = name
}

// SetLogDisable logging
// By default, logs are enabled
func SetLogDisable() {
	gConf.enabled = false
}

// SetLogEnable set logging enabled
func SetLogEnable() {
	gConf.enabled = true
}

// SetFilenamePrefix sets filename prefix for the logfiles and symlinks of the logfiles.
//
// Filename format for logfiles is `PREFIX`.`SEVERITY_LEVEL`.`DATE_TIME`.log
//
// Filename format for symlinks is `PREFIX`.`SEVERITY_LEVEL`
//
// 3 kinds of placeholders can be used in the prefix: %P, %H and %U.
//
// %P means program name, %H means hostname, %U means username.
//
// The default prefix for a log filename is logger.DefFilenamePrefix ("%P.%H.%U").
// The default prefix for a symlink is logger.DefSymlinkPrefix ("%P.%U").
func SetFilenamePrefix(logfilenamePrefix, symlinkPrefix string) {
	gConf.setFilenamePrefix(logfilenamePrefix, symlinkPrefix)
}

// Trace logs down a log with trace level.
// If parameter logTrace of logger.Init() is set to be false, no trace logs will be logged down.
func Trace(format string, args ...interface{}) {
	if gConf.logTrace() {
		log(logLevelTrace, format, args)
	}
}

// Console output only to console
func Console(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

// Info logs down a log with info level.
func Info(format string, args ...interface{}) {
	log(logLevelInfo, format, args)
}

// Update logs down a log with update level.
func Update(format string, args ...interface{}) {
	log(logLevelUpdate, format, args)
}

// Warn logs down a log with warning level.
func Warn(format string, args ...interface{}) {
	log(logLevelWarn, format, args)
}

// Error logs down a log with error level.
func Error(format string, args ...interface{}) {
	log(logLevelError, format, args)
}

// Panic logs down a log with panic level and then panic("panic log") is called.
func Panic(format string, args ...interface{}) {
	log(logLevelPanic, format, args)
}

// Abort logs down a log with abort level and then os.Exit(-1) is called.
func Abort(format string, args ...interface{}) {
	log(logLevelAbort, format, args)
}

// Query logs down a log with query level
func Query(format string, args ...interface{}) {
	log(logLevelQuery, format, args)
}

// Debug logs down a log with debug level
func Debug(format string, args ...interface{}) {
	if gConf.logDebug() {
		log(logLevelDebug, format, args)
	}
}

type Logger struct{}

func (l Logger) Println(v ...interface{}) {
	Trace("", v...)
}

func (l Logger) Printf(format string, v ...interface{}) {
	Trace(format, v...)
}

// Gorm structure used for Gorm SQL query logging
type Gorm struct{}

// Print function used in Gorm to output
func (l Gorm) Print(args ...interface{}) {
	var messages []interface{}

	switch args[0] {
	case "sql":
		messages = []interface{}{
			args[3],
			args[4],
			args[2],
			args[5],
		}

		Query("Query=[%v], Values=%v Duration=[%v], Rows=[%v]", messages...)
	case "log":
		messages = []interface{}{
			args[1],
			args[2],
		}

		Query("Source=[%v], Error=%v", messages...)
	}
}

// logger
type logger struct {
	file   *os.File
	level  int
	day    int
	size   int64
	purged time.Time
	lock   sync.Mutex
}

func (l *logger) log(t time.Time, data []byte) {
	y, m, d := t.Date()

	l.lock.Lock()
	defer l.lock.Unlock()

	// Purge once in 24 hours
	fmt.Println("Purged", l.purged, 24*time.Hour, -time.Until(l.purged) > (24*time.Hour), l.purged.IsZero())

	if l.purged.IsZero() || -time.Until(l.purged) > (24*time.Hour) {
		gConf.purgeLock.Lock()
		hasLocked := true

		defer func() {
			if hasLocked {
				gConf.purgeLock.Unlock()
			}
		}()

		filepath.Walk(gConf.logPath, func(path string, info os.FileInfo, e error) error {
			if e != nil {
				return e
			}

			if !info.Mode().IsRegular() {
				return nil
			}

			if filepath.Ext(info.Name()) != ".log" {
				return nil
			}

			fmt.Println(-time.Until(info.ModTime()) > (time.Hour * 24 * time.Duration(gConf.maxdays)))
			if -time.Until(info.ModTime()) > (time.Hour * 24 * time.Duration(gConf.maxdays)) {
				e = os.Remove(path)
				if e != nil {
					l.errlog(t, nil, e)
				}
			}

			return e
		})

		l.purged = time.Now()
		gConf.purgeLock.Unlock()
		hasLocked = false
	}

	filename := fmt.Sprintf("%s%s_%d%02d%02d.log", gConf.pathPrefix, gLogLevelNames[l.level], y, m, d)
	newfile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		l.errlog(t, data, err)
		return
	}

	l.file.Close()
	l.file = newfile
	l.day = d
	l.size = 0

	err = os.RemoveAll(gFullSymlinks[l.level])
	if err != nil {
		l.errlog(t, nil, err)
	}
	err = os.Symlink(path.Base(filename), gFullSymlinks[l.level])
	if err != nil {
		l.errlog(t, nil, err)
	}

	n, _ := l.file.Write(data)
	l.size += int64(n)
}

// (l *logger).errlog() should only be used within (l *logger).log()
func (l *logger) errlog(t time.Time, originLog []byte, err error) {
	buf := gBufPool.getBuffer()

	genLogPrefix(buf, l.level, 2, t)
	buf.WriteString(err.Error())
	buf.WriteByte('\n')
	if l.file != nil {
		l.file.Write(buf.Bytes())
		if len(originLog) > 0 {
			l.file.Write(originLog)
		}
	} else {
		fmt.Fprint(os.Stderr, buf.String())
		if len(originLog) > 0 {
			fmt.Fprint(os.Stderr, string(originLog))
		}
	}

	gBufPool.putBuffer(buf)
}

// init is called after all the variable declarations in the package have evaluated their initializers,
// and those are evaluated only after all the imported packages have been initialized.
// Besides initializations that cannot be expressed as declarations, a common use of init functions is to verify
// or repair correctness of the program state before real execution begins.
func init() {
	tmpProgname := strings.Split(gProgname, "\\") // for compatible with `go run` under Windows
	gProgname = tmpProgname[len(tmpProgname)-1]

	gConf.setFilenamePrefix(DefFilenamePrefix, DefSymlinkPrefix)
}

func genLogPrefix(buf *buffer, logLevel, skip int, t time.Time) {
	h, m, s := t.Clock()

	// time
	buf.tmp[0] = logLevelChar[logLevel]
	buf.twoDigits(1, h)
	buf.tmp[3] = ':'
	buf.twoDigits(4, m)
	buf.tmp[6] = ':'
	buf.twoDigits(7, s)
	buf.Write(buf.tmp[:9])

	var pc uintptr
	var ok bool
	if gConf.logFilenameLineNum() {
		var file string
		var line int
		pc, file, line, ok = runtime.Caller(skip)
		if ok {
			buf.WriteByte(' ')
			buf.WriteString(path.Base(file))
			buf.tmp[0] = ':'
			n := buf.someDigits(1, line)
			buf.Write(buf.tmp[:n+1])
		}
	}
	if gConf.logFuncName() {
		if !ok {
			pc, _, _, ok = runtime.Caller(skip)
		}
		if ok {
			buf.WriteByte(' ')
			buf.WriteString(runtime.FuncForPC(pc).Name())
		}
	}
	if infoHostName != "" {
		buf.WriteByte(' ')
		buf.WriteString(infoHostName)
	}
	if infoUserName != "" {
		buf.WriteByte(' ')
		buf.WriteString(infoUserName)
	}

	buf.WriteString("] ")
}

func log(logLevel int, format string, args []interface{}) {
	if !gConf.isEnabled() {
		fmt.Println("Logger disabled")
		return
	}

	buf := gBufPool.getBuffer()

	t := time.Now()
	genLogPrefix(buf, logLevel, 3, t)
	fmt.Fprintf(buf, format, args...)
	buf.WriteByte('\n')
	output := buf.Bytes()
	if gConf.logThrough() {
		for i := logLevel; i != logLevelTrace; i-- {
			gLoggers[i].log(t, output)
		}
		if gConf.logTrace() {
			gLoggers[logLevelTrace].log(t, output)
		}
	} else {
		gLoggers[logLevel].log(t, output)
	}
	if gConf.logToConsole() {
		fmt.Print(string(output))
	}

	gBufPool.putBuffer(buf)
}

var gProgname = path.Base(os.Args[0])

var gLogLevelNames = [logLevelMax]string{
	"trace", "info", "warn", "error", "update", "panic", "abort", "query", "debug",
}

var gSymlinks [logLevelMax]string
var isSymlink map[string]bool
var gFullSymlinks [logLevelMax]string
var gBufPool bufferPool
var gLoggers [logLevelMax]logger
