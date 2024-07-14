// Package "main" provides primitives to interact with the AsyncAPI specification.
//
// Code generated by github.com/lerenn/asyncapi-codegen version v0.42.4 DO NOT EDIT.
package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
)

// AppController is the structure that provides publishing capabilities to the
// developer and and connect the broker with the App
type AppController struct {
	controller
}

// NewAppController links the App to the broker
func NewAppController(bc extensions.BrokerController, options ...ControllerOption) (*AppController, error) {
	// Check if broker controller has been provided
	if bc == nil {
		return nil, extensions.ErrNilBrokerController
	}

	// Create default controller
	controller := controller{
		broker:        bc,
		subscriptions: make(map[string]extensions.BrokerChannelSubscription),
		logger:        extensions.DummyLogger{},
		middlewares:   make([]extensions.Middleware, 0),
		errorHandler:  extensions.DefaultErrorHandler(),
	}

	// Apply options
	for _, option := range options {
		option(&controller)
	}

	return &AppController{controller: controller}, nil
}

func (c AppController) wrapMiddlewares(
	middlewares []extensions.Middleware,
	callback extensions.NextMiddleware,
) func(ctx context.Context, msg *extensions.BrokerMessage) error {
	var called bool

	// If there is no more middleware
	if len(middlewares) == 0 {
		return func(ctx context.Context, msg *extensions.BrokerMessage) error {
			// Call the callback if it exists and it has not been called already
			if callback != nil && !called {
				called = true
				return callback(ctx)
			}

			// Nil can be returned, as the callback has already been called
			return nil
		}
	}

	// Get the next function to call from next middlewares or callback
	next := c.wrapMiddlewares(middlewares[1:], callback)

	// Wrap middleware into a check function that will call execute the middleware
	// and call the next wrapped middleware if the returned function has not been
	// called already
	return func(ctx context.Context, msg *extensions.BrokerMessage) error {
		// Call the middleware and the following if it has not been done already
		if !called {
			// Create the next call with the context and the message
			nextWithArgs := func(ctx context.Context) error {
				return next(ctx, msg)
			}

			// Call the middleware and register it as already called
			called = true
			if err := middlewares[0](ctx, msg, nextWithArgs); err != nil {
				return err
			}

			// If next has already been called in middleware, it should not be executed again
			return nextWithArgs(ctx)
		}

		// Nil can be returned, as the next middleware has already been called
		return nil
	}
}

func (c AppController) executeMiddlewares(ctx context.Context, msg *extensions.BrokerMessage, callback extensions.NextMiddleware) error {
	// Wrap middleware to have 'next' function when calling them
	wrapped := c.wrapMiddlewares(c.middlewares, callback)

	// Execute wrapped middlewares
	return wrapped(ctx, msg)
}

func addAppContextValues(ctx context.Context, path string) context.Context {
	ctx = context.WithValue(ctx, extensions.ContextKeyIsVersion, "0.0.1")
	ctx = context.WithValue(ctx, extensions.ContextKeyIsProvider, "app")
	return context.WithValue(ctx, extensions.ContextKeyIsChannel, path)
}

// Close will clean up any existing resources on the controller
func (c *AppController) Close(ctx context.Context) {
	// Unsubscribing remaining channels
}

// PublishOrderCancellation will publish messages to 'CANCELLATIONS' channel
func (c *AppController) PublishOrderCancellation(
	ctx context.Context,
	msg CANCELLATIONSMessage,
) error {
	// Get channel path
	path := "CANCELLATIONS"

	// Set context
	ctx = addAppContextValues(ctx, path)
	ctx = context.WithValue(ctx, extensions.ContextKeyIsDirection, "publication")

	// Convert to BrokerMessage
	brokerMsg, err := msg.toBrokerMessage()
	if err != nil {
		return err
	}

	// Set broker message to context
	ctx = context.WithValue(ctx, extensions.ContextKeyIsBrokerMessage, brokerMsg.String())

	// Publish the message on event-broker through middlewares
	return c.executeMiddlewares(ctx, &brokerMsg, func(ctx context.Context) error {
		return c.broker.Publish(ctx, path, brokerMsg)
	})
}

