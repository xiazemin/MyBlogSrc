/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

// servStaticCmd represents the servStatic command
var servStaticCmd = &cobra.Command{
	Use:   "servStatic",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("servStatic called")
		http.Handle("/MyBlog/", http.StripPrefix("/MyBlog/", http.FileServer(http.Dir("/Users/xiazemin/MyBlogSrc/_site/"))))
		fmt.Println("listen at:http://127.0.0.1:4000/MyBlog/")
		if err := http.ListenAndServe(":4000", nil); err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(servStaticCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// servStaticCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// servStaticCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
