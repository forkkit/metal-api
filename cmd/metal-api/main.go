package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"go.uber.org/zap"

	"git.f-i-ts.de/cloud-native/metal/metal-api/cmd/metal-api/internal/datastore"
	"git.f-i-ts.de/cloud-native/metal/metal-api/cmd/metal-api/internal/metal"
	"git.f-i-ts.de/cloud-native/metal/metal-api/cmd/metal-api/internal/netbox"
	"git.f-i-ts.de/cloud-native/metal/metal-api/cmd/metal-api/internal/service"
	"git.f-i-ts.de/cloud-native/metal/metal-api/cmd/metal-api/internal/utils"
	"git.f-i-ts.de/cloud-native/metallib/bus"
	"git.f-i-ts.de/cloud-native/metallib/rest"
	"git.f-i-ts.de/cloud-native/metallib/version"
	"git.f-i-ts.de/cloud-native/metallib/zapup"
	restful "github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/go-openapi/spec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	cfgFileType = "yaml"
	moduleName  = "metal-api"
)

var (
	cfgFile  string
	ds       *datastore.RethinkStore
	producer bus.Publisher
	nbproxy  *netbox.APIProxy
	logger   = zapup.MustRootLogger().Sugar()
	debug    = false
)

var rootCmd = &cobra.Command{
	Use:     moduleName,
	Short:   "an api to offer pure metal",
	Version: version.V.String(),
	Run: func(cmd *cobra.Command, args []string) {
		initLogging()
		initDataStore()
		initEventBus()
		initNetboxProxy()
		initSignalHandlers()
		run()
	},
}

var dumpSwagger = &cobra.Command{
	Use:     "dump-swagger",
	Short:   "dump the current swagger configuration",
	Version: version.V.String(),
	Run: func(cmd *cobra.Command, args []string) {
		dumpSwaggerJSON()
	},
}

func main() {
	rootCmd.AddCommand(dumpSwagger)
	if err := rootCmd.Execute(); err != nil {
		logger.Error("failed executing root command", "error", err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "alternative path to config file")

	rootCmd.Flags().StringP("bind-addr", "", "127.0.0.1", "the bind addr of the api server")
	rootCmd.Flags().IntP("port", "", 8080, "the port to serve on")

	rootCmd.Flags().StringP("db", "", "rethinkdb", "the database adapter to use")
	rootCmd.Flags().StringP("db-name", "", "metalapi", "the database name to use")
	rootCmd.Flags().StringP("db-addr", "", "", "the database address string to use")
	rootCmd.Flags().StringP("db-user", "", "", "the database user to use")
	rootCmd.Flags().StringP("db-password", "", "", "the database password to use")

	rootCmd.Flags().StringP("nsqd-addr", "", "nsqd:4150", "the address of the nsqd")
	rootCmd.Flags().StringP("nsqd-http-addr", "", "nsqd:4151", "the address of the nsqd rest endpoint")
	rootCmd.Flags().StringP("nsqlookupd-addr", "", "nsqlookupd:4160", "the address of the nsqlookupd as a commalist")

	rootCmd.Flags().StringP("netbox-addr", "", "localhost:8001", "the address of netbox proxy")
	rootCmd.Flags().StringP("netbox-api-token", "", "", "the api token to access the netbox proxy")

	viper.BindPFlags(rootCmd.Flags())
}

func initConfig() {
	viper.SetEnvPrefix("METAL_API")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	viper.SetConfigType(cfgFileType)

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err != nil {
			logger.Error("Config file path set explicitly, but unreadable", "error", err)
		}
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath("/etc/" + moduleName)
		viper.AddConfigPath("$HOME/." + moduleName)
		viper.AddConfigPath(".")
		if err := viper.ReadInConfig(); err != nil {
			usedCfg := viper.ConfigFileUsed()
			if usedCfg != "" {
				logger.Error("Config file unreadable", "config-file", usedCfg, "error", err)
			}
		}
	}

	usedCfg := viper.ConfigFileUsed()
	if usedCfg != "" {
		logger.Info("Read config file", "config-file", usedCfg)
	}
}

func initLogging() {
	debug = logger.Desugar().Core().Enabled(zap.DebugLevel)
}

func initSignalHandlers() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		logger.Error("Received keyboard interrupt, shutting down...")
		if ds != nil {
			logger.Info("Closing connection to datastore")
			err := ds.Close()
			if err != nil {
				logger.Info("Unable to properly shutdown datastore", "error", err)
				os.Exit(1)
			}
			os.Exit(0)
		}
	}()
}

