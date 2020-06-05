package cmd

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"

	"github.com/open-policy-agent/contrib/data_filter_mongodb/internal/mongo"
	"github.com/open-policy-agent/contrib/data_filter_mongodb/internal/opa"
	"github.com/spf13/cobra"
)

var policyFile string
var configFile string

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the Rego to MongoDB query translator",
	Long:  `Run the Rego to MongoDB query translator. Make sure mongo db is already running`,
	Run: func(cmd *cobra.Command, args []string) {
		RunMongoOPA()
	},
}

func RunMongoOPA() {
	var err error
	client := &mongo.Mongo{}
	client.ClientOptions, client.Database = GetDBClient(configFile)
	client.Logger, err = zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed %v", err)
	}

	client.Mongo, err = client.CreateConnection()
	if err != nil {
		log.Fatalf("failed to establish connection with mongo db %s", err)
	}

	client.CreateTestData()
	err = opa.New(client, policyFile).Run(context.Background())
	if err != nil {
		log.Fatalf("failed to create mongo opa server %s", err)
	}
}

func GetDBClient(configFile string) (*options.ClientOptions, string) {
	var cfg server
	var params []string
	params = append(params, configFile)
	_, expandENV := parseConfigFileParameter(params)
	if configFile != "" {
		if err := LoadConfig(configFile, expandENV, &cfg); err != nil {
			fmt.Fprintf(os.Stderr, "error loading config from %s: %v\n", configFile, err)
			os.Exit(1)
		}
	}

	mongodbURI := "mongodb://" + cfg.Cfg.Address
	var clientOptions *options.ClientOptions
	if cfg.Cfg.Password == "" && cfg.Cfg.Username == "" {
		clientOptions = options.Client().ApplyURI(mongodbURI)
	} else {
		credentials := options.Credential{
			Username: cfg.Cfg.Username,
			Password: cfg.Cfg.Password,
		}
		clientOptions = options.Client().ApplyURI(mongodbURI).SetAuth(credentials)
	}
	return clientOptions, cfg.Cfg.Database
}

// Parse -config.file and -config.expand-env option via separate flag set, to avoid polluting default one and calling flag.Parse on it twice.
func parseConfigFileParameter(args []string) (configFile string, expandEnv bool) {
	// ignore errors and any output here. Any flag errors will be reported by main flag.Parse() call.
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.SetOutput(ioutil.Discard)

	// usage not used in these functions.
	fs.StringVar(&configFile, configFileOption, "", "")

	// Try to find -config.file and -config.expand-env option in the flags. As Parsing stops on the first error, eg. unknown flag, we simply
	// try remaining parameters until we find config flag, or there are no params left.
	// (ContinueOnError just means that flag.Parse doesn't call panic or os.Exit, but it returns error, which we ignore)
	for len(args) > 0 {
		_ = fs.Parse(args)
		args = args[1:]
	}

	return
}

// LoadConfig read YAML-formatted config from filename into cfg.
func LoadConfig(filename string, expandENV bool, cfg *server) error {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.Wrap(err, "Error reading config file")
	}

	if expandENV {
		buf = expandEnv(buf)
	}

	err = yaml.UnmarshalStrict(buf, cfg)
	if err != nil {
		return errors.Wrap(err, "Error parsing config file")
	}

	return nil
}

const (
	configFileOption = "config.file"
)

type server struct {
	Cfg Config `yaml:"cfg,omitempty"`
}

type Config struct {
	Address  string `yaml:"address,omitempty"`
	Database string `yaml:"database,omitempty"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}

// expandEnv replaces ${var} or $var in config according to the values of the current environment variables.
// The replacement is case-sensitive. References to undefined variables are replaced by the empty string.
// A default value can be given by using the form ${var:default value}.
func expandEnv(config []byte) []byte {
	return []byte(os.Expand(string(config), func(key string) string {
		keyAndDefault := strings.SplitN(key, ":", 2)
		key = keyAndDefault[0]

		v := os.Getenv(key)
		if v == "" && len(keyAndDefault) == 2 {
			v = keyAndDefault[1] // Set value to the default.
		}
		return v
	}))
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().StringVarP(&policyFile, "policy.file", "p", "", "set path of OPA policy.")
	runCmd.PersistentFlags().StringVarP(&configFile, "config.file", "c", "", "set path of configuration file.")
}
