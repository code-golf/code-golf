import * as Diff  from "diff";

export default function diffHTML(hole, exp, out, argv) {
    let html = getDiffType(hole) === 'arg'
        ? `<h3 id='diff-arguments'>Args</h3>`
        : `<div class='diff-title-bg'></div>`
    html += `<div class='diff-title-bg'></div>
        <h3 id='diff-output'>Output</h3>
        <div class='diff-title-bg'></div>
        <h3 id='diff-expected'>Expected</h3>`;
    let pos = {
        left: 1,
        right: 1
    };
    const changes = getLineChanges(hole, out, exp);
    let pendingChange = null;
    for (let i = 0; i < changes.length; i++) {
        const change = changes[i]
        pos.isLastDiff = i === changes.length - 1
        if (change.added || change.removed) {
            if (pendingChange === null) {
                pendingChange = change;
            } else {
                html += getDiffRow(hole, pendingChange, change, pos, argv);
                pendingChange = null;
            }
        } else {
            if (pendingChange) {
                html += getDiffRow(hole, pendingChange, {}, pos, argv);
                pendingChange = null;
            }
            html += getDiffLines(hole, change, change, pos, argv);
        }
    }
    if (pendingChange) {
        html += getDiffRow(hole, pendingChange, {}, pos, argv);
    }
    return html
}

function getLineChanges(hole, before, after) {
    const includeArgs = getDiffType(hole) == 'arg'
    if (includeArgs) {
        const out = []
        const splitBefore = lines(before)
        const splitAfter = lines(after)
        const compareOpts = {
            sensitivity: shouldIgnoreCase(hole) ? 'accent' : 'base'
        }
        for (let i=0; i<Math.max(splitBefore.length, splitAfter.length); i++) {
            const a = splitBefore[i]
            const b = splitAfter[i]
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

function getDiffRow(hole, change1, change2, pos, argv) {
    change2.value ??= ''
    change2.count ??= 0
    const left = change1.removed ? change1 : change2
    const right = change1.added ? change1 : change2
    return getDiffLines(hole, left, right, pos, argv)
}

function getDiffLines(hole, left, right, pos, argv) {
    const leftSplit = lines(left.value);
    const rightSplit = lines(right.value);
    if (!(pos.isLastDiff && hole === "quine")) {
        // ignore trailing newline
        if (leftSplit[leftSplit.length - 1] === '') leftSplit.pop();
        if (rightSplit[rightSplit.length - 1] === '') rightSplit.pop();
    }
    let s = ''
    const diffOpts = {
        ignoreCase: shouldIgnoreCase(hole)
    }
    for (let i=0; i<Math.max(leftSplit.length, rightSplit.length); i++) {
        const leftLine = leftSplit[i];
        const rightLine = rightSplit[i];
        const charDiff = Diff.diffChars(leftLine ?? '', rightLine ?? '', diffOpts);
        // subtract 1 because the lines start counting at 1 instead of 0
        const arg = argv[i + pos.right - 1]
        if (arg !== undefined) {
            s += `<div class='diff-arg'>${arg}</div>`
        }
        if (leftLine !== undefined) {
            s += `<div class='diff-left-num'>${i + pos.left}</div>
                <div class='diff-left${left.removed?' diff-removal':''}'>${renderCharDiff(charDiff, false)}</div>`
        }
        if (rightLine !== undefined) {
            s += `<div class='diff-right-num'>${i + pos.right}</div>
                <div class='diff-right${right.added?' diff-addition':''}'>${renderCharDiff(charDiff, true)}</div>`
        }
    }
    pos.left += left.count;
    pos.right += right.count;
    return s
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

function getDiffType(hole) {
    switch (hole) {
        case 'arabic-to-roman':
        case 'arrows':
        case 'css-colors':
        case 'emojify':
        case 'fractions':
        case 'intersection':
        case 'levenshtein-distance':
        case 'musical-chords':
        case 'ordinal-numbers':
        case 'roman-to-arabic':
        case 'spelling-numbers':
        case 'united-states':
            // { | category = Transform }
            // - {Pangram Grep, QR Decoder, Seven Segment, Morse Decoder, Morse Encoder}
            // + {Fractions, Intersection}
            return 'arg'
        default:
            return 'line'
    }
}

function shouldIgnoreCase(hole) {
    return hole === "css-colors"
}

function lines(s) {
    return s.split(/\r\n|\n/)
}
