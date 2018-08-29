
function handleArrivals(arrivals) {
  const directions = {}
  arrivals
  .sort((a, b) => a.time - b.time)
  .forEach(arr => {
    if (!directions[arr.direction]) directions[arr.direction] = {}
    if (!directions[arr.direction][arr.train]) directions[arr.direction][arr.train] = []
    directions[arr.direction][arr.train].push(arr.time)
  })

  const results = document.createElement('div')
  results.classList.add('trains')

  const now = Math.round(new Date().valueOf() / 1000)

  Object.keys(directions).sort()
  .forEach(direction => {
    const dirEl = document.createElement('div')
    dirEl.setAttribute('class', 'direction')
    const dirHead = document.createElement('h2')
    dirHead.textContent = direction
    dirEl.appendChild(dirHead)

    const trains = directions[direction]

    Object.keys(trains).sort()
    .forEach(train => {
      const trainEl = document.createElement('div')
      trainEl.setAttribute('class', 'train')
      const trainName = document.createElement('div')
      trainName.setAttribute('class', 'trainName')
      trainName.textContent = train
      trainEl.appendChild(trainName)

      const times = trains[train].map(timestamp => {
        return Math.round((timestamp - now) / 60)
      })

      const trainTimes = document.createElement('div')
      trainTimes.setAttribute('class', 'trainTimes')
      trainTimes.textContent = times.join(', ')
      trainEl.appendChild(trainTimes)

      dirEl.appendChild(trainEl)
    })
    results.appendChild(dirEl)
  })
  document.body.appendChild(results)
}

function runStationPage(id) {
  fetch(`/api/stations/${id}/arrivals`)
  .then(function(response) {
    response.json().then(handleArrivals)
  })

  fetch(`/api/stations/${id}`)
  .then(response => {
    response.json().then(station => {
      const header = document.createElement('h1')
      header.textContent = station.name
      document.body.prepend(header)
    })
  })
}

var match = document.location.pathname.match(/\/stations\/(\d+)/)
if (match) {
  var stationId = match[1]
  runStationPage(stationId)
} else {
  runSearchPage()
}
