"""Handlers for delegating AWS Cloudformation hooks to OPA"""

import logging
from typing import Any, MutableMapping, Optional

from botocore.exceptions import ClientError
import requests

from cloudformation_cli_python_lib import (
    Hook,
    HookInvocationPoint,
    OperationStatus,
    ProgressEvent,
    SessionProxy,
)

from .models import HookHandlerRequest, TypeConfigurationModel

# Use this logger to forward log messages to CloudWatch Logs.
LOG = logging.getLogger(__name__)
LOG.setLevel(logging.INFO)

TYPE_NAME = "OpenPolicyAgent::OPA::Hook"

hook = Hook(TYPE_NAME, TypeConfigurationModel)
test_entrypoint = hook.test_entrypoint

def get_secret(name: str, session: Optional[SessionProxy]) -> str:
    """Get the (optional) secret to use as bearer token for authenticating against OPA"""

    client = session.client("secretsmanager")

    try:
        resp = client.get_secret_value(SecretId=name)
    # pylint: disable=invalid-name
    except ClientError as e:
        LOG.error("Failed fetching secret %s", name)
        LOG.error(e)
        raise e

    if 'SecretString' in resp:
        return resp['SecretString']

    raise Exception("SecretString not found in secret")

def opa_query(
        request: HookHandlerRequest,
        session: Optional[SessionProxy],
        type_configuration: TypeConfigurationModel,
        action: str,
) -> ProgressEvent:
    """Query OPA and return a ProgressEvent based on the decision"""

    progress: ProgressEvent = ProgressEvent(
        status=OperationStatus.IN_PROGRESS
    )

    # Querying the default decision, so don't wrap in "input" attribute
    opa_input = {
        "action": action,
        "hook": request.hookContext.hookTypeName,
        "resource": {
            "id": request.hookContext.targetLogicalId,
            "name": request.hookContext.targetName,
            "type": request.hookContext.targetType,
            "properties": request.hookContext.targetModel.get("resourceProperties")
        }
    }

    headers = {}
    secret = type_configuration.opaAuthTokenSecret
    if secret is not None and secret != "":
        token = get_secret(type_configuration.opaAuthTokenSecret, session)
        headers = {"Authorization": f"Bearer {token}"}

    try:
        connect_timeout = 10 
        resp = requests.post(type_configuration.opaUrl, json=opa_input, headers=headers, timeout=connect_timeout)
    except requests.exceptions.ConnectTimeout:
        LOG.error("Timeout connecting to OPA at %s in %s seconds", type_configuration.opaUrl, connect_timeout)
        progress.status = OperationStatus.FAILED
        return progress
    except requests.ConnectionError:
        LOG.error("Failed connecting to OPA at %s", type_configuration.opaUrl)
        progress.status = OperationStatus.FAILED
        return progress

    if resp.status_code == 200:
        body = resp.json()
        if not "allow" in body:
            LOG.error("OPA returned empty/undefined result")
            progress.status = OperationStatus.FAILED
        else:
            if body["allow"] is True:
                progress.status = OperationStatus.SUCCESS
            else:
                message = " | ".join(body["violations"])
                LOG.info("OPA denied the request with message: %s", message)
                progress.status = OperationStatus.FAILED
                progress.message = message

    else:
        LOG.error("OPA returned status code: %d", resp.status_code)
        LOG.error(resp.json())
        progress.status = OperationStatus.FAILED

    return progress

# pylint: disable=unused-argument,missing-function-docstring
@hook.handler(HookInvocationPoint.CREATE_PRE_PROVISION)
@hook.handler(HookInvocationPoint.UPDATE_PRE_PROVISION)
@hook.handler(HookInvocationPoint.DELETE_PRE_PROVISION)
def pre_handler(
        session: Optional[SessionProxy],
        request: HookHandlerRequest,
        callback_context: MutableMapping[str, Any],
        type_configuration: TypeConfigurationModel
) -> ProgressEvent:

    LOG.info("Hook triggered for target %s %s",
        request.hookContext.targetName,
        request.hookContext.targetLogicalId
    )

    action = request.hookContext.invocationPoint[0:6]

    return opa_query(request, session, type_configuration, action)
