package main

import (
	"context"
	"time"

	"github.com/k0kubun/pp/v3"
	"github.com/quenbyako/asyncapi-playground/nsq"
)

func main() {
	userController, err := nsq.NewController("localhost:4150")
	if err != nil {
		panic(err)
	}

	appController, err := nsq.NewController("localhost:4150")
	if err != nil {
		panic(err)
	}

	user, err := NewUserController(userController)
	if err != nil {
		panic(err)
	}

	if err := user.SubscribeAll(context.Background(), subscriber{}); err != nil {
		panic(err)
	}

	go func() {
		app, err := NewAppController(appController)
		if err != nil {
			panic(err)
		}

		time.Sleep(500 * time.Millisecond)

		app.PublishOrderNotification(context.Background(), NewOrderMessage{
			Headers: NewOrderMessageHeaders{
				OrderSource: ptr("1234"),
			},
		})

		time.Sleep(500 * time.Millisecond)

		app.PublishOrderNotification(context.Background(), NewOrderMessage{
			Headers: NewOrderMessageHeaders{
				OrderSource: ptr("Google"),
			},
		})

		time.Sleep(500 * time.Millisecond)

		app.PublishOrderNotification(context.Background(), NewOrderMessage{
			Headers: NewOrderMessageHeaders{
				OrderSource: ptr("Apple"),
			},
		})
	}()

	select {}
}

type subscriber struct{}

func (subscriber) OrderCancellation(ctx context.Context, msg CANCELLATIONSMessage) error {
	pp.Println(msg)
	return nil
}
func (subscriber) OrderNotification(ctx context.Context, msg NewOrderMessage) error {
	pp.Println(msg)
	return nil
}

func ptr[T any](t T) *T { return &t }
