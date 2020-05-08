package server

import (
	"errors"
	"testing"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/CESARBR/knot-babeltower/pkg/mocks"
	"github.com/CESARBR/knot-babeltower/pkg/network"
	"github.com/CESARBR/knot-babeltower/pkg/thing/controllers"
)

func TestStart(t *testing.T) {
	type fields struct {
		logger          logging.Logger
		amqp            network.AmqpReceiver
		thingController controllers.ThingController
	}
	type args struct {
		started chan bool
		msgChan chan network.InMsg
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		expectedErr bool
	}{
		{
			"when started channel not provided an error should be returned",
			fields{
				&mocks.FakeLogger{},
				&mocks.FakeAmqpReceiver{},
				&mocks.FakeController{},
			},
			args{
				nil,
				make(chan network.InMsg),
			},
			true,
		},
		{
			"when msg channel not provided an error should be returned",
			fields{
				&mocks.FakeLogger{},
				&mocks.FakeAmqpReceiver{},
				&mocks.FakeController{},
			},
			args{
				make(chan bool, 1),
				nil,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := &MsgHandler{
				logger:          tt.fields.logger,
				amqp:            tt.fields.amqp,
				thingController: tt.fields.thingController,
			}
			err := mc.start(tt.args.started, tt.args.msgChan)
			if (err != nil) != tt.expectedErr {
				t.Errorf("msgHandler.start() error = %v, expectedErr %v", err, tt.expectedErr)
			}
		})
	}
}