// PublishOrderNotification will publish messages to 'ORDERS' channel
func (c *AppController) PublishOrderNotification(
	ctx context.Context,
	msg NewOrderMessage,
) error {
	// Get channel path
	path := "ORDERS"

	// Set context
	ctx = addAppContextValues(ctx, path)
	ctx = context.WithValue(ctx, extensions.ContextKeyIsDirection, "publication")

	// Convert to BrokerMessage
	brokerMsg, err := msg.toBrokerMessage()
	if err != nil {
		return err
	}

	// Set broker message to context
	ctx = context.WithValue(ctx, extensions.ContextKeyIsBrokerMessage, brokerMsg.String())

	// Publish the message on event-broker through middlewares
	return c.executeMiddlewares(ctx, &brokerMsg, func(ctx context.Context) error {
		return c.broker.Publish(ctx, path, brokerMsg)
	})
}

// UserSubscriber represents all handlers that are expecting messages for User
type UserSubscriber interface {
	// OrderCancellation subscribes to messages placed on the 'CANCELLATIONS' channel
	OrderCancellation(ctx context.Context, msg CANCELLATIONSMessage) error

	// OrderNotification subscribes to messages placed on the 'ORDERS' channel
	OrderNotification(ctx context.Context, msg NewOrderMessage) error
}

// UserController is the structure that provides publishing capabilities to the
// developer and and connect the broker with the User
type UserController struct {
	controller
}

// NewUserController links the User to the broker
func NewUserController(bc extensions.BrokerController, options ...ControllerOption) (*UserController, error) {
	// Check if broker controller has been provided
	if bc == nil {
		return nil, extensions.ErrNilBrokerController
	}

	// Create default controller
	controller := controller{
		broker:        bc,
		subscriptions: make(map[string]extensions.BrokerChannelSubscription),
		logger:        extensions.DummyLogger{},
		middlewares:   make([]extensions.Middleware, 0),
		errorHandler:  extensions.DefaultErrorHandler(),
	}

	// Apply options
	for _, option := range options {
		option(&controller)
	}

	return &UserController{controller: controller}, nil
}

func (c UserController) wrapMiddlewares(
	middlewares []extensions.Middleware,
	callback extensions.NextMiddleware,
) func(ctx context.Context, msg *extensions.BrokerMessage) error {
	var called bool

	// If there is no more middleware
	if len(middlewares) == 0 {
		return func(ctx context.Context, msg *extensions.BrokerMessage) error {
			// Call the callback if it exists and it has not been called already
			if callback != nil && !called {
				called = true
				return callback(ctx)
			}

			// Nil can be returned, as the callback has already been called
			return nil
		}
	}

	// Get the next function to call from next middlewares or callback
	next := c.wrapMiddlewares(middlewares[1:], callback)

	// Wrap middleware into a check function that will call execute the middleware
	// and call the next wrapped middleware if the returned function has not been
	// called already
	return func(ctx context.Context, msg *extensions.BrokerMessage) error {
		// Call the middleware and the following if it has not been done already
		if !called {
			// Create the next call with the context and the message
			nextWithArgs := func(ctx context.Context) error {
				return next(ctx, msg)
			}

			// Call the middleware and register it as already called
			called = true
			if err := middlewares[0](ctx, msg, nextWithArgs); err != nil {
				return err
			}

			// If next has already been called in middleware, it should not be executed again
			return nextWithArgs(ctx)
		}

		// Nil can be returned, as the next middleware has already been called
		return nil
	}
}

func (c UserController) executeMiddlewares(ctx context.Context, msg *extensions.BrokerMessage, callback extensions.NextMiddleware) error {
	// Wrap middleware to have 'next' function when calling them
	wrapped := c.wrapMiddlewares(c.middlewares, callback)

	// Execute wrapped middlewares
	return wrapped(ctx, msg)
}

