ffiliates. All Rights Reserved.
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

"""Utilities for working with parameter group resources"""

import datetime
import time
import typing

import boto3
import pytest


def get(parameter_group_name):
    """Returns a dict containing the parameter group record from the MemoryDB
    API.
    If no such parameter group exists, returns None.
    """
    c = boto3.client('memorydb')
    try:
        resp = c.describe_parameter_groups(
            ParameterGroupName=parameter_group_name,
        )
        assert len(resp['ParameterGroups']) == 1
        return resp['ParameterGroups'][0]
    except c.exceptions.ParameterGroupNotFoundFault:
        return None


def get_tags(parameter_group_arn):
    """Returns a dict containing the parameter group's tag records from the
    MemoryDB API.
    If no such parameter group exists, returns None.
    """
    c = boto3.client('memorydb')
    try:
        resp = c.list_tags_for_resource(
            ResourceName=parameter_group_arn,
        )
        return resp['TagList']
    except c.exceptions.ParameterGroupNotFoundFault:
        return None