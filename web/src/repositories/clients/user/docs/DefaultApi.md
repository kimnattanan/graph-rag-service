# GraphRagServiceUser.DefaultApi

All URIs are relative to *https://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**createUser**](DefaultApi.md#createUser) | **POST** /users | Create a user (admin only)
[**deleteUser**](DefaultApi.md#deleteUser) | **DELETE** /users/{userId} | Delete a user (admin only)
[**getCurrentUser**](DefaultApi.md#getCurrentUser) | **GET** /users/me | Get the authenticated user profile
[**getCurrentUserPermissions**](DefaultApi.md#getCurrentUserPermissions) | **GET** /users/me/permissions | Get permissions for the authenticated user
[**getUser**](DefaultApi.md#getUser) | **GET** /users/{userId} | Get a user by ID (admin only)
[**listUsers**](DefaultApi.md#listUsers) | **GET** /users | List users (admin only)
[**login**](DefaultApi.md#login) | **POST** /auth/login | Authenticate and receive access tokens
[**refreshToken**](DefaultApi.md#refreshToken) | **POST** /auth/refresh | Exchange a refresh token for a new access token
[**register**](DefaultApi.md#register) | **POST** /auth/register | Register a new user account
[**updateUser**](DefaultApi.md#updateUser) | **PUT** /users/{userId} | Update a user (admin only)



## createUser

> User createUser(createUserRequest)

Create a user (admin only)

### Example

```javascript
import GraphRagServiceUser from 'graph_rag_service_user';
let defaultClient = GraphRagServiceUser.ApiClient.instance;
// Configure Bearer (JWT) access token for authorization: bearerAuth
let bearerAuth = defaultClient.authentications['bearerAuth'];
bearerAuth.accessToken = "YOUR ACCESS TOKEN"

let apiInstance = new GraphRagServiceUser.DefaultApi();
let createUserRequest = new GraphRagServiceUser.CreateUserRequest(); // CreateUserRequest | 
apiInstance.createUser(createUserRequest, (error, data, response) => {
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
 **createUserRequest** | [**CreateUserRequest**](CreateUserRequest.md)|  | 

### Return type

[**User**](User.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


## deleteUser

> deleteUser(userId)

Delete a user (admin only)

### Example

```javascript
import GraphRagServiceUser from 'graph_rag_service_user';
let defaultClient = GraphRagServiceUser.ApiClient.instance;
// Configure Bearer (JWT) access token for authorization: bearerAuth
let bearerAuth = defaultClient.authentications['bearerAuth'];
bearerAuth.accessToken = "YOUR ACCESS TOKEN"

let apiInstance = new GraphRagServiceUser.DefaultApi();
let userId = null; // String | 
apiInstance.deleteUser(userId, (error, data, response) => {
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
 **userId** | [**String**](.md)|  | 

### Return type

null (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## getCurrentUser

> User getCurrentUser()

Get the authenticated user profile

### Example

```javascript
import GraphRagServiceUser from 'graph_rag_service_user';
let defaultClient = GraphRagServiceUser.ApiClient.instance;
// Configure Bearer (JWT) access token for authorization: bearerAuth
let bearerAuth = defaultClient.authentications['bearerAuth'];
bearerAuth.accessToken = "YOUR ACCESS TOKEN"

let apiInstance = new GraphRagServiceUser.DefaultApi();
apiInstance.getCurrentUser((error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**User**](User.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## getCurrentUserPermissions

> Permissions getCurrentUserPermissions()

Get permissions for the authenticated user

### Example

```javascript
import GraphRagServiceUser from 'graph_rag_service_user';
let defaultClient = GraphRagServiceUser.ApiClient.instance;
// Configure Bearer (JWT) access token for authorization: bearerAuth
let bearerAuth = defaultClient.authentications['bearerAuth'];
bearerAuth.accessToken = "YOUR ACCESS TOKEN"

let apiInstance = new GraphRagServiceUser.DefaultApi();
apiInstance.getCurrentUserPermissions((error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**Permissions**](Permissions.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## getUser

> User getUser(userId)

Get a user by ID (admin only)

### Example

```javascript
import GraphRagServiceUser from 'graph_rag_service_user';
let defaultClient = GraphRagServiceUser.ApiClient.instance;
// Configure Bearer (JWT) access token for authorization: bearerAuth
let bearerAuth = defaultClient.authentications['bearerAuth'];
bearerAuth.accessToken = "YOUR ACCESS TOKEN"

let apiInstance = new GraphRagServiceUser.DefaultApi();
let userId = null; // String | 
apiInstance.getUser(userId, (error, data, response) => {
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
 **userId** | [**String**](.md)|  | 

### Return type

[**User**](User.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## listUsers

> UserList listUsers(opts)

List users (admin only)

### Example

```javascript
import GraphRagServiceUser from 'graph_rag_service_user';
let defaultClient = GraphRagServiceUser.ApiClient.instance;
// Configure Bearer (JWT) access token for authorization: bearerAuth
let bearerAuth = defaultClient.authentications['bearerAuth'];
bearerAuth.accessToken = "YOUR ACCESS TOKEN"

let apiInstance = new GraphRagServiceUser.DefaultApi();
let opts = {
  'limit': 20, // Number | 
  'offset': 0 // Number | 
};
apiInstance.listUsers(opts, (error, data, response) => {
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

[**UserList**](UserList.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## login

> AuthResult login(loginRequest)

Authenticate and receive access tokens

### Example

```javascript
import GraphRagServiceUser from 'graph_rag_service_user';

let apiInstance = new GraphRagServiceUser.DefaultApi();
let loginRequest = new GraphRagServiceUser.LoginRequest(); // LoginRequest | 
apiInstance.login(loginRequest, (error, data, response) => {
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
 **loginRequest** | [**LoginRequest**](LoginRequest.md)|  | 

### Return type

[**AuthResult**](AuthResult.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


## refreshToken

> AuthResult refreshToken(refreshRequest)

Exchange a refresh token for a new access token

### Example

```javascript
import GraphRagServiceUser from 'graph_rag_service_user';

let apiInstance = new GraphRagServiceUser.DefaultApi();
let refreshRequest = new GraphRagServiceUser.RefreshRequest(); // RefreshRequest | 
apiInstance.refreshToken(refreshRequest, (error, data, response) => {
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
 **refreshRequest** | [**RefreshRequest**](RefreshRequest.md)|  | 

### Return type

[**AuthResult**](AuthResult.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


## register

> AuthResult register(registerRequest)

Register a new user account

### Example

```javascript
import GraphRagServiceUser from 'graph_rag_service_user';

let apiInstance = new GraphRagServiceUser.DefaultApi();
let registerRequest = new GraphRagServiceUser.RegisterRequest(); // RegisterRequest | 
apiInstance.register(registerRequest, (error, data, response) => {
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
 **registerRequest** | [**RegisterRequest**](RegisterRequest.md)|  | 

### Return type

[**AuthResult**](AuthResult.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


## updateUser

> User updateUser(userId, updateUserRequest)

Update a user (admin only)

### Example

```javascript
import GraphRagServiceUser from 'graph_rag_service_user';
let defaultClient = GraphRagServiceUser.ApiClient.instance;
// Configure Bearer (JWT) access token for authorization: bearerAuth
let bearerAuth = defaultClient.authentications['bearerAuth'];
bearerAuth.accessToken = "YOUR ACCESS TOKEN"

let apiInstance = new GraphRagServiceUser.DefaultApi();
let userId = null; // String | 
let updateUserRequest = new GraphRagServiceUser.UpdateUserRequest(); // UpdateUserRequest | 
apiInstance.updateUser(userId, updateUserRequest, (error, data, response) => {
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
 **userId** | [**String**](.md)|  | 
 **updateUserRequest** | [**UpdateUserRequest**](UpdateUserRequest.md)|  | 

### Return type

[**User**](User.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

