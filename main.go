package main

import (
	"bytes"
	"fmt"

	"github.com/CloudyKit/jet"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

var View = jet.NewHTMLSet("./views")

var logLevelMap = map[string]logrus.Level{
	"trace": logrus.TraceLevel,
	"debug": logrus.DebugLevel,
	"info":  logrus.InfoLevel,
	"warn":  logrus.WarnLevel,
	"error": logrus.ErrorLevel,
}

type arguments struct {
	LogLevel       string
	BindAddress    string
	BindPort       int
	StaticContents string
}

func runServer(args arguments) error {
	level, ok := logLevelMap[args.LogLevel]
	if !ok {
		return fmt.Errorf("Invalid log level: %s", args.LogLevel)
	}
	logger.SetLevel(level)
	logger.SetFormatter(&logrus.JSONFormatter{})

	logger.WithFields(logrus.Fields{
		"args": args,
	}).Info("Given options")

	r := gin.Default()

	r.Use(static.Serve("/", static.LocalFile(args.StaticContents, false)))

	r.GET("/hello", func(c *gin.Context) {
		templateName := "home.jet"
		t, err := View.GetTemplate(templateName)
		if err != nil {
			// template could not be loaded
		}
		var w bytes.Buffer
		vars := make(jet.VarMap)
		if err = t.Execute(&w, vars, nil); err != nil {
			// error when executing template
		}

		c.Writer.WriteHeader(200)
		c.Writer.Write(w.Bytes())
	})

	if err := r.Run(fmt.Sprintf("%s:%d", args.BindAddress, args.BindPort)); err != nil {
		return err
	}

	return nil
}

func main() {
	args := arguments{
		LogLevel:       "info",
		BindAddress:    "0.0.0.0",
		BindPort:       8990,
		StaticContents: "./static",
	}

	if err := runServer(args); err != nil {
		logger.WithError(err).Fatal("Server exits with error")
	}
}
