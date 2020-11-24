package launcher

import (
	"crypto/sha1"
	"io"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/launcher/cmd"
	launcherConfig "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/launcher/config"
	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/launcher/service"
	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/cache"
	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"
	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/utils/log"
)

type ApplicationDescription struct {
	Usage            string
	ShortDescription string
	LongDescription  string
}

type ApplicationOption func(app *Application)

type Application struct {
	logger      *logrus.Entry
	description *ApplicationDescription
	services    []service.Interface
	config      *launcherConfig.StandardConfig
	events      *Events
}

func NewApplication(options ...ApplicationOption) *Application {

	app := &Application{
		description: &ApplicationDescription{},
		services:    make([]service.Interface, 0),
		config:      &launcherConfig.StandardConfig{},
		events:      &Events{},
	}

	for _, option := range options {
		option(app)
	}

	return app
}

func (app *Application) Launch() {

	app.logger.Debug("start to launch application")

	app.logger.Debug("init root command")

	cmd.InitRootCommand(
		cmd.SetCommandUsage(app.description.Usage),
		cmd.SetCommandShortDescription(app.description.ShortDescription),
		cmd.SetCommandLongDescription(app.description.LongDescription),
	)

	cmd.GetRootCommand().Run = func(cmd *cobra.Command, args []string) {

		app.logger.Debug("start unmarshal configuration")
		err := viper.Unmarshal(app.config)
		if err != nil {
			logrus.WithError(err).Error("unmarshal config error")
			os.Exit(1)
			return
		}

		app.logger.Debug("unmarshal configuration completed")
		app.logger.WithField("config", app.config).Info("loaded configuration")

		app.init()

		app.start()

		app.waitSignal()
	}

	app.logger.Debug("initialized root command")

	app.logger.Debug("execute root command")

	cmd.Execute()
}

func (app *Application) start() {

	app.startMySQLClient()
	app.startRedisClient()

	app.logger.Debug("start services")
	for _, svc := range app.services {
		if err := svc.OnStart(); err != nil {

			logrus.WithError(err).WithField("service", svc.GetServiceName()).Error("start service error")
			os.Exit(1)
			return
		}
	}
	app.logger.Debug("services started")

	if app.events.OnStart != nil {

		app.logger.Debug("load on start customer function")
		app.events.OnStart(app)
		app.logger.Debug("loaded on start customer function")
	}
}

func (app *Application) close() {

	for _, svc := range app.services {
		if err := svc.OnStop(); err != nil {

			logrus.WithError(err).WithField("service", svc.GetServiceName()).Error("stop service error")
			continue
		}
	}

	if app.events.OnClose != nil {

		app.events.OnClose(app)
	}

	os.Exit(0)
}

type WaitingToDo func()

