{{/* The name of the application this chart installs */}}
{{- define "ack-memorydb-controller.app.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "ack-memorydb-controller.app.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/* The name and version as used by the chart label */}}
{{- define "ack-memorydb-controller.chart.name-version" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/* The name of the service account to use */}}
{{- define "ack-memorydb-controller.service-account.name" -}}
    {{ default "default" .Values.serviceAccount.name }}
{{- end -}}

{{- define "ack-memorydb-controller.watch-namespace" -}}
{{- if eq .Values.installScope "namespace" -}}
{{ .Values.watchNamespace | default .Release.Namespace }}
{{- end -}}
{{- end -}}

{{/* The mount path for the shared credentials file */}}
{{- define "ack-memorydb-controller.aws.credentials.secret_mount_path" -}}
{{- "/var/run/secrets/aws" -}}
{{- end -}}

{{/* The path the shared credentials file is mounted */}}
{{- define "ack-memorydb-controller.aws.credentials.path" -}}
{{ $secret_mount_path := include "ack-memorydb-controller.aws.credentials.secret_mount_path" . }}
{{- printf "%s/%s" $secret_mount_path .Values.aws.credentials.secretKey -}}
{{- end -}}

{{/* The rules a of ClusterRole or Role */}}
{{- define "ack-memorydb-controller.rbac-rules" -}}
rules:
- apiGroups:
  - ""
  resources:
  - configmaps
  verbs:
  - get
  - list
  - patch
  - watch
- apiGroups:
  - ""
  resources:
  - namespaces
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - ""
  resources:
  - secrets
  verbs:
  - get
  - list
  - patch
  - watch
- apiGroups:
  - ec2.services.k8s.aws
  resources:
  - securitygroups
  verbs:
  - get
  - list
- apiGroups:
  - ec2.services.k8s.aws
  resources:
  - securitygroups/status
  verbs:
  - get
  - list
- apiGroups:
  - ec2.services.k8s.aws
  resources:
  - subnets
  verbs:
  - get
  - list
- apiGroups:
  - ec2.services.k8s.aws
  resources:
  - subnets/status
  verbs:
  - get
  - list
- apiGroups:
  - kms.services.k8s.aws
  resources:
  - keys
  verbs:
  - get
  - list
- apiGroups:
  - kms.services.k8s.aws
  resources:
  - keys/status
  verbs:
  - get
  - list
- apiGroups:
  - memorydb.services.k8s.aws
  resources:
  - acls
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
- apiGroups:
  - memorydb.services.k8s.aws
  resources:
  - acls/status
  verbs:
  - get
  - patch
  - update
- apiGroups:
  - memorydb.services.k8s.aws
  resources:
  - clusters
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
- apiGroups:
  - memorydb.services.k8s.aws
  resources:
  - clusters/status
  verbs:
  - get
  - patch
  - update
- apiGroups:
  - memorydb.services.k8s.aws
  resources:
  - parametergroups
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
- apiGroups:
  - memorydb.services.k8s.aws
  resources:
  - parametergroups/status
  verbs:
  - get
  - patch
  - update
- apiGroups:
  - memorydb.services.k8s.aws
  resources:
  - snapshots
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
- apiGroups:
  - memorydb.services.k8s.aws
  resources:
  - snapshots/status
  verbs:
  - get
  - patch
  - update
- apiGroups:
  - memorydb.services.k8s.aws
  resources:
  - subnetgroups
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
- apiGroups:
  - memorydb.services.k8s.aws
  resources:
  - subnetgroups/status
  verbs:
  - get
  - patch
  - update
- apiGroups:
  - memorydb.services.k8s.aws
  resources:
  - users
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
- apiGroups:
  - memorydb.services.k8s.aws
  resources:
  - users/status
  verbs:
  - get
  - patch
  - update
- apiGroups:
  - services.k8s.aws
  resources:
  - adoptedresources
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
- apiGroups:
  - services.k8s.aws
  resources:
  - adoptedresources/status
  verbs:
  - get
  - patch
  - update
- apiGroups:
  - services.k8s.aws
  resources:
  - fieldexports
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
- apiGroups:
  - services.k8s.aws
  resources:
  - fieldexports/status
  verbs:
  - get
  - patch
  - update
- apiGroups:
  - sns.services.k8s.aws
  resources:
  - topics
  verbs:
  - get
  - list
- apiGroups:
  - sns.services.k8s.aws
  resources:
  - topics/status
  verbs:
  - get
  - list
{{- end }}