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
from acktest.resources import random_suffix_name
import boto3


@dataclass
class Topic(Bootstrappable):
    # Output
    topic_arn: str = field(init=False)

    def bootstrap(self):
        """Create a SNS topic.
        """
        super().bootstrap()
        topic_name1 = random_suffix_name("ack-sns-topic", 32)
        sns = boto3.client("sns")
        self.topic_arn = sns.create_topic(Name=topic_name1)['TopicArn']

    def cleanup(self):
        """Delete SNS topics.
        """
        super().cleanup()
        sns = boto3.client("sns")
        sns.delete_topic(TopicArn=self.topic_arn)