func (app *Application) waitSignal() {

	chanSignal := make(chan os.Signal, 1)
	signal.Notify(chanSignal, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	for {
		select {
		case sig := <-chanSignal:
			logrus.Infof("Received signal: %d", sig)
			app.close()

			goto exit
		}
	}

exit:
	logrus.Infof("loop exited")
	return
}

func (app *Application) startMySQLClient() {

	app.logger.Debug("start to connect mysql clients")
	database.SetLogger(app.logger.Logger)

	for key, mysqlConfig := range app.config.MySQL {

		err := database.Connect(key, database.NewMySQLConfig(
			database.MySQLHost(mysqlConfig.GetHost()),
			database.MySQLPort(mysqlConfig.GetPort()),
			database.MySQLUsername(mysqlConfig.GetUsername()),
			database.MySQLPassword(mysqlConfig.GetPassword()),
			database.MySQLDatabase(mysqlConfig.GetDatabase()),
			database.MySQLLogMode(mysqlConfig.GetLogMode()),
		))

		if err != nil {
			logrus.WithError(err).WithField("key", key).Error("connect mysql error")
		}
	}

	app.logger.Debug("mysql client connected")
}

func (app *Application) startRedisClient() {

	app.logger.Debug("start to connect redis clients")
	cache.SetLogger(app.logger.Logger)

	for key, redisConfig := range app.config.Redis {

		err := cache.Connect(key, cache.NewRedisConfig(
			cache.RedisHost(redisConfig.GetHost()),
			cache.RedisPort(redisConfig.GetPort()),
			cache.RedisPassword(redisConfig.GetPassword()),
			cache.RedisDatabase(redisConfig.GetDatabase()),
		))

		if err != nil {
			logrus.WithError(err).WithField("key", key).Error("connect mysql error")
		}
	}

	app.logger.Debug("redis clients connected")
}

func (app *Application) init() {

	app.logger.Debug("start to init application")
	app.initRandomSeed()
	app.initServiceId()
	app.initLogger()
	app.initWebService()
	app.initRPCService()

	if app.events.OnInit != nil {
		app.events.OnInit(app)
	}

	app.logger.Debug("init completed")
}

func (app *Application) initWebService() {

	if !app.config.Web.Enable {

		app.logger.Info("web service disabled")
		return
	}

	app.logger.Info("start to init web service")

	app.services = append(app.services,
		service.NewGinService(app.logger,
			service.NewGinConfig(
				service.GinConfigListenConfig(
					service.NewGinListenConfig(
						service.GinListenConfigIP(app.config.Web.IP),
						service.GinListenConfigPort(app.config.Web.Port),
						service.GinListenConfigReadWriteTimeout(app.config.Web.ReadWriteTimeout),
					),
				),
				service.GinConfigWebServiceMode(service.WebServiceMode(app.config.Web.Mode)),
			),
		),
	)

	app.logger.Debug("init web service completed")
}

func (app *Application) GetWebService() *service.Gin {

	for _, svc := range app.services {

		ginService, ok := svc.(*service.Gin)
		if ok {
			return ginService
		}
	}

	return nil
}

func (app *Application) GetRPCService() *service.RPC {

	for _, svc := range app.services {

		rpcService, ok := svc.(*service.RPC)
		if ok {
			return rpcService
		}
	}

	return nil
}

func (app *Application) GetServiceId() uint16 {

	return app.config.ServiceId
}

func (app *Application) initRPCService() {

	if !app.config.RPC.Enable {

		app.logger.Info("rpc service disabled")
		return
	}

	app.logger.Info("start to init rpc service")

	app.services = append(app.services,
		service.NewRPCService(app.logger,
			service.NewRPCConfig(
				service.RPCConfigListenConfig(
					service.NewRPCListenConfig(
						service.RPCListenConfigIP(app.config.RPC.IP),
						service.RPCListenConfigPort(app.config.RPC.Port),
					),
				),
			),
		),
	)

	app.logger.Debug("init rpc service completed")
}

func (app *Application) initLogger() {

	app.logger.Info("start to init logger")
	app.logger.Logger.Level = app.config.Log.GetLogLevel()
	log.RegisterFilePath(app.logger.Logger)
	app.logger = app.logger.WithField("serviceId", app.config.ServiceId)
	app.logger.Info("init logger succeed")
}

func (app *Application) initServiceId() {

	var err error
	var serviceId uint16
	var serviceIdUint64 uint64
	var serviceIdString string

	defer func() {
		if serviceId > 0 {
			app.config.ServiceId = serviceId
		}
		app.logger.Infof("serviceId: %d", app.config.GetServiceId())
	}()

	serviceIdString, err = cmd.GetRootCommand().Flags().GetString(cmd.FlagServiceId)
	if err != nil || len(serviceIdString) <= 0 {
		app.logger.WithError(err).Warn("get service id failed")
		app.config.ServiceId = 0
		return
	}

	app.logger.WithField("serviceIdString", serviceIdString).Info("input serviceId")

	serviceIdUint64, err = strconv.ParseUint(serviceIdString, 10, 64)

	if err == nil {
		// 超过了也硬转过去
		serviceId = uint16(serviceIdUint64)
		return
	}

	hash := sha1.New()
	_, err = io.WriteString(hash, serviceIdString)
	if err != nil {
		app.logger.WithError(err).Error("write sha1 stream error")
		return
	}

	hashResult := hash.Sum(nil)

	app.logger.WithField("service", serviceIdString).Infof("hash result: %x", hashResult)

	// 取头两位hash结果作为serviceId值，0作为高位，1作为低位
	serviceId = uint16(hashResult[0])<<0x8 + uint16(hashResult[1])
}

func (app *Application) initRandomSeed() {

	rand.Seed(time.Now().UnixNano())
}

func SetApplicationDescription(description *ApplicationDescription) ApplicationOption {

	return func(app *Application) {
		if description != nil {
			app.description = description
		}
	}
}

func SetApplicationLogger(logger *logrus.Entry) ApplicationOption {

	return func(app *Application) {
		if logger != nil {
			app.logger = logger
		}
	}
}

type Event func(app *Application)

type Events struct {
	OnInit  Event
	OnStart Event
	OnClose Event
}

type ApplicationEventOption func(events *Events)

func SetOnInitEvent(event Event) ApplicationEventOption {

	return func(events *Events) {

		events.OnInit = event
	}
}

func SetOnStartEvent(event Event) ApplicationEventOption {

	return func(events *Events) {

		events.OnStart = event
	}
}

func SetOnCloseEvent(event Event) ApplicationEventOption {

	return func(events *Events) {

		events.OnClose = event
	}
}

func NewApplicationEvents(options ...ApplicationEventOption) *Events {

	events := &Events{}

	for _, option := range options {

		option(events)
	}

	return events
}

func SetApplicationEvents(events *Events) ApplicationOption {

	return func(app *Application) {

		if events != nil {
			app.events = events
		}
	}
}
