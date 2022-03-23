import { createElement } from './_inject';
import * as Diff from 'diff';

export default (hole, exp, out, argv) => {
    if (stringsEqual(exp, out, shouldIgnoreCase(hole))) return null;

    // Show args? Exclude holes with one big argument like QR Decoder.
    const isArgDiff = argv.length > 1 && argv.length == lines(exp).length;

    const {rows, maxLineNum} = diffHTMLRows(hole, exp, out, argv, isArgDiff);

    return <table>
        {makeCols(isArgDiff, maxLineNum, argv)}
        <tr>
            {isArgDiff ? <th class="right">Args</th> : <th/>}
            <th>Output</th>
            {isArgDiff ? null : <th/>}
            <th>Expected</th>
        </tr>
        {rows}
    </table>;
};

function pushToDiff(diff, entry, join) {
    // Mutate the given diff by pushing `entry`
    // If `entry` has the same type as the previous entry, then merge them together
    const last = diff[diff.length - 1];
    if (
        last &&
        ((entry.removed && last.removed) ||
            (entry.added && last.added) ||
            (!entry.added && !last.added && !entry.removed && !last.removed))
    ) {
        // The value keeps a trailing newline when join="\n"
        last.value += entry.value + join;
        last.count += entry.count;
    }
    else {
        diff.push(entry);
    }
}

function diffWrapper(join, left, right, diffOpts) {
    // join = "\n" for line diff or "" for char diff
    // pass in left,right =  list of tokens;
    //   for char diff, this is a list of chars
    //   for line diff, this is a list of lines
    // Wrapper for performance
    // Include characters until the first difference, then include 1000 characters
    // after that, and treat the rest as a single block
    const d = firstDifference(left, right, diffOpts.ignoreCase);
    const length = Math.min(1000, Math.max(left.length - d, right.length - d));
    // Concatenate a newline on line diff because Diff.diffLines counts
    // lines without trailing newlines as changed
    const diff = (join === '' ? Diff.diffChars : Diff.diffLines)(
        left.slice(d, d + length).join(join) + join,
        right.slice(d, d + length).join(join) + join,
        diffOpts,
    );
    const head = left.slice(0, d);
    if (head !== '') {
        const fst = diff[0];
        if (fst && !fst.added && !fst.removed) {
            fst.count += head.length;
            fst.value += head.join(join) + join;
        }
        else {
            diff.unshift({
                count: head.length,
                value: head.join(join),
            });
        }
    }
    const leftTail = left.slice(d + length);
    const ltString = leftTail.join(join);
    const rightTail = right.slice(d + length);
    const rtString = rightTail.join(join);
    if (stringsEqual(ltString, rtString, diffOpts.shouldIgnoreCase)) {
        pushToDiff(
            diff,
            {
                count: leftTail.length,
                value: ltString,
            },
            join,
        );
    }
    else {
        if (ltString !== '') {
            pushToDiff(
                diff,
                {
                    added: undefined,
                    removed: true,
                    count: leftTail.length,
                    value: ltString,
                },
                join,
            );
        }
        if (rtString !== '') {
            pushToDiff(
                diff,
                {
                    added: true,
                    removed: undefined,
                    count: rightTail.length,
                    value: rtString,
                },
                join,
            );
        }
    }
    return diff;
}

function firstDifference(left, right, ignoreCase) {
    for (let i=0; i<left.length || i<right.length; i++) {
        if (!stringsEqual(left[i], right[i], ignoreCase)) {
            return i;
        }
    }
    return Math.min(left.length, right.length) + 1;
}

function makeCols(isArgDiff, maxLineNum, argv) {
    const col       = width => <col style={`width:${width}px`}/>;
    const cols      = [];
    const numLength = String(maxLineNum).length + 1;
    const charWidth = 11;

    if (isArgDiff) {
        const longestArgLength = Math.max(6, ...argv.map(arg => arg.length));
        cols.push(col(Math.min(longestArgLength * charWidth, 350)));
    }
    else {
        cols.push(col(numLength * charWidth));
    }

    cols.push(<col class="diff-col-text"/>);

    if (!isArgDiff)
        cols.push(col(numLength * charWidth));

    cols.push(<col class="diff-col-text"/>);

    return cols;
}