func addUserContextValues(ctx context.Context, path string) context.Context {
	ctx = context.WithValue(ctx, extensions.ContextKeyIsVersion, "0.0.1")
	ctx = context.WithValue(ctx, extensions.ContextKeyIsProvider, "user")
	return context.WithValue(ctx, extensions.ContextKeyIsChannel, path)
}

// Close will clean up any existing resources on the controller
func (c *UserController) Close(ctx context.Context) {
	// Unsubscribing remaining channels
	c.UnsubscribeAll(ctx)

	c.logger.Info(ctx, "Closed user controller")
}

// SubscribeAll will subscribe to channels without parameters on which the app is expecting messages.
// For channels with parameters, they should be subscribed independently.
func (c *UserController) SubscribeAll(ctx context.Context, as UserSubscriber) error {
	if as == nil {
		return extensions.ErrNilUserSubscriber
	}

	if err := c.SubscribeOrderCancellation(ctx, as.OrderCancellation); err != nil {
		return err
	}
	if err := c.SubscribeOrderNotification(ctx, as.OrderNotification); err != nil {
		return err
	}

	return nil
}

// UnsubscribeAll will unsubscribe all remaining subscribed channels
func (c *UserController) UnsubscribeAll(ctx context.Context) {
	c.UnsubscribeOrderCancellation(ctx)
	c.UnsubscribeOrderNotification(ctx)
}

// SubscribeOrderCancellation will subscribe to new messages from 'CANCELLATIONS' channel.
//
// Callback function 'fn' will be called each time a new message is received.
func (c *UserController) SubscribeOrderCancellation(
	ctx context.Context,
	fn func(ctx context.Context, msg CANCELLATIONSMessage) error,
) error {
	// Get channel path
	path := "CANCELLATIONS"

	// Set context
	ctx = addUserContextValues(ctx, path)
	ctx = context.WithValue(ctx, extensions.ContextKeyIsDirection, "reception")

	// Check if there is already a subscription
	_, exists := c.subscriptions[path]
	if exists {
		err := fmt.Errorf("%w: %q channel is already subscribed", extensions.ErrAlreadySubscribedChannel, path)
		c.logger.Error(ctx, err.Error())
		return err
	}

	// Subscribe to broker channel
	sub, err := c.broker.Subscribe(ctx, path)
	if err != nil {
		c.logger.Error(ctx, err.Error())
		return err
	}
	c.logger.Info(ctx, "Subscribed to channel")

	// Asynchronously listen to new messages and pass them to app subscriber
	go func() {
		for {
			// Wait for next message
			acknowledgeableBrokerMessage, open := <-sub.MessagesChannel()

			// If subscription is closed and there is no more message
			// (i.e. uninitialized message), then exit the function
			if !open && acknowledgeableBrokerMessage.IsUninitialized() {
				return
			}

			// Set broker message to context
			ctx = context.WithValue(ctx, extensions.ContextKeyIsBrokerMessage, acknowledgeableBrokerMessage.String())

			// Execute middlewares before handling the message
			if err := c.executeMiddlewares(ctx, &acknowledgeableBrokerMessage.BrokerMessage, func(ctx context.Context) error {
				// Process message
				msg, err := brokerMessageToCANCELLATIONSMessage(acknowledgeableBrokerMessage.BrokerMessage)
				if err != nil {
					return err
				}

				// Execute the subscription function
				if err := fn(ctx, msg); err != nil {
					return err
				}

				acknowledgeableBrokerMessage.Ack()

				return nil
			}); err != nil {
				c.errorHandler(ctx, path, &acknowledgeableBrokerMessage, err)
				// On error execute the acknowledgeableBrokerMessage nack() function and
				// let the BrokerAcknowledgment decide what is the right nack behavior for the broker
				acknowledgeableBrokerMessage.Nak()
			}
		}
	}()

	// Add the cancel channel to the inside map
	c.subscriptions[path] = sub

	return nil
}

