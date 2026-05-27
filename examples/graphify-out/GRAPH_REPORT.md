# Graph Report - .  (2026-05-25)

## Corpus Check
- Corpus is ~16,894 words - fits in a single context window. You may not need a graph.

## Summary
- 334 nodes · 432 edges · 33 communities detected
- Extraction: 99% EXTRACTED · 1% INFERRED · 0% AMBIGUOUS · INFERRED: 4 edges (avg confidence: 0.75)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_Core Framework Types & APIs|Core Framework Types & APIs]]
- [[_COMMUNITY_NetHTTP Adapter|NetHTTP Adapter]]
- [[_COMMUNITY_Gin+Hertz Adapters & Server Factory|Gin+Hertz Adapters & Server Factory]]
- [[_COMMUNITY_Router Implementation|Router Implementation]]
- [[_COMMUNITY_Example Applications (main)|Example Applications (main)]]
- [[_COMMUNITY_Gin Adapter Tests|Gin Adapter Tests]]
- [[_COMMUNITY_Hertz Router|Hertz Router]]
- [[_COMMUNITY_Hertz Adapter Tests|Hertz Adapter Tests]]
- [[_COMMUNITY_Gin Router|Gin Router]]
- [[_COMMUNITY_App Type|App Type]]
- [[_COMMUNITY_Config-in-Constructor Design|Config-in-Constructor Design]]
- [[_COMMUNITY_Gin Server|Gin Server]]
- [[_COMMUNITY_NetHTTP Server|NetHTTP Server]]
- [[_COMMUNITY_HTTPX Router Examples|HTTPX Router Examples]]
- [[_COMMUNITY_Core Type Definitions|Core Type Definitions]]
- [[_COMMUNITY_Hertz Server|Hertz Server]]
- [[_COMMUNITY_Config & WebSocket Test|Config & WebSocket Test]]
- [[_COMMUNITY_Server Interfaces|Server Interfaces]]
- [[_COMMUNITY_WebSocket Client|WebSocket Client]]
- [[_COMMUNITY_Adapter Registry|Adapter Registry]]
- [[_COMMUNITY_WebSocket Upgrade|WebSocket Upgrade]]
- [[_COMMUNITY_Gin Router Test|Gin Router Test]]
- [[_COMMUNITY_Hertz Router Test|Hertz Router Test]]
- [[_COMMUNITY_External Server Test|External Server Test]]
- [[_COMMUNITY_Route Type|Route Type]]
- [[_COMMUNITY_Upgrade Tests|Upgrade Tests]]
- [[_COMMUNITY_NetHTTP Server Tests|NetHTTP Server Tests]]
- [[_COMMUNITY_App Tests|App Tests]]
- [[_COMMUNITY_wx README|wx README]]
- [[_COMMUNITY_wx WebSocket Tasks|wx WebSocket Tasks]]
- [[_COMMUNITY_ws-client Spec|ws-client Spec]]
- [[_COMMUNITY_Hertz Example|Hertz Example]]
- [[_COMMUNITY_NetHTTP Example|NetHTTP Example]]

## God Nodes (most connected - your core abstractions)
1. `main()` - 15 edges
2. `App` - 11 edges
3. `routerImpl` - 11 edges
4. `GinServer` - 11 edges
5. `netHttpServer` - 10 edges
6. `groupRouter` - 9 edges
7. `hertzServer` - 9 edges
8. `hertzRouterGroup` - 9 edges
9. `netHttpRouterGroup` - 9 edges
10. `ginRouter` - 9 edges

## Surprising Connections (you probably didn't know these)
- `httpx - Unified HTTP Framework Adapter Layer` --conceptually_related_to--> `Graceful Shutdown`  [INFERRED]
  README.md → docs/superpowers/specs/2026-05-18-http-library-design.md
- `App Type` --conceptually_related_to--> `Graceful Shutdown`  [INFERRED]
  README.md → docs/superpowers/specs/2026-05-18-http-library-design.md
- `Demo Example` --references--> `App Type`  [EXTRACTED]
  examples/README.md → README.md
