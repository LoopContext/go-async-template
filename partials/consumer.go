{% from "./headers.go" import headers %}

{% macro consumer(asyncapi) %}
// consumer headers
{%- for channelName, channel in asyncapi.channels() -%}
{%- if channel.hasSubscribe() -%}
{%- set subscribeOperation = channel.subscribe()  -%}
{%- set channelInfo = [channel, subscribeOperation] | channelInfo -%}
{%- if channelInfo.hasHeaders %}
{{- headers(channelInfo.headerStructName, channelInfo.headers) -}}
{% endif %}
{%- endif %}
{%- endfor %}

// consumer service (to be implemented)
type ConsumerService interface {
	{% for channelName, channel in asyncapi.channels() -%}
	{%- if channel.hasSubscribe() -%}
	{%- set subscribeOperation = channel.subscribe()  -%}
	{%- set channelInfo = [channel, subscribeOperation] | channelInfo -%}
	{{- channelInfo.funcName -}}(ctx core.Context, key []byte, payload *{{channelInfo.payloadStructName}}{%- if channelInfo.hasHeaders -%}, headers *{{channelInfo.headerStructName}}{%- endif -%}) error
	{% endif %}
	{%- endfor %}
}

// build handlers
{% for channelName, channel in asyncapi.channels() -%}
{%- if channel.hasSubscribe() -%}
{%- set subscribeOperation = channel.subscribe()  -%}
{%- set channelInfo = [channel, subscribeOperation] | channelInfo -%}
type {{channelInfo.handlerStructName}} struct {
	consumerService ConsumerService
}

func (ch *{{- channelInfo.handlerStructName -}}) OnMessageReceived(ctx core.Context, message consumer.Message) error {
	// read payload
	payload := &{{channelInfo.payloadStructName}}{}
	err := message.BindBody(payload)
	if err != nil {
		return err
	}
	{%- if channelInfo.hasHeaders %}
	// read headers
	headers := {{channelInfo.headerStructName}}{}
	{% for headerName, header in channelInfo.headers.properties() %}
	headers.{{headerName | fixPropName}} = message.GetStringHeader("{{headerName}}")
	{%- endfor %}
	{%- endif %}
	return ch.consumerService.{{- channelInfo.funcName -}}(ctx, message.GetByteKey(), payload{%- if channelInfo.hasHeaders -%}, &headers{%- endif -%})
}
{% endif %}
{%- endfor %}

// build consumer manager
type Consumer struct {
	consumerMgr consumer.Manager
}

func (c *Consumer) StartConsumers() {
	err := c.consumerMgr.StartConsumers()
	if err != nil {
		panic(err)
	}
}

func NewConsumer(configStore core.ConfigStore, logger log.SimpleLogger, consumerService ConsumerService, sessionMgr session.Manager, typeHelper common.AsyncTypeHelper, tracer tracer.SimpleTracer) *Consumer {
	consumerMgr := consumer.NewManager(configStore, logger, sessionMgr, typeHelper, tracer)
	{% for channelName, channel in asyncapi.channels() -%}
	{%- if channel.hasSubscribe() -%}
	{%- set subscribeOperation = channel.subscribe()  -%}
	{%- set channelInfo = [channel, subscribeOperation] | channelInfo -%}
	consumerMgr.RegisterHandler("{{- channelName -}}",
		&{{channelInfo.handlerStructName}}{consumerService: consumerService})
	{% endif %}
	{%- endfor %}
	return &Consumer{consumerMgr: consumerMgr}
}
{% endmacro %}
