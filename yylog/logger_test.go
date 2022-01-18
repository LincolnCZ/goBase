package yylog

import (
	"errors"
	"log"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func initAsyncFile() {
	InitLog(ProcessName("loggertest"), SetTarget("asyncfile"), SetEncode("yyjson"), LogFilePath("./"))
}

func TestInitLog(t *testing.T) {
	type args struct {
		options []LogOption
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"stdout",
			args{
				[]LogOption{ProcessName("loggertest"), SetTarget("stdout"), HostName("lulu")},
			}, false,
		},
		{
			"asyncfile",
			args{
				[]LogOption{ProcessName("loggertest"), SetTarget("asyncfile"), SetEncode("yyjson"), LogFilePath("./")},
			}, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InitLog(tt.args.options...); (err != nil) != tt.wantErr {
				t.Errorf("InitLog() error = %v, wantErr %v", err, tt.wantErr)
			}
			defaultLog.writelog("info", "testinfo", zap.Any("mykey", "myvalue"))
		})
	}
}

func TestGetLogger(t *testing.T) {
	tests := []struct {
		name string
		want *yyloger
	}{
		{
			"getlogger",
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetLogger(); got == nil {
				t.Errorf("GetLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_yyloger_GetZLog(t *testing.T) {

	type args struct {
		opts []zap.Option
	}
	tests := []struct {
		name string
		ylog *yyloger
		args args
		want *zap.Logger
	}{
		{
			"getzlog1",
			defaultLog,
			args{
				[]zap.Option{},
			}, nil,
		},
		{
			"getzlog2",
			defaultLog,
			args{
				[]zap.Option{zap.AddCallerSkip(2)},
			}, nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := tt.ylog.GetZLog(tt.args.opts...); got == nil {
				t.Errorf("yyloger.GetZLog() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_yyloger_Clone(t *testing.T) {

	type args struct {
		opts []zap.Option
	}
	tests := []struct {
		name   string
		ylog   *yyloger
		args   args
		unwant *yyloger
	}{
		{
			"clone1",
			defaultLog,
			args{
				[]zap.Option{},
			}, nil,
		},
		{
			"clone2",
			defaultLog,
			args{
				[]zap.Option{zap.AddCallerSkip(2)},
			}, nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := tt.ylog.Clone(tt.args.opts...); got == tt.unwant {
				t.Errorf("yyloger.Clone() = %v, unwant %v", got, tt.unwant)
			}
		})
	}
}

func Test_yyloger_Write(t *testing.T) {

	type args struct {
		p []byte
	}
	tests := []struct {
		name    string
		ylog    *yyloger
		args    args
		wantN   int
		wantErr bool
	}{
		{
			"write1",
			defaultLog,
			args{
				[]byte("mydata"),
			}, 6, false,
		},
		{
			"write2",
			defaultLog,
			args{
				[]byte("mydatadd"),
			}, 8, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotN, err := tt.ylog.Write(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("yyloger.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("yyloger.Write() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}

func Test_yyloger_Log(t *testing.T) {

	type args struct {
		v []interface{}
	}
	tests := []struct {
		name string
		ylog *yyloger
		args args
	}{
		{
			"log",
			defaultLog,
			args{
				[]interface{}{"err", 9889, "string", errors.New("myerr")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ylog.Log(tt.args.v...)
		})
	}
}

func Test_yyloger_Logf(t *testing.T) {

	type args struct {
		format string
		v      []interface{}
	}
	tests := []struct {
		name string
		ylog *yyloger
		args args
	}{
		{
			"logf",
			defaultLog,
			args{
				"msg err %v rescode %u errmsg %s",
				[]interface{}{errors.New("myerr"), 7, "nogood"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ylog.Logf(tt.args.format, tt.args.v...)
		})
	}
}

func Test_yyloger_writelog(t *testing.T) {

	type args struct {
		level  string
		msg    string
		fields []zapcore.Field
	}
	tests := []struct {
		name string
		ylog *yyloger
		args args
	}{
		{
			"writelog_debug",
			defaultLog.Clone(),
			args{
				"debug",
				"testdebug",
				[]zapcore.Field{zap.String("mykey", "myvalue"), zap.Any("key", "value")},
			},
		},
		{
			"writelog_info",
			defaultLog.Clone(),
			args{
				"info",
				"testinfo",
				[]zapcore.Field{zap.String("mykey", "myvalue"), zap.Any("key", "value")},
			},
		},
		{
			"writelog_err",
			defaultLog.Clone(),
			args{
				"error",
				"testerr",
				[]zapcore.Field{zap.String("mykey", "myvalue"), zap.Any("key", "value")},
			},
		},
		{
			"writelog_nolevel",
			defaultLog.Clone(),
			args{
				"what",
				"testnolevel",
				[]zapcore.Field{zap.String("mykey", "myvalue"), zap.Any("key", "value")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.ylog.writelog(tt.args.level, tt.args.msg, tt.args.fields...)
		})
	}
}

func TestLog(t *testing.T) {
	type args struct {
		level string
		v     []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"log",
			args{
				"info",
				[]interface{}{errors.New("llll"), 888, "data"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Log(tt.args.level, tt.args.v...)
		})
	}
}

func TestLogF(t *testing.T) {
	type args struct {
		level  string
		format string
		v      []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"logf",
			args{
				"debug",
				"mymsg err %v rescode %d %s",
				[]interface{}{errors.New("llll"), 888, "data"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			LogF(tt.args.level, tt.args.format, tt.args.v...)
		})
	}
}

func TestDebug(t *testing.T) {
	type args struct {
		msg    string
		fields []zapcore.Field
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"debug",
			args{
				"getdebug",
				[]zapcore.Field{zap.String("mykey", "myvalue"), zap.Any("key", "value")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Debug(tt.args.msg, tt.args.fields...)
		})
	}
}

func TestInfo(t *testing.T) {
	type args struct {
		msg    string
		fields []zapcore.Field
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"info",
			args{
				"getinfo",
				[]zapcore.Field{zap.String("mykey", "myvalue"), zap.Any("key", "value")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Info(tt.args.msg, tt.args.fields...)
		})
	}
}

func TestWarn(t *testing.T) {
	type args struct {
		msg    string
		fields []zapcore.Field
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"warn",
			args{
				"getwarn",
				[]zapcore.Field{zap.String("mykey", "myvalue"), zap.Any("key", "value")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Warn(tt.args.msg, tt.args.fields...)
		})
	}
}

func TestError(t *testing.T) {
	type args struct {
		msg    string
		fields []zapcore.Field
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"error",
			args{
				"geterror",
				[]zapcore.Field{zap.String("mykey", "myvalue"), zap.Any("key", "value")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Error(tt.args.msg, tt.args.fields...)
		})
	}
}

func TestPanic(t *testing.T) {
	type args struct {
		msg    string
		fields []zapcore.Field
	}
	tests := []struct {
		name string
		args args
	}{

		{
			"panic",
			args{
				"getpanic",
				[]zapcore.Field{zap.String("mykey", "myvalue"), zap.Any("key", "value")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					return
				}
				t.Error("not panic")

			}()
			Panic(tt.args.msg, tt.args.fields...)
		})
	}
}

func TestFatal(t *testing.T) {
	type args struct {
		msg    string
		fields []zapcore.Field
	}
	tests := []struct {
		name string
		args args
	}{
		/*
			{
				"fatal",
				args{
					"getfail",
					[]zapcore.Field{zap.String("mykey","myvalue"),zap.Any("key","value")},
				},
			},
		*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Fatal(tt.args.msg, tt.args.fields...)
		})
	}
}

func TestSync(t *testing.T) {
	initAsyncFile()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			"sync",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Sync(); (err != nil) != tt.wantErr {
				t.Errorf("Sync() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetLogLevel(t *testing.T) {
	type args struct {
		level string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"setloglevle",
			args{
				"debug",
			}, false,
		},

		{
			"setloglevle",
			args{
				"off",
			}, false,
		},
		{
			"setloglevle",
			args{
				"unknown",
			}, true,
		},
		{
			"setloglevle",
			args{
				"all",
			}, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetLogLevel(tt.args.level); (err != nil) != tt.wantErr {
				t.Errorf("SetLogLevel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetStdLogger(t *testing.T) {
	tests := []struct {
		name string
		want *log.Logger
	}{
		{
			"stdloger",
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetStdLogger()
			if got == nil {
				t.Errorf("GetStdLogger() = %v, want %v", got, tt.want)
			}

			got.Printf("test %s \n", "logtest")
			log.Printf("test %s \n", "msgdata")
		})
	}
}
