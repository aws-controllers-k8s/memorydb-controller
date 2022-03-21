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
from acktest import resources
from acktest.resources import random_suffix_name
from acktest.k8s import resource as k8s


@dataclass
class Secret(Bootstrappable):
    # Inputs
    name: str

    def bootstrap(self):
        """Create a new secret.
        """
        super().bootstrap()
        k8s.create_opaque_secret("default", self.name, "password", random_suffix_name("password", 32))

    def cleanup(self):
        """Delete the secret.
        """
        super().cleanup()
        k8s.delete_secret("default", self.name)