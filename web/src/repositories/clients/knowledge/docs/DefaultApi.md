# GraphRagServiceKnowledge.DefaultApi

All URIs are relative to *https://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**createDocument**](DefaultApi.md#createDocument) | **POST** /documents | Create a markdown document and start indexing (requires knowledge:write)
[**deleteDocument**](DefaultApi.md#deleteDocument) | **DELETE** /documents/{documentId} | Delete a document and its graph data (requires knowledge:write)
[**getDocument**](DefaultApi.md#getDocument) | **GET** /documents/{documentId} | Get a document by ID (requires knowledge:write)
[**getDocumentIndexStatus**](DefaultApi.md#getDocumentIndexStatus) | **GET** /documents/{documentId}/index-status | Get indexing status for a document (requires knowledge:write)
[**listDocuments**](DefaultApi.md#listDocuments) | **GET** /documents | List knowledge documents (requires knowledge:write)
[**reindexDocument**](DefaultApi.md#reindexDocument) | **POST** /documents/{documentId}/reindex | Trigger reindexing for a document (requires knowledge:write)
[**retrieve**](DefaultApi.md#retrieve) | **POST** /retrieve | Retrieve relevant chunks from the knowledge graph (requires conversation:ask)
[**updateDocument**](DefaultApi.md#updateDocument) | **PUT** /documents/{documentId} | Update a document and reindex (requires knowledge:write)



## createDocument

> Document createDocument(createDocumentRequest)

Create a markdown document and start indexing (requires knowledge:write)

### Example

```javascript
import GraphRagServiceKnowledge from 'graph_rag_service_knowledge';
let defaultClient = GraphRagServiceKnowledge.ApiClient.instance;
// Configure Bearer (JWT) access token for authorization: bearerAuth
let bearerAuth = defaultClient.authentications['bearerAuth'];
bearerAuth.accessToken = "YOUR ACCESS TOKEN"

let apiInstance = new GraphRagServiceKnowledge.DefaultApi();
let createDocumentRequest = new GraphRagServiceKnowledge.CreateDocumentRequest(); // CreateDocumentRequest | 
apiInstance.createDocument(createDocumentRequest, (error, data, response) => {
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
 **createDocumentRequest** | [**CreateDocumentRequest**](CreateDocumentRequest.md)|  | 

### Return type

[**Document**](Document.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


## deleteDocument

> deleteDocument(documentId)

Delete a document and its graph data (requires knowledge:write)

### Example

```javascript
import GraphRagServiceKnowledge from 'graph_rag_service_knowledge';
let defaultClient = GraphRagServiceKnowledge.ApiClient.instance;
// Configure Bearer (JWT) access token for authorization: bearerAuth
let bearerAuth = defaultClient.authentications['bearerAuth'];
bearerAuth.accessToken = "YOUR ACCESS TOKEN"

let apiInstance = new GraphRagServiceKnowledge.DefaultApi();
let documentId = null; // String | 
apiInstance.deleteDocument(documentId, (error, data, response) => {
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
 **documentId** | [**String**](.md)|  | 

### Return type

null (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## getDocument

> Document getDocument(documentId)

Get a document by ID (requires knowledge:write)

### Example

```javascript
import GraphRagServiceKnowledge from 'graph_rag_service_knowledge';
let defaultClient = GraphRagServiceKnowledge.ApiClient.instance;
// Configure Bearer (JWT) access token for authorization: bearerAuth
let bearerAuth = defaultClient.authentications['bearerAuth'];
bearerAuth.accessToken = "YOUR ACCESS TOKEN"

let apiInstance = new GraphRagServiceKnowledge.DefaultApi();
let documentId = null; // String | 
apiInstance.getDocument(documentId, (error, data, response) => {
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
 **documentId** | [**String**](.md)|  | 

### Return type

[**Document**](Document.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## getDocumentIndexStatus

> IndexJob getDocumentIndexStatus(documentId)

Get indexing status for a document (requires knowledge:write)

### Example

```javascript
import GraphRagServiceKnowledge from 'graph_rag_service_knowledge';
let defaultClient = GraphRagServiceKnowledge.ApiClient.instance;
// Configure Bearer (JWT) access token for authorization: bearerAuth
let bearerAuth = defaultClient.authentications['bearerAuth'];
bearerAuth.accessToken = "YOUR ACCESS TOKEN"

let apiInstance = new GraphRagServiceKnowledge.DefaultApi();
let documentId = null; // String | 
apiInstance.getDocumentIndexStatus(documentId, (error, data, response) => {
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
 **documentId** | [**String**](.md)|  | 

### Return type

[**IndexJob**](IndexJob.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## listDocuments

> DocumentList listDocuments(opts)

List knowledge documents (requires knowledge:write)

### Example

```javascript
import GraphRagServiceKnowledge from 'graph_rag_service_knowledge';
let defaultClient = GraphRagServiceKnowledge.ApiClient.instance;
// Configure Bearer (JWT) access token for authorization: bearerAuth
let bearerAuth = defaultClient.authentications['bearerAuth'];
bearerAuth.accessToken = "YOUR ACCESS TOKEN"

let apiInstance = new GraphRagServiceKnowledge.DefaultApi();
let opts = {
  'limit': 20, // Number | 
  'offset': 0, // Number | 
  'tag': "tag_example", // String | Filter by tag
  'indexStatus': new GraphRagServiceKnowledge.IndexStatus() // IndexStatus | 
};
apiInstance.listDocuments(opts, (error, data, response) => {
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
 **tag** | **String**| Filter by tag | [optional] 
 **indexStatus** | [**IndexStatus**](.md)|  | [optional] 

### Return type

[**DocumentList**](DocumentList.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## reindexDocument

> IndexJob reindexDocument(documentId)

Trigger reindexing for a document (requires knowledge:write)

### Example

```javascript
import GraphRagServiceKnowledge from 'graph_rag_service_knowledge';
let defaultClient = GraphRagServiceKnowledge.ApiClient.instance;
// Configure Bearer (JWT) access token for authorization: bearerAuth
let bearerAuth = defaultClient.authentications['bearerAuth'];
bearerAuth.accessToken = "YOUR ACCESS TOKEN"

let apiInstance = new GraphRagServiceKnowledge.DefaultApi();
let documentId = null; // String | 
apiInstance.reindexDocument(documentId, (error, data, response) => {
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
 **documentId** | [**String**](.md)|  | 

### Return type

[**IndexJob**](IndexJob.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## retrieve

> RetrieveResponse retrieve(retrieveRequest)

Retrieve relevant chunks from the knowledge graph (requires conversation:ask)

Called by the Conversation service during message generation. Returns structured context, not an LLM answer.

### Example

```javascript
import GraphRagServiceKnowledge from 'graph_rag_service_knowledge';
let defaultClient = GraphRagServiceKnowledge.ApiClient.instance;
// Configure Bearer (JWT) access token for authorization: bearerAuth
let bearerAuth = defaultClient.authentications['bearerAuth'];
bearerAuth.accessToken = "YOUR ACCESS TOKEN"

let apiInstance = new GraphRagServiceKnowledge.DefaultApi();
let retrieveRequest = new GraphRagServiceKnowledge.RetrieveRequest(); // RetrieveRequest | 
apiInstance.retrieve(retrieveRequest, (error, data, response) => {
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
 **retrieveRequest** | [**RetrieveRequest**](RetrieveRequest.md)|  | 

### Return type

[**RetrieveResponse**](RetrieveResponse.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


## updateDocument

> Document updateDocument(documentId, updateDocumentRequest)

Update a document and reindex (requires knowledge:write)

### Example

```javascript
import GraphRagServiceKnowledge from 'graph_rag_service_knowledge';
let defaultClient = GraphRagServiceKnowledge.ApiClient.instance;
// Configure Bearer (JWT) access token for authorization: bearerAuth
let bearerAuth = defaultClient.authentications['bearerAuth'];
bearerAuth.accessToken = "YOUR ACCESS TOKEN"

let apiInstance = new GraphRagServiceKnowledge.DefaultApi();
let documentId = null; // String | 
let updateDocumentRequest = new GraphRagServiceKnowledge.UpdateDocumentRequest(); // UpdateDocumentRequest | 
apiInstance.updateDocument(documentId, updateDocumentRequest, (error, data, response) => {
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
 **documentId** | [**String**](.md)|  | 
 **updateDocumentRequest** | [**UpdateDocumentRequest**](UpdateDocumentRequest.md)|  | 

### Return type

[**Document**](Document.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

