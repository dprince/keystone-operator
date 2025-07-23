/*
Copyright 2022.

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

package v1beta1

import (
	"fmt"

	keystonev2 "github.com/openstack-k8s-operators/keystone-operator/api/v1beta2"
	"github.com/openstack-k8s-operators/lib-common/modules/common/service"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts this KeystoneAPI to the Hub version (v1beta2).
func (src *KeystoneAPI) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*keystonev2.KeystoneAPI)

	// Copy metadata
	dst.ObjectMeta = src.ObjectMeta

	// Convert Spec
	if err := convertKeystoneAPISpecTo(&src.Spec, &dst.Spec); err != nil {
		return fmt.Errorf("error converting KeystoneAPI spec to v1beta2: %w", err)
	}

	// Convert Status - most fields are the same
	dst.Status.ReadyCount = src.Status.ReadyCount
	dst.Status.Hash = make(map[string]string)
	for k, v := range src.Status.Hash {
		dst.Status.Hash[k] = v
	}
	dst.Status.APIEndpoints = make(map[string]string)
	for k, v := range src.Status.APIEndpoints {
		dst.Status.APIEndpoints[k] = v
	}
	dst.Status.Conditions = src.Status.Conditions
	dst.Status.DatabaseHostname = src.Status.DatabaseHostname
	dst.Status.NetworkAttachments = make(map[string][]string)
	for k, v := range src.Status.NetworkAttachments {
		dst.Status.NetworkAttachments[k] = v
	}
	dst.Status.TransportURLSecret = src.Status.TransportURLSecret
	dst.Status.ObservedGeneration = src.Status.ObservedGeneration
	dst.Status.LastAppliedTopology = src.Status.LastAppliedTopology

	return nil
}

// ConvertFrom converts from the Hub version (v1beta2) to this version.
func (dst *KeystoneAPI) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*keystonev2.KeystoneAPI)

	// Copy metadata
	dst.ObjectMeta = src.ObjectMeta

	// Convert Spec
	if err := convertKeystoneAPISpecFrom(&src.Spec, &dst.Spec); err != nil {
		return fmt.Errorf("error converting KeystoneAPI spec from v1beta2: %w", err)
	}

	// Convert Status - most fields are the same
	dst.Status.ReadyCount = src.Status.ReadyCount
	dst.Status.Hash = make(map[string]string)
	for k, v := range src.Status.Hash {
		dst.Status.Hash[k] = v
	}
	dst.Status.APIEndpoints = make(map[string]string)
	for k, v := range src.Status.APIEndpoints {
		dst.Status.APIEndpoints[k] = v
	}
	dst.Status.Conditions = src.Status.Conditions
	dst.Status.DatabaseHostname = src.Status.DatabaseHostname
	dst.Status.NetworkAttachments = make(map[string][]string)
	for k, v := range src.Status.NetworkAttachments {
		dst.Status.NetworkAttachments[k] = v
	}
	dst.Status.TransportURLSecret = src.Status.TransportURLSecret
	dst.Status.ObservedGeneration = src.Status.ObservedGeneration
	dst.Status.LastAppliedTopology = src.Status.LastAppliedTopology
	// v1beta1 doesn't have Region in status, so we store it during conversion
	if src.Status.Region != "" {
		dst.Status.Region = src.Status.Region
	}

	return nil
}

// convertKeystoneAPISpecTo converts v1beta1 KeystoneAPISpec to v1beta2 KeystoneAPISpec
func convertKeystoneAPISpecTo(src *KeystoneAPISpec, dst *keystonev2.KeystoneAPISpec) error {
	// Copy container image
	if src.ContainerImage != "" {
		dst.ContainerImage = src.ContainerImage
	}
	// Convert core spec
	return convertKeystoneAPISpecCoreTo(&src.KeystoneAPISpecCore, &dst.KeystoneAPISpecCore)
}

// convertKeystoneAPISpecFrom converts v1beta2 KeystoneAPISpec to v1beta1 KeystoneAPISpec
func convertKeystoneAPISpecFrom(src *keystonev2.KeystoneAPISpec, dst *KeystoneAPISpec) error {
	// Copy container image
	if src.ContainerImage != "" {
		dst.ContainerImage = src.ContainerImage
	}

	// Convert core spec
	return convertKeystoneAPISpecCoreFrom(&src.KeystoneAPISpecCore, &dst.KeystoneAPISpecCore)
}

// convertKeystoneAPISpecCoreTo converts v1beta1 KeystoneAPISpecCore to v1beta2 KeystoneAPISpecCore
func convertKeystoneAPISpecCoreTo(src *KeystoneAPISpecCore, dst *keystonev2.KeystoneAPISpecCore) error {
	// Main difference: DatabaseInstance -> DatabaseName
	dst.DatabaseName = src.DatabaseInstance

	// Copy all other fields
	dst.DatabaseAccount = src.DatabaseAccount
	dst.MemcachedInstance = src.MemcachedInstance
	dst.Region = src.Region
	dst.AdminProject = src.AdminProject
	dst.AdminUser = src.AdminUser
	dst.Replicas = src.Replicas
	dst.Secret = src.Secret
	dst.EnableSecureRBAC = src.EnableSecureRBAC
	dst.TrustFlushArgs = src.TrustFlushArgs
	dst.TrustFlushSchedule = src.TrustFlushSchedule
	dst.TrustFlushSuspend = src.TrustFlushSuspend
	dst.FernetRotationDays = src.FernetRotationDays
	dst.FernetMaxActiveKeys = src.FernetMaxActiveKeys
	dst.PasswordSelectors = keystonev2.PasswordSelector{
		Admin: src.PasswordSelectors.Admin,
	}
	dst.NodeSelector = src.NodeSelector
	dst.PreserveJobs = src.PreserveJobs
	dst.CustomServiceConfig = src.CustomServiceConfig
	dst.DefaultConfigOverwrite = make(map[string]string)
	for k, v := range src.DefaultConfigOverwrite {
		dst.DefaultConfigOverwrite[k] = v
	}

	// Convert HttpdCustomization
	dst.HttpdCustomization = keystonev2.HttpdCustomization{
		ProcessNumber:      src.HttpdCustomization.ProcessNumber,
		CustomConfigSecret: src.HttpdCustomization.CustomConfigSecret,
	}

	dst.Resources = src.Resources
	dst.NetworkAttachments = make([]string, len(src.NetworkAttachments))
	copy(dst.NetworkAttachments, src.NetworkAttachments)

	// Convert Override
	dst.Override = keystonev2.APIOverrideSpec{
		Service: make(map[service.Endpoint]service.RoutedOverrideSpec),
	}
	// Convert service override specs
	for k, v := range src.Override.Service {
		dst.Override.Service[k] = v
	}

	dst.RabbitMqClusterName = src.RabbitMqClusterName

	// Convert TLS - types are the same between versions
	dst.TLS = src.TLS

	dst.APITimeout = src.APITimeout
	dst.TopologyRef = src.TopologyRef

	// Convert ExtraMounts
	dst.ExtraMounts = make([]keystonev2.KeystoneExtraMounts, len(src.ExtraMounts))
	for i, mount := range src.ExtraMounts {
		dst.ExtraMounts[i] = keystonev2.KeystoneExtraMounts{
			Name:      mount.Name,
			Region:    mount.Region,
			VolMounts: mount.VolMounts,
		}
	}

	dst.FederatedRealmConfig = src.FederatedRealmConfig

	// v1beta2 has a new field FederationMountPath that doesn't exist in v1beta1
	// We'll set it to the default value if not specified
	dst.FederationMountPath = "/etc/httpd/conf"

	return nil
}

// convertKeystoneAPISpecCoreFrom converts v1beta2 KeystoneAPISpecCore to v1beta1 KeystoneAPISpecCore
func convertKeystoneAPISpecCoreFrom(src *keystonev2.KeystoneAPISpecCore, dst *KeystoneAPISpecCore) error {
	// Main difference: DatabaseName -> DatabaseInstance
	dst.DatabaseInstance = src.DatabaseName

	// Copy all other fields
	dst.DatabaseAccount = src.DatabaseAccount
	dst.MemcachedInstance = src.MemcachedInstance
	dst.Region = src.Region
	dst.AdminProject = src.AdminProject
	dst.AdminUser = src.AdminUser
	dst.Replicas = src.Replicas
	dst.Secret = src.Secret
	dst.EnableSecureRBAC = src.EnableSecureRBAC
	dst.TrustFlushArgs = src.TrustFlushArgs
	dst.TrustFlushSchedule = src.TrustFlushSchedule
	dst.TrustFlushSuspend = src.TrustFlushSuspend
	dst.FernetRotationDays = src.FernetRotationDays
	dst.FernetMaxActiveKeys = src.FernetMaxActiveKeys
	dst.PasswordSelectors = PasswordSelector{
		Admin: src.PasswordSelectors.Admin,
	}
	dst.NodeSelector = src.NodeSelector
	dst.PreserveJobs = src.PreserveJobs
	dst.CustomServiceConfig = src.CustomServiceConfig
	dst.DefaultConfigOverwrite = make(map[string]string)
	for k, v := range src.DefaultConfigOverwrite {
		dst.DefaultConfigOverwrite[k] = v
	}

	// Convert HttpdCustomization
	dst.HttpdCustomization = HttpdCustomization{
		ProcessNumber:      src.HttpdCustomization.ProcessNumber,
		CustomConfigSecret: src.HttpdCustomization.CustomConfigSecret,
	}

	dst.Resources = src.Resources
	dst.NetworkAttachments = make([]string, len(src.NetworkAttachments))
	copy(dst.NetworkAttachments, src.NetworkAttachments)

	// Convert Override
	dst.Override = APIOverrideSpec{
		Service: make(map[service.Endpoint]service.RoutedOverrideSpec),
	}
	for k, v := range src.Override.Service {
		dst.Override.Service[k] = v
	}

	dst.RabbitMqClusterName = src.RabbitMqClusterName

	// Convert TLS - types are the same between versions
	dst.TLS = src.TLS

	dst.APITimeout = src.APITimeout
	if src.TopologyRef != nil {
		dst.TopologyRef = src.TopologyRef
	}

	// Convert ExtraMounts
	dst.ExtraMounts = make([]KeystoneExtraMounts, len(src.ExtraMounts))
	for i, mount := range src.ExtraMounts {
		dst.ExtraMounts[i] = KeystoneExtraMounts{
			Name:      mount.Name,
			Region:    mount.Region,
			VolMounts: mount.VolMounts,
		}
	}

	dst.FederatedRealmConfig = src.FederatedRealmConfig

	// Note: v1beta2's FederationMountPath field is ignored during conversion to v1beta1
	// since v1beta1 doesn't have this field

	return nil
}