function diffHTMLRows(hole, exp, out, argv, isArgDiff) {
    const rows = [];
    const pos = {
        left: 1,
        right: 1,
        isLastDiff: false,
    };
    const changes = getLineChanges(hole, out, exp, isArgDiff);
    let pendingChange = null;
    for (let i = 0; i < changes.length; i++) {
        const change = changes[i];
        pos.isLastDiff = i === changes.length - 1;
        if (change.added || change.removed) {
            if (pendingChange === null) {
                pendingChange = change;
            }
            else {
                rows.push(...getDiffRow(hole, pendingChange, change, pos, argv, isArgDiff));
                pendingChange = null;
            }
        }
        else {
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

function stringsEqual(a, b, ignoreCase) {
    // https://stackoverflow.com/a/2140723/7481517
    return (
        a !== undefined &&
        a !== null &&
        0 ===
        a.localeCompare(
            b,
            undefined,
            ignoreCase
                ? {
                    sensitivity: 'accent',
                }
                : undefined,
        )
    );
}

function getLineChanges(hole, before, after, isArgDiff) {
    if (isArgDiff) {
        const out = [];
        const splitBefore = lines(before);
        const splitAfter = lines(after);
        let currentUnchanged = [];
        const ignoreCase = shouldIgnoreCase(hole);
        for (let i=0; i<Math.max(splitBefore.length, splitAfter.length); i++) {
            const a = splitBefore[i] ?? '';
            const b = splitAfter[i] ?? '';
            const linesEqual = stringsEqual(a, b, ignoreCase);
            if (linesEqual) {
                currentUnchanged.push(a);
            }
            else {
                if (currentUnchanged.length > 0) {
                    out.push({
                        count: currentUnchanged.length,
                        value: currentUnchanged.join('\n') + '\n',
                    });
                    currentUnchanged = [];
                }
                for (const [k,v] of [['removed', a], ['added', b]]) {
                    if (v !== undefined) {
                        pushToDiff(
                            out,
                            {
                                count: 1,
                                [k]: true,
                                value: v + '\n',
                            },
                            '\n',
                        );
                    }
                }
            }
        }
        if (currentUnchanged.length > 0) {
            out.push({
                count: currentUnchanged.length,
                value: currentUnchanged.join('\n') + '\n',
            });
        }
        return out;
    }
    else {
        return diffWrapper('\n', lines(before), lines(after), {
            ignoreCase: shouldIgnoreCase(hole),
        });
    }
}

function getDiffRow(hole, change1, change2, pos, argv, isArgDiff) {
    change2.value ??= '';
    change2.count ??= 0;
    const left = change1.removed ? change1 : change2;
    const right = change1.added ? change1 : change2;
    return getDiffLines(hole, left, right, pos, argv, isArgDiff);
}

function getDiffLines(hole, left, right, pos, argv, isArgDiff) {
    const leftSplit = lines(left.value);
    const rightSplit = lines(right.value);
    if (!(pos.isLastDiff && hole === 'quine')) {
        // ignore trailing newline
        if (leftSplit[leftSplit.length - 1] === '') leftSplit.pop();
        if (rightSplit[rightSplit.length - 1] === '') rightSplit.pop();
    }
    const diffOpts = {
        ignoreCase: shouldIgnoreCase(hole),
    };
    const rows = [];
    const numLines = Math.max(leftSplit.length, rightSplit.length);
    const isUnchanged = !left.removed && !right.added;
    // Skip the middle of any block of more than 7 unchanged lines, or 41 changed lines
    const padding = isUnchanged ? 3 : 20;
    const skipMiddle = numLines > 2 * padding + 1;
    for (let i=0; i<numLines; i++) {
        const row = <tr/>;
        rows.push(row);

        if (skipMiddle) {
            // In the middle; skip the line
            if (padding < i && i < numLines - padding) continue;

            // At the start of the middle; add a line saying lines omitted
            if (i === padding) {
                row.append(<td class="diff-skip" colspan={isArgDiff ? 3 : 4}>
                    {`@@ ${numLines - 2 * padding} lines omitted @@`}
                </td>);
                continue;
            }
        }

        const leftLine = leftSplit[i];
        const rightLine = rightSplit[i];
        const charDiff = diffWrapper(
            '',
            [...leftLine ?? ''],
            [...rightLine ?? ''],
            diffOpts,
        );

        // Subtract 1 because argv is 0-based and lines are 1-based.
        if (isArgDiff)
            row.append(<td class="right">{argv[i + pos.right - 1]}</td>);

        if (leftLine !== undefined) {
            if (!isArgDiff)
                row.append(<td class="diff-left-num">{i + pos.left}</td>);

            row.append(
                renderCharDiff(
                    'diff-left' + (left.removed?' diff-removal':''),
                    charDiff,
                    false,
                ),
            );
        }
        else {
            row.append(<td/>, <td/>);
        }
        if (rightLine !== undefined) {
            if (!isArgDiff)
                row.append(<td class="diff-right-num">{i + pos.right}</td>);

            row.append(
                renderCharDiff(
                    'diff-right' + (right.added ? ' diff-addition' : ''),
                    charDiff,
                    true,
                ),
            );
        }
        else {
            row.append(<td/>, <td/>);
        }
    }
    pos.left += left.count;
    pos.right += right.count;
    return rows;
}

function renderCharDiff(className, charDiff, isRight) {
    const td = <td class={className}/>;

    for (const change of charDiff)
        if (change.added && isRight)
            td.append(<span class="diff-char-addition">{change.value}</span>);
        else if (change.removed && !isRight)
            td.append(<span class="diff-char-removal">{change.value}</span>);
        else if (!change.added && !change.removed)
            td.append(change.value);

    return td;
}

function shouldIgnoreCase(hole) {
    return hole === 'css-colors';
}

function lines(s) {
    return s.split(/\r?\n/);
}
