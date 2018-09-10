const GRID_X = 3 * 28;
const GRID_Y = 2 * 7;

function range(length) {
  return Array.from({ length: length })
}

function setupGrid(width, height) {
  const grid = document.querySelector('#grid')

  range(height).forEach((_, y) => {
    const row = document.createElement('div')
    row.setAttribute('class', 'row')

    range(width).forEach((_, x) => {
      const pixel = document.createElement('div')
      pixel.setAttribute('class', 'pixel')
      pixel.setAttribute('id', `pixel-${x}-${y}`)
      row.appendChild(pixel)
    })
    grid.appendChild(row)
  })
}

function print(grid) {
  grid.forEach((row, y) => {
    row.forEach((on, x) => {
      const pixel = document.querySelector(`#pixel-${x}-${y}`)
      pixel.setAttribute('class', on ? 'pixel on' : 'pixel')
    })
  })
}

setupGrid(GRID_X, GRID_Y)

function randomGrid(width, height) {
  return range(height).map((_, y) => {
    return range(width).map((_, x) => {
      return Math.round(Math.random())
    })
  })
}

setInterval(() => {
  print(randomGrid(GRID_X, GRID_Y))
}, 500)
