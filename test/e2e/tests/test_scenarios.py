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

"""
Tests for custom resources.
Uses declarative tests framework for custom resources.

To add test: add scenario yaml to scenarios/ directory.
"""

from e2e.declarative_test_fwk import runner, loader, helper, model

import pytest
import boto3
import logging

from e2e import service_marker, scenarios_directory, resource_directory, CRD_VERSION, CRD_GROUP, SERVICE_NAME
from e2e.bootstrap_resources import get_bootstrap_resources
from e2e.replacement_values import REPLACEMENT_VALUES

from acktest.k8s import resource as k8s

@helper.register_resource_helper(resource_kind="ParameterGroup", resource_plural="ParameterGroups")
class ParameterGroupHelper(helper.ResourceHelper):
    """
    Helper for parameter group scenarios.
    Overrides methods as required for custom resources.
    """

    def assert_extra_items_k8s(self, expected_name, expected_value, k8s_resource) -> bool:
        # parameterNameValues are processed by custom logic
        # returns AssertError if value of any parameter doesn't match the value in k8s resource
        if expected_name == "parameterNameValues":
            if not expected_value:
                assert expected_value == k8s_resource.get(expected_name)
                return True

            latest_parameters = k8s_resource.get(expected_name)
            for expected_parameter in expected_value:
                parameter_exist = False
                for latest_parameter in latest_parameters:
                    if expected_parameter == latest_parameter:
                        parameter_exist = True
                        break
                if parameter_exist:
                    continue
                assert expected_parameter == None
            return True
        return False

    def get_aws_resource(self, resource_name):
        parameter_groups = self.mdb.describe_parameter_groups(ParameterGroupName=resource_name)
        parameter_group = parameter_groups['ParameterGroups'][0]
        return parameter_group

    def assert_extra_items_aws(self, expected_key, expected_value, aws_resource) -> bool:
        # Parameters are processed by custom logic
        # returns AssertError if value of any parameter doesn't match the value in aws resource
        if expected_key == "Parameters":
            parameter_group_name = aws_resource.get("Name")
            latest_parameters = self.mdb.describe_parameters(ParameterGroupName=parameter_group_name)['Parameters']

            for expected_parameter in expected_value:
                parameter_exist = False
                parameter_name = expected_parameter.get("Name")
                parameter_value = expected_parameter.get("Value")
                for latest_parameter in latest_parameters:
                    if ("Name", parameter_name) in latest_parameter.items():
                        assert parameter_value == latest_parameter.get("Value")
                        parameter_exist = True
                        break
                if parameter_exist:
                    continue
                assert expected_parameter == None
            return True
        return False

    def wait_for_delete(self, reference: k8s.CustomResourceReference):
        logging.debug(f"ParameterGroup - wait_for_delete()")

@helper.register_resource_helper(resource_kind="Snapshot", resource_plural="Snapshots")
class SnapshotHelper(helper.ResourceHelper):
    """
    Helper for snapshot scenarios.
    Overrides methods as required for custom resources.
    """

    def get_aws_resource(self, resource_name):
        snapshots = self.mdb.describe_snapshots(SnapshotName=resource_name)
        snapshot = snapshots['Snapshots'][0]
        return snapshot

    def assert_extra_items_aws(self, expected_key, expected_value, aws_resource) -> bool:
        # cluster name is processed by custom logic
        # returns AssertError if expected cluster name doesn't match cluster name in aws resource
        if expected_key == "ClusterConfiguration":
            if expected_value.items() <= aws_resource.get(expected_key).items():
                return True
            else:
                assert value == None
                return True
        return False

    def wait_for_delete(self, reference: k8s.CustomResourceReference):
        logging.debug(f"Snapshot - wait_for_delete()")

@helper.register_resource_helper(resource_kind="User", resource_plural="Users")
class UserHelper(helper.ResourceHelper):
    """
    Helper for user scenarios.
    Overrides methods as required for custom resources.
    """

    def assert_extra_items_k8s(self, expected_name, expected_value, k8s_resource) -> bool:
        # events are processed by custom logic
        # returns AssertError if events are empty in k8s resource
        if expected_name == "events":
            assert k8s_resource.get(expected_name) is not None
            return True
        return False

    def get_aws_resource(self, resource_name):
        users = self.mdb.describe_users(UserName=resource_name)
        user = users['Users'][0]
        return user

    def wait_for_delete(self, reference: k8s.CustomResourceReference):
        logging.debug(f"User - wait_for_delete()")


