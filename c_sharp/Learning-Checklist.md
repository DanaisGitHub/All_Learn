# Learning Checklist for Mastering the 21D Backend

## Core C# and .NET

- [ ] Get comfortable with C# syntax differences from Go (types, generics, properties, using, nullability).
- [ ] Learn async/await, Task/Task<T>, cancellation tokens, and avoiding blocking calls.
- [ ] Understand project/solution structure: .sln vs .csproj, target frameworks, PackageReference, conditional compilation symbols (DEBUG, QA, RELEASE).
- [ ] Use the dotnet CLI: restore, build, run, test, add package, add reference, sln add/remove.
- [ ] Understand NuGet package management and versioning.
- [ ] Learn JSON serialization via Newtonsoft.Json (settings, attributes, converters).
- [ ] Learn DI patterns in .NET: IServiceCollection, singleton vs scoped vs transient (Functions uses singleton-heavy patterns).

## Azure Functions (Isolated Worker)

- [ ] Understand the isolated worker model vs in-process: HostBuilder, ConfigureFunctionsWebApplication, WorkerOptions.Serializer.
- [ ] Learn function triggers/attributes (HttpTrigger, TimerTrigger) and bindings.
- [ ] Understand middleware in isolated Functions (IFunctionsWorkerMiddleware) and how ControllerMiddleware works here.
- [ ] Local development with Azure Functions Core Tools (`func start`), host.json/local.settings.json usage.

## Project-Specific Architecture

- [ ] LibraryDI: how services are registered, secrets fetched, and AppDBUnit created (21DLibrary/LibraryDI.cs).
- [ ] Auth/JWT flow: JwtManager, JwtSettings, claims used (device guid mandatory, user guid optional), validation paths.
- [ ] HTTP controllers: BaseController helpers (ExecuteAfterAuthentication, HubSpot sync, validators), UserController flows, FolderController usage.
- [ ] Vimeo sync worker: Program.cs init, VimeoFunctions recursion, timer cron syntax, seeding requirements (root folders in app.VideoFolder).
- [ ] DTOs and serialization: Request/Response classes, JsonProperty usage, JsonUtils settings.
- [ ] Error handling: BaseHTTPException subclasses (BadRequestException, ForbiddenException, etc.), how middleware maps them to responses.
- [ ] Logging: AppLogger behavior, dual logging to ILogger + DB (LogDAO), how ControllerMiddleware logs unexpected errors.

## Data Access & SQL

- [ ] SQL Server fundamentals: schemas, stored procedures, transactions, parameterization.
- [ ] ADO.NET patterns used in BaseDAO (ExecuteNonQuery, GetEntity, GetList, transactions).
- [ ] Connection pooling settings (MinDBPoolSize/MaxDBPoolSize) and how they are applied.
- [ ] Schema expectations: key stored procs like app.User_*, app.Device_*, app.Video*, app.RawEvent_Insert, log.*.
- [ ] Entities/DAOs mapping: UserDAO, DeviceDAO, UserDeviceDAO, UserLoginAuditDAO, VideoDAO, VideoFolderDAO, RawEventDAO.

## External Integrations

- [ ] Azure Key Vault: how SecretClient is created (AzureUtils.GetSecretClient), managed identity, GetSecretValueSync extensions.
- [ ] HubSpotClient: auth via bearer token, search/create/update contacts, conversations, form submission, chat tokens.
- [ ] VimeoClient: initialization with API key, folder/video listing, paging helpers.
- [ ] FirebaseClient: Firebase Admin SDK basics, token validation, user lookup.

## Configuration & Secrets

- [ ] appsettings.json variants (debug/qa/release) and conditional loading in Program.cs.
- [ ] local.settings.json for local secrets; never commit it.
- [ ] Environment variables vs configuration binding (.Get<T>() for JwtSettings).
- [ ] Secret naming conventions in Key Vault (DBConnectionStringApp, JWTSecret, HubSpotApiKey, VimeoApiKey, GoogleCredentials).

## Testing Strategy

