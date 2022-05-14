---
title: NetworkPolicy升级
layout: post
category: k8s
author: 夏泽民
---
The v1.16 release will stop serving the following deprecated API versions in favor of newer and more stable API versions:

NetworkPolicy in the extensions/v1beta1 API version is no longer served
Migrate to use the networking.k8s.io/v1 API version, available since v1.8. Existing persisted data can be retrieved/updated via the new version.
PodSecurityPolicy in the extensions/v1beta1 API version
Migrate to use the policy/v1beta1 API, available since v1.10. Existing persisted data can be retrieved/updated via the new version.
DaemonSet in the extensions/v1beta1 and apps/v1beta2 API versions is no longer served
Migrate to use the apps/v1 API version, available since v1.9. Existing persisted data can be retrieved/updated via the new version.
Notable changes:
spec.templateGeneration is removed
spec.selector is now required and immutable after creation; use the existing template labels as the selector for seamless upgrades
spec.updateStrategy.type now defaults to RollingUpdate (the default in extensions/v1beta1 was OnDelete)
Deployment in the extensions/v1beta1, apps/v1beta1, and apps/v1beta2 API versions is no longer served
Migrate to use the apps/v1 API version, available since v1.9. Existing persisted data can be retrieved/updated via the new version.
Notable changes:
spec.rollbackTo is removed
spec.selector is now required and immutable after creation; use the existing template labels as the selector for seamless upgrades
spec.progressDeadlineSeconds now defaults to 600 seconds (the default in extensions/v1beta1 was no deadline)
spec.revisionHistoryLimit now defaults to 10 (the default in apps/v1beta1 was 2, the default in extensions/v1beta1 was to retain all)
maxSurge and maxUnavailable now default to 25% (the default in extensions/v1beta1 was 1)
StatefulSet in the apps/v1beta1 and apps/v1beta2 API versions is no longer served
Migrate to use the apps/v1 API version, available since v1.9. Existing persisted data can be retrieved/updated via the new version.
Notable changes:
spec.selector is now required and immutable after creation; use the existing template labels as the selector for seamless upgrades
spec.updateStrategy.type now defaults to RollingUpdate (the default in apps/v1beta1 was OnDelete)
ReplicaSet in the extensions/v1beta1, apps/v1beta1, and apps/v1beta2 API versions is no longer served
Migrate to use the apps/v1 API version, available since v1.9. Existing persisted data can be retrieved/updated via the new version.
Notable changes:
spec.selector is now required and immutable after creation; use the existing template labels as the selector for seamless upgrades
The v1.22 release will stop serving the following deprecated API versions in favor of newer and more stable API versions:

Ingress in the extensions/v1beta1 API version will no longer be served
Migrate to use the networking.k8s.io/v1beta1 API version, available since v1.14. Existing persisted data can be retrieved/updated via the new version.
https://kubernetes.io/blog/2019/07/18/api-deprecations-in-1-16/
<!-- more -->

W0205 15:14:07.482439       1 warnings.go:67] extensions/v1beta1 Ingress is deprecated in v1.14+, unavailable in v1.22+; use networking.k8s.io/v1 Ingress
time="2021-02-05T15:14:07Z" level=info msg="Updated ingress status" namespace=default ingress=cheddar
W0205 15:18:19.104225       1 warnings.go:67] networking.k8s.io/v1beta1 IngressClass is deprecated in v1.19+, unavailable in v1.22+; use networking.k8s.io/v1 IngressClassList

https://stackoverflow.com/questions/66080909/logs-complaining-extensions-v1beta1-ingress-is-deprecated

https://github.com/kubernetes/kubernetes/issues/94761
