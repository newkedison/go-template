package main

import (
	// begin gin
	"context"
	"net"
	"net/http"
	"strconv"
	"time"
	// end gin
	// begin tcp
	"net"
	"strconv"
	// end tcp
	"os"
	"path"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	buildin_log "log"

	"TEMPLATE/utils/service"

	// begin gin
	"github.com/gin-gonic/gin"
	"github.com/toorop/gin-logrus"
	// end gin
)

func getCurrentPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

// Do all your work in this function
func Main(exit <-chan struct{}) {
	// log some info if you want
	if service.Interactive() {
		log.Debug("Running in terminal.")
	} else {
		log.Debug("Running under service manager.")
	}
	log.Debug("Platform:", service.Platform())
	log.Info("Log level:", log.GetLevel())

	// begin basic
	// 1. do some init work
	// 2. start an goroutine to do long time work
	// end basic
	// begin gin
	// init gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(ginlogrus.Logger(log.StandardLogger()), gin.Recovery())
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello TEMPLATE user")
	})

	addr := net.JoinHostPort(globalConfig.Listener.IP,
		strconv.Itoa(globalConfig.Listener.Port))
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// start a new goroutine to do long time work
	go func() {
		// service connections
		log.Info("Listening on ", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen on %s fail: %s\n", addr, err)
		}
	}()
	// end gin
	// begin tcp
	addr := net.JoinHostPort(globalConfig.Listener.IP,
		strconv.Itoa(globalConfig.Listener.Port))
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen on %s fail: %s\n", addr, err)
	}
	go func() {
		log.Info("Listening on ", addr)
		for {
			conn, err := l.Accept()
			if err != nil {
				log.Warn("Accept fail:", err)
				continue
			}
			go func() {
				for {
					buf := make([]byte, 4096)
					n, err := conn.Read(buf)
					if err != nil {
						conn.Close()
						break
					}
					conn.Write(buf[:n])
				}
			}()
		}
	}()
	// end tcp

	// wating for the exit signal
	<-exit
	// begin basic
	// 3. clean resource, as quickly as possiable
	// end basic
	// begin gin
	// shutdown gin server
	log.Info("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Info("Server exiting")
	// end gin
	// begin tcp
	l.Close()
	// end tcp
}

func main() {
	var err error
	if globalConfig, err = readConfig(
		path.Join(getCurrentPath(), DefaultConfigFile)); err != nil {
		buildin_log.Fatalln("Read config fail:", err)
	}
	if err := initLogger(globalConfig); err != nil {
		buildin_log.Fatalln("Config logger fail:", err)
	}

	service.Init(service.ServiceOption{
		Name:        globalConfig.Service.ServiceName,
		DisplayName: globalConfig.Service.DisplayName,
		Description: globalConfig.Service.Description,
	})
	err = service.Run(Main)
	if err != nil {
		log.Error(err)
	}
}
