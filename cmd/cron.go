package cmd

import (
	"encoding/json"
	"log"
	"time"

	"github.com/neilli-sable/sideupload/infrastructure/adaptor"
	"github.com/neilli-sable/sideupload/infrastructure/setting"
	cron "github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// cronCmd represents the cron command
var cronCmd = &cobra.Command{
	Use:   "cron",
	Short: "upstart cron mode",
	Long:  `upstart cron mode.`,
	Run: func(cmd *cobra.Command, args []string) {
		opt := &setting.Setting{}
		if err := viper.Unmarshal(opt); err != nil {
			panic(err)
		}

		// start log
		log.Println("startup", "appName", cmd.Use)
		bytes, _ := json.Marshal(opt)
		log.Println(string(bytes))

		cronStart(opt)
		for {
			time.Sleep(10000000000000)
		}
	},
}

func cronStart(opt *setting.Setting) {
	c := cron.New(cron.WithSeconds())

	c.AddFunc(opt.Sideupload.CronWithSecond, func() {
		log.Printf("cron start at %v", time.Now())

		usecase := adaptor.UsecaseFactory(opt)
		defer usecase.Close()

		targets, err := usecase.GetTargets(opt.Sideupload.TargetDir)
		if err != nil {
			panic(err)
		}
		log.Printf("get %d targets", len(targets))

		archives, err := usecase.CompressTargets(targets)
		if err != nil {
			panic(err)
		}
		log.Printf("ready %d archives", len(archives))

		err = usecase.UploadArchives(archives)
		if err != nil {
			panic(err)
		}
		log.Printf("upload done!!")

		count, err := usecase.DeleteOldArchives()
		if err != nil {
			panic(err)
		}
		log.Printf("deleted %d files", count)
	})

	c.Start()
}

func init() {
	rootCmd.AddCommand(cronCmd)
}
