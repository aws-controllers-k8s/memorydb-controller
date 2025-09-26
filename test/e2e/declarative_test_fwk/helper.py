# Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"). You may
# not use this file except in compliance with the License. A copy of the
# License is located at
#
#	 http://aws.amazon.com/apache2.0/
#
# or in the "license" file accompanying this file. This file is distributed
# on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
# express or implied. See the License for the specific language governing
# permissions and limitations under the License.

"""Helper for Declarative tests framework for custom resources
"""

from e2e.declarative_test_fwk import model

import boto3
import logging
from typing import Tuple
from time import sleep
from acktest.k8s import resource as k8s

# holds custom resource helper references
TEST_HELPERS = dict()


def register_resource_helper(resource_kind: str, resource_plural: str):
    """Decorator to discover Custom Resource Helper

    Args:
        resource_kind: custom resource kind
        resource_plural: custom resource kind plural

    Returns:
        wrapper
    """

    def registrar(cls):
        global TEST_HELPERS
        if issubclass(cls, ResourceHelper):
            TEST_HELPERS[resource_kind.lower()] = cls
            cls.resource_plural = resource_plural.lower()
            logging.info(f"Registered ResourceHelper: {cls.__name__} for custom resource kind: {resource_kind}")
        else:
            msg = f"Unable to register helper for {resource_kind} resource: {cls} is not a subclass of ResourceHelper"
            logging.error(msg)
            raise Exception(msg)
    return registrar


