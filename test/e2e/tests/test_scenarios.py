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

from e2e.declarative_test_fwk import runner, loader, helper

import pytest
import boto3
import logging

from e2e import service_marker, scenarios_directory, resource_directory, CRD_VERSION, CRD_GROUP, SERVICE_NAME

from acktest.k8s import resource as k8s

from e2e.replacement_values import REPLACEMENT_VALUES


@helper.register_resource_helper(resource_kind="ParameterGroup", resource_plural="ParameterGroups")
class ParameterGroupHelper(helper.ResourceHelper):
    """
    Helper for parameter group scenarios.
    Overrides methods as required for custom resources.
    """

    def wait_for_delete(self, reference: k8s.CustomResourceReference):
        logging.debug(f"ParameterGroup - wait_for_delete()")


@helper.register_resource_helper(resource_kind="User", resource_plural="Users")
class UserHelper(helper.ResourceHelper):
    """
    Helper for user scenarios.
    Overrides methods as required for custom resources.
    """

    def wait_for_delete(self, reference: k8s.CustomResourceReference):
        logging.debug(f"User - wait_for_delete()")


@helper.register_resource_helper(resource_kind="SubnetGroup", resource_plural="SubnetGroups")
class SubnetGroupHelper(helper.ResourceHelper):
    """
    Helper for subnet group scenarios.
    Overrides methods as required for custom resources.
    """

    def wait_for_delete(self, reference: k8s.CustomResourceReference):
        logging.debug(f"SubnetGroup - wait_for_delete()")


@helper.register_resource_helper(resource_kind="ACL", resource_plural="ACLs")
class ACLHelper(helper.ResourceHelper):
    """
    Helper for ACL scenarios.
    Overrides methods as required for custom resources.
    """

    def wait_for_delete(self, reference: k8s.CustomResourceReference):
        logging.debug(f"ACL - wait_for_delete()")


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
