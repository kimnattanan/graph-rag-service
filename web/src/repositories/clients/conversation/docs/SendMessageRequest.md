# GraphRagServiceConversation.SendMessageRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **String** | Client-generated id for the user message | 
**content** | **String** |  | 
**topK** | **Number** | Number of knowledge chunks to retrieve | [optional] [default to 5]
**tags** | **[String]** | Optional tag filter passed to Knowledge retrieval | [optional] 
**historyCapacity** | **Number** | Number of most recent messages in the conversation to include as LLM context | [optional] [default to 3]


