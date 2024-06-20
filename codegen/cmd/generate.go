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
	baseURL       string
	outputPath    string
	etcdEndpoints []string
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate client code",
	Long:  `Generate client code for a specified service registered in etcd.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := clientv3.New(clientv3.Config{
			Endpoints: etcdEndpoints,
		})
		if err != nil {
			log.Fatalf("Failed to create etcd client: %v", err)
		}

		err = internal.GenerateClientCode(serviceName, baseURL, outputPath, client)
		if err != nil {
			log.Fatalf("Failed to generate client code: %v", err)
		}
		fmt.Println("Client code generated successfully.")
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringVar(&serviceName, "service", "", "Name of the service")
	generateCmd.Flags().StringVar(&baseURL, "base-url", "", "Base URL of the service")
	generateCmd.Flags().StringVar(&outputPath, "output", "", "Output path for the generated code")
	generateCmd.Flags().StringSliceVar(&etcdEndpoints, "etcd-endpoints", []string{"localhost:2379"}, "Etcd endpoints")

	generateCmd.MarkFlagRequired("service")
	generateCmd.MarkFlagRequired("base-url")
	generateCmd.MarkFlagRequired("output")
}
