{{/*
Expand the name of the chart.
*/}}
{{- define "kube-pod-rightsizer.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "kube-pod-rightsizer.fullname" -}}
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
{{- define "kube-pod-rightsizer.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "kube-pod-rightsizer.labels" -}}
helm.sh/chart: {{ include "kube-pod-rightsizer.chart" . }}
{{ include "kube-pod-rightsizer.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "kube-pod-rightsizer.selectorLabels" -}}
app.kubernetes.io/name: {{ include "kube-pod-rightsizer.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "kube-pod-rightsizer.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "kube-pod-rightsizer.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Create the namespace name
*/}}
{{- define "kube-pod-rightsizer.namespace" -}}
{{- default .Release.Namespace .Values.namespaceOverride }}
{{- end }}
