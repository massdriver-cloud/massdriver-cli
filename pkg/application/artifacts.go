package application

var MassdriverArtifacts = []string{
	"aws-dynamodb-table",
	"aws-efs-file-system",
	"aws-eventbridge",
	"aws-iam-role",
	"aws-s3-bucket",
	"aws-sns-topic",
	"aws-sqs-queue",
	"aws-vpc",
	"azure-data-lake-storage",
	"azure-databricks-workspace",
	"azure-service-principal",
	"azure-virtual-network",
	"draft-node",
	"elasticsearch-authentication",
	"env-file",
	"gcp-bucket-https",
	"gcp-cloud-function",
	"gcp-firebase-authentication",
	"gcp-global-network",
	"gcp-pubsub-subscription",
	"gcp-pubsub-topic",
	"gcp-service-account",
	"gcp-subnetwork",
	"kafka-authentication",
	"kubernetes-cluster",
	"kubernetes-cluster",
	"mongo-authentication",
	"mysql-authentication",
	"postgresql-authentication",
	"redis-authentication",
}

func GetArtifacts() ([]string, error) {
	// TODO this should list these artifacts from the massdriver API
	return MassdriverArtifacts, nil
}