func TestOnMsgReceived(t *testing.T) {
	type fields struct {
		logger          *mocks.FakeLogger
		amqp            network.AmqpReceiver
		thingController *mocks.FakeController
	}
	type args struct {
		started chan bool
		msgChan chan network.InMsg
		msg     network.InMsg
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		expectedErr bool
		mockArgs    map[string]string
	}{
		{
			"happy path device exchange without RPC should return no error",
			fields{
				&mocks.FakeLogger{},
				&mocks.FakeAmqpReceiver{},
				&mocks.FakeController{},
			},
			args{
				make(chan bool, 1),
				make(chan network.InMsg, 10),
				network.InMsg{
					Exchange: exchangeDevices,
					Body:     []byte{1, 2, 3},
					Headers: map[string]interface{}{
						"Authorization": "test-token",
					},
				},
			},
			false,
			map[string]string{
				bindingKeyRegisterDevice:   "Register",
				bindingKeyUnregisterDevice: "Unregister",
				bindingKeyRequestData:      "RequestData",
				bindingKeyUpdateData:       "UpdateData",
				bindingKeySchemaSent:       "UpdateSchema",
			},
		},
		{
			"happy path device exchange with RPC should return no error",
			fields{
				&mocks.FakeLogger{},
				&mocks.FakeAmqpReceiver{},
				&mocks.FakeController{},
			},
			args{
				make(chan bool, 1),
				make(chan network.InMsg, 10),
				network.InMsg{
					Exchange:      exchangeDevices,
					Body:          []byte{1, 2, 3},
					CorrelationID: "test-corrId",
					ReplyTo:       "test-reply_to",
					Headers: map[string]interface{}{
						"Authorization": "test-token",
					},
				},
			},
			false,
			map[string]string{
				bindingKeyAuthDevice:  "AuthDevice",
				bindingKeyListDevices: "ListDevices",
			},
		},
		{
			"happy path when send data to fanout exchange should ignore routing and call correct function",
			fields{
				&mocks.FakeLogger{},
				&mocks.FakeAmqpReceiver{},
				&mocks.FakeController{},
			},
			args{
				make(chan bool, 1),
				make(chan network.InMsg, 10),
				network.InMsg{
					Exchange: exchangeDataSent,
					Body:     []byte{1, 2, 3},
					Headers: map[string]interface{}{
						"Authorization": "test-token",
					},
				},
			},
			false,
			map[string]string{
				"any.key": "PublishData",
			},
		},
		{
			"missing correlation id on device exchange with RPC is optional and should return no error",
			fields{
				&mocks.FakeLogger{},
				&mocks.FakeAmqpReceiver{},
				&mocks.FakeController{},
			},
			args{
				make(chan bool, 1),
				make(chan network.InMsg, 10),
				network.InMsg{
					Exchange: exchangeDevices,
					Body:     []byte{1, 2, 3},
					ReplyTo:  "test-replyTo",
					Headers: map[string]interface{}{
						"Authorization": "test-token",
					},
				},
			},
			false,
			map[string]string{
				bindingKeyAuthDevice:  "AuthDevice",
				bindingKeyListDevices: "ListDevices",
			},
		},
		{
			"missing reply_to on device exchange with RPC should return missing reply to",
			fields{
				&mocks.FakeLogger{},
				&mocks.FakeAmqpReceiver{},
				&mocks.FakeController{Err: errors.New("missing replyTo")},
			},
			args{
				make(chan bool, 1),
				make(chan network.InMsg, 10),
				network.InMsg{
					Exchange:      exchangeDevices,
					Body:          []byte{1, 2, 3},
					CorrelationID: "test-corrId",
					Headers: map[string]interface{}{
						"Authorization": "test-token",
					},
				},
			},
			true,
			map[string]string{
				bindingKeyAuthDevice:  "AuthDevice",
				bindingKeyListDevices: "ListDevices",
			},
		},
		{
			"empty header should return missing authorization token",
			fields{
				&mocks.FakeLogger{},
				&mocks.FakeAmqpReceiver{},
				&mocks.FakeController{},
			},
			args{
				make(chan bool, 1),
				make(chan network.InMsg, 10),
				network.InMsg{Exchange: exchangeDevices, RoutingKey: bindingKeyRegisterDevice, Body: []byte{1, 2, 3}, Headers: map[string]interface{}{}},
			},
			true,
			nil,
		},
		{
			"unexpected exchange should return operation unsuported",
			fields{
				&mocks.FakeLogger{},
				&mocks.FakeAmqpReceiver{},
				&mocks.FakeController{},
			},
			args{
				make(chan bool, 1),
				make(chan network.InMsg, 10),
				network.InMsg{Exchange: "test", RoutingKey: bindingKeyRegisterDevice, Body: []byte{1, 2, 3}, Headers: map[string]interface{}{"Authorization": "test-token"}},
			},
			true,
			nil,
		},
		{
			"when message is from client and unexpected routing key should return operation unsuported",
			fields{
				&mocks.FakeLogger{},
				&mocks.FakeAmqpReceiver{},
				&mocks.FakeController{},
			},
			args{
				make(chan bool, 1),
				make(chan network.InMsg, 10),
				network.InMsg{
					Exchange:   exchangeDevices,
					RoutingKey: "key",
					Body:       []byte{1, 2, 3},
					Headers: map[string]interface{}{
						"Authorization": "test-token",
					},
				},
			},
			true,
			nil,
		},
		{
			"when header is not provided should return an error",
			fields{
				&mocks.FakeLogger{},
				&mocks.FakeAmqpReceiver{},
				&mocks.FakeController{},
			},
			args{
				make(chan bool, 1),
				make(chan network.InMsg, 10),
				network.InMsg{Exchange: exchangeDevices, RoutingKey: bindingKeyRegisterDevice, Body: []byte{1, 2, 3}, Headers: nil},
			},
			true,
			nil,
		},
		{
			"when body is not provided should return an error",
			fields{
				&mocks.FakeLogger{},
				&mocks.FakeAmqpReceiver{},
				&mocks.FakeController{},
			},
			args{
				make(chan bool, 1),
				make(chan network.InMsg, 10),
				network.InMsg{
					Exchange:   exchangeDevices,
					RoutingKey: bindingKeyRegisterDevice,
					Body:       nil,
					Headers: map[string]interface{}{
						"Authorization": "test-token",
					},
				},
			},
			true,
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := &MsgHandler{
				logger:          tt.fields.logger,
				amqp:            tt.fields.amqp,
				thingController: tt.fields.thingController,
			}
			if tt.mockArgs != nil {
				for key, value := range tt.mockArgs {
					if len(value) > 0 {
						tt.fields.thingController.On(value).Return(tt.fields.thingController.Err).Once()
					}
					tt.args.msg.RoutingKey = key
					tt.args.msgChan <- tt.args.msg
					err := mc.onMsgReceived(tt.args.msgChan)
					if (err != nil) != tt.expectedErr {
						t.Errorf("msgHandler.onMsgReceived() error = %v, expectedErr %v", err, tt.expectedErr)
					}
					if len(value) > 0 {
						tt.fields.thingController.AssertExpectations(t)
					}
				}
			} else {
				tt.args.msgChan <- tt.args.msg
				err := mc.onMsgReceived(tt.args.msgChan)
				if (err != nil) != tt.expectedErr {
					t.Errorf("msgHandler.onMsgReceived() error = %v, expectedErr %v", err, tt.expectedErr)
				}
			}
		})
	}
}

