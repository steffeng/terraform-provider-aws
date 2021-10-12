package aws

import (
	"fmt"
	"log"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/imagebuilder"
	"github.com/hashicorp/aws-sdk-go-base/tfawserr"
	"github.com/hashicorp/go-multierror"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/provider"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep"
)

func init() {
	resource.AddTestSweepers("aws_imagebuilder_image_pipeline", &resource.Sweeper{
		Name: "aws_imagebuilder_image_pipeline",
		F:    testSweepImageBuilderImagePipelines,
	})
}

func testSweepImageBuilderImagePipelines(region string) error {
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}
	conn := client.(*conns.AWSClient).ImageBuilderConn

	var sweeperErrs *multierror.Error

	input := &imagebuilder.ListImagePipelinesInput{}

	err = conn.ListImagePipelinesPages(input, func(page *imagebuilder.ListImagePipelinesOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, imagePipeline := range page.ImagePipelineList {
			if imagePipeline == nil {
				continue
			}

			arn := aws.StringValue(imagePipeline.Arn)

			r := ResourceImagePipeline()
			d := r.Data(nil)
			d.SetId(arn)

			err := r.Delete(d, client)

			if err != nil {
				sweeperErr := fmt.Errorf("error deleting Image Builder Image Pipeline (%s): %w", arn, err)
				log.Printf("[ERROR] %s", sweeperErr)
				sweeperErrs = multierror.Append(sweeperErrs, sweeperErr)
				continue
			}
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping Image Builder Image Pipeline sweep for %s: %s", region, err)
		return sweeperErrs.ErrorOrNil() // In case we have completed some pages, but had errors
	}

	if err != nil {
		sweeperErrs = multierror.Append(sweeperErrs, fmt.Errorf("error listing Image Builder Image Pipelines: %w", err))
	}

	return sweeperErrs.ErrorOrNil()
}

func TestAccAwsImageBuilderImagePipeline_basic(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	imageRecipeResourceName := "aws_imagebuilder_image_recipe.test"
	infrastructureConfigurationResourceName := "aws_imagebuilder_infrastructure_configuration.test"
	resourceName := "aws_imagebuilder_image_pipeline.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, imagebuilder.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAwsImageBuilderImagePipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsImageBuilderImagePipelineConfigName(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsImageBuilderImagePipelineExists(resourceName),
					acctest.CheckResourceAttrRegionalARN(resourceName, "arn", "imagebuilder", fmt.Sprintf("image-pipeline/%s", rName)),
					acctest.CheckResourceAttrRFC3339(resourceName, "date_created"),
					resource.TestCheckResourceAttr(resourceName, "date_last_run", ""),
					resource.TestCheckResourceAttr(resourceName, "date_next_run", ""),
					acctest.CheckResourceAttrRFC3339(resourceName, "date_updated"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "distribution_configuration_arn", ""),
					resource.TestCheckResourceAttr(resourceName, "enhanced_image_metadata_enabled", "true"),
					resource.TestCheckResourceAttrPair(resourceName, "image_recipe_arn", imageRecipeResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "image_tests_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "image_tests_configuration.0.image_tests_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "image_tests_configuration.0.timeout_minutes", "720"),
					resource.TestCheckResourceAttrPair(resourceName, "infrastructure_configuration_arn", infrastructureConfigurationResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "platform", imagebuilder.PlatformLinux),
					resource.TestCheckResourceAttr(resourceName, "schedule.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "status", imagebuilder.PipelineStatusEnabled),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
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

func TestAccAwsImageBuilderImagePipeline_disappears(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_imagebuilder_image_pipeline.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, imagebuilder.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAwsImageBuilderImagePipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsImageBuilderImagePipelineConfigName(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsImageBuilderImagePipelineExists(resourceName),
					acctest.CheckResourceDisappears(acctest.Provider, ResourceImagePipeline(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAwsImageBuilderImagePipeline_Description(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_imagebuilder_image_pipeline.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, imagebuilder.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAwsImageBuilderImagePipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsImageBuilderImagePipelineConfigDescription(rName, "description1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsImageBuilderImagePipelineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", "description1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAwsImageBuilderImagePipelineConfigDescription(rName, "description2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsImageBuilderImagePipelineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", "description2"),
				),
			},
		},
	})
}

func TestAccAwsImageBuilderImagePipeline_DistributionConfigurationArn(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	distributionConfigurationResourceName := "aws_imagebuilder_distribution_configuration.test"
	resourceName := "aws_imagebuilder_image_pipeline.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, imagebuilder.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAwsImageBuilderImagePipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsImageBuilderImagePipelineConfigDistributionConfigurationArn1(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsImageBuilderImagePipelineExists(resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "distribution_configuration_arn", distributionConfigurationResourceName, "arn"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAwsImageBuilderImagePipelineConfigDistributionConfigurationArn2(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsImageBuilderImagePipelineExists(resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "distribution_configuration_arn", distributionConfigurationResourceName, "arn"),
				),
			},
		},
	})
}

