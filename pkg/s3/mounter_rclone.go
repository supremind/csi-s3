package s3

import (
	"fmt"
	"os"
)

// Implements Mounter
type rcloneMounter struct {
	url             string
	region          string
	accessKeyID     string
	secretAccessKey string
}

const (
	rcloneCmd = "rclone"
)

func newRcloneMounter(cfg *Config) (Mounter, error) {
	return &rcloneMounter{
		url:             cfg.Endpoint,
		region:          cfg.Region,
		accessKeyID:     cfg.AccessKeyID,
		secretAccessKey: cfg.SecretAccessKey,
	}, nil
}

func (rclone *rcloneMounter) Stage(*volume, string) error {
	return nil
}

func (rclone *rcloneMounter) Unstage(*volume, string) error {
	return nil
}

func (rclone *rcloneMounter) Mount(vol *volume, source string, target string) error {
	args := []string{
		"mount",
		fmt.Sprintf(":s3:%s/%s", vol.Bucket, vol.Prefix),
		fmt.Sprintf("%s", target),
		"--daemon",
		"--s3-provider=AWS",
		"--s3-env-auth=true",
		fmt.Sprintf("--s3-region=%s", rclone.region),
		fmt.Sprintf("--s3-endpoint=%s", rclone.url),
		"--allow-other",
		// TODO: make this configurable
		"--vfs-cache-mode=writes",
	}
	os.Setenv("AWS_ACCESS_KEY_ID", rclone.accessKeyID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", rclone.secretAccessKey)
	return fuseMount(target, rcloneCmd, args)
}
