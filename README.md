#这是一个简单的go log输出库
实现了对控制台和文件输出，并且增加了缓冲
特点描述：
  1：runtime.Caller 调用一次，效率高
  2：支持缓冲写入，速度很快。
### 如何使用
```go
    
    package main
    
    import (
    	"gologger/mylogger"
    )
    
    func main() {
        /*
        debug:日志级别，日志级别分别为（debug，trace，info，warning，error，fatal，默认为info）
        /opt/log/:为日志路径
        testLog.log:为日志名称
        10240000:单个日志最大限制
        */
    	flogger := mylogger.NewFileLogger("debug", "/opt/log/", "testLog.log", 10240000)
    	i := 1
    	flogger.Debug("Debug:%s,%s,%d", "Debug", "ToTestDebug", i)
    	flogger.Trace("Trace:%s,%s,%d", "Trace", "ToTestTrace", i)
    	flogger.Info("Info:%s,%s,%d", "Info", "ToTestInfo", i)
    	flogger.Warning("Warning:%s,%s,%d","Warning","ToTestWarning",i)
    	flogger.Error("Error:%s,%s,%d","Error","ToTestError",i)
    	flogger.Fatal("Fatal:%s,%s,%d ","Fatal","ToTestFatal",i)
    }

  
