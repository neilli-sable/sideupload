package setting

// Setting ...
type Setting struct {
	Sideupload Sideupload
}

// Sideupload ...
type Sideupload struct {
	ListenPort     uint
	TargetDir      string
	CronWithSecond string
	BackupStorage  BackupStorage
}

// BackupStorage ...
type BackupStorage struct {
	CustomEndpoint string
	Region         string
	BucketName     string
	Prefix         string
	AccessKey      string
	SecretKey      string
}
