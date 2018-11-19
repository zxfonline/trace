package trace

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/zxfonline/fileutil"
	"github.com/zxfonline/timefix"
)

var (
	traceLogs  map[string]*TraceLog
	TimePeriod = 10 * time.Minute
)

type TraceLog struct {
	Name           string //name唯一
	TotalLog       *log.Logger
	TotalLogFile   *os.File
	DetailLog      *log.Logger
	DetailLogFile  *os.File
	FilePath       string
	FileNamePrefix string
}

//InitTraceLog 初始化跟踪日志
func initTraceLog() {
	traceLogs = make(map[string]*TraceLog)
	go writeloop()
}

//RegisterTraceLog 注册跟踪日志日志
func RegisterTraceLog(family string, filePath, fileNamePrefix string) error {
	fileName := fileutil.TransPath(filepath.Join(filePath, fileNamePrefix+"_"+"total"+"_"+time.Now().Format("2006-01-02")+".log"))

	logFile1, err := fileutil.OpenFile(fileName, fileutil.DefaultFileFlag, fileutil.DefaultFileMode)
	if err != nil {
		log.Printf("open file err:%v\n", err)
		return err
	}
	mylog1 := log.New(logFile1, "", 0)

	fileName = fileutil.TransPath(filepath.Join(filePath, fileNamePrefix+"_"+"detail"+"_"+time.Now().Format("2006-01-02")+".log"))

	logFile2, err := fileutil.OpenFile(fileName, fileutil.DefaultFileFlag, fileutil.DefaultFileMode)
	if err != nil {
		log.Printf("open file err:%v\n", err)
		return err
	}
	mylog2 := log.New(logFile2, "", 0)

	traceLog := &TraceLog{Name: family, TotalLog: mylog1, TotalLogFile: logFile1, DetailLog: mylog2, DetailLogFile: logFile2, FilePath: filePath, FileNamePrefix: fileNamePrefix}
	traceLogs[family] = traceLog

	return nil
}

func writeloop() {
	pm := time.NewTimer(TimePeriod)
	baset := time.Now()
	pm1 := time.NewTimer(time.Duration(timefix.NextMidnight(baset, 1).Unix()-baset.Unix()) * time.Second)
	for {
		select {
		case <-pm.C:
			saveTraceLog()
			pm.Reset(TimePeriod)
		case <-pm1.C:
			now := time.Now()
			for _, traceLog := range traceLogs {
				fileName := fileutil.TransPath(filepath.Join(traceLog.FilePath, traceLog.FileNamePrefix+"_"+"total"+"_"+now.Format("2006-01-02")+".log"))

				logFile1, err := fileutil.OpenFile(fileName, fileutil.DefaultFileFlag, fileutil.DefaultFileMode)
				if err != nil {
					log.Printf("[ERROR] "+"open file err:%v\n", err)
					continue
				}

				fileName = fileutil.TransPath(filepath.Join(traceLog.FilePath, traceLog.FileNamePrefix+"_"+"detail"+"_"+now.Format("2006-01-02")+".log"))
				logFile2, err := fileutil.OpenFile(fileName, fileutil.DefaultFileFlag, fileutil.DefaultFileMode)
				if err != nil {
					log.Printf("[ERROR] "+"open file err:%v\n", err)
					continue
				}

				traceLog.TotalLog.SetOutput(logFile1)
				traceLog.TotalLogFile.Close()
				traceLog.TotalLogFile = logFile1

				traceLog.DetailLog.SetOutput(logFile2)
				traceLog.DetailLogFile.Close()
				traceLog.DetailLogFile = logFile2
			}
			pm1.Reset(time.Duration(timefix.NextMidnight(now, 1).Unix()-now.Unix()) * time.Second)
		}
	}
}

func saveTraceLog() {
	defer func() {
		if x := recover(); x != nil {
			fmt.Printf("Recovered %v\n.", x)
		}
	}()
	for _, traceLog := range traceLogs {
		traceLog.TotalLog.Println(GetFamilyTotalString(traceLog.Name))
		for i := 0; i <= 9; i++ {
			if str := GetFamilyDetailString(traceLog.Name, i); len(str) > 0 {
				traceLog.DetailLog.Println(str)
			}
		}
	}
}
