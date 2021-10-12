// Code generated by internal/generate/tags/main.go; DO NOT EDIT.
package acm

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/acm"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

// ListTags lists acm service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func ListTags(conn *acm.ACM, identifier string) (tftags.KeyValueTags, error) {
	input := &acm.ListTagsForCertificateInput{
		CertificateArn: aws.String(identifier),
	}

	output, err := conn.ListTagsForCertificate(input)

	if err != nil {
		return tftags.New(nil), err
	}

	return KeyValueTags(output.Tags), nil
}

// []*SERVICE.Tag handling

// Tags returns acm service tags.
func Tags(tags tftags.KeyValueTags) []*acm.Tag {
	result := make([]*acm.Tag, 0, len(tags))

	for k, v := range tags.Map() {
		tag := &acm.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		result = append(result, tag)
	}

	return result
}

// KeyValueTags creates tftags.KeyValueTags from acm service tags.
func KeyValueTags(tags []*acm.Tag) tftags.KeyValueTags {
	m := make(map[string]*string, len(tags))

	for _, tag := range tags {
		m[aws.StringValue(tag.Key)] = tag.Value
	}

	return tftags.New(m)
}

// UpdateTags updates acm service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func UpdateTags(conn *acm.ACM, identifier string, oldTagsMap interface{}, newTagsMap interface{}) error {
	oldTags := tftags.New(oldTagsMap)
	newTags := tftags.New(newTagsMap)

	if removedTags := oldTags.Removed(newTags); len(removedTags) > 0 {
		input := &acm.RemoveTagsFromCertificateInput{
			CertificateArn: aws.String(identifier),
			Tags:           Tags(removedTags.IgnoreAWS()),
		}

		_, err := conn.RemoveTagsFromCertificate(input)

		if err != nil {
			return fmt.Errorf("error untagging resource (%s): %w", identifier, err)
		}
	}

	if updatedTags := oldTags.Updated(newTags); len(updatedTags) > 0 {
		input := &acm.AddTagsToCertificateInput{
			CertificateArn: aws.String(identifier),
			Tags:           Tags(updatedTags.IgnoreAWS()),
		}

		_, err := conn.AddTagsToCertificate(input)

		if err != nil {
			return fmt.Errorf("error tagging resource (%s): %w", identifier, err)
		}
	}

	return nil
}
