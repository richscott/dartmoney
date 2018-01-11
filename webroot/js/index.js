const $ = require('./jquery-3.2.1');
window.$ = window.jQuery = $;
const _ = require('./underscore');
window._ = _;
const Popper = require('./popper');
window.Popper = Popper;
require('./bootstrap-4/bootstrap');
const accounting = require('./accounting');

$('body').ready(function() {
  $.ajax({
    url: '/api/portfolio/666',
    dataType: 'json',
    success: function(data, textStatus, jqXHR) {
      console.error(`portfolio data =\n${JSON.stringify(data, null, '  ')}`)
      const sortedSyms = _.keys(data).sort()

      _.each(sortedSyms, (symbol) => {
        const equity = data[symbol]
        const newRow = `
          <div class="row" id="equity-${symbol}">
            <div class="col-2 symbol">${symbol.toUpperCase()}</div>
            <div class="col-4 name">${equity.name}</div>
            <div class="col-2 price" style="text-align:right">&#151;</div>
            <div class="col-2 shares" style="text-align:right">${equity.shares}</div>
            <div class="col-2 value" style="text-align:right">&#151;</div>
          </div>
        `
        $('#portfolioList').append(newRow)
      })

      $.ajax({
        url: `/api/quotes/${sortedSyms.join(',')}`,
        dataType: 'json',
        success: function(data, textStatus, jqXHR) {
          console.log(`quotes data =\n${JSON.stringify(data, null, '  ')}`)

          _.each(data['Stock Quotes'], (quote) => {
            const symbol = quote['1. symbol'].toLowerCase()
            const price = quote['2. price']
            $(`#equity-${symbol} > div.price`)
              .empty()
              .append(accounting.formatMoney(price))

            const numShares = parseInt($(`#equity-${symbol} > div.shares`).text(), 10)
            $(`#equity-${symbol} > div.value`)
              .empty()
              .append(accounting.formatMoney(price * numShares))
          })
        }
      })
    }
  })
})
