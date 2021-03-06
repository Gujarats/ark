package ark

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

// Create new access/secret keys
// following the best practice AWS
func CreateNewAccessKey(sess *session.Session, userName string) (*iam.AccessKey, error) {
	svc := iam.New(sess)

	listAccessKeys, err := svc.ListAccessKeys(&iam.ListAccessKeysInput{
		UserName: aws.String(userName),
	})
	if err != nil {
		return nil, err
	}

	err = deleteInactiveKey(svc, userName, listAccessKeys)
	if err != nil {
		return nil, err
	}

	result, err := svc.CreateAccessKey(&iam.CreateAccessKeyInput{
		UserName: aws.String(userName),
	})
	if err != nil {
		return nil, err
	}

	return result.AccessKey, nil
}

// if All active then delete the first key
// if found inactive then delete it
// for avoid creating new keys failure if IAM user is equal MAxIAMUser
// final result must be one accesKey with Inactive status
func deleteInactiveKey(svc *iam.IAM, userName string, listAccessKeys *iam.ListAccessKeysOutput) error {
	var foundInactive bool
	if len(listAccessKeys.AccessKeyMetadata) == MaxIAMUser {
		for _, accessKey := range listAccessKeys.AccessKeyMetadata {

			if *accessKey.Status == Inactive {
				foundInactive = true
				_, err := svc.DeleteAccessKey(&iam.DeleteAccessKeyInput{
					AccessKeyId: accessKey.AccessKeyId,
					UserName:    aws.String(userName),
				})

				if err != nil {
					return err
				}
			}

			if *accessKey.Status == Active {
				deactivateKey(svc, userName, *accessKey.AccessKeyId)
			}
		}

		if !foundInactive {
			_, err := svc.DeleteAccessKey(&iam.DeleteAccessKeyInput{
				AccessKeyId: listAccessKeys.AccessKeyMetadata[0].AccessKeyId,
				UserName:    aws.String(userName),
			})

			deactivateKey(svc, userName, *listAccessKeys.AccessKeyMetadata[1].AccessKeyId)

			if err != nil {
				return err
			}
		}
	} else if len(listAccessKeys.AccessKeyMetadata) == 1 {
		deactivateKey(svc, userName, *listAccessKeys.AccessKeyMetadata[0].AccessKeyId)
	}

	return nil
}

// before creating new one deactivate last key
func deactivateKey(svc *iam.IAM, userName string, accessKeyId string) error {
	_, err := svc.UpdateAccessKey(&iam.UpdateAccessKeyInput{
		UserName:    aws.String(userName),
		Status:      aws.String(Inactive),
		AccessKeyId: aws.String(accessKeyId),
	})

	if err != nil {
		return err
	}

	return nil
}
