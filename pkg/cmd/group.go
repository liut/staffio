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
)

// groupCmd represents the group command
var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "Operate a group",
	Long:  `Add or Kick a user into a group, new group will be create`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.ParseFlags(args)
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			fmt.Println("empty group name")
			return
		}
		svc := backends.NewService()
		group, err := svc.GetGroup(name)
		if err != nil && err != backends.ErrStoreNotFound {
			fmt.Printf("get group ERR %s\n", err)
			return
		}

		if uid, _ := cmd.Flags().GetString("add-member"); uid != "" {
			if err == backends.ErrStoreNotFound {
				group = &backends.Group{
					Name:    name,
					Members: []string{uid},
				}
			} else {
				group.Members = append(group.Members, uid)
			}
			err = svc.SaveGroup(group)
			if err != nil {
				fmt.Printf("save group ERR %s\n", err)
			} else {
				fmt.Println("save group OK")
			}
			return
		}

		if uid, _ := cmd.Flags().GetString("kick-member"); uid != "" && err == nil {
			var members []string
			for _, m := range group.Members {
				if m != uid {
					members = append(members, m)
				}
			}
			if len(members) == len(group.Members) {
				return
			}
			group.Members = members
			err = svc.SaveGroup(group)
			if err != nil {
				fmt.Printf("save group ERR %s\n", err)
			} else {
				fmt.Println("save group OK")
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(groupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// groupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	groupCmd.Flags().StringP("name", "g", "", "Group name")
	groupCmd.Flags().StringP("add-member", "a", "", "UID of member will add")
	groupCmd.Flags().StringP("kick-member", "t", "", "UID of member will kick")
	addstaffCmd.MarkFlagRequired("name") //nolint
}
