package s3

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/ray-project/kuberay/historyserver/pkg/collector/types"
)

const DefaultS3Bucket = "ray-historyserver"

// s3Config sets configuration for s3 and s3 compatible storage options
type s3Config struct {
	S3ForcePathStyle *bool
	DisableSSL       *bool
	Endpoint         string
	Bucket           string
	Region           string
	ID               string
	Secret           string
	Token            string
	types.RayCollectorConfig

	// Provider is the s3 compatible storage option provider
	Provider string
}

func getBucketWithDefault(provider string) string {
	var bucket string
	if provider == "gcs" {
		bucket = os.Getenv("GCS_BUCKET")
	} else {
		bucket = os.Getenv("S3_BUCKET")
	}

	if bucket == "" {
		return DefaultS3Bucket
	}
	return bucket
}

func (c *s3Config) complete(rcc *types.RayCollectorConfig, jd map[string]interface{}) {
	c.Provider = "aws"
	c.RayCollectorConfig = *rcc
	c.ID = os.Getenv("AWS_S3ID")
	c.Secret = os.Getenv("AWS_S3SECRET")
	c.Token = os.Getenv("AWS_S3TOKEN")
	c.Bucket = getBucketWithDefault(c.Provider)
	if len(jd) == 0 {
		c.Endpoint = os.Getenv("S3_ENDPOINT")
		c.Region = os.Getenv("S3_REGION")
		if os.Getenv("S3FORCE_PATH_STYLE") != "" {
			c.S3ForcePathStyle = aws.Bool(os.Getenv("S3FORCE_PATH_STYLE") == "true")
		}
		if os.Getenv("S3DISABLE_SSL") != "" {
			c.DisableSSL = aws.Bool(os.Getenv("S3DISABLE_SSL") == "true")
		}
	} else {
		if bucket, ok := jd["s3Bucket"]; ok {
			c.Bucket = bucket.(string)
		}
		if endpoint, ok := jd["s3Endpoint"]; ok {
			c.Endpoint = endpoint.(string)
		}
		if region, ok := jd["s3Region"]; ok {
			c.Region = region.(string)
		}
		if forcePathStyle, ok := jd["s3ForcePathStyle"]; ok {
			c.S3ForcePathStyle = aws.Bool(forcePathStyle.(string) == "true")
		}
		if s3disableSSL, ok := jd["s3DisableSSL"]; ok {
			c.DisableSSL = aws.Bool(s3disableSSL.(string) == "true")
		}
	}
}

func (c *s3Config) completeHSConfig(rcc *types.RayHistoryServerConfig, jd map[string]interface{}) {
	c.Provider = "aws"
	c.RayCollectorConfig = types.RayCollectorConfig{
		RootDir: rcc.RootDir,
	}
	c.ID = os.Getenv("AWS_S3ID")
	c.Secret = os.Getenv("AWS_S3SECRET")
	c.Token = os.Getenv("AWS_S3TOKEN")
	c.Bucket = getBucketWithDefault(c.Provider) // Use default if S3_BUCKET not set
	if len(jd) == 0 {
		c.Endpoint = os.Getenv("S3_ENDPOINT")
		c.Region = os.Getenv("S3_REGION")
		if os.Getenv("S3FORCE_PATH_STYLE") != "" {
			c.S3ForcePathStyle = aws.Bool(os.Getenv("S3FORCE_PATH_STYLE") == "true")
		}
		if os.Getenv("S3DISABLE_SSL") != "" {
			c.DisableSSL = aws.Bool(os.Getenv("S3DISABLE_SSL") == "true")
		}
	} else {
		if bucket, ok := jd["s3Bucket"]; ok {
			c.Bucket = bucket.(string)
		}
		if endpoint, ok := jd["s3Endpoint"]; ok {
			c.Endpoint = endpoint.(string)
		}
		if region, ok := jd["s3Region"]; ok {
			c.Region = region.(string)
		}
		if forcePathStyle, ok := jd["s3ForcePathStyle"]; ok {
			c.S3ForcePathStyle = aws.Bool(forcePathStyle.(string) == "true")
		}
		if s3disableSSL, ok := jd["s3DisableSSL"]; ok {
			c.DisableSSL = aws.Bool(s3disableSSL.(string) == "true")
		}
	}
}

func (c *s3Config) completeHSConfigForGCS(rcc *types.RayHistoryServerConfig, jd map[string]interface{}) {
	c.Provider = "gcs"
	c.RayCollectorConfig = types.RayCollectorConfig{
		RootDir: rcc.RootDir,
	}
	c.ID = os.Getenv("GCS_ACCESSKEY")
	c.Secret = os.Getenv("GCS_SECRETKEY")

	c.Bucket = getBucketWithDefault(c.Provider)
	if len(jd) == 0 {
		c.Endpoint = os.Getenv("GCS_ENDPOINT")
		if os.Getenv("GCS_DISABLE_SSL") != "" {
			c.DisableSSL = aws.Bool(os.Getenv("GCS_DISABLE_SSL") == "true")
		}
	} else {
		if bucket, ok := jd["gcsBucket"]; ok {
			c.Bucket = bucket.(string)
		}
		if endpoint, ok := jd["gcsEndpoint"]; ok {
			c.Endpoint = endpoint.(string)
		}
		if gcsDisableSSL, ok := jd["gcsDisableSSL"]; ok {
			c.DisableSSL = aws.Bool(gcsDisableSSL.(string) == "true")
		}
	}

	// Set the remaining required GCS fields, token is not needed for GCS
	c.Region = "auto"                   // Any string will work, GCS is global but S3 requires Region to be filled
	c.S3ForcePathStyle = aws.Bool(true) // Has to be true for GCS compatibility
}

func (c *s3Config) completeCollectorConfigForGCS(rcc *types.RayCollectorConfig, jd map[string]interface{}) {
	c.Provider = "gcs"
	c.RayCollectorConfig = *rcc
	c.ID = os.Getenv("GCS_ACCESSKEY")
	c.Secret = os.Getenv("GCS_SECRETKEY")

	c.Bucket = getBucketWithDefault(c.Provider)
	if len(jd) == 0 {
		c.Endpoint = os.Getenv("GCS_ENDPOINT")
		if os.Getenv("GCS_DISABLE_SSL") != "" {
			c.DisableSSL = aws.Bool(os.Getenv("GCS_DISABLE_SSL") == "true")
		}
	} else {
		if bucket, ok := jd["gcsBucket"]; ok {
			c.Bucket = bucket.(string)
		}
		if endpoint, ok := jd["gcsEndpoint"]; ok {
			c.Endpoint = endpoint.(string)
		}
		if gcsDisableSSL, ok := jd["gcsDisableSSL"]; ok {
			c.DisableSSL = aws.Bool(gcsDisableSSL.(string) == "true")
		}
	}

	c.Region = "auto"                   // Any string will work, GCS is global but S3 requires Region to be filled
	c.S3ForcePathStyle = aws.Bool(true) // Has to be true for GCS compatibility
}
