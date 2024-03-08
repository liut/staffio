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
	"fmt"

	"github.com/spf13/cobra"

	"github.com/liut/staffio/pkg/backends"
	"github.com/liut/staffio/pkg/models"
)

// addstaffCmd represents the addstaff command
var addstaffCmd = &cobra.Command{
	Use:   "addstaff",
	Short: "Add a simple user for develop",
	Long:  `Add a simple user for develop`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.ParseFlags(args)
		uid, _ := cmd.Flags().GetString("uid")
		password, _ := cmd.Flags().GetString("password")
		if uid == "" {
			fmt.Println("empty uid")
			return
		}
		cn, _ := cmd.Flags().GetString("name")
		if cn == "" {
			cn = uid
		}
		sn, _ := cmd.LocalFlags().GetString("sn")
		if sn == "" {
			sn = cn
		}
		svc := backends.NewService()
		// check ready
		if err := svc.Ready(); err != nil {
			panic(err)
		}

		staff := &models.Staff{
			UID:        uid,
			CommonName: cn,
			Surname:    sn,
		}
		isNew, err := svc.Save(staff)
		if err != nil {
			logger().Warnw("save staff fail", "staff", staff, "err", err)
			return
		}
		logger().Infow("saved", "staff", staff)
		if isNew && password != "" {
			err = svc.PasswordReset(uid, password)
			if err != nil {
				logger().Infow("reset password fail", "uid", uid, "err", err)
				return
			}
			logger().Infow("reset password OK", "uid", uid)
		}
	},
}

func init() {
	RootCmd.AddCommand(addstaffCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addstaffCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	addstaffCmd.Flags().StringP("uid", "u", "", "staff uid")
	addstaffCmd.Flags().StringP("password", "p", "", "staff password")
	addstaffCmd.Flags().StringP("name", "n", "", "staff name")
	addstaffCmd.Flags().String("sn", "", "surname")
	addstaffCmd.MarkFlagRequired("uid")  //nolint
	addstaffCmd.MarkFlagRequired("name") //nolint
	addstaffCmd.MarkFlagRequired("sn")   //nolint
}
