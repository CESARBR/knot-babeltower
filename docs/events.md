# Babeltower Events API - v2.1.1

This document describes the events `babeltower` is able to receive and send. They are gruped based on the external clients point of view, i.e. publishing or subscribing to the topics. In each section, it is provided information about the header, payload and protocol binding details of the event.

## Content

- [Publish](#publish) (external clients can publish to):
  - [device.register](#device-register)
  - [device.unregister](#device-unregister)
  - [device.config.sent](#device-config-sent)
  - [device.list](#device-list)
  - [device.auth](#device-auth)
  - [data.sent](#data-sent)
  - [data.request](#data-request)
  - [data.update](#data-update)

- [Subscribe](#Subscribe) (external clients can subscribe to):
  - [device.registered](#device-registered)
  - [device.unregistered](#device-unregistered)
  - [device.config.updated](#device-config-updated)
  - [data.published](#data-published)
  - [data.util](#data-util)
  - [data.[sessionId].published](#data-session-published)
  - [device.[id].data.request](#device-<id>-data-request)
  - [device.[id].data.update](#device-<id>-data-update)

-----------------------------------------------------------------

## Publish

This section describes the events that this service can receive from the external applications.

### **device.register** <a name="device-register"></a>

Event-command to register a new thing on the things registry. The operation response is sent through [`device.registered`](#device-registered) event.

<details>
  <summary>Headers</summary>

  - `token` **String** user's token

</details>

<details>
  <summary>Payload</summary>

  JSON in the following format:

  - `id` **String** thing's ID
  - `name` **String** thing's name

  Example:

  ```json
  {
    "id": "fbe64efa6c7f717e",
    "name": "KNoT Thing"
  }
  ```
</details>

<details>
  <summary>AMQP Binding</summary>

  - AMQP:
    - Exchange:
      - Type: direct
      - Name: device
      - Durable: `true`
      - Auto-delete: `false`
    - Routing key: device.register

</details>

### **device.unregister** <a name="device-unregister"></a>

Event-command to remove a thing from the things registry. The operation response is sent through [`device.unregistered`](#device-registered) event.

<details>
  <summary>Headers</summary>

  - `token` **String** user's token

</details>

<details>
  <summary>Payload</summary>

  JSON in the following format:

  - `id` **String** thing's ID

  Example:

  ```json
  {
    "id": "fbe64efa6c7f717e"
  }
  ```
</details>

<details>
  <summary>AMQP Binding</summary>

  - Exchange:
    - Type: direct
    - Name: device
    - Durable: `true`
    - Auto-delete: `false`
  - Routing key: device.unregister

</details>

### **device.config.sent** <a name="device-config-sent"></a>

Event that represents a device sending its config to the services that are interested. After receiving this event, `babeltower` updates the thing's config on the registry and send a [`device.config.updated`](#device-config-updated) event.

<details>
  <summary>Headers</summary>

  - `token` **String** user's token

</details>

<details>
  <summary>Payload</summary>

  JSON in the following format:

  - `id` **String** thing's ID
  - `config` **Array** config items, each one formed by:
    - `sensorId` **Number** sensor ID
    - `schema` **JSON Object** schema item, each one formed by:
      - `typeId` **Number** semantic value type (voltage, current, temperature, etc)
      - `unit` **Number** sensor unit (V, A, W, etc)
      - `valueType` **Number** data value type (boolean, integer, etc)
      - `name` **String** sensor name
    - `event` **JSON Object** event item, each one formed by:
      - `change` **Boolean** enable sending sensor data when its value changes
      - `timeSec` **Number** - **Optional** time interval in seconds that indicates when data must be sent to the cloud
      - `lowerThreshold` **(Depends on schema's valueType)** - **Optional** send data to the cloud if it's lower than this threshold
      - `upperThreshold` **(Depends on schema's valueType)** - **Optional** send data to the cloud if it's upper than this threshold

  The semantic specification that defines `valueType`, `unit` and `typeId` properties can be find [here](https://knot-devel.cesar.org.br/doc/thing/unit-type-value.html).

  Example:

  ```json
  {
    "id": "fbe64efa6c7f717e",
    "config": [{
      "sensorId": 1,
      "schema": {
        "typeId": 0xFFF1,
        "unit": 0,
        "valueType": 3,
        "name": "Door lock",
      },
      "event": {
         "change": true,
         "timeSec": 10,
         "lowerThreshold": 1000,
         "upperThreshold": 3000
      }
    }]
  }
  ```
</details>

<details>
  <summary>AMQP Binding</summary>

  - Exchange:
    - Type: direct
    - Name: device
    - Durable: `true`
    - Auto-delete: `false`
  - Routing key: device.config.sent

</details>

### **device.list** <a name="device-list"></a>

Event-command to list the registered things. It follows the request/reply pattern. After obtaining the things, `babeltower` will send a reply message by using the `reply_to` property, which was received in the request header, as reply message's `routing_key`. Because of that, considering the **requestor** has created and sent this `reply_to` in the request, it can also subscribe to receive events that arrive in a queue associated with the `reply_to`. Therefore, the reply is received by the application that has sent the request, in a **one-to-one** manner.

<details>
  <summary>Headers</summary>

  - `token` **String** user's token

</details>


<details>
  <summary>Payload</summary>

  JSON in the following format:

  - Empty object

  Example:

  ```json
  {}
  ```

</details>

<details>
  <summary>AMQP Binding</summary>

  - Exchange:
    - Type: direct
    - Name: device
    - Durable: `true`
    - Auto-delete: `false`
  - Routing key: `device.list`
  - Reply To: <queueName> reply's queue name
  - Correlation Id: <corrID> ID to correlate reply-request after message arrived in the queue

</details>

### **device.auth** <a name="device-auth"></a>

Event-command to verify if a thing is authenticated based on its credentials. It follows the request/reply pattern. After authenticating the device, `babeltower` will send a reply message by using the `reply_to` property, which is received in the request header, as reply message's `routing_key`. Because of that, considering the **requestor** has created and sent this `reply_to` in the request, it can also subscribe to receive events that arrive in a queue associated with the `reply_to`. Therefore, the reply is received by the application that has sent the request, in a **one-to-one** manner.

<details>
  <summary>Headers</summary>

  - `token` **String** user's token

</details>

<details>
  <summary>Payload</summary>

  JSON in the following format:

  - `id` **String** device's ID
  - `token` **String** device's token

  Example:

  ```json
  {
    "id": "fbe64efa6c7f717e",
    "token": "0c20c12e2ac058d0513d81dc58e33b2f9ff8c83d"
  }
  ```
</details>

<details>
  <summary>AMQP Binding</summary>

  - Exchange:
    - Type: direct
    - Name: device
    - Durable: `true`
    - Auto-delete: `false`
  - Routing key: `device.auth`
  - Reply To: <queueName> reply's queue name
  - Correlation Id: <corrID> ID to correlate reply-request after message arrived in the queue

</details>

### **data.sent** <a name="data-sent"></a>

Event that represents a device sending the data gathered from its sensors to the services that are interested. After receiving this event, `babeltower` makes the necessary semantic validation and send a [`data.published`](#data-published) event.

<details>
  <details>
    <summary>Headers</summary>

    - `token` **String** user's token

  </details>

  <summary>Payload</summary>

  JSON in the following format:

  - `id` **String** thing's ID
  - `data` **Array** data items to be published, each one formed by:
    - `sensorId` **Number** sensor ID
    - `value` **Number|Boolean|String** sensor value
    - `timestamp` **String** sensor timestamp

  Example:

  ```json
  {
    "id": "fbe64efa6c7f717e",
    "data": [
      {
        "sensorId": 1,
        "value": false,
        "timestamp": "2022-12-06T10:00:00.0-3000"
      },
      {
        "sensorId": 2,
        "value": 1000,
        "timestamp": "2022-12-06T10:00:00.0-3000"
      }
    ]
  }
  ```
</details>

<details>
  <summary>AMQP Binding</summary>

  - Exchange:
    - Type: fanout
    - Name: data.sent
    - Durable: `true`
    - Auto-delete: `false`

</details>

### **data.request** <a name="data-request"></a>

Event-command to request data from a thing's sensor. After receiving this event, `babeltower` makes the necessary semantic validation and send a [`device.<id>.data.request`](#device-[id]-data-request) event to be routed to the service which control the thing.

<details>
  <summary>Headers</summary>

  - `token` **String** user's token

</details>

<details>
  <summary>Payload</summary>

  JSON in the following format:

  - `id` **String** thing's ID
  - `sensorIds` **Array (Number)** IDs of the sensor to send last value

  Example:

  ```json
  {
    "id": "fbe64efa6c7f717e",
    "data": [1]
  }
  ```
</details>

<details>
  <summary>AMQP Binding</summary>

  - Exchange:
    - Type: direct
    - Name: device
    - Durable: `true`
    - Auto-delete: `false`
  - Routing key: data.request

</details>

### **data.update** <a name="data-update"></a>

Event-command to update a thing's sensor data. After receiving this event, `babeltower` makes the necessary semantic validation and send a [`device.<id>.data.update`](#device-[id]-data-update) event to be routed to the service which control the thing.

<details>
  <summary>Headers</summary>

  - `token` **String** user's token

</details>

<details>
  <summary>Payload</summary>

  JSON in the following format:

  - `id` **String** thing's ID
  - `data` **Array (Object)** updates for sensors/actuators, each one formed by:
    - `sensorId` **Number** ID of the sensor to update
    - `value` **Number|Boolean|String** data to be written
    - `timestamp` **String** sensor timestamp

  Example:

  ```json
  {
    "id": "fbe64efa6c7f717e",
    "data": [{
        "sensorId": 1,
        "value": true,
        "timestamp": "2022-12-06T10:00:00.0-3000"
    }]
  }
  ```
</details>

<details>
  <summary>AMQP Binding</summary>

  - Exchange:
    - Type: direct
    - Name: device
    - Durable: `true`
    - Auto-delete: `false`
  - Routing key: data.update

</details>

## Subscribe

The external consumer applications can subscribe to the events described in this section to receive them and take the appropriate action.

### **device.registered** <a name="device-registered"></a>

Event that represents a thing was registered.

<details>
  <summary>Payload</summary>

  JSON in the following format:

  - `id` **String** thing's ID
  - `token` **String** thing's token
  - `error` **String** described the occurred error

  Success example:

  ```json
  {
    "id": "fbe64efa6c7f717e",
    "token": "5b67ce6bef21701331152d6297e1bd2b22f91787",
    "error": null
  }
  ```

  Error example:

  ```json
  {
    "id": "3aa21010cda96fe9",
    "token": "",
    "error": "device already exists"
  }
  ```

</details>

<details>
  <summary>AMQP Binding</summary>

  - Exchange:
    - Type: direct
    - Name: device
    - Durable: `true`
    - Auto-delete: `false`
  - Routing key: device.registered

</details>

### **device.unregistered** <a name="device-unregistered"></a>

Event that represents a thing was removed.

<details>
  <summary>Payload</summary>

  JSON in the following format:

  - `id` **String** thing's ID
  - `error` **String** described the occurred error

  Success example:

  ```json
  {
    "id": "fbe64efa6c7f717e",
    "error": null
  }
  ```

  Error example:

  ```json
  {
    "id": "3aa21010cda96fe9",
    "error": "forbidden",
  }
  ```
</details>

<details>
  <summary>AMQP Binding</summary>

  - Exchange:
    - Type: direct
    - Name: device
    - Durable: `true`
    - Auto-delete: `false`
  - Routing key: device.unregistered

</details>

### **device.config.updated** <a name="device-config-updated"></a>

Event that represents a thing's config was updated.

<details>
  <summary>Payload</summary>

  JSON in the following format:

  - `id` **String** thing's ID
  - `config` **Array** - **Optional** config items, each one formed by:
    - `sensorId` **Number** sensor ID
    - `schema` **JSON Object** schema item, each one formed by:
      - `typeId` **Number** semantic value type (voltage, current, temperature, etc)
      - `unit` **Number** sensor unit (V, A, W, etc)
      - `valueType` **Number** data value type (boolean, integer, etc)
      - `name` **String** sensor name
    - `event` **JSON Object** - **Optional** event item, each one formed by:
      - `change` **Boolean** enable sending sensor data when its value changes
      - `timeSec` **Number** - **Optional** time interval in seconds that indicates when data must be sent to the cloud
      - `lowerThreshold` **(Depends on schema's valueType)** - **Optional** send data to the cloud if it's lower than this threshold
      - `upperThreshold` **(Depends on schema's valueType)** - **Optional** send data to the cloud if it's upper than this threshold
  - `changed` **Boolean** inform if the update has changed something in the thing's current configuration
  - `error` **String** a string with detailed error message

  The semantic specification that defines `valueType`, `unit` and `typeId` properties can be find [here](https://knot-devel.cesar.org.br/doc/thing/unit-type-value.html).


  Success example:

  ```json
  {
    "id": "fbe64efa6c7f717e",
    "config": [{
      "sensorId": 1,
      "schema": {
        "typeId": 0xFFF1,
        "unit": 0,
        "valueType": 3,
        "name": "Door lock"
      },
      "event": {
         "change": true,
         "timeSec": 10,
         "lowerThreshold": 1000,
         "upperThreshold": 3000
      }
    }],
    "changed": true,
    "error": null
  }
  ```

  Error example:

  ```json
  {
    "id": "3aa21010cda96fe9",
    "error": "invalid config"
  }
  ```
</details>

<details>
  <summary>AMQP Binding</summary>

  - Exchange:
    - Type: direct
    - Name: device
    - Durable: `true`
    - Auto-delete: `false`
  - Routing key: device.config.updated

</details>

### **data.published** <a name="data-published"></a>

Event that represents a data published from a thing's sensor.

<details>
  <summary>Payload</summary>

  JSON in the following format:

  - `id` **String** thing's ID
  - `data` **Array** data items to be published, each one formed by:
    - `sensorId` **Number** sensor ID
    - `value` **Number|Boolean|String** sensor value
    - `timestamp` **String** sensor timestamp

  Example:

  ```json
  {
    "id": "fbe64efa6c7f717e",
    "data": [
      {
        "sensorId": 1,
        "value": false,
        "timestamp": "2022-12-06T10:00:00.0-3000"
      },
      {
        "sensorId": 2,
        "value": 1000,
        "timestamp": "2022-12-06T10:00:00.0-3000"
      }
    ]
  }
  ```
</details>

<details>
  <summary>AMQP Binding</summary>

  - Exchange:
    - Type: fanout
    - Name: data.published
    - Durable: `true`
    - Auto-delete: `false`

</details>

### **data.util** <a name="data-util"></a>
The event represents data published from a thing's sensor segmented based on the user ID associated with the thing. Trusted external clients can specify the user ID associated with the devices they want to capture. The user ID can be adapted to the needs of each project. For example, you can use a value derived from the user token or the token itself. The only requirement is that external clients specify this user ID as the routing key when subscribing to data.util so that they receive data from devices associated with this user.

<details>
  <summary>Payload</summary>

  JSON in the following format:

  - `id` **String** thing's ID
  - `data` **Array** data items to be published, each one formed by:
    - `sensorId` **Number** sensor ID
    - `value` **Number|Boolean|String** sensor value
    - `timestamp` **String** sensor timestamp

  Example:

  ```json
  {
    "id": "fbe64efa6c7f717e",
    "data": [
      {
        "sensorId": 1,
        "value": false,
        "timestamp": "2022-12-06T10:00:00.0-3000"
      },
      {
        "sensorId": 2,
        "value": 1000,
        "timestamp": "2022-12-06T10:00:00.0-3000"
      }
    ]
  }
  ```
</details>

<details>
  <summary>AMQP Binding</summary>

  - Exchange:
    - Type: direct
    - Name: data.util
    - Durable: `true`
    - Auto-delete: `false`
  - Routing Key: `user ID`

</details>

### **data.[sessionId].published** <a name="data-session-published"></a>

Event that represents a data published from a thing's sensor to a user session. You can obtain a `sessionId` by sending a request to the endpoint `POST /sessions` with a valid authorization token. The endpoint specification can be easily viewed in the browser by accessing the address `http://<address>:<port>/swagger/index.html`.

<details>
  <summary>Payload</summary>

  JSON in the following format:

  - `id` **String** thing's ID
  - `data` **Array** data items to be published, each one formed by:
    - `sensorId` **Number** sensor ID
    - `value` **Number|Boolean|String** sensor value
    - `timestamp` **String** sensor timestamp

  Example:

  ```json
  {
    "id": "fbe64efa6c7f717e",
    "data": [
      {
        "sensorId": 1,
        "value": false,
        "timestamp": "2022-12-06T10:00:00.0-3000"
      },
      {
        "sensorId": 2,
        "value": 1000,
        "timestamp": "2022-12-06T10:00:00.0-3000"
      }
    ]
  }
  ```
</details>

<details>
  <summary>AMQP Binding</summary>

  - Exchange:
    - Type: fanout
    - Name: data.[sessionId].published
    - Durable: `true`
    - Auto-delete: `false`

</details>

### **device.[id].data.request** <a name="device-<id>-data-request"></a>

Event-command to request a specific thing's sensor data after validating if the sensor exists in thing's schema and the `value` is in a valid format.

<details>
  <summary>Headers</summary>

  - `token` **String** user's token

</details>

<details>
  <summary>Payload</summary>

  JSON in the following format:

  - `id` **String** thing's ID
  - `data` **Array** data items to be published, each one formed by:
    - `sensorId` **Number** sensor ID
    - `value` **Number|Boolean|String** sensor value
    - `timestamp` **String** sensor timestamp

  Example:

  ```json
  {
    "id": "fbe64efa6c7f717e",
    "data": [
      {
        "sensorId": 1,
        "value": false,
        "timestamp": "2022-12-06T10:00:00.0-3000"
      },
      {
        "sensorId": 2,
        "value": 1000,
        "timestamp": "2022-12-06T10:00:00.0-3000"
      }
    ]
  }
  ```
</details>

<details>
  <summary>AMQP Binding</summary>

  - Exchange:
    - Type: direct
    - Name: device
    - Durable: `true`
    - Auto-delete: `false`
  - Routing Key: `device.<id>.data.request`

</details>

### **device.[id].data.update** <a name="device-<id>-data-update"></a>

Event-command to update a specific thing's sensor data after validating if the `data` is compatible with thing's schema .

<details>
  <summary>Headers</summary>

  - `token` **String** user's token

</details>

<details>
  <summary>Payload</summary>

  JSON in the following format:

  - `id` **String** thing's ID
  - `data` **Array** data items to be published, each one formed by:
    - `sensorId` **Number** sensor ID
    - `value` **Number|Boolean|String** sensor value
    - `timestamp` **String** sensor timestamp

  Example:

  ```json
  {
    "id": "fbe64efa6c7f717e",
    "data": [
      {
        "sensorId": 1,
        "value": false,
        "timestamp": "2022-12-06T10:00:00.0-3000"
      },
      {
        "sensorId": 2,
        "value": 1000,
        "timestamp": "2022-12-06T10:00:00.0-3000"
      }
    ]
  }
  ```
</details>

<details>
  <summary>AMQP Binding</summary>

  - Exchange:
    - Type: direct
    - Name: device
    - Durable: `true`
    - Auto-delete: `false`
  - Routing Key: `device.<id>.data.update`

</details>