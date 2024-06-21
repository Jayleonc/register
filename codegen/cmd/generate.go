package cmd

import (
	"fmt"
	"github.com/Jayleonc/register/codegen/internal"
	"github.com/spf13/cobra"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
)

var (
	serviceName   string
	outputPath    string
	etcdEndpoints string
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate client code",
	Long:  `Generate client code for a specified service registered in etcd.`,
	Run: func(cmd *cobra.Command, args []string) {
		var endpoints []string
		if etcdEndpoints != "" {
			endpoints = append(endpoints, etcdEndpoints)
		}

		fmt.Println("endpoints:", endpoints)

		client, err := clientv3.New(clientv3.Config{
			Endpoints: endpoints,
		})
		if err != nil {
			log.Fatalf("Failed to create etcd client: %v", err)
		}

		// 如果没有提供输出路径，则使用默认值
		if outputPath == "" {
			outputPath = fmt.Sprintf("./internal/client/%s", serviceName)
		}

		err = internal.GenerateClientCode(serviceName, outputPath, client)
		if err != nil {
			log.Fatalf("Failed to generate client code: %v", err)
		}
		fmt.Println("Client code generated successfully.")
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringVar(&serviceName, "service", "", "Name of the service")
	generateCmd.Flags().StringVar(&outputPath, "output", "", "Output path for the generated code")
	generateCmd.Flags().StringVar(&etcdEndpoints, "etcd-endpoints", "localhost:2379", "Etcd endpoints")

	generateCmd.MarkFlagRequired("service")
}
