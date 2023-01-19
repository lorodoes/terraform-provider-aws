// Code generated by internal/generate/tags/main.go; DO NOT EDIT.
package logs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
)

// ListLogGroupTags lists logs service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func ListLogGroupTags(conn cloudwatchlogsiface.CloudWatchLogsAPI, identifier string) (tftags.KeyValueTags, error) {
	return ListLogGroupTagsWithContext(context.Background(), conn, identifier)
}

func ListLogGroupTagsWithContext(ctx context.Context, conn cloudwatchlogsiface.CloudWatchLogsAPI, identifier string) (tftags.KeyValueTags, error) {
	input := &cloudwatchlogs.ListTagsLogGroupInput{
		LogGroupName: aws.String(identifier),
	}

	output, err := conn.ListTagsLogGroupWithContext(ctx, input)

	if err != nil {
		return tftags.New(nil), err
	}

	return KeyValueTags(output.Tags), nil
}

// UpdateLogGroupTags updates logs service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func UpdateLogGroupTags(conn cloudwatchlogsiface.CloudWatchLogsAPI, identifier string, oldTags interface{}, newTags interface{}) error {
	return UpdateLogGroupTagsWithContext(context.Background(), conn, identifier, oldTags, newTags)
}
func UpdateLogGroupTagsWithContext(ctx context.Context, conn cloudwatchlogsiface.CloudWatchLogsAPI, identifier string, oldTagsMap interface{}, newTagsMap interface{}) error {
	oldTags := tftags.New(oldTagsMap)
	newTags := tftags.New(newTagsMap)

	if removedTags := oldTags.Removed(newTags); len(removedTags) > 0 {
		input := &cloudwatchlogs.UntagLogGroupInput{
			LogGroupName: aws.String(identifier),
			Tags:         aws.StringSlice(removedTags.IgnoreAWS().Keys()),
		}

		_, err := conn.UntagLogGroupWithContext(ctx, input)

		if err != nil {
			return fmt.Errorf("untagging resource (%s): %w", identifier, err)
		}
	}

	if updatedTags := oldTags.Updated(newTags); len(updatedTags) > 0 {
		input := &cloudwatchlogs.TagLogGroupInput{
			LogGroupName: aws.String(identifier),
			Tags:         Tags(updatedTags.IgnoreAWS()),
		}

		_, err := conn.TagLogGroupWithContext(ctx, input)

		if err != nil {
			return fmt.Errorf("tagging resource (%s): %w", identifier, err)
		}
	}

	return nil
}
