// Copyright © 2020 liut <liutao@liut.cc>
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
	"log"

	"github.com/spf13/cobra"

	"github.com/liut/staffio/pkg/backends"
	"github.com/liut/staffio/pkg/models/oauth"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "client add or revoke",
	Long: `
Client add or revoke or list

add client:
	--add name [--uri return URI]
revoke client:
	--revoke code
list:
	--list [--limit=5]
	`,
	Run: clientRun,
}

func init() {
	RootCmd.AddCommand(clientCmd)

	clientCmd.Flags().Bool("list", false, "List clients")
	clientCmd.Flags().String("add", "", "client name")
	clientCmd.Flags().String("uri", "", "client return URI")
	clientCmd.Flags().String("revoke", "", "A client Code(ID)")
	clientCmd.Flags().Int("limit", 5, "list limit")
}

func clientRun(cmd *cobra.Command, args []string) {
	name, _ := cmd.Flags().GetString("add")
	uri, _ := cmd.Flags().GetString("uri")
	service := backends.NewService()

	if revoke, err := cmd.Flags().GetString("revoke"); err == nil && revoke != "" {
		if err = service.OSIN().RemoveClient(revoke); err != nil {
			log.Printf("remove client failed, err %s", err)
			return
		}
		fmt.Printf("remove client %q OK\n", revoke)
	}

	if list, err := cmd.Flags().GetBool("list"); err == nil && list {
		limit, _ := cmd.Flags().GetInt("limit")
		if limit < 1 || limit > 50 {
			limit = 50
		}
		spec := &oauth.ClientSpec{Limit: limit}
		data, err := service.OSIN().LoadClients(spec)
		if err != nil {
			log.Printf("list failed, err %s", err)
			return
		}
		for _, c := range data {
			fmt.Printf("client % 19q (id %q, secret %q)\n", c.GetName(), c.GetId(), c.GetSecret())
		}
	}

	if name != "" {
		client := backends.GenNewClient(name, uri)
		fmt.Printf("generate new client %s, id %s, secret %s \n", client.GetName(), client.GetId(), client.GetSecret())
		if err := service.OSIN().SaveClient(client); err != nil {
			log.Printf("save client failed, error %s", err)
			return
		}
		fmt.Print("save client done\n")
	}

}
