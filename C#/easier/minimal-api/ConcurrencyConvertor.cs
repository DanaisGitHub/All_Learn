using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;

namespace minimal_api
{
    [ApiController]
    [Route("api/[controller]")]
    public class ConcurrencyConvertor : ControllerBase
    {
        private static Dictionary<string, double> currencyRatesByUSD = new();

        private void setCurrencyRates()
        {
            currencyRatesByUSD = new Dictionary<string, double>
            {
                {"USD",1.0},
                {"GBP",0.8},
                {"EUR",1.2}
            };
        }

        public double convertValues(string currency1, double value1, string currency2)
        {
            value1 


        }

        private double convertToUSD(string currency, double value)
        {
            currency = currency.ToUpper();
            if (currency.Equals("USD"))
            {
                return value;
            }

            if (!currencyRatesByUSD.ContainsKey(currency))
            {
                throw new Exception($"currency '{currency}' not found");
            }

            var rate = currencyRatesByUSD[currency];

            return value * rate;

        }

        private


    }
}