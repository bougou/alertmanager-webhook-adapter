{{/*
Expand the name of the chart.
*/}}
{{- define "awa.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "awa.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "awa.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "awa.labels" -}}
helm.sh/chart: {{ include "awa.chart" . }}
{{ include "awa.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- with .Values.commonLabels }}
{{ toYaml . }}
{{- end }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "awa.selectorLabels" -}}
app.kubernetes.io/name: {{ include "awa.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the namespace to deploy into.
*/}}
{{- define "awa.namespace" -}}
{{- if .Values.namespaceOverride }}
{{- .Values.namespaceOverride }}
{{- else }}
{{- .Release.Namespace }}
{{- end }}
{{- end }}

{{/*
Return the proper image repository.
*/}}
{{- define "awa.image.repository" -}}
{{- .Values.image.repository | default .Values.image.name | default "bougou/alertmanager-webhook-adapter" }}
{{- end }}

{{/*
Return the proper image tag.
*/}}
{{- define "awa.image.tag" -}}
{{- .Values.image.tag | default .Chart.AppVersion | default "latest" }}
{{- end }}

{{/*
Return the full image reference.
*/}}
{{- define "awa.image" -}}
{{- printf "%s:%s" (include "awa.image.repository" .) (include "awa.image.tag" .) }}
{{- end }}

{{/*
Return the service account name.
*/}}
{{- define "awa.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "awa.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Return application signature (backward compatible).
*/}}
{{- define "awa.config.signature" -}}
{{- .Values.config.signature | default .Values.signature | default "未知" }}
{{- end }}

{{/*
Return template language (backward compatible).
*/}}
{{- define "awa.config.lang" -}}
{{- .Values.config.lang | default .Values.lang }}
{{- end }}

{{/*
Return timezone (backward compatible).
*/}}
{{- define "awa.config.timezone" -}}
{{- .Values.config.timezone | default .Values.timezone | default "UTC" }}
{{- end }}

{{/*
Return template directory path.
*/}}
{{- define "awa.config.tmplDir" -}}
{{- if include "awa.templates.mountEnabled" . }}
{{- .Values.config.tmplDir | default .Values.templates.mountPath }}
{{- else }}
{{- .Values.config.tmplDir }}
{{- end }}
{{- end }}

{{/*
Return the ConfigMap name for custom templates.
*/}}
{{- define "awa.templates.configMapName" -}}
{{- if .Values.templates.existingConfigMap }}
{{- .Values.templates.existingConfigMap }}
{{- else }}
{{- include "awa.fullname" . }}-templates
{{- end }}
{{- end }}

{{/*
Return whether custom templates should be mounted.
*/}}
{{- define "awa.templates.mountEnabled" -}}
{{- if .Values.templates.enabled -}}
{{- if or .Values.templates.existingConfigMap (gt (len .Values.templates.data) 0) -}}
true
{{- end -}}
{{- end -}}
{{- end }}

{{/*
Return container args.
*/}}
{{- define "awa.containerArgs" -}}
- --listen-address={{ .Values.config.listenAddress | default ":8090" }}
- --signature={{ include "awa.config.signature" . }}
{{- with include "awa.config.lang" . }}
{{- if ne . "" }}
- --tmpl-lang={{ . }}
{{- end }}
{{- end }}
{{- with include "awa.config.tmplDir" . }}
{{- if ne . "" }}
- --tmpl-dir={{ . }}
{{- end }}
{{- end }}
{{- with .Values.config.tmplName }}
- --tmpl-name={{ . }}
{{- end }}
{{- with .Values.config.tmplDefault }}
- --tmpl-default={{ . }}
{{- end }}
{{- if .Values.config.debug }}
- --debug
{{- end }}
{{- with .Values.extraArgs }}
{{ toYaml . }}
{{- end }}
{{- end }}
