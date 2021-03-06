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
"""Stores the values used by each of the integration tests for replacing the
MemoryDB-specific test variables.
"""
from e2e.bootstrap_resources import get_bootstrap_resources

REPLACEMENT_VALUES = {
    "SECRET1": get_bootstrap_resources().Secret1.name,
    "SECRET2": get_bootstrap_resources().Secret2.name,
    "SUBNET1": get_bootstrap_resources().Subnets.subnets[0],
    "SUBNET2": get_bootstrap_resources().Subnets.subnets[1],
    "TOPIC1": get_bootstrap_resources().Topic1.topic_arn,
    "TOPIC2": get_bootstrap_resources().Topic2.topic_arn,
    "KMSKEY": get_bootstrap_resources().KMSKey.key,
    "SNAPSHOT_CLUSTER_NAME1": get_bootstrap_resources().Cluster1.clusterName,
    "SNAPSHOT_CLUSTER_NAME2": get_bootstrap_resources().Cluster2.clusterName
}
