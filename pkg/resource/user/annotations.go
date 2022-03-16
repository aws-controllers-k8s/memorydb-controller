// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package user

import (
	svcapitypes "github.com/aws-controllers-k8s/memorydb-controller/apis/v1alpha1"
)

const (
	// AnnotationLastRequestedAccessString is an annotation whose value is the string
	// passed in as input to either the create or modify API called most recently
	AnnotationLastRequestedAccessString = svcapitypes.AnnotationPrefix + "last-requested-access-string"
)

// setAnnotationsFields sets the required annotations for user resource
// This should only be called upon a successful create or modify call.
func (rm *resourceManager) setAnnotationsFields(
	r *resource,
	ko *svcapitypes.User,
) {
	if ko.ObjectMeta.Annotations == nil {
		ko.ObjectMeta.Annotations = make(map[string]string)
	}

	setLastRequestedAccessString(r, ko.ObjectMeta.Annotations)
}

// setLastRequestedAccessString copies desired.Spec.AccessString into the annotation
// of the object.
func setLastRequestedAccessString(
	r *resource,
	annotations map[string]string,
) {
	if r.ko.Spec.AccessString != nil {
		annotations[AnnotationLastRequestedAccessString] = *r.ko.Spec.AccessString
	}
}
