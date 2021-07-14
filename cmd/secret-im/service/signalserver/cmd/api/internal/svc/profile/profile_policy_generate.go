package profile

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"time"
)

type PostPolicyGenerator struct {
	region string
	bucket string
	awsAccessId string

}

// 创建策略生成器
func NewPostPolicyGenerator(region, bucket, awsAccessID string) *PostPolicyGenerator {
	return &PostPolicyGenerator{
		region:      region,
		bucket:      bucket,
		awsAccessId: awsAccessID,
	}
}
// 2014-02-09T21:23:56.870Z
func (p *PostPolicyGenerator) CreateFor(now time.Time, object string,maxSizeInBytes int) (string,string)  {
	expiration     := now.Add(time.Minute*30).Format("2006-01-02T15:04:05.999Z");
	credentialDate := now.Format("20060102");
	requestDate    := now.Format("20060102T150405Z" );
	credential     := fmt.Sprintf("%s/%s/%s/s3/aws4_request", p.awsAccessId, credentialDate, p.region);

	policy := fmt.Sprintf("{ \"expiration\": \"%s\",\n"+
		"  \"conditions\": [\n"+
		"    {\"bucket\": \"%s\"},\n"+
		"    {\"key\": \"%s\"},\n"+
		"    {\"acl\": \"private\"},\n"+
		"    [\"starts-with\", \"$Content-Type\", \"\"],\n"+
		"    [\"content-length-range\", 1, "+strconv.Itoa(maxSizeInBytes)+"],\n"+
		"\n"+
		"    {\"x-amz-credential\": \"%s\"},\n"+
		"    {\"x-amz-algorithm\": \"AWS4-HMAC-SHA256\"},\n"+
		"    {\"x-amz-date\": \"%s\" }\n"+
		"  ]\n"+
		"}", expiration, p.bucket, object, credential, requestDate)
	return credential, base64.StdEncoding.EncodeToString([]byte(policy))
}