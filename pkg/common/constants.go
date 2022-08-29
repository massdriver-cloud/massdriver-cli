package common

const (
	MassdriverURL             = "https://api.massdriver.cloud/"
	MassdriverYamlFilename    = "massdriver.yaml"
	ArtifactsSchemaFilename   = "schema-artifacts.json"
	ConnectionsSchemaFilename = "schema-connections.json"
	ParamsSchemaFilename      = "schema-params.json"
	UISchemaFilename          = "schema-ui.json"
	MaxBundleSizeMB           = 10
	MaxFileSizeMB             = 1
)

// named constants for common unix file permissions logic
// adapted from https://stackoverflow.com/questions/28969455/how-to-properly-instantiate-os-filemode/42718395#42718395
const (
	Read       = 04
	Write      = 02
	Execute    = 01
	UserShift  = 6
	GroupShift = 3
	OtherShift = 0

	UserRead    = Read << UserShift
	UserWrite   = Write << UserShift
	UserExecute = Execute << UserShift
	UserRW      = UserRead | UserWrite
	UserRWX     = UserRW | UserExecute

	GroupRead    = Read << GroupShift
	GroupWrite   = Write << GroupShift
	GroupExecute = Execute << GroupShift
	GroupRW      = GroupRead | GroupWrite
	GroupRWX     = GroupRW | GroupExecute

	OtherRead    = Read << OtherShift
	OtherWrite   = Write << OtherShift
	OtherExecute = Execute << OtherShift
	OtherRW      = OtherRead | OtherWrite
	OtherRWX     = OtherRW | OtherExecute

	AllRead    = UserRead | GroupRead | OtherRead
	AllWrite   = UserWrite | GroupWrite | OtherWrite
	AllExecute = UserExecute | GroupExecute | OtherExecute
	AllRX      = AllRead | AllExecute
	AllRW      = AllRead | AllWrite
	AllRWX     = AllRW | AllExecute
)

var FileIgnores = []string{
	".terraform",
	".tfstate",
	".tfvars",
	".md",
	".git",
	".DS_Store",
}

var FileAllows = []string{
	MassdriverYamlFilename,
	ArtifactsSchemaFilename,
	ConnectionsSchemaFilename,
	ParamsSchemaFilename,
	UISchemaFilename,
	"src",
}
