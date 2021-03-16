{% set indent1 = indentLevel | indent1 -%}

{% macro dto(schemaName, schema) %}
type {{schemaName | upperFirst}} struct {
{% for propName, prop in schema.properties() -%}
{%- set required = schema.required() -%}
{%- set typeInfo = [propName, required, prop] | fixType -%}
{{indent1}}{{propName | fixPropName}} {{typeInfo}} `json:"{{propName}},omitempty"`
{% endfor -%}
}
{% endmacro %}
