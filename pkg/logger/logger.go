package logger

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"runtime"

	"github.com/kyokomi/emoji"
	. "github.com/logrusorgru/aurora"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var z *zap.Logger = nil
var aurora Aurora = nil

func InitAurora() {
	if aurora == nil {
		aurora = NewAurora(viper.GetBool("logging_color"))
	}
}

func GetAurora() Aurora {
	return aurora
}

func ZapLogger() error {
	var err error
	if z == nil {
		// TODO: test permission for open logfile.
		cfg := zap.NewProductionConfig()
		cfg.OutputPaths = []string{}
		cfg.Level = level2AtomicLevel(viper.GetString("logging_level"))
		cfg.ErrorOutputPaths = []string{}
		if viper.GetBool("logging_json") {
			cfg.Encoding = "json"
		} else {
			cfg.Encoding = "console"
		}
		cfg.DisableCaller = true
		cfg.DisableStacktrace = true
		cfg.EncoderConfig.TimeKey = "time"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		z, err = cfg.Build()
		if err != nil {
			fmt.Fprint(os.Stderr, "Error on initialize file logger: "+err.Error()+"\n")
			return err
		}
	}

	return nil
}

func level2Number(level string) int {
	switch level {
	case "error":
		return 0
	case "warning":
		return 1
	case "info":
		return 2
	case "debug":
		return 3
	default:
		return 2
	}
}

func log2File(level, msg string) {
	switch level {
	case "error":
		z.Error(msg)
	case "warning":
		z.Warn(msg)
	case "info":
		z.Info(msg)
	default:
		z.Debug(msg)
	}
}

func level2AtomicLevel(level string) zap.AtomicLevel {
	switch level {
	case "error":
		return zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "warning":
		return zap.NewAtomicLevelAt(zap.WarnLevel)
	case "debug":
		return zap.NewAtomicLevelAt(zap.DebugLevel)
	default:
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	}
}

func Msg(level string, withoutColor, ln bool, msg ...interface{}) {
	var message string
	var confLevel, msgLevel int

	if viper.GetBool("debug") {
		confLevel = 3
	} else {
		confLevel = level2Number(viper.GetString("logging_level"))
	}
	msgLevel = level2Number(level)
	if msgLevel > confLevel {
		return
	}

	for _, m := range msg {
		message += " " + fmt.Sprintf("%v", m)
	}

	var levelMsg string

	if withoutColor || !viper.GetBool("logging_color") {
		levelMsg = message
	} else {
		switch level {
		case "warning":
			levelMsg = Yellow(":construction: warning" + message).BgBlack().String()
		case "debug":
			levelMsg = White(message).BgBlack().String()
		case "info":
			levelMsg = message
		case "error":
			levelMsg = Red(message).String()
		}
	}

	if viper.GetBool("emoji") {
		levelMsg = emoji.Sprint(levelMsg)
	} else {
		re := regexp.MustCompile(`[:][\w]+[:]`)
		levelMsg = re.ReplaceAllString(levelMsg, "")
	}

	if z != nil {
		log2File(level, message)
	}

	if ln {
		fmt.Println(levelMsg)
	} else {
		fmt.Print(levelMsg)
	}

}

func Warning(mess ...interface{}) {
	Msg("warning", false, true, mess...)
}

func Debug(mess ...interface{}) {
	pc, file, line, ok := runtime.Caller(1)
	if ok {
		mess = append([]interface{}{fmt.Sprintf("DEBUG (%s:#%d:%v)",
			path.Base(file), line, runtime.FuncForPC(pc).Name())}, mess...)
	}
	Msg("debug", false, true, mess...)
}

func DebugC(mess ...interface{}) {
	Msg("debug", true, true, mess...)
}

func Info(mess ...interface{}) {
	Msg("info", false, true, mess...)
}

func InfoC(mess ...interface{}) {
	Msg("info", true, true, mess...)
}

func Error(mess ...interface{}) {
	Msg("error", false, true, mess...)
}

func Fatal(mess ...interface{}) {
	Error(mess...)
	os.Exit(1)
}
