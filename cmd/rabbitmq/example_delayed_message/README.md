Rabbitmq's delayed message exchange example

* An exchange is created with the x-delayed-message type.
* Declare a queue user-published-queue that we bind to the routing key user.event.publish.
* Send an event with into user.event.publish with a 'delay' of 10 sec.
* Consume the queue user-published-queue and we can see that the event sent is only received after the delay of 10 sec.
