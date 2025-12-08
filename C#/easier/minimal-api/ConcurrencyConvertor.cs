using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace minimal_api {
    public static class CurrencyConvertor {
        private static Dictionary<string, double> currencyRatesByUSD = new(){
                {"USD",1.0},
                {"GBP",0.8},
                {"EUR",1.2}
            };

        public static double ConvertValues(CurrencyValues currencyValues) {
            var currency1 = currencyValues.FromCurrency;
            var currency2 = currencyValues.ToCurrency;
            var value1 = currencyValues.FromValue;

            double usdValue = convertToUSD(currency1, value1);
            double convertedValue = uSDToXCurrency(usdValue, currency2);
            return convertedValue;
        }

        private static double convertToUSD(string currency, double value) {
            currency = currency.ToUpper();
            if (currency.Equals("USD")) {
                return value;
            }
            if (!currencyRatesByUSD.ContainsKey(currency)) {
                throw new Exception($"currency '{currency}' not found");
            }
            double rate = currencyRatesByUSD[currency];
            return value * rate;
        }

        private static double uSDToXCurrency(double value, string newCurr) {
            double rate;
            lock (currencyRatesByUSD) {
                if (!currencyRatesByUSD.ContainsKey(newCurr)) {
                    throw new Exception($"currency '{newCurr}' not found");
                }
                rate = currencyRatesByUSD[newCurr];
            }

            return rate * value;

        }
    }
}