
using System.Collections.Immutable;
using System.Collections.ObjectModel;
using System.Data;
using System.Security.Cryptography;
using System.Security.Cryptography.X509Certificates;

public class ItemServer : IItemService {

    private readonly List<ItemResponse> store = new List<ItemResponse>();


    public ItemServer() {}
    public ItemServer(List<ItemResponse> initStore) {
        this.store = initStore;
    }
    public Task<ItemResponse> CreateAsync(CreateItemRequest request, CancellationToken ct = default) {
        ct.ThrowIfCancellationRequested();
        ItemResponse item = makeItem(request);
        return addToStore(item);
        
    }

    public Task<ItemResponse?> GetAsync(Guid id, CancellationToken ct = default) {
        ct.ThrowIfCancellationRequested();

        return readFromStore(id);
    }

    public Task<IReadOnlyList<ItemResponse>> ListAsync(string? categoryFilter, CancellationToken ct = default) {
        ct.ThrowIfCancellationRequested();
        return getLists(categoryFilter ?? "");
    }

    private ItemResponse makeItem(CreateItemRequest request) {
        request.Deconstruct(out string name, out string category);
        if (string.IsNullOrEmpty(name)) {
            throw new NoNullAllowedException("value of name is empty");
        }
        if (string.IsNullOrEmpty(category)) {
            throw new NoNullAllowedException("value of category is empty");
        }
        Guid guid = Guid.NewGuid();
        ItemResponse item = new ItemResponse(guid, name, category);
        return item;
    }

    private Task<ItemResponse> addToStore(ItemResponse item) {
        lock(store){
        this.store.Add(item);
        }
        return Task.FromResult(item);
        
    }

    private Task<ItemResponse?> readFromStore(Guid id) {
        lock(store){
        return Task.FromResult(store.Find(item => item.Id == id));
        }
    }

    private Task<IReadOnlyList<ItemResponse>> getLists(string filter) {
        lock(store){
        List<ItemResponse> lists = string.IsNullOrWhiteSpace(filter) ? store : store.FindAll(item => item.Category.Contains(filter, StringComparison.OrdinalIgnoreCase));
        IReadOnlyList<ItemResponse> readOnlyList = lists.AsReadOnly();
        return Task.FromResult(readOnlyList);
        }
    }


}