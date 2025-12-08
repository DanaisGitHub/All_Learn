using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;

public record CurrencyValues {
    [FromQuery]
    public required string FromCurrency { get; init; }
    [FromQuery]
    public required string ToCurrency { get; init; }
    [FromQuery]
    public required double FromValue { get; init; }
}

namespace minimal_api {
    [ApiController]
    [Route("api/[controller]")]
    public class ConcurrencyController : ControllerBase {

        [HttpGet]
        public IActionResult ConvertValues([FromQuery] CurrencyValues values) {

            //var convertedValue = CurrencyConvertor.ConvertValues(values);

            return Ok(
                value:"hello"
            );


        }


        // [HttpGet]
        // public IActionResult Hello() {
        //     return Ok(
        //         value:"hello"
        //     );
            
        // }


    }
}