// UnsubscribeOrderCancellation will unsubscribe messages from 'CANCELLATIONS' channel.
// A timeout can be set in context to avoid blocking operation, if needed.
func (c *UserController) UnsubscribeOrderCancellation(ctx context.Context) {
	// Get channel path
	path := "CANCELLATIONS"

	// Check if there subscribers for this channel
	sub, exists := c.subscriptions[path]
	if !exists {
		return
	}

	// Set context
	ctx = addUserContextValues(ctx, path)

	// Stop the subscription
	sub.Cancel(ctx)

	// Remove if from the subscribers
	delete(c.subscriptions, path)

	c.logger.Info(ctx, "Unsubscribed from channel")
}

// SubscribeOrderNotification will subscribe to new messages from 'ORDERS' channel.
//
// Callback function 'fn' will be called each time a new message is received.
func (c *UserController) SubscribeOrderNotification(
	ctx context.Context,
	fn func(ctx context.Context, msg NewOrderMessage) error,
) error {
	// Get channel path
	path := "ORDERS"

	// Set context
	ctx = addUserContextValues(ctx, path)
	ctx = context.WithValue(ctx, extensions.ContextKeyIsDirection, "reception")

	// Check if there is already a subscription
	_, exists := c.subscriptions[path]
	if exists {
		err := fmt.Errorf("%w: %q channel is already subscribed", extensions.ErrAlreadySubscribedChannel, path)
		c.logger.Error(ctx, err.Error())
		return err
	}

	// Subscribe to broker channel
	sub, err := c.broker.Subscribe(ctx, path)
	if err != nil {
		c.logger.Error(ctx, err.Error())
		return err
	}
	c.logger.Info(ctx, "Subscribed to channel")

	// Asynchronously listen to new messages and pass them to app subscriber
	go func() {
		for {
			// Wait for next message
			acknowledgeableBrokerMessage, open := <-sub.MessagesChannel()

			// If subscription is closed and there is no more message
			// (i.e. uninitialized message), then exit the function
			if !open && acknowledgeableBrokerMessage.IsUninitialized() {
				return
			}

			// Set broker message to context
			ctx = context.WithValue(ctx, extensions.ContextKeyIsBrokerMessage, acknowledgeableBrokerMessage.String())

			// Execute middlewares before handling the message
			if err := c.executeMiddlewares(ctx, &acknowledgeableBrokerMessage.BrokerMessage, func(ctx context.Context) error {
				// Process message
				msg, err := brokerMessageToNewOrderMessage(acknowledgeableBrokerMessage.BrokerMessage)
				if err != nil {
					return err
				}

				// Execute the subscription function
				if err := fn(ctx, msg); err != nil {
					return err
				}

				acknowledgeableBrokerMessage.Ack()

				return nil
			}); err != nil {
				c.errorHandler(ctx, path, &acknowledgeableBrokerMessage, err)
				// On error execute the acknowledgeableBrokerMessage nack() function and
				// let the BrokerAcknowledgment decide what is the right nack behavior for the broker
				acknowledgeableBrokerMessage.Nak()
			}
		}
	}()

	// Add the cancel channel to the inside map
	c.subscriptions[path] = sub

	return nil
}

// UnsubscribeOrderNotification will unsubscribe messages from 'ORDERS' channel.
// A timeout can be set in context to avoid blocking operation, if needed.
func (c *UserController) UnsubscribeOrderNotification(ctx context.Context) {
	// Get channel path
	path := "ORDERS"

	// Check if there subscribers for this channel
	sub, exists := c.subscriptions[path]
	if !exists {
		return
	}

	// Set context
	ctx = addUserContextValues(ctx, path)

	// Stop the subscription
	sub.Cancel(ctx)

	// Remove if from the subscribers
	delete(c.subscriptions, path)

	c.logger.Info(ctx, "Unsubscribed from channel")
}

// AsyncAPIVersion is the version of the used AsyncAPI document
const AsyncAPIVersion = "0.0.1"

// controller is the controller that will be used to communicate with the broker
// It will be used internally by AppController and UserController
type controller struct {
	// broker is the broker controller that will be used to communicate
	broker extensions.BrokerController
	// subscriptions is a map of all subscriptions
	subscriptions map[string]extensions.BrokerChannelSubscription
	// logger is the logger that will be used² to log operations on controller
	logger extensions.Logger
	// middlewares are the middlewares that will be executed when sending or
	// receiving messages
	middlewares []extensions.Middleware
	// handler to handle errors from consumers and middlewares
	errorHandler extensions.ErrorHandler
}

