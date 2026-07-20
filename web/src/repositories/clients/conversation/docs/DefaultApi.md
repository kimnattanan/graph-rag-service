# GraphRagServiceConversation.DefaultApi

All URIs are relative to *https://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**createConversation**](DefaultApi.md#createConversation) | **POST** /conversations | Create an empty conversation (requires conversation:ask)
[**deleteConversation**](DefaultApi.md#deleteConversation) | **DELETE** /conversations/{conversationId} | Delete a conversation and its messages (requires conversation:ask)
[**getConversation**](DefaultApi.md#getConversation) | **GET** /conversations/{conversationId} | Get a conversation with its messages (requires conversation:ask)
[**getMessage**](DefaultApi.md#getMessage) | **GET** /conversations/{conversationId}/messages/{messageId} | Get a single message (requires conversation:ask)
[**listConversations**](DefaultApi.md#listConversations) | **GET** /conversations | List the authenticated user&#39;s conversations (requires conversation:ask)
[**listMessages**](DefaultApi.md#listMessages) | **GET** /conversations/{conversationId}/messages | List messages in a conversation (requires conversation:ask)
[**sendMessage**](DefaultApi.md#sendMessage) | **POST** /conversations/{conversationId}/messages | Send a user message and receive an assistant reply (requires conversation:ask)



## createConversation

> Conversation createConversation(opts)

Create an empty conversation (requires conversation:ask)

### Example

```javascript
import GraphRagServiceConversation from 'graph_rag_service_conversation';
let defaultClient = GraphRagServiceConversation.ApiClient.instance;
// Configure Bearer (JWT) access token for authorization: bearerAuth
let bearerAuth = defaultClient.authentications['bearerAuth'];
bearerAuth.accessToken = "YOUR ACCESS TOKEN"

let apiInstance = new GraphRagServiceConversation.DefaultApi();
let opts = {
  'createConversationRequest': new GraphRagServiceConversation.CreateConversationRequest() // CreateConversationRequest | 
};
apiInstance.createConversation(opts, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **createConversationRequest** | [**CreateConversationRequest**](CreateConversationRequest.md)|  | [optional] 

### Return type

[**Conversation**](Conversation.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


## deleteConversation

> deleteConversation(conversationId)

Delete a conversation and its messages (requires conversation:ask)

### Example

```javascript
import GraphRagServiceConversation from 'graph_rag_service_conversation';
let defaultClient = GraphRagServiceConversation.ApiClient.instance;
// Configure Bearer (JWT) access token for authorization: bearerAuth
let bearerAuth = defaultClient.authentications['bearerAuth'];
bearerAuth.accessToken = "YOUR ACCESS TOKEN"

let apiInstance = new GraphRagServiceConversation.DefaultApi();
let conversationId = null; // String | 
apiInstance.deleteConversation(conversationId, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **conversationId** | [**String**](.md)|  | 

### Return type

null (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## getConversation

> ConversationDetail getConversation(conversationId)

Get a conversation with its messages (requires conversation:ask)

### Example

```javascript
import GraphRagServiceConversation from 'graph_rag_service_conversation';
let defaultClient = GraphRagServiceConversation.ApiClient.instance;
// Configure Bearer (JWT) access token for authorization: bearerAuth
let bearerAuth = defaultClient.authentications['bearerAuth'];
bearerAuth.accessToken = "YOUR ACCESS TOKEN"

let apiInstance = new GraphRagServiceConversation.DefaultApi();
let conversationId = null; // String | 
apiInstance.getConversation(conversationId, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **conversationId** | [**String**](.md)|  | 

### Return type

[**ConversationDetail**](ConversationDetail.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## getMessage

> Message getMessage(conversationId, messageId)

Get a single message (requires conversation:ask)

### Example

```javascript
import GraphRagServiceConversation from 'graph_rag_service_conversation';
let defaultClient = GraphRagServiceConversation.ApiClient.instance;
// Configure Bearer (JWT) access token for authorization: bearerAuth
let bearerAuth = defaultClient.authentications['bearerAuth'];
bearerAuth.accessToken = "YOUR ACCESS TOKEN"

let apiInstance = new GraphRagServiceConversation.DefaultApi();
let conversationId = null; // String | 
let messageId = null; // String | 
apiInstance.getMessage(conversationId, messageId, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **conversationId** | [**String**](.md)|  | 
 **messageId** | [**String**](.md)|  | 

### Return type

[**Message**](Message.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## listConversations

> ConversationList listConversations(opts)

List the authenticated user&#39;s conversations (requires conversation:ask)

### Example

```javascript
import GraphRagServiceConversation from 'graph_rag_service_conversation';
let defaultClient = GraphRagServiceConversation.ApiClient.instance;
// Configure Bearer (JWT) access token for authorization: bearerAuth
let bearerAuth = defaultClient.authentications['bearerAuth'];
bearerAuth.accessToken = "YOUR ACCESS TOKEN"

let apiInstance = new GraphRagServiceConversation.DefaultApi();
let opts = {
  'limit': 20, // Number | 
  'offset': 0 // Number | 
};
apiInstance.listConversations(opts, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **limit** | **Number**|  | [optional] [default to 20]
 **offset** | **Number**|  | [optional] [default to 0]

### Return type

[**ConversationList**](ConversationList.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## listMessages

> MessageList listMessages(conversationId, opts)

List messages in a conversation (requires conversation:ask)

### Example

```javascript
import GraphRagServiceConversation from 'graph_rag_service_conversation';
let defaultClient = GraphRagServiceConversation.ApiClient.instance;
// Configure Bearer (JWT) access token for authorization: bearerAuth
let bearerAuth = defaultClient.authentications['bearerAuth'];
bearerAuth.accessToken = "YOUR ACCESS TOKEN"

let apiInstance = new GraphRagServiceConversation.DefaultApi();
let conversationId = null; // String | 
let opts = {
  'limit': 50, // Number | 
  'offset': 0 // Number | 
};
apiInstance.listMessages(conversationId, opts, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **conversationId** | [**String**](.md)|  | 
 **limit** | **Number**|  | [optional] [default to 50]
 **offset** | **Number**|  | [optional] [default to 0]

### Return type

[**MessageList**](MessageList.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## sendMessage

> SendMessageResponse sendMessage(conversationId, sendMessageRequest)

Send a user message and receive an assistant reply (requires conversation:ask)

Stores the user message, retrieves context from Knowledge, calls the LLM, stores the assistant message, and returns both.

### Example

```javascript
import GraphRagServiceConversation from 'graph_rag_service_conversation';
let defaultClient = GraphRagServiceConversation.ApiClient.instance;
// Configure Bearer (JWT) access token for authorization: bearerAuth
let bearerAuth = defaultClient.authentications['bearerAuth'];
bearerAuth.accessToken = "YOUR ACCESS TOKEN"

let apiInstance = new GraphRagServiceConversation.DefaultApi();
let conversationId = null; // String | 
let sendMessageRequest = new GraphRagServiceConversation.SendMessageRequest(); // SendMessageRequest | 
apiInstance.sendMessage(conversationId, sendMessageRequest, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **conversationId** | [**String**](.md)|  | 
 **sendMessageRequest** | [**SendMessageRequest**](SendMessageRequest.md)|  | 

### Return type

[**SendMessageResponse**](SendMessageResponse.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

