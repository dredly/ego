package main

import (
	"flag"
	"os"

	"github.com/dredly/ego/internal/db"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	logCfg := zap.NewDevelopmentConfig()
	logCfg.EncoderConfig = zapcore.EncoderConfig{ 
		TimeKey: "", 
		LevelKey: "", 
		NameKey: "", 
		CallerKey: "", 
		MessageKey: "M", 
		StacktraceKey: "",
	}

	unsugared, err := logCfg.Build()
    if err != nil {
        panic("Failed to create logger: " + err.Error())
    }
	defer unsugared.Sync()
	logger := unsugared.Sugar()

	createCmd := flag.NewFlagSet("create", flag.ExitOnError)

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addName := addCmd.String("name", "", "name of the player to add to the leaderboard")

	recordCmd := flag.NewFlagSet("record", flag.ExitOnError)
	recordWinner := recordCmd.String("w", "", "name of the player who won")
	recordLoser := recordCmd.String("l", "", "name of the player who lost")

	if len(os.Args) < 2 {
		logger.Fatal("expected a subcommand")
    }

	switch os.Args[1] {
	case "create":
        createCmd.Parse(os.Args[2:])
		conn, err := db.New()
		if err != nil {
			logger.Fatalf("failed to get db connection: %w", err)
		}
		err = conn.Initialise()
		if err != nil {
			logger.Fatalf("failed to initialise db: %v", err)
		}
		logger.Info("Successfully initialised leaderboard")
    case "add":
        addCmd.Parse(os.Args[2:])
		logger.Infof("added new player %s\n", *addName)
	case "record":
		recordCmd.Parse(os.Args[2:])
		logger.Infof("recorded %s win over %s\n", *recordWinner, *recordLoser)
    default:
		logger.Fatalf("unrecognised subcommand: %s\n", os.Args[1])
	}
}