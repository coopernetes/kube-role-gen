#!/usr/bin/env python3
import yaml
import json
import sys
import subprocess
import typing
import logging
import os
from collections import defaultdict

logging.basicConfig(format='%(asctime)s - %(levelname)s - %(message)s', level=os.environ.get("LOGLEVEL", "INFO"))

def run_command(args: typing.List) -> str:
    p = subprocess.run(args, stdout=subprocess.PIPE)
    p.check_returncode()
    return p.stdout

def convert_to_rbac_role(resources: dict) -> dict:
    """
    Convert a dictionary of resource names & their supported verbs, obtained by
    the Kubernetes API to a properly formatted Kubernetes Role object. The list
    of rules will contain every resource supported by the connected API server.
    Each rule in the Role is a unique by API group and supported verbs for any
    given resource.
    """
    logging.info("Converting resources to rbac.authorization.k8s.io/v1/Role")
    rbac_rules = _create_role_yaml()
    unique_resources_by_verb = defaultdict(dict)
    for api_group, resource_dict in resources.items():
        # key = api_group + verbs (as string)
        # value = resource names
        for resource_name, verbs in resource_dict.items():
            # represent the list of verbs as a json string, so we can use it
            # as a unique key for resources w/ matching supported verbs
            new_key = f"{api_group}_{json.dumps(verbs)}"
            if not unique_resources_by_verb[new_key]:
                unique_resources_by_verb[new_key] = []
            unique_resources_by_verb[new_key].append(resource_name)

    for key, value in unique_resources_by_verb.items():
        logging.debug(key)
        api_group = key.split('_')[0]
        if api_group == 'core':
            api_group = ""
        verbs = json.loads(key.split('_')[1])
        rbac_rules['rules'].append({
            'apiGroups': [api_group],
            'resources': value,
            'verbs': verbs
        })
    return rbac_rules


def _create_role_yaml() -> dict:
    return {
        'apiVersion': 'rbac.authorization.k8s.io/v1',
        'kind': 'Role',
        'metadata': {'name': 'foo-role'},
        'rules': []
    }


def main():
    rbac_resources = {}
    kube_core_api = json.loads(run_command(["kubectl", "get", "--raw",
                                            "/api/v1"]))
    kube_api_list = json.loads(run_command(["kubectl", "get", "--raw",
                                            "/apis"]))
 
    logging.info("Gathering core API resource details")
    rbac_resources['core'] = {}
    for resource in kube_core_api['resources']:
        rbac_resources['core'][resource['name']] = resource['verbs']
 
    logging.info("Gathering API groups & resource details")
    resource_count = 0
    for api in kube_api_list['groups']:
        # legacy resources in extensions apigroup are being deprecated in
        # Kubernetes 1.16 & 1.20. Skip it to enforce use of non-deprecated apis
        if api['name'] == "extensions":
            continue
        rbac_resources[api['name']] = defaultdict(dict)
        for version in api['versions']:
            response = json.loads(run_command(["kubectl", "get", "--raw",
                                       f"/apis/{version['groupVersion']}"]))
            resources = response['resources']
            for resource in resources:
                resource_count += 1
                rbac_resources[api['name']][resource['name']] = resource['verbs']

    logging.info(f"Resource discovery complete. Found {resource_count} resources in {len(rbac_resources.keys())} API groups")
    print(yaml.safe_dump(convert_to_rbac_role(rbac_resources),
                         default_flow_style=False))
    return 0


if __name__ == '__main__':
    sys.exit(main())
