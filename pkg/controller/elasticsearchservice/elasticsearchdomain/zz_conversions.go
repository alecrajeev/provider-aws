/*
Copyright 2021 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by ack-generate. DO NOT EDIT.

package elasticsearchdomain

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	svcsdk "github.com/aws/aws-sdk-go/service/elasticsearchservice"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	svcapitypes "github.com/crossplane/provider-aws/apis/elasticsearchservice/v1alpha1"
)

// NOTE(muvaf): We return pointers in case the function needs to start with an
// empty object, hence need to return a new pointer.

// GenerateDescribeElasticsearchDomainInput returns input for read
// operation.
func GenerateDescribeElasticsearchDomainInput(cr *svcapitypes.ElasticsearchDomain) *svcsdk.DescribeElasticsearchDomainInput {
	res := &svcsdk.DescribeElasticsearchDomainInput{}

	if cr.Spec.ForProvider.DomainName != nil {
		res.SetDomainName(*cr.Spec.ForProvider.DomainName)
	}

	return res
}

// GenerateElasticsearchDomain returns the current state in the form of *svcapitypes.ElasticsearchDomain.
func GenerateElasticsearchDomain(resp *svcsdk.DescribeElasticsearchDomainOutput) *svcapitypes.ElasticsearchDomain {
	cr := &svcapitypes.ElasticsearchDomain{}

	if resp.DomainStatus.ARN != nil {
		cr.Status.AtProvider.ARN = resp.DomainStatus.ARN
	} else {
		cr.Status.AtProvider.ARN = nil
	}
	if resp.DomainStatus.AccessPolicies != nil {
		cr.Spec.ForProvider.AccessPolicies = resp.DomainStatus.AccessPolicies
	} else {
		cr.Spec.ForProvider.AccessPolicies = nil
	}
	if resp.DomainStatus.AdvancedOptions != nil {
		f2 := map[string]*string{}
		for f2key, f2valiter := range resp.DomainStatus.AdvancedOptions {
			var f2val string
			f2val = *f2valiter
			f2[f2key] = &f2val
		}
		cr.Spec.ForProvider.AdvancedOptions = f2
	} else {
		cr.Spec.ForProvider.AdvancedOptions = nil
	}
	if resp.DomainStatus.AdvancedSecurityOptions != nil {
		f3 := &svcapitypes.AdvancedSecurityOptionsInput{}
		if resp.DomainStatus.AdvancedSecurityOptions.Enabled != nil {
			f3.Enabled = resp.DomainStatus.AdvancedSecurityOptions.Enabled
		}
		if resp.DomainStatus.AdvancedSecurityOptions.InternalUserDatabaseEnabled != nil {
			f3.InternalUserDatabaseEnabled = resp.DomainStatus.AdvancedSecurityOptions.InternalUserDatabaseEnabled
		}
		if resp.DomainStatus.AdvancedSecurityOptions.SAMLOptions != nil {
			f3f2 := &svcapitypes.SAMLOptionsInput{}
			if resp.DomainStatus.AdvancedSecurityOptions.SAMLOptions.Enabled != nil {
				f3f2.Enabled = resp.DomainStatus.AdvancedSecurityOptions.SAMLOptions.Enabled
			}
			if resp.DomainStatus.AdvancedSecurityOptions.SAMLOptions.Idp != nil {
				f3f2f1 := &svcapitypes.SAMLIDp{}
				if resp.DomainStatus.AdvancedSecurityOptions.SAMLOptions.Idp.EntityId != nil {
					f3f2f1.EntityID = resp.DomainStatus.AdvancedSecurityOptions.SAMLOptions.Idp.EntityId
				}
				if resp.DomainStatus.AdvancedSecurityOptions.SAMLOptions.Idp.MetadataContent != nil {
					f3f2f1.MetadataContent = resp.DomainStatus.AdvancedSecurityOptions.SAMLOptions.Idp.MetadataContent
				}
				f3f2.IDp = f3f2f1
			}
			if resp.DomainStatus.AdvancedSecurityOptions.SAMLOptions.RolesKey != nil {
				f3f2.RolesKey = resp.DomainStatus.AdvancedSecurityOptions.SAMLOptions.RolesKey
			}
			if resp.DomainStatus.AdvancedSecurityOptions.SAMLOptions.SessionTimeoutMinutes != nil {
				f3f2.SessionTimeoutMinutes = resp.DomainStatus.AdvancedSecurityOptions.SAMLOptions.SessionTimeoutMinutes
			}
			if resp.DomainStatus.AdvancedSecurityOptions.SAMLOptions.SubjectKey != nil {
				f3f2.SubjectKey = resp.DomainStatus.AdvancedSecurityOptions.SAMLOptions.SubjectKey
			}
			f3.SAMLOptions = f3f2
		}
		cr.Spec.ForProvider.AdvancedSecurityOptions = f3
	} else {
		cr.Spec.ForProvider.AdvancedSecurityOptions = nil
	}
	if resp.DomainStatus.CognitoOptions != nil {
		f4 := &svcapitypes.CognitoOptions{}
		if resp.DomainStatus.CognitoOptions.Enabled != nil {
			f4.Enabled = resp.DomainStatus.CognitoOptions.Enabled
		}
		if resp.DomainStatus.CognitoOptions.IdentityPoolId != nil {
			f4.IdentityPoolID = resp.DomainStatus.CognitoOptions.IdentityPoolId
		}
		if resp.DomainStatus.CognitoOptions.RoleArn != nil {
			f4.RoleARN = resp.DomainStatus.CognitoOptions.RoleArn
		}
		if resp.DomainStatus.CognitoOptions.UserPoolId != nil {
			f4.UserPoolID = resp.DomainStatus.CognitoOptions.UserPoolId
		}
		cr.Spec.ForProvider.CognitoOptions = f4
	} else {
		cr.Spec.ForProvider.CognitoOptions = nil
	}
	if resp.DomainStatus.Created != nil {
		cr.Status.AtProvider.Created = resp.DomainStatus.Created
	} else {
		cr.Status.AtProvider.Created = nil
	}
	if resp.DomainStatus.Deleted != nil {
		cr.Status.AtProvider.Deleted = resp.DomainStatus.Deleted
	} else {
		cr.Status.AtProvider.Deleted = nil
	}
	if resp.DomainStatus.DomainEndpointOptions != nil {
		f7 := &svcapitypes.DomainEndpointOptions{}
		if resp.DomainStatus.DomainEndpointOptions.CustomEndpoint != nil {
			f7.CustomEndpoint = resp.DomainStatus.DomainEndpointOptions.CustomEndpoint
		}
		if resp.DomainStatus.DomainEndpointOptions.CustomEndpointCertificateArn != nil {
			f7.CustomEndpointCertificateARN = resp.DomainStatus.DomainEndpointOptions.CustomEndpointCertificateArn
		}
		if resp.DomainStatus.DomainEndpointOptions.CustomEndpointEnabled != nil {
			f7.CustomEndpointEnabled = resp.DomainStatus.DomainEndpointOptions.CustomEndpointEnabled
		}
		if resp.DomainStatus.DomainEndpointOptions.EnforceHTTPS != nil {
			f7.EnforceHTTPS = resp.DomainStatus.DomainEndpointOptions.EnforceHTTPS
		}
		if resp.DomainStatus.DomainEndpointOptions.TLSSecurityPolicy != nil {
			f7.TLSSecurityPolicy = resp.DomainStatus.DomainEndpointOptions.TLSSecurityPolicy
		}
		cr.Spec.ForProvider.DomainEndpointOptions = f7
	} else {
		cr.Spec.ForProvider.DomainEndpointOptions = nil
	}
	if resp.DomainStatus.DomainId != nil {
		cr.Status.AtProvider.DomainID = resp.DomainStatus.DomainId
	} else {
		cr.Status.AtProvider.DomainID = nil
	}
	if resp.DomainStatus.DomainName != nil {
		cr.Spec.ForProvider.DomainName = resp.DomainStatus.DomainName
	} else {
		cr.Spec.ForProvider.DomainName = nil
	}
	if resp.DomainStatus.EBSOptions != nil {
		f10 := &svcapitypes.EBSOptions{}
		if resp.DomainStatus.EBSOptions.EBSEnabled != nil {
			f10.EBSEnabled = resp.DomainStatus.EBSOptions.EBSEnabled
		}
		if resp.DomainStatus.EBSOptions.Iops != nil {
			f10.IOPS = resp.DomainStatus.EBSOptions.Iops
		}
		if resp.DomainStatus.EBSOptions.VolumeSize != nil {
			f10.VolumeSize = resp.DomainStatus.EBSOptions.VolumeSize
		}
		if resp.DomainStatus.EBSOptions.VolumeType != nil {
			f10.VolumeType = resp.DomainStatus.EBSOptions.VolumeType
		}
		cr.Spec.ForProvider.EBSOptions = f10
	} else {
		cr.Spec.ForProvider.EBSOptions = nil
	}
	if resp.DomainStatus.ElasticsearchClusterConfig != nil {
		f11 := &svcapitypes.ElasticsearchClusterConfig{}
		if resp.DomainStatus.ElasticsearchClusterConfig.DedicatedMasterCount != nil {
			f11.DedicatedMasterCount = resp.DomainStatus.ElasticsearchClusterConfig.DedicatedMasterCount
		}
		if resp.DomainStatus.ElasticsearchClusterConfig.DedicatedMasterEnabled != nil {
			f11.DedicatedMasterEnabled = resp.DomainStatus.ElasticsearchClusterConfig.DedicatedMasterEnabled
		}
		if resp.DomainStatus.ElasticsearchClusterConfig.DedicatedMasterType != nil {
			f11.DedicatedMasterType = resp.DomainStatus.ElasticsearchClusterConfig.DedicatedMasterType
		}
		if resp.DomainStatus.ElasticsearchClusterConfig.InstanceCount != nil {
			f11.InstanceCount = resp.DomainStatus.ElasticsearchClusterConfig.InstanceCount
		}
		if resp.DomainStatus.ElasticsearchClusterConfig.InstanceType != nil {
			f11.InstanceType = resp.DomainStatus.ElasticsearchClusterConfig.InstanceType
		}
		if resp.DomainStatus.ElasticsearchClusterConfig.WarmCount != nil {
			f11.WarmCount = resp.DomainStatus.ElasticsearchClusterConfig.WarmCount
		}
		if resp.DomainStatus.ElasticsearchClusterConfig.WarmEnabled != nil {
			f11.WarmEnabled = resp.DomainStatus.ElasticsearchClusterConfig.WarmEnabled
		}
		if resp.DomainStatus.ElasticsearchClusterConfig.WarmType != nil {
			f11.WarmType = resp.DomainStatus.ElasticsearchClusterConfig.WarmType
		}
		if resp.DomainStatus.ElasticsearchClusterConfig.ZoneAwarenessConfig != nil {
			f11f8 := &svcapitypes.ZoneAwarenessConfig{}
			if resp.DomainStatus.ElasticsearchClusterConfig.ZoneAwarenessConfig.AvailabilityZoneCount != nil {
				f11f8.AvailabilityZoneCount = resp.DomainStatus.ElasticsearchClusterConfig.ZoneAwarenessConfig.AvailabilityZoneCount
			}
			f11.ZoneAwarenessConfig = f11f8
		}
		if resp.DomainStatus.ElasticsearchClusterConfig.ZoneAwarenessEnabled != nil {
			f11.ZoneAwarenessEnabled = resp.DomainStatus.ElasticsearchClusterConfig.ZoneAwarenessEnabled
		}
		cr.Spec.ForProvider.ElasticsearchClusterConfig = f11
	} else {
		cr.Spec.ForProvider.ElasticsearchClusterConfig = nil
	}
	if resp.DomainStatus.ElasticsearchVersion != nil {
		cr.Spec.ForProvider.ElasticsearchVersion = resp.DomainStatus.ElasticsearchVersion
	} else {
		cr.Spec.ForProvider.ElasticsearchVersion = nil
	}
	if resp.DomainStatus.EncryptionAtRestOptions != nil {
		f13 := &svcapitypes.EncryptionAtRestOptions{}
		if resp.DomainStatus.EncryptionAtRestOptions.Enabled != nil {
			f13.Enabled = resp.DomainStatus.EncryptionAtRestOptions.Enabled
		}
		if resp.DomainStatus.EncryptionAtRestOptions.KmsKeyId != nil {
			f13.KMSKeyID = resp.DomainStatus.EncryptionAtRestOptions.KmsKeyId
		}
		cr.Spec.ForProvider.EncryptionAtRestOptions = f13
	} else {
		cr.Spec.ForProvider.EncryptionAtRestOptions = nil
	}
	if resp.DomainStatus.Endpoint != nil {
		cr.Status.AtProvider.Endpoint = resp.DomainStatus.Endpoint
	} else {
		cr.Status.AtProvider.Endpoint = nil
	}
	if resp.DomainStatus.Endpoints != nil {
		f15 := map[string]*string{}
		for f15key, f15valiter := range resp.DomainStatus.Endpoints {
			var f15val string
			f15val = *f15valiter
			f15[f15key] = &f15val
		}
		cr.Status.AtProvider.Endpoints = f15
	} else {
		cr.Status.AtProvider.Endpoints = nil
	}
	if resp.DomainStatus.LogPublishingOptions != nil {
		f16 := map[string]*svcapitypes.LogPublishingOption{}
		for f16key, f16valiter := range resp.DomainStatus.LogPublishingOptions {
			f16val := &svcapitypes.LogPublishingOption{}
			if f16valiter.CloudWatchLogsLogGroupArn != nil {
				f16val.CloudWatchLogsLogGroupARN = f16valiter.CloudWatchLogsLogGroupArn
			}
			if f16valiter.Enabled != nil {
				f16val.Enabled = f16valiter.Enabled
			}
			f16[f16key] = f16val
		}
		cr.Spec.ForProvider.LogPublishingOptions = f16
	} else {
		cr.Spec.ForProvider.LogPublishingOptions = nil
	}
	if resp.DomainStatus.NodeToNodeEncryptionOptions != nil {
		f17 := &svcapitypes.NodeToNodeEncryptionOptions{}
		if resp.DomainStatus.NodeToNodeEncryptionOptions.Enabled != nil {
			f17.Enabled = resp.DomainStatus.NodeToNodeEncryptionOptions.Enabled
		}
		cr.Spec.ForProvider.NodeToNodeEncryptionOptions = f17
	} else {
		cr.Spec.ForProvider.NodeToNodeEncryptionOptions = nil
	}
	if resp.DomainStatus.Processing != nil {
		cr.Status.AtProvider.Processing = resp.DomainStatus.Processing
	} else {
		cr.Status.AtProvider.Processing = nil
	}
	if resp.DomainStatus.ServiceSoftwareOptions != nil {
		f19 := &svcapitypes.ServiceSoftwareOptions{}
		if resp.DomainStatus.ServiceSoftwareOptions.AutomatedUpdateDate != nil {
			f19.AutomatedUpdateDate = &metav1.Time{*resp.DomainStatus.ServiceSoftwareOptions.AutomatedUpdateDate}
		}
		if resp.DomainStatus.ServiceSoftwareOptions.Cancellable != nil {
			f19.Cancellable = resp.DomainStatus.ServiceSoftwareOptions.Cancellable
		}
		if resp.DomainStatus.ServiceSoftwareOptions.CurrentVersion != nil {
			f19.CurrentVersion = resp.DomainStatus.ServiceSoftwareOptions.CurrentVersion
		}
		if resp.DomainStatus.ServiceSoftwareOptions.Description != nil {
			f19.Description = resp.DomainStatus.ServiceSoftwareOptions.Description
		}
		if resp.DomainStatus.ServiceSoftwareOptions.NewVersion != nil {
			f19.NewVersion = resp.DomainStatus.ServiceSoftwareOptions.NewVersion
		}
		if resp.DomainStatus.ServiceSoftwareOptions.OptionalDeployment != nil {
			f19.OptionalDeployment = resp.DomainStatus.ServiceSoftwareOptions.OptionalDeployment
		}
		if resp.DomainStatus.ServiceSoftwareOptions.UpdateAvailable != nil {
			f19.UpdateAvailable = resp.DomainStatus.ServiceSoftwareOptions.UpdateAvailable
		}
		if resp.DomainStatus.ServiceSoftwareOptions.UpdateStatus != nil {
			f19.UpdateStatus = resp.DomainStatus.ServiceSoftwareOptions.UpdateStatus
		}
		cr.Status.AtProvider.ServiceSoftwareOptions = f19
	} else {
		cr.Status.AtProvider.ServiceSoftwareOptions = nil
	}
	if resp.DomainStatus.SnapshotOptions != nil {
		f20 := &svcapitypes.SnapshotOptions{}
		if resp.DomainStatus.SnapshotOptions.AutomatedSnapshotStartHour != nil {
			f20.AutomatedSnapshotStartHour = resp.DomainStatus.SnapshotOptions.AutomatedSnapshotStartHour
		}
		cr.Spec.ForProvider.SnapshotOptions = f20
	} else {
		cr.Spec.ForProvider.SnapshotOptions = nil
	}
	if resp.DomainStatus.UpgradeProcessing != nil {
		cr.Status.AtProvider.UpgradeProcessing = resp.DomainStatus.UpgradeProcessing
	} else {
		cr.Status.AtProvider.UpgradeProcessing = nil
	}
	if resp.DomainStatus.VPCOptions != nil {
		f22 := &svcapitypes.VPCOptions{}
		if resp.DomainStatus.VPCOptions.SecurityGroupIds != nil {
			f22f1 := []*string{}
			for _, f22f1iter := range resp.DomainStatus.VPCOptions.SecurityGroupIds {
				var f22f1elem string
				f22f1elem = *f22f1iter
				f22f1 = append(f22f1, &f22f1elem)
			}
			f22.SecurityGroupIDs = f22f1
		}
		if resp.DomainStatus.VPCOptions.SubnetIds != nil {
			f22f2 := []*string{}
			for _, f22f2iter := range resp.DomainStatus.VPCOptions.SubnetIds {
				var f22f2elem string
				f22f2elem = *f22f2iter
				f22f2 = append(f22f2, &f22f2elem)
			}
			f22.SubnetIDs = f22f2
		}
		cr.Spec.ForProvider.VPCOptions = f22
	} else {
		cr.Spec.ForProvider.VPCOptions = nil
	}

	return cr
}

// GenerateCreateElasticsearchDomainInput returns a create input.
func GenerateCreateElasticsearchDomainInput(cr *svcapitypes.ElasticsearchDomain) *svcsdk.CreateElasticsearchDomainInput {
	res := &svcsdk.CreateElasticsearchDomainInput{}

	if cr.Spec.ForProvider.AccessPolicies != nil {
		res.SetAccessPolicies(*cr.Spec.ForProvider.AccessPolicies)
	}
	if cr.Spec.ForProvider.AdvancedOptions != nil {
		f1 := map[string]*string{}
		for f1key, f1valiter := range cr.Spec.ForProvider.AdvancedOptions {
			var f1val string
			f1val = *f1valiter
			f1[f1key] = &f1val
		}
		res.SetAdvancedOptions(f1)
	}
	if cr.Spec.ForProvider.AdvancedSecurityOptions != nil {
		f2 := &svcsdk.AdvancedSecurityOptionsInput{}
		if cr.Spec.ForProvider.AdvancedSecurityOptions.Enabled != nil {
			f2.SetEnabled(*cr.Spec.ForProvider.AdvancedSecurityOptions.Enabled)
		}
		if cr.Spec.ForProvider.AdvancedSecurityOptions.InternalUserDatabaseEnabled != nil {
			f2.SetInternalUserDatabaseEnabled(*cr.Spec.ForProvider.AdvancedSecurityOptions.InternalUserDatabaseEnabled)
		}
		if cr.Spec.ForProvider.AdvancedSecurityOptions.MasterUserOptions != nil {
			f2f2 := &svcsdk.MasterUserOptions{}
			if cr.Spec.ForProvider.AdvancedSecurityOptions.MasterUserOptions.MasterUserARN != nil {
				f2f2.SetMasterUserARN(*cr.Spec.ForProvider.AdvancedSecurityOptions.MasterUserOptions.MasterUserARN)
			}
			if cr.Spec.ForProvider.AdvancedSecurityOptions.MasterUserOptions.MasterUserName != nil {
				f2f2.SetMasterUserName(*cr.Spec.ForProvider.AdvancedSecurityOptions.MasterUserOptions.MasterUserName)
			}
			if cr.Spec.ForProvider.AdvancedSecurityOptions.MasterUserOptions.MasterUserPassword != nil {
				f2f2.SetMasterUserPassword(*cr.Spec.ForProvider.AdvancedSecurityOptions.MasterUserOptions.MasterUserPassword)
			}
			f2.SetMasterUserOptions(f2f2)
		}
		if cr.Spec.ForProvider.AdvancedSecurityOptions.SAMLOptions != nil {
			f2f3 := &svcsdk.SAMLOptionsInput{}
			if cr.Spec.ForProvider.AdvancedSecurityOptions.SAMLOptions.Enabled != nil {
				f2f3.SetEnabled(*cr.Spec.ForProvider.AdvancedSecurityOptions.SAMLOptions.Enabled)
			}
			if cr.Spec.ForProvider.AdvancedSecurityOptions.SAMLOptions.IDp != nil {
				f2f3f1 := &svcsdk.SAMLIdp{}
				if cr.Spec.ForProvider.AdvancedSecurityOptions.SAMLOptions.IDp.EntityID != nil {
					f2f3f1.SetEntityId(*cr.Spec.ForProvider.AdvancedSecurityOptions.SAMLOptions.IDp.EntityID)
				}
				if cr.Spec.ForProvider.AdvancedSecurityOptions.SAMLOptions.IDp.MetadataContent != nil {
					f2f3f1.SetMetadataContent(*cr.Spec.ForProvider.AdvancedSecurityOptions.SAMLOptions.IDp.MetadataContent)
				}
				f2f3.SetIdp(f2f3f1)
			}
			if cr.Spec.ForProvider.AdvancedSecurityOptions.SAMLOptions.MasterBackendRole != nil {
				f2f3.SetMasterBackendRole(*cr.Spec.ForProvider.AdvancedSecurityOptions.SAMLOptions.MasterBackendRole)
			}
			if cr.Spec.ForProvider.AdvancedSecurityOptions.SAMLOptions.MasterUserName != nil {
				f2f3.SetMasterUserName(*cr.Spec.ForProvider.AdvancedSecurityOptions.SAMLOptions.MasterUserName)
			}
			if cr.Spec.ForProvider.AdvancedSecurityOptions.SAMLOptions.RolesKey != nil {
				f2f3.SetRolesKey(*cr.Spec.ForProvider.AdvancedSecurityOptions.SAMLOptions.RolesKey)
			}
			if cr.Spec.ForProvider.AdvancedSecurityOptions.SAMLOptions.SessionTimeoutMinutes != nil {
				f2f3.SetSessionTimeoutMinutes(*cr.Spec.ForProvider.AdvancedSecurityOptions.SAMLOptions.SessionTimeoutMinutes)
			}
			if cr.Spec.ForProvider.AdvancedSecurityOptions.SAMLOptions.SubjectKey != nil {
				f2f3.SetSubjectKey(*cr.Spec.ForProvider.AdvancedSecurityOptions.SAMLOptions.SubjectKey)
			}
			f2.SetSAMLOptions(f2f3)
		}
		res.SetAdvancedSecurityOptions(f2)
	}
	if cr.Spec.ForProvider.CognitoOptions != nil {
		f3 := &svcsdk.CognitoOptions{}
		if cr.Spec.ForProvider.CognitoOptions.Enabled != nil {
			f3.SetEnabled(*cr.Spec.ForProvider.CognitoOptions.Enabled)
		}
		if cr.Spec.ForProvider.CognitoOptions.IdentityPoolID != nil {
			f3.SetIdentityPoolId(*cr.Spec.ForProvider.CognitoOptions.IdentityPoolID)
		}
		if cr.Spec.ForProvider.CognitoOptions.RoleARN != nil {
			f3.SetRoleArn(*cr.Spec.ForProvider.CognitoOptions.RoleARN)
		}
		if cr.Spec.ForProvider.CognitoOptions.UserPoolID != nil {
			f3.SetUserPoolId(*cr.Spec.ForProvider.CognitoOptions.UserPoolID)
		}
		res.SetCognitoOptions(f3)
	}
	if cr.Spec.ForProvider.DomainEndpointOptions != nil {
		f4 := &svcsdk.DomainEndpointOptions{}
		if cr.Spec.ForProvider.DomainEndpointOptions.CustomEndpoint != nil {
			f4.SetCustomEndpoint(*cr.Spec.ForProvider.DomainEndpointOptions.CustomEndpoint)
		}
		if cr.Spec.ForProvider.DomainEndpointOptions.CustomEndpointCertificateARN != nil {
			f4.SetCustomEndpointCertificateArn(*cr.Spec.ForProvider.DomainEndpointOptions.CustomEndpointCertificateARN)
		}
		if cr.Spec.ForProvider.DomainEndpointOptions.CustomEndpointEnabled != nil {
			f4.SetCustomEndpointEnabled(*cr.Spec.ForProvider.DomainEndpointOptions.CustomEndpointEnabled)
		}
		if cr.Spec.ForProvider.DomainEndpointOptions.EnforceHTTPS != nil {
			f4.SetEnforceHTTPS(*cr.Spec.ForProvider.DomainEndpointOptions.EnforceHTTPS)
		}
		if cr.Spec.ForProvider.DomainEndpointOptions.TLSSecurityPolicy != nil {
			f4.SetTLSSecurityPolicy(*cr.Spec.ForProvider.DomainEndpointOptions.TLSSecurityPolicy)
		}
		res.SetDomainEndpointOptions(f4)
	}
	if cr.Spec.ForProvider.DomainName != nil {
		res.SetDomainName(*cr.Spec.ForProvider.DomainName)
	}
	if cr.Spec.ForProvider.EBSOptions != nil {
		f6 := &svcsdk.EBSOptions{}
		if cr.Spec.ForProvider.EBSOptions.EBSEnabled != nil {
			f6.SetEBSEnabled(*cr.Spec.ForProvider.EBSOptions.EBSEnabled)
		}
		if cr.Spec.ForProvider.EBSOptions.IOPS != nil {
			f6.SetIops(*cr.Spec.ForProvider.EBSOptions.IOPS)
		}
		if cr.Spec.ForProvider.EBSOptions.VolumeSize != nil {
			f6.SetVolumeSize(*cr.Spec.ForProvider.EBSOptions.VolumeSize)
		}
		if cr.Spec.ForProvider.EBSOptions.VolumeType != nil {
			f6.SetVolumeType(*cr.Spec.ForProvider.EBSOptions.VolumeType)
		}
		res.SetEBSOptions(f6)
	}
	if cr.Spec.ForProvider.ElasticsearchClusterConfig != nil {
		f7 := &svcsdk.ElasticsearchClusterConfig{}
		if cr.Spec.ForProvider.ElasticsearchClusterConfig.DedicatedMasterCount != nil {
			f7.SetDedicatedMasterCount(*cr.Spec.ForProvider.ElasticsearchClusterConfig.DedicatedMasterCount)
		}
		if cr.Spec.ForProvider.ElasticsearchClusterConfig.DedicatedMasterEnabled != nil {
			f7.SetDedicatedMasterEnabled(*cr.Spec.ForProvider.ElasticsearchClusterConfig.DedicatedMasterEnabled)
		}
		if cr.Spec.ForProvider.ElasticsearchClusterConfig.DedicatedMasterType != nil {
			f7.SetDedicatedMasterType(*cr.Spec.ForProvider.ElasticsearchClusterConfig.DedicatedMasterType)
		}
		if cr.Spec.ForProvider.ElasticsearchClusterConfig.InstanceCount != nil {
			f7.SetInstanceCount(*cr.Spec.ForProvider.ElasticsearchClusterConfig.InstanceCount)
		}
		if cr.Spec.ForProvider.ElasticsearchClusterConfig.InstanceType != nil {
			f7.SetInstanceType(*cr.Spec.ForProvider.ElasticsearchClusterConfig.InstanceType)
		}
		if cr.Spec.ForProvider.ElasticsearchClusterConfig.WarmCount != nil {
			f7.SetWarmCount(*cr.Spec.ForProvider.ElasticsearchClusterConfig.WarmCount)
		}
		if cr.Spec.ForProvider.ElasticsearchClusterConfig.WarmEnabled != nil {
			f7.SetWarmEnabled(*cr.Spec.ForProvider.ElasticsearchClusterConfig.WarmEnabled)
		}
		if cr.Spec.ForProvider.ElasticsearchClusterConfig.WarmType != nil {
			f7.SetWarmType(*cr.Spec.ForProvider.ElasticsearchClusterConfig.WarmType)
		}
		if cr.Spec.ForProvider.ElasticsearchClusterConfig.ZoneAwarenessConfig != nil {
			f7f8 := &svcsdk.ZoneAwarenessConfig{}
			if cr.Spec.ForProvider.ElasticsearchClusterConfig.ZoneAwarenessConfig.AvailabilityZoneCount != nil {
				f7f8.SetAvailabilityZoneCount(*cr.Spec.ForProvider.ElasticsearchClusterConfig.ZoneAwarenessConfig.AvailabilityZoneCount)
			}
			f7.SetZoneAwarenessConfig(f7f8)
		}
		if cr.Spec.ForProvider.ElasticsearchClusterConfig.ZoneAwarenessEnabled != nil {
			f7.SetZoneAwarenessEnabled(*cr.Spec.ForProvider.ElasticsearchClusterConfig.ZoneAwarenessEnabled)
		}
		res.SetElasticsearchClusterConfig(f7)
	}
	if cr.Spec.ForProvider.ElasticsearchVersion != nil {
		res.SetElasticsearchVersion(*cr.Spec.ForProvider.ElasticsearchVersion)
	}
	if cr.Spec.ForProvider.EncryptionAtRestOptions != nil {
		f9 := &svcsdk.EncryptionAtRestOptions{}
		if cr.Spec.ForProvider.EncryptionAtRestOptions.Enabled != nil {
			f9.SetEnabled(*cr.Spec.ForProvider.EncryptionAtRestOptions.Enabled)
		}
		if cr.Spec.ForProvider.EncryptionAtRestOptions.KMSKeyID != nil {
			f9.SetKmsKeyId(*cr.Spec.ForProvider.EncryptionAtRestOptions.KMSKeyID)
		}
		res.SetEncryptionAtRestOptions(f9)
	}
	if cr.Spec.ForProvider.LogPublishingOptions != nil {
		f10 := map[string]*svcsdk.LogPublishingOption{}
		for f10key, f10valiter := range cr.Spec.ForProvider.LogPublishingOptions {
			f10val := &svcsdk.LogPublishingOption{}
			if f10valiter.CloudWatchLogsLogGroupARN != nil {
				f10val.SetCloudWatchLogsLogGroupArn(*f10valiter.CloudWatchLogsLogGroupARN)
			}
			if f10valiter.Enabled != nil {
				f10val.SetEnabled(*f10valiter.Enabled)
			}
			f10[f10key] = f10val
		}
		res.SetLogPublishingOptions(f10)
	}
	if cr.Spec.ForProvider.NodeToNodeEncryptionOptions != nil {
		f11 := &svcsdk.NodeToNodeEncryptionOptions{}
		if cr.Spec.ForProvider.NodeToNodeEncryptionOptions.Enabled != nil {
			f11.SetEnabled(*cr.Spec.ForProvider.NodeToNodeEncryptionOptions.Enabled)
		}
		res.SetNodeToNodeEncryptionOptions(f11)
	}
	if cr.Spec.ForProvider.SnapshotOptions != nil {
		f12 := &svcsdk.SnapshotOptions{}
		if cr.Spec.ForProvider.SnapshotOptions.AutomatedSnapshotStartHour != nil {
			f12.SetAutomatedSnapshotStartHour(*cr.Spec.ForProvider.SnapshotOptions.AutomatedSnapshotStartHour)
		}
		res.SetSnapshotOptions(f12)
	}
	if cr.Spec.ForProvider.VPCOptions != nil {
		f13 := &svcsdk.VPCOptions{}
		if cr.Spec.ForProvider.VPCOptions.SecurityGroupIDs != nil {
			f13f0 := []*string{}
			for _, f13f0iter := range cr.Spec.ForProvider.VPCOptions.SecurityGroupIDs {
				var f13f0elem string
				f13f0elem = *f13f0iter
				f13f0 = append(f13f0, &f13f0elem)
			}
			f13.SetSecurityGroupIds(f13f0)
		}
		if cr.Spec.ForProvider.VPCOptions.SubnetIDs != nil {
			f13f1 := []*string{}
			for _, f13f1iter := range cr.Spec.ForProvider.VPCOptions.SubnetIDs {
				var f13f1elem string
				f13f1elem = *f13f1iter
				f13f1 = append(f13f1, &f13f1elem)
			}
			f13.SetSubnetIds(f13f1)
		}
		res.SetVPCOptions(f13)
	}

	return res
}

// GenerateDeleteElasticsearchDomainInput returns a deletion input.
func GenerateDeleteElasticsearchDomainInput(cr *svcapitypes.ElasticsearchDomain) *svcsdk.DeleteElasticsearchDomainInput {
	res := &svcsdk.DeleteElasticsearchDomainInput{}

	if cr.Spec.ForProvider.DomainName != nil {
		res.SetDomainName(*cr.Spec.ForProvider.DomainName)
	}

	return res
}

// IsNotFound returns whether the given error is of type NotFound or not.
func IsNotFound(err error) bool {
	awsErr, ok := err.(awserr.Error)
	return ok && awsErr.Code() == "UNKNOWN"
}
