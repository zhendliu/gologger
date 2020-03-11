package main

import (
	"gologger/mylogger"
)

func main() {

	flogger := mylogger.NewFileLogger("debug", ".", "testLog.log", 10240000)
	i := 1
	flogger.Debug("Debug:%s,%s,%d", "Debug", "ToTestDebug", i)
	flogger.Trace("Trace:%s,%s,%d", "Trace", "ToTestTrace", i)
	flogger.Info("Info:%s,%s,%d", "Info", "ToTestInfo", i)
	flogger.Warning("Warning:%s,%s,%d","Warning","ToTestWarning",i)
	flogger.Error("Error:%s,%s,%d","Error","ToTestError",i)
	flogger.Fatal("Fatal:%s,%s,%d ","Fatal","ToTestFatal",i)


}
