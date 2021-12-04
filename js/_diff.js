import * as Diff from 'diff';

function diffChars(left, right, diffOpts) {
    // Wrapper for performance
    // Include characters until the first difference, then include 1000 characters
    // after that, and treat the rest as a single block
    const d = firstDifference(left, right);
    const length = Math.min(1000, left.length - d, right.length - d);
    const diff = Diff.diffChars(
        left.substring(d, d + length),
        right.substring(d, d + length),
        diffOpts
    );
    const head = left.substring(0, d);
    if (head !== "") {
        diff.unshift({
            count: head.length,
            value: head
        })
    }
    const leftTail = left.substring(d + length);
    const rightTail = right.substring(d + length);
    if (leftTail === rightTail) {
        diff.push({
            count: leftTail.length,
            value: leftTail
        });
    } else {
        if (leftTail !== "") {
            diff.push({
                added: undefined,
                removed: true,
                count: leftTail.length,
                value: leftTail,
            });
        }
        if (rightTail !== "") {
            diff.push({
                added: true,
                removed: undefined,
                count: rightTail.length,
                value: rightTail
            });
        }
    }
    return diff;
}

function firstDifference(left, right) {
    for (let i=0; i<left.length && i<right.length; i++) {
        if (left[i] != right[i]) {
            return i;
        }
    }
    return Math.min(left.length, right.length) + 1;
}

export function attachDiff(element, hole, exp, out, argv) {
    const isArgDiff = shouldArgDiff(hole, exp, argv);
    const header = getHeader(isArgDiff);
    // Limit `out` to 1001 lines to avoid slow computation of the line diff
    // Limit must be >1000 for Van Eck Sequence
    let outLines = lines(out)
    if (outLines.length > 1000) {
        outLines = outLines.slice(0, 1000);
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
        <div class="diff-gap" style="height:${padEnd}px"></div>`
    // ${rows.slice(rowStart, rowEnd).join("")}
    for (let row of rows) {
        for (let newNode of row) {
            element.insertBefore(
                newNode,
                element.lastElementChild
            );
        }
    }
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

function elFromText(nodeType, className, innerText) {
    const out = document.createElement(nodeType);
    out.className = className;
    out.innerText = innerText;
    return out;
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
        let s = []
        const leftLine = leftSplit[i];
        const rightLine = rightSplit[i];
        const charDiff = diffChars(leftLine ?? '', rightLine ?? '', diffOpts);
        // subtract 1 because the lines start counting at 1 instead of 0
        const arg = argv[i + pos.right - 1]
        if (arg !== undefined && isArgDiff) {
            s.push(elFromText("div", "diff-arg", arg))
        }
        if (leftLine !== undefined) {
            isArgDiff && s.push(
                elFromText("div", "diff-left-num", String(i + pos.left))
            )
            s.push(
                renderCharDiff(
                    'diff-left' + (left.removed?' diff-removal':''),
                    charDiff,
                    false
                )
            )
        }
        if (rightLine !== undefined) {
            isArgDiff && s.push(elFromText("div", "diff-right-num", String(i + pos.right)));
            s.push(
                renderCharDiff(
                    'diff-right' + (right.added?' diff-addition':''),
                    charDiff,
                    true
                )
            )
        }
        rows.push(s)
    }
    pos.left += left.count;
    pos.right += right.count;
    return rows
}

function renderCharDiff(className, charDiff, isRight) {
    const out = document.createElement("div");
    out.className = className;
    for (let change of charDiff) {
        if (change.added && isRight) {
            out.appendChild(
                elFromText("span", "diff-char-addition", change.value)
            );
        } else if (change.removed && !isRight) {
            out.appendChild(
                elFromText("span", "diff-char-removal", change.value)
            );
        } else if (!change.added && !change.removed) {
            out.appendChild(
                document.createTextNode(change.value)
            );
        }
    }
    return out
}

function shouldArgDiff(hole, exp, argv) {
    const expectedLines = lines(exp)
    // The subtracted part removes 1 line in the case of a trailing newline
    const numExpectedLines = expectedLines.length - (lines[lines.length - 1] === '' ? 1 : 0)
    // Exclude holes such as qr-decoder, morse-decoder, and morse-encoder, which have only one (big) arg
    return numExpectedLines === argv.length && argv.length > 1
}

function shouldIgnoreCase(hole) {
    return hole === "css-colors"
}

function lines(s) {
    return s.split(/\r\n|\n/)
}
