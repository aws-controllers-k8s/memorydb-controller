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

from acktest.bootstrapping import Bootstrappable
import boto3


@dataclass
class KMS(Bootstrappable):
    # Output
    key: str = field(init=False)

    def bootstrap(self):
        """Create a KMS key.
        """
        super().bootstrap()
        kms = boto3.client("kms")

        response = kms.create_key(Description="Key for ACK MemoryDB tests")
        self.key = response['KeyMetadata']['KeyId']

    def cleanup(self):
        """Delete KMS key.
        KMS does not allow immediate key deletion; 7 days is the shortest deletion window
        """
        super().cleanup()
        kms = boto3.client("kms")
        kms.schedule_key_deletion(KeyId=self.key, PendingWindowInDays=7)