// ControllerOption is the type of the options that can be passed
// when creating a new Controller
type ControllerOption func(controller *controller)

// WithLogger attaches a logger to the controller
func WithLogger(logger extensions.Logger) ControllerOption {
	return func(controller *controller) {
		controller.logger = logger
	}
}

// WithMiddlewares attaches middlewares that will be executed when sending or receiving messages
func WithMiddlewares(middlewares ...extensions.Middleware) ControllerOption {
	return func(controller *controller) {
		controller.middlewares = middlewares
	}
}

// WithErrorHandler attaches a errorhandler to handle errors from subscriber functions
func WithErrorHandler(handler extensions.ErrorHandler) ControllerOption {
	return func(controller *controller) {
		controller.errorHandler = handler
	}
}

type MessageWithCorrelationID interface {
	CorrelationID() string
	SetCorrelationID(id string)
}

type Error struct {
	Channel string
	Err     error
}

func (e *Error) Error() string {
	return fmt.Sprintf("channel %q: err %v", e.Channel, e.Err)
}

// CANCELLATIONSMessageHeaders is a schema from the AsyncAPI specification required in messages
type CANCELLATIONSMessageHeaders struct {
	// Description: unique ID of the application used to register the order cancellation
	OrderSource *string `json:"orderSource"`
}

// CANCELLATIONSMessagePayload is a schema from the AsyncAPI specification required in messages
type CANCELLATIONSMessagePayload struct {
	// Description: id of the order that was cancelled
	Orderid string `json:"orderid" validate:"required"`
}

// CANCELLATIONSMessage is the message expected for 'CANCELLATIONSMessage' channel.
// NOTE: More info about how the cancellation notifications are **created** and **used**.
type CANCELLATIONSMessage struct {
	// Headers will be used to fill the message headers
	Headers CANCELLATIONSMessageHeaders

	// Payload will be inserted in the message payload
	Payload CANCELLATIONSMessagePayload
}

func NewCANCELLATIONSMessage() CANCELLATIONSMessage {
	var msg CANCELLATIONSMessage

	return msg
}

// brokerMessageToCANCELLATIONSMessage will fill a new CANCELLATIONSMessage with data from generic broker message
func brokerMessageToCANCELLATIONSMessage(bMsg extensions.BrokerMessage) (CANCELLATIONSMessage, error) {
	var msg CANCELLATIONSMessage

	// Unmarshal payload to expected message payload format
	err := json.Unmarshal(bMsg.Payload, &msg.Payload)
	if err != nil {
		return msg, err
	}

	// Get each headers from broker message
	for k, v := range bMsg.Headers {
		switch {
		case k == "orderSource": // Retrieving OrderSource header
			h := string(v)
			msg.Headers.OrderSource = &h
		default:
			// TODO: log unknown error
		}
	}

	// TODO: run checks on msg type

	return msg, nil
}

// toBrokerMessage will generate a generic broker message from CANCELLATIONSMessage data
func (msg CANCELLATIONSMessage) toBrokerMessage() (extensions.BrokerMessage, error) {
	// TODO: implement checks on message

	// Marshal payload to JSON
	payload, err := json.Marshal(msg.Payload)
	if err != nil {
		return extensions.BrokerMessage{}, err
	}

	// Add each headers to broker message
	headers := make(map[string][]byte, 1)

	// Adding OrderSource header
	if msg.Headers.OrderSource != nil {
		headers["orderSource"] = []byte(*msg.Headers.OrderSource)
	}

	return extensions.BrokerMessage{
		Headers: headers,
		Payload: payload,
	}, nil
}

// NewOrderMessageHeaders is a schema from the AsyncAPI specification required in messages
type NewOrderMessageHeaders struct {
	// Description: unique ID of the application used to create the order
	OrderSource *string `json:"orderSource"`
}

