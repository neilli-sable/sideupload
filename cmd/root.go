package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/neilli-sable/sideupload/infrastructure/setting"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "sideupload",
	Short: "sideupload is mock server that something response",
	Long:  `sideupload is mock server that something response`,
	Run: func(cmd *cobra.Command, args []string) {
		opt := &setting.Setting{}
		if err := viper.Unmarshal(opt); err == nil {
			fmt.Println(opt)
		}

		// start log
		log.Println("startup", "appName", cmd.Use)
		bytes, _ := json.Marshal(opt)
		log.Println(string(bytes))

		serverStart(opt)
	},
}

// Execute コマンド実行
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func serverStart(opt *setting.Setting) {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			encoder := json.NewEncoder(w)
			encoder.Encode("OK")
		})
	})

	err := http.ListenAndServe(fmt.Sprintf(":%d", opt.Sideupload.ListenPort), r)
	if err != nil {
		panic(err.Error())
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.sideupload.toml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		dir, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(dir)
		viper.SetConfigName("sideupload")
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
