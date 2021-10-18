import * as Diff from 'diff';

export function attachDiff(element, hole, exp, out, argv) {
    const isArgDiff = getDiffType(hole, exp, argv);
    const header = getHeader(isArgDiff);
    // Limit `out` to 999 lines to avoid slow computation of the line diff
    let outLines = lines(out)
    if (outLines.length > 998) {
        outLines = outLines.slice(0, 998);
        outLines.push("[Truncated for performance]");
        out = outLines.join("\n")
    }
    element.classList.toggle("diff-arg-type", isArgDiff);
    const rows = diffHTMLRows(hole, exp, out, argv, isArgDiff);
    updateDiff(element, rows, header);
    element.onscroll = () => {
        updateDiff(element, rows, header);
    };
}

function getHeader(isArgDiff) {
  const numHeader = isArgDiff ? '' : `<div class='diff-title-bg'></div>`
  return (
    (
      isArgDiff
        ? `<h3 id='diff-arguments'>Args</h3>`
        : `<div class='diff-title-bg'></div>`
    ) +
    numHeader +
    `<h3 id='diff-output'>Output</h3>` +
    numHeader +
    `<h3 id='diff-expected'>Expected</h3>`
  );
}

// pixels
const rowHeight = 20;

function updateDiff(element, rows, header) {
    // -1 for the header row
    const rowTop = element.scrollTop / rowHeight - 1;
    // Buffer of approx 50 rows in each direction
    // Selected rows include rowStart but exclude rowEnd
    // We can only show a few rows for memory reasons:
    // Chrome limits a CSS grid to 1000 rows
    const rowStart = Math.max(Math.floor(rowTop - 50), 0);
    const padStart = rowStart * rowHeight;
    const rowEnd = Math.min(Math.floor(rowTop + 60), rows.length);
    const padEnd = (rows.length - rowEnd - 1) * rowHeight;
    element.innerHTML = `${header}
        <div class="diff-gap" style="height:${padStart}px"></div>
        ${rows.slice(rowStart, rowEnd).join("")}
        <div class="diff-gap" style="height:${padEnd}px"></div>`
}

function diffHTMLRows(hole, exp, out, argv, isArgDiff) {
    let rows = []
    let pos = {
        left: 1,
        right: 1,
        isLastDiff: false
    };
    const changes = getLineChanges(hole, out, exp, isArgDiff);
    let pendingChange = null;
    for (let i = 0; i < changes.length; i++) {
        const change = changes[i]
        pos.isLastDiff = i === changes.length - 1
        if (change.added || change.removed) {
            if (pendingChange === null) {
                pendingChange = change;
            } else {
                rows.push(...getDiffRow(hole, pendingChange, change, pos, argv, isArgDiff));
                pendingChange = null;
            }
        } else {
            if (pendingChange) {
                rows.push(...getDiffRow(hole, pendingChange, {}, pos, argv, isArgDiff));
                pendingChange = null;
            }
            rows.push(...getDiffLines(hole, change, change, pos, argv, isArgDiff));
        }
    }
    if (pendingChange) {
        rows.push(...getDiffRow(hole, pendingChange, {}, pos, argv, isArgDiff));
    }
    return rows
}

function getLineChanges(hole, before, after, isArgDiff) {
    if (isArgDiff) {
        const out = []
        const splitBefore = lines(before)
        const splitAfter = lines(after)
        const compareOpts = {
            sensitivity: shouldIgnoreCase(hole) ? 'accent' : 'base'
        }
        for (let i=0; i<Math.max(splitBefore.length, splitAfter.length); i++) {
            const a = splitBefore[i] ?? '';
            const b = splitAfter[i] ?? '';
            // https://stackoverflow.com/a/2140723/7481517
            const linesEqual = 0 === a.localeCompare(b, undefined, compareOpts);
            if (linesEqual) {
                out.push({
                    count: 1,
                    value: a + '\n'
                })
            } else {
                for (let [k,v] of [['removed', a], ['added', b]]) {
                    if (v !== undefined) {
                        const prev = out[out.length - 1]
                        if (prev && prev[k]) {
                            prev.count++;
                            prev.value += v + '\n'
                        } else {
                            out.push({
                                count: 1,
                                [k]: true,
                                value: v + '\n'
                            })
                        }
                    }
                }
            }
        }
        return out
    } else {
        return Diff.diffLines(before, after, {
            ignoreCase: shouldIgnoreCase(hole)
        })
    }
}

function getDiffRow(hole, change1, change2, pos, argv, isArgDiff) {
    change2.value ??= ''
    change2.count ??= 0
    const left = change1.removed ? change1 : change2
    const right = change1.added ? change1 : change2
    return getDiffLines(hole, left, right, pos, argv, isArgDiff)
}

function getDiffLines(hole, left, right, pos, argv, isArgDiff) {
    const leftSplit = lines(left.value);
    const rightSplit = lines(right.value);
    if (!(pos.isLastDiff && hole === "quine")) {
        // ignore trailing newline
        if (leftSplit[leftSplit.length - 1] === '') leftSplit.pop();
        if (rightSplit[rightSplit.length - 1] === '') rightSplit.pop();
    }
    const diffOpts = {
        ignoreCase: shouldIgnoreCase(hole)
    }
    let rows = []
    const numLines = Math.max(leftSplit.length, rightSplit.length)
    for (let i=0; i<numLines; i++) {
        let s = ''
        const leftLine = leftSplit[i];
        const rightLine = rightSplit[i];
        const charDiff = Diff.diffChars(leftLine ?? '', rightLine ?? '', diffOpts);
        // subtract 1 because the lines start counting at 1 instead of 0
        const arg = argv[i + pos.right - 1]
        if (arg !== undefined && isArgDiff) {
            s += `<div class='diff-arg'>${arg}</div>`
        }
        if (leftLine !== undefined) {
            s += isArgDiff ? '' : `<div class='diff-left-num'>${i + pos.left}</div>`
            s += `<div class='diff-left${left.removed?' diff-removal':''}'>${renderCharDiff(charDiff, false)}</div>`
        }
        if (rightLine !== undefined) {
            s += isArgDiff ? '' : `<div class='diff-right-num'>${i + pos.right}</div>`
            s += `<div class='diff-right${right.added?' diff-addition':''}'>${renderCharDiff(charDiff, true)}</div>`
        }
        rows.push(s)
    }
    pos.left += left.count;
    pos.right += right.count;
    return rows
}

function renderCharDiff(charDiff, isRight) {
    let html = ''
    for (let change of charDiff) {
        if (change.added && isRight) {
            html += `<span class='diff-char-addition'>${change.value}</span>`
        } else if (change.removed && !isRight) {
            html += `<span class='diff-char-removal'>${change.value}</span>`
        } else if (!change.added && !change.removed) {
            html += change.value;
        }
    }
    return html
}

const argBlocklist = ['qr-decoder']

function getDiffType(hole, exp, argv) {
    const expectedLines = lines(exp)
    // The subtracted part removes 1 line in the case of a trailing newline
    const numExpectedLines = expectedLines.length - (lines[lines.length - 1] === '' ? 1 : 0)
    return numExpectedLines === argv.length && !argBlocklist.includes(hole)
}

function shouldIgnoreCase(hole) {
    return hole === "css-colors"
}

function lines(s) {
    return s.split(/\r\n|\n/)
}
