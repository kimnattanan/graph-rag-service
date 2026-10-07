# GraphRagServiceUser.User

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **String** |  | 
**email** | **String** |  | 
**username** | **String** |  | 
**role** | **String** |  | 
**permissions** | **[String]** | Derived from role. Present in the JWT for other services. | 
**createdAt** | **Date** |  | 
**updatedAt** | **Date** |  | 



## Enum: RoleEnum


* `admin` (value: `"admin"`)

* `user` (value: `"user"`)





## Enum: [PermissionsEnum]


* `knowledge:write` (value: `"knowledge:write"`)

* `conversation:ask` (value: `"conversation:ask"`)




