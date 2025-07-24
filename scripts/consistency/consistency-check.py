import argparse
import requests
import yaml
from datetime import datetime, timezone


def load_config(path):
    """Load configuration options from a yaml file

    Args:
        path (string): path to the config.yml file

    Returns:
        dict: A dictionary with all configuration options
    """
    with open(path, "r", encoding="utf-8") as f:
        return yaml.safe_load(f) or {}


def publish_results(status: str, message: str, api: dict):
    """Publish check results to a remote argo-web-api endpoint

    Args:
        status (str): status value of the check - OK or CRITICAL
        message (str): message containing details about the check performed
        api (dict): a subset of the configuration dictionary containing all the api-related config options
    """
    timeout = api.get("timeout", 10)
    resp = requests.post(
        f"https://{api['endpoint']}/api/v3/consistency/auto-check",
        headers={"Accept": "application/json", "x-api-key": api["access_token"]},
        json={"status": status, "message": message},
        verify=api.get("verify", True),
        timeout=timeout,
    )
    resp.raise_for_status()


def get_num_of_endpoints(report: str, tenant: dict, api: dict) -> int:
    """Get the number of endpoints for a specific tenant and a specific report from a remote argo-web-api instance

    Args:
        report (str): name of the report
        tenant (dict): a subset of the configuration dictionary containing all config options for a specific tenant
        api (dict): a subset of the configuration dictionary containing all the api-related config options

    Returns:
        int: number of endpoints found
    """
    today = datetime.now(timezone.utc).date().isoformat()
    headers = {"Accept": "application/json", "x-api-key": tenant["access_token"]}
    timeout = api.get("timeout", 10)
    resp = requests.get(
        f"https://{api['endpoint']}/api/v2/results/{report}/endpoints?start_time={today}T00:00:00Z&end_time={today}T23:59:59Z",
        headers=headers,
        verify=api.get("verify", True),
        timeout=timeout,
    )
    resp.raise_for_status()
    return len(resp.json()["results"])


def get_num_of_flapping_endpoints(report: str, tenant: dict, api: dict) -> int:
    """Get the number of flapping endpoints for a specific tenant and a specific report from a remote argo-web-api instance

    Args:
        report (str): name of the report
        tenant (dict): a subset of the configuration dictionary containing all config options for a specific tenant
        api (dict): a subset of the configuration dictionary containing all the api-related config options

    Returns:
        int: number of flapping endpoints found
    """
    today = datetime.now(timezone.utc).date().isoformat()
    headers = {"Accept": "application/json", "x-api-key": tenant["access_token"]}
    timeout = api.get("timeout", 10)
    resp = requests.get(
        f"https://{api['endpoint']}/api/v2/trends/{report}/flapping/endpoints?date={today}",
        headers=headers,
        verify=api.get("verify", True),
        timeout=timeout,
    )
    resp.raise_for_status()
    return len(resp.json()["data"])


def check_tenant(threshold: int, tenant: dict, api: dict) -> tuple[bool, str]:
    """Check data for a specific tenant in an argo-web-api instance and compare flapping endpoints found against a specific threshold

    Args:
        threshold (int): threshold percentage to compare against
        tenant (dict): a subset of the configuration dictionary containing all config options for a specific tenant
        api (dict): a subset of the configuration dictionary containing all the api-related config options

    Returns:
        tuple[bool, str]: a tuple with a boolean that designates if check is ok and a message string containing details
    """
    ok = True
    messages = []
    for report in tenant["reports"]:
        endpoint_num = get_num_of_endpoints(report, tenant, api)
        flapping_num = get_num_of_flapping_endpoints(report, tenant, api)
        if flapping_num == 0:
            messages.append(f"{tenant['name']}-{report}: No flapping endpoints.\n")
        elif endpoint_num > 0:
            # find flapping percentage
            perc = flapping_num * 100 / endpoint_num
            comp = "<="
            if perc > threshold:
                comp = ">"
                ok = False
            messages.append(
                f"{tenant['name']}-{report}: Flapping percent {int(perc)}% {comp} 50% ({flapping_num} of {endpoint_num})."
            )
    return (ok, "\n".join(messages))


def check_instance(config: dict) -> tuple[bool, str]:
    """Check a remote instance of argo-web-api for consistency

    Args:
        config (dict): the configuration dictionary

    Returns:
        tuple[bool, str]: a tuple with a boolean that designates if check is ok and a message string containing details
    """
    ok = True
    messages = []
    # for each tenant
    for tenant_config in config["check_tenants"]:
        tenant_ok, tenant_message = check_tenant(
            config["check_threshold"], tenant_config, config["api"]
        )
        ok = ok and tenant_ok
        messages.append(tenant_message)
    return ok, "".join(messages)


def main(args):

    config = load_config(args.config)
    # Use configuration to check designated instance - return if everything is ok (boolean) and message
    ok, message = check_instance(config)

    status = "OK" if ok else "CRITICAL"

    publish_results(status, message, config["api"])

    return 0


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "-c", "--config", required=True, help="path to .yml configuration file"
    )
    args = parser.parse_args()
    exit(main(args))
