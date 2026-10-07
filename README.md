Knowledge:


Graph:

Document
 -[TAGGED_AS]-> Tag
 -[HAS_CHUNK]-> Chunk
Chunk
 -[MENTIONS]-> Entity


Job Work:
1. extract chunks + entities from pending documents
2. delete orphaned Tag/Chunk/Entity nodes (zero in-degree non-Document nodes)


Flow:


Admin                         Knowledge                         Memgraph
  |                               |                                 |
  |-- POST /documents ----------->|                                 |
  |                               |-- MERGE Document+Tags --------->|
  |<-- Document pending ----------|                                 |
  |                               |                                 |
  |                               |-- parse / chunk / embed         |
  |                               |-- extract entities              |
  |                               |-- MERGE Chunk/Entity ---------->|
  |                               |-- SET index_status=ready ------>|
  |                               |                                 |
User                          Conversation                      Knowledge
  |-- POST /messages ------------>|                                 |
  |                               |-- POST /retrieve -------------->|
  |                               |                                 |-- vector search
  |                               |                                 |-- expand hops
  |                               |<-- chunks + graphPath ----------|
  |                               |-- prompt + LLM                  |
  |<-- assistant + sources -------|                                 |