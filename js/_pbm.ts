// Decoder for ASCII Portable {Grey,Pix}Map images (P2, P3) to canvas.
//
// https://en.wikipedia.org/wiki/Netpbm
// http://netpbm.sourceforge.net/doc/pgm.html
// http://netpbm.sourceforge.net/doc/ppm.html

export default (pbm: any) => {
    if (!/^P[23](?:\s+\d+){3,}\s*$/.exec(pbm)) return;  // Basic sanity check.

    const [format, width, height, max, ...data] = pbm.trim().split(/\s+/);

    const stride = format == 'P2' ? 1 : 3;

    // Maximum color value valid range.
    if (!(0 < max && max < 65536)) return;

    // We have as much data as the claimed dimensions.
    if (width * height * stride != data.length) return;

    const canvas  = document.createElement('canvas');
    // getContext is not called earlier, so we know ctx can't be null
    const ctx     = canvas.getContext('2d') as CanvasRenderingContext2D;
    canvas.width  = width;
    canvas.height = height;

    for (let i = 0; i < data.length; i += stride) {
        const rgb = stride == 1
            ? [data[i], data[i], data[i]] : [data[i], data[i+1], data[i+2]];

        ctx.fillStyle = `rgb(${rgb.map(c => c / max * 255).join(',')})`;

        ctx.fillRect(Math.floor(i / stride / width), i / stride % width, 1, 1);
    }

    return canvas;
};
