package cmd

import (
	"context"
	"os"
	"path/filepath"
	"time"

	_ "github.com/ihezebin/soup"
	"github.com/ihezebin/soup/logger"
	"github.com/ihezebin/soup/runner"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/ihezebin/changes2queue/component/cache"
	"github.com/ihezebin/changes2queue/config"
	"github.com/ihezebin/changes2queue/domain/repository"
	"github.com/ihezebin/changes2queue/domain/service"
	"github.com/ihezebin/changes2queue/server"
	"github.com/ihezebin/changes2queue/worker"
	"github.com/ihezebin/changes2queue/worker/example"
)

var (
	configPath string
)

func Run(ctx context.Context) error {

	app := &cli.App{
		Name:    "changes2queue",
		Version: "v1.0.1",
		Usage:   "Rapid construction template of Web service based on DDD architecture",
		Authors: []*cli.Author{
			{Name: "hezebin", Email: "ihezebin@qq.com"},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Destination: &configPath,
				Name:        "config", Aliases: []string{"c"},
				Value: "./config/config.toml",
				Usage: "config file path (default find file from pwd and exec dir",
			},
		},
		Before: func(c *cli.Context) error {
			if configPath == "" {
				return errors.New("config path is empty")
			}

			conf, err := config.Load(configPath)
			if err != nil {
				return errors.Wrapf(err, "load config error, path: %s", configPath)
			}
			logger.Debugf(ctx, "load config: %+v", conf.String())

			if err = initComponents(ctx, conf); err != nil {
				return errors.Wrap(err, "init components error")
			}

			return nil
		},
		Action: func(c *cli.Context) error {
			httpServer, err := server.NewServer(ctx, config.GetConfig())
			if err != nil {
				return errors.Wrap(err, "new http server err")
			}

			tasks := make([]runner.Task, 0)
			tasks = append(tasks, worker.NewWorKeeper(example.NewExampleWorker()))
			tasks = append(tasks, httpServer)

			runner.NewRunner(tasks...).Run(ctx)

			return nil
		},
	}

	return app.Run(os.Args)
}

func initComponents(ctx context.Context, conf *config.Config) error {
	// init logger
	if conf.Logger != nil {
		logger.ResetLoggerWithOptions(
			logger.WithServiceName(conf.ServiceName),
			logger.WithPrettyCallerHook(),
			logger.WithTimestampHook(),
			logger.WithLevel(conf.Logger.Level),
			//logger.WithLocalFsHook(filepath.Join(conf.Pwd, conf.Logger.Filename)),
			// 每天切割，保留 3 天的日志
			logger.WithRotateLogsHook(filepath.Join(conf.Pwd, conf.Logger.Filename), time.Hour*24, time.Hour*24*3),
		)
	}

	// init storage
	// if err := storage.InitMySQLClient(ctx, conf.MysqlDsn); err != nil {
	// 	return errors.Wrap(err, "init mysql storage client error")
	// }
	// if err := storage.InitMongoClient(ctx, conf.MongoDsn); err != nil {
	// 	return errors.Wrap(err, "init mongo storage client error")
	// }

	// init cache
	cache.InitMemoryCache(time.Minute*5, time.Minute)

	// init repository
	repository.Init()

	// init service
	service.Init()

	return nil
}