- `WebSocket Support` --references--> `Server Interface`  [EXTRACTED]
  docs/superpowers/specs/2026-05-18-http-library-design.md → README.md
- `HTTP Library Implementation Plan` --implements--> `HandlerFunc Type`  [EXTRACTED]
  docs/superpowers/plans/2026-05-18-http-library-implementation.md → README.md

## Hyperedges (group relationships)
- **config-in-constructor Design Decisions** — New_function, LoadConfig, getAdapter, App_struct, Config_struct, AdapterFactory [EXTRACTED 1.00]
- **WebSocket Client Implementation** — client_go, gorilla_websocket, Message_struct, ws_endpoint [EXTRACTED 1.00]
- **WebSocket Test Infrastructure** — websocket_test, external_server_test, os_exec, viper, config_yaml, dynamic_port_loading_capability, random_message_generation_capability [EXTRACTED 0.85]
- **Core Type Hierarchy for HTTP Framework** — app_type, router_interface, route_type, handlerfunc_type, handlercontext_interface, middlewarefunc_type [EXTRACTED 1.00]
- **Adapter Implementations** — server_interface, adapter_factory, adapter_pattern [EXTRACTED 1.00]
- **Configuration-Driven Startup Pattern** — app_type, config_type, adapter_registry, adapter_factory, graceful_shutdown [EXTRACTED 0.85]

## Communities

### Community 0 - "Core Framework Types & APIs"
Cohesion: 0.08
Nodes (32): AdapterFactory, Adapter Pattern, Adapter Registry, App Type, Basic Example (Gin), Config Type, Demo Example, End-to-End Tests (+24 more)

### Community 1 - "NetHTTP Adapter"
Cohesion: 0.1
Nodes (4): netHttpHandlerContext, netHttpRouter, netHttpRouterGroup, ServerOption

### Community 2 - "Gin+Hertz Adapters & Server Factory"
Cohesion: 0.09
Nodes (10): ginHandlerContext, ServerOption, hertzHandlerContext, ServerOption, init(), loadTLSConfig(), newGinRouter(), NewServer() (+2 more)

### Community 3 - "Router Implementation"
Cohesion: 0.14
Nodes (5): groupRouter, Router, routeRecord, RouterGroup, routerImpl

### Community 4 - "Example Applications (main)"
Cohesion: 0.18
Nodes (13): Message, basicServer(), gracefulShutdownExample(), handlerContextExample(), loadConfig(), main(), middlewareExample(), newChatRoom() (+5 more)

### Community 5 - "Gin Adapter Tests"
Cohesion: 0.11
Nodes (0): 

### Community 6 - "Hertz Router"
Cohesion: 0.21
Nodes (2): hertzRouter, hertzRouterGroup

### Community 7 - "Hertz Adapter Tests"
Cohesion: 0.12
Nodes (0): 

### Community 8 - "Gin Router"
Cohesion: 0.26
Nodes (2): ginRouter, ginRouterGroup

### Community 9 - "App Type"
Cohesion: 0.18
Nodes (1): App

### Community 10 - "Config-in-Constructor Design"
Cohesion: 0.21
Nodes (13): AdapterFactory Interface, App Structure, Config Struct, Hertz Adapter, LoadConfig Function, New() Constructor Function, Run() Method, config-in-constructor Capability (+5 more)

### Community 11 - "Gin Server"
Cohesion: 0.21
Nodes (2): GinServer, toGinMiddleware()

### Community 12 - "NetHTTP Server"
Cohesion: 0.27
Nodes (1): netHttpServer

### Community 13 - "HTTPX Router Examples"
Cohesion: 0.29
Nodes (7): Article, User, setupAdminRoutes(), setupAPIV1Routes(), setupArticlesRoutes(), setupRoutes(), setupUsersRoutes()

### Community 14 - "Core Type Definitions"
Cohesion: 0.22
Nodes (6): HandlerContext, HandlerFunc, MiddlewareFunc, StartOption, TLSConfig, WSConfig

### Community 15 - "Hertz Server"
Cohesion: 0.28
Nodes (1): hertzServer

