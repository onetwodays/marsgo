/*
Package ut_log是对seelog的封装
Xml configuring to be able to change logger parameters without recompilation
Changing configurations on the fly without app restart
Possibility to set different log configurations for different project files and functions
Adjustable message formatting
Simultaneous log output to multiple streams
Choosing logger priority strategy to minimize performance hit
Different output writers
Console writer
File writer
Buffered writer (Chunk writer)
Rolling log writer (Logging with rotation)
SMTP writer
Others... (See Wiki)
Log message wrappers (JSON, XML, etc.)
Global variables and functions for easy usage in standalone apps
Functions for flexible usage in libraries
Installation:go get -u github.com/cihub/seelog

Message context:
%Level - log level (Trace, Debug, Info, Warn, Error, Critical)
%Lev - short log level (Trc, Dbg, Inf, Wrn, Err, Crt)
%LEVEL - capitalized log level (TRACE, DEBUG, INFO, WARN, ERROR, CRITICAL)
%LEV - short capitalized log level (TRC, DBG, INF, WRN, ERR, CRT)
%l - super compact log level (t, d, i, w, e, c)
%Msg - message text (string)
%FullPath - full caller file path
%File - caller filename only
%RelFile - caller path relative to the application runtime directory
%Func - caller function name
%FuncShort - caller function name part after the last dot
%Line - line number where logger was called
Date and time
%Ns - time.Now().UnixNano()
%Date - shortcut for ‘2006-01-02’
%Time - shortcut for ‘15:04:05’
%Date(...) - date with format, specified in parentheses. Uses standard time.Format, so check http://golang.org/src/pkg/time/format.go for identifiers list. Use it like that: "%Date(2006-01-02)" (or any other format)
%UTCNs - time.Now().UTC().UnixNano()
%UTCDate - shortcut for ‘2006-01-02’ (UTC)
%UTCTime - shortcut for ‘15:04:05’ (UTC)
%UTCDate(...) - UTC date with format, specified in parentheses. Uses standard time.Format, so check http://golang.org/src/pkg/time/format.go for identifiers list. Use it like that: "%UTCDate(2006-01-02)" (or any other format)
Special symbols
%EscN - terminal ANSI CSI n [;k] m escape. Check Colored output for details
%n - newline
%t - tab
*/

package ut_log

import (
	"fmt"
	"github.com/cihub/seelog"
	"strings"
)

/*
"trace"
"debug"
"info"
"warn"
"error"
"critical"
*/
const xMLCONFIG string = `
<seelog type="asynctimer" asyncinterval="100" minlevel="debug" maxlevel="critical">
    <!-- 
	<exceptions>
		<exception funcpattern="*main.test*Something*" minlevel="info"/>
		<exception filepattern="*main.go" minlevel="error"/>
    </exceptions>
 	-->
    <outputs formatid="main">
        <console/>
        <rollingfile formatid="main" type="size" filename="./log/APPNAME_VERSION.log" maxsize="100000000" maxrolls="10" />
        <filter levels="warn,error,critical" formatid="important">
        	<rollingfile  type="date" filename="./log/APPNAME_VERSION.log" datepattern="20060102"  maxrolls="10" />   
        </filter>

    </outputs>
    <formats>
		<format id="json" format='{"time":%Ns,"lev":"%Lev","msg":"%Msg","path":"%RelFile","func":"%Func","line":"%Line"}'/>
        <format id="main"   format="[%Date(2006-01-02 15:04:05.999999)][%Lev] [%File] [%FuncShort] [%Line] %Msg%n"/>
        <format id="important" format="[%Date(2006-01-02 15:04:05.999999)][%Lev] [%FullPath] [%Func] [%Line] %Msg%n"/>
    </formats>
</seelog>

`
/*
func init() {
	loadConfigXmlString(UT_LOG_CONFIG)
}
*/


//Loading config
func loadConfigXmlString(xml string) error  {
	logger, err := seelog.LoggerFromConfigAsString(xml)
	if err != nil {
		fmt.Println("load seelog config_xml_string:", err)
		return err
	}else{
		seelog.ReplaceLogger(logger)
	}
	return nil
}

//Note: Flush() in the defer block of the main function,
//what guarantees that all the messages left in a log message queue will
//be normally processed independently
//of whether the application panics or not
func InitLog(appName,version string)  {
	loadConfigXmlString(strings.Replace(xMLCONFIG,"APPNAME_VERSION",appName+"_"+version,-1))
}

// Flush immediately processes all currently queued messages and all currently buffered messages.
// It is a blocking call which returns only after the queue is empty and all the buffers are empty.
//
// If Flush is called for a synchronous logger (type='sync'), it only flushes buffers (e.g. '<buffered>' receivers)
// , because there is no queue.
//
// Call this method when your app is going to shut down not to lose any log messages.
func Flush()  {
	seelog.Flush()
}

