package application

type Application struct {
	Uuid        string                 `json:"uuid" yaml:"uuid"`
	Schema      string                 `json:"schema" yaml:"schema"`
	Title       string                 `json:"title" yaml:"title"`
	Description string                 `json:"description" yaml:"description"`
	Ref         string                 `json:"ref" yaml:"ref"`
	Type        string                 `json:"type" yaml:"type"`
	Access      string                 `json:"access" yaml:"access"`
	Params      map[string]interface{} `json:"params" yaml:"params"`

	Deployment   ApplicationDeployment     `json:"deployment" yaml:"deployment"`
	Dependencies []ApplicationDependencies `json:"dependencies" yaml:"dependencies"`
}

type ApplicationDeployment struct {
	Type  string `json:"type" yaml:"type"`
	Chart string `json:"chart" yaml:"chart"`
}

type ApplicationDependencies struct {
	Type     string                        `json:"type" yaml:"type"`
	Field    string                        `json:"field" yaml:"field"`
	Required *bool                         `json:"required" yaml:"required,omitempty"`
	Env      []ApplicationDependenciesEnvs `json:"env" yaml:"env"`
	Policy   string                        `json:"policy" yaml:"policy"`
}

type ApplicationDependenciesEnvs struct {
	Name  string `json:"name" yaml:"name"`
	Value string `json:"value" yaml:"value"`
}

