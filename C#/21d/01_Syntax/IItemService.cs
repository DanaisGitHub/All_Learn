using System.Dynamic;
using Microsoft.VisualBasic;

/// <summary>
/// A Record that is the an item of the stored data
/// </summary>
/// <param name="Name">the Name</param>
/// <param name="Category">the Category</param>
public sealed record CreateItemRequest(string Name, string Category);

public sealed record ItemResponse(Guid Id, string Name, string Category);

public interface IItemService {
  Task<ItemResponse> CreateAsync(CreateItemRequest request, CancellationToken ct = default);
  Task<ItemResponse?> GetAsync(Guid id, CancellationToken ct = default);
  Task<IReadOnlyList<ItemResponse>> ListAsync(string? categoryFilter, CancellationToken ct = default);
}