func TestAccAwsImageBuilderImagePipeline_EnhancedImageMetadataEnabled(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_imagebuilder_image_pipeline.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, imagebuilder.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAwsImageBuilderImagePipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsImageBuilderImagePipelineConfigEnhancedImageMetadataEnabled(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsImageBuilderImagePipelineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enhanced_image_metadata_enabled", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAwsImageBuilderImagePipelineConfigEnhancedImageMetadataEnabled(rName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsImageBuilderImagePipelineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enhanced_image_metadata_enabled", "true"),
				),
			},
		},
	})
}

func TestAccAwsImageBuilderImagePipeline_ImageRecipeArn(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	imageRecipeResourceName := "aws_imagebuilder_image_recipe.test"
	imageRecipeResourceName2 := "aws_imagebuilder_image_recipe.test2"
	resourceName := "aws_imagebuilder_image_pipeline.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, imagebuilder.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAwsImageBuilderImagePipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsImageBuilderImagePipelineConfigName(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsImageBuilderImagePipelineExists(resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "image_recipe_arn", imageRecipeResourceName, "arn"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAwsImageBuilderImagePipelineConfigImageRecipeArn2(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsImageBuilderImagePipelineExists(resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "image_recipe_arn", imageRecipeResourceName2, "arn"),
				),
			},
		},
	})
}

func TestAccAwsImageBuilderImagePipeline_ImageTestsConfiguration_ImageTestsEnabled(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_imagebuilder_image_pipeline.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, imagebuilder.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAwsImageBuilderImagePipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsImageBuilderImagePipelineConfigImageTestsConfigurationImageTestsEnabled(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsImageBuilderImagePipelineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "image_tests_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "image_tests_configuration.0.image_tests_enabled", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAwsImageBuilderImagePipelineConfigImageTestsConfigurationImageTestsEnabled(rName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsImageBuilderImagePipelineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "image_tests_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "image_tests_configuration.0.image_tests_enabled", "true"),
				),
			},
		},
	})
}

func TestAccAwsImageBuilderImagePipeline_ImageTestsConfiguration_TimeoutMinutes(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_imagebuilder_image_pipeline.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, imagebuilder.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAwsImageBuilderImagePipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsImageBuilderImagePipelineConfigImageTestsConfigurationTimeoutMinutes(rName, 721),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsImageBuilderImagePipelineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "image_tests_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "image_tests_configuration.0.timeout_minutes", "721"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAwsImageBuilderImagePipelineConfigImageTestsConfigurationTimeoutMinutes(rName, 722),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsImageBuilderImagePipelineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "image_tests_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "image_tests_configuration.0.timeout_minutes", "722"),
				),
			},
		},
	})
}

