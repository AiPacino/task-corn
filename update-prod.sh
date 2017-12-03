#!/bin/bash
kubectl --kubeconfig $KUBECONFIG_PROD get deployment/cailianpress-cron -o wide -n default
kubectl --kubeconfig $KUBECONFIG_PROD -n default set image deployment/cailianpress-cron cailianpress-cron=registry.cn-beijing.aliyuncs.com/lanjing/cailianpress-cron:$1
kubectl --kubeconfig $KUBECONFIG_PROD get deployment/cailianpress-cron -o wide -n default
