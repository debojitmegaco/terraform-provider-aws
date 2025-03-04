package wafregional_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/service/wafregional"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccWAFRegionalWebACLDataSource_basic(t *testing.T) {
	ctx := acctest.Context(t)
	name := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_wafregional_web_acl.web_acl"
	datasourceName := "data.aws_wafregional_web_acl.web_acl"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, wafregional.EndpointsID) },
		ErrorCheck:               acctest.ErrorCheck(t, wafregional.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccWebACLDataSourceConfig_nonExistent,
				ExpectError: regexp.MustCompile(`web ACLs not found`),
			},
			{
				Config: testAccWebACLDataSourceConfig_name(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(datasourceName, "name", resourceName, "name"),
				),
			},
		},
	})
}

func testAccWebACLDataSourceConfig_name(name string) string {
	return fmt.Sprintf(`
resource "aws_wafregional_web_acl" "web_acl" {
  name        = %[1]q
  metric_name = "tfWebACL"

  default_action {
    type = "ALLOW"
  }
}

data "aws_wafregional_web_acl" "web_acl" {
  name = aws_wafregional_web_acl.web_acl.name
}
`, name)
}

const testAccWebACLDataSourceConfig_nonExistent = `
data "aws_wafregional_web_acl" "web_acl" {
  name = "tf-acc-test-does-not-exist"
}
`
