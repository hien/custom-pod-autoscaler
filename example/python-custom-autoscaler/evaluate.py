# Copyright 2019 The Custom Pod Autoscaler Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import json
import sys
import math

# JSON piped into this script example:
# {
#   "resource": "hello-kubernetes",
#   "run_type": "scaler",
#   "metrics": [
#     {
#       "resource": "hello-kubernetes",
#       "value": "5"
#     }
#   ]
# }


def main():
    # Parse metrics JSON into a dict
    metrics = json.loads(sys.stdin.read())
    evaluate(metrics)


def evaluate(metrics):
    try:
        value = int(metrics["metrics"][0]["value"])

        # Build JSON dict with target_replicas
        evaluation = {}
        evaluation["target_replicas"] = value * 2

        # Output JSON to stdout
        sys.stdout.write(json.dumps(evaluation))
    except ValueError as err:
        # If not an integer, output error
        sys.stderr.write(f"Invalid metric value: {err}")
        exit(1)


if __name__ == "__main__":
    main()
