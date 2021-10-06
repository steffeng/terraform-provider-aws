package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/efs"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/aws/internal/service/efs/finder"
	"github.com/hashicorp/terraform-provider-aws/aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/provider"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	tfefs "github.com/hashicorp/terraform-provider-aws/internal/service/efs"
	tfefs "github.com/hashicorp/terraform-provider-aws/internal/service/efs"
	tfefs "github.com/hashicorp/terraform-provider-aws/internal/service/efs"
)

func TestAccAWSEFSBackupPolicy_basic(t *testing.T) {
	var v efs.BackupPolicy
	resourceName := "aws_efs_backup_policy.test"
	rName := sdkacctest.RandomWithPrefix("tf-acc-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, efs.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckEfsBackupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSEFSBackupPolicyConfig(rName, "ENABLED"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEFSBackupPolicyExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "backup_policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "backup_policy.0.status", "ENABLED"),
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

func TestAccAWSEFSBackupPolicy_disappears_fs(t *testing.T) {
	var v efs.BackupPolicy
	resourceName := "aws_efs_backup_policy.test"
	fsResourceName := "aws_efs_file_system.test"
	rName := sdkacctest.RandomWithPrefix("tf-acc-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, efs.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckEfsBackupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSEFSBackupPolicyConfig(rName, "ENABLED"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEFSBackupPolicyExists(resourceName, &v),
					acctest.CheckResourceDisappears(acctest.Provider, ResourceFileSystem(), fsResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAWSEFSBackupPolicy_update(t *testing.T) {
	var v efs.BackupPolicy
	resourceName := "aws_efs_backup_policy.test"
	rName := sdkacctest.RandomWithPrefix("tf-acc-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, efs.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckEfsBackupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSEFSBackupPolicyConfig(rName, "DISABLED"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEFSBackupPolicyExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "backup_policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "backup_policy.0.status", "DISABLED"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAWSEFSBackupPolicyConfig(rName, "ENABLED"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEFSBackupPolicyExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "backup_policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "backup_policy.0.status", "ENABLED"),
				),
			},
			{
				Config: testAccAWSEFSBackupPolicyConfig(rName, "DISABLED"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEFSBackupPolicyExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "backup_policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "backup_policy.0.status", "DISABLED"),
				),
			},
		},
	})
}

func testAccCheckEFSBackupPolicyExists(name string, v *efs.BackupPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).EFSConn

		output, err := tfefs.FindBackupPolicyByID(conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccCheckEfsBackupPolicyDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).EFSConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_efs_backup_policy" {
			continue
		}

		output, err := tfefs.FindBackupPolicyByID(conn, rs.Primary.ID)

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		if aws.StringValue(output.Status) == efs.StatusDisabled {
			continue
		}

		return fmt.Errorf("Transfer Server %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccAWSEFSBackupPolicyConfig(rName, status string) string {
	return fmt.Sprintf(`
resource "aws_efs_file_system" "test" {
  creation_token = %[1]q
}

resource "aws_efs_backup_policy" "test" {
  file_system_id = aws_efs_file_system.test.id

  backup_policy {
    status = %[2]q
  }
}
`, rName, status)
}