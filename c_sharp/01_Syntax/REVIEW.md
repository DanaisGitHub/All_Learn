# Task 1 Review

## Findings
- `ItemService.cs:35-38` handles nulls by auto-filling defaults instead of rejecting invalid input. For the nullability goal, validate `request` and its properties with `ArgumentNullException.ThrowIfNull` rather than coercing.
- `ItemService.cs:42-60` locks return live views of the backing `List`; callers can observe later mutations. Return a snapshot (e.g., copy to array) and keep storage separate from DTOs via mapping.
- Unused `using` directives (`ItemService.cs:2-5`, `IItemService.cs:1-2`) add noise—remove them.
- DI wiring and tests are missing; the service isn’t registered in a container and no unit tests cover create/get/filter/null/cancellation cases.
- Naming/shape: consider `sealed class ItemService`, PascalCase private methods, and an internal model type instead of storing DTOs directly.

## DI (Dependency Injection)
Declare dependencies in constructors and let a container provide them. Register with something like `services.AddSingleton<IItemService, ItemService>();` and consume `IItemService` via injection instead of `new`. This separates wiring from logic and improves testability.