func TestAccAwsImageBuilderImagePipeline_InfrastructureConfigurationArn(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	infrastructureConfigurationResourceName := "aws_imagebuilder_infrastructure_configuration.test"
	infrastructureConfigurationResourceName2 := "aws_imagebuilder_infrastructure_configuration.test2"
	resourceName := "aws_imagebuilder_image_pipeline.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, imagebuilder.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAwsImageBuilderImagePipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsImageBuilderImagePipelineConfigName(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsImageBuilderImagePipelineExists(resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "infrastructure_configuration_arn", infrastructureConfigurationResourceName, "arn"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAwsImageBuilderImagePipelineConfigInfrastructureConfigurationArn2(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsImageBuilderImagePipelineExists(resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "infrastructure_configuration_arn", infrastructureConfigurationResourceName2, "arn"),
				),
			},
		},
	})
}

func TestAccAwsImageBuilderImagePipeline_Schedule_PipelineExecutionStartCondition(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_imagebuilder_image_pipeline.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, imagebuilder.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAwsImageBuilderImagePipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsImageBuilderImagePipelineConfigSchedulePipelineExecutionStartCondition(rName, imagebuilder.PipelineExecutionStartConditionExpressionMatchAndDependencyUpdatesAvailable),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsImageBuilderImagePipelineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "schedule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "schedule.0.pipeline_execution_start_condition", imagebuilder.PipelineExecutionStartConditionExpressionMatchAndDependencyUpdatesAvailable),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAwsImageBuilderImagePipelineConfigSchedulePipelineExecutionStartCondition(rName, imagebuilder.PipelineExecutionStartConditionExpressionMatchOnly),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsImageBuilderImagePipelineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "schedule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "schedule.0.pipeline_execution_start_condition", imagebuilder.PipelineExecutionStartConditionExpressionMatchOnly),
				),
			},
		},
	})
}

func TestAccAwsImageBuilderImagePipeline_Schedule_ScheduleExpression(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_imagebuilder_image_pipeline.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, imagebuilder.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAwsImageBuilderImagePipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsImageBuilderImagePipelineConfigScheduleScheduleExpression(rName, "cron(1 0 * * ? *)"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsImageBuilderImagePipelineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "schedule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "schedule.0.schedule_expression", "cron(1 0 * * ? *)"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAwsImageBuilderImagePipelineConfigScheduleScheduleExpression(rName, "cron(2 0 * * ? *)"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsImageBuilderImagePipelineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "schedule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "schedule.0.schedule_expression", "cron(2 0 * * ? *)"),
				),
			},
		},
	})
}

func TestAccAwsImageBuilderImagePipeline_Status(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_imagebuilder_image_pipeline.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, imagebuilder.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAwsImageBuilderImagePipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsImageBuilderImagePipelineConfigStatus(rName, imagebuilder.PipelineStatusDisabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsImageBuilderImagePipelineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "status", imagebuilder.PipelineStatusDisabled),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAwsImageBuilderImagePipelineConfigStatus(rName, imagebuilder.PipelineStatusEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsImageBuilderImagePipelineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "status", imagebuilder.PipelineStatusEnabled),
				),
			},
		},
	})
}

func TestAccAwsImageBuilderImagePipeline_Tags(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_imagebuilder_image_pipeline.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, imagebuilder.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAwsImageBuilderImagePipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsImageBuilderImagePipelineConfigTags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsImageBuilderImagePipelineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAwsImageBuilderImagePipelineConfigTags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsImageBuilderImagePipelineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccAwsImageBuilderImagePipelineConfigTags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsImageBuilderImagePipelineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func testAccCheckAwsImageBuilderImagePipelineDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).ImageBuilderConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_imagebuilder_image_pipeline" {
			continue
		}

		input := &imagebuilder.GetImagePipelineInput{
			ImagePipelineArn: aws.String(rs.Primary.ID),
		}

		output, err := conn.GetImagePipeline(input)

		if tfawserr.ErrCodeEquals(err, imagebuilder.ErrCodeResourceNotFoundException) {
			continue
		}

		if err != nil {
			return fmt.Errorf("error getting Image Builder Image Pipeline (%s): %w", rs.Primary.ID, err)
		}

		if output != nil {
			return fmt.Errorf("Image Builder Image Pipeline (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckAwsImageBuilderImagePipelineExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).ImageBuilderConn

		input := &imagebuilder.GetImagePipelineInput{
			ImagePipelineArn: aws.String(rs.Primary.ID),
		}

		_, err := conn.GetImagePipeline(input)

		if err != nil {
			return fmt.Errorf("error getting Image Builder Image Pipeline (%s): %w", rs.Primary.ID, err)
		}

		return nil
	}
}

