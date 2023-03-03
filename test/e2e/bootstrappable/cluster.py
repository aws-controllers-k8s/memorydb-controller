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

from dataclasses import dataclass

import boto3
import time
from acktest.bootstrapping import Bootstrappable


@dataclass
class Cluster(Bootstrappable):
    # Output
    clusterName : str

    def create_cluster(self):
        mdb = boto3.client("memorydb")
        mdb.create_cluster(ClusterName=self.clusterName,
                           Description='Cluster for Ack snapshot resource testing', SnapshotRetentionLimit=0,
                           NodeType='db.r6g.large', ACLName='open-access', NumShards=1, NumReplicasPerShard=0)
        timeout = time.time() + 30*60  # 30 minutes from now
        available_status = "Available"
        while True:
            clusters = mdb.describe_clusters(ClusterName=self.clusterName)
            cluster = clusters['Clusters'][0]
            if cluster.get("Status").casefold() == available_status.casefold():
                return True
            if time.time() > timeout:
                break
            time.sleep(60)
        raise ValueError('cluster not created within expected time')

    def delete_cluster(self):
        mdb = boto3.client("memorydb")
        mdb.delete_cluster(ClusterName=self.clusterName)

    def bootstrap(self):
        """Find supported subnets.
        """
        super().bootstrap()
        # Try to create the subnet using all the available subnets
        self.create_cluster()

    def cleanup(self):
        """Nothing to do here as we did not create new subnets.
        """
        super().cleanup()

        self.delete_cluster()