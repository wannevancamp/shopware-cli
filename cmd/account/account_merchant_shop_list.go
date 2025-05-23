package account

import (
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/shopware/shopware-cli/internal/table"
)

var accountCompanyMerchantShopListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all shops",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, _ []string) error {
		table := table.NewWriter(os.Stdout)
		table.Header([]string{"ID", "Domain", "Usage"})

		shops, err := services.AccountClient.Merchant().Shops(cmd.Context())
		if err != nil {
			return err
		}

		for _, shop := range shops {
			_ = table.Append([]string{
				strconv.FormatInt(int64(shop.Id), 10),
				shop.Domain,
				shop.Environment.Name,
			})
		}

		_ = table.Render()

		return nil
	},
}

func init() {
	accountCompanyMerchantShopCmd.AddCommand(accountCompanyMerchantShopListCmd)
}
