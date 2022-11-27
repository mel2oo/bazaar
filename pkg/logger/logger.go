package logger

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	Level      string `mapstructure:"level"`
	MaxSize    int    `mapstructure:"max-size"`
	MaxAge     int    `mapstructure:"max-age"`
	MaxBackups int    `mapstructure:"max-backup"`
	Output     string `mapstructure:"output"`
	Stdout     bool   `mapstructure:"stdout"`
}

func Init(servername string, c Config, dbg bool, conditions []string) error {

	// setup log formatter
	logrus.SetFormatter(&nested.Formatter{
		TimestampFormat: time.RFC3339,
		HideKeys:        true,
	})

	level, err := logrus.ParseLevel(c.Level)
	if err != nil {
		return err
	}
	logrus.SetLevel(level)

	// disable writing to stdout
	if !c.Stdout {
		logrus.SetOutput(ioutil.Discard)
	}

	// initialize log rotate hook
	rhook, err := NewHook(servername, c, level, dbg, conditions)
	if err != nil {
		return err
	}
	logrus.AddHook(rhook)

	return err
}

// File represents the rotate file hook.
type File struct {
	app       string
	level     logrus.Level
	formatter logrus.Formatter
	w         io.Writer

	dbg        bool
	conditions []string
}

// NewHook builds a new rotate file hook.
func NewHook(app string, c Config, level logrus.Level, dbg bool, conditions []string) (logrus.Hook, error) {
	name := os.Getenv("LOG_NAME")
	if len(name) == 0 {
		name = filepath.Base(os.Args[0])
	}

	return &File{
		app:   app,
		level: level,
		formatter: &nested.Formatter{
			TimestampFormat: time.RFC3339,
			NoColors:        true,
			HideKeys:        true,
		},
		w: &lumberjack.Logger{
			Filename:   filepath.Join(c.Output, name+".log"),
			MaxSize:    c.MaxSize,
			MaxAge:     c.MaxAge,
			MaxBackups: c.MaxBackups,
		},
		dbg:        dbg,
		conditions: conditions,
	}, nil
}

// Levels determines log levels that for which the logs are written.
func (hook *File) Levels() []logrus.Level {
	return logrus.AllLevels[:hook.level+1]
}

// Fire is called by logrus when it is about to write the log entry.
func (hook *File) Fire(entry *logrus.Entry) (err error) {
	entry.Data["app"] = hook.app
	b, err := hook.formatter.Format(entry)
	if err != nil {
		return err
	}

	_, err = hook.w.Write(b)
	if err != nil {
		return err
	}

	if hook.dbg {
		for _, con := range hook.conditions {
			if strings.Contains(string(b), con) {
				logrus.Errorf("break server")
				os.Exit(5)
			}
		}
	}

	return err
}
