package common

var MassdriverArtifactDefinitions = []string{
	"massdriver/aws-dynamodb-table",
	"massdriver/aws-efs-file-system",
	"massdriver/aws-eventbridge",
	"massdriver/aws-iam-role",
	"massdriver/aws-s3-bucket",
	"massdriver/aws-sns-topic",
	"massdriver/aws-sqs-queue",
	"massdriver/aws-vpc",
	"massdriver/azure-communication-service",
	"massdriver/azure-data-lake-storage",
	"massdriver/azure-databricks-workspace",
	"massdriver/azure-fhir-service",
	"massdriver/azure-service-principal",
	"massdriver/azure-storage-account",
	"massdriver/azure-virtual-network",
	"massdriver/cosmosdb-sql-authentication",
	// draft-node
	"massdriver/elasticsearch-authentication",
	// env-file
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
	"massdriver/sftp-authentication",
}

func ListMassdriverArtifactDefinitions() ([]string, error) {
	// TODO this should list these artifacts from the massdriver API
	return MassdriverArtifactDefinitions, nil
}