// Tracef formats message according to format specifier
// and writes to default logger with log level = Trace.
func Tracef(format string, params ...interface{}) {
	seelog.Tracef(format,params)
}

// Debugf formats message according to format specifier
// and writes to default logger with log level = Debug.
func Debugf(format string, params ...interface{}) {
	seelog.Debugf(format,params...)
}

// Infof formats message according to format specifier
// and writes to default logger with log level = Info.
func Infof(format string, params ...interface{}) {
	seelog.Infof(format,params...)
}

// Warnf formats message according to format specifier and writes to default logger with log level = Warn
func Warnf(format string, params ...interface{}) error {
	return seelog.Warnf(format,params...)
}

// Errorf formats message according to format specifier and writes to default logger with log level = Error
func Errorf(format string, params ...interface{}) error {
	return seelog.Errorf(format,params...)
}

// Criticalf formats message according to format specifier and writes to default logger with log level = Critical
func Criticalf(format string, params ...interface{}) error {
	return seelog.Criticalf(format,params...)
}

// Trace formats message using the default formats for its operands and writes to default logger with log level = Trace
func Trace(v ...interface{}) {
	seelog.Trace(v)
}

// Debug formats message using the default formats for its operands and writes to default logger with log level = Debug
func Debug(v ...interface{}) {
	seelog.Debug(v)
}

// Info formats message using the default formats for its operands and writes to default logger with log level = Info
func Info(v ...interface{}) {
	seelog.Info(v)
}

// Warn formats message using the default formats for its operands and writes to default logger with log level = Warn
func Warn(v ...interface{}) error {
	return seelog.Warn(v)
}

// Error formats message using the default formats for its operands and writes to default logger with log level = Error
func Error(v ...interface{}) error {
	return seelog.Error(v)
}

// Critical formats message using the default formats for its operands and writes to default logger with log level = Critical
func Critical(v ...interface{}) error {
	return seelog.Critical(v)
}








