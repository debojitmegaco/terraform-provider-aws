package athena_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/athena"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func TestAccAthenaNamedQuery_basic(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_athena_named_query.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, athena.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNamedQueryDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccNamedQueryConfig_basic(sdkacctest.RandInt(), sdkacctest.RandString(5)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNamedQueryExists(ctx, resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAthenaNamedQuery_withWorkGroup(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_athena_named_query.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, athena.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNamedQueryDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccNamedQueryConfig_workGroup(sdkacctest.RandInt(), sdkacctest.RandString(5)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNamedQueryExists(ctx, resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckNamedQueryDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).AthenaConn(ctx)
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_athena_named_query" {
				continue
			}

			input := &athena.GetNamedQueryInput{
				NamedQueryId: aws.String(rs.Primary.ID),
			}

			resp, err := conn.GetNamedQueryWithContext(ctx, input)
			if err != nil {
				if tfawserr.ErrMessageContains(err, athena.ErrCodeInvalidRequestException, rs.Primary.ID) {
					return nil
				}
				return err
			}
			if resp.NamedQuery != nil {
				return fmt.Errorf("Athena Named Query (%s) found", rs.Primary.ID)
			}
		}
		return nil
	}
}

func testAccCheckNamedQueryExists(ctx context.Context, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).AthenaConn(ctx)

		input := &athena.GetNamedQueryInput{
			NamedQueryId: aws.String(rs.Primary.ID),
		}

		_, err := conn.GetNamedQueryWithContext(ctx, input)
		return err
	}
}

func testAccNamedQueryConfig_basic(rInt int, rName string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket        = "tf-test-athena-db-%s-%d"
  force_destroy = true
}

resource "aws_athena_database" "test" {
  name   = "%s"
  bucket = aws_s3_bucket.test.bucket
}

resource "aws_athena_named_query" "test" {
  name        = "tf-athena-named-query-%s"
  database    = aws_athena_database.test.name
  query       = "SELECT * FROM ${aws_athena_database.test.name} limit 10;"
  description = "tf test"
}
`, rName, rInt, rName, rName)
}

func testAccNamedQueryConfig_workGroup(rInt int, rName string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket        = "tf-test-athena-db-%s-%d"
  force_destroy = true
}

resource "aws_athena_workgroup" "test" {
  name = "tf-athena-workgroup-%s-%d"
}

resource "aws_athena_database" "test" {
  name   = "%s"
  bucket = aws_s3_bucket.test.bucket
}

resource "aws_athena_named_query" "test" {
  name        = "tf-athena-named-query-%s"
  workgroup   = aws_athena_workgroup.test.id
  database    = aws_athena_database.test.name
  query       = "SELECT * FROM ${aws_athena_database.test.name} limit 10;"
  description = "tf test"
}
`, rName, rInt, rName, rInt, rName, rName)
}