func initNetboxProxy() {
	nbproxy = netbox.New()
}

func initEventBus() {
	for {
		nsqd := viper.GetString("nsqd-addr")
		httpnsqd := viper.GetString("nsqd-http-addr")
		p, err := bus.NewPublisher(zapup.MustRootLogger(), nsqd, httpnsqd)
		if err != nil {
			logger.Errorw("cannot create nsq publisher", "error", err)
			time.Sleep(3 * time.Second)
			continue
		}
		logger.Infow("nsq connected", "nsqd", nsqd)
		if err := p.CreateTopic(string(metal.TopicDevice)); err != nil {
			logger.Errorw("cannot create TopicDevice", "error", err)
			time.Sleep(3 * time.Second)
			continue
		}
		producer = p
		break
	}
}

func initDataStore() {
	dbAdapter := viper.GetString("db")
	if dbAdapter == "rethinkdb" {
		ds = datastore.New(
			logger.Desugar(),
			viper.GetString("db-addr"),
			viper.GetString("db-name"),
			viper.GetString("db-user"),
			viper.GetString("db-password"),
		)
	} else {
		logger.Error("database not supported", "db", dbAdapter)
	}

	err := ds.Connect()

	if err != nil {
		logger.Error("cannot connect to db in root command metal-api/internal/main.initDatastore()", "error", err)
		panic(err)
	}

}

func initRestServices() *restfulspec.Config {
	lg := logger.Desugar()
	restful.DefaultContainer.Add(service.NewSite(lg, ds))
	restful.DefaultContainer.Add(service.NewImage(lg, ds))
	restful.DefaultContainer.Add(service.NewSize(lg, ds))
	restful.DefaultContainer.Add(service.NewDevice(lg, ds, producer, nbproxy))
	restful.DefaultContainer.Add(service.NewSwitch(lg, ds, nbproxy))
	restful.DefaultContainer.Add(rest.NewHealth(lg))
	restful.DefaultContainer.Add(rest.NewVersion(moduleName))
	restful.DefaultContainer.Filter(utils.RestfulLogger(lg, debug))

	config := restfulspec.Config{
		WebServices:                   restful.RegisteredWebServices(), // you control what services are visible
		APIPath:                       "/apidocs.json",
		PostBuildSwaggerObjectHandler: enrichSwaggerObject}
	restful.DefaultContainer.Add(restfulspec.NewOpenAPIService(config))
	return &config
}

func dumpSwaggerJSON() {
	cfg := initRestServices()
	actual := restfulspec.BuildSwagger(*cfg)
	js, err := json.MarshalIndent(actual, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", js)
}

func run() {
	initRestServices()

	// enable OPTIONS-request so clients can query CORS information
	restful.DefaultContainer.Filter(restful.DefaultContainer.OPTIONSFilter)

	// enable CORS for the UI to work.
	// if we will add support for api-tokens as headers, we had to add them
	// here to. note: the token's should not contain the product (aka. metal)
	// because customers should have ONE token for many products.
	// ExposeHeaders:  []string{"X-FITS-TOKEN"},
	cors := restful.CrossOriginResourceSharing{
		AllowedHeaders: []string{"Content-Type", "Accept", "Authorization"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		CookiesAllowed: false,
		Container:      restful.DefaultContainer}
	restful.DefaultContainer.Filter(cors.Filter)

	addr := fmt.Sprintf("%s:%d", viper.GetString("bind-addr"), viper.GetInt("port"))
	logger.Info("start metal api", "version", version.V.String())
	http.ListenAndServe(addr, nil)
}

func enrichSwaggerObject(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       moduleName,
			Description: "Resource for managing pure metal",
			Contact: &spec.ContactInfo{
				Name:  "Devops Team",
				Email: "devops@f-i-ts.de",
				URL:   "http://www.f-i-ts.de",
			},
			License: &spec.License{
				Name: "MIT",
				URL:  "http://mit.org",
			},
			Version: "1.0.0",
		},
	}
	swo.Tags = []spec.Tag{
		spec.Tag{TagProps: spec.TagProps{
			Name:        "image",
			Description: "Managing image entities"}},
		spec.Tag{TagProps: spec.TagProps{
			Name:        "size",
			Description: "Managing size entities"}},
		spec.Tag{TagProps: spec.TagProps{
			Name:        "device",
			Description: "Managing devices"}},
		spec.Tag{TagProps: spec.TagProps{
			Name:        "switch",
			Description: "Managing switches"}},
	}
}