func testAccAwsImageBuilderImagePipelineConfigBase(rName string) string {
	return fmt.Sprintf(`
data "aws_region" "current" {}

data "aws_partition" "current" {}

resource "aws_iam_instance_profile" "test" {
  name = aws_iam_role.role.name
  role = aws_iam_role.role.name
}

resource "aws_iam_role" "role" {
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "ec2.${data.aws_partition.current.dns_suffix}"
      }
      Sid = ""
    }]
  })
  name = %[1]q
}

resource "aws_imagebuilder_component" "test" {
  data = yamlencode({
    phases = [{
      name = "build"
      steps = [{
        action = "ExecuteBash"
        inputs = {
          commands = ["echo 'hello world'"]
        }
        name      = "example"
        onFailure = "Continue"
      }]
    }]
    schemaVersion = 1.0
  })
  name     = %[1]q
  platform = "Linux"
  version  = "1.0.0"
}

resource "aws_imagebuilder_image_recipe" "test" {
  component {
    component_arn = aws_imagebuilder_component.test.arn
  }

  name         = %[1]q
  parent_image = "arn:${data.aws_partition.current.partition}:imagebuilder:${data.aws_region.current.name}:aws:image/amazon-linux-2-x86/x.x.x"
  version      = "1.0.0"
}

resource "aws_imagebuilder_infrastructure_configuration" "test" {
  instance_profile_name = aws_iam_instance_profile.test.name
  name                  = %[1]q
}
`, rName)
}

func testAccAwsImageBuilderImagePipelineConfigDescription(rName string, description string) string {
	return acctest.ConfigCompose(
		testAccAwsImageBuilderImagePipelineConfigBase(rName),
		fmt.Sprintf(`
resource "aws_imagebuilder_image_pipeline" "test" {
  description                      = %[2]q
  image_recipe_arn                 = aws_imagebuilder_image_recipe.test.arn
  infrastructure_configuration_arn = aws_imagebuilder_infrastructure_configuration.test.arn
  name                             = %[1]q
}
`, rName, description))
}

func testAccAwsImageBuilderImagePipelineConfigDistributionConfigurationArn1(rName string) string {
	return acctest.ConfigCompose(
		testAccAwsImageBuilderImagePipelineConfigBase(rName),
		fmt.Sprintf(`
resource "aws_imagebuilder_distribution_configuration" "test" {
  name = "%[1]s-1"

  distribution {
    ami_distribution_configuration {
      name = "{{ imagebuilder:buildDate }}"
    }

    region = data.aws_region.current.name
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_imagebuilder_image_pipeline" "test" {
  distribution_configuration_arn   = aws_imagebuilder_distribution_configuration.test.arn
  image_recipe_arn                 = aws_imagebuilder_image_recipe.test.arn
  infrastructure_configuration_arn = aws_imagebuilder_infrastructure_configuration.test.arn
  name                             = %[1]q
}
`, rName))
}

