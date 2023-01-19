package kms

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
)

func DataSourceCiphertext() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCiphertextRead,

		Schema: map[string]*schema.Schema{
			"plaintext": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},

			"key_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"context": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"ciphertext_blob": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCiphertextRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).KMSConn()

	keyID := d.Get("key_id").(string)
	req := &kms.EncryptInput{
		KeyId:     aws.String(keyID),
		Plaintext: []byte(d.Get("plaintext").(string)),
	}

	if ec := d.Get("context"); ec != nil {
		req.EncryptionContext = flex.ExpandStringMap(ec.(map[string]interface{}))
	}

	log.Printf("[DEBUG] KMS encrypting with KMS Key: %s", keyID)
	resp, err := conn.Encrypt(req)
	if err != nil {
		return fmt.Errorf("encrypting with KMS Key (%s): %w", keyID, err)
	}

	d.SetId(aws.StringValue(resp.KeyId))

	d.Set("ciphertext_blob", base64.StdEncoding.EncodeToString(resp.CiphertextBlob))

	return nil
}
