package main

import (
	"flag"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"
)

type MMSI uint32 // 9 digits = 1 billion values = 30 bits
type Message struct {
	completed time.Time  // of last received sentence
	sentences []Sentence // one or more AIS sentences
	source    string     // AIS listener
	msg_type  uint8
	content   []byte
}

var (
	Log = NewLogger(os.Stderr, LOG_DEBUG, 10*time.Second)
	// For input sentence or message "errors"
	AisLog = NewLogger(os.Stdout, LOG_DEBUG, 0)
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	signalChan := make(chan os.Signal, 1)
	// Intercept ^C and `timeout`s.
	// Catching SIGPIPE has no effect if it was what Log wrote to that broke, as it's, well, broken.
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGPIPE)

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			Log.Fatal(err.Error())
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	merger := make(chan Message, 200)
	go func() {
		for _ = range merger {

		}
	}()

	Log.AddPeriodicLogger("from_main", 120*time.Second, func(l *Logger, _ time.Duration) {
		c := l.Compose(LOG_DEBUG)
		c.Writeln("waiting to be merged: %d/%d", len(merger), cap(merger))
		c.Writeln("source connections: %d", listener_connections)
		c.Close()
	})

	Read("ECC", "http://aishub.ais.ecc.no/raw", 5*time.Second, merger)
	Read("kystverket", "tcp://153.44.253.27:5631", 5*time.Second, merger)
	//Read("http_timeout", "http://127.0.0.1:12340", 8*time.Second, merger)
	//Read("tcp_timeout", "tcp://127.0.0.1:12341", 2*time.Second, merger)
	//Read("redirect", "http://localhost:12342", 0*time.Second, merger)
	//Read("redirect_loop", "http://localhost:12343", 0*time.Second, merger)
	//Read("http_flood", "http://localhost:12344", 2*time.Second, merger)
	//Read("tcp_flood", "tcp://localhost:12345", 2*time.Second, merger)
	//Read("file", "minute_ecc.log", 0*time.Second, merger)

	// Here we wait for CTRL-C or some other kill signal
	_ = <-signalChan
	Log.Info("\n...Stopping...")
	Log.RunPeriodicLoggers(time.Now().Add(1 * time.Hour))
	AisLog.RunPeriodicLoggers(time.Now().Add(1 * time.Hour))
}
