const $ = require('./jquery-3.2.1');
const _ = require('./underscore');
const Popper = require('./popper');
const accounting = require('./accounting');

window._ = _;
window.Popper = Popper;
window.$ = $;
window.jQuery = $;

require('./bootstrap-4/bootstrap');

const updateQuotes = sortedSyms =>
  $.ajax({
    url: `/api/quotes/${sortedSyms.join(',')}`,
    dataType: 'json',
    success: (data /* , textStatus, jqXHR */) => {
      _.each(data['Stock Quotes'], (quote) => {
        const symbol = quote['1. symbol'].toLowerCase()
        const price = quote['2. price']
        const $symbolRow = $(`#equity-${symbol.replace('.', '\\.')}`)
        $symbolRow.find('div.price')
          .empty()
          .append(accounting.formatMoney(price))

        const numShares = parseInt($symbolRow.find('div.shares').text(), 10)
        $symbolRow.find('div.value')
          .empty()
          .append(accounting.formatMoney(price * numShares))
      })
    }
  })

const buildPortfolioTable = () =>
  $.ajax({
    url: '/api/portfolio/666',
    dataType: 'json',
    success: (data /* , textStatus, jqXHR */) => {
      const sortedSyms = data.map(d => d.symbol).sort()

      _.each(sortedSyms, (symbol) => {
        const equity = data.find(d => d.symbol == symbol)
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
      updateQuotes(sortedSyms)
    }
  })

const buyEquity = () => {
  const numShares = parseInt($('#buyShares').val(), 10) || 0

  $.ajax({
    url: '/api/portfolio/666',
    method: 'POST',
    dataType: 'json',
    data: {
      symbol: $('#buySymbol').val(),
      numShares
    },
    success: (data /* , textStatus, jqXHR */) => {
      console.log(`Successfully bought equity; data = ${JSON.stringify(data)}`)
    }
  })
}

$('body').ready(() => {
  buildPortfolioTable()
})
