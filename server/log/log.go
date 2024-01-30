package log

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/petermattis/goid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log = &zap.SugaredLogger{}

type ZapConf struct {
	FileName      string
	ErrorFileName string
	MaxSize       int
	MaxAge        int
	MaxBackups    int
	LocalTime     bool
	Compress      bool
}

// 构造一个zaplogger
func NewZapLogger(conf ZapConf) *zap.SugaredLogger {
	var coreArr []zapcore.Core

	// 自定义终端时间输出格式
	customTermTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(termTimeFormat))
	}

	// 所有等级均记录，具体的等级操作在root.go中进行分辨
	allPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev <= zap.FatalLevel
	})

	// 文件输出样式
	encoderConfig := zap.NewProductionEncoderConfig() //NewJSONEncoder()输出json格式，NewConsoleEncoder()输出普通文本格式
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeTime = customTermTimeEncoder //指定时间格式
	encoderConfig.ConsoleSeparator = " "
	fileEncoder := NewConsoleEncoder(encoderConfig)

	logger := &lumberjack.Logger{
		Filename:   conf.FileName,
		MaxSize:    conf.MaxSize,
		MaxAge:     conf.MaxAge,
		MaxBackups: conf.MaxBackups,
		LocalTime:  conf.LocalTime,
		Compress:   conf.Compress,
	}
	fileWriteSyncer := zapcore.AddSync(logger)
	termEncoder := NewConsoleEncoder(encoderConfig)
	fileCore := zapcore.NewCore(fileEncoder, zapcore.NewMultiWriteSyncer(fileWriteSyncer), allPriority)
	termCore := zapcore.NewCore(termEncoder, zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), allPriority)

	// 终端输出样式
	encoderConfig.EncodeTime = customTermTimeEncoder

	coreArr = append(coreArr, fileCore)
	coreArr = append(coreArr, termCore)
	log := zap.New(zapcore.NewTee(coreArr...))

	return log.Sugar()
}

// len(ctx) == 1时 输入是{msg="...%s",ctx="string"}也就是使用printf格式化
// 其他时候使用 K-V 的结构格式化
func Trace(msg string, ctx ...interface{}) {
	if LogLevel >= int(LvlTrace) {
		msg = getCaller(msg)
		if len(ctx) == 1 {
			log.Debugf(msg, ctx[0])
		}
		log.Debugw(msg, ctx...)
	}
}

func Debug(msg string, ctx ...interface{}) {
	if LogLevel >= int(LvlDebug) {
		msg = getCaller(msg)
		if len(ctx) == 1 {
			log.Debugf(msg, ctx[0])
		}
		log.Debugw(msg, ctx...)
	}
}

func Info(msg string, ctx ...interface{}) {
	if LogLevel >= int(LvlInfo) {
		msg = " " + getCaller(msg)
		if len(ctx) == 1 {
			log.Infof(msg, ctx[0])
		}
		log.Infow(msg, ctx...)
	}
}

func Warn(msg string, ctx ...interface{}) {
	if LogLevel >= int(LvlWarn) {
		msg = " " + getCaller(msg)
		if len(ctx) == 1 {
			log.Warnf(msg, ctx[0])
		}
		log.Warnw(msg, ctx...)
	}
}

func Error(msg string, ctx ...interface{}) {
	if LogLevel >= int(LvlError) {
		msg = getCaller(msg)
		if len(ctx) == 1 {
			log.Warnf(msg, ctx[0])
		}
		log.Errorw(msg, ctx...)
	}
}

func Crit(msg string, ctx ...interface{}) {
	msg = " " + getCaller(msg)
	if len(ctx) == 1 {
		log.Fatalf(msg, ctx[0])
	}
	log.Fatalw(msg, ctx...)
}

func getCaller(msg string) string {
	pc := make([]uintptr, skipLevel)
	n := runtime.Callers(skipLevel, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	if frame.File != "" {
		path := frame.File
		index1 := strings.LastIndex(path, "/") + 1  // 获取最后一个 "/" 的位置
		index2 := strings.Index(path[index1:], ".") // 获取第一个 "." 在子串中的位置
		if index2 < 0 {
			// 没有找到 "."，返回整个字符串
			index2 = len(path) - index1
		} else {
			// 找到了 "."，返回子串
			index2 += index1
		}
		call := path[index1:index2]
		goid := fmt.Sprintf("%06d", goid.Get())

		return "[" + goid + "|" + call + "]" + " " + msg + "\t\t"
	}
	return msg + "\t\t"
}