func testAccAwsImageBuilderImagePipelineConfigDistributionConfigurationArn2(rName string) string {
	return acctest.ConfigCompose(
		testAccAwsImageBuilderImagePipelineConfigBase(rName),
		fmt.Sprintf(`
resource "aws_imagebuilder_distribution_configuration" "test" {
  name = "%[1]s-2"

  distribution {
    ami_distribution_configuration {
      name = "{{ imagebuilder:buildDate }}"
    }

    region = data.aws_region.current.name
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_imagebuilder_image_pipeline" "test" {
  distribution_configuration_arn   = aws_imagebuilder_distribution_configuration.test.arn
  image_recipe_arn                 = aws_imagebuilder_image_recipe.test.arn
  infrastructure_configuration_arn = aws_imagebuilder_infrastructure_configuration.test.arn
  name                             = %[1]q
}
`, rName))
}

func testAccAwsImageBuilderImagePipelineConfigEnhancedImageMetadataEnabled(rName string, enhancedImageMetadataEnabled bool) string {
	return acctest.ConfigCompose(
		testAccAwsImageBuilderImagePipelineConfigBase(rName),
		fmt.Sprintf(`
resource "aws_imagebuilder_image_pipeline" "test" {
  enhanced_image_metadata_enabled  = %[2]t
  image_recipe_arn                 = aws_imagebuilder_image_recipe.test.arn
  infrastructure_configuration_arn = aws_imagebuilder_infrastructure_configuration.test.arn
  name                             = %[1]q
}
`, rName, enhancedImageMetadataEnabled))
}

func testAccAwsImageBuilderImagePipelineConfigImageRecipeArn2(rName string) string {
	return acctest.ConfigCompose(
		testAccAwsImageBuilderImagePipelineConfigBase(rName),
		fmt.Sprintf(`
resource "aws_imagebuilder_image_recipe" "test2" {
  component {
    component_arn = aws_imagebuilder_component.test.arn
  }

  name         = "%[1]s-2"
  parent_image = "arn:${data.aws_partition.current.partition}:imagebuilder:${data.aws_region.current.name}:aws:image/amazon-linux-2-x86/x.x.x"
  version      = "1.0.0"
}

resource "aws_imagebuilder_image_pipeline" "test" {
  image_recipe_arn                 = aws_imagebuilder_image_recipe.test2.arn
  infrastructure_configuration_arn = aws_imagebuilder_infrastructure_configuration.test.arn
  name                             = %[1]q
}
`, rName))
}

func testAccAwsImageBuilderImagePipelineConfigImageTestsConfigurationImageTestsEnabled(rName string, imageTestsEnabled bool) string {
	return acctest.ConfigCompose(
		testAccAwsImageBuilderImagePipelineConfigBase(rName),
		fmt.Sprintf(`
resource "aws_imagebuilder_image_pipeline" "test" {
  image_recipe_arn                 = aws_imagebuilder_image_recipe.test.arn
  infrastructure_configuration_arn = aws_imagebuilder_infrastructure_configuration.test.arn
  name                             = %[1]q

  image_tests_configuration {
    image_tests_enabled = %[2]t
  }
}
`, rName, imageTestsEnabled))
}

func testAccAwsImageBuilderImagePipelineConfigImageTestsConfigurationTimeoutMinutes(rName string, timeoutMinutes int) string {
	return acctest.ConfigCompose(
		testAccAwsImageBuilderImagePipelineConfigBase(rName),
		fmt.Sprintf(`
resource "aws_imagebuilder_image_pipeline" "test" {
  image_recipe_arn                 = aws_imagebuilder_image_recipe.test.arn
  infrastructure_configuration_arn = aws_imagebuilder_infrastructure_configuration.test.arn
  name                             = %[1]q

  image_tests_configuration {
    timeout_minutes = %[2]d
  }
}
`, rName, timeoutMinutes))
}

func testAccAwsImageBuilderImagePipelineConfigInfrastructureConfigurationArn2(rName string) string {
	return acctest.ConfigCompose(
		testAccAwsImageBuilderImagePipelineConfigBase(rName),
		fmt.Sprintf(`
resource "aws_imagebuilder_infrastructure_configuration" "test2" {
  instance_profile_name = aws_iam_instance_profile.test.name
  name                  = "%[1]s-2"
}

resource "aws_imagebuilder_image_pipeline" "test" {
  image_recipe_arn                 = aws_imagebuilder_image_recipe.test.arn
  infrastructure_configuration_arn = aws_imagebuilder_infrastructure_configuration.test2.arn
  name                             = %[1]q
}
`, rName))
}

