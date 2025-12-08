var builder = WebApplication.CreateBuilder(args); 
// wft is this?
{
    builder.Services.AddControllers();
}
var app = builder.Build();
app.MapControllers();

app.MapGet("/", () => "Hello World!\n");


app.Run();