@helper.register_resource_helper(resource_kind="SubnetGroup", resource_plural="SubnetGroups")
class SubnetGroupHelper(helper.ResourceHelper):
    """
    Helper for subnet group scenarios.
    Overrides methods as required for custom resources.
    """

    def get_aws_resource(self, resource_name):
        subnet_groups = self.mdb.describe_subnet_groups(SubnetGroupName=resource_name)
        subnet_group = subnet_groups['SubnetGroups'][0]
        return subnet_group


    def assert_extra_items_aws(self, expected_key, expected_value, aws_resource) -> bool:
        # Subnets are processed by custom logic
        # returns AssertError if any expected subnet is not in subnets field of aws_resource
        if expected_key == "Subnets":
            latest_subnets = aws_resource.get(expected_key)
            for expected_subnet in expected_value:
                subnet_exist = False
                for latest_subnet in latest_subnets:
                    if expected_subnet.items() <= latest_subnet.items():
                        subnet_exist = True
                        break
                if subnet_exist:
                    continue
                assert expected_subnet == None
            return True
        return False

    def wait_for_delete(self, reference: k8s.CustomResourceReference):
        logging.debug(f"SubnetGroup - wait_for_delete()")


@helper.register_resource_helper(resource_kind="ACL", resource_plural="ACLs")
class ACLHelper(helper.ResourceHelper):
    """
    Helper for ACL scenarios.
    Overrides methods as required for custom resources.
    """

    def assert_extra_items_k8s(self, expected_name, expected_value, k8s_resource) -> bool:
        # events are processed by custom logic
        # returns AssertError if events are empty in k8s resource
        if expected_name == "events":
            assert k8s_resource.get(expected_name) is not None
            return True
        return False

    def get_aws_resource(self, resource_name):
        acls = self.mdb.describe_acls(ACLName=resource_name)
        acl = acls['ACLs'][0]
        return acl

    def wait_for_delete(self, reference: k8s.CustomResourceReference):
        logging.debug(f"ACL - wait_for_delete()")


@helper.register_resource_helper(resource_kind="Cluster", resource_plural="Clusters")
class ClusterHelper(helper.ResourceHelper):
    """
    Helper for Cluster scenarios.
    Overrides methods as required for custom resources.
    """

    def assert_extra_items_k8s(self, expected_name, expected_value, k8s_resource) -> bool:
        # events are processed by custom logic
        # returns AssertError if events are empty in k8s resource
        if expected_name == "events":
            assert k8s_resource.get(expected_name) is not None
            return True
        return False

    def get_aws_resource(self, resource_name):
        clusters = self.mdb.describe_clusters(ClusterName=resource_name)
        cluster = clusters['Clusters'][0]
        return cluster

    def wait_for_delete(self, reference: k8s.CustomResourceReference):
        logging.debug(f"Cluster - wait_for_delete()")


@pytest.fixture(scope="session")
def input_replacements():
    """
    provides input replacements for test scenarios.
    """
    resource_replacements = REPLACEMENT_VALUES
    replacements = {
        "CRD_VERSION": CRD_VERSION,
        "CRD_GROUP": CRD_GROUP,
        "SERVICE_NAME": SERVICE_NAME
    }
    yield {**resource_replacements, **replacements}


@pytest.fixture(params=loader.list_scenarios(scenarios_directory), ids=loader.idfn)
def scenario(request, input_replacements):
    """
    Parameterized pytest fixture
    Provides test scenarios to execute
    Supports parallel execution of test scenarios
    """
    scenario_file_path = request.param
    scenario = loader.load_scenario(scenario_file_path, resource_directory, input_replacements)
    yield scenario
    runner.teardown(scenario)


@service_marker
class TestScenarios:
    """
    Declarative scenarios based test suite
    """
    def test_scenario(self, scenario):
        runner.run(scenario)