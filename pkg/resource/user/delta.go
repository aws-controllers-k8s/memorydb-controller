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

// Code generated by ack-generate. DO NOT EDIT.

package user

import (
	"bytes"
	"reflect"

	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
)

// Hack to avoid import errors during build...
var (
	_ = &bytes.Buffer{}
	_ = &reflect.Method{}
)

// newResourceDelta returns a new `ackcompare.Delta` used to compare two
// resources
func newResourceDelta(
	a *resource,
	b *resource,
) *ackcompare.Delta {
	delta := ackcompare.NewDelta()
	if (a == nil && b != nil) ||
		(a != nil && b == nil) {
		delta.Add("", a, b)
		return delta
	}

	if ackcompare.HasNilDifference(a.ko.Spec.AccessString, b.ko.Spec.AccessString) {
		delta.Add("Spec.AccessString", a.ko.Spec.AccessString, b.ko.Spec.AccessString)
	} else if a.ko.Spec.AccessString != nil && b.ko.Spec.AccessString != nil {
		if *a.ko.Spec.AccessString != *b.ko.Spec.AccessString {
			delta.Add("Spec.AccessString", a.ko.Spec.AccessString, b.ko.Spec.AccessString)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.AuthenticationMode, b.ko.Spec.AuthenticationMode) {
		delta.Add("Spec.AuthenticationMode", a.ko.Spec.AuthenticationMode, b.ko.Spec.AuthenticationMode)
	} else if a.ko.Spec.AuthenticationMode != nil && b.ko.Spec.AuthenticationMode != nil {
		if !ackcompare.SliceStringPEqual(a.ko.Spec.AuthenticationMode.Passwords, b.ko.Spec.AuthenticationMode.Passwords) {
			delta.Add("Spec.AuthenticationMode.Passwords", a.ko.Spec.AuthenticationMode.Passwords, b.ko.Spec.AuthenticationMode.Passwords)
		}
		if ackcompare.HasNilDifference(a.ko.Spec.AuthenticationMode.Type, b.ko.Spec.AuthenticationMode.Type) {
			delta.Add("Spec.AuthenticationMode.Type", a.ko.Spec.AuthenticationMode.Type, b.ko.Spec.AuthenticationMode.Type)
		} else if a.ko.Spec.AuthenticationMode.Type != nil && b.ko.Spec.AuthenticationMode.Type != nil {
			if *a.ko.Spec.AuthenticationMode.Type != *b.ko.Spec.AuthenticationMode.Type {
				delta.Add("Spec.AuthenticationMode.Type", a.ko.Spec.AuthenticationMode.Type, b.ko.Spec.AuthenticationMode.Type)
			}
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.Name, b.ko.Spec.Name) {
		delta.Add("Spec.Name", a.ko.Spec.Name, b.ko.Spec.Name)
	} else if a.ko.Spec.Name != nil && b.ko.Spec.Name != nil {
		if *a.ko.Spec.Name != *b.ko.Spec.Name {
			delta.Add("Spec.Name", a.ko.Spec.Name, b.ko.Spec.Name)
		}
	}
	if !reflect.DeepEqual(a.ko.Spec.Tags, b.ko.Spec.Tags) {
		delta.Add("Spec.Tags", a.ko.Spec.Tags, b.ko.Spec.Tags)
	}

	return delta
}