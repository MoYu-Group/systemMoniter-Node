package logger

import (
	"os"

	"github.com/natefinch/lumberjack"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLogger 初始化Logger
func Init() (err error) {
	//配置分割文档的属性
	writeSyncer := getLogWriter(
		viper.GetString("log.filename"),
		viper.GetInt("log.max_size"),
		viper.GetInt("log.max_backups"),
		viper.GetInt("log.max_go"),
	)
	//配置输出的格式
	encoder := getEncoder()
	//创建日志级别
	var l = new(zapcore.Level)
	err = l.UnmarshalText([]byte(viper.GetString("log.level")))
	if err != nil {
		return
	}
	//定制logger
	multiWriteSyncerzapcore := zapcore.NewMultiWriteSyncer(writeSyncer, zapcore.AddSync(os.Stdout))
	core := zapcore.NewCore(encoder, multiWriteSyncerzapcore, l)

	lg := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(lg) // 替换zap包中全局的logger实例，后续在其他包中只需使用zap.L()调用即可
	return
}

//更改时间编码并添加调用者详细信息
func getEncoder() zapcore.Encoder {
	//编码器(如何写入日志)
	encoderConfig := zap.NewProductionEncoderConfig()
	//修改时间
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	//输出时间的key名
	encoderConfig.TimeKey = "time"
	//在日志文件中使用大写字母记录日志级别
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	//执行消耗的时间转化成浮点型的秒
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	//以包/文件:行号 格式化调用堆栈
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

//使用Lumberjack进行日志切割归档
//对Lumberjack进行配置属性
func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}
	return zapcore.AddSync(lumberJackLogger)
}
