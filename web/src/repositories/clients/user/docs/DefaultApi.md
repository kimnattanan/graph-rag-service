# GraphRagServiceUser.DefaultApi

All URIs are relative to *https://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**deleteAccount**](DefaultApi.md#deleteAccount) | **DELETE** /users/me | Delete the authenticated user account
[**getCurrentUser**](DefaultApi.md#getCurrentUser) | **GET** /users/me | Get the authenticated user profile
[**login**](DefaultApi.md#login) | **POST** /auth/login | Authenticate and receive an access token
[**logout**](DefaultApi.md#logout) | **POST** /auth/logout | Log out the current session
[**register**](DefaultApi.md#register) | **POST** /auth/register | Register a new user account



## deleteAccount

> deleteAccount()

Delete the authenticated user account

Command. Deletes the account and its sessions.

### Example

```javascript
import GraphRagServiceUser from 'graph_rag_service_user';
let defaultClient = GraphRagServiceUser.ApiClient.instance;
// Configure Bearer (JWT) access token for authorization: bearerAuth
let bearerAuth = defaultClient.authentications['bearerAuth'];
bearerAuth.accessToken = "YOUR ACCESS TOKEN"

let apiInstance = new GraphRagServiceUser.DefaultApi();
apiInstance.deleteAccount((error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
});
```

### Parameters

This endpoint does not need any parameter.

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


## login

> AuthResult login(loginRequest)

Authenticate and receive an access token

Validates credentials, opens a session, and returns a JWT.

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


## logout

> logout()

Log out the current session

Command. Revokes the session embedded in the access token.

### Example

```javascript
import GraphRagServiceUser from 'graph_rag_service_user';
let defaultClient = GraphRagServiceUser.ApiClient.instance;
// Configure Bearer (JWT) access token for authorization: bearerAuth
let bearerAuth = defaultClient.authentications['bearerAuth'];
bearerAuth.accessToken = "YOUR ACCESS TOKEN"

let apiInstance = new GraphRagServiceUser.DefaultApi();
apiInstance.logout((error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
});
```

### Parameters

This endpoint does not need any parameter.

### Return type

null (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## register

> register(registerRequest)

Register a new user account

Command. Creates an account with role user. Returns no body.

### Example

```javascript
import GraphRagServiceUser from 'graph_rag_service_user';

let apiInstance = new GraphRagServiceUser.DefaultApi();
let registerRequest = new GraphRagServiceUser.RegisterRequest(); // RegisterRequest | 
apiInstance.register(registerRequest, (error, data, response) => {
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
 **registerRequest** | [**RegisterRequest**](RegisterRequest.md)|  | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

