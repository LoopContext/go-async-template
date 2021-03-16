{% set indent1 = indentLevel | indent1 -%}

{% macro headers(headersObjName, messageHeaders) %}
type {{headersObjName}} struct {
{% for propName, prop in messageHeaders.properties() -%}
{%- set required = messageHeaders.required() -%}
{%- set typeInfo = [propName, required, prop] | fixType -%}
{{indent1}}{{propName | fixPropName}} string
{% endfor -%}
}
{% endmacro %}
