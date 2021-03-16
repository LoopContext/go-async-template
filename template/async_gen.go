package {{params.package}}

{%- set asyncAttrs = asyncapi | checkAttrs %}

{% from "../partials/producer.go" import producer -%}
{% from "../partials/consumer.go" import consumer -%}
{% from "../partials/dto.go" import dto -%}

import (
	"bitbucket.org/integrapartners/prism-async-go/common"
	{% if asyncAttrs.hasPublish -%}"bitbucket.org/integrapartners/prism-async-go/producer"{%- endif %}
	{% if asyncAttrs.hasSubscribe -%}"bitbucket.org/integrapartners/prism-async-go/consumer"{%- endif %}
	{% if asyncAttrs.hasSubscribe -%}"bitbucket.org/integrapartners/prism-core-go/session"{%- endif %}
	"bitbucket.org/integrapartners/prism-core-go/core"
	"bitbucket.org/integrapartners/prism-core-go/log"
	"bitbucket.org/integrapartners/prism-core-go/tracer"
	{% if asyncAttrs.importTime -%}"time"{%- endif %}
)

// load schemas
{%- for schemaName, schema in asyncapi.components().schemas() -%}
{{- dto(schemaName, schema) -}}
{% endfor -%}

// load producers
{%- if asyncAttrs.hasPublish -%}
	{{- producer(asyncapi) -}}
{%- endif %}

// load consumer
{%- if asyncAttrs.hasSubscribe -%}
	{{- consumer(asyncapi) -}}
{%- endif %}