var KubernetesClusterConnection = `{"$id":"https://schemas.massdriver.cloud/definitions/artifacts/kubernetes-cluster.json","$md":{"defaultTargetConnectionGroup":"credentials","defaultTargetConnectionGroupLabel":"Kubernetes","importing":{"fileUploadArtifactDataPath":["data","authentication"],"fileUploadType":"yaml","group":"authentication"}},"$schema":"http://json-schema.org/draft-07/schema","additionalProperties":false,"description":"Kubernetes cluster authentication and cloud-specific configuration.","properties":{"data":{"additionalProperties":false,"properties":{"authentication":{"additionalProperties":false,"properties":{"cluster":{"additionalProperties":false,"properties":{"certificate-authority-data":{"type":"string"},"server":{"type":"string"}},"required":["server","certificate-authority-data"],"type":"object"},"user":{"additionalProperties":false,"properties":{"token":{"type":"string"}},"required":["token"],"type":"object"}},"required":["cluster","user"],"type":"object"},"infrastructure":{"additionalProperties":true,"description":"Cloud specific Kubernetes configuration data.","oneOf":[{"$id":"https://schemas.massdriver.cloud/definitions/types/aws-infrastructure-eks.json","$schema":"http://json-schema.org/draft-07/schema","additionalProperties":false,"description":"","properties":{"arn":{"$id":"https://schemas.massdriver.cloud/definitions/types/aws-arn.json","$schema":"http://json-schema.org/draft-07/schema","description":"Amazon Resource Name","examples":["arn:aws:rds::ACCOUNT_NUMBER:db/prod","arn:aws:ec2::ACCOUNT_NUMBER:vpc/vpc-foo"],"title":"AWS ARN","type":"string"},"oidc_provider_arn":{"$id":"https://schemas.massdriver.cloud/definitions/types/aws-arn.json","$schema":"http://json-schema.org/draft-07/schema","description":"Amazon Resource Name","examples":["arn:aws:rds::ACCOUNT_NUMBER:db/prod","arn:aws:ec2::ACCOUNT_NUMBER:vpc/vpc-foo"],"title":"AWS ARN","type":"string"}},"required":["arn","oidc_provider_arn"],"title":"AWS EKS infrastructure config","type":"object"},{"$id":"https://schemas.massdriver.cloud/definitions/types/azure-infrastructure-aks.json","$schema":"http://json-schema.org/draft-07/schema","additionalProperties":false,"description":"","properties":{"id":{"$id":"https://schemas.massdriver.cloud/definitions/types/azure-resource-id.json","$schema":"http://json-schema.org/draft-07/schema","description":"Azure Resource ID","examples":["/subscriptions/12345678-1234-1234-abcd-1234567890ab/resourceGroups/resource-group-name/providers/Microsoft.Network/virtualNetworks/network-name"],"title":"Azure Resource ID","type":"string"}},"required":["id"],"title":"Azure AKS infrastructure config","type":"object"},{"$id":"https://schemas.massdriver.cloud/definitions/types/gcp-infrastructure-gke.json","$schema":"http://json-schema.org/draft-07/schema","additionalProperties":false,"description":"","properties":{"grn":{"$id":"https://schemas.massdriver.cloud/definitions/types/gcp-grn.json","$schema":"http://json-schema.org/draft-07/schema","description":"GCP Resource Name (GRN)","examples":["projects/my-project/locations/us-central1/instances/my-instance"],"title":"GCP Resource Name (GRN)","type":"string"}},"required":["grn"],"title":"GCP GKE infrastructure config","type":"object"}],"title":"Cloud provider configuration","type":"object"}},"required":["authentication"],"type":"object"},"specs":{"additionalProperties":false,"properties":{"kubernetes":{"$id":"https://schemas.massdriver.cloud/definitions/specs/kubernetes.json","$schema":"http://json-schema.org/draft-07/schema","additionalProperties":false,"description":"Kubernetes distribution and version specifications.","properties":{"cloud":{"enum":["aws","gcp","azure"],"type":"string"},"distribution":{"enum":["eks","gke","aks"],"type":"string"},"platform_version":{"type":"string"},"version":{"type":"string"}},"required":["version","cloud","distribution"],"title":"Kubernetes","type":"object"}},"type":"object"}},"required":["data","specs"],"title":"Kubernetes Cluster","type":"object"}`
var CloudAuthenticationConnection = `{"oneOf":[{"$id":"https://schemas.massdriver.cloud/definitions/artifacts/aws-iam-role.json","$md":{"defaultTargetConnectionGroup":"credentials","defaultTargetConnectionGroupLabel":"AWS","diagram":{"isLinkable":false},"importing":{"group":"authentication"}},"$schema":"http://json-schema.org/draft-07/schema","additionalProperties":false,"description":"","properties":{"data":{"additionalProperties":false,"properties":{"arn":{"$id":"https://schemas.massdriver.cloud/definitions/types/aws-arn.json","$schema":"http://json-schema.org/draft-07/schema","description":"Amazon Resource Name","examples":["arn:aws:rds::ACCOUNT_NUMBER:db/prod","arn:aws:ec2::ACCOUNT_NUMBER:vpc/vpc-foo"],"title":"AWS ARN","type":"string"},"external_id":{"description":"An external ID is a piece of data that can be passed to the AssumeRole API of the Security Token Service (STS). You can then use the external ID in the condition element in a role’s trust policy, allowing the role to be assumed only when a certain value is present in the external ID.","title":"External ID","type":"string"}},"required":["arn"],"title":"Artifact Data","type":"object"},"specs":{"additionalProperties":false,"properties":{"aws":{"$id":"https://schemas.massdriver.cloud/definitions/specs/aws.json","$schema":"http://json-schema.org/draft-07/schema","additionalProperties":false,"description":"","properties":{"region":{"$id":"https://schemas.massdriver.cloud/definitions/types/aws-region.json","$schema":"http://json-schema.org/draft-07/schema","default":"us-west-2","description":"AWS regional data center to provision in.","enum":["us-west-2","us-east-1","us-east-2"],"examples":["us-west-2"],"title":"AWS Region","type":"string"},"resource":{"title":"AWS Resource Type","type":"string"},"service":{"title":"AWS Service","type":"string"},"zone":{"$id":"https://schemas.massdriver.cloud/definitions/types/aws-zone.json","$schema":"http://json-schema.org/draft-07/schema","description":"AWS Availability Zone","examples":[],"title":"AWS Zone","type":"string"}},"required":[],"title":"AWS Artifact Specs","type":"object"}},"title":"Artifact Specs","type":"object"}},"required":["data","specs"],"title":"AWS IAM Role","type":"object"},{"$id":"https://schemas.massdriver.cloud/definitions/artifacts/gcp-service-account.json","$md":{"defaultTargetConnectionGroup":"credentials","defaultTargetConnectionGroupLabel":"GCP Service Account","diagram":{"isLinkable":false},"importing":{"fileUploadType":"json","group":"authentication"}},"$schema":"http://json-schema.org/draft-07/schema","additionalProperties":false,"description":"GCP Service Account","properties":{"data":{"additionalProperties":false,"properties":{"auth_provider_x509_cert_url":{"default":"https://www.googleapis.com/oauth2/v1/certs","description":"","title":"Auth Provider x509 Certificate URL","type":"string"},"auth_uri":{"default":"https://accounts.google.com/o/oauth2/auth","description":"","title":"Auth URI","type":"string"},"client_email":{"description":"","title":"Client e-mail","type":"string"},"client_id":{"description":"","title":"Client ID","type":"string"},"client_x509_cert_url":{"description":"","title":"Client x509 Certificate URL","type":"string"},"private_key":{"description":"","title":"Private Key","type":"string"},"private_key_id":{"description":"","title":"Private Key ID","type":"string"},"project_id":{"description":"","title":"Project ID","type":"string"},"token_uri":{"default":"https://oauth2.googleapis.com/token","description":"","title":"Token URI","type":"string"},"type":{"default":"service_account","description":"","title":"Type","type":"string"}},"required":["auth_provider_x509_cert_url","auth_uri","client_email","client_id","client_x509_cert_url","private_key","private_key_id","project_id","token_uri","type"],"title":"Artifact Data","type":"object"},"specs":{"additionalProperties":false,"properties":{"gcp":{"$id":"https://schemas.massdriver.cloud/definitions/specs/gcp.json","$schema":"http://json-schema.org/draft-07/schema","additionalProperties":false,"description":"","properties":{"project":{"title":"GCP Project","type":"string"},"region":{"$id":"https://schemas.massdriver.cloud/definitions/types/gcp-region.json","$schema":"http://json-schema.org/draft-07/schema","description":"GCP region","enum":["us-east1","us-east4","us-west1","us-west2","us-west3","us-west4","us-central1"],"examples":["us-west2"],"title":"GCP Region","type":"string"},"resource":{"title":"GCP Resource Type","type":"string"},"service":{"title":"GCP Service","type":"string"},"zone":{"$id":"https://schemas.massdriver.cloud/definitions/types/gcp-zone.json","$schema":"http://json-schema.org/draft-07/schema","description":"GCP Zone","examples":[],"title":"GCP Zone","type":"string"}},"required":[],"title":"GCP Artifact Specs","type":"object"}},"title":"Artifact Specs","type":"object"}},"required":["data","specs"],"title":"GCP Service Account","type":"object"},{"$id":"https://schemas.massdriver.cloud/definitions/artifacts/azure-service-principal.json","$md":{"defaultTargetConnectionGroup":"credentials","defaultTargetConnectionGroupLabel":"Azure","diagram":{"isLinkable":false},"importing":{"group":"authentication"}},"$schema":"http://json-schema.org/draft-07/schema","additionalProperties":false,"description":"","properties":{"data":{"additionalProperties":false,"properties":{"client_id":{"title":"Client ID","type":"string"},"client_secret":{"title":"Client Secret","type":"string"},"subscription_id":{"title":"Subscription ID","type":"string"},"tenant_id":{"title":"Tenant ID","type":"string"}},"required":["client_id","tenant_id","client_secret","subscription_id"],"title":"Artifact Data","type":"object"},"specs":{"additionalProperties":false,"properties":{},"title":"Artifact Specs","type":"object"}},"required":["data","specs"],"title":"Azure Service Principal","type":"object"}]}`
