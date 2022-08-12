package common

var MassdriverArtifacts = []string{
	"massdriver/aws-dynamodb-table",
	"massdriver/aws-efs-file-system",
	"massdriver/aws-eventbridge",
	"massdriver/aws-iam-role",
	"massdriver/aws-s3-bucket",
	"massdriver/aws-sns-topic",
	"massdriver/aws-sqs-queue",
	"massdriver/aws-vpc",
	"massdriver/azure-data-lake-storage",
	"massdriver/azure-databricks-workspace",
	"massdriver/azure-service-principal",
	"massdriver/azure-virtual-network",
	"massdriver/draft-node",
	"massdriver/elasticsearch-authentication",
	"massdriver/env-file",
	"massdriver/gcp-bucket-https",
	"massdriver/gcp-cloud-function",
	"massdriver/gcp-firebase-authentication",
	"massdriver/gcp-global-network",
	"massdriver/gcp-pubsub-subscription",
	"massdriver/gcp-pubsub-topic",
	"massdriver/gcp-service-account",
	"massdriver/gcp-subnetwork",
	"massdriver/kafka-authentication",
	"massdriver/kubernetes-cluster",
	"massdriver/mongo-authentication",
	"massdriver/mysql-authentication",
	"massdriver/postgresql-authentication",
	"massdriver/redis-authentication",
}

func GetMassdriverArtifacts() ([]string, error) {
	// TODO this should list these artifacts from the massdriver API
	return MassdriverArtifacts, nil
}
