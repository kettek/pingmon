const glyphs: string[][] = [
  /*
  [
    "00000",
    "00000",
    "00000",
    "00000",
    "00000",
  ],
  */
  // 0
  [
    "01110",
    "10001",
    "10001",
    "10001",
    "01110",
  ],
  // 1
  [
    "01100",
    "00100",
    "00100",
    "00100",
    "01110",
  ],
  // 2
  [
    "02220",
    "20002",
    "00020",
    "00200",
    "22222",
  ],
  // 3
  [
    "03330",
    "30003",
    "00330",
    "30003",
    "03330",
  ],
  // 4
  [
    "10001",
    "10001",
    "11111",
    "00001",
    "00001",
  ],
  // 5
  [
    "11111",
    "10000",
    "11110",
    "00001",
    "11110",
  ],
  // 6
  [
    "01100",
    "10000",
    "11110",
    "10001",
    "01110",
  ],
  // 7
  [
    "01111",
    "00001",
    "00010",
    "00100",
    "00100",
  ],
  // 8
  [
    "01110",
    "10001",
    "01110",
    "10001",
    "01110",
  ],
  // 9
  [
    "01110",
    "10001",
    "01111",
    "00001",
    "00110",
  ],
  // /
  [
    "001",
    "001",
    "010",
    "010",
    "100",
    "100",
  ],
]


const faviconCache: {[key: string]: string} = {}
/**
 * 
 * @param img Source image to overlay numbers onto
 * @param min Lefthand number
 * @param max Righthand number
 */
export function MakeFavicon(img: HTMLImageElement, min: number, max: number): string {
  const minStr = `${min}`, maxStr = `${max}`
  const k = `${min}/${max}`
  if (faviconCache[k]) return faviconCache[k]
  const canvas: HTMLCanvasElement = document.createElement('canvas')
  canvas.width = img.naturalWidth
  canvas.height = img.naturalHeight
  const ctx = canvas.getContext('2d')

  ctx.imageSmoothingEnabled = false
  ctx.drawImage(img, 0, 0)

  // If max === min, then just return the icon without text.
  if (max === min) {
    return faviconCache[k] = canvas.toDataURL('image/png')
  }

  // Otherwise let's draw our lil text.
  if (max < 10 && min < 10) { // Draw at the bottom if we have room.
    let x = 1, y = 9
    ;[x] = DrawGlyph(ctx, min, x, y+1)
    ;[x] = DrawGlyph(ctx, 10, x+1, y)
    ;[x] = DrawGlyph(ctx, max, x+1, y+1)
  } else { // Otherwise split it to the top-left for min and bottom-right for max.
    let x = 0, y = 0
    for (let i = 0; i < minStr.length; i++) {
      ;[x] = DrawGlyph(ctx, Number(minStr[i]), x, y)
    }
    x = 5, y = 9
    for (let i = 0; i < maxStr.length; i++) {
      ;[x] = DrawGlyph(ctx, Number(maxStr[i]), x, y)
      x++
    }
  }

  return faviconCache[k] = canvas.toDataURL('image/png')
}

function DrawGlyph(ctx: CanvasRenderingContext2D, index: number, offsetX: number, offsetY: number): [number, number] {
  if (index < 0 || index >= glyphs.length) return

  let blackPixel = ctx.createImageData(1, 1)
  blackPixel.data[0] = 0
  blackPixel.data[1] = 0
  blackPixel.data[2] = 0
  blackPixel.data[3] = 255
  let whitePixel = ctx.createImageData(1, 1)
  whitePixel.data[0] = 255
  whitePixel.data[1] = 255
  whitePixel.data[2] = 255
  whitePixel.data[3] = 255

  let drawFatPixel = (pixel: ImageData, x: number, y: number) => {
    ctx.putImageData(pixel, x, y)
    ctx.putImageData(pixel, x+1, y)
    ctx.putImageData(pixel, x+1, y+1)
    ctx.putImageData(pixel, x, y+1)
  }

  // Draw outline.
  let maxRow = glyphs[index].length
  let maxColumn = 0
  for (let row = 0; row < glyphs[index].length; row++) {
    for (let column = 0; column < glyphs[index][row].length; column++) {
      maxColumn = glyphs[index][row].length
      if (glyphs[index][row][column] === '0') continue
      let x = (offsetX+column) * 2
      let y = (offsetY+row) * 2
      drawFatPixel(blackPixel, x-1, y)
      drawFatPixel(blackPixel, x+1, y)
      drawFatPixel(blackPixel, x, y-1)
      drawFatPixel(blackPixel, x, y+1)
    }
  }

  // Draw center.
  for (let row = 0; row < glyphs[index].length; row++) {
    for (let column = 0; column < glyphs[index][row].length; column++) {
      if (glyphs[index][row][column] === '0') continue
      let x = (offsetX+column) * 2
      let y = (offsetY+row) * 2
      drawFatPixel(whitePixel, x, y)
    }
  }

  return [offsetX+maxColumn,offsetY+maxRow]
}