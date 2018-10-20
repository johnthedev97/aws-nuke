package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iot"
)

type IoTPolicy struct {
	svc  *iot.IoT
	name *string
}

func init() {
	register("IoTPolicy", ListIoTPolicies)
}

func ListIoTPolicies(sess *session.Session) ([]Resource, error) {
	svc := iot.New(sess)
	resources := []Resource{}

	params := &iot.ListPoliciesInput{
		PageSize: aws.Int64(25),
	}
	for {
		output, err := svc.ListPolicies(params)
		if err != nil {
			return nil, err
		}

		for _, policy := range output.Policies {
			resources = append(resources, &IoTPolicy{
				svc:  svc,
				name: policy.PolicyName,
			})
		}
		if output.NextMarker == nil {
			break
		}

		params.Marker = output.NextMarker
	}

	return resources, nil
}

func (f *IoTPolicy) DetachAllPrincipals() error {
	params := &iot.ListTargetsForPolicyInput{
		PolicyName: f.name,
		PageSize:   aws.Int64(25),
	}
	for {
		output, err := f.svc.ListTargetsForPolicy(params)
		if err != nil {
			return err
		}

		for _, target := range output.Targets {
			_, err = f.svc.DetachPrincipalPolicy(&iot.DetachPrincipalPolicyInput{
				PolicyName: f.name,
				Principal:  target,
			})
			if err != nil {
				return err
			}
		}
		if output.NextMarker == nil {
			break
		}

		params.Marker = output.NextMarker
	}
	return nil
}

func (f *IoTPolicy) Remove() error {
	err := f.DetachAllPrincipals()
	if err != nil {
		return err
	}
	_, err = f.svc.DeletePolicy(&iot.DeletePolicyInput{
		PolicyName: f.name,
	})

	return err
}

func (f *IoTPolicy) String() string {
	return *f.name
}