/*

File writer
Function: Writes received messages to a file

Element name: file

Allowed attributes:

formatid - format that will be used by this receiver
path - path to the log file
Example:

<seelog>
    <outputs>
        <file path="log.log"/>
    </outputs>
</seelog>
Important:

Do not use any special symbols that are not allowed in filenames.
Do not use the same file in multi-process environment to avoid log inconsistency.
Console writer
Function: writes received messages to std out

Element name: Console

Allowed attributes:

formatid - format that will be used by this receiver
Example:

<seelog>
    <outputs>
        <console/>
    </outputs>
</seelog>
Rolling file writer (or "Rotation file writer")
Function: Writes received messages to a file, until date changes or file exceeds a specified limit. After that the current log file is renamed and writer starts to log into a new file. If you roll by size, you can set a limit for such renamed files count, if you want, and then the rolling writer would delete older ones when the files count exceed the specified limit.

Element name: rollingfile

Allowed attributes:

formatid - format that will be used by this receiver
filename - path to the log file. On creation, this path is split into folder path and actual file name. The latter will be used as a common part for all files, that act in rolling:
In case of 'date' rolling, the file names will be formed this way: filename.DATE_PATTERN. For example (if fullname flag is not set, see below): app.log, app.log.11.Aug.14, app.log.11.Aug.15, ...

In case of 'size' rolling, the file names will be formed this way: "filename.#". For example: app.log, app.log.1, app.log.2, app.log.3, ...

type - rotation type: "date" or "size"
namemode - naming mode of rolled files. Possible values: "postfix", "prefix". If set to "postfix", rolled files look like 'file.log.1', 'file.log.02.01.2006'. If set to 'prefix': '1.file.log', '02.01.2006.file.log'. Default is 'postfix'.
maxrolls - Maximal count of renamed files. When this limit is exceeded, older rolls are removed. Should be >= zero. NOTE: Before version 2.3 this attribute was not applicable to the 'date' rolling writer.
archivetype - an attribute used to specify the type of the archive where old rolls are stored instead of removal. Possible values: "none", "zip", "gzip". If set to "none", no archivation is performed (old rolls are just deleted).
archiveexploded - an attribute used to specify whether logs should be exploded or grouped inside the same archive file.
archivepath - an attribute used when archivetype is not set to "none". Specifies the path of the archive where old rolls are stored.
Attributes, allowed with 'size' type:

maxsize - This is the size limit (in bytes) exceeding which results in a roll.
Attributes, allowed with 'date' type:

datepattern - This is the pattern that would be used in 'time.LocalTime().Format' to form a file name. The "date" (actually, this means both date & time) rolling occurs when 'time.LocalTime().Format(rollFileWriter.datePattern)' returns something different than the current file name. This means that you create daily rotation using a format with a day identifier like "02.01.2006". Or you can create an hourly rotation using a format with an hour identifier like "02.01.2006.15".
fullname - a boolean attribute that affects the current file name. If set to "true", then current file name will use the same naming convention as the roll names. For example, log files list will look like app.log.11.Aug.13, app.log.11.Aug.14, app.log.11.Aug.15, ... instead of app.log, app.log.11.Aug.14, app.log.11.Aug.15, .... Default is "false".
Important:

Do not use any special symbols that are not allowed in filenames.
Do not use the same file in multi-process environment to avoid log inconsistency.
Example:

<seelog>
	<outputs>
		<rollingfile type="size" filename="logs/roll.log" maxsize="1000" maxrolls="5" />
	</outputs>
</seelog>
<seelog>
    <outputs>
        <rollingfile type="date" filename="logs/roll.log" datepattern="02.01.2006" maxrolls="7" />
    </outputs>
</seelog>
Buffered writer
Function: Acts like a buffer wrapping other writer. Buffered writer stores data in memory and flushes it every flush period or when buffer is full.

Element name: buffered

Allowed attributes:

formatid - format that will be used by this receiver
size - buffer size (in bytes)
flushperiod - interval between buffer flushes (in milliseconds)
Example:

<seelog>
    <outputs>
        <buffered size="10000" flushperiod="1000">
            <file path="bufFile.log"/>
        </buffered>
    </outputs>
</seelog>
<seelog>
    <outputs>
        <buffered size="10000" flushperiod="1000" formatid="someFormat">
            <rollingfile type="date" filename="logs/roll.log" datepattern="02.01.2006" />
        </buffered>
    </outputs>
    <formats>
        ...
    </formats>
</seelog>
NOTE: This writer accumulates data written using a specific format and then flushes it into the inner writer. So inner writers couldn't have its own format: set 'formatid' only for buffered element, like in the last example.

SMTP writer
Function: Sends emails to the specified recipients using a password-protected (but generally unsecured) email account at a given post server.

Element name: smtp

Allowed attributes:

senderaddress - email address of the sender
sendername - name of the sender
hostname - host name of a post server (typically - mail.XXX.YYY)
hostport - TCP port of the post server (typically - 587)
username - user name used to login to the post server
password - password to the post server
subject - subject of the email
Subelement: recipient
Allowed attributes:

address - email address of the recipient (who receives messages from notifiers)
Subelement: header
Allowed attributes:

name - name of the header
value - value of the header
Subelement: cacertdirpath
Allowed attributes:

path - path to a directory with PEM certificate files.
Example:

<seelog>
  <outputs>
   <smtp senderaddress="noreply-notification-service@none.org" sendername="Automatic notification service" hostname="mail.none.org" hostport="587" username="nns" password="123">
    <recipient address="john-smith@none.com"/>
    <recipient address="hans-meier@none.com"/>
   </smtp>
  </outputs>
 </seelog>
<seelog>
  <outputs>
   <filter levels="error,critical">
    <smtp senderaddress="nns@none.org" sendername="ANS" hostname="mail.none.org" hostport="587" username="nns" password="123" subject="test">
     <cacertdirpath path="cacdp1"/>
     <recipient address="hans-meier@none.com"/>
     <header name="Priority" value="Urgent" />
     <header name="Importance" value="high" />
     <header name="Sensitivity" value="Company-Confidential" />
     <header name="Auto-Submitted" value="auto-generated" />
    </smtp>
   </filter>
  </outputs>
 </seelog>
NOTE 1: Look at the first above example to get an idea of how to use SMTP writer. Keep in mind that email noreply-notification-service@none.org must not be considered as security reliable. It's potentially exposed to hacking attacks due to explicit password publication in configs and we strongly recommend you not to use personal or corporate post accounts for SMTP writer. The best practice suggests dedicating an isolated post account especially for email notification service.

NOTE 2: The second example demonstrates the most sensible way to use this writer - notifications from the application on (rare) extraordinary situations. Technically, you can set other filter levels, yet get ready for being flooded with diagnostic emails.

Conn writer
Function: Writes received messages to a network connection.

Element name: conn

Allowed attributes:

formatid - format that will be used by this receiver
net - network ( "tcp", "udp", "tcp4", "udp4", ... )
addr - network address ( ":1000", "127.0.0.1:8888", ... )
reconnectonmsg - If true connection will be opened on each write otherwise on first write. Default - false.
usetls - if true, TLS will be used.
insecureskipverify - sets the InsecureSkipVerify flag of the tls.Config. Use it if useTLS is set.
Example:

<seelog type="sync">
	<outputs>
		<conn net="tcp" addr=":8888" />
	</outputs>
</seelog>
<outputs>
  <conn formatid="syslog" net="tcp4" addr="server.address:5514" tls="true" insecureskipverify="true" />
</outputs>
<formats>
  <format id="syslog" format="%CustomSyslogHeader(20) %Msg%n"/>
</formats>
For %CustomSyslogHeader example code check the custom formatters section.
 */