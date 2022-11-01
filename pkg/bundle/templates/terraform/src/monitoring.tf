# Bundles are better with alarms
# https://docs.massdriver.cloud/bundles/monitoring
# Uncomment the below to add an alarm channel for your bundle. See the docs
# for examples of how to add an alarm using the channel.

# AWS
# module "alarm_channel" {
#   source      = "github.com/massdriver-cloud/terraform-modules//aws-alarm-channel?ref=aa08797"
#   md_metadata = var.md_metadata
# }

# Azure
# module "alarm_channel" {
#   source      = "github.com/massdriver-cloud/terraform-modules//azure-alarm-channel?ref=aa08797"
#   md_metadata = var.md_metadata
# }

# GCP
# module "alarm_channel" {
#   source      = "github.com/massdriver-cloud/terraform-modules//gcp-alarm-channel?ref=aa08797"
#   md_metadata = var.md_metadata
# }

# TODO: add an alarm
# See https://docs.massdriver.cloud/bundles/monitoring for more information
