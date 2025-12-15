Since you are an experienced Go dev, you don't need to learn "how to code." You need to learn the **.NET ecosystem**, **Dependency Injection (DI)**, and **Entity Framework (EF) Core**.

Do these 4 projects in order. They map your Go knowledge to the C# way of doing things.

### Project 1: The "Minimal API" (The Go Transition)

**Goal:** Get comfortable with C# syntax, Program.cs structure, and HTTP routing without the boilerplate.
**Context:** .NET "Minimal APIs" look almost exactly like Go's `Gin` or `Echo`.

- **Build:** A currency converter API.
- **Requirements:**
    - Hardcode exchange rates in a `Dictionary`.
    - Accept `GET /convert?from=USD&to=EUR&amount=100`.
    - Return JSON `{"amount": 95.50}`.
    - Use a `record` type for the response (immutable data structure, like a struct).
- **What you learn:** Top-level statements, Records vs Classes, JSON serialization (built-in), simple routing.

### Project 2: The "Enterprise" CRUD API

**Goal:** Learn the "Microsoft Way." Dependency Injection, Controllers, and EF Core (ORM).
**Context:** Go prefers explicit SQL; .NET runs on ORMs. You must learn EF Core to be employable.

- **Build:** An Inventory Management API.
- **Requirements:**
    - **Database:** Use SQLite or LocalDB.
    - **ORM:** Use Entity Framework Core. Code-First approach (Define C# classes $\rightarrow$ Generate DB migration).
    - **Architecture:** Create an `IInventoryService` interface and an `InventoryService` implementation. Inject it into a Controller.
    - **Endpoints:** GET, POST, PUT, DELETE for Items.
    - **Validation:** Use `Data Annotations` (attributes on properties) to validate input (e.g., `[Required]`, `[Range(0, 100)]`).
- **What you learn:** `DbContext`, Migrations, Interface-based DI, Async/Await database calls (crucial difference from Goroutines), Controller-Service-Repository pattern.

### Project 3: The Azure Function "Glue"

**Goal:** Serverless logic and Event-Driven Architecture.
**Context:** Moving logic out of the main API into background workers.

- **Build:** An Image Resizer & Logger.
- **Requirements:**
    - Create a local **Azure Storage Account** (use generic emulator or Azurite).
    - **Function 1 (HTTP Trigger):** Accepts a file upload. Saves it to a Blob Container named `raw-images`.
    - **Function 2 (Blob Trigger):** Listens to `raw-images`. When a file is added, it triggers, logs the size, mimics "resizing" (just copy it to a `processed` container), and writes a metadata entry to a generic Table Storage.
- **What you learn:** `Triggers` (Events that start code), `Bindings` (Declarative input/output), `Stream` manipulation in C#, `local.settings.json` configuration.

### Project 4: The Full System (The Capstone)

**Goal:** Tie it all together using Asynchronous Messaging.
**Context:** Handling high throughput by offloading work to a queue (The standard Cloud pattern).

- **Build:** A "Report Generation" System.
- **Requirements:**
    1. **Web API (.NET Core):** Endpoint `POST /reports/request`. Accepts a date range. It does **not** generate the report. It places a message `{ "reportId": "guid", "range": "..." }` onto an **Azure Queue** and returns `202 Accepted` immediately.
    2. **Worker (Azure Function):** Queue Trigger listening to that queue. Dequeues the message, waits 5 seconds (simulating work), calculates dummy data, and updates a SQL Database record status to "Completed".
    3. **Polling Endpoint:** Add `GET /reports/{id}` to the Web API to check the status in the DB.
- **What you learn:** Queue storage patterns, Producer/Consumer pattern, dealing with eventual consistency, handling "fire and forget" tasks properly in C#.

### Critical Syntax mappings for you:

- **Goroutines** $\rightarrow$ `Task.Run` (But rarely used manually in web apps).
- **Channels** $\rightarrow$ `System.Threading.Channels` (Advanced) or just Queues.
- **Defer** $\rightarrow$ `using` statement (IDisposable) or `try/finally`.
- **Struct tags** $\rightarrow$ Attributes `[JsonPropertyName("id")]`.
- **`err != nil`** $\rightarrow$ `try { } catch (Exception ex) { }` (Exceptions are control flow in C#, get used to it).

Start Project 1 right now. Create a folder, run `dotnet new web`, and open VS Code. Move.