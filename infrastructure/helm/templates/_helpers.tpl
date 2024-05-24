{{/*
Expand the name of the chart.
*/}}
{{- define "clamav-scanner.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "clamav-scanner.fullname" -}}
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
{{- define "clamav-scanner.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "clamav-scanner.labels" -}}
app: {{ include "clamav-scanner.name" . }}
component: service
{{ include "clamav-scanner.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
helm.sh/chart: {{ include "clamav-scanner.chart" . }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Extra labels
*/}}
{{- define "clamav-scanner.extraLabels" -}}
environment: {{  .Values.target_env }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "clamav-scanner.selectorLabels" -}}
app.kubernetes.io/name: {{ include "clamav-scanner.name" . }}
app.kubernetes.io/instance: {{ include "clamav-scanner.name" . }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "clamav-scanner.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "clamav-scanner.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}