func TestSubscribeToMessagesCalls(t *testing.T) {
	type fields struct {
		logger          logging.Logger
		amqp            *mocks.FakeAmqpReceiver
		thingController *mocks.FakeController
	}
	type mockArgs struct {
		queue          string
		exchange       string
		exchangeType   string
		key            string
		expectedReturn error
	}
	type args struct {
		msgChan       chan network.InMsg
		subscribeArgs []mockArgs
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		expectedErr  bool
		verifyCalled bool
	}{
		{
			"Happy path correct call",
			fields{
				&mocks.FakeLogger{},
				&mocks.FakeAmqpReceiver{},
				&mocks.FakeController{},
			},
			args{
				make(chan network.InMsg),
				[]mockArgs{
					{queueNameCommands, exchangeDevices, exchangeDevicesType, bindingKeyRegisterDevice, nil},
					{queueNameCommands, exchangeDevices, exchangeDevicesType, bindingKeyUnregisterDevice, nil},
					{queueNameCommands, exchangeDevices, exchangeDevicesType, bindingKeyRequestData, nil},
					{queueNameCommands, exchangeDevices, exchangeDevicesType, bindingKeyUpdateData, nil},
					{queueNameCommands, exchangeDevices, exchangeDevicesType, bindingKeySchemaSent, nil},
					{queueNameCommands, exchangeDevices, exchangeDevicesType, bindingKeyAuthDevice, nil},
					{queueNameCommands, exchangeDevices, exchangeDevicesType, bindingKeyListDevices, nil},
					{queueNameEvents, exchangeDataSent, exchangeDataSentType, bindingKeyEmpty, nil},
				},
			},
			false,
			true,
		},
		{
			"when first subscribe mock returns error should return error and not call the other subscribes",
			fields{
				&mocks.FakeLogger{},
				&mocks.FakeAmqpReceiver{},
				&mocks.FakeController{},
			},
			args{
				make(chan network.InMsg),
				[]mockArgs{
					{queueNameCommands, exchangeDevices, exchangeDevicesType, bindingKeyRegisterDevice, errors.New("missing routing key argument on subscribe")},
				},
			},
			true,
			true,
		},
		{
			"when any middle subscribe mock returns error should return error and not call the following subscribes",
			fields{
				&mocks.FakeLogger{},
				&mocks.FakeAmqpReceiver{},
				&mocks.FakeController{},
			},
			args{
				make(chan network.InMsg),
				[]mockArgs{
					{queueNameCommands, exchangeDevices, exchangeDevicesType, bindingKeyRegisterDevice, nil},
					{queueNameCommands, exchangeDevices, exchangeDevicesType, bindingKeyUnregisterDevice, nil},
					{queueNameCommands, exchangeDevices, exchangeDevicesType, bindingKeyRequestData, nil},
					{queueNameCommands, exchangeDevices, exchangeDevicesType, bindingKeyUpdateData, errors.New("missing routing key argument on subscribe")},
				},
			},
			true,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := &MsgHandler{
				logger:          tt.fields.logger,
				amqp:            tt.fields.amqp,
				thingController: tt.fields.thingController,
			}
			if tt.verifyCalled {
				for _, item := range tt.args.subscribeArgs {
					tt.fields.amqp.On("OnMessage", tt.args.msgChan, item.queue, item.exchange, item.exchangeType, item.key).Once().Return(item.expectedReturn)
				}
			}
			if err := mc.subscribeToMessages(tt.args.msgChan); (err != nil) != tt.expectedErr {
				t.Errorf("MsgHandler.subscribeToMessages() error = %v, expectedErr %v", err, tt.expectedErr)
			}

			if tt.verifyCalled {
				tt.fields.amqp.AssertExpectations(t)
			}
		})
	}
}