- [ ] NUnit basics and attributes (Test, TestCase, Explicit, SetUp).
- [ ] TestServer usage in BaseTest and Startup to reuse DI graph.
- [ ] Writing integration tests that call controllers directly (no HTTP layer) and asserting DB effects.
- [ ] Using ConsoleLoggerFactory for test logging.

## Azure Operations

- [ ] Azure Functions deployment (zip deploy, `func azure functionapp publish`, or CI/CD) for both API and Vimeo worker.
- [ ] Managed Identity setup for Function Apps to access Key Vault and DB.
- [ ] Application Insights basics: configuring connection string, querying logs with Kusto (KQL), sampling, and custom traces.
- [ ] Timer trigger behavior in Azure (retries, schedule format, RunOnStartup implications in production).
- [ ] Network considerations: VNET integration if used, outbound access for HubSpot/Vimeo/Firebase.

## CI/CD (GitHub Actions assumed)

- [ ] GitHub Actions fundamentals: workflows, jobs, runners, secrets, environments.
- [ ] Typical .NET pipeline steps: checkout, setup-dotnet, restore, build, test, publish artifacts, deploy Functions.
- [ ] Storing secrets securely in GitHub (Key Vault integration or GitHub secrets) and passing to deployments.

## Security & Compliance

- [ ] JWT best practices (expiry, signing keys, audience/issuer validation).
- [ ] Secrets hygiene (no secrets in repo, rotate in Key Vault, least privilege identities).
- [ ] Input validation patterns (reuse BaseController validators, avoid over-trusting client data).
- [ ] HTTPS-only endpoints and CORS considerations for Functions.

## Observability & Troubleshooting

- [ ] Using AppLogger outputs and DB log tables to trace issues.
- [ ] Leveraging Application Insights logs/metrics to diagnose failures (dependency calls to HubSpot/Vimeo/Firebase/SQL).
- [ ] Understanding ControllerMiddleware error envelopes and how to surface stack traces in DEBUG/QA only.

## Performance & Reliability

- [ ] HttpClient reuse (already handled in DI) and rate limiting if needed (RateLimitedHttpHandler).
- [ ] SQL performance basics: avoiding N+1 patterns, using stored procedures efficiently, reviewing execution plans when needed.
- [ ] Timer trigger retry/backoff expectations; handling poison messages is not applicable but failures will retry.

## Local Development Environment

- [ ] Install dotnet 10 SDK, Azure Functions Core Tools, Azure CLI.
- [ ] Configure az login for Key Vault access and Function App management.
- [ ] Set up local.settings.json with required secrets for API and Vimeo worker.
- [ ] Run API locally (`dotnet run --project 21DAppServer/21DAppServer.csproj`) and worker (`dotnet run --project 21DVimeoSync/21DVimeoSync.csproj`).

## Domain-Specific Workflows

- [ ] User lifecycle: signup -> HubSpot upsert -> JWT issuance -> device link -> audits.
- [ ] Device-only JWT flow for anonymous video access.
- [ ] HubSpot sync back into local DB via BaseController.CheckLocalDataMatchesHubSpotData.
- [ ] Vimeo sync seeding and recursion logic; how FolderController reads the synced data.
- [ ] Event ingestion path: RecordEvent -> RawEventDAO -> analytics tables.

## Data Hygiene & Seeding

- [ ] How to seed `app.VideoFolder` with root folders before running Vimeo sync.
- [ ] How to reset or soft-delete users (UserController.DeleteUser) and implications for tokens.
- [ ] How login audits and raw events accumulate; consider retention/cleanup strategy.

## Stretch Topics (as you mature)

- [ ] Writing custom middleware in the isolated worker.
- [ ] Adding new Function triggers (queues, blobs) if the platform expands.
- [ ] Advanced App Insights dashboards/alerts for uptime and failure rates.
- [ ] Structured logging/telemetry correlation across HTTP and timer workloads.
- [ ] Schema evolution and migration strategy if stored procs change (coordination with DBA or scripts).