class ResourceHelper:
    """Provides generic verb (create, patch, delete) methods for custom resources.
    Keep its methods stateless. Methods are on instance to allow specialization.
    """

    DEFAULT_WAIT_SECS = 30

    def __init__(self):
        self.mdb = boto3.client("memorydb")

    def create(self, input_data: dict, input_replacements: dict = {}) -> Tuple[k8s.CustomResourceReference, dict]:
        """Creates custom resource inside Kubernetes cluster per the specifications in input data.

        Args:
            input_data: custom resource details
            input_replacements: input replacements

        Returns:
            k8s.CustomResourceReference, created custom resource
        """

        reference = self.custom_resource_reference(input_data, input_replacements)
        _ = k8s.create_custom_resource(reference, input_data)
        resource = k8s.wait_resource_consumed_by_controller(reference, wait_periods=10)
        assert resource is not None
        return reference, resource

    def patch(self, input_data: dict, input_replacements: dict = {}) -> Tuple[k8s.CustomResourceReference, dict]:
        """Patches custom resource inside Kubernetes cluster per the specifications in input data.

        Args:
            input_data: custom resource patch details
            input_replacements: input replacements

        Returns:
            k8s.CustomResourceReference, created custom resource
        """

        reference = self.custom_resource_reference(input_data, input_replacements)
        _ = k8s.patch_custom_resource(reference, input_data)
        sleep(self.DEFAULT_WAIT_SECS)  # required as controller has likely not placed the resource in modifying
        resource = k8s.wait_resource_consumed_by_controller(reference, wait_periods=10)
        assert resource is not None
        return reference, resource

    def delete(self, reference: k8s.CustomResourceReference) -> None:
        """Deletes custom resource inside Kubernetes cluster and waits for delete completion

        Args:
            reference: custom resource reference

        Returns:
            None
        """

        resource = k8s.get_resource(reference)
        if not resource:
            logging.warning(f"ResourceReference {reference} not found. Not invoking k8s delete api.")
            return

        k8s.delete_custom_resource(reference, wait_periods=30, period_length=60)  # throws exception if wait fails
        sleep(self.DEFAULT_WAIT_SECS)
        self.wait_for_delete(reference)  # throws exception if wait fails

    def assert_k8s_resource(self, expectations_k8s: model.ExpectK8SDict, reference: k8s.CustomResourceReference) -> None:
        """Compare the supplied expectations with custom resource reference inside Kubernetes cluster

        :param expectations_k8s: expectations to assert
        :param reference: custom resource reference

        :return: None
        """
        self._assert_conditions(expectations_k8s, reference, wait=False)
        # conditions expectations met, now check current resource against expectations_k8s
        resource = k8s.get_resource(reference)
        # status and spec assertions with resources from Kubernetes
        self.assert_items_k8s(expectations_k8s.get("status"), resource.get("status"))
        self.assert_items_k8s(expectations_k8s.get("spec"), resource.get("spec"))

    def assert_aws_resource(self, expectations_aws: dict, reference: k8s.CustomResourceReference) -> None:
        """Compare the supplied expectations with aws resource from api calls by boto3

         Args:
            expectations_aws: expectations to assert
            reference: custom resource reference

        Returns:
            None
        """
        resource = k8s.get_resource(reference)
        # assertions with resources from boto3
        # spec of k8s resource here is for getting name of resource
        # boto3 call APIs using name of resource
        self.assert_items_aws(expectations_aws, resource.get("spec"))

    def wait_for(self, wait_expectations: dict, reference: k8s.CustomResourceReference) -> None:
        """Waits for custom resource reference details inside Kubernetes cluster to match supplied config,
        currently supports wait on "status.conditions",
        it can be enhanced later for wait on any/other properties.

        Args:
            wait_expectations: properties to wait for
            reference:  custom resource reference

        Returns:
            None
        """

        # wait for conditions
        self._assert_conditions(wait_expectations, reference, wait=True)

    def _assert_conditions(self, expectations: dict, reference: k8s.CustomResourceReference, wait: bool = True) -> None:
        expect_conditions: dict = {}
        if "status" in expectations and "conditions" in expectations["status"]:
            expect_conditions = expectations["status"]["conditions"]

        default_wait_periods = 60
        # period_length = 1 will result in condition check every second
        default_period_length = 1
        for (condition_name, expected_value) in expect_conditions.items():
            if type(expected_value) is str:
                # Example: ACK.Terminal: "True"
                if wait:
                    assert k8s.wait_on_condition(reference, condition_name, expected_value,
                                                 wait_periods=default_wait_periods, period_length=default_period_length)
                else:
                    k8s_resource_condition = k8s.get_resource_condition(reference, condition_name)
                    assert k8s_resource_condition is not None
                    assert expected_value == k8s_resource_condition.get("status"), f"Condition status mismatch. Expected condition: {condition_name} - {expected_value} but found {k8s_resource_condition}"

            elif type(expected_value) is dict:
                # Example:
                # Ready:
                #     status: "False"
                #     message: "Expected message ..."
                #     timeout: 60 # seconds
                condition_value = expected_value.get("status")
                condition_message = expected_value.get("message")
                condition_reason = expected_value.get("reason")
                # default wait 60 seconds
                wait_timeout = expected_value.get("timeout", default_wait_periods)

                if wait:
                    assert k8s.wait_on_condition(reference, condition_name, condition_value,
                                                 wait_periods=wait_timeout, period_length=default_period_length, desired_condition_reason=condition_reason)

                k8s_resource_condition = k8s.get_resource_condition(reference,
                                                              condition_name)
                assert k8s_resource_condition is not None
                assert condition_value == k8s_resource_condition.get("status"), f"Condition status mismatch. Expected condition: {condition_name} - {expected_value} but found {k8s_resource_condition}"
                if condition_message is not None:
                    assert condition_message == k8s_resource_condition.get("message"), f"Condition message mismatch. Expected condition: {condition_name} - {expected_value} but found {k8s_resource_condition}"
                if condition_reason is not None:
                    assert condition_reason in k8s_resource_condition.get("reason"), f"Condition reason mismatch. Expected condition: {condition_name} - {expected_value} but found {k8s_resource_condition}"

            else:
                raise Exception(f"Condition {condition_name} is provided with invalid value: {expected_value} ")

    def assert_items_aws(self, expectations: dict, k8s_resource: dict) -> None:
        """Asserts boto3 response against supplied expectations

        Args:
            expectations: dictionary with items (expected) to assert in k8s resource
            k8s_resource: the current status/spec of the k8s resource

        Returns:
            None
        """
        if not expectations:
            # nothing to assert as there are no expectations
            return
        if not k8s_resource:
            # there are expectations but no given k8s resource state to validate
            # following assert will fail and assert introspection will provide useful information for debugging
            assert expectations == k8s_resource

        resource_name = ""
        for (key, value) in k8s_resource.items():
            if key == "name":
                resource_name = value
        # get aws resource
        aws_resource = self.get_aws_resource(resource_name)
        latestTags = self.mdb.list_tags(ResourceArn=aws_resource.get("ARN"))['TagList']

        for (key, value) in expectations.items():
            if key == "Tags":
                self.assert_tags(value, latestTags)
                continue
            # validation for specific fields
            if self.assert_extra_items_aws(key, value, aws_resource):
                continue
            assert value == aws_resource.get(key)

    def get_aws_resource(self, resource_name):
        """get aws resource from boto3
        Override it for each resource type

        Args:
            resource_name: the name of resource to use for calling memorydb apis

        Returns:
            aws resource
        """
        return

    def assert_extra_items_aws(self, expected_name, expected_value, aws_resource) -> bool:
        """Asserts specific fields that boto3 response against supplied expectation
        Override it as needed for validations of specific field

        Args:
            expected_name: expected name of a specific field to assert in aws resource
            expected_value: expected value of a specific field to assert in aws resource
            aws_resource: specific resource from aws memorydb

        Returns:
            The return value. If resource has expected_name as a field, execute custom logic to validate that specific
            field and return true. Otherwise, if resource doesn't have expected_name as a field, return false.
        """
        return False

    def assert_items_k8s(self, expectations: dict, k8s_resource: dict) -> None:
        """Asserts k8s resource against supplied expectations

        Args:
            expectations: dictionary with items (expected) to assert in k8s_resource
            k8s_resource: the current status/spec of the k8s resource

        Returns:
            None
        """

        if not expectations:
            # nothing to assert as there are no expectations
            return
        if not k8s_resource:
            # there are expectations but no given k8s_resource state to validate
            # following assert will fail and assert introspection will provide useful information for debugging
            assert expectations == k8s_resource

        for (key, value) in expectations.items():
            # conditions are processed separately
            if key == "conditions":
                continue
            # tags are processed by custom logic
            if key == "tags":
                self.assert_tags(value, k8s_resource.get(key))
                continue
            # assert any specific extra fields
            if self.assert_extra_items_k8s(key, value, k8s_resource):
                continue
            logging.info(f"Asserting {key} with value {value}")
            logging.info(f"Resource: {k8s_resource}")
            logging.info(f"Resource key: {k8s_resource.get(key)}")
            logging.info(f"//////////////// Expected value: {value}")
            assert value == k8s_resource.get(key)

    def assert_extra_items_k8s(self, expected_name, expected_value, k8s_resource) -> bool:
        """Asserts extra fields from k8s resource against fields from expectations
        Override it as needed for custom verifications

        Args:
            expected_name: expected name of a specific field to assert in k8s resource
            expected_value: expected value of a specific field to assert in k8s resource
            k8s_resource: the current status/spec of the k8s resource

        Returns:
            The return value. If resource has expected_name as a field, execute custom logic to validate that specific
            field and return true. Otherwise, if resource doesn't have expected_name as a field, return false.
        """
        return False

    def assert_tags(self, expected_tags, latest_tags):
        """Asserts tags from resource against tags from expectations
        If latest_tags contains all tags in expected_tags, no AssertError.
        If any tag in expected_tags is not in latest_tags, return an AssertError

        Args:
            expected_tags: tags from expectations
            latest_tags: tags from resource

        Returns:
            None
        """
        for expected_tag in expected_tags:
            tag_exist = False
            for latest_tag in latest_tags:
                if expected_tag == latest_tag:
                    tag_exist = True
                    break
            if tag_exist:
                continue
            assert expected_tag == None

    def custom_resource_reference(self, input_data: dict, input_replacements: dict = {}) -> k8s.CustomResourceReference:
        """Helper method to provide k8s.CustomResourceReference for supplied input

        Args:
            input_data: custom resource input data
            input_replacements: input replacements

        Returns:
            k8s.CustomResourceReference
        """

        resource_name = input_data.get("metadata").get("name")
        crd_group = input_replacements.get("CRD_GROUP")
        crd_version = input_replacements.get("CRD_VERSION")

        reference = k8s.CustomResourceReference(
            crd_group, crd_version, self.resource_plural, resource_name, namespace="default")
        return reference

    def wait_for_delete(self, reference: k8s.CustomResourceReference) -> None:
        """Override this method to implement custom wail logic on resource delete.

        Args:
            reference: custom resource reference

        Returns:
            None
        """

        logging.debug(f"No-op wait_for_delete()")


def get_resource_helper(resource_kind: str) -> ResourceHelper:
    """Provides ResourceHelper for supplied custom resource kind
    If no helper is registered for the supplied resource kind then returns default ResourceHelper

    Args:
        resource_kind: custom resource kind string

    Returns:
        custom resource helper instance
    """

    helper_cls = TEST_HELPERS.get(resource_kind.lower())
    if helper_cls:
        return helper_cls()
    return ResourceHelper()
