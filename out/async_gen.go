package async

import (
	"bitbucket.org/integrapartners/prism-async-go/common"
	"bitbucket.org/integrapartners/prism-async-go/producer"
	"bitbucket.org/integrapartners/prism-async-go/consumer"
	"bitbucket.org/integrapartners/prism-core-go/session"
	"bitbucket.org/integrapartners/prism-core-go/core"
	"bitbucket.org/integrapartners/prism-core-go/log"
	"bitbucket.org/integrapartners/prism-core-go/tracer"
	
)

// load schemas
type UploadMemberPayload struct {
    ExtAltId *string `json:"extAltId,omitempty"`
    LobId int64 `json:"lobId,omitempty"`
    InsuranceNum string `json:"insuranceNum,omitempty"`
    BenefitType string `json:"benefitType,omitempty"`
    MedicareId string `json:"medicareId,omitempty"`
    MedicaidId string `json:"medicaidId,omitempty"`
    FirstName string `json:"firstName,omitempty"`
    LastName string `json:"lastName,omitempty"`
    Dob string `json:"dob,omitempty"`
    Gender string `json:"gender,omitempty"`
    CvrgStartDT common.Date `json:"cvrgStartDT,omitempty"`
    CvrgEndDT *common.Date `json:"cvrgEndDT,omitempty"`
    Language string `json:"language,omitempty"`
    Ssn *int64 `json:"ssn,omitempty"`
    Contacts []*Contacts `json:"contacts,omitempty"`
    Addresses []*Addresses `json:"addresses,omitempty"`
}

type Contacts struct {
    ContactType string `json:"contactType,omitempty"`
    ContactValue string `json:"contactValue,omitempty"`
}

type Addresses struct {
    Relation string `json:"relation,omitempty"`
    AddressType string `json:"addressType,omitempty"`
    Address1 string `json:"address1,omitempty"`
    Address2 string `json:"address2,omitempty"`
    City string `json:"city,omitempty"`
    State string `json:"state,omitempty"`
    Zip string `json:"zip,omitempty"`
}
// load producers
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

// Headers object
type PublishUploadMemberHeaders struct {
    IdempotencyKey string
}


// Publish information about triggered schedule
func (p *Producer) PublishUploadMember(ctx core.Context,
	messageKey []byte,
	payload *UploadMemberPayload,
	headers PublishUploadMemberHeaders) (partition int32, offset int64, err error) {

	message := p.msgProducer.BuildMessage().
		AddByteKey(messageKey).
		AddStringHeader("idempotency-key", headers.IdempotencyKey).
		AddBody(&payload).
		Build()
	return p.msgProducer.SendMessage(ctx, "event.member.edi.enrollment.v1", message)
}



// load consumer
// consumer headers
type ReceiveUploadMemberHeaders struct {
    IdempotencyKey string
}


// consumer service (to be implemented)
type ConsumerService interface {
	ReceiveUploadMember(ctx core.Context, key []byte, payload *UploadMemberPayload, headers *ReceiveUploadMemberHeaders) error
	
}

// build handlers
type ReceiveUploadMemberHandler struct {
	consumerService ConsumerService
}

func (ch *ReceiveUploadMemberHandler) OnMessageReceived(ctx core.Context, message consumer.Message) error {
	// read payload
	payload := &UploadMemberPayload{}
	err := message.BindBody(payload)
	if err != nil {
		return err
	}
	// read headers
	headers := ReceiveUploadMemberHeaders{}
	
	headers.IdempotencyKey = message.GetStringHeader("idempotency-key")
	return ch.consumerService.ReceiveUploadMember(ctx, message.GetByteKey(), payload, &headers)
}


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
	consumerMgr.RegisterHandler("event.member.edi.enrollment.v1",
		&ReceiveUploadMemberHandler{consumerService: consumerService})
	
	return &Consumer{consumerMgr: consumerMgr}
}