// NewOrderMessagePayload is a schema from the AsyncAPI specification required in messages
type NewOrderMessagePayload struct{}

// NewOrderMessage is the message expected for 'NewOrderMessage' channel.
// NOTE: More info about how the order notifications are **created** and **used**.
type NewOrderMessage struct {
	// Headers will be used to fill the message headers
	Headers NewOrderMessageHeaders

	// Payload will be inserted in the message payload
	Payload NewOrderMessagePayload
}

func NewNewOrderMessage() NewOrderMessage {
	var msg NewOrderMessage

	return msg
}

// brokerMessageToNewOrderMessage will fill a new NewOrderMessage with data from generic broker message
func brokerMessageToNewOrderMessage(bMsg extensions.BrokerMessage) (NewOrderMessage, error) {
	var msg NewOrderMessage

	// Unmarshal payload to expected message payload format
	err := json.Unmarshal(bMsg.Payload, &msg.Payload)
	if err != nil {
		return msg, err
	}

	// Get each headers from broker message
	for k, v := range bMsg.Headers {
		switch {
		case k == "orderSource": // Retrieving OrderSource header
			h := string(v)
			msg.Headers.OrderSource = &h
		default:
			// TODO: log unknown error
		}
	}

	// TODO: run checks on msg type

	return msg, nil
}

// toBrokerMessage will generate a generic broker message from NewOrderMessage data
func (msg NewOrderMessage) toBrokerMessage() (extensions.BrokerMessage, error) {
	// TODO: implement checks on message

	// Marshal payload to JSON
	payload, err := json.Marshal(msg.Payload)
	if err != nil {
		return extensions.BrokerMessage{}, err
	}

	// Add each headers to broker message
	headers := make(map[string][]byte, 1)

	// Adding OrderSource header
	if msg.Headers.OrderSource != nil {
		headers["orderSource"] = []byte(*msg.Headers.OrderSource)
	}

	return extensions.BrokerMessage{
		Headers: headers,
		Payload: payload,
	}, nil
}

// NewOrderSchema is a schema from the AsyncAPI specification required in messages
type NewOrderSchema struct {
	Catalogid string  `json:"catalogid" validate:"required"`
	Cost      float64 `json:"cost" validate:"required"`

	// Description: username for the order - should be the email address of the customer
	Customer string   `json:"customer" validate:"required"`
	Discount *float64 `json:"discount"`

	// Description: unique order id on the Modo Jeans order system
	Id       string `json:"id" validate:"required"`
	Itemid   string `json:"itemid" validate:"required"`
	Quantity int32  `json:"quantity" validate:"required"`
	Version  string `json:"version" validate:"required,oneof=v1 v2"`
}

// NewOrderV1Schema is a schema from the AsyncAPI specification required in messages
type NewOrderV1Schema struct {
	Cost float64 `json:"cost" validate:"required"`

	// Description: unique order id on the Modo Jeans order system
	Id      string `json:"id" validate:"required"`
	Itemid  string `json:"itemid" validate:"required"`
	Version string `json:"version" validate:"required,oneof=v1"`
}

// NewOrderV2Schema is a schema from the AsyncAPI specification required in messages
type NewOrderV2Schema struct {
	Catalogid string  `json:"catalogid" validate:"required"`
	Cost      float64 `json:"cost" validate:"required"`

	// Description: username for the order - should be the email address of the customer
	Customer string   `json:"customer" validate:"required"`
	Discount *float64 `json:"discount"`

	// Description: unique order id on the Modo Jeans order system
	Id       string `json:"id" validate:"required"`
	Quantity int32  `json:"quantity" validate:"required"`
	Version  string `json:"version" validate:"required,oneof=v2"`
}

const (
	// CANCELLATIONSPath is the constant representing the 'CANCELLATIONS' channel path.
	CANCELLATIONSPath = "CANCELLATIONS"
	// ORDERSPath is the constant representing the 'ORDERS' channel path.
	ORDERSPath = "ORDERS"
)

// ChannelsPaths is an array of all channels paths
var ChannelsPaths = []string{
	CANCELLATIONSPath,
	ORDERSPath,
}