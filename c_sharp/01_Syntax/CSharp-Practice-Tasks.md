# C# Practice Tasks for Backend Readiness

1) In-memory API/service (async + DI + nullability)
- Build a tiny service that manages an in-memory list (no DB). Define request/response DTOs as `record` types.
- Enable nullable reference types and handle nulls explicitly (annotations + `ThrowIfNull`).
- Define `IItemService` and an implementation registered via DI; expose async `Task` methods that accept a cancellation token.
- Add unit tests (xUnit/NUnit) covering creation, retrieval, filtering, and a null-handling edge case.
- Goal: practice idiomatic async, DI wiring, and null safety without infrastructure.
2) OOP modeling (inheritance + interfaces + value types)
- Model a small hierarchy: an `abstract` base (e.g., `Document`), two derived classes with overrides, and an interface with a default implementation (e.g., `IRenderable.Render()` providing a default string).
- Add a `record struct` for a lightweight value (e.g., `DocumentId`), and show copy vs mutation semantics compared to a `class`.
- Demonstrate polymorphic dispatch (virtual/override), pattern matching (`switch`/`is`), and equality differences (`record` vs `class`).
- Tests should assert dispatch works, equality behaves as expected, and structs stay value-like (copies do not share state).
- Goal: solidify C# object model differences (virtual dispatch, records, value vs reference semantics).
3) Generics, constraints, and variance
- Create a small generic pipeline or cache with constraints: `where T : class`, `where T : struct`, `where T : notnull`, `where T : new()`; show how constraints change usage.
- Define a variant interface `ITransformer<in TIn, out TOut>` and demonstrate safe covariance/contravariance with a couple of types.
- Include a test showing mutation through a reference type propagates, while a struct copy does not; also show compile-time restriction enforced by a constraint (e.g., cannot pass null to `notnull`).
- Goal: get comfortable with C# generics syntax and variance rules that surface often in backend abstractions.
4) Properties and resource handling
- Demonstrate auto-properties with `init` and `required`, plus a custom getter/setter that validates and throws `ArgumentOutOfRangeException` on bad input.
- Implement a disposable component with both `IDisposable` and `IAsyncDisposable` (e.g., a fake stream). Use `using` and `await using` to ensure deterministic cleanup.
- Add a test that asserts disposal/async disposal was called (e.g., flags or counters) and that resources are cleaned up even when exceptions occur.
- Goal: internalize property syntax options and the `using`/`await using` patterns relied on in .NET backend code.

Turn-in suggestion

- Keep code/tests in a small sample project folder (e.g., `samples/SyntaxPractice`) and include brief README notes on what you observed for each task.
