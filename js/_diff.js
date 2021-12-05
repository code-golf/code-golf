import * as Diff from 'diff';

function diffWrapper(join, left, right, diffOpts) {
    // join = "\n" for line diff or "" for char diff
    // pass in left,right =  list of tokens;
    //   for char diff, this is a list of chars
    //   for line diff, this is a list of lines
    // Wrapper for performance
    // Include characters until the first difference, then include 1000 characters
    // after that, and treat the rest as a single block
    const d = firstDifference(left, right);
    const length = Math.min(1000, Math.max(left.length - d, right.length - d));
    const diff = (join === "" ? Diff.diffChars : Diff.diffLines)(
        left.slice(d, d + length).join(join),
        right.slice(d, d + length).join(join),
        diffOpts
    );
    const head = left.slice(0, d);
    if (head !== "") {
        diff.unshift({
            count: head.length,
            value: head.join(join)
        })
    }
    const leftTail = left.slice(d + length);
    const ltString = leftTail.join(join);
    const rightTail = right.slice(d + length);
    const rtString = rightTail.join(join);
    if (ltString === rtString) {
        diff.push({
            count: leftTail.length,
            value: ltString
        });
    } else {
        if (ltString !== "") {
            diff.push({
                added: undefined,
                removed: true,
                count: leftTail.length,
                value: ltString
            });
        }
        if (rtString !== "") {
            diff.push({
                added: true,
                removed: undefined,
                count: rightTail.length,
                value: rtString
            });
        }
    }
    return diff;
}

function firstDifference(left, right) {
    for (let i=0; i<left.length || i<right.length; i++) {
        if (left[i] !== right[i]) {
            return i;
        }
    }
    return Math.min(left.length, right.length) + 1;
}

function colFromWidth(className, width) {
    const out = document.createElement("col");
    out.className = className;
    out.style.width = width;
    return out;
}

export function attachDiff(element, hole, exp, out, argv) {
    const isArgDiff = shouldArgDiff(hole, exp, argv);
    element.classList.toggle("diff-arg-type", isArgDiff);
    const {rows, maxLineNum} = diffHTMLRows(hole, exp, out, argv, isArgDiff);

    element.innerHTML = "";
    const table = document.createElement("table");
    table.appendChild(getColgroup(isArgDiff, maxLineNum, argv));
    const tbody = document.createElement("tbody");
    tbody.appendChild(getHeader(isArgDiff));
    for (let row of rows) {
      tbody.appendChild(row);
    }
    table.appendChild(tbody);
    element.appendChild(table);
}

function getColgroup(isArgDiff, maxLineNum, argv) {
    const colgroup = document.createElement("colgroup");
    const numLength = String(maxLineNum).length + 1;
    const charWidth = 11;
    if (isArgDiff) {
        const longestArgLength = Math.max(6, ...argv.map((arg) => arg.length));
        colgroup.appendChild(
          colFromWidth("diff-col-arg", Math.min(longestArgLength * charWidth, 350) + "px")
        );
    } else {
        colgroup.appendChild(
          colFromWidth("diff-col-left-num", numLength * charWidth + "px")
        );
    }
    colgroup.appendChild(colFromWidth("diff-col-left", "auto"));
    if (!isArgDiff) {
        colgroup.appendChild(
          colFromWidth("diff-col-right-num", numLength * charWidth + "px")
        );
    }
    colgroup.appendChild(colFromWidth("diff-col-right", "auto"));
    return colgroup;
}

function getHeader(isArgDiff) {
    const header = document.createElement("tr");
    isArgDiff && header.appendChild(
        elFromText("th", "diff-header-args", "Args")
    );
    isArgDiff || header.appendChild(elFromText("th", "", ""));
    header.appendChild(elFromText("th", "diff-header-output", "Output"));
    isArgDiff || header.appendChild(elFromText("th", "diff-title-bg", ""));
    header.appendChild(elFromText("th", "diff-header-expected", "Expected"));
    return header;
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
    return {
        rows,
        maxLineNum: Math.max(pos.left, pos.right),
    };
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
        return diffWrapper("\n", lines(before), lines(after), {
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
    const isUnchanged = !left.removed && !right.added;
    // Skip the middle of any block of more than 7 unchanged lines, or 41 changed lines
    const padding = isUnchanged ? 3 : 20;
    const skipMiddle = numLines > 2 * padding + 1;
    for (let i=0; i<numLines; i++) {
        let row = document.createElement("tr");
        if (skipMiddle && i === padding) {
            const td = elFromText("td", "diff-skip", `@@ ${numLines - 2 * padding} lines omitted @@`);
            td.colSpan = isArgDiff ? "3" : "4";
            row.appendChild(td);
            rows.push(row);
            continue;
        }
        if (skipMiddle && padding <= i && i < numLines - padding) {
            continue;
        }
        const leftLine = leftSplit[i];
        const rightLine = rightSplit[i];
        const charDiff = diffWrapper(
            "",
            [...leftLine ?? ''],
            [...rightLine ?? ''],
            diffOpts
        );
        // subtract 1 because the lines start counting at 1 instead of 0
        const arg = argv[i + pos.right - 1]
        if (isArgDiff) {
            row.appendChild(elFromText("td", "diff-arg", arg ?? ""));
        }
        if (leftLine !== undefined) {
            isArgDiff || row.appendChild(
                elFromText("td", "diff-left-num", String(i + pos.left))
            )
            row.appendChild(
                renderCharDiff(
                    'diff-left' + (left.removed?' diff-removal':''),
                    charDiff,
                    false
                )
            )
        } else {
            row.appendChild(elFromText("td", "", ""));
            row.appendChild(elFromText("td", "", ""));
        }
        if (rightLine !== undefined) {
            isArgDiff || row.appendChild(
                elFromText("td", "diff-right-num", String(i + pos.right))
            );
            row.appendChild(
              renderCharDiff(
                "diff-right" + (right.added ? " diff-addition" : ""),
                charDiff,
                true
              )
            );
        } else {
            row.appendChild(elFromText("td", "", ""));
            row.appendChild(elFromText("td", "", ""));
        }
        rows.push(row)
    }
    pos.left += left.count;
    pos.right += right.count;
    return rows
}

function renderCharDiff(className, charDiff, isRight) {
    const out = document.createElement("td");
    out.className = className;
    const contents = document.createElement("span");
    out.appendChild(contents);
    for (let change of charDiff) {
        if (change.added && isRight) {
            contents.appendChild(
                elFromText("span", "diff-char-addition", change.value)
            );
        } else if (change.removed && !isRight) {
            contents.appendChild(
                elFromText("span", "diff-char-removal", change.value)
            );
        } else if (!change.added && !change.removed) {
            contents.appendChild(
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