### Community 16 - "Config & WebSocket Test"
Cohesion: 0.29
Nodes (6): Config, config.yaml, dynamic-port-loading Capability, random-message-generation Capability, viper Configuration Library, websocket_test.go

### Community 17 - "Server Interfaces"
Cohesion: 0.33
Nodes (3): GracefulServer, Option, Server

### Community 18 - "WebSocket Client"
Cohesion: 0.33
Nodes (6): Message JSON Structure, WebSocket Client Program, gorilla/websocket Library, ws-client Capability, WebSocket /ws Endpoint, wx-websocket-client Change Proposal

### Community 19 - "Adapter Registry"
Cohesion: 0.5
Nodes (1): AdapterFactory

### Community 20 - "WebSocket Upgrade"
Cohesion: 0.67
Nodes (2): UpgradeHTTP(), Upgrader

### Community 21 - "Gin Router Test"
Cohesion: 0.83
Nodes (3): ginDoReq(), startGinServer(), TestGinAdapter()

### Community 22 - "Hertz Router Test"
Cohesion: 0.83
Nodes (3): hertzDoReq(), startHertzServer(), TestHertzAdapter()

### Community 23 - "External Server Test"
Cohesion: 0.5
Nodes (4): External Server Test Pattern, os/exec Package, websocket-external-server-test Capability, websocket-external-test Specification

### Community 24 - "Route Type"
Cohesion: 0.67
Nodes (1): Route

### Community 25 - "Upgrade Tests"
Cohesion: 0.67
Nodes (0): 

### Community 26 - "NetHTTP Server Tests"
Cohesion: 0.67
Nodes (0): 

### Community 27 - "App Tests"
Cohesion: 1.0
Nodes (0): 

### Community 28 - "wx README"
Cohesion: 1.0
Nodes (1): wx WebSocket Example README

### Community 29 - "wx WebSocket Tasks"
Cohesion: 1.0
Nodes (1): wx-websocket-client Tasks

### Community 30 - "ws-client Spec"
Cohesion: 1.0
Nodes (1): ws-client Specification

### Community 31 - "Hertz Example"
Cohesion: 1.0
Nodes (1): Hertz Example

### Community 32 - "NetHTTP Example"
Cohesion: 1.0
Nodes (1): NetHTTP Example

## Knowledge Gaps
- **50 isolated node(s):** `Config`, `Server`, `GracefulServer`, `Option`, `HandlerFunc` (+45 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **Thin community `App Tests`** (2 nodes): `app_test.go`, `TestLoadConfig_Default()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `wx README`** (1 nodes): `wx WebSocket Example README`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `wx WebSocket Tasks`** (1 nodes): `wx-websocket-client Tasks`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `ws-client Spec`** (1 nodes): `ws-client Specification`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Hertz Example`** (1 nodes): `Hertz Example`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `NetHTTP Example`** (1 nodes): `NetHTTP Example`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `NewServer()` connect `Gin+Hertz Adapters & Server Factory` to `NetHTTP Adapter`?**
  _High betweenness centrality (0.043) - this node is a cross-community bridge._
- **Why does `init()` connect `Gin+Hertz Adapters & Server Factory` to `NetHTTP Adapter`?**
  _High betweenness centrality (0.041) - this node is a cross-community bridge._
- **Why does `GinServer` connect `Gin Server` to `Gin+Hertz Adapters & Server Factory`?**
  _High betweenness centrality (0.019) - this node is a cross-community bridge._
- **What connects `Config`, `Server`, `GracefulServer` to the rest of the system?**
  _50 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Core Framework Types & APIs` be split into smaller, more focused modules?**
  _Cohesion score 0.08 - nodes in this community are weakly interconnected._
- **Should `NetHTTP Adapter` be split into smaller, more focused modules?**
  _Cohesion score 0.1 - nodes in this community are weakly interconnected._
- **Should `Gin+Hertz Adapters & Server Factory` be split into smaller, more focused modules?**
  _Cohesion score 0.09 - nodes in this community are weakly interconnected._