func testAccAwsImageBuilderImagePipelineConfigName(rName string) string {
	return acctest.ConfigCompose(
		testAccAwsImageBuilderImagePipelineConfigBase(rName),
		fmt.Sprintf(`
resource "aws_imagebuilder_image_pipeline" "test" {
  image_recipe_arn                 = aws_imagebuilder_image_recipe.test.arn
  infrastructure_configuration_arn = aws_imagebuilder_infrastructure_configuration.test.arn
  name                             = %[1]q
}
`, rName))
}

func testAccAwsImageBuilderImagePipelineConfigSchedulePipelineExecutionStartCondition(rName string, pipelineExecutionStartCondition string) string {
	return acctest.ConfigCompose(
		testAccAwsImageBuilderImagePipelineConfigBase(rName),
		fmt.Sprintf(`
resource "aws_imagebuilder_image_pipeline" "test" {
  image_recipe_arn                 = aws_imagebuilder_image_recipe.test.arn
  infrastructure_configuration_arn = aws_imagebuilder_infrastructure_configuration.test.arn
  name                             = %[1]q

  schedule {
    pipeline_execution_start_condition = %[2]q
    schedule_expression                = "cron(0 0 * * ? *)"
  }
}
`, rName, pipelineExecutionStartCondition))
}

func testAccAwsImageBuilderImagePipelineConfigScheduleScheduleExpression(rName string, scheduleExpression string) string {
	return acctest.ConfigCompose(
		testAccAwsImageBuilderImagePipelineConfigBase(rName),
		fmt.Sprintf(`
resource "aws_imagebuilder_image_pipeline" "test" {
  image_recipe_arn                 = aws_imagebuilder_image_recipe.test.arn
  infrastructure_configuration_arn = aws_imagebuilder_infrastructure_configuration.test.arn
  name                             = %[1]q

  schedule {
    schedule_expression = %[2]q
  }
}
`, rName, scheduleExpression))
}

func testAccAwsImageBuilderImagePipelineConfigStatus(rName string, status string) string {
	return acctest.ConfigCompose(
		testAccAwsImageBuilderImagePipelineConfigBase(rName),
		fmt.Sprintf(`
resource "aws_imagebuilder_image_pipeline" "test" {
  image_recipe_arn                 = aws_imagebuilder_image_recipe.test.arn
  infrastructure_configuration_arn = aws_imagebuilder_infrastructure_configuration.test.arn
  name                             = %[1]q
  status                           = %[2]q
}
`, rName, status))
}

func testAccAwsImageBuilderImagePipelineConfigTags1(rName string, tagKey1 string, tagValue1 string) string {
	return acctest.ConfigCompose(
		testAccAwsImageBuilderImagePipelineConfigBase(rName),
		fmt.Sprintf(`
resource "aws_imagebuilder_image_pipeline" "test" {
  image_recipe_arn                 = aws_imagebuilder_image_recipe.test.arn
  infrastructure_configuration_arn = aws_imagebuilder_infrastructure_configuration.test.arn
  name                             = %[1]q

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1))
}

func testAccAwsImageBuilderImagePipelineConfigTags2(rName string, tagKey1 string, tagValue1 string, tagKey2 string, tagValue2 string) string {
	return acctest.ConfigCompose(
		testAccAwsImageBuilderImagePipelineConfigBase(rName),
		fmt.Sprintf(`
resource "aws_imagebuilder_image_pipeline" "test" {
  image_recipe_arn                 = aws_imagebuilder_image_recipe.test.arn
  infrastructure_configuration_arn = aws_imagebuilder_infrastructure_configuration.test.arn
  name                             = %[1]q

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2))
}
