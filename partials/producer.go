{% from "./headers.go" import headers %}

{% macro producer(asyncapi) %}
// producer
type Producer struct {
	msgProducer producer.MessageProducer
}

func NewProducer(configStore core.ConfigStore, logger log.SimpleLogger, typeHelper common.AsyncTypeHelper, tracer tracer.SimpleTracer) (*Producer, error) {
	msgProducer, err := producer.NewMessageProducer(configStore, logger, typeHelper, tracer)
	if err != nil {
		return nil, err
	}
	return &Producer{msgProducer: msgProducer}, nil
}

{% for channelName, channel in asyncapi.channels() -%}
{%- if channel.hasPublish() -%}
{%- set publishOperation = channel.publish() -%}
{%- set channelInfo = [channel, publishOperation] | channelInfo -%}

// Headers object
{%- if channelInfo.hasHeaders %}
{{- headers(channelInfo.headerStructName, channelInfo.headers) -}}
{% endif %}

// {{publishOperation.summary()}}
func (p *Producer){{' '}}{{- channelInfo.funcName -}}(ctx core.Context,
	messageKey []byte,
	payload *{{channelInfo.payloadStructName}}
{%- if channelInfo.hasHeaders -%}
,
	headers {{channelInfo.headerStructName}}
{%- endif -%}
) (partition int32, offset int64, err error) {

	message := p.msgProducer.BuildMessage().
		AddByteKey(messageKey).
		{%- if channelInfo.hasHeaders %}
		{% for propName, prop in channelInfo.headers.properties() -%}
			{%- set typeInfo = [propName, required, prop] | fixType -%}
			AddStringHeader("{{propName}}", headers.{{propName | fixPropName}}).
		{% endfor -%}
		{%- endif -%}
		AddBody(&payload).
		Build()
	return p.msgProducer.SendMessage(ctx, "{{- channelName -}}", message)
}
{% endif %}
{% endfor -%}

{% endmacro %}
