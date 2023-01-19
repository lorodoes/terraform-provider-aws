package ec2_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfec2 "github.com/hashicorp/terraform-provider-aws/internal/service/ec2"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccVPCEndpointConnectionNotification_basic(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_vpc_endpoint_connection_notification.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckVPCEndpointConnectionNotificationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCEndpointConnectionNotificationConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCEndpointConnectionNotificationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "connection_events.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "notification_type", "Topic"),
					resource.TestCheckResourceAttr(resourceName, "state", "Enabled"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccVPCEndpointConnectionNotificationConfig_modified(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCEndpointConnectionNotificationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "connection_events.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "notification_type", "Topic"),
					resource.TestCheckResourceAttr(resourceName, "state", "Enabled"),
				),
			},
		},
	})
}

func testAccCheckVPCEndpointConnectionNotificationDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_vpc_endpoint_connection_notification" {
			continue
		}

		_, err := tfec2.FindVPCConnectionNotificationByID(conn, rs.Primary.ID)

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("EC2 VPC Endpoint Connection Notification %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccCheckVPCEndpointConnectionNotificationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No EC2 VPC Endpoint Connection Notification ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn()

		_, err := tfec2.FindVPCConnectionNotificationByID(conn, rs.Primary.ID)

		return err
	}
}

func testAccVPCEndpointConnectionNotificationConfig_base(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigVPCWithSubnets(rName, 2), fmt.Sprintf(`
data "aws_partition" "current" {}

resource "aws_lb" "nlb_test" {
  name = %[1]q

  subnets = aws_subnet.test[*].id

  load_balancer_type         = "network"
  internal                   = true
  idle_timeout               = 60
  enable_deletion_protection = false
}

data "aws_caller_identity" "current" {}

resource "aws_vpc_endpoint_service" "test" {
  acceptance_required = false

  network_load_balancer_arns = [
    aws_lb.nlb_test.id,
  ]

  allowed_principals = [
    data.aws_caller_identity.current.arn
  ]

  tags = {
    Name = %[1]q
  }
}

resource "aws_sns_topic" "test" {
  name = %[1]q

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "vpce.${data.aws_partition.current.dns_suffix}"
      },
      "Action": "SNS:Publish",
      "Resource": "arn:${data.aws_partition.current.partition}:sns:*:*:%[1]s"
    }
  ]
}
POLICY
}
`, rName))
}

func testAccVPCEndpointConnectionNotificationConfig_basic(rName string) string {
	return acctest.ConfigCompose(testAccVPCEndpointConnectionNotificationConfig_base(rName), `
resource "aws_vpc_endpoint_connection_notification" "test" {
  vpc_endpoint_service_id     = aws_vpc_endpoint_service.test.id
  connection_notification_arn = aws_sns_topic.test.arn
  connection_events           = ["Accept", "Reject"]
}
`)
}

func testAccVPCEndpointConnectionNotificationConfig_modified(rName string) string {
	return acctest.ConfigCompose(testAccVPCEndpointConnectionNotificationConfig_base(rName), `
resource "aws_vpc_endpoint_connection_notification" "test" {
  vpc_endpoint_service_id     = aws_vpc_endpoint_service.test.id
  connection_notification_arn = aws_sns_topic.test.arn
  connection_events           = ["Accept"]
}
`)
}
