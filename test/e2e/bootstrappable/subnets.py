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

from dataclasses import dataclass, field
from typing import List

import boto3
from acktest.bootstrapping import Bootstrappable
from acktest.resources import random_suffix_name


@dataclass
class Subnets(Bootstrappable):
    # Output
    subnets: List[str] = field(init=False, default_factory=lambda: [])

    def bootstrap(self):
        """Find supported subnets.
        """
        super().bootstrap()
        ec2 = boto3.client("ec2")
        # Find the default VPC
        vpc_response = ec2.describe_vpcs(Filters=[{"Name": "isDefault", "Values": ["true"]}])
        default_vpc_id = vpc_response['Vpcs'][0]['VpcId']

        # Find all the subnets from default VPC
        subnet_response = ec2.describe_subnets(Filters=[{"Name": "vpc-id", "Values": [default_vpc_id]}])

        ec2_subnets = []
        for subnet in subnet_response['Subnets']:
            ec2_subnets.append(subnet['SubnetId'])

        # Try to create the subnet using all the available subnets
        mdb = boto3.client("memorydb")
        subnet_name = random_suffix_name("sub", 10)
        try:
            mdb.create_subnet_group(SubnetGroupName=subnet_name, Description='Determine valid subnets',
                                    SubnetIds=ec2_subnets)
        except mdb.exceptions.SubnetNotAllowedFault as ex:
            message = str(ex)
            exp_message = "Supported availability zones are "
            index = message.index(exp_message)
            # Format of the message is like "Supported availability zones are [us-east-1c, us-east-1d, us-east-1b]."
            valid_az = message[index + len(exp_message) + 1: len(message) - 2].split(", ")

            for subnet in subnet_response['Subnets']:
                if subnet['AvailabilityZone'] in valid_az:
                    self.subnets.append(subnet['SubnetId'])
            return

        # We were able to create subnet group using all the subnets, so delete the subnet group that was created
        self.subnets = ec2_subnets
        mdb.delete_subnet_group(SubnetGroupName=subnet_name)

    def cleanup(self):
        """Nothing to do here as we did not create new subnets.
        """
        super().cleanup()