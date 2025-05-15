# Prometheus AlertManager Slack Integration Guide

This guide explains how to set up Prometheus AlertManager to send alerts to Slack.

## Prerequisites

- A Kubernetes cluster with Prometheus and AlertManager installed;
- A Slack workspace with permissions to create apps;
- A Slack bot token starting with `xoxb-`;

## Setup Steps for Slack Integration

1. Create a [Slack App](https://api.slack.com/apps);

2. Configure the App:
- Add a [bot token](https://api.slack.com/concepts/token-types#bot);
- Add a slack channel;
- Add a slack user;  

3. Install the App to your Slack workspace:
- Click on "Install App to Slack";
- Select your workspace;
- Click on "Install";

4. Add the user to the channel:
- Click on "Add User to Channel"
- Select the channel
- Click on "Add"

### 1. Create Slack token secret

Create a Kubernetes secret with the Slack API URL and token in the `monitoring` namespace:

```bash
TOKEN=xoxb-etcetera-etcetera
kubectl create secret generic slack-token \
  --from-literal=url="https://slack.com/api/chat.postMessage" \
  --from-literal=token="$TOKEN" \
  -n monitoring
```

### 2. AlertManager Configuration Namespace

Ensure the `AlertManagerConfig` resource is in the **same namespace** (`monitoring`) where AlertManager is deployed.  
If it is in another namespace (e.g., `bitcoin-investment-tracker`), AlertManager **will not detect** the configuration.

Verify the `alertmanagerConfigSelector`:

```bash
kubectl -n monitoring get alertmanager my-kube-prometheus-stack-alertmanager -o yaml | grep alertmanagerConfigSelector -A 3
```

Example output:

```yaml
alertmanagerConfigSelector:
  matchLabels:
    release: my-kube-prometheus-stack
automountServiceAccountToken: true
```

## 3. Matching Alert Labels

Your AlertManager routing configuration was expecting alerts with the label `namespace="monitoring"`:

```yaml
matchers:
- alertname="TestBitcoinPriceDrop"
- namespace="monitoring"
```

To ensure alerts are routed correctly, we added the corresponding label in the Prometheus alert rule:

```yaml
labels:
  severity: warning
  namespace: monitoring  # This label ensures the alert matches AlertManager routing rules
```

---

## ✅ Summary

- `AlertManagerConfig` must be in the same namespace as AlertManager (`monitoring`).
- Labels in alert rules must match those defined in AlertManager routing config.
- Slack API credentials should be stored in a Kubernetes secret (`slack-token`).
- Use the `chat.postMessage` API with your Slack bot token.
