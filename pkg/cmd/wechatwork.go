// Copyright © 2019 liut <liutao@liut.cc>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/liut/staffio/pkg/backends/wechatwork"
)

// wechatworkCmd represents the wechatwork command
var wechatworkCmd = &cobra.Command{
	Use:   "wechatwork",
	Short: "Sync with wechat work",
	Run:   wechatworkRun,
}

func init() {
	RootCmd.AddCommand(wechatworkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// wechatworkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// wechatworkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	wechatworkCmd.Flags().String("act", "", "action: query | sync | sync-all")
	wechatworkCmd.Flags().String("uid", "", "query uid")
}

func wechatworkRun(cmd *cobra.Command, args []string) {
	action, _ := cmd.Flags().GetString("act")
	uid, _ := cmd.Flags().GetString("uid")

	wechatwork.SyncDepartment(action, uid)
}
