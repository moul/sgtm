package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"os"
	"syscall"

	"github.com/oklog/run"
	ff "github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"moul.io/bounce/pkg/bounce"
	"moul.io/srand"
)

func main() {
	if err := app(os.Args); err != nil {
		log.Fatalf("error: %v", err)
		os.Exit(1)
	}
}

var opts bounce.Opts

func app(args []string) error {
	opts = bounce.DefaultOpts()
	rootFlags := flag.NewFlagSet("root", flag.ExitOnError)
	rootFlags.BoolVar(&opts.DevMode, "dev-mode", opts.DevMode, "start in developer mode")
	/// discord
	rootFlags.BoolVar(&opts.EnableDiscord, "enable-discord", opts.EnableDiscord, "enable discord bot")
	rootFlags.StringVar(&opts.DiscordToken, "discord-token", opts.DiscordToken, "discord bot token")
	rootFlags.StringVar(&opts.DiscordAdminChannel, "discord-admin-channel", opts.DiscordAdminChannel, "discord channel ID for admin messages")
	/// server
	rootFlags.BoolVar(&opts.EnableServer, "enable-server", opts.EnableServer, "enable HTTP+gRPC Server")
	rootFlags.StringVar(&opts.ServerBind, "server-bind", opts.ServerBind, "server bind (HTTP + gRPC)")
	rootFlags.StringVar(&opts.ServerCORSAllowedOrigins, "server-cors-allowed-origins", opts.ServerCORSAllowedOrigins, "allowed CORS origins")
	rootFlags.DurationVar(&opts.ServerRequestTimeout, "server-request-timeout", opts.ServerRequestTimeout, "server request timeout")
	rootFlags.DurationVar(&opts.ServerShutdownTimeout, "server-shutdown-timeout", opts.ServerShutdownTimeout, "server shutdown timeout")
	rootFlags.BoolVar(&opts.ServerWithPprof, "server-with-pprof", opts.ServerWithPprof, "enable pprof on HTTP server")

	root := &ffcli.Command{
		FlagSet: rootFlags,
		Options: []ff.Option{
			ff.WithEnvVarPrefix("BOUNCE"),
			ff.WithConfigFile("config.txt"),
			ff.WithConfigFileParser(ff.PlainParser),
		},
		Subcommands: []*ffcli.Command{
			{Name: "run", Exec: runCmd},
		},
		Exec: func(context.Context, []string) error {
			return flag.ErrHelp
		},
	}

	return root.ParseAndRun(context.Background(), args[1:])
}

func runCmd(ctx context.Context, _ []string) error {
	// init
	rand.Seed(srand.Secure())
	gr := run.Group{}
	gr.Add(run.SignalHandler(ctx, syscall.SIGTERM, syscall.SIGINT, os.Interrupt, os.Kill))

	// bearer
	config := zap.NewDevelopmentConfig()
	config.Level.SetLevel(zap.DebugLevel)
	config.DisableStacktrace = true
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := config.Build()
	if err != nil {
		return err
	}
	opts.Logger = logger
	//opts.Context = ctx

	// init service
	svc := bounce.New(opts)
	defer svc.Close()

	// start
	if opts.EnableDiscord {
		gr.Add(svc.StartDiscord, svc.CloseDiscord)
	}
	if opts.EnableServer {
		gr.Add(svc.StartServer, svc.CloseServer)
	}
	return gr.Run()
}