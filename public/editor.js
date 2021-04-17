// node_modules/@codemirror/text/dist/index.js
var extend = "lc,34,7n,7,7b,19,,,,2,,2,,,20,b,1c,l,g,,2t,7,2,6,2,2,,4,z,,u,r,2j,b,1m,9,9,,o,4,,9,,3,,5,17,3,3b,f,,w,1j,,,,4,8,4,,3,7,a,2,t,,1m,,,,2,4,8,,9,,a,2,q,,2,2,1l,,4,2,4,2,2,3,3,,u,2,3,,b,2,1l,,4,5,,2,4,,k,2,m,6,,,1m,,,2,,4,8,,7,3,a,2,u,,1n,,,,c,,9,,14,,3,,1l,3,5,3,,4,7,2,b,2,t,,1m,,2,,2,,3,,5,2,7,2,b,2,s,2,1l,2,,,2,4,8,,9,,a,2,t,,20,,4,,2,3,,,8,,29,,2,7,c,8,2q,,2,9,b,6,22,2,r,,,,,,1j,e,,5,,2,5,b,,10,9,,2u,4,,6,,2,2,2,p,2,4,3,g,4,d,,2,2,6,,f,,jj,3,qa,3,t,3,t,2,u,2,1s,2,,7,8,,2,b,9,,19,3,3b,2,y,,3a,3,4,2,9,,6,3,63,2,2,,1m,,,7,,,,,2,8,6,a,2,,1c,h,1r,4,1c,7,,,5,,14,9,c,2,w,4,2,2,,3,1k,,,2,3,,,3,1m,8,2,2,48,3,,d,,7,4,,6,,3,2,5i,1m,,5,ek,,5f,x,2da,3,3x,,2o,w,fe,6,2x,2,n9w,4,,a,w,2,28,2,7k,,3,,4,,p,2,5,,47,2,q,i,d,,12,8,p,b,1a,3,1c,,2,4,2,2,13,,1v,6,2,2,2,2,c,,8,,1b,,1f,,,3,2,2,5,2,,,16,2,8,,6m,,2,,4,,fn4,,kh,g,g,g,a6,2,gt,,6a,,45,5,1ae,3,,2,5,4,14,3,4,,4l,2,fx,4,ar,2,49,b,4w,,1i,f,1k,3,1d,4,2,2,1x,3,10,5,,8,1q,,c,2,1g,9,a,4,2,,2n,3,2,,,2,6,,4g,,3,8,l,2,1l,2,,,,,m,,e,7,3,5,5f,8,2,3,,,n,,29,,2,6,,,2,,,2,,2,6j,,2,4,6,2,,2,r,2,2d,8,2,,,2,2y,,,,2,6,,,2t,3,2,4,,5,77,9,,2,6t,,a,2,,,4,,40,4,2,2,4,,w,a,14,6,2,4,8,,9,6,2,3,1a,d,,2,ba,7,,6,,,2a,m,2,7,,2,,2,3e,6,3,,,2,,7,,,20,2,3,,,,9n,2,f0b,5,1n,7,t4,,1r,4,29,,f5k,2,43q,,,3,4,5,8,8,2,7,u,4,44,3,1iz,1j,4,1e,8,,e,,m,5,,f,11s,7,,h,2,7,,2,,5,79,7,c5,4,15s,7,31,7,240,5,gx7k,2o,3k,6o".split(",").map((s) => s ? parseInt(s, 36) : 1);
for (let i = 1; i < extend.length; i++)
  extend[i] += extend[i - 1];
function isExtendingChar(code) {
  for (let i = 1; i < extend.length; i += 2)
    if (extend[i] > code)
      return extend[i - 1] <= code;
  return false;
}
function isRegionalIndicator(code) {
  return code >= 127462 && code <= 127487;
}
var ZWJ = 8205;
function findClusterBreak(str, pos, forward = true) {
  return (forward ? nextClusterBreak : prevClusterBreak)(str, pos);
}
function nextClusterBreak(str, pos) {
  if (pos == str.length)
    return pos;
  if (pos && surrogateLow(str.charCodeAt(pos)) && surrogateHigh(str.charCodeAt(pos - 1)))
    pos--;
  let prev = codePointAt(str, pos);
  pos += codePointSize(prev);
  while (pos < str.length) {
    let next = codePointAt(str, pos);
    if (prev == ZWJ || next == ZWJ || isExtendingChar(next)) {
      pos += codePointSize(next);
      prev = next;
    } else if (isRegionalIndicator(next)) {
      let countBefore = 0, i = pos - 2;
      while (i >= 0 && isRegionalIndicator(codePointAt(str, i))) {
        countBefore++;
        i -= 2;
      }
      if (countBefore % 2 == 0)
        break;
      else
        pos += 2;
    } else {
      break;
    }
  }
  return pos;
}
function prevClusterBreak(str, pos) {
  while (pos > 0) {
    let found = nextClusterBreak(str, pos - 2);
    if (found < pos)
      return found;
    pos--;
  }
  return 0;
}
function surrogateLow(ch) {
  return ch >= 56320 && ch < 57344;
}
function surrogateHigh(ch) {
  return ch >= 55296 && ch < 56320;
}
function codePointAt(str, pos) {
  let code0 = str.charCodeAt(pos);
  if (!surrogateHigh(code0) || pos + 1 == str.length)
    return code0;
  let code1 = str.charCodeAt(pos + 1);
  if (!surrogateLow(code1))
    return code0;
  return (code0 - 55296 << 10) + (code1 - 56320) + 65536;
}
function fromCodePoint(code) {
  if (code <= 65535)
    return String.fromCharCode(code);
  code -= 65536;
  return String.fromCharCode((code >> 10) + 55296, (code & 1023) + 56320);
}
function codePointSize(code) {
  return code < 65536 ? 1 : 2;
}
function countColumn(string3, n, tabSize) {
  for (let i = 0; i < string3.length; ) {
    if (string3.charCodeAt(i) == 9) {
      n += tabSize - n % tabSize;
      i++;
    } else {
      n++;
      i = findClusterBreak(string3, i);
    }
  }
  return n;
}
function findColumn(string3, n, col, tabSize) {
  for (let i = 0; i < string3.length; ) {
    if (n >= col)
      return {offset: i, leftOver: 0};
    n += string3.charCodeAt(i) == 9 ? tabSize - n % tabSize : 1;
    i = findClusterBreak(string3, i);
  }
  return {offset: string3.length, leftOver: col - n};
}
var Text = class {
  constructor() {
  }
  lineAt(pos) {
    if (pos < 0 || pos > this.length)
      throw new RangeError(`Invalid position ${pos} in document of length ${this.length}`);
    return this.lineInner(pos, false, 1, 0);
  }
  line(n) {
    if (n < 1 || n > this.lines)
      throw new RangeError(`Invalid line number ${n} in ${this.lines}-line document`);
    return this.lineInner(n, true, 1, 0);
  }
  replace(from, to, text) {
    let parts = [];
    this.decompose(0, from, parts, 2);
    if (text.length)
      text.decompose(0, text.length, parts, 1 | 2);
    this.decompose(to, this.length, parts, 1);
    return TextNode.from(parts, this.length - (to - from) + text.length);
  }
  append(other) {
    return this.replace(this.length, this.length, other);
  }
  slice(from, to = this.length) {
    let parts = [];
    this.decompose(from, to, parts, 0);
    return TextNode.from(parts, to - from);
  }
  eq(other) {
    if (other == this)
      return true;
    if (other.length != this.length || other.lines != this.lines)
      return false;
    let a = new RawTextCursor(this), b = new RawTextCursor(other);
    for (; ; ) {
      a.next();
      b.next();
      if (a.lineBreak != b.lineBreak || a.done != b.done || a.value != b.value)
        return false;
      if (a.done)
        return true;
    }
  }
  iter(dir = 1) {
    return new RawTextCursor(this, dir);
  }
  iterRange(from, to = this.length) {
    return new PartialTextCursor(this, from, to);
  }
  toString() {
    return this.sliceString(0);
  }
  toJSON() {
    let lines = [];
    this.flatten(lines);
    return lines;
  }
  static of(text) {
    if (text.length == 0)
      throw new RangeError("A document must have at least one line");
    if (text.length == 1 && !text[0])
      return Text.empty;
    return text.length <= 32 ? new TextLeaf(text) : TextNode.from(TextLeaf.split(text, []));
  }
};
if (typeof Symbol != "undefined")
  Text.prototype[Symbol.iterator] = function() {
    return this.iter();
  };
var TextLeaf = class extends Text {
  constructor(text, length = textLength(text)) {
    super();
    this.text = text;
    this.length = length;
  }
  get lines() {
    return this.text.length;
  }
  get children() {
    return null;
  }
  lineInner(target, isLine, line, offset) {
    for (let i = 0; ; i++) {
      let string3 = this.text[i], end = offset + string3.length;
      if ((isLine ? line : end) >= target)
        return new Line(offset, end, line, string3);
      offset = end + 1;
      line++;
    }
  }
  decompose(from, to, target, open) {
    let text = from <= 0 && to >= this.length ? this : new TextLeaf(sliceText(this.text, from, to), Math.min(to, this.length) - Math.max(0, from));
    if (open & 1) {
      let prev = target.pop();
      let joined = appendText(text.text, prev.text.slice(), 0, text.length);
      if (joined.length <= 32) {
        target.push(new TextLeaf(joined, prev.length + text.length));
      } else {
        let mid = joined.length >> 1;
        target.push(new TextLeaf(joined.slice(0, mid)), new TextLeaf(joined.slice(mid)));
      }
    } else {
      target.push(text);
    }
  }
  replace(from, to, text) {
    if (!(text instanceof TextLeaf))
      return super.replace(from, to, text);
    let lines = appendText(this.text, appendText(text.text, sliceText(this.text, 0, from)), to);
    let newLen = this.length + text.length - (to - from);
    if (lines.length <= 32)
      return new TextLeaf(lines, newLen);
    return TextNode.from(TextLeaf.split(lines, []), newLen);
  }
  sliceString(from, to = this.length, lineSep = "\n") {
    let result = "";
    for (let pos = 0, i = 0; pos <= to && i < this.text.length; i++) {
      let line = this.text[i], end = pos + line.length;
      if (pos > from && i)
        result += lineSep;
      if (from < end && to > pos)
        result += line.slice(Math.max(0, from - pos), to - pos);
      pos = end + 1;
    }
    return result;
  }
  flatten(target) {
    for (let line of this.text)
      target.push(line);
  }
  static split(text, target) {
    let part = [], len = -1;
    for (let line of text) {
      part.push(line);
      len += line.length + 1;
      if (part.length == 32) {
        target.push(new TextLeaf(part, len));
        part = [];
        len = -1;
      }
    }
    if (len > -1)
      target.push(new TextLeaf(part, len));
    return target;
  }
};
var TextNode = class extends Text {
  constructor(children, length) {
    super();
    this.children = children;
    this.length = length;
    this.lines = 0;
    for (let child of children)
      this.lines += child.lines;
  }
  lineInner(target, isLine, line, offset) {
    for (let i = 0; ; i++) {
      let child = this.children[i], end = offset + child.length, endLine = line + child.lines - 1;
      if ((isLine ? endLine : end) >= target)
        return child.lineInner(target, isLine, line, offset);
      offset = end + 1;
      line = endLine + 1;
    }
  }
  decompose(from, to, target, open) {
    for (let i = 0, pos = 0; pos <= to && i < this.children.length; i++) {
      let child = this.children[i], end = pos + child.length;
      if (from <= end && to >= pos) {
        let childOpen = open & ((pos <= from ? 1 : 0) | (end >= to ? 2 : 0));
        if (pos >= from && end <= to && !childOpen)
          target.push(child);
        else
          child.decompose(from - pos, to - pos, target, childOpen);
      }
      pos = end + 1;
    }
  }
  replace(from, to, text) {
    if (text.lines < this.lines)
      for (let i = 0, pos = 0; i < this.children.length; i++) {
        let child = this.children[i], end = pos + child.length;
        if (from >= pos && to <= end) {
          let updated = child.replace(from - pos, to - pos, text);
          let totalLines = this.lines - child.lines + updated.lines;
          if (updated.lines < totalLines >> 5 - 1 && updated.lines > totalLines >> 5 + 1) {
            let copy = this.children.slice();
            copy[i] = updated;
            return new TextNode(copy, this.length - (to - from) + text.length);
          }
          return super.replace(pos, end, updated);
        }
        pos = end + 1;
      }
    return super.replace(from, to, text);
  }
  sliceString(from, to = this.length, lineSep = "\n") {
    let result = "";
    for (let i = 0, pos = 0; i < this.children.length && pos <= to; i++) {
      let child = this.children[i], end = pos + child.length;
      if (pos > from && i)
        result += lineSep;
      if (from < end && to > pos)
        result += child.sliceString(from - pos, to - pos, lineSep);
      pos = end + 1;
    }
    return result;
  }
  flatten(target) {
    for (let child of this.children)
      child.flatten(target);
  }
  static from(children, length = children.reduce((l, ch) => l + ch.length + 1, -1)) {
    let lines = 0;
    for (let ch of children)
      lines += ch.lines;
    if (lines < 32) {
      let flat = [];
      for (let ch of children)
        ch.flatten(flat);
      return new TextLeaf(flat, length);
    }
    let chunk = Math.max(32, lines >> 5), maxChunk = chunk << 1, minChunk = chunk >> 1;
    let chunked = [], currentLines = 0, currentLen = -1, currentChunk = [];
    function add(child) {
      let last;
      if (child.lines > maxChunk && child instanceof TextNode) {
        for (let node of child.children)
          add(node);
      } else if (child.lines > minChunk && (currentLines > minChunk || !currentLines)) {
        flush();
        chunked.push(child);
      } else if (child instanceof TextLeaf && currentLines && (last = currentChunk[currentChunk.length - 1]) instanceof TextLeaf && child.lines + last.lines <= 32) {
        currentLines += child.lines;
        currentLen += child.length + 1;
        currentChunk[currentChunk.length - 1] = new TextLeaf(last.text.concat(child.text), last.length + 1 + child.length);
      } else {
        if (currentLines + child.lines > chunk)
          flush();
        currentLines += child.lines;
        currentLen += child.length + 1;
        currentChunk.push(child);
      }
    }
    function flush() {
      if (currentLines == 0)
        return;
      chunked.push(currentChunk.length == 1 ? currentChunk[0] : TextNode.from(currentChunk, currentLen));
      currentLen = -1;
      currentLines = currentChunk.length = 0;
    }
    for (let child of children)
      add(child);
    flush();
    return chunked.length == 1 ? chunked[0] : new TextNode(chunked, length);
  }
};
Text.empty = new TextLeaf([""], 0);
function textLength(text) {
  let length = -1;
  for (let line of text)
    length += line.length + 1;
  return length;
}
function appendText(text, target, from = 0, to = 1e9) {
  for (let pos = 0, i = 0, first = true; i < text.length && pos <= to; i++) {
    let line = text[i], end = pos + line.length;
    if (end >= from) {
      if (end > to)
        line = line.slice(0, to - pos);
      if (pos < from)
        line = line.slice(from - pos);
      if (first) {
        target[target.length - 1] += line;
        first = false;
      } else
        target.push(line);
    }
    pos = end + 1;
  }
  return target;
}
function sliceText(text, from, to) {
  return appendText(text, [""], from, to);
}
var RawTextCursor = class {
  constructor(text, dir = 1) {
    this.dir = dir;
    this.done = false;
    this.lineBreak = false;
    this.value = "";
    this.nodes = [text];
    this.offsets = [dir > 0 ? 0 : text instanceof TextLeaf ? text.text.length : text.children.length];
  }
  next(skip = 0) {
    for (; ; ) {
      let last = this.nodes.length - 1;
      if (last < 0) {
        this.done = true;
        this.value = "";
        this.lineBreak = false;
        return this;
      }
      let top2 = this.nodes[last], offset = this.offsets[last];
      let size = top2 instanceof TextLeaf ? top2.text.length : top2.children.length;
      if (offset == (this.dir > 0 ? size : 0)) {
        this.nodes.pop();
        this.offsets.pop();
      } else if (!this.lineBreak && offset != (this.dir > 0 ? 0 : size)) {
        this.lineBreak = true;
        if (skip == 0) {
          this.value = "\n";
          return this;
        }
        skip--;
      } else if (top2 instanceof TextLeaf) {
        let next = top2.text[offset - (this.dir < 0 ? 1 : 0)];
        this.offsets[last] = offset += this.dir;
        this.lineBreak = false;
        if (next.length > Math.max(0, skip)) {
          this.value = skip == 0 ? next : this.dir > 0 ? next.slice(skip) : next.slice(0, next.length - skip);
          return this;
        }
        skip -= next.length;
      } else {
        let next = top2.children[this.dir > 0 ? offset : offset - 1];
        this.offsets[last] = offset + this.dir;
        this.lineBreak = false;
        if (skip > next.length) {
          skip -= next.length;
        } else {
          this.nodes.push(next);
          this.offsets.push(this.dir > 0 ? 0 : next instanceof TextLeaf ? next.text.length : next.children.length);
        }
      }
    }
  }
};
var PartialTextCursor = class {
  constructor(text, start, end) {
    this.value = "";
    this.cursor = new RawTextCursor(text, start > end ? -1 : 1);
    if (start > end) {
      this.skip = text.length - start;
      this.limit = start - end;
    } else {
      this.skip = start;
      this.limit = end - start;
    }
  }
  next(skip = 0) {
    if (this.limit <= 0) {
      this.limit = -1;
    } else {
      let {value, lineBreak, done} = this.cursor.next(this.skip + skip);
      this.skip = 0;
      this.value = value;
      let len = lineBreak ? 1 : value.length;
      if (len > this.limit)
        this.value = this.cursor.dir > 0 ? value.slice(0, this.limit) : value.slice(len - this.limit);
      if (done || this.value.length == 0)
        this.limit = -1;
      else
        this.limit -= this.value.length;
    }
    return this;
  }
  get lineBreak() {
    return this.cursor.lineBreak;
  }
  get done() {
    return this.limit < 0;
  }
};
var Line = class {
  constructor(from, to, number2, text) {
    this.from = from;
    this.to = to;
    this.number = number2;
    this.text = text;
  }
  get length() {
    return this.to - this.from;
  }
};

// node_modules/@codemirror/state/dist/index.js
var DefaultSplit = /\r\n?|\n/;
var MapMode = /* @__PURE__ */ function(MapMode2) {
  MapMode2[MapMode2["Simple"] = 0] = "Simple";
  MapMode2[MapMode2["TrackDel"] = 1] = "TrackDel";
  MapMode2[MapMode2["TrackBefore"] = 2] = "TrackBefore";
  MapMode2[MapMode2["TrackAfter"] = 3] = "TrackAfter";
  return MapMode2;
}(MapMode || (MapMode = {}));
var ChangeDesc = class {
  constructor(sections) {
    this.sections = sections;
  }
  get length() {
    let result = 0;
    for (let i = 0; i < this.sections.length; i += 2)
      result += this.sections[i];
    return result;
  }
  get newLength() {
    let result = 0;
    for (let i = 0; i < this.sections.length; i += 2) {
      let ins = this.sections[i + 1];
      result += ins < 0 ? this.sections[i] : ins;
    }
    return result;
  }
  get empty() {
    return this.sections.length == 0 || this.sections.length == 2 && this.sections[1] < 0;
  }
  iterGaps(f) {
    for (let i = 0, posA = 0, posB = 0; i < this.sections.length; ) {
      let len = this.sections[i++], ins = this.sections[i++];
      if (ins < 0) {
        f(posA, posB, len);
        posB += len;
      } else {
        posB += ins;
      }
      posA += len;
    }
  }
  iterChangedRanges(f, individual = false) {
    iterChanges(this, f, individual);
  }
  get invertedDesc() {
    let sections = [];
    for (let i = 0; i < this.sections.length; ) {
      let len = this.sections[i++], ins = this.sections[i++];
      if (ins < 0)
        sections.push(len, ins);
      else
        sections.push(ins, len);
    }
    return new ChangeDesc(sections);
  }
  composeDesc(other) {
    return this.empty ? other : other.empty ? this : composeSets(this, other);
  }
  mapDesc(other, before = false) {
    return other.empty ? this : mapSet(this, other, before);
  }
  mapPos(pos, assoc = -1, mode = MapMode.Simple) {
    let posA = 0, posB = 0;
    for (let i = 0; i < this.sections.length; ) {
      let len = this.sections[i++], ins = this.sections[i++], endA = posA + len;
      if (ins < 0) {
        if (endA > pos)
          return posB + (pos - posA);
        posB += len;
      } else {
        if (mode != MapMode.Simple && endA >= pos && (mode == MapMode.TrackDel && posA < pos && endA > pos || mode == MapMode.TrackBefore && posA < pos || mode == MapMode.TrackAfter && endA > pos))
          return null;
        if (endA > pos || endA == pos && assoc < 0 && !len)
          return pos == posA || assoc < 0 ? posB : posB + ins;
        posB += ins;
      }
      posA = endA;
    }
    if (pos > posA)
      throw new RangeError(`Position ${pos} is out of range for changeset of length ${posA}`);
    return posB;
  }
  touchesRange(from, to = from) {
    for (let i = 0, pos = 0; i < this.sections.length && pos <= to; ) {
      let len = this.sections[i++], ins = this.sections[i++], end = pos + len;
      if (ins >= 0 && pos <= to && end >= from)
        return pos < from && end > to ? "cover" : true;
      pos = end;
    }
    return false;
  }
  toString() {
    let result = "";
    for (let i = 0; i < this.sections.length; ) {
      let len = this.sections[i++], ins = this.sections[i++];
      result += (result ? " " : "") + len + (ins >= 0 ? ":" + ins : "");
    }
    return result;
  }
  toJSON() {
    return this.sections;
  }
  static fromJSON(json) {
    if (!Array.isArray(json) || json.length % 2 || json.some((a) => typeof a != "number"))
      throw new RangeError("Invalid JSON representation of ChangeDesc");
    return new ChangeDesc(json);
  }
};
var ChangeSet = class extends ChangeDesc {
  constructor(sections, inserted) {
    super(sections);
    this.inserted = inserted;
  }
  apply(doc2) {
    if (this.length != doc2.length)
      throw new RangeError("Applying change set to a document with the wrong length");
    iterChanges(this, (fromA, toA, fromB, _toB, text) => doc2 = doc2.replace(fromB, fromB + (toA - fromA), text), false);
    return doc2;
  }
  mapDesc(other, before = false) {
    return mapSet(this, other, before, true);
  }
  invert(doc2) {
    let sections = this.sections.slice(), inserted = [];
    for (let i = 0, pos = 0; i < sections.length; i += 2) {
      let len = sections[i], ins = sections[i + 1];
      if (ins >= 0) {
        sections[i] = ins;
        sections[i + 1] = len;
        let index = i >> 1;
        while (inserted.length < index)
          inserted.push(Text.empty);
        inserted.push(len ? doc2.slice(pos, pos + len) : Text.empty);
      }
      pos += len;
    }
    return new ChangeSet(sections, inserted);
  }
  compose(other) {
    return this.empty ? other : other.empty ? this : composeSets(this, other, true);
  }
  map(other, before = false) {
    return other.empty ? this : mapSet(this, other, before, true);
  }
  iterChanges(f, individual = false) {
    iterChanges(this, f, individual);
  }
  get desc() {
    return new ChangeDesc(this.sections);
  }
  filter(ranges) {
    let resultSections = [], resultInserted = [], filteredSections = [];
    let iter = new SectionIter(this);
    done:
      for (let i = 0, pos = 0; ; ) {
        let next = i == ranges.length ? 1e9 : ranges[i++];
        while (pos < next || pos == next && iter.len == 0) {
          if (iter.done)
            break done;
          let len = Math.min(iter.len, next - pos);
          addSection(filteredSections, len, -1);
          let ins = iter.ins == -1 ? -1 : iter.off == 0 ? iter.ins : 0;
          addSection(resultSections, len, ins);
          if (ins > 0)
            addInsert(resultInserted, resultSections, iter.text);
          iter.forward(len);
          pos += len;
        }
        let end = ranges[i++];
        while (pos < end) {
          if (iter.done)
            break done;
          let len = Math.min(iter.len, end - pos);
          addSection(resultSections, len, -1);
          addSection(filteredSections, len, iter.ins == -1 ? -1 : iter.off == 0 ? iter.ins : 0);
          iter.forward(len);
          pos += len;
        }
      }
    return {
      changes: new ChangeSet(resultSections, resultInserted),
      filtered: new ChangeDesc(filteredSections)
    };
  }
  toJSON() {
    let parts = [];
    for (let i = 0; i < this.sections.length; i += 2) {
      let len = this.sections[i], ins = this.sections[i + 1];
      if (ins < 0)
        parts.push(len);
      else if (ins == 0)
        parts.push([len]);
      else
        parts.push([len].concat(this.inserted[i >> 1].toJSON()));
    }
    return parts;
  }
  static of(changes, length, lineSep) {
    let sections = [], inserted = [], pos = 0;
    let total = null;
    function flush(force = false) {
      if (!force && !sections.length)
        return;
      if (pos < length)
        addSection(sections, length - pos, -1);
      let set = new ChangeSet(sections, inserted);
      total = total ? total.compose(set.map(total)) : set;
      sections = [];
      inserted = [];
      pos = 0;
    }
    function process2(spec) {
      if (Array.isArray(spec)) {
        for (let sub of spec)
          process2(sub);
      } else if (spec instanceof ChangeSet) {
        if (spec.length != length)
          throw new RangeError(`Mismatched change set length (got ${spec.length}, expected ${length})`);
        flush();
        total = total ? total.compose(spec.map(total)) : spec;
      } else {
        let {from, to = from, insert: insert2} = spec;
        if (from > to || from < 0 || to > length)
          throw new RangeError(`Invalid change range ${from} to ${to} (in doc of length ${length})`);
        let insText = !insert2 ? Text.empty : typeof insert2 == "string" ? Text.of(insert2.split(lineSep || DefaultSplit)) : insert2;
        let insLen = insText.length;
        if (from == to && insLen == 0)
          return;
        if (from < pos)
          flush();
        if (from > pos)
          addSection(sections, from - pos, -1);
        addSection(sections, to - from, insLen);
        addInsert(inserted, sections, insText);
        pos = to;
      }
    }
    process2(changes);
    flush(!total);
    return total;
  }
  static empty(length) {
    return new ChangeSet(length ? [length, -1] : [], []);
  }
  static fromJSON(json) {
    if (!Array.isArray(json))
      throw new RangeError("Invalid JSON representation of ChangeSet");
    let sections = [], inserted = [];
    for (let i = 0; i < json.length; i++) {
      let part = json[i];
      if (typeof part == "number") {
        sections.push(part, -1);
      } else if (!Array.isArray(part) || typeof part[0] != "number" || part.some((e, i2) => i2 && typeof e != "string")) {
        throw new RangeError("Invalid JSON representation of ChangeSet");
      } else if (part.length == 1) {
        sections.push(part[0], 0);
      } else {
        while (inserted.length < i)
          inserted.push(Text.empty);
        inserted[i] = Text.of(part.slice(1));
        sections.push(part[0], inserted[i].length);
      }
    }
    return new ChangeSet(sections, inserted);
  }
};
function addSection(sections, len, ins, forceJoin = false) {
  if (len == 0 && ins <= 0)
    return;
  let last = sections.length - 2;
  if (last >= 0 && ins <= 0 && ins == sections[last + 1])
    sections[last] += len;
  else if (len == 0 && sections[last] == 0)
    sections[last + 1] += ins;
  else if (forceJoin) {
    sections[last] += len;
    sections[last + 1] += ins;
  } else
    sections.push(len, ins);
}
function addInsert(values, sections, value) {
  if (value.length == 0)
    return;
  let index = sections.length - 2 >> 1;
  if (index < values.length) {
    values[values.length - 1] = values[values.length - 1].append(value);
  } else {
    while (values.length < index)
      values.push(Text.empty);
    values.push(value);
  }
}
function iterChanges(desc, f, individual) {
  let inserted = desc.inserted;
  for (let posA = 0, posB = 0, i = 0; i < desc.sections.length; ) {
    let len = desc.sections[i++], ins = desc.sections[i++];
    if (ins < 0) {
      posA += len;
      posB += len;
    } else {
      let endA = posA, endB = posB, text = Text.empty;
      for (; ; ) {
        endA += len;
        endB += ins;
        if (ins && inserted)
          text = text.append(inserted[i - 2 >> 1]);
        if (individual || i == desc.sections.length || desc.sections[i + 1] < 0)
          break;
        len = desc.sections[i++];
        ins = desc.sections[i++];
      }
      f(posA, endA, posB, endB, text);
      posA = endA;
      posB = endB;
    }
  }
}
function mapSet(setA, setB, before, mkSet = false) {
  let sections = [], insert2 = mkSet ? [] : null;
  let a = new SectionIter(setA), b = new SectionIter(setB);
  for (let posA = 0, posB = 0; ; ) {
    if (a.ins == -1) {
      posA += a.len;
      a.next();
    } else if (b.ins == -1 && posB < posA) {
      let skip = Math.min(b.len, posA - posB);
      b.forward(skip);
      addSection(sections, skip, -1);
      posB += skip;
    } else if (b.ins >= 0 && (a.done || posB < posA || posB == posA && (b.len < a.len || b.len == a.len && !before))) {
      addSection(sections, b.ins, -1);
      while (posA > posB && !a.done && posA + a.len < posB + b.len) {
        posA += a.len;
        a.next();
      }
      posB += b.len;
      b.next();
    } else if (a.ins >= 0) {
      let len = 0, end = posA + a.len;
      for (; ; ) {
        if (b.ins >= 0 && posB > posA && posB + b.len < end) {
          len += b.ins;
          posB += b.len;
          b.next();
        } else if (b.ins == -1 && posB < end) {
          let skip = Math.min(b.len, end - posB);
          len += skip;
          b.forward(skip);
          posB += skip;
        } else {
          break;
        }
      }
      addSection(sections, len, a.ins);
      if (insert2)
        addInsert(insert2, sections, a.text);
      posA = end;
      a.next();
    } else if (a.done && b.done) {
      return insert2 ? new ChangeSet(sections, insert2) : new ChangeDesc(sections);
    } else {
      throw new Error("Mismatched change set lengths");
    }
  }
}
function composeSets(setA, setB, mkSet = false) {
  let sections = [];
  let insert2 = mkSet ? [] : null;
  let a = new SectionIter(setA), b = new SectionIter(setB);
  for (let open = false; ; ) {
    if (a.done && b.done) {
      return insert2 ? new ChangeSet(sections, insert2) : new ChangeDesc(sections);
    } else if (a.ins == 0) {
      addSection(sections, a.len, 0, open);
      a.next();
    } else if (b.len == 0 && !b.done) {
      addSection(sections, 0, b.ins, open);
      if (insert2)
        addInsert(insert2, sections, b.text);
      b.next();
    } else if (a.done || b.done) {
      throw new Error("Mismatched change set lengths");
    } else {
      let len = Math.min(a.len2, b.len), sectionLen = sections.length;
      if (a.ins == -1) {
        let insB = b.ins == -1 ? -1 : b.off ? 0 : b.ins;
        addSection(sections, len, insB, open);
        if (insert2 && insB)
          addInsert(insert2, sections, b.text);
      } else if (b.ins == -1) {
        addSection(sections, a.off ? 0 : a.len, len, open);
        if (insert2)
          addInsert(insert2, sections, a.textBit(len));
      } else {
        addSection(sections, a.off ? 0 : a.len, b.off ? 0 : b.ins, open);
        if (insert2 && !b.off)
          addInsert(insert2, sections, b.text);
      }
      open = (a.ins > len || b.ins >= 0 && b.len > len) && (open || sections.length > sectionLen);
      a.forward2(len);
      b.forward(len);
    }
  }
}
var SectionIter = class {
  constructor(set) {
    this.set = set;
    this.i = 0;
    this.next();
  }
  next() {
    let {sections} = this.set;
    if (this.i < sections.length) {
      this.len = sections[this.i++];
      this.ins = sections[this.i++];
    } else {
      this.len = 0;
      this.ins = -2;
    }
    this.off = 0;
  }
  get done() {
    return this.ins == -2;
  }
  get len2() {
    return this.ins < 0 ? this.len : this.ins;
  }
  get text() {
    let {inserted} = this.set, index = this.i - 2 >> 1;
    return index >= inserted.length ? Text.empty : inserted[index];
  }
  textBit(len) {
    let {inserted} = this.set, index = this.i - 2 >> 1;
    return index >= inserted.length && !len ? Text.empty : inserted[index].slice(this.off, len == null ? void 0 : this.off + len);
  }
  forward(len) {
    if (len == this.len)
      this.next();
    else {
      this.len -= len;
      this.off += len;
    }
  }
  forward2(len) {
    if (this.ins == -1)
      this.forward(len);
    else if (len == this.ins)
      this.next();
    else {
      this.ins -= len;
      this.off += len;
    }
  }
};
var SelectionRange = class {
  constructor(from, to, flags) {
    this.from = from;
    this.to = to;
    this.flags = flags;
  }
  get anchor() {
    return this.flags & 16 ? this.to : this.from;
  }
  get head() {
    return this.flags & 16 ? this.from : this.to;
  }
  get empty() {
    return this.from == this.to;
  }
  get assoc() {
    return this.flags & 4 ? -1 : this.flags & 8 ? 1 : 0;
  }
  get bidiLevel() {
    let level = this.flags & 3;
    return level == 3 ? null : level;
  }
  get goalColumn() {
    let value = this.flags >> 5;
    return value == 33554431 ? void 0 : value;
  }
  map(change, assoc = -1) {
    let from = change.mapPos(this.from, assoc), to = change.mapPos(this.to, assoc);
    return from == this.from && to == this.to ? this : new SelectionRange(from, to, this.flags);
  }
  extend(from, to = from) {
    if (from <= this.anchor && to >= this.anchor)
      return EditorSelection.range(from, to);
    let head = Math.abs(from - this.anchor) > Math.abs(to - this.anchor) ? from : to;
    return EditorSelection.range(this.anchor, head);
  }
  eq(other) {
    return this.anchor == other.anchor && this.head == other.head;
  }
  toJSON() {
    return {anchor: this.anchor, head: this.head};
  }
  static fromJSON(json) {
    if (!json || typeof json.anchor != "number" || typeof json.head != "number")
      throw new RangeError("Invalid JSON representation for SelectionRange");
    return EditorSelection.range(json.anchor, json.head);
  }
};
var EditorSelection = class {
  constructor(ranges, mainIndex = 0) {
    this.ranges = ranges;
    this.mainIndex = mainIndex;
  }
  map(change, assoc = -1) {
    if (change.empty)
      return this;
    return EditorSelection.create(this.ranges.map((r) => r.map(change, assoc)), this.mainIndex);
  }
  eq(other) {
    if (this.ranges.length != other.ranges.length || this.mainIndex != other.mainIndex)
      return false;
    for (let i = 0; i < this.ranges.length; i++)
      if (!this.ranges[i].eq(other.ranges[i]))
        return false;
    return true;
  }
  get main() {
    return this.ranges[this.mainIndex];
  }
  asSingle() {
    return this.ranges.length == 1 ? this : new EditorSelection([this.main]);
  }
  addRange(range, main = true) {
    return EditorSelection.create([range].concat(this.ranges), main ? 0 : this.mainIndex + 1);
  }
  replaceRange(range, which = this.mainIndex) {
    let ranges = this.ranges.slice();
    ranges[which] = range;
    return EditorSelection.create(ranges, this.mainIndex);
  }
  toJSON() {
    return {ranges: this.ranges.map((r) => r.toJSON()), main: this.mainIndex};
  }
  static fromJSON(json) {
    if (!json || !Array.isArray(json.ranges) || typeof json.main != "number" || json.main >= json.ranges.length)
      throw new RangeError("Invalid JSON representation for EditorSelection");
    return new EditorSelection(json.ranges.map((r) => SelectionRange.fromJSON(r)), json.main);
  }
  static single(anchor, head = anchor) {
    return new EditorSelection([EditorSelection.range(anchor, head)], 0);
  }
  static create(ranges, mainIndex = 0) {
    if (ranges.length == 0)
      throw new RangeError("A selection needs at least one range");
    for (let pos = 0, i = 0; i < ranges.length; i++) {
      let range = ranges[i];
      if (range.empty ? range.from <= pos : range.from < pos)
        return normalized(ranges.slice(), mainIndex);
      pos = range.to;
    }
    return new EditorSelection(ranges, mainIndex);
  }
  static cursor(pos, assoc = 0, bidiLevel, goalColumn) {
    return new SelectionRange(pos, pos, (assoc == 0 ? 0 : assoc < 0 ? 4 : 8) | (bidiLevel == null ? 3 : Math.min(2, bidiLevel)) | (goalColumn !== null && goalColumn !== void 0 ? goalColumn : 33554431) << 5);
  }
  static range(anchor, head, goalColumn) {
    let goal = (goalColumn !== null && goalColumn !== void 0 ? goalColumn : 33554431) << 5;
    return head < anchor ? new SelectionRange(head, anchor, 16 | goal) : new SelectionRange(anchor, head, goal);
  }
};
function normalized(ranges, mainIndex = 0) {
  let main = ranges[mainIndex];
  ranges.sort((a, b) => a.from - b.from);
  mainIndex = ranges.indexOf(main);
  for (let i = 1; i < ranges.length; i++) {
    let range = ranges[i], prev = ranges[i - 1];
    if (range.empty ? range.from <= prev.to : range.from < prev.to) {
      let from = prev.from, to = Math.max(range.to, prev.to);
      if (i <= mainIndex)
        mainIndex--;
      ranges.splice(--i, 2, range.anchor > range.head ? EditorSelection.range(to, from) : EditorSelection.range(from, to));
    }
  }
  return new EditorSelection(ranges, mainIndex);
}
function checkSelection(selection, docLength) {
  for (let range of selection.ranges)
    if (range.to > docLength)
      throw new RangeError("Selection points outside of document");
}
var nextID = 0;
var Facet = class {
  constructor(combine, compareInput, compare2, isStatic, extensions2) {
    this.combine = combine;
    this.compareInput = compareInput;
    this.compare = compare2;
    this.isStatic = isStatic;
    this.extensions = extensions2;
    this.id = nextID++;
    this.default = combine([]);
  }
  static define(config2 = {}) {
    return new Facet(config2.combine || ((a) => a), config2.compareInput || ((a, b) => a === b), config2.compare || (!config2.combine ? sameArray : (a, b) => a === b), !!config2.static, config2.enables);
  }
  of(value) {
    return new FacetProvider([], this, 0, value);
  }
  compute(deps, get) {
    if (this.isStatic)
      throw new Error("Can't compute a static facet");
    return new FacetProvider(deps, this, 1, get);
  }
  computeN(deps, get) {
    if (this.isStatic)
      throw new Error("Can't compute a static facet");
    return new FacetProvider(deps, this, 2, get);
  }
  from(field, get) {
    if (!get)
      get = (x) => x;
    return this.compute([field], (state) => get(state.field(field)));
  }
};
function sameArray(a, b) {
  return a == b || a.length == b.length && a.every((e, i) => e === b[i]);
}
var FacetProvider = class {
  constructor(dependencies, facet, type2, value) {
    this.dependencies = dependencies;
    this.facet = facet;
    this.type = type2;
    this.value = value;
    this.id = nextID++;
  }
  dynamicSlot(addresses) {
    var _a;
    let getter = this.value;
    let compare2 = this.facet.compareInput;
    let idx = addresses[this.id] >> 1, multi = this.type == 2;
    let depDoc = false, depSel = false, depAddrs = [];
    for (let dep of this.dependencies) {
      if (dep == "doc")
        depDoc = true;
      else if (dep == "selection")
        depSel = true;
      else if ((((_a = addresses[dep.id]) !== null && _a !== void 0 ? _a : 1) & 1) == 0)
        depAddrs.push(addresses[dep.id]);
    }
    return (state, tr) => {
      if (!tr || tr.reconfigured) {
        state.values[idx] = getter(state);
        return 1;
      } else {
        let depChanged = depDoc && tr.docChanged || depSel && (tr.docChanged || tr.selection) || depAddrs.some((addr) => (ensureAddr(state, addr) & 1) > 0);
        if (!depChanged)
          return 0;
        let newVal = getter(state), oldVal = tr.startState.values[idx];
        if (multi ? compareArray(newVal, oldVal, compare2) : compare2(newVal, oldVal))
          return 0;
        state.values[idx] = newVal;
        return 1;
      }
    };
  }
};
function compareArray(a, b, compare2) {
  if (a.length != b.length)
    return false;
  for (let i = 0; i < a.length; i++)
    if (!compare2(a[i], b[i]))
      return false;
  return true;
}
function dynamicFacetSlot(addresses, facet, providers) {
  let providerAddrs = providers.map((p) => addresses[p.id]);
  let providerTypes = providers.map((p) => p.type);
  let dynamic = providerAddrs.filter((p) => !(p & 1));
  let idx = addresses[facet.id] >> 1;
  return (state, tr) => {
    let oldAddr = !tr ? null : tr.reconfigured ? tr.startState.config.address[facet.id] : idx << 1;
    let changed = oldAddr == null;
    for (let dynAddr of dynamic) {
      if (ensureAddr(state, dynAddr) & 1)
        changed = true;
    }
    if (!changed)
      return 0;
    let values = [];
    for (let i = 0; i < providerAddrs.length; i++) {
      let value = getAddr(state, providerAddrs[i]);
      if (providerTypes[i] == 2)
        for (let val of value)
          values.push(val);
      else
        values.push(value);
    }
    let newVal = facet.combine(values);
    if (oldAddr != null && facet.compare(newVal, getAddr(tr.startState, oldAddr)))
      return 0;
    state.values[idx] = newVal;
    return 1;
  };
}
function maybeIndex(state, id2) {
  let found = state.config.address[id2];
  return found == null ? null : found >> 1;
}
var initField = /* @__PURE__ */ Facet.define({static: true});
var StateField = class {
  constructor(id2, createF, updateF, compareF, spec) {
    this.id = id2;
    this.createF = createF;
    this.updateF = updateF;
    this.compareF = compareF;
    this.spec = spec;
    this.provides = void 0;
  }
  static define(config2) {
    let field = new StateField(nextID++, config2.create, config2.update, config2.compare || ((a, b) => a === b), config2);
    if (config2.provide)
      field.provides = config2.provide(field);
    return field;
  }
  create(state) {
    let init = state.facet(initField).find((i) => i.field == this);
    return ((init === null || init === void 0 ? void 0 : init.create) || this.createF)(state);
  }
  slot(addresses) {
    let idx = addresses[this.id] >> 1;
    return (state, tr) => {
      if (!tr) {
        state.values[idx] = this.create(state);
        return 1;
      }
      let oldVal, changed = 0;
      if (tr.reconfigured) {
        let oldIdx = maybeIndex(tr.startState, this.id);
        oldVal = oldIdx == null ? this.create(tr.startState) : tr.startState.values[oldIdx];
        changed = 1;
      } else {
        oldVal = tr.startState.values[idx];
      }
      let value = this.updateF(oldVal, tr);
      if (!changed && !this.compareF(oldVal, value))
        changed = 1;
      if (changed)
        state.values[idx] = value;
      return changed;
    };
  }
  init(create) {
    return [this, initField.of({field: this, create})];
  }
  get extension() {
    return this;
  }
};
var Prec_ = {fallback: 3, default: 2, extend: 1, override: 0};
function prec(value) {
  return (ext) => new PrecExtension(ext, value);
}
var Prec = {
  fallback: /* @__PURE__ */ prec(Prec_.fallback),
  default: /* @__PURE__ */ prec(Prec_.default),
  extend: /* @__PURE__ */ prec(Prec_.extend),
  override: /* @__PURE__ */ prec(Prec_.override)
};
var PrecExtension = class {
  constructor(inner, prec2) {
    this.inner = inner;
    this.prec = prec2;
  }
};
var Compartment = class {
  of(ext) {
    return new CompartmentInstance(this, ext);
  }
  reconfigure(content2) {
    return Compartment.reconfigure.of({compartment: this, extension: content2});
  }
  get(state) {
    return state.config.compartments.get(this);
  }
};
var CompartmentInstance = class {
  constructor(compartment, inner) {
    this.compartment = compartment;
    this.inner = inner;
  }
};
var Configuration = class {
  constructor(base3, compartments, dynamicSlots, address, staticValues) {
    this.base = base3;
    this.compartments = compartments;
    this.dynamicSlots = dynamicSlots;
    this.address = address;
    this.staticValues = staticValues;
    this.statusTemplate = [];
    while (this.statusTemplate.length < dynamicSlots.length)
      this.statusTemplate.push(0);
  }
  staticFacet(facet) {
    let addr = this.address[facet.id];
    return addr == null ? facet.default : this.staticValues[addr >> 1];
  }
  static resolve(base3, compartments, oldState) {
    let fields = [];
    let facets = Object.create(null);
    let newCompartments = new Map();
    for (let ext of flatten(base3, compartments, newCompartments)) {
      if (ext instanceof StateField)
        fields.push(ext);
      else
        (facets[ext.facet.id] || (facets[ext.facet.id] = [])).push(ext);
    }
    let address = Object.create(null);
    let staticValues = [];
    let dynamicSlots = [];
    for (let field of fields) {
      address[field.id] = dynamicSlots.length << 1;
      dynamicSlots.push((a) => field.slot(a));
    }
    for (let id2 in facets) {
      let providers = facets[id2], facet = providers[0].facet;
      if (providers.every((p) => p.type == 0)) {
        address[facet.id] = staticValues.length << 1 | 1;
        let value = facet.combine(providers.map((p) => p.value));
        let oldAddr = oldState ? oldState.config.address[facet.id] : null;
        if (oldAddr != null) {
          let oldVal = getAddr(oldState, oldAddr);
          if (facet.compare(value, oldVal))
            value = oldVal;
        }
        staticValues.push(value);
      } else {
        for (let p of providers) {
          if (p.type == 0) {
            address[p.id] = staticValues.length << 1 | 1;
            staticValues.push(p.value);
          } else {
            address[p.id] = dynamicSlots.length << 1;
            dynamicSlots.push((a) => p.dynamicSlot(a));
          }
        }
        address[facet.id] = dynamicSlots.length << 1;
        dynamicSlots.push((a) => dynamicFacetSlot(a, facet, providers));
      }
    }
    return new Configuration(base3, newCompartments, dynamicSlots.map((f) => f(address)), address, staticValues);
  }
};
function flatten(extension, compartments, newCompartments) {
  let result = [[], [], [], []];
  let seen = new Map();
  function inner(ext, prec2) {
    let known = seen.get(ext);
    if (known != null) {
      if (known >= prec2)
        return;
      let found = result[known].indexOf(ext);
      if (found > -1)
        result[known].splice(found, 1);
      if (ext instanceof CompartmentInstance)
        newCompartments.delete(ext.compartment);
    }
    seen.set(ext, prec2);
    if (Array.isArray(ext)) {
      for (let e of ext)
        inner(e, prec2);
    } else if (ext instanceof CompartmentInstance) {
      if (newCompartments.has(ext.compartment))
        throw new RangeError(`Duplicate use of compartment in extensions`);
      let content2 = compartments.get(ext.compartment) || ext.inner;
      newCompartments.set(ext.compartment, content2);
      inner(content2, prec2);
    } else if (ext instanceof PrecExtension) {
      inner(ext.inner, ext.prec);
    } else if (ext instanceof StateField) {
      result[prec2].push(ext);
      if (ext.provides)
        inner(ext.provides, prec2);
    } else if (ext instanceof FacetProvider) {
      result[prec2].push(ext);
      if (ext.facet.extensions)
        inner(ext.facet.extensions, prec2);
    } else {
      let content2 = ext.extension;
      if (!content2)
        throw new Error(`Unrecognized extension value in extension set (${ext}). This sometimes happens because multiple instances of @codemirror/state are loaded, breaking instanceof checks.`);
      inner(content2, prec2);
    }
  }
  inner(extension, Prec_.default);
  return result.reduce((a, b) => a.concat(b));
}
function ensureAddr(state, addr) {
  if (addr & 1)
    return 2;
  let idx = addr >> 1;
  let status = state.status[idx];
  if (status == 4)
    throw new Error("Cyclic dependency between fields and/or facets");
  if (status & 2)
    return status;
  state.status[idx] = 4;
  let changed = state.config.dynamicSlots[idx](state, state.applying);
  return state.status[idx] = 2 | changed;
}
function getAddr(state, addr) {
  return addr & 1 ? state.config.staticValues[addr >> 1] : state.values[addr >> 1];
}
var languageData = /* @__PURE__ */ Facet.define();
var allowMultipleSelections = /* @__PURE__ */ Facet.define({
  combine: (values) => values.some((v) => v),
  static: true
});
var lineSeparator = /* @__PURE__ */ Facet.define({
  combine: (values) => values.length ? values[0] : void 0,
  static: true
});
var changeFilter = /* @__PURE__ */ Facet.define();
var transactionFilter = /* @__PURE__ */ Facet.define();
var transactionExtender = /* @__PURE__ */ Facet.define();
var Annotation = class {
  constructor(type2, value) {
    this.type = type2;
    this.value = value;
  }
  static define() {
    return new AnnotationType();
  }
};
var AnnotationType = class {
  of(value) {
    return new Annotation(this, value);
  }
};
var StateEffectType = class {
  constructor(map) {
    this.map = map;
  }
  of(value) {
    return new StateEffect(this, value);
  }
};
var StateEffect = class {
  constructor(type2, value) {
    this.type = type2;
    this.value = value;
  }
  map(mapping) {
    let mapped = this.type.map(this.value, mapping);
    return mapped === void 0 ? void 0 : mapped == this.value ? this : new StateEffect(this.type, mapped);
  }
  is(type2) {
    return this.type == type2;
  }
  static define(spec = {}) {
    return new StateEffectType(spec.map || ((v) => v));
  }
  static mapEffects(effects, mapping) {
    if (!effects.length)
      return effects;
    let result = [];
    for (let effect of effects) {
      let mapped = effect.map(mapping);
      if (mapped)
        result.push(mapped);
    }
    return result;
  }
};
StateEffect.reconfigure = /* @__PURE__ */ StateEffect.define();
StateEffect.appendConfig = /* @__PURE__ */ StateEffect.define();
var Transaction = class {
  constructor(startState, changes, selection, effects, annotations, scrollIntoView2) {
    this.startState = startState;
    this.changes = changes;
    this.selection = selection;
    this.effects = effects;
    this.annotations = annotations;
    this.scrollIntoView = scrollIntoView2;
    this._doc = null;
    this._state = null;
    if (selection)
      checkSelection(selection, changes.newLength);
    if (!annotations.some((a) => a.type == Transaction.time))
      this.annotations = annotations.concat(Transaction.time.of(Date.now()));
  }
  get newDoc() {
    return this._doc || (this._doc = this.changes.apply(this.startState.doc));
  }
  get newSelection() {
    return this.selection || this.startState.selection.map(this.changes);
  }
  get state() {
    if (!this._state)
      this.startState.applyTransaction(this);
    return this._state;
  }
  annotation(type2) {
    for (let ann of this.annotations)
      if (ann.type == type2)
        return ann.value;
    return void 0;
  }
  get docChanged() {
    return !this.changes.empty;
  }
  get reconfigured() {
    return this.startState.config != this.state.config;
  }
};
Transaction.time = /* @__PURE__ */ Annotation.define();
Transaction.userEvent = /* @__PURE__ */ Annotation.define();
Transaction.addToHistory = /* @__PURE__ */ Annotation.define();
Transaction.remote = /* @__PURE__ */ Annotation.define();
function joinRanges(a, b) {
  let result = [];
  for (let iA = 0, iB = 0; ; ) {
    let from, to;
    if (iA < a.length && (iB == b.length || b[iB] >= a[iA])) {
      from = a[iA++];
      to = a[iA++];
    } else if (iB < b.length) {
      from = b[iB++];
      to = b[iB++];
    } else
      return result;
    if (!result.length || result[result.length - 1] < from)
      result.push(from, to);
    else if (result[result.length - 1] < to)
      result[result.length - 1] = to;
  }
}
function mergeTransaction(a, b, sequential) {
  var _a;
  let mapForA, mapForB, changes;
  if (sequential) {
    mapForA = b.changes;
    mapForB = ChangeSet.empty(b.changes.length);
    changes = a.changes.compose(b.changes);
  } else {
    mapForA = b.changes.map(a.changes);
    mapForB = a.changes.mapDesc(b.changes, true);
    changes = a.changes.compose(mapForA);
  }
  return {
    changes,
    selection: b.selection ? b.selection.map(mapForB) : (_a = a.selection) === null || _a === void 0 ? void 0 : _a.map(mapForA),
    effects: StateEffect.mapEffects(a.effects, mapForA).concat(StateEffect.mapEffects(b.effects, mapForB)),
    annotations: a.annotations.length ? a.annotations.concat(b.annotations) : b.annotations,
    scrollIntoView: a.scrollIntoView || b.scrollIntoView
  };
}
function resolveTransactionInner(state, spec, docSize) {
  let sel = spec.selection;
  return {
    changes: spec.changes instanceof ChangeSet ? spec.changes : ChangeSet.of(spec.changes || [], docSize, state.facet(lineSeparator)),
    selection: sel && (sel instanceof EditorSelection ? sel : EditorSelection.single(sel.anchor, sel.head)),
    effects: asArray(spec.effects),
    annotations: asArray(spec.annotations),
    scrollIntoView: !!spec.scrollIntoView
  };
}
function resolveTransaction(state, specs, filter) {
  let s = resolveTransactionInner(state, specs.length ? specs[0] : {}, state.doc.length);
  if (specs.length && specs[0].filter === false)
    filter = false;
  for (let i = 1; i < specs.length; i++) {
    if (specs[i].filter === false)
      filter = false;
    let seq = !!specs[i].sequential;
    s = mergeTransaction(s, resolveTransactionInner(state, specs[i], seq ? s.changes.newLength : state.doc.length), seq);
  }
  let tr = new Transaction(state, s.changes, s.selection, s.effects, s.annotations, s.scrollIntoView);
  return extendTransaction(filter ? filterTransaction(tr) : tr);
}
function filterTransaction(tr) {
  let state = tr.startState;
  let result = true;
  for (let filter of state.facet(changeFilter)) {
    let value = filter(tr);
    if (value === false) {
      result = false;
      break;
    }
    if (Array.isArray(value))
      result = result === true ? value : joinRanges(result, value);
  }
  if (result !== true) {
    let changes, back;
    if (result === false) {
      back = tr.changes.invertedDesc;
      changes = ChangeSet.empty(state.doc.length);
    } else {
      let filtered = tr.changes.filter(result);
      changes = filtered.changes;
      back = filtered.filtered.invertedDesc;
    }
    tr = new Transaction(state, changes, tr.selection && tr.selection.map(back), StateEffect.mapEffects(tr.effects, back), tr.annotations, tr.scrollIntoView);
  }
  let filters = state.facet(transactionFilter);
  for (let i = filters.length - 1; i >= 0; i--) {
    let filtered = filters[i](tr);
    if (filtered instanceof Transaction)
      tr = filtered;
    else if (Array.isArray(filtered) && filtered.length == 1 && filtered[0] instanceof Transaction)
      tr = filtered[0];
    else
      tr = resolveTransaction(state, asArray(filtered), false);
  }
  return tr;
}
function extendTransaction(tr) {
  let state = tr.startState, extenders = state.facet(transactionExtender), spec = tr;
  for (let i = extenders.length - 1; i >= 0; i--) {
    let extension = extenders[i](tr);
    if (extension && Object.keys(extension).length)
      spec = mergeTransaction(tr, resolveTransactionInner(state, extension, tr.changes.newLength), true);
  }
  return spec == tr ? tr : new Transaction(state, tr.changes, tr.selection, spec.effects, spec.annotations, spec.scrollIntoView);
}
var none = [];
function asArray(value) {
  return value == null ? none : Array.isArray(value) ? value : [value];
}
var CharCategory = /* @__PURE__ */ function(CharCategory2) {
  CharCategory2[CharCategory2["Word"] = 0] = "Word";
  CharCategory2[CharCategory2["Space"] = 1] = "Space";
  CharCategory2[CharCategory2["Other"] = 2] = "Other";
  return CharCategory2;
}(CharCategory || (CharCategory = {}));
var nonASCIISingleCaseWordChar = /[\u00df\u0587\u0590-\u05f4\u0600-\u06ff\u3040-\u309f\u30a0-\u30ff\u3400-\u4db5\u4e00-\u9fcc\uac00-\ud7af]/;
var wordChar;
try {
  wordChar = /* @__PURE__ */ new RegExp("[\\p{Alphabetic}\\p{Number}_]", "u");
} catch (_) {
}
function hasWordChar(str) {
  if (wordChar)
    return wordChar.test(str);
  for (let i = 0; i < str.length; i++) {
    let ch = str[i];
    if (/\w/.test(ch) || ch > "\x80" && (ch.toUpperCase() != ch.toLowerCase() || nonASCIISingleCaseWordChar.test(ch)))
      return true;
  }
  return false;
}
function makeCategorizer(wordChars) {
  return (char) => {
    if (!/\S/.test(char))
      return CharCategory.Space;
    if (hasWordChar(char))
      return CharCategory.Word;
    for (let i = 0; i < wordChars.length; i++)
      if (char.indexOf(wordChars[i]) > -1)
        return CharCategory.Word;
    return CharCategory.Other;
  };
}
var EditorState = class {
  constructor(config2, doc2, selection, tr = null) {
    this.config = config2;
    this.doc = doc2;
    this.selection = selection;
    this.applying = null;
    this.status = config2.statusTemplate.slice();
    if (tr && tr.startState.config == config2) {
      this.values = tr.startState.values.slice();
    } else {
      this.values = config2.dynamicSlots.map((_) => null);
      if (tr)
        for (let id2 in config2.address) {
          let cur2 = config2.address[id2], prev = tr.startState.config.address[id2];
          if (prev != null && (cur2 & 1) == 0)
            this.values[cur2 >> 1] = getAddr(tr.startState, prev);
        }
    }
    this.applying = tr;
    if (tr)
      tr._state = this;
    for (let i = 0; i < this.config.dynamicSlots.length; i++)
      ensureAddr(this, i << 1);
    this.applying = null;
  }
  field(field, require2 = true) {
    let addr = this.config.address[field.id];
    if (addr == null) {
      if (require2)
        throw new RangeError("Field is not present in this state");
      return void 0;
    }
    ensureAddr(this, addr);
    return getAddr(this, addr);
  }
  update(...specs) {
    return resolveTransaction(this, specs, true);
  }
  applyTransaction(tr) {
    let conf = this.config, {base: base3, compartments} = conf;
    for (let effect of tr.effects) {
      if (effect.is(Compartment.reconfigure)) {
        if (conf) {
          compartments = new Map();
          conf.compartments.forEach((val, key) => compartments.set(key, val));
          conf = null;
        }
        compartments.set(effect.value.compartment, effect.value.extension);
      } else if (effect.is(StateEffect.reconfigure)) {
        conf = null;
        base3 = effect.value;
      } else if (effect.is(StateEffect.appendConfig)) {
        conf = null;
        base3 = asArray(base3).concat(effect.value);
      }
    }
    new EditorState(conf || Configuration.resolve(base3, compartments, this), tr.newDoc, tr.newSelection, tr);
  }
  replaceSelection(text) {
    if (typeof text == "string")
      text = this.toText(text);
    return this.changeByRange((range) => ({
      changes: {from: range.from, to: range.to, insert: text},
      range: EditorSelection.cursor(range.from + text.length)
    }));
  }
  changeByRange(f) {
    let sel = this.selection;
    let result1 = f(sel.ranges[0]);
    let changes = this.changes(result1.changes), ranges = [result1.range];
    let effects = asArray(result1.effects);
    for (let i = 1; i < sel.ranges.length; i++) {
      let result = f(sel.ranges[i]);
      let newChanges = this.changes(result.changes), newMapped = newChanges.map(changes);
      for (let j = 0; j < i; j++)
        ranges[j] = ranges[j].map(newMapped);
      let mapBy = changes.mapDesc(newChanges, true);
      ranges.push(result.range.map(mapBy));
      changes = changes.compose(newMapped);
      effects = StateEffect.mapEffects(effects, newMapped).concat(StateEffect.mapEffects(asArray(result.effects), mapBy));
    }
    return {
      changes,
      selection: EditorSelection.create(ranges, sel.mainIndex),
      effects
    };
  }
  changes(spec = []) {
    if (spec instanceof ChangeSet)
      return spec;
    return ChangeSet.of(spec, this.doc.length, this.facet(EditorState.lineSeparator));
  }
  toText(string3) {
    return Text.of(string3.split(this.facet(EditorState.lineSeparator) || DefaultSplit));
  }
  sliceDoc(from = 0, to = this.doc.length) {
    return this.doc.sliceString(from, to, this.lineBreak);
  }
  facet(facet) {
    let addr = this.config.address[facet.id];
    if (addr == null)
      return facet.default;
    ensureAddr(this, addr);
    return getAddr(this, addr);
  }
  toJSON(fields) {
    let result = {
      doc: this.sliceDoc(),
      selection: this.selection.toJSON()
    };
    if (fields)
      for (let prop in fields) {
        let value = fields[prop];
        if (value instanceof StateField)
          result[prop] = value.spec.toJSON(this.field(fields[prop]), this);
      }
    return result;
  }
  static fromJSON(json, config2 = {}, fields) {
    if (!json || typeof json.doc != "string")
      throw new RangeError("Invalid JSON representation for EditorState");
    let fieldInit = [];
    if (fields)
      for (let prop in fields) {
        let field = fields[prop], value = json[prop];
        fieldInit.push(field.init((state) => field.spec.fromJSON(value, state)));
      }
    return EditorState.create({
      doc: json.doc,
      selection: EditorSelection.fromJSON(json.selection),
      extensions: config2.extensions ? fieldInit.concat([config2.extensions]) : fieldInit
    });
  }
  static create(config2 = {}) {
    let configuration = Configuration.resolve(config2.extensions || [], new Map());
    let doc2 = config2.doc instanceof Text ? config2.doc : Text.of((config2.doc || "").split(configuration.staticFacet(EditorState.lineSeparator) || DefaultSplit));
    let selection = !config2.selection ? EditorSelection.single(0) : config2.selection instanceof EditorSelection ? config2.selection : EditorSelection.single(config2.selection.anchor, config2.selection.head);
    checkSelection(selection, doc2.length);
    if (!configuration.staticFacet(allowMultipleSelections))
      selection = selection.asSingle();
    return new EditorState(configuration, doc2, selection);
  }
  get tabSize() {
    return this.facet(EditorState.tabSize);
  }
  get lineBreak() {
    return this.facet(EditorState.lineSeparator) || "\n";
  }
  phrase(phrase) {
    for (let map of this.facet(EditorState.phrases))
      if (Object.prototype.hasOwnProperty.call(map, phrase))
        return map[phrase];
    return phrase;
  }
  languageDataAt(name2, pos) {
    let values = [];
    for (let provider of this.facet(languageData)) {
      for (let result of provider(this, pos)) {
        if (Object.prototype.hasOwnProperty.call(result, name2))
          values.push(result[name2]);
      }
    }
    return values;
  }
  charCategorizer(at) {
    return makeCategorizer(this.languageDataAt("wordChars", at).join(""));
  }
  wordAt(pos) {
    let {text, from, length} = this.doc.lineAt(pos);
    let cat = this.charCategorizer(pos);
    let start = pos - from, end = pos - from;
    while (start > 0) {
      let prev = findClusterBreak(text, start, false);
      if (cat(text.slice(prev, start)) != CharCategory.Word)
        break;
      start = prev;
    }
    while (end < length) {
      let next = findClusterBreak(text, end);
      if (cat(text.slice(end, next)) != CharCategory.Word)
        break;
      end = next;
    }
    return start == end ? EditorSelection.range(start + from, end + from) : null;
  }
};
EditorState.allowMultipleSelections = allowMultipleSelections;
EditorState.tabSize = /* @__PURE__ */ Facet.define({
  combine: (values) => values.length ? values[0] : 4
});
EditorState.lineSeparator = lineSeparator;
EditorState.phrases = /* @__PURE__ */ Facet.define();
EditorState.languageData = languageData;
EditorState.changeFilter = changeFilter;
EditorState.transactionFilter = transactionFilter;
EditorState.transactionExtender = transactionExtender;
Compartment.reconfigure = /* @__PURE__ */ StateEffect.define();
function combineConfig(configs, defaults4, combine = {}) {
  let result = {};
  for (let config2 of configs)
    for (let key of Object.keys(config2)) {
      let value = config2[key], current = result[key];
      if (current === void 0)
        result[key] = value;
      else if (current === value || value === void 0)
        ;
      else if (Object.hasOwnProperty.call(combine, key))
        result[key] = combine[key](current, value);
      else
        throw new Error("Config merge conflict for field " + key);
    }
  for (let key in defaults4)
    if (result[key] === void 0)
      result[key] = defaults4[key];
  return result;
}

// node_modules/style-mod/src/style-mod.js
var C = "\u037C";
var COUNT = typeof Symbol == "undefined" ? "__" + C : Symbol.for(C);
var SET = typeof Symbol == "undefined" ? "__styleSet" + Math.floor(Math.random() * 1e8) : Symbol("styleSet");
var top = typeof globalThis != "undefined" ? globalThis : typeof window != "undefined" ? window : {};
var StyleModule = class {
  constructor(spec, options) {
    this.rules = [];
    let {finish} = options || {};
    function splitSelector(selector) {
      return /^@/.test(selector) ? [selector] : selector.split(/,\s*/);
    }
    function render(selectors, spec2, target, isKeyframes) {
      let local = [], isAt = /^@(\w+)\b/.exec(selectors[0]), keyframes = isAt && isAt[1] == "keyframes";
      if (isAt && spec2 == null)
        return target.push(selectors[0] + ";");
      for (let prop in spec2) {
        let value = spec2[prop];
        if (/&/.test(prop)) {
          render(prop.split(/,\s*/).map((part) => selectors.map((sel) => part.replace(/&/, sel))).reduce((a, b) => a.concat(b)), value, target);
        } else if (value && typeof value == "object") {
          if (!isAt)
            throw new RangeError("The value of a property (" + prop + ") should be a primitive value.");
          render(splitSelector(prop), value, local, keyframes);
        } else if (value != null) {
          local.push(prop.replace(/_.*/, "").replace(/[A-Z]/g, (l) => "-" + l.toLowerCase()) + ": " + value + ";");
        }
      }
      if (local.length || keyframes) {
        target.push((finish && !isAt && !isKeyframes ? selectors.map(finish) : selectors).join(", ") + " {" + local.join(" ") + "}");
      }
    }
    for (let prop in spec)
      render(splitSelector(prop), spec[prop], this.rules);
  }
  getRules() {
    return this.rules.join("\n");
  }
  static newName() {
    let id2 = top[COUNT] || 1;
    top[COUNT] = id2 + 1;
    return C + id2.toString(36);
  }
  static mount(root, modules) {
    (root[SET] || new StyleSet(root)).mount(Array.isArray(modules) ? modules : [modules]);
  }
};
var adoptedSet = null;
var StyleSet = class {
  constructor(root) {
    if (!root.head && root.adoptedStyleSheets && typeof CSSStyleSheet != "undefined") {
      if (adoptedSet) {
        root.adoptedStyleSheets = [adoptedSet.sheet].concat(root.adoptedStyleSheets);
        return root[SET] = adoptedSet;
      }
      this.sheet = new CSSStyleSheet();
      root.adoptedStyleSheets = [this.sheet].concat(root.adoptedStyleSheets);
      adoptedSet = this;
    } else {
      this.styleTag = (root.ownerDocument || root).createElement("style");
      let target = root.head || root;
      target.insertBefore(this.styleTag, target.firstChild);
    }
    this.modules = [];
    root[SET] = this;
  }
  mount(modules) {
    let sheet = this.sheet;
    let pos = 0, j = 0;
    for (let i = 0; i < modules.length; i++) {
      let mod = modules[i], index = this.modules.indexOf(mod);
      if (index < j && index > -1) {
        this.modules.splice(index, 1);
        j--;
        index = -1;
      }
      if (index == -1) {
        this.modules.splice(j++, 0, mod);
        if (sheet)
          for (let k = 0; k < mod.rules.length; k++)
            sheet.insertRule(mod.rules[k], pos++);
      } else {
        while (j < index)
          pos += this.modules[j++].rules.length;
        pos += mod.rules.length;
        j++;
      }
    }
    if (!sheet) {
      let text = "";
      for (let i = 0; i < this.modules.length; i++)
        text += this.modules[i].getRules() + "\n";
      this.styleTag.textContent = text;
    }
  }
};

// node_modules/@codemirror/rangeset/dist/index.js
var RangeValue = class {
  eq(other) {
    return this == other;
  }
  range(from, to = from) {
    return new Range(from, to, this);
  }
};
RangeValue.prototype.startSide = RangeValue.prototype.endSide = 0;
RangeValue.prototype.point = false;
RangeValue.prototype.mapMode = MapMode.TrackDel;
var Range = class {
  constructor(from, to, value) {
    this.from = from;
    this.to = to;
    this.value = value;
  }
};
function cmpRange(a, b) {
  return a.from - b.from || a.value.startSide - b.value.startSide;
}
var Chunk = class {
  constructor(from, to, value, maxPoint) {
    this.from = from;
    this.to = to;
    this.value = value;
    this.maxPoint = maxPoint;
  }
  get length() {
    return this.to[this.to.length - 1];
  }
  findIndex(pos, end, side = end * 1e9, startAt = 0) {
    if (pos <= 0)
      return startAt;
    let arr = end < 0 ? this.to : this.from;
    for (let lo = startAt, hi = arr.length; ; ) {
      if (lo == hi)
        return lo;
      let mid = lo + hi >> 1;
      let diff = arr[mid] - pos || (end < 0 ? this.value[mid].startSide : this.value[mid].endSide) - side;
      if (mid == lo)
        return diff >= 0 ? lo : hi;
      if (diff >= 0)
        hi = mid;
      else
        lo = mid + 1;
    }
  }
  between(offset, from, to, f) {
    for (let i = this.findIndex(from, -1), e = this.findIndex(to, 1, void 0, i); i < e; i++)
      if (f(this.from[i] + offset, this.to[i] + offset, this.value[i]) === false)
        return false;
  }
  map(offset, changes) {
    let value = [], from = [], to = [], newPos = -1, maxPoint = -1;
    for (let i = 0; i < this.value.length; i++) {
      let val = this.value[i], curFrom = this.from[i] + offset, curTo = this.to[i] + offset, newFrom, newTo;
      if (curFrom == curTo) {
        let mapped = changes.mapPos(curFrom, val.startSide, val.mapMode);
        if (mapped == null)
          continue;
        newFrom = newTo = mapped;
      } else {
        newFrom = changes.mapPos(curFrom, val.startSide);
        newTo = changes.mapPos(curTo, val.endSide);
        if (newFrom > newTo || newFrom == newTo && val.startSide > 0 && val.endSide <= 0)
          continue;
      }
      if ((newTo - newFrom || val.endSide - val.startSide) < 0)
        continue;
      if (newPos < 0)
        newPos = newFrom;
      if (val.point)
        maxPoint = Math.max(maxPoint, newTo - newFrom);
      value.push(val);
      from.push(newFrom - newPos);
      to.push(newTo - newPos);
    }
    return {mapped: value.length ? new Chunk(from, to, value, maxPoint) : null, pos: newPos};
  }
};
var RangeSet = class {
  constructor(chunkPos, chunk, nextLayer = RangeSet.empty, maxPoint) {
    this.chunkPos = chunkPos;
    this.chunk = chunk;
    this.nextLayer = nextLayer;
    this.maxPoint = maxPoint;
  }
  get length() {
    let last = this.chunk.length - 1;
    return last < 0 ? 0 : Math.max(this.chunkEnd(last), this.nextLayer.length);
  }
  get size() {
    if (this == RangeSet.empty)
      return 0;
    let size = this.nextLayer.size;
    for (let chunk of this.chunk)
      size += chunk.value.length;
    return size;
  }
  chunkEnd(index) {
    return this.chunkPos[index] + this.chunk[index].length;
  }
  update(updateSpec) {
    let {add = [], sort = false, filterFrom = 0, filterTo = this.length} = updateSpec;
    let filter = updateSpec.filter;
    if (add.length == 0 && !filter)
      return this;
    if (sort)
      add.slice().sort(cmpRange);
    if (this == RangeSet.empty)
      return add.length ? RangeSet.of(add) : this;
    let cur2 = new LayerCursor(this, null, -1).goto(0), i = 0, spill = [];
    let builder = new RangeSetBuilder();
    while (cur2.value || i < add.length) {
      if (i < add.length && (cur2.from - add[i].from || cur2.startSide - add[i].value.startSide) >= 0) {
        let range = add[i++];
        if (!builder.addInner(range.from, range.to, range.value))
          spill.push(range);
      } else if (cur2.rangeIndex == 1 && cur2.chunkIndex < this.chunk.length && (i == add.length || this.chunkEnd(cur2.chunkIndex) < add[i].from) && (!filter || filterFrom > this.chunkEnd(cur2.chunkIndex) || filterTo < this.chunkPos[cur2.chunkIndex]) && builder.addChunk(this.chunkPos[cur2.chunkIndex], this.chunk[cur2.chunkIndex])) {
        cur2.nextChunk();
      } else {
        if (!filter || filterFrom > cur2.to || filterTo < cur2.from || filter(cur2.from, cur2.to, cur2.value)) {
          if (!builder.addInner(cur2.from, cur2.to, cur2.value))
            spill.push(new Range(cur2.from, cur2.to, cur2.value));
        }
        cur2.next();
      }
    }
    return builder.finishInner(this.nextLayer == RangeSet.empty && !spill.length ? RangeSet.empty : this.nextLayer.update({add: spill, filter, filterFrom, filterTo}));
  }
  map(changes) {
    if (changes.length == 0 || this == RangeSet.empty)
      return this;
    let chunks = [], chunkPos = [], maxPoint = -1;
    for (let i = 0; i < this.chunk.length; i++) {
      let start = this.chunkPos[i], chunk = this.chunk[i];
      let touch = changes.touchesRange(start, start + chunk.length);
      if (touch === false) {
        maxPoint = Math.max(maxPoint, chunk.maxPoint);
        chunks.push(chunk);
        chunkPos.push(changes.mapPos(start));
      } else if (touch === true) {
        let {mapped, pos} = chunk.map(start, changes);
        if (mapped) {
          maxPoint = Math.max(maxPoint, mapped.maxPoint);
          chunks.push(mapped);
          chunkPos.push(pos);
        }
      }
    }
    let next = this.nextLayer.map(changes);
    return chunks.length == 0 ? next : new RangeSet(chunkPos, chunks, next, maxPoint);
  }
  between(from, to, f) {
    if (this == RangeSet.empty)
      return;
    for (let i = 0; i < this.chunk.length; i++) {
      let start = this.chunkPos[i], chunk = this.chunk[i];
      if (to >= start && from <= start + chunk.length && chunk.between(start, from - start, to - start, f) === false)
        return;
    }
    this.nextLayer.between(from, to, f);
  }
  iter(from = 0) {
    return HeapCursor.from([this]).goto(from);
  }
  static iter(sets, from = 0) {
    return HeapCursor.from(sets).goto(from);
  }
  static compare(oldSets, newSets, textDiff, comparator, minPointSize = -1) {
    let a = oldSets.filter((set) => set.maxPoint >= 500 || set != RangeSet.empty && newSets.indexOf(set) < 0 && set.maxPoint >= minPointSize);
    let b = newSets.filter((set) => set.maxPoint >= 500 || set != RangeSet.empty && oldSets.indexOf(set) < 0 && set.maxPoint >= minPointSize);
    let sharedChunks = findSharedChunks(a, b);
    let sideA = new SpanCursor(a, sharedChunks, minPointSize);
    let sideB = new SpanCursor(b, sharedChunks, minPointSize);
    textDiff.iterGaps((fromA, fromB, length) => compare(sideA, fromA, sideB, fromB, length, comparator));
    if (textDiff.empty && textDiff.length == 0)
      compare(sideA, 0, sideB, 0, 0, comparator);
  }
  static spans(sets, from, to, iterator, minPointSize = -1) {
    let cursor = new SpanCursor(sets, null, minPointSize).goto(from), pos = from;
    let open = cursor.openStart;
    for (; ; ) {
      let curTo = Math.min(cursor.to, to);
      if (cursor.point) {
        iterator.point(pos, curTo, cursor.point, cursor.activeForPoint(cursor.to), open);
        open = cursor.openEnd(curTo) + (cursor.to > curTo ? 1 : 0);
      } else if (curTo > pos) {
        iterator.span(pos, curTo, cursor.active, open);
        open = cursor.openEnd(curTo);
      }
      if (cursor.to > to)
        break;
      pos = cursor.to;
      cursor.next();
    }
    return open;
  }
  static of(ranges, sort = false) {
    let build = new RangeSetBuilder();
    for (let range of ranges instanceof Range ? [ranges] : sort ? ranges.slice().sort(cmpRange) : ranges)
      build.add(range.from, range.to, range.value);
    return build.finish();
  }
};
RangeSet.empty = new RangeSet([], [], null, -1);
RangeSet.empty.nextLayer = RangeSet.empty;
var RangeSetBuilder = class {
  constructor() {
    this.chunks = [];
    this.chunkPos = [];
    this.chunkStart = -1;
    this.last = null;
    this.lastFrom = -1e9;
    this.lastTo = -1e9;
    this.from = [];
    this.to = [];
    this.value = [];
    this.maxPoint = -1;
    this.setMaxPoint = -1;
    this.nextLayer = null;
  }
  finishChunk(newArrays) {
    this.chunks.push(new Chunk(this.from, this.to, this.value, this.maxPoint));
    this.chunkPos.push(this.chunkStart);
    this.chunkStart = -1;
    this.setMaxPoint = Math.max(this.setMaxPoint, this.maxPoint);
    this.maxPoint = -1;
    if (newArrays) {
      this.from = [];
      this.to = [];
      this.value = [];
    }
  }
  add(from, to, value) {
    if (!this.addInner(from, to, value))
      (this.nextLayer || (this.nextLayer = new RangeSetBuilder())).add(from, to, value);
  }
  addInner(from, to, value) {
    let diff = from - this.lastTo || value.startSide - this.last.endSide;
    if (diff <= 0 && (from - this.lastFrom || value.startSide - this.last.startSide) < 0)
      throw new Error("Ranges must be added sorted by `from` position and `startSide`");
    if (diff < 0)
      return false;
    if (this.from.length == 250)
      this.finishChunk(true);
    if (this.chunkStart < 0)
      this.chunkStart = from;
    this.from.push(from - this.chunkStart);
    this.to.push(to - this.chunkStart);
    this.last = value;
    this.lastFrom = from;
    this.lastTo = to;
    this.value.push(value);
    if (value.point)
      this.maxPoint = Math.max(this.maxPoint, to - from);
    return true;
  }
  addChunk(from, chunk) {
    if ((from - this.lastTo || chunk.value[0].startSide - this.last.endSide) < 0)
      return false;
    if (this.from.length)
      this.finishChunk(true);
    this.setMaxPoint = Math.max(this.setMaxPoint, chunk.maxPoint);
    this.chunks.push(chunk);
    this.chunkPos.push(from);
    let last = chunk.value.length - 1;
    this.last = chunk.value[last];
    this.lastFrom = chunk.from[last] + from;
    this.lastTo = chunk.to[last] + from;
    return true;
  }
  finish() {
    return this.finishInner(RangeSet.empty);
  }
  finishInner(next) {
    if (this.from.length)
      this.finishChunk(false);
    if (this.chunks.length == 0)
      return next;
    let result = new RangeSet(this.chunkPos, this.chunks, this.nextLayer ? this.nextLayer.finishInner(next) : next, this.setMaxPoint);
    this.from = null;
    return result;
  }
};
function findSharedChunks(a, b) {
  let inA = new Map();
  for (let set of a)
    for (let i = 0; i < set.chunk.length; i++)
      if (set.chunk[i].maxPoint < 500)
        inA.set(set.chunk[i], set.chunkPos[i]);
  let shared = new Set();
  for (let set of b)
    for (let i = 0; i < set.chunk.length; i++)
      if (inA.get(set.chunk[i]) == set.chunkPos[i])
        shared.add(set.chunk[i]);
  return shared;
}
var LayerCursor = class {
  constructor(layer, skip, minPoint, rank = 0) {
    this.layer = layer;
    this.skip = skip;
    this.minPoint = minPoint;
    this.rank = rank;
  }
  get startSide() {
    return this.value ? this.value.startSide : 0;
  }
  get endSide() {
    return this.value ? this.value.endSide : 0;
  }
  goto(pos, side = -1e9) {
    this.chunkIndex = this.rangeIndex = 0;
    this.gotoInner(pos, side, false);
    return this;
  }
  gotoInner(pos, side, forward) {
    while (this.chunkIndex < this.layer.chunk.length) {
      let next = this.layer.chunk[this.chunkIndex];
      if (!(this.skip && this.skip.has(next) || this.layer.chunkEnd(this.chunkIndex) < pos || next.maxPoint < this.minPoint))
        break;
      this.chunkIndex++;
      forward = false;
    }
    let rangeIndex = this.chunkIndex == this.layer.chunk.length ? 0 : this.layer.chunk[this.chunkIndex].findIndex(pos - this.layer.chunkPos[this.chunkIndex], -1, side);
    if (!forward || this.rangeIndex < rangeIndex)
      this.rangeIndex = rangeIndex;
    this.next();
  }
  forward(pos, side) {
    if ((this.to - pos || this.endSide - side) < 0)
      this.gotoInner(pos, side, true);
  }
  next() {
    for (; ; ) {
      if (this.chunkIndex == this.layer.chunk.length) {
        this.from = this.to = 1e9;
        this.value = null;
        break;
      } else {
        let chunkPos = this.layer.chunkPos[this.chunkIndex], chunk = this.layer.chunk[this.chunkIndex];
        let from = chunkPos + chunk.from[this.rangeIndex];
        this.from = from;
        this.to = chunkPos + chunk.to[this.rangeIndex];
        this.value = chunk.value[this.rangeIndex];
        if (++this.rangeIndex == chunk.value.length) {
          this.chunkIndex++;
          if (this.skip) {
            while (this.chunkIndex < this.layer.chunk.length && this.skip.has(this.layer.chunk[this.chunkIndex]))
              this.chunkIndex++;
          }
          this.rangeIndex = 0;
        }
        if (this.minPoint < 0 || this.value.point && this.to - this.from >= this.minPoint)
          break;
      }
    }
  }
  nextChunk() {
    this.chunkIndex++;
    this.rangeIndex = 0;
    this.next();
  }
  compare(other) {
    return this.from - other.from || this.startSide - other.startSide || this.to - other.to || this.endSide - other.endSide;
  }
};
var HeapCursor = class {
  constructor(heap) {
    this.heap = heap;
  }
  static from(sets, skip = null, minPoint = -1) {
    let heap = [];
    for (let i = 0; i < sets.length; i++) {
      for (let cur2 = sets[i]; cur2 != RangeSet.empty; cur2 = cur2.nextLayer) {
        if (cur2.maxPoint >= minPoint)
          heap.push(new LayerCursor(cur2, skip, minPoint, i));
      }
    }
    return heap.length == 1 ? heap[0] : new HeapCursor(heap);
  }
  get startSide() {
    return this.value ? this.value.startSide : 0;
  }
  goto(pos, side = -1e9) {
    for (let cur2 of this.heap)
      cur2.goto(pos, side);
    for (let i = this.heap.length >> 1; i >= 0; i--)
      heapBubble(this.heap, i);
    this.next();
    return this;
  }
  forward(pos, side) {
    for (let cur2 of this.heap)
      cur2.forward(pos, side);
    for (let i = this.heap.length >> 1; i >= 0; i--)
      heapBubble(this.heap, i);
    if ((this.to - pos || this.value.endSide - side) < 0)
      this.next();
  }
  next() {
    if (this.heap.length == 0) {
      this.from = this.to = 1e9;
      this.value = null;
      this.rank = -1;
    } else {
      let top2 = this.heap[0];
      this.from = top2.from;
      this.to = top2.to;
      this.value = top2.value;
      this.rank = top2.rank;
      if (top2.value)
        top2.next();
      heapBubble(this.heap, 0);
    }
  }
};
function heapBubble(heap, index) {
  for (let cur2 = heap[index]; ; ) {
    let childIndex = (index << 1) + 1;
    if (childIndex >= heap.length)
      break;
    let child = heap[childIndex];
    if (childIndex + 1 < heap.length && child.compare(heap[childIndex + 1]) >= 0) {
      child = heap[childIndex + 1];
      childIndex++;
    }
    if (cur2.compare(child) < 0)
      break;
    heap[childIndex] = cur2;
    heap[index] = child;
    index = childIndex;
  }
}
var SpanCursor = class {
  constructor(sets, skip, minPoint) {
    this.minPoint = minPoint;
    this.active = [];
    this.activeTo = [];
    this.activeRank = [];
    this.minActive = -1;
    this.point = null;
    this.pointFrom = 0;
    this.pointRank = 0;
    this.to = -1e9;
    this.endSide = 0;
    this.openStart = -1;
    this.cursor = HeapCursor.from(sets, skip, minPoint);
  }
  goto(pos, side = -1e9) {
    this.cursor.goto(pos, side);
    this.active.length = this.activeTo.length = this.activeRank.length = 0;
    this.minActive = -1;
    this.to = pos;
    this.endSide = side;
    this.openStart = -1;
    this.next();
    return this;
  }
  forward(pos, side) {
    while (this.minActive > -1 && (this.activeTo[this.minActive] - pos || this.active[this.minActive].endSide - side) < 0)
      this.removeActive(this.minActive);
    this.cursor.forward(pos, side);
  }
  removeActive(index) {
    remove(this.active, index);
    remove(this.activeTo, index);
    remove(this.activeRank, index);
    this.minActive = findMinIndex(this.active, this.activeTo);
  }
  addActive(trackOpen) {
    let i = 0, {value, to, rank} = this.cursor;
    while (i < this.activeRank.length && this.activeRank[i] <= rank)
      i++;
    insert(this.active, i, value);
    insert(this.activeTo, i, to);
    insert(this.activeRank, i, rank);
    if (trackOpen)
      insert(trackOpen, i, this.cursor.from);
    this.minActive = findMinIndex(this.active, this.activeTo);
  }
  next() {
    let from = this.to;
    this.point = null;
    let trackOpen = this.openStart < 0 ? [] : null, trackExtra = 0;
    for (; ; ) {
      let a = this.minActive;
      if (a > -1 && (this.activeTo[a] - this.cursor.from || this.active[a].endSide - this.cursor.startSide) < 0) {
        if (this.activeTo[a] > from) {
          this.to = this.activeTo[a];
          this.endSide = this.active[a].endSide;
          break;
        }
        this.removeActive(a);
        if (trackOpen)
          remove(trackOpen, a);
      } else if (!this.cursor.value) {
        this.to = this.endSide = 1e9;
        break;
      } else if (this.cursor.from > from) {
        this.to = this.cursor.from;
        this.endSide = this.cursor.startSide;
        break;
      } else {
        let nextVal = this.cursor.value;
        if (!nextVal.point) {
          this.addActive(trackOpen);
          this.cursor.next();
        } else {
          this.point = nextVal;
          this.pointFrom = this.cursor.from;
          this.pointRank = this.cursor.rank;
          this.to = this.cursor.to;
          this.endSide = nextVal.endSide;
          if (this.cursor.from < from)
            trackExtra = 1;
          this.cursor.next();
          if (this.to > from)
            this.forward(this.to, this.endSide);
          break;
        }
      }
    }
    if (trackOpen) {
      let openStart = 0;
      while (openStart < trackOpen.length && trackOpen[openStart] < from)
        openStart++;
      this.openStart = openStart + trackExtra;
    }
  }
  activeForPoint(to) {
    if (!this.active.length)
      return this.active;
    let active = [];
    for (let i = 0; i < this.active.length; i++) {
      if (this.activeRank[i] > this.pointRank)
        break;
      if (this.activeTo[i] > to || this.activeTo[i] == to && this.active[i].endSide > this.point.endSide)
        active.push(this.active[i]);
    }
    return active;
  }
  openEnd(to) {
    let open = 0;
    while (open < this.activeTo.length && this.activeTo[open] > to)
      open++;
    return open;
  }
};
function compare(a, startA, b, startB, length, comparator) {
  a.goto(startA);
  b.goto(startB);
  let endB = startB + length;
  let pos = startB, dPos = startB - startA;
  for (; ; ) {
    let diff = a.to + dPos - b.to || a.endSide - b.endSide;
    let end = diff < 0 ? a.to + dPos : b.to, clipEnd = Math.min(end, endB);
    if (a.point || b.point) {
      if (!(a.point && b.point && (a.point == b.point || a.point.eq(b.point))))
        comparator.comparePoint(pos, clipEnd, a.point, b.point);
    } else {
      if (clipEnd > pos && !sameValues(a.active, b.active))
        comparator.compareRange(pos, clipEnd, a.active, b.active);
    }
    if (end > endB)
      break;
    pos = end;
    if (diff <= 0)
      a.next();
    if (diff >= 0)
      b.next();
  }
}
function sameValues(a, b) {
  if (a.length != b.length)
    return false;
  for (let i = 0; i < a.length; i++)
    if (a[i] != b[i] && !a[i].eq(b[i]))
      return false;
  return true;
}
function remove(array, index) {
  for (let i = index, e = array.length - 1; i < e; i++)
    array[i] = array[i + 1];
  array.pop();
}
function insert(array, index, value) {
  for (let i = array.length - 1; i >= index; i--)
    array[i + 1] = array[i];
  array[index] = value;
}
function findMinIndex(value, array) {
  let found = -1, foundPos = 1e9;
  for (let i = 0; i < array.length; i++)
    if ((array[i] - foundPos || value[i].endSide - value[found].endSide) < 0) {
      found = i;
      foundPos = array[i];
    }
  return found;
}

// node_modules/w3c-keyname/index.es.js
var base = {
  8: "Backspace",
  9: "Tab",
  10: "Enter",
  12: "NumLock",
  13: "Enter",
  16: "Shift",
  17: "Control",
  18: "Alt",
  20: "CapsLock",
  27: "Escape",
  32: " ",
  33: "PageUp",
  34: "PageDown",
  35: "End",
  36: "Home",
  37: "ArrowLeft",
  38: "ArrowUp",
  39: "ArrowRight",
  40: "ArrowDown",
  44: "PrintScreen",
  45: "Insert",
  46: "Delete",
  59: ";",
  61: "=",
  91: "Meta",
  92: "Meta",
  106: "*",
  107: "+",
  108: ",",
  109: "-",
  110: ".",
  111: "/",
  144: "NumLock",
  145: "ScrollLock",
  160: "Shift",
  161: "Shift",
  162: "Control",
  163: "Control",
  164: "Alt",
  165: "Alt",
  173: "-",
  186: ";",
  187: "=",
  188: ",",
  189: "-",
  190: ".",
  191: "/",
  192: "`",
  219: "[",
  220: "\\",
  221: "]",
  222: "'",
  229: "q"
};
var shift = {
  48: ")",
  49: "!",
  50: "@",
  51: "#",
  52: "$",
  53: "%",
  54: "^",
  55: "&",
  56: "*",
  57: "(",
  59: ":",
  61: "+",
  173: "_",
  186: ":",
  187: "+",
  188: "<",
  189: "_",
  190: ">",
  191: "?",
  192: "~",
  219: "{",
  220: "|",
  221: "}",
  222: '"',
  229: "Q"
};
var chrome = typeof navigator != "undefined" && /Chrome\/(\d+)/.exec(navigator.userAgent);
var safari = typeof navigator != "undefined" && /Apple Computer/.test(navigator.vendor);
var gecko = typeof navigator != "undefined" && /Gecko\/\d+/.test(navigator.userAgent);
var mac = typeof navigator != "undefined" && /Mac/.test(navigator.platform);
var ie = typeof navigator != "undefined" && /MSIE \d|Trident\/(?:[7-9]|\d{2,})\..*rv:(\d+)/.exec(navigator.userAgent);
var brokenModifierNames = chrome && (mac || +chrome[1] < 57) || gecko && mac;
for (var i = 0; i < 10; i++)
  base[48 + i] = base[96 + i] = String(i);
for (var i = 1; i <= 24; i++)
  base[i + 111] = "F" + i;
for (var i = 65; i <= 90; i++) {
  base[i] = String.fromCharCode(i + 32);
  shift[i] = String.fromCharCode(i);
}
for (var code in base)
  if (!shift.hasOwnProperty(code))
    shift[code] = base[code];
function keyName(event) {
  var ignoreKey = brokenModifierNames && (event.ctrlKey || event.altKey || event.metaKey) || (safari || ie) && event.shiftKey && event.key && event.key.length == 1;
  var name2 = !ignoreKey && event.key || (event.shiftKey ? shift : base)[event.keyCode] || event.key || "Unidentified";
  if (name2 == "Esc")
    name2 = "Escape";
  if (name2 == "Del")
    name2 = "Delete";
  if (name2 == "Left")
    name2 = "ArrowLeft";
  if (name2 == "Up")
    name2 = "ArrowUp";
  if (name2 == "Right")
    name2 = "ArrowRight";
  if (name2 == "Down")
    name2 = "ArrowDown";
  return name2;
}

// node_modules/@codemirror/view/dist/index.js
var [nav, doc] = typeof navigator != "undefined" ? [navigator, document] : [{userAgent: "", vendor: "", platform: ""}, {documentElement: {style: {}}}];
var ie_edge = /* @__PURE__ */ /Edge\/(\d+)/.exec(nav.userAgent);
var ie_upto10 = /* @__PURE__ */ /MSIE \d/.test(nav.userAgent);
var ie_11up = /* @__PURE__ */ /Trident\/(?:[7-9]|\d{2,})\..*rv:(\d+)/.exec(nav.userAgent);
var ie2 = !!(ie_upto10 || ie_11up || ie_edge);
var gecko2 = !ie2 && /* @__PURE__ */ /gecko\/(\d+)/i.test(nav.userAgent);
var chrome2 = !ie2 && /* @__PURE__ */ /Chrome\/(\d+)/.exec(nav.userAgent);
var webkit = "webkitFontSmoothing" in doc.documentElement.style;
var safari2 = !ie2 && /* @__PURE__ */ /Apple Computer/.test(nav.vendor);
var browser = {
  mac: /* @__PURE__ */ /Mac/.test(nav.platform),
  ie: ie2,
  ie_version: ie_upto10 ? doc.documentMode || 6 : ie_11up ? +ie_11up[1] : ie_edge ? +ie_edge[1] : 0,
  gecko: gecko2,
  gecko_version: gecko2 ? +(/* @__PURE__ */ /Firefox\/(\d+)/.exec(nav.userAgent) || [0, 0])[1] : 0,
  chrome: !!chrome2,
  chrome_version: chrome2 ? +chrome2[1] : 0,
  ios: safari2 && (/* @__PURE__ */ /Mobile\/\w+/.test(nav.userAgent) || nav.maxTouchPoints > 2),
  android: /* @__PURE__ */ /Android\b/.test(nav.userAgent),
  webkit,
  safari: safari2,
  webkit_version: webkit ? +(/* @__PURE__ */ /\bAppleWebKit\/(\d+)/.exec(navigator.userAgent) || [0, 0])[1] : 0,
  tabSize: doc.documentElement.style.tabSize != null ? "tab-size" : "-moz-tab-size"
};
function getSelection(root) {
  return root.getSelection ? root.getSelection() : document.getSelection();
}
function selectionCollapsed(domSel) {
  let collapsed = domSel.isCollapsed;
  if (collapsed && browser.chrome && domSel.rangeCount && !domSel.getRangeAt(0).collapsed)
    collapsed = false;
  return collapsed;
}
function contains(dom, node) {
  return node ? dom.contains(node.nodeType != 1 ? node.parentNode : node) : false;
}
function hasSelection(dom, selection) {
  if (!selection.anchorNode)
    return false;
  try {
    return contains(dom, selection.anchorNode);
  } catch (_) {
    return false;
  }
}
function clientRectsFor(dom) {
  if (dom.nodeType == 3)
    return textRange(dom, 0, dom.nodeValue.length).getClientRects();
  else if (dom.nodeType == 1)
    return dom.getClientRects();
  else
    return [];
}
function isEquivalentPosition(node, off, targetNode, targetOff) {
  return targetNode ? scanFor(node, off, targetNode, targetOff, -1) || scanFor(node, off, targetNode, targetOff, 1) : false;
}
function domIndex(node) {
  for (var index = 0; ; index++) {
    node = node.previousSibling;
    if (!node)
      return index;
  }
}
function scanFor(node, off, targetNode, targetOff, dir) {
  for (; ; ) {
    if (node == targetNode && off == targetOff)
      return true;
    if (off == (dir < 0 ? 0 : maxOffset(node))) {
      if (node.nodeName == "DIV")
        return false;
      let parent = node.parentNode;
      if (!parent || parent.nodeType != 1)
        return false;
      off = domIndex(node) + (dir < 0 ? 0 : 1);
      node = parent;
    } else if (node.nodeType == 1) {
      node = node.childNodes[off + (dir < 0 ? -1 : 0)];
      off = dir < 0 ? maxOffset(node) : 0;
    } else {
      return false;
    }
  }
}
function maxOffset(node) {
  return node.nodeType == 3 ? node.nodeValue.length : node.childNodes.length;
}
var Rect0 = {left: 0, right: 0, top: 0, bottom: 0};
function flattenRect(rect, left) {
  let x = left ? rect.left : rect.right;
  return {left: x, right: x, top: rect.top, bottom: rect.bottom};
}
function windowRect(win) {
  return {
    left: 0,
    right: win.innerWidth,
    top: 0,
    bottom: win.innerHeight
  };
}
var ScrollSpace = 5;
function scrollRectIntoView(dom, rect) {
  let doc2 = dom.ownerDocument, win = doc2.defaultView;
  for (let cur2 = dom.parentNode; cur2; ) {
    if (cur2.nodeType == 1) {
      let bounding, top2 = cur2 == document.body;
      if (top2) {
        bounding = windowRect(win);
      } else {
        if (cur2.scrollHeight <= cur2.clientHeight && cur2.scrollWidth <= cur2.clientWidth) {
          cur2 = cur2.parentNode;
          continue;
        }
        let rect2 = cur2.getBoundingClientRect();
        bounding = {
          left: rect2.left,
          right: rect2.left + cur2.clientWidth,
          top: rect2.top,
          bottom: rect2.top + cur2.clientHeight
        };
      }
      let moveX = 0, moveY = 0;
      if (rect.top < bounding.top)
        moveY = -(bounding.top - rect.top + ScrollSpace);
      else if (rect.bottom > bounding.bottom)
        moveY = rect.bottom - bounding.bottom + ScrollSpace;
      if (rect.left < bounding.left)
        moveX = -(bounding.left - rect.left + ScrollSpace);
      else if (rect.right > bounding.right)
        moveX = rect.right - bounding.right + ScrollSpace;
      if (moveX || moveY) {
        if (top2) {
          win.scrollBy(moveX, moveY);
        } else {
          if (moveY) {
            let start = cur2.scrollTop;
            cur2.scrollTop += moveY;
            moveY = cur2.scrollTop - start;
          }
          if (moveX) {
            let start = cur2.scrollLeft;
            cur2.scrollLeft += moveX;
            moveX = cur2.scrollLeft - start;
          }
          rect = {
            left: rect.left - moveX,
            top: rect.top - moveY,
            right: rect.right - moveX,
            bottom: rect.bottom - moveY
          };
        }
      }
      if (top2)
        break;
      cur2 = cur2.parentNode;
    } else if (cur2.nodeType == 11) {
      cur2 = cur2.host;
    } else {
      break;
    }
  }
}
var DOMSelection = class {
  constructor() {
    this.anchorNode = null;
    this.anchorOffset = 0;
    this.focusNode = null;
    this.focusOffset = 0;
  }
  eq(domSel) {
    return this.anchorNode == domSel.anchorNode && this.anchorOffset == domSel.anchorOffset && this.focusNode == domSel.focusNode && this.focusOffset == domSel.focusOffset;
  }
  set(domSel) {
    this.anchorNode = domSel.anchorNode;
    this.anchorOffset = domSel.anchorOffset;
    this.focusNode = domSel.focusNode;
    this.focusOffset = domSel.focusOffset;
  }
};
var preventScrollSupported = null;
function focusPreventScroll(dom) {
  if (dom.setActive)
    return dom.setActive();
  if (preventScrollSupported)
    return dom.focus(preventScrollSupported);
  let stack = [];
  for (let cur2 = dom; cur2; cur2 = cur2.parentNode) {
    stack.push(cur2, cur2.scrollTop, cur2.scrollLeft);
    if (cur2 == cur2.ownerDocument)
      break;
  }
  dom.focus(preventScrollSupported == null ? {
    get preventScroll() {
      preventScrollSupported = {preventScroll: true};
      return true;
    }
  } : void 0);
  if (!preventScrollSupported) {
    preventScrollSupported = false;
    for (let i = 0; i < stack.length; ) {
      let elt = stack[i++], top2 = stack[i++], left = stack[i++];
      if (elt.scrollTop != top2)
        elt.scrollTop = top2;
      if (elt.scrollLeft != left)
        elt.scrollLeft = left;
    }
  }
}
var scratchRange;
function textRange(node, from, to = from) {
  let range = scratchRange || (scratchRange = document.createRange());
  range.setEnd(node, to);
  range.setStart(node, from);
  return range;
}
var DOMPos = class {
  constructor(node, offset, precise = true) {
    this.node = node;
    this.offset = offset;
    this.precise = precise;
  }
  static before(dom, precise) {
    return new DOMPos(dom.parentNode, domIndex(dom), precise);
  }
  static after(dom, precise) {
    return new DOMPos(dom.parentNode, domIndex(dom) + 1, precise);
  }
};
var none$3 = [];
var ContentView = class {
  constructor() {
    this.parent = null;
    this.dom = null;
    this.dirty = 2;
  }
  get editorView() {
    if (!this.parent)
      throw new Error("Accessing view in orphan content view");
    return this.parent.editorView;
  }
  get overrideDOMText() {
    return null;
  }
  get posAtStart() {
    return this.parent ? this.parent.posBefore(this) : 0;
  }
  get posAtEnd() {
    return this.posAtStart + this.length;
  }
  posBefore(view) {
    let pos = this.posAtStart;
    for (let child of this.children) {
      if (child == view)
        return pos;
      pos += child.length + child.breakAfter;
    }
    throw new RangeError("Invalid child in posBefore");
  }
  posAfter(view) {
    return this.posBefore(view) + view.length;
  }
  coordsAt(_pos, _side) {
    return null;
  }
  sync(track) {
    if (this.dirty & 2) {
      let parent = this.dom, pos = null;
      for (let child of this.children) {
        if (child.dirty) {
          let next2 = pos ? pos.nextSibling : parent.firstChild;
          if (next2 && !child.dom && !ContentView.get(next2))
            child.reuseDOM(next2);
          child.sync(track);
          child.dirty = 0;
        }
        if (track && track.node == parent && pos != child.dom)
          track.written = true;
        syncNodeInto(parent, pos, child.dom);
        pos = child.dom;
      }
      let next = pos ? pos.nextSibling : parent.firstChild;
      if (next && track && track.node == parent)
        track.written = true;
      while (next)
        next = rm(next);
    } else if (this.dirty & 1) {
      for (let child of this.children)
        if (child.dirty) {
          child.sync(track);
          child.dirty = 0;
        }
    }
  }
  reuseDOM(_dom) {
    return false;
  }
  localPosFromDOM(node, offset) {
    let after;
    if (node == this.dom) {
      after = this.dom.childNodes[offset];
    } else {
      let bias = maxOffset(node) == 0 ? 0 : offset == 0 ? -1 : 1;
      for (; ; ) {
        let parent = node.parentNode;
        if (parent == this.dom)
          break;
        if (bias == 0 && parent.firstChild != parent.lastChild) {
          if (node == parent.firstChild)
            bias = -1;
          else
            bias = 1;
        }
        node = parent;
      }
      if (bias < 0)
        after = node;
      else
        after = node.nextSibling;
    }
    if (after == this.dom.firstChild)
      return 0;
    while (after && !ContentView.get(after))
      after = after.nextSibling;
    if (!after)
      return this.length;
    for (let i = 0, pos = 0; ; i++) {
      let child = this.children[i];
      if (child.dom == after)
        return pos;
      pos += child.length + child.breakAfter;
    }
  }
  domBoundsAround(from, to, offset = 0) {
    let fromI = -1, fromStart = -1, toI = -1, toEnd = -1;
    for (let i = 0, pos = offset; i < this.children.length; i++) {
      let child = this.children[i], end = pos + child.length;
      if (pos < from && end > to)
        return child.domBoundsAround(from, to, pos);
      if (end >= from && fromI == -1) {
        fromI = i;
        fromStart = pos;
      }
      if (end >= to && end != pos && toI == -1) {
        toI = i;
        toEnd = end;
        break;
      }
      pos = end + child.breakAfter;
    }
    return {from: fromStart, to: toEnd < 0 ? offset + this.length : toEnd, startDOM: (fromI ? this.children[fromI - 1].dom.nextSibling : null) || this.dom.firstChild, endDOM: toI < this.children.length - 1 && toI >= 0 ? this.children[toI + 1].dom : null};
  }
  markDirty(andParent = false) {
    if (this.dirty & 2)
      return;
    this.dirty |= 2;
    this.markParentsDirty(andParent);
  }
  markParentsDirty(childList) {
    for (let parent = this.parent; parent; parent = parent.parent) {
      if (childList)
        parent.dirty |= 2;
      if (parent.dirty & 1)
        return;
      parent.dirty |= 1;
      childList = false;
    }
  }
  setParent(parent) {
    if (this.parent != parent) {
      this.parent = parent;
      if (this.dirty)
        this.markParentsDirty(true);
    }
  }
  setDOM(dom) {
    this.dom = dom;
    dom.cmView = this;
  }
  get rootView() {
    for (let v = this; ; ) {
      let parent = v.parent;
      if (!parent)
        return v;
      v = parent;
    }
  }
  replaceChildren(from, to, children = none$3) {
    this.markDirty();
    for (let i = from; i < to; i++)
      this.children[i].parent = null;
    this.children.splice(from, to - from, ...children);
    for (let i = 0; i < children.length; i++)
      children[i].setParent(this);
  }
  ignoreMutation(_rec) {
    return false;
  }
  ignoreEvent(_event) {
    return false;
  }
  childCursor(pos = this.length) {
    return new ChildCursor(this.children, pos, this.children.length);
  }
  childPos(pos, bias = 1) {
    return this.childCursor().findPos(pos, bias);
  }
  toString() {
    let name2 = this.constructor.name.replace("View", "");
    return name2 + (this.children.length ? "(" + this.children.join() + ")" : this.length ? "[" + (name2 == "Text" ? this.text : this.length) + "]" : "") + (this.breakAfter ? "#" : "");
  }
  static get(node) {
    return node.cmView;
  }
};
ContentView.prototype.breakAfter = 0;
function rm(dom) {
  let next = dom.nextSibling;
  dom.parentNode.removeChild(dom);
  return next;
}
function syncNodeInto(parent, after, dom) {
  let next = after ? after.nextSibling : parent.firstChild;
  if (dom.parentNode == parent)
    while (next != dom)
      next = rm(next);
  else
    parent.insertBefore(dom, next);
}
var ChildCursor = class {
  constructor(children, pos, i) {
    this.children = children;
    this.pos = pos;
    this.i = i;
    this.off = 0;
  }
  findPos(pos, bias = 1) {
    for (; ; ) {
      if (pos > this.pos || pos == this.pos && (bias > 0 || this.i == 0 || this.children[this.i - 1].breakAfter)) {
        this.off = pos - this.pos;
        return this;
      }
      let next = this.children[--this.i];
      this.pos -= next.length + next.breakAfter;
    }
  }
};
var none$2 = [];
var InlineView = class extends ContentView {
  become(_other) {
    return false;
  }
  getSide() {
    return 0;
  }
};
InlineView.prototype.children = none$2;
var MaxJoinLen = 256;
var TextView = class extends InlineView {
  constructor(text) {
    super();
    this.text = text;
  }
  get length() {
    return this.text.length;
  }
  createDOM(textDOM) {
    this.setDOM(textDOM || document.createTextNode(this.text));
  }
  sync(track) {
    if (!this.dom)
      this.createDOM();
    if (this.dom.nodeValue != this.text) {
      if (track && track.node == this.dom)
        track.written = true;
      this.dom.nodeValue = this.text;
    }
  }
  reuseDOM(dom) {
    if (dom.nodeType != 3)
      return false;
    this.createDOM(dom);
    return true;
  }
  merge(from, to, source) {
    if (source && (!(source instanceof TextView) || this.length - (to - from) + source.length > MaxJoinLen))
      return false;
    this.text = this.text.slice(0, from) + (source ? source.text : "") + this.text.slice(to);
    this.markDirty();
    return true;
  }
  slice(from) {
    return new TextView(this.text.slice(from));
  }
  localPosFromDOM(node, offset) {
    return node == this.dom ? offset : offset ? this.text.length : 0;
  }
  domAtPos(pos) {
    return new DOMPos(this.dom, pos);
  }
  domBoundsAround(_from, _to, offset) {
    return {from: offset, to: offset + this.length, startDOM: this.dom, endDOM: this.dom.nextSibling};
  }
  coordsAt(pos, side) {
    return textCoords(this.dom, pos, side);
  }
};
var MarkView = class extends InlineView {
  constructor(mark, children = [], length = 0) {
    super();
    this.mark = mark;
    this.children = children;
    this.length = length;
    for (let ch of children)
      ch.setParent(this);
  }
  createDOM() {
    let dom = document.createElement(this.mark.tagName);
    if (this.mark.class)
      dom.className = this.mark.class;
    if (this.mark.attrs)
      for (let name2 in this.mark.attrs)
        dom.setAttribute(name2, this.mark.attrs[name2]);
    this.setDOM(dom);
  }
  sync(track) {
    if (!this.dom)
      this.createDOM();
    super.sync(track);
  }
  merge(from, to, source, openStart, openEnd) {
    if (source && (!(source instanceof MarkView && source.mark.eq(this.mark)) || from && openStart <= 0 || to < this.length && openEnd <= 0))
      return false;
    mergeInlineChildren(this, from, to, source ? source.children : none$2, openStart - 1, openEnd - 1);
    this.markDirty();
    return true;
  }
  slice(from) {
    return new MarkView(this.mark, sliceInlineChildren(this.children, from), this.length - from);
  }
  domAtPos(pos) {
    return inlineDOMAtPos(this.dom, this.children, pos);
  }
  coordsAt(pos, side) {
    return coordsInChildren(this, pos, side);
  }
};
function textCoords(text, pos, side) {
  let length = text.nodeValue.length;
  if (pos > length)
    pos = length;
  let from = pos, to = pos, flatten2 = 0;
  if (pos == 0 && side < 0 || pos == length && side >= 0) {
    if (!(browser.chrome || browser.gecko)) {
      if (pos) {
        from--;
        flatten2 = 1;
      } else {
        to++;
        flatten2 = -1;
      }
    }
  } else {
    if (side < 0)
      from--;
    else
      to++;
  }
  let rects = textRange(text, from, to).getClientRects();
  if (!rects.length)
    return Rect0;
  let rect = rects[(flatten2 ? flatten2 < 0 : side >= 0) ? 0 : rects.length - 1];
  if (browser.safari && !flatten2 && rect.width == 0)
    rect = Array.prototype.find.call(rects, (r) => r.width) || rect;
  return flatten2 ? flattenRect(rect, flatten2 < 0) : rect;
}
var WidgetView = class extends InlineView {
  constructor(widget, length, side) {
    super();
    this.widget = widget;
    this.length = length;
    this.side = side;
  }
  static create(widget, length, side) {
    return new (widget.customView || WidgetView)(widget, length, side);
  }
  slice(from) {
    return WidgetView.create(this.widget, this.length - from, this.side);
  }
  sync() {
    if (!this.dom || !this.widget.updateDOM(this.dom)) {
      this.setDOM(this.widget.toDOM(this.editorView));
      this.dom.contentEditable = "false";
    }
  }
  getSide() {
    return this.side;
  }
  merge(from, to, source, openStart, openEnd) {
    if (source && (!(source instanceof WidgetView) || !this.widget.compare(source.widget) || from > 0 && openStart <= 0 || to < this.length && openEnd <= 0))
      return false;
    this.length = from + (source ? source.length : 0) + (this.length - to);
    return true;
  }
  become(other) {
    if (other.length == this.length && other instanceof WidgetView && other.side == this.side) {
      if (this.widget.constructor == other.widget.constructor) {
        if (!this.widget.eq(other.widget))
          this.markDirty(true);
        this.widget = other.widget;
        return true;
      }
    }
    return false;
  }
  ignoreMutation() {
    return true;
  }
  ignoreEvent(event) {
    return this.widget.ignoreEvent(event);
  }
  get overrideDOMText() {
    if (this.length == 0)
      return Text.empty;
    let top2 = this;
    while (top2.parent)
      top2 = top2.parent;
    let view = top2.editorView, text = view && view.state.doc, start = this.posAtStart;
    return text ? text.slice(start, start + this.length) : Text.empty;
  }
  domAtPos(pos) {
    return pos == 0 ? DOMPos.before(this.dom) : DOMPos.after(this.dom, pos == this.length);
  }
  domBoundsAround() {
    return null;
  }
  coordsAt(pos, side) {
    let rects = this.dom.getClientRects(), rect = null;
    if (!rects.length)
      return Rect0;
    for (let i = pos > 0 ? rects.length - 1 : 0; ; i += pos > 0 ? -1 : 1) {
      rect = rects[i];
      if (pos > 0 ? i == 0 : i == rects.length - 1 || rect.top < rect.bottom)
        break;
    }
    return pos == 0 && side > 0 || pos == this.length && side <= 0 ? rect : flattenRect(rect, pos == 0);
  }
};
var CompositionView = class extends WidgetView {
  domAtPos(pos) {
    return new DOMPos(this.widget.text, pos);
  }
  sync() {
    if (!this.dom)
      this.setDOM(this.widget.toDOM());
  }
  localPosFromDOM(node, offset) {
    return !offset ? 0 : node.nodeType == 3 ? Math.min(offset, this.length) : this.length;
  }
  ignoreMutation() {
    return false;
  }
  get overrideDOMText() {
    return null;
  }
  coordsAt(pos, side) {
    return textCoords(this.widget.text, pos, side);
  }
};
function mergeInlineChildren(parent, from, to, elts, openStart, openEnd) {
  let cur2 = parent.childCursor();
  let {i: toI, off: toOff} = cur2.findPos(to, 1);
  let {i: fromI, off: fromOff} = cur2.findPos(from, -1);
  let dLen = from - to;
  for (let view of elts)
    dLen += view.length;
  parent.length += dLen;
  let {children} = parent;
  if (fromI == toI && fromOff) {
    let start = children[fromI];
    if (elts.length == 1 && start.merge(fromOff, toOff, elts[0], openStart, openEnd))
      return;
    if (elts.length == 0) {
      start.merge(fromOff, toOff, null, openStart, openEnd);
      return;
    }
    let after = start.slice(toOff);
    if (after.merge(0, 0, elts[elts.length - 1], 0, openEnd))
      elts[elts.length - 1] = after;
    else
      elts.push(after);
    toI++;
    openEnd = toOff = 0;
  }
  if (toOff) {
    let end = children[toI];
    if (elts.length && end.merge(0, toOff, elts[elts.length - 1], 0, openEnd)) {
      elts.pop();
      openEnd = 0;
    } else {
      end.merge(0, toOff, null, 0, 0);
    }
  } else if (toI < children.length && elts.length && children[toI].merge(0, 0, elts[elts.length - 1], 0, openEnd)) {
    elts.pop();
    openEnd = 0;
  }
  if (fromOff) {
    let start = children[fromI];
    if (elts.length && start.merge(fromOff, start.length, elts[0], openStart, 0)) {
      elts.shift();
      openStart = 0;
    } else {
      start.merge(fromOff, start.length, null, 0, 0);
    }
    fromI++;
  } else if (fromI && elts.length) {
    let end = children[fromI - 1];
    if (end.merge(end.length, end.length, elts[0], openStart, 0)) {
      elts.shift();
      openStart = 0;
    }
  }
  while (fromI < toI && elts.length && children[toI - 1].become(elts[elts.length - 1])) {
    elts.pop();
    toI--;
    openEnd = 0;
  }
  while (fromI < toI && elts.length && children[fromI].become(elts[0])) {
    elts.shift();
    fromI++;
    openStart = 0;
  }
  if (!elts.length && fromI && toI < children.length && openStart && openEnd && children[toI].merge(0, 0, children[fromI - 1], openStart, openEnd))
    fromI--;
  if (elts.length || fromI != toI)
    parent.replaceChildren(fromI, toI, elts);
}
function sliceInlineChildren(children, from) {
  let result = [], off = 0;
  for (let elt of children) {
    let end = off + elt.length;
    if (end > from)
      result.push(off < from ? elt.slice(from - off) : elt);
    off = end;
  }
  return result;
}
function inlineDOMAtPos(dom, children, pos) {
  let i = 0;
  for (let off = 0; i < children.length; i++) {
    let child = children[i], end = off + child.length;
    if (end == off && child.getSide() <= 0)
      continue;
    if (pos > off && pos < end && child.dom.parentNode == dom)
      return child.domAtPos(pos - off);
    if (pos <= off)
      break;
    off = end;
  }
  for (; i > 0; i--) {
    let before = children[i - 1].dom;
    if (before.parentNode == dom)
      return DOMPos.after(before);
  }
  return new DOMPos(dom, 0);
}
function joinInlineInto(parent, view, open) {
  let last, {children} = parent;
  if (open > 0 && view instanceof MarkView && children.length && (last = children[children.length - 1]) instanceof MarkView && last.mark.eq(view.mark)) {
    joinInlineInto(last, view.children[0], open - 1);
  } else {
    children.push(view);
    view.setParent(parent);
  }
  parent.length += view.length;
}
function coordsInChildren(view, pos, side) {
  for (let off = 0, i = 0; i < view.children.length; i++) {
    let child = view.children[i], end = off + child.length;
    if (end == off && child.getSide() <= 0)
      continue;
    if (side <= 0 || end == view.length ? end >= pos : end > pos)
      return child.coordsAt(pos - off, side);
    off = end;
  }
  return (view.dom.lastChild || view.dom).getBoundingClientRect();
}
function combineAttrs(source, target) {
  for (let name2 in source) {
    if (name2 == "class" && target.class)
      target.class += " " + source.class;
    else if (name2 == "style" && target.style)
      target.style += ";" + source.style;
    else
      target[name2] = source[name2];
  }
  return target;
}
function attrsEq(a, b) {
  if (a == b)
    return true;
  if (!a || !b)
    return false;
  let keysA = Object.keys(a), keysB = Object.keys(b);
  if (keysA.length != keysB.length)
    return false;
  for (let key of keysA) {
    if (keysB.indexOf(key) == -1 || a[key] !== b[key])
      return false;
  }
  return true;
}
function updateAttrs(dom, prev, attrs) {
  if (prev) {
    for (let name2 in prev)
      if (!(attrs && name2 in attrs))
        dom.removeAttribute(name2);
  }
  if (attrs) {
    for (let name2 in attrs)
      if (!(prev && prev[name2] == attrs[name2]))
        dom.setAttribute(name2, attrs[name2]);
  }
}
var WidgetType = class {
  eq(_widget) {
    return false;
  }
  updateDOM(_dom) {
    return false;
  }
  compare(other) {
    return this == other || this.constructor == other.constructor && this.eq(other);
  }
  get estimatedHeight() {
    return -1;
  }
  ignoreEvent(_event) {
    return true;
  }
  get customView() {
    return null;
  }
};
var BlockType = /* @__PURE__ */ function(BlockType2) {
  BlockType2[BlockType2["Text"] = 0] = "Text";
  BlockType2[BlockType2["WidgetBefore"] = 1] = "WidgetBefore";
  BlockType2[BlockType2["WidgetAfter"] = 2] = "WidgetAfter";
  BlockType2[BlockType2["WidgetRange"] = 3] = "WidgetRange";
  return BlockType2;
}(BlockType || (BlockType = {}));
var Decoration = class extends RangeValue {
  constructor(startSide, endSide, widget, spec) {
    super();
    this.startSide = startSide;
    this.endSide = endSide;
    this.widget = widget;
    this.spec = spec;
  }
  get heightRelevant() {
    return false;
  }
  static mark(spec) {
    return new MarkDecoration(spec);
  }
  static widget(spec) {
    let side = spec.side || 0;
    if (spec.block)
      side += (2e8 + 1) * (side > 0 ? 1 : -1);
    return new PointDecoration(spec, side, side, !!spec.block, spec.widget || null, false);
  }
  static replace(spec) {
    let block = !!spec.block;
    let {start, end} = getInclusive(spec);
    let startSide = block ? -2e8 * (start ? 2 : 1) : 1e8 * (start ? -1 : 1);
    let endSide = block ? 2e8 * (end ? 2 : 1) : 1e8 * (end ? 1 : -1);
    return new PointDecoration(spec, startSide, endSide, block, spec.widget || null, true);
  }
  static line(spec) {
    return new LineDecoration(spec);
  }
  static set(of, sort = false) {
    return RangeSet.of(of, sort);
  }
  hasHeight() {
    return this.widget ? this.widget.estimatedHeight > -1 : false;
  }
};
Decoration.none = RangeSet.empty;
var MarkDecoration = class extends Decoration {
  constructor(spec) {
    let {start, end} = getInclusive(spec);
    super(1e8 * (start ? -1 : 1), 1e8 * (end ? 1 : -1), null, spec);
    this.tagName = spec.tagName || "span";
    this.class = spec.class || "";
    this.attrs = spec.attributes || null;
  }
  eq(other) {
    return this == other || other instanceof MarkDecoration && this.tagName == other.tagName && this.class == other.class && attrsEq(this.attrs, other.attrs);
  }
  range(from, to = from) {
    if (from >= to)
      throw new RangeError("Mark decorations may not be empty");
    return super.range(from, to);
  }
};
MarkDecoration.prototype.point = false;
var LineDecoration = class extends Decoration {
  constructor(spec) {
    super(-1e8, -1e8, null, spec);
  }
  eq(other) {
    return other instanceof LineDecoration && attrsEq(this.spec.attributes, other.spec.attributes);
  }
  range(from, to = from) {
    if (to != from)
      throw new RangeError("Line decoration ranges must be zero-length");
    return super.range(from, to);
  }
};
LineDecoration.prototype.mapMode = MapMode.TrackBefore;
LineDecoration.prototype.point = true;
var PointDecoration = class extends Decoration {
  constructor(spec, startSide, endSide, block, widget, isReplace) {
    super(startSide, endSide, widget, spec);
    this.block = block;
    this.isReplace = isReplace;
    this.mapMode = !block ? MapMode.TrackDel : startSide < 0 ? MapMode.TrackBefore : MapMode.TrackAfter;
  }
  get type() {
    return this.startSide < this.endSide ? BlockType.WidgetRange : this.startSide < 0 ? BlockType.WidgetBefore : BlockType.WidgetAfter;
  }
  get heightRelevant() {
    return this.block || !!this.widget && this.widget.estimatedHeight >= 5;
  }
  eq(other) {
    return other instanceof PointDecoration && widgetsEq(this.widget, other.widget) && this.block == other.block && this.startSide == other.startSide && this.endSide == other.endSide;
  }
  range(from, to = from) {
    if (this.isReplace && (from > to || from == to && this.startSide > 0 && this.endSide < 0))
      throw new RangeError("Invalid range for replacement decoration");
    if (!this.isReplace && to != from)
      throw new RangeError("Widget decorations can only have zero-length ranges");
    return super.range(from, to);
  }
};
PointDecoration.prototype.point = true;
function getInclusive(spec) {
  let {inclusiveStart: start, inclusiveEnd: end} = spec;
  if (start == null)
    start = spec.inclusive;
  if (end == null)
    end = spec.inclusive;
  return {start: start || false, end: end || false};
}
function widgetsEq(a, b) {
  return a == b || !!(a && b && a.compare(b));
}
function addRange(from, to, ranges, margin = 0) {
  let last = ranges.length - 1;
  if (last >= 0 && ranges[last] + margin > from)
    ranges[last] = Math.max(ranges[last], to);
  else
    ranges.push(from, to);
}
var LineView = class extends ContentView {
  constructor() {
    super(...arguments);
    this.children = [];
    this.length = 0;
    this.prevAttrs = void 0;
    this.attrs = null;
    this.breakAfter = 0;
  }
  merge(from, to, source, takeDeco, openStart, openEnd) {
    if (source) {
      if (!(source instanceof LineView))
        return false;
      if (!this.dom)
        source.transferDOM(this);
    }
    if (takeDeco)
      this.setDeco(source ? source.attrs : null);
    mergeInlineChildren(this, from, to, source ? source.children : none$1, openStart, openEnd);
    return true;
  }
  split(at) {
    let end = new LineView();
    end.breakAfter = this.breakAfter;
    if (this.length == 0)
      return end;
    let {i, off} = this.childPos(at);
    if (off) {
      end.append(this.children[i].slice(off), 0);
      this.children[i].merge(off, this.children[i].length, null, 0, 0);
      i++;
    }
    for (let j = i; j < this.children.length; j++)
      end.append(this.children[j], 0);
    while (i > 0 && this.children[i - 1].length == 0) {
      this.children[i - 1].parent = null;
      i--;
    }
    this.children.length = i;
    this.markDirty();
    this.length = at;
    return end;
  }
  transferDOM(other) {
    if (!this.dom)
      return;
    other.setDOM(this.dom);
    other.prevAttrs = this.prevAttrs === void 0 ? this.attrs : this.prevAttrs;
    this.prevAttrs = void 0;
    this.dom = null;
  }
  setDeco(attrs) {
    if (!attrsEq(this.attrs, attrs)) {
      if (this.dom) {
        this.prevAttrs = this.attrs;
        this.markDirty();
      }
      this.attrs = attrs;
    }
  }
  append(child, openStart) {
    joinInlineInto(this, child, openStart);
  }
  addLineDeco(deco) {
    let attrs = deco.spec.attributes;
    if (attrs)
      this.attrs = combineAttrs(attrs, this.attrs || {});
  }
  domAtPos(pos) {
    return inlineDOMAtPos(this.dom, this.children, pos);
  }
  sync(track) {
    if (!this.dom) {
      this.setDOM(document.createElement("div"));
      this.dom.className = "cm-line";
      this.prevAttrs = this.attrs ? null : void 0;
    }
    if (this.prevAttrs !== void 0) {
      updateAttrs(this.dom, this.prevAttrs, this.attrs);
      this.dom.classList.add("cm-line");
      this.prevAttrs = void 0;
    }
    super.sync(track);
    let last = this.dom.lastChild;
    if (!last || last.nodeName != "BR" && ContentView.get(last) instanceof WidgetView) {
      let hack = document.createElement("BR");
      hack.cmIgnore = true;
      this.dom.appendChild(hack);
    }
  }
  measureTextSize() {
    if (this.children.length == 0 || this.length > 20)
      return null;
    let totalWidth = 0;
    for (let child of this.children) {
      if (!(child instanceof TextView))
        return null;
      let rects = clientRectsFor(child.dom);
      if (rects.length != 1)
        return null;
      totalWidth += rects[0].width;
    }
    return {lineHeight: this.dom.getBoundingClientRect().height, charWidth: totalWidth / this.length};
  }
  coordsAt(pos, side) {
    return coordsInChildren(this, pos, side);
  }
  match(_other) {
    return false;
  }
  get type() {
    return BlockType.Text;
  }
  static find(docView, pos) {
    for (let i = 0, off = 0; ; i++) {
      let block = docView.children[i], end = off + block.length;
      if (end >= pos) {
        if (block instanceof LineView)
          return block;
        if (block.length)
          return null;
      }
      off = end + block.breakAfter;
    }
  }
};
var none$1 = [];
var BlockWidgetView = class extends ContentView {
  constructor(widget, length, type2) {
    super();
    this.widget = widget;
    this.length = length;
    this.type = type2;
    this.breakAfter = 0;
  }
  merge(from, to, source, _takeDeco, openStart, openEnd) {
    if (source && (!(source instanceof BlockWidgetView) || !this.widget.compare(source.widget) || from > 0 && openStart <= 0 || to < this.length && openEnd <= 0))
      return false;
    this.length = from + (source ? source.length : 0) + (this.length - to);
    return true;
  }
  domAtPos(pos) {
    return pos == 0 ? DOMPos.before(this.dom) : DOMPos.after(this.dom, pos == this.length);
  }
  split(at) {
    let len = this.length - at;
    this.length = at;
    return new BlockWidgetView(this.widget, len, this.type);
  }
  get children() {
    return none$1;
  }
  sync() {
    if (!this.dom || !this.widget.updateDOM(this.dom)) {
      this.setDOM(this.widget.toDOM(this.editorView));
      this.dom.contentEditable = "false";
    }
  }
  get overrideDOMText() {
    return this.parent ? this.parent.view.state.doc.slice(this.posAtStart, this.posAtEnd) : Text.empty;
  }
  domBoundsAround() {
    return null;
  }
  match(other) {
    if (other instanceof BlockWidgetView && other.type == this.type && other.widget.constructor == this.widget.constructor) {
      if (!other.widget.eq(this.widget))
        this.markDirty(true);
      this.widget = other.widget;
      this.length = other.length;
      this.breakAfter = other.breakAfter;
      return true;
    }
    return false;
  }
  ignoreMutation() {
    return true;
  }
  ignoreEvent(event) {
    return this.widget.ignoreEvent(event);
  }
};
var ContentBuilder = class {
  constructor(doc2, pos, end) {
    this.doc = doc2;
    this.pos = pos;
    this.end = end;
    this.content = [];
    this.curLine = null;
    this.breakAtStart = 0;
    this.openStart = -1;
    this.openEnd = -1;
    this.text = "";
    this.textOff = 0;
    this.cursor = doc2.iter();
    this.skip = pos;
  }
  posCovered() {
    if (this.content.length == 0)
      return !this.breakAtStart && this.doc.lineAt(this.pos).from != this.pos;
    let last = this.content[this.content.length - 1];
    return !last.breakAfter && !(last instanceof BlockWidgetView && last.type == BlockType.WidgetBefore);
  }
  getLine() {
    if (!this.curLine)
      this.content.push(this.curLine = new LineView());
    return this.curLine;
  }
  addWidget(view) {
    this.curLine = null;
    this.content.push(view);
  }
  finish() {
    if (!this.posCovered())
      this.getLine();
  }
  wrapMarks(view, active) {
    for (let i = active.length - 1; i >= 0; i--)
      view = new MarkView(active[i], [view], view.length);
    return view;
  }
  buildText(length, active, openStart) {
    while (length > 0) {
      if (this.textOff == this.text.length) {
        let {value, lineBreak, done} = this.cursor.next(this.skip);
        this.skip = 0;
        if (done)
          throw new Error("Ran out of text content when drawing inline views");
        if (lineBreak) {
          if (!this.posCovered())
            this.getLine();
          if (this.content.length)
            this.content[this.content.length - 1].breakAfter = 1;
          else
            this.breakAtStart = 1;
          this.curLine = null;
          length--;
          continue;
        } else {
          this.text = value;
          this.textOff = 0;
        }
      }
      let take = Math.min(this.text.length - this.textOff, length, 512);
      this.getLine().append(this.wrapMarks(new TextView(this.text.slice(this.textOff, this.textOff + take)), active), openStart);
      this.textOff += take;
      length -= take;
      openStart = 0;
    }
  }
  span(from, to, active, openStart) {
    this.buildText(to - from, active, openStart);
    this.pos = to;
    if (this.openStart < 0)
      this.openStart = openStart;
  }
  point(from, to, deco, active, openStart) {
    let len = to - from;
    if (deco instanceof PointDecoration) {
      if (deco.block) {
        let {type: type2} = deco;
        if (type2 == BlockType.WidgetAfter && !this.posCovered())
          this.getLine();
        this.addWidget(new BlockWidgetView(deco.widget || new NullWidget("div"), len, type2));
      } else {
        let widget = this.wrapMarks(WidgetView.create(deco.widget || new NullWidget("span"), len, deco.startSide), active);
        this.getLine().append(widget, openStart);
      }
    } else if (this.doc.lineAt(this.pos).from == this.pos) {
      this.getLine().addLineDeco(deco);
    }
    if (len) {
      if (this.textOff + len <= this.text.length) {
        this.textOff += len;
      } else {
        this.skip += len - (this.text.length - this.textOff);
        this.text = "";
        this.textOff = 0;
      }
      this.pos = to;
    }
    if (this.openStart < 0)
      this.openStart = openStart;
  }
  static build(text, from, to, decorations2) {
    let builder = new ContentBuilder(text, from, to);
    builder.openEnd = RangeSet.spans(decorations2, from, to, builder);
    if (builder.openStart < 0)
      builder.openStart = builder.openEnd;
    builder.finish();
    return builder;
  }
};
var NullWidget = class extends WidgetType {
  constructor(tag) {
    super();
    this.tag = tag;
  }
  eq(other) {
    return other.tag == this.tag;
  }
  toDOM() {
    return document.createElement(this.tag);
  }
  updateDOM(elt) {
    return elt.nodeName.toLowerCase() == this.tag;
  }
};
var none2 = [];
var clickAddsSelectionRange = /* @__PURE__ */ Facet.define();
var dragMovesSelection$1 = /* @__PURE__ */ Facet.define();
var mouseSelectionStyle = /* @__PURE__ */ Facet.define();
var exceptionSink = /* @__PURE__ */ Facet.define();
var updateListener = /* @__PURE__ */ Facet.define();
var inputHandler = /* @__PURE__ */ Facet.define();
function logException(state, exception, context) {
  let handler = state.facet(exceptionSink);
  if (handler.length)
    handler[0](exception);
  else if (window.onerror)
    window.onerror(String(exception), context, void 0, void 0, exception);
  else if (context)
    console.error(context + ":", exception);
  else
    console.error(exception);
}
var editable = /* @__PURE__ */ Facet.define({combine: (values) => values.length ? values[0] : true});
var PluginFieldProvider = class {
  constructor(field, get) {
    this.field = field;
    this.get = get;
  }
};
var PluginField = class {
  from(get) {
    return new PluginFieldProvider(this, get);
  }
  static define() {
    return new PluginField();
  }
};
PluginField.decorations = /* @__PURE__ */ PluginField.define();
PluginField.scrollMargins = /* @__PURE__ */ PluginField.define();
var nextPluginID = 0;
var viewPlugin = /* @__PURE__ */ Facet.define();
var ViewPlugin = class {
  constructor(id2, create, fields) {
    this.id = id2;
    this.create = create;
    this.fields = fields;
    this.extension = viewPlugin.of(this);
  }
  static define(create, spec) {
    let {eventHandlers, provide, decorations: decorations2} = spec || {};
    let fields = [];
    if (provide)
      for (let provider of Array.isArray(provide) ? provide : [provide])
        fields.push(provider);
    if (eventHandlers)
      fields.push(domEventHandlers.from((value) => ({plugin: value, handlers: eventHandlers})));
    if (decorations2)
      fields.push(PluginField.decorations.from(decorations2));
    return new ViewPlugin(nextPluginID++, create, fields);
  }
  static fromClass(cls, spec) {
    return ViewPlugin.define((view) => new cls(view), spec);
  }
};
var domEventHandlers = /* @__PURE__ */ PluginField.define();
var PluginInstance = class {
  constructor(spec) {
    this.spec = spec;
    this.mustUpdate = null;
    this.value = null;
  }
  takeField(type2, target) {
    for (let {field, get} of this.spec.fields)
      if (field == type2)
        target.push(get(this.value));
  }
  update(view) {
    if (!this.value) {
      try {
        this.value = this.spec.create(view);
      } catch (e) {
        logException(view.state, e, "CodeMirror plugin crashed");
        return PluginInstance.dummy;
      }
    } else if (this.mustUpdate) {
      let update = this.mustUpdate;
      this.mustUpdate = null;
      if (!this.value.update)
        return this;
      try {
        this.value.update(update);
      } catch (e) {
        logException(update.state, e, "CodeMirror plugin crashed");
        if (this.value.destroy)
          try {
            this.value.destroy();
          } catch (_) {
          }
        return PluginInstance.dummy;
      }
    }
    return this;
  }
  destroy(view) {
    var _a;
    if ((_a = this.value) === null || _a === void 0 ? void 0 : _a.destroy) {
      try {
        this.value.destroy();
      } catch (e) {
        logException(view.state, e, "CodeMirror plugin crashed");
      }
    }
  }
};
PluginInstance.dummy = /* @__PURE__ */ new PluginInstance(/* @__PURE__ */ ViewPlugin.define(() => ({})));
var editorAttributes = /* @__PURE__ */ Facet.define({
  combine: (values) => values.reduce((a, b) => combineAttrs(b, a), {})
});
var contentAttributes = /* @__PURE__ */ Facet.define({
  combine: (values) => values.reduce((a, b) => combineAttrs(b, a), {})
});
var decorations = /* @__PURE__ */ Facet.define();
var styleModule = /* @__PURE__ */ Facet.define();
var ChangedRange = class {
  constructor(fromA, toA, fromB, toB) {
    this.fromA = fromA;
    this.toA = toA;
    this.fromB = fromB;
    this.toB = toB;
  }
  join(other) {
    return new ChangedRange(Math.min(this.fromA, other.fromA), Math.max(this.toA, other.toA), Math.min(this.fromB, other.fromB), Math.max(this.toB, other.toB));
  }
  addToSet(set) {
    let i = set.length, me = this;
    for (; i > 0; i--) {
      let range = set[i - 1];
      if (range.fromA > me.toA)
        continue;
      if (range.toA < me.fromA)
        break;
      me = me.join(range);
      set.splice(i - 1, 1);
    }
    set.splice(i, 0, me);
    return set;
  }
  static extendWithRanges(diff, ranges) {
    if (ranges.length == 0)
      return diff;
    let result = [];
    for (let dI = 0, rI = 0, posA = 0, posB = 0; ; dI++) {
      let next = dI == diff.length ? null : diff[dI], off = posA - posB;
      let end = next ? next.fromB : 1e9;
      while (rI < ranges.length && ranges[rI] < end) {
        let from = ranges[rI], to = ranges[rI + 1];
        let fromB = Math.max(posB, from), toB = Math.min(end, to);
        if (fromB <= toB)
          new ChangedRange(fromB + off, toB + off, fromB, toB).addToSet(result);
        if (to > end)
          break;
        else
          rI += 2;
      }
      if (!next)
        return result;
      new ChangedRange(next.fromA, next.toA, next.fromB, next.toB).addToSet(result);
      posA = next.toA;
      posB = next.toB;
    }
  }
};
var ViewUpdate = class {
  constructor(view, state, transactions = none2) {
    this.view = view;
    this.state = state;
    this.transactions = transactions;
    this.flags = 0;
    this.startState = view.state;
    this.changes = ChangeSet.empty(this.startState.doc.length);
    for (let tr of transactions)
      this.changes = this.changes.compose(tr.changes);
    let changedRanges = [];
    this.changes.iterChangedRanges((fromA, toA, fromB, toB) => changedRanges.push(new ChangedRange(fromA, toA, fromB, toB)));
    this.changedRanges = changedRanges;
    let focus = view.hasFocus;
    if (focus != view.inputState.notifiedFocused) {
      view.inputState.notifiedFocused = focus;
      this.flags |= 1;
    }
    if (this.docChanged)
      this.flags |= 2;
  }
  get viewportChanged() {
    return (this.flags & 4) > 0;
  }
  get heightChanged() {
    return (this.flags & 2) > 0;
  }
  get geometryChanged() {
    return this.docChanged || (this.flags & (16 | 2)) > 0;
  }
  get focusChanged() {
    return (this.flags & 1) > 0;
  }
  get docChanged() {
    return this.transactions.some((tr) => tr.docChanged);
  }
  get selectionSet() {
    return this.transactions.some((tr) => tr.selection);
  }
  get empty() {
    return this.flags == 0 && this.transactions.length == 0;
  }
};
var DocView = class extends ContentView {
  constructor(view) {
    super();
    this.view = view;
    this.compositionDeco = Decoration.none;
    this.decorations = [];
    this.minWidth = 0;
    this.minWidthFrom = 0;
    this.minWidthTo = 0;
    this.impreciseAnchor = null;
    this.impreciseHead = null;
    this.setDOM(view.contentDOM);
    this.children = [new LineView()];
    this.children[0].setParent(this);
    this.updateInner([new ChangedRange(0, 0, 0, view.state.doc.length)], this.updateDeco(), 0);
  }
  get root() {
    return this.view.root;
  }
  get editorView() {
    return this.view;
  }
  get length() {
    return this.view.state.doc.length;
  }
  update(update) {
    let changedRanges = update.changedRanges;
    if (this.minWidth > 0 && changedRanges.length) {
      if (!changedRanges.every(({fromA, toA}) => toA < this.minWidthFrom || fromA > this.minWidthTo)) {
        this.minWidth = 0;
      } else {
        this.minWidthFrom = update.changes.mapPos(this.minWidthFrom, 1);
        this.minWidthTo = update.changes.mapPos(this.minWidthTo, 1);
      }
    }
    if (this.view.inputState.composing < 0)
      this.compositionDeco = Decoration.none;
    else if (update.transactions.length)
      this.compositionDeco = computeCompositionDeco(this.view, update.changes);
    let forceSelection = (browser.ie || browser.chrome) && !this.compositionDeco.size && update && update.state.doc.lines != update.startState.doc.lines;
    let prevDeco = this.decorations, deco = this.updateDeco();
    let decoDiff = findChangedDeco(prevDeco, deco, update.changes);
    changedRanges = ChangedRange.extendWithRanges(changedRanges, decoDiff);
    let pointerSel = update.transactions.some((tr) => tr.annotation(Transaction.userEvent) == "pointerselection");
    if (this.dirty == 0 && changedRanges.length == 0 && !(update.flags & (4 | 8)) && update.state.selection.main.from >= this.view.viewport.from && update.state.selection.main.to <= this.view.viewport.to) {
      this.updateSelection(forceSelection, pointerSel);
      return false;
    } else {
      this.updateInner(changedRanges, deco, update.startState.doc.length, forceSelection, pointerSel);
      return true;
    }
  }
  updateInner(changes, deco, oldLength, forceSelection = false, pointerSel = false) {
    this.updateChildren(changes, deco, oldLength);
    this.view.observer.ignore(() => {
      this.dom.style.height = this.view.viewState.domHeight + "px";
      this.dom.style.minWidth = this.minWidth ? this.minWidth + "px" : "";
      let track = browser.chrome ? {node: getSelection(this.view.root).focusNode, written: false} : void 0;
      this.sync(track);
      this.dirty = 0;
      if (track === null || track === void 0 ? void 0 : track.written)
        forceSelection = true;
      this.updateSelection(forceSelection, pointerSel);
      this.dom.style.height = "";
    });
  }
  updateChildren(changes, deco, oldLength) {
    let cursor = this.childCursor(oldLength);
    for (let i = changes.length - 1; ; i--) {
      let next = i >= 0 ? changes[i] : null;
      if (!next)
        break;
      let {fromA, toA, fromB, toB} = next;
      let {content: content2, breakAtStart, openStart, openEnd} = ContentBuilder.build(this.view.state.doc, fromB, toB, deco);
      let {i: toI, off: toOff} = cursor.findPos(toA, 1);
      let {i: fromI, off: fromOff} = cursor.findPos(fromA, -1);
      this.replaceRange(fromI, fromOff, toI, toOff, content2, breakAtStart, openStart, openEnd);
    }
  }
  replaceRange(fromI, fromOff, toI, toOff, content2, breakAtStart, openStart, openEnd) {
    let before = this.children[fromI], last = content2.length ? content2[content2.length - 1] : null;
    let breakAtEnd = last ? last.breakAfter : breakAtStart;
    if (fromI == toI && !breakAtStart && !breakAtEnd && content2.length < 2 && before.merge(fromOff, toOff, content2.length ? last : null, fromOff == 0, openStart, openEnd))
      return;
    let after = this.children[toI];
    if (toOff < after.length || after.children.length && after.children[after.children.length - 1].length == 0) {
      if (fromI == toI) {
        after = after.split(toOff);
        toOff = 0;
      }
      if (!breakAtEnd && last && after.merge(0, toOff, last, true, 0, openEnd)) {
        content2[content2.length - 1] = after;
      } else {
        if (toOff || after.children.length && after.children[0].length == 0)
          after.merge(0, toOff, null, false, 0, openEnd);
        content2.push(after);
      }
    } else if (after.breakAfter) {
      if (last)
        last.breakAfter = 1;
      else
        breakAtStart = 1;
    }
    toI++;
    before.breakAfter = breakAtStart;
    if (fromOff > 0) {
      if (!breakAtStart && content2.length && before.merge(fromOff, before.length, content2[0], false, openStart, 0)) {
        before.breakAfter = content2.shift().breakAfter;
      } else if (fromOff < before.length || before.children.length && before.children[before.children.length - 1].length == 0) {
        before.merge(fromOff, before.length, null, false, openStart, 0);
      }
      fromI++;
    }
    while (fromI < toI && content2.length) {
      if (this.children[toI - 1].match(content2[content2.length - 1]))
        toI--, content2.pop();
      else if (this.children[fromI].match(content2[0]))
        fromI++, content2.shift();
      else
        break;
    }
    if (fromI < toI || content2.length)
      this.replaceChildren(fromI, toI, content2);
  }
  updateSelection(force = false, fromPointer = false) {
    if (!(fromPointer || this.mayControlSelection()))
      return;
    let main = this.view.state.selection.main;
    let anchor = this.domAtPos(main.anchor);
    let head = main.empty ? anchor : this.domAtPos(main.head);
    if (browser.gecko && main.empty && betweenUneditable(anchor)) {
      let dummy = document.createTextNode("");
      this.view.observer.ignore(() => anchor.node.insertBefore(dummy, anchor.node.childNodes[anchor.offset] || null));
      anchor = head = new DOMPos(dummy, 0);
      force = true;
    }
    let domSel = getSelection(this.root);
    if (force || !domSel.focusNode || browser.gecko && main.empty && nextToUneditable(domSel.focusNode, domSel.focusOffset) || !isEquivalentPosition(anchor.node, anchor.offset, domSel.anchorNode, domSel.anchorOffset) || !isEquivalentPosition(head.node, head.offset, domSel.focusNode, domSel.focusOffset)) {
      this.view.observer.ignore(() => {
        if (main.empty) {
          if (browser.gecko) {
            let nextTo = nextToUneditable(anchor.node, anchor.offset);
            if (nextTo && nextTo != (1 | 2)) {
              let text = nearbyTextNode(anchor.node, anchor.offset, nextTo == 1 ? 1 : -1);
              if (text)
                anchor = new DOMPos(text, nextTo == 1 ? 0 : text.nodeValue.length);
            }
          }
          domSel.collapse(anchor.node, anchor.offset);
          if (main.bidiLevel != null && domSel.cursorBidiLevel != null)
            domSel.cursorBidiLevel = main.bidiLevel;
        } else if (domSel.extend) {
          domSel.collapse(anchor.node, anchor.offset);
          domSel.extend(head.node, head.offset);
        } else {
          let range = document.createRange();
          if (main.anchor > main.head)
            [anchor, head] = [head, anchor];
          range.setEnd(head.node, head.offset);
          range.setStart(anchor.node, anchor.offset);
          domSel.removeAllRanges();
          domSel.addRange(range);
        }
      });
    }
    this.impreciseAnchor = anchor.precise ? null : new DOMPos(domSel.anchorNode, domSel.anchorOffset);
    this.impreciseHead = head.precise ? null : new DOMPos(domSel.focusNode, domSel.focusOffset);
  }
  enforceCursorAssoc() {
    let cursor = this.view.state.selection.main;
    let sel = getSelection(this.root);
    if (!cursor.empty || !cursor.assoc || !sel.modify)
      return;
    let line = LineView.find(this, cursor.head);
    if (!line)
      return;
    let lineStart = line.posAtStart;
    if (cursor.head == lineStart || cursor.head == lineStart + line.length)
      return;
    let before = this.coordsAt(cursor.head, -1), after = this.coordsAt(cursor.head, 1);
    if (!before || !after || before.bottom > after.top)
      return;
    let dom = this.domAtPos(cursor.head + cursor.assoc);
    sel.collapse(dom.node, dom.offset);
    sel.modify("move", cursor.assoc < 0 ? "forward" : "backward", "lineboundary");
  }
  mayControlSelection() {
    return this.view.state.facet(editable) ? this.root.activeElement == this.dom : hasSelection(this.dom, getSelection(this.root));
  }
  nearest(dom) {
    for (let cur2 = dom; cur2; ) {
      let domView = ContentView.get(cur2);
      if (domView && domView.rootView == this)
        return domView;
      cur2 = cur2.parentNode;
    }
    return null;
  }
  posFromDOM(node, offset) {
    let view = this.nearest(node);
    if (!view)
      throw new RangeError("Trying to find position for a DOM position outside of the document");
    return view.localPosFromDOM(node, offset) + view.posAtStart;
  }
  domAtPos(pos) {
    let {i, off} = this.childCursor().findPos(pos, -1);
    for (; i < this.children.length - 1; ) {
      let child = this.children[i];
      if (off < child.length || child instanceof LineView)
        break;
      i++;
      off = 0;
    }
    return this.children[i].domAtPos(off);
  }
  coordsAt(pos, side) {
    for (let off = this.length, i = this.children.length - 1; ; i--) {
      let child = this.children[i], start = off - child.breakAfter - child.length;
      if (pos > start || pos == start && (child.type == BlockType.Text || !i || this.children[i - 1].breakAfter))
        return child.coordsAt(pos - start, side);
      off = start;
    }
  }
  measureVisibleLineHeights() {
    let result = [], {from, to} = this.view.viewState.viewport;
    let minWidth = Math.max(this.view.scrollDOM.clientWidth, this.minWidth) + 1;
    for (let pos = 0, i = 0; i < this.children.length; i++) {
      let child = this.children[i], end = pos + child.length;
      if (end > to)
        break;
      if (pos >= from) {
        result.push(child.dom.getBoundingClientRect().height);
        let width = child.dom.scrollWidth;
        if (width > minWidth) {
          this.minWidth = minWidth = width;
          this.minWidthFrom = pos;
          this.minWidthTo = end;
        }
      }
      pos = end + child.breakAfter;
    }
    return result;
  }
  measureTextSize() {
    for (let child of this.children) {
      if (child instanceof LineView) {
        let measure = child.measureTextSize();
        if (measure)
          return measure;
      }
    }
    let dummy = document.createElement("div"), lineHeight, charWidth;
    dummy.className = "cm-line";
    dummy.textContent = "abc def ghi jkl mno pqr stu";
    this.view.observer.ignore(() => {
      this.dom.appendChild(dummy);
      let rect = clientRectsFor(dummy.firstChild)[0];
      lineHeight = dummy.getBoundingClientRect().height;
      charWidth = rect ? rect.width / 27 : 7;
      dummy.remove();
    });
    return {lineHeight, charWidth};
  }
  childCursor(pos = this.length) {
    let i = this.children.length;
    if (i)
      pos -= this.children[--i].length;
    return new ChildCursor(this.children, pos, i);
  }
  computeBlockGapDeco() {
    let deco = [], vs = this.view.viewState;
    for (let pos = 0, i = 0; ; i++) {
      let next = i == vs.viewports.length ? null : vs.viewports[i];
      let end = next ? next.from - 1 : this.length;
      if (end > pos) {
        let height = vs.lineAt(end, 0).bottom - vs.lineAt(pos, 0).top;
        deco.push(Decoration.replace({widget: new BlockGapWidget(height), block: true, inclusive: true}).range(pos, end));
      }
      if (!next)
        break;
      pos = next.to + 1;
    }
    return Decoration.set(deco);
  }
  updateDeco() {
    return this.decorations = [
      this.computeBlockGapDeco(),
      this.view.viewState.lineGapDeco,
      this.compositionDeco,
      ...this.view.state.facet(decorations),
      ...this.view.pluginField(PluginField.decorations)
    ];
  }
  scrollPosIntoView(pos, side) {
    let rect = this.coordsAt(pos, side);
    if (!rect)
      return;
    let mLeft = 0, mRight = 0, mTop = 0, mBottom = 0;
    for (let margins of this.view.pluginField(PluginField.scrollMargins))
      if (margins) {
        let {left, right, top: top2, bottom} = margins;
        if (left != null)
          mLeft = Math.max(mLeft, left);
        if (right != null)
          mRight = Math.max(mRight, right);
        if (top2 != null)
          mTop = Math.max(mTop, top2);
        if (bottom != null)
          mBottom = Math.max(mBottom, bottom);
      }
    scrollRectIntoView(this.dom, {
      left: rect.left - mLeft,
      top: rect.top - mTop,
      right: rect.right + mRight,
      bottom: rect.bottom + mBottom
    });
  }
};
function betweenUneditable(pos) {
  return pos.node.nodeType == 1 && pos.node.firstChild && (pos.offset == 0 || pos.node.childNodes[pos.offset - 1].contentEditable == "false") && (pos.offset < pos.node.childNodes.length || pos.node.childNodes[pos.offset].contentEditable == "false");
}
var BlockGapWidget = class extends WidgetType {
  constructor(height) {
    super();
    this.height = height;
  }
  toDOM() {
    let elt = document.createElement("div");
    this.updateDOM(elt);
    return elt;
  }
  eq(other) {
    return other.height == this.height;
  }
  updateDOM(elt) {
    elt.style.height = this.height + "px";
    return true;
  }
  get estimatedHeight() {
    return this.height;
  }
};
function computeCompositionDeco(view, changes) {
  let sel = getSelection(view.root);
  let textNode = sel.focusNode && nearbyTextNode(sel.focusNode, sel.focusOffset, 0);
  if (!textNode)
    return Decoration.none;
  let cView = view.docView.nearest(textNode);
  let from, to, topNode = textNode;
  if (cView instanceof InlineView) {
    while (cView.parent instanceof InlineView)
      cView = cView.parent;
    from = cView.posAtStart;
    to = from + cView.length;
    topNode = cView.dom;
  } else if (cView instanceof LineView) {
    while (topNode.parentNode != cView.dom)
      topNode = topNode.parentNode;
    let prev = topNode.previousSibling;
    while (prev && !ContentView.get(prev))
      prev = prev.previousSibling;
    from = to = prev ? ContentView.get(prev).posAtEnd : cView.posAtStart;
  } else {
    return Decoration.none;
  }
  let newFrom = changes.mapPos(from, 1), newTo = Math.max(newFrom, changes.mapPos(to, -1));
  let text = textNode.nodeValue, {state} = view;
  if (newTo - newFrom < text.length) {
    if (state.sliceDoc(newFrom, Math.min(state.doc.length, newFrom + text.length)) == text)
      newTo = newFrom + text.length;
    else if (state.sliceDoc(Math.max(0, newTo - text.length), newTo) == text)
      newFrom = newTo - text.length;
    else
      return Decoration.none;
  } else if (state.sliceDoc(newFrom, newTo) != text) {
    return Decoration.none;
  }
  return Decoration.set(Decoration.replace({widget: new CompositionWidget(topNode, textNode)}).range(newFrom, newTo));
}
var CompositionWidget = class extends WidgetType {
  constructor(top2, text) {
    super();
    this.top = top2;
    this.text = text;
  }
  eq(other) {
    return this.top == other.top && this.text == other.text;
  }
  toDOM() {
    return this.top;
  }
  ignoreEvent() {
    return false;
  }
  get customView() {
    return CompositionView;
  }
};
function nearbyTextNode(node, offset, side) {
  for (; ; ) {
    if (node.nodeType == 3)
      return node;
    if (node.nodeType == 1 && offset > 0 && side <= 0) {
      node = node.childNodes[offset - 1];
      offset = maxOffset(node);
    } else if (node.nodeType == 1 && offset < node.childNodes.length && side >= 0) {
      node = node.childNodes[offset];
      offset = 0;
    } else {
      return null;
    }
  }
}
function nextToUneditable(node, offset) {
  if (node.nodeType != 1)
    return 0;
  return (offset && node.childNodes[offset - 1].contentEditable == "false" ? 1 : 0) | (offset < node.childNodes.length && node.childNodes[offset].contentEditable == "false" ? 2 : 0);
}
var DecorationComparator$1 = class {
  constructor() {
    this.changes = [];
  }
  compareRange(from, to) {
    addRange(from, to, this.changes);
  }
  comparePoint(from, to) {
    addRange(from, to, this.changes);
  }
};
function findChangedDeco(a, b, diff) {
  let comp = new DecorationComparator$1();
  RangeSet.compare(a, b, diff, comp);
  return comp.changes;
}
var Direction = /* @__PURE__ */ function(Direction2) {
  Direction2[Direction2["LTR"] = 0] = "LTR";
  Direction2[Direction2["RTL"] = 1] = "RTL";
  return Direction2;
}(Direction || (Direction = {}));
var LTR = Direction.LTR;
var RTL = Direction.RTL;
function dec(str) {
  let result = [];
  for (let i = 0; i < str.length; i++)
    result.push(1 << +str[i]);
  return result;
}
var LowTypes = /* @__PURE__ */ dec("88888888888888888888888888888888888666888888787833333333337888888000000000000000000000000008888880000000000000000000000000088888888888888888888888888888888888887866668888088888663380888308888800000000000000000000000800000000000000000000000000000008");
var ArabicTypes = /* @__PURE__ */ dec("4444448826627288999999999992222222222222222222222222222222222222222222222229999999999999999999994444444444644222822222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222999999949999999229989999223333333333");
var Brackets = /* @__PURE__ */ Object.create(null);
var BracketStack = [];
for (let p of ["()", "[]", "{}"]) {
  let l = /* @__PURE__ */ p.charCodeAt(0), r = /* @__PURE__ */ p.charCodeAt(1);
  Brackets[l] = r;
  Brackets[r] = -l;
}
function charType(ch) {
  return ch <= 247 ? LowTypes[ch] : 1424 <= ch && ch <= 1524 ? 2 : 1536 <= ch && ch <= 1785 ? ArabicTypes[ch - 1536] : 1774 <= ch && ch <= 2220 ? 4 : 8192 <= ch && ch <= 8203 ? 256 : ch == 8204 ? 256 : 1;
}
var BidiRE = /[\u0590-\u05f4\u0600-\u06ff\u0700-\u08ac]/;
var BidiSpan = class {
  constructor(from, to, level) {
    this.from = from;
    this.to = to;
    this.level = level;
  }
  get dir() {
    return this.level % 2 ? RTL : LTR;
  }
  side(end, dir) {
    return this.dir == dir == end ? this.to : this.from;
  }
  static find(order, index, level, assoc) {
    let maybe = -1;
    for (let i = 0; i < order.length; i++) {
      let span = order[i];
      if (span.from <= index && span.to >= index) {
        if (span.level == level)
          return i;
        if (maybe < 0 || (assoc != 0 ? assoc < 0 ? span.from < index : span.to > index : order[maybe].level > span.level))
          maybe = i;
      }
    }
    if (maybe < 0)
      throw new RangeError("Index out of range");
    return maybe;
  }
};
var types = [];
function computeOrder(line, direction) {
  let len = line.length, outerType = direction == LTR ? 1 : 2, oppositeType = direction == LTR ? 2 : 1;
  if (!line || outerType == 1 && !BidiRE.test(line))
    return trivialOrder(len);
  for (let i = 0, prev = outerType, prevStrong = outerType; i < len; i++) {
    let type2 = charType(line.charCodeAt(i));
    if (type2 == 512)
      type2 = prev;
    else if (type2 == 8 && prevStrong == 4)
      type2 = 16;
    types[i] = type2 == 4 ? 2 : type2;
    if (type2 & 7)
      prevStrong = type2;
    prev = type2;
  }
  for (let i = 0, prev = outerType, prevStrong = outerType; i < len; i++) {
    let type2 = types[i];
    if (type2 == 128) {
      if (i < len - 1 && prev == types[i + 1] && prev & 24)
        type2 = types[i] = prev;
      else
        types[i] = 256;
    } else if (type2 == 64) {
      let end = i + 1;
      while (end < len && types[end] == 64)
        end++;
      let replace = i && prev == 8 || end < len && types[end] == 8 ? prevStrong == 1 ? 1 : 8 : 256;
      for (let j = i; j < end; j++)
        types[j] = replace;
      i = end - 1;
    } else if (type2 == 8 && prevStrong == 1) {
      types[i] = 1;
    }
    prev = type2;
    if (type2 & 7)
      prevStrong = type2;
  }
  for (let i = 0, sI = 0, context = 0, ch, br, type2; i < len; i++) {
    if (br = Brackets[ch = line.charCodeAt(i)]) {
      if (br < 0) {
        for (let sJ = sI - 3; sJ >= 0; sJ -= 3) {
          if (BracketStack[sJ + 1] == -br) {
            let flags = BracketStack[sJ + 2];
            let type3 = flags & 2 ? outerType : !(flags & 4) ? 0 : flags & 1 ? oppositeType : outerType;
            if (type3)
              types[i] = types[BracketStack[sJ]] = type3;
            sI = sJ;
            break;
          }
        }
      } else if (BracketStack.length == 189) {
        break;
      } else {
        BracketStack[sI++] = i;
        BracketStack[sI++] = ch;
        BracketStack[sI++] = context;
      }
    } else if ((type2 = types[i]) == 2 || type2 == 1) {
      let embed = type2 == outerType;
      context = embed ? 0 : 1;
      for (let sJ = sI - 3; sJ >= 0; sJ -= 3) {
        let cur2 = BracketStack[sJ + 2];
        if (cur2 & 2)
          break;
        if (embed) {
          BracketStack[sJ + 2] |= 2;
        } else {
          if (cur2 & 4)
            break;
          BracketStack[sJ + 2] |= 4;
        }
      }
    }
  }
  for (let i = 0; i < len; i++) {
    if (types[i] == 256) {
      let end = i + 1;
      while (end < len && types[end] == 256)
        end++;
      let beforeL = (i ? types[i - 1] : outerType) == 1;
      let afterL = (end < len ? types[end] : outerType) == 1;
      let replace = beforeL == afterL ? beforeL ? 1 : 2 : outerType;
      for (let j = i; j < end; j++)
        types[j] = replace;
      i = end - 1;
    }
  }
  let order = [];
  if (outerType == 1) {
    for (let i = 0; i < len; ) {
      let start = i, rtl = types[i++] != 1;
      while (i < len && rtl == (types[i] != 1))
        i++;
      if (rtl) {
        for (let j = i; j > start; ) {
          let end = j, l = types[--j] != 2;
          while (j > start && l == (types[j - 1] != 2))
            j--;
          order.push(new BidiSpan(j, end, l ? 2 : 1));
        }
      } else {
        order.push(new BidiSpan(start, i, 0));
      }
    }
  } else {
    for (let i = 0; i < len; ) {
      let start = i, rtl = types[i++] == 2;
      while (i < len && rtl == (types[i] == 2))
        i++;
      order.push(new BidiSpan(start, i, rtl ? 1 : 2));
    }
  }
  return order;
}
function trivialOrder(length) {
  return [new BidiSpan(0, length, 0)];
}
var movedOver = "";
function moveVisually(line, order, dir, start, forward) {
  var _a;
  let startIndex = start.head - line.from, spanI = -1;
  if (startIndex == 0) {
    if (!forward || !line.length)
      return null;
    if (order[0].level != dir) {
      startIndex = order[0].side(false, dir);
      spanI = 0;
    }
  } else if (startIndex == line.length) {
    if (forward)
      return null;
    let last = order[order.length - 1];
    if (last.level != dir) {
      startIndex = last.side(true, dir);
      spanI = order.length - 1;
    }
  }
  if (spanI < 0)
    spanI = BidiSpan.find(order, startIndex, (_a = start.bidiLevel) !== null && _a !== void 0 ? _a : -1, start.assoc);
  let span = order[spanI];
  if (startIndex == span.side(forward, dir)) {
    span = order[spanI += forward ? 1 : -1];
    startIndex = span.side(!forward, dir);
  }
  let indexForward = forward == (span.dir == dir);
  let nextIndex = findClusterBreak(line.text, startIndex, indexForward);
  movedOver = line.text.slice(Math.min(startIndex, nextIndex), Math.max(startIndex, nextIndex));
  if (nextIndex != span.side(forward, dir))
    return EditorSelection.cursor(nextIndex + line.from, indexForward ? -1 : 1, span.level);
  let nextSpan = spanI == (forward ? order.length - 1 : 0) ? null : order[spanI + (forward ? 1 : -1)];
  if (!nextSpan && span.level != dir)
    return EditorSelection.cursor(forward ? line.to : line.from, forward ? -1 : 1, dir);
  if (nextSpan && nextSpan.level < span.level)
    return EditorSelection.cursor(nextSpan.side(!forward, dir) + line.from, forward ? 1 : -1, nextSpan.level);
  return EditorSelection.cursor(nextIndex + line.from, forward ? -1 : 1, span.level);
}
function groupAt(state, pos, bias = 1) {
  let categorize = state.charCategorizer(pos);
  let line = state.doc.lineAt(pos), linePos = pos - line.from;
  if (line.length == 0)
    return EditorSelection.cursor(pos);
  if (linePos == 0)
    bias = 1;
  else if (linePos == line.length)
    bias = -1;
  let from = linePos, to = linePos;
  if (bias < 0)
    from = findClusterBreak(line.text, linePos, false);
  else
    to = findClusterBreak(line.text, linePos);
  let cat = categorize(line.text.slice(from, to));
  while (from > 0) {
    let prev = findClusterBreak(line.text, from, false);
    if (categorize(line.text.slice(prev, from)) != cat)
      break;
    from = prev;
  }
  while (to < line.length) {
    let next = findClusterBreak(line.text, to);
    if (categorize(line.text.slice(to, next)) != cat)
      break;
    to = next;
  }
  return EditorSelection.range(from + line.from, to + line.from);
}
function getdx(x, rect) {
  return rect.left > x ? rect.left - x : Math.max(0, x - rect.right);
}
function getdy(y, rect) {
  return rect.top > y ? rect.top - y : Math.max(0, y - rect.bottom);
}
function yOverlap(a, b) {
  return a.top < b.bottom - 1 && a.bottom > b.top + 1;
}
function upTop(rect, top2) {
  return top2 < rect.top ? {top: top2, left: rect.left, right: rect.right, bottom: rect.bottom} : rect;
}
function upBot(rect, bottom) {
  return bottom > rect.bottom ? {top: rect.top, left: rect.left, right: rect.right, bottom} : rect;
}
function domPosAtCoords(parent, x, y) {
  let closest, closestRect, closestX, closestY;
  let above, below, aboveRect, belowRect;
  for (let child = parent.firstChild; child; child = child.nextSibling) {
    let rects = clientRectsFor(child);
    for (let i = 0; i < rects.length; i++) {
      let rect = rects[i];
      if (closestRect && yOverlap(closestRect, rect))
        rect = upTop(upBot(rect, closestRect.bottom), closestRect.top);
      let dx = getdx(x, rect), dy = getdy(y, rect);
      if (dx == 0 && dy == 0)
        return child.nodeType == 3 ? domPosInText(child, x, y) : domPosAtCoords(child, x, y);
      if (!closest || closestY > dy || closestY == dy && closestX > dx) {
        closest = child;
        closestRect = rect;
        closestX = dx;
        closestY = dy;
      }
      if (dx == 0) {
        if (y > rect.bottom && (!aboveRect || aboveRect.bottom < rect.bottom)) {
          above = child;
          aboveRect = rect;
        } else if (y < rect.top && (!belowRect || belowRect.top > rect.top)) {
          below = child;
          belowRect = rect;
        }
      } else if (aboveRect && yOverlap(aboveRect, rect)) {
        aboveRect = upBot(aboveRect, rect.bottom);
      } else if (belowRect && yOverlap(belowRect, rect)) {
        belowRect = upTop(belowRect, rect.top);
      }
    }
  }
  if (aboveRect && aboveRect.bottom >= y) {
    closest = above;
    closestRect = aboveRect;
  } else if (belowRect && belowRect.top <= y) {
    closest = below;
    closestRect = belowRect;
  }
  if (!closest)
    return {node: parent, offset: 0};
  let clipX = Math.max(closestRect.left, Math.min(closestRect.right, x));
  if (closest.nodeType == 3)
    return domPosInText(closest, clipX, y);
  if (!closestX && closest.contentEditable == "true")
    return domPosAtCoords(closest, clipX, y);
  let offset = Array.prototype.indexOf.call(parent.childNodes, closest) + (x >= (closestRect.left + closestRect.right) / 2 ? 1 : 0);
  return {node: parent, offset};
}
function domPosInText(node, x, y) {
  let len = node.nodeValue.length;
  let closestOffset = -1, closestDY = 1e9, generalSide = 0;
  for (let i = 0; i < len; i++) {
    let rects = textRange(node, i, i + 1).getClientRects();
    for (let j = 0; j < rects.length; j++) {
      let rect = rects[j];
      if (rect.top == rect.bottom)
        continue;
      if (!generalSide)
        generalSide = x - rect.left;
      let dy = (rect.top > y ? rect.top - y : y - rect.bottom) - 1;
      if (rect.left - 1 <= x && rect.right + 1 >= x && dy < closestDY) {
        let right = x >= (rect.left + rect.right) / 2, after = right;
        if (browser.chrome || browser.gecko) {
          let rectBefore = textRange(node, i).getBoundingClientRect();
          if (rectBefore.left == rect.right)
            after = !right;
        }
        if (dy <= 0)
          return {node, offset: i + (after ? 1 : 0)};
        closestOffset = i + (after ? 1 : 0);
        closestDY = dy;
      }
    }
  }
  return {node, offset: closestOffset > -1 ? closestOffset : generalSide > 0 ? node.nodeValue.length : 0};
}
function posAtCoords(view, {x, y}, bias = -1) {
  let content2 = view.contentDOM.getBoundingClientRect(), block;
  let halfLine = view.defaultLineHeight / 2;
  for (let bounced = false; ; ) {
    block = view.blockAtHeight(y, content2.top);
    if (block.top > y || block.bottom < y) {
      bias = block.top > y ? -1 : 1;
      y = Math.min(block.bottom - halfLine, Math.max(block.top + halfLine, y));
      if (bounced)
        return -1;
      else
        bounced = true;
    }
    if (block.type == BlockType.Text)
      break;
    y = bias > 0 ? block.bottom + halfLine : block.top - halfLine;
  }
  let lineStart = block.from;
  if (lineStart < view.viewport.from)
    return view.viewport.from == 0 ? 0 : null;
  if (lineStart > view.viewport.to)
    return view.viewport.to == view.state.doc.length ? view.state.doc.length : null;
  x = Math.max(content2.left + 1, Math.min(content2.right - 1, x));
  let root = view.root, element = root.elementFromPoint(x, y);
  let node, offset = -1;
  if (element && view.contentDOM.contains(element) && !(view.docView.nearest(element) instanceof WidgetView)) {
    if (root.caretPositionFromPoint) {
      let pos = root.caretPositionFromPoint(x, y);
      if (pos)
        ({offsetNode: node, offset} = pos);
    } else if (root.caretRangeFromPoint) {
      let range = root.caretRangeFromPoint(x, y);
      if (range) {
        ({startContainer: node, startOffset: offset} = range);
        if (browser.safari && isSuspiciousCaretResult(node, offset, x))
          node = void 0;
      }
    }
  }
  if (!node || !view.docView.dom.contains(node)) {
    let line = LineView.find(view.docView, lineStart);
    ({node, offset} = domPosAtCoords(line.dom, x, y));
  }
  return view.docView.posFromDOM(node, offset);
}
function isSuspiciousCaretResult(node, offset, x) {
  let len;
  if (node.nodeType != 3 || offset != (len = node.nodeValue.length))
    return false;
  for (let next = node.nextSibling; next; next = node.nextSibling)
    if (next.nodeType != 1 || next.nodeName != "BR")
      return false;
  return textRange(node, len - 1, len).getBoundingClientRect().left > x;
}
function moveToLineBoundary(view, start, forward, includeWrap) {
  let line = view.state.doc.lineAt(start.head);
  let coords = !includeWrap || !view.lineWrapping ? null : view.coordsAtPos(start.assoc < 0 && start.head > line.from ? start.head - 1 : start.head);
  if (coords) {
    let editorRect = view.dom.getBoundingClientRect();
    let pos = view.posAtCoords({
      x: forward == (view.textDirection == Direction.LTR) ? editorRect.right - 1 : editorRect.left + 1,
      y: (coords.top + coords.bottom) / 2
    });
    if (pos != null)
      return EditorSelection.cursor(pos, forward ? -1 : 1);
  }
  let lineView = LineView.find(view.docView, start.head);
  let end = lineView ? forward ? lineView.posAtEnd : lineView.posAtStart : forward ? line.to : line.from;
  return EditorSelection.cursor(end, forward ? -1 : 1);
}
function moveByChar(view, start, forward, by) {
  let line = view.state.doc.lineAt(start.head), spans = view.bidiSpans(line);
  for (let cur2 = start, check = null; ; ) {
    let next = moveVisually(line, spans, view.textDirection, cur2, forward), char = movedOver;
    if (!next) {
      if (line.number == (forward ? view.state.doc.lines : 1))
        return cur2;
      char = "\n";
      line = view.state.doc.line(line.number + (forward ? 1 : -1));
      spans = view.bidiSpans(line);
      next = EditorSelection.cursor(forward ? line.from : line.to);
    }
    if (!check) {
      if (!by)
        return next;
      check = by(char);
    } else if (!check(char)) {
      return cur2;
    }
    cur2 = next;
  }
}
function byGroup(view, pos, start) {
  let categorize = view.state.charCategorizer(pos);
  let cat = categorize(start);
  return (next) => {
    let nextCat = categorize(next);
    if (cat == CharCategory.Space)
      cat = nextCat;
    return cat == nextCat;
  };
}
function moveVertically(view, start, forward, distance) {
  var _a;
  let startPos = start.head, dir = forward ? 1 : -1;
  if (startPos == (forward ? view.state.doc.length : 0))
    return EditorSelection.cursor(startPos);
  let startCoords = view.coordsAtPos(startPos);
  if (startCoords) {
    let rect = view.dom.getBoundingClientRect();
    let goal2 = (_a = start.goalColumn) !== null && _a !== void 0 ? _a : startCoords.left - rect.left;
    let resolvedGoal = rect.left + goal2;
    let dist = distance !== null && distance !== void 0 ? distance : view.defaultLineHeight >> 1;
    for (let startY = dir < 0 ? startCoords.top : startCoords.bottom, extra = 0; extra < 50; extra += 10) {
      let pos = posAtCoords(view, {x: resolvedGoal, y: startY + (dist + extra) * dir}, dir);
      if (pos == null)
        break;
      if (pos != startPos)
        return EditorSelection.cursor(pos, void 0, void 0, goal2);
    }
  }
  let {doc: doc2} = view.state, line = doc2.lineAt(startPos), tabSize = view.state.tabSize;
  let goal = start.goalColumn, goalCol = 0;
  if (goal == null) {
    for (const iter = doc2.iterRange(line.from, startPos); !iter.next().done; )
      goalCol = countColumn(iter.value, goalCol, tabSize);
    goal = goalCol * view.defaultCharacterWidth;
  } else {
    goalCol = Math.round(goal / view.defaultCharacterWidth);
  }
  if (dir < 0 && line.from == 0)
    return EditorSelection.cursor(0);
  else if (dir > 0 && line.to == doc2.length)
    return EditorSelection.cursor(line.to);
  let otherLine = doc2.line(line.number + dir);
  let result = otherLine.from;
  let seen = 0;
  for (const iter = doc2.iterRange(otherLine.from, otherLine.to); seen >= goalCol && !iter.next().done; ) {
    const {offset, leftOver} = findColumn(iter.value, seen, goalCol, tabSize);
    seen = goalCol - leftOver;
    result += offset;
  }
  return EditorSelection.cursor(result, void 0, void 0, goal);
}
var InputState = class {
  constructor(view) {
    this.lastKeyCode = 0;
    this.lastKeyTime = 0;
    this.lastIOSEnter = 0;
    this.lastIOSBackspace = 0;
    this.lastSelectionOrigin = null;
    this.lastSelectionTime = 0;
    this.lastEscPress = 0;
    this.scrollHandlers = [];
    this.registeredEvents = [];
    this.customHandlers = [];
    this.composing = -1;
    this.compositionEndedAt = 0;
    this.mouseSelection = null;
    for (let type2 in handlers) {
      let handler = handlers[type2];
      view.contentDOM.addEventListener(type2, (event) => {
        if (type2 == "keydown" && this.keydown(view, event))
          return;
        if (!eventBelongsToEditor(view, event) || this.ignoreDuringComposition(event))
          return;
        if (this.mustFlushObserver(event))
          view.observer.forceFlush();
        if (this.runCustomHandlers(type2, view, event))
          event.preventDefault();
        else
          handler(view, event);
      });
      this.registeredEvents.push(type2);
    }
    this.notifiedFocused = view.hasFocus;
    this.ensureHandlers(view);
  }
  setSelectionOrigin(origin) {
    this.lastSelectionOrigin = origin;
    this.lastSelectionTime = Date.now();
  }
  ensureHandlers(view) {
    let handlers2 = this.customHandlers = view.pluginField(domEventHandlers);
    for (let set of handlers2) {
      for (let type2 in set.handlers)
        if (this.registeredEvents.indexOf(type2) < 0 && type2 != "scroll") {
          this.registeredEvents.push(type2);
          view.contentDOM.addEventListener(type2, (event) => {
            if (!eventBelongsToEditor(view, event))
              return;
            if (this.runCustomHandlers(type2, view, event))
              event.preventDefault();
          });
        }
    }
  }
  runCustomHandlers(type2, view, event) {
    for (let set of this.customHandlers) {
      let handler = set.handlers[type2], handled = false;
      if (handler) {
        try {
          handled = handler.call(set.plugin, event, view);
        } catch (e) {
          logException(view.state, e);
        }
        if (handled || event.defaultPrevented) {
          if (browser.android && type2 == "keydown" && event.keyCode == 13)
            view.observer.flushSoon();
          return true;
        }
      }
    }
    return false;
  }
  runScrollHandlers(view, event) {
    for (let set of this.customHandlers) {
      let handler = set.handlers.scroll;
      if (handler) {
        try {
          handler.call(set.plugin, event, view);
        } catch (e) {
          logException(view.state, e);
        }
      }
    }
  }
  keydown(view, event) {
    this.lastKeyCode = event.keyCode;
    this.lastKeyTime = Date.now();
    if (this.screenKeyEvent(view, event))
      return;
    if (browser.ios && (event.keyCode == 13 || event.keyCode == 8) && !(event.ctrlKey || event.altKey || event.metaKey) && !event.synthetic) {
      this[event.keyCode == 13 ? "lastIOSEnter" : "lastIOSBackspace"] = Date.now();
      return true;
    }
    return false;
  }
  ignoreDuringComposition(event) {
    if (!/^key/.test(event.type))
      return false;
    if (this.composing > 0)
      return true;
    if (browser.safari && event.timeStamp - this.compositionEndedAt < 500) {
      this.compositionEndedAt = 0;
      return true;
    }
    return false;
  }
  screenKeyEvent(view, event) {
    let protectedTab = event.keyCode == 9 && Date.now() < this.lastEscPress + 2e3;
    if (event.keyCode == 27)
      this.lastEscPress = Date.now();
    else if (modifierCodes.indexOf(event.keyCode) < 0)
      this.lastEscPress = 0;
    return protectedTab;
  }
  mustFlushObserver(event) {
    return event.type == "keydown" && event.keyCode != 229 || event.type == "compositionend" && !browser.ios;
  }
  startMouseSelection(view, event, style) {
    if (this.mouseSelection)
      this.mouseSelection.destroy();
    this.mouseSelection = new MouseSelection(this, view, event, style);
  }
  update(update) {
    if (this.mouseSelection)
      this.mouseSelection.update(update);
    this.lastKeyCode = this.lastSelectionTime = 0;
  }
  destroy() {
    if (this.mouseSelection)
      this.mouseSelection.destroy();
  }
};
var modifierCodes = [16, 17, 18, 20, 91, 92, 224, 225];
var MouseSelection = class {
  constructor(inputState, view, startEvent, style) {
    this.inputState = inputState;
    this.view = view;
    this.startEvent = startEvent;
    this.style = style;
    let doc2 = view.contentDOM.ownerDocument;
    doc2.addEventListener("mousemove", this.move = this.move.bind(this));
    doc2.addEventListener("mouseup", this.up = this.up.bind(this));
    this.extend = startEvent.shiftKey;
    this.multiple = view.state.facet(EditorState.allowMultipleSelections) && addsSelectionRange(view, startEvent);
    this.dragMove = dragMovesSelection(view, startEvent);
    this.dragging = isInPrimarySelection(view, startEvent) ? null : false;
    if (this.dragging === false) {
      startEvent.preventDefault();
      this.select(startEvent);
    }
  }
  move(event) {
    if (event.buttons == 0)
      return this.destroy();
    if (this.dragging !== false)
      return;
    this.select(event);
  }
  up(event) {
    if (this.dragging == null)
      this.select(this.startEvent);
    if (!this.dragging)
      event.preventDefault();
    this.destroy();
  }
  destroy() {
    let doc2 = this.view.contentDOM.ownerDocument;
    doc2.removeEventListener("mousemove", this.move);
    doc2.removeEventListener("mouseup", this.up);
    this.inputState.mouseSelection = null;
  }
  select(event) {
    let selection = this.style.get(event, this.extend, this.multiple);
    if (!selection.eq(this.view.state.selection) || selection.main.assoc != this.view.state.selection.main.assoc)
      this.view.dispatch({
        selection,
        annotations: Transaction.userEvent.of("pointerselection"),
        scrollIntoView: true
      });
  }
  update(update) {
    if (update.docChanged && this.dragging)
      this.dragging = this.dragging.map(update.changes);
    this.style.update(update);
  }
};
function addsSelectionRange(view, event) {
  let facet = view.state.facet(clickAddsSelectionRange);
  return facet.length ? facet[0](event) : browser.mac ? event.metaKey : event.ctrlKey;
}
function dragMovesSelection(view, event) {
  let facet = view.state.facet(dragMovesSelection$1);
  return facet.length ? facet[0](event) : browser.mac ? !event.altKey : !event.ctrlKey;
}
function isInPrimarySelection(view, event) {
  let {main} = view.state.selection;
  if (main.empty)
    return false;
  let sel = getSelection(view.root);
  if (sel.rangeCount == 0)
    return true;
  let rects = sel.getRangeAt(0).getClientRects();
  for (let i = 0; i < rects.length; i++) {
    let rect = rects[i];
    if (rect.left <= event.clientX && rect.right >= event.clientX && rect.top <= event.clientY && rect.bottom >= event.clientY)
      return true;
  }
  return false;
}
function eventBelongsToEditor(view, event) {
  if (!event.bubbles)
    return true;
  if (event.defaultPrevented)
    return false;
  for (let node = event.target, cView; node != view.contentDOM; node = node.parentNode)
    if (!node || node.nodeType == 11 || (cView = ContentView.get(node)) && cView.ignoreEvent(event))
      return false;
  return true;
}
var handlers = /* @__PURE__ */ Object.create(null);
var brokenClipboardAPI = browser.ie && browser.ie_version < 15 || browser.ios && browser.webkit_version < 604;
function capturePaste(view) {
  let parent = view.dom.parentNode;
  if (!parent)
    return;
  let target = parent.appendChild(document.createElement("textarea"));
  target.style.cssText = "position: fixed; left: -10000px; top: 10px";
  target.focus();
  setTimeout(() => {
    view.focus();
    target.remove();
    doPaste(view, target.value);
  }, 50);
}
function doPaste(view, input) {
  let {state} = view, changes, i = 1, text = state.toText(input);
  let byLine = text.lines == state.selection.ranges.length;
  let linewise = lastLinewiseCopy && state.selection.ranges.every((r) => r.empty) && lastLinewiseCopy == text.toString();
  if (linewise) {
    let lastLine = -1;
    changes = state.changeByRange((range) => {
      let line = state.doc.lineAt(range.from);
      if (line.from == lastLine)
        return {range};
      lastLine = line.from;
      let insert2 = state.toText((byLine ? text.line(i++).text : input) + state.lineBreak);
      return {
        changes: {from: line.from, insert: insert2},
        range: EditorSelection.cursor(range.from + insert2.length)
      };
    });
  } else if (byLine) {
    changes = state.changeByRange((range) => {
      let line = text.line(i++);
      return {
        changes: {from: range.from, to: range.to, insert: line.text},
        range: EditorSelection.cursor(range.from + line.length)
      };
    });
  } else {
    changes = state.replaceSelection(text);
  }
  view.dispatch(changes, {
    annotations: Transaction.userEvent.of("paste"),
    scrollIntoView: true
  });
}
handlers.keydown = (view, event) => {
  view.inputState.setSelectionOrigin("keyboardselection");
};
var lastTouch = 0;
function mouseLikeTouchEvent(e) {
  return e.touches.length == 1 && e.touches[0].radiusX <= 1 && e.touches[0].radiusY <= 1;
}
handlers.touchstart = (view, e) => {
  if (!mouseLikeTouchEvent(e))
    lastTouch = Date.now();
  view.inputState.setSelectionOrigin("pointerselection");
};
handlers.touchmove = (view) => {
  view.inputState.setSelectionOrigin("pointerselection");
};
handlers.mousedown = (view, event) => {
  view.observer.flush();
  if (lastTouch > Date.now() - 2e3)
    return;
  let style = null;
  for (let makeStyle of view.state.facet(mouseSelectionStyle)) {
    style = makeStyle(view, event);
    if (style)
      break;
  }
  if (!style && event.button == 0)
    style = basicMouseSelection(view, event);
  if (style) {
    if (view.root.activeElement != view.contentDOM)
      view.observer.ignore(() => focusPreventScroll(view.contentDOM));
    view.inputState.startMouseSelection(view, event, style);
  }
};
function rangeForClick(view, pos, bias, type2) {
  if (type2 == 1) {
    return EditorSelection.cursor(pos, bias);
  } else if (type2 == 2) {
    return groupAt(view.state, pos, bias);
  } else {
    let visual = LineView.find(view.docView, pos), line = view.state.doc.lineAt(visual ? visual.posAtEnd : pos);
    let from = visual ? visual.posAtStart : line.from, to = visual ? visual.posAtEnd : line.to;
    if (to < view.state.doc.length && to == line.to)
      to++;
    return EditorSelection.range(from, to);
  }
}
var insideY = (y, rect) => y >= rect.top && y <= rect.bottom;
var inside = (x, y, rect) => insideY(y, rect) && x >= rect.left && x <= rect.right;
function findPositionSide(view, pos, x, y) {
  let line = LineView.find(view.docView, pos);
  if (!line)
    return 1;
  let off = pos - line.posAtStart;
  if (off == 0)
    return 1;
  if (off == line.length)
    return -1;
  let before = line.coordsAt(off, -1);
  if (before && inside(x, y, before))
    return -1;
  let after = line.coordsAt(off, 1);
  if (after && inside(x, y, after))
    return 1;
  return before && insideY(y, before) ? -1 : 1;
}
function queryPos(view, event) {
  let pos = view.posAtCoords({x: event.clientX, y: event.clientY});
  if (pos == null)
    return null;
  return {pos, bias: findPositionSide(view, pos, event.clientX, event.clientY)};
}
var BadMouseDetail = browser.ie && browser.ie_version <= 11;
var lastMouseDown = null;
var lastMouseDownCount = 0;
function getClickType(event) {
  if (!BadMouseDetail)
    return event.detail;
  let last = lastMouseDown;
  lastMouseDown = event;
  return lastMouseDownCount = !last || last.timeStamp > Date.now() - 400 && Math.abs(last.clientX - event.clientX) < 2 && Math.abs(last.clientY - event.clientY) < 2 ? (lastMouseDownCount + 1) % 3 : 1;
}
function basicMouseSelection(view, event) {
  let start = queryPos(view, event), type2 = getClickType(event);
  let startSel = view.state.selection;
  let last = start, lastEvent = event;
  return {
    update(update) {
      if (update.changes) {
        if (start)
          start.pos = update.changes.mapPos(start.pos);
        startSel = startSel.map(update.changes);
      }
    },
    get(event2, extend2, multiple) {
      let cur2;
      if (event2.clientX == lastEvent.clientX && event2.clientY == lastEvent.clientY)
        cur2 = last;
      else {
        cur2 = last = queryPos(view, event2);
        lastEvent = event2;
      }
      if (!cur2 || !start)
        return startSel;
      let range = rangeForClick(view, cur2.pos, cur2.bias, type2);
      if (start.pos != cur2.pos && !extend2) {
        let startRange = rangeForClick(view, start.pos, start.bias, type2);
        let from = Math.min(startRange.from, range.from), to = Math.max(startRange.to, range.to);
        range = from < range.from ? EditorSelection.range(from, to) : EditorSelection.range(to, from);
      }
      if (extend2)
        return startSel.replaceRange(startSel.main.extend(range.from, range.to));
      else if (multiple)
        return startSel.addRange(range);
      else
        return EditorSelection.create([range]);
    }
  };
}
handlers.dragstart = (view, event) => {
  let {selection: {main}} = view.state;
  let {mouseSelection} = view.inputState;
  if (mouseSelection)
    mouseSelection.dragging = main;
  if (event.dataTransfer) {
    event.dataTransfer.setData("Text", view.state.sliceDoc(main.from, main.to));
    event.dataTransfer.effectAllowed = "copyMove";
  }
};
handlers.drop = (view, event) => {
  if (!event.dataTransfer || !view.state.facet(editable))
    return;
  let dropPos = view.posAtCoords({x: event.clientX, y: event.clientY});
  let text = event.dataTransfer.getData("Text");
  if (dropPos == null || !text)
    return;
  event.preventDefault();
  let {mouseSelection} = view.inputState;
  let del = mouseSelection && mouseSelection.dragging && mouseSelection.dragMove ? {from: mouseSelection.dragging.from, to: mouseSelection.dragging.to} : null;
  let ins = {from: dropPos, insert: text};
  let changes = view.state.changes(del ? [del, ins] : ins);
  view.focus();
  view.dispatch({
    changes,
    selection: {anchor: changes.mapPos(dropPos, -1), head: changes.mapPos(dropPos, 1)},
    annotations: Transaction.userEvent.of("drop")
  });
};
handlers.paste = (view, event) => {
  if (!view.state.facet(editable))
    return;
  view.observer.flush();
  let data = brokenClipboardAPI ? null : event.clipboardData;
  let text = data && data.getData("text/plain");
  if (text) {
    doPaste(view, text);
    event.preventDefault();
  } else {
    capturePaste(view);
  }
};
function captureCopy(view, text) {
  let parent = view.dom.parentNode;
  if (!parent)
    return;
  let target = parent.appendChild(document.createElement("textarea"));
  target.style.cssText = "position: fixed; left: -10000px; top: 10px";
  target.value = text;
  target.focus();
  target.selectionEnd = text.length;
  target.selectionStart = 0;
  setTimeout(() => {
    target.remove();
    view.focus();
  }, 50);
}
function copiedRange(state) {
  let content2 = [], ranges = [], linewise = false;
  for (let range of state.selection.ranges)
    if (!range.empty) {
      content2.push(state.sliceDoc(range.from, range.to));
      ranges.push(range);
    }
  if (!content2.length) {
    let upto = -1;
    for (let {from} of state.selection.ranges) {
      let line = state.doc.lineAt(from);
      if (line.number > upto) {
        content2.push(line.text);
        ranges.push({from: line.from, to: Math.min(state.doc.length, line.to + 1)});
      }
      upto = line.number;
    }
    linewise = true;
  }
  return {text: content2.join(state.lineBreak), ranges, linewise};
}
var lastLinewiseCopy = null;
handlers.copy = handlers.cut = (view, event) => {
  let {text, ranges, linewise} = copiedRange(view.state);
  if (!text)
    return;
  lastLinewiseCopy = linewise ? text : null;
  let data = brokenClipboardAPI ? null : event.clipboardData;
  if (data) {
    event.preventDefault();
    data.clearData();
    data.setData("text/plain", text);
  } else {
    captureCopy(view, text);
  }
  if (event.type == "cut" && view.state.facet(editable))
    view.dispatch({
      changes: ranges,
      scrollIntoView: true,
      annotations: Transaction.userEvent.of("cut")
    });
};
handlers.focus = handlers.blur = (view) => {
  setTimeout(() => {
    if (view.hasFocus != view.inputState.notifiedFocused)
      view.update([]);
  }, 10);
};
handlers.beforeprint = (view) => {
  view.viewState.printing = true;
  view.requestMeasure();
  setTimeout(() => {
    view.viewState.printing = false;
    view.requestMeasure();
  }, 2e3);
};
function forceClearComposition(view) {
  if (view.docView.compositionDeco.size)
    view.update([]);
}
handlers.compositionstart = handlers.compositionupdate = (view) => {
  if (view.inputState.composing < 0) {
    if (view.docView.compositionDeco.size) {
      view.observer.flush();
      forceClearComposition(view);
    }
    view.inputState.composing = 0;
  }
};
handlers.compositionend = (view) => {
  view.inputState.composing = -1;
  view.inputState.compositionEndedAt = Date.now();
  setTimeout(() => {
    if (view.inputState.composing < 0)
      forceClearComposition(view);
  }, 50);
};
var wrappingWhiteSpace = ["pre-wrap", "normal", "pre-line"];
var HeightOracle = class {
  constructor() {
    this.doc = Text.empty;
    this.lineWrapping = false;
    this.direction = Direction.LTR;
    this.heightSamples = {};
    this.lineHeight = 14;
    this.charWidth = 7;
    this.lineLength = 30;
    this.heightChanged = false;
  }
  heightForGap(from, to) {
    let lines = this.doc.lineAt(to).number - this.doc.lineAt(from).number + 1;
    if (this.lineWrapping)
      lines += Math.ceil((to - from - lines * this.lineLength * 0.5) / this.lineLength);
    return this.lineHeight * lines;
  }
  heightForLine(length) {
    if (!this.lineWrapping)
      return this.lineHeight;
    let lines = 1 + Math.max(0, Math.ceil((length - this.lineLength) / (this.lineLength - 5)));
    return lines * this.lineHeight;
  }
  setDoc(doc2) {
    this.doc = doc2;
    return this;
  }
  mustRefresh(lineHeights, whiteSpace, direction) {
    let newHeight = false;
    for (let i = 0; i < lineHeights.length; i++) {
      let h = lineHeights[i];
      if (h < 0) {
        i++;
      } else if (!this.heightSamples[Math.floor(h * 10)]) {
        newHeight = true;
        this.heightSamples[Math.floor(h * 10)] = true;
      }
    }
    return newHeight || wrappingWhiteSpace.indexOf(whiteSpace) > -1 != this.lineWrapping || this.direction != direction;
  }
  refresh(whiteSpace, direction, lineHeight, charWidth, lineLength, knownHeights) {
    let lineWrapping = wrappingWhiteSpace.indexOf(whiteSpace) > -1;
    let changed = Math.round(lineHeight) != Math.round(this.lineHeight) || this.lineWrapping != lineWrapping || this.direction != direction;
    this.lineWrapping = lineWrapping;
    this.direction = direction;
    this.lineHeight = lineHeight;
    this.charWidth = charWidth;
    this.lineLength = lineLength;
    if (changed) {
      this.heightSamples = {};
      for (let i = 0; i < knownHeights.length; i++) {
        let h = knownHeights[i];
        if (h < 0)
          i++;
        else
          this.heightSamples[Math.floor(h * 10)] = true;
      }
    }
    return changed;
  }
};
var MeasuredHeights = class {
  constructor(from, heights) {
    this.from = from;
    this.heights = heights;
    this.index = 0;
  }
  get more() {
    return this.index < this.heights.length;
  }
};
var BlockInfo = class {
  constructor(from, length, top2, height, type2) {
    this.from = from;
    this.length = length;
    this.top = top2;
    this.height = height;
    this.type = type2;
  }
  get to() {
    return this.from + this.length;
  }
  get bottom() {
    return this.top + this.height;
  }
  join(other) {
    let detail = (Array.isArray(this.type) ? this.type : [this]).concat(Array.isArray(other.type) ? other.type : [other]);
    return new BlockInfo(this.from, this.length + other.length, this.top, this.height + other.height, detail);
  }
};
var QueryType = /* @__PURE__ */ function(QueryType2) {
  QueryType2[QueryType2["ByPos"] = 0] = "ByPos";
  QueryType2[QueryType2["ByHeight"] = 1] = "ByHeight";
  QueryType2[QueryType2["ByPosNoHeight"] = 2] = "ByPosNoHeight";
  return QueryType2;
}(QueryType || (QueryType = {}));
var Epsilon = 1e-4;
var HeightMap = class {
  constructor(length, height, flags = 2) {
    this.length = length;
    this.height = height;
    this.flags = flags;
  }
  get outdated() {
    return (this.flags & 2) > 0;
  }
  set outdated(value) {
    this.flags = (value ? 2 : 0) | this.flags & ~2;
  }
  setHeight(oracle, height) {
    if (this.height != height) {
      if (Math.abs(this.height - height) > Epsilon)
        oracle.heightChanged = true;
      this.height = height;
    }
  }
  replace(_from, _to, nodes) {
    return HeightMap.of(nodes);
  }
  decomposeLeft(_to, result) {
    result.push(this);
  }
  decomposeRight(_from, result) {
    result.push(this);
  }
  applyChanges(decorations2, oldDoc, oracle, changes) {
    let me = this;
    for (let i = changes.length - 1; i >= 0; i--) {
      let {fromA, toA, fromB, toB} = changes[i];
      let start = me.lineAt(fromA, QueryType.ByPosNoHeight, oldDoc, 0, 0);
      let end = start.to >= toA ? start : me.lineAt(toA, QueryType.ByPosNoHeight, oldDoc, 0, 0);
      toB += end.to - toA;
      toA = end.to;
      while (i > 0 && start.from <= changes[i - 1].toA) {
        fromA = changes[i - 1].fromA;
        fromB = changes[i - 1].fromB;
        i--;
        if (fromA < start.from)
          start = me.lineAt(fromA, QueryType.ByPosNoHeight, oldDoc, 0, 0);
      }
      fromB += start.from - fromA;
      fromA = start.from;
      let nodes = NodeBuilder.build(oracle, decorations2, fromB, toB);
      me = me.replace(fromA, toA, nodes);
    }
    return me.updateHeight(oracle, 0);
  }
  static empty() {
    return new HeightMapText(0, 0);
  }
  static of(nodes) {
    if (nodes.length == 1)
      return nodes[0];
    let i = 0, j = nodes.length, before = 0, after = 0;
    for (; ; ) {
      if (i == j) {
        if (before > after * 2) {
          let split = nodes[i - 1];
          if (split.break)
            nodes.splice(--i, 1, split.left, null, split.right);
          else
            nodes.splice(--i, 1, split.left, split.right);
          j += 1 + split.break;
          before -= split.size;
        } else if (after > before * 2) {
          let split = nodes[j];
          if (split.break)
            nodes.splice(j, 1, split.left, null, split.right);
          else
            nodes.splice(j, 1, split.left, split.right);
          j += 2 + split.break;
          after -= split.size;
        } else {
          break;
        }
      } else if (before < after) {
        let next = nodes[i++];
        if (next)
          before += next.size;
      } else {
        let next = nodes[--j];
        if (next)
          after += next.size;
      }
    }
    let brk = 0;
    if (nodes[i - 1] == null) {
      brk = 1;
      i--;
    } else if (nodes[i] == null) {
      brk = 1;
      j++;
    }
    return new HeightMapBranch(HeightMap.of(nodes.slice(0, i)), brk, HeightMap.of(nodes.slice(j)));
  }
};
HeightMap.prototype.size = 1;
var HeightMapBlock = class extends HeightMap {
  constructor(length, height, type2) {
    super(length, height);
    this.type = type2;
  }
  blockAt(_height, _doc, top2, offset) {
    return new BlockInfo(offset, this.length, top2, this.height, this.type);
  }
  lineAt(_value, _type, doc2, top2, offset) {
    return this.blockAt(0, doc2, top2, offset);
  }
  forEachLine(_from, _to, doc2, top2, offset, f) {
    f(this.blockAt(0, doc2, top2, offset));
  }
  updateHeight(oracle, offset = 0, _force = false, measured) {
    if (measured && measured.from <= offset && measured.more)
      this.setHeight(oracle, measured.heights[measured.index++]);
    this.outdated = false;
    return this;
  }
  toString() {
    return `block(${this.length})`;
  }
};
var HeightMapText = class extends HeightMapBlock {
  constructor(length, height) {
    super(length, height, BlockType.Text);
    this.collapsed = 0;
    this.widgetHeight = 0;
  }
  replace(_from, _to, nodes) {
    let node = nodes[0];
    if (nodes.length == 1 && (node instanceof HeightMapText || node instanceof HeightMapGap && node.flags & 4) && Math.abs(this.length - node.length) < 10) {
      if (node instanceof HeightMapGap)
        node = new HeightMapText(node.length, this.height);
      else
        node.height = this.height;
      if (!this.outdated)
        node.outdated = false;
      return node;
    } else {
      return HeightMap.of(nodes);
    }
  }
  updateHeight(oracle, offset = 0, force = false, measured) {
    if (measured && measured.from <= offset && measured.more)
      this.setHeight(oracle, measured.heights[measured.index++]);
    else if (force || this.outdated)
      this.setHeight(oracle, Math.max(this.widgetHeight, oracle.heightForLine(this.length - this.collapsed)));
    this.outdated = false;
    return this;
  }
  toString() {
    return `line(${this.length}${this.collapsed ? -this.collapsed : ""}${this.widgetHeight ? ":" + this.widgetHeight : ""})`;
  }
};
var HeightMapGap = class extends HeightMap {
  constructor(length) {
    super(length, 0);
  }
  lines(doc2, offset) {
    let firstLine = doc2.lineAt(offset).number, lastLine = doc2.lineAt(offset + this.length).number;
    return {firstLine, lastLine, lineHeight: this.height / (lastLine - firstLine + 1)};
  }
  blockAt(height, doc2, top2, offset) {
    let {firstLine, lastLine, lineHeight} = this.lines(doc2, offset);
    let line = Math.max(0, Math.min(lastLine - firstLine, Math.floor((height - top2) / lineHeight)));
    let {from, length} = doc2.line(firstLine + line);
    return new BlockInfo(from, length, top2 + lineHeight * line, lineHeight, BlockType.Text);
  }
  lineAt(value, type2, doc2, top2, offset) {
    if (type2 == QueryType.ByHeight)
      return this.blockAt(value, doc2, top2, offset);
    if (type2 == QueryType.ByPosNoHeight) {
      let {from: from2, to} = doc2.lineAt(value);
      return new BlockInfo(from2, to - from2, 0, 0, BlockType.Text);
    }
    let {firstLine, lineHeight} = this.lines(doc2, offset);
    let {from, length, number: number2} = doc2.lineAt(value);
    return new BlockInfo(from, length, top2 + lineHeight * (number2 - firstLine), lineHeight, BlockType.Text);
  }
  forEachLine(from, to, doc2, top2, offset, f) {
    let {firstLine, lineHeight} = this.lines(doc2, offset);
    for (let pos = Math.max(from, offset), end = Math.min(offset + this.length, to); pos <= end; ) {
      let line = doc2.lineAt(pos);
      if (pos == from)
        top2 += lineHeight * (line.number - firstLine);
      f(new BlockInfo(line.from, line.length, top2, lineHeight, BlockType.Text));
      top2 += lineHeight;
      pos = line.to + 1;
    }
  }
  replace(from, to, nodes) {
    let after = this.length - to;
    if (after > 0) {
      let last = nodes[nodes.length - 1];
      if (last instanceof HeightMapGap)
        nodes[nodes.length - 1] = new HeightMapGap(last.length + after);
      else
        nodes.push(null, new HeightMapGap(after - 1));
    }
    if (from > 0) {
      let first = nodes[0];
      if (first instanceof HeightMapGap)
        nodes[0] = new HeightMapGap(from + first.length);
      else
        nodes.unshift(new HeightMapGap(from - 1), null);
    }
    return HeightMap.of(nodes);
  }
  decomposeLeft(to, result) {
    result.push(new HeightMapGap(to - 1), null);
  }
  decomposeRight(from, result) {
    result.push(null, new HeightMapGap(this.length - from - 1));
  }
  updateHeight(oracle, offset = 0, force = false, measured) {
    let end = offset + this.length;
    if (measured && measured.from <= offset + this.length && measured.more) {
      let nodes = [], pos = Math.max(offset, measured.from);
      if (measured.from > offset)
        nodes.push(new HeightMapGap(measured.from - offset - 1).updateHeight(oracle, offset));
      while (pos <= end && measured.more) {
        let len = oracle.doc.lineAt(pos).length;
        if (nodes.length)
          nodes.push(null);
        let line = new HeightMapText(len, measured.heights[measured.index++]);
        line.outdated = false;
        nodes.push(line);
        pos += len + 1;
      }
      if (pos <= end)
        nodes.push(null, new HeightMapGap(end - pos).updateHeight(oracle, pos));
      oracle.heightChanged = true;
      return HeightMap.of(nodes);
    } else if (force || this.outdated) {
      this.setHeight(oracle, oracle.heightForGap(offset, offset + this.length));
      this.outdated = false;
    }
    return this;
  }
  toString() {
    return `gap(${this.length})`;
  }
};
var HeightMapBranch = class extends HeightMap {
  constructor(left, brk, right) {
    super(left.length + brk + right.length, left.height + right.height, brk | (left.outdated || right.outdated ? 2 : 0));
    this.left = left;
    this.right = right;
    this.size = left.size + right.size;
  }
  get break() {
    return this.flags & 1;
  }
  blockAt(height, doc2, top2, offset) {
    let mid = top2 + this.left.height;
    return height < mid || this.right.height == 0 ? this.left.blockAt(height, doc2, top2, offset) : this.right.blockAt(height, doc2, mid, offset + this.left.length + this.break);
  }
  lineAt(value, type2, doc2, top2, offset) {
    let rightTop = top2 + this.left.height, rightOffset = offset + this.left.length + this.break;
    let left = type2 == QueryType.ByHeight ? value < rightTop || this.right.height == 0 : value < rightOffset;
    let base3 = left ? this.left.lineAt(value, type2, doc2, top2, offset) : this.right.lineAt(value, type2, doc2, rightTop, rightOffset);
    if (this.break || (left ? base3.to < rightOffset : base3.from > rightOffset))
      return base3;
    let subQuery = type2 == QueryType.ByPosNoHeight ? QueryType.ByPosNoHeight : QueryType.ByPos;
    if (left)
      return base3.join(this.right.lineAt(rightOffset, subQuery, doc2, rightTop, rightOffset));
    else
      return this.left.lineAt(rightOffset, subQuery, doc2, top2, offset).join(base3);
  }
  forEachLine(from, to, doc2, top2, offset, f) {
    let rightTop = top2 + this.left.height, rightOffset = offset + this.left.length + this.break;
    if (this.break) {
      if (from < rightOffset)
        this.left.forEachLine(from, to, doc2, top2, offset, f);
      if (to >= rightOffset)
        this.right.forEachLine(from, to, doc2, rightTop, rightOffset, f);
    } else {
      let mid = this.lineAt(rightOffset, QueryType.ByPos, doc2, top2, offset);
      if (from < mid.from)
        this.left.forEachLine(from, mid.from - 1, doc2, top2, offset, f);
      if (mid.to >= from && mid.from <= to)
        f(mid);
      if (to > mid.to)
        this.right.forEachLine(mid.to + 1, to, doc2, rightTop, rightOffset, f);
    }
  }
  replace(from, to, nodes) {
    let rightStart = this.left.length + this.break;
    if (to < rightStart)
      return this.balanced(this.left.replace(from, to, nodes), this.right);
    if (from > this.left.length)
      return this.balanced(this.left, this.right.replace(from - rightStart, to - rightStart, nodes));
    let result = [];
    if (from > 0)
      this.decomposeLeft(from, result);
    let left = result.length;
    for (let node of nodes)
      result.push(node);
    if (from > 0)
      mergeGaps(result, left - 1);
    if (to < this.length) {
      let right = result.length;
      this.decomposeRight(to, result);
      mergeGaps(result, right);
    }
    return HeightMap.of(result);
  }
  decomposeLeft(to, result) {
    let left = this.left.length;
    if (to <= left)
      return this.left.decomposeLeft(to, result);
    result.push(this.left);
    if (this.break) {
      left++;
      if (to >= left)
        result.push(null);
    }
    if (to > left)
      this.right.decomposeLeft(to - left, result);
  }
  decomposeRight(from, result) {
    let left = this.left.length, right = left + this.break;
    if (from >= right)
      return this.right.decomposeRight(from - right, result);
    if (from < left)
      this.left.decomposeRight(from, result);
    if (this.break && from < right)
      result.push(null);
    result.push(this.right);
  }
  balanced(left, right) {
    if (left.size > 2 * right.size || right.size > 2 * left.size)
      return HeightMap.of(this.break ? [left, null, right] : [left, right]);
    this.left = left;
    this.right = right;
    this.height = left.height + right.height;
    this.outdated = left.outdated || right.outdated;
    this.size = left.size + right.size;
    this.length = left.length + this.break + right.length;
    return this;
  }
  updateHeight(oracle, offset = 0, force = false, measured) {
    let {left, right} = this, rightStart = offset + left.length + this.break, rebalance = null;
    if (measured && measured.from <= offset + left.length && measured.more)
      rebalance = left = left.updateHeight(oracle, offset, force, measured);
    else
      left.updateHeight(oracle, offset, force);
    if (measured && measured.from <= rightStart + right.length && measured.more)
      rebalance = right = right.updateHeight(oracle, rightStart, force, measured);
    else
      right.updateHeight(oracle, rightStart, force);
    if (rebalance)
      return this.balanced(left, right);
    this.height = this.left.height + this.right.height;
    this.outdated = false;
    return this;
  }
  toString() {
    return this.left + (this.break ? " " : "-") + this.right;
  }
};
function mergeGaps(nodes, around) {
  let before, after;
  if (nodes[around] == null && (before = nodes[around - 1]) instanceof HeightMapGap && (after = nodes[around + 1]) instanceof HeightMapGap)
    nodes.splice(around - 1, 3, new HeightMapGap(before.length + 1 + after.length));
}
var relevantWidgetHeight = 5;
var NodeBuilder = class {
  constructor(pos, oracle) {
    this.pos = pos;
    this.oracle = oracle;
    this.nodes = [];
    this.lineStart = -1;
    this.lineEnd = -1;
    this.covering = null;
    this.writtenTo = pos;
  }
  get isCovered() {
    return this.covering && this.nodes[this.nodes.length - 1] == this.covering;
  }
  span(_from, to) {
    if (this.lineStart > -1) {
      let end = Math.min(to, this.lineEnd), last = this.nodes[this.nodes.length - 1];
      if (last instanceof HeightMapText)
        last.length += end - this.pos;
      else if (end > this.pos || !this.isCovered)
        this.nodes.push(new HeightMapText(end - this.pos, -1));
      this.writtenTo = end;
      if (to > end) {
        this.nodes.push(null);
        this.writtenTo++;
        this.lineStart = -1;
      }
    }
    this.pos = to;
  }
  point(from, to, deco) {
    if (from < to || deco.heightRelevant) {
      let height = deco.widget ? Math.max(0, deco.widget.estimatedHeight) : 0;
      let len = to - from;
      if (deco.block) {
        this.addBlock(new HeightMapBlock(len, height, deco.type));
      } else if (len || height >= relevantWidgetHeight) {
        this.addLineDeco(height, len);
      }
    } else if (to > from) {
      this.span(from, to);
    }
    if (this.lineEnd > -1 && this.lineEnd < this.pos)
      this.lineEnd = this.oracle.doc.lineAt(this.pos).to;
  }
  enterLine() {
    if (this.lineStart > -1)
      return;
    let {from, to} = this.oracle.doc.lineAt(this.pos);
    this.lineStart = from;
    this.lineEnd = to;
    if (this.writtenTo < from) {
      if (this.writtenTo < from - 1 || this.nodes[this.nodes.length - 1] == null)
        this.nodes.push(this.blankContent(this.writtenTo, from - 1));
      this.nodes.push(null);
    }
    if (this.pos > from)
      this.nodes.push(new HeightMapText(this.pos - from, -1));
    this.writtenTo = this.pos;
  }
  blankContent(from, to) {
    let gap = new HeightMapGap(to - from);
    if (this.oracle.doc.lineAt(from).to == to)
      gap.flags |= 4;
    return gap;
  }
  ensureLine() {
    this.enterLine();
    let last = this.nodes.length ? this.nodes[this.nodes.length - 1] : null;
    if (last instanceof HeightMapText)
      return last;
    let line = new HeightMapText(0, -1);
    this.nodes.push(line);
    return line;
  }
  addBlock(block) {
    this.enterLine();
    if (block.type == BlockType.WidgetAfter && !this.isCovered)
      this.ensureLine();
    this.nodes.push(block);
    this.writtenTo = this.pos = this.pos + block.length;
    if (block.type != BlockType.WidgetBefore)
      this.covering = block;
  }
  addLineDeco(height, length) {
    let line = this.ensureLine();
    line.length += length;
    line.collapsed += length;
    line.widgetHeight = Math.max(line.widgetHeight, height);
    this.writtenTo = this.pos = this.pos + length;
  }
  finish(from) {
    let last = this.nodes.length == 0 ? null : this.nodes[this.nodes.length - 1];
    if (this.lineStart > -1 && !(last instanceof HeightMapText) && !this.isCovered)
      this.nodes.push(new HeightMapText(0, -1));
    else if (this.writtenTo < this.pos || last == null)
      this.nodes.push(this.blankContent(this.writtenTo, this.pos));
    let pos = from;
    for (let node of this.nodes) {
      if (node instanceof HeightMapText)
        node.updateHeight(this.oracle, pos);
      pos += node ? node.length : 1;
    }
    return this.nodes;
  }
  static build(oracle, decorations2, from, to) {
    let builder = new NodeBuilder(from, oracle);
    RangeSet.spans(decorations2, from, to, builder, 0);
    return builder.finish(from);
  }
};
function heightRelevantDecoChanges(a, b, diff) {
  let comp = new DecorationComparator();
  RangeSet.compare(a, b, diff, comp, 0);
  return comp.changes;
}
var DecorationComparator = class {
  constructor() {
    this.changes = [];
  }
  compareRange() {
  }
  comparePoint(from, to, a, b) {
    if (from < to || a && a.heightRelevant || b && b.heightRelevant)
      addRange(from, to, this.changes, 5);
  }
};
function visiblePixelRange(dom, paddingTop) {
  let rect = dom.getBoundingClientRect();
  let left = Math.max(0, rect.left), right = Math.min(innerWidth, rect.right);
  let top2 = Math.max(0, rect.top), bottom = Math.min(innerHeight, rect.bottom);
  for (let parent = dom.parentNode; parent; ) {
    if (parent.nodeType == 1) {
      if ((parent.scrollHeight > parent.clientHeight || parent.scrollWidth > parent.clientWidth) && window.getComputedStyle(parent).overflow != "visible") {
        let parentRect = parent.getBoundingClientRect();
        left = Math.max(left, parentRect.left);
        right = Math.min(right, parentRect.right);
        top2 = Math.max(top2, parentRect.top);
        bottom = Math.min(bottom, parentRect.bottom);
      }
      parent = parent.parentNode;
    } else if (parent.nodeType == 11) {
      parent = parent.host;
    } else {
      break;
    }
  }
  return {
    left: left - rect.left,
    right: right - rect.left,
    top: top2 - (rect.top + paddingTop),
    bottom: bottom - (rect.top + paddingTop)
  };
}
var LineGap = class {
  constructor(from, to, size) {
    this.from = from;
    this.to = to;
    this.size = size;
  }
  static same(a, b) {
    if (a.length != b.length)
      return false;
    for (let i = 0; i < a.length; i++) {
      let gA = a[i], gB = b[i];
      if (gA.from != gB.from || gA.to != gB.to || gA.size != gB.size)
        return false;
    }
    return true;
  }
  draw(wrapping) {
    return Decoration.replace({widget: new LineGapWidget(this.size, wrapping)}).range(this.from, this.to);
  }
};
var LineGapWidget = class extends WidgetType {
  constructor(size, vertical) {
    super();
    this.size = size;
    this.vertical = vertical;
  }
  eq(other) {
    return other.size == this.size && other.vertical == this.vertical;
  }
  toDOM() {
    let elt = document.createElement("div");
    if (this.vertical) {
      elt.style.height = this.size + "px";
    } else {
      elt.style.width = this.size + "px";
      elt.style.height = "2px";
      elt.style.display = "inline-block";
    }
    return elt;
  }
  get estimatedHeight() {
    return this.vertical ? this.size : -1;
  }
};
var ViewState = class {
  constructor(state) {
    this.state = state;
    this.pixelViewport = {left: 0, right: window.innerWidth, top: 0, bottom: 0};
    this.inView = true;
    this.paddingTop = 0;
    this.paddingBottom = 0;
    this.contentWidth = 0;
    this.heightOracle = new HeightOracle();
    this.scaler = IdScaler;
    this.scrollTo = null;
    this.printing = false;
    this.visibleRanges = [];
    this.mustEnforceCursorAssoc = false;
    this.heightMap = HeightMap.empty().applyChanges(state.facet(decorations), Text.empty, this.heightOracle.setDoc(state.doc), [new ChangedRange(0, 0, 0, state.doc.length)]);
    this.viewport = this.getViewport(0, null);
    this.updateForViewport();
    this.lineGaps = this.ensureLineGaps([]);
    this.lineGapDeco = Decoration.set(this.lineGaps.map((gap) => gap.draw(false)));
    this.computeVisibleRanges();
  }
  updateForViewport() {
    let viewports = [this.viewport], {main} = this.state.selection;
    for (let i = 0; i <= 1; i++) {
      let pos = i ? main.head : main.anchor;
      if (!viewports.some(({from, to}) => pos >= from && pos <= to)) {
        let {from, to} = this.lineAt(pos, 0);
        viewports.push(new Viewport(from, to));
      }
    }
    this.viewports = viewports.sort((a, b) => a.from - b.from);
    this.scaler = this.heightMap.height <= 7e6 ? IdScaler : new BigScaler(this.heightOracle.doc, this.heightMap, this.viewports);
  }
  update(update, scrollTo2 = null) {
    let prev = this.state;
    this.state = update.state;
    let newDeco = this.state.facet(decorations);
    let contentChanges = update.changedRanges;
    let heightChanges = ChangedRange.extendWithRanges(contentChanges, heightRelevantDecoChanges(update.startState.facet(decorations), newDeco, update ? update.changes : ChangeSet.empty(this.state.doc.length)));
    let prevHeight = this.heightMap.height;
    this.heightMap = this.heightMap.applyChanges(newDeco, prev.doc, this.heightOracle.setDoc(this.state.doc), heightChanges);
    if (this.heightMap.height != prevHeight)
      update.flags |= 2;
    let viewport = heightChanges.length ? this.mapViewport(this.viewport, update.changes) : this.viewport;
    if (scrollTo2 && (scrollTo2.head < viewport.from || scrollTo2.head > viewport.to) || !this.viewportIsAppropriate(viewport))
      viewport = this.getViewport(0, scrollTo2);
    if (!viewport.eq(this.viewport)) {
      this.viewport = viewport;
      update.flags |= 4;
    }
    this.updateForViewport();
    if (this.lineGaps.length || this.viewport.to - this.viewport.from > 15e3)
      update.flags |= this.updateLineGaps(this.ensureLineGaps(this.mapLineGaps(this.lineGaps, update.changes)));
    this.computeVisibleRanges();
    if (scrollTo2)
      this.scrollTo = scrollTo2;
    if (!this.mustEnforceCursorAssoc && update.selectionSet && update.view.lineWrapping && update.state.selection.main.empty && update.state.selection.main.assoc)
      this.mustEnforceCursorAssoc = true;
  }
  measure(docView, repeated) {
    let dom = docView.dom, whiteSpace = "", direction = Direction.LTR;
    if (!repeated) {
      let style = window.getComputedStyle(dom);
      whiteSpace = style.whiteSpace, direction = style.direction == "rtl" ? Direction.RTL : Direction.LTR;
      this.paddingTop = parseInt(style.paddingTop) || 0;
      this.paddingBottom = parseInt(style.paddingBottom) || 0;
    }
    let pixelViewport = this.printing ? {top: -1e8, bottom: 1e8, left: -1e8, right: 1e8} : visiblePixelRange(dom, this.paddingTop);
    let dTop = pixelViewport.top - this.pixelViewport.top, dBottom = pixelViewport.bottom - this.pixelViewport.bottom;
    this.pixelViewport = pixelViewport;
    this.inView = this.pixelViewport.bottom > this.pixelViewport.top && this.pixelViewport.right > this.pixelViewport.left;
    if (!this.inView)
      return 0;
    let lineHeights = docView.measureVisibleLineHeights();
    let refresh = false, bias = 0, result = 0, oracle = this.heightOracle;
    if (!repeated) {
      let contentWidth = docView.dom.clientWidth;
      if (oracle.mustRefresh(lineHeights, whiteSpace, direction) || oracle.lineWrapping && Math.abs(contentWidth - this.contentWidth) > oracle.charWidth) {
        let {lineHeight, charWidth} = docView.measureTextSize();
        refresh = oracle.refresh(whiteSpace, direction, lineHeight, charWidth, contentWidth / charWidth, lineHeights);
        if (refresh) {
          docView.minWidth = 0;
          result |= 16;
        }
      }
      if (this.contentWidth != contentWidth) {
        this.contentWidth = contentWidth;
        result |= 16;
      }
      if (dTop > 0 && dBottom > 0)
        bias = Math.max(dTop, dBottom);
      else if (dTop < 0 && dBottom < 0)
        bias = Math.min(dTop, dBottom);
    }
    oracle.heightChanged = false;
    this.heightMap = this.heightMap.updateHeight(oracle, 0, refresh, new MeasuredHeights(this.viewport.from, lineHeights));
    if (oracle.heightChanged)
      result |= 2;
    if (!this.viewportIsAppropriate(this.viewport, bias) || this.scrollTo && (this.scrollTo.head < this.viewport.from || this.scrollTo.head > this.viewport.to)) {
      let newVP = this.getViewport(bias, this.scrollTo);
      if (newVP.from != this.viewport.from || newVP.to != this.viewport.to) {
        this.viewport = newVP;
        result |= 4;
      }
    }
    this.updateForViewport();
    if (this.lineGaps.length || this.viewport.to - this.viewport.from > 15e3)
      result |= this.updateLineGaps(this.ensureLineGaps(refresh ? [] : this.lineGaps));
    this.computeVisibleRanges();
    if (this.mustEnforceCursorAssoc) {
      this.mustEnforceCursorAssoc = false;
      docView.enforceCursorAssoc();
    }
    return result;
  }
  get visibleTop() {
    return this.scaler.fromDOM(this.pixelViewport.top, 0);
  }
  get visibleBottom() {
    return this.scaler.fromDOM(this.pixelViewport.bottom, 0);
  }
  getViewport(bias, scrollTo2) {
    let marginTop = 0.5 - Math.max(-0.5, Math.min(0.5, bias / 1e3 / 2));
    let map = this.heightMap, doc2 = this.state.doc, {visibleTop, visibleBottom} = this;
    let viewport = new Viewport(map.lineAt(visibleTop - marginTop * 1e3, QueryType.ByHeight, doc2, 0, 0).from, map.lineAt(visibleBottom + (1 - marginTop) * 1e3, QueryType.ByHeight, doc2, 0, 0).to);
    if (scrollTo2) {
      if (scrollTo2.head < viewport.from) {
        let {top: newTop} = map.lineAt(scrollTo2.head, QueryType.ByPos, doc2, 0, 0);
        viewport = new Viewport(map.lineAt(newTop - 1e3 / 2, QueryType.ByHeight, doc2, 0, 0).from, map.lineAt(newTop + (visibleBottom - visibleTop) + 1e3 / 2, QueryType.ByHeight, doc2, 0, 0).to);
      } else if (scrollTo2.head > viewport.to) {
        let {bottom: newBottom} = map.lineAt(scrollTo2.head, QueryType.ByPos, doc2, 0, 0);
        viewport = new Viewport(map.lineAt(newBottom - (visibleBottom - visibleTop) - 1e3 / 2, QueryType.ByHeight, doc2, 0, 0).from, map.lineAt(newBottom + 1e3 / 2, QueryType.ByHeight, doc2, 0, 0).to);
      }
    }
    return viewport;
  }
  mapViewport(viewport, changes) {
    let from = changes.mapPos(viewport.from, -1), to = changes.mapPos(viewport.to, 1);
    return new Viewport(this.heightMap.lineAt(from, QueryType.ByPos, this.state.doc, 0, 0).from, this.heightMap.lineAt(to, QueryType.ByPos, this.state.doc, 0, 0).to);
  }
  viewportIsAppropriate({from, to}, bias = 0) {
    let {top: top2} = this.heightMap.lineAt(from, QueryType.ByPos, this.state.doc, 0, 0);
    let {bottom} = this.heightMap.lineAt(to, QueryType.ByPos, this.state.doc, 0, 0);
    let {visibleTop, visibleBottom} = this;
    return (from == 0 || top2 <= visibleTop - Math.max(10, Math.min(-bias, 250))) && (to == this.state.doc.length || bottom >= visibleBottom + Math.max(10, Math.min(bias, 250))) && (top2 > visibleTop - 2 * 1e3 && bottom < visibleBottom + 2 * 1e3);
  }
  mapLineGaps(gaps, changes) {
    if (!gaps.length || changes.empty)
      return gaps;
    let mapped = [];
    for (let gap of gaps)
      if (!changes.touchesRange(gap.from, gap.to))
        mapped.push(new LineGap(changes.mapPos(gap.from), changes.mapPos(gap.to), gap.size));
    return mapped;
  }
  ensureLineGaps(current) {
    let gaps = [];
    if (this.heightOracle.direction != Direction.LTR)
      return gaps;
    this.heightMap.forEachLine(this.viewport.from, this.viewport.to, this.state.doc, 0, 0, (line) => {
      if (line.length < 1e4)
        return;
      let structure = lineStructure(line.from, line.to, this.state);
      if (structure.total < 1e4)
        return;
      let viewFrom, viewTo;
      if (this.heightOracle.lineWrapping) {
        if (line.from != this.viewport.from)
          viewFrom = line.from;
        else
          viewFrom = findPosition(structure, (this.visibleTop - line.top) / line.height);
        if (line.to != this.viewport.to)
          viewTo = line.to;
        else
          viewTo = findPosition(structure, (this.visibleBottom - line.top) / line.height);
      } else {
        let totalWidth = structure.total * this.heightOracle.charWidth;
        viewFrom = findPosition(structure, this.pixelViewport.left / totalWidth);
        viewTo = findPosition(structure, this.pixelViewport.right / totalWidth);
      }
      let sel = this.state.selection.main;
      if (sel.from <= viewFrom && sel.to >= line.from)
        viewFrom = sel.from;
      if (sel.from <= line.to && sel.to >= viewTo)
        viewTo = sel.to;
      let gapTo = viewFrom - 1e4, gapFrom = viewTo + 1e4;
      if (gapTo > line.from + 5e3)
        gaps.push(find(current, (gap) => gap.from == line.from && gap.to > gapTo - 5e3 && gap.to < gapTo + 5e3) || new LineGap(line.from, gapTo, this.gapSize(line, gapTo, true, structure)));
      if (gapFrom < line.to - 5e3)
        gaps.push(find(current, (gap) => gap.to == line.to && gap.from > gapFrom - 5e3 && gap.from < gapFrom + 5e3) || new LineGap(gapFrom, line.to, this.gapSize(line, gapFrom, false, structure)));
    });
    return gaps;
  }
  gapSize(line, pos, start, structure) {
    if (this.heightOracle.lineWrapping) {
      let height = line.height * findFraction(structure, pos);
      return start ? height : line.height - height;
    } else {
      let ratio = findFraction(structure, pos);
      return structure.total * this.heightOracle.charWidth * (start ? ratio : 1 - ratio);
    }
  }
  updateLineGaps(gaps) {
    if (!LineGap.same(gaps, this.lineGaps)) {
      this.lineGaps = gaps;
      this.lineGapDeco = Decoration.set(gaps.map((gap) => gap.draw(this.heightOracle.lineWrapping)));
      return 8;
    }
    return 0;
  }
  computeVisibleRanges() {
    let deco = this.state.facet(decorations);
    if (this.lineGaps.length)
      deco = deco.concat(this.lineGapDeco);
    let ranges = [];
    RangeSet.spans(deco, this.viewport.from, this.viewport.to, {
      span(from, to) {
        ranges.push({from, to});
      },
      point() {
      }
    }, 20);
    this.visibleRanges = ranges;
  }
  lineAt(pos, editorTop) {
    editorTop += this.paddingTop;
    return scaleBlock(this.heightMap.lineAt(pos, QueryType.ByPos, this.state.doc, editorTop, 0), this.scaler, editorTop);
  }
  lineAtHeight(height, editorTop) {
    editorTop += this.paddingTop;
    return scaleBlock(this.heightMap.lineAt(this.scaler.fromDOM(height, editorTop), QueryType.ByHeight, this.state.doc, editorTop, 0), this.scaler, editorTop);
  }
  blockAtHeight(height, editorTop) {
    editorTop += this.paddingTop;
    return scaleBlock(this.heightMap.blockAt(this.scaler.fromDOM(height, editorTop), this.state.doc, editorTop, 0), this.scaler, editorTop);
  }
  forEachLine(from, to, f, editorTop) {
    editorTop += this.paddingTop;
    return this.heightMap.forEachLine(from, to, this.state.doc, editorTop, 0, this.scaler.scale == 1 ? f : (b) => f(scaleBlock(b, this.scaler, editorTop)));
  }
  get contentHeight() {
    return this.domHeight + this.paddingTop + this.paddingBottom;
  }
  get domHeight() {
    return this.scaler.toDOM(this.heightMap.height, this.paddingTop);
  }
};
var Viewport = class {
  constructor(from, to) {
    this.from = from;
    this.to = to;
  }
  eq(b) {
    return this.from == b.from && this.to == b.to;
  }
};
function lineStructure(from, to, state) {
  let ranges = [], pos = from, total = 0;
  RangeSet.spans(state.facet(decorations), from, to, {
    span() {
    },
    point(from2, to2) {
      if (from2 > pos) {
        ranges.push({from: pos, to: from2});
        total += from2 - pos;
      }
      pos = to2;
    }
  }, 20);
  if (pos < to) {
    ranges.push({from: pos, to});
    total += to - pos;
  }
  return {total, ranges};
}
function findPosition({total, ranges}, ratio) {
  if (ratio <= 0)
    return ranges[0].from;
  if (ratio >= 1)
    return ranges[ranges.length - 1].to;
  let dist = Math.floor(total * ratio);
  for (let i = 0; ; i++) {
    let {from, to} = ranges[i], size = to - from;
    if (dist <= size)
      return from + dist;
    dist -= size;
  }
}
function findFraction(structure, pos) {
  let counted = 0;
  for (let {from, to} of structure.ranges) {
    if (pos <= to) {
      counted += pos - from;
      break;
    }
    counted += to - from;
  }
  return counted / structure.total;
}
function find(array, f) {
  for (let val of array)
    if (f(val))
      return val;
  return void 0;
}
var IdScaler = {
  toDOM(n) {
    return n;
  },
  fromDOM(n) {
    return n;
  },
  scale: 1
};
var BigScaler = class {
  constructor(doc2, heightMap, viewports) {
    let vpHeight = 0, base3 = 0, domBase = 0;
    this.viewports = viewports.map(({from, to}) => {
      let top2 = heightMap.lineAt(from, QueryType.ByPos, doc2, 0, 0).top;
      let bottom = heightMap.lineAt(to, QueryType.ByPos, doc2, 0, 0).bottom;
      vpHeight += bottom - top2;
      return {from, to, top: top2, bottom, domTop: 0, domBottom: 0};
    });
    this.scale = (7e6 - vpHeight) / (heightMap.height - vpHeight);
    for (let obj of this.viewports) {
      obj.domTop = domBase + (obj.top - base3) * this.scale;
      domBase = obj.domBottom = obj.domTop + (obj.bottom - obj.top);
      base3 = obj.bottom;
    }
  }
  toDOM(n, top2) {
    n -= top2;
    for (let i = 0, base3 = 0, domBase = 0; ; i++) {
      let vp = i < this.viewports.length ? this.viewports[i] : null;
      if (!vp || n < vp.top)
        return domBase + (n - base3) * this.scale + top2;
      if (n <= vp.bottom)
        return vp.domTop + (n - vp.top) + top2;
      base3 = vp.bottom;
      domBase = vp.domBottom;
    }
  }
  fromDOM(n, top2) {
    n -= top2;
    for (let i = 0, base3 = 0, domBase = 0; ; i++) {
      let vp = i < this.viewports.length ? this.viewports[i] : null;
      if (!vp || n < vp.domTop)
        return base3 + (n - domBase) / this.scale + top2;
      if (n <= vp.domBottom)
        return vp.top + (n - vp.domTop) + top2;
      base3 = vp.bottom;
      domBase = vp.domBottom;
    }
  }
};
function scaleBlock(block, scaler, top2) {
  if (scaler.scale == 1)
    return block;
  let bTop = scaler.toDOM(block.top, top2), bBottom = scaler.toDOM(block.bottom, top2);
  return new BlockInfo(block.from, block.length, bTop, bBottom - bTop, Array.isArray(block.type) ? block.type.map((b) => scaleBlock(b, scaler, top2)) : block.type);
}
var theme = /* @__PURE__ */ Facet.define({combine: (strs) => strs.join(" ")});
var darkTheme = /* @__PURE__ */ Facet.define({combine: (values) => values.indexOf(true) > -1});
var baseThemeID = /* @__PURE__ */ StyleModule.newName();
var baseLightID = /* @__PURE__ */ StyleModule.newName();
var baseDarkID = /* @__PURE__ */ StyleModule.newName();
var lightDarkIDs = {"&light": "." + baseLightID, "&dark": "." + baseDarkID};
function buildTheme(main, spec, scopes) {
  return new StyleModule(spec, {
    finish(sel) {
      return /&/.test(sel) ? sel.replace(/&\w*/, (m) => {
        if (m == "&")
          return main;
        if (!scopes || !scopes[m])
          throw new RangeError(`Unsupported selector: ${m}`);
        return scopes[m];
      }) : main + " " + sel;
    }
  });
}
var baseTheme = /* @__PURE__ */ buildTheme("." + baseThemeID, {
  "&": {
    position: "relative !important",
    boxSizing: "border-box",
    "&.cm-focused": {
      outline: "1px dotted #212121"
    },
    display: "flex !important",
    flexDirection: "column"
  },
  ".cm-scroller": {
    display: "flex !important",
    alignItems: "flex-start !important",
    fontFamily: "monospace",
    lineHeight: 1.4,
    flexGrow: 2,
    overflowX: "auto",
    position: "relative",
    zIndex: 0
  },
  ".cm-content": {
    margin: 0,
    flexGrow: 2,
    minHeight: "100%",
    display: "block",
    whiteSpace: "pre",
    wordWrap: "normal",
    boxSizing: "border-box",
    padding: "4px 0",
    outline: "none"
  },
  ".cm-lineWrapping": {
    whiteSpace: "pre-wrap",
    overflowWrap: "anywhere"
  },
  "&light .cm-content": {caretColor: "black"},
  "&dark .cm-content": {caretColor: "white"},
  ".cm-line": {
    display: "block",
    padding: "0 2px 0 4px"
  },
  ".cm-selectionLayer": {
    zIndex: -1,
    contain: "size style"
  },
  ".cm-selectionBackground": {
    position: "absolute"
  },
  "&light .cm-selectionBackground": {
    background: "#d9d9d9"
  },
  "&dark .cm-selectionBackground": {
    background: "#222"
  },
  "&light.cm-focused .cm-selectionBackground": {
    background: "#d7d4f0"
  },
  "&dark.cm-focused .cm-selectionBackground": {
    background: "#233"
  },
  ".cm-cursorLayer": {
    zIndex: 100,
    contain: "size style",
    pointerEvents: "none"
  },
  "&.cm-focused .cm-cursorLayer": {
    animation: "steps(1) cm-blink 1.2s infinite"
  },
  "@keyframes cm-blink": {"0%": {}, "50%": {visibility: "hidden"}, "100%": {}},
  "@keyframes cm-blink2": {"0%": {}, "50%": {visibility: "hidden"}, "100%": {}},
  ".cm-cursor": {
    position: "absolute",
    borderLeft: "1.2px solid black",
    marginLeft: "-0.6px",
    pointerEvents: "none",
    display: "none"
  },
  "&dark .cm-cursor": {
    borderLeftColor: "#444"
  },
  "&.cm-focused .cm-cursor": {
    display: "block"
  },
  "&light .cm-activeLine": {backgroundColor: "#f3f9ff"},
  "&dark .cm-activeLine": {backgroundColor: "#223039"},
  "&light .cm-specialChar": {color: "red"},
  "&dark .cm-specialChar": {color: "#f78"},
  ".cm-tab": {
    display: "inline-block",
    overflow: "hidden",
    verticalAlign: "bottom"
  },
  ".cm-placeholder": {
    color: "#888",
    display: "inline-block"
  },
  ".cm-button": {
    verticalAlign: "middle",
    color: "inherit",
    fontSize: "70%",
    padding: ".2em 1em",
    borderRadius: "3px"
  },
  "&light .cm-button": {
    backgroundImage: "linear-gradient(#eff1f5, #d9d9df)",
    border: "1px solid #888",
    "&:active": {
      backgroundImage: "linear-gradient(#b4b4b4, #d0d3d6)"
    }
  },
  "&dark .cm-button": {
    backgroundImage: "linear-gradient(#393939, #111)",
    border: "1px solid #888",
    "&:active": {
      backgroundImage: "linear-gradient(#111, #333)"
    }
  },
  ".cm-textfield": {
    verticalAlign: "middle",
    color: "inherit",
    fontSize: "70%",
    border: "1px solid silver",
    padding: ".2em .5em"
  },
  "&light .cm-textfield": {
    backgroundColor: "white"
  },
  "&dark .cm-textfield": {
    border: "1px solid #555",
    backgroundColor: "inherit"
  }
}, lightDarkIDs);
var observeOptions = {
  childList: true,
  characterData: true,
  subtree: true,
  characterDataOldValue: true
};
var useCharData = browser.ie && browser.ie_version <= 11;
var DOMObserver = class {
  constructor(view, onChange, onScrollChanged) {
    this.view = view;
    this.onChange = onChange;
    this.onScrollChanged = onScrollChanged;
    this.active = false;
    this.ignoreSelection = new DOMSelection();
    this.delayedFlush = -1;
    this.queue = [];
    this.scrollTargets = [];
    this.intersection = null;
    this.intersecting = false;
    this.parentCheck = -1;
    this.dom = view.contentDOM;
    this.observer = new MutationObserver((mutations) => {
      for (let mut of mutations)
        this.queue.push(mut);
      if ((browser.ie && browser.ie_version <= 11 || browser.ios && view.composing) && mutations.some((m) => m.type == "childList" && m.removedNodes.length || m.type == "characterData" && m.oldValue.length > m.target.nodeValue.length))
        this.flushSoon();
      else
        this.flush();
    });
    if (useCharData)
      this.onCharData = (event) => {
        this.queue.push({
          target: event.target,
          type: "characterData",
          oldValue: event.prevValue
        });
        this.flushSoon();
      };
    this.onSelectionChange = this.onSelectionChange.bind(this);
    this.start();
    this.onScroll = this.onScroll.bind(this);
    window.addEventListener("scroll", this.onScroll);
    if (typeof IntersectionObserver == "function") {
      this.intersection = new IntersectionObserver((entries) => {
        if (this.parentCheck < 0)
          this.parentCheck = setTimeout(this.listenForScroll.bind(this), 1e3);
        if (entries[entries.length - 1].intersectionRatio > 0 != this.intersecting) {
          this.intersecting = !this.intersecting;
          if (this.intersecting != this.view.inView)
            this.onScrollChanged(document.createEvent("Event"));
        }
      }, {});
      this.intersection.observe(this.dom);
    }
    this.listenForScroll();
  }
  onScroll(e) {
    if (this.intersecting) {
      this.flush();
      this.onScrollChanged(e);
    }
  }
  onSelectionChange(event) {
    let {view} = this, sel = getSelection(view.root);
    if (view.state.facet(editable) ? view.root.activeElement != this.dom : !hasSelection(view.dom, sel))
      return;
    let context = sel.anchorNode && view.docView.nearest(sel.anchorNode);
    if (context && context.ignoreEvent(event))
      return;
    if (browser.ie && browser.ie_version <= 11 && !view.state.selection.main.empty && sel.focusNode && isEquivalentPosition(sel.focusNode, sel.focusOffset, sel.anchorNode, sel.anchorOffset))
      this.flushSoon();
    else
      this.flush();
  }
  listenForScroll() {
    this.parentCheck = -1;
    let i = 0, changed = null;
    for (let dom = this.dom; dom; ) {
      if (dom.nodeType == 1) {
        if (!changed && i < this.scrollTargets.length && this.scrollTargets[i] == dom)
          i++;
        else if (!changed)
          changed = this.scrollTargets.slice(0, i);
        if (changed)
          changed.push(dom);
        dom = dom.parentNode;
      } else if (dom.nodeType == 11) {
        dom = dom.host;
      } else {
        break;
      }
    }
    if (i < this.scrollTargets.length && !changed)
      changed = this.scrollTargets.slice(0, i);
    if (changed) {
      for (let dom of this.scrollTargets)
        dom.removeEventListener("scroll", this.onScroll);
      for (let dom of this.scrollTargets = changed)
        dom.addEventListener("scroll", this.onScroll);
    }
  }
  ignore(f) {
    if (!this.active)
      return f();
    try {
      this.stop();
      return f();
    } finally {
      this.start();
      this.clear();
    }
  }
  start() {
    if (this.active)
      return;
    this.observer.observe(this.dom, observeOptions);
    this.dom.ownerDocument.addEventListener("selectionchange", this.onSelectionChange);
    if (useCharData)
      this.dom.addEventListener("DOMCharacterDataModified", this.onCharData);
    this.active = true;
  }
  stop() {
    if (!this.active)
      return;
    this.active = false;
    this.observer.disconnect();
    this.dom.ownerDocument.removeEventListener("selectionchange", this.onSelectionChange);
    if (useCharData)
      this.dom.removeEventListener("DOMCharacterDataModified", this.onCharData);
  }
  clearSelection() {
    this.ignoreSelection.set(getSelection(this.view.root));
  }
  clear() {
    this.observer.takeRecords();
    this.queue.length = 0;
    this.clearSelection();
  }
  flushSoon() {
    if (this.delayedFlush < 0)
      this.delayedFlush = window.setTimeout(() => {
        this.delayedFlush = -1;
        this.flush();
      }, 20);
  }
  forceFlush() {
    if (this.delayedFlush >= 0) {
      window.clearTimeout(this.delayedFlush);
      this.delayedFlush = -1;
      this.flush();
    }
  }
  flush() {
    if (this.delayedFlush >= 0)
      return;
    let records = this.queue;
    for (let mut of this.observer.takeRecords())
      records.push(mut);
    if (records.length)
      this.queue = [];
    let selection = getSelection(this.view.root);
    let newSel = !this.ignoreSelection.eq(selection) && hasSelection(this.dom, selection);
    if (records.length == 0 && !newSel)
      return;
    let from = -1, to = -1, typeOver = false;
    for (let record of records) {
      let range = this.readMutation(record);
      if (!range)
        continue;
      if (range.typeOver)
        typeOver = true;
      if (from == -1) {
        ({from, to} = range);
      } else {
        from = Math.min(range.from, from);
        to = Math.max(range.to, to);
      }
    }
    let startState = this.view.state;
    if (from > -1 || newSel)
      this.onChange(from, to, typeOver);
    if (this.view.state == startState) {
      if (this.view.docView.dirty) {
        this.ignore(() => this.view.docView.sync());
        this.view.docView.dirty = 0;
      }
      this.view.docView.updateSelection();
    }
    this.clearSelection();
  }
  readMutation(rec) {
    let cView = this.view.docView.nearest(rec.target);
    if (!cView || cView.ignoreMutation(rec))
      return null;
    cView.markDirty();
    if (rec.type == "childList") {
      let childBefore = findChild(cView, rec.previousSibling || rec.target.previousSibling, -1);
      let childAfter = findChild(cView, rec.nextSibling || rec.target.nextSibling, 1);
      return {
        from: childBefore ? cView.posAfter(childBefore) : cView.posAtStart,
        to: childAfter ? cView.posBefore(childAfter) : cView.posAtEnd,
        typeOver: false
      };
    } else {
      return {from: cView.posAtStart, to: cView.posAtEnd, typeOver: rec.target.nodeValue == rec.oldValue};
    }
  }
  destroy() {
    this.stop();
    if (this.intersection)
      this.intersection.disconnect();
    for (let dom of this.scrollTargets)
      dom.removeEventListener("scroll", this.onScroll);
    window.removeEventListener("scroll", this.onScroll);
    clearTimeout(this.parentCheck);
  }
};
function findChild(cView, dom, dir) {
  while (dom) {
    let curView = ContentView.get(dom);
    if (curView && curView.parent == cView)
      return curView;
    let parent = dom.parentNode;
    dom = parent != cView.dom ? parent : dir > 0 ? dom.nextSibling : dom.previousSibling;
  }
  return null;
}
function applyDOMChange(view, start, end, typeOver) {
  let change, newSel;
  let sel = view.state.selection.main, bounds;
  if (start > -1 && (bounds = view.docView.domBoundsAround(start, end, 0))) {
    let {from, to} = bounds;
    let selPoints = view.docView.impreciseHead || view.docView.impreciseAnchor ? [] : selectionPoints(view.contentDOM, view.root);
    let reader = new DOMReader(selPoints, view);
    reader.readRange(bounds.startDOM, bounds.endDOM);
    newSel = selectionFromPoints(selPoints, from);
    let preferredPos = sel.from, preferredSide = null;
    if (view.inputState.lastKeyCode === 8 && view.inputState.lastKeyTime > Date.now() - 100 || browser.android && reader.text.length < to - from) {
      preferredPos = sel.to;
      preferredSide = "end";
    }
    let diff = findDiff(view.state.sliceDoc(from, to), reader.text, preferredPos - from, preferredSide);
    if (diff)
      change = {
        from: from + diff.from,
        to: from + diff.toA,
        insert: view.state.toText(reader.text.slice(diff.from, diff.toB))
      };
  } else if (view.hasFocus || !view.state.facet(editable)) {
    let domSel = getSelection(view.root);
    let {impreciseHead: iHead, impreciseAnchor: iAnchor} = view.docView;
    let head = iHead && iHead.node == domSel.focusNode && iHead.offset == domSel.focusOffset || !contains(view.contentDOM, domSel.focusNode) ? view.state.selection.main.head : view.docView.posFromDOM(domSel.focusNode, domSel.focusOffset);
    let anchor = iAnchor && iAnchor.node == domSel.anchorNode && iAnchor.offset == domSel.anchorOffset || !contains(view.contentDOM, domSel.anchorNode) ? view.state.selection.main.anchor : selectionCollapsed(domSel) ? head : view.docView.posFromDOM(domSel.anchorNode, domSel.anchorOffset);
    if (head != sel.head || anchor != sel.anchor)
      newSel = EditorSelection.single(anchor, head);
  }
  if (!change && !newSel)
    return;
  if (!change && typeOver && !sel.empty && newSel && newSel.main.empty)
    change = {from: sel.from, to: sel.to, insert: view.state.doc.slice(sel.from, sel.to)};
  if (change) {
    let startState = view.state;
    if (browser.android && (change.from == sel.from && change.to == sel.to && change.insert.length == 1 && change.insert.lines == 2 && dispatchKey(view, "Enter", 13) || change.from == sel.from - 1 && change.to == sel.to && change.insert.length == 0 && dispatchKey(view, "Backspace", 8) || change.from == sel.from && change.to == sel.to + 1 && change.insert.length == 0 && dispatchKey(view, "Delete", 46)) || browser.ios && (view.inputState.lastIOSEnter > Date.now() - 225 && change.insert.lines > 1 && dispatchKey(view, "Enter", 13) || view.inputState.lastIOSBackspace > Date.now() - 225 && !change.insert.length && dispatchKey(view, "Backspace", 8)))
      return;
    let text = change.insert.toString();
    if (view.state.facet(inputHandler).some((h) => h(view, change.from, change.to, text)))
      return;
    if (view.inputState.composing >= 0)
      view.inputState.composing++;
    let tr;
    if (change.from >= sel.from && change.to <= sel.to && change.to - change.from >= (sel.to - sel.from) / 3 && (!newSel || newSel.main.empty && newSel.main.from == change.from + change.insert.length)) {
      let before = sel.from < change.from ? startState.sliceDoc(sel.from, change.from) : "";
      let after = sel.to > change.to ? startState.sliceDoc(change.to, sel.to) : "";
      tr = startState.replaceSelection(view.state.toText(before + change.insert.sliceString(0, void 0, view.state.lineBreak) + after));
    } else {
      let changes = startState.changes(change);
      tr = {
        changes,
        selection: newSel && !startState.selection.main.eq(newSel.main) && newSel.main.to <= changes.newLength ? startState.selection.replaceRange(newSel.main) : void 0
      };
    }
    view.dispatch(tr, {scrollIntoView: true, annotations: Transaction.userEvent.of("input")});
  } else if (newSel && !newSel.main.eq(sel)) {
    let scrollIntoView2 = false, annotations;
    if (view.inputState.lastSelectionTime > Date.now() - 50) {
      if (view.inputState.lastSelectionOrigin == "keyboardselection")
        scrollIntoView2 = true;
      else
        annotations = Transaction.userEvent.of(view.inputState.lastSelectionOrigin);
    }
    view.dispatch({selection: newSel, scrollIntoView: scrollIntoView2, annotations});
  }
}
function findDiff(a, b, preferredPos, preferredSide) {
  let minLen = Math.min(a.length, b.length);
  let from = 0;
  while (from < minLen && a.charCodeAt(from) == b.charCodeAt(from))
    from++;
  if (from == minLen && a.length == b.length)
    return null;
  let toA = a.length, toB = b.length;
  while (toA > 0 && toB > 0 && a.charCodeAt(toA - 1) == b.charCodeAt(toB - 1)) {
    toA--;
    toB--;
  }
  if (preferredSide == "end") {
    let adjust = Math.max(0, from - Math.min(toA, toB));
    preferredPos -= toA + adjust - from;
  }
  if (toA < from && a.length < b.length) {
    let move = preferredPos <= from && preferredPos >= toA ? from - preferredPos : 0;
    from -= move;
    toB = from + (toB - toA);
    toA = from;
  } else if (toB < from) {
    let move = preferredPos <= from && preferredPos >= toB ? from - preferredPos : 0;
    from -= move;
    toA = from + (toA - toB);
    toB = from;
  }
  return {from, toA, toB};
}
var DOMReader = class {
  constructor(points, view) {
    this.points = points;
    this.view = view;
    this.text = "";
    this.lineBreak = view.state.lineBreak;
  }
  readRange(start, end) {
    if (!start)
      return;
    let parent = start.parentNode;
    for (let cur2 = start; ; ) {
      this.findPointBefore(parent, cur2);
      this.readNode(cur2);
      let next = cur2.nextSibling;
      if (next == end)
        break;
      let view = ContentView.get(cur2), nextView = ContentView.get(next);
      if ((view ? view.breakAfter : isBlockElement(cur2)) || (nextView ? nextView.breakAfter : isBlockElement(next)) && !(cur2.nodeName == "BR" && !cur2.cmIgnore))
        this.text += this.lineBreak;
      cur2 = next;
    }
    this.findPointBefore(parent, end);
  }
  readNode(node) {
    if (node.cmIgnore)
      return;
    let view = ContentView.get(node);
    let fromView = view && view.overrideDOMText;
    let text;
    if (fromView != null)
      text = fromView.sliceString(0, void 0, this.lineBreak);
    else if (node.nodeType == 3)
      text = node.nodeValue;
    else if (node.nodeName == "BR")
      text = node.nextSibling ? this.lineBreak : "";
    else if (node.nodeType == 1)
      this.readRange(node.firstChild, null);
    if (text != null) {
      this.findPointIn(node, text.length);
      this.text += text;
      if (browser.chrome && this.view.inputState.lastKeyCode == 13 && !node.nextSibling && /\n\n$/.test(this.text))
        this.text = this.text.slice(0, -1);
    }
  }
  findPointBefore(node, next) {
    for (let point of this.points)
      if (point.node == node && node.childNodes[point.offset] == next)
        point.pos = this.text.length;
  }
  findPointIn(node, maxLen) {
    for (let point of this.points)
      if (point.node == node)
        point.pos = this.text.length + Math.min(point.offset, maxLen);
  }
};
function isBlockElement(node) {
  return node.nodeType == 1 && /^(DIV|P|LI|UL|OL|BLOCKQUOTE|DD|DT|H\d|SECTION|PRE)$/.test(node.nodeName);
}
var DOMPoint = class {
  constructor(node, offset) {
    this.node = node;
    this.offset = offset;
    this.pos = -1;
  }
};
function selectionPoints(dom, root) {
  let result = [];
  if (root.activeElement != dom)
    return result;
  let {anchorNode, anchorOffset, focusNode, focusOffset} = getSelection(root);
  if (anchorNode) {
    result.push(new DOMPoint(anchorNode, anchorOffset));
    if (focusNode != anchorNode || focusOffset != anchorOffset)
      result.push(new DOMPoint(focusNode, focusOffset));
  }
  return result;
}
function selectionFromPoints(points, base3) {
  if (points.length == 0)
    return null;
  let anchor = points[0].pos, head = points.length == 2 ? points[1].pos : anchor;
  return anchor > -1 && head > -1 ? EditorSelection.single(anchor + base3, head + base3) : null;
}
function dispatchKey(view, name2, code) {
  let options = {key: name2, code: name2, keyCode: code, which: code, cancelable: true};
  let down = new KeyboardEvent("keydown", options);
  down.synthetic = true;
  view.contentDOM.dispatchEvent(down);
  let up = new KeyboardEvent("keyup", options);
  up.synthetic = true;
  view.contentDOM.dispatchEvent(up);
  return down.defaultPrevented || up.defaultPrevented;
}
var EditorView = class {
  constructor(config2 = {}) {
    this.plugins = [];
    this.editorAttrs = {};
    this.contentAttrs = {};
    this.bidiCache = [];
    this.updateState = 2;
    this.measureScheduled = -1;
    this.measureRequests = [];
    this.contentDOM = document.createElement("div");
    this.scrollDOM = document.createElement("div");
    this.scrollDOM.tabIndex = -1;
    this.scrollDOM.className = "cm-scroller";
    this.scrollDOM.appendChild(this.contentDOM);
    this.announceDOM = document.createElement("div");
    this.announceDOM.style.cssText = "position: absolute; top: -10000px";
    this.announceDOM.setAttribute("aria-live", "polite");
    this.dom = document.createElement("div");
    this.dom.appendChild(this.announceDOM);
    this.dom.appendChild(this.scrollDOM);
    this._dispatch = config2.dispatch || ((tr) => this.update([tr]));
    this.dispatch = this.dispatch.bind(this);
    this.root = config2.root || document;
    this.viewState = new ViewState(config2.state || EditorState.create());
    this.plugins = this.state.facet(viewPlugin).map((spec) => new PluginInstance(spec).update(this));
    this.observer = new DOMObserver(this, (from, to, typeOver) => {
      applyDOMChange(this, from, to, typeOver);
    }, (event) => {
      this.inputState.runScrollHandlers(this, event);
      this.measure();
    });
    this.inputState = new InputState(this);
    this.docView = new DocView(this);
    this.mountStyles();
    this.updateAttrs();
    this.updateState = 0;
    ensureGlobalHandler();
    this.requestMeasure();
    if (config2.parent)
      config2.parent.appendChild(this.dom);
  }
  get state() {
    return this.viewState.state;
  }
  get viewport() {
    return this.viewState.viewport;
  }
  get visibleRanges() {
    return this.viewState.visibleRanges;
  }
  get inView() {
    return this.viewState.inView;
  }
  get composing() {
    return this.inputState.composing > 0;
  }
  dispatch(...input) {
    this._dispatch(input.length == 1 && input[0] instanceof Transaction ? input[0] : this.state.update(...input));
  }
  update(transactions) {
    if (this.updateState != 0)
      throw new Error("Calls to EditorView.update are not allowed while an update is in progress");
    let redrawn = false, update;
    let state = this.state;
    for (let tr of transactions) {
      if (tr.startState != state)
        throw new RangeError("Trying to update state with a transaction that doesn't start from the previous state.");
      state = tr.state;
    }
    if (state.facet(EditorState.phrases) != this.state.facet(EditorState.phrases))
      return this.setState(state);
    update = new ViewUpdate(this, state, transactions);
    try {
      this.updateState = 2;
      let scrollTo2 = transactions.some((tr) => tr.scrollIntoView) ? state.selection.main : null;
      this.viewState.update(update, scrollTo2);
      this.bidiCache = CachedOrder.update(this.bidiCache, update.changes);
      if (!update.empty)
        this.updatePlugins(update);
      redrawn = this.docView.update(update);
      if (this.state.facet(styleModule) != this.styleModules)
        this.mountStyles();
      this.updateAttrs();
      this.showAnnouncements(transactions);
    } finally {
      this.updateState = 0;
    }
    if (redrawn || scrollTo || this.viewState.mustEnforceCursorAssoc)
      this.requestMeasure();
    if (!update.empty)
      for (let listener of this.state.facet(updateListener))
        listener(update);
  }
  setState(newState) {
    if (this.updateState != 0)
      throw new Error("Calls to EditorView.setState are not allowed while an update is in progress");
    this.updateState = 2;
    try {
      for (let plugin of this.plugins)
        plugin.destroy(this);
      this.viewState = new ViewState(newState);
      this.plugins = newState.facet(viewPlugin).map((spec) => new PluginInstance(spec).update(this));
      this.docView = new DocView(this);
      this.inputState.ensureHandlers(this);
      this.mountStyles();
      this.updateAttrs();
      this.bidiCache = [];
    } finally {
      this.updateState = 0;
    }
    this.requestMeasure();
  }
  updatePlugins(update) {
    let prevSpecs = update.startState.facet(viewPlugin), specs = update.state.facet(viewPlugin);
    if (prevSpecs != specs) {
      let newPlugins = [];
      for (let spec of specs) {
        let found = prevSpecs.indexOf(spec);
        if (found < 0) {
          newPlugins.push(new PluginInstance(spec));
        } else {
          let plugin = this.plugins[found];
          plugin.mustUpdate = update;
          newPlugins.push(plugin);
        }
      }
      for (let plugin of this.plugins)
        if (plugin.mustUpdate != update)
          plugin.destroy(this);
      this.plugins = newPlugins;
      this.inputState.ensureHandlers(this);
    } else {
      for (let p of this.plugins)
        p.mustUpdate = update;
    }
    for (let i = 0; i < this.plugins.length; i++)
      this.plugins[i] = this.plugins[i].update(this);
  }
  measure() {
    if (this.measureScheduled > -1)
      cancelAnimationFrame(this.measureScheduled);
    this.measureScheduled = -1;
    let updated = null;
    try {
      for (let i = 0; ; i++) {
        this.updateState = 1;
        let changed = this.viewState.measure(this.docView, i > 0);
        let measuring = this.measureRequests;
        if (!changed && !measuring.length && this.viewState.scrollTo == null)
          break;
        this.measureRequests = [];
        if (i > 5) {
          console.warn("Viewport failed to stabilize");
          break;
        }
        let measured = measuring.map((m) => {
          try {
            return m.read(this);
          } catch (e) {
            logException(this.state, e);
            return BadMeasure;
          }
        });
        let update = new ViewUpdate(this, this.state);
        update.flags |= changed;
        if (!updated)
          updated = update;
        else
          updated.flags |= changed;
        this.updateState = 2;
        if (!update.empty)
          this.updatePlugins(update);
        this.updateAttrs();
        if (changed)
          this.docView.update(update);
        for (let i2 = 0; i2 < measuring.length; i2++)
          if (measured[i2] != BadMeasure) {
            try {
              measuring[i2].write(measured[i2], this);
            } catch (e) {
              logException(this.state, e);
            }
          }
        if (this.viewState.scrollTo) {
          this.docView.scrollPosIntoView(this.viewState.scrollTo.head, this.viewState.scrollTo.assoc);
          this.viewState.scrollTo = null;
        }
        if (!(changed & 4) && this.measureRequests.length == 0)
          break;
      }
    } finally {
      this.updateState = 0;
    }
    this.measureScheduled = -1;
    if (updated && !updated.empty)
      for (let listener of this.state.facet(updateListener))
        listener(updated);
  }
  get themeClasses() {
    return baseThemeID + " " + (this.state.facet(darkTheme) ? baseDarkID : baseLightID) + " " + this.state.facet(theme);
  }
  updateAttrs() {
    let editorAttrs = combineAttrs(this.state.facet(editorAttributes), {
      class: "cm-editor cm-wrap" + (this.hasFocus ? " cm-focused " : " ") + this.themeClasses
    });
    updateAttrs(this.dom, this.editorAttrs, editorAttrs);
    this.editorAttrs = editorAttrs;
    let contentAttrs = combineAttrs(this.state.facet(contentAttributes), {
      spellcheck: "false",
      contenteditable: String(this.state.facet(editable)),
      class: "cm-content",
      style: `${browser.tabSize}: ${this.state.tabSize}`,
      role: "textbox",
      "aria-multiline": "true"
    });
    updateAttrs(this.contentDOM, this.contentAttrs, contentAttrs);
    this.contentAttrs = contentAttrs;
  }
  showAnnouncements(trs) {
    let first = true;
    for (let tr of trs)
      for (let effect of tr.effects)
        if (effect.is(EditorView.announce)) {
          if (first)
            this.announceDOM.textContent = "";
          first = false;
          let div = this.announceDOM.appendChild(document.createElement("div"));
          div.textContent = effect.value;
        }
  }
  mountStyles() {
    this.styleModules = this.state.facet(styleModule);
    StyleModule.mount(this.root, this.styleModules.concat(baseTheme).reverse());
  }
  readMeasured() {
    if (this.updateState == 2)
      throw new Error("Reading the editor layout isn't allowed during an update");
    if (this.updateState == 0 && this.measureScheduled > -1)
      this.measure();
  }
  requestMeasure(request) {
    if (this.measureScheduled < 0)
      this.measureScheduled = requestAnimationFrame(() => this.measure());
    if (request) {
      if (request.key != null)
        for (let i = 0; i < this.measureRequests.length; i++) {
          if (this.measureRequests[i].key === request.key) {
            this.measureRequests[i] = request;
            return;
          }
        }
      this.measureRequests.push(request);
    }
  }
  pluginField(field) {
    let result = [];
    for (let plugin of this.plugins)
      plugin.update(this).takeField(field, result);
    return result;
  }
  plugin(plugin) {
    for (let inst of this.plugins)
      if (inst.spec == plugin)
        return inst.update(this).value;
    return null;
  }
  blockAtHeight(height, docTop) {
    this.readMeasured();
    return this.viewState.blockAtHeight(height, ensureTop(docTop, this.contentDOM));
  }
  visualLineAtHeight(height, docTop) {
    this.readMeasured();
    return this.viewState.lineAtHeight(height, ensureTop(docTop, this.contentDOM));
  }
  viewportLines(f, docTop) {
    let {from, to} = this.viewport;
    this.viewState.forEachLine(from, to, f, ensureTop(docTop, this.contentDOM));
  }
  visualLineAt(pos, docTop = 0) {
    return this.viewState.lineAt(pos, docTop);
  }
  get contentHeight() {
    return this.viewState.contentHeight;
  }
  moveByChar(start, forward, by) {
    return moveByChar(this, start, forward, by);
  }
  moveByGroup(start, forward) {
    return moveByChar(this, start, forward, (initial) => byGroup(this, start.head, initial));
  }
  moveToLineBoundary(start, forward, includeWrap = true) {
    return moveToLineBoundary(this, start, forward, includeWrap);
  }
  moveVertically(start, forward, distance) {
    return moveVertically(this, start, forward, distance);
  }
  scrollPosIntoView(pos) {
    this.viewState.scrollTo = EditorSelection.cursor(pos);
    this.requestMeasure();
  }
  domAtPos(pos) {
    return this.docView.domAtPos(pos);
  }
  posAtDOM(node, offset = 0) {
    return this.docView.posFromDOM(node, offset);
  }
  posAtCoords(coords) {
    this.readMeasured();
    return posAtCoords(this, coords);
  }
  coordsAtPos(pos, side = 1) {
    this.readMeasured();
    let rect = this.docView.coordsAt(pos, side);
    if (!rect || rect.left == rect.right)
      return rect;
    let line = this.state.doc.lineAt(pos), order = this.bidiSpans(line);
    let span = order[BidiSpan.find(order, pos - line.from, -1, side)];
    return flattenRect(rect, span.dir == Direction.LTR == side > 0);
  }
  get defaultCharacterWidth() {
    return this.viewState.heightOracle.charWidth;
  }
  get defaultLineHeight() {
    return this.viewState.heightOracle.lineHeight;
  }
  get textDirection() {
    return this.viewState.heightOracle.direction;
  }
  get lineWrapping() {
    return this.viewState.heightOracle.lineWrapping;
  }
  bidiSpans(line) {
    if (line.length > MaxBidiLine)
      return trivialOrder(line.length);
    let dir = this.textDirection;
    for (let entry of this.bidiCache)
      if (entry.from == line.from && entry.dir == dir)
        return entry.order;
    let order = computeOrder(line.text, this.textDirection);
    this.bidiCache.push(new CachedOrder(line.from, line.to, dir, order));
    return order;
  }
  get hasFocus() {
    return document.hasFocus() && this.root.activeElement == this.contentDOM;
  }
  focus() {
    this.observer.ignore(() => {
      focusPreventScroll(this.contentDOM);
      this.docView.updateSelection();
    });
  }
  destroy() {
    for (let plugin of this.plugins)
      plugin.destroy(this);
    this.inputState.destroy();
    this.dom.remove();
    this.observer.destroy();
    if (this.measureScheduled > -1)
      cancelAnimationFrame(this.measureScheduled);
  }
  static domEventHandlers(handlers2) {
    return ViewPlugin.define(() => ({}), {eventHandlers: handlers2});
  }
  static theme(spec, options) {
    let prefix2 = StyleModule.newName();
    let result = [theme.of(prefix2), styleModule.of(buildTheme(`.${prefix2}`, spec))];
    if (options && options.dark)
      result.push(darkTheme.of(true));
    return result;
  }
  static baseTheme(spec) {
    return Prec.fallback(styleModule.of(buildTheme("." + baseThemeID, spec, lightDarkIDs)));
  }
};
EditorView.styleModule = styleModule;
EditorView.inputHandler = inputHandler;
EditorView.exceptionSink = exceptionSink;
EditorView.updateListener = updateListener;
EditorView.editable = editable;
EditorView.mouseSelectionStyle = mouseSelectionStyle;
EditorView.dragMovesSelection = dragMovesSelection$1;
EditorView.clickAddsSelectionRange = clickAddsSelectionRange;
EditorView.decorations = decorations;
EditorView.contentAttributes = contentAttributes;
EditorView.editorAttributes = editorAttributes;
EditorView.lineWrapping = /* @__PURE__ */ EditorView.contentAttributes.of({class: "cm-lineWrapping"});
EditorView.announce = /* @__PURE__ */ StateEffect.define();
var MaxBidiLine = 4096;
function ensureTop(given, dom) {
  return given == null ? dom.getBoundingClientRect().top : given;
}
var resizeDebounce = -1;
function ensureGlobalHandler() {
  window.addEventListener("resize", () => {
    if (resizeDebounce == -1)
      resizeDebounce = setTimeout(handleResize, 50);
  });
}
function handleResize() {
  resizeDebounce = -1;
  let found = document.querySelectorAll(".cm-content");
  for (let i = 0; i < found.length; i++) {
    let docView = ContentView.get(found[i]);
    if (docView)
      docView.editorView.requestMeasure();
  }
}
var BadMeasure = {};
var CachedOrder = class {
  constructor(from, to, dir, order) {
    this.from = from;
    this.to = to;
    this.dir = dir;
    this.order = order;
  }
  static update(cache, changes) {
    if (changes.empty)
      return cache;
    let result = [], lastDir = cache.length ? cache[cache.length - 1].dir : Direction.LTR;
    for (let i = Math.max(0, cache.length - 10); i < cache.length; i++) {
      let entry = cache[i];
      if (entry.dir == lastDir && !changes.touchesRange(entry.from, entry.to))
        result.push(new CachedOrder(changes.mapPos(entry.from, 1), changes.mapPos(entry.to, -1), entry.dir, entry.order));
    }
    return result;
  }
};
var currentPlatform = typeof navigator == "undefined" ? "key" : /* @__PURE__ */ /Mac/.test(navigator.platform) ? "mac" : /* @__PURE__ */ /Win/.test(navigator.platform) ? "win" : /* @__PURE__ */ /Linux|X11/.test(navigator.platform) ? "linux" : "key";
function normalizeKeyName(name2, platform) {
  const parts = name2.split(/-(?!$)/);
  let result = parts[parts.length - 1];
  if (result == "Space")
    result = " ";
  let alt, ctrl, shift2, meta2;
  for (let i = 0; i < parts.length - 1; ++i) {
    const mod = parts[i];
    if (/^(cmd|meta|m)$/i.test(mod))
      meta2 = true;
    else if (/^a(lt)?$/i.test(mod))
      alt = true;
    else if (/^(c|ctrl|control)$/i.test(mod))
      ctrl = true;
    else if (/^s(hift)?$/i.test(mod))
      shift2 = true;
    else if (/^mod$/i.test(mod)) {
      if (platform == "mac")
        meta2 = true;
      else
        ctrl = true;
    } else
      throw new Error("Unrecognized modifier name: " + mod);
  }
  if (alt)
    result = "Alt-" + result;
  if (ctrl)
    result = "Ctrl-" + result;
  if (meta2)
    result = "Meta-" + result;
  if (shift2)
    result = "Shift-" + result;
  return result;
}
function modifiers(name2, event, shift2) {
  if (event.altKey)
    name2 = "Alt-" + name2;
  if (event.ctrlKey)
    name2 = "Ctrl-" + name2;
  if (event.metaKey)
    name2 = "Meta-" + name2;
  if (shift2 !== false && event.shiftKey)
    name2 = "Shift-" + name2;
  return name2;
}
var handleKeyEvents = /* @__PURE__ */ EditorView.domEventHandlers({
  keydown(event, view) {
    return runHandlers(getKeymap(view.state), event, view, "editor");
  }
});
var keymap = /* @__PURE__ */ Facet.define({enables: handleKeyEvents});
var Keymaps = /* @__PURE__ */ new WeakMap();
function getKeymap(state) {
  let bindings = state.facet(keymap);
  let map = Keymaps.get(bindings);
  if (!map)
    Keymaps.set(bindings, map = buildKeymap(bindings.reduce((a, b) => a.concat(b), [])));
  return map;
}
var storedPrefix = null;
var PrefixTimeout = 4e3;
function buildKeymap(bindings, platform = currentPlatform) {
  let bound = Object.create(null);
  let isPrefix = Object.create(null);
  let checkPrefix = (name2, is) => {
    let current = isPrefix[name2];
    if (current == null)
      isPrefix[name2] = is;
    else if (current != is)
      throw new Error("Key binding " + name2 + " is used both as a regular binding and as a multi-stroke prefix");
  };
  let add = (scope, key, command, preventDefault) => {
    let scopeObj = bound[scope] || (bound[scope] = Object.create(null));
    let parts = key.split(/ (?!$)/).map((k) => normalizeKeyName(k, platform));
    for (let i = 1; i < parts.length; i++) {
      let prefix2 = parts.slice(0, i).join(" ");
      checkPrefix(prefix2, true);
      if (!scopeObj[prefix2])
        scopeObj[prefix2] = {
          preventDefault: true,
          commands: [(view) => {
            let ourObj = storedPrefix = {view, prefix: prefix2, scope};
            setTimeout(() => {
              if (storedPrefix == ourObj)
                storedPrefix = null;
            }, PrefixTimeout);
            return true;
          }]
        };
    }
    let full = parts.join(" ");
    checkPrefix(full, false);
    let binding = scopeObj[full] || (scopeObj[full] = {preventDefault: false, commands: []});
    binding.commands.push(command);
    if (preventDefault)
      binding.preventDefault = true;
  };
  for (let b of bindings) {
    let name2 = b[platform] || b.key;
    if (!name2)
      continue;
    for (let scope of b.scope ? b.scope.split(" ") : ["editor"]) {
      add(scope, name2, b.run, b.preventDefault);
      if (b.shift)
        add(scope, "Shift-" + name2, b.shift, b.preventDefault);
    }
  }
  return bound;
}
function runHandlers(map, event, view, scope) {
  let name2 = keyName(event), isChar = name2.length == 1 && name2 != " ";
  let prefix2 = "", fallthrough = false;
  if (storedPrefix && storedPrefix.view == view && storedPrefix.scope == scope) {
    prefix2 = storedPrefix.prefix + " ";
    if (fallthrough = modifierCodes.indexOf(event.keyCode) < 0)
      storedPrefix = null;
  }
  let runFor = (binding) => {
    if (binding) {
      for (let cmd of binding.commands)
        if (cmd(view))
          return true;
      if (binding.preventDefault)
        fallthrough = true;
    }
    return false;
  };
  let scopeObj = map[scope], baseName;
  if (scopeObj) {
    if (runFor(scopeObj[prefix2 + modifiers(name2, event, !isChar)]))
      return true;
    if (isChar && (event.shiftKey || event.altKey || event.metaKey) && (baseName = base[event.keyCode]) && baseName != name2) {
      if (runFor(scopeObj[prefix2 + modifiers(baseName, event, true)]))
        return true;
    } else if (isChar && event.shiftKey) {
      if (runFor(scopeObj[prefix2 + modifiers(name2, event, true)]))
        return true;
    }
  }
  return fallthrough;
}
var CanHidePrimary = !browser.ios;
var themeSpec = {
  ".cm-line": {
    "& ::selection": {backgroundColor: "transparent !important"},
    "&::selection": {backgroundColor: "transparent !important"}
  }
};
if (CanHidePrimary)
  themeSpec[".cm-line"].caretColor = "transparent !important";
var UnicodeRegexpSupport = /x/.unicode != null ? "gu" : "g";

// node_modules/lezer-tree/dist/tree.es.js
var DefaultBufferLength = 1024;
var nextPropID = 0;
var CachedNode = new WeakMap();
var NodeProp = class {
  constructor({deserialize} = {}) {
    this.id = nextPropID++;
    this.deserialize = deserialize || (() => {
      throw new Error("This node type doesn't define a deserialize function");
    });
  }
  static string() {
    return new NodeProp({deserialize: (str) => str});
  }
  static number() {
    return new NodeProp({deserialize: Number});
  }
  static flag() {
    return new NodeProp({deserialize: () => true});
  }
  set(propObj, value) {
    propObj[this.id] = value;
    return propObj;
  }
  add(match) {
    if (typeof match != "function")
      match = NodeType.match(match);
    return (type2) => {
      let result = match(type2);
      return result === void 0 ? null : [this, result];
    };
  }
};
NodeProp.closedBy = new NodeProp({deserialize: (str) => str.split(" ")});
NodeProp.openedBy = new NodeProp({deserialize: (str) => str.split(" ")});
NodeProp.group = new NodeProp({deserialize: (str) => str.split(" ")});
var noProps = Object.create(null);
var NodeType = class {
  constructor(name2, props, id2, flags = 0) {
    this.name = name2;
    this.props = props;
    this.id = id2;
    this.flags = flags;
  }
  static define(spec) {
    let props = spec.props && spec.props.length ? Object.create(null) : noProps;
    let flags = (spec.top ? 1 : 0) | (spec.skipped ? 2 : 0) | (spec.error ? 4 : 0) | (spec.name == null ? 8 : 0);
    let type2 = new NodeType(spec.name || "", props, spec.id, flags);
    if (spec.props)
      for (let src of spec.props) {
        if (!Array.isArray(src))
          src = src(type2);
        if (src)
          src[0].set(props, src[1]);
      }
    return type2;
  }
  prop(prop) {
    return this.props[prop.id];
  }
  get isTop() {
    return (this.flags & 1) > 0;
  }
  get isSkipped() {
    return (this.flags & 2) > 0;
  }
  get isError() {
    return (this.flags & 4) > 0;
  }
  get isAnonymous() {
    return (this.flags & 8) > 0;
  }
  is(name2) {
    if (typeof name2 == "string") {
      if (this.name == name2)
        return true;
      let group = this.prop(NodeProp.group);
      return group ? group.indexOf(name2) > -1 : false;
    }
    return this.id == name2;
  }
  static match(map) {
    let direct = Object.create(null);
    for (let prop in map)
      for (let name2 of prop.split(" "))
        direct[name2] = map[prop];
    return (node) => {
      for (let groups = node.prop(NodeProp.group), i = -1; i < (groups ? groups.length : 0); i++) {
        let found = direct[i < 0 ? node.name : groups[i]];
        if (found)
          return found;
      }
    };
  }
};
NodeType.none = new NodeType("", Object.create(null), 0, 8);
var NodeSet = class {
  constructor(types4) {
    this.types = types4;
    for (let i = 0; i < types4.length; i++)
      if (types4[i].id != i)
        throw new RangeError("Node type ids should correspond to array positions when creating a node set");
  }
  extend(...props) {
    let newTypes = [];
    for (let type2 of this.types) {
      let newProps = null;
      for (let source of props) {
        let add = source(type2);
        if (add) {
          if (!newProps)
            newProps = Object.assign({}, type2.props);
          add[0].set(newProps, add[1]);
        }
      }
      newTypes.push(newProps ? new NodeType(type2.name, newProps, type2.id, type2.flags) : type2);
    }
    return new NodeSet(newTypes);
  }
};
var Tree = class {
  constructor(type2, children, positions, length) {
    this.type = type2;
    this.children = children;
    this.positions = positions;
    this.length = length;
  }
  toString() {
    let children = this.children.map((c2) => c2.toString()).join();
    return !this.type.name ? children : (/\W/.test(this.type.name) && !this.type.isError ? JSON.stringify(this.type.name) : this.type.name) + (children.length ? "(" + children + ")" : "");
  }
  cursor(pos, side = 0) {
    let scope = pos != null && CachedNode.get(this) || this.topNode;
    let cursor = new TreeCursor(scope);
    if (pos != null) {
      cursor.moveTo(pos, side);
      CachedNode.set(this, cursor._tree);
    }
    return cursor;
  }
  fullCursor() {
    return new TreeCursor(this.topNode, true);
  }
  get topNode() {
    return new TreeNode(this, 0, 0, null);
  }
  resolve(pos, side = 0) {
    return this.cursor(pos, side).node;
  }
  iterate(spec) {
    let {enter, leave, from = 0, to = this.length} = spec;
    for (let c2 = this.cursor(); ; ) {
      let mustLeave = false;
      if (c2.from <= to && c2.to >= from && (c2.type.isAnonymous || enter(c2.type, c2.from, c2.to) !== false)) {
        if (c2.firstChild())
          continue;
        if (!c2.type.isAnonymous)
          mustLeave = true;
      }
      for (; ; ) {
        if (mustLeave && leave)
          leave(c2.type, c2.from, c2.to);
        mustLeave = c2.type.isAnonymous;
        if (c2.nextSibling())
          break;
        if (!c2.parent())
          return;
        mustLeave = true;
      }
    }
  }
  balance(maxBufferLength = DefaultBufferLength) {
    return this.children.length <= BalanceBranchFactor ? this : balanceRange(this.type, NodeType.none, this.children, this.positions, 0, this.children.length, 0, maxBufferLength, this.length, 0);
  }
  static build(data) {
    return buildTree(data);
  }
};
Tree.empty = new Tree(NodeType.none, [], [], 0);
function withHash(tree, hash2) {
  if (hash2)
    tree.contextHash = hash2;
  return tree;
}
var TreeBuffer = class {
  constructor(buffer, length, set, type2 = NodeType.none) {
    this.buffer = buffer;
    this.length = length;
    this.set = set;
    this.type = type2;
  }
  toString() {
    let result = [];
    for (let index = 0; index < this.buffer.length; ) {
      result.push(this.childString(index));
      index = this.buffer[index + 3];
    }
    return result.join(",");
  }
  childString(index) {
    let id2 = this.buffer[index], endIndex = this.buffer[index + 3];
    let type2 = this.set.types[id2], result = type2.name;
    if (/\W/.test(result) && !type2.isError)
      result = JSON.stringify(result);
    index += 4;
    if (endIndex == index)
      return result;
    let children = [];
    while (index < endIndex) {
      children.push(this.childString(index));
      index = this.buffer[index + 3];
    }
    return result + "(" + children.join(",") + ")";
  }
  findChild(startIndex, endIndex, dir, after) {
    let {buffer} = this, pick = -1;
    for (let i = startIndex; i != endIndex; i = buffer[i + 3]) {
      if (after != -1e8) {
        let start = buffer[i + 1], end = buffer[i + 2];
        if (dir > 0) {
          if (end > after)
            pick = i;
          if (end > after)
            break;
        } else {
          if (start < after)
            pick = i;
          if (end >= after)
            break;
        }
      } else {
        pick = i;
        if (dir > 0)
          break;
      }
    }
    return pick;
  }
};
var TreeNode = class {
  constructor(node, from, index, _parent) {
    this.node = node;
    this.from = from;
    this.index = index;
    this._parent = _parent;
  }
  get type() {
    return this.node.type;
  }
  get name() {
    return this.node.type.name;
  }
  get to() {
    return this.from + this.node.length;
  }
  nextChild(i, dir, after, full = false) {
    for (let parent = this; ; ) {
      for (let {children, positions} = parent.node, e = dir > 0 ? children.length : -1; i != e; i += dir) {
        let next = children[i], start = positions[i] + parent.from;
        if (after != -1e8 && (dir < 0 ? start >= after : start + next.length <= after))
          continue;
        if (next instanceof TreeBuffer) {
          let index = next.findChild(0, next.buffer.length, dir, after == -1e8 ? -1e8 : after - start);
          if (index > -1)
            return new BufferNode(new BufferContext(parent, next, i, start), null, index);
        } else if (full || (!next.type.isAnonymous || hasChild(next))) {
          let inner = new TreeNode(next, start, i, parent);
          return full || !inner.type.isAnonymous ? inner : inner.nextChild(dir < 0 ? next.children.length - 1 : 0, dir, after);
        }
      }
      if (full || !parent.type.isAnonymous)
        return null;
      i = parent.index + dir;
      parent = parent._parent;
      if (!parent)
        return null;
    }
  }
  get firstChild() {
    return this.nextChild(0, 1, -1e8);
  }
  get lastChild() {
    return this.nextChild(this.node.children.length - 1, -1, -1e8);
  }
  childAfter(pos) {
    return this.nextChild(0, 1, pos);
  }
  childBefore(pos) {
    return this.nextChild(this.node.children.length - 1, -1, pos);
  }
  nextSignificantParent() {
    let val = this;
    while (val.type.isAnonymous && val._parent)
      val = val._parent;
    return val;
  }
  get parent() {
    return this._parent ? this._parent.nextSignificantParent() : null;
  }
  get nextSibling() {
    return this._parent ? this._parent.nextChild(this.index + 1, 1, -1) : null;
  }
  get prevSibling() {
    return this._parent ? this._parent.nextChild(this.index - 1, -1, -1) : null;
  }
  get cursor() {
    return new TreeCursor(this);
  }
  resolve(pos, side = 0) {
    return this.cursor.moveTo(pos, side).node;
  }
  getChild(type2, before = null, after = null) {
    let r = getChildren(this, type2, before, after);
    return r.length ? r[0] : null;
  }
  getChildren(type2, before = null, after = null) {
    return getChildren(this, type2, before, after);
  }
  toString() {
    return this.node.toString();
  }
};
function getChildren(node, type2, before, after) {
  let cur2 = node.cursor, result = [];
  if (!cur2.firstChild())
    return result;
  if (before != null) {
    while (!cur2.type.is(before))
      if (!cur2.nextSibling())
        return result;
  }
  for (; ; ) {
    if (after != null && cur2.type.is(after))
      return result;
    if (cur2.type.is(type2))
      result.push(cur2.node);
    if (!cur2.nextSibling())
      return after == null ? result : [];
  }
}
var BufferContext = class {
  constructor(parent, buffer, index, start) {
    this.parent = parent;
    this.buffer = buffer;
    this.index = index;
    this.start = start;
  }
};
var BufferNode = class {
  constructor(context, _parent, index) {
    this.context = context;
    this._parent = _parent;
    this.index = index;
    this.type = context.buffer.set.types[context.buffer.buffer[index]];
  }
  get name() {
    return this.type.name;
  }
  get from() {
    return this.context.start + this.context.buffer.buffer[this.index + 1];
  }
  get to() {
    return this.context.start + this.context.buffer.buffer[this.index + 2];
  }
  child(dir, after) {
    let {buffer} = this.context;
    let index = buffer.findChild(this.index + 4, buffer.buffer[this.index + 3], dir, after == -1e8 ? -1e8 : after - this.context.start);
    return index < 0 ? null : new BufferNode(this.context, this, index);
  }
  get firstChild() {
    return this.child(1, -1e8);
  }
  get lastChild() {
    return this.child(-1, -1e8);
  }
  childAfter(pos) {
    return this.child(1, pos);
  }
  childBefore(pos) {
    return this.child(-1, pos);
  }
  get parent() {
    return this._parent || this.context.parent.nextSignificantParent();
  }
  externalSibling(dir) {
    return this._parent ? null : this.context.parent.nextChild(this.context.index + dir, dir, -1);
  }
  get nextSibling() {
    let {buffer} = this.context;
    let after = buffer.buffer[this.index + 3];
    if (after < (this._parent ? buffer.buffer[this._parent.index + 3] : buffer.buffer.length))
      return new BufferNode(this.context, this._parent, after);
    return this.externalSibling(1);
  }
  get prevSibling() {
    let {buffer} = this.context;
    let parentStart = this._parent ? this._parent.index + 4 : 0;
    if (this.index == parentStart)
      return this.externalSibling(-1);
    return new BufferNode(this.context, this._parent, buffer.findChild(parentStart, this.index, -1, -1e8));
  }
  get cursor() {
    return new TreeCursor(this);
  }
  resolve(pos, side = 0) {
    return this.cursor.moveTo(pos, side).node;
  }
  toString() {
    return this.context.buffer.childString(this.index);
  }
  getChild(type2, before = null, after = null) {
    let r = getChildren(this, type2, before, after);
    return r.length ? r[0] : null;
  }
  getChildren(type2, before = null, after = null) {
    return getChildren(this, type2, before, after);
  }
};
var TreeCursor = class {
  constructor(node, full = false) {
    this.full = full;
    this.buffer = null;
    this.stack = [];
    this.index = 0;
    this.bufferNode = null;
    if (node instanceof TreeNode) {
      this.yieldNode(node);
    } else {
      this._tree = node.context.parent;
      this.buffer = node.context;
      for (let n = node._parent; n; n = n._parent)
        this.stack.unshift(n.index);
      this.bufferNode = node;
      this.yieldBuf(node.index);
    }
  }
  get name() {
    return this.type.name;
  }
  yieldNode(node) {
    if (!node)
      return false;
    this._tree = node;
    this.type = node.type;
    this.from = node.from;
    this.to = node.to;
    return true;
  }
  yieldBuf(index, type2) {
    this.index = index;
    let {start, buffer} = this.buffer;
    this.type = type2 || buffer.set.types[buffer.buffer[index]];
    this.from = start + buffer.buffer[index + 1];
    this.to = start + buffer.buffer[index + 2];
    return true;
  }
  yield(node) {
    if (!node)
      return false;
    if (node instanceof TreeNode) {
      this.buffer = null;
      return this.yieldNode(node);
    }
    this.buffer = node.context;
    return this.yieldBuf(node.index, node.type);
  }
  toString() {
    return this.buffer ? this.buffer.buffer.childString(this.index) : this._tree.toString();
  }
  enter(dir, after) {
    if (!this.buffer)
      return this.yield(this._tree.nextChild(dir < 0 ? this._tree.node.children.length - 1 : 0, dir, after, this.full));
    let {buffer} = this.buffer;
    let index = buffer.findChild(this.index + 4, buffer.buffer[this.index + 3], dir, after == -1e8 ? -1e8 : after - this.buffer.start);
    if (index < 0)
      return false;
    this.stack.push(this.index);
    return this.yieldBuf(index);
  }
  firstChild() {
    return this.enter(1, -1e8);
  }
  lastChild() {
    return this.enter(-1, -1e8);
  }
  childAfter(pos) {
    return this.enter(1, pos);
  }
  childBefore(pos) {
    return this.enter(-1, pos);
  }
  parent() {
    if (!this.buffer)
      return this.yieldNode(this.full ? this._tree._parent : this._tree.parent);
    if (this.stack.length)
      return this.yieldBuf(this.stack.pop());
    let parent = this.full ? this.buffer.parent : this.buffer.parent.nextSignificantParent();
    this.buffer = null;
    return this.yieldNode(parent);
  }
  sibling(dir) {
    if (!this.buffer)
      return !this._tree._parent ? false : this.yield(this._tree._parent.nextChild(this._tree.index + dir, dir, -1e8, this.full));
    let {buffer} = this.buffer, d = this.stack.length - 1;
    if (dir < 0) {
      let parentStart = d < 0 ? 0 : this.stack[d] + 4;
      if (this.index != parentStart)
        return this.yieldBuf(buffer.findChild(parentStart, this.index, -1, -1e8));
    } else {
      let after = buffer.buffer[this.index + 3];
      if (after < (d < 0 ? buffer.buffer.length : buffer.buffer[this.stack[d] + 3]))
        return this.yieldBuf(after);
    }
    return d < 0 ? this.yield(this.buffer.parent.nextChild(this.buffer.index + dir, dir, -1e8, this.full)) : false;
  }
  nextSibling() {
    return this.sibling(1);
  }
  prevSibling() {
    return this.sibling(-1);
  }
  atLastNode(dir) {
    let index, parent, {buffer} = this;
    if (buffer) {
      if (dir > 0) {
        if (this.index < buffer.buffer.buffer.length)
          return false;
      } else {
        for (let i = 0; i < this.index; i++)
          if (buffer.buffer.buffer[i + 3] < this.index)
            return false;
      }
      ({index, parent} = buffer);
    } else {
      ({index, _parent: parent} = this._tree);
    }
    for (; parent; {index, _parent: parent} = parent) {
      for (let i = index + dir, e = dir < 0 ? -1 : parent.node.children.length; i != e; i += dir) {
        let child = parent.node.children[i];
        if (this.full || !child.type.isAnonymous || child instanceof TreeBuffer || hasChild(child))
          return false;
      }
    }
    return true;
  }
  move(dir) {
    if (this.enter(dir, -1e8))
      return true;
    for (; ; ) {
      if (this.sibling(dir))
        return true;
      if (this.atLastNode(dir) || !this.parent())
        return false;
    }
  }
  next() {
    return this.move(1);
  }
  prev() {
    return this.move(-1);
  }
  moveTo(pos, side = 0) {
    while (this.from == this.to || (side < 1 ? this.from >= pos : this.from > pos) || (side > -1 ? this.to <= pos : this.to < pos))
      if (!this.parent())
        break;
    for (; ; ) {
      if (side < 0 ? !this.childBefore(pos) : !this.childAfter(pos))
        break;
      if (this.from == this.to || (side < 1 ? this.from >= pos : this.from > pos) || (side > -1 ? this.to <= pos : this.to < pos)) {
        this.parent();
        break;
      }
    }
    return this;
  }
  get node() {
    if (!this.buffer)
      return this._tree;
    let cache = this.bufferNode, result = null, depth = 0;
    if (cache && cache.context == this.buffer) {
      scan:
        for (let index = this.index, d = this.stack.length; d >= 0; ) {
          for (let c2 = cache; c2; c2 = c2._parent)
            if (c2.index == index) {
              if (index == this.index)
                return c2;
              result = c2;
              depth = d + 1;
              break scan;
            }
          index = this.stack[--d];
        }
    }
    for (let i = depth; i < this.stack.length; i++)
      result = new BufferNode(this.buffer, result, this.stack[i]);
    return this.bufferNode = new BufferNode(this.buffer, result, this.index);
  }
  get tree() {
    return this.buffer ? null : this._tree.node;
  }
};
function hasChild(tree) {
  return tree.children.some((ch) => !ch.type.isAnonymous || ch instanceof TreeBuffer || hasChild(ch));
}
var FlatBufferCursor = class {
  constructor(buffer, index) {
    this.buffer = buffer;
    this.index = index;
  }
  get id() {
    return this.buffer[this.index - 4];
  }
  get start() {
    return this.buffer[this.index - 3];
  }
  get end() {
    return this.buffer[this.index - 2];
  }
  get size() {
    return this.buffer[this.index - 1];
  }
  get pos() {
    return this.index;
  }
  next() {
    this.index -= 4;
  }
  fork() {
    return new FlatBufferCursor(this.buffer, this.index);
  }
};
var BalanceBranchFactor = 8;
function buildTree(data) {
  var _a;
  let {buffer, nodeSet: nodeSet2, topID = 0, maxBufferLength = DefaultBufferLength, reused = [], minRepeatType = nodeSet2.types.length} = data;
  let cursor = Array.isArray(buffer) ? new FlatBufferCursor(buffer, buffer.length) : buffer;
  let types4 = nodeSet2.types;
  let contextHash = 0;
  function takeNode(parentStart, minPos, children2, positions2, inRepeat) {
    let {id: id2, start, end, size} = cursor;
    let startPos = start - parentStart;
    if (size < 0) {
      if (size == -1) {
        children2.push(reused[id2]);
        positions2.push(startPos);
      } else {
        contextHash = id2;
      }
      cursor.next();
      return;
    }
    let type2 = types4[id2], node, buffer2;
    if (end - start <= maxBufferLength && (buffer2 = findBufferSize(cursor.pos - minPos, inRepeat))) {
      let data2 = new Uint16Array(buffer2.size - buffer2.skip);
      let endPos = cursor.pos - buffer2.size, index = data2.length;
      while (cursor.pos > endPos)
        index = copyToBuffer(buffer2.start, data2, index, inRepeat);
      node = new TreeBuffer(data2, end - buffer2.start, nodeSet2, inRepeat < 0 ? NodeType.none : types4[inRepeat]);
      startPos = buffer2.start - parentStart;
    } else {
      let endPos = cursor.pos - size;
      cursor.next();
      let localChildren = [], localPositions = [];
      let localInRepeat = id2 >= minRepeatType ? id2 : -1;
      while (cursor.pos > endPos) {
        if (cursor.id == localInRepeat)
          cursor.next();
        else
          takeNode(start, endPos, localChildren, localPositions, localInRepeat);
      }
      localChildren.reverse();
      localPositions.reverse();
      if (localInRepeat > -1 && localChildren.length > BalanceBranchFactor)
        node = balanceRange(type2, type2, localChildren, localPositions, 0, localChildren.length, 0, maxBufferLength, end - start, contextHash);
      else
        node = withHash(new Tree(type2, localChildren, localPositions, end - start), contextHash);
    }
    children2.push(node);
    positions2.push(startPos);
  }
  function findBufferSize(maxSize, inRepeat) {
    let fork = cursor.fork();
    let size = 0, start = 0, skip = 0, minStart = fork.end - maxBufferLength;
    let result = {size: 0, start: 0, skip: 0};
    scan:
      for (let minPos = fork.pos - maxSize; fork.pos > minPos; ) {
        if (fork.id == inRepeat) {
          result.size = size;
          result.start = start;
          result.skip = skip;
          skip += 4;
          size += 4;
          fork.next();
          continue;
        }
        let nodeSize = fork.size, startPos = fork.pos - nodeSize;
        if (nodeSize < 0 || startPos < minPos || fork.start < minStart)
          break;
        let localSkipped = fork.id >= minRepeatType ? 4 : 0;
        let nodeStart2 = fork.start;
        fork.next();
        while (fork.pos > startPos) {
          if (fork.size < 0)
            break scan;
          if (fork.id >= minRepeatType)
            localSkipped += 4;
          fork.next();
        }
        start = nodeStart2;
        size += nodeSize;
        skip += localSkipped;
      }
    if (inRepeat < 0 || size == maxSize) {
      result.size = size;
      result.start = start;
      result.skip = skip;
    }
    return result.size > 4 ? result : void 0;
  }
  function copyToBuffer(bufferStart, buffer2, index, inRepeat) {
    let {id: id2, start, end, size} = cursor;
    cursor.next();
    if (id2 == inRepeat)
      return index;
    let startIndex = index;
    if (size > 4) {
      let endPos = cursor.pos - (size - 4);
      while (cursor.pos > endPos)
        index = copyToBuffer(bufferStart, buffer2, index, inRepeat);
    }
    if (id2 < minRepeatType) {
      buffer2[--index] = startIndex;
      buffer2[--index] = end - bufferStart;
      buffer2[--index] = start - bufferStart;
      buffer2[--index] = id2;
    }
    return index;
  }
  let children = [], positions = [];
  while (cursor.pos > 0)
    takeNode(data.start || 0, 0, children, positions, -1);
  let length = (_a = data.length) !== null && _a !== void 0 ? _a : children.length ? positions[0] + children[0].length : 0;
  return new Tree(types4[topID], children.reverse(), positions.reverse(), length);
}
function balanceRange(outerType, innerType, children, positions, from, to, start, maxBufferLength, length, contextHash) {
  let localChildren = [], localPositions = [];
  if (length <= maxBufferLength) {
    for (let i = from; i < to; i++) {
      localChildren.push(children[i]);
      localPositions.push(positions[i] - start);
    }
  } else {
    let maxChild = Math.max(maxBufferLength, Math.ceil(length * 1.5 / BalanceBranchFactor));
    for (let i = from; i < to; ) {
      let groupFrom = i, groupStart = positions[i];
      i++;
      for (; i < to; i++) {
        let nextEnd = positions[i] + children[i].length;
        if (nextEnd - groupStart > maxChild)
          break;
      }
      if (i == groupFrom + 1) {
        let only = children[groupFrom];
        if (only instanceof Tree && only.type == innerType && only.length > maxChild << 1) {
          for (let j = 0; j < only.children.length; j++) {
            localChildren.push(only.children[j]);
            localPositions.push(only.positions[j] + groupStart - start);
          }
          continue;
        }
        localChildren.push(only);
      } else if (i == groupFrom + 1) {
        localChildren.push(children[groupFrom]);
      } else {
        let inner = balanceRange(innerType, innerType, children, positions, groupFrom, i, groupStart, maxBufferLength, positions[i - 1] + children[i - 1].length - groupStart, contextHash);
        if (innerType != NodeType.none && !containsType(inner.children, innerType))
          inner = withHash(new Tree(NodeType.none, inner.children, inner.positions, inner.length), contextHash);
        localChildren.push(inner);
      }
      localPositions.push(groupStart - start);
    }
  }
  return withHash(new Tree(outerType, localChildren, localPositions, length), contextHash);
}
function containsType(nodes, type2) {
  for (let elt of nodes)
    if (elt.type == type2)
      return true;
  return false;
}
var TreeFragment = class {
  constructor(from, to, tree, offset, open) {
    this.from = from;
    this.to = to;
    this.tree = tree;
    this.offset = offset;
    this.open = open;
  }
  get openStart() {
    return (this.open & 1) > 0;
  }
  get openEnd() {
    return (this.open & 2) > 0;
  }
  static applyChanges(fragments, changes, minGap = 128) {
    if (!changes.length)
      return fragments;
    let result = [];
    let fI = 1, nextF = fragments.length ? fragments[0] : null;
    let cI = 0, pos = 0, off = 0;
    for (; ; ) {
      let nextC = cI < changes.length ? changes[cI++] : null;
      let nextPos = nextC ? nextC.fromA : 1e9;
      if (nextPos - pos >= minGap)
        while (nextF && nextF.from < nextPos) {
          let cut = nextF;
          if (pos >= cut.from || nextPos <= cut.to || off) {
            let fFrom = Math.max(cut.from, pos) - off, fTo = Math.min(cut.to, nextPos) - off;
            cut = fFrom >= fTo ? null : new TreeFragment(fFrom, fTo, cut.tree, cut.offset + off, (cI > 0 ? 1 : 0) | (nextC ? 2 : 0));
          }
          if (cut)
            result.push(cut);
          if (nextF.to > nextPos)
            break;
          nextF = fI < fragments.length ? fragments[fI++] : null;
        }
      if (!nextC)
        break;
      pos = nextC.toA;
      off = nextC.toA - nextC.toB;
    }
    return result;
  }
  static addTree(tree, fragments = [], partial = false) {
    let result = [new TreeFragment(0, tree.length, tree, 0, partial ? 2 : 0)];
    for (let f of fragments)
      if (f.to > tree.length)
        result.push(f);
    return result;
  }
};
function stringInput(input) {
  return new StringInput(input);
}
var StringInput = class {
  constructor(string3, length = string3.length) {
    this.string = string3;
    this.length = length;
  }
  get(pos) {
    return pos < 0 || pos >= this.length ? -1 : this.string.charCodeAt(pos);
  }
  lineAfter(pos) {
    if (pos < 0)
      return "";
    let end = this.string.indexOf("\n", pos);
    return this.string.slice(pos, end < 0 ? this.length : Math.min(end, this.length));
  }
  read(from, to) {
    return this.string.slice(from, Math.min(this.length, to));
  }
  clip(at) {
    return new StringInput(this.string, at);
  }
};

// node_modules/@codemirror/language/dist/index.js
var languageDataProp = new NodeProp();
function defineLanguageFacet(baseData) {
  return Facet.define({
    combine: baseData ? (values) => values.concat(baseData) : void 0
  });
}
var Language = class {
  constructor(data, parser6, topNode, extraExtensions = []) {
    this.data = data;
    this.topNode = topNode;
    if (!EditorState.prototype.hasOwnProperty("tree"))
      Object.defineProperty(EditorState.prototype, "tree", {get() {
        return syntaxTree(this);
      }});
    this.parser = parser6;
    this.extension = [
      language.of(this),
      EditorState.languageData.of((state, pos) => state.facet(languageDataFacetAt(state, pos)))
    ].concat(extraExtensions);
  }
  isActiveAt(state, pos) {
    return languageDataFacetAt(state, pos) == this.data;
  }
  findRegions(state) {
    let lang = state.facet(language);
    if ((lang === null || lang === void 0 ? void 0 : lang.data) == this.data)
      return [{from: 0, to: state.doc.length}];
    if (!lang || !lang.allowsNesting)
      return [];
    let result = [];
    syntaxTree(state).iterate({
      enter: (type2, from, to) => {
        if (type2.isTop && type2.prop(languageDataProp) == this.data) {
          result.push({from, to});
          return false;
        }
        return void 0;
      }
    });
    return result;
  }
  get allowsNesting() {
    return true;
  }
  parseString(code) {
    let doc2 = Text.of(code.split("\n"));
    let parse = this.parser.startParse(new DocInput(doc2), 0, new EditorParseContext(this.parser, EditorState.create({doc: doc2}), [], Tree.empty, {from: 0, to: code.length}, [], null));
    let tree;
    while (!(tree = parse.advance())) {
    }
    return tree;
  }
};
Language.setState = StateEffect.define();
function languageDataFacetAt(state, pos) {
  let topLang = state.facet(language);
  if (!topLang)
    return null;
  if (!topLang.allowsNesting)
    return topLang.data;
  let tree = syntaxTree(state);
  let target = tree.resolve(pos, -1);
  while (target) {
    let facet = target.type.prop(languageDataProp);
    if (facet)
      return facet;
    target = target.parent;
  }
  return topLang.data;
}
var LezerLanguage = class extends Language {
  constructor(data, parser6) {
    super(data, parser6, parser6.topNode);
    this.parser = parser6;
  }
  static define(spec) {
    let data = defineLanguageFacet(spec.languageData);
    return new LezerLanguage(data, spec.parser.configure({
      props: [languageDataProp.add((type2) => type2.isTop ? data : void 0)]
    }));
  }
  configure(options) {
    return new LezerLanguage(this.data, this.parser.configure(options));
  }
  get allowsNesting() {
    return this.parser.hasNested;
  }
};
function syntaxTree(state) {
  let field = state.field(Language.state, false);
  return field ? field.tree : Tree.empty;
}
var DocInput = class {
  constructor(doc2, length = doc2.length) {
    this.doc = doc2;
    this.length = length;
    this.cursorPos = 0;
    this.string = "";
    this.prevString = "";
    this.cursor = doc2.iter();
  }
  syncTo(pos) {
    if (pos < this.cursorPos) {
      this.cursor = this.doc.iter();
      this.cursorPos = 0;
    }
    this.prevString = pos == this.cursorPos ? this.string : "";
    this.string = this.cursor.next(pos - this.cursorPos).value;
    this.cursorPos = pos + this.string.length;
    return this.cursorPos - this.string.length;
  }
  get(pos) {
    if (pos >= this.length)
      return -1;
    let stringStart = this.cursorPos - this.string.length;
    if (pos < stringStart || pos >= this.cursorPos) {
      if (pos < stringStart && pos >= stringStart - this.prevString.length)
        return this.prevString.charCodeAt(pos - (stringStart - this.prevString.length));
      stringStart = this.syncTo(pos);
    }
    return this.string.charCodeAt(pos - stringStart);
  }
  lineAfter(pos) {
    if (pos >= this.length || pos < 0)
      return "";
    let stringStart = this.cursorPos - this.string.length;
    if (pos < stringStart || pos >= this.cursorPos)
      stringStart = this.syncTo(pos);
    return this.cursor.lineBreak ? "" : this.string.slice(pos - stringStart, Math.min(this.length - stringStart, this.string.length));
  }
  read(from, to) {
    let stringStart = this.cursorPos - this.string.length;
    if (from < stringStart || to >= this.cursorPos)
      return this.doc.sliceString(from, to);
    else
      return this.string.slice(from - stringStart, to - stringStart);
  }
  clip(at) {
    return new DocInput(this.doc, at);
  }
};
var EditorParseContext = class {
  constructor(parser6, state, fragments = [], tree, viewport, skipped, scheduleOn) {
    this.parser = parser6;
    this.state = state;
    this.fragments = fragments;
    this.tree = tree;
    this.viewport = viewport;
    this.skipped = skipped;
    this.scheduleOn = scheduleOn;
    this.parse = null;
    this.tempSkipped = [];
  }
  work(time, upto) {
    if (this.tree != Tree.empty && (upto == null ? this.tree.length == this.state.doc.length : this.tree.length >= upto)) {
      this.takeTree();
      return true;
    }
    if (!this.parse)
      this.parse = this.parser.startParse(new DocInput(this.state.doc), 0, this);
    let endTime = Date.now() + time;
    for (; ; ) {
      let done = this.parse.advance();
      if (done) {
        this.fragments = this.withoutTempSkipped(TreeFragment.addTree(done));
        this.parse = null;
        this.tree = done;
        return true;
      } else if (upto != null && this.parse.pos >= upto) {
        this.takeTree();
        return true;
      }
      if (Date.now() > endTime)
        return false;
    }
  }
  takeTree() {
    if (this.parse && this.parse.pos > this.tree.length) {
      this.tree = this.parse.forceFinish();
      this.fragments = this.withoutTempSkipped(TreeFragment.addTree(this.tree, this.fragments, true));
    }
  }
  withoutTempSkipped(fragments) {
    for (let r; r = this.tempSkipped.pop(); )
      fragments = cutFragments(fragments, r.from, r.to);
    return fragments;
  }
  changes(changes, newState) {
    let {fragments, tree, viewport, skipped} = this;
    this.takeTree();
    if (!changes.empty) {
      let ranges = [];
      changes.iterChangedRanges((fromA, toA, fromB, toB) => ranges.push({fromA, toA, fromB, toB}));
      fragments = TreeFragment.applyChanges(fragments, ranges);
      tree = Tree.empty;
      viewport = {from: changes.mapPos(viewport.from, -1), to: changes.mapPos(viewport.to, 1)};
      if (this.skipped.length) {
        skipped = [];
        for (let r of this.skipped) {
          let from = changes.mapPos(r.from, 1), to = changes.mapPos(r.to, -1);
          if (from < to)
            skipped.push({from, to});
        }
      }
    }
    return new EditorParseContext(this.parser, newState, fragments, tree, viewport, skipped, this.scheduleOn);
  }
  updateViewport(viewport) {
    this.viewport = viewport;
    let startLen = this.skipped.length;
    for (let i = 0; i < this.skipped.length; i++) {
      let {from, to} = this.skipped[i];
      if (from < viewport.to && to > viewport.from) {
        this.fragments = cutFragments(this.fragments, from, to);
        this.skipped.splice(i--, 1);
      }
    }
    return this.skipped.length < startLen;
  }
  reset() {
    if (this.parse) {
      this.takeTree();
      this.parse = null;
    }
  }
  skipUntilInView(from, to) {
    this.skipped.push({from, to});
  }
  static getSkippingParser(until) {
    return {
      startParse(input, startPos, context) {
        return {
          pos: startPos,
          advance() {
            let ecx = context;
            ecx.tempSkipped.push({from: startPos, to: input.length});
            if (until)
              ecx.scheduleOn = ecx.scheduleOn ? Promise.all([ecx.scheduleOn, until]) : until;
            this.pos = input.length;
            return new Tree(NodeType.none, [], [], input.length - startPos);
          },
          forceFinish() {
            return this.advance();
          }
        };
      }
    };
  }
  movedPast(pos) {
    return this.tree.length < pos && this.parse && this.parse.pos >= pos;
  }
};
EditorParseContext.skippingParser = EditorParseContext.getSkippingParser();
function cutFragments(fragments, from, to) {
  return TreeFragment.applyChanges(fragments, [{fromA: from, toA: to, fromB: from, toB: to}]);
}
var LanguageState = class {
  constructor(context) {
    this.context = context;
    this.tree = context.tree;
  }
  apply(tr) {
    if (!tr.docChanged)
      return this;
    let newCx = this.context.changes(tr.changes, tr.state);
    let upto = this.context.tree.length == tr.startState.doc.length ? void 0 : Math.max(tr.changes.mapPos(this.context.tree.length), newCx.viewport.to);
    if (!newCx.work(25, upto))
      newCx.takeTree();
    return new LanguageState(newCx);
  }
  static init(state) {
    let parseState = new EditorParseContext(state.facet(language).parser, state, [], Tree.empty, {from: 0, to: state.doc.length}, [], null);
    if (!parseState.work(25))
      parseState.takeTree();
    return new LanguageState(parseState);
  }
};
Language.state = StateField.define({
  create: LanguageState.init,
  update(value, tr) {
    for (let e of tr.effects)
      if (e.is(Language.setState))
        return e.value;
    if (tr.startState.facet(language) != tr.state.facet(language))
      return LanguageState.init(tr.state);
    return value.apply(tr);
  }
});
var requestIdle = typeof window != "undefined" && window.requestIdleCallback || ((callback, {timeout}) => setTimeout(callback, timeout));
var cancelIdle = typeof window != "undefined" && window.cancelIdleCallback || clearTimeout;
var parseWorker = ViewPlugin.fromClass(class ParseWorker {
  constructor(view) {
    this.view = view;
    this.working = -1;
    this.chunkEnd = -1;
    this.chunkBudget = -1;
    this.work = this.work.bind(this);
    this.scheduleWork();
  }
  update(update) {
    let cx = this.view.state.field(Language.state).context;
    if (update.viewportChanged) {
      if (cx.updateViewport(update.view.viewport))
        cx.reset();
      if (this.view.viewport.to > cx.tree.length)
        this.scheduleWork();
    }
    if (update.docChanged) {
      if (this.view.hasFocus)
        this.chunkBudget += 50;
      this.scheduleWork();
    }
    this.checkAsyncSchedule(cx);
  }
  scheduleWork(force = false) {
    if (this.working > -1)
      return;
    let {state} = this.view, field = state.field(Language.state);
    if (!force && field.tree.length >= state.doc.length)
      return;
    this.working = requestIdle(this.work, {timeout: 500});
  }
  work(deadline) {
    this.working = -1;
    let now = Date.now();
    if (this.chunkEnd < now && (this.chunkEnd < 0 || this.view.hasFocus)) {
      this.chunkEnd = now + 3e4;
      this.chunkBudget = 3e3;
    }
    if (this.chunkBudget <= 0)
      return;
    let {state, viewport: {to: vpTo}} = this.view, field = state.field(Language.state);
    if (field.tree.length >= vpTo + 1e6)
      return;
    let time = Math.min(this.chunkBudget, deadline ? Math.max(25, deadline.timeRemaining()) : 100);
    let done = field.context.work(time, vpTo + 1e6);
    this.chunkBudget -= Date.now() - now;
    if (done || this.chunkBudget <= 0 || field.context.movedPast(vpTo)) {
      field.context.takeTree();
      this.view.dispatch({effects: Language.setState.of(new LanguageState(field.context))});
    }
    if (!done && this.chunkBudget > 0)
      this.scheduleWork();
    this.checkAsyncSchedule(field.context);
  }
  checkAsyncSchedule(cx) {
    if (cx.scheduleOn) {
      cx.scheduleOn.then(() => this.scheduleWork(true));
      cx.scheduleOn = null;
    }
  }
  destroy() {
    if (this.working >= 0)
      cancelIdle(this.working);
  }
}, {
  eventHandlers: {focus() {
    this.scheduleWork();
  }}
});
var language = Facet.define({
  combine(languages2) {
    return languages2.length ? languages2[0] : null;
  },
  enables: [Language.state, parseWorker]
});
var LanguageSupport = class {
  constructor(language2, support = []) {
    this.language = language2;
    this.support = support;
    this.extension = [language2, support];
  }
};
var indentService = Facet.define();
var indentUnit = Facet.define({
  combine: (values) => {
    if (!values.length)
      return "  ";
    if (!/^(?: +|\t+)$/.test(values[0]))
      throw new Error("Invalid indent unit: " + JSON.stringify(values[0]));
    return values[0];
  }
});
function getIndentUnit(state) {
  let unit = state.facet(indentUnit);
  return unit.charCodeAt(0) == 9 ? state.tabSize * unit.length : unit.length;
}
function indentString(state, cols) {
  let result = "", ts = state.tabSize;
  if (state.facet(indentUnit).charCodeAt(0) == 9)
    while (cols >= ts) {
      result += "	";
      cols -= ts;
    }
  for (let i = 0; i < cols; i++)
    result += " ";
  return result;
}
function getIndentation(context, pos) {
  if (context instanceof EditorState)
    context = new IndentContext(context);
  for (let service of context.state.facet(indentService)) {
    let result = service(context, pos);
    if (result != null)
      return result;
  }
  let tree = syntaxTree(context.state);
  return tree ? syntaxIndentation(context, tree, pos) : null;
}
var IndentContext = class {
  constructor(state, options = {}) {
    this.state = state;
    this.options = options;
    this.unit = getIndentUnit(state);
  }
  textAfterPos(pos) {
    var _a, _b2;
    let sim = (_a = this.options) === null || _a === void 0 ? void 0 : _a.simulateBreak;
    if (pos == sim && ((_b2 = this.options) === null || _b2 === void 0 ? void 0 : _b2.simulateDoubleBreak))
      return "";
    return this.state.sliceDoc(pos, Math.min(pos + 100, sim != null && sim > pos ? sim : 1e9, this.state.doc.lineAt(pos).to));
  }
  column(pos) {
    var _a;
    let line = this.state.doc.lineAt(pos), text = line.text.slice(0, pos - line.from);
    let result = this.countColumn(text, pos - line.from);
    let override = ((_a = this.options) === null || _a === void 0 ? void 0 : _a.overrideIndentation) ? this.options.overrideIndentation(line.from) : -1;
    if (override > -1)
      result += override - this.countColumn(text, text.search(/\S/));
    return result;
  }
  countColumn(line, pos) {
    return countColumn(pos < 0 ? line : line.slice(0, pos), 0, this.state.tabSize);
  }
  lineIndent(line) {
    var _a;
    let override = (_a = this.options) === null || _a === void 0 ? void 0 : _a.overrideIndentation;
    if (override) {
      let overriden = override(line.from);
      if (overriden > -1)
        return overriden;
    }
    return this.countColumn(line.text, line.text.search(/\S/));
  }
};
var indentNodeProp = new NodeProp();
function syntaxIndentation(cx, ast, pos) {
  let tree = ast.resolve(pos);
  for (let scan = tree, scanPos = pos; ; ) {
    let last = scan.childBefore(scanPos);
    if (!last)
      break;
    if (last.type.isError && last.from == last.to) {
      tree = scan;
      scanPos = last.from;
    } else {
      scan = last;
      scanPos = scan.to + 1;
    }
  }
  return indentFrom(tree, pos, cx);
}
function ignoreClosed(cx) {
  var _a, _b2;
  return cx.pos == ((_a = cx.options) === null || _a === void 0 ? void 0 : _a.simulateBreak) && ((_b2 = cx.options) === null || _b2 === void 0 ? void 0 : _b2.simulateDoubleBreak);
}
function indentStrategy(tree) {
  let strategy = tree.type.prop(indentNodeProp);
  if (strategy)
    return strategy;
  let first = tree.firstChild, close;
  if (first && (close = first.type.prop(NodeProp.closedBy))) {
    let last = tree.lastChild, closed = last && close.indexOf(last.name) > -1;
    return (cx) => delimitedStrategy(cx, true, 1, void 0, closed && !ignoreClosed(cx) ? last.from : void 0);
  }
  return tree.parent == null ? topIndent : null;
}
function indentFrom(node, pos, base3) {
  for (; node; node = node.parent) {
    let strategy = indentStrategy(node);
    if (strategy)
      return strategy(new TreeIndentContext(base3, pos, node));
  }
  return null;
}
function topIndent() {
  return 0;
}
var TreeIndentContext = class extends IndentContext {
  constructor(base3, pos, node) {
    super(base3.state, base3.options);
    this.base = base3;
    this.pos = pos;
    this.node = node;
  }
  get textAfter() {
    return this.textAfterPos(this.pos);
  }
  get baseIndent() {
    let line = this.state.doc.lineAt(this.node.from);
    for (; ; ) {
      let atBreak = this.node.resolve(line.from);
      while (atBreak.parent && atBreak.parent.from == atBreak.from)
        atBreak = atBreak.parent;
      if (isParent(atBreak, this.node))
        break;
      line = this.state.doc.lineAt(atBreak.from);
    }
    return this.lineIndent(line);
  }
  continue() {
    let parent = this.node.parent;
    return parent ? indentFrom(parent, this.pos, this.base) : 0;
  }
};
function isParent(parent, of) {
  for (let cur2 = of; cur2; cur2 = cur2.parent)
    if (parent == cur2)
      return true;
  return false;
}
function bracketedAligned(context) {
  var _a;
  let tree = context.node;
  let openToken = tree.childAfter(tree.from), last = tree.lastChild;
  if (!openToken)
    return null;
  let sim = (_a = context.options) === null || _a === void 0 ? void 0 : _a.simulateBreak;
  let openLine = context.state.doc.lineAt(openToken.from);
  let lineEnd = sim == null || sim <= openLine.from ? openLine.to : Math.min(openLine.to, sim);
  for (let pos = openToken.to; ; ) {
    let next = tree.childAfter(pos);
    if (!next || next == last)
      return null;
    if (!next.type.isSkipped)
      return next.from < lineEnd ? openToken : null;
    pos = next.to;
  }
}
function delimitedIndent({closing: closing3, align = true, units = 1}) {
  return (context) => delimitedStrategy(context, align, units, closing3);
}
function delimitedStrategy(context, align, units, closing3, closedAt) {
  let after = context.textAfter, space3 = after.match(/^\s*/)[0].length;
  let closed = closing3 && after.slice(space3, space3 + closing3.length) == closing3 || closedAt == context.pos + space3;
  let aligned = align ? bracketedAligned(context) : null;
  if (aligned)
    return closed ? context.column(aligned.from) : context.column(aligned.to);
  return context.baseIndent + (closed ? 0 : context.unit * units);
}
var flatIndent = (context) => context.baseIndent;
function continuedIndent({except, units = 1} = {}) {
  return (context) => {
    let matchExcept = except && except.test(context.textAfter);
    return context.baseIndent + (matchExcept ? 0 : units * context.unit);
  };
}
var DontIndentBeyond = 200;
function indentOnInput() {
  return EditorState.transactionFilter.of((tr) => {
    if (!tr.docChanged || tr.annotation(Transaction.userEvent) != "input")
      return tr;
    let rules = tr.startState.languageDataAt("indentOnInput", tr.startState.selection.main.head);
    if (!rules.length)
      return tr;
    let doc2 = tr.newDoc, {head} = tr.newSelection.main, line = doc2.lineAt(head);
    if (head > line.from + DontIndentBeyond)
      return tr;
    let lineStart = doc2.sliceString(line.from, head);
    if (!rules.some((r) => r.test(lineStart)))
      return tr;
    let {state} = tr, last = -1, changes = [];
    for (let {head: head2} of state.selection.ranges) {
      let line2 = state.doc.lineAt(head2);
      if (line2.from == last)
        continue;
      last = line2.from;
      let indent2 = getIndentation(state, line2.from);
      if (indent2 == null)
        continue;
      let cur2 = /^\s*/.exec(line2.text)[0];
      let norm = indentString(state, indent2);
      if (cur2 != norm)
        changes.push({from: line2.from, to: line2.from + cur2.length, insert: norm});
    }
    return changes.length ? [tr, {changes}] : tr;
  });
}
var foldService = Facet.define();
var foldNodeProp = new NodeProp();
function foldInside(node) {
  let first = node.firstChild, last = node.lastChild;
  return first && first.to < last.from ? {from: first.to, to: last.type.isError ? node.to : last.from} : null;
}

// node_modules/@codemirror/closebrackets/dist/index.js
var defaults = {
  brackets: ["(", "[", "{", "'", '"'],
  before: `)]}'":;>`
};
var closeBracketEffect = StateEffect.define({
  map(value, mapping) {
    let mapped = mapping.mapPos(value, -1, MapMode.TrackAfter);
    return mapped == null ? void 0 : mapped;
  }
});
var skipBracketEffect = StateEffect.define({
  map(value, mapping) {
    return mapping.mapPos(value);
  }
});
var closedBracket = new class extends RangeValue {
}();
closedBracket.startSide = 1;
closedBracket.endSide = -1;
var bracketState = StateField.define({
  create() {
    return RangeSet.empty;
  },
  update(value, tr) {
    if (tr.selection) {
      let lineStart = tr.state.doc.lineAt(tr.selection.main.head).from;
      let prevLineStart = tr.startState.doc.lineAt(tr.startState.selection.main.head).from;
      if (lineStart != tr.changes.mapPos(prevLineStart, -1))
        value = RangeSet.empty;
    }
    value = value.map(tr.changes);
    for (let effect of tr.effects) {
      if (effect.is(closeBracketEffect))
        value = value.update({add: [closedBracket.range(effect.value, effect.value + 1)]});
      else if (effect.is(skipBracketEffect))
        value = value.update({filter: (from) => from != effect.value});
    }
    return value;
  }
});
function closeBrackets() {
  return [EditorView.inputHandler.of(handleInput), bracketState];
}
var definedClosing = "()[]{}<>";
function closing(ch) {
  for (let i = 0; i < definedClosing.length; i += 2)
    if (definedClosing.charCodeAt(i) == ch)
      return definedClosing.charAt(i + 1);
  return fromCodePoint(ch < 128 ? ch : ch + 1);
}
function config(state, pos) {
  return state.languageDataAt("closeBrackets", pos)[0] || defaults;
}
function handleInput(view, from, to, insert2) {
  if (view.composing)
    return false;
  let sel = view.state.selection.main;
  if (insert2.length > 2 || insert2.length == 2 && codePointSize(codePointAt(insert2, 0)) == 1 || from != sel.from || to != sel.to)
    return false;
  let tr = insertBracket(view.state, insert2);
  if (!tr)
    return false;
  view.dispatch(tr);
  return true;
}
var deleteBracketPair = ({state, dispatch}) => {
  let conf = config(state, state.selection.main.head);
  let tokens2 = conf.brackets || defaults.brackets;
  let dont = null, changes = state.changeByRange((range) => {
    if (range.empty) {
      let before = prevChar(state.doc, range.head);
      for (let token of tokens2) {
        if (token == before && nextChar(state.doc, range.head) == closing(codePointAt(token, 0)))
          return {
            changes: {from: range.head - token.length, to: range.head + token.length},
            range: EditorSelection.cursor(range.head - token.length),
            annotations: Transaction.userEvent.of("delete")
          };
      }
    }
    return {range: dont = range};
  });
  if (!dont)
    dispatch(state.update(changes, {scrollIntoView: true}));
  return !dont;
};
var closeBracketsKeymap = [
  {key: "Backspace", run: deleteBracketPair}
];
function insertBracket(state, bracket2) {
  let conf = config(state, state.selection.main.head);
  let tokens2 = conf.brackets || defaults.brackets;
  for (let tok of tokens2) {
    let closed = closing(codePointAt(tok, 0));
    if (bracket2 == tok)
      return closed == tok ? handleSame(state, tok, tokens2.indexOf(tok + tok + tok) > -1) : handleOpen(state, tok, closed, conf.before || defaults.before);
    if (bracket2 == closed && closedBracketAt(state, state.selection.main.from))
      return handleClose(state, tok, closed);
  }
  return null;
}
function closedBracketAt(state, pos) {
  let found = false;
  state.field(bracketState).between(0, state.doc.length, (from) => {
    if (from == pos)
      found = true;
  });
  return found;
}
function nextChar(doc2, pos) {
  let next = doc2.sliceString(pos, pos + 2);
  return next.slice(0, codePointSize(codePointAt(next, 0)));
}
function prevChar(doc2, pos) {
  let prev = doc2.sliceString(pos - 2, pos);
  return codePointSize(codePointAt(prev, 0)) == prev.length ? prev : prev.slice(1);
}
function handleOpen(state, open, close, closeBefore) {
  let dont = null, changes = state.changeByRange((range) => {
    if (!range.empty)
      return {
        changes: [{insert: open, from: range.from}, {insert: close, from: range.to}],
        effects: closeBracketEffect.of(range.to + open.length),
        range: EditorSelection.range(range.anchor + open.length, range.head + open.length)
      };
    let next = nextChar(state.doc, range.head);
    if (!next || /\s/.test(next) || closeBefore.indexOf(next) > -1)
      return {
        changes: {insert: open + close, from: range.head},
        effects: closeBracketEffect.of(range.head + open.length),
        range: EditorSelection.cursor(range.head + open.length)
      };
    return {range: dont = range};
  });
  return dont ? null : state.update(changes, {
    scrollIntoView: true,
    annotations: Transaction.userEvent.of("input")
  });
}
function handleClose(state, _open, close) {
  let dont = null, moved = state.selection.ranges.map((range) => {
    if (range.empty && nextChar(state.doc, range.head) == close)
      return EditorSelection.cursor(range.head + close.length);
    return dont = range;
  });
  return dont ? null : state.update({
    selection: EditorSelection.create(moved, state.selection.mainIndex),
    scrollIntoView: true,
    effects: state.selection.ranges.map(({from}) => skipBracketEffect.of(from))
  });
}
function handleSame(state, token, allowTriple) {
  let dont = null, changes = state.changeByRange((range) => {
    if (!range.empty)
      return {
        changes: [{insert: token, from: range.from}, {insert: token, from: range.to}],
        effects: closeBracketEffect.of(range.to + token.length),
        range: EditorSelection.range(range.anchor + token.length, range.head + token.length)
      };
    let pos = range.head, next = nextChar(state.doc, pos);
    if (next == token) {
      if (nodeStart(state, pos)) {
        return {
          changes: {insert: token + token, from: pos},
          effects: closeBracketEffect.of(pos + token.length),
          range: EditorSelection.cursor(pos + token.length)
        };
      } else if (closedBracketAt(state, pos)) {
        let isTriple = allowTriple && state.sliceDoc(pos, pos + token.length * 3) == token + token + token;
        return {
          range: EditorSelection.cursor(pos + token.length * (isTriple ? 3 : 1)),
          effects: skipBracketEffect.of(pos)
        };
      }
    } else if (allowTriple && state.sliceDoc(pos - 2 * token.length, pos) == token + token && nodeStart(state, pos - 2 * token.length)) {
      return {
        changes: {insert: token + token + token + token, from: pos},
        effects: closeBracketEffect.of(pos + token.length),
        range: EditorSelection.cursor(pos + token.length)
      };
    } else if (state.charCategorizer(pos)(next) != CharCategory.Word) {
      let prev = state.sliceDoc(pos - 1, pos);
      if (prev != token && state.charCategorizer(pos)(prev) != CharCategory.Word)
        return {
          changes: {insert: token + token, from: pos},
          effects: closeBracketEffect.of(pos + token.length),
          range: EditorSelection.cursor(pos + token.length)
        };
    }
    return {range: dont = range};
  });
  return dont ? null : state.update(changes, {
    scrollIntoView: true,
    annotations: Transaction.userEvent.of("input")
  });
}
function nodeStart(state, pos) {
  let tree = syntaxTree(state).resolve(pos + 1);
  return tree.parent && tree.from == pos;
}

// node_modules/@codemirror/matchbrackets/dist/index.js
var baseTheme2 = EditorView.baseTheme({
  ".cm-matchingBracket": {color: "#0b0"},
  ".cm-nonmatchingBracket": {color: "#a22"}
});
var DefaultScanDist = 1e4;
var DefaultBrackets = "()[]{}";
var bracketMatchingConfig = Facet.define({
  combine(configs) {
    return combineConfig(configs, {
      afterCursor: true,
      brackets: DefaultBrackets,
      maxScanDistance: DefaultScanDist
    });
  }
});
var matchingMark = Decoration.mark({class: "cm-matchingBracket"});
var nonmatchingMark = Decoration.mark({class: "cm-nonmatchingBracket"});
var bracketMatchingState = StateField.define({
  create() {
    return Decoration.none;
  },
  update(deco, tr) {
    if (!tr.docChanged && !tr.selection)
      return deco;
    let decorations2 = [];
    let config2 = tr.state.facet(bracketMatchingConfig);
    for (let range of tr.state.selection.ranges) {
      if (!range.empty)
        continue;
      let match = matchBrackets(tr.state, range.head, -1, config2) || range.head > 0 && matchBrackets(tr.state, range.head - 1, 1, config2) || config2.afterCursor && (matchBrackets(tr.state, range.head, 1, config2) || range.head < tr.state.doc.length && matchBrackets(tr.state, range.head + 1, -1, config2));
      if (!match)
        continue;
      let mark = match.matched ? matchingMark : nonmatchingMark;
      decorations2.push(mark.range(match.start.from, match.start.to));
      if (match.end)
        decorations2.push(mark.range(match.end.from, match.end.to));
    }
    return Decoration.set(decorations2, true);
  },
  provide: (f) => EditorView.decorations.from(f)
});
var bracketMatchingUnique = [
  bracketMatchingState,
  baseTheme2
];
function bracketMatching(config2 = {}) {
  return [bracketMatchingConfig.of(config2), bracketMatchingUnique];
}
function matchingNodes(node, dir, brackets) {
  let byProp = node.prop(dir < 0 ? NodeProp.openedBy : NodeProp.closedBy);
  if (byProp)
    return byProp;
  if (node.name.length == 1) {
    let index = brackets.indexOf(node.name);
    if (index > -1 && index % 2 == (dir < 0 ? 1 : 0))
      return [brackets[index + dir]];
  }
  return null;
}
function matchBrackets(state, pos, dir, config2 = {}) {
  let maxScanDistance = config2.maxScanDistance || DefaultScanDist, brackets = config2.brackets || DefaultBrackets;
  let tree = syntaxTree(state), sub = tree.resolve(pos, dir), matches;
  if (matches = matchingNodes(sub.type, dir, brackets))
    return matchMarkedBrackets(state, pos, dir, sub, matches, brackets);
  else
    return matchPlainBrackets(state, pos, dir, tree, sub.type, maxScanDistance, brackets);
}
function matchMarkedBrackets(_state, _pos, dir, token, matching2, brackets) {
  let parent = token.parent, firstToken = {from: token.from, to: token.to};
  let depth = 0, cursor = parent === null || parent === void 0 ? void 0 : parent.cursor;
  if (cursor && (dir < 0 ? cursor.childBefore(token.from) : cursor.childAfter(token.to)))
    do {
      if (dir < 0 ? cursor.to <= token.from : cursor.from >= token.to) {
        if (depth == 0 && matching2.indexOf(cursor.type.name) > -1) {
          return {start: firstToken, end: {from: cursor.from, to: cursor.to}, matched: true};
        } else if (matchingNodes(cursor.type, dir, brackets)) {
          depth++;
        } else if (matchingNodes(cursor.type, -dir, brackets)) {
          depth--;
          if (depth == 0)
            return {start: firstToken, end: {from: cursor.from, to: cursor.to}, matched: false};
        }
      }
    } while (dir < 0 ? cursor.prevSibling() : cursor.nextSibling());
  return {start: firstToken, matched: false};
}
function matchPlainBrackets(state, pos, dir, tree, tokenType, maxScanDistance, brackets) {
  let startCh = dir < 0 ? state.sliceDoc(pos - 1, pos) : state.sliceDoc(pos, pos + 1);
  let bracket2 = brackets.indexOf(startCh);
  if (bracket2 < 0 || bracket2 % 2 == 0 != dir > 0)
    return null;
  let startToken = {from: dir < 0 ? pos - 1 : pos, to: dir > 0 ? pos + 1 : pos};
  let iter = state.doc.iterRange(pos, dir > 0 ? state.doc.length : 0), depth = 0;
  for (let distance = 0; !iter.next().done && distance <= maxScanDistance; ) {
    let text = iter.value;
    if (dir < 0)
      distance += text.length;
    let basePos = pos + distance * dir;
    for (let pos2 = dir > 0 ? 0 : text.length - 1, end = dir > 0 ? text.length : -1; pos2 != end; pos2 += dir) {
      let found = brackets.indexOf(text[pos2]);
      if (found < 0 || tree.resolve(basePos + pos2, 1).type != tokenType)
        continue;
      if (found % 2 == 0 == dir > 0) {
        depth++;
      } else if (depth == 1) {
        return {start: startToken, end: {from: basePos + pos2, to: basePos + pos2 + 1}, matched: found >> 1 == bracket2 >> 1};
      } else {
        depth--;
      }
    }
    if (dir > 0)
      distance += text.length;
  }
  return iter.done ? {start: startToken, matched: false} : null;
}

// node_modules/@codemirror/commands/dist/index.js
function updateSel(sel, by) {
  return EditorSelection.create(sel.ranges.map(by), sel.mainIndex);
}
function setSel(state, selection) {
  return state.update({selection, scrollIntoView: true, annotations: Transaction.userEvent.of("keyboardselection")});
}
function moveSel({state, dispatch}, how) {
  let selection = updateSel(state.selection, how);
  if (selection.eq(state.selection))
    return false;
  dispatch(setSel(state, selection));
  return true;
}
function rangeEnd(range, forward) {
  return EditorSelection.cursor(forward ? range.to : range.from);
}
function cursorByChar(view, forward) {
  return moveSel(view, (range) => range.empty ? view.moveByChar(range, forward) : rangeEnd(range, forward));
}
var cursorCharLeft = (view) => cursorByChar(view, view.textDirection != Direction.LTR);
var cursorCharRight = (view) => cursorByChar(view, view.textDirection == Direction.LTR);
function cursorByGroup(view, forward) {
  return moveSel(view, (range) => range.empty ? view.moveByGroup(range, forward) : rangeEnd(range, forward));
}
var cursorGroupLeft = (view) => cursorByGroup(view, view.textDirection != Direction.LTR);
var cursorGroupRight = (view) => cursorByGroup(view, view.textDirection == Direction.LTR);
var cursorGroupForward = (view) => cursorByGroup(view, true);
var cursorGroupBackward = (view) => cursorByGroup(view, false);
function cursorByLine(view, forward) {
  return moveSel(view, (range) => range.empty ? view.moveVertically(range, forward) : rangeEnd(range, forward));
}
var cursorLineUp = (view) => cursorByLine(view, false);
var cursorLineDown = (view) => cursorByLine(view, true);
function cursorByPage(view, forward) {
  return moveSel(view, (range) => range.empty ? view.moveVertically(range, forward, view.dom.clientHeight) : rangeEnd(range, forward));
}
var cursorPageUp = (view) => cursorByPage(view, false);
var cursorPageDown = (view) => cursorByPage(view, true);
function moveByLineBoundary(view, start, forward) {
  let line = view.visualLineAt(start.head), moved = view.moveToLineBoundary(start, forward);
  if (moved.head == start.head && moved.head != (forward ? line.to : line.from))
    moved = view.moveToLineBoundary(start, forward, false);
  if (!forward && moved.head == line.from && line.length) {
    let space3 = /^\s*/.exec(view.state.sliceDoc(line.from, Math.min(line.from + 100, line.to)))[0].length;
    if (space3 && start.head != line.from + space3)
      moved = EditorSelection.cursor(line.from + space3);
  }
  return moved;
}
var cursorLineBoundaryForward = (view) => moveSel(view, (range) => moveByLineBoundary(view, range, true));
var cursorLineBoundaryBackward = (view) => moveSel(view, (range) => moveByLineBoundary(view, range, false));
var cursorLineStart = (view) => moveSel(view, (range) => EditorSelection.cursor(view.visualLineAt(range.head).from, 1));
var cursorLineEnd = (view) => moveSel(view, (range) => EditorSelection.cursor(view.visualLineAt(range.head).to, -1));
function extendSel(view, how) {
  let selection = updateSel(view.state.selection, (range) => {
    let head = how(range);
    return EditorSelection.range(range.anchor, head.head, head.goalColumn);
  });
  if (selection.eq(view.state.selection))
    return false;
  view.dispatch(setSel(view.state, selection));
  return true;
}
function selectByChar(view, forward) {
  return extendSel(view, (range) => view.moveByChar(range, forward));
}
var selectCharLeft = (view) => selectByChar(view, view.textDirection != Direction.LTR);
var selectCharRight = (view) => selectByChar(view, view.textDirection == Direction.LTR);
function selectByGroup(view, forward) {
  return extendSel(view, (range) => view.moveByGroup(range, forward));
}
var selectGroupLeft = (view) => selectByGroup(view, view.textDirection != Direction.LTR);
var selectGroupRight = (view) => selectByGroup(view, view.textDirection == Direction.LTR);
var selectGroupForward = (view) => selectByGroup(view, true);
var selectGroupBackward = (view) => selectByGroup(view, false);
function selectByLine(view, forward) {
  return extendSel(view, (range) => view.moveVertically(range, forward));
}
var selectLineUp = (view) => selectByLine(view, false);
var selectLineDown = (view) => selectByLine(view, true);
function selectByPage(view, forward) {
  return extendSel(view, (range) => view.moveVertically(range, forward, view.dom.clientHeight));
}
var selectPageUp = (view) => selectByPage(view, false);
var selectPageDown = (view) => selectByPage(view, true);
var selectLineBoundaryForward = (view) => extendSel(view, (range) => moveByLineBoundary(view, range, true));
var selectLineBoundaryBackward = (view) => extendSel(view, (range) => moveByLineBoundary(view, range, false));
var selectLineStart = (view) => extendSel(view, (range) => EditorSelection.cursor(view.visualLineAt(range.head).from));
var selectLineEnd = (view) => extendSel(view, (range) => EditorSelection.cursor(view.visualLineAt(range.head).to));
var cursorDocStart = ({state, dispatch}) => {
  dispatch(setSel(state, {anchor: 0}));
  return true;
};
var cursorDocEnd = ({state, dispatch}) => {
  dispatch(setSel(state, {anchor: state.doc.length}));
  return true;
};
var selectDocStart = ({state, dispatch}) => {
  dispatch(setSel(state, {anchor: state.selection.main.anchor, head: 0}));
  return true;
};
var selectDocEnd = ({state, dispatch}) => {
  dispatch(setSel(state, {anchor: state.selection.main.anchor, head: state.doc.length}));
  return true;
};
var selectAll = ({state, dispatch}) => {
  dispatch(state.update({selection: {anchor: 0, head: state.doc.length}, annotations: Transaction.userEvent.of("keyboardselection")}));
  return true;
};
function deleteBy({state, dispatch}, by) {
  let changes = state.changeByRange((range) => {
    let {from, to} = range;
    if (from == to) {
      let towards = by(from);
      from = Math.min(from, towards);
      to = Math.max(to, towards);
    }
    return from == to ? {range} : {changes: {from, to}, range: EditorSelection.cursor(from)};
  });
  if (changes.changes.empty)
    return false;
  dispatch(state.update(changes, {scrollIntoView: true, annotations: Transaction.userEvent.of("delete")}));
  return true;
}
var deleteByChar = (target, forward, codePoint) => deleteBy(target, (pos) => {
  let {state} = target, line = state.doc.lineAt(pos), before;
  if (!forward && pos > line.from && pos < line.from + 200 && !/[^ \t]/.test(before = line.text.slice(0, pos - line.from))) {
    if (before[before.length - 1] == "	")
      return pos - 1;
    let col = countColumn(before, 0, state.tabSize), drop = col % getIndentUnit(state) || getIndentUnit(state);
    for (let i = 0; i < drop && before[before.length - 1 - i] == " "; i++)
      pos--;
    return pos;
  }
  let targetPos;
  if (codePoint) {
    let next = line.text.slice(pos - line.from + (forward ? 0 : -2), pos - line.from + (forward ? 2 : 0));
    let size = next ? codePointSize(codePointAt(next, 0)) : 1;
    targetPos = forward ? Math.min(state.doc.length, pos + size) : Math.max(0, pos - size);
  } else {
    targetPos = findClusterBreak(line.text, pos - line.from, forward) + line.from;
  }
  if (targetPos == pos && line.number != (forward ? state.doc.lines : 1))
    targetPos += forward ? 1 : -1;
  return targetPos;
});
var deleteCodePointBackward = (view) => deleteByChar(view, false, true);
var deleteCharBackward = (view) => deleteByChar(view, false, false);
var deleteCharForward = (view) => deleteByChar(view, true, false);
var deleteByGroup = (target, forward) => deleteBy(target, (start) => {
  let pos = start, {state} = target, line = state.doc.lineAt(pos);
  let categorize = state.charCategorizer(pos);
  for (let cat = null; ; ) {
    if (pos == (forward ? line.to : line.from)) {
      if (pos == start && line.number != (forward ? state.doc.lines : 1))
        pos += forward ? 1 : -1;
      break;
    }
    let next = findClusterBreak(line.text, pos - line.from, forward) + line.from;
    let nextChar2 = line.text.slice(Math.min(pos, next) - line.from, Math.max(pos, next) - line.from);
    let nextCat = categorize(nextChar2);
    if (cat != null && nextCat != cat)
      break;
    if (nextChar2 != " " || pos != start)
      cat = nextCat;
    pos = next;
  }
  return pos;
});
var deleteGroupBackward = (target) => deleteByGroup(target, false);
var deleteGroupForward = (target) => deleteByGroup(target, true);
var deleteToLineEnd = (view) => deleteBy(view, (pos) => {
  let lineEnd = view.visualLineAt(pos).to;
  if (pos < lineEnd)
    return lineEnd;
  return Math.min(view.state.doc.length, pos + 1);
});
var deleteToLineStart = (view) => deleteBy(view, (pos) => {
  let lineStart = view.visualLineAt(pos).from;
  if (pos > lineStart)
    return lineStart;
  return Math.max(0, pos - 1);
});
var splitLine = ({state, dispatch}) => {
  let changes = state.changeByRange((range) => {
    return {
      changes: {from: range.from, to: range.to, insert: Text.of(["", ""])},
      range: EditorSelection.cursor(range.from)
    };
  });
  dispatch(state.update(changes, {scrollIntoView: true, annotations: Transaction.userEvent.of("input")}));
  return true;
};
var transposeChars = ({state, dispatch}) => {
  let changes = state.changeByRange((range) => {
    if (!range.empty || range.from == 0 || range.from == state.doc.length)
      return {range};
    let pos = range.from, line = state.doc.lineAt(pos);
    let from = pos == line.from ? pos - 1 : findClusterBreak(line.text, pos - line.from, false) + line.from;
    let to = pos == line.to ? pos + 1 : findClusterBreak(line.text, pos - line.from, true) + line.from;
    return {
      changes: {from, to, insert: state.doc.slice(pos, to).append(state.doc.slice(from, pos))},
      range: EditorSelection.cursor(to)
    };
  });
  if (changes.changes.empty)
    return false;
  dispatch(state.update(changes, {scrollIntoView: true}));
  return true;
};
function isBetweenBrackets(state, pos) {
  if (/\(\)|\[\]|\{\}/.test(state.sliceDoc(pos - 1, pos + 1)))
    return {from: pos, to: pos};
  let context = syntaxTree(state).resolve(pos);
  let before = context.childBefore(pos), after = context.childAfter(pos), closedBy;
  if (before && after && before.to <= pos && after.from >= pos && (closedBy = before.type.prop(NodeProp.closedBy)) && closedBy.indexOf(after.name) > -1 && state.doc.lineAt(before.to).from == state.doc.lineAt(after.from).from)
    return {from: before.to, to: after.from};
  return null;
}
var insertNewlineAndIndent = ({state, dispatch}) => {
  let changes = state.changeByRange(({from, to}) => {
    let explode = from == to && isBetweenBrackets(state, from);
    let cx = new IndentContext(state, {simulateBreak: from, simulateDoubleBreak: !!explode});
    let indent2 = getIndentation(cx, from);
    if (indent2 == null)
      indent2 = /^\s*/.exec(state.doc.lineAt(from).text)[0].length;
    let line = state.doc.lineAt(from);
    while (to < line.to && /\s/.test(line.text.slice(to - line.from, to + 1 - line.from)))
      to++;
    if (explode)
      ({from, to} = explode);
    else if (from > line.from && from < line.from + 100 && !/\S/.test(line.text.slice(0, from)))
      from = line.from;
    let insert2 = ["", indentString(state, indent2)];
    if (explode)
      insert2.push(indentString(state, cx.lineIndent(line)));
    return {
      changes: {from, to, insert: Text.of(insert2)},
      range: EditorSelection.cursor(from + 1 + insert2[1].length)
    };
  });
  dispatch(state.update(changes, {scrollIntoView: true}));
  return true;
};
function changeBySelectedLine(state, f) {
  let atLine = -1;
  return state.changeByRange((range) => {
    let changes = [];
    for (let pos = range.from; pos <= range.to; ) {
      let line = state.doc.lineAt(pos);
      if (line.number > atLine && (range.empty || range.to > line.from)) {
        f(line, changes, range);
        atLine = line.number;
      }
      pos = line.to + 1;
    }
    let changeSet = state.changes(changes);
    return {
      changes,
      range: EditorSelection.range(changeSet.mapPos(range.anchor, 1), changeSet.mapPos(range.head, 1))
    };
  });
}
var indentSelection = ({state, dispatch}) => {
  let updated = Object.create(null);
  let context = new IndentContext(state, {overrideIndentation: (start) => {
    let found = updated[start];
    return found == null ? -1 : found;
  }});
  let changes = changeBySelectedLine(state, (line, changes2, range) => {
    let indent2 = getIndentation(context, line.from);
    if (indent2 == null)
      return;
    let cur2 = /^\s*/.exec(line.text)[0];
    let norm = indentString(state, indent2);
    if (cur2 != norm || range.from < line.from + cur2.length) {
      updated[line.from] = indent2;
      changes2.push({from: line.from, to: line.from + cur2.length, insert: norm});
    }
  });
  if (!changes.changes.empty)
    dispatch(state.update(changes));
  return true;
};
var indentMore = ({state, dispatch}) => {
  dispatch(state.update(changeBySelectedLine(state, (line, changes) => {
    changes.push({from: line.from, insert: state.facet(indentUnit)});
  })));
  return true;
};
var insertTab = ({state, dispatch}) => {
  if (state.selection.ranges.some((r) => !r.empty))
    return indentMore({state, dispatch});
  dispatch(state.update(state.replaceSelection("	"), {scrollIntoView: true, annotations: Transaction.userEvent.of("input")}));
  return true;
};
var emacsStyleKeymap = [
  {key: "Ctrl-b", run: cursorCharLeft, shift: selectCharLeft},
  {key: "Ctrl-f", run: cursorCharRight, shift: selectCharRight},
  {key: "Ctrl-p", run: cursorLineUp, shift: selectLineUp},
  {key: "Ctrl-n", run: cursorLineDown, shift: selectLineDown},
  {key: "Ctrl-a", run: cursorLineStart, shift: selectLineStart},
  {key: "Ctrl-e", run: cursorLineEnd, shift: selectLineEnd},
  {key: "Ctrl-d", run: deleteCharForward},
  {key: "Ctrl-h", run: deleteCharBackward},
  {key: "Ctrl-k", run: deleteToLineEnd},
  {key: "Alt-d", run: deleteGroupForward},
  {key: "Ctrl-Alt-h", run: deleteGroupBackward},
  {key: "Ctrl-o", run: splitLine},
  {key: "Ctrl-t", run: transposeChars},
  {key: "Alt-f", run: cursorGroupForward, shift: selectGroupForward},
  {key: "Alt-b", run: cursorGroupBackward, shift: selectGroupBackward},
  {key: "Alt-<", run: cursorDocStart},
  {key: "Alt->", run: cursorDocEnd},
  {key: "Ctrl-v", run: cursorPageDown},
  {key: "Alt-v", run: cursorPageUp}
];
var standardKeymap = /* @__PURE__ */ [
  {key: "ArrowLeft", run: cursorCharLeft, shift: selectCharLeft},
  {key: "Mod-ArrowLeft", mac: "Alt-ArrowLeft", run: cursorGroupLeft, shift: selectGroupLeft},
  {mac: "Cmd-ArrowLeft", run: cursorLineStart, shift: selectLineStart},
  {key: "ArrowRight", run: cursorCharRight, shift: selectCharRight},
  {key: "Mod-ArrowRight", mac: "Alt-ArrowRight", run: cursorGroupRight, shift: selectGroupRight},
  {mac: "Cmd-ArrowRight", run: cursorLineEnd, shift: selectLineEnd},
  {key: "ArrowUp", run: cursorLineUp, shift: selectLineUp},
  {mac: "Cmd-ArrowUp", run: cursorDocStart, shift: selectDocStart},
  {mac: "Ctrl-ArrowUp", run: cursorPageUp, shift: selectPageUp},
  {key: "ArrowDown", run: cursorLineDown, shift: selectLineDown},
  {mac: "Cmd-ArrowDown", run: cursorDocEnd, shift: selectDocEnd},
  {mac: "Ctrl-ArrowDown", run: cursorPageDown, shift: selectPageDown},
  {key: "PageUp", run: cursorPageUp, shift: selectPageUp},
  {key: "PageDown", run: cursorPageDown, shift: selectPageDown},
  {key: "Home", run: cursorLineBoundaryBackward, shift: selectLineBoundaryBackward},
  {key: "Mod-Home", run: cursorDocStart, shift: selectDocStart},
  {key: "End", run: cursorLineBoundaryForward, shift: selectLineBoundaryForward},
  {key: "Mod-End", run: cursorDocEnd, shift: selectDocEnd},
  {key: "Enter", run: insertNewlineAndIndent},
  {key: "Mod-a", run: selectAll},
  {key: "Backspace", run: deleteCodePointBackward, shift: deleteCodePointBackward},
  {key: "Delete", run: deleteCharForward, shift: deleteCharForward},
  {key: "Mod-Backspace", mac: "Alt-Backspace", run: deleteGroupBackward},
  {key: "Mod-Delete", mac: "Alt-Delete", run: deleteGroupForward},
  {mac: "Mod-Backspace", run: deleteToLineStart},
  {mac: "Mod-Delete", run: deleteToLineEnd}
].concat(/* @__PURE__ */ emacsStyleKeymap.map((b) => ({mac: b.key, run: b.run, shift: b.shift})));
var defaultTabBinding = {key: "Tab", run: insertTab, shift: indentSelection};

// node_modules/@codemirror/gutter/dist/index.js
var GutterMarker = class extends RangeValue {
  compare(other) {
    return this == other || this.constructor == other.constructor && this.eq(other);
  }
  toDOM(_view) {
    return null;
  }
  at(pos) {
    return this.range(pos);
  }
};
GutterMarker.prototype.elementClass = "";
GutterMarker.prototype.mapMode = MapMode.TrackBefore;
var defaults2 = {
  class: "",
  renderEmptyElements: false,
  elementStyle: "",
  markers: () => RangeSet.empty,
  lineMarker: () => null,
  initialSpacer: null,
  updateSpacer: null,
  domEventHandlers: {}
};
var activeGutters = Facet.define();
function gutter(config2) {
  return [gutters(), activeGutters.of(Object.assign(Object.assign({}, defaults2), config2))];
}
var baseTheme3 = EditorView.baseTheme({
  ".cm-gutters": {
    display: "flex",
    height: "100%",
    boxSizing: "border-box",
    left: 0
  },
  "&light .cm-gutters": {
    backgroundColor: "#f5f5f5",
    color: "#999",
    borderRight: "1px solid #ddd"
  },
  "&dark .cm-gutters": {
    backgroundColor: "#333338",
    color: "#ccc"
  },
  ".cm-gutter": {
    display: "flex !important",
    flexDirection: "column",
    flexShrink: 0,
    boxSizing: "border-box",
    height: "100%",
    overflow: "hidden"
  },
  ".cm-gutterElement": {
    boxSizing: "border-box"
  },
  ".cm-lineNumbers .cm-gutterElement": {
    padding: "0 3px 0 5px",
    minWidth: "20px",
    textAlign: "right",
    whiteSpace: "nowrap"
  }
});
var unfixGutters = Facet.define({
  combine: (values) => values.some((x) => x)
});
function gutters(config2) {
  let result = [
    gutterView,
    baseTheme3
  ];
  if (config2 && config2.fixed === false)
    result.push(unfixGutters.of(true));
  return result;
}
var gutterView = ViewPlugin.fromClass(class {
  constructor(view) {
    this.view = view;
    this.dom = document.createElement("div");
    this.dom.className = "cm-gutters";
    this.dom.setAttribute("aria-hidden", "true");
    this.gutters = view.state.facet(activeGutters).map((conf) => new SingleGutterView(view, conf));
    for (let gutter2 of this.gutters)
      this.dom.appendChild(gutter2.dom);
    this.fixed = !view.state.facet(unfixGutters);
    if (this.fixed) {
      this.dom.style.position = "sticky";
    }
    view.scrollDOM.insertBefore(this.dom, view.contentDOM);
  }
  update(update) {
    if (!this.updateGutters(update))
      return;
    let contexts = this.gutters.map((gutter2) => new UpdateContext(gutter2, this.view.viewport));
    this.view.viewportLines((line) => {
      let text;
      if (Array.isArray(line.type)) {
        for (let b of line.type)
          if (b.type == BlockType.Text) {
            text = b;
            break;
          }
      } else {
        text = line.type == BlockType.Text ? line : void 0;
      }
      if (!text)
        return;
      for (let cx of contexts)
        cx.line(this.view, text);
    }, 0);
    for (let cx of contexts)
      cx.finish();
    this.dom.style.minHeight = this.view.contentHeight + "px";
    if (update.state.facet(unfixGutters) != !this.fixed) {
      this.fixed = !this.fixed;
      this.dom.style.position = this.fixed ? "sticky" : "";
    }
  }
  updateGutters(update) {
    let prev = update.startState.facet(activeGutters), cur2 = update.state.facet(activeGutters);
    let change = update.docChanged || update.heightChanged || update.viewportChanged;
    if (prev == cur2) {
      for (let gutter2 of this.gutters)
        if (gutter2.update(update))
          change = true;
    } else {
      change = true;
      let gutters2 = [];
      for (let conf of cur2) {
        let known = prev.indexOf(conf);
        if (known < 0) {
          gutters2.push(new SingleGutterView(this.view, conf));
        } else {
          this.gutters[known].update(update);
          gutters2.push(this.gutters[known]);
        }
      }
      for (let g of this.gutters)
        g.dom.remove();
      for (let g of gutters2)
        this.dom.appendChild(g.dom);
      this.gutters = gutters2;
    }
    return change;
  }
  destroy() {
    this.dom.remove();
  }
}, {
  provide: PluginField.scrollMargins.from((value) => {
    if (value.gutters.length == 0 || !value.fixed)
      return null;
    return value.view.textDirection == Direction.LTR ? {left: value.dom.offsetWidth} : {right: value.dom.offsetWidth};
  })
});
function asArray2(val) {
  return Array.isArray(val) ? val : [val];
}
var UpdateContext = class {
  constructor(gutter2, viewport) {
    this.gutter = gutter2;
    this.localMarkers = [];
    this.i = 0;
    this.height = 0;
    this.cursor = RangeSet.iter(gutter2.markers, viewport.from);
  }
  line(view, line) {
    if (this.localMarkers.length)
      this.localMarkers = [];
    while (this.cursor.value && this.cursor.from <= line.from) {
      if (this.cursor.from == line.from)
        this.localMarkers.push(this.cursor.value);
      this.cursor.next();
    }
    let forLine = this.gutter.config.lineMarker(view, line, this.localMarkers);
    if (forLine)
      this.localMarkers.unshift(forLine);
    let gutter2 = this.gutter;
    if (this.localMarkers.length == 0 && !gutter2.config.renderEmptyElements)
      return;
    let above = line.top - this.height;
    if (this.i == gutter2.elements.length) {
      let newElt = new GutterElement(view, line.height, above, this.localMarkers);
      gutter2.elements.push(newElt);
      gutter2.dom.appendChild(newElt.dom);
    } else {
      let markers = this.localMarkers, elt = gutter2.elements[this.i];
      if (sameMarkers(markers, elt.markers)) {
        markers = elt.markers;
        this.localMarkers.length = 0;
      }
      elt.update(view, line.height, above, markers);
    }
    this.height = line.bottom;
    this.i++;
  }
  finish() {
    let gutter2 = this.gutter;
    while (gutter2.elements.length > this.i)
      gutter2.dom.removeChild(gutter2.elements.pop().dom);
  }
};
var SingleGutterView = class {
  constructor(view, config2) {
    this.view = view;
    this.config = config2;
    this.elements = [];
    this.spacer = null;
    this.dom = document.createElement("div");
    this.dom.className = "cm-gutter" + (this.config.class ? " " + this.config.class : "");
    for (let prop in config2.domEventHandlers) {
      this.dom.addEventListener(prop, (event) => {
        let line = view.visualLineAtHeight(event.clientY, view.contentDOM.getBoundingClientRect().top);
        if (config2.domEventHandlers[prop](view, line, event))
          event.preventDefault();
      });
    }
    this.markers = asArray2(config2.markers(view));
    if (config2.initialSpacer) {
      this.spacer = new GutterElement(view, 0, 0, [config2.initialSpacer(view)]);
      this.dom.appendChild(this.spacer.dom);
      this.spacer.dom.style.cssText += "visibility: hidden; pointer-events: none";
    }
  }
  update(update) {
    let prevMarkers = this.markers;
    this.markers = asArray2(this.config.markers(update.view));
    if (this.spacer && this.config.updateSpacer) {
      let updated = this.config.updateSpacer(this.spacer.markers[0], update);
      if (updated != this.spacer.markers[0])
        this.spacer.update(update.view, 0, 0, [updated]);
    }
    return this.markers != prevMarkers;
  }
};
var GutterElement = class {
  constructor(view, height, above, markers) {
    this.height = -1;
    this.above = 0;
    this.dom = document.createElement("div");
    this.update(view, height, above, markers);
  }
  update(view, height, above, markers) {
    if (this.height != height)
      this.dom.style.height = (this.height = height) + "px";
    if (this.above != above)
      this.dom.style.marginTop = (this.above = above) ? above + "px" : "";
    if (this.markers != markers) {
      this.markers = markers;
      for (let ch; ch = this.dom.lastChild; )
        ch.remove();
      let cls = "cm-gutterElement";
      for (let m of markers) {
        let dom = m.toDOM(view);
        if (dom)
          this.dom.appendChild(dom);
        let c2 = m.elementClass;
        if (c2)
          cls += " " + c2;
      }
      this.dom.className = cls;
    }
  }
};
function sameMarkers(a, b) {
  if (a.length != b.length)
    return false;
  for (let i = 0; i < a.length; i++)
    if (!a[i].compare(b[i]))
      return false;
  return true;
}
var lineNumberMarkers = Facet.define();
var lineNumberConfig = Facet.define({
  combine(values) {
    return combineConfig(values, {formatNumber: String, domEventHandlers: {}}, {
      domEventHandlers(a, b) {
        let result = Object.assign({}, a);
        for (let event in b) {
          let exists = result[event], add = b[event];
          result[event] = exists ? (view, line, event2) => exists(view, line, event2) || add(view, line, event2) : add;
        }
        return result;
      }
    });
  }
});
var NumberMarker = class extends GutterMarker {
  constructor(number2) {
    super();
    this.number = number2;
  }
  eq(other) {
    return this.number == other.number;
  }
  toDOM() {
    return document.createTextNode(this.number);
  }
};
function formatNumber(view, number2) {
  return view.state.facet(lineNumberConfig).formatNumber(number2, view.state);
}
var lineNumberGutter = gutter({
  class: "cm-lineNumbers",
  markers(view) {
    return view.state.facet(lineNumberMarkers);
  },
  lineMarker(view, line, others) {
    if (others.length)
      return null;
    return new NumberMarker(formatNumber(view, view.state.doc.lineAt(line.from).number));
  },
  initialSpacer(view) {
    return new NumberMarker(formatNumber(view, maxLineNumber(view.state.doc.lines)));
  },
  updateSpacer(spacer, update) {
    let max = formatNumber(update.view, maxLineNumber(update.view.state.doc.lines));
    return max == spacer.number ? spacer : new NumberMarker(max);
  }
});
function lineNumbers(config2 = {}) {
  return [
    lineNumberConfig.of(config2),
    lineNumberGutter
  ];
}
function maxLineNumber(lines) {
  let last = 9;
  while (last < lines)
    last = last * 10 + 9;
  return last;
}

// node_modules/@codemirror/highlight/dist/index.js
var nextTagID = 0;
var Tag = class {
  constructor(set, base3, modified) {
    this.set = set;
    this.base = base3;
    this.modified = modified;
    this.id = nextTagID++;
  }
  static define(parent) {
    if (parent === null || parent === void 0 ? void 0 : parent.base)
      throw new Error("Can not derive from a modified tag");
    let tag = new Tag([], null, []);
    tag.set.push(tag);
    if (parent)
      for (let t2 of parent.set)
        tag.set.push(t2);
    return tag;
  }
  static defineModifier() {
    let mod = new Modifier();
    return (tag) => {
      if (tag.modified.indexOf(mod) > -1)
        return tag;
      return Modifier.get(tag.base || tag, tag.modified.concat(mod).sort((a, b) => a.id - b.id));
    };
  }
};
var nextModifierID = 0;
var Modifier = class {
  constructor() {
    this.instances = [];
    this.id = nextModifierID++;
  }
  static get(base3, mods) {
    if (!mods.length)
      return base3;
    let exists = mods[0].instances.find((t2) => t2.base == base3 && sameArray2(mods, t2.modified));
    if (exists)
      return exists;
    let set = [], tag = new Tag(set, base3, mods);
    for (let m of mods)
      m.instances.push(tag);
    let configs = permute(mods);
    for (let parent of base3.set)
      for (let config2 of configs)
        set.push(Modifier.get(parent, config2));
    return tag;
  }
};
function sameArray2(a, b) {
  return a.length == b.length && a.every((x, i) => x == b[i]);
}
function permute(array) {
  let result = [array];
  for (let i = 0; i < array.length; i++) {
    for (let a of permute(array.slice(0, i).concat(array.slice(i + 1))))
      result.push(a);
  }
  return result;
}
function styleTags(spec) {
  let byName = Object.create(null);
  for (let prop in spec) {
    let tags2 = spec[prop];
    if (!Array.isArray(tags2))
      tags2 = [tags2];
    for (let part of prop.split(" "))
      if (part) {
        let pieces = [], mode = 2, rest = part;
        for (let pos = 0; ; ) {
          if (rest == "..." && pos > 0 && pos + 3 == part.length) {
            mode = 1;
            break;
          }
          let m = /^"(?:[^"\\]|\\.)*?"|[^\/!]+/.exec(rest);
          if (!m)
            throw new RangeError("Invalid path: " + part);
          pieces.push(m[0] == "*" ? null : m[0][0] == '"' ? JSON.parse(m[0]) : m[0]);
          pos += m[0].length;
          if (pos == part.length)
            break;
          let next = part[pos++];
          if (pos == part.length && next == "!") {
            mode = 0;
            break;
          }
          if (next != "/")
            throw new RangeError("Invalid path: " + part);
          rest = part.slice(pos);
        }
        let last = pieces.length - 1, inner = pieces[last];
        if (!inner)
          throw new RangeError("Invalid path: " + part);
        let rule = new Rule(tags2, mode, last > 0 ? pieces.slice(0, last) : null);
        byName[inner] = rule.sort(byName[inner]);
      }
  }
  return ruleNodeProp.add(byName);
}
var ruleNodeProp = new NodeProp();
var highlightStyle = Facet.define({
  combine(stylings) {
    return stylings.length ? HighlightStyle.combinedMatch(stylings) : null;
  }
});
var fallbackHighlightStyle = Facet.define({
  combine(values) {
    return values.length ? values[0].match : null;
  }
});
function noHighlight() {
  return null;
}
function getHighlightStyle(state) {
  return state.facet(highlightStyle) || state.facet(fallbackHighlightStyle) || noHighlight;
}
var Rule = class {
  constructor(tags2, mode, context, next) {
    this.tags = tags2;
    this.mode = mode;
    this.context = context;
    this.next = next;
  }
  sort(other) {
    if (!other || other.depth < this.depth) {
      this.next = other;
      return this;
    }
    other.next = this.sort(other.next);
    return other;
  }
  get depth() {
    return this.context ? this.context.length : 0;
  }
};
var HighlightStyle = class {
  constructor(spec, options) {
    this.map = Object.create(null);
    let modSpec;
    function def(spec2) {
      let cls = StyleModule.newName();
      (modSpec || (modSpec = Object.create(null)))["." + cls] = spec2;
      return cls;
    }
    this.all = typeof options.all == "string" ? options.all : options.all ? def(options.all) : null;
    for (let style of spec) {
      let cls = (style.class || def(Object.assign({}, style, {tag: null}))) + (this.all ? " " + this.all : "");
      let tags2 = style.tag;
      if (!Array.isArray(tags2))
        this.map[tags2.id] = cls;
      else
        for (let tag of tags2)
          this.map[tag.id] = cls;
    }
    this.module = modSpec ? new StyleModule(modSpec) : null;
    this.scope = options.scope || null;
    this.match = this.match.bind(this);
    let ext = [treeHighlighter];
    if (this.module)
      ext.push(EditorView.styleModule.of(this.module));
    this.extension = ext.concat(highlightStyle.of(this));
    this.fallback = ext.concat(fallbackHighlightStyle.of(this));
  }
  match(tag, scope) {
    if (this.scope && scope != this.scope)
      return null;
    for (let t2 of tag.set) {
      let match = this.map[t2.id];
      if (match !== void 0) {
        if (t2 != tag)
          this.map[tag.id] = match;
        return match;
      }
    }
    return this.map[tag.id] = this.all;
  }
  static combinedMatch(styles) {
    if (styles.length == 1)
      return styles[0].match;
    let cache = styles.some((s) => s.scope) ? void 0 : Object.create(null);
    return (tag, scope) => {
      let cached = cache && cache[tag.id];
      if (cached !== void 0)
        return cached;
      let result = null;
      for (let style of styles) {
        let value = style.match(tag, scope);
        if (value)
          result = result ? result + " " + value : value;
      }
      if (cache)
        cache[tag.id] = result;
      return result;
    };
  }
  static define(specs, options) {
    return new HighlightStyle(specs, options || {});
  }
  static get(state, tag, scope) {
    return getHighlightStyle(state)(tag, scope || NodeType.none);
  }
};
var TreeHighlighter = class {
  constructor(view) {
    this.markCache = Object.create(null);
    this.tree = syntaxTree(view.state);
    this.decorations = this.buildDeco(view, getHighlightStyle(view.state));
  }
  update(update) {
    let tree = syntaxTree(update.state), style = getHighlightStyle(update.state);
    let styleChange = style != update.startState.facet(highlightStyle);
    if (tree.length < update.view.viewport.to && !styleChange) {
      this.decorations = this.decorations.map(update.changes);
    } else if (tree != this.tree || update.viewportChanged || styleChange) {
      this.tree = tree;
      this.decorations = this.buildDeco(update.view, style);
    }
  }
  buildDeco(view, match) {
    if (match == noHighlight || !this.tree.length)
      return Decoration.none;
    let builder = new RangeSetBuilder();
    for (let {from, to} of view.visibleRanges) {
      highlightTreeRange(this.tree, from, to, match, (from2, to2, style) => {
        builder.add(from2, to2, this.markCache[style] || (this.markCache[style] = Decoration.mark({class: style})));
      });
    }
    return builder.finish();
  }
};
var treeHighlighter = Prec.fallback(ViewPlugin.fromClass(TreeHighlighter, {
  decorations: (v) => v.decorations
}));
var nodeStack = [""];
function highlightTreeRange(tree, from, to, style, span) {
  let spanStart = from, spanClass = "";
  let cursor = tree.topNode.cursor;
  function node(inheritedClass, depth, scope) {
    let {type: type2, from: start, to: end} = cursor;
    if (start >= to || end <= from)
      return;
    nodeStack[depth] = type2.name;
    if (type2.isTop)
      scope = type2;
    let cls = inheritedClass;
    let rule = type2.prop(ruleNodeProp), opaque = false;
    while (rule) {
      if (!rule.context || matchContext(rule.context, nodeStack, depth)) {
        for (let tag of rule.tags) {
          let st = style(tag, scope);
          if (st) {
            if (cls)
              cls += " ";
            cls += st;
            if (rule.mode == 1)
              inheritedClass += (inheritedClass ? " " : "") + st;
            else if (rule.mode == 0)
              opaque = true;
          }
        }
        break;
      }
      rule = rule.next;
    }
    if (cls != spanClass) {
      if (start > spanStart && spanClass)
        span(spanStart, cursor.from, spanClass);
      spanStart = start;
      spanClass = cls;
    }
    if (!opaque && cursor.firstChild()) {
      do {
        let end2 = cursor.to;
        node(inheritedClass, depth + 1, scope);
        if (spanClass != cls) {
          let pos = Math.min(to, end2);
          if (pos > spanStart && spanClass)
            span(spanStart, pos, spanClass);
          spanStart = pos;
          spanClass = cls;
        }
      } while (cursor.nextSibling());
      cursor.parent();
    }
  }
  node("", 0, tree.type);
}
function matchContext(context, stack, depth) {
  if (context.length > depth - 1)
    return false;
  for (let d = depth - 1, i = context.length - 1; i >= 0; i--, d--) {
    let check = context[i];
    if (check && check != stack[d])
      return false;
  }
  return true;
}
var t = Tag.define;
var comment = t();
var name = t();
var typeName = t(name);
var literal = t();
var string = t(literal);
var number = t(literal);
var content = t();
var heading = t(content);
var keyword = t();
var operator = t();
var punctuation = t();
var bracket = t(punctuation);
var meta = t();
var tags = {
  comment,
  lineComment: t(comment),
  blockComment: t(comment),
  docComment: t(comment),
  name,
  variableName: t(name),
  typeName,
  tagName: t(typeName),
  propertyName: t(name),
  className: t(name),
  labelName: t(name),
  namespace: t(name),
  macroName: t(name),
  literal,
  string,
  docString: t(string),
  character: t(string),
  number,
  integer: t(number),
  float: t(number),
  bool: t(literal),
  regexp: t(literal),
  escape: t(literal),
  color: t(literal),
  url: t(literal),
  keyword,
  self: t(keyword),
  null: t(keyword),
  atom: t(keyword),
  unit: t(keyword),
  modifier: t(keyword),
  operatorKeyword: t(keyword),
  controlKeyword: t(keyword),
  definitionKeyword: t(keyword),
  operator,
  derefOperator: t(operator),
  arithmeticOperator: t(operator),
  logicOperator: t(operator),
  bitwiseOperator: t(operator),
  compareOperator: t(operator),
  updateOperator: t(operator),
  definitionOperator: t(operator),
  typeOperator: t(operator),
  controlOperator: t(operator),
  punctuation,
  separator: t(punctuation),
  bracket,
  angleBracket: t(bracket),
  squareBracket: t(bracket),
  paren: t(bracket),
  brace: t(bracket),
  content,
  heading,
  heading1: t(heading),
  heading2: t(heading),
  heading3: t(heading),
  heading4: t(heading),
  heading5: t(heading),
  heading6: t(heading),
  contentSeparator: t(content),
  list: t(content),
  quote: t(content),
  emphasis: t(content),
  strong: t(content),
  link: t(content),
  monospace: t(content),
  inserted: t(),
  deleted: t(),
  changed: t(),
  invalid: t(),
  meta,
  documentMeta: t(meta),
  annotation: t(meta),
  processingInstruction: t(meta),
  definition: Tag.defineModifier(),
  constant: Tag.defineModifier(),
  function: Tag.defineModifier(),
  standard: Tag.defineModifier(),
  local: Tag.defineModifier(),
  special: Tag.defineModifier()
};
var defaultHighlightStyle = HighlightStyle.define([
  {
    tag: tags.link,
    textDecoration: "underline"
  },
  {
    tag: tags.heading,
    textDecoration: "underline",
    fontWeight: "bold"
  },
  {
    tag: tags.emphasis,
    fontStyle: "italic"
  },
  {
    tag: tags.strong,
    fontWeight: "bold"
  },
  {
    tag: tags.keyword,
    color: "#708"
  },
  {
    tag: [tags.atom, tags.bool, tags.url, tags.contentSeparator, tags.labelName],
    color: "#219"
  },
  {
    tag: [tags.literal, tags.inserted],
    color: "#164"
  },
  {
    tag: [tags.string, tags.deleted],
    color: "#a11"
  },
  {
    tag: [tags.regexp, tags.escape, tags.special(tags.string)],
    color: "#e40"
  },
  {
    tag: tags.definition(tags.variableName),
    color: "#00f"
  },
  {
    tag: tags.local(tags.variableName),
    color: "#30a"
  },
  {
    tag: [tags.typeName, tags.namespace],
    color: "#085"
  },
  {
    tag: tags.className,
    color: "#167"
  },
  {
    tag: [tags.special(tags.variableName), tags.macroName],
    color: "#256"
  },
  {
    tag: tags.definition(tags.propertyName),
    color: "#00c"
  },
  {
    tag: tags.comment,
    color: "#940"
  },
  {
    tag: tags.meta,
    color: "#7a757a"
  },
  {
    tag: tags.invalid,
    color: "#f00"
  }
]);
var classHighlightStyle = HighlightStyle.define([
  {tag: tags.link, class: "cmt-link"},
  {tag: tags.heading, class: "cmt-heading"},
  {tag: tags.emphasis, class: "cmt-emphasis"},
  {tag: tags.strong, class: "cmt-strong"},
  {tag: tags.keyword, class: "cmt-keyword"},
  {tag: tags.atom, class: "cmt-atom"},
  {tag: tags.bool, class: "cmt-bool"},
  {tag: tags.url, class: "cmt-url"},
  {tag: tags.labelName, class: "cmt-labelName"},
  {tag: tags.inserted, class: "cmt-inserted"},
  {tag: tags.deleted, class: "cmt-deleted"},
  {tag: tags.literal, class: "cmt-literal"},
  {tag: tags.string, class: "cmt-string"},
  {tag: tags.number, class: "cmt-number"},
  {tag: [tags.regexp, tags.escape, tags.special(tags.string)], class: "cmt-string2"},
  {tag: tags.variableName, class: "cmt-variableName"},
  {tag: tags.local(tags.variableName), class: "cmt-variableName cmt-local"},
  {tag: tags.definition(tags.variableName), class: "cmt-variableName cmt-definition"},
  {tag: tags.special(tags.variableName), class: "cmt-variableName2"},
  {tag: tags.typeName, class: "cmt-typeName"},
  {tag: tags.namespace, class: "cmt-namespace"},
  {tag: tags.macroName, class: "cmt-macroName"},
  {tag: tags.propertyName, class: "cmt-propertyName"},
  {tag: tags.operator, class: "cmt-operator"},
  {tag: tags.comment, class: "cmt-comment"},
  {tag: tags.meta, class: "cmt-meta"},
  {tag: tags.invalid, class: "cmt-invalid"},
  {tag: tags.punctuation, class: "cmt-punctuation"}
]);

// node_modules/@codemirror/stream-parser/dist/index.js
function countCol(string3, end, tabSize, startIndex = 0, startValue = 0) {
  if (end == null) {
    end = string3.search(/[^\s\u00a0]/);
    if (end == -1)
      end = string3.length;
  }
  return countColumn(string3.slice(startIndex, end), startValue, tabSize);
}
var StringStream = class {
  constructor(string3, tabSize, indentUnit2) {
    this.string = string3;
    this.tabSize = tabSize;
    this.indentUnit = indentUnit2;
    this.pos = 0;
    this.start = 0;
    this.lastColumnPos = 0;
    this.lastColumnValue = 0;
  }
  eol() {
    return this.pos >= this.string.length;
  }
  sol() {
    return this.pos == 0;
  }
  peek() {
    return this.string.charAt(this.pos) || void 0;
  }
  next() {
    if (this.pos < this.string.length)
      return this.string.charAt(this.pos++);
  }
  eat(match) {
    let ch = this.string.charAt(this.pos);
    let ok;
    if (typeof match == "string")
      ok = ch == match;
    else
      ok = ch && (match instanceof RegExp ? match.test(ch) : match(ch));
    if (ok) {
      ++this.pos;
      return ch;
    }
  }
  eatWhile(match) {
    let start = this.pos;
    while (this.eat(match)) {
    }
    return this.pos > start;
  }
  eatSpace() {
    let start = this.pos;
    while (/[\s\u00a0]/.test(this.string.charAt(this.pos)))
      ++this.pos;
    return this.pos > start;
  }
  skipToEnd() {
    this.pos = this.string.length;
  }
  skipTo(ch) {
    let found = this.string.indexOf(ch, this.pos);
    if (found > -1) {
      this.pos = found;
      return true;
    }
  }
  backUp(n) {
    this.pos -= n;
  }
  column() {
    if (this.lastColumnPos < this.start) {
      this.lastColumnValue = countCol(this.string, this.start, this.tabSize, this.lastColumnPos, this.lastColumnValue);
      this.lastColumnPos = this.start;
    }
    return this.lastColumnValue;
  }
  indentation() {
    return countCol(this.string, null, this.tabSize);
  }
  match(pattern, consume, caseInsensitive) {
    if (typeof pattern == "string") {
      let cased = (str) => caseInsensitive ? str.toLowerCase() : str;
      let substr = this.string.substr(this.pos, pattern.length);
      if (cased(substr) == cased(pattern)) {
        if (consume !== false)
          this.pos += pattern.length;
        return true;
      } else
        return null;
    } else {
      let match = this.string.slice(this.pos).match(pattern);
      if (match && match.index > 0)
        return null;
      if (match && consume !== false)
        this.pos += match[0].length;
      return match;
    }
  }
  current() {
    return this.string.slice(this.start, this.pos);
  }
};
function fullParser(spec) {
  return {
    token: spec.token,
    blankLine: spec.blankLine || (() => {
    }),
    startState: spec.startState || (() => true),
    copyState: spec.copyState || defaultCopyState,
    indent: spec.indent || (() => null),
    languageData: spec.languageData || {}
  };
}
function defaultCopyState(state) {
  if (typeof state != "object")
    return state;
  let newState = {};
  for (let prop in state) {
    let val = state[prop];
    newState[prop] = val instanceof Array ? val.slice() : val;
  }
  return newState;
}
var StreamLanguage = class extends Language {
  constructor(parser6) {
    let data = defineLanguageFacet(parser6.languageData);
    let p = fullParser(parser6);
    let startParse = (input, startPos, context) => new Parse(this, input, startPos, context);
    super(data, {startParse}, docID(data), [indentService.of((cx, pos) => this.getIndent(cx, pos))]);
    this.streamParser = p;
    this.stateAfter = new WeakMap();
  }
  static define(spec) {
    return new StreamLanguage(spec);
  }
  getIndent(cx, pos) {
    let tree = syntaxTree(cx.state), at = tree.resolve(pos);
    while (at && at.type != this.topNode)
      at = at.parent;
    if (!at)
      return null;
    let start = findState(this, tree, 0, at.from, pos), statePos, state;
    if (start) {
      state = start.state;
      statePos = start.pos + 1;
    } else {
      state = this.streamParser.startState(cx.unit);
      statePos = 0;
    }
    if (pos - statePos > 1e4)
      return null;
    while (statePos < pos) {
      let line = cx.state.doc.lineAt(statePos), end = Math.min(pos, line.to);
      if (line.length) {
        let stream = new StringStream(line.text, cx.state.tabSize, cx.unit);
        while (stream.pos < end - line.from)
          readToken(this.streamParser.token, stream, state);
      } else {
        this.streamParser.blankLine(state, cx.unit);
      }
      if (end == pos)
        break;
      statePos = line.to + 1;
    }
    let {text} = cx.state.doc.lineAt(pos);
    return this.streamParser.indent(state, /^\s*(.*)/.exec(text)[1], cx);
  }
  get allowsNesting() {
    return false;
  }
};
function findState(lang, tree, off, startPos, before) {
  let state = off >= startPos && off + tree.length <= before && lang.stateAfter.get(tree);
  if (state)
    return {state: lang.streamParser.copyState(state), pos: off + tree.length};
  for (let i = tree.children.length - 1; i >= 0; i--) {
    let child = tree.children[i], pos = off + tree.positions[i];
    let found = child instanceof Tree && pos < before && findState(lang, child, pos, startPos, before);
    if (found)
      return found;
  }
  return null;
}
function cutTree(lang, tree, from, to, inside2) {
  if (inside2 && from <= 0 && to >= tree.length)
    return tree;
  if (!inside2 && tree.type == lang.topNode)
    inside2 = true;
  for (let i = tree.children.length - 1; i >= 0; i--) {
    let pos = tree.positions[i] + from, child = tree.children[i], inner;
    if (pos < to && child instanceof Tree) {
      if (!(inner = cutTree(lang, child, from - pos, to - pos, inside2)))
        break;
      return !inside2 ? inner : new Tree(tree.type, tree.children.slice(0, i).concat(inner), tree.positions.slice(0, i + 1), pos + inner.length);
    }
  }
  return null;
}
function findStartInFragments(lang, fragments, startPos, state) {
  for (let f of fragments) {
    let found = f.from <= startPos && f.to > startPos && findState(lang, f.tree, 0 - f.offset, startPos, f.to), tree;
    if (found && (tree = cutTree(lang, f.tree, startPos + f.offset, found.pos + f.offset, false)))
      return {state: found.state, tree};
  }
  return {state: lang.streamParser.startState(getIndentUnit(state)), tree: Tree.empty};
}
var Parse = class {
  constructor(lang, input, startPos, context) {
    this.lang = lang;
    this.input = input;
    this.startPos = startPos;
    this.context = context;
    this.chunks = [];
    this.chunkPos = [];
    this.chunk = [];
    let {state, tree} = findStartInFragments(lang, context.fragments, startPos, context.state);
    this.state = state;
    this.pos = this.chunkStart = startPos + tree.length;
    if (tree.length) {
      this.chunks.push(tree);
      this.chunkPos.push(0);
    }
    if (this.pos < context.viewport.from - 1e5) {
      this.state = this.lang.streamParser.startState(getIndentUnit(context.state));
      context.skipUntilInView(this.pos, context.viewport.from);
      this.pos = context.viewport.from;
    }
  }
  advance() {
    let end = Math.min(this.context.viewport.to, this.input.length, this.chunkStart + 2048);
    while (this.pos < end)
      this.parseLine();
    if (this.chunkStart < this.pos)
      this.finishChunk();
    if (end < this.input.length && this.pos < this.context.viewport.to)
      return null;
    this.context.skipUntilInView(this.pos, this.input.length);
    return this.finish();
  }
  parseLine() {
    let line = this.input.lineAfter(this.pos), {streamParser} = this.lang;
    let stream = new StringStream(line, this.context ? this.context.state.tabSize : 4, getIndentUnit(this.context.state));
    if (stream.eol()) {
      streamParser.blankLine(this.state, stream.indentUnit);
    } else {
      while (!stream.eol()) {
        let token = readToken(streamParser.token, stream, this.state);
        if (token)
          this.chunk.push(tokenID(token), this.pos + stream.start, this.pos + stream.pos, 4);
      }
    }
    this.pos += line.length;
    if (this.pos < this.input.length)
      this.pos++;
  }
  finishChunk() {
    let tree = Tree.build({
      buffer: this.chunk,
      start: this.chunkStart,
      length: this.pos - this.chunkStart,
      nodeSet,
      topID: 0,
      maxBufferLength: 2048
    });
    this.lang.stateAfter.set(tree, this.lang.streamParser.copyState(this.state));
    this.chunks.push(tree);
    this.chunkPos.push(this.chunkStart - this.startPos);
    this.chunk = [];
    this.chunkStart = this.pos;
  }
  finish() {
    return new Tree(this.lang.topNode, this.chunks, this.chunkPos, this.pos - this.startPos).balance();
  }
  forceFinish() {
    return this.finish();
  }
};
function readToken(token, stream, state) {
  stream.start = stream.pos;
  for (let i = 0; i < 10; i++) {
    let result = token(stream, state);
    if (stream.pos > stream.start)
      return result;
  }
  throw new Error("Stream parser failed to advance stream.");
}
var tokenTable = /* @__PURE__ */ Object.create(null);
var typeArray = [NodeType.none];
var nodeSet = /* @__PURE__ */ new NodeSet(typeArray);
var warned = [];
function tokenID(tag) {
  return !tag ? 0 : tokenTable[tag] || (tokenTable[tag] = createTokenType(tag));
}
for (let [legacyName, name2] of [
  ["variable", "variableName"],
  ["variable-2", "variableName.special"],
  ["string-2", "string.special"],
  ["def", "variableName.definition"],
  ["tag", "typeName"],
  ["attribute", "propertyName"],
  ["type", "typeName"],
  ["builtin", "variableName.standard"],
  ["qualifier", "modifier"],
  ["error", "invalid"],
  ["header", "heading"],
  ["property", "propertyName"]
])
  tokenTable[legacyName] = /* @__PURE__ */ tokenID(name2);
function warnForPart(part, msg) {
  if (warned.indexOf(part) > -1)
    return;
  warned.push(part);
  console.warn(msg);
}
function createTokenType(tagStr) {
  let tag = null;
  for (let part of tagStr.split(".")) {
    let value = tags[part];
    if (!value) {
      warnForPart(part, `Unknown highlighting tag ${part}`);
    } else if (typeof value == "function") {
      if (!tag)
        warnForPart(part, `Modifier ${part} used at start of tag`);
      else
        tag = value(tag);
    } else {
      if (tag)
        warnForPart(part, `Tag ${part} used as modifier`);
      else
        tag = value;
    }
  }
  if (!tag)
    return 0;
  let name2 = tagStr.replace(/ /g, "_"), type2 = NodeType.define({
    id: typeArray.length,
    name: name2,
    props: [styleTags({[name2]: tag})]
  });
  typeArray.push(type2);
  return type2.id;
}
function docID(data) {
  let type2 = NodeType.define({id: typeArray.length, name: "Document", props: [languageDataProp.add(() => data)]});
  typeArray.push(type2);
  return type2;
}

// node_modules/@codemirror/legacy-modes/mode/brainfuck.js
var reserve = "><+-.,[]".split("");
var brainfuck = {
  startState: function() {
    return {
      commentLine: false,
      left: 0,
      right: 0,
      commentLoop: false
    };
  },
  token: function(stream, state) {
    if (stream.eatSpace())
      return null;
    if (stream.sol()) {
      state.commentLine = false;
    }
    var ch = stream.next().toString();
    if (reserve.indexOf(ch) !== -1) {
      if (state.commentLine === true) {
        if (stream.eol()) {
          state.commentLine = false;
        }
        return "comment";
      }
      if (ch === "]" || ch === "[") {
        if (ch === "[") {
          state.left++;
        } else {
          state.right++;
        }
        return "bracket";
      } else if (ch === "+" || ch === "-") {
        return "keyword";
      } else if (ch === "<" || ch === ">") {
        return "atom";
      } else if (ch === "." || ch === ",") {
        return "def";
      }
    } else {
      state.commentLine = true;
      if (stream.eol()) {
        state.commentLine = false;
      }
      return "comment";
    }
    if (stream.eol()) {
      state.commentLine = false;
    }
  }
};

// node_modules/@codemirror/legacy-modes/mode/clike.js
function Context(indented, column, type2, info, align, prev) {
  this.indented = indented;
  this.column = column;
  this.type = type2;
  this.info = info;
  this.align = align;
  this.prev = prev;
}
function pushContext(state, col, type2, info) {
  var indent2 = state.indented;
  if (state.context && state.context.type == "statement" && type2 != "statement")
    indent2 = state.context.indented;
  return state.context = new Context(indent2, col, type2, info, null, state.context);
}
function popContext(state) {
  var t2 = state.context.type;
  if (t2 == ")" || t2 == "]" || t2 == "}")
    state.indented = state.context.indented;
  return state.context = state.context.prev;
}
function typeBefore(stream, state, pos) {
  if (state.prevToken == "variable" || state.prevToken == "type")
    return true;
  if (/\S(?:[^- ]>|[*\]])\s*$|\*$/.test(stream.string.slice(0, pos)))
    return true;
  if (state.typeAtEndOfLine && stream.column() == stream.indentation())
    return true;
}
function isTopScope(context) {
  for (; ; ) {
    if (!context || context.type == "top")
      return true;
    if (context.type == "}" && context.prev.info != "namespace")
      return false;
    context = context.prev;
  }
}
function clike(parserConfig) {
  var statementIndentUnit = parserConfig.statementIndentUnit, dontAlignCalls = parserConfig.dontAlignCalls, keywords11 = parserConfig.keywords || {}, types4 = parserConfig.types || {}, builtin = parserConfig.builtin || {}, blockKeywords = parserConfig.blockKeywords || {}, defKeywords = parserConfig.defKeywords || {}, atoms4 = parserConfig.atoms || {}, hooks = parserConfig.hooks || {}, multiLineStrings = parserConfig.multiLineStrings, indentStatements = parserConfig.indentStatements !== false, indentSwitch = parserConfig.indentSwitch !== false, namespaceSeparator = parserConfig.namespaceSeparator, isPunctuationChar = parserConfig.isPunctuationChar || /[\[\]{}\(\),;\:\.]/, numberStart = parserConfig.numberStart || /[\d\.]/, number2 = parserConfig.number || /^(?:0x[a-f\d]+|0b[01]+|(?:\d+\.?\d*|\.\d+)(?:e[-+]?\d+)?)(u|ll?|l|f)?/i, isOperatorChar3 = parserConfig.isOperatorChar || /[+\-*&%=<>!?|\/]/, isIdentifierChar = parserConfig.isIdentifierChar || /[\w\$_\xa1-\uffff]/, isReservedIdentifier = parserConfig.isReservedIdentifier || false;
  var curPunc3, isDefKeyword;
  function tokenBase9(stream, state) {
    var ch = stream.next();
    if (hooks[ch]) {
      var result = hooks[ch](stream, state);
      if (result !== false)
        return result;
    }
    if (ch == '"' || ch == "'") {
      state.tokenize = tokenString5(ch);
      return state.tokenize(stream, state);
    }
    if (numberStart.test(ch)) {
      stream.backUp(1);
      if (stream.match(number2))
        return "number";
      stream.next();
    }
    if (isPunctuationChar.test(ch)) {
      curPunc3 = ch;
      return null;
    }
    if (ch == "/") {
      if (stream.eat("*")) {
        state.tokenize = tokenComment5;
        return tokenComment5(stream, state);
      }
      if (stream.eat("/")) {
        stream.skipToEnd();
        return "comment";
      }
    }
    if (isOperatorChar3.test(ch)) {
      while (!stream.match(/^\/[\/*]/, false) && stream.eat(isOperatorChar3)) {
      }
      return "operator";
    }
    stream.eatWhile(isIdentifierChar);
    if (namespaceSeparator)
      while (stream.match(namespaceSeparator))
        stream.eatWhile(isIdentifierChar);
    var cur2 = stream.current();
    if (contains2(keywords11, cur2)) {
      if (contains2(blockKeywords, cur2))
        curPunc3 = "newstatement";
      if (contains2(defKeywords, cur2))
        isDefKeyword = true;
      return "keyword";
    }
    if (contains2(types4, cur2))
      return "type";
    if (contains2(builtin, cur2) || isReservedIdentifier && isReservedIdentifier(cur2)) {
      if (contains2(blockKeywords, cur2))
        curPunc3 = "newstatement";
      return "builtin";
    }
    if (contains2(atoms4, cur2))
      return "atom";
    return "variable";
  }
  function tokenString5(quote) {
    return function(stream, state) {
      var escaped = false, next, end = false;
      while ((next = stream.next()) != null) {
        if (next == quote && !escaped) {
          end = true;
          break;
        }
        escaped = !escaped && next == "\\";
      }
      if (end || !(escaped || multiLineStrings))
        state.tokenize = null;
      return "string";
    };
  }
  function tokenComment5(stream, state) {
    var maybeEnd = false, ch;
    while (ch = stream.next()) {
      if (ch == "/" && maybeEnd) {
        state.tokenize = null;
        break;
      }
      maybeEnd = ch == "*";
    }
    return "comment";
  }
  function maybeEOL(stream, state) {
    if (parserConfig.typeFirstDefinitions && stream.eol() && isTopScope(state.context))
      state.typeAtEndOfLine = typeBefore(stream, state, stream.pos);
  }
  return {
    startState: function(indentUnit2) {
      return {
        tokenize: null,
        context: new Context(-indentUnit2, 0, "top", null, false),
        indented: 0,
        startOfLine: true,
        prevToken: null
      };
    },
    token: function(stream, state) {
      var ctx = state.context;
      if (stream.sol()) {
        if (ctx.align == null)
          ctx.align = false;
        state.indented = stream.indentation();
        state.startOfLine = true;
      }
      if (stream.eatSpace()) {
        maybeEOL(stream, state);
        return null;
      }
      curPunc3 = isDefKeyword = null;
      var style = (state.tokenize || tokenBase9)(stream, state);
      if (style == "comment" || style == "meta")
        return style;
      if (ctx.align == null)
        ctx.align = true;
      if (curPunc3 == ";" || curPunc3 == ":" || curPunc3 == "," && stream.match(/^\s*(?:\/\/.*)?$/, false))
        while (state.context.type == "statement")
          popContext(state);
      else if (curPunc3 == "{")
        pushContext(state, stream.column(), "}");
      else if (curPunc3 == "[")
        pushContext(state, stream.column(), "]");
      else if (curPunc3 == "(")
        pushContext(state, stream.column(), ")");
      else if (curPunc3 == "}") {
        while (ctx.type == "statement")
          ctx = popContext(state);
        if (ctx.type == "}")
          ctx = popContext(state);
        while (ctx.type == "statement")
          ctx = popContext(state);
      } else if (curPunc3 == ctx.type)
        popContext(state);
      else if (indentStatements && ((ctx.type == "}" || ctx.type == "top") && curPunc3 != ";" || ctx.type == "statement" && curPunc3 == "newstatement")) {
        pushContext(state, stream.column(), "statement", stream.current());
      }
      if (style == "variable" && (state.prevToken == "def" || parserConfig.typeFirstDefinitions && typeBefore(stream, state, stream.start) && isTopScope(state.context) && stream.match(/^\s*\(/, false)))
        style = "def";
      if (hooks.token) {
        var result = hooks.token(stream, state, style);
        if (result !== void 0)
          style = result;
      }
      if (style == "def" && parserConfig.styleDefs === false)
        style = "variable";
      state.startOfLine = false;
      state.prevToken = isDefKeyword ? "def" : style || curPunc3;
      maybeEOL(stream, state);
      return style;
    },
    indent: function(state, textAfter, context) {
      if (state.tokenize != tokenBase9 && state.tokenize != null || state.typeAtEndOfLine)
        return null;
      var ctx = state.context, firstChar = textAfter && textAfter.charAt(0);
      var closing3 = firstChar == ctx.type;
      if (ctx.type == "statement" && firstChar == "}")
        ctx = ctx.prev;
      if (parserConfig.dontIndentStatements)
        while (ctx.type == "statement" && parserConfig.dontIndentStatements.test(ctx.info))
          ctx = ctx.prev;
      if (hooks.indent) {
        var hook = hooks.indent(state, ctx, textAfter, context.unit);
        if (typeof hook == "number")
          return hook;
      }
      var switchBlock = ctx.prev && ctx.prev.info == "switch";
      if (parserConfig.allmanIndentation && /[{(]/.test(firstChar)) {
        while (ctx.type != "top" && ctx.type != "}")
          ctx = ctx.prev;
        return ctx.indented;
      }
      if (ctx.type == "statement")
        return ctx.indented + (firstChar == "{" ? 0 : statementIndentUnit || context.unit);
      if (ctx.align && (!dontAlignCalls || ctx.type != ")"))
        return ctx.column + (closing3 ? 0 : 1);
      if (ctx.type == ")" && !closing3)
        return ctx.indented + (statementIndentUnit || context.unit);
      return ctx.indented + (closing3 ? 0 : context.unit) + (!closing3 && switchBlock && !/^(?:case|default)\b/.test(textAfter) ? context.unit : 0);
    },
    languageData: {
      indentOnInput: indentSwitch ? /^\s*(?:case .*?:|default:|\{\}?|\})$/ : /^\s*[{}]$/,
      commentTokens: {line: "//", block: {open: "/*", close: "*/"}},
      autocomplete: Object.keys(keywords11).concat(Object.keys(types4)).concat(Object.keys(builtin)).concat(Object.keys(atoms4)),
      ...parserConfig.languageData
    }
  };
}
function words(str) {
  var obj = {}, words4 = str.split(" ");
  for (var i = 0; i < words4.length; ++i)
    obj[words4[i]] = true;
  return obj;
}
function contains2(words4, word) {
  if (typeof words4 === "function") {
    return words4(word);
  } else {
    return words4.propertyIsEnumerable(word);
  }
}
var cKeywords = "auto if break case register continue return default do sizeof static else struct switch extern typedef union for goto while enum const volatile inline restrict asm fortran";
var cppKeywords = "alignas alignof and and_eq audit axiom bitand bitor catch class compl concept constexpr const_cast decltype delete dynamic_cast explicit export final friend import module mutable namespace new noexcept not not_eq operator or or_eq override private protected public reinterpret_cast requires static_assert static_cast template this thread_local throw try typeid typename using virtual xor xor_eq";
var objCKeywords = "bycopy byref in inout oneway out self super atomic nonatomic retain copy readwrite readonly strong weak assign typeof nullable nonnull null_resettable _cmd @interface @implementation @end @protocol @encode @property @synthesize @dynamic @class @public @package @private @protected @required @optional @try @catch @finally @import @selector @encode @defs @synchronized @autoreleasepool @compatibility_alias @available";
var objCBuiltins = "FOUNDATION_EXPORT FOUNDATION_EXTERN NS_INLINE NS_FORMAT_FUNCTION  NS_RETURNS_RETAINEDNS_ERROR_ENUM NS_RETURNS_NOT_RETAINED NS_RETURNS_INNER_POINTER NS_DESIGNATED_INITIALIZER NS_ENUM NS_OPTIONS NS_REQUIRES_NIL_TERMINATION NS_ASSUME_NONNULL_BEGIN NS_ASSUME_NONNULL_END NS_SWIFT_NAME NS_REFINED_FOR_SWIFT";
var basicCTypes = words("int long char short double float unsigned signed void bool");
var basicObjCTypes = words("SEL instancetype id Class Protocol BOOL");
function cTypes(identifier2) {
  return contains2(basicCTypes, identifier2) || /.+_t$/.test(identifier2);
}
function objCTypes(identifier2) {
  return cTypes(identifier2) || contains2(basicObjCTypes, identifier2);
}
var cBlockKeywords = "case do else for if switch while struct enum union";
var cDefKeywords = "struct enum union";
function cppHook(stream, state) {
  if (!state.startOfLine)
    return false;
  for (var ch, next = null; ch = stream.peek(); ) {
    if (ch == "\\" && stream.match(/^.$/)) {
      next = cppHook;
      break;
    } else if (ch == "/" && stream.match(/^\/[\/\*]/, false)) {
      break;
    }
    stream.next();
  }
  state.tokenize = next;
  return "meta";
}
function pointerHook(_stream, state) {
  if (state.prevToken == "type")
    return "type";
  return false;
}
function cIsReservedIdentifier(token) {
  if (!token || token.length < 2)
    return false;
  if (token[0] != "_")
    return false;
  return token[1] == "_" || token[1] !== token[1].toLowerCase();
}
function cpp14Literal(stream) {
  stream.eatWhile(/[\w\.']/);
  return "number";
}
function cpp11StringHook(stream, state) {
  stream.backUp(1);
  if (stream.match(/^(?:R|u8R|uR|UR|LR)/)) {
    var match = stream.match(/^"([^\s\\()]{0,16})\(/);
    if (!match) {
      return false;
    }
    state.cpp11RawStringDelim = match[1];
    state.tokenize = tokenRawString;
    return tokenRawString(stream, state);
  }
  if (stream.match(/^(?:u8|u|U|L)/)) {
    if (stream.match(/^["']/, false)) {
      return "string";
    }
    return false;
  }
  stream.next();
  return false;
}
function cppLooksLikeConstructor(word) {
  var lastTwo = /(\w+)::~?(\w+)$/.exec(word);
  return lastTwo && lastTwo[1] == lastTwo[2];
}
function tokenAtString(stream, state) {
  var next;
  while ((next = stream.next()) != null) {
    if (next == '"' && !stream.eat('"')) {
      state.tokenize = null;
      break;
    }
  }
  return "string";
}
function tokenRawString(stream, state) {
  var delim = state.cpp11RawStringDelim.replace(/[^\w\s]/g, "\\$&");
  var match = stream.match(new RegExp(".*?\\)" + delim + '"'));
  if (match)
    state.tokenize = null;
  else
    stream.skipToEnd();
  return "string";
}
var c = clike({
  keywords: words(cKeywords),
  types: cTypes,
  blockKeywords: words(cBlockKeywords),
  defKeywords: words(cDefKeywords),
  typeFirstDefinitions: true,
  atoms: words("NULL true false"),
  isReservedIdentifier: cIsReservedIdentifier,
  hooks: {
    "#": cppHook,
    "*": pointerHook
  }
});
var cpp = clike({
  keywords: words(cKeywords + " " + cppKeywords),
  types: cTypes,
  blockKeywords: words(cBlockKeywords + " class try catch"),
  defKeywords: words(cDefKeywords + " class namespace"),
  typeFirstDefinitions: true,
  atoms: words("true false NULL nullptr"),
  dontIndentStatements: /^template$/,
  isIdentifierChar: /[\w\$_~\xa1-\uffff]/,
  isReservedIdentifier: cIsReservedIdentifier,
  hooks: {
    "#": cppHook,
    "*": pointerHook,
    u: cpp11StringHook,
    U: cpp11StringHook,
    L: cpp11StringHook,
    R: cpp11StringHook,
    "0": cpp14Literal,
    "1": cpp14Literal,
    "2": cpp14Literal,
    "3": cpp14Literal,
    "4": cpp14Literal,
    "5": cpp14Literal,
    "6": cpp14Literal,
    "7": cpp14Literal,
    "8": cpp14Literal,
    "9": cpp14Literal,
    token: function(stream, state, style) {
      if (style == "variable" && stream.peek() == "(" && (state.prevToken == ";" || state.prevToken == null || state.prevToken == "}") && cppLooksLikeConstructor(stream.current()))
        return "def";
    }
  },
  namespaceSeparator: "::"
});
var java = clike({
  keywords: words("abstract assert break case catch class const continue default do else enum extends final finally for goto if implements import instanceof interface native new package private protected public return static strictfp super switch synchronized this throw throws transient try volatile while @interface"),
  types: words("byte short int long float double boolean char void Boolean Byte Character Double Float Integer Long Number Object Short String StringBuffer StringBuilder Void"),
  blockKeywords: words("catch class do else finally for if switch try while"),
  defKeywords: words("class interface enum @interface"),
  typeFirstDefinitions: true,
  atoms: words("true false null"),
  number: /^(?:0x[a-f\d_]+|0b[01_]+|(?:[\d_]+\.?\d*|\.\d+)(?:e[-+]?[\d_]+)?)(u|ll?|l|f)?/i,
  hooks: {
    "@": function(stream) {
      if (stream.match("interface", false))
        return false;
      stream.eatWhile(/[\w\$_]/);
      return "meta";
    }
  }
});
var csharp = clike({
  keywords: words("abstract as async await base break case catch checked class const continue default delegate do else enum event explicit extern finally fixed for foreach goto if implicit in interface internal is lock namespace new operator out override params private protected public readonly ref return sealed sizeof stackalloc static struct switch this throw try typeof unchecked unsafe using virtual void volatile while add alias ascending descending dynamic from get global group into join let orderby partial remove select set value var yield"),
  types: words("Action Boolean Byte Char DateTime DateTimeOffset Decimal Double Func Guid Int16 Int32 Int64 Object SByte Single String Task TimeSpan UInt16 UInt32 UInt64 bool byte char decimal double short int long object sbyte float string ushort uint ulong"),
  blockKeywords: words("catch class do else finally for foreach if struct switch try while"),
  defKeywords: words("class interface namespace struct var"),
  typeFirstDefinitions: true,
  atoms: words("true false null"),
  hooks: {
    "@": function(stream, state) {
      if (stream.eat('"')) {
        state.tokenize = tokenAtString;
        return tokenAtString(stream, state);
      }
      stream.eatWhile(/[\w\$_]/);
      return "meta";
    }
  }
});
function tokenTripleString(stream, state) {
  var escaped = false;
  while (!stream.eol()) {
    if (!escaped && stream.match('"""')) {
      state.tokenize = null;
      break;
    }
    escaped = stream.next() == "\\" && !escaped;
  }
  return "string";
}
function tokenNestedComment(depth) {
  return function(stream, state) {
    var ch;
    while (ch = stream.next()) {
      if (ch == "*" && stream.eat("/")) {
        if (depth == 1) {
          state.tokenize = null;
          break;
        } else {
          state.tokenize = tokenNestedComment(depth - 1);
          return state.tokenize(stream, state);
        }
      } else if (ch == "/" && stream.eat("*")) {
        state.tokenize = tokenNestedComment(depth + 1);
        return state.tokenize(stream, state);
      }
    }
    return "comment";
  };
}
var scala = clike({
  keywords: words("abstract case catch class def do else extends final finally for forSome if implicit import lazy match new null object override package private protected return sealed super this throw trait try type val var while with yield _ assert assume require print println printf readLine readBoolean readByte readShort readChar readInt readLong readFloat readDouble"),
  types: words("AnyVal App Application Array BufferedIterator BigDecimal BigInt Char Console Either Enumeration Equiv Error Exception Fractional Function IndexedSeq Int Integral Iterable Iterator List Map Numeric Nil NotNull Option Ordered Ordering PartialFunction PartialOrdering Product Proxy Range Responder Seq Serializable Set Specializable Stream StringBuilder StringContext Symbol Throwable Traversable TraversableOnce Tuple Unit Vector Boolean Byte Character CharSequence Class ClassLoader Cloneable Comparable Compiler Double Exception Float Integer Long Math Number Object Package Pair Process Runtime Runnable SecurityManager Short StackTraceElement StrictMath String StringBuffer System Thread ThreadGroup ThreadLocal Throwable Triple Void"),
  multiLineStrings: true,
  blockKeywords: words("catch class enum do else finally for forSome if match switch try while"),
  defKeywords: words("class enum def object package trait type val var"),
  atoms: words("true false null"),
  indentStatements: false,
  indentSwitch: false,
  isOperatorChar: /[+\-*&%=<>!?|\/#:@]/,
  hooks: {
    "@": function(stream) {
      stream.eatWhile(/[\w\$_]/);
      return "meta";
    },
    '"': function(stream, state) {
      if (!stream.match('""'))
        return false;
      state.tokenize = tokenTripleString;
      return state.tokenize(stream, state);
    },
    "'": function(stream) {
      stream.eatWhile(/[\w\$_\xa1-\uffff]/);
      return "atom";
    },
    "=": function(stream, state) {
      var cx = state.context;
      if (cx.type == "}" && cx.align && stream.eat(">")) {
        state.context = new Context(cx.indented, cx.column, cx.type, cx.info, null, cx.prev);
        return "operator";
      } else {
        return false;
      }
    },
    "/": function(stream, state) {
      if (!stream.eat("*"))
        return false;
      state.tokenize = tokenNestedComment(1);
      return state.tokenize(stream, state);
    }
  },
  languageData: {
    closeBrackets: {brackets: ["(", "[", "{", "'", '"', '"""']}
  }
});
function tokenKotlinString(tripleString) {
  return function(stream, state) {
    var escaped = false, next, end = false;
    while (!stream.eol()) {
      if (!tripleString && !escaped && stream.match('"')) {
        end = true;
        break;
      }
      if (tripleString && stream.match('"""')) {
        end = true;
        break;
      }
      next = stream.next();
      if (!escaped && next == "$" && stream.match("{"))
        stream.skipTo("}");
      escaped = !escaped && next == "\\" && !tripleString;
    }
    if (end || !tripleString)
      state.tokenize = null;
    return "string";
  };
}
var kotlin = clike({
  keywords: words("package as typealias class interface this super val operator var fun for is in This throw return annotation break continue object if else while do try when !in !is as? file import where by get set abstract enum open inner override private public internal protected catch finally out final vararg reified dynamic companion constructor init sealed field property receiver param sparam lateinit data inline noinline tailrec external annotation crossinline const operator infix suspend actual expect setparam"),
  types: words("Boolean Byte Character CharSequence Class ClassLoader Cloneable Comparable Compiler Double Exception Float Integer Long Math Number Object Package Pair Process Runtime Runnable SecurityManager Short StackTraceElement StrictMath String StringBuffer System Thread ThreadGroup ThreadLocal Throwable Triple Void Annotation Any BooleanArray ByteArray Char CharArray DeprecationLevel DoubleArray Enum FloatArray Function Int IntArray Lazy LazyThreadSafetyMode LongArray Nothing ShortArray Unit"),
  intendSwitch: false,
  indentStatements: false,
  multiLineStrings: true,
  number: /^(?:0x[a-f\d_]+|0b[01_]+|(?:[\d_]+(\.\d+)?|\.\d+)(?:e[-+]?[\d_]+)?)(u|ll?|l|f)?/i,
  blockKeywords: words("catch class do else finally for if where try while enum"),
  defKeywords: words("class val var object interface fun"),
  atoms: words("true false null this"),
  hooks: {
    "@": function(stream) {
      stream.eatWhile(/[\w\$_]/);
      return "meta";
    },
    "*": function(_stream, state) {
      return state.prevToken == "." ? "variable" : "operator";
    },
    '"': function(stream, state) {
      state.tokenize = tokenKotlinString(stream.match('""'));
      return state.tokenize(stream, state);
    },
    "/": function(stream, state) {
      if (!stream.eat("*"))
        return false;
      state.tokenize = tokenNestedComment(1);
      return state.tokenize(stream, state);
    },
    indent: function(state, ctx, textAfter, indentUnit2) {
      var firstChar = textAfter && textAfter.charAt(0);
      if ((state.prevToken == "}" || state.prevToken == ")") && textAfter == "")
        return state.indented;
      if (state.prevToken == "operator" && textAfter != "}" && state.context.type != "}" || state.prevToken == "variable" && firstChar == "." || (state.prevToken == "}" || state.prevToken == ")") && firstChar == ".")
        return indentUnit2 * 2 + ctx.indented;
      if (ctx.align && ctx.type == "}")
        return ctx.indented + (state.context.type == (textAfter || "").charAt(0) ? 0 : indentUnit2);
    }
  },
  languageData: {
    closeBrackets: {brackets: ["(", "[", "{", "'", '"', '"""']}
  }
});
var shader = clike({
  keywords: words("sampler1D sampler2D sampler3D samplerCube sampler1DShadow sampler2DShadow const attribute uniform varying break continue discard return for while do if else struct in out inout"),
  types: words("float int bool void vec2 vec3 vec4 ivec2 ivec3 ivec4 bvec2 bvec3 bvec4 mat2 mat3 mat4"),
  blockKeywords: words("for while do if else struct"),
  builtin: words("radians degrees sin cos tan asin acos atan pow exp log exp2 sqrt inversesqrt abs sign floor ceil fract mod min max clamp mix step smoothstep length distance dot cross normalize ftransform faceforward reflect refract matrixCompMult lessThan lessThanEqual greaterThan greaterThanEqual equal notEqual any all not texture1D texture1DProj texture1DLod texture1DProjLod texture2D texture2DProj texture2DLod texture2DProjLod texture3D texture3DProj texture3DLod texture3DProjLod textureCube textureCubeLod shadow1D shadow2D shadow1DProj shadow2DProj shadow1DLod shadow2DLod shadow1DProjLod shadow2DProjLod dFdx dFdy fwidth noise1 noise2 noise3 noise4"),
  atoms: words("true false gl_FragColor gl_SecondaryColor gl_Normal gl_Vertex gl_MultiTexCoord0 gl_MultiTexCoord1 gl_MultiTexCoord2 gl_MultiTexCoord3 gl_MultiTexCoord4 gl_MultiTexCoord5 gl_MultiTexCoord6 gl_MultiTexCoord7 gl_FogCoord gl_PointCoord gl_Position gl_PointSize gl_ClipVertex gl_FrontColor gl_BackColor gl_FrontSecondaryColor gl_BackSecondaryColor gl_TexCoord gl_FogFragCoord gl_FragCoord gl_FrontFacing gl_FragData gl_FragDepth gl_ModelViewMatrix gl_ProjectionMatrix gl_ModelViewProjectionMatrix gl_TextureMatrix gl_NormalMatrix gl_ModelViewMatrixInverse gl_ProjectionMatrixInverse gl_ModelViewProjectionMatrixInverse gl_TextureMatrixTranspose gl_ModelViewMatrixInverseTranspose gl_ProjectionMatrixInverseTranspose gl_ModelViewProjectionMatrixInverseTranspose gl_TextureMatrixInverseTranspose gl_NormalScale gl_DepthRange gl_ClipPlane gl_Point gl_FrontMaterial gl_BackMaterial gl_LightSource gl_LightModel gl_FrontLightModelProduct gl_BackLightModelProduct gl_TextureColor gl_EyePlaneS gl_EyePlaneT gl_EyePlaneR gl_EyePlaneQ gl_FogParameters gl_MaxLights gl_MaxClipPlanes gl_MaxTextureUnits gl_MaxTextureCoords gl_MaxVertexAttribs gl_MaxVertexUniformComponents gl_MaxVaryingFloats gl_MaxVertexTextureImageUnits gl_MaxTextureImageUnits gl_MaxFragmentUniformComponents gl_MaxCombineTextureImageUnits gl_MaxDrawBuffers"),
  indentSwitch: false,
  hooks: {"#": cppHook}
});
var nesC = clike({
  keywords: words(cKeywords + " as atomic async call command component components configuration event generic implementation includes interface module new norace nx_struct nx_union post provides signal task uses abstract extends"),
  types: cTypes,
  blockKeywords: words(cBlockKeywords),
  atoms: words("null true false"),
  hooks: {"#": cppHook}
});
var objectiveC = clike({
  keywords: words(cKeywords + " " + objCKeywords),
  types: objCTypes,
  builtin: words(objCBuiltins),
  blockKeywords: words(cBlockKeywords + " @synthesize @try @catch @finally @autoreleasepool @synchronized"),
  defKeywords: words(cDefKeywords + " @interface @implementation @protocol @class"),
  dontIndentStatements: /^@.*$/,
  typeFirstDefinitions: true,
  atoms: words("YES NO NULL Nil nil true false nullptr"),
  isReservedIdentifier: cIsReservedIdentifier,
  hooks: {
    "#": cppHook,
    "*": pointerHook
  }
});
var objectiveCpp = clike({
  keywords: words(cKeywords + " " + objCKeywords + " " + cppKeywords),
  types: objCTypes,
  builtin: words(objCBuiltins),
  blockKeywords: words(cBlockKeywords + " @synthesize @try @catch @finally @autoreleasepool @synchronized class try catch"),
  defKeywords: words(cDefKeywords + " @interface @implementation @protocol @class class namespace"),
  dontIndentStatements: /^@.*$|^template$/,
  typeFirstDefinitions: true,
  atoms: words("YES NO NULL Nil nil true false nullptr"),
  isReservedIdentifier: cIsReservedIdentifier,
  hooks: {
    "#": cppHook,
    "*": pointerHook,
    u: cpp11StringHook,
    U: cpp11StringHook,
    L: cpp11StringHook,
    R: cpp11StringHook,
    "0": cpp14Literal,
    "1": cpp14Literal,
    "2": cpp14Literal,
    "3": cpp14Literal,
    "4": cpp14Literal,
    "5": cpp14Literal,
    "6": cpp14Literal,
    "7": cpp14Literal,
    "8": cpp14Literal,
    "9": cpp14Literal,
    token: function(stream, state, style) {
      if (style == "variable" && stream.peek() == "(" && (state.prevToken == ";" || state.prevToken == null || state.prevToken == "}") && cppLooksLikeConstructor(stream.current()))
        return "def";
    }
  },
  namespaceSeparator: "::"
});
var squirrel = clike({
  keywords: words("base break clone continue const default delete enum extends function in class foreach local resume return this throw typeof yield constructor instanceof static"),
  types: cTypes,
  blockKeywords: words("case catch class else for foreach if switch try while"),
  defKeywords: words("function local class"),
  typeFirstDefinitions: true,
  atoms: words("true false null"),
  hooks: {"#": cppHook}
});
var stringTokenizer = null;
function tokenCeylonString(type2) {
  return function(stream, state) {
    var escaped = false, next, end = false;
    while (!stream.eol()) {
      if (!escaped && stream.match('"') && (type2 == "single" || stream.match('""'))) {
        end = true;
        break;
      }
      if (!escaped && stream.match("``")) {
        stringTokenizer = tokenCeylonString(type2);
        end = true;
        break;
      }
      next = stream.next();
      escaped = type2 == "single" && !escaped && next == "\\";
    }
    if (end)
      state.tokenize = null;
    return "string";
  };
}
var ceylon = clike({
  keywords: words("abstracts alias assembly assert assign break case catch class continue dynamic else exists extends finally for function given if import in interface is let module new nonempty object of out outer package return satisfies super switch then this throw try value void while"),
  types: function(word) {
    var first = word.charAt(0);
    return first === first.toUpperCase() && first !== first.toLowerCase();
  },
  blockKeywords: words("case catch class dynamic else finally for function if interface module new object switch try while"),
  defKeywords: words("class dynamic function interface module object package value"),
  builtin: words("abstract actual aliased annotation by default deprecated doc final formal late license native optional sealed see serializable shared suppressWarnings tagged throws variable"),
  isPunctuationChar: /[\[\]{}\(\),;\:\.`]/,
  isOperatorChar: /[+\-*&%=<>!?|^~:\/]/,
  numberStart: /[\d#$]/,
  number: /^(?:#[\da-fA-F_]+|\$[01_]+|[\d_]+[kMGTPmunpf]?|[\d_]+\.[\d_]+(?:[eE][-+]?\d+|[kMGTPmunpf]|)|)/i,
  multiLineStrings: true,
  typeFirstDefinitions: true,
  atoms: words("true false null larger smaller equal empty finished"),
  indentSwitch: false,
  styleDefs: false,
  hooks: {
    "@": function(stream) {
      stream.eatWhile(/[\w\$_]/);
      return "meta";
    },
    '"': function(stream, state) {
      state.tokenize = tokenCeylonString(stream.match('""') ? "triple" : "single");
      return state.tokenize(stream, state);
    },
    "`": function(stream, state) {
      if (!stringTokenizer || !stream.match("`"))
        return false;
      state.tokenize = stringTokenizer;
      stringTokenizer = null;
      return state.tokenize(stream, state);
    },
    "'": function(stream) {
      stream.eatWhile(/[\w\$_\xa1-\uffff]/);
      return "atom";
    },
    token: function(_stream, state, style) {
      if ((style == "variable" || style == "type") && state.prevToken == ".") {
        return "variableName.special";
      }
    }
  },
  languageData: {
    closeBrackets: {brackets: ["(", "[", "{", "'", '"', '"""']}
  }
});
function pushInterpolationStack(state) {
  (state.interpolationStack || (state.interpolationStack = [])).push(state.tokenize);
}
function popInterpolationStack(state) {
  return (state.interpolationStack || (state.interpolationStack = [])).pop();
}
function sizeInterpolationStack(state) {
  return state.interpolationStack ? state.interpolationStack.length : 0;
}
function tokenDartString(quote, stream, state, raw) {
  var tripleQuoted = false;
  if (stream.eat(quote)) {
    if (stream.eat(quote))
      tripleQuoted = true;
    else
      return "string";
  }
  function tokenStringHelper(stream2, state2) {
    var escaped = false;
    while (!stream2.eol()) {
      if (!raw && !escaped && stream2.peek() == "$") {
        pushInterpolationStack(state2);
        state2.tokenize = tokenInterpolation;
        return "string";
      }
      var next = stream2.next();
      if (next == quote && !escaped && (!tripleQuoted || stream2.match(quote + quote))) {
        state2.tokenize = null;
        break;
      }
      escaped = !raw && !escaped && next == "\\";
    }
    return "string";
  }
  state.tokenize = tokenStringHelper;
  return tokenStringHelper(stream, state);
}
function tokenInterpolation(stream, state) {
  stream.eat("$");
  if (stream.eat("{")) {
    state.tokenize = null;
  } else {
    state.tokenize = tokenInterpolationIdentifier;
  }
  return null;
}
function tokenInterpolationIdentifier(stream, state) {
  stream.eatWhile(/[\w_]/);
  state.tokenize = popInterpolationStack(state);
  return "variable";
}
var dart = clike({
  keywords: words("this super static final const abstract class extends external factory implements mixin get native set typedef with enum throw rethrow assert break case continue default in return new deferred async await covariant try catch finally do else for if switch while import library export part of show hide is as extension on yield late required"),
  blockKeywords: words("try catch finally do else for if switch while"),
  builtin: words("void bool num int double dynamic var String Null Never"),
  atoms: words("true false null"),
  hooks: {
    "@": function(stream) {
      stream.eatWhile(/[\w\$_\.]/);
      return "meta";
    },
    "'": function(stream, state) {
      return tokenDartString("'", stream, state, false);
    },
    '"': function(stream, state) {
      return tokenDartString('"', stream, state, false);
    },
    r: function(stream, state) {
      var peek = stream.peek();
      if (peek == "'" || peek == '"') {
        return tokenDartString(stream.next(), stream, state, true);
      }
      return false;
    },
    "}": function(_stream, state) {
      if (sizeInterpolationStack(state) > 0) {
        state.tokenize = popInterpolationStack(state);
        return null;
      }
      return false;
    },
    "/": function(stream, state) {
      if (!stream.eat("*"))
        return false;
      state.tokenize = tokenNestedComment(1);
      return state.tokenize(stream, state);
    },
    token: function(stream, _, style) {
      if (style == "variable") {
        var isUpper = RegExp("^[_$]*[A-Z][a-zA-Z0-9_$]*$", "g");
        if (isUpper.test(stream.current())) {
          return "type";
        }
      }
    }
  }
});

// node_modules/@codemirror/legacy-modes/mode/cobol.js
var BUILTIN = "builtin";
var COMMENT = "comment";
var STRING = "string";
var ATOM = "atom";
var NUMBER = "number";
var KEYWORD = "keyword";
var MODTAG = "header";
var COBOLLINENUM = "def";
var PERIOD = "link";
function makeKeywords(str) {
  var obj = {}, words4 = str.split(" ");
  for (var i = 0; i < words4.length; ++i)
    obj[words4[i]] = true;
  return obj;
}
var atoms = makeKeywords("TRUE FALSE ZEROES ZEROS ZERO SPACES SPACE LOW-VALUE LOW-VALUES ");
var keywords = makeKeywords("ACCEPT ACCESS ACQUIRE ADD ADDRESS ADVANCING AFTER ALIAS ALL ALPHABET ALPHABETIC ALPHABETIC-LOWER ALPHABETIC-UPPER ALPHANUMERIC ALPHANUMERIC-EDITED ALSO ALTER ALTERNATE AND ANY ARE AREA AREAS ARITHMETIC ASCENDING ASSIGN AT ATTRIBUTE AUTHOR AUTO AUTO-SKIP AUTOMATIC B-AND B-EXOR B-LESS B-NOT B-OR BACKGROUND-COLOR BACKGROUND-COLOUR BEEP BEFORE BELL BINARY BIT BITS BLANK BLINK BLOCK BOOLEAN BOTTOM BY CALL CANCEL CD CF CH CHARACTER CHARACTERS CLASS CLOCK-UNITS CLOSE COBOL CODE CODE-SET COL COLLATING COLUMN COMMA COMMIT COMMITMENT COMMON COMMUNICATION COMP COMP-0 COMP-1 COMP-2 COMP-3 COMP-4 COMP-5 COMP-6 COMP-7 COMP-8 COMP-9 COMPUTATIONAL COMPUTATIONAL-0 COMPUTATIONAL-1 COMPUTATIONAL-2 COMPUTATIONAL-3 COMPUTATIONAL-4 COMPUTATIONAL-5 COMPUTATIONAL-6 COMPUTATIONAL-7 COMPUTATIONAL-8 COMPUTATIONAL-9 COMPUTE CONFIGURATION CONNECT CONSOLE CONTAINED CONTAINS CONTENT CONTINUE CONTROL CONTROL-AREA CONTROLS CONVERTING COPY CORR CORRESPONDING COUNT CRT CRT-UNDER CURRENCY CURRENT CURSOR DATA DATE DATE-COMPILED DATE-WRITTEN DAY DAY-OF-WEEK DB DB-ACCESS-CONTROL-KEY DB-DATA-NAME DB-EXCEPTION DB-FORMAT-NAME DB-RECORD-NAME DB-SET-NAME DB-STATUS DBCS DBCS-EDITED DE DEBUG-CONTENTS DEBUG-ITEM DEBUG-LINE DEBUG-NAME DEBUG-SUB-1 DEBUG-SUB-2 DEBUG-SUB-3 DEBUGGING DECIMAL-POINT DECLARATIVES DEFAULT DELETE DELIMITED DELIMITER DEPENDING DESCENDING DESCRIBED DESTINATION DETAIL DISABLE DISCONNECT DISPLAY DISPLAY-1 DISPLAY-2 DISPLAY-3 DISPLAY-4 DISPLAY-5 DISPLAY-6 DISPLAY-7 DISPLAY-8 DISPLAY-9 DIVIDE DIVISION DOWN DROP DUPLICATE DUPLICATES DYNAMIC EBCDIC EGI EJECT ELSE EMI EMPTY EMPTY-CHECK ENABLE END END. END-ACCEPT END-ACCEPT. END-ADD END-CALL END-COMPUTE END-DELETE END-DISPLAY END-DIVIDE END-EVALUATE END-IF END-INVOKE END-MULTIPLY END-OF-PAGE END-PERFORM END-READ END-RECEIVE END-RETURN END-REWRITE END-SEARCH END-START END-STRING END-SUBTRACT END-UNSTRING END-WRITE END-XML ENTER ENTRY ENVIRONMENT EOP EQUAL EQUALS ERASE ERROR ESI EVALUATE EVERY EXCEEDS EXCEPTION EXCLUSIVE EXIT EXTEND EXTERNAL EXTERNALLY-DESCRIBED-KEY FD FETCH FILE FILE-CONTROL FILE-STREAM FILES FILLER FINAL FIND FINISH FIRST FOOTING FOR FOREGROUND-COLOR FOREGROUND-COLOUR FORMAT FREE FROM FULL FUNCTION GENERATE GET GIVING GLOBAL GO GOBACK GREATER GROUP HEADING HIGH-VALUE HIGH-VALUES HIGHLIGHT I-O I-O-CONTROL ID IDENTIFICATION IF IN INDEX INDEX-1 INDEX-2 INDEX-3 INDEX-4 INDEX-5 INDEX-6 INDEX-7 INDEX-8 INDEX-9 INDEXED INDIC INDICATE INDICATOR INDICATORS INITIAL INITIALIZE INITIATE INPUT INPUT-OUTPUT INSPECT INSTALLATION INTO INVALID INVOKE IS JUST JUSTIFIED KANJI KEEP KEY LABEL LAST LD LEADING LEFT LEFT-JUSTIFY LENGTH LENGTH-CHECK LESS LIBRARY LIKE LIMIT LIMITS LINAGE LINAGE-COUNTER LINE LINE-COUNTER LINES LINKAGE LOCAL-STORAGE LOCALE LOCALLY LOCK MEMBER MEMORY MERGE MESSAGE METACLASS MODE MODIFIED MODIFY MODULES MOVE MULTIPLE MULTIPLY NATIONAL NATIVE NEGATIVE NEXT NO NO-ECHO NONE NOT NULL NULL-KEY-MAP NULL-MAP NULLS NUMBER NUMERIC NUMERIC-EDITED OBJECT OBJECT-COMPUTER OCCURS OF OFF OMITTED ON ONLY OPEN OPTIONAL OR ORDER ORGANIZATION OTHER OUTPUT OVERFLOW OWNER PACKED-DECIMAL PADDING PAGE PAGE-COUNTER PARSE PERFORM PF PH PIC PICTURE PLUS POINTER POSITION POSITIVE PREFIX PRESENT PRINTING PRIOR PROCEDURE PROCEDURE-POINTER PROCEDURES PROCEED PROCESS PROCESSING PROGRAM PROGRAM-ID PROMPT PROTECTED PURGE QUEUE QUOTE QUOTES RANDOM RD READ READY REALM RECEIVE RECONNECT RECORD RECORD-NAME RECORDS RECURSIVE REDEFINES REEL REFERENCE REFERENCE-MONITOR REFERENCES RELATION RELATIVE RELEASE REMAINDER REMOVAL RENAMES REPEATED REPLACE REPLACING REPORT REPORTING REPORTS REPOSITORY REQUIRED RERUN RESERVE RESET RETAINING RETRIEVAL RETURN RETURN-CODE RETURNING REVERSE-VIDEO REVERSED REWIND REWRITE RF RH RIGHT RIGHT-JUSTIFY ROLLBACK ROLLING ROUNDED RUN SAME SCREEN SD SEARCH SECTION SECURE SECURITY SEGMENT SEGMENT-LIMIT SELECT SEND SENTENCE SEPARATE SEQUENCE SEQUENTIAL SET SHARED SIGN SIZE SKIP1 SKIP2 SKIP3 SORT SORT-MERGE SORT-RETURN SOURCE SOURCE-COMPUTER SPACE-FILL SPECIAL-NAMES STANDARD STANDARD-1 STANDARD-2 START STARTING STATUS STOP STORE STRING SUB-QUEUE-1 SUB-QUEUE-2 SUB-QUEUE-3 SUB-SCHEMA SUBFILE SUBSTITUTE SUBTRACT SUM SUPPRESS SYMBOLIC SYNC SYNCHRONIZED SYSIN SYSOUT TABLE TALLYING TAPE TENANT TERMINAL TERMINATE TEST TEXT THAN THEN THROUGH THRU TIME TIMES TITLE TO TOP TRAILING TRAILING-SIGN TRANSACTION TYPE TYPEDEF UNDERLINE UNEQUAL UNIT UNSTRING UNTIL UP UPDATE UPON USAGE USAGE-MODE USE USING VALID VALIDATE VALUE VALUES VARYING VLR WAIT WHEN WHEN-COMPILED WITH WITHIN WORDS WORKING-STORAGE WRITE XML XML-CODE XML-EVENT XML-NTEXT XML-TEXT ZERO ZERO-FILL ");
var builtins = makeKeywords("- * ** / + < <= = > >= ");
var tests = {
  digit: /\d/,
  digit_or_colon: /[\d:]/,
  hex: /[0-9a-f]/i,
  sign: /[+-]/,
  exponent: /e/i,
  keyword_char: /[^\s\(\[\;\)\]]/,
  symbol: /[\w*+\-]/
};
function isNumber(ch, stream) {
  if (ch === "0" && stream.eat(/x/i)) {
    stream.eatWhile(tests.hex);
    return true;
  }
  if ((ch == "+" || ch == "-") && tests.digit.test(stream.peek())) {
    stream.eat(tests.sign);
    ch = stream.next();
  }
  if (tests.digit.test(ch)) {
    stream.eat(ch);
    stream.eatWhile(tests.digit);
    if (stream.peek() == ".") {
      stream.eat(".");
      stream.eatWhile(tests.digit);
    }
    if (stream.eat(tests.exponent)) {
      stream.eat(tests.sign);
      stream.eatWhile(tests.digit);
    }
    return true;
  }
  return false;
}
var cobol = {
  startState: function() {
    return {
      indentStack: null,
      indentation: 0,
      mode: false
    };
  },
  token: function(stream, state) {
    if (state.indentStack == null && stream.sol()) {
      state.indentation = 6;
    }
    if (stream.eatSpace()) {
      return null;
    }
    var returnType = null;
    switch (state.mode) {
      case "string":
        var next = false;
        while ((next = stream.next()) != null) {
          if (next == '"' || next == "'") {
            state.mode = false;
            break;
          }
        }
        returnType = STRING;
        break;
      default:
        var ch = stream.next();
        var col = stream.column();
        if (col >= 0 && col <= 5) {
          returnType = COBOLLINENUM;
        } else if (col >= 72 && col <= 79) {
          stream.skipToEnd();
          returnType = MODTAG;
        } else if (ch == "*" && col == 6) {
          stream.skipToEnd();
          returnType = COMMENT;
        } else if (ch == '"' || ch == "'") {
          state.mode = "string";
          returnType = STRING;
        } else if (ch == "'" && !tests.digit_or_colon.test(stream.peek())) {
          returnType = ATOM;
        } else if (ch == ".") {
          returnType = PERIOD;
        } else if (isNumber(ch, stream)) {
          returnType = NUMBER;
        } else {
          if (stream.current().match(tests.symbol)) {
            while (col < 71) {
              if (stream.eat(tests.symbol) === void 0) {
                break;
              } else {
                col++;
              }
            }
          }
          if (keywords && keywords.propertyIsEnumerable(stream.current().toUpperCase())) {
            returnType = KEYWORD;
          } else if (builtins && builtins.propertyIsEnumerable(stream.current().toUpperCase())) {
            returnType = BUILTIN;
          } else if (atoms && atoms.propertyIsEnumerable(stream.current().toUpperCase())) {
            returnType = ATOM;
          } else
            returnType = null;
        }
    }
    return returnType;
  },
  indent: function(state) {
    if (state.indentStack == null)
      return state.indentation;
    return state.indentStack.indent;
  }
};

// node_modules/@codemirror/legacy-modes/mode/commonlisp.js
var specialForm = /^(block|let*|return-from|catch|load-time-value|setq|eval-when|locally|symbol-macrolet|flet|macrolet|tagbody|function|multiple-value-call|the|go|multiple-value-prog1|throw|if|progn|unwind-protect|labels|progv|let|quote)$/;
var assumeBody = /^with|^def|^do|^prog|case$|^cond$|bind$|when$|unless$/;
var numLiteral = /^(?:[+\-]?(?:\d+|\d*\.\d+)(?:[efd][+\-]?\d+)?|[+\-]?\d+(?:\/[+\-]?\d+)?|#b[+\-]?[01]+|#o[+\-]?[0-7]+|#x[+\-]?[\da-f]+)/;
var symbol = /[^\s'`,@()\[\]";]/;
var type;
function readSym(stream) {
  var ch;
  while (ch = stream.next()) {
    if (ch == "\\")
      stream.next();
    else if (!symbol.test(ch)) {
      stream.backUp(1);
      break;
    }
  }
  return stream.current();
}
function base2(stream, state) {
  if (stream.eatSpace()) {
    type = "ws";
    return null;
  }
  if (stream.match(numLiteral))
    return "number";
  var ch = stream.next();
  if (ch == "\\")
    ch = stream.next();
  if (ch == '"')
    return (state.tokenize = inString)(stream, state);
  else if (ch == "(") {
    type = "open";
    return "bracket";
  } else if (ch == ")" || ch == "]") {
    type = "close";
    return "bracket";
  } else if (ch == ";") {
    stream.skipToEnd();
    type = "ws";
    return "comment";
  } else if (/['`,@]/.test(ch))
    return null;
  else if (ch == "|") {
    if (stream.skipTo("|")) {
      stream.next();
      return "variableName";
    } else {
      stream.skipToEnd();
      return "error";
    }
  } else if (ch == "#") {
    var ch = stream.next();
    if (ch == "(") {
      type = "open";
      return "bracket";
    } else if (/[+\-=\.']/.test(ch))
      return null;
    else if (/\d/.test(ch) && stream.match(/^\d*#/))
      return null;
    else if (ch == "|")
      return (state.tokenize = inComment)(stream, state);
    else if (ch == ":") {
      readSym(stream);
      return "meta";
    } else if (ch == "\\") {
      stream.next();
      readSym(stream);
      return "string.special";
    } else
      return "error";
  } else {
    var name2 = readSym(stream);
    if (name2 == ".")
      return null;
    type = "symbol";
    if (name2 == "nil" || name2 == "t" || name2.charAt(0) == ":")
      return "atom";
    if (state.lastType == "open" && (specialForm.test(name2) || assumeBody.test(name2)))
      return "keyword";
    if (name2.charAt(0) == "&")
      return "variableName.special";
    return "variableName";
  }
}
function inString(stream, state) {
  var escaped = false, next;
  while (next = stream.next()) {
    if (next == '"' && !escaped) {
      state.tokenize = base2;
      break;
    }
    escaped = !escaped && next == "\\";
  }
  return "string";
}
function inComment(stream, state) {
  var next, last;
  while (next = stream.next()) {
    if (next == "#" && last == "|") {
      state.tokenize = base2;
      break;
    }
    last = next;
  }
  type = "ws";
  return "comment";
}
var commonLisp = {
  startState: function() {
    return {ctx: {prev: null, start: 0, indentTo: 0}, lastType: null, tokenize: base2};
  },
  token: function(stream, state) {
    if (stream.sol() && typeof state.ctx.indentTo != "number")
      state.ctx.indentTo = state.ctx.start + 1;
    type = null;
    var style = state.tokenize(stream, state);
    if (type != "ws") {
      if (state.ctx.indentTo == null) {
        if (type == "symbol" && assumeBody.test(stream.current()))
          state.ctx.indentTo = state.ctx.start + stream.indentUnit;
        else
          state.ctx.indentTo = "next";
      } else if (state.ctx.indentTo == "next") {
        state.ctx.indentTo = stream.column();
      }
      state.lastType = type;
    }
    if (type == "open")
      state.ctx = {prev: state.ctx, start: stream.column(), indentTo: null};
    else if (type == "close")
      state.ctx = state.ctx.prev || state.ctx;
    return style;
  },
  indent: function(state) {
    var i = state.ctx.indentTo;
    return typeof i == "number" ? i : state.ctx.start + 1;
  },
  languageData: {
    commentTokens: {line: ";;", block: {open: "#|", close: "|#"}},
    closeBrackets: {brackets: ["(", "[", "{", '"']}
  }
};

// node_modules/@codemirror/legacy-modes/mode/crystal.js
function wordRegExp(words4, end) {
  return new RegExp((end ? "" : "^") + "(?:" + words4.join("|") + ")" + (end ? "$" : "\\b"));
}
function chain(tokenize2, stream, state) {
  state.tokenize.push(tokenize2);
  return tokenize2(stream, state);
}
var operators = /^(?:[-+/%|&^]|\*\*?|[<>]{2})/;
var conditionalOperators = /^(?:[=!]~|===|<=>|[<>=!]=?|[|&]{2}|~)/;
var indexingOperators = /^(?:\[\][?=]?)/;
var anotherOperators = /^(?:\.(?:\.{2})?|->|[?:])/;
var idents = /^[a-z_\u009F-\uFFFF][a-zA-Z0-9_\u009F-\uFFFF]*/;
var types2 = /^[A-Z_\u009F-\uFFFF][a-zA-Z0-9_\u009F-\uFFFF]*/;
var keywords2 = wordRegExp([
  "abstract",
  "alias",
  "as",
  "asm",
  "begin",
  "break",
  "case",
  "class",
  "def",
  "do",
  "else",
  "elsif",
  "end",
  "ensure",
  "enum",
  "extend",
  "for",
  "fun",
  "if",
  "include",
  "instance_sizeof",
  "lib",
  "macro",
  "module",
  "next",
  "of",
  "out",
  "pointerof",
  "private",
  "protected",
  "rescue",
  "return",
  "require",
  "select",
  "sizeof",
  "struct",
  "super",
  "then",
  "type",
  "typeof",
  "uninitialized",
  "union",
  "unless",
  "until",
  "when",
  "while",
  "with",
  "yield",
  "__DIR__",
  "__END_LINE__",
  "__FILE__",
  "__LINE__"
]);
var atomWords = wordRegExp(["true", "false", "nil", "self"]);
var indentKeywordsArray = [
  "def",
  "fun",
  "macro",
  "class",
  "module",
  "struct",
  "lib",
  "enum",
  "union",
  "do",
  "for"
];
var indentKeywords = wordRegExp(indentKeywordsArray);
var indentExpressionKeywordsArray = ["if", "unless", "case", "while", "until", "begin", "then"];
var indentExpressionKeywords = wordRegExp(indentExpressionKeywordsArray);
var dedentKeywordsArray = ["end", "else", "elsif", "rescue", "ensure"];
var dedentKeywords = wordRegExp(dedentKeywordsArray);
var dedentPunctualsArray = ["\\)", "\\}", "\\]"];
var dedentPunctuals = new RegExp("^(?:" + dedentPunctualsArray.join("|") + ")$");
var nextTokenizer = {
  def: tokenFollowIdent,
  fun: tokenFollowIdent,
  macro: tokenMacroDef,
  class: tokenFollowType,
  module: tokenFollowType,
  struct: tokenFollowType,
  lib: tokenFollowType,
  enum: tokenFollowType,
  union: tokenFollowType
};
var matching = {"[": "]", "{": "}", "(": ")", "<": ">"};
function tokenBase(stream, state) {
  if (stream.eatSpace()) {
    return null;
  }
  if (state.lastToken != "\\" && stream.match("{%", false)) {
    return chain(tokenMacro("%", "%"), stream, state);
  }
  if (state.lastToken != "\\" && stream.match("{{", false)) {
    return chain(tokenMacro("{", "}"), stream, state);
  }
  if (stream.peek() == "#") {
    stream.skipToEnd();
    return "comment";
  }
  var matched;
  if (stream.match(idents)) {
    stream.eat(/[?!]/);
    matched = stream.current();
    if (stream.eat(":")) {
      return "atom";
    } else if (state.lastToken == ".") {
      return "property";
    } else if (keywords2.test(matched)) {
      if (indentKeywords.test(matched)) {
        if (!(matched == "fun" && state.blocks.indexOf("lib") >= 0) && !(matched == "def" && state.lastToken == "abstract")) {
          state.blocks.push(matched);
          state.currentIndent += 1;
        }
      } else if ((state.lastStyle == "operator" || !state.lastStyle) && indentExpressionKeywords.test(matched)) {
        state.blocks.push(matched);
        state.currentIndent += 1;
      } else if (matched == "end") {
        state.blocks.pop();
        state.currentIndent -= 1;
      }
      if (nextTokenizer.hasOwnProperty(matched)) {
        state.tokenize.push(nextTokenizer[matched]);
      }
      return "keyword";
    } else if (atomWords.test(matched)) {
      return "atom";
    }
    return "variable";
  }
  if (stream.eat("@")) {
    if (stream.peek() == "[") {
      return chain(tokenNest("[", "]", "meta"), stream, state);
    }
    stream.eat("@");
    stream.match(idents) || stream.match(types2);
    return "propertyName";
  }
  if (stream.match(types2)) {
    return "tag";
  }
  if (stream.eat(":")) {
    if (stream.eat('"')) {
      return chain(tokenQuote('"', "atom", false), stream, state);
    } else if (stream.match(idents) || stream.match(types2) || stream.match(operators) || stream.match(conditionalOperators) || stream.match(indexingOperators)) {
      return "atom";
    }
    stream.eat(":");
    return "operator";
  }
  if (stream.eat('"')) {
    return chain(tokenQuote('"', "string", true), stream, state);
  }
  if (stream.peek() == "%") {
    var style = "string";
    var embed = true;
    var delim;
    if (stream.match("%r")) {
      style = "string.special";
      delim = stream.next();
    } else if (stream.match("%w")) {
      embed = false;
      delim = stream.next();
    } else if (stream.match("%q")) {
      embed = false;
      delim = stream.next();
    } else {
      if (delim = stream.match(/^%([^\w\s=])/)) {
        delim = delim[1];
      } else if (stream.match(/^%[a-zA-Z0-9_\u009F-\uFFFF]*/)) {
        return "meta";
      } else {
        return "operator";
      }
    }
    if (matching.hasOwnProperty(delim)) {
      delim = matching[delim];
    }
    return chain(tokenQuote(delim, style, embed), stream, state);
  }
  if (matched = stream.match(/^<<-('?)([A-Z]\w*)\1/)) {
    return chain(tokenHereDoc(matched[2], !matched[1]), stream, state);
  }
  if (stream.eat("'")) {
    stream.match(/^(?:[^']|\\(?:[befnrtv0'"]|[0-7]{3}|u(?:[0-9a-fA-F]{4}|\{[0-9a-fA-F]{1,6}\})))/);
    stream.eat("'");
    return "atom";
  }
  if (stream.eat("0")) {
    if (stream.eat("x")) {
      stream.match(/^[0-9a-fA-F]+/);
    } else if (stream.eat("o")) {
      stream.match(/^[0-7]+/);
    } else if (stream.eat("b")) {
      stream.match(/^[01]+/);
    }
    return "number";
  }
  if (stream.eat(/^\d/)) {
    stream.match(/^\d*(?:\.\d+)?(?:[eE][+-]?\d+)?/);
    return "number";
  }
  if (stream.match(operators)) {
    stream.eat("=");
    return "operator";
  }
  if (stream.match(conditionalOperators) || stream.match(anotherOperators)) {
    return "operator";
  }
  if (matched = stream.match(/[({[]/, false)) {
    matched = matched[0];
    return chain(tokenNest(matched, matching[matched], null), stream, state);
  }
  if (stream.eat("\\")) {
    stream.next();
    return "meta";
  }
  stream.next();
  return null;
}
function tokenNest(begin, end, style, started) {
  return function(stream, state) {
    if (!started && stream.match(begin)) {
      state.tokenize[state.tokenize.length - 1] = tokenNest(begin, end, style, true);
      state.currentIndent += 1;
      return style;
    }
    var nextStyle = tokenBase(stream, state);
    if (stream.current() === end) {
      state.tokenize.pop();
      state.currentIndent -= 1;
      nextStyle = style;
    }
    return nextStyle;
  };
}
function tokenMacro(begin, end, started) {
  return function(stream, state) {
    if (!started && stream.match("{" + begin)) {
      state.currentIndent += 1;
      state.tokenize[state.tokenize.length - 1] = tokenMacro(begin, end, true);
      return "meta";
    }
    if (stream.match(end + "}")) {
      state.currentIndent -= 1;
      state.tokenize.pop();
      return "meta";
    }
    return tokenBase(stream, state);
  };
}
function tokenMacroDef(stream, state) {
  if (stream.eatSpace()) {
    return null;
  }
  var matched;
  if (matched = stream.match(idents)) {
    if (matched == "def") {
      return "keyword";
    }
    stream.eat(/[?!]/);
  }
  state.tokenize.pop();
  return "def";
}
function tokenFollowIdent(stream, state) {
  if (stream.eatSpace()) {
    return null;
  }
  if (stream.match(idents)) {
    stream.eat(/[!?]/);
  } else {
    stream.match(operators) || stream.match(conditionalOperators) || stream.match(indexingOperators);
  }
  state.tokenize.pop();
  return "def";
}
function tokenFollowType(stream, state) {
  if (stream.eatSpace()) {
    return null;
  }
  stream.match(types2);
  state.tokenize.pop();
  return "def";
}
function tokenQuote(end, style, embed) {
  return function(stream, state) {
    var escaped = false;
    while (stream.peek()) {
      if (!escaped) {
        if (stream.match("{%", false)) {
          state.tokenize.push(tokenMacro("%", "%"));
          return style;
        }
        if (stream.match("{{", false)) {
          state.tokenize.push(tokenMacro("{", "}"));
          return style;
        }
        if (embed && stream.match("#{", false)) {
          state.tokenize.push(tokenNest("#{", "}", "meta"));
          return style;
        }
        var ch = stream.next();
        if (ch == end) {
          state.tokenize.pop();
          return style;
        }
        escaped = embed && ch == "\\";
      } else {
        stream.next();
        escaped = false;
      }
    }
    return style;
  };
}
function tokenHereDoc(phrase, embed) {
  return function(stream, state) {
    if (stream.sol()) {
      stream.eatSpace();
      if (stream.match(phrase)) {
        state.tokenize.pop();
        return "string";
      }
    }
    var escaped = false;
    while (stream.peek()) {
      if (!escaped) {
        if (stream.match("{%", false)) {
          state.tokenize.push(tokenMacro("%", "%"));
          return "string";
        }
        if (stream.match("{{", false)) {
          state.tokenize.push(tokenMacro("{", "}"));
          return "string";
        }
        if (embed && stream.match("#{", false)) {
          state.tokenize.push(tokenNest("#{", "}", "meta"));
          return "string";
        }
        escaped = embed && stream.next() == "\\";
      } else {
        stream.next();
        escaped = false;
      }
    }
    return "string";
  };
}
var crystal = {
  startState: function() {
    return {
      tokenize: [tokenBase],
      currentIndent: 0,
      lastToken: null,
      lastStyle: null,
      blocks: []
    };
  },
  token: function(stream, state) {
    var style = state.tokenize[state.tokenize.length - 1](stream, state);
    var token = stream.current();
    if (style && style != "comment") {
      state.lastToken = token;
      state.lastStyle = style;
    }
    return style;
  },
  indent: function(state, textAfter, cx) {
    textAfter = textAfter.replace(/^\s*(?:\{%)?\s*|\s*(?:%\})?\s*$/g, "");
    if (dedentKeywords.test(textAfter) || dedentPunctuals.test(textAfter)) {
      return cx.unit * (state.currentIndent - 1);
    }
    return cx.unit * state.currentIndent;
  },
  languageData: {
    indentOnInput: wordRegExp(dedentPunctualsArray.concat(dedentKeywordsArray), true),
    commentTokens: {line: "#"}
  }
};

// node_modules/@codemirror/legacy-modes/mode/fortran.js
function words2(array) {
  var keys = {};
  for (var i = 0; i < array.length; ++i) {
    keys[array[i]] = true;
  }
  return keys;
}
var keywords3 = words2([
  "abstract",
  "accept",
  "allocatable",
  "allocate",
  "array",
  "assign",
  "asynchronous",
  "backspace",
  "bind",
  "block",
  "byte",
  "call",
  "case",
  "class",
  "close",
  "common",
  "contains",
  "continue",
  "cycle",
  "data",
  "deallocate",
  "decode",
  "deferred",
  "dimension",
  "do",
  "elemental",
  "else",
  "encode",
  "end",
  "endif",
  "entry",
  "enumerator",
  "equivalence",
  "exit",
  "external",
  "extrinsic",
  "final",
  "forall",
  "format",
  "function",
  "generic",
  "go",
  "goto",
  "if",
  "implicit",
  "import",
  "include",
  "inquire",
  "intent",
  "interface",
  "intrinsic",
  "module",
  "namelist",
  "non_intrinsic",
  "non_overridable",
  "none",
  "nopass",
  "nullify",
  "open",
  "optional",
  "options",
  "parameter",
  "pass",
  "pause",
  "pointer",
  "print",
  "private",
  "program",
  "protected",
  "public",
  "pure",
  "read",
  "recursive",
  "result",
  "return",
  "rewind",
  "save",
  "select",
  "sequence",
  "stop",
  "subroutine",
  "target",
  "then",
  "to",
  "type",
  "use",
  "value",
  "volatile",
  "where",
  "while",
  "write"
]);
var builtins2 = words2([
  "abort",
  "abs",
  "access",
  "achar",
  "acos",
  "adjustl",
  "adjustr",
  "aimag",
  "aint",
  "alarm",
  "all",
  "allocated",
  "alog",
  "amax",
  "amin",
  "amod",
  "and",
  "anint",
  "any",
  "asin",
  "associated",
  "atan",
  "besj",
  "besjn",
  "besy",
  "besyn",
  "bit_size",
  "btest",
  "cabs",
  "ccos",
  "ceiling",
  "cexp",
  "char",
  "chdir",
  "chmod",
  "clog",
  "cmplx",
  "command_argument_count",
  "complex",
  "conjg",
  "cos",
  "cosh",
  "count",
  "cpu_time",
  "cshift",
  "csin",
  "csqrt",
  "ctime",
  "c_funloc",
  "c_loc",
  "c_associated",
  "c_null_ptr",
  "c_null_funptr",
  "c_f_pointer",
  "c_null_char",
  "c_alert",
  "c_backspace",
  "c_form_feed",
  "c_new_line",
  "c_carriage_return",
  "c_horizontal_tab",
  "c_vertical_tab",
  "dabs",
  "dacos",
  "dasin",
  "datan",
  "date_and_time",
  "dbesj",
  "dbesj",
  "dbesjn",
  "dbesy",
  "dbesy",
  "dbesyn",
  "dble",
  "dcos",
  "dcosh",
  "ddim",
  "derf",
  "derfc",
  "dexp",
  "digits",
  "dim",
  "dint",
  "dlog",
  "dlog",
  "dmax",
  "dmin",
  "dmod",
  "dnint",
  "dot_product",
  "dprod",
  "dsign",
  "dsinh",
  "dsin",
  "dsqrt",
  "dtanh",
  "dtan",
  "dtime",
  "eoshift",
  "epsilon",
  "erf",
  "erfc",
  "etime",
  "exit",
  "exp",
  "exponent",
  "extends_type_of",
  "fdate",
  "fget",
  "fgetc",
  "float",
  "floor",
  "flush",
  "fnum",
  "fputc",
  "fput",
  "fraction",
  "fseek",
  "fstat",
  "ftell",
  "gerror",
  "getarg",
  "get_command",
  "get_command_argument",
  "get_environment_variable",
  "getcwd",
  "getenv",
  "getgid",
  "getlog",
  "getpid",
  "getuid",
  "gmtime",
  "hostnm",
  "huge",
  "iabs",
  "iachar",
  "iand",
  "iargc",
  "ibclr",
  "ibits",
  "ibset",
  "ichar",
  "idate",
  "idim",
  "idint",
  "idnint",
  "ieor",
  "ierrno",
  "ifix",
  "imag",
  "imagpart",
  "index",
  "int",
  "ior",
  "irand",
  "isatty",
  "ishft",
  "ishftc",
  "isign",
  "iso_c_binding",
  "is_iostat_end",
  "is_iostat_eor",
  "itime",
  "kill",
  "kind",
  "lbound",
  "len",
  "len_trim",
  "lge",
  "lgt",
  "link",
  "lle",
  "llt",
  "lnblnk",
  "loc",
  "log",
  "logical",
  "long",
  "lshift",
  "lstat",
  "ltime",
  "matmul",
  "max",
  "maxexponent",
  "maxloc",
  "maxval",
  "mclock",
  "merge",
  "move_alloc",
  "min",
  "minexponent",
  "minloc",
  "minval",
  "mod",
  "modulo",
  "mvbits",
  "nearest",
  "new_line",
  "nint",
  "not",
  "or",
  "pack",
  "perror",
  "precision",
  "present",
  "product",
  "radix",
  "rand",
  "random_number",
  "random_seed",
  "range",
  "real",
  "realpart",
  "rename",
  "repeat",
  "reshape",
  "rrspacing",
  "rshift",
  "same_type_as",
  "scale",
  "scan",
  "second",
  "selected_int_kind",
  "selected_real_kind",
  "set_exponent",
  "shape",
  "short",
  "sign",
  "signal",
  "sinh",
  "sin",
  "sleep",
  "sngl",
  "spacing",
  "spread",
  "sqrt",
  "srand",
  "stat",
  "sum",
  "symlnk",
  "system",
  "system_clock",
  "tan",
  "tanh",
  "time",
  "tiny",
  "transfer",
  "transpose",
  "trim",
  "ttynam",
  "ubound",
  "umask",
  "unlink",
  "unpack",
  "verify",
  "xor",
  "zabs",
  "zcos",
  "zexp",
  "zlog",
  "zsin",
  "zsqrt"
]);
var dataTypes = words2([
  "c_bool",
  "c_char",
  "c_double",
  "c_double_complex",
  "c_float",
  "c_float_complex",
  "c_funptr",
  "c_int",
  "c_int16_t",
  "c_int32_t",
  "c_int64_t",
  "c_int8_t",
  "c_int_fast16_t",
  "c_int_fast32_t",
  "c_int_fast64_t",
  "c_int_fast8_t",
  "c_int_least16_t",
  "c_int_least32_t",
  "c_int_least64_t",
  "c_int_least8_t",
  "c_intmax_t",
  "c_intptr_t",
  "c_long",
  "c_long_double",
  "c_long_double_complex",
  "c_long_long",
  "c_ptr",
  "c_short",
  "c_signed_char",
  "c_size_t",
  "character",
  "complex",
  "double",
  "integer",
  "logical",
  "real"
]);
var isOperatorChar = /[+\-*&=<>\/\:]/;
var litOperator = new RegExp("(.and.|.or.|.eq.|.lt.|.le.|.gt.|.ge.|.ne.|.not.|.eqv.|.neqv.)", "i");
function tokenBase2(stream, state) {
  if (stream.match(litOperator)) {
    return "operator";
  }
  var ch = stream.next();
  if (ch == "!") {
    stream.skipToEnd();
    return "comment";
  }
  if (ch == '"' || ch == "'") {
    state.tokenize = tokenString(ch);
    return state.tokenize(stream, state);
  }
  if (/[\[\]\(\),]/.test(ch)) {
    return null;
  }
  if (/\d/.test(ch)) {
    stream.eatWhile(/[\w\.]/);
    return "number";
  }
  if (isOperatorChar.test(ch)) {
    stream.eatWhile(isOperatorChar);
    return "operator";
  }
  stream.eatWhile(/[\w\$_]/);
  var word = stream.current().toLowerCase();
  if (keywords3.hasOwnProperty(word)) {
    return "keyword";
  }
  if (builtins2.hasOwnProperty(word) || dataTypes.hasOwnProperty(word)) {
    return "builtin";
  }
  return "variable";
}
function tokenString(quote) {
  return function(stream, state) {
    var escaped = false, next, end = false;
    while ((next = stream.next()) != null) {
      if (next == quote && !escaped) {
        end = true;
        break;
      }
      escaped = !escaped && next == "\\";
    }
    if (end || !escaped)
      state.tokenize = null;
    return "string";
  };
}
var fortran = {
  startState: function() {
    return {tokenize: null};
  },
  token: function(stream, state) {
    if (stream.eatSpace())
      return null;
    var style = (state.tokenize || tokenBase2)(stream, state);
    if (style == "comment" || style == "meta")
      return style;
    return style;
  }
};

// node_modules/@codemirror/legacy-modes/mode/mllike.js
function mlLike(parserConfig) {
  var words4 = {
    as: "keyword",
    do: "keyword",
    else: "keyword",
    end: "keyword",
    exception: "keyword",
    fun: "keyword",
    functor: "keyword",
    if: "keyword",
    in: "keyword",
    include: "keyword",
    let: "keyword",
    of: "keyword",
    open: "keyword",
    rec: "keyword",
    struct: "keyword",
    then: "keyword",
    type: "keyword",
    val: "keyword",
    while: "keyword",
    with: "keyword"
  };
  var extraWords = parserConfig.extraWords || {};
  for (var prop in extraWords) {
    if (extraWords.hasOwnProperty(prop)) {
      words4[prop] = parserConfig.extraWords[prop];
    }
  }
  var hintWords = [];
  for (var k in words4) {
    hintWords.push(k);
  }
  function tokenBase9(stream, state) {
    var ch = stream.next();
    if (ch === '"') {
      state.tokenize = tokenString5;
      return state.tokenize(stream, state);
    }
    if (ch === "{") {
      if (stream.eat("|")) {
        state.longString = true;
        state.tokenize = tokenLongString;
        return state.tokenize(stream, state);
      }
    }
    if (ch === "(") {
      if (stream.eat("*")) {
        state.commentLevel++;
        state.tokenize = tokenComment5;
        return state.tokenize(stream, state);
      }
    }
    if (ch === "~" || ch === "?") {
      stream.eatWhile(/\w/);
      return "variableName.special";
    }
    if (ch === "`") {
      stream.eatWhile(/\w/);
      return "quote";
    }
    if (ch === "/" && parserConfig.slashComments && stream.eat("/")) {
      stream.skipToEnd();
      return "comment";
    }
    if (/\d/.test(ch)) {
      if (ch === "0" && stream.eat(/[bB]/)) {
        stream.eatWhile(/[01]/);
      }
      if (ch === "0" && stream.eat(/[xX]/)) {
        stream.eatWhile(/[0-9a-fA-F]/);
      }
      if (ch === "0" && stream.eat(/[oO]/)) {
        stream.eatWhile(/[0-7]/);
      } else {
        stream.eatWhile(/[\d_]/);
        if (stream.eat(".")) {
          stream.eatWhile(/[\d]/);
        }
        if (stream.eat(/[eE]/)) {
          stream.eatWhile(/[\d\-+]/);
        }
      }
      return "number";
    }
    if (/[+\-*&%=<>!?|@\.~:]/.test(ch)) {
      return "operator";
    }
    if (/[\w\xa1-\uffff]/.test(ch)) {
      stream.eatWhile(/[\w\xa1-\uffff]/);
      var cur2 = stream.current();
      return words4.hasOwnProperty(cur2) ? words4[cur2] : "variable";
    }
    return null;
  }
  function tokenString5(stream, state) {
    var next, end = false, escaped = false;
    while ((next = stream.next()) != null) {
      if (next === '"' && !escaped) {
        end = true;
        break;
      }
      escaped = !escaped && next === "\\";
    }
    if (end && !escaped) {
      state.tokenize = tokenBase9;
    }
    return "string";
  }
  ;
  function tokenComment5(stream, state) {
    var prev, next;
    while (state.commentLevel > 0 && (next = stream.next()) != null) {
      if (prev === "(" && next === "*")
        state.commentLevel++;
      if (prev === "*" && next === ")")
        state.commentLevel--;
      prev = next;
    }
    if (state.commentLevel <= 0) {
      state.tokenize = tokenBase9;
    }
    return "comment";
  }
  function tokenLongString(stream, state) {
    var prev, next;
    while (state.longString && (next = stream.next()) != null) {
      if (prev === "|" && next === "}")
        state.longString = false;
      prev = next;
    }
    if (!state.longString) {
      state.tokenize = tokenBase9;
    }
    return "string";
  }
  return {
    startState: function() {
      return {tokenize: tokenBase9, commentLevel: 0, longString: false};
    },
    token: function(stream, state) {
      if (stream.eatSpace())
        return null;
      return state.tokenize(stream, state);
    },
    languageData: {
      autocomplete: hintWords,
      commentTokens: {
        line: parserConfig.slashComments ? "//" : void 0,
        block: {open: "(*", close: "*)"}
      }
    }
  };
}
var oCaml = mlLike({
  extraWords: {
    and: "keyword",
    assert: "keyword",
    begin: "keyword",
    class: "keyword",
    constraint: "keyword",
    done: "keyword",
    downto: "keyword",
    external: "keyword",
    function: "keyword",
    initializer: "keyword",
    lazy: "keyword",
    match: "keyword",
    method: "keyword",
    module: "keyword",
    mutable: "keyword",
    new: "keyword",
    nonrec: "keyword",
    object: "keyword",
    private: "keyword",
    sig: "keyword",
    to: "keyword",
    try: "keyword",
    value: "keyword",
    virtual: "keyword",
    when: "keyword",
    raise: "builtin",
    failwith: "builtin",
    true: "builtin",
    false: "builtin",
    asr: "builtin",
    land: "builtin",
    lor: "builtin",
    lsl: "builtin",
    lsr: "builtin",
    lxor: "builtin",
    mod: "builtin",
    or: "builtin",
    raise_notrace: "builtin",
    trace: "builtin",
    exit: "builtin",
    print_string: "builtin",
    print_endline: "builtin",
    int: "type",
    float: "type",
    bool: "type",
    char: "type",
    string: "type",
    unit: "type",
    List: "builtin"
  }
});
var fSharp = mlLike({
  extraWords: {
    abstract: "keyword",
    assert: "keyword",
    base: "keyword",
    begin: "keyword",
    class: "keyword",
    default: "keyword",
    delegate: "keyword",
    "do!": "keyword",
    done: "keyword",
    downcast: "keyword",
    downto: "keyword",
    elif: "keyword",
    extern: "keyword",
    finally: "keyword",
    for: "keyword",
    function: "keyword",
    global: "keyword",
    inherit: "keyword",
    inline: "keyword",
    interface: "keyword",
    internal: "keyword",
    lazy: "keyword",
    "let!": "keyword",
    match: "keyword",
    member: "keyword",
    module: "keyword",
    mutable: "keyword",
    namespace: "keyword",
    new: "keyword",
    null: "keyword",
    override: "keyword",
    private: "keyword",
    public: "keyword",
    "return!": "keyword",
    return: "keyword",
    select: "keyword",
    static: "keyword",
    to: "keyword",
    try: "keyword",
    upcast: "keyword",
    "use!": "keyword",
    use: "keyword",
    void: "keyword",
    when: "keyword",
    "yield!": "keyword",
    yield: "keyword",
    atomic: "keyword",
    break: "keyword",
    checked: "keyword",
    component: "keyword",
    const: "keyword",
    constraint: "keyword",
    constructor: "keyword",
    continue: "keyword",
    eager: "keyword",
    event: "keyword",
    external: "keyword",
    fixed: "keyword",
    method: "keyword",
    mixin: "keyword",
    object: "keyword",
    parallel: "keyword",
    process: "keyword",
    protected: "keyword",
    pure: "keyword",
    sealed: "keyword",
    tailcall: "keyword",
    trait: "keyword",
    virtual: "keyword",
    volatile: "keyword",
    List: "builtin",
    Seq: "builtin",
    Map: "builtin",
    Set: "builtin",
    Option: "builtin",
    int: "builtin",
    string: "builtin",
    not: "builtin",
    true: "builtin",
    false: "builtin",
    raise: "builtin",
    failwith: "builtin"
  },
  slashComments: true
});
var sml = mlLike({
  extraWords: {
    abstype: "keyword",
    and: "keyword",
    andalso: "keyword",
    case: "keyword",
    datatype: "keyword",
    fn: "keyword",
    handle: "keyword",
    infix: "keyword",
    infixr: "keyword",
    local: "keyword",
    nonfix: "keyword",
    op: "keyword",
    orelse: "keyword",
    raise: "keyword",
    withtype: "keyword",
    eqtype: "keyword",
    sharing: "keyword",
    sig: "keyword",
    signature: "keyword",
    structure: "keyword",
    where: "keyword",
    true: "keyword",
    false: "keyword",
    int: "builtin",
    real: "builtin",
    string: "builtin",
    char: "builtin",
    bool: "builtin"
  },
  slashComments: true
});

// node_modules/@codemirror/legacy-modes/mode/go.js
var keywords4 = {
  break: true,
  case: true,
  chan: true,
  const: true,
  continue: true,
  default: true,
  defer: true,
  else: true,
  fallthrough: true,
  for: true,
  func: true,
  go: true,
  goto: true,
  if: true,
  import: true,
  interface: true,
  map: true,
  package: true,
  range: true,
  return: true,
  select: true,
  struct: true,
  switch: true,
  type: true,
  var: true,
  bool: true,
  byte: true,
  complex64: true,
  complex128: true,
  float32: true,
  float64: true,
  int8: true,
  int16: true,
  int32: true,
  int64: true,
  string: true,
  uint8: true,
  uint16: true,
  uint32: true,
  uint64: true,
  int: true,
  uint: true,
  uintptr: true,
  error: true,
  rune: true
};
var atoms2 = {
  true: true,
  false: true,
  iota: true,
  nil: true,
  append: true,
  cap: true,
  close: true,
  complex: true,
  copy: true,
  delete: true,
  imag: true,
  len: true,
  make: true,
  new: true,
  panic: true,
  print: true,
  println: true,
  real: true,
  recover: true
};
var isOperatorChar2 = /[+\-*&^%:=<>!|\/]/;
var curPunc;
function tokenBase3(stream, state) {
  var ch = stream.next();
  if (ch == '"' || ch == "'" || ch == "`") {
    state.tokenize = tokenString2(ch);
    return state.tokenize(stream, state);
  }
  if (/[\d\.]/.test(ch)) {
    if (ch == ".") {
      stream.match(/^[0-9]+([eE][\-+]?[0-9]+)?/);
    } else if (ch == "0") {
      stream.match(/^[xX][0-9a-fA-F]+/) || stream.match(/^0[0-7]+/);
    } else {
      stream.match(/^[0-9]*\.?[0-9]*([eE][\-+]?[0-9]+)?/);
    }
    return "number";
  }
  if (/[\[\]{}\(\),;\:\.]/.test(ch)) {
    curPunc = ch;
    return null;
  }
  if (ch == "/") {
    if (stream.eat("*")) {
      state.tokenize = tokenComment;
      return tokenComment(stream, state);
    }
    if (stream.eat("/")) {
      stream.skipToEnd();
      return "comment";
    }
  }
  if (isOperatorChar2.test(ch)) {
    stream.eatWhile(isOperatorChar2);
    return "operator";
  }
  stream.eatWhile(/[\w\$_\xa1-\uffff]/);
  var cur2 = stream.current();
  if (keywords4.propertyIsEnumerable(cur2)) {
    if (cur2 == "case" || cur2 == "default")
      curPunc = "case";
    return "keyword";
  }
  if (atoms2.propertyIsEnumerable(cur2))
    return "atom";
  return "variable";
}
function tokenString2(quote) {
  return function(stream, state) {
    var escaped = false, next, end = false;
    while ((next = stream.next()) != null) {
      if (next == quote && !escaped) {
        end = true;
        break;
      }
      escaped = !escaped && quote != "`" && next == "\\";
    }
    if (end || !(escaped || quote == "`"))
      state.tokenize = tokenBase3;
    return "string";
  };
}
function tokenComment(stream, state) {
  var maybeEnd = false, ch;
  while (ch = stream.next()) {
    if (ch == "/" && maybeEnd) {
      state.tokenize = tokenBase3;
      break;
    }
    maybeEnd = ch == "*";
  }
  return "comment";
}
function Context2(indented, column, type2, align, prev) {
  this.indented = indented;
  this.column = column;
  this.type = type2;
  this.align = align;
  this.prev = prev;
}
function pushContext2(state, col, type2) {
  return state.context = new Context2(state.indented, col, type2, null, state.context);
}
function popContext2(state) {
  if (!state.context.prev)
    return;
  var t2 = state.context.type;
  if (t2 == ")" || t2 == "]" || t2 == "}")
    state.indented = state.context.indented;
  return state.context = state.context.prev;
}
var go = {
  startState: function(indentUnit2) {
    return {
      tokenize: null,
      context: new Context2(-indentUnit2, 0, "top", false),
      indented: 0,
      startOfLine: true
    };
  },
  token: function(stream, state) {
    var ctx = state.context;
    if (stream.sol()) {
      if (ctx.align == null)
        ctx.align = false;
      state.indented = stream.indentation();
      state.startOfLine = true;
      if (ctx.type == "case")
        ctx.type = "}";
    }
    if (stream.eatSpace())
      return null;
    curPunc = null;
    var style = (state.tokenize || tokenBase3)(stream, state);
    if (style == "comment")
      return style;
    if (ctx.align == null)
      ctx.align = true;
    if (curPunc == "{")
      pushContext2(state, stream.column(), "}");
    else if (curPunc == "[")
      pushContext2(state, stream.column(), "]");
    else if (curPunc == "(")
      pushContext2(state, stream.column(), ")");
    else if (curPunc == "case")
      ctx.type = "case";
    else if (curPunc == "}" && ctx.type == "}")
      popContext2(state);
    else if (curPunc == ctx.type)
      popContext2(state);
    state.startOfLine = false;
    return style;
  },
  indent: function(state, textAfter, cx) {
    if (state.tokenize != tokenBase3 && state.tokenize != null)
      return null;
    var ctx = state.context, firstChar = textAfter && textAfter.charAt(0);
    if (ctx.type == "case" && /^(?:case|default)\b/.test(textAfter)) {
      state.context.type = "}";
      return ctx.indented;
    }
    var closing3 = firstChar == ctx.type;
    if (ctx.align)
      return ctx.column + (closing3 ? 0 : 1);
    else
      return ctx.indented + (closing3 ? 0 : cx.unit);
  },
  languageData: {
    indentOnInput: /^\s([{}]|case |default\s*:)$/,
    commentTokens: {line: "//", block: {open: "/*", close: "*/"}}
  }
};

// node_modules/@codemirror/legacy-modes/mode/haskell.js
function switchState(source, setState, f) {
  setState(f);
  return f(source, setState);
}
var smallRE = /[a-z_]/;
var largeRE = /[A-Z]/;
var digitRE = /\d/;
var hexitRE = /[0-9A-Fa-f]/;
var octitRE = /[0-7]/;
var idRE = /[a-z_A-Z0-9'\xa1-\uffff]/;
var symbolRE = /[-!#$%&*+.\/<=>?@\\^|~:]/;
var specialRE = /[(),;[\]`{}]/;
var whiteCharRE = /[ \t\v\f]/;
function normal(source, setState) {
  if (source.eatWhile(whiteCharRE)) {
    return null;
  }
  var ch = source.next();
  if (specialRE.test(ch)) {
    if (ch == "{" && source.eat("-")) {
      var t2 = "comment";
      if (source.eat("#")) {
        t2 = "meta";
      }
      return switchState(source, setState, ncomment(t2, 1));
    }
    return null;
  }
  if (ch == "'") {
    if (source.eat("\\")) {
      source.next();
    } else {
      source.next();
    }
    if (source.eat("'")) {
      return "string";
    }
    return "error";
  }
  if (ch == '"') {
    return switchState(source, setState, stringLiteral);
  }
  if (largeRE.test(ch)) {
    source.eatWhile(idRE);
    if (source.eat(".")) {
      return "qualifier";
    }
    return "type";
  }
  if (smallRE.test(ch)) {
    source.eatWhile(idRE);
    return "variable";
  }
  if (digitRE.test(ch)) {
    if (ch == "0") {
      if (source.eat(/[xX]/)) {
        source.eatWhile(hexitRE);
        return "integer";
      }
      if (source.eat(/[oO]/)) {
        source.eatWhile(octitRE);
        return "number";
      }
    }
    source.eatWhile(digitRE);
    var t2 = "number";
    if (source.match(/^\.\d+/)) {
      t2 = "number";
    }
    if (source.eat(/[eE]/)) {
      t2 = "number";
      source.eat(/[-+]/);
      source.eatWhile(digitRE);
    }
    return t2;
  }
  if (ch == "." && source.eat("."))
    return "keyword";
  if (symbolRE.test(ch)) {
    if (ch == "-" && source.eat(/-/)) {
      source.eatWhile(/-/);
      if (!source.eat(symbolRE)) {
        source.skipToEnd();
        return "comment";
      }
    }
    source.eatWhile(symbolRE);
    return "variable";
  }
  return "error";
}
function ncomment(type2, nest) {
  if (nest == 0) {
    return normal;
  }
  return function(source, setState) {
    var currNest = nest;
    while (!source.eol()) {
      var ch = source.next();
      if (ch == "{" && source.eat("-")) {
        ++currNest;
      } else if (ch == "-" && source.eat("}")) {
        --currNest;
        if (currNest == 0) {
          setState(normal);
          return type2;
        }
      }
    }
    setState(ncomment(type2, currNest));
    return type2;
  };
}
function stringLiteral(source, setState) {
  while (!source.eol()) {
    var ch = source.next();
    if (ch == '"') {
      setState(normal);
      return "string";
    }
    if (ch == "\\") {
      if (source.eol() || source.eat(whiteCharRE)) {
        setState(stringGap);
        return "string";
      }
      if (source.eat("&")) {
      } else {
        source.next();
      }
    }
  }
  setState(normal);
  return "error";
}
function stringGap(source, setState) {
  if (source.eat("\\")) {
    return switchState(source, setState, stringLiteral);
  }
  source.next();
  setState(normal);
  return "error";
}
var wellKnownWords = function() {
  var wkw = {};
  function setType(t2) {
    return function() {
      for (var i = 0; i < arguments.length; i++)
        wkw[arguments[i]] = t2;
    };
  }
  setType("keyword")("case", "class", "data", "default", "deriving", "do", "else", "foreign", "if", "import", "in", "infix", "infixl", "infixr", "instance", "let", "module", "newtype", "of", "then", "type", "where", "_");
  setType("keyword")("..", ":", "::", "=", "\\", "<-", "->", "@", "~", "=>");
  setType("builtin")("!!", "$!", "$", "&&", "+", "++", "-", ".", "/", "/=", "<", "<*", "<=", "<$>", "<*>", "=<<", "==", ">", ">=", ">>", ">>=", "^", "^^", "||", "*", "*>", "**");
  setType("builtin")("Applicative", "Bool", "Bounded", "Char", "Double", "EQ", "Either", "Enum", "Eq", "False", "FilePath", "Float", "Floating", "Fractional", "Functor", "GT", "IO", "IOError", "Int", "Integer", "Integral", "Just", "LT", "Left", "Maybe", "Monad", "Nothing", "Num", "Ord", "Ordering", "Rational", "Read", "ReadS", "Real", "RealFloat", "RealFrac", "Right", "Show", "ShowS", "String", "True");
  setType("builtin")("abs", "acos", "acosh", "all", "and", "any", "appendFile", "asTypeOf", "asin", "asinh", "atan", "atan2", "atanh", "break", "catch", "ceiling", "compare", "concat", "concatMap", "const", "cos", "cosh", "curry", "cycle", "decodeFloat", "div", "divMod", "drop", "dropWhile", "either", "elem", "encodeFloat", "enumFrom", "enumFromThen", "enumFromThenTo", "enumFromTo", "error", "even", "exp", "exponent", "fail", "filter", "flip", "floatDigits", "floatRadix", "floatRange", "floor", "fmap", "foldl", "foldl1", "foldr", "foldr1", "fromEnum", "fromInteger", "fromIntegral", "fromRational", "fst", "gcd", "getChar", "getContents", "getLine", "head", "id", "init", "interact", "ioError", "isDenormalized", "isIEEE", "isInfinite", "isNaN", "isNegativeZero", "iterate", "last", "lcm", "length", "lex", "lines", "log", "logBase", "lookup", "map", "mapM", "mapM_", "max", "maxBound", "maximum", "maybe", "min", "minBound", "minimum", "mod", "negate", "not", "notElem", "null", "odd", "or", "otherwise", "pi", "pred", "print", "product", "properFraction", "pure", "putChar", "putStr", "putStrLn", "quot", "quotRem", "read", "readFile", "readIO", "readList", "readLn", "readParen", "reads", "readsPrec", "realToFrac", "recip", "rem", "repeat", "replicate", "return", "reverse", "round", "scaleFloat", "scanl", "scanl1", "scanr", "scanr1", "seq", "sequence", "sequence_", "show", "showChar", "showList", "showParen", "showString", "shows", "showsPrec", "significand", "signum", "sin", "sinh", "snd", "span", "splitAt", "sqrt", "subtract", "succ", "sum", "tail", "take", "takeWhile", "tan", "tanh", "toEnum", "toInteger", "toRational", "truncate", "uncurry", "undefined", "unlines", "until", "unwords", "unzip", "unzip3", "userError", "words", "writeFile", "zip", "zip3", "zipWith", "zipWith3");
  return wkw;
}();
var haskell = {
  startState: function() {
    return {f: normal};
  },
  copyState: function(s) {
    return {f: s.f};
  },
  token: function(stream, state) {
    var t2 = state.f(stream, function(s) {
      state.f = s;
    });
    var w = stream.current();
    return wellKnownWords.hasOwnProperty(w) ? wellKnownWords[w] : t2;
  },
  languageData: {
    commentTokens: {line: "--", block: {open: "{-", close: "-}"}}
  }
};

// node_modules/lezer/dist/index.es.js
var Stack = class {
  constructor(p, stack, state, reducePos, pos, score2, buffer, bufferBase, curContext, parent) {
    this.p = p;
    this.stack = stack;
    this.state = state;
    this.reducePos = reducePos;
    this.pos = pos;
    this.score = score2;
    this.buffer = buffer;
    this.bufferBase = bufferBase;
    this.curContext = curContext;
    this.parent = parent;
  }
  toString() {
    return `[${this.stack.filter((_, i) => i % 3 == 0).concat(this.state)}]@${this.pos}${this.score ? "!" + this.score : ""}`;
  }
  static start(p, state, pos = 0) {
    let cx = p.parser.context;
    return new Stack(p, [], state, pos, pos, 0, [], 0, cx ? new StackContext(cx, cx.start) : null, null);
  }
  get context() {
    return this.curContext ? this.curContext.context : null;
  }
  pushState(state, start) {
    this.stack.push(this.state, start, this.bufferBase + this.buffer.length);
    this.state = state;
  }
  reduce(action) {
    let depth = action >> 19, type2 = action & 65535;
    let {parser: parser6} = this.p;
    let dPrec = parser6.dynamicPrecedence(type2);
    if (dPrec)
      this.score += dPrec;
    if (depth == 0) {
      if (type2 < parser6.minRepeatTerm)
        this.storeNode(type2, this.reducePos, this.reducePos, 4, true);
      this.pushState(parser6.getGoto(this.state, type2, true), this.reducePos);
      this.reduceContext(type2);
      return;
    }
    let base3 = this.stack.length - (depth - 1) * 3 - (action & 262144 ? 6 : 0);
    let start = this.stack[base3 - 2];
    let bufferBase = this.stack[base3 - 1], count = this.bufferBase + this.buffer.length - bufferBase;
    if (type2 < parser6.minRepeatTerm || action & 131072) {
      let pos = parser6.stateFlag(this.state, 1) ? this.pos : this.reducePos;
      this.storeNode(type2, start, pos, count + 4, true);
    }
    if (action & 262144) {
      this.state = this.stack[base3];
    } else {
      let baseStateID = this.stack[base3 - 3];
      this.state = parser6.getGoto(baseStateID, type2, true);
    }
    while (this.stack.length > base3)
      this.stack.pop();
    this.reduceContext(type2);
  }
  storeNode(term, start, end, size = 4, isReduce = false) {
    if (term == 0) {
      let cur2 = this, top2 = this.buffer.length;
      if (top2 == 0 && cur2.parent) {
        top2 = cur2.bufferBase - cur2.parent.bufferBase;
        cur2 = cur2.parent;
      }
      if (top2 > 0 && cur2.buffer[top2 - 4] == 0 && cur2.buffer[top2 - 1] > -1) {
        if (start == end)
          return;
        if (cur2.buffer[top2 - 2] >= start) {
          cur2.buffer[top2 - 2] = end;
          return;
        }
      }
    }
    if (!isReduce || this.pos == end) {
      this.buffer.push(term, start, end, size);
    } else {
      let index = this.buffer.length;
      if (index > 0 && this.buffer[index - 4] != 0)
        while (index > 0 && this.buffer[index - 2] > end) {
          this.buffer[index] = this.buffer[index - 4];
          this.buffer[index + 1] = this.buffer[index - 3];
          this.buffer[index + 2] = this.buffer[index - 2];
          this.buffer[index + 3] = this.buffer[index - 1];
          index -= 4;
          if (size > 4)
            size -= 4;
        }
      this.buffer[index] = term;
      this.buffer[index + 1] = start;
      this.buffer[index + 2] = end;
      this.buffer[index + 3] = size;
    }
  }
  shift(action, next, nextEnd) {
    if (action & 131072) {
      this.pushState(action & 65535, this.pos);
    } else if ((action & 262144) == 0) {
      let start = this.pos, nextState = action, {parser: parser6} = this.p;
      if (nextEnd > this.pos || next <= parser6.maxNode) {
        this.pos = nextEnd;
        if (!parser6.stateFlag(nextState, 1))
          this.reducePos = nextEnd;
      }
      this.pushState(nextState, start);
      if (next <= parser6.maxNode)
        this.buffer.push(next, start, nextEnd, 4);
      this.shiftContext(next);
    } else {
      if (next <= this.p.parser.maxNode)
        this.buffer.push(next, this.pos, nextEnd, 4);
      this.pos = nextEnd;
    }
  }
  apply(action, next, nextEnd) {
    if (action & 65536)
      this.reduce(action);
    else
      this.shift(action, next, nextEnd);
  }
  useNode(value, next) {
    let index = this.p.reused.length - 1;
    if (index < 0 || this.p.reused[index] != value) {
      this.p.reused.push(value);
      index++;
    }
    let start = this.pos;
    this.reducePos = this.pos = start + value.length;
    this.pushState(next, start);
    this.buffer.push(index, start, this.reducePos, -1);
    if (this.curContext)
      this.updateContext(this.curContext.tracker.reuse(this.curContext.context, value, this.p.input, this));
  }
  split() {
    let parent = this;
    let off = parent.buffer.length;
    while (off > 0 && parent.buffer[off - 2] > parent.reducePos)
      off -= 4;
    let buffer = parent.buffer.slice(off), base3 = parent.bufferBase + off;
    while (parent && base3 == parent.bufferBase)
      parent = parent.parent;
    return new Stack(this.p, this.stack.slice(), this.state, this.reducePos, this.pos, this.score, buffer, base3, this.curContext, parent);
  }
  recoverByDelete(next, nextEnd) {
    let isNode = next <= this.p.parser.maxNode;
    if (isNode)
      this.storeNode(next, this.pos, nextEnd);
    this.storeNode(0, this.pos, nextEnd, isNode ? 8 : 4);
    this.pos = this.reducePos = nextEnd;
    this.score -= 200;
  }
  canShift(term) {
    for (let sim = new SimulatedStack(this); ; ) {
      let action = this.p.parser.stateSlot(sim.top, 4) || this.p.parser.hasAction(sim.top, term);
      if ((action & 65536) == 0)
        return true;
      if (action == 0)
        return false;
      sim.reduce(action);
    }
  }
  get ruleStart() {
    for (let state = this.state, base3 = this.stack.length; ; ) {
      let force = this.p.parser.stateSlot(state, 5);
      if (!(force & 65536))
        return 0;
      base3 -= 3 * (force >> 19);
      if ((force & 65535) < this.p.parser.minRepeatTerm)
        return this.stack[base3 + 1];
      state = this.stack[base3];
    }
  }
  startOf(types4, before) {
    let state = this.state, frame = this.stack.length, {parser: parser6} = this.p;
    for (; ; ) {
      let force = parser6.stateSlot(state, 5);
      let depth = force >> 19, term = force & 65535;
      if (types4.indexOf(term) > -1) {
        let base3 = frame - 3 * (force >> 19), pos = this.stack[base3 + 1];
        if (before == null || before > pos)
          return pos;
      }
      if (frame == 0)
        return null;
      if (depth == 0) {
        frame -= 3;
        state = this.stack[frame];
      } else {
        frame -= 3 * (depth - 1);
        state = parser6.getGoto(this.stack[frame - 3], term, true);
      }
    }
  }
  recoverByInsert(next) {
    if (this.stack.length >= 300)
      return [];
    let nextStates = this.p.parser.nextStates(this.state);
    if (nextStates.length > 4 << 1 || this.stack.length >= 120) {
      let best = [];
      for (let i = 0, s; i < nextStates.length; i += 2) {
        if ((s = nextStates[i + 1]) != this.state && this.p.parser.hasAction(s, next))
          best.push(nextStates[i], s);
      }
      if (this.stack.length < 120)
        for (let i = 0; best.length < 4 << 1 && i < nextStates.length; i += 2) {
          let s = nextStates[i + 1];
          if (!best.some((v, i2) => i2 & 1 && v == s))
            best.push(nextStates[i], s);
        }
      nextStates = best;
    }
    let result = [];
    for (let i = 0; i < nextStates.length && result.length < 4; i += 2) {
      let s = nextStates[i + 1];
      if (s == this.state)
        continue;
      let stack = this.split();
      stack.storeNode(0, stack.pos, stack.pos, 4, true);
      stack.pushState(s, this.pos);
      stack.shiftContext(nextStates[i]);
      stack.score -= 200;
      result.push(stack);
    }
    return result;
  }
  forceReduce() {
    let reduce = this.p.parser.stateSlot(this.state, 5);
    if ((reduce & 65536) == 0)
      return false;
    if (!this.p.parser.validAction(this.state, reduce)) {
      this.storeNode(0, this.reducePos, this.reducePos, 4, true);
      this.score -= 100;
    }
    this.reduce(reduce);
    return true;
  }
  forceAll() {
    while (!this.p.parser.stateFlag(this.state, 2) && this.forceReduce()) {
    }
    return this;
  }
  get deadEnd() {
    if (this.stack.length != 3)
      return false;
    let {parser: parser6} = this.p;
    return parser6.data[parser6.stateSlot(this.state, 1)] == 65535 && !parser6.stateSlot(this.state, 4);
  }
  restart() {
    this.state = this.stack[0];
    this.stack.length = 0;
  }
  sameState(other) {
    if (this.state != other.state || this.stack.length != other.stack.length)
      return false;
    for (let i = 0; i < this.stack.length; i += 3)
      if (this.stack[i] != other.stack[i])
        return false;
    return true;
  }
  get parser() {
    return this.p.parser;
  }
  dialectEnabled(dialectID) {
    return this.p.parser.dialect.flags[dialectID];
  }
  shiftContext(term) {
    if (this.curContext)
      this.updateContext(this.curContext.tracker.shift(this.curContext.context, term, this.p.input, this));
  }
  reduceContext(term) {
    if (this.curContext)
      this.updateContext(this.curContext.tracker.reduce(this.curContext.context, term, this.p.input, this));
  }
  emitContext() {
    let cx = this.curContext;
    if (!cx.tracker.strict)
      return;
    let last = this.buffer.length - 1;
    if (last < 0 || this.buffer[last] != -2)
      this.buffer.push(cx.hash, this.reducePos, this.reducePos, -2);
  }
  updateContext(context) {
    if (context != this.curContext.context) {
      let newCx = new StackContext(this.curContext.tracker, context);
      if (newCx.hash != this.curContext.hash)
        this.emitContext();
      this.curContext = newCx;
    }
  }
};
var StackContext = class {
  constructor(tracker, context) {
    this.tracker = tracker;
    this.context = context;
    this.hash = tracker.hash(context);
  }
};
var Recover;
(function(Recover2) {
  Recover2[Recover2["Token"] = 200] = "Token";
  Recover2[Recover2["Reduce"] = 100] = "Reduce";
  Recover2[Recover2["MaxNext"] = 4] = "MaxNext";
  Recover2[Recover2["MaxInsertStackDepth"] = 300] = "MaxInsertStackDepth";
  Recover2[Recover2["DampenInsertStackDepth"] = 120] = "DampenInsertStackDepth";
})(Recover || (Recover = {}));
var SimulatedStack = class {
  constructor(stack) {
    this.stack = stack;
    this.top = stack.state;
    this.rest = stack.stack;
    this.offset = this.rest.length;
  }
  reduce(action) {
    let term = action & 65535, depth = action >> 19;
    if (depth == 0) {
      if (this.rest == this.stack.stack)
        this.rest = this.rest.slice();
      this.rest.push(this.top, 0, 0);
      this.offset += 3;
    } else {
      this.offset -= (depth - 1) * 3;
    }
    let goto = this.stack.p.parser.getGoto(this.rest[this.offset - 3], term, true);
    this.top = goto;
  }
};
var StackBufferCursor = class {
  constructor(stack, pos, index) {
    this.stack = stack;
    this.pos = pos;
    this.index = index;
    this.buffer = stack.buffer;
    if (this.index == 0)
      this.maybeNext();
  }
  static create(stack) {
    return new StackBufferCursor(stack, stack.bufferBase + stack.buffer.length, stack.buffer.length);
  }
  maybeNext() {
    let next = this.stack.parent;
    if (next != null) {
      this.index = this.stack.bufferBase - next.bufferBase;
      this.stack = next;
      this.buffer = next.buffer;
    }
  }
  get id() {
    return this.buffer[this.index - 4];
  }
  get start() {
    return this.buffer[this.index - 3];
  }
  get end() {
    return this.buffer[this.index - 2];
  }
  get size() {
    return this.buffer[this.index - 1];
  }
  next() {
    this.index -= 4;
    this.pos -= 4;
    if (this.index == 0)
      this.maybeNext();
  }
  fork() {
    return new StackBufferCursor(this.stack, this.pos, this.index);
  }
};
var Token = class {
  constructor() {
    this.start = -1;
    this.value = -1;
    this.end = -1;
  }
  accept(value, end) {
    this.value = value;
    this.end = end;
  }
};
var TokenGroup = class {
  constructor(data, id2) {
    this.data = data;
    this.id = id2;
  }
  token(input, token, stack) {
    readToken2(this.data, input, token, stack, this.id);
  }
};
TokenGroup.prototype.contextual = TokenGroup.prototype.fallback = TokenGroup.prototype.extend = false;
var ExternalTokenizer = class {
  constructor(token, options = {}) {
    this.token = token;
    this.contextual = !!options.contextual;
    this.fallback = !!options.fallback;
    this.extend = !!options.extend;
  }
};
function readToken2(data, input, token, stack, group) {
  let state = 0, groupMask = 1 << group, dialect2 = stack.p.parser.dialect;
  scan:
    for (let pos = token.start; ; ) {
      if ((groupMask & data[state]) == 0)
        break;
      let accEnd = data[state + 1];
      for (let i = state + 3; i < accEnd; i += 2)
        if ((data[i + 1] & groupMask) > 0) {
          let term = data[i];
          if (dialect2.allows(term) && (token.value == -1 || token.value == term || stack.p.parser.overrides(term, token.value))) {
            token.accept(term, pos);
            break;
          }
        }
      let next = input.get(pos++);
      for (let low = 0, high = data[state + 2]; low < high; ) {
        let mid = low + high >> 1;
        let index = accEnd + mid + (mid << 1);
        let from = data[index], to = data[index + 1];
        if (next < from)
          high = mid;
        else if (next >= to)
          low = mid + 1;
        else {
          state = data[index + 2];
          continue scan;
        }
      }
      break;
    }
}
function decodeArray(input, Type2 = Uint16Array) {
  if (typeof input != "string")
    return input;
  let array = null;
  for (let pos = 0, out = 0; pos < input.length; ) {
    let value = 0;
    for (; ; ) {
      let next = input.charCodeAt(pos++), stop = false;
      if (next == 126) {
        value = 65535;
        break;
      }
      if (next >= 92)
        next--;
      if (next >= 34)
        next--;
      let digit = next - 32;
      if (digit >= 46) {
        digit -= 46;
        stop = true;
      }
      value += digit;
      if (stop)
        break;
      value *= 46;
    }
    if (array)
      array[out++] = value;
    else
      array = new Type2(value);
  }
  return array;
}
var verbose = typeof process != "undefined" && /\bparse\b/.test(process.env.LOG);
var stackIDs = null;
function cutAt(tree, pos, side) {
  let cursor = tree.cursor(pos);
  for (; ; ) {
    if (!(side < 0 ? cursor.childBefore(pos) : cursor.childAfter(pos)))
      for (; ; ) {
        if ((side < 0 ? cursor.to <= pos : cursor.from >= pos) && !cursor.type.isError)
          return side < 0 ? Math.max(0, Math.min(cursor.to - 1, pos - 5)) : Math.min(tree.length, Math.max(cursor.from + 1, pos + 5));
        if (side < 0 ? cursor.prevSibling() : cursor.nextSibling())
          break;
        if (!cursor.parent())
          return side < 0 ? 0 : tree.length;
      }
  }
}
var FragmentCursor = class {
  constructor(fragments) {
    this.fragments = fragments;
    this.i = 0;
    this.fragment = null;
    this.safeFrom = -1;
    this.safeTo = -1;
    this.trees = [];
    this.start = [];
    this.index = [];
    this.nextFragment();
  }
  nextFragment() {
    let fr = this.fragment = this.i == this.fragments.length ? null : this.fragments[this.i++];
    if (fr) {
      this.safeFrom = fr.openStart ? cutAt(fr.tree, fr.from + fr.offset, 1) - fr.offset : fr.from;
      this.safeTo = fr.openEnd ? cutAt(fr.tree, fr.to + fr.offset, -1) - fr.offset : fr.to;
      while (this.trees.length) {
        this.trees.pop();
        this.start.pop();
        this.index.pop();
      }
      this.trees.push(fr.tree);
      this.start.push(-fr.offset);
      this.index.push(0);
      this.nextStart = this.safeFrom;
    } else {
      this.nextStart = 1e9;
    }
  }
  nodeAt(pos) {
    if (pos < this.nextStart)
      return null;
    while (this.fragment && this.safeTo <= pos)
      this.nextFragment();
    if (!this.fragment)
      return null;
    for (; ; ) {
      let last = this.trees.length - 1;
      if (last < 0) {
        this.nextFragment();
        return null;
      }
      let top2 = this.trees[last], index = this.index[last];
      if (index == top2.children.length) {
        this.trees.pop();
        this.start.pop();
        this.index.pop();
        continue;
      }
      let next = top2.children[index];
      let start = this.start[last] + top2.positions[index];
      if (start > pos) {
        this.nextStart = start;
        return null;
      } else if (start == pos && start + next.length <= this.safeTo) {
        return start == pos && start >= this.safeFrom ? next : null;
      }
      if (next instanceof TreeBuffer) {
        this.index[last]++;
        this.nextStart = start + next.length;
      } else {
        this.index[last]++;
        if (start + next.length >= pos) {
          this.trees.push(next);
          this.start.push(start);
          this.index.push(0);
        }
      }
    }
  }
};
var CachedToken = class extends Token {
  constructor() {
    super(...arguments);
    this.extended = -1;
    this.mask = 0;
    this.context = 0;
  }
  clear(start) {
    this.start = start;
    this.value = this.extended = -1;
  }
};
var dummyToken = new Token();
var TokenCache = class {
  constructor(parser6) {
    this.tokens = [];
    this.mainToken = dummyToken;
    this.actions = [];
    this.tokens = parser6.tokenizers.map((_) => new CachedToken());
  }
  getActions(stack, input) {
    let actionIndex = 0;
    let main = null;
    let {parser: parser6} = stack.p, {tokenizers} = parser6;
    let mask = parser6.stateSlot(stack.state, 3);
    let context = stack.curContext ? stack.curContext.hash : 0;
    for (let i = 0; i < tokenizers.length; i++) {
      if ((1 << i & mask) == 0)
        continue;
      let tokenizer = tokenizers[i], token = this.tokens[i];
      if (main && !tokenizer.fallback)
        continue;
      if (tokenizer.contextual || token.start != stack.pos || token.mask != mask || token.context != context) {
        this.updateCachedToken(token, tokenizer, stack, input);
        token.mask = mask;
        token.context = context;
      }
      if (token.value != 0) {
        let startIndex = actionIndex;
        if (token.extended > -1)
          actionIndex = this.addActions(stack, token.extended, token.end, actionIndex);
        actionIndex = this.addActions(stack, token.value, token.end, actionIndex);
        if (!tokenizer.extend) {
          main = token;
          if (actionIndex > startIndex)
            break;
        }
      }
    }
    while (this.actions.length > actionIndex)
      this.actions.pop();
    if (!main) {
      main = dummyToken;
      main.start = stack.pos;
      if (stack.pos == input.length)
        main.accept(stack.p.parser.eofTerm, stack.pos);
      else
        main.accept(0, stack.pos + 1);
    }
    this.mainToken = main;
    return this.actions;
  }
  updateCachedToken(token, tokenizer, stack, input) {
    token.clear(stack.pos);
    tokenizer.token(input, token, stack);
    if (token.value > -1) {
      let {parser: parser6} = stack.p;
      for (let i = 0; i < parser6.specialized.length; i++)
        if (parser6.specialized[i] == token.value) {
          let result = parser6.specializers[i](input.read(token.start, token.end), stack);
          if (result >= 0 && stack.p.parser.dialect.allows(result >> 1)) {
            if ((result & 1) == 0)
              token.value = result >> 1;
            else
              token.extended = result >> 1;
            break;
          }
        }
    } else if (stack.pos == input.length) {
      token.accept(stack.p.parser.eofTerm, stack.pos);
    } else {
      token.accept(0, stack.pos + 1);
    }
  }
  putAction(action, token, end, index) {
    for (let i = 0; i < index; i += 3)
      if (this.actions[i] == action)
        return index;
    this.actions[index++] = action;
    this.actions[index++] = token;
    this.actions[index++] = end;
    return index;
  }
  addActions(stack, token, end, index) {
    let {state} = stack, {parser: parser6} = stack.p, {data} = parser6;
    for (let set = 0; set < 2; set++) {
      for (let i = parser6.stateSlot(state, set ? 2 : 1); ; i += 3) {
        if (data[i] == 65535) {
          if (data[i + 1] == 1) {
            i = pair(data, i + 2);
          } else {
            if (index == 0 && data[i + 1] == 2)
              index = this.putAction(pair(data, i + 1), token, end, index);
            break;
          }
        }
        if (data[i] == token)
          index = this.putAction(pair(data, i + 1), token, end, index);
      }
    }
    return index;
  }
};
var Rec;
(function(Rec2) {
  Rec2[Rec2["Distance"] = 5] = "Distance";
  Rec2[Rec2["MaxRemainingPerStep"] = 3] = "MaxRemainingPerStep";
  Rec2[Rec2["MinBufferLengthPrune"] = 200] = "MinBufferLengthPrune";
  Rec2[Rec2["ForceReduceLimit"] = 10] = "ForceReduceLimit";
})(Rec || (Rec = {}));
var Parse2 = class {
  constructor(parser6, input, startPos, context) {
    this.parser = parser6;
    this.input = input;
    this.startPos = startPos;
    this.context = context;
    this.pos = 0;
    this.recovering = 0;
    this.nextStackID = 9812;
    this.nested = null;
    this.nestEnd = 0;
    this.nestWrap = null;
    this.reused = [];
    this.tokens = new TokenCache(parser6);
    this.topTerm = parser6.top[1];
    this.stacks = [Stack.start(this, parser6.top[0], this.startPos)];
    let fragments = context === null || context === void 0 ? void 0 : context.fragments;
    this.fragments = fragments && fragments.length ? new FragmentCursor(fragments) : null;
  }
  advance() {
    if (this.nested) {
      let result = this.nested.advance();
      this.pos = this.nested.pos;
      if (result) {
        this.finishNested(this.stacks[0], result);
        this.nested = null;
      }
      return null;
    }
    let stacks = this.stacks, pos = this.pos;
    let newStacks = this.stacks = [];
    let stopped, stoppedTokens;
    let maybeNest;
    for (let i = 0; i < stacks.length; i++) {
      let stack = stacks[i], nest;
      for (; ; ) {
        if (stack.pos > pos) {
          newStacks.push(stack);
        } else if (nest = this.checkNest(stack)) {
          if (!maybeNest || maybeNest.stack.score < stack.score)
            maybeNest = nest;
        } else if (this.advanceStack(stack, newStacks, stacks)) {
          continue;
        } else {
          if (!stopped) {
            stopped = [];
            stoppedTokens = [];
          }
          stopped.push(stack);
          let tok = this.tokens.mainToken;
          stoppedTokens.push(tok.value, tok.end);
        }
        break;
      }
    }
    if (maybeNest) {
      this.startNested(maybeNest);
      return null;
    }
    if (!newStacks.length) {
      let finished = stopped && findFinished(stopped);
      if (finished)
        return this.stackToTree(finished);
      if (this.parser.strict) {
        if (verbose && stopped)
          console.log("Stuck with token " + this.parser.getName(this.tokens.mainToken.value));
        throw new SyntaxError("No parse at " + pos);
      }
      if (!this.recovering)
        this.recovering = 5;
    }
    if (this.recovering && stopped) {
      let finished = this.runRecovery(stopped, stoppedTokens, newStacks);
      if (finished)
        return this.stackToTree(finished.forceAll());
    }
    if (this.recovering) {
      let maxRemaining = this.recovering == 1 ? 1 : this.recovering * 3;
      if (newStacks.length > maxRemaining) {
        newStacks.sort((a, b) => b.score - a.score);
        while (newStacks.length > maxRemaining)
          newStacks.pop();
      }
      if (newStacks.some((s) => s.reducePos > pos))
        this.recovering--;
    } else if (newStacks.length > 1) {
      outer:
        for (let i = 0; i < newStacks.length - 1; i++) {
          let stack = newStacks[i];
          for (let j = i + 1; j < newStacks.length; j++) {
            let other = newStacks[j];
            if (stack.sameState(other) || stack.buffer.length > 200 && other.buffer.length > 200) {
              if ((stack.score - other.score || stack.buffer.length - other.buffer.length) > 0) {
                newStacks.splice(j--, 1);
              } else {
                newStacks.splice(i--, 1);
                continue outer;
              }
            }
          }
        }
    }
    this.pos = newStacks[0].pos;
    for (let i = 1; i < newStacks.length; i++)
      if (newStacks[i].pos < this.pos)
        this.pos = newStacks[i].pos;
    return null;
  }
  advanceStack(stack, stacks, split) {
    let start = stack.pos, {input, parser: parser6} = this;
    let base3 = verbose ? this.stackID(stack) + " -> " : "";
    if (this.fragments) {
      let strictCx = stack.curContext && stack.curContext.tracker.strict, cxHash = strictCx ? stack.curContext.hash : 0;
      for (let cached = this.fragments.nodeAt(start); cached; ) {
        let match = this.parser.nodeSet.types[cached.type.id] == cached.type ? parser6.getGoto(stack.state, cached.type.id) : -1;
        if (match > -1 && cached.length && (!strictCx || (cached.contextHash || 0) == cxHash)) {
          stack.useNode(cached, match);
          if (verbose)
            console.log(base3 + this.stackID(stack) + ` (via reuse of ${parser6.getName(cached.type.id)})`);
          return true;
        }
        if (!(cached instanceof Tree) || cached.children.length == 0 || cached.positions[0] > 0)
          break;
        let inner = cached.children[0];
        if (inner instanceof Tree)
          cached = inner;
        else
          break;
      }
    }
    let defaultReduce = parser6.stateSlot(stack.state, 4);
    if (defaultReduce > 0) {
      stack.reduce(defaultReduce);
      if (verbose)
        console.log(base3 + this.stackID(stack) + ` (via always-reduce ${parser6.getName(defaultReduce & 65535)})`);
      return true;
    }
    let actions = this.tokens.getActions(stack, input);
    for (let i = 0; i < actions.length; ) {
      let action = actions[i++], term = actions[i++], end = actions[i++];
      let last = i == actions.length || !split;
      let localStack = last ? stack : stack.split();
      localStack.apply(action, term, end);
      if (verbose)
        console.log(base3 + this.stackID(localStack) + ` (via ${(action & 65536) == 0 ? "shift" : `reduce of ${parser6.getName(action & 65535)}`} for ${parser6.getName(term)} @ ${start}${localStack == stack ? "" : ", split"})`);
      if (last)
        return true;
      else if (localStack.pos > start)
        stacks.push(localStack);
      else
        split.push(localStack);
    }
    return false;
  }
  advanceFully(stack, newStacks) {
    let pos = stack.pos;
    for (; ; ) {
      let nest = this.checkNest(stack);
      if (nest)
        return nest;
      if (!this.advanceStack(stack, null, null))
        return false;
      if (stack.pos > pos) {
        pushStackDedup(stack, newStacks);
        return true;
      }
    }
  }
  runRecovery(stacks, tokens2, newStacks) {
    let finished = null, restarted = false;
    let maybeNest;
    for (let i = 0; i < stacks.length; i++) {
      let stack = stacks[i], token = tokens2[i << 1], tokenEnd = tokens2[(i << 1) + 1];
      let base3 = verbose ? this.stackID(stack) + " -> " : "";
      if (stack.deadEnd) {
        if (restarted)
          continue;
        restarted = true;
        stack.restart();
        if (verbose)
          console.log(base3 + this.stackID(stack) + " (restarted)");
        let done = this.advanceFully(stack, newStacks);
        if (done) {
          if (done !== true)
            maybeNest = done;
          continue;
        }
      }
      let force = stack.split(), forceBase = base3;
      for (let j = 0; force.forceReduce() && j < 10; j++) {
        if (verbose)
          console.log(forceBase + this.stackID(force) + " (via force-reduce)");
        let done = this.advanceFully(force, newStacks);
        if (done) {
          if (done !== true)
            maybeNest = done;
          break;
        }
        if (verbose)
          forceBase = this.stackID(force) + " -> ";
      }
      for (let insert2 of stack.recoverByInsert(token)) {
        if (verbose)
          console.log(base3 + this.stackID(insert2) + " (via recover-insert)");
        this.advanceFully(insert2, newStacks);
      }
      if (this.input.length > stack.pos) {
        if (tokenEnd == stack.pos) {
          tokenEnd++;
          token = 0;
        }
        stack.recoverByDelete(token, tokenEnd);
        if (verbose)
          console.log(base3 + this.stackID(stack) + ` (via recover-delete ${this.parser.getName(token)})`);
        pushStackDedup(stack, newStacks);
      } else if (!finished || finished.score < stack.score) {
        finished = stack;
      }
    }
    if (finished)
      return finished;
    if (maybeNest) {
      for (let s of this.stacks)
        if (s.score > maybeNest.stack.score) {
          maybeNest = void 0;
          break;
        }
    }
    if (maybeNest)
      this.startNested(maybeNest);
    return null;
  }
  forceFinish() {
    let stack = this.stacks[0].split();
    if (this.nested)
      this.finishNested(stack, this.nested.forceFinish());
    return this.stackToTree(stack.forceAll());
  }
  stackToTree(stack, pos = stack.pos) {
    if (this.parser.context)
      stack.emitContext();
    return Tree.build({
      buffer: StackBufferCursor.create(stack),
      nodeSet: this.parser.nodeSet,
      topID: this.topTerm,
      maxBufferLength: this.parser.bufferLength,
      reused: this.reused,
      start: this.startPos,
      length: pos - this.startPos,
      minRepeatType: this.parser.minRepeatTerm
    });
  }
  checkNest(stack) {
    let info = this.parser.findNested(stack.state);
    if (!info)
      return null;
    let spec = info.value;
    if (typeof spec == "function")
      spec = spec(this.input, stack);
    return spec ? {stack, info, spec} : null;
  }
  startNested(nest) {
    let {stack, info, spec} = nest;
    this.stacks = [stack];
    this.nestEnd = this.scanForNestEnd(stack, info.end, spec.filterEnd);
    this.nestWrap = typeof spec.wrapType == "number" ? this.parser.nodeSet.types[spec.wrapType] : spec.wrapType || null;
    if (spec.startParse) {
      this.nested = spec.startParse(this.input.clip(this.nestEnd), stack.pos, this.context);
    } else {
      this.finishNested(stack);
    }
  }
  scanForNestEnd(stack, endToken, filter) {
    for (let pos = stack.pos; pos < this.input.length; pos++) {
      dummyToken.start = pos;
      dummyToken.value = -1;
      endToken.token(this.input, dummyToken, stack);
      if (dummyToken.value > -1 && (!filter || filter(this.input.read(pos, dummyToken.end))))
        return pos;
    }
    return this.input.length;
  }
  finishNested(stack, tree) {
    if (this.nestWrap)
      tree = new Tree(this.nestWrap, tree ? [tree] : [], tree ? [0] : [], this.nestEnd - stack.pos);
    else if (!tree)
      tree = new Tree(NodeType.none, [], [], this.nestEnd - stack.pos);
    let info = this.parser.findNested(stack.state);
    stack.useNode(tree, this.parser.getGoto(stack.state, info.placeholder, true));
    if (verbose)
      console.log(this.stackID(stack) + ` (via unnest)`);
  }
  stackID(stack) {
    let id2 = (stackIDs || (stackIDs = new WeakMap())).get(stack);
    if (!id2)
      stackIDs.set(stack, id2 = String.fromCodePoint(this.nextStackID++));
    return id2 + stack;
  }
};
function pushStackDedup(stack, newStacks) {
  for (let i = 0; i < newStacks.length; i++) {
    let other = newStacks[i];
    if (other.pos == stack.pos && other.sameState(stack)) {
      if (newStacks[i].score < stack.score)
        newStacks[i] = stack;
      return;
    }
  }
  newStacks.push(stack);
}
var Dialect = class {
  constructor(source, flags, disabled) {
    this.source = source;
    this.flags = flags;
    this.disabled = disabled;
  }
  allows(term) {
    return !this.disabled || this.disabled[term] == 0;
  }
};
var id = (x) => x;
var ContextTracker = class {
  constructor(spec) {
    this.start = spec.start;
    this.shift = spec.shift || id;
    this.reduce = spec.reduce || id;
    this.reuse = spec.reuse || id;
    this.hash = spec.hash;
    this.strict = spec.strict !== false;
  }
};
var Parser = class {
  constructor(spec) {
    this.bufferLength = DefaultBufferLength;
    this.strict = false;
    this.cachedDialect = null;
    if (spec.version != 13)
      throw new RangeError(`Parser version (${spec.version}) doesn't match runtime version (${13})`);
    let tokenArray = decodeArray(spec.tokenData);
    let nodeNames = spec.nodeNames.split(" ");
    this.minRepeatTerm = nodeNames.length;
    this.context = spec.context;
    for (let i = 0; i < spec.repeatNodeCount; i++)
      nodeNames.push("");
    let nodeProps = [];
    for (let i = 0; i < nodeNames.length; i++)
      nodeProps.push([]);
    function setProp(nodeID, prop, value) {
      nodeProps[nodeID].push([prop, prop.deserialize(String(value))]);
    }
    if (spec.nodeProps)
      for (let propSpec of spec.nodeProps) {
        let prop = propSpec[0];
        for (let i = 1; i < propSpec.length; ) {
          let next = propSpec[i++];
          if (next >= 0) {
            setProp(next, prop, propSpec[i++]);
          } else {
            let value = propSpec[i + -next];
            for (let j = -next; j > 0; j--)
              setProp(propSpec[i++], prop, value);
            i++;
          }
        }
      }
    this.specialized = new Uint16Array(spec.specialized ? spec.specialized.length : 0);
    this.specializers = [];
    if (spec.specialized)
      for (let i = 0; i < spec.specialized.length; i++) {
        this.specialized[i] = spec.specialized[i].term;
        this.specializers[i] = spec.specialized[i].get;
      }
    this.states = decodeArray(spec.states, Uint32Array);
    this.data = decodeArray(spec.stateData);
    this.goto = decodeArray(spec.goto);
    let topTerms = Object.keys(spec.topRules).map((r) => spec.topRules[r][1]);
    this.nodeSet = new NodeSet(nodeNames.map((name2, i) => NodeType.define({
      name: i >= this.minRepeatTerm ? void 0 : name2,
      id: i,
      props: nodeProps[i],
      top: topTerms.indexOf(i) > -1,
      error: i == 0,
      skipped: spec.skippedNodes && spec.skippedNodes.indexOf(i) > -1
    })));
    this.maxTerm = spec.maxTerm;
    this.tokenizers = spec.tokenizers.map((value) => typeof value == "number" ? new TokenGroup(tokenArray, value) : value);
    this.topRules = spec.topRules;
    this.nested = (spec.nested || []).map(([name2, value, endToken, placeholder]) => {
      return {name: name2, value, end: new TokenGroup(decodeArray(endToken), 0), placeholder};
    });
    this.dialects = spec.dialects || {};
    this.dynamicPrecedences = spec.dynamicPrecedences || null;
    this.tokenPrecTable = spec.tokenPrec;
    this.termNames = spec.termNames || null;
    this.maxNode = this.nodeSet.types.length - 1;
    this.dialect = this.parseDialect();
    this.top = this.topRules[Object.keys(this.topRules)[0]];
  }
  parse(input, startPos = 0, context = {}) {
    if (typeof input == "string")
      input = stringInput(input);
    let cx = new Parse2(this, input, startPos, context);
    for (; ; ) {
      let done = cx.advance();
      if (done)
        return done;
    }
  }
  startParse(input, startPos = 0, context = {}) {
    if (typeof input == "string")
      input = stringInput(input);
    return new Parse2(this, input, startPos, context);
  }
  getGoto(state, term, loose = false) {
    let table = this.goto;
    if (term >= table[0])
      return -1;
    for (let pos = table[term + 1]; ; ) {
      let groupTag = table[pos++], last = groupTag & 1;
      let target = table[pos++];
      if (last && loose)
        return target;
      for (let end = pos + (groupTag >> 1); pos < end; pos++)
        if (table[pos] == state)
          return target;
      if (last)
        return -1;
    }
  }
  hasAction(state, terminal) {
    let data = this.data;
    for (let set = 0; set < 2; set++) {
      for (let i = this.stateSlot(state, set ? 2 : 1), next; ; i += 3) {
        if ((next = data[i]) == 65535) {
          if (data[i + 1] == 1)
            next = data[i = pair(data, i + 2)];
          else if (data[i + 1] == 2)
            return pair(data, i + 2);
          else
            break;
        }
        if (next == terminal || next == 0)
          return pair(data, i + 1);
      }
    }
    return 0;
  }
  stateSlot(state, slot) {
    return this.states[state * 6 + slot];
  }
  stateFlag(state, flag) {
    return (this.stateSlot(state, 0) & flag) > 0;
  }
  findNested(state) {
    let flags = this.stateSlot(state, 0);
    return flags & 4 ? this.nested[flags >> 10] : null;
  }
  validAction(state, action) {
    if (action == this.stateSlot(state, 4))
      return true;
    for (let i = this.stateSlot(state, 1); ; i += 3) {
      if (this.data[i] == 65535) {
        if (this.data[i + 1] == 1)
          i = pair(this.data, i + 2);
        else
          return false;
      }
      if (action == pair(this.data, i + 1))
        return true;
    }
  }
  nextStates(state) {
    let result = [];
    for (let i = this.stateSlot(state, 1); ; i += 3) {
      if (this.data[i] == 65535) {
        if (this.data[i + 1] == 1)
          i = pair(this.data, i + 2);
        else
          break;
      }
      if ((this.data[i + 2] & 65536 >> 16) == 0) {
        let value = this.data[i + 1];
        if (!result.some((v, i2) => i2 & 1 && v == value))
          result.push(this.data[i], value);
      }
    }
    return result;
  }
  overrides(token, prev) {
    let iPrev = findOffset(this.data, this.tokenPrecTable, prev);
    return iPrev < 0 || findOffset(this.data, this.tokenPrecTable, token) < iPrev;
  }
  configure(config2) {
    let copy = Object.assign(Object.create(Parser.prototype), this);
    if (config2.props)
      copy.nodeSet = this.nodeSet.extend(...config2.props);
    if (config2.top) {
      let info = this.topRules[config2.top];
      if (!info)
        throw new RangeError(`Invalid top rule name ${config2.top}`);
      copy.top = info;
    }
    if (config2.tokenizers)
      copy.tokenizers = this.tokenizers.map((t2) => {
        let found = config2.tokenizers.find((r) => r.from == t2);
        return found ? found.to : t2;
      });
    if (config2.dialect)
      copy.dialect = this.parseDialect(config2.dialect);
    if (config2.nested)
      copy.nested = this.nested.map((obj) => {
        if (!Object.prototype.hasOwnProperty.call(config2.nested, obj.name))
          return obj;
        return {name: obj.name, value: config2.nested[obj.name], end: obj.end, placeholder: obj.placeholder};
      });
    if (config2.strict != null)
      copy.strict = config2.strict;
    if (config2.bufferLength != null)
      copy.bufferLength = config2.bufferLength;
    return copy;
  }
  getName(term) {
    return this.termNames ? this.termNames[term] : String(term <= this.maxNode && this.nodeSet.types[term].name || term);
  }
  get eofTerm() {
    return this.maxNode + 1;
  }
  get hasNested() {
    return this.nested.length > 0;
  }
  get topNode() {
    return this.nodeSet.types[this.top[1]];
  }
  dynamicPrecedence(term) {
    let prec2 = this.dynamicPrecedences;
    return prec2 == null ? 0 : prec2[term] || 0;
  }
  parseDialect(dialect2) {
    if (this.cachedDialect && this.cachedDialect.source == dialect2)
      return this.cachedDialect;
    let values = Object.keys(this.dialects), flags = values.map(() => false);
    if (dialect2)
      for (let part of dialect2.split(" ")) {
        let id2 = values.indexOf(part);
        if (id2 >= 0)
          flags[id2] = true;
      }
    let disabled = null;
    for (let i = 0; i < values.length; i++)
      if (!flags[i]) {
        for (let j = this.dialects[values[i]], id2; (id2 = this.data[j++]) != 65535; )
          (disabled || (disabled = new Uint8Array(this.maxTerm + 1)))[id2] = 1;
      }
    return this.cachedDialect = new Dialect(dialect2, flags, disabled);
  }
  static deserialize(spec) {
    return new Parser(spec);
  }
};
function pair(data, off) {
  return data[off] | data[off + 1] << 16;
}
function findOffset(data, start, term) {
  for (let i = start, next; (next = data[i]) != 65535; i++)
    if (next == term)
      return i - start;
  return -1;
}
function findFinished(stacks) {
  let best = null;
  for (let stack of stacks) {
    if (stack.pos == stack.p.input.length && stack.p.parser.stateFlag(stack.state, 2) && (!best || best.score < stack.score))
      best = stack;
  }
  return best;
}

// node_modules/lezer-java/dist/index.es.js
var spec_identifier = {__proto__: null, true: 34, false: 34, null: 40, void: 44, byte: 46, short: 46, int: 46, long: 46, char: 46, float: 46, double: 46, boolean: 46, extends: 60, super: 62, class: 74, this: 76, new: 82, public: 98, protected: 100, private: 102, abstract: 104, static: 106, final: 108, strictfp: 110, default: 112, synchronized: 114, native: 116, transient: 118, volatile: 120, throws: 148, implements: 158, interface: 164, enum: 174, instanceof: 234, open: 263, module: 265, requires: 270, transitive: 272, exports: 274, to: 276, opens: 278, uses: 280, provides: 282, with: 284, package: 288, import: 292, if: 304, else: 306, while: 310, for: 314, assert: 326, switch: 330, case: 336, do: 340, break: 344, continue: 350, return: 356, throw: 362, try: 366, catch: 370, finally: 378};
var parser = Parser.deserialize({
  version: 13,
  states: "#'fQ]QPOOO&nQQO'#H[O)OQQO'#CbOOQO'#Cb'#CbO)VQPO'#CaOOQO'#Ha'#HaOOQO'#Ct'#CtO*oQPO'#D^O+YQQO'#HhOOQO'#Hh'#HhO-nQQO'#HcO-uQQO'#HcOOQO'#Hc'#HcOOQO'#Hb'#HbO-|QPO'#DTO0PQPO'#GlO1dQPO'#D^O2tQPO'#DyO)VQPO'#EZO2{QPO'#EZOOQO'#DU'#DUO4nQQO'#H_O6rQQO'#EdO6yQPO'#EcO7OQPO'#EeOOQO'#H`'#H`O5UQQO'#H`O8RQQO'#FfO8YQPO'#EvO8_QPO'#E{O8_QPO'#E}OOQO'#H_'#H_OOQO'#HW'#HWOOQO'#Gf'#GfOOQO'#HV'#HVO9lQPO'#FgOOQO'#HU'#HUOOQO'#Ge'#GeQ]QPOOOOQO'#Hn'#HnO9qQPO'#HnO9vQPO'#DzO9vQPO'#EUO9vQPO'#EPO:OQPO'#HkO:aQQO'#EeO)VQPO'#C`O:iQPO'#C`O)VQPO'#FaO:nQPO'#FcO:yQPO'#FiO:yQPO'#FlO;OQPO'#FnO8_QPO'#FtO:yQPO'#FvO]QPO'#F{O;TQPO'#F}O;]QPO'#GQO;eQPO'#GTO:yQPO'#GVO8_QPO'#GWO;lQPO'#GYOOQO'#H['#H[O<]QQO,58{OOQO'#HY'#HYOOQO'#Hd'#HdO>aQPO,59dO?fQPO,59xOOQO-E:d-E:dO)VQPO,58zO@VQPO,58zO)VQPO,5;{O@[QPO'#DOO@aQPO'#DOOOQO'#Gh'#GhOAjQQO,59iOOQO'#Dl'#DlOBuQPO'#HpOCPQPO'#DkOC_QPO'#HoOCgQPO,5<]OClQPO,59]ODVQPO'#CwOOQO,59b,59bOD^QPO,59aOFfQQO'#CbO)_QPO'#D^OG_QQO'#HhOGrQQO,59oOGyQPO'#DuOHXQPO'#HwOHaQPO,5:_OHfQPO,5:_OH|QPO,5;lOIXQPO'#IOOIdQPO,5;cOIiQPO,5=WOOQO-E:j-E:jOOQO,5:e,5:eOJ|QPO,5:eOKTQPO,5:uOKYQPO,5<]O)VQPO,5:uO9vQPO,5:fO9vQPO,5:pO9vQPO,5:kOKyQPO,59pOLQQPO,5:|OM_QPO,5;PO8_QPO,59TOMmQPO'#DWOOQO,5:},5:}OOQO'#Ek'#EkOOQO'#Em'#EmO8_QPO,5;TO8_QPO,5;TO8_QPO,5;TO8_QPO,5;TO8_QPO,5;TO8_QPO,5;TO8_QPO,5;dOOQO,5;g,5;gOOQO,5<Q,5<QOMtQPO,5;`ONVQPO,5;bOMtQPO'#CxON^QQO'#HhONlQQO,5;iO]QPO,5<ROOQO-E:c-E:cOOQO,5>Y,5>YO! |QPO,5:fO!![QPO,5:pO!!dQPO,5:kO!!oQPO,5>VOGyQPO,5>VOKhQPO,59UO!!zQQO,58zO!#SQQO,5;{O!#[QQO,5;}O)VQPO,5;}O8_QPO'#DTO]QPO,5<TO]QPO,5<WO!#dQPO'#FpO]QPO,5<YO]QPO,5<^O!$^QQO,5<`O!$hQPO,5<bO!$mQPO,5<gOOQO'#GP'#GPOOQO,5<i,5<iO!$rQPO,5<iOOQO'#GS'#GSOOQO,5<l,5<lO!$wQPO,5<lO!$|QQO,5<oOOQO,5<o,5<oO;oQPO,5<qO!%TQQO,5<rO!%[QPO'#GcO!&_QPO,5<tO;oQPO,5<|O)VQPO,58}O!*VQPO'#ChOOQO1G.k1G.kO!*aQPO,59iO!!zQQO1G.fO)VQPO1G.fO!+bQQO1G1gOOQO'#Gi'#GiO!,hQQO,59jO@[QPO,59jOOQO-E:f-E:fO!-hQPO,5>[O!.PQPO,5:VO9vQPO'#GnO!.WQPO,5>ZOOQO1G1w1G1wOOQO1G.w1G.wO!.qQPO'#CxO!/^QPO'#HhO!/kQPO'#CyO!/yQPO'#HgO!0RQPO,59cOOQO1G.{1G.{OD^QPO1G.{O!0iQPO,59dO!0vQQO'#H[O!1XQQO'#CbOOQO,5:a,5:aO9vQPO,5:bOOQO,5:`,5:`O!1jQQO,5:`OOQO1G/Z1G/ZO!1oQPO,5:aO!2QQPO'#GqO!2eQPO,5>cOOQO1G/y1G/yO!2mQPO'#DuO!3OQPO'#D^O!3VQPO1G/yOMtQPO'#GoO!3[QPO1G1WO8_QPO1G1WO9vQPO'#GwO!3dQPO,5>jOOQO1G0}1G0}OOQO1G0P1G0PO!3lQPO'#E[OOQO1G0a1G0aO!4]QPO1G1wOKTQPO1G0aO! |QPO1G0QO!![QPO1G0[O!!dQPO1G0VOOQO1G/[1G/[O!4bQQO1G.pO6yQPO1G0iO)VQPO1G0iO:OQPO'#HkO!6UQQO1G.pOOQO1G.p1G.pO!7XQQO1G0hOOQO1G0k1G0kO!7`QPO1G0kO!7kQQO1G.oO!8OQQO'#HlO!8]QPO,59rO!9iQQO1G0oO!:}QQO1G0oO!<YQQO1G0oO!<gQQO1G0oO!=iQQO1G0oO!>PQQO1G0oO!>^QQO1G1OO!>eQQO'#HhOOQO1G0z1G0zO!?hQQO1G0|OOQO1G0|1G0|OOQO1G1m1G1mOK]QPO'#DpO!AfQPO'#DZOMtQPO'#D{OMtQPO'#D|OOQO1G0Q1G0QO!AmQPO1G0QO!ArQPO1G0QO!AzQPO1G0QO!BVQPO'#EWOOQO1G0[1G0[O!BjQPO1G0[O!BoQPO'#ESOMtQPO'#EROOQO1G0V1G0VO!CiQPO1G0VO!CnQPO1G0VO!CvQPO'#EgO!C}QPO'#EgOOQO'#Gv'#GvO!DVQQO1G0lO!EvQQO1G3qO6yQPO1G3qO!GuQPO'#FVOOQO1G.f1G.fOOQO1G1g1G1gO!G|QPO1G1iOOQO1G1i1G1iO!HXQQO1G1iO!HaQPO1G1oOOQO1G1r1G1rO)_QPO'#D^O+YQQO,5<_OGyQPO,5<_O!LRQPO,5<[O!LYQPO,5<[OOQO1G1t1G1tOOQO1G1x1G1xOOQO1G1z1G1zO8_QPO1G1zO# vQPO'#FxOOQO1G1|1G1|O:yQPO1G2ROOQO1G2T1G2TOOQO1G2W1G2WOOQO1G2Z1G2ZOOQO1G2]1G2]OOQO1G2^1G2^O#!uQQO'#H[O#!|QQO'#CbO+YQQO'#HhO##wQQOOO#$eQQO'#EdO#$SQQO'#H`OGyQPO'#GdO#$lQPO,5<}OOQO'#HO'#HOO#$tQPO1G2`O#(lQPO'#G[O;oQPO'#G`OOQO1G2`1G2`O#(qQPO1G2hOOQO1G.i1G.iO#-sQQO'#EdO#.QQQO'#H^O#.bQPO'#FROOQO'#H^'#H^O#.lQPO'#H^O#/ZQPO'#IRO#/cQPO,59SO#/hQPO,59jOOQO7+$Q7+$QO!!zQQO7+$QOOQO7+'R7+'ROOQO-E:g-E:gO#0|QQO1G/UO#1|QPO'#DnO#2WQQO'#HqOOQO'#Hq'#HqOOQO1G/q1G/qOOQO,5=Y,5=YOOQO-E:l-E:lO#2hQSO,58{O#2oQPO,59eOOQO,59e,59eOMtQPO'#HjOCqQPO'#GgO#2}QPO,5>ROOQO1G.}1G.}OOQO7+$g7+$gOOQO1G/z1G/zO#3VQQO1G/zOOQO1G/|1G/|O#3[QPO1G/zOOQO1G/{1G/{O9vQPO1G/|OOQO,5=],5=]OOQO-E:o-E:oOOQO7+%e7+%eOOQO,5=Z,5=ZOOQO-E:m-E:mO8_QPO7+&rOOQO7+&r7+&rOOQO,5=c,5=cOOQO-E:u-E:uO#3aQPO'#ETO#3oQPO'#ETOOQO'#Gu'#GuO#4WQPO,5:vOOQO,5:v,5:vOOQO7+'c7+'cOOQO7+%{7+%{OOQO7+%l7+%lO!AmQPO7+%lO!ArQPO7+%lO!AzQPO7+%lOOQO7+%v7+%vO!BjQPO7+%vOOQO7+%q7+%qO!CiQPO7+%qO!CnQPO7+%qOOQO7+&T7+&TOOQO'#Ed'#EdO6yQPO7+&TO6yQPO,5>VO#4wQPO7+$[OOQO7+&S7+&SOOQO7+&V7+&VO8_QPO'#GjO#5VQPO,5>WOOQO1G/^1G/^O8_QPO7+&jO#5bQQO,59dO#6eQPO'#DqOK]QPO'#DqO#6pQPO'#HtO#6xQPO,5:[O#7cQQO'#HdO#8OQQO'#CtOKYQPO'#HsO#8nQPO'#DoO#8xQPO'#HsO#9ZQPO'#DoO#9cQPO'#H|O#9hQPO'#E_OOQO'#Hm'#HmOOQO'#Gk'#GkO#9pQPO,59uOOQO,59u,59uO#9wQPO'#HnOOQO,5:g,5:gO#;_QPO'#HyOOQO'#EO'#EOOOQO,5:h,5:hO#;jQPO'#EXO9vQPO'#EXO#;{QPO'#HzO#<WQPO,5:rOKYQPO'#HsOMtQPO'#HsO#<`QPO'#DoOOQO'#Gs'#GsO#<gQPO,5:nOOQO,5:n,5:nOOQO,5:m,5:mOOQO,5;R,5;RO#=aQQO,5;RO#=hQPO,5;ROOQO-E:t-E:tOOQO7+&W7+&WOOQO7+)]7+)]O#=oQQO7+)]OOQO'#Gz'#GzO#?]QPO,5;qOOQO,5;q,5;qO#?dQPO'#FWO)VQPO'#FWO)VQPO'#FWO)VQPO'#FWO#?rQPO7+'TO#?wQPO7+'TOOQO7+'T7+'TO]QPO7+'ZO#@SQPO1G1yOKYQPO1G1yO#@bQQO1G1vOMmQPO1G1vO#@iQPO1G1vO#@pQQO7+'fOOQO'#G}'#G}O#@wQPO,5<dOOQO,5<d,5<dO#AOQPO'#HnO8_QPO'#FyO#AWQPO7+'mO#A]QPO,5=OOKYQPO,5=OO#AbQPO1G2iO#BhQPO1G2iOOQO1G2i1G2iOOQO-E:|-E:|OOQO7+'z7+'zO!2QQPO'#G^O;oQPO,5<vOOQO,5<z,5<zO#BpQPO7+(SOOQO7+(S7+(SO#FhQPO,59TO#FoQPO'#IQO#FwQPO,5;mO)VQPO'#GyO#F|QPO,5>mOOQO1G.n1G.nO#GUQPO1G/UOOQO<<Gl<<GlO#GlQPO'#HrO#GtQPO,5:YOOQO1G/P1G/POOQO,5>U,5>UOOQO,5=R,5=ROOQO-E:e-E:eO#GyQPO7+%fOOQO7+%f7+%fOOQO7+%h7+%hOOQO<<J^<<J^O#HaQPO'#H[O#HhQPO'#CbO#HoQPO,5:oO#HtQPO,5:wO#3aQPO,5:oOOQO-E:s-E:sOOQO1G0b1G0bOOQO<<IW<<IWO!AmQPO<<IWO!ArQPO<<IWOOQO<<Ib<<IbOOQO<<I]<<I]O!CiQPO<<I]OOQO<<Io<<IoO#HyQQO<<GvO6yQPO<<IoO)VQPO<<IoOOQO<<Gv<<GvO#JmQQO,5=UOOQO-E:h-E:hO#JzQQO<<JUOOQO,5:],5:]OMtQPO'#DrO#K_QPO,5:]OK]QPO'#GpO#KjQPO,5>`OOQO1G/v1G/vO#KrQPO'#HpO#KyQPO,59wO#LOQPO,5>_OKYQPO,59wO#LZQPO,5:ZO#9hQPO,5:yOKYQPO,5>_OMtQPO,5>_O#9cQPO,5>hOOQO,5:Z,5:ZOHfQPO'#DsOOQO,5>h,5>hO#LcQPO'#E`OOQO,5:y,5:yO$ ^QPO,5:yOMtQPO'#DwOOQO-E:i-E:iOOQO1G/a1G/aOOQO,5:x,5:xOMtQPO'#GrO$ cQPO,5>eOOQO,5:s,5:sO$ nQPO,5:sO$ |QPO,5:sO$!_QPO'#GtO$!uQPO,5>fO$#QQPO'#EYOOQO1G0^1G0^O$#XQPO1G0^OKYQPO,5:oOOQO-E:q-E:qOOQO1G0Y1G0YOOQO1G0m1G0mO$#^QQO1G0mOOQO<<Lw<<LwOOQO-E:x-E:xOOQO1G1]1G1]O$#eQQO,5;rOOQO'#G{'#G{O#?dQPO,5;rOOQO'#IS'#ISO$#mQQO,5;rO$$OQQO,5;rOOQO<<Jo<<JoO$$WQPO<<JoOOQO<<Ju<<JuO8_QPO7+'eO$$]QPO7+'eOMmQPO7+'bO$$kQPO7+'bO$$pQQO7+'bOOQO<<KQ<<KQOOQO-E:{-E:{OOQO1G2O1G2OOOQO,5<e,5<eO$$wQQO,5<eOOQO<<KX<<KXO8_QPO1G2jO$%OQPO1G2jOOQO,5=l,5=lOOQO7+(T7+(TO$%TQPO7+(TOOQO-E;O-E;OO$&oQSO'#HcO$&ZQSO'#HcO$&vQPO'#G_O9vQPO,5<xOGyQPO,5<xOOQO1G2b1G2bOOQO<<Kn<<KnO$'XQQO1G.oOOQO1G1Y1G1YO$'cQPO'#GxO$'pQPO,5>lOOQO1G1X1G1XO$'xQPO'#FSOOQO,5=e,5=eOOQO-E:w-E:wO$'}QPO'#GmO$([QPO,5>^OOQO1G/t1G/tOOQO<<IQ<<IQOOQO1G0Z1G0ZO$(dQPO1G0cO$(iQPO1G0ZO$(nQPO1G0cOOQOAN>rAN>rO!AmQPOAN>rOOQOAN>wAN>wOOQOAN?ZAN?ZO6yQPOAN?ZO$(sQPO,5:^OOQO1G/w1G/wOOQO,5=[,5=[OOQO-E:n-E:nO$)OQPO,5>bOOQO1G/c1G/cOOQO1G3y1G3yO$)aQPO1G/cOOQO1G/u1G/uOOQO1G0e1G0eO$ ^QPO1G0eO#9cQPO'#HvO$)fQPO1G3yOKYQPO1G3yOOQO1G4S1G4SO$)qQPO'#DuO)_QPO'#D^OOQO,5:z,5:zO$)xQPO,5:zO$)xQPO,5:zO$*PQQO'#H_O$+_QQO'#H`O$+iQQO'#EaO$+tQPO'#EaOOQO,5:c,5:cOOQO,5=^,5=^OOQO-E:p-E:pOOQO1G0_1G0_O$+|QPO1G0_OOQO,5=`,5=`OOQO-E:r-E:rO$,[QPO,5:tOOQO7+%x7+%xOOQO7+&X7+&XOOQO1G1^1G1^O$,cQQO1G1^OOQO-E:y-E:yO$,kQQO'#ITO$,fQPO1G1^O$#sQPO1G1^O)VQPO1G1^OOQOAN@ZAN@ZO$,vQQO<<KPO8_QPO<<KPO$,}QPO<<J|OOQO<<J|<<J|OMmQPO<<J|OOQO1G2P1G2PO$-SQQO7+(UO8_QPO7+(UOOQO<<Ko<<KoP!%[QPO'#HQOGyQPO'#HPO$-^QPO,5<yO$-iQPO1G2dO9vQPO1G2dOOQO,5=d,5=dOOQO-E:v-E:vO#FhQPO,5;nOOQO,5=X,5=XOOQO-E:k-E:kO$-nQPO7+%}OOQO7+%u7+%uO$-|QPO7+%}OOQOG24^G24^OOQOG24uG24uOOQO,59j,59jO$.RQPO1G/xO$.^QPO1G3|OOQO7+$}7+$}OOQO7+&P7+&POOQO7+)e7+)eO$.oQPO7+)eO!0WQPO,5:`OOQO1G0f1G0fO$.zQPO1G0fO$/RQPO,59pO$/gQPO,5:{O6yQPO,5:{OOQO7+%y7+%yOOQO7+&x7+&xO)VQPO'#G|O$/lQPO,5>oO$/tQPO7+&xO$/yQQO'#IUOOQOAN@kAN@kO$0UQQOAN@kOOQOAN@hAN@hO$0]QPOAN@hO$0bQQO<<KpO$0lQPO,5=kOOQO-E:}-E:}OOQO7+(O7+(OO$0}QPO7+(OO$1SQPO<<IiOOQO<<Ii<<IiO#FhQPO<<IiO$1SQPO<<IiOOQO1G/U1G/UOOQO<<MP<<MPOOQO7+&Q7+&QO$1bQPO1G0iO$1mQQO1G0gOOQO1G0g1G0gO$1uQPO1G0gO$1zQQO,5=hOOQO-E:z-E:zOOQO<<Jd<<JdO$2VQPO,5>pOOQOG26VG26VOOQOG26SG26SOOQO<<Kj<<KjOOQOAN?TAN?TO#FhQPOAN?TO$2_QPOAN?TO$2dQPOAN?TO6yQPO7+&RO$2rQPO7+&ROOQO7+&R7+&RO$2wQPOG24oOOQOG24oG24oO#FhQPOG24oO$2|QPO<<ImOOQO<<Im<<ImOOQOLD*ZLD*ZO$3RQPOLD*ZOOQOAN?XAN?XOOQO!$'Mu!$'MuO$3WQQO'#H[O)VQPO'#CaO@[QPO'#DOO@[QPO'#DOO$3nQQO,59iO$3xQQO'#CbOMtQPO'#CxO@[QPO,59jO@[QPO,59jO$4]QQO'#HhO$5dQQO,59dO$6kQPO'#DOO$6sQPO'#DOOMtQPO,5;`O$6{QQO1G.oO$8QQQO1G0oO$9PQQO1G0oO$9^QQO1G0oO$:YQQO1G0oO$:jQQO1G0|O$:qQQO<<JUOMtQPO'#CxOLQQPO,59TOLQQPO,5;TOLQQPO,5;TOLQQPO,5;TOLQQPO,5;TO$:xQPO,5;bOLQQPO7+&jO$;PQQO'#EdO$;^QQO'#H`O$;kQPO'#EvOLQQPO'#E{OLQQPO'#E}OLQQPO,5;TOLQQPO,5;TOLQQPO1G1WO$;pQQO1G0oO$<QQQO1G1OOLQQPO7+&rO$<XQPO,5;lO8_QPO,5;dO$<dQPO1G1WO-|QPO'#DT",
  stateData: "$<o~OPOSQOS%wOS~OZ^O_TO`TOaTObTOcTOdTOf[Og[Oo}OuyOviOy|O|aO!OtO!RwO!SwO!TwO!UwO!VwO!WwO!XwO!YxO!Z!]O![wO!]wO!^wO!t{O!yzO#dnO#pmO#rnO#snO#w!PO#x!OO$U!QO$W!RO$^!SO$a!TO$c!UO$i!VO$k!WO$p!XO$r!YO$u!ZO$x![O${!^O$}!_O%{SO%}QO&PPO&obO~OWhXW&OXZ&OXshXs&OX!a&OX#[&OX#^&OX#`&OX#b&OX#c&OX#d&OX#e&OX#f&OX#g&OX#i&OX#m&OX#p&OX%{hX%}hX&PhX&X&OX&YhX&Y&OX&i&OX&qhX&q&OX&s!`XY&OX~O!O&OX#n&OXt&OXp&OX{&OX~P$qOWUXW&WXZUXsUXs&WX!OUX!aUX#[UX#^UX#`UX#bUX#cUX#dUX#eUX#fUX#gUX#iUX#mUX#pUX%{&WX%}&WX&P&WX&XUX&YUX&Y&WX&iUX&qUX&q&WX&s!`X~O#n$[X~P'RO%}RO&P!`O~Of[Og[O!RwO!SwO!TwO!UwO!VwO!WwO!XwO!YxO![wO!]wO!^wO%{SO%}!cO&PUOf!QXg!QX%}!QX&P!QX~O#w!hO#x!gO$U!iOu!QX!t!QX!y!QX&o!QX~P)_OW!sOs!jO%{SO%}!nO&P!nO&q&[X~OW!vOs&VX%{&VX%}&VX&P&VX&q&VXY&VXv&VX&i&VX&l&VXZ&VXp&VX&X&VX!O&VX#^&VX#`&VX#b&VX#c&VX#d&VX#e&VX#f&VX#g&VX#i&VX#m&VX#p&VX|&VX!q&VX#n&VXt&VX{&VX~O&Y!tO~P+nO&Y&VX~P+nOZ^O_TO`TOaTObTOcTOdTOf[Og[Oo}OviOy|O!RwO!SwO!TwO!UwO!VwO!WwO!XwO!YxO![wO!]wO!^wO#dnO#pmO#rnO#snO%{SO%}!wO&P0eOY&kP~O%{SOf%`Xg%`Xu%`X!R%`X!S%`X!T%`X!U%`X!V%`X!W%`X!X%`X!Y%`X![%`X!]%`X!^%`X!t%`X!y%`X%}%`X&P%`X&o%`X&Y%`X~O!RwO!SwO!TwO!UwO!VwO!WwO!XwO!YxO![wO!]wO!^wOf!QXg!QXu!QX!t!QX!y!QX%}!QX&P!QX&o!QX&Y!QX~O{#UO~P]Of[Og[Ou#ZO!t#]O!y#[O%}!cO&PUO&o#YO~Os#_O&q#`O!O&RX#^&RX#`&RX#b&RX#c&RX#d&RX#e&RX#f&RX#g&RX#i&RX#m&RX#p&RX&X&RX&Y&RX&i&RX~OW#^OY&RX#n&RXt&RXp&RX{&RX~P3gO!a#aO#[#aOW&SXs&SX!O&SX#^&SX#`&SX#b&SX#c&SX#d&SX#e&SX#f&SX#g&SX#i&SX#m&SX#p&SX&X&SX&Y&SX&i&SX&q&SXY&SX#n&SXp&SX{&SX~OZ#WX~P5UOZ#bO~O&q#`O~O#^#fO#`#gO#b#hO#c#hO#d#iO#e#jO#f#kO#g#kO#i#oO#m#lO#p#mO&X#dO&Y#dO&i#eO~O!O#nO~P7TO&s#pO~OZ^O_TO`TOaTObTOcTOdTOf[Og[Oo}OviOy|O#dnO#pmO#rnO#snO%{SO%}0jO&PPO~O#n#tO~O!Z#vO~O%}!nO&P!nO~Of[Og[O%}!cO&PUO&Y!tO~OW#|O&q#`O~O#x!gO~O!V$QO%}RO&P!`O~OZ$RO~OZ$UO~O!O$]O%}$[O~O!O$`O%}$_O~O!O$cO~P8_OZ$fO|aO~OW$iOZ$jOfTagTa%{Ta%}Ta&PTa~OuTa!RTa!STa!TTa!UTa!VTa!WTa!XTa!YTa![Ta!]Ta!^Ta!tTa!yTa#wTa#xTa$UTa&oTasTaYTa&YTapTa{Ta!OTa~P;tO%{SOpla&XlaYla&ila!Ola~Os0gO&qla|la!qla~P={O!RwO!SwO!TwO!UwO!VwO!WwO!XwO!YxO![wO!]wO!^wO~Of!Qag!Qau!Qa!t!Qa!y!Qa%}!Qa&P!Qa&o!Qa&Y!Qa~P>qO#x$nO~Ot$pO~Os$rO%{SO~O%{qa&iqa#^qa#`qa#bqa#cqa#dqa#eqa#fqa#gqa#iqa#mqa#pqa&Xqa&Yqa~Os!jOWqa%}qa&Pqa&qqaYqavqa&lqa!Oqa#nqapqa{qa~P@iOs0gO%{SOp&dX!O&dX!a&dX~OY&dX#n&dX~PBdO!a$uOp!_X!O!_XY!_X~Op$vO!O&cX~O!O$xO~Ou$yO~Of[Og[O%{0fO%}!cO&PUO&]$|O~O&X&ZP~PCqO%{SO%}!cO&PUO~OWUXW&WXYUXZUXsUXs&WX!aUX#[UX#^UX#`UX#bUX#cUX#dUX#eUX#fUX#gUX#iUX#mUX#pUX%{&WX%}&WX&P&WX&XUX&YUX&Y&WX&iUX&qUX&q&WX&s!`X~OY!`XY&WXp!`Xv&WX&i&WX&l&WX~PDiOv%WO%{SO%}%TO&P%SO&l%VO~OW!sOs!jOY&[X&i&[X&q&[X~PF|OY%YO~P7TOf[Og[O%}!cO&PUO~Op%[OY&kX~OY%^O~Of[Og[O%{SO%}!cO&PUOY&kP~P>qOY%dO&i%bO&q#`O~Op%eO&s#pOY&rX~OY%gO~O%{SOf%`ag%`au%`a!R%`a!S%`a!T%`a!U%`a!V%`a!W%`a!X%`a!Y%`a![%`a!]%`a!^%`a!t%`a!y%`a%}%`a&P%`a&o%`a&Y%`a~O{%hO~P]O|%iO~Os0gO%{SO%}!nO&P!nO~Oo%uOv%vO%}RO&P!`O&Y!tO~Oy%tO~PKhOZ1bO_TO`TOaTObTOcTOdTOf[Og[Oo}OviOy|O#d1WO#p1VO#r1WO#s1WO%{SO%}0jO&PPO~Oy%xO%}RO&P!`O&Y!tO~OY&`P~P8_Of[Og[O%{SO%}!cO&PUO~O|aO~P8_OW!sOs!jO%{SO&q&[X~O#p#mO!O#qa#^#qa#`#qa#b#qa#c#qa#d#qa#e#qa#f#qa#g#qa#i#qa#m#qa&X#qa&Y#qa&i#qaY#qa#n#qat#qap#qa{#qa~On&]O|&[O!q&^O&Y&ZO~O|&cO!q&^O~On&gO|&fO&Y&ZO~OZ#bOs&kO%{SO~OW$iO|&qO~OW$iO!O&sO~OW&tO!O&uO~O!RwO!SwO!TwO!UwO!VwO!WwO!XwO!YxO![wO!]wO!^wO!O&`P~P8_O!O'QO#n'RO~P7TO|'SO~O$a'UO~O!O'VO~O!O'WO~O!O'XO~P7TO!O'ZO~P7TOZ$RO_TO`TOaTObTOcTOdTOf[Og[Oo}OviOy|O%{SO%}']O&P'[O~P>qO%P'fO%T'gOZ$|a_$|a`$|aa$|ab$|ac$|ad$|af$|ag$|ao$|au$|av$|ay$|a|$|a!O$|a!R$|a!S$|a!T$|a!U$|a!V$|a!W$|a!X$|a!Y$|a!Z$|a![$|a!]$|a!^$|a!t$|a!y$|a#d$|a#p$|a#r$|a#s$|a#w$|a#x$|a$U$|a$W$|a$^$|a$a$|a$c$|a$i$|a$k$|a$p$|a$r$|a$u$|a$x$|a${$|a$}$|a%u$|a%{$|a%}$|a&P$|a&o$|a{$|a$_$|a$n$|a~O|'mOY&uP~P8_Os0gO%{qa&qqa%}qa&Pqapqa&XqaYqavqa&iqa&lqa|qa!qqa&yqa!Oqa~OW$iO!O'uO~Ot$pOsra%{ra%}ra&Pra&qraYravra&ira&lra!Ora&Xra#nrapra~OWra#^ra#`ra#bra#cra#dra#era#fra#gra#ira#mra#pra&Yra{ra~P!+jOs0gO%{SOp&da!O&da!a&daY&da#n&da~O|'xO~P8_Op$vO!O&ca~Of[Og[O%{0fO%}!cO&PUO~O&](PO~P!.`O%{SOp&[X&X&[XY&[X&i&[X!O&[X~Os0gO|&[X!q&[X~P!.xOn(ROo(ROpmX&XmX~Op(SO&X&ZX~O&X(UO~Os0gOv(WO%{SO%}RO&P!`O~OYla&ila&qla~P!0WOW&OXY!`Xp!`Xs!`X%{!`X~OWUXY!`Xp!`Xs!`X%{!`X~OW(ZO~Os0gO%{SO%}!nO&P!nO&l(]O~Of[Og[O%{SO%}!cO&PUO~P>qOp%[OY&ka~Os0gO%{SO%}!nO&P!nO&l%VO~O%{SO~P1dOY(`O~OY(cO&i%bO~Op%eOY&ra~Of[Og[OuyO{(kO!t{O%{SO%}!cO&PUO&obO~P>qO!O(lO~OW^iZ#WXs^i!O^i!a^i#[^i#^^i#`^i#b^i#c^i#d^i#e^i#f^i#g^i#i^i#m^i#p^i&X^i&Y^i&i^i&q^iY^i#n^it^ip^i{^i~OW({O~O#^1XO#`0|O#b0}O#c0}O#d1OO#e1PO#f1YO#g1YO#i0rO#m1`O#p#mO&X#dO&Y#dO&i#eO~Ot(|O~P!6ZOy(}O%}RO&P!`O~O!O]iY]i#n]ip]i{]i~P7TOp)OOY&`X!O&`X~P7TOY)QO~O#p#mO!O#]i#^#]i#`#]i#b#]i#c#]i#d#]i#e#]i#i#]i#m#]i&X#]i&Y#]i&i#]iY#]i#n#]ip#]i{#]i~O#f#kO#g#kO~P!8bO#^#fO#e#jO#f#kO#g#kO#i#oO#p#mO&X#dO&Y#dO!O#]i#`#]i#b#]i#c#]i#m#]i&i#]iY#]i#n#]ip#]i{#]i~O#d#iO~P!9sO#^#fO#e#jO#f#kO#g#kO#i#oO#p#mO&X#dO&Y#dO!O#]i#b#]i#c#]i#m#]iY#]i#n#]ip#]i{#]i~O#`#gO#d#iO&i#eO~P!;UO#d#]i~P!9sO#p#mO!O#]i#`#]i#b#]i#c#]i#d#]i#e#]i#m#]i&i#]iY#]i#n#]ip#]i{#]i~O#^#fO#f#kO#g#kO#i#oO&X#dO&Y#dO~P!<nO#f#]i#g#]it#]i~P!8bO#n)RO~P7TOs!jO#^&[X#`&[X#b&[X#c&[X#d&[X#e&[X#f&[X#g&[X#i&[X#m&[X#p&[X&Y&[X#n&[X{&[X~P!.xO!O#jiY#ji#n#jip#ji{#ji~P7TOf[Og[OuyO|aO!O)aO!RwO!SwO!TwO!UwO!V)eO!WwO!XwO!YxO![wO!]wO!^wO!t{O!yzO%{SO%})XO&P)YO&Y&ZO&obO~O{)dO~P!?{O|&[O~O|&[O!q&^O~On&]O|&[O!q&^O~O%{SO%}!nO&P!nO{&nP!O&nP~P>qO|&cO~Of[Og[OuyO{)sO!O)qO!t{O!yzO%{SO%}!cO&PUO&Y&ZO&obO~P>qO|&fO~On&gO|&fO~Ot)uO~PLQOs)wO%{SO~Os&kO|'xO%{SOW#Yi!O#Yi#^#Yi#`#Yi#b#Yi#c#Yi#d#Yi#e#Yi#f#Yi#g#Yi#i#Yi#m#Yi#p#Yi&X#Yi&Y#Yi&i#Yi&q#YiY#Yi#n#Yit#Yip#Yi{#Yi~O|&[OW&_is&_i!O&_i#^&_i#`&_i#b&_i#c&_i#d&_i#e&_i#f&_i#g&_i#i&_i#m&_i#p&_i&X&_i&Y&_i&i&_i&q&_iY&_i#n&_it&_ip&_i{&_i~O#{*PO#}*QO$P*QO$Q*RO$R*SO~O{*OO~P!GdO$X*TO%}RO&P!`O~OW*UO!O*VO~O$_*WOZ$]i_$]i`$]ia$]ib$]ic$]id$]if$]ig$]io$]iu$]iv$]iy$]i|$]i!O$]i!R$]i!S$]i!T$]i!U$]i!V$]i!W$]i!X$]i!Y$]i!Z$]i![$]i!]$]i!^$]i!t$]i!y$]i#d$]i#p$]i#r$]i#s$]i#w$]i#x$]i$U$]i$W$]i$^$]i$a$]i$c$]i$i$]i$k$]i$p$]i$r$]i$u$]i$x$]i${$]i$}$]i%u$]i%{$]i%}$]i&P$]i&o$]i{$]i$n$]i~O!O*[O~P8_O!O*]O~OZ^O_TO`TOaTObTOcTOdTOf[Og[Oo}OuyOviOy|O|aO!OtO!RwO!SwO!TwO!UwO!VwO!WwO!XwO!Y*bO!Z!]O![wO!]wO!^wO!t{O!yzO#dnO#pmO#rnO#snO#w!PO#x!OO$U!QO$W!RO$^!SO$a!TO$c!UO$i!VO$k!WO$n*cO$p!XO$r!YO$u!ZO$x![O${!^O$}!_O%{SO%}QO&PPO&obO~O{*aO~P!L_OWhXW&OXY&OXZ&OXshXs&OX%{hX%}hX&PhX&YhX&qhX&q&OX~O!O&OX~P# }OWUXW&WXYUXZUXsUXs&WX!OUX%{&WX%}&WX&P&WX&Y&WX&qUX&q&WX~OW#^Os#_O&q#`O~OW&SXY%WXs&SX!O%WX&q&SX~OZ#WX~P#$SOY*iO!O*gO~O%P'fO%T'gOZ$|i_$|i`$|ia$|ib$|ic$|id$|if$|ig$|io$|iu$|iv$|iy$|i|$|i!O$|i!R$|i!S$|i!T$|i!U$|i!V$|i!W$|i!X$|i!Y$|i!Z$|i![$|i!]$|i!^$|i!t$|i!y$|i#d$|i#p$|i#r$|i#s$|i#w$|i#x$|i$U$|i$W$|i$^$|i$a$|i$c$|i$i$|i$k$|i$p$|i$r$|i$u$|i$x$|i${$|i$}$|i%u$|i%{$|i%}$|i&P$|i&o$|i{$|i$_$|i$n$|i~OZ*lO~O%P'fO%T'gOZ%Ui_%Ui`%Uia%Uib%Uic%Uid%Uif%Uig%Uio%Uiu%Uiv%Uiy%Ui|%Ui!O%Ui!R%Ui!S%Ui!T%Ui!U%Ui!V%Ui!W%Ui!X%Ui!Y%Ui!Z%Ui![%Ui!]%Ui!^%Ui!t%Ui!y%Ui#d%Ui#p%Ui#r%Ui#s%Ui#w%Ui#x%Ui$U%Ui$W%Ui$^%Ui$a%Ui$c%Ui$i%Ui$k%Ui$p%Ui$r%Ui$u%Ui$x%Ui${%Ui$}%Ui%u%Ui%{%Ui%}%Ui&P%Ui&o%Ui{%Ui$_%Ui$n%Ui~OW&SXZ#WXs&SX#^&SX#`&SX#b&SX#c&SX#d&SX#e&SX#f&SX#g&SX#i&SX#m&SX#p&SX&X&SX&Y&SX&i&SX&q&SX~O!a*qO#[#aOY&SX~P#,iOY&QXp&QX{&QX!O&QX~P7TO|'mO{&tP~P8_OY&QXf%YXg%YX%{%YX%}%YX&P%YXp&QX{&QX!O&QX~Op*tOY&uX~OY*vO~O!ara|ra!qra&yra!lra!Yra~P!+jOt$pOsri%{ri%}ri&Pri&qriYrivri&iri&lri!Ori&Xri#nripri~OWri#^ri#`ri#bri#cri#dri#eri#fri#gri#iri#mri#pri&Yri{ri~P#0OO|'xO{&fP~P8_Op&eX!O&eX{&eXY&eX~P7TO&]Ta~P;tOn(ROo(ROpma&Xma~Op(SO&X&Za~OW+PO~Ov+QO~Os0gO%{SO%}+UO&P+TO~Of[Og[Ou#ZO!t#]O%}!cO&PUO&o#YO~Of[Og[OuyO{+ZO!t{O%{SO%}!cO&PUO&obO~P>qOv+fO%}RO&P!`O&Y!tO~Op)OOY&`a!O&`a~Os!jO#^la#`la#bla#cla#dla#ela#fla#gla#ila#mla#pla&Yla#nla{la~P={On+kOp!eX&X!eX~Op+mO&X&hX~O&X+oO~OW&WXs&WX%{&WX%}&WX&P&WX&Y&WX~OZ!`X~P#6}OWhXshX%{hX%}hX&PhX&YhX~OZ!`X~P#7jOf[Og[Ou#ZO!t#]O!y#[O&Y&ZO&o#YO~O%})XO&P)YO~P#8VOf[Og[O%{SO%})XO&P)YO~O|aO!O+yO~OZ+zO~O|+|O!l,PO~O{,RO~P!?{O|aOf&bXg&bXu&bX!R&bX!S&bX!T&bX!U&bX!V&bX!W&bX!X&bX!Y&bX![&bX!]&bX!^&bX!t&bX!y&bX%{&bX%}&bX&P&bX&Y&bX&o&bX~Op,TO|&mX!O&mX~OZ#bO|&[Op!{X{!{X!O!{X~Op,YO{&nX!O&nX~O{,]O!O,[O~O&Y&ZO~P2{Of[Og[OuyO{,aO!O)qO!t{O!yzO%{SO%}!cO&PUO&Y&ZO&obO~P>qOt,bO~P!6ZOt,bO~PLQO|&[OW&_qs&_q!O&_q#^&_q#`&_q#b&_q#c&_q#d&_q#e&_q#f&_q#g&_q#i&_q#m&_q#p&_q&X&_q&Y&_q&i&_q&q&_qY&_q#n&_qt&_qp&_q{&_q~O{,fO~P!GdO!V,jO#|,jO%}RO&P!`O~O!O,mO~O$X,nO%}RO&P!`O~O!a$uO#n,pOp!_X!O!_X~O!O,rO~P7TO!O,rO~P8_O!O,uO~P7TO{,wO~P!L_O!Z#vO#n,xO~O!O,zO~O!a,{O~OY-OOZ$RO_TO`TOaTObTOcTOdTOf[Og[Oo}OviOy|O%{SO%}']O&P'[O~P>qOY-OO!O-PO~O%P'fO%T'gOZ%Uq_%Uq`%Uqa%Uqb%Uqc%Uqd%Uqf%Uqg%Uqo%Uqu%Uqv%Uqy%Uq|%Uq!O%Uq!R%Uq!S%Uq!T%Uq!U%Uq!V%Uq!W%Uq!X%Uq!Y%Uq!Z%Uq![%Uq!]%Uq!^%Uq!t%Uq!y%Uq#d%Uq#p%Uq#r%Uq#s%Uq#w%Uq#x%Uq$U%Uq$W%Uq$^%Uq$a%Uq$c%Uq$i%Uq$k%Uq$p%Uq$r%Uq$u%Uq$x%Uq${%Uq$}%Uq%u%Uq%{%Uq%}%Uq&P%Uq&o%Uq{%Uq$_%Uq$n%Uq~O|'mO~P8_Op-[O{&tX~O{-^O~Op*tOY&ua~O!ari|ri!qri&yri!lri!Yri~P#0OOp-bO{&fX~O{-dO~Ov-eO~Op!`Xs!`X!O!`X!a!`X%{!`X~OZ&OX~P#HOOZUX~P#HOO!O-fO~OZ-gO~OW^yZ#WXs^y!O^y!a^y#[^y#^^y#`^y#b^y#c^y#d^y#e^y#f^y#g^y#i^y#m^y#p^y&X^y&Y^y&i^y&q^yY^y#n^yt^yp^y{^y~OY%^ap%^a!O%^a~P7TO!O#lyY#ly#n#lyp#ly{#ly~P7TOn+kOp!ea&X!ea~Op+mO&X&ha~OZ+zO~PBdO!O-tO~O!l,PO|&ga!O&ga~O|aO!O-wO~OZ^O_TO`TOaTObTOcTOdTOf[Og[Oo.VOuyOv.UOy|O{.QO|aO!OtO!Z!]O!t{O!yzO#dnO#pmO#rnO#snO#w!PO#x!OO$U!QO$W!RO$^!SO$a!TO$c!UO$i!VO$k!WO$p!XO$r!YO$u!ZO$x![O${!^O$}!_O%{SO%}QO&PPO&Y!tO&obO~P>qO|+|O~Op,TO|&ma!O&ma~O|&[Op!{a{!{a!O!{a~OZ#bO|&[Op!{a{!{a!O!{a~O%{SO%}!nO&P!nOp%hX{%hX!O%hX~P>qOp,YO{&na!O&na~O{!|X~P!?{O{.aO~Ot.bO~P!6ZOW$iO!O.cO~OW$iO$O.hO%}RO&P!`O!O&wP~OW$iO$S.iO~O!O.jO~O!a$uO#n.lOp!_X!O!_X~OY.nO~O!O.oO~P7TO#n.pO~P7TO!a.rO~OY.sOZ$RO_TO`TOaTObTOcTOdTOf[Og[Oo}OviOy|O%{SO%}']O&P'[O~P>qOW!vOs&VX%{&VX%}&VX&P&VX&y&VX~O&Y!tO~P$&ZOs0gO%{SO&y.uO%}%RX&P%RX~OY&QXp&QX~P7TO|'mOp%lX{%lX~P8_Op-[O{&ta~O!a.{O~O|'xOp%aX{%aX~P8_Op-bO{&fa~OY/OO~O!O/PO~OZ/QO~O&i%bOp!fa&X!fa~Os0gO%{SO|&ja!O&ja!l&ja~O!O/WO~O!l,PO|&gi!O&gi~Os0gO~PF|O{/]O~P]OW/_O~P3gOW&SXs&SX#^&SX#`&SX#b&SX#c&SX#d&SX#e&SX#f&SX#g&SX#i&SX#m&SX#p&SX&X&SX&Y&SX&i&SX&q&SX~OZ#bO!O&SX~P$*WOW#|OZ#bO&q#`O~Oo/aOv/aO~O|&[Op!{i{!{i!O!{i~O{!|a~P!?{OW$iO!O/cO~OW$iOp/dO!O&wX~OY/hO~P7TOY/jO~OY%Wq!O%Wq~P7TO&y.uO%}%Ra&P%Ra~OY/oO~Os0gO!O/rO!Y/sO%{SO~OY/tO~O&i%bOp!fi&X!fi~Os0gO%{SO|&ji!O&ji!l&ji~O!l,PO|&gq!O&gq~O{/wO~P]Oo/yOv%vOy%tO%}RO&P!`O&Y!tO~O!O/zO~Op/dO!O&wa~O!O0OO~OW$iOp/dO!O&xX~OY0QO~P7TOY0RO~OY%Wy!O%Wy~P7TOs0gO%{SO%}%sa&P%sa&y%sa~OY0SO~Os0gO!O0TO!Y0UO%{SO~Oo0XO%}RO&P!`O~OW({OZ#bO~O!O0ZO~OW$iOp%pa!O%pa~Op/dO!O&xa~O!O0]O~Os0gO!O0]O!Y0^O%{SO~O!O0`O~O!O0aO~O!O0cO~O!O0dO~OYhXY!`Xp!`XvhX&ihX&lhX~P$qOs0hOtqa~P@iO#nUXYUXtUXpUX{UX~P'ROs0hO%{SOt&[X#^&[X#`&[X#b&[X#c&[X#d&[X#e&[X#f&[X#g&[X#i&[X#m&[X#p&[X&X&[X&Y&[X&i&[X~Os0hO%{SOtla#^la#`la#bla#cla#dla#ela#fla#gla#ila#mla#pla&Xla&Yla&ila~Os0lO%{SO~Os0mO%{SO~Ot]i~P!6ZO#^1XO#e1PO#f1YO#g1YO#i0rO#p#mO&X#dO&Y#dOt#]i#`#]i#b#]i#c#]i#m#]i&i#]i~O#d1OO~P$7SO#^1XO#e1PO#f1YO#g1YO#i0rO#p#mO&X#dO&Y#dOt#]i#b#]i#c#]i#m#]i~O#`0|O#d1OO&i#eO~P$8XO#d#]i~P$7SO#f1YO#g1YO#p#mOt#]i#`#]i#b#]i#c#]i#d#]i#e#]i#m#]i&i#]i~O#^1XO#i0rO&X#dO&Y#dO~P$9eOt#ji~P!6ZOt#ly~P!6ZO|aO~PLQO!a0{O#[0{Ot&SX~P#,iO!a0{O#[0{Ot&SX~P$*WO&s1QO~O#^#]i#i#]i&X#]i&Y#]i~P$9eO#n1RO~P7TOY1ZO&i%bO&q#`O~OY1^O&i%bO~O`#e~",
  goto: "#1q&yPPPP&z'_+T.iP'_PP.}/R0vPPPPPP2sPP4l7n:j=f>O@TPPP@ZCQPPPPC}2sPFVPPGQPGwG}PPPPPPPPPPPPIXInPMTM]MgNPNVN]!!^!!c!!c!!lP!!{!$S!$u!%PP!%f!$SP!%l!%v!&V!&_P!&|!'W!'^!$S!'a!'gGwGw!'k!'u!'x2s!)u2s2s!+}P/RP!,RP!,|PPPPPP/RP/R!-q/RPP/RP/RPP/R!/h!/rPP!/x!0RPPPPPPPP&zP&zPP!0V!0V!0j!0VPP!0VP!0VP!0}!1Q!0V!1h!0VP!0VP!1k!1nP!0VP!0VP!1r!0VP!1u!0VP!0V!0VP!0VP!1xP!2O!2R!2XP!0V!2e!2h!2p!3S!7l!7r!8}!9g!9m!9w!:|!;S!;Y!;h!;n!;t!;z!<Q!<W!<^!<d!<j!<p!<v!<|!=S!=^!=d!=n!=tPPP!=z!0V!>oP!BgP!CkP!F]!Fs!Jb2s!L_#!`#%`PP#([#(_P#*z#+Q#,}#-^#-d#.e#.{#/t#/}#0Q#0^P#0a#0mP#0u#0|P#1PP#1YP#1^#1a#1d#1h#1nsrOav!X#V#t$S$T$V$W'S*W*`+|.R.S/^(gpOVW^_amnv!V!X![!^!d!k!o!t!v!x!y#O#S#V#X#_#a#b#f#g#h#i#j#k#l#o#p#q#r#t#z$R$S$T$U$V$W$f$j$t$u$z${%Q%R%Z%[%_%`%b%d%i&U&Z&[&]&^&c&f&g&k&l&n&y&z&|'R'S'^'m'x(R(S(c(g(j)O)R)S)U)Z)])c)n)o)r)w*W*Y*[*]*`*c*f*g*l*q+X+k+m+p+s+v+w+z+|,P,T,Y,[,_,p,r,{-P-T-[-b-s-|.O.P.R.S.`.l.o.r.t.{/O/V/[/^/m/q/s/t0U0W0^0k0n0o0p0q0r0z0{0|0}1O1P1Q1R1V1W1X1Y1Z1^1`1b#pfO^amnv!V!X![!^#V#a#b#f#g#h#i#j#k#l#p#t$R$S$T$U$V$W$u%d&|'R'S'm'x(c)O)R*W*[*]*`*c*q+|,p,r,{-[-b.R.S.l.o.r.{/^/s0U0^1`1bt!bS!O!Q!R!g!i$Q$n*P*Q*R*S,i,k.h.i/d0fQ#WbS%X!y.OQ%l#YU%q#^#|/_Q%x#`W'`$f*g-P.tU'j$i&t*UQ'k$jS(X%R/[U(x%s+e/xQ(}%yQ+W(gQ+c({Q-_*tQ-i+Xq1S#_&k)w0{0|0}1O1P1Q1R1V1W1X1Y1Z1^u!bS!O!Q!R!g!i$Q$n*P*Q*R*S,i,k.h.i/d0fT$k!a(O$eoO^amnv!V!X![!^#V#_#a#b#f#g#h#i#j#k#l#p#t$R$S$T$U$V$W$j$u%d&k&|'R'S'm'x(c)O)R)w*W*[*]*`*c*q+|,p,r,{-[-b.R.S.l.o.r.{/^/s0U0^0{0|0}1O1P1Q1R1V1W1X1Y1Z1^1`1b#rjO^amnv!V!X![!^#V#a#b#f#g#h#i#j#k#l#p#t$R$S$T$U$V$W$j$u%d&|'R'S'm'x(c)O)R*W*[*]*`*c*q+|,p,r,{-[-b.R.S.l.o.r.{/^/s0U0^1`1bW'a$f*g-P.tq1T#_&k)w0{0|0}1O1P1Q1R1V1W1X1Y1Z1^$miO^amnv!V!X![!^#V#_#a#b#f#g#h#i#j#k#l#p#t$R$S$T$U$V$W$f$j$u%d&k&|'R'S'm'x(c)O)R)w*W*[*]*`*c*g*q+|,p,r,{-P-[-b.R.S.l.o.r.t.{/^/s0U0^0{0|0}1O1P1Q1R1V1W1X1Y1Z1^1`1b&hYOV^acmnv|!V!X![!^!t!x!{#O#V#_#a#b#f#g#h#i#j#k#l#o#p#q#t#{$R$S$T$U$V$W$f$j$u$z%[%b%d%i%t&[&]&^&f&g&k&y&{&|'R'S'b'm'x(R(S(c(h(j)O)R)[)])c)o)p)r)w*W*[*]*`*c*g*q+k+w+z+|,P,T,[,p,r,{-P-[-b.P.R.S.`.l.o.r.t.{/^/s0U0^0k0r0z0{0|0}1O1P1Q1R1V1W1X1Y1Z1^1`1bQ%P!vQ(V%QV-R*l-V.u&hYOV^acmnv|!V!X![!^!t!x!{#O#V#_#a#b#f#g#h#i#j#k#l#o#p#q#t#{$R$S$T$U$V$W$f$j$u$z%[%b%d%i%t&[&]&^&f&g&k&y&{&|'R'S'b'm'x(R(S(c(h(j)O)R)[)])c)o)p)r)w*W*[*]*`*c*g*q+k+w+z+|,P,T,[,p,r,{-P-[-b.P.R.S.`.l.o.r.t.{/^/s0U0^0k0r0z0{0|0}1O1P1Q1R1V1W1X1Y1Z1^1`1bV-R*l-V.u&hZOV^acmnv|!V!X![!^!t!x!{#O#V#_#a#b#f#g#h#i#j#k#l#o#p#q#t#{$R$S$T$U$V$W$f$j$u$z%[%b%d%i%t&[&]&^&f&g&k&y&{&|'R'S'b'm'x(R(S(c(h(j)O)R)[)])c)o)p)r)w*W*[*]*`*c*g*q+k+w+z+|,P,T,[,p,r,{-P-[-b.P.R.S.`.l.o.r.t.{/^/s0U0^0k0r0z0{0|0}1O1P1Q1R1V1W1X1Y1Z1^1`1bV-S*l-V.uS!uY-RS#{|%tS%s#^#|Q%y#`Q+e({Q.W+|R/x/_%VXO^amnv!V!X![!^!t#V#_#a#b#f#g#h#i#j#k#l#o#p#t$R$S$T$U$V$W$f$j$u%b%d&]&^&g&k&|'R'S'm'x(R(S(c)O)R)w*W*[*]*`*c*g*q+k+|,P,T,p,r,{-P-[-b.R.S.l.o.r.t.{/^/s0U0^0r0{0|0}1O1P1Q1R1V1W1X1Y1Z1^1`1bQ$}!tR*}(S&i]OV^acmnv!V!X![!^!t!x!{#O#V#_#a#b#f#g#h#i#j#k#l#o#p#q#t$R$S$T$U$V$W$f$j$u$z%[%b%d%i&[&]&^&f&g&k&y&{&|'R'S'b'm'x(R(S(c(h(j)O)R)[)])c)o)p)r)w*W*[*]*`*c*g*l*q+k+w+z+|,P,T,[,p,r,{-P-V-[-b.P.R.S.`.l.o.r.t.u.{/^/s0U0^0k0r0z0{0|0}1O1P1Q1R1V1W1X1Y1Z1^1`1b!u!lW!d!m!o!y#X#r$l$t${%R%Z%_&U&z'^(g)S)Z)n*Y*f+X+p+s+v,_-T-s-|.O/O/V/[/m/q/t0W0i0n0o$liO^amnv!V!X![!^#V#_#a#b#f#g#h#i#j#k#l#p#t$R$S$T$U$V$W$f$j$u%d&k&|'R'S'm'x(c)O)R)w*W*[*]*`*c*g*q+|,p,r,{-P-[-b.R.S.l.o.r.t.{/^/s0U0^0{0|0}1O1P1Q1R1V1W1X1Y1Z1^1`1bQ$S!SQ$T!TQ$Y!WQ$d!]R*d'UQ#cgS&o#z(zQ(w%rQ){&pQ+b(yQ,W)jQ-m+dQ.],XQ/S-nS/`.U.VQ/{/aQ0Y/yR0_0XQ&_#wW(n%m&`&a&bQ)z&oU+[(o(p(qQ,V)jQ,d){S-j+]+^S.[,W,XQ/R-kR/b.]X)a&[)c,[.`rcOav!X#V#t$S$T$V$W'S*W*`+|.R.S/^Y!{^#O%[+z1bQ&{$UW'b$f*g-P.tS(h%i(jW)[&[)c,[.`S)k&c,YS)p&f)rR-V*ld!qW#X&z(g)Z)n*Y+X+s,_Q'|$vQ(Y%VR+R(]#nlOamnv!V!X![!^#V#a#b#f#g#h#i#j#k#l#p#t$R$S$T$U$V$W$j$u%d&|'R'S'm'x(c)O)R*W*[*]*`*c*q+|,p,r,{-[-b.R.S.l.o.r.{/^/s0U0^1`r!oW!y#X$v%V%Z%_&z'^(](g*Y*f+X-U.O.xS#Q^1bQ#wyQ#xzQ#y{Q%m#ZQ%n#[Q%o#]Q(e%eS)T&Z+mY)_&[)[)c,[.`S)j&c,YQ+l)UW+p)Z)n+s,_Q+x)]Q,X)kS-z+v-|q1U#_&k)w0{0|0}1O1P1Q1R1V1W1X1Y1Z1^U'z$u'x-bR)y&nW)a&[)c,[.`T)q&f)rQ&b#wQ&j#yQ(q%mQ(v%oY)]&[)[)c,[.`V)o&f)p)rQ)V&ZR-q+mQ+j)TR-p+l#rlO^amnv!V!X![!^#V#a#b#f#g#h#i#j#k#l#p#t$R$S$T$U$V$W$j$u%d&|'R'S'm'x(c)O)R*W*[*]*`*c*q+|,p,r,{-[-b.R.S.l.o.r.{/^/s0U0^1`1bQ+{)_S-s+p-zQ-}+xq1U#_&k)w0{0|0}1O1P1Q1R1V1W1X1Y1Z1^V#O^+z1bW!|^#O+z1bR(^%[Q,O)`Q-u+rQ-y+uQ/Y-{R/v/ZrtOav!X#V#t$S$T$V$W'S*W*`+|.R.S/^Q$g!_S&X#p1QQ'Y$dQ'i$hW)a&[)c,[.`Q*n'gQ+y)^Q,S)eQ-W*mR-w+trrOav!X#V#t$S$T$V$W'S*W*`+|.R.S/^S(i%i(jW)a&[)c,[.`T)q&f)rQ&a#wS(p%m&bR+^(qQ&`#wQ&e#xU(o%m&a&bQ(s%nS+](p(qR-k+^Q)i&^R)t&gQ&i#yS(u%o&jR+a(vQ&h#yU(t%o&i&jS+`(u(vR-l+aS(i%i(jT)q&f)rrrOav!X#V#t$S$T$V$W'S*W*`+|.R.S/^W)a&[)c,[.`T)q&f)rQ&d#xS(r%n&eR+_(sQ)l&cR.^,YR,^)mQ%j#WR(m%lT(i%i(jQ+})`S-x+u,OR/X-yR.S+|Wi$f*g-P.t#rjO^amnv!V!X![!^#V#a#b#f#g#h#i#j#k#l#p#t$R$S$T$U$V$W$j$u%d&|'R'S'm'x(c)O)R*W*[*]*`*c*q+|,p,r,{-[-b.R.S.l.o.r.{/^/s0U0^1`1bq1T#_&k)w0{0|0}1O1P1Q1R1V1W1X1Y1Z1^$lgO^amnv!V!X![!^#V#_#a#b#f#g#h#i#j#k#l#p#t$R$S$T$U$V$W$f$j$u%d&k&|'R'S'm'x(c)O)R)w*W*[*]*`*c*g*q+|,p,r,{-P-[-b.R.S.l.o.r.t.{/^/s0U0^0{0|0}1O1P1Q1R1V1W1X1Y1Z1^1`1bU%r#^#|/_S(y%s/xQ+d({R-n+eT&m#z&n!W#fk!z$X$b$e%z%{&O&P&Q&R&T&W'l'y*Z*^+g+i,t,y-Y.k.q/i/l1]e1X%w)v,c0s0t0u0v0w0x0y!Q#gk!z$X$b$e%z%{&P&T&W'l'y*Z*^+g+i,t,y-Y.k.q/i/l1]_0|%w)v,c0s0u0x0y#rlO^amnv!V!X![!^#V#a#b#f#g#h#i#j#k#l#p#t$R$S$T$U$V$W$j$u%d&|'R'S'm'x(c)O)R*W*[*]*`*c*q+|,p,r,{-[-b.R.S.l.o.r.{/^/s0U0^1`1bq1U#_&k)w0{0|0}1O1P1Q1R1V1W1X1Y1Z1^a'n$j'm*q-[.{/s0U0^Q'p$jR-`*tQ&r#}Q's$mR*x'tT)|&q)}stOav!X#V#t$S$T$V$W'S*W*`+|.R.S/^ssOav!X#V#t$S$T$V$W'S*W*`+|.R.S/^R$V!UrtOav!X#V#t$S$T$V$W'S*W*`+|.R.S/^R&|$UR$W!UR'T$YT*_'S*`R$^!YR$a!ZX'd$g'e'i*oR*m'fQ-U*lR.x-VQ'h$gQ*k'eQ*p'iR-X*oR$h!_Q'c$fV,}*g-P.tQvOQ#VaW#uv#V.R/^Q.R+|R/^.SrVOav!X#V#t$S$T$V$W'S*W*`+|.R.S/^r!fV!k!x#S#q$z%Q%`&l&y)U+w.P0k0p0q0z^!kW!y#r&U&z'^)SS!x^1bQ#S_#z#qmn!V![!^#_#a#b#f#g#h#i#j#k#l#p$R$j$u%b%d&]&^&g&k&|'R'm'x(R(c)O)R)w*[*]*c*q+k,P,T,p,r,{-[-b.l.o.r.{/s0U0^0{0|0}1O1P1Q1R1V1W1X1Y1Z1^1`S$z!t(SQ%Q!vj%`#O%[%i&[&c&f(j)c)r*l,Y,[.`S&l#z&nY&y$U$f*g-P.tS)U&Z+mS+w)])oQ.P+zQ0k#o![0p!d!o#X$t${%R%Z%_(g)Z)n*Y*f+X+p+s+v,_-T-s-|.O/O/V/[/m/q/t0WS0q0n0oR0z0rQ(T$}R+O(T^!mW!y#r&U&z'^)Sx$l!d#X${%R%Z%_(g)Z)n*Y*f+X+s+v,_-T-|.O/[/m^$s!m$l$t/V/q0W0iS$t!o+pQ/V-sQ/q/OQ0W/tT0i0n0oQ$q!jQ'r0gW'v$q'r'w*wQ'w$rQ*w0lQ/T0hR/u0mQ)P%{R+h)PQ)c&[S,Q)c.`R.`,[!n`O^av!X#O#V#t$S$T$U$V$W$f%[%i&[&c&f'S(j)c)r*W*`*g*l+z+|,Y,[-P.R.S.`.t/^1bY!eV!x%`&y.PT#T`!eQ-c*yR.}-cQ$w!qR'}$wQ%c#PU(b%c/U1aQ/U-oR1a1_Q+n)VR-r+nQ%]!|R(_%]Q,U)gR.Z,UQ)r&fR,`)rQ,Z)lR._,ZQ(j%iR+Y(jQ&n#zR)x&nQ%f#QR(f%fQ-]*rR.z-]Q*u'pR-a*uQ)}&qR,e)}Q,i*PR.e,iQ/e.fS/}/e0PR0P/gQ*`'SR,v*`Q'e$gS*j'e*oR*o'iQ.v-TR/n.vQ*h'cR-Q*h`uOav#V+|.R.S/^Q$Z!XQ&Y#tQ&w$SQ&x$TQ'O$VQ'P$WS*_'S*`R,o*W(UqOVW^_amnv!V!X![!^!d!k!o!t!v!x!y#O#S#V#X#_#a#b#f#g#h#i#j#k#l#o#p#q#r#t#z$R$S$T$U$V$W$f$t$u$z${%Q%R%Z%[%_%`%b%d%i&U&Z&[&]&^&c&f&g&k&l&n&y&z&|'R'S'^'x(R(S(c(g(j)O)R)S)U)Z)])c)n)o)r)w*W*Y*[*]*`*c*f*g*l+X+k+m+p+s+v+w+z+|,P,T,Y,[,_,p,r,{-P-T-b-s-|.O.P.R.S.`.l.o.r.t/O/V/[/^/m/q/t0W0k0n0o0p0q0r0z0{0|0}1O1P1Q1R1V1W1X1Y1Z1^1`1ba'o$j'm*q-[.{/s0U0^Q!aSQ#}!OQ$O!QQ$P!RQ$m!gQ$o!iQ&v$QQ't$nQ(O0fS,g*P*RQ,k*QQ,l*SQ.d,iS.f,k.hQ/g.iR/|/d&_ROS^abmnv!O!Q!R!V!X![!^!g!i!y#V#Y#^#_#`#a#b#f#g#h#i#j#k#l#p#t#|$Q$R$S$T$U$V$W$f$i$j$n$u%R%d%s%y&k&t&|'R'S'm'x(c(g({)O)R)w*P*Q*R*S*U*W*[*]*`*c*g*q*t+X+e+|,i,k,p,r,{-P-[-b.O.R.S.h.i.l.o.r.t.{/[/^/_/d/s/x0U0^0f0{0|0}1O1P1Q1R1V1W1X1Y1Z1^1`1bQ'q$jQ*r'mS-Z*q.{Q.y-[Q0V/sQ0[0UR0b0^rkOav!X#V#t$S$T$V$W'S*W*`+|.R.S/^U!z^$R1bS#mm1VS#sn1WQ$X!VQ$b![Q$e!^Q%w#_Q%z#aY%{#b$U*[,r.oQ%}#fQ&O#gQ&P#hQ&Q#iQ&R#jS&S#k1YQ&T#lQ&W#p^'l$j'm-[.{/s0U0^U'y$u'x-bS(d%d1ZQ)v&kQ*Z&|Q*^'RS+S(c1^Q+g)OQ+i)RQ,c)wQ,t*]Q,y*cQ-Y*qQ.k,pQ.q,{Q/i.lQ/l.rQ0s0{Q0t0|Q0u0}Q0v1OQ0w1PQ0x1QQ0y1RQ1[1XR1]1`$beO^amnv!V!X![!^#V#_#a#b#f#g#h#i#j#k#l#p#t$R$S$T$U$V$W$j$u%d&k&|'R'S'm'x(c)O)R)w*W*[*]*`*c*q,p,r,{-[-b.R.S.l.o.r.{/^/s0U0^0{0|0}1O1P1Q1R1V1W1X1Y1Z1^1`1bW'_$f*g-P.tR.T+|rWOav!X#V#t$S$T$V$W'S*W*`+|.R.S/^W!dV#q$z&yS!y^1bQ#Xc#j#rmn!V![!^#_#a#b#f#g#h#i#j#k#l#p$R$j$u%d&k&|'R'm'x(c)O)R)w*[*]*c*q,p,r,{-[-b.l.o.r.{/s0U0^0{0|0}1O1P1Q1R1V1W1X1Y1Z1^1`d${!t%b&]&^&g(R(S+k,P,TQ%R!xQ%Z!{S%_#O%[Q&U#oQ&z$UW'^$f*g-P.tS(g%i(jQ)S0kW)Z&[)c,[.`S)n&f)rQ*Y&{Q*f'bQ+X(hQ+s)[S+v)])oQ,_)pS-T*l-VQ-|+wQ.O+zQ/[.PQ/m.uQ0n0rR0o0z&h]OV^acmnv!V!X![!^!t!x!{#O#V#_#a#b#f#g#h#i#j#k#l#o#p#q#t$R$S$T$U$V$W$f$j$u$z%[%b%d%i&[&]&^&f&g&k&y&{&|'R'S'b'm'x(R(S(c(h(j)O)R)[)])c)o)p)r)w*W*[*]*`*c*g*l*q+k+w+z+|,P,T,[,p,r,{-P-V-[-b.P.R.S.`.l.o.r.t.u.{/^/s0U0^0k0r0z0{0|0}1O1P1Q1R1V1W1X1Y1Z1^1`1bQ#z|Q&p#{R(z%t&sUOV^acmnv|!V!X![!^!t!v!x!{#O#V#_#a#b#f#g#h#i#j#k#l#o#p#q#t#{$R$S$T$U$V$W$f$j$u$z%Q%[%b%d%i%t&[&]&^&f&g&k&y&{&|'R'S'b'm'x(R(S(c(h(j)O)R)[)])c)o)p)r)w*W*[*]*`*c*g*l*q+k+w+z+|,P,T,[,p,r,{-P-V-[-b.P.R.S.`.l.o.r.t.u.{/^/s0U0^0k0r0z0{0|0}1O1P1Q1R1V1W1X1Y1Z1^1`1bR%O!t$hhOamnv!V!X![!^#V#_#a#b#f#g#h#i#j#k#l#p#t$R$S$T$U$V$W$f$j$u%d&k&|'R'S'm'x(c)O)R)w*W*[*]*`*c*g*q+|,p,r,{-P-[-b.R.S.l.o.r.t.{/^/s0U0^0{0|0}1O1P1Q1R1V1W1X1Y1Z1^1`Q#P^Q$}!tS&V#o0rQ(a%bQ)f&]U)g&^&g,PQ*|(RQ*}(SQ-o+kQ.Y,TR1_1bQ(Q$|R*{(P$ldO^amnv!V!X![!^#V#_#a#b#f#g#h#i#j#k#l#p#t$R$S$T$U$V$W$f$j$u%d&k&|'R'S'm'x(c)O)R)w*W*[*]*`*c*g*q+|,p,r,{-P-[-b.R.S.l.o.r.t.{/^/s0U0^0{0|0}1O1P1Q1R1V1W1X1Y1Z1^1`1bT%p#^/_Q%|#bQ&}$UQ,s*[Q.m,rR/k.oX)b&[)c,[.`!}_OV^`av!X!e!x#O#V#t$S$T$U$V$W$f%[%`%i&[&c&f&y'S(j)c)r*W*`*g*l+z+|,Y,[-P.P.R.S.`.t/^1bS!rW&zS%k#X*YS+V(g)nQ+q)ZS-h+X,_R-v+sf!pW#X$v%V(](g)Z)n+X+s,_U%U!y%_.OQ([%ZQ*X&zQ*e'^Q,q*YQ,|*fQ.w-UR/p.xQ'{$uQ*y'xR.|-bR*z'x[)^&[&f)c)r,[.`T+t)[)pR)W&ZW+r)Z)n+s,_Q-{+vR/Z-|U!}^+z1bR%a#OS)h&^&gR.X,PR)m&cW)`&[)c,[.`R+u)[T#R^1bR*s'mR'q$jT,h*P,iQ.g,kR/f.hR/f.i",
  nodeNames: "\u26A0 LineComment BlockComment Program ModuleDeclaration MarkerAnnotation Identifier ScopedIdentifier . Annotation ) ( AnnotationArgumentList AssignmentExpression FieldAccess IntegerLiteral FloatingPointLiteral BooleanLiteral CharacterLiteral StringLiteral null ClassLiteral void PrimitiveType TypeName ScopedTypeName GenericType TypeArguments AnnotatedType Wildcard extends super , ArrayType Dimension [ ] class this ParenthesizedExpression ObjectCreationExpression new ArgumentList } { ClassBody ; FieldDeclaration Modifiers public protected private abstract static final strictfp default synchronized native transient volatile VariableDeclarator Definition AssignOp ArrayInitializer MethodDeclaration TypeParameters TypeParameter TypeBound FormalParameters ReceiverParameter FormalParameter SpreadParameter Throws throws Block ClassDeclaration Superclass SuperInterfaces implements InterfaceTypeList InterfaceDeclaration interface ExtendsInterfaces InterfaceBody ConstantDeclaration EnumDeclaration enum EnumBody EnumConstant EnumBodyDeclarations AnnotationTypeDeclaration AnnotationTypeBody AnnotationTypeElementDeclaration StaticInitializer ConstructorDeclaration ConstructorBody ExplicitConstructorInvocation ArrayAccess MethodInvocation MethodName MethodReference ArrayCreationExpression Dimension AssignOp BinaryExpression CompareOp CompareOp LogicOp BitOp BitOp LogicOp ArithOp ArithOp ArithOp BitOp InstanceofExpression instanceof LambdaExpression InferredParameters TernaryExpression LogicOp : UpdateExpression UpdateOp UnaryExpression LogicOp BitOp CastExpression ElementValueArrayInitializer ElementValuePair open module ModuleBody ModuleDirective requires transitive exports to opens uses provides with PackageDeclaration package ImportDeclaration import Asterisk ExpressionStatement LabeledStatement Label IfStatement if else WhileStatement while ForStatement for ForSpec LocalVariableDeclaration EnhancedForStatement ForSpec AssertStatement assert SwitchStatement switch SwitchBlock SwitchLabel case DoStatement do BreakStatement break Label ContinueStatement continue Label ReturnStatement return SynchronizedStatement ThrowStatement throw TryStatement try CatchClause catch CatchFormalParameter CatchType FinallyClause finally TryWithResourcesStatement ResourceSpecification Resource",
  maxTerm: 271,
  nodeProps: [
    [NodeProp.group, -26, 4, 46, 75, 76, 81, 86, 91, 143, 145, 148, 149, 151, 154, 156, 159, 160, 162, 164, 169, 171, 174, 177, 179, 180, 182, 190, "Statement", -24, 6, 13, 14, 15, 16, 17, 18, 19, 20, 21, 38, 39, 40, 98, 99, 101, 102, 105, 116, 118, 120, 123, 125, 128, "Expression", -7, 22, 23, 24, 25, 26, 28, 33, "Type"],
    [NodeProp.openedBy, 10, "(", 43, "{"],
    [NodeProp.closedBy, 11, ")", 44, "}"]
  ],
  skippedNodes: [0, 1, 2],
  repeatNodeCount: 28,
  tokenData: "Cr~R{X^#xpq#xqr$mrs$ztu%ruv&Wvw&ewx&uxy(]yz(bz{(g{|(q|})R}!O)W!O!P)k!P!Q-S!Q!R.b!R![3S![!]?[!]!^?i!^!_?n!_!`@R!`!a@Z!a!b@q!b!c@x!c!}BX!}#OBm#P#QBr#Q#RBw#R#S%r#T#o%r#o#pCP#p#qCU#q#rCh#r#sCm#y#z#x$f$g#x#BY#BZ#x$IS$I_#x$I|$JO#x$JT$JU#x$KV$KW#x&FU&FV#x~#}Y%w~X^#xpq#x#y#z#x$f$g#x#BY#BZ#x$IS$I_#x$I|$JO#x$JT$JU#x$KV$KW#x&FU&FV#xR$rP#rP!_!`$uQ$zO#^Q~$}UOY$zZr$zrs%as#O$z#O#P%f#P~$z~%fOc~~%iROY$zYZ$zZ~$z~%wT%}~tu%r!Q![%r!c!}%r#R#S%r#T#o%r~&]P#f~!_!`&`Q&eO#[Q~&jQ&i~vw&p!_!`&`~&uO#`~~&xTOY'XZw'Xx#O'X#O#P(P#P~'X~'[UOY'XZw'Xwx'nx#O'X#O#P's#P~'X~'sOb~~'vROY'XYZ'XZ~'X~(SROY'XYZ'XZ~'X~(bOZ~~(gOY~R(nP$XP#eQ!_!`&`~(vQ#d~{|(|!_!`&`~)RO#p~~)WOp~~)]R#d~}!O(|!_!`&`!`!a)f~)kO&s~~)pQWU!O!P)v!Q![*R~)yP!O!P)|~*RO&l~P*WW`P!Q![*R!f!g*p!g!h*u!h!i*p#R#S,_#W#X*p#X#Y*u#Y#Z*pP*uO`PP*xR{|+R}!O+R!Q![+XP+UP!Q![+XP+^U`P!Q![+X!f!g*p!h!i*p#R#S+p#W#X*p#Y#Z*pP+sP!Q![+vP+{U`P!Q![+v!f!g*p!h!i*p#R#S+p#W#X*p#Y#Z*pP,bP!Q![,eP,jW`P!Q![,e!f!g*p!g!h*u!h!i*p#R#S,_#W#X*p#X#Y*u#Y#Z*p~-XR#eQz{-b!P!Q.V!_!`&`~-eROz-bz{-n{~-b~-qTOz-bz{-n{!P-b!P!Q.Q!Q~-b~.VOQ~~.[QP~OY.VZ~.V~.ga_~!O!P/l!Q![3S!d!e6g!f!g*p!g!h3z!h!i*p!n!o5d!q!r7s!z!{8s#R#S5i#U#V6g#W#X*p#X#Y3z#Y#Z*p#`#a5d#c#d7s#l#m8sP/qV`P!Q![0W!f!g*p!g!h0u!h!i*p#W#X*p#X#Y0u#Y#Z*pP0]W`P!Q![0W!f!g*p!g!h0u!h!i*p#R#S2_#W#X*p#X#Y0u#Y#Z*pP0xR{|1R}!O1R!Q![1XP1UP!Q![1XP1^U`P!Q![1X!f!g*p!h!i*p#R#S1p#W#X*p#Y#Z*pP1sP!Q![1vP1{U`P!Q![1v!f!g*p!h!i*p#R#S1p#W#X*p#Y#Z*pP2bP!Q![2eP2jW`P!Q![2e!f!g*p!g!h0u!h!i*p#R#S2_#W#X*p#X#Y0u#Y#Z*p~3XZ_~!O!P/l!Q![3S!f!g*p!g!h3z!h!i*p!n!o5d#R#S5i#W#X*p#X#Y3z#Y#Z*p#`#a5dP3}R{|4W}!O4W!Q![4^P4ZP!Q![4^P4cU`P!Q![4^!f!g*p!h!i*p#R#S4u#W#X*p#Y#Z*pP4xP!Q![4{P5QU`P!Q![4{!f!g*p!h!i*p#R#S4u#W#X*p#Y#Z*p~5iO_~~5lP!Q![5o~5tZ_~!O!P/l!Q![5o!f!g*p!g!h3z!h!i*p!n!o5d#R#S5i#W#X*p#X#Y3z#Y#Z*p#`#a5d~6jQ!Q!R6p!R!S6p~6uT_~!Q!R6p!R!S6p!n!o5d#R#S7U#`#a5d~7XQ!Q!R7_!R!S7_~7dT_~!Q!R7_!R!S7_!n!o5d#R#S7U#`#a5d~7vP!Q!Y7y~8OS_~!Q!Y7y!n!o5d#R#S8[#`#a5d~8_P!Q!Y8b~8gS_~!Q!Y8b!n!o5d#R#S8[#`#a5d~8vS!O!P9S!Q![<Q!c!i<Q#T#Z<QP9VR!Q![9`!c!i9`#T#Z9`P9cU!Q![9`!c!i9`!r!s9u#R#S;_#T#Z9`#d#e9uP9xR{|:R}!O:R!Q![:XP:UP!Q![:XP:^U`P!Q![:X!f!g*p!h!i*p#R#S:p#W#X*p#Y#Z*pP:sP!Q![:vP:{U`P!Q![:v!f!g*p!h!i*p#R#S:p#W#X*p#Y#Z*pP;bR!Q![;k!c!i;k#T#Z;kP;nU!Q![;k!c!i;k!r!s9u#R#S;_#T#Z;k#d#e9u~<VX_~!O!P<r!Q![<Q!c!i<Q!n!o5d!r!s9u#R#S>^#T#Z<Q#`#a5d#d#e9uP<uT!Q![=U!c!i=U!r!s9u#T#Z=U#d#e9uP=XU!Q![=U!c!i=U!r!s9u#R#S=k#T#Z=U#d#e9uP=nR!Q![=w!c!i=w#T#Z=wP=zU!Q![=w!c!i=w!r!s9u#R#S=k#T#Z=w#d#e9u~>aR!Q![>j!c!i>j#T#Z>j~>oX_~!O!P<r!Q![>j!c!i>j!n!o5d!r!s9u#R#S>^#T#Z>j#`#a5d#d#e9u~?aP#n~![!]?d~?iO&q~~?nO!O~~?sQ&Y~!^!_?y!_!`$u~@OP#g~!_!`&`~@WP!a~!_!`$u~@`Q&X~!_!`$u!`!a@f~@kQ#g~!_!`&`!`!a?yV@xO&]T#mQ~@}P%{~#]#^AQ~ATP#b#cAW~AZP#h#iA^~AaP#X#YAd~AgP#f#gAj~AmP#Y#ZAp~AsP#T#UAv~AyP#V#WA|~BPP#X#YBS~BXO&o~~B^T&P~tuBX!Q![BX!c!}BX#R#SBX#T#oBX~BrOs~~BwOt~QB|P#bQ!_!`&`~CUO|~VC]Q&yT#bQ!_!`&`#p#qCcQChO#cQ~CmO{~~CrO#s~",
  tokenizers: [0, 1, 2],
  topRules: {Program: [0, 3]},
  dynamicPrecedences: {"26": 1, "230": -1, "238": -1},
  specialized: [{term: 229, get: (value) => spec_identifier[value] || -1}],
  tokenPrec: 7618
});

// node_modules/@codemirror/lang-java/dist/index.js
var javaLanguage = LezerLanguage.define({
  parser: parser.configure({
    props: [
      indentNodeProp.add({
        IfStatement: continuedIndent({except: /^\s*({|else\b)/}),
        TryStatement: continuedIndent({except: /^\s*({|catch|finally)\b/}),
        LabeledStatement: flatIndent,
        SwitchBlock: (context) => {
          let after = context.textAfter, closed = /^\s*\}/.test(after), isCase = /^\s*(case|default)\b/.test(after);
          return context.baseIndent + (closed ? 0 : isCase ? 1 : 2) * context.unit;
        },
        BlockComment: () => -1,
        Statement: continuedIndent({except: /^{/})
      }),
      foldNodeProp.add({
        ["Block SwitchBlock ClassBody ElementValueArrayInitializer ModuleBody EnumBody ConstructorBody InterfaceBody ArrayInitializer"]: foldInside,
        BlockComment(tree) {
          return {from: tree.from + 2, to: tree.to - 2};
        }
      }),
      styleTags({
        null: tags.null,
        instanceof: tags.operatorKeyword,
        this: tags.self,
        "new super assert open to with void": tags.keyword,
        "class interface extends implements module package import enum": tags.definitionKeyword,
        "switch while for if else case default do break continue return try catch finally throw": tags.controlKeyword,
        ["requires exports opens uses provides public private protected static transitive abstract final strictfp synchronized native transient volatile throws"]: tags.modifier,
        IntegerLiteral: tags.integer,
        FloatLiteral: tags.float,
        StringLiteral: tags.string,
        CharacterLiteral: tags.character,
        LineComment: tags.lineComment,
        BlockComment: tags.blockComment,
        BooleanLiteral: tags.bool,
        PrimitiveType: tags.standard(tags.typeName),
        TypeName: tags.typeName,
        Identifier: tags.variableName,
        "MethodName/Identifier": tags.function(tags.variableName),
        Definition: tags.definition(tags.variableName),
        ArithOp: tags.arithmeticOperator,
        LogicOp: tags.logicOperator,
        BitOp: tags.bitwiseOperator,
        CompareOp: tags.compareOperator,
        AssignOp: tags.definitionOperator,
        UpdateOp: tags.updateOperator,
        Asterisk: tags.punctuation,
        Label: tags.labelName,
        "( )": tags.paren,
        "[ ]": tags.squareBracket,
        "{ }": tags.brace,
        ".": tags.derefOperator,
        ", ;": tags.separator
      })
    ]
  }),
  languageData: {
    commentTokens: {line: "//", block: {open: "/*", close: "*/"}},
    indentOnInput: /^\s*(?:case |default:|\{|\})$/
  }
});
function java2() {
  return new LanguageSupport(javaLanguage);
}

// node_modules/lezer-javascript/dist/index.es.js
var noSemi = 269;
var incdec = 1;
var incdecPrefix = 2;
var templateContent = 270;
var templateDollarBrace = 271;
var templateEnd = 272;
var insertSemi = 273;
var TSExtends = 3;
var Dialect_ts = 1;
var newline = [10, 13, 8232, 8233];
var space = [9, 11, 12, 32, 133, 160, 5760, 8192, 8193, 8194, 8195, 8196, 8197, 8198, 8199, 8200, 8201, 8202, 8239, 8287, 12288];
var braceR = 125;
var braceL = 123;
var semicolon = 59;
var slash = 47;
var star = 42;
var plus = 43;
var minus = 45;
var dollar = 36;
var backtick = 96;
var backslash = 92;
function newlineBefore(input, pos) {
  for (let i = pos - 1; i >= 0; i--) {
    let prev = input.get(i);
    if (newline.indexOf(prev) > -1)
      return true;
    if (space.indexOf(prev) < 0)
      break;
  }
  return false;
}
var insertSemicolon = new ExternalTokenizer((input, token, stack) => {
  let pos = token.start, next = input.get(pos);
  if ((next == braceR || next == -1 || newlineBefore(input, pos)) && stack.canShift(insertSemi))
    token.accept(insertSemi, token.start);
}, {contextual: true, fallback: true});
var noSemicolon = new ExternalTokenizer((input, token, stack) => {
  let pos = token.start, next = input.get(pos++);
  if (space.indexOf(next) > -1 || newline.indexOf(next) > -1)
    return;
  if (next == slash) {
    let after = input.get(pos++);
    if (after == slash || after == star)
      return;
  }
  if (next != braceR && next != semicolon && next != -1 && !newlineBefore(input, token.start) && stack.canShift(noSemi))
    token.accept(noSemi, token.start);
}, {contextual: true});
var incdecToken = new ExternalTokenizer((input, token, stack) => {
  let pos = token.start, next = input.get(pos);
  if ((next == plus || next == minus) && next == input.get(pos + 1)) {
    let mayPostfix = !newlineBefore(input, token.start) && stack.canShift(incdec);
    token.accept(mayPostfix ? incdec : incdecPrefix, pos + 2);
  }
}, {contextual: true});
var template = new ExternalTokenizer((input, token) => {
  let pos = token.start, afterDollar = false;
  for (; ; ) {
    let next = input.get(pos++);
    if (next < 0) {
      if (pos - 1 > token.start)
        token.accept(templateContent, pos - 1);
      break;
    } else if (next == backtick) {
      if (pos == token.start + 1)
        token.accept(templateEnd, pos);
      else
        token.accept(templateContent, pos - 1);
      break;
    } else if (next == braceL && afterDollar) {
      if (pos == token.start + 2)
        token.accept(templateDollarBrace, pos);
      else
        token.accept(templateContent, pos - 2);
      break;
    } else if (next == 10 && pos > token.start + 1) {
      token.accept(templateContent, pos);
      break;
    } else if (next == backslash && pos != input.length) {
      pos++;
    }
    afterDollar = next == dollar;
  }
});
function tsExtends(value, stack) {
  return value == "extends" && stack.dialectEnabled(Dialect_ts) ? TSExtends : -1;
}
var spec_identifier2 = {__proto__: null, export: 16, as: 21, from: 25, default: 30, async: 35, function: 36, this: 46, true: 54, false: 54, void: 58, typeof: 62, null: 76, super: 78, new: 112, await: 129, yield: 131, delete: 132, class: 142, extends: 144, public: 181, private: 181, protected: 181, readonly: 183, in: 202, instanceof: 204, import: 236, keyof: 287, unique: 291, infer: 297, is: 331, abstract: 351, implements: 353, type: 355, let: 358, var: 360, const: 362, interface: 369, enum: 373, namespace: 379, module: 381, declare: 385, global: 389, for: 410, of: 419, while: 422, with: 426, do: 430, if: 434, else: 436, switch: 440, case: 446, try: 452, catch: 454, finally: 456, return: 460, throw: 464, break: 468, continue: 472, debugger: 476};
var spec_word = {__proto__: null, async: 99, get: 101, set: 103, public: 151, private: 151, protected: 151, static: 153, abstract: 155, readonly: 159, new: 335};
var spec_LessThan = {__proto__: null, "<": 119};
var parser2 = Parser.deserialize({
  version: 13,
  states: "$8xO]QYOOO&zQ!LdO'#CgO'ROSO'#DRO)ZQYO'#DWO)kQYO'#DcO)rQYO'#DmO-iQYO'#DsOOQO'#ET'#ETO-|QWO'#ESO.RQWO'#ESO.ZQ!LdO'#IgO2dQ!LdO'#IhO3QQWO'#EpO3VQpO'#FVOOQ!LS'#Ex'#ExO3_O!bO'#ExO3mQWO'#F^O4wQWO'#F]OOQ!LS'#Ih'#IhOOQ!LS'#Ig'#IgOOQQ'#JR'#JRO4|QWO'#HeO5RQ!LYO'#HfOOQQ'#I['#I[OOQQ'#Hg'#HgQ]QYOOO)rQYO'#DeO5ZQWO'#GQO5`Q#tO'#ClO5nQWO'#ERO5yQ#tO'#EwO6eQWO'#GQO6jQWO'#GUO6uQWO'#GUO7TQWO'#GYO7TQWO'#GZO7TQWO'#G]O5ZQWO'#G`O7tQWO'#GcO9SQWO'#CcO9dQWO'#GpO9lQWO'#GvO9lQWO'#GxO]QYO'#GzO9lQWO'#G|O9lQWO'#HPO9qQWO'#HVO9vQ!LZO'#HZO)rQYO'#H]O:RQ!LZO'#H_O:^Q!LZO'#HaO5RQ!LYO'#HcO)rQYO'#IjOOOS'#Hh'#HhO:iOSO,59mOOQ!LS,59m,59mO<zQbO'#CgO=UQYO'#HiO=cQWO'#IlO?bQbO'#IlO'^QYO'#IlO?iQWO,59rO@PQ&jO'#D]O@xQWO'#ETOAVQWO'#IvOAbQWO'#IuOAjQWO,5:qOAoQWO'#ItOAvQWO'#DtO5`Q#tO'#EROBUQWO'#EROBaQ`O'#EwOOQ!LS,59},59}OBiQYO,59}ODgQ!LdO,5:XOETQWO,5:_OEnQ!LYO'#IsO6jQWO'#IrOEuQWO'#IrOE}QWO,5:pOFSQWO'#IrOFbQYO,5:nOH_QWO'#EPOIfQWO,5:nOJrQWO'#DgOJyQYO'#DlOKTQ&jO,5:wO)rQYO,5:wOOQQ'#Eh'#EhOOQQ'#Ej'#EjO)rQYO,5:xO)rQYO,5:xO)rQYO,5:xO)rQYO,5:xO)rQYO,5:xO)rQYO,5:xO)rQYO,5:xO)rQYO,5:xO)rQYO,5:xO)rQYO,5:xO)rQYO,5:xOOQQ'#En'#EnOKYQYO,5;XOOQ!LS,5;^,5;^OOQ!LS,5;_,5;_OMVQWO,5;_OOQ!LS,5;`,5;`O)rQYO'#HsOM[Q!LYO,5;yOH_QWO,5:xO)rQYO,5;[ONXQpO'#IzOMvQpO'#IzON`QpO'#IzONqQpO,5;gOOQO,5;q,5;qO!!|QYO'#FXOOOO'#Hr'#HrO3_O!bO,5;dO!#TQpO'#FZOOQ!LS,5;d,5;dO!#qQ,UO'#CqOOQ!LS'#Ct'#CtO!$UQWO'#CtO!$lQ#tO,5;vO!$sQWO,5;xO!%|QWO'#FhO!&ZQWO'#FiO!&`QWO'#FmO!'bQ&jO'#FqO!(TQ,UO'#IeOOQ!LS'#Ie'#IeO!(_QWO'#IdO!(mQWO'#IcOOQ!LS'#Cr'#CrOOQ!LS'#Cx'#CxO!(uQWO'#CzOIkQWO'#F`OIkQWO'#FbO!(zQWO'#FdOIaQWO'#FeO!)PQWO'#FkOIkQWO'#FpO!)UQWO'#EUO!)mQWO,5;wO]QYO,5>POOQQ'#I_'#I_OOQQ,5>Q,5>QOOQQ-E;e-E;eO!+iQ!LdO,5:POOQ!LQ'#Co'#CoO!,YQ#tO,5<lOOQO'#Ce'#CeO!,kQWO'#CpO!,sQ!LYO'#I`O4wQWO'#I`O9qQWO,59WO!-RQpO,59WO!-ZQ#tO,59WO5`Q#tO,59WO!-fQWO,5:nO!-nQWO'#GoO!-vQWO'#JVO!.OQYO,5;aOKTQ&jO,5;cO!/{QWO,5=YO!0QQWO,5=YO!0VQWO,5=YO5RQ!LYO,5=YO5ZQWO,5<lO!0eQWO'#EVO!0vQ&jO'#EWOOQ!LQ'#It'#ItO!1XQ!LYO'#JSO5RQ!LYO,5<pO7TQWO,5<wOOQO'#Cq'#CqO!1dQpO,5<tO!1lQ#tO,5<uO!1wQWO,5<wO!1|Q`O,5<zO9qQWO'#GeO5ZQWO'#GgO!2UQWO'#GgO5`Q#tO'#GjO!2ZQWO'#GjOOQQ,5<},5<}O!2`QWO'#GkO!2hQWO'#ClO!2mQWO,58}O!2wQWO,58}O!4vQYO,58}OOQQ,58},58}O!5TQ!LYO,58}O)rQYO,58}O!5`QYO'#GrOOQQ'#Gs'#GsOOQQ'#Gt'#GtO]QYO,5=[O!5pQWO,5=[O)rQYO'#DsO]QYO,5=bO]QYO,5=dO!5uQWO,5=fO]QYO,5=hO!5zQWO,5=kO!6PQYO,5=qOOQQ,5=u,5=uO)rQYO,5=uO5RQ!LYO,5=wOOQQ,5=y,5=yO!9}QWO,5=yOOQQ,5={,5={O!9}QWO,5={OOQQ,5=},5=}O!:SQ`O,5?UOOOS-E;f-E;fOOQ!LS1G/X1G/XO!:XQbO,5>TO)rQYO,5>TOOQO-E;g-E;gO!:cQWO,5?WO!:kQbO,5?WO!:rQWO,5?aOOQ!LS1G/^1G/^O!:zQpO'#DPOOQO'#In'#InO)rQYO'#InO!;iQpO'#InO!<WQpO'#D^O!<iQ&jO'#D^O!>qQYO'#D^O!>xQWO'#ImO!?QQWO,59wO!?VQWO'#EXO!?eQWO'#IwO!?mQWO,5:rO!@TQ&jO'#D^O)rQYO,5?bO!@_QWO'#HnO!:rQWO,5?aOOQ!LQ1G0]1G0]O!AeQ&jO'#DwOOQ!LS,5:`,5:`O)rQYO,5:`OH_QWO,5:`O!AlQWO,5:`O9qQWO,5:mO!-RQpO,5:mO!-ZQ#tO,5:mO5`Q#tO,5:mOOQ!LS1G/i1G/iOOQ!LS1G/y1G/yOOQ!LQ'#EO'#EOO)rQYO,5?_O!AwQ!LYO,5?_O!BYQ!LYO,5?_O!BaQWO,5?^O!BiQWO'#HpO!BaQWO,5?^OOQ!LQ1G0[1G0[O6jQWO,5?^OOQ!LS1G0Y1G0YO!CTQ!LbO,5:kOOQ!LS'#Fg'#FgO!CqQ!LdO'#IeOFbQYO1G0YO!EpQ#tO'#IoO!EzQWO,5:RO!FPQbO'#IpO)rQYO'#IpO!FZQWO,5:WOOQ!LS'#DP'#DPOOQ!LS1G0c1G0cO!F`QWO1G0cO!HqQ!LdO1G0dO!HxQ!LdO1G0dO!K]Q!LdO1G0dO!KdQ!LdO1G0dO!MkQ!LdO1G0dO!NOQ!LdO1G0dO#!oQ!LdO1G0dO#!vQ!LdO1G0dO#%ZQ!LdO1G0dO#%bQ!LdO1G0dO#'VQ!LdO1G0dO#*PQ7^O'#CgO#+zQ7^O1G0sO#-xQ7^O'#IhOOQ!LS1G0y1G0yO#.SQ!LdO,5>_OOQ!LS-E;q-E;qO#.sQ!LdO1G0dO#0uQ!LdO1G0vO#1fQpO,5;iO#1kQpO,5;jO#1pQpO'#FQO#2UQWO'#FPOOQO'#I{'#I{OOQO'#Hq'#HqO#2ZQpO1G1ROOQ!LS1G1R1G1ROOQO1G1[1G1[O#2iQ7^O'#IgO#4cQWO,5;sO! PQYO,5;sOOOO-E;p-E;pOOQ!LS1G1O1G1OOOQ!LS,5;u,5;uO#4hQpO,5;uOOQ!LS,59`,59`O)rQYO1G1bOKTQ&jO'#HuO#4mQWO,5<ZOOQ!LS,5<W,5<WOOQO'#F{'#F{OIkQWO,5<fOOQO'#F}'#F}OIkQWO,5<hOIkQWO,5<jOOQO1G1d1G1dO#4xQ`O'#CoO#5]Q`O,5<SO#5dQWO'#JOO5ZQWO'#JOO#5rQWO,5<UOIkQWO,5<TO#5wQ`O'#FgO#6UQ`O'#JPO#6`QWO'#JPOH_QWO'#JPO#6eQWO,5<XOOQ!LQ'#Db'#DbO#6jQWO'#FjO#6uQpO'#FrO!']Q&jO'#FrO!']Q&jO'#FtO#7WQWO'#FuO!)PQWO'#FxOOQO'#Hw'#HwO#7]Q&jO,5<]OOQ!LS,5<],5<]O#7dQ&jO'#FrO#7rQ&jO'#FsO#7zQ&jO'#FsOOQ!LS,5<k,5<kOIkQWO,5?OOIkQWO,5?OO#8PQWO'#HxO#8[QWO,5>}OOQ!LS'#Cg'#CgO#9OQ#tO,59fOOQ!LS,59f,59fO#9qQ#tO,5;zO#:dQ#tO,5;|O#:nQWO,5<OOOQ!LS,5<P,5<PO#:sQWO,5<VO#:xQ#tO,5<[OFbQYO1G1cO#;YQWO1G1cOOQQ1G3k1G3kOOQ!LS1G/k1G/kOMVQWO1G/kOOQQ1G2W1G2WOH_QWO1G2WO)rQYO1G2WOH_QWO1G2WO#;_QWO1G2WO#;mQWO,59[O#<sQWO'#EPOOQ!LQ,5>z,5>zO#<}Q!LYO,5>zOOQQ1G.r1G.rO9qQWO1G.rO!-RQpO1G.rO!-ZQ#tO1G.rO#=]QWO1G0YO#=bQWO'#CgO#=mQWO'#JWO#=uQWO,5=ZO#=zQWO'#JWO#>PQWO'#IQO#>_QWO,5?qO#@ZQbO1G0{OOQ!LS1G0}1G0}O5ZQWO1G2tO#@bQWO1G2tO#@gQWO1G2tO#@lQWO1G2tOOQQ1G2t1G2tO#@qQ#tO1G2WO6jQWO'#IuO6jQWO'#EXO6jQWO'#HzO#ASQ!LYO,5?nOOQQ1G2[1G2[O!1wQWO1G2cOH_QWO1G2`O#A_QWO1G2`OOQQ1G2a1G2aOH_QWO1G2aO#AdQWO1G2aO#AlQ&jO'#G_OOQQ1G2c1G2cO!']Q&jO'#H|O!1|Q`O1G2fOOQQ1G2f1G2fOOQQ,5=P,5=PO#AtQ#tO,5=RO5ZQWO,5=RO#7WQWO,5=UO4wQWO,5=UO!-RQpO,5=UO!-ZQ#tO,5=UO5`Q#tO,5=UO#BVQWO'#JUO#BbQWO,5=VOOQQ1G.i1G.iO#BgQ!LYO1G.iO#BrQWO1G.iO!(uQWO1G.iO5RQ!LYO1G.iO#BwQbO,5?sO#CRQWO,5?sO#C^QYO,5=^O#CeQWO,5=^O6jQWO,5?sOOQQ1G2v1G2vO]QYO1G2vOOQQ1G2|1G2|OOQQ1G3O1G3OO9lQWO1G3QO#CjQYO1G3SO#GbQYO'#HROOQQ1G3V1G3VO9qQWO1G3]O#GoQWO1G3]O5RQ!LYO1G3aOOQQ1G3c1G3cOOQ!LQ'#Fn'#FnO5RQ!LYO1G3eO5RQ!LYO1G3gOOOS1G4p1G4pO#IkQ!LdO,5;yO#JOQbO1G3oO#JYQWO1G4rO#JbQWO1G4{O#JjQWO,5?YO! PQYO,5:sO6jQWO,5:sO9qQWO,59xO! PQYO,59xO!-RQpO,59xO#LcQ7^O,59xOOQO,5:s,5:sO#LmQ&jO'#HjO#MTQWO,5?XOOQ!LS1G/c1G/cO#M]Q&jO'#HoO#MqQWO,5?cOOQ!LQ1G0^1G0^O!<iQ&jO,59xO#MyQbO1G4|OOQO,5>Y,5>YO6jQWO,5>YOOQO-E;l-E;lO#NTQ!LrO'#D|O!']Q&jO'#DxOOQO'#Hm'#HmO#NoQ&jO,5:cOOQ!LS,5:c,5:cO#NvQ&jO'#DxO$ UQ&jO'#D|O$ jQ&jO'#D|O!']Q&jO'#D|O$ tQWO1G/zO$ yQ`O1G/zOOQ!LS1G/z1G/zO)rQYO1G/zOH_QWO1G/zOOQ!LS1G0X1G0XO9qQWO1G0XO!-RQpO1G0XO!-ZQ#tO1G0XO$!QQ!LdO1G4yO)rQYO1G4yO$!bQ!LYO1G4yO$!sQWO1G4xO6jQWO,5>[OOQO,5>[,5>[O$!{QWO,5>[OOQO-E;n-E;nO$!sQWO1G4xOOQ!LS,5;y,5;yO$#ZQ!LdO,59fO$%YQ!LdO,5;zO$'[Q!LdO,5;|O$)^Q!LdO,5<[OOQ!LS7+%t7+%tO$+fQWO'#HkO$+pQWO,5?ZOOQ!LS1G/m1G/mO$+xQYO'#HlO$,VQWO,5?[O$,_QbO,5?[OOQ!LS1G/r1G/rOOQ!LS7+%}7+%}O$,iQ7^O,5:XO)rQYO7+&_O$,sQ7^O,5:POOQO1G1T1G1TOOQO1G1U1G1UO$,zQMhO,5;lO! PQYO,5;kOOQO-E;o-E;oOOQ!LS7+&m7+&mOOQO7+&v7+&vOOOO1G1_1G1_O$-VQWO1G1_OOQ!LS1G1a1G1aO$-[Q!LdO7+&|OOQ!LS,5>a,5>aO$-{QWO,5>aOOQ!LS1G1u1G1uP$.QQWO'#HuPOQ!LS-E;s-E;sO$.qQ#tO1G2QO$/dQ#tO1G2SO$/nQ#tO1G2UOOQ!LS1G1n1G1nO$/uQWO'#HtO$0TQWO,5?jO$0TQWO,5?jO$0]QWO,5?jO$0hQWO,5?jOOQO1G1p1G1pO$0vQ#tO1G1oO$1WQWO'#HvO$1hQWO,5?kOH_QWO,5?kO$1pQ`O,5?kOOQ!LS1G1s1G1sO5RQ!LYO,5<^O5RQ!LYO,5<_O$1zQWO,5<_O#7RQWO,5<_O!-RQpO,5<^O$2PQWO,5<`O5RQ!LYO,5<aO$1zQWO,5<dOOQO-E;u-E;uOOQ!LS1G1w1G1wO!']Q&jO,5<^O$2XQWO,5<_O!']Q&jO,5<`O!']Q&jO,5<_O$2dQ#tO1G4jO$2nQ#tO1G4jOOQO,5>d,5>dOOQO-E;v-E;vOKTQ&jO,59hO)rQYO,59hO$2{QWO1G1jOIkQWO1G1qOOQ!LS7+&}7+&}OFbQYO7+&}OOQ!LS7+%V7+%VO$3QQ`O'#JQO$ tQWO7+'rO$3[QWO7+'rO$3dQ`O7+'rOOQQ7+'r7+'rOH_QWO7+'rO)rQYO7+'rOH_QWO7+'rOOQO1G.v1G.vO$3nQ!LbO'#CgO$4OQ!LbO,5<bO$4mQWO,5<bOOQ!LQ1G4f1G4fOOQQ7+$^7+$^O9qQWO7+$^O!-RQpO7+$^OFbQYO7+%tO$4rQWO'#IPO$4}QWO,5?rOOQO1G2u1G2uO5ZQWO,5?rOOQO,5>l,5>lOOQO-E<O-E<OOOQ!LS7+&g7+&gO$5VQWO7+(`O5RQ!LYO7+(`O5ZQWO7+(`O$5[QWO7+(`O$5aQWO7+'rOOQ!LQ,5>f,5>fOOQ!LQ-E;x-E;xOOQQ7+'}7+'}O$5oQ!LbO7+'zOH_QWO7+'zO$5yQ`O7+'{OOQQ7+'{7+'{OH_QWO7+'{O$6QQWO'#JTO$6]QWO,5<yOOQO,5>h,5>hOOQO-E;z-E;zOOQQ7+(Q7+(QO$7SQ&jO'#GhOOQQ1G2m1G2mOH_QWO1G2mO)rQYO1G2mOH_QWO1G2mO$7ZQWO1G2mO$7iQ#tO1G2mO5RQ!LYO1G2pO#7WQWO1G2pO4wQWO1G2pO!-RQpO1G2pO!-ZQ#tO1G2pO$7zQWO'#IOO$8VQWO,5?pO$8_Q&jO,5?pOOQ!LQ1G2q1G2qOOQQ7+$T7+$TO$8dQWO7+$TO5RQ!LYO7+$TO$8iQWO7+$TO)rQYO1G5_O)rQYO1G5`O$8nQYO1G2xO$8uQWO1G2xO$8zQYO1G2xO$9RQ!LYO1G5_OOQQ7+(b7+(bO5RQ!LYO7+(lO]QYO7+(nOOQQ'#JZ'#JZOOQQ'#IR'#IRO$9]QYO,5=mOOQQ,5=m,5=mO)rQYO'#HSO$9jQWO'#HUOOQQ7+(w7+(wO$9oQYO7+(wO6jQWO7+(wOOQQ7+({7+({OOQQ7+)P7+)POOQQ7+)R7+)ROOQO1G4t1G4tO$=jQ7^O1G0_O$=tQWO1G0_OOQO1G/d1G/dO$>PQ7^O1G/dO9qQWO1G/dO! PQYO'#D^OOQO,5>U,5>UOOQO-E;h-E;hOOQO,5>Z,5>ZOOQO-E;m-E;mO!-RQpO1G/dOOQO1G3t1G3tO9qQWO,5:dOOQO,5:h,5:hO!.OQYO,5:hO$>ZQ!LYO,5:hO$>fQ!LYO,5:hO!-RQpO,5:dOOQO-E;k-E;kOOQ!LS1G/}1G/}O!']Q&jO,5:dO$>tQ!LrO,5:hO$?`Q&jO,5:dO!']Q&jO,5:hO$?nQ&jO,5:hO$@SQ!LYO,5:hOOQ!LS7+%f7+%fO$ tQWO7+%fO$ yQ`O7+%fOOQ!LS7+%s7+%sO9qQWO7+%sO!-RQpO7+%sO$@hQ!LdO7+*eO)rQYO7+*eOOQO1G3v1G3vO6jQWO1G3vO$@xQWO7+*dO$AQQ!LdO1G2QO$CSQ!LdO1G2SO$EUQ!LdO1G1oO$G^Q#tO,5>VOOQO-E;i-E;iO$GhQbO,5>WO)rQYO,5>WOOQO-E;j-E;jO$GrQWO1G4vO$ItQ7^O1G0dO$KoQ7^O1G0dO$MjQ7^O1G0dO$MqQ7^O1G0dO% `Q7^O1G0dO% sQ7^O1G0dO%#zQ7^O1G0dO%$RQ7^O1G0dO%%|Q7^O1G0dO%&TQ7^O1G0dO%'xQ7^O1G0dO%(VQ!LdO<<IyO%(vQ7^O1G0dO%*fQ7^O'#IeO%,cQ7^O1G0vO! PQYO'#FSOOQO'#I|'#I|OOQO1G1W1G1WO%,jQWO1G1VO%,oQ7^O,5>_OOOO7+&y7+&yOOQ!LS1G3{1G3{OIkQWO7+'pO%,|QWO,5>`O5ZQWO,5>`OOQO-E;r-E;rO%-[QWO1G5UO%-[QWO1G5UO%-dQWO1G5UO%-oQ`O,5>bO%-yQWO,5>bOH_QWO,5>bOOQO-E;t-E;tO%.OQ`O1G5VO%.YQWO1G5VOOQO1G1x1G1xOOQO1G1y1G1yO5RQ!LYO1G1yO$1zQWO1G1yO5RQ!LYO1G1xO%.bQWO1G1zOH_QWO1G1zOOQO1G1{1G1{O5RQ!LYO1G2OO!-RQpO1G1xO#7RQWO1G1yO%.gQWO1G1zO%.oQWO1G1yOIkQWO7+*UOOQ!LS1G/S1G/SO%.zQWO1G/SOOQ!LS7+'U7+'UO%/PQ#tO7+']OOQ!LS<<Ji<<JiOH_QWO'#HyO%/aQWO,5?lOOQQ<<K^<<K^OH_QWO<<K^O$ tQWO<<K^O%/iQWO<<K^O%/qQ`O<<K^OH_QWO1G1|OOQQ<<Gx<<GxO9qQWO<<GxOOQ!LS<<I`<<I`OOQO,5>k,5>kO%/{QWO,5>kOOQO-E;}-E;}O%0QQWO1G5^O%0YQWO<<KzOOQQ<<Kz<<KzO%0_QWO<<KzO5RQ!LYO<<KzO)rQYO<<K^OH_QWO<<K^OOQQ<<Kf<<KfO$5oQ!LbO<<KfOOQQ<<Kg<<KgO$5yQ`O<<KgO%0dQ&jO'#H{O%0oQWO,5?oO! PQYO,5?oOOQQ1G2e1G2eO#NTQ!LrO'#D|O!']Q&jO'#GiOOQO'#H}'#H}O%0wQ&jO,5=SOOQQ,5=S,5=SO#7rQ&jO'#D|O%1OQ&jO'#D|O%1dQ&jO'#D|O%1nQ&jO'#GiO%1|QWO7+(XO%2RQWO7+(XO%2ZQ`O7+(XOOQQ7+(X7+(XOH_QWO7+(XO)rQYO7+(XOH_QWO7+(XO%2eQWO7+(XOOQQ7+([7+([O5RQ!LYO7+([O#7WQWO7+([O4wQWO7+([O!-RQpO7+([O%2sQWO,5>jOOQO-E;|-E;|OOQO'#Gl'#GlO%3OQWO1G5[O5RQ!LYO<<GoOOQQ<<Go<<GoO%3WQWO<<GoO%3]QWO7+*yO%3bQWO7+*zOOQQ7+(d7+(dO%3gQWO7+(dO%3lQYO7+(dO%3sQWO7+(dO)rQYO7+*yO)rQYO7+*zOOQQ<<LW<<LWOOQQ<<LY<<LYOOQQ-E<P-E<POOQQ1G3X1G3XO%3xQWO,5=nOOQQ,5=p,5=pO9qQWO<<LcO%3}QWO<<LcO! PQYO7+%yOOQO7+%O7+%OO%4SQ7^O1G4|O9qQWO7+%OOOQO1G0O1G0OO%4^Q!LdO1G0SOOQO1G0S1G0SO!.OQYO1G0SO%4hQ!LYO1G0SO9qQWO1G0OO!-RQpO1G0OO%4sQ!LYO1G0SO!']Q&jO1G0OO%5RQ!LYO1G0SO%5gQ!LrO1G0SO%5qQ&jO1G0OO!']Q&jO1G0SOOQ!LS<<IQ<<IQOOQ!LS<<I_<<I_O9qQWO<<I_O%6PQ!LdO<<NPOOQO7+)b7+)bO%6aQ!LdO7+']O%8iQbO1G3rO%8sQ7^O,5;yO%8}Q7^O,59fO%:zQ7^O,5;zO%<wQ7^O,5;|O%>tQ7^O,5<[O%@dQ7^O7+&|O%@kQWO,5;nOOQO7+&q7+&qO%@pQ#tO<<K[OOQO1G3z1G3zO%AQQWO1G3zO%A]QWO1G3zO%AkQWO7+*pO%AkQWO7+*pOH_QWO1G3|O%AsQ`O1G3|O%A}QWO7+*qOOQO7+'e7+'eO5RQ!LYO7+'eOOQO7+'d7+'dO$1zQWO7+'fO%BVQ`O7+'fOOQO7+'j7+'jO5RQ!LYO7+'dO$1zQWO7+'eO%B^QWO7+'fOH_QWO7+'fO#7RQWO7+'eO%BcQ#tO<<MpOOQ!LS7+$n7+$nO%BmQ`O,5>eOOQO-E;w-E;wO$ tQWOAN@xOOQQAN@xAN@xOH_QWOAN@xO%BwQ!LbO7+'hOOQQAN=dAN=dO5ZQWO1G4VO%CUQWO7+*xO5RQ!LYOANAfO%C^QWOANAfOOQQANAfANAfO%CcQWOAN@xO%CkQ`OAN@xOOQQANAQANAQOOQQANARANARO%CuQWO,5>gOOQO-E;y-E;yO%DQQ7^O1G5ZO#7WQWO,5=TO4wQWO,5=TO!-RQpO,5=TOOQO-E;{-E;{OOQQ1G2n1G2nO$>tQ!LrO,5:hO!']Q&jO,5=TO%D[Q&jO,5=TO%DjQ&jO,5:hOOQQ<<Ks<<KsOH_QWO<<KsO%1|QWO<<KsO%EOQWO<<KsO%EWQ`O<<KsO)rQYO<<KsOH_QWO<<KsOOQQ<<Kv<<KvO5RQ!LYO<<KvO#7WQWO<<KvO4wQWO<<KvO%EbQ&jO1G4UO%EgQWO7+*vOOQQAN=ZAN=ZO5RQ!LYOAN=ZOOQQ<<Ne<<NeOOQQ<<Nf<<NfOOQQ<<LO<<LOO%EoQWO<<LOO%EtQYO<<LOO%E{QWO<<NeO%FQQWO<<NfOOQQ1G3Y1G3YOOQQANA}ANA}O9qQWOANA}O%FVQ7^O<<IeOOQO<<Hj<<HjOOQO7+%n7+%nO%4^Q!LdO7+%nO!.OQYO7+%nOOQO7+%j7+%jO9qQWO7+%jO%FaQ!LYO7+%nO!-RQpO7+%jO%FlQ!LYO7+%nO!']Q&jO7+%jO%FzQ!LYO7+%nOOQ!LSAN>yAN>yO%G`Q!LdO<<K[O%IhQ7^O<<IyO%IoQ7^O1G1oO%K_Q7^O1G2QO%M[Q7^O1G2SOOQO1G1Y1G1YOOQO7+)f7+)fO& XQWO7+)fO& dQWO<<N[O& lQ`O7+)hOOQO<<KP<<KPO5RQ!LYO<<KQO$1zQWO<<KQOOQO<<KO<<KOO5RQ!LYO<<KPO& vQ`O<<KQO$1zQWO<<KPOOQQG26dG26dO$ tQWOG26dOOQO7+)q7+)qOOQQG27QG27QO5RQ!LYOG27QOH_QWOG26dO! PQYO1G4RO& }QWO7+*uO5RQ!LYO1G2oO#7WQWO1G2oO4wQWO1G2oO!-RQpO1G2oO!']Q&jO1G2oO%5gQ!LrO1G0SO&!VQ&jO1G2oO%1|QWOANA_OOQQANA_ANA_OH_QWOANA_O&!eQWOANA_O&!mQ`OANA_OOQQANAbANAbO5RQ!LYOANAbO#7WQWOANAbOOQO'#Gm'#GmOOQO7+)p7+)pOOQQG22uG22uOOQQANAjANAjO&!wQWOANAjOOQQANDPANDPOOQQANDQANDQO&!|QYOG27iOOQO<<IY<<IYO%4^Q!LdO<<IYOOQO<<IU<<IUO!.OQYO<<IYO9qQWO<<IUO&&wQ!LYO<<IYO!-RQpO<<IUO&'SQ!LYO<<IYO&'bQ7^O7+']OOQO<<MQ<<MQOOQOAN@lAN@lO5RQ!LYOAN@lOOQOAN@kAN@kO$1zQWOAN@lO5RQ!LYOAN@kOOQQLD,OLD,OOOQQLD,lLD,lO$ tQWOLD,OO&)QQ7^O7+)mOOQO7+(Z7+(ZO5RQ!LYO7+(ZO#7WQWO7+(ZO4wQWO7+(ZO!-RQpO7+(ZO!']Q&jO7+(ZOOQQG26yG26yO%1|QWOG26yOH_QWOG26yOOQQG26|G26|O5RQ!LYOG26|OOQQG27UG27UO9qQWOLD-TOOQOAN>tAN>tO%4^Q!LdOAN>tOOQOAN>pAN>pO!.OQYOAN>tO9qQWOAN>pO&)[Q!LYOAN>tO&)gQ7^O<<K[OOQOG26WG26WO5RQ!LYOG26WOOQOG26VG26VOOQQ!$( j!$( jOOQO<<Ku<<KuO5RQ!LYO<<KuO#7WQWO<<KuO4wQWO<<KuO!-RQpO<<KuOOQQLD,eLD,eO%1|QWOLD,eOOQQLD,hLD,hOOQQ!$(!o!$(!oOOQOG24`G24`O%4^Q!LdOG24`OOQOG24[G24[O!.OQYOG24`OOQOLD+rLD+rOOQOANAaANAaO5RQ!LYOANAaO#7WQWOANAaO4wQWOANAaOOQQ!$(!P!$(!POOQOLD)zLD)zO%4^Q!LdOLD)zOOQOG26{G26{O5RQ!LYOG26{O#7WQWOG26{OOQO!$'Mf!$'MfOOQOLD,gLD,gO5RQ!LYOLD,gOOQO!$(!R!$(!ROKYQYO'#DmO&+VQ!LdO'#IgO&+jQ!LdO'#IgOKYQYO'#DeO&+qQ!LdO'#CgO&,[QbO'#CgO&,lQYO,5:nOFbQYO,5:nOKYQYO,5:xOKYQYO,5:xOKYQYO,5:xOKYQYO,5:xOKYQYO,5:xOKYQYO,5:xOKYQYO,5:xOKYQYO,5:xOKYQYO,5:xOKYQYO,5:xOKYQYO,5:xO! PQYO'#HsO&.iQWO,5;yO&.qQWO,5:xOKYQYO,5;[O!(uQWO'#CzO!(uQWO'#CzOH_QWO'#F`O&.qQWO'#F`OH_QWO'#FbO&.qQWO'#FbOH_QWO'#FpO&.qQWO'#FpO! PQYO,5?bO&,lQYO1G0YOFbQYO1G0YO&/xQ7^O'#CgO&0SQ7^O'#IgO&0^Q7^O'#IgOKYQYO1G1bOH_QWO,5<fO&.qQWO,5<fOH_QWO,5<hO&.qQWO,5<hOH_QWO,5<TO&.qQWO,5<TO&,lQYO1G1cOFbQYO1G1cO&,lQYO1G1cO&,lQYO1G0YOKYQYO7+&_OH_QWO1G1qO&.qQWO1G1qO&,lQYO7+&}OFbQYO7+&}O&,lQYO7+&}O&,lQYO7+%tOFbQYO7+%tO&,lQYO7+%tOH_QWO7+'pO&.qQWO7+'pO&0eQWO'#ESO&0jQWO'#ESO&0oQWO'#ESO&0wQWO'#ESO&1PQWO'#EpO!.OQYO'#DeO!.OQYO'#DmO&1UQWO'#IvO&1aQWO'#ItO&1lQWO,5:nO&1qQWO,5:nO!.OQYO,5:xO!.OQYO,5:xO!.OQYO,5:xO!.OQYO,5:xO!.OQYO,5:xO!.OQYO,5:xO!.OQYO,5:xO!.OQYO,5:xO!.OQYO,5:xO!.OQYO,5:xO!.OQYO,5:xO!.OQYO,5;[O&1vQ#tO,5;vO&1}QWO'#FiO&2SQWO'#FiO&2XQWO,5;wO&2aQWO,5;wO&2iQWO,5;wO&2qQ!LdO,5:PO&3OQWO,5:nO&3TQWO,5:nO&3]QWO,5:nO&3eQWO,5:nO&5aQ!LdO1G0dO&5nQ!LdO1G0dO&7uQ!LdO1G0dO&7|Q!LdO1G0dO&9}Q!LdO1G0dO&:UQ!LdO1G0dO&<VQ!LdO1G0dO&<^Q!LdO1G0dO&>_Q!LdO1G0dO&>fQ!LdO1G0dO&>mQ7^O1G0sO&>tQ!LdO1G0vO!.OQYO1G1bO&?RQWO,5<VO&?WQWO,5<VO&?]QWO1G1cO&?bQWO1G1cO&?gQWO1G1cO&?lQWO1G0YO&?qQWO1G0YO&?vQWO1G0YO!.OQYO7+&_O&?{Q!LdO7+&|O&@YQ#tO1G2UO&@aQ#tO1G2UO&@hQ!LdO<<IyO&,lQYO,5:nO&BiQ!LdO'#IhO&B|QWO'#EpO3mQWO'#F^O4wQWO'#F]O4wQWO'#F]O4wQWO'#F]OBUQWO'#EROBUQWO'#EROBUQWO'#EROKYQYO,5;XO&CRQ#tO,5;vO!)PQWO'#FkO!)PQWO'#FkO&CYQ7^O1G0sOIkQWO,5<jOIkQWO,5<jO! PQYO'#DmO! PQYO'#DeO! PQYO,5:xO! PQYO,5:xO! PQYO,5:xO! PQYO,5:xO! PQYO,5:xO! PQYO,5:xO! PQYO,5:xO! PQYO,5:xO! PQYO,5:xO! PQYO,5:xO! PQYO,5:xO! PQYO,5;[O! PQYO1G1bO! PQYO7+&_O&CaQWO'#ESO&CfQWO'#ESO&CnQWO'#EpO&CsQ#tO,5;vO&CzQ7^O1G0sO3mQWO'#F^OKYQYO,5;XO&DRQ7^O'#IhO&DcQ7^O,5:PO&DpQ7^O1G0dO&FqQ7^O1G0dO&FxQ7^O1G0dO&HmQ7^O1G0dO&IQQ7^O1G0dO&K_Q7^O1G0dO&KfQ7^O1G0dO&MgQ7^O1G0dO&MnQ7^O1G0dO' cQ7^O1G0dO' vQ7^O1G0vO'!TQ7^O7+&|O'!bQ7^O<<IyO3mQWO'#F^OKYQYO,5;X",
  stateData: "'#b~O&}OSSOSTOS~OPTOQTOWwO]bO^gOamOblOgbOiTOjbOkbOmTOoTOtROvbOwbOxbO!OSO!YjO!_UO!bTO!cTO!dTO!eTO!fTO!ikO#jnO#n]O$uoO$wrO$ypO$zpO${qO%OsO%QtO%TuO%UuO%WvO%exO%kyO%mzO%o{O%q|O%t}O%z!OO&O!PO&Q!QO&S!RO&U!SO&W!TO'PPO']QO'q`O~OPZXYZX^ZXiZXqZXrZXtZX|ZX![ZX!]ZX!_ZX!eZX!tZX#OcX#RZX#SZX#TZX#UZX#VZX#WZX#XZX#YZX#ZZX#]ZX#_ZX#`ZX#eZX&{ZX']ZX'eZX'lZX'mZX~O!W$bX~P$tO&x!VO&y!UO&z!XO~OPTOQTO]bOa!hOb!gOgbOiTOjbOkbOmTOoTOtROvbOwbOxbO!O!`O!YjO!_UO!bTO!cTO!dTO!eTO!fTO!i!fO#j!iO#n]O'P!YO']QO'q`O~O{!^O|!ZOy'`Py'iP~P'^O}!jO~P]OPTOQTO]bOa!hOb!gOgbOiTOjbOkbOmTOoTOtROvbOwbOxbO!O!`O!YjO!_UO!bTO!cTO!dTO!eTO!fTO!i!fO#j!iO#n]O'P8ZO']QO'q`O~OPTOQTO]bOa!hOb!gOgbOiTOjbOkbOmTOoTOtROvbOwbOxbO!O!`O!YjO!_UO!bTO!cTO!dTO!eTO!fTO!i!fO#j!iO#n]O']QO'q`O~O{!oO!|!rO!}!oO'P8[O!^'fP~P+oO#O!sO~O!W!tO#O!sO~OP#ZOY#aOi#OOq!xOr!xOt!yO|#_O![#QO!]!vO!_!wO!e#ZO#R!|O#S!}O#T!}O#U!}O#V#PO#W#QO#X#QO#Y#QO#Z#RO#]#TO#_#VO#`#WO']QO'e#XO'l!zO'm!{O^'ZX&{'ZX!^'ZXy'ZX!O'ZX$v'ZX!W'ZX~O!t#bO#e#bOP'[XY'[X^'[Xi'[Xq'[Xr'[Xt'[X|'[X!['[X!]'[X!_'[X!e'[X#R'[X#S'[X#T'[X#U'[X#V'[X#W'[X#Y'[X#Z'[X#]'[X#_'[X#`'[X']'[X'e'[X'l'[X'm'[X~O#X'[X&{'[Xy'[X!^'[X'_'[X!O'[X$v'[X!W'[X~P0gO!t#bO~O#p#cO#w#gO~O!O#hO#n]O#z#iO#|#kO~O]#nOg#zOi#oOj#nOk#nOm#{Oo#|Ot#tO!O#uO!Y$RO!_#rO!}$SO#j$PO$T#}O$V$OO$Y$QO'P#mO'T'VP~O!_$TO~O!W$VO~O^$WO&{$WO~O'P$[O~O!_$TO'P$[O'Q$^O'U$_O~Ob$eO!_$TO'P$[O~O]$nOq$jO!O$gO!_$iO$w$mO'P$[O'Q$^O['yP~O!i$oO~Ot$pO!O$qO'P$[O~Ot$pO!O$qO%Q$uO'P$[O~O'P$vO~O$wrO$ypO$zpO${qO%OsO%QtO%TuO%UuO~Oa%POb%OO!i$|O$u$}O%Y${O~P7YOa%SOblO!O%RO!ikO$uoO$ypO$zpO${qO%OsO%QtO%TuO%UuO%WvO~O_%VO!t%YO$w%TO'Q$^O~P8XO!_%ZO!b%_O~O!_%`O~O!OSO~O^$WO&w%hO&{$WO~O^$WO&w%kO&{$WO~O^$WO&w%mO&{$WO~O&x!VO&y!UO&z%qO~OPZXYZXiZXqZXrZXtZX|ZX|cX![ZX!]ZX!_ZX!eZX!tZX!tcX#OcX#RZX#SZX#TZX#UZX#VZX#WZX#XZX#YZX#ZZX#]ZX#_ZX#`ZX#eZX']ZX'eZX'lZX'mZX~OyZXycX~P:tO{%sOy&]X|&]X~P)rO|!ZOy'`X~OP#ZOY#aOi#OOq!xOr!xOt!yO|!ZO![#QO!]!vO!_!wO!e#ZO#R!|O#S!}O#T!}O#U!}O#V#PO#W#QO#X#QO#Y#QO#Z#RO#]#TO#_#VO#`#WO']QO'e#XO'l!zO'm!{O~Oy'`X~P=kOy%xO~Ot%{O!R&VO!S&OO!T&OO'Q$^O~O]%|Oj%|O{&PO'Y%yO}'aP}'kP~P?nOy'hX|'hX!W'hX!^'hX'e'hX~O!t'hX#O!wX}'hX~P@gO!t&WOy'jX|'jX~O|&XOy'iX~Oy&ZO~O!t#bO~P@gOR&_O!O&[O!j&^O'P$[O~Ob&dO!_$TO'P$[O~Oq$jO!_$iO~O}&eO~P]Oq!xOr!xOt!yO!]!vO!_!wO']QOP!aaY!aai!aa|!aa![!aa!e!aa#R!aa#S!aa#T!aa#U!aa#V!aa#W!aa#X!aa#Y!aa#Z!aa#]!aa#_!aa#`!aa'e!aa'l!aa'm!aa~O^!aa&{!aay!aa!^!aa'_!aa!O!aa$v!aa!W!aa~PBpO!^&fO~O!W!tO!t&hO'e&gO|'gX^'gX&{'gX~O!^'gX~PEYO|&lO!^'fX~O!^&nO~Ot$pO!O$qO!}&oO'P$[O~OPTOQTO]bOa!hOb!gOgbOiTOjbOkbOmTOoTOtROvbOwbOxbO!OSO!YjO!_UO!bTO!cTO!dTO!eTO!fTO!i!fO#j!iO#n]O'P8ZO']QO'q`O~O]#nOg#zOi#oOj#nOk#nOm#{Oo8nOt#tO!O#uO!Y;OO!_#rO!}8tO#j$PO$T8pO$V8rO$Y$QO'P&rO~O#O&tO~O]#nOg#zOi#oOj#nOk#nOm#{Oo#|Ot#tO!O#uO!Y$RO!_#rO!}$SO#j$PO$T#}O$V$OO$Y$QO'P&rO~O'T'cP~PIkO{&xO!^'dP~P)rO'Y&zO~OP8VOQ8VO]bOa:yOb!gOgbOi8VOjbOkbOm8VOo8VOtROvbOwbOxbO!O!`O!Y8YO!_UO!b8VO!c8VO!d8VO!e8VO!f8VO!i!fO#j!iO#n]O'P'YO']QO'q:uO~O!_!wO~O|#_O^$Ra&{$Ra!^$Ray$Ra!O$Ra$v$Ra!W$Ra~O!W'bO!O'nX#m'nX#p'nX#w'nX~Oq'cO~PMvOq'cO!O'nX#m'nX#p'nX#w'nX~O!O'eO#m'iO#p'dO#w'jO~OP;TOQ;TO]bOa:{Ob!gOgbOi;TOjbOkbOm;TOo;TOtROvbOwbOxbO!O!`O!Y;UO!_UO!b;TO!c;TO!d;TO!e;TO!f;TO!i!fO#j!iO#n]O'P'YO']QO'q;{O~O{'mO~P! PO#p#cO#w'pO~Oq$ZXt$ZX!]$ZX'e$ZX'l$ZX'm$ZX~OReX|eX!teX'TeX'T$ZX~P!#]Oj'rO~Oq'tOt'uO'e#XO'l'wO'm'yO~O'T'sO~P!$ZO'T'|O~O]#nOg#zOi#oOj#nOk#nOm#{Oo8nOt#tO!O#uO!Y;OO!_#rO!}8tO#j$PO$T8pO$V8rO$Y$QO~O{(QO'P'}O!^'rP~P!$xO#O(SO~O{(WO'P(TOy'sP~P!$xO^(aOi(fOt(^O!R(dO!S(]O!T(]O!_(ZO!q(eO$m(`O'Q$^O'Y(YO~O}(cO~P!&mO!]!vOq'XXt'XX'e'XX'l'XX'm'XX|'XX!t'XX~O'T'XX#c'XX~P!'iOR(iO!t(hO|'WX'T'WX~O|(jO'T'VX~O'P(lO~O!_(qO~O!_(ZO~Ot$pO{!oO!O$qO!|!rO!}!oO'P$[O!^'fP~O!W!tO#O(uO~OP#ZOY#aOi#OOq!xOr!xOt!yO![#QO!]!vO!_!wO!e#ZO#R!|O#S!}O#T!}O#U!}O#V#PO#W#QO#X#QO#Y#QO#Z#RO#]#TO#_#VO#`#WO']QO'e#XO'l!zO'm!{O~O^!Xa|!Xa&{!Xay!Xa!^!Xa'_!Xa!O!Xa$v!Xa!W!Xa~P!)uOR(}O!O&[O!j(|O$v({O'U$_O~O'P$vO'T'VP~O!W)QO!O'SX^'SX&{'SX~O!_$TO'U$_O~O!_$TO'P$[O'U$_O~O!W!tO#O&tO~O'P)YO}'zP~O|)^O['yX~OP9jOQ9jO]bOa:zOb!gOgbOi9jOjbOkbOm9jOo9jOtROvbOwbOxbO!O!`O!Y9iO!_UO!b9jO!c9jO!d9jO!e9jO!f9jO!i!fO#j!iO#n]O'P8ZO']QO'q;jO~OY)bO~O[)cO~O!O$gO'P$[O'Q$^O['yP~Ot$pO{)hO!O$qO'P$[Oy'iP~O]&SOj&SO{)iO'Y&zO}'kP~O|)jO^'vX&{'vX~O!t)nO'U$_O~OR)qO!O#uO'U$_O~O!O)sO~Oq)uO!OSO~O!i)zO~Ob*PO~O'P(lO}'xP~Ob$eO~O$wrO'P$vO~P8XOY*VO[*UO~OPTOQTO]bOamOblOgbOiTOjbOkbOmTOoTOtROvbOwbOxbO!YjO!_UO!bTO!cTO!dTO!eTO!fTO!ikO#n]O$uoO']QO'q`O~O!O!`O#j!iO'P8ZO~P!3PO[*UO^$WO&{$WO~O^*ZO$y*]O$z*]O${*]O~P)rO!_%ZO~O%k*bO~O!O*dO~O%{*gO%|*fOP%yaQ%yaW%ya]%ya^%yaa%yab%yag%yai%yaj%yak%yam%yao%yat%yav%yaw%yax%ya!O%ya!Y%ya!_%ya!b%ya!c%ya!d%ya!e%ya!f%ya!i%ya#j%ya#n%ya$u%ya$w%ya$y%ya$z%ya${%ya%O%ya%Q%ya%T%ya%U%ya%W%ya%e%ya%k%ya%m%ya%o%ya%q%ya%t%ya%z%ya&O%ya&Q%ya&S%ya&U%ya&W%ya&v%ya'P%ya']%ya'q%ya}%ya%r%ya_%ya%w%ya~O'P*jO~O'_*mO~Oy&]a|&]a~P!)uO|!ZOy'`a~Oy'`a~P=kO|&XOy'ia~O|sX|!UX}sX}!UX!WsX!W!UX!_!UX!tsX'U!UX~O!W*tO!t*sO|!{X|'bX}!{X}'bX!W'bX!_'bX'U'bX~O!W*vO!_$TO'U$_O|!QX}!QX~O]%zOj%zOt%{O'Y(YO~OP;TOQ;TO]bOa:{Ob!gOgbOi;TOjbOkbOm;TOo;TOtROvbOwbOxbO!O!`O!Y;UO!_UO!b;TO!c;TO!d;TO!e;TO!f;TO!i!fO#j!iO#n]O']QO'q;{O~O'P8yO~P!<wO|*zO}'aX~O}*|O~O!W*tO!t*sO|!{X}!{X~O|*}O}'kX~O}+PO~O]%zOj%zOt%{O'Q$^O'Y(YO~O!S+QO!T+QO~P!?rOt$pO{+TO!O$qO'P$[Oy&bX|&bX~O^+XO!R+[O!S+WO!T+WO!m+^O!n+]O!o+]O!q+_O'Q$^O'Y(YO~O}+ZO~P!@sOR+dO!O&[O!j+cO~O!t+jO|'ga!^'ga^'ga&{'ga~O!W!tO~P!AwO|&lO!^'fa~Ot$pO{+mO!O$qO!|+oO!}+mO'P$[O|&dX!^&dX~O#O!sa|!sa!^!sa!t!sa!O!sa^!sa&{!say!sa~P!$ZO#O'XXP'XXY'XX^'XXi'XXr'XX!['XX!_'XX!e'XX#R'XX#S'XX#T'XX#U'XX#V'XX#W'XX#X'XX#Y'XX#Z'XX#]'XX#_'XX#`'XX&{'XX']'XX!^'XXy'XX!O'XX$v'XX'_'XX!W'XX~P!'iO|+xO'T'cX~P!$ZO'T+zO~O|+{O!^'dX~P!)uO!^,OO~Oy,PO~OP#ZOq!xOr!xOt!yO!]!vO!_!wO!e#ZO']QOY#Qi^#Qii#Qi|#Qi![#Qi#S#Qi#T#Qi#U#Qi#V#Qi#W#Qi#X#Qi#Y#Qi#Z#Qi#]#Qi#_#Qi#`#Qi&{#Qi'e#Qi'l#Qi'm#Qiy#Qi!^#Qi'_#Qi!O#Qi$v#Qi!W#Qi~O#R#Qi~P!FeO#R!|O~P!FeOP#ZOq!xOr!xOt!yO!]!vO!_!wO!e#ZO#R!|O#S!}O#T!}O#U!}O']QOY#Qi^#Qi|#Qi![#Qi#V#Qi#W#Qi#X#Qi#Y#Qi#Z#Qi#]#Qi#_#Qi#`#Qi&{#Qi'e#Qi'l#Qi'm#Qiy#Qi!^#Qi'_#Qi!O#Qi$v#Qi!W#Qi~Oi#Qi~P!IPOi#OO~P!IPOP#ZOi#OOq!xOr!xOt!yO!]!vO!_!wO!e#ZO#R!|O#S!}O#T!}O#U!}O#V#PO']QO^#Qi|#Qi#Z#Qi#]#Qi#_#Qi#`#Qi&{#Qi'e#Qi'l#Qi'm#Qiy#Qi!^#Qi'_#Qi!O#Qi$v#Qi!W#Qi~OY#Qi![#Qi#W#Qi#X#Qi#Y#Qi~P!KkOY#aO![#QO#W#QO#X#QO#Y#QO~P!KkOP#ZOY#aOi#OOq!xOr!xOt!yO![#QO!]!vO!_!wO!e#ZO#R!|O#S!}O#T!}O#U!}O#V#PO#W#QO#X#QO#Y#QO#Z#RO']QO^#Qi|#Qi#]#Qi#_#Qi#`#Qi&{#Qi'e#Qi'm#Qiy#Qi!^#Qi'_#Qi!O#Qi$v#Qi!W#Qi~O'l#Qi~P!NcO'l!zO~P!NcOP#ZOY#aOi#OOq!xOr!xOt!yO![#QO!]!vO!_!wO!e#ZO#R!|O#S!}O#T!}O#U!}O#V#PO#W#QO#X#QO#Y#QO#Z#RO#]#TO']QO'l!zO^#Qi|#Qi#_#Qi#`#Qi&{#Qi'e#Qiy#Qi!^#Qi'_#Qi!O#Qi$v#Qi!W#Qi~O'm#Qi~P#!}O'm!{O~P#!}OP#ZOY#aOi#OOq!xOr!xOt!yO![#QO!]!vO!_!wO!e#ZO#R!|O#S!}O#T!}O#U!}O#V#PO#W#QO#X#QO#Y#QO#Z#RO#]#TO#_#VO']QO'l!zO'm!{O~O^#Qi|#Qi#`#Qi&{#Qi'e#Qiy#Qi!^#Qi'_#Qi!O#Qi$v#Qi!W#Qi~P#%iOPZXYZXiZXqZXrZXtZX![ZX!]ZX!_ZX!eZX!tZX#OcX#RZX#SZX#TZX#UZX#VZX#WZX#XZX#YZX#ZZX#]ZX#_ZX#`ZX#eZX']ZX'eZX'lZX'mZX|ZX}ZX~O#cZX~P#'|OP#ZOY8lOi8aOq!xOr!xOt!yO![8cO!]!vO!_!wO!e#ZO#R8_O#S8`O#T8`O#U8`O#V8bO#W8cO#X8cO#Y8cO#Z8dO#]8fO#_8hO#`8iO']QO'e#XO'l!zO'm!{O~O#c,RO~P#*WOP'[XY'[Xi'[Xq'[Xr'[Xt'[X!['[X!]'[X!_'[X!e'[X#R'[X#S'[X#T'[X#U'[X#V'[X#W'[X#X'[X#Y'[X#Z'[X#]'[X#_'[X#`'[X#c'[X']'[X'e'[X'l'[X'm'[X~O!t8mO#e8mO~P#,RO^&ga|&ga&{&ga!^&ga'_&gay&ga!O&ga$v&ga!W&ga~P!)uOP#QiY#Qi^#Qii#Qir#Qi|#Qi![#Qi!]#Qi!_#Qi!e#Qi#R#Qi#S#Qi#T#Qi#U#Qi#V#Qi#W#Qi#X#Qi#Y#Qi#Z#Qi#]#Qi#_#Qi#`#Qi&{#Qi']#Qiy#Qi!^#Qi'_#Qi!O#Qi$v#Qi!W#Qi~P!$ZO^#di|#di&{#diy#di!^#di'_#di!O#di$v#di!W#di~P!)uO#p,TO~O#p,UO~O!W'bO!t,VO!O#tX#m#tX#p#tX#w#tX~O{,WO~O!O'eO#m,YO#p'dO#w,ZO~OP#ZOY8lOi;XOq!xOr!xOt!yO|8jO![;ZO!]!vO!_!wO!e#ZO#R;VO#S;WO#T;WO#U;WO#V;YO#W;ZO#X;ZO#Y;ZO#Z;[O#];^O#_;`O#`;aO']QO'e#XO'l!zO'm!{O}'ZX~O},[O~O#w,^O~O],aOj,aOy,bO~O|cX!WcX!^cX!^$ZX'ecX~P!#]O!^,hO~P!$ZO|,iO!W!tO'e&gO!^'rX~O!^,nO~Oy$ZX|$ZX!W$bX~P!#]O|,pOy'sX~P!$ZO!W,rO~Oy,tO~O{(QO'P$[O!^'rP~Oi,xO!W!tO!_$TO'U$_O'e&gO~O!W)QO~O}-OO~P!&mO!S-PO!T-PO'Q$^O'Y(YO~Ot-RO'Y(YO~O!q-SO~O'P$vO|&lX'T&lX~O|(jO'T'Va~Oq-XOr-XOt-YO'ena'lna'mna|na!tna~O'Tna#cna~P#8dOq'tOt'uO'e$Sa'l$Sa'm$Sa|$Sa!t$Sa~O'T$Sa#c$Sa~P#9YOq'tOt'uO'e$Ua'l$Ua'm$Ua|$Ua!t$Ua~O'T$Ua#c$Ua~P#9{O]-ZO~O#O-[O~O'T$da|$da#c$da!t$da~P!$ZO#O-^O~OR-gO!O&[O!j-fO$v-eO~O'T-hO~O]#nOi#oOj#nOk#nOm#{Oo8nOt#tO!O#uO!Y;OO!_#rO!}8tO#j$PO$T8pO$V8rO$Y$QO~Og-jO'P-iO~P#;rO!W)QO!O'Sa^'Sa&{'Sa~O#O-pO~OYZX|cX}cX~O|-qO}'zX~O}-sO~OY-tO~O!O$gO'P$[O[&tX|&tX~O|)^O['ya~OP#ZOY#aOi9qOq!xOr!xOt!yO![9sO!]!vO!_!wO!e#ZO#R9oO#S9pO#T9pO#U9pO#V9rO#W9sO#X9sO#Y9sO#Z9tO#]9vO#_9xO#`9yO']QO'e#XO'l!zO'm!{O~O!^-wO~P#>gO]-yO~OY-zO~O[-{O~OR-gO!O&[O!j-fO$v-eO'U$_O~O|)jO^'va&{'va~O!t.RO~OR.UO!O#uO~O'Y&zO}'wP~OR.`O!O.[O!j._O$v.^O'U$_O~OY.jO|.hO}'xX~O}.kO~O[.mO^$WO&{$WO~O].nO~O#X.pO%i.qO~P0gO!t#bO#X.pO%i.qO~O^.rO~P)rO^.tO~O%r.xOP%piQ%piW%pi]%pi^%pia%pib%pig%pii%pij%pik%pim%pio%pit%piv%piw%pix%pi!O%pi!Y%pi!_%pi!b%pi!c%pi!d%pi!e%pi!f%pi!i%pi#j%pi#n%pi$u%pi$w%pi$y%pi$z%pi${%pi%O%pi%Q%pi%T%pi%U%pi%W%pi%e%pi%k%pi%m%pi%o%pi%q%pi%t%pi%z%pi&O%pi&Q%pi&S%pi&U%pi&W%pi&v%pi'P%pi']%pi'q%pi}%pi_%pi%w%pi~O_/OO}.|O%w.}O~P]O!OSO!_/RO~OP$RaY$Rai$Raq$Rar$Rat$Ra![$Ra!]$Ra!_$Ra!e$Ra#R$Ra#S$Ra#T$Ra#U$Ra#V$Ra#W$Ra#X$Ra#Y$Ra#Z$Ra#]$Ra#_$Ra#`$Ra']$Ra'e$Ra'l$Ra'm$Ra~O|#_O'_$Ra!^$Ra^$Ra&{$Ra~P#GwOy&]i|&]i~P!)uO|!ZOy'`i~O|&XOy'ii~Oy/VO~OP#ZOY8lOi;XOq!xOr!xOt!yO![;ZO!]!vO!_!wO!e#ZO#R;VO#S;WO#T;WO#U;WO#V;YO#W;ZO#X;ZO#Y;ZO#Z;[O#];^O#_;`O#`;aO']QO'e#XO'l!zO'm!{O~O|!Qa}!Qa~P#JoO]%zOj%zO{/]O'Y(YO|&^X}&^X~P?nO|*zO}'aa~O]&SOj&SO{)iO'Y&zO|&cX}&cX~O|*}O}'ka~Oy'ji|'ji~P!)uO^$WO!W!tO!_$TO!e/hO!t/fO&{$WO'U$_O'e&gO~O}/kO~P!@sO!S/lO!T/lO'Q$^O'Y(YO~O!R/nO!S/lO!T/lO!q/oO'Q$^O'Y(YO~O!n/pO!o/pO~P$ UO!O&[O~O!O&[O~P!$ZO|'gi!^'gi^'gi&{'gi~P!)uO!t/yO|'gi!^'gi^'gi&{'gi~O|&lO!^'fi~Ot$pO!O$qO!}/{O'P$[O~O#OnaPnaYna^naina![na!]na!_na!ena#Rna#Sna#Tna#Una#Vna#Wna#Xna#Yna#Zna#]na#_na#`na&{na']na!^nayna!Ona$vna'_na!Wna~P#8dO#O$SaP$SaY$Sa^$Sai$Sar$Sa![$Sa!]$Sa!_$Sa!e$Sa#R$Sa#S$Sa#T$Sa#U$Sa#V$Sa#W$Sa#X$Sa#Y$Sa#Z$Sa#]$Sa#_$Sa#`$Sa&{$Sa']$Sa!^$Say$Sa!O$Sa$v$Sa'_$Sa!W$Sa~P#9YO#O$UaP$UaY$Ua^$Uai$Uar$Ua![$Ua!]$Ua!_$Ua!e$Ua#R$Ua#S$Ua#T$Ua#U$Ua#V$Ua#W$Ua#X$Ua#Y$Ua#Z$Ua#]$Ua#_$Ua#`$Ua&{$Ua']$Ua!^$Uay$Ua!O$Ua$v$Ua'_$Ua!W$Ua~P#9{O#O$daP$daY$da^$dai$dar$da|$da![$da!]$da!_$da!e$da#R$da#S$da#T$da#U$da#V$da#W$da#X$da#Y$da#Z$da#]$da#_$da#`$da&{$da']$da!^$day$da!O$da!t$da$v$da'_$da!W$da~P!$ZO|&_X'T&_X~PIkO|+xO'T'ca~O{0TO|&`X!^&`X~P)rO|+{O!^'da~O|+{O!^'da~P!)uO#c!aa}!aa~PBpO#c!Xa~P#*WO!O0gO#n]O#u0hO~O}0lO~O^$Oq|$Oq&{$Oqy$Oq!^$Oq'_$Oq!O$Oq$v$Oq!W$Oq~P!)uOy0mO~O],aOj,aO~Oq'tOt'uO'm'yO'e$ni'l$ni|$ni!t$ni~O'T$ni#c$ni~P$.YOq'tOt'uO'e$pi'l$pi'm$pi|$pi!t$pi~O'T$pi#c$pi~P$.{O#c0nO~P!$ZO{0pO'P$[O|&hX!^&hX~O|,iO!^'ra~O|,iO!W!tO!^'ra~O|,iO!W!tO'e&gO!^'ra~O'T$]i|$]i#c$]i!t$]i~P!$ZO{0wO'P(TOy&jX|&jX~P!$xO|,pOy'sa~O|,pOy'sa~P!$ZO!W!tO~O!W!tO#X1RO~Oi1VO!W!tO'e&gO~O|'Wi'T'Wi~P!$ZO!t1YO|'Wi'T'Wi~P!$ZO!^1]O~O|1`O!O'tX~P!$ZO!O&[O$v1cO~O!O&[O$v1cO~P!$ZO!O$ZX$kZX^$ZX&{$ZX~P!#]O$k1gOqfXtfX!OfX'efX'lfX'mfX^fX&{fX~O$k1gO~O'P)YO|&sX}&sX~O|-qO}'za~O[1oO~O]1rO~OR1tO!O&[O!j1sO$v1cO~O^$WO&{$WO~P!$ZO!O#uO~P!$ZO|1yO!t1{O}'wX~O}1|O~Ot(^O!R2VO!S2OO!T2OO!m2UO!n2TO!o2TO!q2SO'Q$^O'Y(YO~O}2RO~P$6bOR2^O!O.[O!j2]O$v2[O~OR2^O!O.[O!j2]O$v2[O'U$_O~O'P(lO|&rX}&rX~O|.hO}'xa~O'Y2gO~O]2iO~O[2kO~O!^2nO~P)rO^2pO~O^2pO~P)rO#X2rO%i2sO~PEYO_/OO}2wO%w.}O~P]O!W2yO~O%|2zOP%yqQ%yqW%yq]%yq^%yqa%yqb%yqg%yqi%yqj%yqk%yqm%yqo%yqt%yqv%yqw%yqx%yq!O%yq!Y%yq!_%yq!b%yq!c%yq!d%yq!e%yq!f%yq!i%yq#j%yq#n%yq$u%yq$w%yq$y%yq$z%yq${%yq%O%yq%Q%yq%T%yq%U%yq%W%yq%e%yq%k%yq%m%yq%o%yq%q%yq%t%yq%z%yq&O%yq&Q%yq&S%yq&U%yq&W%yq&v%yq'P%yq']%yq'q%yq}%yq%r%yq_%yq%w%yq~O|!{i}!{i~P#JoO!t2|O|!{i}!{i~O|!Qi}!Qi~P#JoO^$WO!t3TO&{$WO~O^$WO!W!tO!t3TO&{$WO~O^$WO!W!tO!_$TO!e3XO!t3TO&{$WO'U$_O'e&gO~O!S3YO!T3YO'Q$^O'Y(YO~O!R3]O!S3YO!T3YO!q3^O'Q$^O'Y(YO~O^$WO!W!tO!e3XO!t3TO&{$WO'e&gO~O|'gq!^'gq^'gq&{'gq~P!)uO|&lO!^'fq~O#O$niP$niY$ni^$nii$nir$ni![$ni!]$ni!_$ni!e$ni#R$ni#S$ni#T$ni#U$ni#V$ni#W$ni#X$ni#Y$ni#Z$ni#]$ni#_$ni#`$ni&{$ni']$ni!^$niy$ni!O$ni$v$ni'_$ni!W$ni~P$.YO#O$piP$piY$pi^$pii$pir$pi![$pi!]$pi!_$pi!e$pi#R$pi#S$pi#T$pi#U$pi#V$pi#W$pi#X$pi#Y$pi#Z$pi#]$pi#_$pi#`$pi&{$pi']$pi!^$piy$pi!O$pi$v$pi'_$pi!W$pi~P$.{O#O$]iP$]iY$]i^$]ii$]ir$]i|$]i![$]i!]$]i!_$]i!e$]i#R$]i#S$]i#T$]i#U$]i#V$]i#W$]i#X$]i#Y$]i#Z$]i#]$]i#_$]i#`$]i&{$]i']$]i!^$]iy$]i!O$]i!t$]i$v$]i'_$]i!W$]i~P!$ZO|&_a'T&_a~P!$ZO|&`a!^&`a~P!)uO|+{O!^'di~OP#ZOq!xOr!xOt!yO!]!vO!_!wO!e#ZO']QOY#Qii#Qi![#Qi#S#Qi#T#Qi#U#Qi#V#Qi#W#Qi#X#Qi#Y#Qi#Z#Qi#]#Qi#_#Qi#`#Qi#c#Qi'e#Qi'l#Qi'm#Qi|#Qi}#Qi~O#R#Qi~P$GzOP#ZOq!xOr!xOt!yO!]!vO!_!wO!e#ZO']QOY#Qii#Qi![#Qi#S#Qi#T#Qi#U#Qi#V#Qi#W#Qi#X#Qi#Y#Qi#Z#Qi#]#Qi#_#Qi#`#Qi#c#Qi'e#Qi'l#Qi'm#Qi~O#R8_O~P$I{OP#ZOq!xOr!xOt!yO!]!vO!_!wO!e#ZO#R8_O#S8`O#T8`O#U8`O']QOY#Qi![#Qi#V#Qi#W#Qi#X#Qi#Y#Qi#Z#Qi#]#Qi#_#Qi#`#Qi#c#Qi'e#Qi'l#Qi'm#Qi~Oi#Qi~P$KvOi8aO~P$KvOP#ZOi8aOq!xOr!xOt!yO!]!vO!_!wO!e#ZO#R8_O#S8`O#T8`O#U8`O#V8bO']QO#Z#Qi#]#Qi#_#Qi#`#Qi#c#Qi'e#Qi'l#Qi'm#Qi~OY#Qi![#Qi#W#Qi#X#Qi#Y#Qi~P$MxOY8lO![8cO#W8cO#X8cO#Y8cO~P$MxOP#ZOY8lOi8aOq!xOr!xOt!yO![8cO!]!vO!_!wO!e#ZO#R8_O#S8`O#T8`O#U8`O#V8bO#W8cO#X8cO#Y8cO#Z8dO']QO#]#Qi#_#Qi#`#Qi#c#Qi'e#Qi'm#Qi~O'l#Qi~P%!WO'l!zO~P%!WOP#ZOY8lOi8aOq!xOr!xOt!yO![8cO!]!vO!_!wO!e#ZO#R8_O#S8`O#T8`O#U8`O#V8bO#W8cO#X8cO#Y8cO#Z8dO#]8fO']QO'l!zO#_#Qi#`#Qi#c#Qi'e#Qi~O'm#Qi~P%$YO'm!{O~P%$YOP#ZOY8lOi8aOq!xOr!xOt!yO![8cO!]!vO!_!wO!e#ZO#R8_O#S8`O#T8`O#U8`O#V8bO#W8cO#X8cO#Y8cO#Z8dO#]8fO#_8hO']QO'l!zO'm!{O~O#`#Qi#c#Qi'e#Qi~P%&[O^#ay|#ay&{#ayy#ay!^#ay'_#ay!O#ay$v#ay!W#ay~P!)uOP#QiY#Qii#Qir#Qi![#Qi!]#Qi!_#Qi!e#Qi#R#Qi#S#Qi#T#Qi#U#Qi#V#Qi#W#Qi#X#Qi#Y#Qi#Z#Qi#]#Qi#_#Qi#`#Qi#c#Qi']#Qi|#Qi}#Qi~P!$ZO!]!vOP'XXY'XXi'XXq'XXr'XXt'XX!['XX!_'XX!e'XX#R'XX#S'XX#T'XX#U'XX#V'XX#W'XX#X'XX#Y'XX#Z'XX#]'XX#_'XX#`'XX#c'XX']'XX'e'XX'l'XX'm'XX|'XX}'XX~O#c#di~P#*WO}3mO~O|&ga}&ga#c&ga~P#JoO!W!tO'e&gO|&ha!^&ha~O|,iO!^'ri~O|,iO!W!tO!^'ri~Oy&ja|&ja~P!$ZO!W3tO~O|,pOy'si~P!$ZO|,pOy'si~Oy3zO~O!W!tO#X4QO~Oi4RO!W!tO'e&gO~Oy4TO~O'T$_q|$_q#c$_q!t$_q~P!$ZO|1`O!O'ta~O!O&[O$v4YO~O!O&[O$v4YO~P!$ZOY4]O~O|-qO}'zi~O]4_O~O[4`O~O'Y&zO|&oX}&oX~O|1yO}'wa~O}4mO~P$6bO!R4pO!S4oO!T4oO!q/oO'Q$^O'Y(YO~O!n4qO!o4qO~P%1OO!S4oO!T4oO'Q$^O'Y(YO~O!O.[O~O!O.[O$v4sO~O!O.[O$v4sO~P!$ZOR4xO!O.[O!j4wO$v4sO~OY4}O|&ra}&ra~O|.hO}'xi~O]5QO~O!^5RO~O!^5SO~O!^5TO~O!^5TO~P)rO^5VO~O!W5YO~O!^5[O~O|'ji}'ji~P#JoO^$WO&{$WO~P#>gO^$WO!t5aO&{$WO~O^$WO!W!tO!t5aO&{$WO~O^$WO!W!tO!e5fO!t5aO&{$WO'e&gO~O!_$TO'U$_O~P%5RO!S5gO!T5gO'Q$^O'Y(YO~O|'gy!^'gy^'gy&{'gy~P!)uO#O$_qP$_qY$_q^$_qi$_qr$_q|$_q![$_q!]$_q!_$_q!e$_q#R$_q#S$_q#T$_q#U$_q#V$_q#W$_q#X$_q#Y$_q#Z$_q#]$_q#_$_q#`$_q&{$_q']$_q!^$_qy$_q!O$_q!t$_q$v$_q'_$_q!W$_q~P!$ZO|&`i!^&`i~P!)uO|8jO#c$Ra~P#GwOq-XOr-XOt-YOPnaYnaina![na!]na!_na!ena#Rna#Sna#Tna#Una#Vna#Wna#Xna#Yna#Zna#]na#_na#`na#cna']na'ena'lna'mna|na}na~Oq'tOt'uOP$SaY$Sai$Sar$Sa![$Sa!]$Sa!_$Sa!e$Sa#R$Sa#S$Sa#T$Sa#U$Sa#V$Sa#W$Sa#X$Sa#Y$Sa#Z$Sa#]$Sa#_$Sa#`$Sa#c$Sa']$Sa'e$Sa'l$Sa'm$Sa|$Sa}$Sa~Oq'tOt'uOP$UaY$Uai$Uar$Ua![$Ua!]$Ua!_$Ua!e$Ua#R$Ua#S$Ua#T$Ua#U$Ua#V$Ua#W$Ua#X$Ua#Y$Ua#Z$Ua#]$Ua#_$Ua#`$Ua#c$Ua']$Ua'e$Ua'l$Ua'm$Ua|$Ua}$Ua~OP$daY$dai$dar$da![$da!]$da!_$da!e$da#R$da#S$da#T$da#U$da#V$da#W$da#X$da#Y$da#Z$da#]$da#_$da#`$da#c$da']$da|$da}$da~P!$ZO#c$Oq~P#*WO}5oO~O'T$ry|$ry#c$ry!t$ry~P!$ZO!W!tO|&hi!^&hi~O!W!tO'e&gO|&hi!^&hi~O|,iO!^'rq~Oy&ji|&ji~P!$ZO|,pOy'sq~Oy5vO~P!$ZOy5vO~O|'Wy'T'Wy~P!$ZO|&ma!O&ma~P!$ZO!O$jq^$jq&{$jq~P!$ZO|-qO}'zq~O]6PO~O!O&[O$v6QO~O!O&[O$v6QO~P!$ZO!t6RO|&oa}&oa~O|1yO}'wi~P#JoO!S6XO!T6XO'Q$^O'Y(YO~O!R6ZO!S6XO!T6XO!q3^O'Q$^O'Y(YO~O!O.[O$v6^O~O!O.[O$v6^O~P!$ZO'Y6dO~O|.hO}'xq~O!^6gO~O!^6gO~P)rO!^6iO~O!^6jO~O|!{y}!{y~P#JoO^$WO!t6oO&{$WO~O^$WO!W!tO!t6oO&{$WO~O^$WO!W!tO!e6sO!t6oO&{$WO'e&gO~O#O$ryP$ryY$ry^$ryi$ryr$ry|$ry![$ry!]$ry!_$ry!e$ry#R$ry#S$ry#T$ry#U$ry#V$ry#W$ry#X$ry#Y$ry#Z$ry#]$ry#_$ry#`$ry&{$ry']$ry!^$ryy$ry!O$ry!t$ry$v$ry'_$ry!W$ry~P!$ZO#c#ay~P#*WOP$]iY$]ii$]ir$]i![$]i!]$]i!_$]i!e$]i#R$]i#S$]i#T$]i#U$]i#V$]i#W$]i#X$]i#Y$]i#Z$]i#]$]i#_$]i#`$]i#c$]i']$]i|$]i}$]i~P!$ZOq'tOt'uO'm'yOP$niY$nii$nir$ni![$ni!]$ni!_$ni!e$ni#R$ni#S$ni#T$ni#U$ni#V$ni#W$ni#X$ni#Y$ni#Z$ni#]$ni#_$ni#`$ni#c$ni']$ni'e$ni'l$ni|$ni}$ni~Oq'tOt'uOP$piY$pii$pir$pi![$pi!]$pi!_$pi!e$pi#R$pi#S$pi#T$pi#U$pi#V$pi#W$pi#X$pi#Y$pi#Z$pi#]$pi#_$pi#`$pi#c$pi']$pi'e$pi'l$pi'm$pi|$pi}$pi~O!W!tO|&hq!^&hq~O|,iO!^'ry~Oy&jq|&jq~P!$ZOy6yO~P!$ZO|1yO}'wq~O!S7UO!T7UO'Q$^O'Y(YO~O!O.[O$v7XO~O!O.[O$v7XO~P!$ZO!^7[O~O%|7]OP%y!ZQ%y!ZW%y!Z]%y!Z^%y!Za%y!Zb%y!Zg%y!Zi%y!Zj%y!Zk%y!Zm%y!Zo%y!Zt%y!Zv%y!Zw%y!Zx%y!Z!O%y!Z!Y%y!Z!_%y!Z!b%y!Z!c%y!Z!d%y!Z!e%y!Z!f%y!Z!i%y!Z#j%y!Z#n%y!Z$u%y!Z$w%y!Z$y%y!Z$z%y!Z${%y!Z%O%y!Z%Q%y!Z%T%y!Z%U%y!Z%W%y!Z%e%y!Z%k%y!Z%m%y!Z%o%y!Z%q%y!Z%t%y!Z%z%y!Z&O%y!Z&Q%y!Z&S%y!Z&U%y!Z&W%y!Z&v%y!Z'P%y!Z']%y!Z'q%y!Z}%y!Z%r%y!Z_%y!Z%w%y!Z~O^$WO!t7aO&{$WO~O^$WO!W!tO!t7aO&{$WO~OP$_qY$_qi$_qr$_q![$_q!]$_q!_$_q!e$_q#R$_q#S$_q#T$_q#U$_q#V$_q#W$_q#X$_q#Y$_q#Z$_q#]$_q#_$_q#`$_q#c$_q']$_q|$_q}$_q~P!$ZO|&oq}&oq~P#JoO^$WO!t7uO&{$WO~OP$ryY$ryi$ryr$ry![$ry!]$ry!_$ry!e$ry#R$ry#S$ry#T$ry#U$ry#V$ry#W$ry#X$ry#Y$ry#Z$ry#]$ry#_$ry#`$ry#c$ry']$ry|$ry}$ry~P!$ZO|#_O'_'ZX!^'ZX^'ZX&{'ZX~P!)uO'_'ZX~P.ZO'_ZXyZX!^ZX%iZX!OZX$vZX!WZX~P$tO!WcX!^ZX!^cX'ecX~P:tOP;TOQ;TO]bOa:{Ob!gOgbOi;TOjbOkbOm;TOo;TOtROvbOwbOxbO!OSO!Y;UO!_UO!b;TO!c;TO!d;TO!e;TO!f;TO!i!fO#j!iO#n]O'P'YO']QO'q;{O~O|8jO}$Ra~O]#nOg#zOi#oOj#nOk#nOm#{Oo8oOt#tO!O#uO!Y;PO!_#rO!}8uO#j$PO$T8qO$V8sO$Y$QO'P&rO~O}ZX}cX~P:tO|8jO#c'ZX~P#JoO#c'ZX~P#2iO#O8]O~O#O8^O~O!W!tO#O8]O~O!W!tO#O8^O~O!t8mO~O!t8vO|'jX}'jX~O!t;bO|'hX}'hX~O#O8wO~O#O8xO~O'T8|O~P!$ZO#O9RO~O#O9SO~O!W!tO#O9TO~O!W!tO#O9UO~O!W!tO#O9VO~O!^!Xa^!Xa&{!Xa~P#>gO#O9WO~O!W!tO#O8wO~O!W!tO#O8xO~O!W!tO#O9WO~OP#ZOq!xOr!xOt!yO!]!vO!_!wO!e#ZO#R9oO']QOY#Qii#Qi![#Qi!^#Qi#V#Qi#W#Qi#X#Qi#Y#Qi#Z#Qi#]#Qi#_#Qi#`#Qi'e#Qi'l#Qi'm#Qi^#Qi&{#Qi~O#S#Qi#T#Qi#U#Qi~P&3mO#S9pO#T9pO#U9pO~P&3mOP#ZOi9qOq!xOr!xOt!yO!]!vO!_!wO!e#ZO#R9oO#S9pO#T9pO#U9pO']QOY#Qi![#Qi!^#Qi#W#Qi#X#Qi#Y#Qi#Z#Qi#]#Qi#_#Qi#`#Qi'e#Qi'l#Qi'm#Qi^#Qi&{#Qi~O#V#Qi~P&5{O#V9rO~P&5{OP#ZOY#aOi9qOq!xOr!xOt!yO![9sO!]!vO!_!wO!e#ZO#R9oO#S9pO#T9pO#U9pO#V9rO#W9sO#X9sO#Y9sO']QO!^#Qi#]#Qi#_#Qi#`#Qi'e#Qi'l#Qi'm#Qi^#Qi&{#Qi~O#Z#Qi~P&8TO#Z9tO~P&8TOP#ZOY#aOi9qOq!xOr!xOt!yO![9sO!]!vO!_!wO!e#ZO#R9oO#S9pO#T9pO#U9pO#V9rO#W9sO#X9sO#Y9sO#Z9tO']QO'l!zO!^#Qi#_#Qi#`#Qi'e#Qi'm#Qi^#Qi&{#Qi~O#]#Qi~P&:]O#]9vO~P&:]OP#ZOY#aOi9qOq!xOr!xOt!yO![9sO!]!vO!_!wO!e#ZO#R9oO#S9pO#T9pO#U9pO#V9rO#W9sO#X9sO#Y9sO#Z9tO#]9vO']QO'l!zO'm!{O!^#Qi#`#Qi'e#Qi^#Qi&{#Qi~O#_#Qi~P&<eO#_9xO~P&<eO#c9XO~P#*WO!^#di^#di&{#di~P#>gO#O9YO~O#O9ZO~O#O9[O~O#O9]O~O#O9^O~O#O9_O~O#O9`O~O#O9aO~O!^$Oq^$Oq&{$Oq~P#>gO#c9bO~P!$ZO#c9cO~P!$ZO!^#ay^#ay&{#ay~P#>gOP'[XY'[Xi'[Xq'[Xr'[Xt'[X!['[X!]'[X!_'[X!e'[X#R'[X#S'[X#T'[X#U'[X#V'[X#W'[X#X'[X#Y'[X#Z'[X#]'[X#_'[X#`'[X']'[X'e'[X'l'[X'm'[X~O!t9zO#e9zO!^'[X^'[X&{'[X~P&@uO!t9zO~O'T:dO~P!$ZO#c:mO~P#*WO#O:rO~O!W!tO#O:rO~O!t;bO~O'T;cO~P!$ZO#c;dO~P#*WO!t;bO#e;bO|'[X}'[X~P#,RO|!Xa}!Xa#c!Xa~P#JoO#R;VO~P$GzOP#ZOq!xOr!xOt!yO!]!vO!_!wO!e#ZO#R;VO#S;WO#T;WO#U;WO']QOY#Qi|#Qi}#Qi![#Qi#V#Qi#W#Qi#X#Qi#Y#Qi#Z#Qi#]#Qi#_#Qi#`#Qi'e#Qi'l#Qi'm#Qi#c#Qi~Oi#Qi~P&DwOi;XO~P&DwOP#ZOi;XOq!xOr!xOt!yO!]!vO!_!wO!e#ZO#R;VO#S;WO#T;WO#U;WO#V;YO']QO|#Qi}#Qi#Z#Qi#]#Qi#_#Qi#`#Qi'e#Qi'l#Qi'm#Qi#c#Qi~OY#Qi![#Qi#W#Qi#X#Qi#Y#Qi~P&GPOY8lO![;ZO#W;ZO#X;ZO#Y;ZO~P&GPOP#ZOY8lOi;XOq!xOr!xOt!yO![;ZO!]!vO!_!wO!e#ZO#R;VO#S;WO#T;WO#U;WO#V;YO#W;ZO#X;ZO#Y;ZO#Z;[O']QO|#Qi}#Qi#]#Qi#_#Qi#`#Qi'e#Qi'm#Qi#c#Qi~O'l#Qi~P&IeO'l!zO~P&IeOP#ZOY8lOi;XOq!xOr!xOt!yO![;ZO!]!vO!_!wO!e#ZO#R;VO#S;WO#T;WO#U;WO#V;YO#W;ZO#X;ZO#Y;ZO#Z;[O#];^O']QO'l!zO|#Qi}#Qi#_#Qi#`#Qi'e#Qi#c#Qi~O'm#Qi~P&KmO'm!{O~P&KmOP#ZOY8lOi;XOq!xOr!xOt!yO![;ZO!]!vO!_!wO!e#ZO#R;VO#S;WO#T;WO#U;WO#V;YO#W;ZO#X;ZO#Y;ZO#Z;[O#];^O#_;`O']QO'l!zO'm!{O~O|#Qi}#Qi#`#Qi'e#Qi#c#Qi~P&MuO|#di}#di#c#di~P#JoO|$Oq}$Oq#c$Oq~P#JoO|#ay}#ay#c#ay~P#JoO#n~!]!m!o!|!}'q$T$V$Y$k$u$v$w%O%Q%T%U%W%Y~TS#n'q#p'Y'P&}#Sx~",
  goto: "$!x(OPPPPPPP(PP(aP)|PPPP._PP.t4x6k7QP7QPPP7QP7QP8oPP8tP9]PPPP?RPPPP?RBoPPPBuDxP?RPGgPPPPIv?RPPPPPLW?RPP!!T!#QPPP!#UP!#^!$_P?R?R!'x!+y!1w!1w!6WPPP!6_?RPPPPPPPPP!:TP!;uPP?R!=_P?RP?R?R?R?RP?R!?zPP!CoP!G`!Gh!Gl!GlP!ClP!Gp!GpP!KaP!Ke?R?R!Kk# _7QP7QP7Q7QP#!v7Q7Q#$l7Q7Q7Q#&o7Q7Q#']#)W#)W#)[#)W#)dP#)WP7Q#*`7Q#+k7Q7Q._PPP#,yPPP#-c#-cP#-cP#-x#-cPP#.OP#-uP#-u#.b!#Y#-u#/P#/V#/Y(P#/](PP#/d#/d#/dP(PP(PP(PP(PPP(PP#/j#/mP#/m(PPPP(PP(PP(PP(PP(PP(P(P#/q#/{#0R#0a#0g#0m#0w#0}#1X#1_#1m#1s#1y#2a#2v#4Z#4i#4o#4u#4{#5R#5]#5c#5i#5s#5}#6TPPPPPPPP#6ZPP#6}#:{PP#<`#<i#<sP#AS#DVP#K}P#LR#LU#LX#Ld#LgP#Lj#Ln#M]#NQ#NU#NhPP#Nl#Nr#NvP#Ny#N}$ Q$ p$!W$!]$!`$!c$!i$!l$!p$!tmgOSi{!k$V%^%a%b%d*_*d.x.{Q$dlQ$knQ%UwS&O!`*zQ&c!gS(]#u(bQ)W$eQ)d$mQ*O%OQ+Q&VS+W&[+YQ+h&dQ-P(dQ.g*PU/l+[+]+^S2O.[2QS3Y/n/pU4o2T2U2VQ5g3]S6X4p4qR7U6Z$hZORSTUij{!Q!U!Z!^!k!s!w!y!|!}#O#P#Q#R#S#T#U#V#W#_#b$V%V%Y%^%`%a%b%d%h%s%{&W&^&h&t&x's(u(|*Z*_*d+c+j+{,R-Y-^-f-p._.p.q.r.t.x.{.}/y0T1s2]2p2r2s4w5V8^8x9U9]9`x'[#Y8V8Y8_8`8a8b8c8d8e8f8g8h8i8m8|9X:|;k;|Q(m#|Q)]$gQ*Q%RQ*X%ZQ+s8nQ-k)QQ.o*VQ1l-qQ2e.hQ3g8o!O:s$i/f3T5a6o7a7u9i9j9o9p9q9r9s9t9u9v9w9x9y9z:d:m!q;l#h&P'm*s*v,W/]0g1{2|6R8]8j8v8w9T9V9W9[9^9_9a:r;T;U;V;W;X;Y;Z;[;];^;_;`;a;b;c;dpdOSiw{!k$V%T%^%a%b%d*_*d.x.{R*S%V(WVOSTijm{!Q!U!Z!h!k!s!w!y!|!}#O#P#Q#R#S#T#U#V#W#Y#_#b#h$V$i%V%Y%Z%^%`%a%b%d%h%s%{&W&^&h&t&x'm's(u(|*Z*_*d*s*v+c+j+{,R,W-Y-^-f-p._.p.q.r.t.x.{.}/]/f/y0T0g1s1{2]2p2r2s2|3T4w5V5a6R6o7a7u8V8Y8]8^8_8`8a8b8c8d8e8f8g8h8i8j8m8v8w8x8|9T9U9V9W9X9[9]9^9_9`9a9i9j9o9p9q9r9s9t9u9v9w9x9y9z:d:m:r:y:z:{:|;T;U;V;W;X;Y;Z;[;];^;_;`;a;b;c;d;k;|W!aRU!^&PQ$]kQ$clS$hn$mv$rpq!o!r$T$p&X&l&o)h)i)j*]*t+T+m+o/R/{Q$zuQ&`!fQ&b!gS(P#r(ZS)V$d$eQ)Z$gQ)g$oQ)y$|Q)}%OS+g&c&dQ,m(QQ-o)WQ-u)^Q-x)bQ.b)zS.f*O*PQ/w+hQ0o,iQ1k-qQ1n-tQ1q-zQ2d.gQ3q0pR5}4]!W$al!g$c$d$e%}&b&c&d([)V)W*w+V+g+h,y-o/b/i/m/w1U3W3[5e6rQ)O$]Q)o$wQ)r$xQ)|%OQ-|)gQ.a)yU.e)}*O*PQ2_.bS2c.f.gQ4j1}Q4|2dS6V4k4nS7S6W6YQ7l7TR7z7m[#x`$_(j:u;j;{S$wr%TQ$xsQ$ytR)m$u$X#w`!t!v#a#r#t#}$O$S&_'x'z'{(S(W(h(i({(})Q)n)q+d+x,p,r-[-e-g.R.U.^.`0n0w1R1Y1`1c1g1t2[2^3t4Q4Y4s4x6Q6^7X8l8p8q8r8s8t8u8}9O9P9Q9R9S9Y9Z9b9c:u;R;S;j;{V(n#|8n8oU&S!`$q*}Q&{!xQ)a$jQ,`'tQ.V)sQ1Z-XR4f1y(UbORSTUij{!Q!U!Z!^!k!s!w!y!|!}#O#P#Q#R#S#T#U#V#W#Y#_#b#h$V$i%V%Y%Z%^%`%a%b%d%h%s%{&P&W&^&h&t&x'm's(u(|*Z*_*d*s*v+c+j+{,R,W-Y-^-f-p._.p.q.r.t.x.{.}/]/f/y0T0g1s1{2]2p2r2s2|3T4w5V5a6R6o7a7u8V8Y8]8^8_8`8a8b8c8d8e8f8g8h8i8j8m8v8w8x8|9T9U9V9W9X9[9]9^9_9`9a9i9j9o9p9q9r9s9t9u9v9w9x9y9z:d:m:r:|;T;U;V;W;X;Y;Z;[;];^;_;`;a;b;c;d;k;|%]#^Y!]!l$Z%r%v&w&}'O'P'Q'R'S'T'U'V'W'X'Z'^'a'k)`*o*x+R+i+},Q,S,_/W/Z/x0S0W0X0Y0Z0[0]0^0_0`0a0b0c0f0k3O3R3b3e3k4h5]5`5k6m7O7_7s7}8W8X8z8{:R:W:X:Y:Z:[:]:^:_:`:a:b:c:n:q;Q;i;m;n;o;p;q;r;s;t;u;v;w;x;y;z(VbORSTUij{!Q!U!Z!^!k!s!w!y!|!}#O#P#Q#R#S#T#U#V#W#Y#_#b#h$V$i%V%Y%Z%^%`%a%b%d%h%s%{&P&W&^&h&t&x'm's(u(|*Z*_*d*s*v+c+j+{,R,W-Y-^-f-p._.p.q.r.t.x.{.}/]/f/y0T0g1s1{2]2p2r2s2|3T4w5V5a6R6o7a7u8V8Y8]8^8_8`8a8b8c8d8e8f8g8h8i8j8m8v8w8x8|9T9U9V9W9X9[9]9^9_9`9a9i9j9o9p9q9r9s9t9u9v9w9x9y9z:d:m:r:|;T;U;V;W;X;Y;Z;[;];^;_;`;a;b;c;d;k;|Q&Q!`R/^*zY%z!`&O&V*z+QS([#u(bS+V&[+YS,y(](dQ,z(^Q-Q(eQ.X)uS/i+W+[S/m+]+^S/q+_2SQ1U-PQ1W-RQ1X-SS1}.[2QS3W/l/nQ3Z/oQ3[/pS4k2O2VS4n2T2US5e3Y3]Q5h3^S6W4o4pQ6Y4qQ6r5gS7T6X6ZR7m7UlgOSi{!k$V%^%a%b%d*_*d.x.{Q%f!OW&p!s8]8^:rQ)T$bQ)w$zQ)x${Q+e&aW+w&t8w8x9WW-](u9T9U9VQ-m)UQ.Z)vQ/P*fQ/Q*gQ/Y*uQ/u+fW1_-^9[9]9^Q1h-nW1j-p9_9`9aQ2}/[Q3Q/dQ3`/vQ4[1iQ5Z2zQ5^3PQ5b3VQ5i3aQ6k5[Q6n5cQ7`6pQ7q7]R7t7b%S#]Y!]!l%r%v&w&}'O'P'Q'R'S'T'U'V'W'X'Z'^'a'k)`*o*x+R+i+},Q,_/W/Z/x0S0W0X0Y0Z0[0]0^0_0`0a0b0c0f0k3O3R3b3e3k4h5]5`5k6m7O7_7s7}8W8X8z8{:W:X:Y:Z:[:]:^:_:`:a:b:c:n:q;Q;i;n;o;p;q;r;s;t;u;v;w;x;y;zU(g#v&s0eX(y$Z,S:R;m%S#[Y!]!l%r%v&w&}'O'P'Q'R'S'T'U'V'W'X'Z'^'a'k)`*o*x+R+i+},Q,_/W/Z/x0S0W0X0Y0Z0[0]0^0_0`0a0b0c0f0k3O3R3b3e3k4h5]5`5k6m7O7_7s7}8W8X8z8{:W:X:Y:Z:[:]:^:_:`:a:b:c:n:q;Q;i;n;o;p;q;r;s;t;u;v;w;x;y;zQ']#]W(x$Z,S:R;mR-_(y(UbORSTUij{!Q!U!Z!^!k!s!w!y!|!}#O#P#Q#R#S#T#U#V#W#Y#_#b#h$V$i%V%Y%Z%^%`%a%b%d%h%s%{&P&W&^&h&t&x'm's(u(|*Z*_*d*s*v+c+j+{,R,W-Y-^-f-p._.p.q.r.t.x.{.}/]/f/y0T0g1s1{2]2p2r2s2|3T4w5V5a6R6o7a7u8V8Y8]8^8_8`8a8b8c8d8e8f8g8h8i8j8m8v8w8x8|9T9U9V9W9X9[9]9^9_9`9a9i9j9o9p9q9r9s9t9u9v9w9x9y9z:d:m:r:|;T;U;V;W;X;Y;Z;[;];^;_;`;a;b;c;d;k;|Q%ayQ%bzQ%d|Q%e}R.w*bQ&]!fQ(z$]Q+b&`S-d)O)gS/r+`+aW1b-a-b-c-|S3_/s/tU4X1d1e1fU5{4W4b4cQ6{5|R7h6}T+X&[+YS+X&[+YT2P.[2QS&j!n.uQ,l(PQ,w([S/h+V1}Q0t,mS1O,x-QU3X/m/q4nQ3p0oS4O1V1XU5f3Z3[6YQ5q3qQ5z4RR6s5hQ!uXS&i!n.uQ(v$UQ)R$`Q)X$fQ+k&jQ,k(PQ,v([Q,{(_Q-l)SQ.c){S/g+V1}S0s,l,mS0},w-QQ1Q,zQ1T,|Q2a.dW3U/h/m/q4nQ3o0oQ3s0tS3x1O1XQ4P1WQ4z2bW5d3X3Z3[6YS5p3p3qQ5u3zQ5x4OQ6T4iQ6b4{S6q5f5hQ6u5qQ6w5vQ6z5zQ7Q6UQ7Z6cQ7c6sQ7f6yQ7j7RQ7x7kQ8P7yQ8T8QQ9m9fQ9n9gQ:S;fQ:g:OQ:h:PQ:i:QQ:j:TQ:k:UR:l:V$jWORSTUij{!Q!U!Z!^!k!s!w!y!|!}#O#P#Q#R#S#T#U#V#W#_#b$V%V%Y%Z%^%`%a%b%d%h%s%{&W&^&h&t&x's(u(|*Z*_*d+c+j+{,R-Y-^-f-p._.p.q.r.t.x.{.}/y0T1s2]2p2r2s4w5V8^8x9U9]9`S!um!hx9d#Y8V8Y8_8`8a8b8c8d8e8f8g8h8i8m8|9X:|;k;|!O9e$i/f3T5a6o7a7u9i9j9o9p9q9r9s9t9u9v9w9x9y9z:d:mQ9m:yQ9n:zQ:S:{!q;e#h&P'm*s*v,W/]0g1{2|6R8]8j8v8w9T9V9W9[9^9_9a:r;T;U;V;W;X;Y;Z;[;];^;_;`;a;b;c;d$jXORSTUij{!Q!U!Z!^!k!s!w!y!|!}#O#P#Q#R#S#T#U#V#W#_#b$V%V%Y%Z%^%`%a%b%d%h%s%{&W&^&h&t&x's(u(|*Z*_*d+c+j+{,R-Y-^-f-p._.p.q.r.t.x.{.}/y0T1s2]2p2r2s4w5V8^8x9U9]9`Q$Ua!W$`l!g$c$d$e%}&b&c&d([)V)W*w+V+g+h,y-o/b/i/m/w1U3W3[5e6rS$fm!hQ)S$aQ){%OW.d)|)}*O*PU2b.e.f.gQ4i1}S4{2c2dU6U4j4k4nQ6c4|U7R6V6W6YS7k7S7TS7y7l7mQ8Q7zx9f#Y8V8Y8_8`8a8b8c8d8e8f8g8h8i8m8|9X:|;k;|!O9g$i/f3T5a6o7a7u9i9j9o9p9q9r9s9t9u9v9w9x9y9z:d:mQ:O:vQ:P:wQ:Q:xQ:T:yQ:U:zQ:V:{!q;f#h&P'm*s*v,W/]0g1{2|6R8]8j8v8w9T9V9W9[9^9_9a:r;T;U;V;W;X;Y;Z;[;];^;_;`;a;b;c;d$b[OSTij{!Q!U!Z!k!s!w!y!|!}#O#P#Q#R#S#T#U#V#W#_#b$V%V%Y%^%`%a%b%d%h%s%{&W&^&h&t&x's(u(|*Z*_*d+c+j+{,R-Y-^-f-p._.p.q.r.t.x.{.}/y0T1s2]2p2r2s4w5V8^8x9U9]9`U!eRU!^v$rpq!o!r$T$p&X&l&o)h)i)j*]*t+T+m+o/R/{Q*Y%Zx9h#Y8V8Y8_8`8a8b8c8d8e8f8g8h8i8m8|9X:|;k;|Q9l&P!O:t$i/f3T5a6o7a7u9i9j9o9p9q9r9s9t9u9v9w9x9y9z:d:m!o;g#h'm*s*v,W/]0g1{2|6R8]8j8v8w9T9V9W9[9^9_9a:r;T;U;V;W;X;Y;Z;[;];^;_;`;a;b;c;dS&T!`$qR/`*}$hZORSTUij{!Q!U!Z!^!k!s!w!y!|!}#O#P#Q#R#S#T#U#V#W#_#b$V%V%Y%^%`%a%b%d%h%s%{&W&^&h&t&x's(u(|*Z*_*d+c+j+{,R-Y-^-f-p._.p.q.r.t.x.{.}/y0T1s2]2p2r2s4w5V8^8x9U9]9`x'[#Y8V8Y8_8`8a8b8c8d8e8f8g8h8i8m8|9X:|;k;|Q*X%Z!O:s$i/f3T5a6o7a7u9i9j9o9p9q9r9s9t9u9v9w9x9y9z:d:m!q;l#h&P'm*s*v,W/]0g1{2|6R8]8j8v8w9T9V9W9[9^9_9a:r;T;U;V;W;X;Y;Z;[;];^;_;`;a;b;c;d!Q#SY!]$Z%r%v&w'U'V'W'X'^'a*o+R+i+},_/x0S0c3b3e8W8Xh8e'Z,S0_0`0a0b0f3k5k:b;Q;in9u)`3R5`6m7_7s7}:R:^:_:`:a:c:n:qw;]'k*x/W/Z0k3O4h5]7O8z8{;m;t;u;v;w;x;y;z|#UY!]$Z%r%v&w'W'X'^'a*o+R+i+},_/x0S0c3b3e8W8Xd8g'Z,S0a0b0f3k5k:b;Q;ij9w)`3R5`6m7_7s7}:R:`:a:c:n:qs;_'k*x/W/Z0k3O4h5]7O8z8{;m;v;w;x;y;zx#YY!]$Z%r%v&w'^'a*o+R+i+},_/x0S0c3b3e8W8Xp'{#p&u(t,g,o-T-U0Q1^3n4S9{:o:p:};h`:|'Z,S0f3k5k:b;Q;i!^;R&q'`(O(U+a+v,s-`-c.Q.S/t0P0u0y1f1v1x2Y3d3u3{4U4Z4c4v5j5s5y6`Y;S0d3j5l6t7df;k)`3R5`6m7_7s7}:R:c:n:qo;|'k*x/W/Z0k3O4h5]7O8z8{;m;x;y;z(UbORSTUij{!Q!U!Z!^!k!s!w!y!|!}#O#P#Q#R#S#T#U#V#W#Y#_#b#h$V$i%V%Y%Z%^%`%a%b%d%h%s%{&P&W&^&h&t&x'm's(u(|*Z*_*d*s*v+c+j+{,R,W-Y-^-f-p._.p.q.r.t.x.{.}/]/f/y0T0g1s1{2]2p2r2s2|3T4w5V5a6R6o7a7u8V8Y8]8^8_8`8a8b8c8d8e8f8g8h8i8j8m8v8w8x8|9T9U9V9W9X9[9]9^9_9`9a9i9j9o9p9q9r9s9t9u9v9w9x9y9z:d:m:r:|;T;U;V;W;X;Y;Z;[;];^;_;`;a;b;c;d;k;|S#i_#jR0h,V(]^ORSTU_ij{!Q!U!Z!^!k!s!w!y!|!}#O#P#Q#R#S#T#U#V#W#Y#_#b#h#j$V$i%V%Y%Z%^%`%a%b%d%h%s%{&P&W&^&h&t&x'm's(u(|*Z*_*d*s*v+c+j+{,R,V,W-Y-^-f-p._.p.q.r.t.x.{.}/]/f/y0T0g1s1{2]2p2r2s2|3T4w5V5a6R6o7a7u8V8Y8]8^8_8`8a8b8c8d8e8f8g8h8i8j8m8v8w8x8|9T9U9V9W9X9[9]9^9_9`9a9i9j9o9p9q9r9s9t9u9v9w9x9y9z:d:m:r:|;T;U;V;W;X;Y;Z;[;];^;_;`;a;b;c;d;k;|S#d]#kT'd#f'hT#e]#kT'f#f'h(]_ORSTU_ij{!Q!U!Z!^!k!s!w!y!|!}#O#P#Q#R#S#T#U#V#W#Y#_#b#h#j$V$i%V%Y%Z%^%`%a%b%d%h%s%{&P&W&^&h&t&x'm's(u(|*Z*_*d*s*v+c+j+{,R,V,W-Y-^-f-p._.p.q.r.t.x.{.}/]/f/y0T0g1s1{2]2p2r2s2|3T4w5V5a6R6o7a7u8V8Y8]8^8_8`8a8b8c8d8e8f8g8h8i8j8m8v8w8x8|9T9U9V9W9X9[9]9^9_9`9a9i9j9o9p9q9r9s9t9u9v9w9x9y9z:d:m:r:|;T;U;V;W;X;Y;Z;[;];^;_;`;a;b;c;d;k;|T#i_#jQ#l_R'o#j$jaORSTUij{!Q!U!Z!^!k!s!w!y!|!}#O#P#Q#R#S#T#U#V#W#_#b$V%V%Y%Z%^%`%a%b%d%h%s%{&W&^&h&t&x's(u(|*Z*_*d+c+j+{,R-Y-^-f-p._.p.q.r.t.x.{.}/y0T1s2]2p2r2s4w5V8^8x9U9]9`x:v#Y8V8Y8_8`8a8b8c8d8e8f8g8h8i8m8|9X:|;k;|!O:w$i/f3T5a6o7a7u9i9j9o9p9q9r9s9t9u9v9w9x9y9z:d:m!q:x#h&P'm*s*v,W/]0g1{2|6R8]8j8v8w9T9V9W9[9^9_9a:r;T;U;V;W;X;Y;Z;[;];^;_;`;a;b;c;d#{cOSUi{!Q!U!k!s!y#h$V%V%Y%Z%^%`%a%b%d%h%{&^&t'm(u(|*Z*_*d+c,W-Y-^-f-p._.p.q.r.t.x.{.}0g1s2]2p2r2s4w5V8]8^8w8x9T9U9V9W9[9]9^9_9`9a:rx#v`!v#}$O$S'x'z'{(S(h(i+x-[0n1Y:u;R;S;j;{!z&s!t#a#r#t&_(W({(})Q)n)q+d,p,r-e-g.R.U.^.`0w1R1`1c1g1t2[2^3t4Q4Y4s4x6Q6^7X8p8r8t8}9P9R9Y9bQ(r$Qc0e8l8q8s8u9O9Q9S9Z9cx#s`!v#}$O$S'x'z'{(S(h(i+x-[0n1Y:u;R;S;j;{S(_#u(bQ(s$RQ,|(`!z9|!t#a#r#t&_(W({(})Q)n)q+d,p,r-e-g.R.U.^.`0w1R1`1c1g1t2[2^3t4Q4Y4s4x6Q6^7X8p8r8t8}9P9R9Y9bb9}8l8q8s8u9O9Q9S9Z9cQ:e;OR:f;PleOSi{!k$V%^%a%b%d*_*d.x.{Q(V#tQ*k%kQ*l%mR0v,p$W#w`!t!v#a#r#t#}$O$S&_'x'z'{(S(W(h(i({(})Q)n)q+d+x,p,r-[-e-g.R.U.^.`0n0w1R1Y1`1c1g1t2[2^3t4Q4Y4s4x6Q6^7X8l8p8q8r8s8t8u8}9O9P9Q9R9S9Y9Z9b9c:u;R;S;j;{Q)p$xQ.T)rQ1w.SR4e1xT(a#u(bS(a#u(bT2P.[2QQ)R$`Q,{(_Q-l)SQ.c){Q2a.dQ4z2bQ6T4iQ6b4{Q7Q6UQ7Z6cQ7j7RQ7x7kQ8P7yR8T8Qp'x#p&u(t,g,o-T-U0Q1^3n4S9{:o:p:};h!^8}&q'`(O(U+a+v,s-`-c.Q.S/t0P0u0y1f1v1x2Y3d3u3{4U4Z4c4v5j5s5y6`Z9O0d3j5l6t7dr'z#p&u(t,e,g,o-T-U0Q1^3n4S9{:o:p:};h!`9P&q'`(O(U+a+v,s-`-c.Q.S/t/}0P0u0y1f1v1x2Y3d3u3{4U4Z4c4v5j5s5y6`]9Q0d3j5l5m6t7dpdOSiw{!k$V%T%^%a%b%d*_*d.x.{Q%QvR*Z%ZpdOSiw{!k$V%T%^%a%b%d*_*d.x.{R%QvQ)t$yR.P)mqdOSiw{!k$V%T%^%a%b%d*_*d.x.{Q.])yS2Z.a.bW4r2W2X2Y2_U6]4t4u4vU7V6[6_6`Q7n7WR7{7oQ%XwR*T%TR2h.jR6e4}S$hn$mR-u)^Q%^xR*_%_R*e%eT.y*d.{QiOQ!kST$Yi!kQ!WQR%p!WQ![RU%t![%u*pQ%u!]R*p%vQ*{&QR/_*{Q+y&uR0R+yQ+|&wS0U+|0VR0V+}Q+Y&[R/j+YQ&Y!cQ*q%wT+U&Y*qQ+O&TR/a+OQ&m!pQ+l&kU+p&m+l/|R/|+qQ'h#fR,X'hQ#j_R'n#jQ#`YW'_#`*n3f8kQ*n8WS+r8X8{Q3f8zR8k'kQ,j(PW0q,j0r3r5rU0r,k,l,mS3r0s0tR5r3s#s'v#p&q&u'`(O(U(o(p(t+a+t+u+v,e,f,g,o,s-T-U-`-c.Q.S/t/}0O0P0Q0d0u0y1^1f1v1x2Y3d3h3i3j3n3u3{4S4U4Z4c4v5j5l5m5n5s5y6`6t7d9{:o:p:};hQ,q(UU0x,q0z3vQ0z,sR3v0yQ(b#uR,}(bQ(k#yR-W(kQ1a-`R4V1aQ)k$sR.O)kQ1z.VS4g1z6SR6S4hQ)v$zR.Y)vQ2Q.[R4l2QQ.i*QS2f.i5OR5O2hQ-r)ZS1m-r4^R4^1nQ)_$hR-v)_Q.{*dR2v.{WhOSi!kQ%c{Q(w$VQ*^%^Q*`%aQ*a%bQ*c%dQ.v*_S.y*d.{R2u.xQ$XfQ%g!PQ%j!RQ%l!SQ%n!TQ)f$nQ)l$tQ*S%XQ*i%iS.l*T*WQ/S*hQ/T*kQ/U*lS/e+V1}Q0{,uQ0|,vQ1S,{Q1p-yQ1u.QQ2`.cQ2j.nQ2t.wY3S/g/h/m/q4nQ3w0}Q3y1PQ3|1TQ4a1rQ4d1vQ4y2aQ5P2i[5_3R3U3X3Z3[6YQ5t3xQ5w3}Q6O4_Q6a4zQ6f5QW6l5`5d5f5hQ6v5uQ6x5xQ6|6PQ7P6TQ7Y6bU7^6m6q6sQ7e6wQ7g6zQ7i7QQ7p7ZS7r7_7cQ7v7fQ7w7jQ7|7sQ8O7xQ8R7}Q8S8PR8U8TQ$blQ&a!gU)U$c$d$eQ*u%}U+f&b&c&dQ,u([S-n)V)WQ/[*wQ/d+VS/v+g+hQ1P,yQ1i-oQ3P/bS3V/i/mQ3a/wQ3}1US5c3W3[Q6p5eR7b6rW#q`:u;j;{R)P$_Y#y`$_:u;j;{R-V(jQ#p`S&q!t)QQ&u!vQ'`#aQ(O#rQ(U#tQ(o#}Q(p$OQ(t$SQ+a&_Q+t8pQ+u8rQ+v8tQ,e'xQ,f'zQ,g'{Q,o(SQ,s(WQ-T(hQ-U(id-`({-e.^1c2[4Y4s6Q6^7XQ-c(}Q.Q)nQ.S)qQ/t+dQ/}8}Q0O9PQ0P9RQ0Q+xQ0d8lQ0u,pQ0y,rQ1^-[Q1f-gQ1v.RQ1x.UQ2Y.`Q3d9YQ3h8qQ3i8sQ3j8uQ3n0nQ3u0wQ3{1RQ4S1YQ4U1`Q4Z1gQ4c1tQ4v2^Q5j9bQ5l9SQ5m9OQ5n9QQ5s3tQ5y4QQ6`4xQ6t9ZQ7d9cQ9{:uQ:o;RQ:p;SQ:};jR;h;{lfOSi{!k$V%^%a%b%d*_*d.x.{S!mU%`Q%i!QQ%o!UW&p!s8]8^:rQ&|!yQ'l#hS*W%V%YQ*[%ZQ*h%hQ*r%{Q+`&^W+w&t8w8x9WQ,]'mW-](u9T9U9VQ-b(|Q.s*ZQ/s+cQ0j,WQ1[-YW1_-^9[9]9^Q1e-fW1j-p9_9`9aQ2X._Q2l.pQ2m.qQ2o.rQ2q.tQ2x.}Q3l0gQ4b1sQ4u2]Q5U2pQ5W2rQ5X2sQ6_4wR6h5V!vYOSUi{!Q!k!y$V%V%Y%Z%^%`%a%b%d%h%{&^(|*Z*_*d+c-Y-f._.p.q.r.t.x.{.}1s2]2p2r2s4w5VQ!]RS!lT9jQ$ZjQ%r!ZQ%v!^Q&w!wS&}!|9oQ'O!}Q'P#OQ'Q#PQ'R#QQ'S#RQ'T#SQ'U#TQ'V#UQ'W#VQ'X#WQ'Z#YQ'^#_Q'a#bW'k#h'm,W0gQ)`$iQ*o%sS*x&P/]Q+R&WQ+i&hQ+}&xS,Q8V;TQ,S8YQ,_'sQ/W*sQ/Z*vQ/x+jQ0S+{S0W8_;VQ0X8`Q0Y8aQ0Z8bQ0[8cQ0]8dQ0^8eQ0_8fQ0`8gQ0a8hQ0b8iQ0c,RQ0f8mQ0k8jQ3O8vQ3R/fQ3b/yQ3e0TQ3k8|Q4h1{Q5]2|Q5`3TQ5k9XQ6m5aQ7O6RQ7_6oQ7s7aQ7}7u[8W!U8^8x9U9]9`Y8X!s&t(u-^-pY8z8]8w9T9[9_Y8{9V9W9^9a:rQ:R9iQ:W9pQ:X9qQ:Y9rQ:Z9sQ:[9tQ:]9uQ:^9vQ:_9wQ:`9xQ:a9yQ:b:|Q:c9zQ:n:dQ:q:mQ;Q;kQ;i;|Q;m;UQ;n;WQ;o;XQ;p;YQ;q;ZQ;r;[Q;s;]Q;t;^Q;u;_Q;v;`Q;w;aQ;x;bQ;y;cR;z;dT!VQ!WR!_RR&R!`S%}!`*zS*w&O&VR/b+QR&v!vR&y!wT!qU$TS!pU$TU$spq*]S&k!o!rQ+n&lQ+q&oQ-})jS/z+m+oR3c/{[!bR!^$p&X)h+Th!nUpq!o!r$T&l&o)j+m+o/{Q.u*]Q/X*tQ2{/RT9k&P)iT!dR$pS!cR$pS%w!^)hS*y&P)iQ+S&XR/c+TT&U!`$qQ#f]R'q#kT'g#f'hR0i,VT(R#r(ZR(X#tQ-a({Q1d-eQ2W.^Q4W1cQ4t2[Q5|4YQ6[4sQ6}6QQ7W6^R7o7XlgOSi{!k$V%^%a%b%d*_*d.x.{Q%WwR*S%TV$tpq*]R.W)sR*R%RQ$lnR)e$mR)[$gT%[x%_T%]x%_T.z*d.{",
  nodeNames: "\u26A0 ArithOp ArithOp extends LineComment BlockComment Script ExportDeclaration export Star as VariableName from String ; default FunctionDeclaration async function VariableDefinition TypeParamList TypeDefinition ThisType this LiteralType ArithOp Number BooleanLiteral VoidType void TypeofType typeof MemberExpression . ?. PropertyName [ TemplateString null super RegExp ] ArrayExpression Spread , } { ObjectExpression Property async get set PropertyNameDefinition Block : NewExpression new TypeArgList CompareOp < ) ( ArgList UnaryExpression await yield delete LogicOp BitOp ParenthesizedExpression ClassExpression class extends ClassBody MethodDeclaration Privacy static abstract PropertyDeclaration readonly Optional TypeAnnotation Equals FunctionExpression ArrowFunction ParamList ParamList ArrayPattern ObjectPattern PatternProperty Privacy readonly Arrow MemberExpression BinaryExpression ArithOp ArithOp ArithOp ArithOp BitOp CompareOp in instanceof CompareOp BitOp BitOp BitOp LogicOp LogicOp ConditionalExpression LogicOp LogicOp AssignmentExpression UpdateOp PostfixExpression CallExpression TaggedTemplatExpression DynamicImport import ImportMeta JSXElement JSXSelfCloseEndTag JSXStartTag JSXSelfClosingTag JSXIdentifier JSXNamespacedName JSXMemberExpression JSXSpreadAttribute JSXAttribute JSXAttributeValue JSXEscape JSXEndTag JSXOpenTag JSXFragmentTag JSXText JSXEscape JSXStartCloseTag JSXCloseTag PrefixCast ArrowFunction TypeParamList SequenceExpression KeyofType keyof UniqueType unique ImportType InferredType infer TypeName ParenthesizedType FunctionSignature ParamList NewSignature IndexedType TupleType Label ArrayType ReadonlyType ObjectType MethodType PropertyType IndexSignature CallSignature TypePredicate is NewSignature new UnionType LogicOp IntersectionType LogicOp ConditionalType ParameterizedType ClassDeclaration abstract implements type VariableDeclaration let var const TypeAliasDeclaration InterfaceDeclaration interface EnumDeclaration enum EnumBody NamespaceDeclaration namespace module AmbientDeclaration declare GlobalDeclaration global ClassDeclaration ClassBody MethodDeclaration AmbientFunctionDeclaration ExportGroup VariableName VariableName ImportDeclaration ImportGroup ForStatement for ForSpec ForInSpec ForOfSpec of WhileStatement while WithStatement with DoStatement do IfStatement if else SwitchStatement switch SwitchBody CaseLabel case DefaultLabel TryStatement try catch finally ReturnStatement return ThrowStatement throw BreakStatement break ContinueStatement continue DebuggerStatement debugger LabeledStatement ExpressionStatement",
  maxTerm: 321,
  nodeProps: [
    [NodeProp.group, -26, 7, 14, 16, 53, 174, 178, 182, 183, 185, 188, 191, 202, 204, 210, 212, 214, 216, 219, 225, 229, 231, 233, 235, 237, 239, 240, "Statement", -30, 11, 13, 23, 26, 27, 37, 38, 39, 40, 42, 47, 55, 63, 69, 70, 83, 84, 93, 94, 109, 112, 114, 115, 116, 117, 119, 120, 138, 139, 141, "Expression", -21, 22, 24, 28, 30, 142, 144, 146, 147, 149, 150, 151, 153, 154, 155, 157, 158, 159, 168, 170, 172, 173, "Type", -2, 74, 78, "ClassItem"],
    [NodeProp.closedBy, 36, "]", 46, "}", 61, ")", 122, "JSXSelfCloseEndTag JSXEndTag", 136, "JSXEndTag"],
    [NodeProp.openedBy, 41, "[", 45, "{", 60, "(", 121, "JSXStartTag", 131, "JSXStartTag JSXStartCloseTag"]
  ],
  skippedNodes: [0, 4, 5],
  repeatNodeCount: 27,
  tokenData: "!Ck~R!ZOX$tX^%S^p$tpq%Sqr&rrs'zst$ttu/wuv2Xvw2|wx3zxy:byz:rz{;S{|<S|}<g}!O<S!O!P<w!P!QAT!Q!R!0Z!R![!2j![!]!8Y!]!^!8l!^!_!8|!_!`!9y!`!a!;U!a!b!<{!b!c$t!c!}/w!}#O!>^#O#P$t#P#Q!>n#Q#R!?O#R#S/w#S#T!?c#T#o/w#o#p!?s#p#q!?x#q#r!@`#r#s!@r#s#y$t#y#z%S#z$f$t$f$g%S$g#BY/w#BY#BZ!AS#BZ$IS/w$IS$I_!AS$I_$I|/w$I|$JO!AS$JO$JT/w$JT$JU!AS$JU$KV/w$KV$KW!AS$KW&FU/w&FU&FV!AS&FV~/wW$yR#zWO!^$t!_#o$t#p~$t,T%Zg#zW&}+{OX$tX^%S^p$tpq%Sq!^$t!_#o$t#p#y$t#y#z%S#z$f$t$f$g%S$g#BY$t#BY#BZ%S#BZ$IS$t$IS$I_%S$I_$I|$t$I|$JO%S$JO$JT$t$JT$JU%S$JU$KV$t$KV$KW%S$KW&FU$t&FU&FV%S&FV~$t$T&yS#zW!e#{O!^$t!_!`'V!`#o$t#p~$t$O'^S#Z#v#zWO!^$t!_!`'j!`#o$t#p~$t$O'qR#Z#v#zWO!^$t!_#o$t#p~$t'u(RZ#zW]!ROY'zYZ(tZr'zrs*Rs!^'z!^!_*e!_#O'z#O#P,q#P#o'z#o#p*e#p~'z&r(yV#zWOr(trs)`s!^(t!^!_)p!_#o(t#o#p)p#p~(t&r)gR#u&j#zWO!^$t!_#o$t#p~$t&j)sROr)prs)|s~)p&j*RO#u&j'u*[R#u&j#zW]!RO!^$t!_#o$t#p~$t'm*jV]!ROY*eYZ)pZr*ers+Ps#O*e#O#P+W#P~*e'm+WO#u&j]!R'm+ZROr*ers+ds~*e'm+kU#u&j]!ROY+}Zr+}rs,fs#O+}#O#P,k#P~+}!R,SU]!ROY+}Zr+}rs,fs#O+}#O#P,k#P~+}!R,kO]!R!R,nPO~+}'u,vV#zWOr'zrs-]s!^'z!^!_*e!_#o'z#o#p*e#p~'z'u-fZ#u&j#zW]!ROY.XYZ$tZr.Xrs/Rs!^.X!^!_+}!_#O.X#O#P/c#P#o.X#o#p+}#p~.X!Z.`Z#zW]!ROY.XYZ$tZr.Xrs/Rs!^.X!^!_+}!_#O.X#O#P/c#P#o.X#o#p+}#p~.X!Z/YR#zW]!RO!^$t!_#o$t#p~$t!Z/hT#zWO!^.X!^!_+}!_#o.X#o#p+}#p~.X&i0S_#zW#pS'Yp'P%kOt$ttu/wu}$t}!O1R!O!Q$t!Q![/w![!^$t!_!c$t!c!}/w!}#R$t#R#S/w#S#T$t#T#o/w#p$g$t$g~/w[1Y_#zW#pSOt$ttu1Ru}$t}!O1R!O!Q$t!Q![1R![!^$t!_!c$t!c!}1R!}#R$t#R#S1R#S#T$t#T#o1R#p$g$t$g~1R$O2`S#T#v#zWO!^$t!_!`2l!`#o$t#p~$t$O2sR#zW#e#vO!^$t!_#o$t#p~$t%r3TU'm%j#zWOv$tvw3gw!^$t!_!`2l!`#o$t#p~$t$O3nS#zW#_#vO!^$t!_!`2l!`#o$t#p~$t'u4RZ#zW]!ROY3zYZ4tZw3zwx*Rx!^3z!^!_5l!_#O3z#O#P7l#P#o3z#o#p5l#p~3z&r4yV#zWOw4twx)`x!^4t!^!_5`!_#o4t#o#p5`#p~4t&j5cROw5`wx)|x~5`'m5qV]!ROY5lYZ5`Zw5lwx+Px#O5l#O#P6W#P~5l'm6ZROw5lwx6dx~5l'm6kU#u&j]!ROY6}Zw6}wx,fx#O6}#O#P7f#P~6}!R7SU]!ROY6}Zw6}wx,fx#O6}#O#P7f#P~6}!R7iPO~6}'u7qV#zWOw3zwx8Wx!^3z!^!_5l!_#o3z#o#p5l#p~3z'u8aZ#u&j#zW]!ROY9SYZ$tZw9Swx/Rx!^9S!^!_6}!_#O9S#O#P9|#P#o9S#o#p6}#p~9S!Z9ZZ#zW]!ROY9SYZ$tZw9Swx/Rx!^9S!^!_6}!_#O9S#O#P9|#P#o9S#o#p6}#p~9S!Z:RT#zWO!^9S!^!_6}!_#o9S#o#p6}#p~9S%V:iR!_$}#zWO!^$t!_#o$t#p~$tZ:yR!^R#zWO!^$t!_#o$t#p~$t%R;]U'Q!R#U#v#zWOz$tz{;o{!^$t!_!`2l!`#o$t#p~$t$O;vS#R#v#zWO!^$t!_!`2l!`#o$t#p~$t$u<ZSi$m#zWO!^$t!_!`2l!`#o$t#p~$t&i<nR|&a#zWO!^$t!_#o$t#p~$t&i=OVq%n#zWO!O$t!O!P=e!P!Q$t!Q![>Z![!^$t!_#o$t#p~$ty=jT#zWO!O$t!O!P=y!P!^$t!_#o$t#p~$ty>QR{q#zWO!^$t!_#o$t#p~$ty>bZ#zWjqO!Q$t!Q![>Z![!^$t!_!g$t!g!h?T!h#R$t#R#S>Z#S#X$t#X#Y?T#Y#o$t#p~$ty?YZ#zWO{$t{|?{|}$t}!O?{!O!Q$t!Q![@g![!^$t!_#R$t#R#S@g#S#o$t#p~$ty@QV#zWO!Q$t!Q![@g![!^$t!_#R$t#R#S@g#S#o$t#p~$ty@nV#zWjqO!Q$t!Q![@g![!^$t!_#R$t#R#S@g#S#o$t#p~$t,TA[`#zW#S#vOYB^YZ$tZzB^z{HT{!PB^!P!Q!*|!Q!^B^!^!_Da!_!`!+u!`!a!,t!a!}B^!}#O!-s#O#P!/o#P#oB^#o#pDa#p~B^XBe[#zWxPOYB^YZ$tZ!PB^!P!QCZ!Q!^B^!^!_Da!_!}B^!}#OFY#O#PGi#P#oB^#o#pDa#p~B^XCb_#zWxPO!^$t!_#Z$t#Z#[CZ#[#]$t#]#^CZ#^#a$t#a#bCZ#b#g$t#g#hCZ#h#i$t#i#jCZ#j#m$t#m#nCZ#n#o$t#p~$tPDfVxPOYDaZ!PDa!P!QD{!Q!}Da!}#OEd#O#PFP#P~DaPEQUxP#Z#[D{#]#^D{#a#bD{#g#hD{#i#jD{#m#nD{PEgTOYEdZ#OEd#O#PEv#P#QDa#Q~EdPEyQOYEdZ~EdPFSQOYDaZ~DaXF_Y#zWOYFYYZ$tZ!^FY!^!_Ed!_#OFY#O#PF}#P#QB^#Q#oFY#o#pEd#p~FYXGSV#zWOYFYYZ$tZ!^FY!^!_Ed!_#oFY#o#pEd#p~FYXGnV#zWOYB^YZ$tZ!^B^!^!_Da!_#oB^#o#pDa#p~B^,TH[^#zWxPOYHTYZIWZzHTz{Ki{!PHT!P!Q!)j!Q!^HT!^!_Mt!_!}HT!}#O!%e#O#P!(x#P#oHT#o#pMt#p~HT,TI]V#zWOzIWz{Ir{!^IW!^!_Jt!_#oIW#o#pJt#p~IW,TIwX#zWOzIWz{Ir{!PIW!P!QJd!Q!^IW!^!_Jt!_#oIW#o#pJt#p~IW,TJkR#zWT+{O!^$t!_#o$t#p~$t+{JwROzJtz{KQ{~Jt+{KTTOzJtz{KQ{!PJt!P!QKd!Q~Jt+{KiOT+{,TKp^#zWxPOYHTYZIWZzHTz{Ki{!PHT!P!QLl!Q!^HT!^!_Mt!_!}HT!}#O!%e#O#P!(x#P#oHT#o#pMt#p~HT,TLu_#zWT+{xPO!^$t!_#Z$t#Z#[CZ#[#]$t#]#^CZ#^#a$t#a#bCZ#b#g$t#g#hCZ#h#i$t#i#jCZ#j#m$t#m#nCZ#n#o$t#p~$t+{MyYxPOYMtYZJtZzMtz{Ni{!PMt!P!Q!$a!Q!}Mt!}#O! w#O#P!#}#P~Mt+{NnYxPOYMtYZJtZzMtz{Ni{!PMt!P!Q! ^!Q!}Mt!}#O! w#O#P!#}#P~Mt+{! eUT+{xP#Z#[D{#]#^D{#a#bD{#g#hD{#i#jD{#m#nD{+{! zWOY! wYZJtZz! wz{!!d{#O! w#O#P!#k#P#QMt#Q~! w+{!!gYOY! wYZJtZz! wz{!!d{!P! w!P!Q!#V!Q#O! w#O#P!#k#P#QMt#Q~! w+{!#[TT+{OYEdZ#OEd#O#PEv#P#QDa#Q~Ed+{!#nTOY! wYZJtZz! wz{!!d{~! w+{!$QTOYMtYZJtZzMtz{Ni{~Mt+{!$f_xPOzJtz{KQ{#ZJt#Z#[!$a#[#]Jt#]#^!$a#^#aJt#a#b!$a#b#gJt#g#h!$a#h#iJt#i#j!$a#j#mJt#m#n!$a#n~Jt,T!%j[#zWOY!%eYZIWZz!%ez{!&`{!^!%e!^!_! w!_#O!%e#O#P!(W#P#QHT#Q#o!%e#o#p! w#p~!%e,T!&e^#zWOY!%eYZIWZz!%ez{!&`{!P!%e!P!Q!'a!Q!^!%e!^!_! w!_#O!%e#O#P!(W#P#QHT#Q#o!%e#o#p! w#p~!%e,T!'hY#zWT+{OYFYYZ$tZ!^FY!^!_Ed!_#OFY#O#PF}#P#QB^#Q#oFY#o#pEd#p~FY,T!(]X#zWOY!%eYZIWZz!%ez{!&`{!^!%e!^!_! w!_#o!%e#o#p! w#p~!%e,T!(}X#zWOYHTYZIWZzHTz{Ki{!^HT!^!_Mt!_#oHT#o#pMt#p~HT,T!)qc#zWxPOzIWz{Ir{!^IW!^!_Jt!_#ZIW#Z#[!)j#[#]IW#]#^!)j#^#aIW#a#b!)j#b#gIW#g#h!)j#h#iIW#i#j!)j#j#mIW#m#n!)j#n#oIW#o#pJt#p~IW,T!+TV#zWS+{OY!*|YZ$tZ!^!*|!^!_!+j!_#o!*|#o#p!+j#p~!*|+{!+oQS+{OY!+jZ~!+j$P!,O[#zW#e#vxPOYB^YZ$tZ!PB^!P!QCZ!Q!^B^!^!_Da!_!}B^!}#OFY#O#PGi#P#oB^#o#pDa#p~B^]!,}[#mS#zWxPOYB^YZ$tZ!PB^!P!QCZ!Q!^B^!^!_Da!_!}B^!}#OFY#O#PGi#P#oB^#o#pDa#p~B^X!-xY#zWOY!-sYZ$tZ!^!-s!^!_!.h!_#O!-s#O#P!/T#P#QB^#Q#o!-s#o#p!.h#p~!-sP!.kTOY!.hZ#O!.h#O#P!.z#P#QDa#Q~!.hP!.}QOY!.hZ~!.hX!/YV#zWOY!-sYZ$tZ!^!-s!^!_!.h!_#o!-s#o#p!.h#p~!-sX!/tV#zWOYB^YZ$tZ!^B^!^!_Da!_#oB^#o#pDa#p~B^y!0bd#zWjqO!O$t!O!P!1p!P!Q$t!Q![!2j![!^$t!_!g$t!g!h?T!h#R$t#R#S!2j#S#U$t#U#V!4Q#V#X$t#X#Y?T#Y#b$t#b#c!3p#c#d!5`#d#l$t#l#m!6h#m#o$t#p~$ty!1wZ#zWjqO!Q$t!Q![!1p![!^$t!_!g$t!g!h?T!h#R$t#R#S!1p#S#X$t#X#Y?T#Y#o$t#p~$ty!2q_#zWjqO!O$t!O!P!1p!P!Q$t!Q![!2j![!^$t!_!g$t!g!h?T!h#R$t#R#S!2j#S#X$t#X#Y?T#Y#b$t#b#c!3p#c#o$t#p~$ty!3wR#zWjqO!^$t!_#o$t#p~$ty!4VW#zWO!Q$t!Q!R!4o!R!S!4o!S!^$t!_#R$t#R#S!4o#S#o$t#p~$ty!4vW#zWjqO!Q$t!Q!R!4o!R!S!4o!S!^$t!_#R$t#R#S!4o#S#o$t#p~$ty!5eV#zWO!Q$t!Q!Y!5z!Y!^$t!_#R$t#R#S!5z#S#o$t#p~$ty!6RV#zWjqO!Q$t!Q!Y!5z!Y!^$t!_#R$t#R#S!5z#S#o$t#p~$ty!6mZ#zWO!Q$t!Q![!7`![!^$t!_!c$t!c!i!7`!i#R$t#R#S!7`#S#T$t#T#Z!7`#Z#o$t#p~$ty!7gZ#zWjqO!Q$t!Q![!7`![!^$t!_!c$t!c!i!7`!i#R$t#R#S!7`#S#T$t#T#Z!7`#Z#o$t#p~$t%w!8cR!WV#zW#c%hO!^$t!_#o$t#p~$t!P!8sR^w#zWO!^$t!_#o$t#p~$t+c!9XR'Ud![%Y#n&s'qP!P!Q!9b!^!_!9g!_!`!9tW!9gO#|W#v!9lP#V#v!_!`!9o#v!9tO#e#v#v!9yO#W#v%w!:QT!t%o#zWO!^$t!_!`!:a!`!a!:t!a#o$t#p~$t$O!:hS#Z#v#zWO!^$t!_!`'j!`#o$t#p~$t$P!:{R#O#w#zWO!^$t!_#o$t#p~$t%w!;aT'T!s#W#v#wS#zWO!^$t!_!`!;p!`!a!<Q!a#o$t#p~$t$O!;wR#W#v#zWO!^$t!_#o$t#p~$t$O!<XT#V#v#zWO!^$t!_!`2l!`!a!<h!a#o$t#p~$t$O!<oS#V#v#zWO!^$t!_!`2l!`#o$t#p~$t%w!=SV'e%o#zWO!O$t!O!P!=i!P!^$t!_!a$t!a!b!=y!b#o$t#p~$t$`!=pRr$W#zWO!^$t!_#o$t#p~$t$O!>QS#zW#`#vO!^$t!_!`2l!`#o$t#p~$t&e!>eRt&]#zWO!^$t!_#o$t#p~$tZ!>uRyR#zWO!^$t!_#o$t#p~$t$O!?VS#]#v#zWO!^$t!_!`2l!`#o$t#p~$t$P!?jR#zW']#wO!^$t!_#o$t#p~$t~!?xO!O~%r!@PT'l%j#zWO!^$t!_!`2l!`#o$t#p#q!=y#q~$t$u!@iR}$k#zW'_QO!^$t!_#o$t#p~$tX!@yR!fP#zWO!^$t!_#o$t#p~$t,T!Aar#zW#pS'Yp'P%k&}+{OX$tX^%S^p$tpq%Sqt$ttu/wu}$t}!O1R!O!Q$t!Q![/w![!^$t!_!c$t!c!}/w!}#R$t#R#S/w#S#T$t#T#o/w#p#y$t#y#z%S#z$f$t$f$g%S$g#BY/w#BY#BZ!AS#BZ$IS/w$IS$I_!AS$I_$I|/w$I|$JO!AS$JO$JT/w$JT$JU!AS$JU$KV/w$KV$KW!AS$KW&FU/w&FU&FV!AS&FV~/w",
  tokenizers: [noSemicolon, incdecToken, template, 0, 1, 2, 3, 4, 5, 6, 7, 8, insertSemicolon],
  topRules: {Script: [0, 6]},
  dialects: {jsx: 12773, ts: 12775},
  dynamicPrecedences: {"139": 1, "166": 1},
  specialized: [{term: 277, get: (value, stack) => tsExtends(value, stack) << 1 | 1}, {term: 277, get: (value) => spec_identifier2[value] || -1}, {term: 286, get: (value) => spec_word[value] || -1}, {term: 58, get: (value) => spec_LessThan[value] || -1}],
  tokenPrec: 12795
});

// node_modules/@codemirror/tooltip/dist/index.js
var ios = typeof navigator != "undefined" && !/Edge\/(\d+)/.exec(navigator.userAgent) && /Apple Computer/.test(navigator.vendor) && (/Mobile\/\w+/.test(navigator.userAgent) || navigator.maxTouchPoints > 2);
var Outside = "-10000px";
var tooltipPlugin = ViewPlugin.fromClass(class {
  constructor(view) {
    this.view = view;
    this.inView = true;
    this.measureReq = {read: this.readMeasure.bind(this), write: this.writeMeasure.bind(this), key: this};
    this.input = view.state.facet(showTooltip);
    this.tooltips = this.input.filter((t2) => t2);
    this.tooltipViews = this.tooltips.map((tp) => this.createTooltip(tp));
  }
  update(update) {
    let input = update.state.facet(showTooltip);
    if (input == this.input) {
      for (let t2 of this.tooltipViews)
        if (t2.update)
          t2.update(update);
    } else {
      let tooltips = input.filter((x) => x);
      let views = [];
      for (let i = 0; i < tooltips.length; i++) {
        let tip = tooltips[i], known = -1;
        if (!tip)
          continue;
        for (let i2 = 0; i2 < this.tooltips.length; i2++) {
          let other = this.tooltips[i2];
          if (other && other.create == tip.create)
            known = i2;
        }
        if (known < 0) {
          views[i] = this.createTooltip(tip);
        } else {
          let tooltipView = views[i] = this.tooltipViews[known];
          if (tooltipView.update)
            tooltipView.update(update);
        }
      }
      for (let t2 of this.tooltipViews)
        if (views.indexOf(t2) < 0)
          t2.dom.remove();
      this.input = input;
      this.tooltips = tooltips;
      this.tooltipViews = views;
      this.maybeMeasure();
    }
  }
  createTooltip(tooltip) {
    let tooltipView = tooltip.create(this.view);
    tooltipView.dom.classList.add("cm-tooltip");
    if (tooltip.class)
      tooltipView.dom.classList.add(tooltip.class);
    tooltipView.dom.style.top = Outside;
    this.view.dom.appendChild(tooltipView.dom);
    if (tooltipView.mount)
      tooltipView.mount(this.view);
    return tooltipView;
  }
  destroy() {
    for (let {dom} of this.tooltipViews)
      dom.remove();
  }
  readMeasure() {
    return {
      editor: this.view.dom.getBoundingClientRect(),
      pos: this.tooltips.map((t2) => this.view.coordsAtPos(t2.pos)),
      size: this.tooltipViews.map(({dom}) => dom.getBoundingClientRect()),
      innerWidth: window.innerWidth,
      innerHeight: window.innerHeight
    };
  }
  writeMeasure(measured) {
    let {editor} = measured;
    for (let i = 0; i < this.tooltipViews.length; i++) {
      let tooltip = this.tooltips[i], tView = this.tooltipViews[i], {dom} = tView;
      let pos = measured.pos[i], size = measured.size[i];
      if (!pos || pos.bottom <= editor.top || pos.top >= editor.bottom || pos.right <= editor.left || pos.left >= editor.right) {
        dom.style.top = Outside;
        continue;
      }
      let width = size.right - size.left, height = size.bottom - size.top;
      let left = this.view.textDirection == Direction.LTR ? Math.min(pos.left, measured.innerWidth - width) : Math.max(0, pos.left - width);
      let above = !!tooltip.above;
      if (!tooltip.strictSide && (above ? pos.top - (size.bottom - size.top) < 0 : pos.bottom + (size.bottom - size.top) > measured.innerHeight))
        above = !above;
      if (ios) {
        dom.style.top = (above ? pos.top - height : pos.bottom) - editor.top + "px";
        dom.style.left = left - editor.left + "px";
        dom.style.position = "absolute";
      } else {
        dom.style.top = (above ? pos.top - height : pos.bottom) + "px";
        dom.style.left = left + "px";
      }
      dom.classList.toggle("cm-tooltip-above", above);
      dom.classList.toggle("cm-tooltip-below", !above);
      if (tView.positioned)
        tView.positioned();
    }
  }
  maybeMeasure() {
    if (this.tooltips.length) {
      if (this.view.inView || this.inView)
        this.view.requestMeasure(this.measureReq);
      this.inView = this.view.inView;
    }
  }
}, {
  eventHandlers: {
    scroll() {
      this.maybeMeasure();
    }
  }
});
var baseTheme4 = EditorView.baseTheme({
  ".cm-tooltip": {
    position: "fixed",
    zIndex: 100
  },
  "&light .cm-tooltip": {
    border: "1px solid #ddd",
    backgroundColor: "#f5f5f5"
  },
  "&dark .cm-tooltip": {
    backgroundColor: "#333338",
    color: "white"
  }
});
var showTooltip = Facet.define({
  enables: [tooltipPlugin, baseTheme4]
});

// node_modules/@codemirror/autocomplete/dist/index.js
var CompletionContext = class {
  constructor(state, pos, explicit) {
    this.state = state;
    this.pos = pos;
    this.explicit = explicit;
    this.abortListeners = [];
  }
  tokenBefore(types4) {
    let token = syntaxTree(this.state).resolve(this.pos, -1);
    while (token && types4.indexOf(token.name) < 0)
      token = token.parent;
    return token ? {
      from: token.from,
      to: this.pos,
      text: this.state.sliceDoc(token.from, this.pos),
      type: token.type
    } : null;
  }
  matchBefore(expr) {
    let line = this.state.doc.lineAt(this.pos);
    let start = Math.max(line.from, this.pos - 250);
    let str = line.text.slice(start - line.from, this.pos - line.from);
    let found = str.search(ensureAnchor(expr, false));
    return found < 0 ? null : {from: start + found, to: this.pos, text: str.slice(found)};
  }
  get aborted() {
    return this.abortListeners == null;
  }
  addEventListener(type2, listener) {
    if (type2 == "abort" && this.abortListeners)
      this.abortListeners.push(listener);
  }
};
function toSet(chars2) {
  let flat = Object.keys(chars2).join("");
  let words4 = /\w/.test(flat);
  if (words4)
    flat = flat.replace(/\w/g, "");
  return `[${words4 ? "\\w" : ""}${flat.replace(/[^\w\s]/g, "\\$&")}]`;
}
function prefixMatch(options) {
  let first = Object.create(null), rest = Object.create(null);
  for (let {label} of options) {
    first[label[0]] = true;
    for (let i = 1; i < label.length; i++)
      rest[label[i]] = true;
  }
  let source = toSet(first) + toSet(rest) + "*$";
  return [new RegExp("^" + source), new RegExp(source)];
}
function completeFromList(list) {
  let options = list.map((o) => typeof o == "string" ? {label: o} : o);
  let [span, match] = options.every((o) => /^\w+$/.test(o.label)) ? [/\w*$/, /\w+$/] : prefixMatch(options);
  return (context) => {
    let token = context.matchBefore(match);
    return token || context.explicit ? {from: token ? token.from : context.pos, options, span} : null;
  };
}
function ifNotIn(nodes, source) {
  return (context) => {
    for (let pos = syntaxTree(context.state).resolve(context.pos, -1); pos; pos = pos.parent)
      if (nodes.indexOf(pos.name) > -1)
        return null;
    return source(context);
  };
}
var Option = class {
  constructor(completion, source, match) {
    this.completion = completion;
    this.source = source;
    this.match = match;
  }
};
function cur(state) {
  return state.selection.main.head;
}
function ensureAnchor(expr, start) {
  var _a;
  let {source} = expr;
  let addStart = start && source[0] != "^", addEnd = source[source.length - 1] != "$";
  if (!addStart && !addEnd)
    return expr;
  return new RegExp(`${addStart ? "^" : ""}(?:${source})${addEnd ? "$" : ""}`, (_a = expr.flags) !== null && _a !== void 0 ? _a : expr.ignoreCase ? "i" : "");
}
function applyCompletion(view, option) {
  let apply = option.completion.apply || option.completion.label;
  let result = option.source;
  if (typeof apply == "string") {
    view.dispatch({
      changes: {from: result.from, to: result.to, insert: apply},
      selection: {anchor: result.from + apply.length}
    });
  } else {
    apply(view, option.completion, result.from, result.to);
  }
}
var SourceCache = new WeakMap();
function asSource(source) {
  if (!Array.isArray(source))
    return source;
  let known = SourceCache.get(source);
  if (!known)
    SourceCache.set(source, known = completeFromList(source));
  return known;
}
var FuzzyMatcher = class {
  constructor(pattern) {
    this.pattern = pattern;
    this.chars = [];
    this.folded = [];
    this.any = [];
    this.precise = [];
    this.byWord = [];
    for (let p = 0; p < pattern.length; ) {
      let char = codePointAt(pattern, p), size = codePointSize(char);
      this.chars.push(char);
      let part = pattern.slice(p, p + size), upper = part.toUpperCase();
      this.folded.push(codePointAt(upper == part ? part.toLowerCase() : upper, 0));
      p += size;
    }
    this.astral = pattern.length != this.chars.length;
  }
  match(word) {
    if (this.pattern.length == 0)
      return [0];
    if (word.length < this.pattern.length)
      return null;
    let {chars: chars2, folded, any, precise, byWord} = this;
    if (chars2.length == 1) {
      let first = codePointAt(word, 0);
      return first == chars2[0] ? [0, 0, codePointSize(first)] : first == folded[0] ? [-200, 0, codePointSize(first)] : null;
    }
    let direct = word.indexOf(this.pattern);
    if (direct == 0)
      return [0, 0, this.pattern.length];
    let len = chars2.length, anyTo = 0;
    if (direct < 0) {
      for (let i = 0, e = Math.min(word.length, 200); i < e && anyTo < len; ) {
        let next = codePointAt(word, i);
        if (next == chars2[anyTo] || next == folded[anyTo])
          any[anyTo++] = i;
        i += codePointSize(next);
      }
      if (anyTo < len)
        return null;
    }
    let preciseTo = 0;
    let byWordTo = 0, byWordFolded = false;
    let adjacentTo = 0, adjacentStart = -1, adjacentEnd = -1;
    let hasLower = /[a-z]/.test(word);
    for (let i = 0, e = Math.min(word.length, 200), prevType = 0; i < e && byWordTo < len; ) {
      let next = codePointAt(word, i);
      if (direct < 0) {
        if (preciseTo < len && next == chars2[preciseTo])
          precise[preciseTo++] = i;
        if (adjacentTo < len) {
          if (next == chars2[adjacentTo] || next == folded[adjacentTo]) {
            if (adjacentTo == 0)
              adjacentStart = i;
            adjacentEnd = i;
            adjacentTo++;
          } else {
            adjacentTo = 0;
          }
        }
      }
      let ch, type2 = next < 255 ? next >= 48 && next <= 57 || next >= 97 && next <= 122 ? 2 : next >= 65 && next <= 90 ? 1 : 0 : (ch = fromCodePoint(next)) != ch.toLowerCase() ? 1 : ch != ch.toUpperCase() ? 2 : 0;
      if ((type2 == 1 && hasLower || prevType == 0 && type2 != 0) && (chars2[byWordTo] == next || folded[byWordTo] == next && (byWordFolded = true)))
        byWord[byWordTo++] = i;
      prevType = type2;
      i += codePointSize(next);
    }
    if (byWordTo == len && byWord[0] == 0)
      return this.result(-100 + (byWordFolded ? -200 : 0), byWord, word);
    if (adjacentTo == len && adjacentStart == 0)
      return [-200, 0, adjacentEnd];
    if (direct > -1)
      return [-700, direct, direct + this.pattern.length];
    if (adjacentTo == len)
      return [-200 + -700, adjacentStart, adjacentEnd];
    if (byWordTo == len)
      return this.result(-100 + (byWordFolded ? -200 : 0) + -700, byWord, word);
    return chars2.length == 2 ? null : this.result((any[0] ? -700 : 0) + -200 + -1100, any, word);
  }
  result(score2, positions, word) {
    let result = [score2], i = 1;
    for (let pos of positions) {
      let to = pos + (this.astral ? codePointSize(codePointAt(word, pos)) : 1);
      if (i > 1 && result[i - 1] == pos)
        result[i - 1] = to;
      else {
        result[i++] = pos;
        result[i++] = to;
      }
    }
    return result;
  }
};
var completionConfig = Facet.define({
  combine(configs) {
    return combineConfig(configs, {
      activateOnTyping: true,
      override: null,
      maxRenderedOptions: 100,
      defaultKeymap: true
    }, {
      defaultKeymap: (a, b) => a && b
    });
  }
});
var MaxInfoWidth = 300;
var baseTheme5 = EditorView.baseTheme({
  ".cm-tooltip.cm-tooltip-autocomplete": {
    "& > ul": {
      fontFamily: "monospace",
      whiteSpace: "nowrap",
      overflow: "auto",
      maxWidth_fallback: "700px",
      maxWidth: "min(700px, 95vw)",
      maxHeight: "10em",
      listStyle: "none",
      margin: 0,
      padding: 0,
      "& > li": {
        cursor: "pointer",
        padding: "1px 1em 1px 3px",
        lineHeight: 1.2
      },
      "& > li[aria-selected]": {
        background_fallback: "#bdf",
        backgroundColor: "Highlight",
        color_fallback: "white",
        color: "HighlightText"
      }
    }
  },
  ".cm-completionListIncompleteTop:before, .cm-completionListIncompleteBottom:after": {
    content: '"\xB7\xB7\xB7"',
    opacity: 0.5,
    display: "block",
    textAlign: "center"
  },
  ".cm-tooltip.cm-completionInfo": {
    position: "absolute",
    padding: "3px 9px",
    width: "max-content",
    maxWidth: MaxInfoWidth + "px"
  },
  ".cm-completionInfo.cm-completionInfo-left": {right: "100%"},
  ".cm-completionInfo.cm-completionInfo-right": {left: "100%"},
  "&light .cm-snippetField": {backgroundColor: "#00000022"},
  "&dark .cm-snippetField": {backgroundColor: "#ffffff22"},
  ".cm-snippetFieldPosition": {
    verticalAlign: "text-top",
    width: 0,
    height: "1.15em",
    margin: "0 -0.7px -.7em",
    borderLeft: "1.4px dotted #888"
  },
  ".cm-completionMatchedText": {
    textDecoration: "underline"
  },
  ".cm-completionDetail": {
    marginLeft: "0.5em",
    fontStyle: "italic"
  },
  ".cm-completionIcon": {
    fontSize: "90%",
    width: ".8em",
    display: "inline-block",
    textAlign: "center",
    paddingRight: ".6em",
    opacity: "0.6"
  },
  ".cm-completionIcon-function, .cm-completionIcon-method": {
    "&:after": {content: "'\u0192'"}
  },
  ".cm-completionIcon-class": {
    "&:after": {content: "'\u25CB'"}
  },
  ".cm-completionIcon-interface": {
    "&:after": {content: "'\u25CC'"}
  },
  ".cm-completionIcon-variable": {
    "&:after": {content: "'\u{1D465}'"}
  },
  ".cm-completionIcon-constant": {
    "&:after": {content: "'\u{1D436}'"}
  },
  ".cm-completionIcon-type": {
    "&:after": {content: "'\u{1D461}'"}
  },
  ".cm-completionIcon-enum": {
    "&:after": {content: "'\u222A'"}
  },
  ".cm-completionIcon-property": {
    "&:after": {content: "'\u25A1'"}
  },
  ".cm-completionIcon-keyword": {
    "&:after": {content: "'\u{1F511}\uFE0E'"}
  },
  ".cm-completionIcon-namespace": {
    "&:after": {content: "'\u25A2'"}
  },
  ".cm-completionIcon-text": {
    "&:after": {content: "'abc'", fontSize: "50%", verticalAlign: "middle"}
  }
});
function createListBox(options, id2, range) {
  const ul = document.createElement("ul");
  ul.id = id2;
  ul.setAttribute("role", "listbox");
  ul.setAttribute("aria-expanded", "true");
  for (let i = range.from; i < range.to; i++) {
    let {completion, match} = options[i];
    const li = ul.appendChild(document.createElement("li"));
    li.id = id2 + "-" + i;
    let icon = li.appendChild(document.createElement("div"));
    icon.classList.add("cm-completionIcon");
    if (completion.type)
      icon.classList.add("cm-completionIcon-" + completion.type);
    icon.setAttribute("aria-hidden", "true");
    let labelElt = li.appendChild(document.createElement("span"));
    labelElt.className = "cm-completionLabel";
    let {label, detail} = completion, off = 0;
    for (let j = 1; j < match.length; ) {
      let from = match[j++], to = match[j++];
      if (from > off)
        labelElt.appendChild(document.createTextNode(label.slice(off, from)));
      let span = labelElt.appendChild(document.createElement("span"));
      span.appendChild(document.createTextNode(label.slice(from, to)));
      span.className = "cm-completionMatchedText";
      off = to;
    }
    if (off < label.length)
      labelElt.appendChild(document.createTextNode(label.slice(off)));
    if (detail) {
      let detailElt = li.appendChild(document.createElement("span"));
      detailElt.className = "cm-completionDetail";
      detailElt.textContent = detail;
    }
    li.setAttribute("role", "option");
  }
  if (range.from)
    ul.classList.add("cm-completionListIncompleteTop");
  if (range.to < options.length)
    ul.classList.add("cm-completionListIncompleteBottom");
  return ul;
}
function createInfoDialog(option, view) {
  let dom = document.createElement("div");
  dom.className = "cm-tooltip cm-completionInfo";
  let {info} = option.completion;
  if (typeof info == "string") {
    dom.textContent = info;
  } else {
    let content2 = info(option.completion);
    if (content2.then)
      content2.then((node) => dom.appendChild(node), (e) => logException(view.state, e, "completion info"));
    else
      dom.appendChild(content2);
  }
  return dom;
}
function rangeAroundSelected(total, selected, max) {
  if (total <= max)
    return {from: 0, to: total};
  if (selected <= total >> 1) {
    let off2 = Math.floor(selected / max);
    return {from: off2 * max, to: (off2 + 1) * max};
  }
  let off = Math.floor((total - selected) / max);
  return {from: total - (off + 1) * max, to: total - off * max};
}
var CompletionTooltip = class {
  constructor(view, stateField) {
    this.view = view;
    this.stateField = stateField;
    this.info = null;
    this.placeInfo = {
      read: () => this.measureInfo(),
      write: (pos) => this.positionInfo(pos),
      key: this
    };
    let cState = view.state.field(stateField);
    let {options, selected} = cState.open;
    let config2 = view.state.facet(completionConfig);
    this.range = rangeAroundSelected(options.length, selected, config2.maxRenderedOptions);
    this.dom = document.createElement("div");
    this.dom.className = "cm-tooltip-autocomplete";
    this.dom.addEventListener("mousedown", (e) => {
      for (let dom = e.target, match; dom && dom != this.dom; dom = dom.parentNode) {
        if (dom.nodeName == "LI" && (match = /-(\d+)$/.exec(dom.id)) && +match[1] < options.length) {
          applyCompletion(view, options[+match[1]]);
          e.preventDefault();
          return;
        }
      }
    });
    this.list = this.dom.appendChild(createListBox(options, cState.id, this.range));
    this.list.addEventListener("scroll", () => {
      if (this.info)
        this.view.requestMeasure(this.placeInfo);
    });
  }
  mount() {
    this.updateSel();
  }
  update(update) {
    if (update.state.field(this.stateField) != update.startState.field(this.stateField))
      this.updateSel();
  }
  positioned() {
    if (this.info)
      this.view.requestMeasure(this.placeInfo);
  }
  updateSel() {
    let cState = this.view.state.field(this.stateField), open = cState.open;
    if (open.selected < this.range.from || open.selected >= this.range.to) {
      this.range = rangeAroundSelected(open.options.length, open.selected, this.view.state.facet(completionConfig).maxRenderedOptions);
      this.list.remove();
      this.list = this.dom.appendChild(createListBox(open.options, cState.id, this.range));
      this.list.addEventListener("scroll", () => {
        if (this.info)
          this.view.requestMeasure(this.placeInfo);
      });
    }
    if (this.updateSelectedOption(open.selected)) {
      if (this.info) {
        this.info.remove();
        this.info = null;
      }
      let option = open.options[open.selected];
      if (option.completion.info) {
        this.info = this.dom.appendChild(createInfoDialog(option, this.view));
        this.view.requestMeasure(this.placeInfo);
      }
    }
  }
  updateSelectedOption(selected) {
    let set = null;
    for (let opt = this.list.firstChild, i = this.range.from; opt; opt = opt.nextSibling, i++) {
      if (i == selected) {
        if (!opt.hasAttribute("aria-selected")) {
          opt.setAttribute("aria-selected", "true");
          set = opt;
        }
      } else {
        if (opt.hasAttribute("aria-selected"))
          opt.removeAttribute("aria-selected");
      }
    }
    if (set)
      scrollIntoView(this.list, set);
    return set;
  }
  measureInfo() {
    let sel = this.dom.querySelector("[aria-selected]");
    if (!sel)
      return null;
    let rect = this.dom.getBoundingClientRect();
    let top2 = sel.getBoundingClientRect().top - rect.top;
    if (top2 < 0 || top2 > this.list.clientHeight - 10)
      return null;
    let left = this.view.textDirection == Direction.RTL;
    let spaceLeft = rect.left, spaceRight = innerWidth - rect.right;
    if (left && spaceLeft < Math.min(MaxInfoWidth, spaceRight))
      left = false;
    else if (!left && spaceRight < Math.min(MaxInfoWidth, spaceLeft))
      left = true;
    return {top: top2, left};
  }
  positionInfo(pos) {
    if (this.info && pos) {
      this.info.style.top = pos.top + "px";
      this.info.classList.toggle("cm-completionInfo-left", pos.left);
      this.info.classList.toggle("cm-completionInfo-right", !pos.left);
    }
  }
};
function completionTooltip(stateField) {
  return (view) => new CompletionTooltip(view, stateField);
}
function scrollIntoView(container, element) {
  let parent = container.getBoundingClientRect();
  let self = element.getBoundingClientRect();
  if (self.top < parent.top)
    container.scrollTop -= parent.top - self.top;
  else if (self.bottom > parent.bottom)
    container.scrollTop += self.bottom - parent.bottom;
}
var MaxOptions = 300;
function score(option) {
  return (option.boost || 0) * 100 + (option.apply ? 10 : 0) + (option.info ? 5 : 0) + (option.type ? 1 : 0);
}
function sortOptions(active, state) {
  let options = [];
  for (let a of active)
    if (a.hasResult()) {
      let matcher = new FuzzyMatcher(state.sliceDoc(a.from, a.to)), match;
      for (let option of a.result.options)
        if (match = matcher.match(option.label)) {
          if (option.boost != null)
            match[0] += option.boost;
          options.push(new Option(option, a, match));
        }
    }
  options.sort(cmpOption);
  let result = [], prev = null;
  for (let opt of options.sort(cmpOption)) {
    if (result.length == MaxOptions)
      break;
    if (!prev || prev.label != opt.completion.label || prev.detail != opt.completion.detail)
      result.push(opt);
    else if (score(opt.completion) > score(prev))
      result[result.length - 1] = opt;
    prev = opt.completion;
  }
  return result;
}
var CompletionDialog = class {
  constructor(options, attrs, tooltip, timestamp, selected) {
    this.options = options;
    this.attrs = attrs;
    this.tooltip = tooltip;
    this.timestamp = timestamp;
    this.selected = selected;
  }
  setSelected(selected, id2) {
    return selected == this.selected || selected >= this.options.length ? this : new CompletionDialog(this.options, makeAttrs(id2, selected), this.tooltip, this.timestamp, selected);
  }
  static build(active, state, id2, prev) {
    let options = sortOptions(active, state);
    if (!options.length)
      return null;
    let selected = 0;
    if (prev && prev.selected) {
      let selectedValue = prev.options[prev.selected].completion;
      for (let i = 0; i < options.length && !selected; i++) {
        if (options[i].completion == selectedValue)
          selected = i;
      }
    }
    return new CompletionDialog(options, makeAttrs(id2, selected), {
      pos: active.reduce((a, b) => b.hasResult() ? Math.min(a, b.from) : a, 1e8),
      create: completionTooltip(completionState)
    }, prev ? prev.timestamp : Date.now(), selected);
  }
  map(changes) {
    return new CompletionDialog(this.options, this.attrs, Object.assign(Object.assign({}, this.tooltip), {pos: changes.mapPos(this.tooltip.pos)}), this.timestamp, this.selected);
  }
};
var CompletionState = class {
  constructor(active, id2, open) {
    this.active = active;
    this.id = id2;
    this.open = open;
  }
  static start() {
    return new CompletionState(none3, "cm-ac-" + Math.floor(Math.random() * 2e6).toString(36), null);
  }
  update(tr) {
    let {state} = tr, conf = state.facet(completionConfig);
    let sources = conf.override || state.languageDataAt("autocomplete", cur(state)).map(asSource);
    let active = sources.map((source) => {
      let value = this.active.find((s) => s.source == source) || new ActiveSource(source, 0, false);
      return value.update(tr, conf);
    });
    if (active.length == this.active.length && active.every((a, i) => a == this.active[i]))
      active = this.active;
    let open = tr.selection || active.some((a) => a.hasResult() && tr.changes.touchesRange(a.from, a.to)) || !sameResults(active, this.active) ? CompletionDialog.build(active, state, this.id, this.open) : this.open && tr.docChanged ? this.open.map(tr.changes) : this.open;
    if (!open && active.every((a) => a.state != 1) && active.some((a) => a.hasResult()))
      active = active.map((a) => a.hasResult() ? new ActiveSource(a.source, 0, false) : a);
    for (let effect of tr.effects)
      if (effect.is(setSelectedEffect))
        open = open && open.setSelected(effect.value, this.id);
    return active == this.active && open == this.open ? this : new CompletionState(active, this.id, open);
  }
  get tooltip() {
    return this.open ? this.open.tooltip : null;
  }
  get attrs() {
    return this.open ? this.open.attrs : baseAttrs;
  }
};
function sameResults(a, b) {
  if (a == b)
    return true;
  for (let iA = 0, iB = 0; ; ) {
    while (iA < a.length && !a[iA].hasResult)
      iA++;
    while (iB < b.length && !b[iB].hasResult)
      iB++;
    let endA = iA == a.length, endB = iB == b.length;
    if (endA || endB)
      return endA == endB;
    if (a[iA++].result != b[iB++].result)
      return false;
  }
}
function makeAttrs(id2, selected) {
  return {
    "aria-autocomplete": "list",
    "aria-activedescendant": id2 + "-" + selected,
    "aria-owns": id2
  };
}
var baseAttrs = {"aria-autocomplete": "list"};
var none3 = [];
function cmpOption(a, b) {
  let dScore = b.match[0] - a.match[0];
  if (dScore)
    return dScore;
  let lA = a.completion.label, lB = b.completion.label;
  return lA < lB ? -1 : lA == lB ? 0 : 1;
}
var ActiveSource = class {
  constructor(source, state, explicit) {
    this.source = source;
    this.state = state;
    this.explicit = explicit;
  }
  hasResult() {
    return false;
  }
  update(tr, conf) {
    let event = tr.annotation(Transaction.userEvent), value = this;
    if (event == "input" || event == "delete")
      value = value.handleUserEvent(tr, event, conf);
    else if (tr.docChanged)
      value = value.handleChange(tr);
    else if (tr.selection && value.state != 0)
      value = new ActiveSource(value.source, 0, false);
    for (let effect of tr.effects) {
      if (effect.is(startCompletionEffect))
        value = new ActiveSource(value.source, 1, effect.value);
      else if (effect.is(closeCompletionEffect))
        value = new ActiveSource(value.source, 0, false);
      else if (effect.is(setActiveEffect)) {
        for (let active of effect.value)
          if (active.source == value.source)
            value = active;
      }
    }
    return value;
  }
  handleUserEvent(_tr, type2, conf) {
    return type2 == "delete" || !conf.activateOnTyping ? this : new ActiveSource(this.source, 1, false);
  }
  handleChange(tr) {
    return tr.changes.touchesRange(cur(tr.startState)) ? new ActiveSource(this.source, 0, false) : this;
  }
};
var ActiveResult = class extends ActiveSource {
  constructor(source, explicit, result, from, to, span) {
    super(source, 2, explicit);
    this.result = result;
    this.from = from;
    this.to = to;
    this.span = span;
  }
  hasResult() {
    return true;
  }
  handleUserEvent(tr, type2, conf) {
    let from = tr.changes.mapPos(this.from), to = tr.changes.mapPos(this.to, 1);
    let pos = cur(tr.state);
    if ((this.explicit ? pos < from : pos <= from) || pos > to)
      return new ActiveSource(this.source, type2 == "input" && conf.activateOnTyping ? 1 : 0, false);
    if (this.span && (from == to || this.span.test(tr.state.sliceDoc(from, to))))
      return new ActiveResult(this.source, this.explicit, this.result, from, to, this.span);
    return new ActiveSource(this.source, 1, this.explicit);
  }
  handleChange(tr) {
    return tr.changes.touchesRange(this.from, this.to) ? new ActiveSource(this.source, 0, false) : new ActiveResult(this.source, this.explicit, this.result, tr.changes.mapPos(this.from), tr.changes.mapPos(this.to, 1), this.span);
  }
  map(mapping) {
    return new ActiveResult(this.source, this.explicit, this.result, mapping.mapPos(this.from), mapping.mapPos(this.to, 1), this.span);
  }
};
var startCompletionEffect = StateEffect.define();
var closeCompletionEffect = StateEffect.define();
var setActiveEffect = StateEffect.define({
  map(sources, mapping) {
    return sources.map((s) => s.hasResult() && !mapping.empty ? s.map(mapping) : s);
  }
});
var setSelectedEffect = StateEffect.define();
var completionState = StateField.define({
  create() {
    return CompletionState.start();
  },
  update(value, tr) {
    return value.update(tr);
  },
  provide: (f) => [
    showTooltip.from(f, (val) => val.tooltip),
    EditorView.contentAttributes.from(f, (state) => state.attrs)
  ]
});
var CompletionInteractMargin = 75;
function moveCompletionSelection(forward, by = "option") {
  return (view) => {
    let cState = view.state.field(completionState, false);
    if (!cState || !cState.open || Date.now() - cState.open.timestamp < CompletionInteractMargin)
      return false;
    let step = 1, tooltip;
    if (by == "page" && (tooltip = view.dom.querySelector(".cm-tooltip-autocomplete")))
      step = Math.max(2, Math.floor(tooltip.offsetHeight / tooltip.firstChild.offsetHeight));
    let selected = cState.open.selected + step * (forward ? 1 : -1), {length} = cState.open.options;
    if (selected < 0)
      selected = by == "page" ? 0 : length - 1;
    else if (selected >= length)
      selected = by == "page" ? length - 1 : 0;
    view.dispatch({effects: setSelectedEffect.of(selected)});
    return true;
  };
}
var acceptCompletion = (view) => {
  let cState = view.state.field(completionState, false);
  if (!cState || !cState.open || Date.now() - cState.open.timestamp < CompletionInteractMargin)
    return false;
  applyCompletion(view, cState.open.options[cState.open.selected]);
  return true;
};
var startCompletion = (view) => {
  let cState = view.state.field(completionState, false);
  if (!cState)
    return false;
  view.dispatch({effects: startCompletionEffect.of(true)});
  return true;
};
var closeCompletion = (view) => {
  let cState = view.state.field(completionState, false);
  if (!cState || !cState.active.some((a) => a.state != 0))
    return false;
  view.dispatch({effects: closeCompletionEffect.of(null)});
  return true;
};
var RunningQuery = class {
  constructor(source, context) {
    this.source = source;
    this.context = context;
    this.time = Date.now();
    this.updates = [];
    this.done = void 0;
  }
};
var DebounceTime = 50;
var MaxUpdateCount = 50;
var MinAbortTime = 1e3;
var completionPlugin = ViewPlugin.fromClass(class {
  constructor(view) {
    this.view = view;
    this.debounceUpdate = -1;
    this.running = [];
    this.debounceAccept = -1;
    this.composing = 0;
    for (let active of view.state.field(completionState).active)
      if (active.state == 1)
        this.startQuery(active);
  }
  update(update) {
    let cState = update.state.field(completionState);
    if (!update.selectionSet && !update.docChanged && update.startState.field(completionState) == cState)
      return;
    let doesReset = update.transactions.some((tr) => {
      let event = tr.annotation(Transaction.userEvent);
      return (tr.selection || tr.docChanged) && event != "input" && event != "delete";
    });
    for (let i = 0; i < this.running.length; i++) {
      let query = this.running[i];
      if (doesReset || query.updates.length + update.transactions.length > MaxUpdateCount && query.time - Date.now() > MinAbortTime) {
        for (let handler of query.context.abortListeners) {
          try {
            handler();
          } catch (e) {
            logException(this.view.state, e);
          }
        }
        query.context.abortListeners = null;
        this.running.splice(i--, 1);
      } else {
        query.updates.push(...update.transactions);
      }
    }
    if (this.debounceUpdate > -1)
      clearTimeout(this.debounceUpdate);
    this.debounceUpdate = cState.active.some((a) => a.state == 1 && !this.running.some((q) => q.source == a.source)) ? setTimeout(() => this.startUpdate(), DebounceTime) : -1;
    if (this.composing != 0)
      for (let tr of update.transactions) {
        if (tr.annotation(Transaction.userEvent) == "input")
          this.composing = 2;
        else if (this.composing == 2 && tr.selection)
          this.composing = 3;
      }
  }
  startUpdate() {
    this.debounceUpdate = -1;
    let {state} = this.view, cState = state.field(completionState);
    for (let active of cState.active) {
      if (active.state == 1 && !this.running.some((r) => r.source == active.source))
        this.startQuery(active);
    }
  }
  startQuery(active) {
    let {state} = this.view, pos = cur(state);
    let context = new CompletionContext(state, pos, active.explicit);
    let pending = new RunningQuery(active.source, context);
    this.running.push(pending);
    Promise.resolve(active.source(context)).then((result) => {
      if (!pending.context.aborted) {
        pending.done = result || null;
        this.scheduleAccept();
      }
    }, (err) => {
      this.view.dispatch({effects: closeCompletionEffect.of(null)});
      logException(this.view.state, err);
    });
  }
  scheduleAccept() {
    if (this.running.every((q) => q.done !== void 0))
      this.accept();
    else if (this.debounceAccept < 0)
      this.debounceAccept = setTimeout(() => this.accept(), DebounceTime);
  }
  accept() {
    var _a;
    if (this.debounceAccept > -1)
      clearTimeout(this.debounceAccept);
    this.debounceAccept = -1;
    let updated = [];
    let conf = this.view.state.facet(completionConfig);
    for (let i = 0; i < this.running.length; i++) {
      let query = this.running[i];
      if (query.done === void 0)
        continue;
      this.running.splice(i--, 1);
      if (query.done) {
        let active = new ActiveResult(query.source, query.context.explicit, query.done, query.done.from, (_a = query.done.to) !== null && _a !== void 0 ? _a : cur(query.updates.length ? query.updates[0].startState : this.view.state), query.done.span ? ensureAnchor(query.done.span, true) : null);
        for (let tr of query.updates)
          active = active.update(tr, conf);
        if (active.hasResult()) {
          updated.push(active);
          continue;
        }
      }
      let current = this.view.state.field(completionState).active.find((a) => a.source == query.source);
      if (current && current.state == 1) {
        if (query.done == null) {
          let active = new ActiveSource(query.source, 0, false);
          for (let tr of query.updates)
            active = active.update(tr, conf);
          if (active.state != 1)
            updated.push(active);
        } else {
          this.startQuery(current);
        }
      }
    }
    if (updated.length)
      this.view.dispatch({effects: setActiveEffect.of(updated)});
  }
}, {
  eventHandlers: {
    compositionstart() {
      this.composing = 1;
    },
    compositionend() {
      if (this.composing == 3)
        this.view.dispatch({effects: startCompletionEffect.of(false)});
      this.composing = 0;
    }
  }
});
var FieldPos = class {
  constructor(field, line, from, to) {
    this.field = field;
    this.line = line;
    this.from = from;
    this.to = to;
  }
};
var FieldRange = class {
  constructor(field, from, to) {
    this.field = field;
    this.from = from;
    this.to = to;
  }
  map(changes) {
    return new FieldRange(this.field, changes.mapPos(this.from, -1), changes.mapPos(this.to, 1));
  }
};
var Snippet = class {
  constructor(lines, fieldPositions) {
    this.lines = lines;
    this.fieldPositions = fieldPositions;
  }
  instantiate(state, pos) {
    let text = [], lineStart = [pos];
    let lineObj = state.doc.lineAt(pos), baseIndent = /^\s*/.exec(lineObj.text)[0];
    for (let line of this.lines) {
      if (text.length) {
        let indent2 = baseIndent, tabs = /^\t*/.exec(line)[0].length;
        for (let i = 0; i < tabs; i++)
          indent2 += state.facet(indentUnit);
        lineStart.push(pos + indent2.length - tabs);
        line = indent2 + line.slice(tabs);
      }
      text.push(line);
      pos += line.length + 1;
    }
    let ranges = this.fieldPositions.map((pos2) => new FieldRange(pos2.field, lineStart[pos2.line] + pos2.from, lineStart[pos2.line] + pos2.to));
    return {text, ranges};
  }
  static parse(template2) {
    let fields = [];
    let lines = [], positions = [], m;
    for (let line of template2.split(/\r\n?|\n/)) {
      while (m = /[#$]\{(?:(\d+)(?::([^}]*))?|([^}]*))\}/.exec(line)) {
        let seq = m[1] ? +m[1] : null, name2 = m[2] || m[3], found = -1;
        for (let i = 0; i < fields.length; i++) {
          if (name2 ? fields[i].name == name2 : seq != null && fields[i].seq == seq)
            found = i;
        }
        if (found < 0) {
          let i = 0;
          while (i < fields.length && (seq == null || fields[i].seq != null && fields[i].seq < seq))
            i++;
          fields.splice(i, 0, {seq, name: name2 || null});
          found = i;
        }
        positions.push(new FieldPos(found, lines.length, m.index, m.index + name2.length));
        line = line.slice(0, m.index) + name2 + line.slice(m.index + m[0].length);
      }
      lines.push(line);
    }
    return new Snippet(lines, positions);
  }
};
var fieldMarker = Decoration.widget({widget: new class extends WidgetType {
  toDOM() {
    let span = document.createElement("span");
    span.className = "cm-snippetFieldPosition";
    return span;
  }
  ignoreEvent() {
    return false;
  }
}()});
var fieldRange = Decoration.mark({class: "cm-snippetField"});
var ActiveSnippet = class {
  constructor(ranges, active) {
    this.ranges = ranges;
    this.active = active;
    this.deco = Decoration.set(ranges.map((r) => (r.from == r.to ? fieldMarker : fieldRange).range(r.from, r.to)));
  }
  map(changes) {
    return new ActiveSnippet(this.ranges.map((r) => r.map(changes)), this.active);
  }
  selectionInsideField(sel) {
    return sel.ranges.every((range) => this.ranges.some((r) => r.field == this.active && r.from <= range.from && r.to >= range.to));
  }
};
var setActive = StateEffect.define({
  map(value, changes) {
    return value && value.map(changes);
  }
});
var moveToField = StateEffect.define();
var snippetState = StateField.define({
  create() {
    return null;
  },
  update(value, tr) {
    for (let effect of tr.effects) {
      if (effect.is(setActive))
        return effect.value;
      if (effect.is(moveToField) && value)
        return new ActiveSnippet(value.ranges, effect.value);
    }
    if (value && tr.docChanged)
      value = value.map(tr.changes);
    if (value && tr.selection && !value.selectionInsideField(tr.selection))
      value = null;
    return value;
  },
  provide: (f) => EditorView.decorations.from(f, (val) => val ? val.deco : Decoration.none)
});
function fieldSelection(ranges, field) {
  return EditorSelection.create(ranges.filter((r) => r.field == field).map((r) => EditorSelection.range(r.from, r.to)));
}
function snippet(template2) {
  let snippet2 = Snippet.parse(template2);
  return (editor, _completion, from, to) => {
    let {text, ranges} = snippet2.instantiate(editor.state, from);
    let spec = {changes: {from, to, insert: Text.of(text)}};
    if (ranges.length)
      spec.selection = fieldSelection(ranges, 0);
    if (ranges.length > 1) {
      let effects = spec.effects = [setActive.of(new ActiveSnippet(ranges, 0))];
      if (editor.state.field(snippetState, false) === void 0)
        effects.push(StateEffect.appendConfig.of([snippetState, addSnippetKeymap, snippetPointerHandler, baseTheme5]));
    }
    editor.dispatch(editor.state.update(spec));
  };
}
function moveField(dir) {
  return ({state, dispatch}) => {
    let active = state.field(snippetState, false);
    if (!active || dir < 0 && active.active == 0)
      return false;
    let next = active.active + dir, last = dir > 0 && !active.ranges.some((r) => r.field == next + dir);
    dispatch(state.update({
      selection: fieldSelection(active.ranges, next),
      effects: setActive.of(last ? null : new ActiveSnippet(active.ranges, next))
    }));
    return true;
  };
}
var clearSnippet = ({state, dispatch}) => {
  let active = state.field(snippetState, false);
  if (!active)
    return false;
  dispatch(state.update({effects: setActive.of(null)}));
  return true;
};
var nextSnippetField = moveField(1);
var prevSnippetField = moveField(-1);
var defaultSnippetKeymap = [
  {key: "Tab", run: nextSnippetField, shift: prevSnippetField},
  {key: "Escape", run: clearSnippet}
];
var snippetKeymap = Facet.define({
  combine(maps) {
    return maps.length ? maps[0] : defaultSnippetKeymap;
  }
});
var addSnippetKeymap = Prec.override(keymap.compute([snippetKeymap], (state) => state.facet(snippetKeymap)));
function snippetCompletion(template2, completion) {
  return Object.assign(Object.assign({}, completion), {apply: snippet(template2)});
}
var snippetPointerHandler = EditorView.domEventHandlers({
  mousedown(event, view) {
    let active = view.state.field(snippetState, false), pos;
    if (!active || (pos = view.posAtCoords({x: event.clientX, y: event.clientY})) == null)
      return false;
    let match = active.ranges.find((r) => r.from <= pos && r.to >= pos);
    if (!match || match.field == active.active)
      return false;
    view.dispatch({
      selection: fieldSelection(active.ranges, match.field),
      effects: setActive.of(active.ranges.some((r) => r.field > match.field) ? new ActiveSnippet(active.ranges, match.field) : null)
    });
    return true;
  }
});
var completionKeymap = [
  {key: "Ctrl-Space", run: startCompletion},
  {key: "Escape", run: closeCompletion},
  {key: "ArrowDown", run: moveCompletionSelection(true)},
  {key: "ArrowUp", run: moveCompletionSelection(false)},
  {key: "PageDown", run: moveCompletionSelection(true, "page")},
  {key: "PageUp", run: moveCompletionSelection(false, "page")},
  {key: "Enter", run: acceptCompletion}
];
var completionKeymapExt = Prec.override(keymap.computeN([completionConfig], (state) => state.facet(completionConfig).defaultKeymap ? [completionKeymap] : []));

// node_modules/@codemirror/lang-javascript/dist/index.js
var snippets = [
  snippetCompletion("function ${name}(${params}) {\n	${}\n}", {
    label: "function",
    detail: "definition",
    type: "keyword"
  }),
  snippetCompletion("for (let ${index} = 0; ${index} < ${bound}; ${index}++) {\n	${}\n}", {
    label: "for",
    detail: "loop",
    type: "keyword"
  }),
  snippetCompletion("for (let ${name} of ${collection}) {\n	${}\n}", {
    label: "for",
    detail: "of loop",
    type: "keyword"
  }),
  snippetCompletion("try {\n	${}\n} catch (${error}) {\n	${}\n}", {
    label: "try",
    detail: "block",
    type: "keyword"
  }),
  snippetCompletion("class ${name} {\n	constructor(${params}) {\n		${}\n	}\n}", {
    label: "class",
    detail: "definition",
    type: "keyword"
  }),
  snippetCompletion('import {${names}} from "${module}"\n${}', {
    label: "import",
    detail: "named",
    type: "keyword"
  }),
  snippetCompletion('import ${name} from "${module}"\n${}', {
    label: "import",
    detail: "default",
    type: "keyword"
  })
];
var javascriptLanguage = LezerLanguage.define({
  parser: parser2.configure({
    props: [
      indentNodeProp.add({
        IfStatement: continuedIndent({except: /^\s*({|else\b)/}),
        TryStatement: continuedIndent({except: /^\s*({|catch|finally)\b/}),
        LabeledStatement: flatIndent,
        SwitchBody: (context) => {
          let after = context.textAfter, closed = /^\s*\}/.test(after), isCase = /^\s*(case|default)\b/.test(after);
          return context.baseIndent + (closed ? 0 : isCase ? 1 : 2) * context.unit;
        },
        Block: delimitedIndent({closing: "}"}),
        ArrowFunction: (cx) => cx.baseIndent + cx.unit,
        "TemplateString BlockComment": () => -1,
        "Statement Property": continuedIndent({except: /^{/}),
        JSXElement(context) {
          let closed = /^\s*<\//.test(context.textAfter);
          return context.lineIndent(context.state.doc.lineAt(context.node.from)) + (closed ? 0 : context.unit);
        },
        JSXEscape(context) {
          let closed = /\s*\}/.test(context.textAfter);
          return context.lineIndent(context.state.doc.lineAt(context.node.from)) + (closed ? 0 : context.unit);
        },
        "JSXOpenTag JSXSelfClosingTag"(context) {
          return context.column(context.node.from) + context.unit;
        }
      }),
      foldNodeProp.add({
        "Block ClassBody SwitchBody EnumBody ObjectExpression ArrayExpression": foldInside,
        BlockComment(tree) {
          return {from: tree.from + 2, to: tree.to - 2};
        }
      }),
      styleTags({
        "get set async static": tags.modifier,
        "for while do if else switch try catch finally return throw break continue default case": tags.controlKeyword,
        "in of await yield void typeof delete instanceof": tags.operatorKeyword,
        "export import let var const function class extends": tags.definitionKeyword,
        "with debugger from as new": tags.keyword,
        TemplateString: tags.special(tags.string),
        Super: tags.atom,
        BooleanLiteral: tags.bool,
        this: tags.self,
        null: tags.null,
        Star: tags.modifier,
        VariableName: tags.variableName,
        "CallExpression/VariableName": tags.function(tags.variableName),
        VariableDefinition: tags.definition(tags.variableName),
        Label: tags.labelName,
        PropertyName: tags.propertyName,
        "CallExpression/MemberExpression/PropertyName": tags.function(tags.propertyName),
        "FunctionDeclaration/VariableDefinition": tags.function(tags.definition(tags.variableName)),
        "ClassDeclaration/VariableDefinition": tags.definition(tags.className),
        PropertyNameDefinition: tags.definition(tags.propertyName),
        UpdateOp: tags.updateOperator,
        LineComment: tags.lineComment,
        BlockComment: tags.blockComment,
        Number: tags.number,
        String: tags.string,
        ArithOp: tags.arithmeticOperator,
        LogicOp: tags.logicOperator,
        BitOp: tags.bitwiseOperator,
        CompareOp: tags.compareOperator,
        RegExp: tags.regexp,
        Equals: tags.definitionOperator,
        "Arrow : Spread": tags.punctuation,
        "( )": tags.paren,
        "[ ]": tags.squareBracket,
        "{ }": tags.brace,
        ".": tags.derefOperator,
        ", ;": tags.separator,
        TypeName: tags.typeName,
        TypeDefinition: tags.definition(tags.typeName),
        "type enum interface implements namespace module declare": tags.definitionKeyword,
        "abstract global privacy readonly": tags.modifier,
        "is keyof unique infer": tags.operatorKeyword,
        JSXAttributeValue: tags.string,
        JSXText: tags.content,
        "JSXStartTag JSXStartCloseTag JSXSelfCloseEndTag JSXEndTag": tags.angleBracket,
        "JSXIdentifier JSXNameSpacedName": tags.tagName,
        "JSXAttribute/JSXIdentifier JSXAttribute/JSXNameSpacedName": tags.propertyName
      })
    ]
  }),
  languageData: {
    closeBrackets: {brackets: ["(", "[", "{", "'", '"', "`"]},
    commentTokens: {line: "//", block: {open: "/*", close: "*/"}},
    indentOnInput: /^\s*(?:case |default:|\{|\}|<\/)$/,
    wordChars: "$"
  }
});
var typescriptLanguage = javascriptLanguage.configure({dialect: "ts"});
var jsxLanguage = javascriptLanguage.configure({dialect: "jsx"});
var tsxLanguage = javascriptLanguage.configure({dialect: "jsx ts"});
function javascript(config2 = {}) {
  let lang = config2.jsx ? config2.typescript ? tsxLanguage : jsxLanguage : config2.typescript ? typescriptLanguage : javascriptLanguage;
  return new LanguageSupport(lang, javascriptLanguage.data.of({
    autocomplete: ifNotIn(["LineComment", "BlockComment", "String"], completeFromList(snippets))
  }));
}

// node_modules/@codemirror/legacy-modes/mode/julia.js
function wordRegexp(words4, end) {
  if (typeof end === "undefined") {
    end = "\\b";
  }
  return new RegExp("^((" + words4.join(")|(") + "))" + end);
}
var octChar = "\\\\[0-7]{1,3}";
var hexChar = "\\\\x[A-Fa-f0-9]{1,2}";
var sChar = `\\\\[abefnrtv0%?'"\\\\]`;
var uChar = "([^\\u0027\\u005C\\uD800-\\uDFFF]|[\\uD800-\\uDFFF][\\uDC00-\\uDFFF])";
var operators2 = wordRegexp([
  "[<>]:",
  "[<>=]=",
  "<<=?",
  ">>>?=?",
  "=>",
  "->",
  "\\/\\/",
  "[\\\\%*+\\-<>!=\\/^|&\\u00F7\\u22BB]=?",
  "\\?",
  "\\$",
  "~",
  ":",
  "\\u00D7",
  "\\u2208",
  "\\u2209",
  "\\u220B",
  "\\u220C",
  "\\u2218",
  "\\u221A",
  "\\u221B",
  "\\u2229",
  "\\u222A",
  "\\u2260",
  "\\u2264",
  "\\u2265",
  "\\u2286",
  "\\u2288",
  "\\u228A",
  "\\u22C5",
  "\\b(in|isa)\\b(?!.?\\()"
], "");
var delimiters = /^[;,()[\]{}]/;
var identifiers = /^[_A-Za-z\u00A1-\u2217\u2219-\uFFFF][\w\u00A1-\u2217\u2219-\uFFFF]*!*/;
var chars = wordRegexp([octChar, hexChar, sChar, uChar], "'");
var openersList = [
  "begin",
  "function",
  "type",
  "struct",
  "immutable",
  "let",
  "macro",
  "for",
  "while",
  "quote",
  "if",
  "else",
  "elseif",
  "try",
  "finally",
  "catch",
  "do"
];
var closersList = ["end", "else", "elseif", "catch", "finally"];
var keywordsList = [
  "if",
  "else",
  "elseif",
  "while",
  "for",
  "begin",
  "let",
  "end",
  "do",
  "try",
  "catch",
  "finally",
  "return",
  "break",
  "continue",
  "global",
  "local",
  "const",
  "export",
  "import",
  "importall",
  "using",
  "function",
  "where",
  "macro",
  "module",
  "baremodule",
  "struct",
  "type",
  "mutable",
  "immutable",
  "quote",
  "typealias",
  "abstract",
  "primitive",
  "bitstype"
];
var builtinsList = ["true", "false", "nothing", "NaN", "Inf"];
var openers = wordRegexp(openersList);
var closers = wordRegexp(closersList);
var keywords5 = wordRegexp(keywordsList);
var builtins3 = wordRegexp(builtinsList);
var macro = /^@[_A-Za-z][\w]*/;
var symbol2 = /^:[_A-Za-z\u00A1-\uFFFF][\w\u00A1-\uFFFF]*!*/;
var stringPrefixes = /^(`|([_A-Za-z\u00A1-\uFFFF]*"("")?))/;
function inArray(state) {
  return state.nestedArrays > 0;
}
function inGenerator(state) {
  return state.nestedGenerators > 0;
}
function currentScope(state, n) {
  if (typeof n === "undefined") {
    n = 0;
  }
  if (state.scopes.length <= n) {
    return null;
  }
  return state.scopes[state.scopes.length - (n + 1)];
}
function tokenBase4(stream, state) {
  if (stream.match("#=", false)) {
    state.tokenize = tokenComment2;
    return state.tokenize(stream, state);
  }
  var leavingExpr = state.leavingExpr;
  if (stream.sol()) {
    leavingExpr = false;
  }
  state.leavingExpr = false;
  if (leavingExpr) {
    if (stream.match(/^'+/)) {
      return "operator";
    }
  }
  if (stream.match(/\.{4,}/)) {
    return "error";
  } else if (stream.match(/\.{1,3}/)) {
    return "operator";
  }
  if (stream.eatSpace()) {
    return null;
  }
  var ch = stream.peek();
  if (ch === "#") {
    stream.skipToEnd();
    return "comment";
  }
  if (ch === "[") {
    state.scopes.push("[");
    state.nestedArrays++;
  }
  if (ch === "(") {
    state.scopes.push("(");
    state.nestedGenerators++;
  }
  if (inArray(state) && ch === "]") {
    while (state.scopes.length && currentScope(state) !== "[") {
      state.scopes.pop();
    }
    state.scopes.pop();
    state.nestedArrays--;
    state.leavingExpr = true;
  }
  if (inGenerator(state) && ch === ")") {
    while (state.scopes.length && currentScope(state) !== "(") {
      state.scopes.pop();
    }
    state.scopes.pop();
    state.nestedGenerators--;
    state.leavingExpr = true;
  }
  if (inArray(state)) {
    if (state.lastToken == "end" && stream.match(":")) {
      return "operator";
    }
    if (stream.match("end")) {
      return "number";
    }
  }
  var match;
  if (match = stream.match(openers, false)) {
    state.scopes.push(match[0]);
  }
  if (stream.match(closers, false)) {
    state.scopes.pop();
  }
  if (stream.match(/^::(?![:\$])/)) {
    state.tokenize = tokenAnnotation;
    return state.tokenize(stream, state);
  }
  if (!leavingExpr && stream.match(symbol2) || stream.match(/:([<>]:|<<=?|>>>?=?|->|\/\/|\.{2,3}|[\.\\%*+\-<>!\/^|&]=?|[~\?\$])/)) {
    return "builtin";
  }
  if (stream.match(operators2)) {
    return "operator";
  }
  if (stream.match(/^\.?\d/, false)) {
    var imMatcher = RegExp(/^im\b/);
    var numberLiteral = false;
    if (stream.match(/^0x\.[0-9a-f_]+p[\+\-]?[_\d]+/i)) {
      numberLiteral = true;
    }
    if (stream.match(/^0x[0-9a-f_]+/i)) {
      numberLiteral = true;
    }
    if (stream.match(/^0b[01_]+/i)) {
      numberLiteral = true;
    }
    if (stream.match(/^0o[0-7_]+/i)) {
      numberLiteral = true;
    }
    if (stream.match(/^(?:(?:\d[_\d]*)?\.(?!\.)(?:\d[_\d]*)?|\d[_\d]*\.(?!\.)(?:\d[_\d]*))?([Eef][\+\-]?[_\d]+)?/i)) {
      numberLiteral = true;
    }
    if (stream.match(/^\d[_\d]*(e[\+\-]?\d+)?/i)) {
      numberLiteral = true;
    }
    if (numberLiteral) {
      stream.match(imMatcher);
      state.leavingExpr = true;
      return "number";
    }
  }
  if (stream.match("'")) {
    state.tokenize = tokenChar;
    return state.tokenize(stream, state);
  }
  if (stream.match(stringPrefixes)) {
    state.tokenize = tokenStringFactory(stream.current());
    return state.tokenize(stream, state);
  }
  if (stream.match(macro)) {
    return "meta";
  }
  if (stream.match(delimiters)) {
    return null;
  }
  if (stream.match(keywords5)) {
    return "keyword";
  }
  if (stream.match(builtins3)) {
    return "builtin";
  }
  var isDefinition = state.isDefinition || state.lastToken == "function" || state.lastToken == "macro" || state.lastToken == "type" || state.lastToken == "struct" || state.lastToken == "immutable";
  if (stream.match(identifiers)) {
    if (isDefinition) {
      if (stream.peek() === ".") {
        state.isDefinition = true;
        return "variable";
      }
      state.isDefinition = false;
      return "def";
    }
    state.leavingExpr = true;
    return "variable";
  }
  stream.next();
  return "error";
}
function tokenAnnotation(stream, state) {
  stream.match(/.*?(?=[,;{}()=\s]|$)/);
  if (stream.match("{")) {
    state.nestedParameters++;
  } else if (stream.match("}") && state.nestedParameters > 0) {
    state.nestedParameters--;
  }
  if (state.nestedParameters > 0) {
    stream.match(/.*?(?={|})/) || stream.next();
  } else if (state.nestedParameters == 0) {
    state.tokenize = tokenBase4;
  }
  return "builtin";
}
function tokenComment2(stream, state) {
  if (stream.match("#=")) {
    state.nestedComments++;
  }
  if (!stream.match(/.*?(?=(#=|=#))/)) {
    stream.skipToEnd();
  }
  if (stream.match("=#")) {
    state.nestedComments--;
    if (state.nestedComments == 0)
      state.tokenize = tokenBase4;
  }
  return "comment";
}
function tokenChar(stream, state) {
  var isChar = false, match;
  if (stream.match(chars)) {
    isChar = true;
  } else if (match = stream.match(/\\u([a-f0-9]{1,4})(?=')/i)) {
    var value = parseInt(match[1], 16);
    if (value <= 55295 || value >= 57344) {
      isChar = true;
      stream.next();
    }
  } else if (match = stream.match(/\\U([A-Fa-f0-9]{5,8})(?=')/)) {
    var value = parseInt(match[1], 16);
    if (value <= 1114111) {
      isChar = true;
      stream.next();
    }
  }
  if (isChar) {
    state.leavingExpr = true;
    state.tokenize = tokenBase4;
    return "string";
  }
  if (!stream.match(/^[^']+(?=')/)) {
    stream.skipToEnd();
  }
  if (stream.match("'")) {
    state.tokenize = tokenBase4;
  }
  return "error";
}
function tokenStringFactory(delimiter) {
  if (delimiter.substr(-3) === '"""') {
    delimiter = '"""';
  } else if (delimiter.substr(-1) === '"') {
    delimiter = '"';
  }
  function tokenString5(stream, state) {
    if (stream.eat("\\")) {
      stream.next();
    } else if (stream.match(delimiter)) {
      state.tokenize = tokenBase4;
      state.leavingExpr = true;
      return "string";
    } else {
      stream.eat(/[`"]/);
    }
    stream.eatWhile(/[^\\`"]/);
    return "string";
  }
  return tokenString5;
}
var julia = {
  startState: function() {
    return {
      tokenize: tokenBase4,
      scopes: [],
      lastToken: null,
      leavingExpr: false,
      isDefinition: false,
      nestedArrays: 0,
      nestedComments: 0,
      nestedGenerators: 0,
      nestedParameters: 0,
      firstParenPos: -1
    };
  },
  token: function(stream, state) {
    var style = state.tokenize(stream, state);
    var current = stream.current();
    if (current && style) {
      state.lastToken = current;
    }
    return style;
  },
  indent: function(state, textAfter, cx) {
    var delta = 0;
    if (textAfter === "]" || textAfter === ")" || /^end\b/.test(textAfter) || /^else/.test(textAfter) || /^catch\b/.test(textAfter) || /^elseif\b/.test(textAfter) || /^finally/.test(textAfter)) {
      delta = -1;
    }
    return (state.scopes.length + delta) * cx.unit;
  },
  languageData: {
    indentOnInput: /^\s*(end|else|catch|finally)\b$/,
    commentTokens: {line: "#", block: {open: "#=", close: "=#"}},
    closeBrackets: {brackets: ["(", "[", "{", '"']},
    autocomplete: keywordsList.concat(builtinsList)
  }
};

// node_modules/@codemirror/legacy-modes/mode/lua.js
function prefixRE(words4) {
  return new RegExp("^(?:" + words4.join("|") + ")", "i");
}
function wordRE(words4) {
  return new RegExp("^(?:" + words4.join("|") + ")$", "i");
}
var builtins4 = wordRE([
  "_G",
  "_VERSION",
  "assert",
  "collectgarbage",
  "dofile",
  "error",
  "getfenv",
  "getmetatable",
  "ipairs",
  "load",
  "loadfile",
  "loadstring",
  "module",
  "next",
  "pairs",
  "pcall",
  "print",
  "rawequal",
  "rawget",
  "rawset",
  "require",
  "select",
  "setfenv",
  "setmetatable",
  "tonumber",
  "tostring",
  "type",
  "unpack",
  "xpcall",
  "coroutine.create",
  "coroutine.resume",
  "coroutine.running",
  "coroutine.status",
  "coroutine.wrap",
  "coroutine.yield",
  "debug.debug",
  "debug.getfenv",
  "debug.gethook",
  "debug.getinfo",
  "debug.getlocal",
  "debug.getmetatable",
  "debug.getregistry",
  "debug.getupvalue",
  "debug.setfenv",
  "debug.sethook",
  "debug.setlocal",
  "debug.setmetatable",
  "debug.setupvalue",
  "debug.traceback",
  "close",
  "flush",
  "lines",
  "read",
  "seek",
  "setvbuf",
  "write",
  "io.close",
  "io.flush",
  "io.input",
  "io.lines",
  "io.open",
  "io.output",
  "io.popen",
  "io.read",
  "io.stderr",
  "io.stdin",
  "io.stdout",
  "io.tmpfile",
  "io.type",
  "io.write",
  "math.abs",
  "math.acos",
  "math.asin",
  "math.atan",
  "math.atan2",
  "math.ceil",
  "math.cos",
  "math.cosh",
  "math.deg",
  "math.exp",
  "math.floor",
  "math.fmod",
  "math.frexp",
  "math.huge",
  "math.ldexp",
  "math.log",
  "math.log10",
  "math.max",
  "math.min",
  "math.modf",
  "math.pi",
  "math.pow",
  "math.rad",
  "math.random",
  "math.randomseed",
  "math.sin",
  "math.sinh",
  "math.sqrt",
  "math.tan",
  "math.tanh",
  "os.clock",
  "os.date",
  "os.difftime",
  "os.execute",
  "os.exit",
  "os.getenv",
  "os.remove",
  "os.rename",
  "os.setlocale",
  "os.time",
  "os.tmpname",
  "package.cpath",
  "package.loaded",
  "package.loaders",
  "package.loadlib",
  "package.path",
  "package.preload",
  "package.seeall",
  "string.byte",
  "string.char",
  "string.dump",
  "string.find",
  "string.format",
  "string.gmatch",
  "string.gsub",
  "string.len",
  "string.lower",
  "string.match",
  "string.rep",
  "string.reverse",
  "string.sub",
  "string.upper",
  "table.concat",
  "table.insert",
  "table.maxn",
  "table.remove",
  "table.sort"
]);
var keywords6 = wordRE([
  "and",
  "break",
  "elseif",
  "false",
  "nil",
  "not",
  "or",
  "return",
  "true",
  "function",
  "end",
  "if",
  "then",
  "else",
  "do",
  "while",
  "repeat",
  "until",
  "for",
  "in",
  "local"
]);
var indentTokens = wordRE(["function", "if", "repeat", "do", "\\(", "{"]);
var dedentTokens = wordRE(["end", "until", "\\)", "}"]);
var dedentPartial = prefixRE(["end", "until", "\\)", "}", "else", "elseif"]);
function readBracket(stream) {
  var level = 0;
  while (stream.eat("="))
    ++level;
  stream.eat("[");
  return level;
}
function normal2(stream, state) {
  var ch = stream.next();
  if (ch == "-" && stream.eat("-")) {
    if (stream.eat("[") && stream.eat("["))
      return (state.cur = bracketed(readBracket(stream), "comment"))(stream, state);
    stream.skipToEnd();
    return "comment";
  }
  if (ch == '"' || ch == "'")
    return (state.cur = string2(ch))(stream, state);
  if (ch == "[" && /[\[=]/.test(stream.peek()))
    return (state.cur = bracketed(readBracket(stream), "string"))(stream, state);
  if (/\d/.test(ch)) {
    stream.eatWhile(/[\w.%]/);
    return "number";
  }
  if (/[\w_]/.test(ch)) {
    stream.eatWhile(/[\w\\\-_.]/);
    return "variable";
  }
  return null;
}
function bracketed(level, style) {
  return function(stream, state) {
    var curlev = null, ch;
    while ((ch = stream.next()) != null) {
      if (curlev == null) {
        if (ch == "]")
          curlev = 0;
      } else if (ch == "=")
        ++curlev;
      else if (ch == "]" && curlev == level) {
        state.cur = normal2;
        break;
      } else
        curlev = null;
    }
    return style;
  };
}
function string2(quote) {
  return function(stream, state) {
    var escaped = false, ch;
    while ((ch = stream.next()) != null) {
      if (ch == quote && !escaped)
        break;
      escaped = !escaped && ch == "\\";
    }
    if (!escaped)
      state.cur = normal2;
    return "string";
  };
}
var lua = {
  startState: function(basecol) {
    return {basecol: basecol || 0, indentDepth: 0, cur: normal2};
  },
  token: function(stream, state) {
    if (stream.eatSpace())
      return null;
    var style = state.cur(stream, state);
    var word = stream.current();
    if (style == "variable") {
      if (keywords6.test(word))
        style = "keyword";
      else if (builtins4.test(word))
        style = "builtin";
    }
    if (style != "comment" && style != "string") {
      if (indentTokens.test(word))
        ++state.indentDepth;
      else if (dedentTokens.test(word))
        --state.indentDepth;
    }
    return style;
  },
  indent: function(state, textAfter, cx) {
    var closing3 = dedentPartial.test(textAfter);
    return state.basecol + cx.unit * (state.indentDepth - (closing3 ? 1 : 0));
  },
  languageData: {
    commentTokens: {line: "--", block: {open: "--[[", close: "]]--"}}
  }
};

// node_modules/@codemirror/legacy-modes/mode/perl.js
function look(stream, c2) {
  return stream.string.charAt(stream.pos + (c2 || 0));
}
function prefix(stream, c2) {
  if (c2) {
    var x = stream.pos - c2;
    return stream.string.substr(x >= 0 ? x : 0, c2);
  } else {
    return stream.string.substr(0, stream.pos - 1);
  }
}
function suffix(stream, c2) {
  var y = stream.string.length;
  var x = y - stream.pos + 1;
  return stream.string.substr(stream.pos, c2 && c2 < y ? c2 : x);
}
function eatSuffix(stream, c2) {
  var x = stream.pos + c2;
  var y;
  if (x <= 0)
    stream.pos = 0;
  else if (x >= (y = stream.string.length - 1))
    stream.pos = y;
  else
    stream.pos = x;
}
var PERL = {
  "->": 4,
  "++": 4,
  "--": 4,
  "**": 4,
  "=~": 4,
  "!~": 4,
  "*": 4,
  "/": 4,
  "%": 4,
  x: 4,
  "+": 4,
  "-": 4,
  ".": 4,
  "<<": 4,
  ">>": 4,
  "<": 4,
  ">": 4,
  "<=": 4,
  ">=": 4,
  lt: 4,
  gt: 4,
  le: 4,
  ge: 4,
  "==": 4,
  "!=": 4,
  "<=>": 4,
  eq: 4,
  ne: 4,
  cmp: 4,
  "~~": 4,
  "&": 4,
  "|": 4,
  "^": 4,
  "&&": 4,
  "||": 4,
  "//": 4,
  "..": 4,
  "...": 4,
  "?": 4,
  ":": 4,
  "=": 4,
  "+=": 4,
  "-=": 4,
  "*=": 4,
  ",": 4,
  "=>": 4,
  "::": 4,
  not: 4,
  and: 4,
  or: 4,
  xor: 4,
  BEGIN: [5, 1],
  END: [5, 1],
  PRINT: [5, 1],
  PRINTF: [5, 1],
  GETC: [5, 1],
  READ: [5, 1],
  READLINE: [5, 1],
  DESTROY: [5, 1],
  TIE: [5, 1],
  TIEHANDLE: [5, 1],
  UNTIE: [5, 1],
  STDIN: 5,
  STDIN_TOP: 5,
  STDOUT: 5,
  STDOUT_TOP: 5,
  STDERR: 5,
  STDERR_TOP: 5,
  $ARG: 5,
  $_: 5,
  "@ARG": 5,
  "@_": 5,
  $LIST_SEPARATOR: 5,
  '$"': 5,
  $PROCESS_ID: 5,
  $PID: 5,
  $$: 5,
  $REAL_GROUP_ID: 5,
  $GID: 5,
  "$(": 5,
  $EFFECTIVE_GROUP_ID: 5,
  $EGID: 5,
  "$)": 5,
  $PROGRAM_NAME: 5,
  $0: 5,
  $SUBSCRIPT_SEPARATOR: 5,
  $SUBSEP: 5,
  "$;": 5,
  $REAL_USER_ID: 5,
  $UID: 5,
  "$<": 5,
  $EFFECTIVE_USER_ID: 5,
  $EUID: 5,
  "$>": 5,
  $a: 5,
  $b: 5,
  $COMPILING: 5,
  "$^C": 5,
  $DEBUGGING: 5,
  "$^D": 5,
  "${^ENCODING}": 5,
  $ENV: 5,
  "%ENV": 5,
  $SYSTEM_FD_MAX: 5,
  "$^F": 5,
  "@F": 5,
  "${^GLOBAL_PHASE}": 5,
  "$^H": 5,
  "%^H": 5,
  "@INC": 5,
  "%INC": 5,
  $INPLACE_EDIT: 5,
  "$^I": 5,
  "$^M": 5,
  $OSNAME: 5,
  "$^O": 5,
  "${^OPEN}": 5,
  $PERLDB: 5,
  "$^P": 5,
  $SIG: 5,
  "%SIG": 5,
  $BASETIME: 5,
  "$^T": 5,
  "${^TAINT}": 5,
  "${^UNICODE}": 5,
  "${^UTF8CACHE}": 5,
  "${^UTF8LOCALE}": 5,
  $PERL_VERSION: 5,
  "$^V": 5,
  "${^WIN32_SLOPPY_STAT}": 5,
  $EXECUTABLE_NAME: 5,
  "$^X": 5,
  $1: 5,
  $MATCH: 5,
  "$&": 5,
  "${^MATCH}": 5,
  $PREMATCH: 5,
  "$`": 5,
  "${^PREMATCH}": 5,
  $POSTMATCH: 5,
  "$'": 5,
  "${^POSTMATCH}": 5,
  $LAST_PAREN_MATCH: 5,
  "$+": 5,
  $LAST_SUBMATCH_RESULT: 5,
  "$^N": 5,
  "@LAST_MATCH_END": 5,
  "@+": 5,
  "%LAST_PAREN_MATCH": 5,
  "%+": 5,
  "@LAST_MATCH_START": 5,
  "@-": 5,
  "%LAST_MATCH_START": 5,
  "%-": 5,
  $LAST_REGEXP_CODE_RESULT: 5,
  "$^R": 5,
  "${^RE_DEBUG_FLAGS}": 5,
  "${^RE_TRIE_MAXBUF}": 5,
  $ARGV: 5,
  "@ARGV": 5,
  ARGV: 5,
  ARGVOUT: 5,
  $OUTPUT_FIELD_SEPARATOR: 5,
  $OFS: 5,
  "$,": 5,
  $INPUT_LINE_NUMBER: 5,
  $NR: 5,
  "$.": 5,
  $INPUT_RECORD_SEPARATOR: 5,
  $RS: 5,
  "$/": 5,
  $OUTPUT_RECORD_SEPARATOR: 5,
  $ORS: 5,
  "$\\": 5,
  $OUTPUT_AUTOFLUSH: 5,
  "$|": 5,
  $ACCUMULATOR: 5,
  "$^A": 5,
  $FORMAT_FORMFEED: 5,
  "$^L": 5,
  $FORMAT_PAGE_NUMBER: 5,
  "$%": 5,
  $FORMAT_LINES_LEFT: 5,
  "$-": 5,
  $FORMAT_LINE_BREAK_CHARACTERS: 5,
  "$:": 5,
  $FORMAT_LINES_PER_PAGE: 5,
  "$=": 5,
  $FORMAT_TOP_NAME: 5,
  "$^": 5,
  $FORMAT_NAME: 5,
  "$~": 5,
  "${^CHILD_ERROR_NATIVE}": 5,
  $EXTENDED_OS_ERROR: 5,
  "$^E": 5,
  $EXCEPTIONS_BEING_CAUGHT: 5,
  "$^S": 5,
  $WARNING: 5,
  "$^W": 5,
  "${^WARNING_BITS}": 5,
  $OS_ERROR: 5,
  $ERRNO: 5,
  "$!": 5,
  "%OS_ERROR": 5,
  "%ERRNO": 5,
  "%!": 5,
  $CHILD_ERROR: 5,
  "$?": 5,
  $EVAL_ERROR: 5,
  "$@": 5,
  $OFMT: 5,
  "$#": 5,
  "$*": 5,
  $ARRAY_BASE: 5,
  "$[": 5,
  $OLD_PERL_VERSION: 5,
  "$]": 5,
  if: [1, 1],
  elsif: [1, 1],
  else: [1, 1],
  while: [1, 1],
  unless: [1, 1],
  for: [1, 1],
  foreach: [1, 1],
  abs: 1,
  accept: 1,
  alarm: 1,
  atan2: 1,
  bind: 1,
  binmode: 1,
  bless: 1,
  bootstrap: 1,
  break: 1,
  caller: 1,
  chdir: 1,
  chmod: 1,
  chomp: 1,
  chop: 1,
  chown: 1,
  chr: 1,
  chroot: 1,
  close: 1,
  closedir: 1,
  connect: 1,
  continue: [1, 1],
  cos: 1,
  crypt: 1,
  dbmclose: 1,
  dbmopen: 1,
  default: 1,
  defined: 1,
  delete: 1,
  die: 1,
  do: 1,
  dump: 1,
  each: 1,
  endgrent: 1,
  endhostent: 1,
  endnetent: 1,
  endprotoent: 1,
  endpwent: 1,
  endservent: 1,
  eof: 1,
  eval: 1,
  exec: 1,
  exists: 1,
  exit: 1,
  exp: 1,
  fcntl: 1,
  fileno: 1,
  flock: 1,
  fork: 1,
  format: 1,
  formline: 1,
  getc: 1,
  getgrent: 1,
  getgrgid: 1,
  getgrnam: 1,
  gethostbyaddr: 1,
  gethostbyname: 1,
  gethostent: 1,
  getlogin: 1,
  getnetbyaddr: 1,
  getnetbyname: 1,
  getnetent: 1,
  getpeername: 1,
  getpgrp: 1,
  getppid: 1,
  getpriority: 1,
  getprotobyname: 1,
  getprotobynumber: 1,
  getprotoent: 1,
  getpwent: 1,
  getpwnam: 1,
  getpwuid: 1,
  getservbyname: 1,
  getservbyport: 1,
  getservent: 1,
  getsockname: 1,
  getsockopt: 1,
  given: 1,
  glob: 1,
  gmtime: 1,
  goto: 1,
  grep: 1,
  hex: 1,
  import: 1,
  index: 1,
  int: 1,
  ioctl: 1,
  join: 1,
  keys: 1,
  kill: 1,
  last: 1,
  lc: 1,
  lcfirst: 1,
  length: 1,
  link: 1,
  listen: 1,
  local: 2,
  localtime: 1,
  lock: 1,
  log: 1,
  lstat: 1,
  m: null,
  map: 1,
  mkdir: 1,
  msgctl: 1,
  msgget: 1,
  msgrcv: 1,
  msgsnd: 1,
  my: 2,
  new: 1,
  next: 1,
  no: 1,
  oct: 1,
  open: 1,
  opendir: 1,
  ord: 1,
  our: 2,
  pack: 1,
  package: 1,
  pipe: 1,
  pop: 1,
  pos: 1,
  print: 1,
  printf: 1,
  prototype: 1,
  push: 1,
  q: null,
  qq: null,
  qr: null,
  quotemeta: null,
  qw: null,
  qx: null,
  rand: 1,
  read: 1,
  readdir: 1,
  readline: 1,
  readlink: 1,
  readpipe: 1,
  recv: 1,
  redo: 1,
  ref: 1,
  rename: 1,
  require: 1,
  reset: 1,
  return: 1,
  reverse: 1,
  rewinddir: 1,
  rindex: 1,
  rmdir: 1,
  s: null,
  say: 1,
  scalar: 1,
  seek: 1,
  seekdir: 1,
  select: 1,
  semctl: 1,
  semget: 1,
  semop: 1,
  send: 1,
  setgrent: 1,
  sethostent: 1,
  setnetent: 1,
  setpgrp: 1,
  setpriority: 1,
  setprotoent: 1,
  setpwent: 1,
  setservent: 1,
  setsockopt: 1,
  shift: 1,
  shmctl: 1,
  shmget: 1,
  shmread: 1,
  shmwrite: 1,
  shutdown: 1,
  sin: 1,
  sleep: 1,
  socket: 1,
  socketpair: 1,
  sort: 1,
  splice: 1,
  split: 1,
  sprintf: 1,
  sqrt: 1,
  srand: 1,
  stat: 1,
  state: 1,
  study: 1,
  sub: 1,
  substr: 1,
  symlink: 1,
  syscall: 1,
  sysopen: 1,
  sysread: 1,
  sysseek: 1,
  system: 1,
  syswrite: 1,
  tell: 1,
  telldir: 1,
  tie: 1,
  tied: 1,
  time: 1,
  times: 1,
  tr: null,
  truncate: 1,
  uc: 1,
  ucfirst: 1,
  umask: 1,
  undef: 1,
  unlink: 1,
  unpack: 1,
  unshift: 1,
  untie: 1,
  use: 1,
  utime: 1,
  values: 1,
  vec: 1,
  wait: 1,
  waitpid: 1,
  wantarray: 1,
  warn: 1,
  when: 1,
  write: 1,
  y: null
};
var RXstyle = "string.special";
var RXmodifiers = /[goseximacplud]/;
function tokenChain(stream, state, chain3, style, tail) {
  state.chain = null;
  state.style = null;
  state.tail = null;
  state.tokenize = function(stream2, state2) {
    var e = false, c2, i = 0;
    while (c2 = stream2.next()) {
      if (c2 === chain3[i] && !e) {
        if (chain3[++i] !== void 0) {
          state2.chain = chain3[i];
          state2.style = style;
          state2.tail = tail;
        } else if (tail)
          stream2.eatWhile(tail);
        state2.tokenize = tokenPerl;
        return style;
      }
      e = !e && c2 == "\\";
    }
    return style;
  };
  return state.tokenize(stream, state);
}
function tokenSOMETHING(stream, state, string3) {
  state.tokenize = function(stream2, state2) {
    if (stream2.string == string3)
      state2.tokenize = tokenPerl;
    stream2.skipToEnd();
    return "string";
  };
  return state.tokenize(stream, state);
}
function tokenPerl(stream, state) {
  if (stream.eatSpace())
    return null;
  if (state.chain)
    return tokenChain(stream, state, state.chain, state.style, state.tail);
  if (stream.match(/^\-?[\d\.]/, false)) {
    if (stream.match(/^(\-?(\d*\.\d+(e[+-]?\d+)?|\d+\.\d*)|0x[\da-fA-F]+|0b[01]+|\d+(e[+-]?\d+)?)/))
      return "number";
  }
  if (stream.match(/^<<(?=[_a-zA-Z])/)) {
    stream.eatWhile(/\w/);
    return tokenSOMETHING(stream, state, stream.current().substr(2));
  }
  if (stream.sol() && stream.match(/^\=item(?!\w)/)) {
    return tokenSOMETHING(stream, state, "=cut");
  }
  var ch = stream.next();
  if (ch == '"' || ch == "'") {
    if (prefix(stream, 3) == "<<" + ch) {
      var p = stream.pos;
      stream.eatWhile(/\w/);
      var n = stream.current().substr(1);
      if (n && stream.eat(ch))
        return tokenSOMETHING(stream, state, n);
      stream.pos = p;
    }
    return tokenChain(stream, state, [ch], "string");
  }
  if (ch == "q") {
    var c2 = look(stream, -2);
    if (!(c2 && /\w/.test(c2))) {
      c2 = look(stream, 0);
      if (c2 == "x") {
        c2 = look(stream, 1);
        if (c2 == "(") {
          eatSuffix(stream, 2);
          return tokenChain(stream, state, [")"], RXstyle, RXmodifiers);
        }
        if (c2 == "[") {
          eatSuffix(stream, 2);
          return tokenChain(stream, state, ["]"], RXstyle, RXmodifiers);
        }
        if (c2 == "{") {
          eatSuffix(stream, 2);
          return tokenChain(stream, state, ["}"], RXstyle, RXmodifiers);
        }
        if (c2 == "<") {
          eatSuffix(stream, 2);
          return tokenChain(stream, state, [">"], RXstyle, RXmodifiers);
        }
        if (/[\^'"!~\/]/.test(c2)) {
          eatSuffix(stream, 1);
          return tokenChain(stream, state, [stream.eat(c2)], RXstyle, RXmodifiers);
        }
      } else if (c2 == "q") {
        c2 = look(stream, 1);
        if (c2 == "(") {
          eatSuffix(stream, 2);
          return tokenChain(stream, state, [")"], "string");
        }
        if (c2 == "[") {
          eatSuffix(stream, 2);
          return tokenChain(stream, state, ["]"], "string");
        }
        if (c2 == "{") {
          eatSuffix(stream, 2);
          return tokenChain(stream, state, ["}"], "string");
        }
        if (c2 == "<") {
          eatSuffix(stream, 2);
          return tokenChain(stream, state, [">"], "string");
        }
        if (/[\^'"!~\/]/.test(c2)) {
          eatSuffix(stream, 1);
          return tokenChain(stream, state, [stream.eat(c2)], "string");
        }
      } else if (c2 == "w") {
        c2 = look(stream, 1);
        if (c2 == "(") {
          eatSuffix(stream, 2);
          return tokenChain(stream, state, [")"], "bracket");
        }
        if (c2 == "[") {
          eatSuffix(stream, 2);
          return tokenChain(stream, state, ["]"], "bracket");
        }
        if (c2 == "{") {
          eatSuffix(stream, 2);
          return tokenChain(stream, state, ["}"], "bracket");
        }
        if (c2 == "<") {
          eatSuffix(stream, 2);
          return tokenChain(stream, state, [">"], "bracket");
        }
        if (/[\^'"!~\/]/.test(c2)) {
          eatSuffix(stream, 1);
          return tokenChain(stream, state, [stream.eat(c2)], "bracket");
        }
      } else if (c2 == "r") {
        c2 = look(stream, 1);
        if (c2 == "(") {
          eatSuffix(stream, 2);
          return tokenChain(stream, state, [")"], RXstyle, RXmodifiers);
        }
        if (c2 == "[") {
          eatSuffix(stream, 2);
          return tokenChain(stream, state, ["]"], RXstyle, RXmodifiers);
        }
        if (c2 == "{") {
          eatSuffix(stream, 2);
          return tokenChain(stream, state, ["}"], RXstyle, RXmodifiers);
        }
        if (c2 == "<") {
          eatSuffix(stream, 2);
          return tokenChain(stream, state, [">"], RXstyle, RXmodifiers);
        }
        if (/[\^'"!~\/]/.test(c2)) {
          eatSuffix(stream, 1);
          return tokenChain(stream, state, [stream.eat(c2)], RXstyle, RXmodifiers);
        }
      } else if (/[\^'"!~\/(\[{<]/.test(c2)) {
        if (c2 == "(") {
          eatSuffix(stream, 1);
          return tokenChain(stream, state, [")"], "string");
        }
        if (c2 == "[") {
          eatSuffix(stream, 1);
          return tokenChain(stream, state, ["]"], "string");
        }
        if (c2 == "{") {
          eatSuffix(stream, 1);
          return tokenChain(stream, state, ["}"], "string");
        }
        if (c2 == "<") {
          eatSuffix(stream, 1);
          return tokenChain(stream, state, [">"], "string");
        }
        if (/[\^'"!~\/]/.test(c2)) {
          return tokenChain(stream, state, [stream.eat(c2)], "string");
        }
      }
    }
  }
  if (ch == "m") {
    var c2 = look(stream, -2);
    if (!(c2 && /\w/.test(c2))) {
      c2 = stream.eat(/[(\[{<\^'"!~\/]/);
      if (c2) {
        if (/[\^'"!~\/]/.test(c2)) {
          return tokenChain(stream, state, [c2], RXstyle, RXmodifiers);
        }
        if (c2 == "(") {
          return tokenChain(stream, state, [")"], RXstyle, RXmodifiers);
        }
        if (c2 == "[") {
          return tokenChain(stream, state, ["]"], RXstyle, RXmodifiers);
        }
        if (c2 == "{") {
          return tokenChain(stream, state, ["}"], RXstyle, RXmodifiers);
        }
        if (c2 == "<") {
          return tokenChain(stream, state, [">"], RXstyle, RXmodifiers);
        }
      }
    }
  }
  if (ch == "s") {
    var c2 = /[\/>\]})\w]/.test(look(stream, -2));
    if (!c2) {
      c2 = stream.eat(/[(\[{<\^'"!~\/]/);
      if (c2) {
        if (c2 == "[")
          return tokenChain(stream, state, ["]", "]"], RXstyle, RXmodifiers);
        if (c2 == "{")
          return tokenChain(stream, state, ["}", "}"], RXstyle, RXmodifiers);
        if (c2 == "<")
          return tokenChain(stream, state, [">", ">"], RXstyle, RXmodifiers);
        if (c2 == "(")
          return tokenChain(stream, state, [")", ")"], RXstyle, RXmodifiers);
        return tokenChain(stream, state, [c2, c2], RXstyle, RXmodifiers);
      }
    }
  }
  if (ch == "y") {
    var c2 = /[\/>\]})\w]/.test(look(stream, -2));
    if (!c2) {
      c2 = stream.eat(/[(\[{<\^'"!~\/]/);
      if (c2) {
        if (c2 == "[")
          return tokenChain(stream, state, ["]", "]"], RXstyle, RXmodifiers);
        if (c2 == "{")
          return tokenChain(stream, state, ["}", "}"], RXstyle, RXmodifiers);
        if (c2 == "<")
          return tokenChain(stream, state, [">", ">"], RXstyle, RXmodifiers);
        if (c2 == "(")
          return tokenChain(stream, state, [")", ")"], RXstyle, RXmodifiers);
        return tokenChain(stream, state, [c2, c2], RXstyle, RXmodifiers);
      }
    }
  }
  if (ch == "t") {
    var c2 = /[\/>\]})\w]/.test(look(stream, -2));
    if (!c2) {
      c2 = stream.eat("r");
      if (c2) {
        c2 = stream.eat(/[(\[{<\^'"!~\/]/);
        if (c2) {
          if (c2 == "[")
            return tokenChain(stream, state, ["]", "]"], RXstyle, RXmodifiers);
          if (c2 == "{")
            return tokenChain(stream, state, ["}", "}"], RXstyle, RXmodifiers);
          if (c2 == "<")
            return tokenChain(stream, state, [">", ">"], RXstyle, RXmodifiers);
          if (c2 == "(")
            return tokenChain(stream, state, [")", ")"], RXstyle, RXmodifiers);
          return tokenChain(stream, state, [c2, c2], RXstyle, RXmodifiers);
        }
      }
    }
  }
  if (ch == "`") {
    return tokenChain(stream, state, [ch], "builtin");
  }
  if (ch == "/") {
    if (!/~\s*$/.test(prefix(stream)))
      return "operator";
    else
      return tokenChain(stream, state, [ch], RXstyle, RXmodifiers);
  }
  if (ch == "$") {
    var p = stream.pos;
    if (stream.eatWhile(/\d/) || stream.eat("{") && stream.eatWhile(/\d/) && stream.eat("}"))
      return "builtin";
    else
      stream.pos = p;
  }
  if (/[$@%]/.test(ch)) {
    var p = stream.pos;
    if (stream.eat("^") && stream.eat(/[A-Z]/) || !/[@$%&]/.test(look(stream, -2)) && stream.eat(/[=|\\\-#?@;:&`~\^!\[\]*'"$+.,\/<>()]/)) {
      var c2 = stream.current();
      if (PERL[c2])
        return "builtin";
    }
    stream.pos = p;
  }
  if (/[$@%&]/.test(ch)) {
    if (stream.eatWhile(/[\w$]/) || stream.eat("{") && stream.eatWhile(/[\w$]/) && stream.eat("}")) {
      var c2 = stream.current();
      if (PERL[c2])
        return "builtin";
      else
        return "variable";
    }
  }
  if (ch == "#") {
    if (look(stream, -2) != "$") {
      stream.skipToEnd();
      return "comment";
    }
  }
  if (/[:+\-\^*$&%@=<>!?|\/~\.]/.test(ch)) {
    var p = stream.pos;
    stream.eatWhile(/[:+\-\^*$&%@=<>!?|\/~\.]/);
    if (PERL[stream.current()])
      return "operator";
    else
      stream.pos = p;
  }
  if (ch == "_") {
    if (stream.pos == 1) {
      if (suffix(stream, 6) == "_END__") {
        return tokenChain(stream, state, ["\0"], "comment");
      } else if (suffix(stream, 7) == "_DATA__") {
        return tokenChain(stream, state, ["\0"], "builtin");
      } else if (suffix(stream, 7) == "_C__") {
        return tokenChain(stream, state, ["\0"], "string");
      }
    }
  }
  if (/\w/.test(ch)) {
    var p = stream.pos;
    if (look(stream, -2) == "{" && (look(stream, 0) == "}" || stream.eatWhile(/\w/) && look(stream, 0) == "}"))
      return "string";
    else
      stream.pos = p;
  }
  if (/[A-Z]/.test(ch)) {
    var l = look(stream, -2);
    var p = stream.pos;
    stream.eatWhile(/[A-Z_]/);
    if (/[\da-z]/.test(look(stream, 0))) {
      stream.pos = p;
    } else {
      var c2 = PERL[stream.current()];
      if (!c2)
        return "meta";
      if (c2[1])
        c2 = c2[0];
      if (l != ":") {
        if (c2 == 1)
          return "keyword";
        else if (c2 == 2)
          return "def";
        else if (c2 == 3)
          return "atom";
        else if (c2 == 4)
          return "operator";
        else if (c2 == 5)
          return "builtin";
        else
          return "meta";
      } else
        return "meta";
    }
  }
  if (/[a-zA-Z_]/.test(ch)) {
    var l = look(stream, -2);
    stream.eatWhile(/\w/);
    var c2 = PERL[stream.current()];
    if (!c2)
      return "meta";
    if (c2[1])
      c2 = c2[0];
    if (l != ":") {
      if (c2 == 1)
        return "keyword";
      else if (c2 == 2)
        return "def";
      else if (c2 == 3)
        return "atom";
      else if (c2 == 4)
        return "operator";
      else if (c2 == 5)
        return "builtin";
      else
        return "meta";
    } else
      return "meta";
  }
  return null;
}
var perl = {
  startState: function() {
    return {
      tokenize: tokenPerl,
      chain: null,
      style: null,
      tail: null
    };
  },
  token: function(stream, state) {
    return (state.tokenize || tokenPerl)(stream, state);
  },
  languageData: {
    commentTokens: {line: "#"},
    wordChars: "$"
  }
};

// node_modules/@codemirror/legacy-modes/mode/powershell.js
function buildRegexp(patterns, options) {
  options = options || {};
  var prefix2 = options.prefix !== void 0 ? options.prefix : "^";
  var suffix2 = options.suffix !== void 0 ? options.suffix : "\\b";
  for (var i = 0; i < patterns.length; i++) {
    if (patterns[i] instanceof RegExp) {
      patterns[i] = patterns[i].source;
    } else {
      patterns[i] = patterns[i].replace(/[-\/\\^$*+?.()|[\]{}]/g, "\\$&");
    }
  }
  return new RegExp(prefix2 + "(" + patterns.join("|") + ")" + suffix2, "i");
}
var notCharacterOrDash = "(?=[^A-Za-z\\d\\-_]|$)";
var varNames = /[\w\-:]/;
var keywords7 = buildRegexp([
  /begin|break|catch|continue|data|default|do|dynamicparam/,
  /else|elseif|end|exit|filter|finally|for|foreach|from|function|if|in/,
  /param|process|return|switch|throw|trap|try|until|where|while/
], {suffix: notCharacterOrDash});
var punctuation2 = /[\[\]{},;`\\\.]|@[({]/;
var wordOperators = buildRegexp([
  "f",
  /b?not/,
  /[ic]?split/,
  "join",
  /is(not)?/,
  "as",
  /[ic]?(eq|ne|[gl][te])/,
  /[ic]?(not)?(like|match|contains)/,
  /[ic]?replace/,
  /b?(and|or|xor)/
], {prefix: "-"});
var symbolOperators = /[+\-*\/%]=|\+\+|--|\.\.|[+\-*&^%:=!|\/]|<(?!#)|(?!#)>/;
var operators3 = buildRegexp([wordOperators, symbolOperators], {suffix: ""});
var numbers = /^((0x[\da-f]+)|((\d+\.\d+|\d\.|\.\d+|\d+)(e[\+\-]?\d+)?))[ld]?([kmgtp]b)?/i;
var identifiers2 = /^[A-Za-z\_][A-Za-z\-\_\d]*\b/;
var symbolBuiltins = /[A-Z]:|%|\?/i;
var namedBuiltins = buildRegexp([
  /Add-(Computer|Content|History|Member|PSSnapin|Type)/,
  /Checkpoint-Computer/,
  /Clear-(Content|EventLog|History|Host|Item(Property)?|Variable)/,
  /Compare-Object/,
  /Complete-Transaction/,
  /Connect-PSSession/,
  /ConvertFrom-(Csv|Json|SecureString|StringData)/,
  /Convert-Path/,
  /ConvertTo-(Csv|Html|Json|SecureString|Xml)/,
  /Copy-Item(Property)?/,
  /Debug-Process/,
  /Disable-(ComputerRestore|PSBreakpoint|PSRemoting|PSSessionConfiguration)/,
  /Disconnect-PSSession/,
  /Enable-(ComputerRestore|PSBreakpoint|PSRemoting|PSSessionConfiguration)/,
  /(Enter|Exit)-PSSession/,
  /Export-(Alias|Clixml|Console|Counter|Csv|FormatData|ModuleMember|PSSession)/,
  /ForEach-Object/,
  /Format-(Custom|List|Table|Wide)/,
  new RegExp("Get-(Acl|Alias|AuthenticodeSignature|ChildItem|Command|ComputerRestorePoint|Content|ControlPanelItem|Counter|Credential|Culture|Date|Event|EventLog|EventSubscriber|ExecutionPolicy|FormatData|Help|History|Host|HotFix|Item|ItemProperty|Job|Location|Member|Module|PfxCertificate|Process|PSBreakpoint|PSCallStack|PSDrive|PSProvider|PSSession|PSSessionConfiguration|PSSnapin|Random|Service|TraceSource|Transaction|TypeData|UICulture|Unique|Variable|Verb|WinEvent|WmiObject)"),
  /Group-Object/,
  /Import-(Alias|Clixml|Counter|Csv|LocalizedData|Module|PSSession)/,
  /ImportSystemModules/,
  /Invoke-(Command|Expression|History|Item|RestMethod|WebRequest|WmiMethod)/,
  /Join-Path/,
  /Limit-EventLog/,
  /Measure-(Command|Object)/,
  /Move-Item(Property)?/,
  new RegExp("New-(Alias|Event|EventLog|Item(Property)?|Module|ModuleManifest|Object|PSDrive|PSSession|PSSessionConfigurationFile|PSSessionOption|PSTransportOption|Service|TimeSpan|Variable|WebServiceProxy|WinEvent)"),
  /Out-(Default|File|GridView|Host|Null|Printer|String)/,
  /Pause/,
  /(Pop|Push)-Location/,
  /Read-Host/,
  /Receive-(Job|PSSession)/,
  /Register-(EngineEvent|ObjectEvent|PSSessionConfiguration|WmiEvent)/,
  /Remove-(Computer|Event|EventLog|Item(Property)?|Job|Module|PSBreakpoint|PSDrive|PSSession|PSSnapin|TypeData|Variable|WmiObject)/,
  /Rename-(Computer|Item(Property)?)/,
  /Reset-ComputerMachinePassword/,
  /Resolve-Path/,
  /Restart-(Computer|Service)/,
  /Restore-Computer/,
  /Resume-(Job|Service)/,
  /Save-Help/,
  /Select-(Object|String|Xml)/,
  /Send-MailMessage/,
  new RegExp("Set-(Acl|Alias|AuthenticodeSignature|Content|Date|ExecutionPolicy|Item(Property)?|Location|PSBreakpoint|PSDebug|PSSessionConfiguration|Service|StrictMode|TraceSource|Variable|WmiInstance)"),
  /Show-(Command|ControlPanelItem|EventLog)/,
  /Sort-Object/,
  /Split-Path/,
  /Start-(Job|Process|Service|Sleep|Transaction|Transcript)/,
  /Stop-(Computer|Job|Process|Service|Transcript)/,
  /Suspend-(Job|Service)/,
  /TabExpansion2/,
  /Tee-Object/,
  /Test-(ComputerSecureChannel|Connection|ModuleManifest|Path|PSSessionConfigurationFile)/,
  /Trace-Command/,
  /Unblock-File/,
  /Undo-Transaction/,
  /Unregister-(Event|PSSessionConfiguration)/,
  /Update-(FormatData|Help|List|TypeData)/,
  /Use-Transaction/,
  /Wait-(Event|Job|Process)/,
  /Where-Object/,
  /Write-(Debug|Error|EventLog|Host|Output|Progress|Verbose|Warning)/,
  /cd|help|mkdir|more|oss|prompt/,
  /ac|asnp|cat|cd|chdir|clc|clear|clhy|cli|clp|cls|clv|cnsn|compare|copy|cp|cpi|cpp|cvpa|dbp|del|diff|dir|dnsn|ebp/,
  /echo|epal|epcsv|epsn|erase|etsn|exsn|fc|fl|foreach|ft|fw|gal|gbp|gc|gci|gcm|gcs|gdr|ghy|gi|gjb|gl|gm|gmo|gp|gps/,
  /group|gsn|gsnp|gsv|gu|gv|gwmi|h|history|icm|iex|ihy|ii|ipal|ipcsv|ipmo|ipsn|irm|ise|iwmi|iwr|kill|lp|ls|man|md/,
  /measure|mi|mount|move|mp|mv|nal|ndr|ni|nmo|npssc|nsn|nv|ogv|oh|popd|ps|pushd|pwd|r|rbp|rcjb|rcsn|rd|rdr|ren|ri/,
  /rjb|rm|rmdir|rmo|rni|rnp|rp|rsn|rsnp|rujb|rv|rvpa|rwmi|sajb|sal|saps|sasv|sbp|sc|select|set|shcm|si|sl|sleep|sls/,
  /sort|sp|spjb|spps|spsv|start|sujb|sv|swmi|tee|trcm|type|where|wjb|write/
], {prefix: "", suffix: ""});
var variableBuiltins = buildRegexp([
  /[$?^_]|Args|ConfirmPreference|ConsoleFileName|DebugPreference|Error|ErrorActionPreference|ErrorView|ExecutionContext/,
  /FormatEnumerationLimit|Home|Host|Input|MaximumAliasCount|MaximumDriveCount|MaximumErrorCount|MaximumFunctionCount/,
  /MaximumHistoryCount|MaximumVariableCount|MyInvocation|NestedPromptLevel|OutputEncoding|Pid|Profile|ProgressPreference/,
  /PSBoundParameters|PSCommandPath|PSCulture|PSDefaultParameterValues|PSEmailServer|PSHome|PSScriptRoot|PSSessionApplicationName/,
  /PSSessionConfigurationName|PSSessionOption|PSUICulture|PSVersionTable|Pwd|ShellId|StackTrace|VerbosePreference/,
  /WarningPreference|WhatIfPreference/,
  /Event|EventArgs|EventSubscriber|Sender/,
  /Matches|Ofs|ForEach|LastExitCode|PSCmdlet|PSItem|PSSenderInfo|This/,
  /true|false|null/
], {prefix: "\\$", suffix: ""});
var builtins5 = buildRegexp([symbolBuiltins, namedBuiltins, variableBuiltins], {suffix: notCharacterOrDash});
var grammar = {
  keyword: keywords7,
  number: numbers,
  operator: operators3,
  builtin: builtins5,
  punctuation: punctuation2,
  variable: identifiers2
};
function tokenBase5(stream, state) {
  var parent = state.returnStack[state.returnStack.length - 1];
  if (parent && parent.shouldReturnFrom(state)) {
    state.tokenize = parent.tokenize;
    state.returnStack.pop();
    return state.tokenize(stream, state);
  }
  if (stream.eatSpace()) {
    return null;
  }
  if (stream.eat("(")) {
    state.bracketNesting += 1;
    return "punctuation";
  }
  if (stream.eat(")")) {
    state.bracketNesting -= 1;
    return "punctuation";
  }
  for (var key in grammar) {
    if (stream.match(grammar[key])) {
      return key;
    }
  }
  var ch = stream.next();
  if (ch === "'") {
    return tokenSingleQuoteString(stream, state);
  }
  if (ch === "$") {
    return tokenVariable(stream, state);
  }
  if (ch === '"') {
    return tokenDoubleQuoteString(stream, state);
  }
  if (ch === "<" && stream.eat("#")) {
    state.tokenize = tokenComment3;
    return tokenComment3(stream, state);
  }
  if (ch === "#") {
    stream.skipToEnd();
    return "comment";
  }
  if (ch === "@") {
    var quoteMatch = stream.eat(/["']/);
    if (quoteMatch && stream.eol()) {
      state.tokenize = tokenMultiString;
      state.startQuote = quoteMatch[0];
      return tokenMultiString(stream, state);
    } else if (stream.eol()) {
      return "error";
    } else if (stream.peek().match(/[({]/)) {
      return "punctuation";
    } else if (stream.peek().match(varNames)) {
      return tokenVariable(stream, state);
    }
  }
  return "error";
}
function tokenSingleQuoteString(stream, state) {
  var ch;
  while ((ch = stream.peek()) != null) {
    stream.next();
    if (ch === "'" && !stream.eat("'")) {
      state.tokenize = tokenBase5;
      return "string";
    }
  }
  return "error";
}
function tokenDoubleQuoteString(stream, state) {
  var ch;
  while ((ch = stream.peek()) != null) {
    if (ch === "$") {
      state.tokenize = tokenStringInterpolation;
      return "string";
    }
    stream.next();
    if (ch === "`") {
      stream.next();
      continue;
    }
    if (ch === '"' && !stream.eat('"')) {
      state.tokenize = tokenBase5;
      return "string";
    }
  }
  return "error";
}
function tokenStringInterpolation(stream, state) {
  return tokenInterpolation2(stream, state, tokenDoubleQuoteString);
}
function tokenMultiStringReturn(stream, state) {
  state.tokenize = tokenMultiString;
  state.startQuote = '"';
  return tokenMultiString(stream, state);
}
function tokenHereStringInterpolation(stream, state) {
  return tokenInterpolation2(stream, state, tokenMultiStringReturn);
}
function tokenInterpolation2(stream, state, parentTokenize) {
  if (stream.match("$(")) {
    var savedBracketNesting = state.bracketNesting;
    state.returnStack.push({
      shouldReturnFrom: function(state2) {
        return state2.bracketNesting === savedBracketNesting;
      },
      tokenize: parentTokenize
    });
    state.tokenize = tokenBase5;
    state.bracketNesting += 1;
    return "punctuation";
  } else {
    stream.next();
    state.returnStack.push({
      shouldReturnFrom: function() {
        return true;
      },
      tokenize: parentTokenize
    });
    state.tokenize = tokenVariable;
    return state.tokenize(stream, state);
  }
}
function tokenComment3(stream, state) {
  var maybeEnd = false, ch;
  while ((ch = stream.next()) != null) {
    if (maybeEnd && ch == ">") {
      state.tokenize = tokenBase5;
      break;
    }
    maybeEnd = ch === "#";
  }
  return "comment";
}
function tokenVariable(stream, state) {
  var ch = stream.peek();
  if (stream.eat("{")) {
    state.tokenize = tokenVariableWithBraces;
    return tokenVariableWithBraces(stream, state);
  } else if (ch != void 0 && ch.match(varNames)) {
    stream.eatWhile(varNames);
    state.tokenize = tokenBase5;
    return "variable";
  } else {
    state.tokenize = tokenBase5;
    return "error";
  }
}
function tokenVariableWithBraces(stream, state) {
  var ch;
  while ((ch = stream.next()) != null) {
    if (ch === "}") {
      state.tokenize = tokenBase5;
      break;
    }
  }
  return "variable";
}
function tokenMultiString(stream, state) {
  var quote = state.startQuote;
  if (stream.sol() && stream.match(new RegExp(quote + "@"))) {
    state.tokenize = tokenBase5;
  } else if (quote === '"') {
    while (!stream.eol()) {
      var ch = stream.peek();
      if (ch === "$") {
        state.tokenize = tokenHereStringInterpolation;
        return "string";
      }
      stream.next();
      if (ch === "`") {
        stream.next();
      }
    }
  } else {
    stream.skipToEnd();
  }
  return "string";
}
var powerShell = {
  startState: function() {
    return {
      returnStack: [],
      bracketNesting: 0,
      tokenize: tokenBase5
    };
  },
  token: function(stream, state) {
    return state.tokenize(stream, state);
  },
  languageData: {
    commentTokens: {line: "#", block: {open: "<#", close: "#>"}}
  }
};

// node_modules/lezer-python/dist/index.es.js
var printKeyword = 1;
var indent = 162;
var dedent = 163;
var newline2 = 164;
var newlineBracketed = 165;
var newlineEmpty = 166;
var eof = 167;
var ParenthesizedExpression = 21;
var TupleExpression = 47;
var ComprehensionExpression = 48;
var ArrayExpression = 52;
var ArrayComprehensionExpression = 55;
var DictionaryExpression = 56;
var DictionaryComprehensionExpression = 59;
var SetExpression = 60;
var SetComprehensionExpression = 61;
var newline$1 = 10;
var carriageReturn = 13;
var space2 = 32;
var tab = 9;
var hash = 35;
var parenOpen = 40;
var dot = 46;
var bracketed2 = [
  ParenthesizedExpression,
  TupleExpression,
  ComprehensionExpression,
  ArrayExpression,
  ArrayComprehensionExpression,
  DictionaryExpression,
  DictionaryComprehensionExpression,
  SetExpression,
  SetComprehensionExpression
];
var cachedIndent = 0;
var cachedInput = null;
var cachedPos = 0;
function getIndent(input, pos) {
  if (pos == cachedPos && input == cachedInput)
    return cachedIndent;
  cachedInput = input;
  cachedPos = pos;
  return cachedIndent = getIndentInner(input, pos);
}
function getIndentInner(input, pos) {
  for (let indent2 = 0; ; pos++) {
    let ch = input.get(pos);
    if (ch == space2)
      indent2++;
    else if (ch == tab)
      indent2 += 8 - indent2 % 8;
    else if (ch == newline$1 || ch == carriageReturn || ch == hash)
      return -1;
    else
      return indent2;
  }
}
var newlines = new ExternalTokenizer((input, token, stack) => {
  let next = input.get(token.start);
  if (next < 0) {
    token.accept(eof, token.start);
  } else if (next != newline$1 && next != carriageReturn)
    ;
  else if (stack.startOf(bracketed2) != null) {
    token.accept(newlineBracketed, token.start + 1);
  } else if (getIndent(input, token.start + 1) < 0) {
    token.accept(newlineEmpty, token.start + 1);
  } else {
    token.accept(newline2, token.start + 1);
  }
}, {contextual: true, fallback: true});
var indentation = new ExternalTokenizer((input, token, stack) => {
  let prev = input.get(token.start - 1), depth;
  if ((prev == newline$1 || prev == carriageReturn) && (depth = getIndent(input, token.start)) >= 0 && depth != stack.context.depth && stack.startOf(bracketed2) == null)
    token.accept(depth < stack.context.depth ? dedent : indent, token.start);
});
function IndentLevel(parent, depth) {
  this.parent = parent;
  this.depth = depth;
  this.hash = (parent ? parent.hash + parent.hash << 8 : 0) + depth + (depth << 4);
}
var topIndent2 = new IndentLevel(null, 0);
var trackIndent = new ContextTracker({
  start: topIndent2,
  shift(context, term, input, stack) {
    return term == indent ? new IndentLevel(context, getIndent(input, stack.pos)) : term == dedent ? context.parent : context;
  },
  hash(context) {
    return context.hash;
  }
});
var legacyPrint = new ExternalTokenizer((input, token) => {
  let pos = token.start;
  for (let print = "print", i = 0; i < print.length; i++, pos++)
    if (input.get(pos) != print.charCodeAt(i))
      return;
  let end = pos;
  if (/\w/.test(String.fromCharCode(input.get(pos))))
    return;
  for (; ; pos++) {
    let next = input.get(pos);
    if (next == space2 || next == tab)
      continue;
    if (next != parenOpen && next != dot && next != newline$1 && next != carriageReturn && next != hash)
      token.accept(printKeyword, end);
    return;
  }
});
var spec_identifier3 = {__proto__: null, await: 40, or: 48, and: 50, in: 54, not: 56, is: 58, if: 64, else: 66, lambda: 70, yield: 88, from: 90, async: 98, for: 100, None: 152, True: 154, False: 154, del: 168, pass: 172, break: 176, continue: 180, return: 184, raise: 192, import: 196, as: 198, global: 202, nonlocal: 204, assert: 208, elif: 218, while: 222, try: 228, except: 230, finally: 232, with: 236, def: 240, class: 250};
var parser3 = Parser.deserialize({
  version: 13,
  states: "!?|O`Q$IXOOO%cQ$I[O'#GaOOQ$IS'#Cm'#CmOOQ$IS'#Cn'#CnO'RQ$IWO'#ClO(tQ$I[O'#G`OOQ$IS'#Ga'#GaOOQ$IS'#DR'#DROOQ$IS'#G`'#G`O)bQ$IWO'#CqO)rQ$IWO'#DbO*SQ$IWO'#DfOOQ$IS'#Ds'#DsO*gO`O'#DsO*oOpO'#DsO*wO!bO'#DtO+SO#tO'#DtO+_O&jO'#DtO+jO,UO'#DtO-lQ$I[O'#GQOOQ$IS'#GQ'#GQO'RQ$IWO'#GPO/OQ$I[O'#GPOOQ$IS'#E]'#E]O/gQ$IWO'#E^OOQ$IS'#GO'#GOO/qQ$IWO'#F}OOQ$IV'#F}'#F}O/|Q$IWO'#FPOOQ$IS'#Fr'#FrO0RQ$IWO'#FOOOQ$IV'#HZ'#HZOOQ$IV'#F|'#F|OOQ$IT'#FR'#FRQ`Q$IXOOO'RQ$IWO'#CoO0aQ$IWO'#CzO0hQ$IWO'#DOO0vQ$IWO'#GeO1WQ$I[O'#EQO'RQ$IWO'#EROOQ$IS'#ET'#ETOOQ$IS'#EV'#EVOOQ$IS'#EX'#EXO1lQ$IWO'#EZO2SQ$IWO'#E_O/|Q$IWO'#EaO2gQ$I[O'#EaO/|Q$IWO'#EdO/gQ$IWO'#EgO/gQ$IWO'#EkO/gQ$IWO'#EnO2rQ$IWO'#EpO2yQ$IWO'#EuO3UQ$IWO'#EqO/gQ$IWO'#EuO/|Q$IWO'#EwO/|Q$IWO'#E|OOQ$IS'#Cc'#CcOOQ$IS'#Cd'#CdOOQ$IS'#Ce'#CeOOQ$IS'#Cf'#CfOOQ$IS'#Cg'#CgOOQ$IS'#Ch'#ChOOQ$IS'#Cj'#CjO'RQ$IWO,58|O'RQ$IWO,58|O'RQ$IWO,58|O'RQ$IWO,58|O'RQ$IWO,58|O'RQ$IWO,58|O3ZQ$IWO'#DmOOQ$IS,5:W,5:WO3nQ$IWO,5:ZO3{Q%1`O,5:ZO4QQ$I[O,59WO0aQ$IWO,59_O0aQ$IWO,59_O0aQ$IWO,59_O6pQ$IWO,59_O6uQ$IWO,59_O6|Q$IWO,59gO7TQ$IWO'#G`O8ZQ$IWO'#G_OOQ$IS'#G_'#G_OOQ$IS'#DX'#DXO8rQ$IWO,59]O'RQ$IWO,59]O9QQ$IWO,59]O9VQ$IWO,5:PO'RQ$IWO,5:POOQ$IS,59|,59|O9eQ$IWO,59|O9jQ$IWO,5:VO'RQ$IWO,5:VO'RQ$IWO,5:TOOQ$IS,5:Q,5:QO9{Q$IWO,5:QO:QQ$IWO,5:UOOOO'#FZ'#FZO:VO`O,5:_OOQ$IS,5:_,5:_OOOO'#F['#F[O:_OpO,5:_O:gQ$IWO'#DuOOOO'#F]'#F]O:wO!bO,5:`OOQ$IS,5:`,5:`OOOO'#F`'#F`O;SO#tO,5:`OOOO'#Fa'#FaO;_O&jO,5:`OOOO'#Fb'#FbO;jO,UO,5:`OOQ$IS'#Fc'#FcO;uQ$I[O,5:dO>gQ$I[O,5<kO?QQ%GlO,5<kO?qQ$I[O,5<kOOQ$IS,5:x,5:xO@YQ$IXO'#FkOAiQ$IWO,5;TOOQ$IV,5<i,5<iOAtQ$I[O'#HWOB]Q$IWO,5;kOOQ$IS-E9p-E9pOOQ$IV,5;j,5;jO3PQ$IWO'#EwOOQ$IT-E9P-E9POBeQ$I[O,59ZODlQ$I[O,59fOEVQ$IWO'#GbOEbQ$IWO'#GbO/|Q$IWO'#GbOEmQ$IWO'#DQOEuQ$IWO,59jOEzQ$IWO'#GfO'RQ$IWO'#GfO/gQ$IWO,5=POOQ$IS,5=P,5=PO/gQ$IWO'#D|OOQ$IS'#D}'#D}OFiQ$IWO'#FeOFyQ$IWO,58zOGXQ$IWO,58zO)eQ$IWO,5:jOG^Q$I[O'#GhOOQ$IS,5:m,5:mOOQ$IS,5:u,5:uOGqQ$IWO,5:yOHSQ$IWO,5:{OOQ$IS'#Fh'#FhOHbQ$I[O,5:{OHpQ$IWO,5:{OHuQ$IWO'#HYOOQ$IS,5;O,5;OOITQ$IWO'#HVOOQ$IS,5;R,5;RO3UQ$IWO,5;VO3UQ$IWO,5;YOIfQ$I[O'#H[O'RQ$IWO'#H[OIpQ$IWO,5;[O2rQ$IWO,5;[O/gQ$IWO,5;aO/|Q$IWO,5;cOIuQ$IXO'#ElOKOQ$IZO,5;]ONaQ$IWO'#H]O3UQ$IWO,5;aONlQ$IWO,5;cONqQ$IWO,5;hO!#fQ$I[O1G.hO!#mQ$I[O1G.hO!&^Q$I[O1G.hO!&hQ$I[O1G.hO!)RQ$I[O1G.hO!)fQ$I[O1G.hO!)yQ$IWO'#GnO!*XQ$I[O'#GQO/gQ$IWO'#GnO!*cQ$IWO'#GmOOQ$IS,5:X,5:XO!*kQ$IWO,5:XO!*pQ$IWO'#GoO!*{Q$IWO'#GoO!+`Q$IWO1G/uOOQ$IS'#Dq'#DqOOQ$IS1G/u1G/uOOQ$IS1G.y1G.yO!,`Q$I[O1G.yO!,gQ$I[O1G.yO0aQ$IWO1G.yO!-SQ$IWO1G/ROOQ$IS'#DW'#DWO/gQ$IWO,59qOOQ$IS1G.w1G.wO!-ZQ$IWO1G/cO!-kQ$IWO1G/cO!-sQ$IWO1G/dO'RQ$IWO'#GgO!-xQ$IWO'#GgO!-}Q$I[O1G.wO!._Q$IWO,59fO!/eQ$IWO,5=VO!/uQ$IWO,5=VO!/}Q$IWO1G/kO!0SQ$I[O1G/kOOQ$IS1G/h1G/hO!0dQ$IWO,5=QO!1ZQ$IWO,5=QO/gQ$IWO1G/oO!1xQ$IWO1G/qO!1}Q$I[O1G/qO!2_Q$I[O1G/oOOQ$IS1G/l1G/lOOQ$IS1G/p1G/pOOOO-E9X-E9XOOQ$IS1G/y1G/yOOOO-E9Y-E9YO!2oQ$IWO'#GzO/gQ$IWO'#GzO!2}Q$IWO,5:aOOOO-E9Z-E9ZOOQ$IS1G/z1G/zOOOO-E9^-E9^OOOO-E9_-E9_OOOO-E9`-E9`OOQ$IS-E9a-E9aO!3YQ%GlO1G2VO!3yQ$I[O1G2VO'RQ$IWO,5<OOOQ$IS,5<O,5<OOOQ$IS-E9b-E9bOOQ$IS,5<V,5<VOOQ$IS-E9i-E9iOOQ$IV1G0o1G0oO/|Q$IWO'#FgO!4bQ$I[O,5=rOOQ$IS1G1V1G1VO!4yQ$IWO1G1VOOQ$IS'#DS'#DSO/gQ$IWO,5<|OOQ$IS,5<|,5<|O!5OQ$IWO'#FSO!5ZQ$IWO,59lO!5cQ$IWO1G/UO!5mQ$I[O,5=QOOQ$IS1G2k1G2kOOQ$IS,5:h,5:hO!6^Q$IWO'#GPOOQ$IS,5<P,5<POOQ$IS-E9c-E9cO!6oQ$IWO1G.fOOQ$IS1G0U1G0UO!6}Q$IWO,5=SO!7_Q$IWO,5=SO/gQ$IWO1G0eO/gQ$IWO1G0eO/|Q$IWO1G0gOOQ$IS-E9f-E9fO!7pQ$IWO1G0gO!7{Q$IWO1G0gO!8QQ$IWO,5=tO!8`Q$IWO,5=tO!8nQ$IWO,5=qO!9UQ$IWO,5=qO!9gQ$IZO1G0qO!<uQ$IZO1G0tO!@QQ$IWO,5=vO!@[Q$IWO,5=vO!@dQ$I[O,5=vO/gQ$IWO1G0vO!@nQ$IWO1G0vO3UQ$IWO1G0{ONlQ$IWO1G0}OOQ$IV,5;W,5;WO!@sQ$IYO,5;WO!@xQ$IZO1G0wO!DZQ$IWO'#FoO3UQ$IWO1G0wO3UQ$IWO1G0wO!DhQ$IWO,5=wO!DuQ$IWO,5=wO/|Q$IWO,5=wOOQ$IV1G0{1G0{O!D}Q$IWO'#EyO!E`Q%1`O1G0}OOQ$IV1G1S1G1SO3UQ$IWO1G1SOOQ$IS,5=Y,5=YOOQ$IS'#Dn'#DnO/gQ$IWO,5=YO!EhQ$IWO,5=XO!E{Q$IWO,5=XOOQ$IS1G/s1G/sO!FTQ$IWO,5=ZO!FeQ$IWO,5=ZO!FmQ$IWO,5=ZO!GQQ$IWO,5=ZO!GbQ$IWO,5=ZOOQ$IS7+%a7+%aOOQ$IS7+$e7+$eO!5cQ$IWO7+$mO!ITQ$IWO1G.yO!I[Q$IWO1G.yOOQ$IS1G/]1G/]OOQ$IS,5;p,5;pO'RQ$IWO,5;pOOQ$IS7+$}7+$}O!IcQ$IWO7+$}OOQ$IS-E9S-E9SOOQ$IS7+%O7+%OO!IsQ$IWO,5=RO'RQ$IWO,5=ROOQ$IS7+$c7+$cO!IxQ$IWO7+$}O!JQQ$IWO7+%OO!JVQ$IWO1G2qOOQ$IS7+%V7+%VO!JgQ$IWO1G2qO!JoQ$IWO7+%VOOQ$IS,5;o,5;oO'RQ$IWO,5;oO!JtQ$IWO1G2lOOQ$IS-E9R-E9RO!KkQ$IWO7+%ZOOQ$IS7+%]7+%]O!KyQ$IWO1G2lO!LhQ$IWO7+%]O!LmQ$IWO1G2rO!L}Q$IWO1G2rO!MVQ$IWO7+%ZO!M[Q$IWO,5=fO!MrQ$IWO,5=fO!MrQ$IWO,5=fO!NQO!LQO'#DwO!N]OSO'#G{OOOO1G/{1G/{O!NbQ$IWO1G/{O!NjQ%GlO7+'qO# ZQ$I[O1G1jP# tQ$IWO'#FdOOQ$IS,5<R,5<ROOQ$IS-E9e-E9eOOQ$IS7+&q7+&qOOQ$IS1G2h1G2hOOQ$IS,5;n,5;nOOQ$IS-E9Q-E9QOOQ$IS7+$p7+$pO#!RQ$IWO,5<kO#!lQ$IWO,5<kO#!}Q$I[O,5;qO##bQ$IWO1G2nOOQ$IS-E9T-E9TOOQ$IS7+&P7+&PO##rQ$IWO7+&POOQ$IS7+&R7+&RO#$QQ$IWO'#HXO/|Q$IWO7+&RO#$fQ$IWO7+&ROOQ$IS,5<U,5<UO#$qQ$IWO1G3`OOQ$IS-E9h-E9hOOQ$IS,5<Q,5<QO#%PQ$IWO1G3]OOQ$IS-E9d-E9dO#%gQ$IZO7+&]O!DZQ$IWO'#FmO3UQ$IWO7+&]O3UQ$IWO7+&`O#(uQ$I[O,5<YO'RQ$IWO,5<YO#)PQ$IWO1G3bOOQ$IS-E9l-E9lO#)ZQ$IWO1G3bO3UQ$IWO7+&bO/gQ$IWO7+&bOOQ$IV7+&g7+&gO!E`Q%1`O7+&iO#)cQ$IXO1G0rOOQ$IV-E9m-E9mO3UQ$IWO7+&cO3UQ$IWO7+&cOOQ$IV,5<Z,5<ZO#+UQ$IWO,5<ZOOQ$IV7+&c7+&cO#+aQ$IZO7+&cO#.lQ$IWO,5<[O#.wQ$IWO1G3cOOQ$IS-E9n-E9nO#/UQ$IWO1G3cO#/^Q$IWO'#H_O#/lQ$IWO'#H_O/|Q$IWO'#H_OOQ$IS'#H_'#H_O#/wQ$IWO'#H^OOQ$IS,5;e,5;eO#0PQ$IWO,5;eO/gQ$IWO'#E{OOQ$IV7+&i7+&iO3UQ$IWO7+&iOOQ$IV7+&n7+&nOOQ$IS1G2t1G2tOOQ$IS,5;s,5;sO#0UQ$IWO1G2sOOQ$IS-E9V-E9VO#0iQ$IWO,5;tO#0tQ$IWO,5;tO#1XQ$IWO1G2uOOQ$IS-E9W-E9WO#1iQ$IWO1G2uO#1qQ$IWO1G2uO#2RQ$IWO1G2uO#1iQ$IWO1G2uOOQ$IS<<HX<<HXO#2^Q$I[O1G1[OOQ$IS<<Hi<<HiP#2kQ$IWO'#FUO6|Q$IWO1G2mO#2xQ$IWO1G2mO#2}Q$IWO<<HiOOQ$IS<<Hj<<HjO#3_Q$IWO7+(]OOQ$IS<<Hq<<HqO#3oQ$I[O1G1ZP#4`Q$IWO'#FTO#4mQ$IWO7+(^O#4}Q$IWO7+(^O#5VQ$IWO<<HuO#5[Q$IWO7+(WOOQ$IS<<Hw<<HwO#6RQ$IWO,5;rO'RQ$IWO,5;rOOQ$IS-E9U-E9UOOQ$IS<<Hu<<HuOOQ$IS,5;x,5;xO/gQ$IWO,5;xO#6WQ$IWO1G3QOOQ$IS-E9[-E9[O#6nQ$IWO1G3QOOOO'#F_'#F_O#6|O!LQO,5:cOOOO,5=g,5=gOOOO7+%g7+%gO#7XQ$IWO1G2VO#7rQ$IWO1G2VP'RQ$IWO'#FVO/gQ$IWO<<IkO#8TQ$IWO,5=sO#8fQ$IWO,5=sO/|Q$IWO,5=sO#8wQ$IWO<<ImOOQ$IS<<Im<<ImO/|Q$IWO<<ImP/|Q$IWO'#FjP/gQ$IWO'#FfOOQ$IV-E9k-E9kO3UQ$IWO<<IwOOQ$IV,5<X,5<XO3UQ$IWO,5<XOOQ$IV<<Iw<<IwOOQ$IV<<Iz<<IzO#8|Q$I[O1G1tP#9WQ$IWO'#FnO#9_Q$IWO7+(|O#9iQ$IZO<<I|O3UQ$IWO<<I|OOQ$IV<<JT<<JTO3UQ$IWO<<JTOOQ$IV'#Fl'#FlO#<tQ$IZO7+&^OOQ$IV<<I}<<I}O#>mQ$IZO<<I}OOQ$IV1G1u1G1uO/|Q$IWO1G1uO3UQ$IWO<<I}O/|Q$IWO1G1vP/gQ$IWO'#FpO#AxQ$IWO7+(}O#BVQ$IWO7+(}OOQ$IS'#Ez'#EzO/gQ$IWO,5=yO#B_Q$IWO,5=yOOQ$IS,5=y,5=yO#BjQ$IWO,5=xO#B{Q$IWO,5=xOOQ$IS1G1P1G1POOQ$IS,5;g,5;gP#CTQ$IWO'#FXO#CeQ$IWO1G1`O#CxQ$IWO1G1`O#DYQ$IWO1G1`P#DeQ$IWO'#FYO#DrQ$IWO7+(aO#ESQ$IWO7+(aO#ESQ$IWO7+(aO#E[Q$IWO7+(aO#ElQ$IWO7+(XO6|Q$IWO7+(XOOQ$ISAN>TAN>TO#FVQ$IWO<<KxOOQ$ISAN>aAN>aO/gQ$IWO1G1^O#FgQ$I[O1G1^P#FqQ$IWO'#FWOOQ$IS1G1d1G1dP#GOQ$IWO'#F^O#G]Q$IWO7+(lOOOO-E9]-E9]O#GsQ$IWO7+'qOOQ$ISAN?VAN?VO#H^Q$IWO,5<TO#HrQ$IWO1G3_OOQ$IS-E9g-E9gO#ITQ$IWO1G3_OOQ$ISAN?XAN?XO#IfQ$IWOAN?XOOQ$IVAN?cAN?cOOQ$IV1G1s1G1sO3UQ$IWOAN?hO#IkQ$IZOAN?hOOQ$IVAN?oAN?oOOQ$IV-E9j-E9jOOQ$IV<<Ix<<IxO3UQ$IWOAN?iO3UQ$IWO7+'aOOQ$IVAN?iAN?iOOQ$IS7+'b7+'bO#LvQ$IWO<<LiOOQ$IS1G3e1G3eO/gQ$IWO1G3eOOQ$IS,5<],5<]O#MTQ$IWO1G3dOOQ$IS-E9o-E9oO#MfQ$IWO7+&zO#MvQ$IWO7+&zOOQ$IS7+&z7+&zO#NRQ$IWO<<K{O#NcQ$IWO<<K{O#NcQ$IWO<<K{O#NkQ$IWO'#GiOOQ$IS<<Ks<<KsO#NuQ$IWO<<KsOOQ$IS7+&x7+&xO/|Q$IWO1G1oP/|Q$IWO'#FiO$ `Q$IWO7+(yO$ qQ$IWO7+(yOOQ$ISG24sG24sOOQ$IVG25SG25SO3UQ$IWOG25SOOQ$IVG25TG25TOOQ$IV<<J{<<J{OOQ$IS7+)P7+)PP$!SQ$IWO'#FqOOQ$IS<<Jf<<JfO$!bQ$IWO<<JfO$!rQ$IWOANAgO$#SQ$IWOANAgO$#[Q$IWO'#GjOOQ$IS'#Gj'#GjO0hQ$IWO'#DaO$#uQ$IWO,5=TOOQ$ISANA_ANA_OOQ$IS7+'Z7+'ZO$$^Q$IWO<<LeOOQ$IVLD*nLD*nOOQ$ISAN@QAN@QO$$oQ$IWOG27RO$%PQ$IWO,59{OOQ$IS1G2o1G2oO#NkQ$IWO1G/gOOQ$IS7+%R7+%RO6|Q$IWO'#CzO6|Q$IWO,59_O6|Q$IWO,59_O6|Q$IWO,59_O$%UQ$I[O,5<kO6|Q$IWO1G.yO/gQ$IWO1G/UO/gQ$IWO7+$mP$%iQ$IWO'#FdO'RQ$IWO'#GPO$%vQ$IWO,59_O$%{Q$IWO,59_O$&SQ$IWO,59jO$&XQ$IWO1G/RO0hQ$IWO'#DOO6|Q$IWO,59g",
  stateData: "$&o~O$oOS$lOS$kOSQOS~OPhOTeOdsOfXOltOp!SOsuO|vO}!PO!R!VO!S!UO!VYO!ZZO!fdO!mdO!ndO!odO!vxO!xyO!zzO!|{O#O|O#S}O#U!OO#X!QO#Y!QO#[!RO#c!TO#f!WO#j!XO#l!YO#q!ZO#tlO$jqO$zQO${QO%PRO%QVO%e[O%f]O%i^O%l_O%r`O%uaO%wbO~OT!aO]!aO_!bOf!iO!V!kO!d!lO$u![O$v!]O$w!^O$x!_O$y!_O$z!`O${!`O$|!aO$}!aO%O!aO~Oh%TXi%TXj%TXk%TXl%TXm%TXp%TXw%TXx%TX!s%TX#^%TX$j%TX$m%TX%V%TX!O%TX!R%TX!S%TX%W%TX!W%TX![%TX}%TX#V%TXq%TX!j%TX~P$_OdsOfXO!VYO!ZZO!fdO!mdO!ndO!odO$zQO${QO%PRO%QVO%e[O%f]O%i^O%l_O%r`O%uaO%wbO~Ow%SXx%SX#^%SX$j%SX$m%SX%V%SX~Oh!oOi!pOj!nOk!nOl!qOm!rOp!sO!s%SX~P(`OT!yOl-fOs-tO|vO~P'ROT!|Ol-fOs-tO!W!}O~P'ROT#QO_#ROl-fOs-tO![#SO~P'RO%g#VO%h#XO~O%j#YO%k#XO~O!Z#[O%m#]O%q#_O~O!Z#[O%s#`O%t#_O~O!Z#[O%h#_O%v#bO~O!Z#[O%k#_O%x#dO~OT$tX]$tX_$tXf$tXh$tXi$tXj$tXk$tXl$tXm$tXp$tXw$tX!V$tX!d$tX$u$tX$v$tX$w$tX$x$tX$y$tX$z$tX${$tX$|$tX$}$tX%O$tX!O$tX!R$tX!S$tX~O%e[O%f]O%i^O%l_O%r`O%uaO%wbOx$tX!s$tX#^$tX$j$tX$m$tX%V$tX%W$tX!W$tX![$tX}$tX#V$tXq$tX!j$tX~P+uOw#iOx$sX!s$sX#^$sX$j$sX$m$sX%V$sX~Ol-fOs-tO~P'RO#^#lO$j#nO$m#nO~O%QVO~O!R#sO#l!YO#q!ZO#tlO~OltO~P'ROT#xO_#yO%QVOxtP~OT#}Ol-fOs-tO}$OO~P'ROx$QO!s$VO%V$RO#^!tX$j!tX$m!tX~OT#}Ol-fOs-tO#^!}X$j!}X$m!}X~P'ROl-fOs-tO#^#RX$j#RX$m#RX~P'RO!d$]O!m$]O%QVO~OT$gO~P'RO!S$iO#j$jO#l$kO~Ox$lO~OT$zO_$zOl-fOs-tO!O$|O~P'ROl-fOs-tOx%PO~P'RO%d%RO~O_!bOf!iO!V!kO!d!lOT`a]`ah`ai`aj`ak`al`am`ap`aw`ax`a!s`a#^`a$j`a$m`a$u`a$v`a$w`a$x`a$y`a$z`a${`a$|`a$}`a%O`a%V`a!O`a!R`a!S`a%W`a!W`a![`a}`a#V`aq`a!j`a~Ok%WO~Ol%WO~P'ROl-fO~P'ROh-hOi-iOj-gOk-gOl-pOm-qOp-uO!O%SX!R%SX!S%SX%W%SX!W%SX![%SX}%SX#V%SX!j%SX~P(`O%W%YOw%RX!O%RX!R%RX!S%RX!W%RXx%RX~Ow%]O!O%[O!R%aO!S%`O~O!O%[O~Ow%dO!R%aO!S%`O!W%_X~O!W%hO~Ow%iOx%kO!R%aO!S%`O![%YX~O![%oO~O![%pO~O%g#VO%h%rO~O%j#YO%k%rO~OT%uOl-fOs-tO|vO~P'RO!Z#[O%m#]O%q%xO~O!Z#[O%s#`O%t%xO~O!Z#[O%h%xO%v#bO~O!Z#[O%k%xO%x#dO~OT!la]!la_!laf!lah!lai!laj!lak!lal!lam!lap!law!lax!la!V!la!d!la!s!la#^!la$j!la$m!la$u!la$v!la$w!la$x!la$y!la$z!la${!la$|!la$}!la%O!la%V!la!O!la!R!la!S!la%W!la!W!la![!la}!la#V!laq!la!j!la~P#vOw%}Ox$sa!s$sa#^$sa$j$sa$m$sa%V$sa~P$_OT&POltOsuOx$sa!s$sa#^$sa$j$sa$m$sa%V$sa~P'ROw%}Ox$sa!s$sa#^$sa$j$sa$m$sa%V$sa~OPhOTeOltOsuO|vO}!PO!vxO!xyO!zzO!|{O#O|O#S}O#U!OO#X!QO#Y!QO#[!RO#^$_X$j$_X$m$_X~P'RO#^#lO$j&UO$m&UO~O!d&VOf%zX$j%zX#V%zX#^%zX$m%zX#U%zX~Of!iO$j&XO~Ohcaicajcakcalcamcapcawcaxca!sca#^ca$jca$mca%Vca!Oca!Rca!Sca%Wca!Wca![ca}ca#Vcaqca!jca~P$_Opnawnaxna#^na$jna$mna%Vna~Oh!oOi!pOj!nOk!nOl!qOm!rO!sna~PDTO%V&ZOw%UXx%UX~O%QVOw%UXx%UX~Ow&^OxtX~Ox&`O~Ow%iO#^%YX$j%YX$m%YX!O%YXx%YX![%YX!j%YX%V%YX~OT-oOl-fOs-tO|vO~P'RO%V$RO#^Sa$jSa$mSa~O%V$RO~Ow&iO#^%[X$j%[X$m%[Xk%[X~P$_Ow&lO}&kO#^#Ra$j#Ra$m#Ra~O#V&mO#^#Ta$j#Ta$m#Ta~O!d$]O!m$]O#U&oO%QVO~O#U&oO~Ow&qO#^%|X$j%|X$m%|X~Ow&sO#^%yX$j%yX$m%yXx%yX~Ow&wOk&OX~P$_Ok&zO~OPhOTeOltOsuO|vO}!PO!vxO!xyO!zzO!|{O#O|O#S}O#U!OO#X!QO#Y!QO#[!RO$j'PO~P'ROq'TO#g'RO#h'SOP#eaT#ead#eaf#eal#eap#eas#ea|#ea}#ea!R#ea!S#ea!V#ea!Z#ea!f#ea!m#ea!n#ea!o#ea!v#ea!x#ea!z#ea!|#ea#O#ea#S#ea#U#ea#X#ea#Y#ea#[#ea#c#ea#f#ea#j#ea#l#ea#q#ea#t#ea$g#ea$j#ea$z#ea${#ea%P#ea%Q#ea%e#ea%f#ea%i#ea%l#ea%r#ea%u#ea%w#ea$i#ea$m#ea~Ow'UO#V'WOx&PX~Of'YO~Of!iOx$lO~OT!aO]!aO_!bOf!iO!V!kO!d!lO$w!^O$x!_O$y!_O$z!`O${!`O$|!aO$}!aO%O!aOhUiiUijUikUilUimUipUiwUixUi!sUi#^Ui$jUi$mUi$uUi%VUi!OUi!RUi!SUi%WUi!WUi![Ui}Ui#VUiqUi!jUi~O$v!]O~PNyO$vUi~PNyOT!aO]!aO_!bOf!iO!V!kO!d!lO$z!`O${!`O$|!aO$}!aO%O!aOhUiiUijUikUilUimUipUiwUixUi!sUi#^Ui$jUi$mUi$uUi$vUi$wUi%VUi!OUi!RUi!SUi%WUi!WUi![Ui}Ui#VUiqUi!jUi~O$x!_O$y!_O~P!#tO$xUi$yUi~P!#tO_!bOf!iO!V!kO!d!lOhUiiUijUikUilUimUipUiwUixUi!sUi#^Ui$jUi$mUi$uUi$vUi$wUi$xUi$yUi$zUi${Ui%VUi!OUi!RUi!SUi%WUi!WUi![Ui}Ui#VUiqUi!jUi~OT!aO]!aO$|!aO$}!aO%O!aO~P!&rOTUi]Ui$|Ui$}Ui%OUi~P!&rO!R%aO!S%`Ow%bX!O%bX~O%V'_O%W'_O~P+uOw'aO!O%aX~O!O'cO~Ow'dOx'fO!W%cX~Ol-fOs-tOw'dOx'gO!W%cX~P'RO!W'iO~Oj!nOk!nOl!qOm!rOhgipgiwgixgi!sgi#^gi$jgi$mgi%Vgi~Oi!pO~P!+eOigi~P!+eOh-hOi-iOj-gOk-gOl-pOm-qO~Oq'kO~P!,nOT'pOl-fOs-tO!O'qO~P'ROw'rO!O'qO~O!O'tO~O!S'vO~Ow'rO!O'wO!R%aO!S%`O~P$_Oh-hOi-iOj-gOk-gOl-pOm-qO!Ona!Rna!Sna%Wna!Wna![na}na#Vnaqna!jna~PDTOT'pOl-fOs-tO!W%_a~P'ROw'zO!W%_a~O!W'{O~Ow'zO!R%aO!S%`O!W%_a~P$_OT(POl-fOs-tO![%Ya#^%Ya$j%Ya$m%Ya!O%Yax%Ya!j%Ya%V%Ya~P'ROw(QO![%Ya#^%Ya$j%Ya$m%Ya!O%Yax%Ya!j%Ya%V%Ya~O![(TO~Ow(QO!R%aO!S%`O![%Ya~P$_Ow(WO!R%aO!S%`O![%`a~P$_Ow(ZOx%nX![%nX!j%nX~Ox(^O![(`O!j(aO~OT&POltOsuOx$si!s$si#^$si$j$si$m$si%V$si~P'ROw(bOx$si!s$si#^$si$j$si$m$si%V$si~O!d&VOf%za$j%za#V%za#^%za$m%za#U%za~O$j(gO~OT#xO_#yO%QVO~Ow&^Oxta~OltOsuO~P'ROw(QO#^%Ya$j%Ya$m%Ya!O%Yax%Ya![%Ya!j%Ya%V%Ya~P$_Ow(lO#^$sX$j$sX$m$sX%V$sX~O%V$RO#^Si$jSi$mSi~O#^%[a$j%[a$m%[ak%[a~P'ROw(oO#^%[a$j%[a$m%[ak%[a~OT(sOf(uO%QVO~O#U(vO~O%QVO#^%|a$j%|a$m%|a~Ow(xO#^%|a$j%|a$m%|a~Ol-fOs-tO#^%ya$j%ya$m%yax%ya~P'ROw({O#^%ya$j%ya$m%yax%ya~Oq)PO#a)OOP#_iT#_id#_if#_il#_ip#_is#_i|#_i}#_i!R#_i!S#_i!V#_i!Z#_i!f#_i!m#_i!n#_i!o#_i!v#_i!x#_i!z#_i!|#_i#O#_i#S#_i#U#_i#X#_i#Y#_i#[#_i#c#_i#f#_i#j#_i#l#_i#q#_i#t#_i$g#_i$j#_i$z#_i${#_i%P#_i%Q#_i%e#_i%f#_i%i#_i%l#_i%r#_i%u#_i%w#_i$i#_i$m#_i~Oq)QOP#biT#bid#bif#bil#bip#bis#bi|#bi}#bi!R#bi!S#bi!V#bi!Z#bi!f#bi!m#bi!n#bi!o#bi!v#bi!x#bi!z#bi!|#bi#O#bi#S#bi#U#bi#X#bi#Y#bi#[#bi#c#bi#f#bi#j#bi#l#bi#q#bi#t#bi$g#bi$j#bi$z#bi${#bi%P#bi%Q#bi%e#bi%f#bi%i#bi%l#bi%r#bi%u#bi%w#bi$i#bi$m#bi~OT)SOk&Oa~P'ROw)TOk&Oa~Ow)TOk&Oa~P$_Ok)XO~O$h)[O~Oq)_O#g'RO#h)^OP#eiT#eid#eif#eil#eip#eis#ei|#ei}#ei!R#ei!S#ei!V#ei!Z#ei!f#ei!m#ei!n#ei!o#ei!v#ei!x#ei!z#ei!|#ei#O#ei#S#ei#U#ei#X#ei#Y#ei#[#ei#c#ei#f#ei#j#ei#l#ei#q#ei#t#ei$g#ei$j#ei$z#ei${#ei%P#ei%Q#ei%e#ei%f#ei%i#ei%l#ei%r#ei%u#ei%w#ei$i#ei$m#ei~Ol-fOs-tOx$lO~P'ROl-fOs-tOx&Pa~P'ROw)eOx&Pa~OT)iO_)jO!O)mO$|)kO%QVO~Ox$lO&S)oO~OT$zO_$zOl-fOs-tO!O%aa~P'ROw)uO!O%aa~Ol-fOs-tOx)xO!W%ca~P'ROw)yO!W%ca~Ol-fOs-tOw)yOx)|O!W%ca~P'ROl-fOs-tOw)yO!W%ca~P'ROw)yOx)|O!W%ca~Oj-gOk-gOl-pOm-qOhgipgiwgi!Ogi!Rgi!Sgi%Wgi!Wgixgi![gi#^gi$jgi$mgi}gi#Vgiqgi!jgi%Vgi~Oi-iO~P!GmOigi~P!GmOT'pOl-fOs-tO!O*RO~P'ROk*TO~Ow*VO!O*RO~O!O*WO~OT'pOl-fOs-tO!W%_i~P'ROw*XO!W%_i~O!W*YO~OT(POl-fOs-tO![%Yi#^%Yi$j%Yi$m%Yi!O%Yix%Yi!j%Yi%V%Yi~P'ROw*]O!R%aO!S%`O![%`i~Ow*`O![%Yi#^%Yi$j%Yi$m%Yi!O%Yix%Yi!j%Yi%V%Yi~O![*aO~O_*cOl-fOs-tO![%`i~P'ROw*]O![%`i~O![*eO~OT*gOl-fOs-tOx%na![%na!j%na~P'ROw*hOx%na![%na!j%na~O!Z#[O%p*kO![!kX~O![*mO~Ox(^O![*nO~OT&POltOsuOx$sq!s$sq#^$sq$j$sq$m$sq%V$sq~P'ROw$Wix$Wi!s$Wi#^$Wi$j$Wi$m$Wi%V$Wi~P$_OT&POltOsuO~P'ROT&POl-fOs-tO#^$sa$j$sa$m$sa%V$sa~P'ROw*oO#^$sa$j$sa$m$sa%V$sa~Ow#ya#^#ya$j#ya$m#yak#ya~P$_O#^%[i$j%[i$m%[ik%[i~P'ROw*rO#^#Rq$j#Rq$m#Rq~Ow*sO#V*uO#^%{X$j%{X$m%{X!O%{X~OT*wOf*xO%QVO~O%QVO#^%|i$j%|i$m%|i~Ol-fOs-tO#^%yi$j%yi$m%yix%yi~P'ROq*|O#a)OOP#_qT#_qd#_qf#_ql#_qp#_qs#_q|#_q}#_q!R#_q!S#_q!V#_q!Z#_q!f#_q!m#_q!n#_q!o#_q!v#_q!x#_q!z#_q!|#_q#O#_q#S#_q#U#_q#X#_q#Y#_q#[#_q#c#_q#f#_q#j#_q#l#_q#q#_q#t#_q$g#_q$j#_q$z#_q${#_q%P#_q%Q#_q%e#_q%f#_q%i#_q%l#_q%r#_q%u#_q%w#_q$i#_q$m#_q~Ok$baw$ba~P$_OT)SOk&Oi~P'ROw+TOk&Oi~OPhOTeOltOp!SOsuO|vO}!PO!R!VO!S!UO!vxO!xyO!zzO!|{O#O|O#S}O#U!OO#X!QO#Y!QO#[!RO#c!TO#f!WO#j!XO#l!YO#q!ZO#tlO~P'ROw+_Ox$lO#V+_O~O#h+`OP#eqT#eqd#eqf#eql#eqp#eqs#eq|#eq}#eq!R#eq!S#eq!V#eq!Z#eq!f#eq!m#eq!n#eq!o#eq!v#eq!x#eq!z#eq!|#eq#O#eq#S#eq#U#eq#X#eq#Y#eq#[#eq#c#eq#f#eq#j#eq#l#eq#q#eq#t#eq$g#eq$j#eq$z#eq${#eq%P#eq%Q#eq%e#eq%f#eq%i#eq%l#eq%r#eq%u#eq%w#eq$i#eq$m#eq~O#V+aOw$dax$da~Ol-fOs-tOx&Pi~P'ROw+cOx&Pi~Ox$QO%V+eOw&RX!O&RX~O%QVOw&RX!O&RX~Ow+iO!O&QX~O!O+kO~OT$zO_$zOl-fOs-tO!O%ai~P'ROx+nOw#|a!W#|a~Ol-fOs-tOx+oOw#|a!W#|a~P'ROl-fOs-tOx)xO!W%ci~P'ROw+rO!W%ci~Ol-fOs-tOw+rO!W%ci~P'ROw+rOx+uO!W%ci~Ow#xi!O#xi!W#xi~P$_OT'pOl-fOs-tO~P'ROk+wO~OT'pOl-fOs-tO!O+xO~P'ROT'pOl-fOs-tO!W%_q~P'ROw#wi![#wi#^#wi$j#wi$m#wi!O#wix#wi!j#wi%V#wi~P$_OT(POl-fOs-tO~P'RO_*cOl-fOs-tO![%`q~P'ROw+yO![%`q~O![+zO~OT(POl-fOs-tO![%Yq#^%Yq$j%Yq$m%Yq!O%Yqx%Yq!j%Yq%V%Yq~P'ROx+{O~OT*gOl-fOs-tOx%ni![%ni!j%ni~P'ROw,QOx%ni![%ni!j%ni~O!Z#[O%p*kO![!ka~OT&POl-fOs-tO#^$si$j$si$m$si%V$si~P'ROw,SO#^$si$j$si$m$si%V$si~O%QVO#^%{a$j%{a$m%{a!O%{a~Ow,VO#^%{a$j%{a$m%{a!O%{a~O!O,YO~Ok$biw$bi~P$_OT)SO~P'ROT)SOk&Oq~P'ROq,^OP#dyT#dyd#dyf#dyl#dyp#dys#dy|#dy}#dy!R#dy!S#dy!V#dy!Z#dy!f#dy!m#dy!n#dy!o#dy!v#dy!x#dy!z#dy!|#dy#O#dy#S#dy#U#dy#X#dy#Y#dy#[#dy#c#dy#f#dy#j#dy#l#dy#q#dy#t#dy$g#dy$j#dy$z#dy${#dy%P#dy%Q#dy%e#dy%f#dy%i#dy%l#dy%r#dy%u#dy%w#dy$i#dy$m#dy~OPhOTeOltOp!SOsuO|vO}!PO!R!VO!S!UO!vxO!xyO!zzO!|{O#O|O#S}O#U!OO#X!QO#Y!QO#[!RO#c!TO#f!WO#j!XO#l!YO#q!ZO#tlO$i,bO$m,bO~P'RO#h,cOP#eyT#eyd#eyf#eyl#eyp#eys#ey|#ey}#ey!R#ey!S#ey!V#ey!Z#ey!f#ey!m#ey!n#ey!o#ey!v#ey!x#ey!z#ey!|#ey#O#ey#S#ey#U#ey#X#ey#Y#ey#[#ey#c#ey#f#ey#j#ey#l#ey#q#ey#t#ey$g#ey$j#ey$z#ey${#ey%P#ey%Q#ey%e#ey%f#ey%i#ey%l#ey%r#ey%u#ey%w#ey$i#ey$m#ey~Ol-fOs-tOx&Pq~P'ROw,gOx&Pq~O%V+eOw&Ra!O&Ra~OT)iO_)jO$|)kO%QVO!O&Qa~Ow,kO!O&Qa~OT$zO_$zOl-fOs-tO~P'ROl-fOs-tOx,mOw#|i!W#|i~P'ROl-fOs-tOw#|i!W#|i~P'ROx,mOw#|i!W#|i~Ol-fOs-tOx)xO~P'ROl-fOs-tOx)xO!W%cq~P'ROw,pO!W%cq~Ol-fOs-tOw,pO!W%cq~P'ROp,sO!R%aO!S%`O!O%Zq!W%Zq![%Zqw%Zq~P!,nO_*cOl-fOs-tO![%`y~P'ROw#zi![#zi~P$_O_*cOl-fOs-tO~P'ROT*gOl-fOs-tO~P'ROT*gOl-fOs-tOx%nq![%nq!j%nq~P'ROT&POl-fOs-tO#^$sq$j$sq$m$sq%V$sq~P'RO#V,wOw$]a#^$]a$j$]a$m$]a!O$]a~O%QVO#^%{i$j%{i$m%{i!O%{i~Ow,yO#^%{i$j%{i$m%{i!O%{i~O!O,{O~Oq,}OP#d!RT#d!Rd#d!Rf#d!Rl#d!Rp#d!Rs#d!R|#d!R}#d!R!R#d!R!S#d!R!V#d!R!Z#d!R!f#d!R!m#d!R!n#d!R!o#d!R!v#d!R!x#d!R!z#d!R!|#d!R#O#d!R#S#d!R#U#d!R#X#d!R#Y#d!R#[#d!R#c#d!R#f#d!R#j#d!R#l#d!R#q#d!R#t#d!R$g#d!R$j#d!R$z#d!R${#d!R%P#d!R%Q#d!R%e#d!R%f#d!R%i#d!R%l#d!R%r#d!R%u#d!R%w#d!R$i#d!R$m#d!R~Ol-fOs-tOx&Py~P'ROT)iO_)jO$|)kO%QVO!O&Qi~Ol-fOs-tOw#|q!W#|q~P'ROx-TOw#|q!W#|q~Ol-fOs-tOx)xO!W%cy~P'ROw-UO!W%cy~Ol-fOs-YO~P'ROp,sO!R%aO!S%`O!O%Zy!W%Zy![%Zyw%Zy~P!,nO%QVO#^%{q$j%{q$m%{q!O%{q~Ow-^O#^%{q$j%{q$m%{q!O%{q~OT)iO_)jO$|)kO%QVO~Ol-fOs-tOw#|y!W#|y~P'ROl-fOs-tOx)xO!W%c!R~P'ROw-aO!W%c!R~Op%^X!O%^X!R%^X!S%^X!W%^X![%^Xw%^X~P!,nOp,sO!R%aO!S%`O!O%]a!W%]a![%]aw%]a~O%QVO#^%{y$j%{y$m%{y!O%{y~Ol-fOs-tOx)xO!W%c!Z~P'ROx-dO~Ow*oO#^$sa$j$sa$m$sa%V$sa~P$_OT&POl-fOs-tO~P'ROk-kO~Ol-kO~P'ROx-lO~Oq-mO~P!,nO%f%i%u%w%e!Z%m%s%v%x%l%r%l%Q~",
  goto: "!,u&SPPPP&TP&])n*T*k+S+l,VP,qP&]-_-_&]P&]P0pPPPPPP0p3`PP3`P5l5u:yPP:|;[;_PPP&]&]PP;k&]PP&]&]PP&]&]&]&];o<c&]P<fP<i<i@OP@d&]PPP@h@n&TP&T&TP&TP&TP&TP&TP&T&T&TP&TPP&TPP&TP@tP@{ARP@{P@{@{PPP@{PBzPCTCZCaBzP@{CgPCnCtCzDWDjDpDzEQEnEtEzFQF[FbFhFnFtFzG^GhGnGtGzHUH[HbHhHnHxIOIYI`PPPPPPPPPIiIqIzJUJaPPPPPPPPPPPPNv! `!%n!(zPP!)S!)b!)k!*a!*W!*j!*p!*s!*v!*y!+RPPPPPPPPPP!+U!+XPPPPPPPPP!+_!+k!+w!,T!,W!,^!,d!,j!,m]iOr#l$l)[+Z'odOSXYZehrstvx|}!R!S!T!U!X!c!d!e!f!g!h!i!k!n!o!p!r!s!y!|#Q#R#[#i#l#}$O$Q$S$V$g$i$j$l$z%P%W%Z%]%`%d%i%k%u%}&P&[&`&i&k&l&s&w&z'R'U'`'a'd'f'g'k'p'r'v'z(P(Q(W(Z(b(d(l(o({)O)S)T)X)[)e)o)u)x)y)|*S*T*V*X*[*]*`*c*g*h*o*q*r*z+S+T+Z+b+c+f+m+n+o+q+r+u+w+y+{+},P,Q,S,g,i,m,p,s-T-U-a-d-f-g-h-i-k-l-m-n-o-q-uw!cP#h#u$W$f%b%g%m%n&a&y(c(n)R*Q*Z+R+|-jy!dP#h#u$W$f$r%b%g%m%n&a&y(c(n)R*Q*Z+R+|-j{!eP#h#u$W$f$r$s%b%g%m%n&a&y(c(n)R*Q*Z+R+|-j}!fP#h#u$W$f$r$s$t%b%g%m%n&a&y(c(n)R*Q*Z+R+|-j!P!gP#h#u$W$f$r$s$t$u%b%g%m%n&a&y(c(n)R*Q*Z+R+|-j!R!hP#h#u$W$f$r$s$t$u$v%b%g%m%n&a&y(c(n)R*Q*Z+R+|-j!V!hP!m#h#u$W$f$r$s$t$u$v$w%b%g%m%n&a&y(c(n)R*Q*Z+R+|-j'oSOSXYZehrstvx|}!R!S!T!U!X!c!d!e!f!g!h!i!k!n!o!p!r!s!y!|#Q#R#[#i#l#}$O$Q$S$V$g$i$j$l$z%P%W%Z%]%`%d%i%k%u%}&P&[&`&i&k&l&s&w&z'R'U'`'a'd'f'g'k'p'r'v'z(P(Q(W(Z(b(d(l(o({)O)S)T)X)[)e)o)u)x)y)|*S*T*V*X*[*]*`*c*g*h*o*q*r*z+S+T+Z+b+c+f+m+n+o+q+r+u+w+y+{+},P,Q,S,g,i,m,p,s-T-U-a-d-f-g-h-i-k-l-m-n-o-q-u&ZUOXYZhrtv|}!R!S!T!X!i!k!n!o!p!r!s#[#i#l$O$Q$S$V$j$l$z%P%W%Z%]%d%i%k%u%}&[&`&k&l&s&z'R'U'`'a'd'f'g'k'r'z(Q(W(Z(b(d(l({)O)X)[)e)o)u)x)y)|*S*T*V*X*[*]*`*g*h*o*r*z+Z+b+c+f+m+n+o+q+r+u+w+y+{+},P,Q,S,g,i,m,p,s-T-U-a-d-f-g-h-i-k-l-m-n-q-u%eWOXYZhrv|}!R!S!T!X!i!k#[#i#l$O$Q$S$V$j$l$z%P%Z%]%d%i%k%u%}&[&`&k&l&s&z'R'U'`'a'd'f'g'k'r'z(Q(W(Z(b(d(l({)O)X)[)e)o)u)x)y)|*S*V*X*[*]*`*g*h*o*r*z+Z+b+c+f+m+n+o+q+r+u+y+{+},P,Q,S,g,i,m,p-T-U-a-l-m-nQ#{uQ-b-YR-r-t'fdOSXYZehrstvx|}!R!S!T!U!X!c!d!e!f!g!h!k!n!o!p!r!s!y!|#Q#R#[#i#l#}$O$Q$S$V$g$i$j$l$z%P%W%Z%]%`%d%i%k%u%}&P&[&`&i&k&l&s&w&z'R'U'`'d'f'g'k'p'r'v'z(P(Q(W(Z(b(d(l(o({)O)S)T)X)[)e)o)x)y)|*S*T*V*X*[*]*`*c*g*h*o*q*r*z+S+T+Z+b+c+f+n+o+q+r+u+w+y+{+},P,Q,S,g,i,m,p,s-T-U-a-d-f-g-h-i-k-l-m-n-o-q-uW#ol!O!P$^W#wu&^-Y-tQ$`!QQ$p!YQ$q!ZW$y!i'a)u+mS&]#x#yQ&}$kQ(e&VQ(s&mW(t&o(u(v*xU(w&q(x*yQ)g'WW)h'Y+i,k-RS+h)i)jY,U*s,V,x,y-^Q,X*uQ,d+_Q,f+aR-],wR&[#wi!vXY!S!T%]%d'r'z)O*S*V*XR%Z!uQ!zXQ%v#[Q&e$SR&h$VT-X,s-d!U!jP!m#h#u$W$f$r$s$t$u$v$w%b%g%m%n&a&y(c(n)R*Q*Z+R+|-jQ&Y#pR']$qR'`$yR%S!l'ncOSXYZehrstvx|}!R!S!T!U!X!c!d!e!f!g!h!i!k!n!o!p!r!s!y!|#Q#R#[#i#l#}$O$Q$S$V$g$i$j$l$z%P%W%Z%]%`%d%i%k%u%}&P&[&`&i&k&l&s&w&z'R'U'`'a'd'f'g'k'p'r'v'z(P(Q(W(Z(b(d(l(o({)O)S)T)X)[)e)o)u)x)y)|*S*T*V*X*[*]*`*c*g*h*o*q*r*z+S+T+Z+b+c+f+m+n+o+q+r+u+w+y+{+},P,Q,S,g,i,m,p,s-T-U-a-d-f-g-h-i-k-l-m-n-o-q-uT#fc#gS#]_#^S#``#aS#ba#cS#db#eT*k(^*lT(_%v(aQ$UwR+g)hX$Sw$T$U&gZkOr$l)[+ZXoOr)[+ZQ$m!WQ&u$dQ&v$eQ'X$oQ'[$qQ)Y&|Q)`'RQ)b'SQ)c'TQ)p'ZQ)r']Q*})OQ+P)PQ+Q)QQ+U)WS+W)Z)qQ+[)^Q+])_Q+^)aQ,[*|Q,]+OQ,_+VQ,`+XQ,e+`Q,|,^Q-O,cQ-P,dR-_,}WoOr)[+ZR#rnQ'Z$pR)Z&}Q+f)hR,i+gQ)q'ZR+X)ZZmOnr)[+ZQrOR#trQ&_#zR(j&_S%j#P#|S(R%j(UT(U%m&aQ%^!xQ%e!{W's%^%e'x'|Q'x%bR'|%gQ&j$WR(p&jQ(X%nQ*^(ST*d(X*^Q'b${R)v'bS'e%O%PY)z'e){+s,q-VU){'f'g'hU+s)|)}*OS,q+t+uR-V,rQ#W]R%q#WQ#Z^R%s#ZQ#^_R%w#^Q([%tS*i([*jR*j(]Q*l(^R,R*lQ#a`R%y#aQ#caR%z#cQ#ebR%{#eQ#gcR%|#gQ#jfQ&O#hW&R#j&O(m*pQ(m&dR*p-jQ$TwS&f$T&gR&g$UQ&t$bR(|&tQ&W#oR(f&WQ$^!PR&n$^Q*t(tS,W*t,zR,z,XQ&r$`R(y&rQ#mjR&T#mQ+Z)[R,a+ZQ(}&uR*{(}Q&x$fS)U&x)VR)V&yQ'Q$mR)]'QQ'V$nS)f'V+dR+d)gQ+j)lR,l+jWnOr)[+ZR#qnSqOrT+Y)[+ZWpOr)[+ZR'O$lYjOr$l)[+ZR&S#l[wOr#l$l)[+ZR&e$S&YPOXYZhrtv|}!R!S!T!X!i!k!n!o!p!r!s#[#i#l$O$Q$S$V$j$l$z%P%W%Z%]%d%i%k%u%}&[&`&k&l&s&z'R'U'`'a'd'f'g'k'r'z(Q(W(Z(b(d(l({)O)X)[)e)o)u)x)y)|*S*T*V*X*[*]*`*g*h*o*r*z+Z+b+c+f+m+n+o+q+r+u+w+y+{+},P,Q,S,g,i,m,p,s-T-U-a-d-f-g-h-i-k-l-m-n-q-uQ!mSQ#heQ#usU$Wx%`'vS$f!U$iQ$r!cQ$s!dQ$t!eQ$u!fQ$v!gQ$w!hQ%b!yQ%g!|Q%m#QQ%n#RQ&a#}Q&y$gQ(c&PU(n&i(o*qW)R&w)T+S+TQ*Q'pQ*Z(PQ+R)SQ+|*cR-j-oQ!xXQ!{YQ$d!SQ$e!T^'o%]%d'r'z*S*V*XR+O)O[fOr#l$l)[+Zh!uXY!S!T%]%d'r'z)O*S*V*XQ#PZQ#khS#|v|Q$Z}W$b!R$V&z)XS$n!X$jW$x!i'a)u+mQ%O!kQ%t#[`&Q#i%}(b(d(l*o,S-nQ&b$OQ&c$QQ&d$SQ'^$zQ'h%PQ'n%ZW(O%i(Q*[*`Q(S%kQ(]%uQ(h&[S(k&`-lQ(q&kQ(r&lU(z&s({*zQ)a'RY)d'U)e+b+c,gQ)s'`^)w'd)y+q+r,p-U-aQ)}'fQ*O'gS*P'k-mW*b(W*]+y+}W*f(Z*h,P,QQ+l)oQ+p)xQ+t)|Q,O*gQ,T*rQ,h+fQ,n+nQ,o+oQ,r+uQ,v+{Q-Q,iQ-S,mR-`-ThTOr#i#l$l%}&`'k(b(d)[+Z$z!tXYZhv|}!R!S!T!X!i!k#[$O$Q$S$V$j$z%P%Z%]%d%i%k%u&[&k&l&s&z'R'U'`'a'd'f'g'r'z(Q(W(Z(l({)O)X)e)o)u)x)y)|*S*V*X*[*]*`*g*h*o*r*z+b+c+f+m+n+o+q+r+u+y+{+},P,Q,S,g,i,m,p-T-U-a-l-m-nQ#vtW%T!n!r-g-qQ%U!oQ%V!pQ%X!sQ%c-fS'j%W-kQ'l-hQ'm-iQ+v*TQ,u+wS-W,s-dR-s-uU#zu-Y-tR(i&^[gOr#l$l)[+ZX!wX#[$S$VQ#UZQ$PvR$Y|Q%_!xQ%f!{Q%l#PQ'^$xQ'y%bQ'}%gQ(V%mQ(Y%nQ*_(SQ,t+vQ-[,uR-c-ZQ$XxQ'u%`R*U'vQ-Z,sR-e-dR#OYR#TZR$}!iQ${!iV)t'a)u+mR%Q!kR%v#[Q(`%vR*n(aQ$c!RQ&h$VQ)W&zR+V)XQ#plQ$[!OQ$_!PR&p$^Q(s&oQ*v(uQ*w(vR,Z*xR$a!QXpOr)[+ZQ$h!UR&{$iQ$o!XR&|$jR)n'YQ)l'YV,j+i,k-R",
  nodeNames: "\u26A0 print Comment Script AssignStatement * BinaryExpression BitOp BitOp BitOp BitOp ArithOp ArithOp @ ArithOp ** UnaryExpression ArithOp BitOp AwaitExpression await ParenthesizedExpression ( BinaryExpression or and CompareOp in not is UnaryExpression ConditionalExpression if else LambdaExpression lambda ParamList VariableName AssignOp , : NamedExpression AssignOp YieldExpression yield from ) TupleExpression ComprehensionExpression async for LambdaExpression ArrayExpression [ ] ArrayComprehensionExpression DictionaryExpression { } DictionaryComprehensionExpression SetExpression SetComprehensionExpression CallExpression ArgList AssignOp MemberExpression . PropertyName Number String FormatString FormatReplacement FormatConversion FormatSpec ContinuedString Ellipsis None Boolean TypeDef AssignOp UpdateStatement UpdateOp ExpressionStatement DeleteStatement del PassStatement pass BreakStatement break ContinueStatement continue ReturnStatement return YieldStatement PrintStatement RaiseStatement raise ImportStatement import as ScopeStatement global nonlocal AssertStatement assert StatementGroup ; IfStatement Body elif WhileStatement while ForStatement TryStatement try except finally WithStatement with FunctionDefinition def ParamList AssignOp TypeDef ClassDefinition class DecoratedStatement Decorator At",
  maxTerm: 234,
  context: trackIndent,
  nodeProps: [
    [NodeProp.group, -14, 4, 80, 82, 83, 85, 87, 89, 91, 93, 94, 95, 97, 100, 103, "Statement Statement", -22, 6, 16, 19, 21, 37, 47, 48, 52, 55, 56, 59, 60, 61, 62, 65, 68, 69, 70, 74, 75, 76, 77, "Expression", -9, 105, 107, 110, 112, 113, 117, 119, 124, 126, "Statement"]
  ],
  skippedNodes: [0, 2],
  repeatNodeCount: 32,
  tokenData: "(#RMgR!^OX$}XY!5[Y[$}[]!5[]p$}pq!5[qr!7frs!;]st#+otu$}uv%3Tvw%5gwx%6sxy&)oyz&*uz{&+{{|&.k|}&/w}!O&0}!O!P&3d!P!Q&>j!Q!R&AY!R![&GW![!]'$S!]!^'&f!^!_''l!_!`'*[!`!a'+h!a!b$}!b!c'.T!c!d'/c!d!e'1T!e!h'/c!h!i'=R!i!t'/c!t!u'Fg!u!w'/c!w!x';a!x!}'/c!}#O'Hq#O#P'Iw#P#Q'Ji#Q#R'Ko#R#S'/c#S#T$}#T#U'/c#U#V'1T#V#Y'/c#Y#Z'=R#Z#f'/c#f#g'Fg#g#i'/c#i#j';a#j#o'/c#o#p'L{#p#q'Mq#q#r'N}#r#s( {#s$g$}$g~'/c<r%`Z%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux#O$}#O#P!3Q#P#o$}#o#p!4h#p#q$}#q#r!3f#r~$}9[&^Z%p7[%gS%m`%v!bOr'PrsLQsw'Pwx(Px#O'P#O#PNp#P#o'P#o#pKQ#p#q'P#q#rGW#r~'P9['^Z%p7[%gS%jW%m`%v!bOr'Prs&Rsw'Pwx(Px#O'P#O#PFr#P#o'P#o#pKQ#p#q'P#q#rGW#r~'P8z(WZ%p7[%jWOr(yrs)wsw(ywxAjx#O(y#O#PF^#P#o(y#o#p>v#p#q(y#q#r5T#r~(y8z)UZ%p7[%gS%jW%v!bOr(yrs)wsw(ywx(Px#O(y#O#PAU#P#o(y#o#p?p#p#q(y#q#r5T#r~(y8z*QZ%p7[%gS%v!bOr(yrs*ssw(ywx(Px#O(y#O#P@p#P#o(y#o#p?p#p#q(y#q#r5T#r~(y8z*|Z%p7[%gS%v!bOr(yrs+osw(ywx(Px#O(y#O#P4o#P#o(y#o#p?p#p#q(y#q#r5T#r~(y8r+xX%p7[%gS%v!bOw+owx,ex#O+o#O#P4Z#P#o+o#o#p3Z#p#q+o#q#r.k#r~+o8r,jX%p7[Ow+owx-Vx#O+o#O#P3u#P#o+o#o#p2i#p#q+o#q#r.k#r~+o8r-[X%p7[Ow+owx-wx#O+o#O#P.V#P#o+o#o#p0^#p#q+o#q#r.k#r~+o7[-|R%p7[O#o-w#p#q-w#r~-w8r.[T%p7[O#o+o#o#p.k#p#q+o#q#r.k#r~+o!f.rV%gS%v!bOw.kwx/Xx#O.k#O#P3T#P#o.k#o#p3Z#p~.k!f/[VOw.kwx/qx#O.k#O#P2c#P#o.k#o#p2i#p~.k!f/tUOw.kx#O.k#O#P0W#P#o.k#o#p0^#p~.k!f0ZPO~.k!f0cV%gSOw0xwx1^x#O0x#O#P2]#P#o0x#o#p.k#p~0xS0}T%gSOw0xwx1^x#O0x#O#P2]#P~0xS1aTOw0xwx1px#O0x#O#P2V#P~0xS1sSOw0xx#O0x#O#P2P#P~0xS2SPO~0xS2YPO~0xS2`PO~0x!f2fPO~.k!f2nV%gSOw0xwx1^x#O0x#O#P2]#P#o0x#o#p.k#p~0x!f3WPO~.k!f3`V%gSOw0xwx1^x#O0x#O#P2]#P#o0x#o#p.k#p~0x8r3zT%p7[O#o+o#o#p.k#p#q+o#q#r.k#r~+o8r4`T%p7[O#o+o#o#p.k#p#q+o#q#r.k#r~+o8z4tT%p7[O#o(y#o#p5T#p#q(y#q#r5T#r~(y!n5^X%gS%jW%v!bOr5Trs5ysw5Twx7ax#O5T#O#P@j#P#o5T#o#p?p#p~5T!n6QX%gS%v!bOr5Trs6msw5Twx7ax#O5T#O#P@d#P#o5T#o#p?p#p~5T!n6tX%gS%v!bOr5Trs.ksw5Twx7ax#O5T#O#P?j#P#o5T#o#p?p#p~5T!n7fX%jWOr5Trs5ysw5Twx8Rx#O5T#O#P>p#P#o5T#o#p>v#p~5T!n8WX%jWOr5Trs5ysw5Twx8sx#O5T#O#P:^#P#o5T#o#p:d#p~5TW8xT%jWOr8srs9Xs#O8s#O#P:W#P~8sW9[TOr8srs9ks#O8s#O#P:Q#P~8sW9nSOr8ss#O8s#O#P9z#P~8sW9}PO~8sW:TPO~8sW:ZPO~8s!n:aPO~5T!n:kX%gS%jWOr;Wrs;tsw;Wwx<zx#O;W#O#P>j#P#o;W#o#p5T#p~;W[;_V%gS%jWOr;Wrs;tsw;Wwx<zx#O;W#O#P>j#P~;W[;yV%gSOr;Wrs<`sw;Wwx<zx#O;W#O#P>d#P~;W[<eV%gSOr;Wrs0xsw;Wwx<zx#O;W#O#P>^#P~;W[=PV%jWOr;Wrs;tsw;Wwx=fx#O;W#O#P>W#P~;W[=kV%jWOr;Wrs;tsw;Wwx8sx#O;W#O#P>Q#P~;W[>TPO~;W[>ZPO~;W[>aPO~;W[>gPO~;W[>mPO~;W!n>sPO~5T!n>}X%gS%jWOr;Wrs;tsw;Wwx<zx#O;W#O#P>j#P#o;W#o#p5T#p~;W!n?mPO~5T!n?wX%gS%jWOr;Wrs;tsw;Wwx<zx#O;W#O#P>j#P#o;W#o#p5T#p~;W!n@gPO~5T!n@mPO~5T8z@uT%p7[O#o(y#o#p5T#p#q(y#q#r5T#r~(y8zAZT%p7[O#o(y#o#p5T#p#q(y#q#r5T#r~(y8zAqZ%p7[%jWOr(yrs)wsw(ywxBdx#O(y#O#PEx#P#o(y#o#p:d#p#q(y#q#r5T#r~(y7dBkX%p7[%jWOrBdrsCWs#OBd#O#PEd#P#oBd#o#p8s#p#qBd#q#r8s#r~Bd7dC]X%p7[OrBdrsCxs#OBd#O#PEO#P#oBd#o#p8s#p#qBd#q#r8s#r~Bd7dC}X%p7[OrBdrs-ws#OBd#O#PDj#P#oBd#o#p8s#p#qBd#q#r8s#r~Bd7dDoT%p7[O#oBd#o#p8s#p#qBd#q#r8s#r~Bd7dETT%p7[O#oBd#o#p8s#p#qBd#q#r8s#r~Bd7dEiT%p7[O#oBd#o#p8s#p#qBd#q#r8s#r~Bd8zE}T%p7[O#o(y#o#p5T#p#q(y#q#r5T#r~(y8zFcT%p7[O#o(y#o#p5T#p#q(y#q#r5T#r~(y9[FwT%p7[O#o'P#o#pGW#p#q'P#q#rGW#r~'P#OGcX%gS%jW%m`%v!bOrGWrsHOswGWwx7ax#OGW#O#PKz#P#oGW#o#pKQ#p~GW#OHXX%gS%m`%v!bOrGWrsHtswGWwx7ax#OGW#O#PKt#P#oGW#o#pKQ#p~GW#OH}X%gS%m`%v!bOrGWrsIjswGWwx7ax#OGW#O#PJz#P#oGW#o#pKQ#p~GW!vIsV%gS%m`%v!bOwIjwx/Xx#OIj#O#PJY#P#oIj#o#pJ`#p~Ij!vJ]PO~Ij!vJeV%gSOw0xwx1^x#O0x#O#P2]#P#o0x#o#pIj#p~0x#OJ}PO~GW#OKXX%gS%jWOr;Wrs;tsw;Wwx<zx#O;W#O#P>j#P#o;W#o#pGW#p~;W#OKwPO~GW#OK}PO~GW9[L]Z%p7[%gS%m`%v!bOr'PrsMOsw'Pwx(Px#O'P#O#PN[#P#o'P#o#pKQ#p#q'P#q#rGW#r~'P9SMZX%p7[%gS%m`%v!bOwMOwx,ex#OMO#O#PMv#P#oMO#o#pJ`#p#qMO#q#rIj#r~MO9SM{T%p7[O#oMO#o#pIj#p#qMO#q#rIj#r~MO9[NaT%p7[O#o'P#o#pGW#p#q'P#q#rGW#r~'P9[NuT%p7[O#o'P#o#pGW#p#q'P#q#rGW#r~'P<b! aZ%p7[%jW%sp%x#tOr!!Srs)wsw!!Swx!-Qx#O!!S#O#P!2l#P#o!!S#o#p!+d#p#q!!S#q#r!#j#r~!!S<b!!cZ%p7[%gS%jW%sp%v!b%x#tOr!!Srs)wsw!!Swx! Ux#O!!S#O#P!#U#P#o!!S#o#p!,^#p#q!!S#q#r!#j#r~!!S<b!#ZT%p7[O#o!!S#o#p!#j#p#q!!S#q#r!#j#r~!!S&U!#wX%gS%jW%sp%v!b%x#tOr!#jrs5ysw!#jwx!$dx#O!#j#O#P!,W#P#o!#j#o#p!,^#p~!#j&U!$mX%jW%sp%x#tOr!#jrs5ysw!#jwx!%Yx#O!#j#O#P!+^#P#o!#j#o#p!+d#p~!#j&U!%cX%jW%sp%x#tOr!#jrs5ysw!#jwx!&Ox#O!#j#O#P!*d#P#o!#j#o#p!*j#p~!#j$n!&XX%jW%sp%x#tOr!&trs9Xsw!&twx!&Ox#O!&t#O#P!)r#P#o!&t#o#p!)x#p~!&t$n!&}X%jW%sp%x#tOr!&trs9Xsw!&twx!'jx#O!&t#O#P!)Q#P#o!&t#o#p!)W#p~!&t$n!'sX%jW%sp%x#tOr!&trs9Xsw!&twx!&Ox#O!&t#O#P!(`#P#o!&t#o#p!(f#p~!&t$n!(cPO~!&t$n!(kV%jWOr8srs9Xs#O8s#O#P:W#P#o8s#o#p!&t#p~8s$n!)TPO~!&t$n!)]V%jWOr8srs9Xs#O8s#O#P:W#P#o8s#o#p!&t#p~8s$n!)uPO~!&t$n!)}V%jWOr8srs9Xs#O8s#O#P:W#P#o8s#o#p!&t#p~8s&U!*gPO~!#j&U!*qX%gS%jWOr;Wrs;tsw;Wwx<zx#O;W#O#P>j#P#o;W#o#p!#j#p~;W&U!+aPO~!#j&U!+kX%gS%jWOr;Wrs;tsw;Wwx<zx#O;W#O#P>j#P#o;W#o#p!#j#p~;W&U!,ZPO~!#j&U!,eX%gS%jWOr;Wrs;tsw;Wwx<zx#O;W#O#P>j#P#o;W#o#p!#j#p~;W<b!-]Z%p7[%jW%sp%x#tOr!!Srs)wsw!!Swx!.Ox#O!!S#O#P!2W#P#o!!S#o#p!*j#p#q!!S#q#r!#j#r~!!S:z!.ZZ%p7[%jW%sp%x#tOr!.|rsCWsw!.|wx!.Ox#O!.|#O#P!1r#P#o!.|#o#p!)x#p#q!.|#q#r!&t#r~!.|:z!/XZ%p7[%jW%sp%x#tOr!.|rsCWsw!.|wx!/zx#O!.|#O#P!1^#P#o!.|#o#p!)W#p#q!.|#q#r!&t#r~!.|:z!0VZ%p7[%jW%sp%x#tOr!.|rsCWsw!.|wx!.Ox#O!.|#O#P!0x#P#o!.|#o#p!(f#p#q!.|#q#r!&t#r~!.|:z!0}T%p7[O#o!.|#o#p!&t#p#q!.|#q#r!&t#r~!.|:z!1cT%p7[O#o!.|#o#p!&t#p#q!.|#q#r!&t#r~!.|:z!1wT%p7[O#o!.|#o#p!&t#p#q!.|#q#r!&t#r~!.|<b!2]T%p7[O#o!!S#o#p!#j#p#q!!S#q#r!#j#r~!!S<b!2qT%p7[O#o!!S#o#p!#j#p#q!!S#q#r!#j#r~!!S<r!3VT%p7[O#o$}#o#p!3f#p#q$}#q#r!3f#r~$}&f!3uX%gS%jW%m`%sp%v!b%x#tOr!3frsHOsw!3fwx!$dx#O!3f#O#P!4b#P#o!3f#o#p!4h#p~!3f&f!4ePO~!3f&f!4oX%gS%jWOr;Wrs;tsw;Wwx<zx#O;W#O#P>j#P#o;W#o#p!3f#p~;WMg!5oa%p7[%gS%jW$o1s%m`%sp%v!b%x#tOX$}XY!5[Y[$}[]!5[]p$}pq!5[qr$}rs&Rsw$}wx! Ux#O$}#O#P!6t#P#o$}#o#p!4h#p#q$}#q#r!3f#r~$}Mg!6yX%p7[OY$}YZ!5[Z]$}]^!5[^#o$}#o#p!3f#p#q$}#q#r!3f#r~$}<u!7wb%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux!_$}!_!`!9P!`#O$}#O#P!3Q#P#T$}#T#U!:V#U#f$}#f#g!:V#g#h!:V#h#o$}#o#p!4h#p#q$}#q#r!3f#r~$}<u!9dZjR%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux#O$}#O#P!3Q#P#o$}#o#p!4h#p#q$}#q#r!3f#r~$}<u!:jZ!jR%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux#O$}#O#P!3Q#P#o$}#o#p!4h#p#q$}#q#r!3f#r~$}G{!;l_%tp%p7[%gS%e,X%m`%v!bOY!<kYZ'PZ]!<k]^'P^r!<krs#(ysw!<kwx!>yx#O!<k#O#P#+Z#P#o!<k#o#p#'w#p#q!<k#q#r#%s#r~!<kDe!<z_%p7[%gS%jW%e,X%m`%v!bOY!<kYZ'PZ]!<k]^'P^r!<krs!=ysw!<kwx!>yx#O!<k#O#P#%_#P#o!<k#o#p#'w#p#q!<k#q#r#%s#r~!<kDe!>WZ%p7[%gS%e,X%m`%v!bOr'PrsLQsw'Pwx(Px#O'P#O#PNp#P#o'P#o#pKQ#p#q'P#q#rGW#r~'PDT!?S_%p7[%jW%e,XOY!@RYZ(yZ]!@R]^(y^r!@Rrs!A_sw!@Rwx# Rx#O!@R#O#P#$y#P#o!@R#o#p!Lw#p#q!@R#q#r!Bq#r~!@RDT!@`_%p7[%gS%jW%e,X%v!bOY!@RYZ(yZ]!@R]^(y^r!@Rrs!A_sw!@Rwx!>yx#O!@R#O#P!B]#P#o!@R#o#p!NP#p#q!@R#q#r!Bq#r~!@RDT!AjZ%p7[%gS%e,X%v!bOr(yrs*ssw(ywx(Px#O(y#O#P@p#P#o(y#o#p?p#p#q(y#q#r5T#r~(yDT!BbT%p7[O#o!@R#o#p!Bq#p#q!@R#q#r!Bq#r~!@R-w!B|]%gS%jW%e,X%v!bOY!BqYZ5TZ]!Bq]^5T^r!Bqrs!Cusw!Bqwx!Dkx#O!Bq#O#P!My#P#o!Bq#o#p!NP#p~!Bq-w!DOX%gS%e,X%v!bOr5Trs6msw5Twx7ax#O5T#O#P@d#P#o5T#o#p?p#p~5T-w!Dr]%jW%e,XOY!BqYZ5TZ]!Bq]^5T^r!Bqrs!Cusw!Bqwx!Ekx#O!Bq#O#P!Lq#P#o!Bq#o#p!Lw#p~!Bq-w!Er]%jW%e,XOY!BqYZ5TZ]!Bq]^5T^r!Bqrs!Cusw!Bqwx!Fkx#O!Bq#O#P!Gy#P#o!Bq#o#p!HP#p~!Bq,a!FrX%jW%e,XOY!FkYZ8sZ]!Fk]^8s^r!Fkrs!G_s#O!Fk#O#P!Gs#P~!Fk,a!GdT%e,XOr8srs9ks#O8s#O#P:Q#P~8s,a!GvPO~!Fk-w!G|PO~!Bq-w!HY]%gS%jW%e,XOY!IRYZ;WZ]!IR]^;W^r!IRrs!I}sw!IRwx!Jkx#O!IR#O#P!Lk#P#o!IR#o#p!Bq#p~!IR,e!I[Z%gS%jW%e,XOY!IRYZ;WZ]!IR]^;W^r!IRrs!I}sw!IRwx!Jkx#O!IR#O#P!Lk#P~!IR,e!JUV%gS%e,XOr;Wrs<`sw;Wwx<zx#O;W#O#P>d#P~;W,e!JrZ%jW%e,XOY!IRYZ;WZ]!IR]^;W^r!IRrs!I}sw!IRwx!Kex#O!IR#O#P!Le#P~!IR,e!KlZ%jW%e,XOY!IRYZ;WZ]!IR]^;W^r!IRrs!I}sw!IRwx!Fkx#O!IR#O#P!L_#P~!IR,e!LbPO~!IR,e!LhPO~!IR,e!LnPO~!IR-w!LtPO~!Bq-w!MQ]%gS%jW%e,XOY!IRYZ;WZ]!IR]^;W^r!IRrs!I}sw!IRwx!Jkx#O!IR#O#P!Lk#P#o!IR#o#p!Bq#p~!IR-w!M|PO~!Bq-w!NY]%gS%jW%e,XOY!IRYZ;WZ]!IR]^;W^r!IRrs!I}sw!IRwx!Jkx#O!IR#O#P!Lk#P#o!IR#o#p!Bq#p~!IRDT# [_%p7[%jW%e,XOY!@RYZ(yZ]!@R]^(y^r!@Rrs!A_sw!@Rwx#!Zx#O!@R#O#P#$e#P#o!@R#o#p!HP#p#q!@R#q#r!Bq#r~!@RBm#!d]%p7[%jW%e,XOY#!ZYZBdZ]#!Z]^Bd^r#!Zrs##]s#O#!Z#O#P#$P#P#o#!Z#o#p!Fk#p#q#!Z#q#r!Fk#r~#!ZBm##dX%p7[%e,XOrBdrsCxs#OBd#O#PEO#P#oBd#o#p8s#p#qBd#q#r8s#r~BdBm#$UT%p7[O#o#!Z#o#p!Fk#p#q#!Z#q#r!Fk#r~#!ZDT#$jT%p7[O#o!@R#o#p!Bq#p#q!@R#q#r!Bq#r~!@RDT#%OT%p7[O#o!@R#o#p!Bq#p#q!@R#q#r!Bq#r~!@RDe#%dT%p7[O#o!<k#o#p#%s#p#q!<k#q#r#%s#r~!<k.X#&Q]%gS%jW%e,X%m`%v!bOY#%sYZGWZ]#%s]^GW^r#%srs#&ysw#%swx!Dkx#O#%s#O#P#'q#P#o#%s#o#p#'w#p~#%s.X#'UX%gS%e,X%m`%v!bOrGWrsHtswGWwx7ax#OGW#O#PKt#P#oGW#o#pKQ#p~GW.X#'tPO~#%s.X#(Q]%gS%jW%e,XOY!IRYZ;WZ]!IR]^;W^r!IRrs!I}sw!IRwx!Jkx#O!IR#O#P!Lk#P#o!IR#o#p#%s#p~!IRGZ#)WZ%p7[%gS%e,X%m`%v!bOr'Prs#)ysw'Pwx(Px#O'P#O#P#*u#P#o'P#o#pKQ#p#q'P#q#rGW#r~'PGZ#*YX%k#|%p7[%gS%i,X%m`%v!bOwMOwx,ex#OMO#O#PMv#P#oMO#o#pJ`#p#qMO#q#rIj#r~MO9[#*zT%p7[O#o'P#o#pGW#p#q'P#q#rGW#r~'PDe#+`T%p7[O#o!<k#o#p#%s#p#q!<k#q#r#%s#r~!<kMg#,S_Q1s%p7[%gS%jW%m`%sp%v!b%x#tOY#+oYZ$}Z]#+o]^$}^r#+ors#-Rsw#+owx$Bmx#O#+o#O#P%/o#P#o#+o#o#p%2R#p#q#+o#q#r%0c#r~#+oJP#-`_Q1s%p7[%gS%m`%v!bOY#._YZ'PZ]#._]^'P^r#._rs$>Psw#._wx#/mx#O#._#O#P$Ay#P#o#._#o#p$<T#p#q#._#q#r$6T#r~#._JP#.n_Q1s%p7[%gS%jW%m`%v!bOY#._YZ'PZ]#._]^'P^r#._rs#-Rsw#._wx#/mx#O#._#O#P$5a#P#o#._#o#p$<T#p#q#._#q#r$6T#r~#._Io#/v_Q1s%p7[%jWOY#0uYZ(yZ]#0u]^(y^r#0urs#2Rsw#0uwx$-ex#O#0u#O#P$4m#P#o#0u#o#p$(k#p#q#0u#q#r#Eg#r~#0uIo#1S_Q1s%p7[%gS%jW%v!bOY#0uYZ(yZ]#0u]^(y^r#0urs#2Rsw#0uwx#/mx#O#0u#O#P$,q#P#o#0u#o#p$*R#p#q#0u#q#r#Eg#r~#0uIo#2^_Q1s%p7[%gS%v!bOY#0uYZ(yZ]#0u]^(y^r#0urs#3]sw#0uwx#/mx#O#0u#O#P$+}#P#o#0u#o#p$*R#p#q#0u#q#r#Eg#r~#0uIo#3h_Q1s%p7[%gS%v!bOY#0uYZ(yZ]#0u]^(y^r#0urs#4gsw#0uwx#/mx#O#0u#O#P#Ds#P#o#0u#o#p$*R#p#q#0u#q#r#Eg#r~#0uIg#4r]Q1s%p7[%gS%v!bOY#4gYZ+oZ]#4g]^+o^w#4gwx#5kx#O#4g#O#P#DP#P#o#4g#o#p#Bc#p#q#4g#q#r#9a#r~#4gIg#5r]Q1s%p7[OY#4gYZ+oZ]#4g]^+o^w#4gwx#6kx#O#4g#O#P#C]#P#o#4g#o#p#AT#p#q#4g#q#r#9a#r~#4gIg#6r]Q1s%p7[OY#4gYZ+oZ]#4g]^+o^w#4gwx#7kx#O#4g#O#P#8m#P#o#4g#o#p#<a#p#q#4g#q#r#9a#r~#4gHP#7rXQ1s%p7[OY#7kYZ-wZ]#7k]^-w^#o#7k#o#p#8_#p#q#7k#q#r#8_#r~#7k1s#8dRQ1sOY#8_Z]#8_^~#8_Ig#8tXQ1s%p7[OY#4gYZ+oZ]#4g]^+o^#o#4g#o#p#9a#p#q#4g#q#r#9a#r~#4g3Z#9jZQ1s%gS%v!bOY#9aYZ.kZ]#9a]^.k^w#9awx#:]x#O#9a#O#P#A}#P#o#9a#o#p#Bc#p~#9a3Z#:bZQ1sOY#9aYZ.kZ]#9a]^.k^w#9awx#;Tx#O#9a#O#P#@o#P#o#9a#o#p#AT#p~#9a3Z#;YZQ1sOY#9aYZ.kZ]#9a]^.k^w#9awx#8_x#O#9a#O#P#;{#P#o#9a#o#p#<a#p~#9a3Z#<QTQ1sOY#9aYZ.kZ]#9a]^.k^~#9a3Z#<hZQ1s%gSOY#=ZYZ0xZ]#=Z]^0x^w#=Zwx#=}x#O#=Z#O#P#@Z#P#o#=Z#o#p#9a#p~#=Z1w#=bXQ1s%gSOY#=ZYZ0xZ]#=Z]^0x^w#=Zwx#=}x#O#=Z#O#P#@Z#P~#=Z1w#>SXQ1sOY#=ZYZ0xZ]#=Z]^0x^w#=Zwx#>ox#O#=Z#O#P#?u#P~#=Z1w#>tXQ1sOY#=ZYZ0xZ]#=Z]^0x^w#=Zwx#8_x#O#=Z#O#P#?a#P~#=Z1w#?fTQ1sOY#=ZYZ0xZ]#=Z]^0x^~#=Z1w#?zTQ1sOY#=ZYZ0xZ]#=Z]^0x^~#=Z1w#@`TQ1sOY#=ZYZ0xZ]#=Z]^0x^~#=Z3Z#@tTQ1sOY#9aYZ.kZ]#9a]^.k^~#9a3Z#A[ZQ1s%gSOY#=ZYZ0xZ]#=Z]^0x^w#=Zwx#=}x#O#=Z#O#P#@Z#P#o#=Z#o#p#9a#p~#=Z3Z#BSTQ1sOY#9aYZ.kZ]#9a]^.k^~#9a3Z#BjZQ1s%gSOY#=ZYZ0xZ]#=Z]^0x^w#=Zwx#=}x#O#=Z#O#P#@Z#P#o#=Z#o#p#9a#p~#=ZIg#CdXQ1s%p7[OY#4gYZ+oZ]#4g]^+o^#o#4g#o#p#9a#p#q#4g#q#r#9a#r~#4gIg#DWXQ1s%p7[OY#4gYZ+oZ]#4g]^+o^#o#4g#o#p#9a#p#q#4g#q#r#9a#r~#4gIo#DzXQ1s%p7[OY#0uYZ(yZ]#0u]^(y^#o#0u#o#p#Eg#p#q#0u#q#r#Eg#r~#0u3c#Er]Q1s%gS%jW%v!bOY#EgYZ5TZ]#Eg]^5T^r#Egrs#Fksw#Egwx#Hox#O#Eg#O#P$+i#P#o#Eg#o#p$*R#p~#Eg3c#Ft]Q1s%gS%v!bOY#EgYZ5TZ]#Eg]^5T^r#Egrs#Gmsw#Egwx#Hox#O#Eg#O#P$+T#P#o#Eg#o#p$*R#p~#Eg3c#Gv]Q1s%gS%v!bOY#EgYZ5TZ]#Eg]^5T^r#Egrs#9asw#Egwx#Hox#O#Eg#O#P$)m#P#o#Eg#o#p$*R#p~#Eg3c#Hv]Q1s%jWOY#EgYZ5TZ]#Eg]^5T^r#Egrs#Fksw#Egwx#Iox#O#Eg#O#P$(V#P#o#Eg#o#p$(k#p~#Eg3c#Iv]Q1s%jWOY#EgYZ5TZ]#Eg]^5T^r#Egrs#Fksw#Egwx#Jox#O#Eg#O#P#NT#P#o#Eg#o#p#Ni#p~#Eg1{#JvXQ1s%jWOY#JoYZ8sZ]#Jo]^8s^r#Jors#Kcs#O#Jo#O#P#Mo#P~#Jo1{#KhXQ1sOY#JoYZ8sZ]#Jo]^8s^r#Jors#LTs#O#Jo#O#P#MZ#P~#Jo1{#LYXQ1sOY#JoYZ8sZ]#Jo]^8s^r#Jors#8_s#O#Jo#O#P#Lu#P~#Jo1{#LzTQ1sOY#JoYZ8sZ]#Jo]^8s^~#Jo1{#M`TQ1sOY#JoYZ8sZ]#Jo]^8s^~#Jo1{#MtTQ1sOY#JoYZ8sZ]#Jo]^8s^~#Jo3c#NYTQ1sOY#EgYZ5TZ]#Eg]^5T^~#Eg3c#Nr]Q1s%gS%jWOY$ kYZ;WZ]$ k]^;W^r$ krs$!gsw$ kwx$$Zx#O$ k#O#P$'q#P#o$ k#o#p#Eg#p~$ k2P$ tZQ1s%gS%jWOY$ kYZ;WZ]$ k]^;W^r$ krs$!gsw$ kwx$$Zx#O$ k#O#P$'q#P~$ k2P$!nZQ1s%gSOY$ kYZ;WZ]$ k]^;W^r$ krs$#asw$ kwx$$Zx#O$ k#O#P$']#P~$ k2P$#hZQ1s%gSOY$ kYZ;WZ]$ k]^;W^r$ krs#=Zsw$ kwx$$Zx#O$ k#O#P$&w#P~$ k2P$$bZQ1s%jWOY$ kYZ;WZ]$ k]^;W^r$ krs$!gsw$ kwx$%Tx#O$ k#O#P$&c#P~$ k2P$%[ZQ1s%jWOY$ kYZ;WZ]$ k]^;W^r$ krs$!gsw$ kwx#Jox#O$ k#O#P$%}#P~$ k2P$&STQ1sOY$ kYZ;WZ]$ k]^;W^~$ k2P$&hTQ1sOY$ kYZ;WZ]$ k]^;W^~$ k2P$&|TQ1sOY$ kYZ;WZ]$ k]^;W^~$ k2P$'bTQ1sOY$ kYZ;WZ]$ k]^;W^~$ k2P$'vTQ1sOY$ kYZ;WZ]$ k]^;W^~$ k3c$([TQ1sOY#EgYZ5TZ]#Eg]^5T^~#Eg3c$(t]Q1s%gS%jWOY$ kYZ;WZ]$ k]^;W^r$ krs$!gsw$ kwx$$Zx#O$ k#O#P$'q#P#o$ k#o#p#Eg#p~$ k3c$)rTQ1sOY#EgYZ5TZ]#Eg]^5T^~#Eg3c$*[]Q1s%gS%jWOY$ kYZ;WZ]$ k]^;W^r$ krs$!gsw$ kwx$$Zx#O$ k#O#P$'q#P#o$ k#o#p#Eg#p~$ k3c$+YTQ1sOY#EgYZ5TZ]#Eg]^5T^~#Eg3c$+nTQ1sOY#EgYZ5TZ]#Eg]^5T^~#EgIo$,UXQ1s%p7[OY#0uYZ(yZ]#0u]^(y^#o#0u#o#p#Eg#p#q#0u#q#r#Eg#r~#0uIo$,xXQ1s%p7[OY#0uYZ(yZ]#0u]^(y^#o#0u#o#p#Eg#p#q#0u#q#r#Eg#r~#0uIo$-n_Q1s%p7[%jWOY#0uYZ(yZ]#0u]^(y^r#0urs#2Rsw#0uwx$.mx#O#0u#O#P$3y#P#o#0u#o#p#Ni#p#q#0u#q#r#Eg#r~#0uHX$.v]Q1s%p7[%jWOY$.mYZBdZ]$.m]^Bd^r$.mrs$/os#O$.m#O#P$3V#P#o$.m#o#p#Jo#p#q$.m#q#r#Jo#r~$.mHX$/v]Q1s%p7[OY$.mYZBdZ]$.m]^Bd^r$.mrs$0os#O$.m#O#P$2c#P#o$.m#o#p#Jo#p#q$.m#q#r#Jo#r~$.mHX$0v]Q1s%p7[OY$.mYZBdZ]$.m]^Bd^r$.mrs#7ks#O$.m#O#P$1o#P#o$.m#o#p#Jo#p#q$.m#q#r#Jo#r~$.mHX$1vXQ1s%p7[OY$.mYZBdZ]$.m]^Bd^#o$.m#o#p#Jo#p#q$.m#q#r#Jo#r~$.mHX$2jXQ1s%p7[OY$.mYZBdZ]$.m]^Bd^#o$.m#o#p#Jo#p#q$.m#q#r#Jo#r~$.mHX$3^XQ1s%p7[OY$.mYZBdZ]$.m]^Bd^#o$.m#o#p#Jo#p#q$.m#q#r#Jo#r~$.mIo$4QXQ1s%p7[OY#0uYZ(yZ]#0u]^(y^#o#0u#o#p#Eg#p#q#0u#q#r#Eg#r~#0uIo$4tXQ1s%p7[OY#0uYZ(yZ]#0u]^(y^#o#0u#o#p#Eg#p#q#0u#q#r#Eg#r~#0uJP$5hXQ1s%p7[OY#._YZ'PZ]#._]^'P^#o#._#o#p$6T#p#q#._#q#r$6T#r~#._3s$6b]Q1s%gS%jW%m`%v!bOY$6TYZGWZ]$6T]^GW^r$6Trs$7Zsw$6Twx#Hox#O$6T#O#P$=k#P#o$6T#o#p$<T#p~$6T3s$7f]Q1s%gS%m`%v!bOY$6TYZGWZ]$6T]^GW^r$6Trs$8_sw$6Twx#Hox#O$6T#O#P$=V#P#o$6T#o#p$<T#p~$6T3s$8j]Q1s%gS%m`%v!bOY$6TYZGWZ]$6T]^GW^r$6Trs$9csw$6Twx#Hox#O$6T#O#P$;o#P#o$6T#o#p$<T#p~$6T3k$9nZQ1s%gS%m`%v!bOY$9cYZIjZ]$9c]^Ij^w$9cwx#:]x#O$9c#O#P$:a#P#o$9c#o#p$:u#p~$9c3k$:fTQ1sOY$9cYZIjZ]$9c]^Ij^~$9c3k$:|ZQ1s%gSOY#=ZYZ0xZ]#=Z]^0x^w#=Zwx#=}x#O#=Z#O#P#@Z#P#o#=Z#o#p$9c#p~#=Z3s$;tTQ1sOY$6TYZGWZ]$6T]^GW^~$6T3s$<^]Q1s%gS%jWOY$ kYZ;WZ]$ k]^;W^r$ krs$!gsw$ kwx$$Zx#O$ k#O#P$'q#P#o$ k#o#p$6T#p~$ k3s$=[TQ1sOY$6TYZGWZ]$6T]^GW^~$6T3s$=pTQ1sOY$6TYZGWZ]$6T]^GW^~$6TJP$>^_Q1s%p7[%gS%m`%v!bOY#._YZ'PZ]#._]^'P^r#._rs$?]sw#._wx#/mx#O#._#O#P$AV#P#o#._#o#p$<T#p#q#._#q#r$6T#r~#._Iw$?j]Q1s%p7[%gS%m`%v!bOY$?]YZMOZ]$?]]^MO^w$?]wx#5kx#O$?]#O#P$@c#P#o$?]#o#p$:u#p#q$?]#q#r$9c#r~$?]Iw$@jXQ1s%p7[OY$?]YZMOZ]$?]]^MO^#o$?]#o#p$9c#p#q$?]#q#r$9c#r~$?]JP$A^XQ1s%p7[OY#._YZ'PZ]#._]^'P^#o#._#o#p$6T#p#q#._#q#r$6T#r~#._JP$BQXQ1s%p7[OY#._YZ'PZ]#._]^'P^#o#._#o#p$6T#p#q#._#q#r$6T#r~#._MV$Bz_Q1s%p7[%jW%sp%x#tOY$CyYZ!!SZ]$Cy]^!!S^r$Cyrs#2Rsw$Cywx%&{x#O$Cy#O#P%.{#P#o$Cy#o#p%$c#p#q$Cy#q#r$E}#r~$CyMV$D[_Q1s%p7[%gS%jW%sp%v!b%x#tOY$CyYZ!!SZ]$Cy]^!!S^r$Cyrs#2Rsw$Cywx$Bmx#O$Cy#O#P$EZ#P#o$Cy#o#p%%y#p#q$Cy#q#r$E}#r~$CyMV$EbXQ1s%p7[OY$CyYZ!!SZ]$Cy]^!!S^#o$Cy#o#p$E}#p#q$Cy#q#r$E}#r~$Cy6y$F^]Q1s%gS%jW%sp%v!b%x#tOY$E}YZ!#jZ]$E}]^!#j^r$E}rs#Fksw$E}wx$GVx#O$E}#O#P%%e#P#o$E}#o#p%%y#p~$E}6y$Gb]Q1s%jW%sp%x#tOY$E}YZ!#jZ]$E}]^!#j^r$E}rs#Fksw$E}wx$HZx#O$E}#O#P%#}#P#o$E}#o#p%$c#p~$E}6y$Hf]Q1s%jW%sp%x#tOY$E}YZ!#jZ]$E}]^!#j^r$E}rs#Fksw$E}wx$I_x#O$E}#O#P%!g#P#o$E}#o#p%!{#p~$E}5c$Ij]Q1s%jW%sp%x#tOY$JcYZ!&tZ]$Jc]^!&t^r$Jcrs#Kcsw$Jcwx$I_x#O$Jc#O#P% X#P#o$Jc#o#p% m#p~$Jc5c$Jn]Q1s%jW%sp%x#tOY$JcYZ!&tZ]$Jc]^!&t^r$Jcrs#Kcsw$Jcwx$Kgx#O$Jc#O#P$My#P#o$Jc#o#p$N_#p~$Jc5c$Kr]Q1s%jW%sp%x#tOY$JcYZ!&tZ]$Jc]^!&t^r$Jcrs#Kcsw$Jcwx$I_x#O$Jc#O#P$Lk#P#o$Jc#o#p$MP#p~$Jc5c$LpTQ1sOY$JcYZ!&tZ]$Jc]^!&t^~$Jc5c$MWZQ1s%jWOY#JoYZ8sZ]#Jo]^8s^r#Jors#Kcs#O#Jo#O#P#Mo#P#o#Jo#o#p$Jc#p~#Jo5c$NOTQ1sOY$JcYZ!&tZ]$Jc]^!&t^~$Jc5c$NfZQ1s%jWOY#JoYZ8sZ]#Jo]^8s^r#Jors#Kcs#O#Jo#O#P#Mo#P#o#Jo#o#p$Jc#p~#Jo5c% ^TQ1sOY$JcYZ!&tZ]$Jc]^!&t^~$Jc5c% tZQ1s%jWOY#JoYZ8sZ]#Jo]^8s^r#Jors#Kcs#O#Jo#O#P#Mo#P#o#Jo#o#p$Jc#p~#Jo6y%!lTQ1sOY$E}YZ!#jZ]$E}]^!#j^~$E}6y%#U]Q1s%gS%jWOY$ kYZ;WZ]$ k]^;W^r$ krs$!gsw$ kwx$$Zx#O$ k#O#P$'q#P#o$ k#o#p$E}#p~$ k6y%$STQ1sOY$E}YZ!#jZ]$E}]^!#j^~$E}6y%$l]Q1s%gS%jWOY$ kYZ;WZ]$ k]^;W^r$ krs$!gsw$ kwx$$Zx#O$ k#O#P$'q#P#o$ k#o#p$E}#p~$ k6y%%jTQ1sOY$E}YZ!#jZ]$E}]^!#j^~$E}6y%&S]Q1s%gS%jWOY$ kYZ;WZ]$ k]^;W^r$ krs$!gsw$ kwx$$Zx#O$ k#O#P$'q#P#o$ k#o#p$E}#p~$ kMV%'Y_Q1s%p7[%jW%sp%x#tOY$CyYZ!!SZ]$Cy]^!!S^r$Cyrs#2Rsw$Cywx%(Xx#O$Cy#O#P%.X#P#o$Cy#o#p%!{#p#q$Cy#q#r$E}#r~$CyKo%(f_Q1s%p7[%jW%sp%x#tOY%)eYZ!.|Z]%)e]^!.|^r%)ers$/osw%)ewx%(Xx#O%)e#O#P%-e#P#o%)e#o#p% m#p#q%)e#q#r$Jc#r~%)eKo%)r_Q1s%p7[%jW%sp%x#tOY%)eYZ!.|Z]%)e]^!.|^r%)ers$/osw%)ewx%*qx#O%)e#O#P%,q#P#o%)e#o#p$N_#p#q%)e#q#r$Jc#r~%)eKo%+O_Q1s%p7[%jW%sp%x#tOY%)eYZ!.|Z]%)e]^!.|^r%)ers$/osw%)ewx%(Xx#O%)e#O#P%+}#P#o%)e#o#p$MP#p#q%)e#q#r$Jc#r~%)eKo%,UXQ1s%p7[OY%)eYZ!.|Z]%)e]^!.|^#o%)e#o#p$Jc#p#q%)e#q#r$Jc#r~%)eKo%,xXQ1s%p7[OY%)eYZ!.|Z]%)e]^!.|^#o%)e#o#p$Jc#p#q%)e#q#r$Jc#r~%)eKo%-lXQ1s%p7[OY%)eYZ!.|Z]%)e]^!.|^#o%)e#o#p$Jc#p#q%)e#q#r$Jc#r~%)eMV%.`XQ1s%p7[OY$CyYZ!!SZ]$Cy]^!!S^#o$Cy#o#p$E}#p#q$Cy#q#r$E}#r~$CyMV%/SXQ1s%p7[OY$CyYZ!!SZ]$Cy]^!!S^#o$Cy#o#p$E}#p#q$Cy#q#r$E}#r~$CyMg%/vXQ1s%p7[OY#+oYZ$}Z]#+o]^$}^#o#+o#o#p%0c#p#q#+o#q#r%0c#r~#+o7Z%0t]Q1s%gS%jW%m`%sp%v!b%x#tOY%0cYZ!3fZ]%0c]^!3f^r%0crs$7Zsw%0cwx$GVx#O%0c#O#P%1m#P#o%0c#o#p%2R#p~%0c7Z%1rTQ1sOY%0cYZ!3fZ]%0c]^!3f^~%0c7Z%2[]Q1s%gS%jWOY$ kYZ;WZ]$ k]^;W^r$ krs$!gsw$ kwx$$Zx#O$ k#O#P$'q#P#o$ k#o#p%0c#p~$ kGz%3h]$}Q%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux!_$}!_!`%4a!`#O$}#O#P!3Q#P#o$}#o#p!4h#p#q$}#q#r!3f#r~$}Gz%4tZ!s,W%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux#O$}#O#P!3Q#P#o$}#o#p!4h#p#q$}#q#r!3f#r~$}Gz%5z]$wQ%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux!_$}!_!`%4a!`#O$}#O#P!3Q#P#o$}#o#p!4h#p#q$}#q#r!3f#r~$}G{%7S_%q`%p7[%jW%e,X%sp%x#tOY%8RYZ!!SZ]%8R]^!!S^r%8Rrs%9csw%8Rwx&$}x#O%8R#O#P&(X#P#o%8R#o#p&(m#p#q%8R#q#r& u#r~%8RGk%8d_%p7[%gS%jW%e,X%sp%v!b%x#tOY%8RYZ!!SZ]%8R]^!!S^r%8Rrs%9csw%8Rwx%Nax#O%8R#O#P& a#P#o%8R#o#p&#{#p#q%8R#q#r& u#r~%8RDT%9n_%p7[%gS%e,X%v!bOY%:mYZ(yZ]%:m]^(y^r%:mrs%JPsw%:mwx%;yx#O%:m#O#P%M{#P#o%:m#o#p%ER#p#q%:m#q#r%=Z#r~%:mDT%:z_%p7[%gS%jW%e,X%v!bOY%:mYZ(yZ]%:m]^(y^r%:mrs%9csw%:mwx%;yx#O%:m#O#P%<u#P#o%:m#o#p%ER#p#q%:m#q#r%=Z#r~%:mDT%<SZ%p7[%jW%e,XOr(yrs)wsw(ywxAjx#O(y#O#PF^#P#o(y#o#p>v#p#q(y#q#r5T#r~(yDT%<zT%p7[O#o%:m#o#p%=Z#p#q%:m#q#r%=Z#r~%:m-w%=f]%gS%jW%e,X%v!bOY%=ZYZ5TZ]%=Z]^5T^r%=Zrs%>_sw%=Zwx%DXx#O%=Z#O#P%Iy#P#o%=Z#o#p%ER#p~%=Z-w%>h]%gS%e,X%v!bOY%=ZYZ5TZ]%=Z]^5T^r%=Zrs%?asw%=Zwx%DXx#O%=Z#O#P%Is#P#o%=Z#o#p%ER#p~%=Z-w%?j]%gS%e,X%v!bOY%=ZYZ5TZ]%=Z]^5T^r%=Zrs%@csw%=Zwx%DXx#O%=Z#O#P%D{#P#o%=Z#o#p%ER#p~%=Z-o%@lZ%gS%e,X%v!bOY%@cYZ.kZ]%@c]^.k^w%@cwx%A_x#O%@c#O#P%Ay#P#o%@c#o#p%BP#p~%@c-o%AdV%e,XOw.kwx/qx#O.k#O#P2c#P#o.k#o#p2i#p~.k-o%A|PO~%@c-o%BWZ%gS%e,XOY%ByYZ0xZ]%By]^0x^w%Bywx%Cmx#O%By#O#P%DR#P#o%By#o#p%@c#p~%By,]%CQX%gS%e,XOY%ByYZ0xZ]%By]^0x^w%Bywx%Cmx#O%By#O#P%DR#P~%By,]%CrT%e,XOw0xwx1px#O0x#O#P2V#P~0x,]%DUPO~%By-w%D`X%jW%e,XOr5Trs5ysw5Twx8Rx#O5T#O#P>p#P#o5T#o#p>v#p~5T-w%EOPO~%=Z-w%E[]%gS%jW%e,XOY%FTYZ;WZ]%FT]^;W^r%FTrs%GPsw%FTwx%Hsx#O%FT#O#P%Im#P#o%FT#o#p%=Z#p~%FT,e%F^Z%gS%jW%e,XOY%FTYZ;WZ]%FT]^;W^r%FTrs%GPsw%FTwx%Hsx#O%FT#O#P%Im#P~%FT,e%GWZ%gS%e,XOY%FTYZ;WZ]%FT]^;W^r%FTrs%Gysw%FTwx%Hsx#O%FT#O#P%Ig#P~%FT,e%HQZ%gS%e,XOY%FTYZ;WZ]%FT]^;W^r%FTrs%Bysw%FTwx%Hsx#O%FT#O#P%Ia#P~%FT,e%HzV%jW%e,XOr;Wrs;tsw;Wwx=fx#O;W#O#P>W#P~;W,e%IdPO~%FT,e%IjPO~%FT,e%IpPO~%FT-w%IvPO~%=Z-w%I|PO~%=ZDT%J[_%p7[%gS%e,X%v!bOY%:mYZ(yZ]%:m]^(y^r%:mrs%KZsw%:mwx%;yx#O%:m#O#P%Mg#P#o%:m#o#p%ER#p#q%:m#q#r%=Z#r~%:mC{%Kf]%p7[%gS%e,X%v!bOY%KZYZ+oZ]%KZ]^+o^w%KZwx%L_x#O%KZ#O#P%MR#P#o%KZ#o#p%BP#p#q%KZ#q#r%@c#r~%KZC{%LfX%p7[%e,XOw+owx-Vx#O+o#O#P3u#P#o+o#o#p2i#p#q+o#q#r.k#r~+oC{%MWT%p7[O#o%KZ#o#p%@c#p#q%KZ#q#r%@c#r~%KZDT%MlT%p7[O#o%:m#o#p%=Z#p#q%:m#q#r%=Z#r~%:mDT%NQT%p7[O#o%:m#o#p%=Z#p#q%:m#q#r%=Z#r~%:mGk%NnZ%p7[%jW%e,X%sp%x#tOr!!Srs)wsw!!Swx!-Qx#O!!S#O#P!2l#P#o!!S#o#p!+d#p#q!!S#q#r!#j#r~!!SGk& fT%p7[O#o%8R#o#p& u#p#q%8R#q#r& u#r~%8R1_&!U]%gS%jW%e,X%sp%v!b%x#tOY& uYZ!#jZ]& u]^!#j^r& urs%>_sw& uwx&!}x#O& u#O#P&#u#P#o& u#o#p&#{#p~& u1_&#YX%jW%e,X%sp%x#tOr!#jrs5ysw!#jwx!%Yx#O!#j#O#P!+^#P#o!#j#o#p!+d#p~!#j1_&#xPO~& u1_&$U]%gS%jW%e,XOY%FTYZ;WZ]%FT]^;W^r%FTrs%GPsw%FTwx%Hsx#O%FT#O#P%Im#P#o%FT#o#p& u#p~%FTGk&%[Z%p7[%jW%e,X%sp%x#tOr!!Srs)wsw!!Swx&%}x#O!!S#O#P&'P#P#o!!S#o#p&'e#p#q!!S#q#r!#j#r~!!SGk&&^Z%h!f%p7[%jW%f,X%sp%x#tOr!.|rsCWsw!.|wx!.Ox#O!.|#O#P!1r#P#o!.|#o#p!)x#p#q!.|#q#r!&t#r~!.|<b&'UT%p7[O#o!!S#o#p!#j#p#q!!S#q#r!#j#r~!!S&U&'lX%gS%jWOr;Wrs;tsw;Wwx<zx#O;W#O#P>j#P#o;W#o#p!#j#p~;WGk&(^T%p7[O#o%8R#o#p& u#p#q%8R#q#r& u#r~%8R1_&(v]%gS%jW%e,XOY%FTYZ;WZ]%FT]^;W^r%FTrs%GPsw%FTwx%Hsx#O%FT#O#P%Im#P#o%FT#o#p& u#p~%FTG{&*SZf,X%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux#O$}#O#P!3Q#P#o$}#o#p!4h#p#q$}#q#r!3f#r~$}<u&+YZ!OR%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux#O$}#O#P!3Q#P#o$}#o#p!4h#p#q$}#q#r!3f#r~$}G{&,`_T,X%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Uxz$}z{&-_{!_$}!_!`%4a!`#O$}#O#P!3Q#P#o$}#o#p!4h#p#q$}#q#r!3f#r~$}G{&-r]_R%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux!_$}!_!`%4a!`#O$}#O#P!3Q#P#o$}#o#p!4h#p#q$}#q#r!3f#r~$}G{&/O]$z,X%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux!_$}!_!`%4a!`#O$}#O#P!3Q#P#o$}#o#p!4h#p#q$}#q#r!3f#r~$}<u&0[ZwR%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux#O$}#O#P!3Q#P#o$}#o#p!4h#p#q$}#q#r!3f#r~$}Mg&1b^${,X%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux!_$}!_!`%4a!`!a&2^!a#O$}#O#P!3Q#P#o$}#o#p!4h#p#q$}#q#r!3f#r~$}B^&2qZ&S&j%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux#O$}#O#P!3Q#P#o$}#o#p!4h#p#q$}#q#r!3f#r~$}G{&3w_!dQ%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux!O$}!O!P&4v!P!Q$}!Q![&7W![#O$}#O#P!3Q#P#o$}#o#p!4h#p#q$}#q#r!3f#r~$}G{&5X]%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux!O$}!O!P&6Q!P#O$}#O#P!3Q#P#o$}#o#p!4h#p#q$}#q#r!3f#r~$}G{&6eZ!m,X%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux#O$}#O#P!3Q#P#o$}#o#p!4h#p#q$}#q#r!3f#r~$}Gy&7kg!f,V%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux!Q$}!Q![&7W![!g$}!g!h&9S!h!l$}!l!m&=d!m#O$}#O#P!3Q#P#R$}#R#S&7W#S#X$}#X#Y&9S#Y#^$}#^#_&=d#_#o$}#o#p!4h#p#q$}#q#r!3f#r~$}Gy&9ea%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux{$}{|&:j|}$}}!O&:j!O!Q$}!Q![&;t![#O$}#O#P!3Q#P#o$}#o#p!4h#p#q$}#q#r!3f#r~$}Gy&:{]%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux!Q$}!Q![&;t![#O$}#O#P!3Q#P#o$}#o#p!4h#p#q$}#q#r!3f#r~$}Gy&<Xc!f,V%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux!Q$}!Q![&;t![!l$}!l!m&=d!m#O$}#O#P!3Q#P#R$}#R#S&;t#S#^$}#^#_&=d#_#o$}#o#p!4h#p#q$}#q#r!3f#r~$}Gy&=wZ!f,V%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux#O$}#O#P!3Q#P#o$}#o#p!4h#p#q$}#q#r!3f#r~$}G{&>}_$|R%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux!P$}!P!Q&?|!Q!_$}!_!`%4a!`#O$}#O#P!3Q#P#o$}#o#p!4h#p#q$}#q#r!3f#r~$}Gz&@a]%OQ%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux!_$}!_!`%4a!`#O$}#O#P!3Q#P#o$}#o#p!4h#p#q$}#q#r!3f#r~$}Gy&Amu!f,V%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux!O$}!O!P&DQ!P!Q$}!Q![&GW![!d$}!d!e&IY!e!g$}!g!h&9S!h!l$}!l!m&=d!m!q$}!q!r&LS!r!z$}!z!{&Nv!{#O$}#O#P!3Q#P#R$}#R#S&GW#S#U$}#U#V&IY#V#X$}#X#Y&9S#Y#^$}#^#_&=d#_#c$}#c#d&LS#d#l$}#l#m&Nv#m#o$}#o#p!4h#p#q$}#q#r!3f#r~$}Gy&Dc]%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux!Q$}!Q![&E[![#O$}#O#P!3Q#P#o$}#o#p!4h#p#q$}#q#r!3f#r~$}Gy&Eog!f,V%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux!Q$}!Q![&E[![!g$}!g!h&9S!h!l$}!l!m&=d!m#O$}#O#P!3Q#P#R$}#R#S&E[#S#X$}#X#Y&9S#Y#^$}#^#_&=d#_#o$}#o#p!4h#p#q$}#q#r!3f#r~$}Gy&Gki!f,V%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux!O$}!O!P&DQ!P!Q$}!Q![&GW![!g$}!g!h&9S!h!l$}!l!m&=d!m#O$}#O#P!3Q#P#R$}#R#S&GW#S#X$}#X#Y&9S#Y#^$}#^#_&=d#_#o$}#o#p!4h#p#q$}#q#r!3f#r~$}Gy&Ik`%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux!Q$}!Q!R&Jm!R!S&Jm!S#O$}#O#P!3Q#P#R$}#R#S&Jm#S#o$}#o#p!4h#p#q$}#q#r!3f#r~$}Gy&KQ`!f,V%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux!Q$}!Q!R&Jm!R!S&Jm!S#O$}#O#P!3Q#P#R$}#R#S&Jm#S#o$}#o#p!4h#p#q$}#q#r!3f#r~$}Gy&Le_%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux!Q$}!Q!Y&Md!Y#O$}#O#P!3Q#P#R$}#R#S&Md#S#o$}#o#p!4h#p#q$}#q#r!3f#r~$}Gy&Mw_!f,V%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux!Q$}!Q!Y&Md!Y#O$}#O#P!3Q#P#R$}#R#S&Md#S#o$}#o#p!4h#p#q$}#q#r!3f#r~$}Gy' Xc%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux!Q$}!Q!['!d![!c$}!c!i'!d!i#O$}#O#P!3Q#P#R$}#R#S'!d#S#T$}#T#Z'!d#Z#o$}#o#p!4h#p#q$}#q#r!3f#r~$}Gy'!wc!f,V%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux!Q$}!Q!['!d![!c$}!c!i'!d!i#O$}#O#P!3Q#P#R$}#R#S'!d#S#T$}#T#Z'!d#Z#o$}#o#p!4h#p#q$}#q#r!3f#r~$}Mg'$g]x1s%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux!_$}!_!`'%`!`#O$}#O#P!3Q#P#o$}#o#p!4h#p#q$}#q#r!3f#r~$}<u'%sZ%WR%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux#O$}#O#P!3Q#P#o$}#o#p!4h#p#q$}#q#r!3f#r~$}G{'&yZ#^,X%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux#O$}#O#P!3Q#P#o$}#o#p!4h#p#q$}#q#r!3f#r~$}G{'(P_jR%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux!^$}!^!_')O!_!`!9P!`!a!9P!a#O$}#O#P!3Q#P#o$}#o#p!4h#p#q$}#q#r!3f#r~$}Gz')c]$xQ%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux!_$}!_!`%4a!`#O$}#O#P!3Q#P#o$}#o#p!4h#p#q$}#q#r!3f#r~$}G{'*o]%V,X%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux!_$}!_!`!9P!`#O$}#O#P!3Q#P#o$}#o#p!4h#p#q$}#q#r!3f#r~$}G{'+{^jR%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux!_$}!_!`!9P!`!a',w!a#O$}#O#P!3Q#P#o$}#o#p!4h#p#q$}#q#r!3f#r~$}Gz'-[]$yQ%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux!_$}!_!`%4a!`#O$}#O#P!3Q#P#o$}#o#p!4h#p#q$}#q#r!3f#r~$}G{'.j]]Q#tP%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux!_$}!_!`%4a!`#O$}#O#P!3Q#P#o$}#o#p!4h#p#q$}#q#r!3f#r~$}Mg'/xc%p7[%gS%jW%d&j%m`%sp%v!b%x#t%Q,XOr$}rs&Rsw$}wx! Ux!Q$}!Q!['/c![!c$}!c!}'/c!}#O$}#O#P!3Q#P#R$}#R#S'/c#S#T$}#T#o'/c#o#p!4h#p#q$}#q#r!3f#r$g$}$g~'/cMg'1jg%p7[%gS%jW%d&j%m`%sp%v!b%x#t%Q,XOr$}rs'3Rsw$}wx'6mx!Q$}!Q!['/c![!c$}!c!t'/c!t!u';a!u!}'/c!}#O$}#O#P!3Q#P#R$}#R#S'/c#S#T$}#T#f'/c#f#g';a#g#o'/c#o#p!4h#p#q$}#q#r!3f#r$g$}$g~'/cDe'3`_%p7[%gS%e,X%m`%v!bOY!<kYZ'PZ]!<k]^'P^r!<krs'4_sw!<kwx!>yx#O!<k#O#P'6X#P#o!<k#o#p#'w#p#q!<k#q#r#%s#r~!<kDe'4lZ%p7[%gS%e,X%m`%v!bOr'Prs'5_sw'Pwx(Px#O'P#O#PN[#P#o'P#o#pKQ#p#q'P#q#rGW#r~'PD]'5lX%p7[%gS%i,X%m`%v!bOwMOwx,ex#OMO#O#PMv#P#oMO#o#pJ`#p#qMO#q#rIj#r~MODe'6^T%p7[O#o!<k#o#p#%s#p#q!<k#q#r#%s#r~!<kGk'6z_%p7[%jW%e,X%sp%x#tOY%8RYZ!!SZ]%8R]^!!S^r%8Rrs%9csw%8Rwx'7yx#O%8R#O#P'9y#P#o%8R#o#p':_#p#q%8R#q#r& u#r~%8RGk'8WZ%p7[%jW%e,X%sp%x#tOr!!Srs)wsw!!Swx'8yx#O!!S#O#P!2W#P#o!!S#o#p!*j#p#q!!S#q#r!#j#r~!!SFT'9WZ%p7[%jW%f,X%sp%x#tOr!.|rsCWsw!.|wx!.Ox#O!.|#O#P!1r#P#o!.|#o#p!)x#p#q!.|#q#r!&t#r~!.|Gk':OT%p7[O#o%8R#o#p& u#p#q%8R#q#r& u#r~%8R1_':h]%gS%jW%e,XOY%FTYZ;WZ]%FT]^;W^r%FTrs%GPsw%FTwx%Hsx#O%FT#O#P%Im#P#o%FT#o#p& u#p~%FTMg';vc%p7[%gS%jW%d&j%m`%sp%v!b%x#t%Q,XOr$}rs'3Rsw$}wx'6mx!Q$}!Q!['/c![!c$}!c!}'/c!}#O$}#O#P!3Q#P#R$}#R#S'/c#S#T$}#T#o'/c#o#p!4h#p#q$}#q#r!3f#r$g$}$g~'/cMg'=hg%p7[%gS%jW%d&j%m`%sp%v!b%x#t%Q,XOr$}rs'?Psw$}wx'Awx!Q$}!Q!['/c![!c$}!c!t'/c!t!u'Du!u!}'/c!}#O$}#O#P!3Q#P#R$}#R#S'/c#S#T$}#T#f'/c#f#g'Du#g#o'/c#o#p!4h#p#q$}#q#r!3f#r$g$}$g~'/cDe'?^Z%p7[%gS%m`%v!b%r,XOr'Prs'@Psw'Pwx(Px#O'P#O#PNp#P#o'P#o#pKQ#p#q'P#q#rGW#r~'PDe'@[Z%p7[%gS%m`%v!bOr'Prs'@}sw'Pwx(Px#O'P#O#PN[#P#o'P#o#pKQ#p#q'P#q#rGW#r~'PD]'A[X%p7[%gS%w,X%m`%v!bOwMOwx,ex#OMO#O#PMv#P#oMO#o#pJ`#p#qMO#q#rIj#r~MOGk'BUZ%p7[%jW%sp%x#t%l,XOr!!Srs)wsw!!Swx'Bwx#O!!S#O#P!2l#P#o!!S#o#p!+d#p#q!!S#q#r!#j#r~!!SGk'CSZ%p7[%jW%sp%x#tOr!!Srs)wsw!!Swx'Cux#O!!S#O#P!2W#P#o!!S#o#p!*j#p#q!!S#q#r!#j#r~!!SFT'DSZ%p7[%jW%u,X%sp%x#tOr!.|rsCWsw!.|wx!.Ox#O!.|#O#P!1r#P#o!.|#o#p!)x#p#q!.|#q#r!&t#r~!.|Mg'E[c%p7[%gS%jW%d&j%m`%sp%v!b%x#t%Q,XOr$}rs'?Psw$}wx'Awx!Q$}!Q!['/c![!c$}!c!}'/c!}#O$}#O#P!3Q#P#R$}#R#S'/c#S#T$}#T#o'/c#o#p!4h#p#q$}#q#r!3f#r$g$}$g~'/cMg'F|k%p7[%gS%jW%d&j%m`%sp%v!b%x#t%Q,XOr$}rs'3Rsw$}wx'6mx!Q$}!Q!['/c![!c$}!c!h'/c!h!i'Du!i!t'/c!t!u';a!u!}'/c!}#O$}#O#P!3Q#P#R$}#R#S'/c#S#T$}#T#U'/c#U#V';a#V#Y'/c#Y#Z'Du#Z#o'/c#o#p!4h#p#q$}#q#r!3f#r$g$}$g~'/cG{'IUZ!V,X%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux#O$}#O#P!3Q#P#o$}#o#p!4h#p#q$}#q#r!3f#r~$}Mg'I|X%p7[OY$}YZ!5[Z]$}]^!5[^#o$}#o#p!3f#p#q$}#q#r!3f#r~$}<u'J|Z!WR%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux#O$}#O#P!3Q#P#o$}#o#p!4h#p#q$}#q#r!3f#r~$}Gz'LS]$vQ%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux!_$}!_!`%4a!`#O$}#O#P!3Q#P#o$}#o#p!4h#p#q$}#q#r!3f#r~$}Gy'MUX%gS%jW!ZGmOr;Wrs;tsw;Wwx<zx#O;W#O#P>j#P#o;W#o#p!3f#p~;WGz'NU]$uQ%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux!_$}!_!`%4a!`#O$}#O#P!3Q#P#o$}#o#p!4h#p#q$}#q#r!3f#r~$}<u( `X![7_%gS%jW%m`%sp%v!b%x#tOr!3frsHOsw!3fwx!$dx#O!3f#O#P!4b#P#o!3f#o#p!4h#p~!3fGy(!`Z%P,V%p7[%gS%jW%m`%sp%v!b%x#tOr$}rs&Rsw$}wx! Ux#O$}#O#P!3Q#P#o$}#o#p!4h#p#q$}#q#r!3f#r~$}",
  tokenizers: [legacyPrint, indentation, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, newlines],
  topRules: {Script: [0, 3]},
  specialized: [{term: 186, get: (value) => spec_identifier3[value] || -1}],
  tokenPrec: 6594
});

// node_modules/@codemirror/lang-python/dist/index.js
var pythonLanguage = LezerLanguage.define({
  parser: parser3.configure({
    props: [
      indentNodeProp.add({
        Body: continuedIndent()
      }),
      foldNodeProp.add({
        "Body ArrayExpression DictionaryExpression": foldInside
      }),
      styleTags({
        "async '*' '**' FormatConversion": tags.modifier,
        "for while if elif else try except finally return raise break continue with pass assert await yield": tags.controlKeyword,
        "in not and or is del": tags.operatorKeyword,
        "import from def class global nonlocal lambda": tags.definitionKeyword,
        "with as print": tags.keyword,
        self: tags.self,
        Boolean: tags.bool,
        None: tags.null,
        VariableName: tags.variableName,
        "CallExpression/VariableName": tags.function(tags.variableName),
        "FunctionDefinition/VariableName": tags.function(tags.definition(tags.variableName)),
        "ClassDefinition/VariableName": tags.definition(tags.className),
        PropertyName: tags.propertyName,
        "CallExpression/MemberExpression/ProperyName": tags.function(tags.propertyName),
        Comment: tags.lineComment,
        Number: tags.number,
        String: tags.string,
        FormatString: tags.special(tags.string),
        UpdateOp: tags.updateOperator,
        ArithOp: tags.arithmeticOperator,
        BitOp: tags.bitwiseOperator,
        CompareOp: tags.compareOperator,
        AssignOp: tags.definitionOperator,
        Ellipsis: tags.punctuation,
        At: tags.meta,
        "( )": tags.paren,
        "[ ]": tags.squareBracket,
        "{ }": tags.brace,
        ".": tags.derefOperator,
        ", ;": tags.separator
      })
    ]
  }),
  languageData: {
    closeBrackets: {brackets: ["(", "[", "{", "'", '"', "'''", '"""']},
    commentTokens: {line: "#"},
    indentOnInput: /^\s*[\}\]\)]$/
  }
});
function python() {
  return new LanguageSupport(pythonLanguage);
}

// node_modules/@codemirror/legacy-modes/mode/ruby.js
function wordObj(words4) {
  var o = {};
  for (var i = 0, e = words4.length; i < e; ++i)
    o[words4[i]] = true;
  return o;
}
var keywordList = [
  "alias",
  "and",
  "BEGIN",
  "begin",
  "break",
  "case",
  "class",
  "def",
  "defined?",
  "do",
  "else",
  "elsif",
  "END",
  "end",
  "ensure",
  "false",
  "for",
  "if",
  "in",
  "module",
  "next",
  "not",
  "or",
  "redo",
  "rescue",
  "retry",
  "return",
  "self",
  "super",
  "then",
  "true",
  "undef",
  "unless",
  "until",
  "when",
  "while",
  "yield",
  "nil",
  "raise",
  "throw",
  "catch",
  "fail",
  "loop",
  "callcc",
  "caller",
  "lambda",
  "proc",
  "public",
  "protected",
  "private",
  "require",
  "load",
  "require_relative",
  "extend",
  "autoload",
  "__END__",
  "__FILE__",
  "__LINE__",
  "__dir__"
];
var keywords8 = wordObj(keywordList);
var indentWords = wordObj([
  "def",
  "class",
  "case",
  "for",
  "while",
  "until",
  "module",
  "then",
  "catch",
  "loop",
  "proc",
  "begin"
]);
var dedentWords = wordObj(["end", "until"]);
var opening = {"[": "]", "{": "}", "(": ")"};
var closing2 = {"]": "[", "}": "{", ")": "("};
var curPunc2;
function chain2(newtok, stream, state) {
  state.tokenize.push(newtok);
  return newtok(stream, state);
}
function tokenBase6(stream, state) {
  if (stream.sol() && stream.match("=begin") && stream.eol()) {
    state.tokenize.push(readBlockComment);
    return "comment";
  }
  if (stream.eatSpace())
    return null;
  var ch = stream.next(), m;
  if (ch == "`" || ch == "'" || ch == '"') {
    return chain2(readQuoted(ch, "string", ch == '"' || ch == "`"), stream, state);
  } else if (ch == "/") {
    if (regexpAhead(stream))
      return chain2(readQuoted(ch, "string.special", true), stream, state);
    else
      return "operator";
  } else if (ch == "%") {
    var style = "string", embed = true;
    if (stream.eat("s"))
      style = "atom";
    else if (stream.eat(/[WQ]/))
      style = "string";
    else if (stream.eat(/[r]/))
      style = "string.special";
    else if (stream.eat(/[wxq]/)) {
      style = "string";
      embed = false;
    }
    var delim = stream.eat(/[^\w\s=]/);
    if (!delim)
      return "operator";
    if (opening.propertyIsEnumerable(delim))
      delim = opening[delim];
    return chain2(readQuoted(delim, style, embed, true), stream, state);
  } else if (ch == "#") {
    stream.skipToEnd();
    return "comment";
  } else if (ch == "<" && (m = stream.match(/^<([-~])[\`\"\']?([a-zA-Z_?]\w*)[\`\"\']?(?:;|$)/))) {
    return chain2(readHereDoc(m[2], m[1]), stream, state);
  } else if (ch == "0") {
    if (stream.eat("x"))
      stream.eatWhile(/[\da-fA-F]/);
    else if (stream.eat("b"))
      stream.eatWhile(/[01]/);
    else
      stream.eatWhile(/[0-7]/);
    return "number";
  } else if (/\d/.test(ch)) {
    stream.match(/^[\d_]*(?:\.[\d_]+)?(?:[eE][+\-]?[\d_]+)?/);
    return "number";
  } else if (ch == "?") {
    while (stream.match(/^\\[CM]-/)) {
    }
    if (stream.eat("\\"))
      stream.eatWhile(/\w/);
    else
      stream.next();
    return "string";
  } else if (ch == ":") {
    if (stream.eat("'"))
      return chain2(readQuoted("'", "atom", false), stream, state);
    if (stream.eat('"'))
      return chain2(readQuoted('"', "atom", true), stream, state);
    if (stream.eat(/[\<\>]/)) {
      stream.eat(/[\<\>]/);
      return "atom";
    }
    if (stream.eat(/[\+\-\*\/\&\|\:\!]/)) {
      return "atom";
    }
    if (stream.eat(/[a-zA-Z$@_\xa1-\uffff]/)) {
      stream.eatWhile(/[\w$\xa1-\uffff]/);
      stream.eat(/[\?\!\=]/);
      return "atom";
    }
    return "operator";
  } else if (ch == "@" && stream.match(/^@?[a-zA-Z_\xa1-\uffff]/)) {
    stream.eat("@");
    stream.eatWhile(/[\w\xa1-\uffff]/);
    return "propertyName";
  } else if (ch == "$") {
    if (stream.eat(/[a-zA-Z_]/)) {
      stream.eatWhile(/[\w]/);
    } else if (stream.eat(/\d/)) {
      stream.eat(/\d/);
    } else {
      stream.next();
    }
    return "variableName.special";
  } else if (/[a-zA-Z_\xa1-\uffff]/.test(ch)) {
    stream.eatWhile(/[\w\xa1-\uffff]/);
    stream.eat(/[\?\!]/);
    if (stream.eat(":"))
      return "atom";
    return "variable";
  } else if (ch == "|" && (state.varList || state.lastTok == "{" || state.lastTok == "do")) {
    curPunc2 = "|";
    return null;
  } else if (/[\(\)\[\]{}\\;]/.test(ch)) {
    curPunc2 = ch;
    return null;
  } else if (ch == "-" && stream.eat(">")) {
    return "operator";
  } else if (/[=+\-\/*:\.^%<>~|]/.test(ch)) {
    var more = stream.eatWhile(/[=+\-\/*:\.^%<>~|]/);
    if (ch == "." && !more)
      curPunc2 = ".";
    return "operator";
  } else {
    return null;
  }
}
function regexpAhead(stream) {
  var start = stream.pos, depth = 0, next, found = false, escaped = false;
  while ((next = stream.next()) != null) {
    if (!escaped) {
      if ("[{(".indexOf(next) > -1) {
        depth++;
      } else if ("]})".indexOf(next) > -1) {
        depth--;
        if (depth < 0)
          break;
      } else if (next == "/" && depth == 0) {
        found = true;
        break;
      }
      escaped = next == "\\";
    } else {
      escaped = false;
    }
  }
  stream.backUp(stream.pos - start);
  return found;
}
function tokenBaseUntilBrace(depth) {
  if (!depth)
    depth = 1;
  return function(stream, state) {
    if (stream.peek() == "}") {
      if (depth == 1) {
        state.tokenize.pop();
        return state.tokenize[state.tokenize.length - 1](stream, state);
      } else {
        state.tokenize[state.tokenize.length - 1] = tokenBaseUntilBrace(depth - 1);
      }
    } else if (stream.peek() == "{") {
      state.tokenize[state.tokenize.length - 1] = tokenBaseUntilBrace(depth + 1);
    }
    return tokenBase6(stream, state);
  };
}
function tokenBaseOnce() {
  var alreadyCalled = false;
  return function(stream, state) {
    if (alreadyCalled) {
      state.tokenize.pop();
      return state.tokenize[state.tokenize.length - 1](stream, state);
    }
    alreadyCalled = true;
    return tokenBase6(stream, state);
  };
}
function readQuoted(quote, style, embed, unescaped) {
  return function(stream, state) {
    var escaped = false, ch;
    if (state.context.type === "read-quoted-paused") {
      state.context = state.context.prev;
      stream.eat("}");
    }
    while ((ch = stream.next()) != null) {
      if (ch == quote && (unescaped || !escaped)) {
        state.tokenize.pop();
        break;
      }
      if (embed && ch == "#" && !escaped) {
        if (stream.eat("{")) {
          if (quote == "}") {
            state.context = {prev: state.context, type: "read-quoted-paused"};
          }
          state.tokenize.push(tokenBaseUntilBrace());
          break;
        } else if (/[@\$]/.test(stream.peek())) {
          state.tokenize.push(tokenBaseOnce());
          break;
        }
      }
      escaped = !escaped && ch == "\\";
    }
    return style;
  };
}
function readHereDoc(phrase, mayIndent) {
  return function(stream, state) {
    if (mayIndent)
      stream.eatSpace();
    if (stream.match(phrase))
      state.tokenize.pop();
    else
      stream.skipToEnd();
    return "string";
  };
}
function readBlockComment(stream, state) {
  if (stream.sol() && stream.match("=end") && stream.eol())
    state.tokenize.pop();
  stream.skipToEnd();
  return "comment";
}
var ruby = {
  startState: function(indentUnit2) {
    return {
      tokenize: [tokenBase6],
      indented: 0,
      context: {type: "top", indented: -indentUnit2},
      continuedLine: false,
      lastTok: null,
      varList: false
    };
  },
  token: function(stream, state) {
    curPunc2 = null;
    if (stream.sol())
      state.indented = stream.indentation();
    var style = state.tokenize[state.tokenize.length - 1](stream, state), kwtype;
    var thisTok = curPunc2;
    if (style == "variable") {
      var word = stream.current();
      style = state.lastTok == "." ? "property" : keywords8.propertyIsEnumerable(stream.current()) ? "keyword" : /^[A-Z]/.test(word) ? "tag" : state.lastTok == "def" || state.lastTok == "class" || state.varList ? "def" : "variable";
      if (style == "keyword") {
        thisTok = word;
        if (indentWords.propertyIsEnumerable(word))
          kwtype = "indent";
        else if (dedentWords.propertyIsEnumerable(word))
          kwtype = "dedent";
        else if ((word == "if" || word == "unless") && stream.column() == stream.indentation())
          kwtype = "indent";
        else if (word == "do" && state.context.indented < state.indented)
          kwtype = "indent";
      }
    }
    if (curPunc2 || style && style != "comment")
      state.lastTok = thisTok;
    if (curPunc2 == "|")
      state.varList = !state.varList;
    if (kwtype == "indent" || /[\(\[\{]/.test(curPunc2))
      state.context = {prev: state.context, type: curPunc2 || style, indented: state.indented};
    else if ((kwtype == "dedent" || /[\)\]\}]/.test(curPunc2)) && state.context.prev)
      state.context = state.context.prev;
    if (stream.eol())
      state.continuedLine = curPunc2 == "\\" || style == "operator";
    return style;
  },
  indent: function(state, textAfter, cx) {
    if (state.tokenize[state.tokenize.length - 1] != tokenBase6)
      return null;
    var firstChar = textAfter && textAfter.charAt(0);
    var ct = state.context;
    var closed = ct.type == closing2[firstChar] || ct.type == "keyword" && /^(?:end|until|else|elsif|when|rescue)\b/.test(textAfter);
    return ct.indented + (closed ? 0 : cx.unit) + (state.continuedLine ? cx.unit : 0);
  },
  languageData: {
    indentOnInput: /^\s*(?:end|rescue|elsif|else|\})$/,
    commentTokens: {line: "#"},
    autocomplete: keywordList
  }
};

// node_modules/lezer-rust/dist/index.es.js
var closureParamDelim = 1;
var tpOpen = 2;
var tpClose = 3;
var RawString = 4;
var Float = 5;
var _b = 98;
var _e = 101;
var _f = 102;
var _r = 114;
var _E = 69;
var Dot = 46;
var Plus = 43;
var Minus = 45;
var Hash = 35;
var Quote = 34;
var Pipe = 124;
var LessThan = 60;
var GreaterThan = 62;
function isNum(ch) {
  return ch >= 48 && ch <= 57;
}
function isNum_(ch) {
  return isNum(ch) || ch == 95;
}
var literalTokens = new ExternalTokenizer((input, token, stack) => {
  let pos = token.start, next = input.get(pos);
  if (isNum(next)) {
    let isFloat = false;
    do {
      next = input.get(++pos);
    } while (isNum_(next));
    if (next == Dot) {
      isFloat = true;
      next = input.get(++pos);
      if (isNum(next)) {
        do {
          next = input.get(++pos);
        } while (isNum_(next));
      } else if (next == Dot || next > 127 || /\w/.test(String.fromCharCode(next))) {
        return;
      }
    }
    if (next == _e || next == _E) {
      isFloat = true;
      next = input.get(++pos);
      if (next == Plus || next == Minus)
        next = input.get(++pos);
      let startNum = pos;
      while (isNum_(next))
        next = input.get(++pos);
      if (pos == startNum)
        return;
    }
    if (next == _f) {
      if (!/32|64/.test(input.read(pos + 1, pos + 3)))
        return;
      isFloat = true;
      pos += 3;
    }
    if (isFloat)
      token.accept(Float, pos);
  } else if (next == _b || next == _r) {
    if (next == _b)
      next = input.get(++pos);
    if (next != _r)
      return;
    next = input.get(++pos);
    let count = 0;
    while (next == Hash) {
      count++;
      next = input.get(++pos);
    }
    if (next != Quote)
      return;
    next = input.get(++pos);
    content:
      for (; ; ) {
        if (next < 0)
          return;
        let isQuote = next == Quote;
        next = input.get(++pos);
        if (isQuote) {
          for (let i = 0; i < count; i++) {
            if (next != Hash)
              continue content;
            next = input.get(++pos);
          }
          token.accept(RawString, pos);
          return;
        }
      }
  }
});
var closureParam = new ExternalTokenizer((input, token) => {
  if (input.get(token.start) == Pipe)
    token.accept(closureParamDelim, token.start + 1);
});
var tpDelim = new ExternalTokenizer((input, token) => {
  let pos = token.start, next = input.get(pos);
  if (next == LessThan)
    token.accept(tpOpen, pos + 1);
  else if (next == GreaterThan)
    token.accept(tpClose, pos + 1);
});
var spec_identifier4 = {__proto__: null, self: 28, super: 32, crate: 34, impl: 46, true: 72, false: 72, pub: 88, in: 92, const: 96, unsafe: 104, async: 108, move: 110, if: 114, let: 118, ref: 142, mut: 144, _: 198, else: 200, match: 204, as: 248, return: 252, await: 262, break: 270, continue: 276, while: 312, loop: 316, for: 320, macro_rules: 327, mod: 334, extern: 342, struct: 346, where: 364, union: 379, enum: 382, type: 390, default: 395, fn: 396, trait: 412, use: 420, static: 438, dyn: 476};
var parser4 = Parser.deserialize({
  version: 13,
  states: "$3tQ]Q_OOP$wOWOOO&sQWO'#CnO)WQWO'#IaOOQP'#Ia'#IaOOQQ'#If'#IfO)hO`O'#C}OOQR'#Ii'#IiO)sQWO'#IvOOQO'#Hk'#HkO)xQWO'#DpOOQR'#Ix'#IxO)xQWO'#DpO*ZQWO'#DpOOQO'#Iw'#IwO,SQWO'#J`O,ZQWO'#EiOOQV'#Hp'#HpO,cQYO'#F{OOQV'#El'#ElOOQV'#Em'#EmOOQV'#En'#EnO.YQ_O'#EkO0_Q_O'#EoO2gQWOOO4QQ_O'#FPO7hQWO'#J`OOQV'#FY'#FYO7{Q_O'#F^O:WQ_O'#FaOOQO'#F`'#F`O=sQ_O'#FcO=}Q_O'#FbO@VQWO'#FgOOQO'#J`'#J`OOQV'#Ip'#IpOA]Q_O'#IoOEPQWO'#IoOOQV'#Fw'#FwOF[QWO'#JuOFcQWO'#F|OOQO'#IO'#IOOGrQWO'#GhOOQV'#In'#InOOQV'#Im'#ImOOQV'#Hj'#HjQGyQ_OOOKeQ_O'#DUOKlQYO'#CqOOQP'#I`'#I`OOQV'#Hg'#HgQ]Q_OOOLuQWO'#IaONsQYO'#DXO!!eQWO'#JuO!!lQWO'#JuO!!vQ_O'#DfO!%]Q_O'#E}O!(sQ_O'#FWO!,ZQWO'#FZO!.^QXO'#FbO!.cQ_O'#EeO!!vQ_O'#FmO!0uQWO'#FoO!0zQWO'#FoO!1PQ^O'#FqO!1WQWO'#JuO!1_QWO'#FtO!1dQWO'#FxO!2WQWO'#JjO!2_QWO'#GOO!2_QWO'#G`O!2_QWO'#GbO!2_QWO'#GsOOQO'#Ju'#JuO!2dQWO'#GhO!2lQYO'#GpO!2_QWO'#GqO!3uQ^O'#GtO!3|QWO'#GuO!4hQWO'#HOP!4sOpO'#CcPOOO)CDO)CDOOOOO'#Hi'#HiO!5OO`O,59iOOQV,59i,59iO!5ZQYO,5?bOOQO-E;i-E;iOOQO,5:[,5:[OOQP,59Z,59ZO)xQWO,5:[O)xQWO,5:[O!5oQWO,5?lO!5zQYO,5;qO!6PQYO,5;TO!6hQWO,59QO!7kQXO'#CnO!7rQXO'#IaO!8vQWO'#CoO,^QWO'#EiOOQV-E;n-E;nO!9XQWO'#FsOOQV,5<g,5<gO!8vQWO'#CoO!9^QWO'#CoO!9cQWO'#IaO! yQWO'#JuO!9mQWO'#J`O!:TQWO,5;VOOQO'#Io'#IoO!0zQWO'#DaO!<TQWO'#DcO!<]QWO,5;ZO.YQ_O,5;ZOOQO,5;[,5;[OOQV'#Er'#ErOOQV'#Es'#EsOOQV'#Et'#EtOOQV'#Eu'#EuOOQV'#Ev'#EvOOQV'#Ew'#EwOOQV'#Ex'#ExOOQV'#Ey'#EyO.YQ_O,5;]O.YQ_O,5;]O.YQ_O,5;]O.YQ_O,5;]O.YQ_O,5;]O.YQ_O,5;]O.YQ_O,5;]O.YQ_O,5;]O.YQ_O,5;]O.YQ_O,5;fO!<sQ_O,5;kO!@ZQ_O'#FROOQO,5;l,5;lO!BfQWO,5;pO.YQ_O,5;wOKlQYO,5;gO!DRQWO,5;kO!DrQWO,5;xOOQO,5;x,5;xO!EPQWO,5;xO!EUQ_O,5;xO!GaQWO'#CfO!GfQWO,5<QO!GpQ_O,5<QOOQO,5;{,5;{O!J^QXO'#CnO!KoQXO'#IaOOQS'#Dk'#DkOOQP'#Is'#IsO!LiQ[O'#IsO!LqQXO'#DjO!MoQWO'#DnO!MoQWO'#DnO!NQQWO'#DnOOQP'#Iu'#IuO!NVQXO'#IuO# QQ^O'#DoO# [QWO'#DrO# dQ^O'#DzO# nQ^O'#D|O# uQWO'#EPO#!QQXO'#FdOOQP'#ES'#ESOOQP'#Ir'#IrO#!`QXO'#JfOOQP'#Je'#JeO#!hQXO,5;}O#!mQXO'#IaO!1PQ^O'#DyO!1PQ^O'#FdO##gQWO,5;|OOQO,5;|,5;|OKlQYO,5;|O##}QWO'#FhOOQO,5<R,5<ROOQV,5=l,5=lO#&SQYO'#FzOOQV,5<h,5<hO#&ZQWO,5<hO#&bQWO,5=SO!1WQWO,59rO!1dQWO,5<dO#&iQWO,5=iO!2_QWO,5<jO!2_QWO,5<zO!2_QWO,5<|O!2_QWO,5=QO#&pQWO,5=]O#&wQWO,5=SO!2_QWO,5=]O!3|QWO,5=aO#'PQWO,5=jOOQO-E;|-E;|O#'[QWO'#JjOOQV-E;h-E;hO#'sQWO'#HRO#'zQ_O,59pOOQV,59p,59pO#(RQWO,59pO#(WQ_O,59pO#(vQZO'#CuO#+OQZO'#CvOOQV'#C|'#C|O#-kQWO'#HTO#-rQYO'#IeOOQO'#Hh'#HhO#-zQWO'#CwO#-zQWO'#CwO#.]QWO'#CwOOQR'#Id'#IdO#.bQZO'#IcO#0wQYO'#HTO#1eQYO'#H[O#2qQYO'#H_OKlQYO'#H`OOQR'#Hb'#HbO#3}QWO'#HeO#4SQYO,59]OOQR'#Ic'#IcO#4sQZO'#CtO#7OQYO'#HUO#7TQWO'#HTO#7YQYO'#CrO#7yQWO'#H]O#7YQYO'#HcOOQV-E;e-E;eO#8RQWO,59sOOQV,59{,59{O#8aQYO,5=[OOQV,59},59}O!0zQWO,59}O#;TQWO'#IqOOQO'#Iq'#IqO!1PQ^O'#DhO!0zQWO,5:QO#;[QWO,5;iO#;rQWO,5;rO#<YQ_O,5;rOOQO,5;u,5;uO#?sQ_O,5;|O#A{QWO,5;PO!0zQWO,5<XO#BSQWO,5<ZOOQV,5<Z,5<ZO#B_QWO,5<]O!1PQ^O'#EOOOQQ'#D_'#D_O#BgQWO,59rO#BlQWO,5<`O#BqQWO,5<dOOQO,5@U,5@UO#ByQWO,5=iOOQQ'#Cv'#CvO#COQYO,5<jO#CaQYO,5<zO#ClQYO,5<|O#CwQYO,5=_O#DVQYO,5=SO#EoQYO'#GQO#E|QYO,5=[O#FaQWO,5=[O#FoQYO,5=[O#GxQYO,5=]O#HWQWO,5=`O!1PQ^O,5=`O#HfQWO'#CnO#HwQWO'#IaOOQO'#Jy'#JyO#IYQWO'#IQO#I_QWO'#GwOOQO'#Jz'#JzO#IvQWO'#GzOOQO'#G|'#G|OOQO'#Jx'#JxO#I_QWO'#GwO#I}QWO'#GxO#JSQWO,5=aO#JXQWO,5=jO!1dQWO,5=jO#'SQWO,5=jPOOO'#Hf'#HfP#J^OpO,58}POOO,58},58}OOOO-E;g-E;gOOQV1G/T1G/TO#JiQWO1G4|O#JnQ^O'#CyPOQQ'#Cx'#CxOOQO1G/v1G/vOOQP1G.u1G.uO)xQWO1G/vO#MwQ!fO'#EUO#NOQ!fO'#EVO#NVQ!fO'#ETO$ _QWO1G5WO$!RQ_O1G5WOOQO1G1]1G1]O$%uQWO1G0oO$%zQWO'#CiO!7rQXO'#IaO!6PQYO1G.lO!5oQWO,5<_O!8vQWO,59ZO!8vQWO,59ZO!5oQWO,5?lO$&]QWO1G0uO$(jQWO1G0wO$*bQWO1G0wO$*xQWO1G0wO$,|QWO1G0wO$-TQWO1G0wO$/UQWO1G0wO$/]QWO1G0wO$1^QWO1G0wO$1eQWO1G0wO$2|QWO1G1QO$4}QWO1G1VO$5nQ_O'#JcO$7vQWO'#JcOOQO'#Jb'#JbO$8QQWO,5;mOOQO'#Dw'#DwOOQO1G1[1G1[OOQO1G1Y1G1YO$8VQWO1G1cOOQO1G1R1G1RO$8^Q_O'#HrO$:lQWO,5@OO.YQ_O1G1dOOQO1G1d1G1dO$:tQWO1G1dO$;RQWO1G1dO$;WQWO1G1eOOQO1G1l1G1lO$;`QWO1G1lOOQP,5?_,5?_O$;jQ^O,5:kO$<TQXO,5:YO!MoQWO,5:YO!MoQWO,5:YO!1PQ^O,5:gO$=UQWO'#IzOOQO'#Iy'#IyO$=dQWO,5:ZO# QQ^O,5:ZO$=iQWO'#DsOOQP,5:^,5:^O$=zQWO,5:fOOQP,5:h,5:hO!1PQ^O,5:hO!1PQ^O,5:mO$>PQYO,5<OO$>ZQ_O'#HsO$>hQXO,5@QOOQV1G1i1G1iOOQP,5:e,5:eO$>pQXO,5<OO$?OQWO1G1hO$?WQWO'#CnO$?cQWO'#FiOOQO'#Fi'#FiO$?kQWO'#FjO.YQ_O'#FkOOQO'#Ji'#JiO$?pQWO'#JhOOQO'#Jg'#JgO$?xQWO,5<SOOQQ'#Hv'#HvO$?}QYO,5<fOOQV,5<f,5<fO$@UQYO,5<fOOQV1G2S1G2SO$@]QWO1G2nO$@eQWO1G/^O$@jQWO1G2OO#ByQWO1G3TO$@rQYO1G2UO#CaQYO1G2fO#ClQYO1G2hO$ATQYO1G2lO!2_QWO1G2wO#DVQYO1G2nO#GxQYO1G2wO$A]QWO1G2{O$AbQWO1G3UO!1dQWO1G3UO$AgQWO1G3UOOQV1G/[1G/[O$AoQWO1G/[O$AtQ_O1G/[O#7TQWO,5=oO$A{QYO,5?PO$BaQWO,5?PO$BfQZO'#IfOOQO-E;f-E;fOOQR,59c,59cO#-zQWO,59cO#-zQWO,59cOOQR,5=n,5=nO$ERQYO'#HVO$FkQZO,5=oO!5oQWO,5={O$H}QWO,5=oO$IUQZO,5=vO$KeQYO,5=vO$>PQYO,5=vO$KuQWO'#KRO$LQQWO,5=xOOQR,5=y,5=yO$LVQWO,5=zO$>PQYO,5>PO$>PQYO,5>POOQO1G.w1G.wO$>PQYO1G.wO$LbQYO,5=pO$LjQZO,59^OOQR,59^,59^O$>PQYO,5=wO$N|QZO,5=}OOQR,5=},5=}O%#`QWO1G/_O!6PQYO1G/_O#E|QYO1G2vO%#eQWO1G2vO%#sQYO1G2vOOQV1G/i1G/iO%$|QWO,5:SO%%UQ_O1G/lO%*_QWO1G1^O%*uQWO1G1hOOQO1G1h1G1hO$>PQYO1G1hO%+]Q^O'#EgOOQV1G0k1G0kOOQV1G1s1G1sO!!vQ_O1G1sO!0zQWO1G1uO!1PQ^O1G1wO!.cQ_O1G1wOOQP,5:j,5:jO$>PQYO1G/^OOQO'#Cn'#CnO%+jQWO1G1zOOQV1G2O1G2OO%+rQWO'#CnO%+zQWO1G3TO%,PQWO1G3TO%,UQYO'#GQO%,gQWO'#G]O%,xQYO'#G_O%.[QYO'#GXOOQV1G2U1G2UO%/kQWO1G2UO%/pQWO1G2UO$@uQWO1G2UOOQV1G2f1G2fO%/kQWO1G2fO#CdQWO1G2fO%/xQWO'#GdOOQV1G2h1G2hO%0ZQWO1G2hO#CoQWO1G2hO%0`QYO'#GSO$>PQYO1G2lO$AWQWO1G2lOOQV1G2y1G2yO%1lQWO1G2yO%3[Q^O'#GkO%3fQWO1G2nO#DYQWO1G2nO%3tQYO,5<lO%4OQYO,5<lO%4^QYO,5<lO%4{QYO,5<lOOQQ,5<l,5<lO!1WQWO'#JuO%5WQYO,5<lO%5`QWO1G2vOOQV1G2v1G2vO%5hQWO1G2vO$>PQYO1G2vOOQV1G2w1G2wO%5hQWO1G2wO%5mQWO1G2wO#G{QWO1G2wOOQV1G2z1G2zO.YQ_O1G2zO$>PQYO1G2zO%5uQWO1G2zOOQO,5>l,5>lOOQO-E<O-E<OOOQO,5=c,5=cOOQO,5=e,5=eOOQO,5=g,5=gOOQO,5=h,5=hO%6TQWO'#J|OOQO'#J{'#J{O%6]QWO,5=fO%6bQWO,5=cO!1dQWO,5=dOOQV1G2{1G2{O$>PQYO1G3UPOOO-E;d-E;dPOOO1G.i1G.iOOQO7+*h7+*hO%6yQYO'#IdO%7bQYO'#IgO%7mQYO'#IgO%7uQYO'#IgO%8QQYO,59eOOQO7+%b7+%bOOQP7+$a7+$aOOQV,5:p,5:pO%8VQ!fO,5:pO%8^Q!fO'#JTOOQS'#EZ'#EZOOQS'#E['#E[OOQS'#E]'#E]OOQS'#JT'#JTO%;PQWO'#EYOOQS'#Eb'#EbOOQS'#JR'#JROOQS'#Hn'#HnOOQV,5:q,5:qO%;UQ!fO,5:qO%;]Q!fO,5:oOOQV,5:o,5:oOOQV7+'e7+'eOOQV7+&Z7+&ZO%;dQ[O,59TO%;xQ^O,59TO%<cQWO7+$WO%<hQWO1G1yOOQV1G1y1G1yO!8vQWO1G.uOOQP1G5W1G5WO%<mQWO,5?}O%<wQ_O'#HqO%?SQWO,5?}OOQO1G1X1G1XOOQO7+&}7+&}O%?[QWO,5>^OOQO-E;p-E;pO%?iQWO7+'OO%?pQ_O7+'OOOQO7+'O7+'OOOQO7+'P7+'PO%ArQWO7+'POOQO7+'W7+'WOOQP1G0V1G0VO%AzQXO1G/tO!MoQWO1G/tO%B{QXO1G0RO%CsQ^O'#HlO%DTQWO,5?fOOQP1G/u1G/uO%D`QWO1G/uO%DeQWO'#D_OOQO'#Dt'#DtO%DpQWO'#DtO%DuQWO'#I|OOQO'#I{'#I{O%D}QWO,5:_O%ESQWO'#DtO%EXQWO'#DtOOQP1G0Q1G0QOOQP1G0S1G0SOOQP1G0X1G0XO%EaQXO1G1jO%ElQXO'#FeOOQP,5>_,5>_O!1PQ^O'#FeOOQP-E;q-E;qO$>PQYO1G1jOOQO7+'S7+'SOOQO,5<T,5<TO%EzQWO,5<UO%?pQ_O,5<UO%FPQWO,5<VO%FZQWO'#HtO%FlQWO,5@SOOQO1G1n1G1nOOQQ-E;t-E;tOOQV1G2Q1G2QO%FtQYO1G2QO#DVQYO7+(YO$>PQYO7+$xOOQV7+'j7+'jO%F{QWO7+(oO%GQQWO7+(oOOQV7+'p7+'pO%/kQWO7+'pO%GVQWO7+'pO%G_QWO7+'pOOQV7+(Q7+(QO%/kQWO7+(QO#CdQWO7+(QOOQV7+(S7+(SO%0ZQWO7+(SO#CoQWO7+(SO$>PQYO7+(WO%GmQWO7+(WO#GxQYO7+(cO%GrQWO7+(YO#DYQWO7+(YOOQV7+(c7+(cO%5hQWO7+(cO%5mQWO7+(cO#G{QWO7+(cOOQV7+(g7+(gO$>PQYO7+(pO%HQQWO7+(pO!1dQWO7+(pOOQV7+$v7+$vO%HVQWO7+$vO%H[QZO1G3ZO%JnQWO1G4kOOQO1G4k1G4kOOQR1G.}1G.}O#-zQWO1G.}O%JsQWO'#KQOOQO'#HW'#HWO%KUQWO'#HXO%KaQWO'#KQOOQO'#KP'#KPO%KiQWO,5=qO%KnQYO'#H[O%LzQWO'#GmO%MVQYO'#CtO%MaQWO'#GmO$>PQYO1G3ZOOQR1G3g1G3gO#7TQWO1G3ZO%MfQZO1G3bO$>PQYO1G3bO& uQYO'#IVO&!VQWO,5@mOOQR1G3d1G3dOOQR1G3f1G3fO%?pQ_O1G3fOOQR1G3k1G3kO&!_QYO7+$cO&!gQYO'#KOOOQQ'#J}'#J}O&!oQYO1G3[O&!tQZO1G3cOOQQ7+$y7+$yO&%TQWO7+$yO&%YQWO7+(bOOQV7+(b7+(bO%5hQWO7+(bO$>PQYO7+(bO#E|QYO7+(bO&%bQWO7+(bO!.cQ_O1G/nO&%pQWO7+%WO$?OQWO7+'SO&%xQWO'#EhO&&TQ^O'#EhOOQU'#Ho'#HoO&&TQ^O,5;ROOQV,5;R,5;RO&&_QWO,5;RO&&dQ^O,5;RO!0zQWO7+'_OOQV7+'a7+'aO&&qQWO7+'cO&&yQWO7+'cO&'QQWO7+$xO&)uQ!fO7+'fO&)|Q!fO7+'fOOQV7+(o7+(oO!1dQWO7+(oO&*TQYO,5<lO&*`QYO,5<lO!1dQWO'#GWO&*nQWO'#JpO&*|QWO'#G^O!BlQWO'#G^O&+RQWO'#JpOOQO'#Jo'#JoO&+ZQWO,5<wOOQO'#DX'#DXO&+`QYO'#JrO&,oQWO'#JrO$>PQYO'#JrOOQO'#Jq'#JqO&,zQWO,5<yO&-PQWO'#GZO#DQQWO'#G[O&-XQWO'#G[O&-aQWO'#JmOOQO'#Jl'#JlO&-lQYO'#GTOOQO,5<s,5<sO&-qQWO7+'pO&-vQWO'#JtO&.UQWO'#GeO#BlQWO'#GeO&.gQWO'#JtOOQO'#Js'#JsO&.oQWO,5=OO$>PQYO'#GUO&.tQYO'#JkOOQQ,5<n,5<nO&/]QWO7+(WOOQV7+(e7+(eO&/eQ^O'#D|O&0kQWO'#GlO&0sQ^O'#JwOOQO'#Gn'#GnO&0zQWO'#JwOOQO'#Jv'#JvO&1SQWO,5=VO&1XQWO'#IaO&1iQ^O'#GmO&2lQWO'#IrO&2zQWO'#GmOOQV7+(Y7+(YO&3SQWO7+(YO$>PQYO7+(YO&3[QYO'#HxO&3pQYO1G2WOOQQ1G2W1G2WOOQQ,5<m,5<mO$>PQYO,5<qO&3xQWO,5<rO&3}QWO7+(bO&4YQWO7+(fO&4aQWO7+(fOOQV7+(f7+(fO%?pQ_O7+(fO$>PQYO7+(fO&4lQWO'#IRO&4vQWO,5@hOOQO1G3Q1G3QOOQO1G2}1G2}OOQO1G3P1G3POOQO1G3R1G3ROOQO1G3S1G3SOOQO1G3O1G3OO&5OQWO7+(pO$>PQYO,59fO&5ZQ^O'#ISO&6QQYO,5?ROOQR1G/P1G/POOQV1G0[1G0[OOQS-E;l-E;lO&6YQ!bO,5:rO&6_Q!fO,5:tOOQV1G0]1G0]OOQV1G0Z1G0ZOOQO1G.o1G.oO&6fQWO'#KTOOQO'#KS'#KSO&6nQWO1G.oOOQV<<Gr<<GrO&6sQWO1G5iO&6{Q_O,5>]O&9QQWO,5>]OOQO-E;o-E;oOOQO<<Jj<<JjO&9[QWO<<JjOOQO<<Jk<<JkO&9cQXO7+%`O&:dQWO,5>WOOQO-E;j-E;jOOQP7+%a7+%aO!1PQ^O,5:`O&:rQWO'#HmO&;WQWO,5?hOOQP1G/y1G/yOOQO,5:`,5:`O&;`QWO,5:`O%ESQWO,5:`O$>PQYO,5<PO&;eQXO,5<PO&;sQXO7+'UO%?pQ_O1G1pO&<OQWO1G1pOOQO,5>`,5>`OOQO-E;r-E;rOOQV7+'l7+'lO&<YQWO<<KtO#DYQWO<<KtO&<hQWO<<HdOOQV<<LZ<<LZO!1dQWO<<LZOOQV<<K[<<K[O&<sQWO<<K[O%/kQWO<<K[O&<xQWO<<K[OOQV<<Kl<<KlO%/kQWO<<KlOOQV<<Kn<<KnO%0ZQWO<<KnO&=QQWO<<KrO$>PQYO<<KrOOQV<<K}<<K}O%5hQWO<<K}O%5mQWO<<K}O#G{QWO<<K}OOQV<<Kt<<KtO&=YQWO<<KtO$>PQYO<<KtO&=bQWO<<L[O$>PQYO<<L[O&=mQWO<<L[OOQV<<Hb<<HbO$>PQYO7+(uOOQO7+*V7+*VOOQR7+$i7+$iO&=rQWO,5@lOOQO'#Gm'#GmO&=zQWO'#GmO&>VQYO'#IUO&=rQWO,5@lOOQR1G3]1G3]O&?rQYO,5=vO&ARQYO,5=XO&A]QWO,5=XOOQO,5=X,5=XOOQR7+(u7+(uO&AbQZO7+(uO&CtQZO7+(|O&FTQWO,5>qOOQO-E<T-E<TO&F`QWO7+)QOOQO<<G}<<G}O&FgQYO'#ITO&FrQYO,5@jOOQQ7+(v7+(vOOQQ<<He<<HeO$>PQYO<<K|OOQV<<K|<<K|O&3}QWO<<K|O&FzQWO<<K|O%5hQWO<<K|O&GSQWO7+%YOOQV<<Hr<<HrOOQO<<Jn<<JnO%?pQ_O,5;SO&GZQWO,5;SO%?pQ_O'#EjO&G`QWO,5;SOOQU-E;m-E;mO&GkQWO1G0mOOQV1G0m1G0mO&&TQ^O1G0mOOQV<<Jy<<JyO!.cQ_O<<J}OOQV<<J}<<J}OOQV<<Hd<<HdO%?pQ_O<<HdO&GpQWO'#JTO&GxQWO'#FvO&G}QWO<<KQO&HVQ!fO<<KQO&H^QWO<<KQO&HcQWO<<KQO&HkQ!fO<<KQOOQV<<KQ<<KQO&HrQWO<<LZO&HwQWO,5@[O$>PQYO,5<xO&IPQWO,5<xO&IUQWO'#H{O&HwQWO,5@[OOQV1G2c1G2cO&IjQWO,5@^O$>PQYO,5@^O&IuQYO'#H|O&K[QWO,5@^OOQO1G2e1G2eO%,bQWO,5<uOOQO,5<v,5<vO&KdQYO'#HzO&LvQWO,5@XO%,UQYO,5=pO$>PQYO,5<oO&MRQWO,5@`O%?pQ_O,5=PO&MZQWO,5=PO&MfQWO,5=PO&MwQWO'#H}O&MRQWO,5@`OOQV1G2j1G2jO&N]QYO,5<pO%0`QYO,5>PO&NtQYO,5@VOOQV<<Kr<<KrO' ]QWO,5=XO' mQ^O,5:hO'!pQWO,5=XO$>PQYO,5=WO'!xQWO,5@cO'#QQWO,5@cO'#`Q^O'#IPO'!xQWO,5@cOOQO1G2q1G2qO'$rQWO,5=WO'$zQWO<<KtO'%YQYO,5>oO'%eQYO,5>dO'%sQYO,5>dOOQQ,5>d,5>dOOQQ-E;v-E;vOOQQ7+'r7+'rO'&OQYO1G2]O$>PQYO1G2^OOQV<<LQ<<LQO%?pQ_O<<LQO'&ZQWO<<LQO'&bQWO<<LQOOQO,5>m,5>mOOQO-E<P-E<POOQV<<L[<<L[O%?pQ_O<<L[O'&mQYO1G/QO'&xQYO,5>nOOQQ,5>n,5>nO''TQYO,5>nOOQQ-E<Q-E<QOOQS1G0^1G0^O')cQ!fO1G0`O')pQ!fO1G0`O')wQ^O'#IWO'*eQWO,5@oOOQO7+$Z7+$ZO'*mQWO1G3wOOQOAN@UAN@UO'*wQWO1G/zOOQO,5>X,5>XOOQO-E;k-E;kO!1PQ^O1G/zOOQO1G/z1G/zO'+SQWO1G/zO'+XQXO1G1kO$>PQYO1G1kO'+dQWO7+'[OOQVANA`ANA`O'+nQWOANA`O$>PQYOANA`O'+vQWOANA`OOQVAN>OAN>OO%?pQ_OAN>OO',UQWOANAuOOQVAN@vAN@vO',ZQWOAN@vOOQVANAWANAWOOQVANAYANAYOOQVANA^ANA^O',`QWOANA^OOQVANAiANAiO%5hQWOANAiO%5mQWOANAiO',hQWOANA`OOQVANAvANAvO%?pQ_OANAvO',vQWOANAvO$>PQYOANAvOOQR<<La<<LaO'-RQWO1G6WO%JsQWO,5>pOOQO'#HY'#HYO'-ZQWO'#HZOOQO,5>p,5>pOOQO-E<S-E<SO'-fQYO1G2sO'-pQWO1G2sOOQO1G2s1G2sO$>PQYO<<LaOOQR<<Ll<<LlOOQQ,5>o,5>oOOQQ-E<R-E<RO&3}QWOANAhOOQVANAhANAhO%5hQWOANAhO$>PQYOANAhO'-uQWO1G1rO'.iQ^O1G0nO%?pQ_O1G0nO'0_QWO,5;UO'0fQWO1G0nP'0kQWO'#ERP&&TQ^O'#HpOOQV7+&X7+&XO'0vQWO7+&XO&&yQWOAN@iO'0{QWOAN>OO!5oQWO,5<bOOQS,5>a,5>aO'1SQWOAN@lO'1XQWOAN@lOOQS-E;s-E;sOOQVAN@lAN@lO'1aQWOAN@lOOQVANAuANAuO'1iQWO1G5vO'1qQWO1G2dO$>PQYO1G2dO&*nQWO,5>gOOQO,5>g,5>gOOQO-E;y-E;yO'1|QWO1G5xO'2UQWO1G5xO&+`QYO,5>hO'2aQWO,5>hO$>PQYO,5>hOOQO-E;z-E;zO'2lQWO'#JnOOQO1G2a1G2aOOQO,5>f,5>fOOQO-E;x-E;xO&*TQYO,5<lO'2zQYO1G2ZO'3fQWO1G5zO'3nQWO1G2kO%?pQ_O1G2kO'3xQWO1G2kO&-vQWO,5>iOOQO,5>i,5>iOOQO-E;{-E;{OOQQ,5>c,5>cOOQQ-E;u-E;uO'4TQWO1G2sO'4eQWO1G2rO'4pQWO1G5}O'4xQ^O,5>kOOQO'#Go'#GoOOQO,5>k,5>kO'6UQWO,5>kOOQO-E;}-E;}O$>PQYO1G2rO'6dQYO7+'xO'6oQWOANAlOOQVANAlANAlO%?pQ_OANAlO'6vQWOANAvOOQS7+%z7+%zO'6}QWO7+%zO'7YQ!fO7+%zOOQO,5>r,5>rOOQO-E<U-E<UO'7gQWO7+%fO!1PQ^O7+%fO'7rQXO7+'VOOQVG26zG26zO'7}QWOG26zO'8]QWOG26zO$>PQYOG26zO'8eQWOG23jOOQVG27aG27aOOQVG26bG26bOOQVG26xG26xOOQVG27TG27TO%5hQWOG27TO'8lQWOG27bOOQVG27bG27bO%?pQ_OG27bO'8sQWOG27bOOQO1G4[1G4[OOQO7+(_7+(_OOQRANA{ANA{OOQVG27SG27SO%5hQWOG27SO&3}QWOG27SO'9OQ^O7+&YO':iQWO7+'^O';]Q^O7+&YO%?pQ_O7+&YP%?pQ_O,5;SP'<iQWO,5;SP'<nQWO,5;SOOQV<<Is<<IsOOQVG26TG26TOOQVG23jG23jOOQO1G1|1G1|OOQVG26WG26WO'<yQWOG26WP&HfQWO'#HuO'=OQWO7+(OOOQO1G4R1G4RO'=ZQWO7++dO'=cQWO1G4SO$>PQYO1G4SO%,bQWO'#HyO'=nQWO,5@YO'=|QWO7+(VO%?pQ_O7+(VOOQO1G4T1G4TOOQO1G4V1G4VO'>WQWO1G4VO'>fQWO7+(^OOQVG27WG27WO'>qQWOG27WOOQS<<If<<IfO'>xQWO<<IfO'?TQWO<<IQOOQVLD,fLD,fO'?`QWOLD,fO'?hQWOLD,fOOQVLD)ULD)UOOQVLD,oLD,oOOQVLD,|LD,|O'?vQWOLD,|O%?pQ_OLD,|OOQVLD,nLD,nO%5hQWOLD,nO'?}Q^O<<ItO'AhQWO<<JxO'B[Q^O<<ItP'ChQWO1G0nP'DXQ^O1G0nP%?pQ_O1G0nP'EzQWO1G0nOOQVLD+rLD+rO'FPQWO7+)nOOQO,5>e,5>eOOQO-E;w-E;wO'F[QWO<<KqOOQVLD,rLD,rOOQSAN?QAN?QOOQV!$(!Q!$(!QO'FfQWO!$(!QOOQV!$(!h!$(!hO'FnQWO!$(!hOOQV!$(!Y!$(!YO'FuQ^OAN?`POQU7+&Y7+&YP'H`QWO7+&YP'IPQ^O7+&YP%?pQ_O7+&YOOQV!)9El!)9ElOOQV!)9FS!)9FSPOQU<<It<<ItP'JrQWO<<ItP'KcQ^O<<ItPOQUAN?`AN?`O'MUQWO'#CnO'M]QXO'#CnO'NUQWO'#IaO( kQXO'#IaO(!bQWO'#DpO(!bQWO'#DpO!.cQ_O'#EkO(!sQ_O'#EoO(!zQ_O'#FPO(%{Q_O'#FbO(&SQXO'#IaO(&yQ_O'#E}O('|Q_O'#FWO(!bQWO,5:[O(!bQWO,5:[O!.cQ_O,5;ZO!.cQ_O,5;]O!.cQ_O,5;]O!.cQ_O,5;]O!.cQ_O,5;]O!.cQ_O,5;]O!.cQ_O,5;]O!.cQ_O,5;]O!.cQ_O,5;]O!.cQ_O,5;]O!.cQ_O,5;fO()PQ_O,5;kO(,QQWO,5;kO(,bQWO,5;|O(,iQYO'#CuO(,tQYO'#CvO(-PQWO'#CwO(-PQWO'#CwO(-bQYO'#CtO(-mQWO,5;iO(-tQWO,5;rO(-{Q_O,5;rO(/RQ_O,5;|O(!bQWO1G/vO(/YQWO1G0uO(0wQWO1G0wO(1RQWO1G0wO(2vQWO1G0wO(2}QWO1G0wO(4oQWO1G0wO(4vQWO1G0wO(6hQWO1G0wO(6oQWO1G0wO(6vQWO1G1QO(7WQWO1G1VO(7hQYO'#IfO(-PQWO,59cO(-PQWO,59cO(7sQWO1G1^O(7zQWO1G1hO(-PQWO1G.}O(8RQWO'#DpO!.^QXO'#FbO(8WQWO,5;ZO(8_QWO'#Cw",
  stateData: "(8q~O&}OSUOS'OPQ~OPoOQ!QOSVOTVOZeO[lO^RO_RO`ROa!UOd[Og!nOsVOtVOuVOw!POyvO|!VO}mO!Q!dO!U!WO!W!XO!X!^O!Z!YO!]!pO!liO!qgO!tiO#Y!_O#r!ZO#{![O$O!]O$b!`O$d!bO$f!cO$i!eO$m!fO$q!gO$s!hO%T!iO%V!jO%Z!kO%]!lO%^!mO%f!oO%j!qO%s!rO'R`O'UQO'[kO'_UO'hcO'riO(QdO~O'O!sO~OZbX[bXdbXdlXobXwjX}bX!lbX!qbX!tbX#QbX#RbX#pbX'hbX'rbX'sbX'xbX'ybX'zbX'{bX'|bX'}bX(ObX(PbX(QbX(RbX(TbX~OybXXbX!ebX!PbXvbX#TbX~P$|OZ'TX['TXd'TXd'YXo'TXw'lXy'TX}'TX!l'TX!q'TX!t'TX#Q'TX#R'TX#p'TX'h'TX'r'TX's'TX'x'TX'y'TX'z'TX'{'TX'|'TX'}'TX(O'TX(P'TX(Q'TX(R'TX(T'TXv'TX~OX'TX!e'TX!P'TX#T'TX~P'ZOr!uO'^!wO'`!uO~Od!xO~O^RO_RO`ROaRO'UQO~Od!}O~Od#PO[(SXo(SXy(SX}(SX!l(SX!q(SX!t(SX#Q(SX#R(SX#p(SX'h(SX'r(SX's(SX'x(SX'y(SX'z(SX'{(SX'|(SX'}(SX(O(SX(P(SX(Q(SX(R(SX(T(SXv(SX~OZ#OO~P*`OZ#RO[#QO~OQ!QO^#TO_#TO`#TOa#]Od#ZOg!nOyvO|!VO!Q!dO!U#^O!W!lO!]!pO$i!eO$m!fO$q!gO$s!hO%T!iO%V!jO%Z!kO%]!lO%^!mO%f!oO%j!qO%s!rO'R#VO'U#SO~OPoOQ!QOSVOTVOZeO[lOd[OsVOtVOuVOw!PO}mO!U#bO!W#cO!X!^O!Z!YO!liO!qgO!tiO#Y!_O#r!ZO#{![O$O!]O$b!`O$d!bO$f!cO'[kO'_UO'hcO'riO(QdO~P)xOPoOQ!QOSVOTVOZeO[lOd[OsVOtVOuVOw!PO}mO!U#bO!W#cO!X!^O!Z!YO!j#eO!liO!qgO!tiO#Y!_O#r!ZO#{![O$O!]O$b!`O$d!bO$f!cO'[kO'_UO'hcO'riO(QdO~P)xO[#}Oo#xO}#zO!l#yO!q#jO!t#yO#Q#xO#R#uO#p$OO'h#gO'r#yO's#lO'x#hO'y#iO'z#iO'{#kO'|#nO'}#mO(O#|O(P#gO(Q#hO(R#fO(T#hO~OPoOQ!QOSVOTVOZeOd[OsVOtVOuVOw!PO!U#bO!W#cO!X!^O!Z!YO#Y!_O#r!ZO#{![O$O!]O$b!`O$d!bO$f!cO'[kO'_UO[#sXo#sXy#sX}#sX!l#sX!q#sX!t#sX#Q#sX#R#sX#p#sX'h#sX'r#sX's#sX'x#sX'y#sX'z#sX'{#sX'|#sX'}#sX(O#sX(P#sX(Q#sX(R#sX(T#sXX#sX!e#sX!P#sXv#sX#T#sX~P)xOX(SX!e(SX!P(SXw(SX#T(SX~P*`OPoOQ!QOSVOTVOX$ROZeO[lOd[OsVOtVOuVOw!PO}mO!U#bO!W#cO!X!^O!Z!YO!liO!qgO!tiO#Y!_O#r!ZO#{![O$O!]O$b!`O$d!bO$f!cO'R$UO'[kO'_UO'hcO'riO(QdO~P)xOPoOQ!QOSVOTVOZeO[lOd[OsVOtVOuVOw!PO}mO!P$XO!U#bO!W#cO!X!^O!Z!YO!liO!qgO!tiO#Y!_O#r!ZO#{![O$O!]O$b!`O$d!bO$f!cO'R$UO'[kO'_UO'hcO'riO(QdO~P)xOQ!QOSVOTVO[$gO^$pO_$ZO`:QOa:QOd$aOsVOtVOuVO}$eO!i$qO!l$lO!q$hO#V$lO'U$YO'_UO'h$[O~O!j$rOP(XP~P<cOPoOQ!QOSVOTVOZeO[lOd[OsVOtVOuVOw!PO}mO!U#bO!W#cO!X!^O!Z!YO!liO!qgO!tiO#S$uO#Y!_O#r!ZO#{![O$O!]O$b!`O$d!bO$f!cO'[kO'_UO'hcO'riO(QdO~P)xOw$vO~Oo'cX#Q'cX#R'cX#p'cX's'cX'x'cX'y'cX'z'cX'{'cX'|'cX'}'cX(O'cX(P'cX(R'cX(T'cX~OP%tXQ%tXS%tXT%tXZ%tX[%tX^%tX_%tX`%tXa%tXd%tXg%tXs%tXt%tXu%tXw%tXy%tX|%tX}%tX!Q%tX!U%tX!W%tX!X%tX!Z%tX!]%tX!l%tX!q%tX!t%tX#Y%tX#r%tX#{%tX$O%tX$b%tX$d%tX$f%tX$i%tX$m%tX$q%tX$s%tX%T%tX%V%tX%Z%tX%]%tX%^%tX%f%tX%j%tX%s%tX&{%tX'R%tX'U%tX'[%tX'_%tX'h%tX'r%tX(Q%tXv%tX~P@[Oy$xO['cX}'cX!l'cX!q'cX!t'cX'h'cX'r'cX(Q'cXv'cX~P@[Ow$yO!Q(iX!U(iX!W(iX$q(iX%](iX%^(iX~Oy$zO~PEsO!Q$}O!U%UO!W!lO$m%OO$q%PO$s%QO%T%RO%V%SO%Z%TO%]!lO%^%VO%f%WO%j%XO%s%YO~O!Q!lO!U!lO!W!lO$q%[O%]!lO~O%^%VO~PGaOPoOQ!QOSVOTVOZeO[lO^RO_RO`ROa!UOd[Og!nOsVOtVOuVOw!POyvO|!VO}mO!Q!dO!U!WO!W!XO!X!^O!Z!YO!]!pO!liO!qgO!tiO#Y!_O#r!ZO#{![O$O!]O$b!`O$d!bO$f!cO$i!eO$m!fO$q!gO$s!hO%T!iO%V!jO%Z!kO%]!lO%^!mO%f!oO%j!qO%s!rO'R#VO'UQO'[kO'_UO'hcO'riO(QdO~Ov%`O~P]OQ!QOZ%rO[%qO^%vO_%cO`TOaTOd%jOg%yO}%pO!q%oO$f%wO%^%xO&W%{O'U%dO'[%eO(Q%zO~PGaO!Q{X!U{X!W{X$m{X$q{X$s{X%T{X%V{X%Z{X%]{X%^{X%f{X%j{X%s{X~P'ZO!Q{X!U{X!W{X$m{X$q{X$s{X%T{X%V{X%Z{X%]{X%^{X%f{X%j{X%s{X~O}%}O'U{XQ{XZ{X[{X^{X_{X`{Xa{Xd{Xg{X!q{X$f{X&W{X'[{X(Q{X~PMuOg&PO%f%WO!Q(iX!U(iX!W(iX$q(iX%](iX%^(iX~Ow!PO~P! yOw!PO!X&RO~PEvOPoOQ!QOSVOTVOZeO[lO^9xO_9xO`9xOa9xOd9{OsVOtVOuVOw!PO}mO!U#bO!W#cO!X;RO!Z!YO!]&UO!l:OO!q9}O!t:OO#Y!_O#r:RO#{:SO$O!]O$b!`O$d!bO$f!cO'U9vO'[kO'_UO'hcO'r:OO(QdO~OPoOQ!QOSVOTVOZeO[lOd[OsVOtVOuVOw!PO}mO!U#bO!W#cO!X!^O!Z!YO!liO!qgO!tiO#Y!_O#r!ZO#{![O$O!]O$b!`O$d!bO$f!cO'[kO'_UO'hcO'riO(QdOo#qXy#qX#Q#qX#R#qX#p#qX's#qX'x#qX'y#qX'z#qX'{#qX'|#qX'}#qX(O#qX(P#qX(R#qX(T#qXX#qX!e#qX!P#qXv#qX#T#qX~P)xOPoOQ!QOSVOTVOZeO[lOd[OsVOtVOuVOw!PO}mO!U#bO!W#cO!X!^O!Z!YO!liO!qgO!tiO#Y!_O#r!ZO#{![O$O!]O$b!`O$d!bO$f!cO'[kO'_UO'hcO'riO(QdOo#zXy#zX#Q#zX#R#zX#p#zX's#zX'x#zX'y#zX'z#zX'{#zX'|#zX'}#zX(O#zX(P#zX(R#zX(T#zXX#zX!e#zX!P#zXv#zX#T#zX~P)xO'[kO[#}Xo#}Xy#}X}#}X!l#}X!q#}X!t#}X#Q#}X#R#}X#p#}X'h#}X'r#}X's#}X'x#}X'y#}X'z#}X'{#}X'|#}X'}#}X(O#}X(P#}X(Q#}X(R#}X(T#}XX#}X!e#}X!P#}Xv#}Xw#}X#T#}X~OPoO~OPoOQ!QOSVOTVOZeO[lO^9xO_9xO`9xOa9xOd9{OsVOtVOuVOw!PO}mO!U#bO!W#cO!X;RO!Z!YO!l:OO!q9}O!t:OO#Y!_O#r:RO#{:SO$O!]O$b!`O$d!bO$f!cO'U9vO'[kO'_UO'hcO'r:OO(QdO~O!S&_O~Ow!PO~O!j&bO~P<cO'U&cO~PEvOZ&eO~O'U&cO~O'_UOw(^Xy(^X!Q(^X!U(^X!W(^X$q(^X%](^X%^(^X~Oa&hO~P!1iO'U&iO~O_&nO'U&cO~OQ&oOZ&pO[%qO^%vO_%cO`TOaTOd%jOg%yO}%pO!q%oO$f%wO%^%xO&W%{O'U%dO'[%eO(Q%zO~PGaO!j&uO~P<cO^&wO_&wO`&wOa&wOd'POw&|O'U&vO(Q&}O~O!i'UO!j'TO'U&cO~O'O!sO'P'VO'Q'XO~Or!uO'^'ZO'`!uO~OQ']O^'ja_'ja`'jaa'ja'U'ja~O['bOw'cO}'dO~OQ']O~OQ!QO^#TO_#TO`#TOa'jOd#ZO'U#SO~O['kO~OZbXdlXXbXobXPbX!SbX!ebX'sbX!PbX!ObXybX!ZbX#TbXvbX~O}bX~P!6mOZ'TXd'YXX'TXo'TX}'TX#p'TXP'TX!S'TX!e'TX's'TX!P'TX!O'TXy'TX!Z'TX#T'TXv'TX~O^#TO_#TO`#TOa'jO'U#SO~OZ'lO~Od'nO~OZ'TXd'YX~PMuOZ'oOX(SX!e(SX!P(SXw(SX#T(SX~P*`O[#}O}#zO(O#|O(R#fOo#_ay#_a!l#_a!q#_a!t#_a#Q#_a#R#_a#p#_a'h#_a'r#_a's#_a'x#_a'y#_a'z#_a'{#_a'|#_a'}#_a(P#_a(Q#_a(T#_aX#_a!e#_a!P#_av#_aw#_a#T#_a~Ow!PO!X&RO~Oy#caX#ca!e#ca!P#cav#ca#T#ca~P2gOPoOQ!QOSVOTVOZeOd[OsVOtVOuVOw!PO!U#bO!W#cO!X!^O!Z!YO#Y!_O#r!ZO#{![O$O!]O$b!`O$d!bO$f!cO'[kO'_UO[#sao#say#sa}#sa!l#sa!q#sa!t#sa#Q#sa#R#sa#p#sa'h#sa'r#sa's#sa'x#sa'y#sa'z#sa'{#sa'|#sa'}#sa(O#sa(P#sa(Q#sa(R#sa(T#saX#sa!e#sa!P#sav#sa#T#sa~P)xOPoOQ!QOSVOTVOZeO[lOd[OsVOtVOuVOw!PO}mO!U#bO!W#cO!X!^O!Z!YO!liO!qgO!tiO#Y!_O#r!ZO#{![O$O!]O$b!`O$d!bO$f!cO'R#VO'[kO'_UO'hcO'riO(QdO!P(UP~P)xOu(RO#w(SO'U(QO~O[#}O}#zO!q#jO'h#gO's#lO'x#hO'y#iO'z#iO'{#kO'|#nO'}#mO(O#|O(P#gO(Q#hO(R#fO(T#hO!l#sa!t#sa#p#sa'r#sa~Oo#xO#Q#xO#R#uOy#saX#sa!e#sa!P#sav#sa#T#sa~P!BqOy(XO!e(VOX(WX~P2gOX(YO~OPoOQ!QOSVOTVOX(YOZeO[lOd[OsVOtVOuVOw!PO}mO!U#bO!W#cO!X!^O!Z!YO!liO!qgO!tiO#Y!_O#r!ZO#{![O$O!]O$b!`O$d!bO$f!cO'R$UO'[kO'_UO'hcO'riO(QdO~P)xOZ#RO~O!P(^O!e(VO~P2gOPoOQ!QOSVOTVOZeO[lOd[OsVOtVOuVOw!PO}mO!U#bO!W#cO!X!^O!Z!YO!liO!qgO!tiO#Y!_O#r!ZO#{![O$O!]O$b!`O$d!bO$f!cO'R$UO'[kO'_UO'hcO'riO(QdO~P)xOZbXdlXwjX}jX!tbX'rbX~OP!RX!S!RX!e!RX'q!RX's!RX!O!RXo!RXy!RX!P!RXX!RX!Z!RX#T!RXv!RX~P!IxOZ'TXd'YXw'lX}'lX!t'TX'r'TX~OP!`X!S!`X!e!`X's!`X!O!`Xo!`Xy!`X!P!`XX!`X!Z!`X#T!`Xv!`X~P!KZOT(`Ou(`O~O!t(aO'r(aOP!^X!S!^X!e!^X's!^X!O!^Xo!^Xy!^X!P!^XX!^X!Z!^X#T!^Xv!^X~O^9yO_9yO`:QOa:QO'U9wO~Od(dO~O'q(eOP'iX!S'iX!e'iX's'iX!O'iXo'iXy'iX!P'iXX'iX!Z'iX#T'iXv'iX~O!j&bO!P'mP~P<cOw(jO}(iO~O!j&bOX'mP~P<cO!j(nO~P<cOZ'oO!t(aO'r(aO~O!S(pO's(oOP$WX!e$WX~O!e(qOP(YX~OP(sO~OP!aX!S!aX!e!aX's!aX!O!aXo!aXy!aX!P!aXX!aX!Z!aX#T!aXv!aX~P!KZOy$UaX$Ua!e$Ua!P$Uav$Ua#T$Ua~P2gO!l({O'R#VO'U(wOv(ZP~OQ!QO^#TO_#TO`#TOa#]Od#ZOg!nOyvO|!VO!Q!dO!U#^O!W!lO!]!pO$i!eO$m!fO$q!gO$s!hO%T!iO%V!jO%Z!kO%]!lO%^!mO%f!oO%j!qO%s!rO'R`O'U#SO~Ov)SO~P#$]Oy)UO~PEsO%^)VO~PGaOa)YO~P!1iO%f)_O~PEvO_)`O'U&cO~O!i)eO!j)dO'U&cO~O'_UO!Q(^X!U(^X!W(^X$q(^X%](^X%^(^X~Ov%uX~P2gOv)fO~PGyOv)fO~Ov)fO~P]OQiXQ'YXZiXd'YX}iX#piX(PiX~ORiXwiX$fiX$|iX[iXoiXyiX!liX!qiX!tiX#QiX#RiX'hiX'riX'siX'xiX'yiX'ziX'{iX'|iX'}iX(OiX(QiX(RiX(TiX!PiX!eiXXiXPiXviX!SiX#TiX~P#(_OQjXQlXRjXZjXdlX}jX#pjX(PjXwjX$fjX$|jX[jXojXyjX!ljX!qjX!tjX#QjX#RjX'hjX'rjX'sjX'xjX'yjX'zjX'{jX'|jX'}jX(OjX(QjX(RjX(TjX!PjX!ejXXjX!SjXPjXvjX#TjX~O%^)iO~PGaOQ']Od)jO~O^)lO_)lO`)lOa)lO'U%dO~Od)pO~OQ']OZ)tO})rOR'VX#p'VX(P'VXw'VX$f'VX$|'VX['VXo'VXy'VX!l'VX!q'VX!t'VX#Q'VX#R'VX'h'VX'r'VX's'VX'x'VX'y'VX'z'VX'{'VX'|'VX'}'VX(O'VX(Q'VX(R'VX(T'VX!P'VX!e'VXX'VXP'VXv'VX!S'VX#T'VX~OQ!QO^:iO_:eO`TOaTOd:hO%^)iO'U:fO~PGaOQ!QOZ%rO[%qO^%vO_%cO`TOaTOd%jOg%yO}%pO!j)xO!q%oO$f%wO%^%xO&W%{O'U%dO'[%eO(Q%zO~PGaOQ!QOZ%rO[%qO^%vO_%cO`TOaTOd%jOg%yO}%pO!P){O!q%oO$f%wO%^%xO&W%{O'U%dO'[%eO(Q%zO~PGaO(P)}O~OR*PO#p*QO(P*OO~OQhXQ'YXZhXd'YX}hX(PhX~ORhX#phXwhX$fhX$|hX[hXohXyhX!lhX!qhX!thX#QhX#RhX'hhX'rhX'shX'xhX'yhX'zhX'{hX'|hX'}hX(OhX(QhX(RhX(ThX!PhX!ehXXhXPhXvhX!ShX#ThX~P#4_OQ*RO~O})rO~OQ!QO^%vO_%cO`TOaTOd%jO$f%wO%^%xO'U%dO~PGaO!Q*UO!j*UO~O^*XO`*XOa*XO!O*YO~OQ&oOZ*ZO[%qO^%vO_%cO`TOaTOd%jOg%yO}%pO!q%oO$f%wO%^%xO&W%{O'U%dO'[%eO(Q%zO~PGaO[#}Oo:aO}#zO!l:bO!q#jO!t:bO#Q:aO#R:^O#p$OO'h#gO'r:bO's#lO'x#hO'y#iO'z#iO'{#kO'|#nO'}#mO(O#|O(P#gO(Q#hO(R#fO(T#hO~Ow'eX~P#9jOy#qaX#qa!e#qa!P#qav#qa#T#qa~P2gOy#zaX#za!e#za!P#zav#za#T#za~P2gOPoOQ!QOSVOTVOZeO[lOd[OsVOtVOuVOw!PO}mO!S&_O!U#bO!W#cO!X!^O!Z!YO!liO!qgO!tiO#Y!_O#r!ZO#{![O$O!]O$b!`O$d!bO$f!cO'[kO'_UO'hcO'riO(QdOo#zay#za#Q#za#R#za#p#za's#za'x#za'y#za'z#za'{#za'|#za'}#za(O#za(P#za(R#za(T#zaX#za!e#za!P#zav#za#T#za~P)xOPoOQ!QOSVOTVOZeO[lOd[OsVOtVOuVOw!PO}mO!U#bO!W#cO!X!^O!Z!YO!liO!qgO!tiO#S*dO#Y!_O#r!ZO#{![O$O!]O$b!`O$d!bO$f!cO'[kO'_UO'hcO'riO(QdO~P)xOw*eO~P#9jO$b*hO$d*iO$f*jO~O!O*kO's(oO~O!S*mO~O'U*nO~Ow$yOy*pO~O'U*qO~OQ*tOw*uOy*xO}*vO$|*wO~OQ*tOw*uO$|*wO~OQ*tOw+PO$|*wO~OQ*tOo+UOy+WO!S+TO~OQ*tO}+YO~OQ!QOZ%rO[%qO^%vO`TOaTOd%jOg%yO}%pO!U!lO!W!lO!q%oO$f%wO$q%[O%]!lO%^%xO&W%{O'U%dO'[%eO(Q%zO~OR+aO_+]O!Q+bO~P#D_O_%cO!Q!lOw&UX$|&UX(P&UX~P#D_Ow$yO$f+gO$|*wO(P*OO~OQ!QOZ*ZO[%qO^%vO_%cO`TOaTOd%jOg%yO}%pO!q%oO$f%wO%^%xO&W%{O'U%dO'[%eO(Q%zO~PGaOQ*tOw$yO!S+TO$|*wO~Oo+mOy+lO!S+nO's(oO~OdlXy!RX#pbXv!RX!e!RX~Od'YXy(mX#p'TXv(mX!e(mX~Od+pO~O^#TO_#TO`#TOa'jOw&|O'U&vO(Q+uO~Ov(oP~P!3|O#p+zO~Oy+{O~O!S+|O~O'O!sO'P'VO'Q,OO~Od,PO~OSVOTVO_%cOsVOtVOuVOw!PO!Q!lO'_UO~P#D_OS,_OT,_OZ,_O['bO_,ZOd,_Oo,_Os,_Ou,_Ow'cOy,_O}'dO!S,_O!e,_O!l,_O!q,]O!t,_O!{,_O#Q,_O#R,_O#S,_O#T,_O'R,_O'[%eO'_UO'h,[O's,]O'v,`O'x,[O'y,]O'z,]O'{,]O'|,^O'},^O(O,_O(P,aO(Q,aO(R,bO~OX,XO~P#K_Ov,dO~P#K_O!P,gO~P#K_Oo'ti#Q'ti#R'ti#p'ti's'ti'x'ti'y'ti'z'ti'{'ti'|'ti'}'ti(O'ti(P'ti(R'ti(T'ti~Oy,hO['ti}'ti!l'ti!q'ti!t'ti'h'ti'r'ti(Q'tiv'ti~P#N^OP$giQ$giS$giT$giZ$gi[$gi^$gi_$gi`$gia$gid$gig$gis$git$giu$giw$giy$gi|$gi}$gi!Q$gi!U$gi!W$gi!X$gi!Z$gi!]$gi!l$gi!q$gi!t$gi#Y$gi#r$gi#{$gi$O$gi$b$gi$d$gi$f$gi$i$gi$m$gi$q$gi$s$gi%T$gi%V$gi%Z$gi%]$gi%^$gi%f$gi%j$gi%s$gi&{$gi'R$gi'U$gi'[$gi'_$gi'h$gi'r$gi(Q$giv$gi~P#N^OX,iO~Oo,jO},kOX]X!P]X!e]X~Oy#ciX#ci!e#ci!P#civ#ci#T#ci~P2gO[#}O}#zO'x#hO(O#|O(Q#hO(R#fO(T#hOo#eiy#ei!l#ei!q#ei!t#ei#Q#ei#R#ei#p#ei'r#ei's#ei'y#ei'z#ei'{#ei'|#ei'}#eiX#ei!e#ei!P#eiv#ei#T#ei~O'h#ei(P#ei~P$&sO[#}O}#zO(O#|O(R#fOo#eiy#ei!l#ei!q#ei!t#ei#Q#ei#R#ei#p#ei'r#ei's#ei'y#ei'z#ei'{#ei'|#ei'}#eiX#ei!e#ei!P#eiv#ei#T#ei~O'h#ei'x#ei(P#ei(Q#ei(T#eiw#ei~P$(tO'h#gO(P#gO~P$&sO[#}O}#zO'h#gO'x#hO'y#iO'z#iO(O#|O(P#gO(Q#hO(R#fO(T#hOo#eiy#ei!l#ei!t#ei#Q#ei#R#ei#p#ei'r#ei's#ei'{#ei'|#ei'}#eiX#ei!e#ei!P#eiv#ei#T#ei~O!q#ei~P$+SO!q#jO~P$+SO[#}O}#zO!q#jO'h#gO'x#hO'y#iO'z#iO'{#kO(O#|O(P#gO(Q#hO(R#fO(T#hOo#eiy#ei!l#ei!t#ei#Q#ei#R#ei#p#ei'r#ei'|#ei'}#eiX#ei!e#ei!P#eiv#ei#T#ei~O's#ei~P$-[O's#lO~P$-[O[#}O}#zO!q#jO#R#uO'h#gO's#lO'x#hO'y#iO'z#iO'{#kO(O#|O(P#gO(Q#hO(R#fO(T#hOo#eiy#ei!l#ei!t#ei#Q#ei#p#ei'r#ei'|#eiX#ei!e#ei!P#eiv#ei#T#ei~O'}#ei~P$/dO'}#mO~P$/dO[#}O}#zO!q#jO'h#gO's#lO'x#hO'y#iO'z#iO'{#kO'|#nO'}#mO(O#|O(P#gO(Q#hO(R#fO(T#hO!l#ni!t#ni#p#ni'r#ni~Oo#xO#Q#xO#R#uOy#niX#ni!e#ni!P#niv#ni#T#ni~P$1lO[#}O}#zO!q#jO'h#gO's#lO'x#hO'y#iO'z#iO'{#kO'|#nO'}#mO(O#|O(P#gO(Q#hO(R#fO(T#hO!l#si!t#si#p#si'r#si~Oo#xO#Q#xO#R#uOy#siX#si!e#si!P#siv#si#T#si~P$3mOPoOQ!QOSVOTVOZeO[lOd[OsVOtVOuVOw!PO}mO!U#bO!W#cO!X!^O!Z!YO!liO!qgO!tiO#Y!_O#r!ZO#{![O$O!]O$b!`O$d!bO$f!cO'R#VO'[kO'_UO'hcO'riO(QdO~P)xO!e,rO!P(VX~P2gO!P,tO~OX,uO~P2gOPoOQ!QOSVOTVOZeO[lOd[OsVOtVOuVOw!PO}mO!U#bO!W#cO!X!^O!Z!YO!liO!qgO!tiO#Y!_O#r!ZO#{![O$O!]O$b!`O$d!bO$f!cO'[kO'_UO'hcO'riO(QdOX&fX!e&fX!P&fX~P)xO!e(VOX(Wa~Oy,yO!e(VOX(WX~P2gOX,zO~O!P,{O!e(VO~O!P,}O!e(VO~P2gOSVOTVOsVOtVOuVO'_UO'h$[O~P!6POP!baZca!S!ba!e!ba!tca'rca's!ba!O!bao!bay!ba!P!baX!ba!Z!ba#T!bav!ba~O!e-SO's(oO!P'nXX'nX~O!P-UO~O!i-_O!j-^O!l-ZO'U-WOv'oP~OX-`O~O_%cO!Q!lO~P#D_O!j-fOP&gX!e&gX~P<cO!e(qOP(Ya~O!S-hO's(oOP$Wa!e$Wa~Ow!PO(P*OO~OvbX!S!kX!ebX~O'R#VO'U(wO~O!S-lO~O!e-nOv([X~Ov-pO~Ov-rO~P,cOv-rO~P#$]O_-tO'U&cO~O!S-uO~Ow$yOy-vO~OQ*tOw*uOy-yO}*vO$|*wO~OQ*tOo.TO~Oy.^O~O!S._O~O!j.aO'U&cO~Ov.bO~Ov.bO~PGyOQ']O^'Xa_'Xa`'Xaa'Xa'U'Xa~Od.fO~OQ'YXQ'lXR'lXZ'lXd'YX}'lX#p'lX(P'lXw'lX$f'lX$|'lX['lXo'lXy'lX!l'lX!q'lX!t'lX#Q'lX#R'lX'h'lX'r'lX's'lX'x'lX'y'lX'z'lX'{'lX'|'lX'}'lX(O'lX(Q'lX(R'lX(T'lX!P'lX!e'lXX'lXP'lXv'lX!S'lX#T'lX~OQ!QOZ%rO[%qO^.qO_%cO`TOaTOd%jOg%yO}%pO!j.rO!q.oO!t.jO#V.lO$f%wO%^%xO&W%{O'R#VO'U%dO'[%eO(Q%zO!P(sP~PGaO#S.sOR%wa#p%wa(P%waw%wa$f%wa$|%wa[%wao%way%wa}%wa!l%wa!q%wa!t%wa#Q%wa#R%wa'h%wa'r%wa's%wa'x%wa'y%wa'z%wa'{%wa'|%wa'}%wa(O%wa(Q%wa(R%wa(T%wa!P%wa!e%waX%waP%wav%wa!S%wa#T%wa~O%^.uO~PGaO(P*OOR&Oa#p&Oaw&Oa$f&Oa$|&Oa[&Oao&Oay&Oa}&Oa!l&Oa!q&Oa!t&Oa#Q&Oa#R&Oa'h&Oa'r&Oa's&Oa'x&Oa'y&Oa'z&Oa'{&Oa'|&Oa'}&Oa(O&Oa(Q&Oa(R&Oa(T&Oa!P&Oa!e&OaX&OaP&Oav&Oa!S&Oa#T&Oa~O_%cO!Q!lO!j.wO(P)}O~P#D_O!e.xO(P*OO!P(uX~O!P.zO~OX.{Oy.|O(P*OO~O'[%eOR(qP~OQ']O})rORfa#pfa(Pfawfa$ffa$|fa[faofayfa!lfa!qfa!tfa#Qfa#Rfa'hfa'rfa'sfa'xfa'yfa'zfa'{fa'|fa'}fa(Ofa(Qfa(Rfa(Tfa!Pfa!efaXfaPfavfa!Sfa#Tfa~OQ']O})rOR&Va#p&Va(P&Vaw&Va$f&Va$|&Va[&Vao&Vay&Va!l&Va!q&Va!t&Va#Q&Va#R&Va'h&Va'r&Va's&Va'x&Va'y&Va'z&Va'{&Va'|&Va'}&Va(O&Va(Q&Va(R&Va(T&Va!P&Va!e&VaX&VaP&Vav&Va!S&Va#T&Va~O!P/TO~Ow$yO$f/YO$|*wO(P*OO~OQ!QOZ/ZO[%qO^%vO_%cO`TOaTOd%jOg%yO}%pO!q%oO$f%wO%^%xO&W%{O'U%dO'[%eO(Q%zO~PGaOo/]O's(oO~O#W/^OP!YiQ!YiS!YiT!YiZ!Yi[!Yi^!Yi_!Yi`!Yia!Yid!Yig!Yio!Yis!Yit!Yiu!Yiw!Yiy!Yi|!Yi}!Yi!Q!Yi!U!Yi!W!Yi!X!Yi!Z!Yi!]!Yi!l!Yi!q!Yi!t!Yi#Q!Yi#R!Yi#Y!Yi#p!Yi#r!Yi#{!Yi$O!Yi$b!Yi$d!Yi$f!Yi$i!Yi$m!Yi$q!Yi$s!Yi%T!Yi%V!Yi%Z!Yi%]!Yi%^!Yi%f!Yi%j!Yi%s!Yi&{!Yi'R!Yi'U!Yi'[!Yi'_!Yi'h!Yi'r!Yi's!Yi'x!Yi'y!Yi'z!Yi'{!Yi'|!Yi'}!Yi(O!Yi(P!Yi(Q!Yi(R!Yi(T!YiX!Yi!e!Yi!P!Yiv!Yi!i!Yi!j!Yi#V!Yi#T!Yi~Oy#ziX#zi!e#zi!P#ziv#zi#T#zi~P2gOy$UiX$Ui!e$Ui!P$Uiv$Ui#T$Ui~P2gOv/dO!j&bO'R`O~P<cOw/mO}/lO~Oy!RX#pbX~Oy/nO~O#p/oO~OR+aO_+cO!Q/rO'U&iO'[%eO~Oa/yO|!VO'R#VO'U(QOv(cP~OQ!QOZ%rO[%qO^%vO_%cO`TOa/yOd%jOg%yO|!VO}%pO!q%oO$f%wO%^%xO&W%{O'R#VO'U%dO'[%eO(Q%zO!P(eP~PGaOQ!QOZ%rO[%qO^%vO_%cO`TOaTOd%jOg%yO}%pO!q%oO$f0UO%^%xO&W%{O'U%dO'[%eO(Q%zOw(`Py(`P~PGaOw*uO~Oy-yO$|*wO~Oa/yO|!VO'R#VO'U*nOv(gP~Ow+PO~OQ!QOZ%rO[%qO^%vO_%cO`TOaTOd%jOg%yO}%pO!q%oO$f0UO%^%xO&W%{O'U%dO'[%eO(Q%zO(R0_O~PGaOy0cO~OQ!QOSVOTVO[$gO^0kO_$ZO`:QOa:QOd$aOsVOtVOuVO}$eO!i$qO!j0lO!l$lO!q0dO!t0gO'R#VO'U$YO'[%eO'_UO'h$[O~O#V0mO!P(jP~P%1qOw!POy0oO#S0qO$|*wO~OR0tO!e0rO~P#(_OR0tO!S+TO!e0rO(P)}O~OR0tOo0vO!S+TO!e0rOQ'WXZ'WX}'WX#p'WX(P'WX~OR0tOo0vO!e0rO~OR0tO!e0rO~O$f/YO(P*OO~Ow$yO~Ow$yO$|*wO~Oo0|Oy0{O!S0}O's(oO~O!e1OOv(pX~Ov1QO~O^#TO_#TO`#TOa'jOw&|O'U&vO(Q1UO~Oo1XOQ'WXR'WXZ'WX}'WX!e'WX(P'WX~O!e1YO(P*OOR'ZX~O!e1YOR'ZX~O!e1YO(P)}OR'ZX~OR1[O~OX1]O~P#K_O!S1_OS'wXT'wXX'wXZ'wX['wX_'wXd'wXo'wXs'wXu'wXw'wXy'wX}'wX!e'wX!l'wX!q'wX!t'wX!{'wX#Q'wX#R'wX#S'wX#T'wX'R'wX'['wX'_'wX'h'wX's'wX'v'wX'x'wX'y'wX'z'wX'{'wX'|'wX'}'wX(O'wX(P'wX(Q'wX(R'wXv'wX!P'wX~O}1`O~Ov1aO~P#K_O!P1bO~P#K_OSVOTVOsVOtVOuVO'_UO~OSVOTVOsVOtVOuVO'_UO!P(vP~P!6POX1gO~Oy,hO~O!e,rO!P(Va~P2gOPoOQ!QOZeO[lO^RO_RO`ROaROd[Ow!PO}mO!U#bO!W#cO!X!^O!Z!YO!liO!qgO!tiO#Y!_O#r!ZO#{![O$O!]O$b!`O$d!bO$f!cO'R#VO'UQO'[kO'hcO'riO(QdO!P&eX!e&eX~P%;dO!e,rO!P(Va~OX&fa!e&fa!P&fa~P2gOX1lO~P2gOPoOQ!QOZeO[lO^RO_RO`ROaROd[Ow!PO}mO!U#bO!W#cO!X!^O!Z!YO!liO!qgO!tiO#Y!_O#r!ZO#{![O$O!]O$b!`O$d!bO$f!cO'UQO'[kO'hcO'riO(QdO~P%;dO!P1nO!e(VO~OP!biZci!S!bi!e!bi!tci'rci's!bi!O!bio!biy!bi!P!biX!bi!Z!bi#T!biv!bi~O's(oOP!oi!S!oi!e!oi!O!oio!oiy!oi!P!oiX!oi!Z!oi#T!oiv!oi~O!j&bO!P&`X!e&`XX&`X~P<cO!e-SO!P'naX'na~O!P1rO~Ov!RX!S!kX!e!RX~O!S1sO~O!e1tOv'pX~Ov1vO~O'U-WO~O!j1yO'U-WO~O(P*OOP$Wi!e$Wi~O!S1zO's(oOP$XX!e$XX~O!S1}O~Ov$_a!e$_a~P2gO!l({O'R#VO'U(wOv&hX!e&hX~O!e-nOv([a~Ov2RO~P,cOy2VO~O#p2WO~Oy2XO$|*wO~Ow*uOy2XO}*vO$|*wO~Oo2bO~Ow!POy2gO#S2iO$|*wO~O!S2kO~Ov2mO~O#S2nOR%wi#p%wi(P%wiw%wi$f%wi$|%wi[%wio%wiy%wi}%wi!l%wi!q%wi!t%wi#Q%wi#R%wi'h%wi'r%wi's%wi'x%wi'y%wi'z%wi'{%wi'|%wi'}%wi(O%wi(Q%wi(R%wi(T%wi!P%wi!e%wiX%wiP%wiv%wi!S%wi#T%wi~Od2oO~O^2rO!j.rO!q2sO'R#VO'[%eO~O(P*OO!P%{X!e%{X~O!e2tO!P(tX~O!P2vO~OQ!QOZ%rO[%qO^2xO_%cO`TOaTOd%jOg%yO}%pO!j2yO!q%oO$f%wO%^%xO&W%{O'U%dO'[%eO(Q%zO~PGaO^2zO!j2yO(P)}O~O!P%aX!e%aX~P#4_O^2zO~O(P*OOR&Oi#p&Oiw&Oi$f&Oi$|&Oi[&Oio&Oiy&Oi}&Oi!l&Oi!q&Oi!t&Oi#Q&Oi#R&Oi'h&Oi'r&Oi's&Oi'x&Oi'y&Oi'z&Oi'{&Oi'|&Oi'}&Oi(O&Oi(Q&Oi(R&Oi(T&Oi!P&Oi!e&OiX&OiP&Oiv&Oi!S&Oi#T&Oi~O_%cO!Q!lO!P&yX!e&yX~P#D_O!e.xO!P(ua~OR3RO(P*OO~O!e3SOR(rX~OR3UO~O(P*OOR&Pi#p&Piw&Pi$f&Pi$|&Pi[&Pio&Piy&Pi}&Pi!l&Pi!q&Pi!t&Pi#Q&Pi#R&Pi'h&Pi'r&Pi's&Pi'x&Pi'y&Pi'z&Pi'{&Pi'|&Pi'}&Pi(O&Pi(Q&Pi(R&Pi(T&Pi!P&Pi!e&PiX&PiP&Piv&Pi!S&Pi#T&Pi~O!P3VO~O$f3WO(P*OO~Ow$yO$f3WO$|*wO(P*OO~Ow!PO!Z!YO~O!Z3bO#T3`O's(oO~O!j&bO'R#VO~P<cOv3fO~Ov3fO!j&bO'R`O~P<cO!O3iO's(oO~Ow!PO~P#9jOo3lOy3kO(P*OO~OS,_OT,_OZ,_O['bO_3mOd,_Oo,_Os,_Ou,_Ow'cOy,_O}'dO!S,_O!e,_O!l,_O!q,]O!t,_O!{,_O#Q,_O#R,_O#S,_O#T,_O'R,_O'[%eO'_UO'h,[O's,]O'v,`O'x,[O'y,]O'z,]O'{,]O'|,^O'},^O(O,_O(P,aO(Q,aO(R,bO~O!P3qO~P&']Ov3tO~P&']OR0tO!S+TO!e0rO~OR0tOo0vO!S+TO!e0rO~Oa/yO|!VO'R#VO'U(QO~O!S3wO~O!e3yOv(dX~Ov3{O~OQ!QOZ%rO[%qO^%vO_%cO`TOa/yOd%jOg%yO|!VO}%pO!q%oO$f%wO%^%xO&W%{O'R#VO'U%dO'[%eO(Q%zO~PGaO!e4OO(P*OO!P(fX~O!P4QO~O!S4RO(P)}O~O!S+TO(P*OO~O!e4TOw(aXy(aX~OQ4VO~Oy2XO~Oa/yO|!VO'R#VO'U*nO~Oo4YOw*uO}*vOv%XX!e%XX~O!e4]Ov(hX~Ov4_O~O(P4aOy(_Xw(_X$|(_XR(_Xo(_X!e(_X~Oy4cO(P*OO~OQ!QO[$gO^4dO_$ZO`:QOa:QOd$aO}$eO!i$qO!j4eO!l$lO!q$hO#V$lO'U$YO'[%eO'h$[O~P%;dO!S4gO's(oO~O#V4iO~P%1qO!e4jO!P(kX~O!P4lO~O!P%aX!S!aX!e%aX's!aX~P!KZOQ!QO[$gO^4dO_$ZO`:QOa:QOd$aO}$eO!i$qO!j&bO!l$lO!q$hO#V$lO'U$YO'h$[O~P%;dO!e4jO!P(kX!S'fX's'fX~O^2zO!j2yO~Ow!POy2gO~O_4rO!Q/rO'U&iO'[%eOR&lX!e&lX~OR4tO!e0rO~O!S4vO~Ow$yO$|*wO(P*OO~Oy4wO~P2gOo4xOy4wO(P*OO~Ov&uX!e&uX~P!3|O!e1OOv(pa~Oo5OOy4}O(P*OO~OSVOTVO_%cOsVOtVOuVOw!PO!Q!lO'_UOR&vX!e&vX~P#D_O!e1YOR'Za~O!{5UO~O!P5VO~P#K_O!e5XO!P(wX~O!P5ZO~O!e,rO!P(Vi~OPoOQ!QOZeO[lO^RO_RO`ROaROd[Ow!PO}mO!U#bO!W#cO!X!^O!Z!YO!liO!qgO!tiO#Y!_O#r!ZO#{![O$O!]O$b!`O$d!bO$f!cO'R#VO'UQO'[kO'hcO'riO(QdO~P%;dO!P&ea!e&ea~P2gOX5]O~P2gOP!bqZcq!S!bq!e!bq!tcq'rcq's!bq!O!bqo!bqy!bq!P!bqX!bq!Z!bq#T!bqv!bq~O's(oO!P&`a!e&`aX&`a~O!i-_O!j-^O!l5_O'U-WOv&aX!e&aX~O!e1tOv'pa~O!S5aO~O!S5eO's(oOP$Xa!e$Xa~O(P*OOP$Wq!e$Wq~Ov$^i!e$^i~P2gOw!POy5gO#S5iO$|*wO~Oo5lOy5kO(P*OO~Oy5nO~Oy5nO$|*wO~Oy5rO(P*OO~Ow!POy5gO~Oo5yOy5xO(P*OO~O!S5{O~O!e2tO!P(ta~O^2zO!j2yO'[%eO~OQ!QOZ%rO[%qO^.qO_%cO`TOaTOd%jOg%yO}%pO!j.rO!q.oO!t6PO#V6RO$f%wO%^%xO&W%{O'R#VO'U%dO'[%eO(Q%zO!P&xX!e&xX~PGaOQ!QOZ%rO[%qO^6TO_%cO`TOaTOd%jOg%yO}%pO!j6UO!q%oO$f%wO%^%xO&W%{O'U%dO'[%eO(P)}O(Q%zO~PGaO!P%aa!e%aa~P#4_O^6VO~O#S6WOR%wq#p%wq(P%wqw%wq$f%wq$|%wq[%wqo%wqy%wq}%wq!l%wq!q%wq!t%wq#Q%wq#R%wq'h%wq'r%wq's%wq'x%wq'y%wq'z%wq'{%wq'|%wq'}%wq(O%wq(Q%wq(R%wq(T%wq!P%wq!e%wqX%wqP%wqv%wq!S%wq#T%wq~O(P*OOR&Oq#p&Oqw&Oq$f&Oq$|&Oq[&Oqo&Oqy&Oq}&Oq!l&Oq!q&Oq!t&Oq#Q&Oq#R&Oq'h&Oq'r&Oq's&Oq'x&Oq'y&Oq'z&Oq'{&Oq'|&Oq'}&Oq(O&Oq(Q&Oq(R&Oq(T&Oq!P&Oq!e&OqX&OqP&Oqv&Oq!S&Oq#T&Oq~O(P*OO!P&ya!e&ya~OX6XO~P2gO'[%eOR&wX!e&wX~O!e3SOR(ra~O$f6_O(P*OO~Ow![q~P#9jO#T6bO~O!Z3bO#T6bO's(oO~Ov6gO~O!S1_O#T'wX~O#T6kO~Oy6lO!P6mO~O!P6mO~P&']Oy6pO~Ov6pOy6lO~Ov6pO~P&']Oy6rO~O!e3yOv(da~O!S6uO~Oa/yO|!VO'R#VO'U(QOv&oX!e&oX~O!e4OO(P*OO!P(fa~OQ!QOZ%rO[%qO^%vO_%cO`TOa/yOd%jOg%yO|!VO}%pO!q%oO$f%wO%^%xO&W%{O'R#VO'U%dO'[%eO(Q%zO!P&pX!e&pX~PGaO!e4OO!P(fa~OQ!QOZ%rO[%qO^%vO_%cO`TOaTOd%jOg%yO}%pO!q%oO$f0UO%^%xO&W%{O'U%dO'[%eO(Q%zOw&nX!e&nXy&nX~PGaO!e4TOw(aay(aa~O!e4]Ov(ha~Oo7XOv%Xa!e%Xa~Oo7XOw*uO}*vOv%Xa!e%Xa~Oa/yO|!VO'R#VO'U*nOv&qX!e&qX~O(P*OOy$xaw$xa$|$xaR$xao$xa!e$xa~O(P4aOy(_aw(_a$|(_aR(_ao(_a!e(_a~O!P%aa!S!aX!e%aa's!aX~P!KZOQ!QO[$gO^7`O_$ZO`:QOa:QOd$aO}$eO!i$qO!j&bO!l$lO!q$hO#V$lO'U$YO'h$[O~P%;dO^6VO!j6UO~O!e4jO!P(ka~O!e4jO!P(ka!S'fX's'fX~OQ!QO[$gO^0kO_$ZO`:QOa:QOd$aO}$eO!i$qO!j0lO!l$lO!q0dO!t7dO#V7fO'R#VO'U$YO'[%eO'h$[O!P&sX!e&sX~P%;dO!S7hO's(oO~Ow!POy5gO$|*wO(P*OO~O!S+TOR&la!e&la~Oo0vO!S+TOR&la!e&la~Oo0vOR&la!e&la~O(P*OOR$yi!e$yi~Oy7kO~P2gOo7lOy7kO(P*OO~O(P*OORni!eni~O(P*OOR&va!e&va~O(P)}OR&va!e&va~OS,_OT,_OZ,_O_,_Od,_Oo,_Os,_Ou,_Oy,_O!S,_O!e,_O!l,_O!q,]O!t,_O!{,_O#Q,_O#R,_O#S,_O#T,_O'R,_O'[%eO'_UO'h,[O's,]O'x,[O'y,]O'z,]O'{,]O'|,^O'},^O(O,_O~O(P7nO(Q7nO(R7nO~P''`O!P7pO~P#K_OSVOTVOsVOtVOuVO'_UO!P&zX!e&zX~P!6PO!e5XO!P(wa~O!P&ei!e&ei~P2gO's(oOv!hi!e!hi~O!S7tO~O(P*OOP$Xi!e$Xi~Ov$^q!e$^q~P2gOw!POy7vO~Ow!POy7vO#S7yO$|*wO~Oy7{O~Oy7|O~Oy7}O(P*OO~Ow!POy7vO$|*wO(P*OO~Oo8SOy8RO(P*OO~O!e2tO!P(ti~O(P*OO!P%}X!e%}X~O!P%ai!e%ai~P#4_O^8VO~O!e8[O['cXv$`i}'cX!l'cX!q'cX!t'cX'h'cX'r'cX(Q'cX~P@[OQ#[iS#[iT#[i[#[i^#[i_#[i`#[ia#[id#[is#[it#[iu#[iv$`i}#[i!i#[i!j#[i!l#[i!q#[i!t'cX#V#[i'R#[i'U#[i'_#[i'h#[i'r'cX(Q'cX~P@[O#T#^a~P2gO#T8_O~O!Z3bO#T8`O's(oO~Ov8cO~Oy8eO~P2gOy8gO~Oy6lO!P8hO~Ov8gOy6lO~O!e3yOv(di~O(P*OOv%Qi!e%Qi~O!e4OO!P(fi~O!e4OO(P*OO!P(fi~O(P*OO!P&pa!e&pa~O(P8oOw(bX!e(bXy(bX~O(P*OO!S$wiy$wiw$wi$|$wiR$wio$wi!e$wi~O!e4]Ov(hi~Ov%Xi!e%Xi~P2gOo8rOv%Xi!e%Xi~O!P%ai!S!aX!e%ai's!aX~P!KZO(P*OO!P%`i!e%`i~O!e4jO!P(ki~OQ!QO[$gO^0kO_$ZO`:QOa:QOd$aO}$eO!i$qO!j0lO!l$lO!q0dO!t7dO#V8uO'R#VO'U$YO'[%eO'h$[O~P%;dO!P&sa!S'fX!e&sa's'fX~O(P*OOR$zq!e$zq~Oy8wO~P2gOy8RO~P2gO(P8yO(Q8yO(R8yO~O(P8yO(Q8yO(R8yO~P''`O's(oOv!hq!e!hq~O(P*OOP$Xq!e$Xq~Ow!POy8|O$|*wO(P*OO~Ow!POy8|O~Oy9PO~P2gOy9RO~P2gOo9TOy9RO(P*OO~OQ#[qS#[qT#[q[#[q^#[q_#[q`#[qa#[qd#[qs#[qt#[qu#[qv$`q}#[q!i#[q!j#[q!l#[q!q#[q#V#[q'R#[q'U#[q'_#[q'h#[q~O!e9WO['cXv$`q}'cX!l'cX!q'cX!t'cX'h'cX'r'cX(Q'cX~P@[Oo'cX!t'cX#Q'cX#R'cX#p'cX'r'cX's'cX'x'cX'y'cX'z'cX'{'cX'|'cX'}'cX(O'cX(P'cX(Q'cX(R'cX(T'cX~P'9OO#T9]O~O!Z3bO#T9]O's(oO~Oy9_O~O(P*OOv%Qq!e%Qq~O!e4OO!P(fq~O(P*OO!P&pi!e&pi~O(P8oOw(ba!e(bay(ba~Ov%Xq!e%Xq~P2gO!P&si!S'fX!e&si's'fX~O(P*OO!P%`q!e%`q~Oy9dO~P2gO(P9eO(Q9eO(R9eO~O's(oOv!hy!e!hy~Ow!POy9fO~Ow!POy9fO$|*wO(P*OO~Oy9hO~P2gOQ#[yS#[yT#[y[#[y^#[y_#[y`#[ya#[yd#[ys#[yt#[yu#[yv$`y}#[y!i#[y!j#[y!l#[y!q#[y#V#[y'R#[y'U#[y'_#[y'h#[y~O!e9kO['cXv$`y}'cX!l'cX!q'cX!t'cX'h'cX'r'cX(Q'cX~P@[Oo'cX!t'cX#Q'cX#R'cX#p'cX'r'cX's'cX'x'cX'y'cX'z'cX'{'cX'|'cX'}'cX(O'cX(P'cX(Q'cX(R'cX(T'cX~P'?}O!e9lO['cX}'cX!l'cX!q'cX!t'cX'h'cX'r'cX(Q'cX~P@[OQ#[iS#[iT#[i[#[i^#[i_#[i`#[ia#[id#[is#[it#[iu#[i}#[i!i#[i!j#[i!l#[i!q#[i!t'cX#V#[i'R#[i'U#[i'_#[i'h#[i'r'cX(Q'cX~P@[O#T9oO~O(P*OO!P&pq!e&pq~Ov%Xy!e%Xy~P2gOw!POy9pO~Oy9qO~P2gOQ#[!RS#[!RT#[!R[#[!R^#[!R_#[!R`#[!Ra#[!Rd#[!Rs#[!Rt#[!Ru#[!Rv$`!R}#[!R!i#[!R!j#[!R!l#[!R!q#[!R#V#[!R'R#[!R'U#[!R'_#[!R'h#[!R~O!e9rO['cX}'cX!l'cX!q'cX!t'cX'h'cX'r'cX(Q'cX~P@[OQ#[qS#[qT#[q[#[q^#[q_#[q`#[qa#[qd#[qs#[qt#[qu#[q}#[q!i#[q!j#[q!l#[q!q#[q!t'cX#V#[q'R#[q'U#[q'_#[q'h#[q'r'cX(Q'cX~P@[O!e9uO['cX}'cX!l'cX!q'cX!t'cX'h'cX'r'cX(Q'cX~P@[OQ#[yS#[yT#[y[#[y^#[y_#[y`#[ya#[yd#[ys#[yt#[yu#[y}#[y!i#[y!j#[y!l#[y!q#[y!t'cX#V#[y'R#[y'U#[y'_#[y'h#[y'r'cX(Q'cX~P@[OwbX~P$|OwjX}jX!tbX'rbX~P!6mOZ'TXd'YXo'TXw'lX!t'TX'r'TX's'TX~O['TXd'TXw'TX}'TX!l'TX!q'TX#Q'TX#R'TX#p'TX'h'TX'x'TX'y'TX'z'TX'{'TX'|'TX'}'TX(O'TX(P'TX(Q'TX(R'TX(T'TX~P'MmOP'TX}'lX!S'TX!e'TX!O'TXy'TX!P'TXX'TX!Z'TX#T'TXv'TX~P'MmO^9xO_9xO`9xOa9xO'U9vO~O!j:VO~P!.cOPoOQ!QOZeO^9xO_9xO`9xOa9xOd9{O!U#bO!W#cO!X;RO!Z!YO#Y!_O#r:RO#{:SO$O!]O$b!`O$d!bO$f!cO'U9vO'[kO[#sXo#sXw#sX}#sX!l#sX!q#sX!t#sX#Q#sX#R#sX#p#sX'h#sX'r#sX's#sX'x#sX'y#sX'z#sX'{#sX'|#sX'}#sX(O#sX(P#sX(Q#sX(R#sX(T#sX~P%;dO#S$uO~P!.cO}'lXP'TX!S'TX!e'TX!O'TXy'TX!P'TXX'TX!Z'TX#T'TXv'TX~P'MmOo#qX#Q#qX#R#qX#p#qX's#qX'x#qX'y#qX'z#qX'{#qX'|#qX'}#qX(O#qX(P#qX(R#qX(T#qX~P!.cOo#zX#Q#zX#R#zX#p#zX's#zX'x#zX'y#zX'z#zX'{#zX'|#zX'}#zX(O#zX(P#zX(R#zX(T#zX~P!.cOPoOQ!QOZeO^9xO_9xO`9xOa9xOd9{O!U#bO!W#cO!X;RO!Z!YO#Y!_O#r:RO#{:SO$O!]O$b!`O$d!bO$f!cO'U9vO'[kO[#sao#saw#sa}#sa!l#sa!q#sa!t#sa#Q#sa#R#sa#p#sa'h#sa'r#sa's#sa'x#sa'y#sa'z#sa'{#sa'|#sa'}#sa(O#sa(P#sa(Q#sa(R#sa(T#sa~P%;dOo:aO#Q:aO#R:^Ow#sa~P!BqOw$Ua~P#9jOQ'YXd'YX}iX~OQlXdlX}jX~O^:zO_:zO`:zOa:zO'U:fO~OQ'YXd'YX}hX~Ow#qa~P#9jOw#za~P#9jO!S&_Oo#za#Q#za#R#za#p#za's#za'x#za'y#za'z#za'{#za'|#za'}#za(O#za(P#za(R#za(T#za~P!.cO#S*dO~P!.cOw#ci~P#9jO[#}O}#zO'x#hO(O#|O(Q#hO(R#fO(T#hOo#eiw#ei!l#ei!q#ei!t#ei#Q#ei#R#ei#p#ei'r#ei's#ei'y#ei'z#ei'{#ei'|#ei'}#ei~O'h#ei(P#ei~P(/aO'h#gO(P#gO~P(/aO[#}O}#zO'h#gO'x#hO'y#iO'z#iO(O#|O(P#gO(Q#hO(R#fO(T#hOo#eiw#ei!l#ei!t#ei#Q#ei#R#ei#p#ei'r#ei's#ei'{#ei'|#ei'}#ei~O!q#ei~P(1]O!q#jO~P(1]O[#}O}#zO!q#jO'h#gO'x#hO'y#iO'z#iO'{#kO(O#|O(P#gO(Q#hO(R#fO(T#hOo#eiw#ei!l#ei!t#ei#Q#ei#R#ei#p#ei'r#ei'|#ei'}#ei~O's#ei~P(3UO's#lO~P(3UO[#}O}#zO!q#jO#R:^O'h#gO's#lO'x#hO'y#iO'z#iO'{#kO(O#|O(P#gO(Q#hO(R#fO(T#hOo#eiw#ei!l#ei!t#ei#Q#ei#p#ei'r#ei'|#ei~O'}#ei~P(4}O'}#mO~P(4}Oo:aO#Q:aO#R:^Ow#ni~P$1lOo:aO#Q:aO#R:^Ow#si~P$3mOQ'YXd'YX}'lX~Ow#zi~P#9jOw$Ui~P#9jOd:UO~Ow#ca~P#9jOd:|O~OU'x_'v'Q'P'_s!{'_'U'[~",
  goto: "$L^(xPPPPPPP(yPP)QPP)`PPPP)l-rP0r5oP7a7a9U7a?VDoEQPEWHaPPPPPPKqP! b! pPPPPP!!hP!%QP!%QPP!'QP!)TP!)Y!*P!*w!*w!*w!)Y!+nP!)Y!.c!.fPP!.lP!)Y!)Y!)Y!)YP!)Y!)YP!)Y!)Y!/[!/[!/y!0hP!0hKaKaKaPPPP!0hPP!%QP!0v!0y!1P!2Q!2^!4^!4^!6[!8^!2^!2^!:Y!;w!=h!?T!@n!BV!Cl!D}!2^!2^P!2^P!2^!2^!F^!2^P!G}!2^!2^P!I}!2^P!2^!8^!8^!2^!8^!2^!LU!N^!Na!8^!2^!Nd!Ng!Ng!Ng!Nk!%QP!%QP!%QP! b! bP!Nu! b! bP# R#!g! bP! bP#!v##{#$T#$s#$w#$}#$}#%VP#']#']#'c#(X#(e! bP! bP#(u#)U! bP! bPP#)b#)p#)|#*f#)v! b! bP! b! b! bP#*l#*l#*r#*x#*l#*l! b! bP#+V#+`#+j#+j#-b#/U#/b#/b#/e#/e5o5o5o5o5o5o5o5oP5o#/h#/n#0Y#2e#2k#2z#6x#7O#7U#7h#7r#9c#9m#9|#:S#:Y#:d#:n#:t#;R#;X#;_#;i#;w#<R#>a#>m#>z#?Q#?Y#?a#?k#?qPPPPPPP#?w#CTP#GS#Kn#Mi$ h$'UP$'XPPP$*`$*i$*{$0V$2e$2n$4gP!)Y$5a$8u$;l$?W$?a$?f$?iPPP$?l$BcP$BsPPPPPPPPPP$CXP$Eg$Ej$Em$Es$Ev$Ey$E|$FP$FV$Ha$Hd$Hg$Hj$Hm$Hp$Hs$Hv$Hy$H|$IP$KV$KY$K]#*l$Ki$Ko$Kr$Ku$Ky$K}$LQ$LT$LW$LZQ!tPT'V!s'Wi!SOlm!P!T$T$W$y%b)T*e/fQ'h#QQ,l'kQ1d,kR7q5X(SSOY[bfgilmop!O!P!T!Y!Z![!_!`!c!p!q!|!}#Q#U#Z#e#o#p#q#r#s#t#u#v#w#x#y#z#}$T$W$`$a$e$g$h$q$r$y%X%_%b&U&Y&[&b&u&z&|'P'a'k'm'n'|(V(X(a(c(d(e(i(n(o(q({)R)T)h*Y*e*h*j*k+Y+m+y,k,o,r,y-Q-S-f-l-s.|/]/a/c/f0d0f0l0|1O1i1s1}3`3b3g3i3l4Y4e4j4x5O5X5a5l5y6b6f7X7c7l7t8S8_8`8r9T9]9o9z9{9|9}:O:P:R:S:T:U:V:W:X:Y:Z:[:]:^:_:`:a:b:l:m:nS(y$v-nQ*o&eQ*s&hQ-j(xQ-x)YW0Y+P0X4]7ZR4[0Z&{!RObfgilmop!O!P!T!Y!Z![!_!`!c!p#Q#e#o#p#q#r#s#t#u#v#w#x#y#z#}$T$W$e$g$h$q$r$y%_%b&U&Y&[&b&u'k'|(V(X(a(e(i(n(o(q({)R)T)h*Y*e*h*j*k+Y+m,k,r,y-S-f-l-s.|/]/a/c/f0d0f0l0|1i1s1}3`3b3g3i3l4Y4e4j4x5O5X5a5l5y6b6f7X7c7l7t8S8_8`8r9T9]9o9|9}:O:P:R:S:V:W:X:Y:Z:[:]:^:_:`:a:b:l:m#r]Ofgilmp!O!P!T!Z![#e#o#p#q#r#s#t#u#v#w#x#y#z#}$T$W%_%b&Y&['|(V(X({)h+m,r,y-l.|0|1i1}3`3b3l4Y4x5O5l5y6b7X7l8S8_8`8r9T9]9of#[b#Q$y'k(a)R)T*Y,k-s5X!h$bo!c!p$e$g$h$q$r&U&b&u(e(i(n(o(q*e*j+Y-S-f/a/c/f0d0f0l1s3g4e4j5a6f7c7t$b%k!Q!n$O$u%o%p%q%y%{&P&o&p&r'](p)r)w)x)}*O*Q*U*Z*]*d*m*v*w+T+U+g+n+|-h-u.T._.o.s.w.x/Y/Z/z/|0_0q0v0}1X1Y1z2b2i2k2n2t2w3W3w3}4O4T4W4a4g4v5e5i5{6W6_6u6{6}7h7y8n!W;Q!Y!_!`*h*k/]3i9|9}:O:P:R:S:V:W:X:Y:Z:[:]:^:_:`:a:b:l:mR;T%n$_%u!Q!n$O$u%o%p%q&P&o&p&r'](p)r)w)x)}*O*Q*U*Z*]*d*m*v*w+T+U+g+n+|-h-u.T._.o.s.w.x/Y/Z/z/|0_0q0v0}1X1Y1z2b2i2k2n2t2w3W3w3}4O4T4W4a4g4v5e5i5{6W6_6u6{6}7h7y8n$e%l!Q!n$O$u%n%o%p%q%y%{&P&o&p&r'](p)r)w)x)}*O*Q*U*Z*]*d*m*v*w+T+U+g+n+|-h-u.T._.o.s.w.x/Y/Z/z/|0_0q0v0}1X1Y1z2b2i2k2n2t2w3W3w3}4O4T4W4a4g4v5e5i5{6W6_6u6{6}7h7y8n'hZOY[fgilmop!O!P!T!Y!Z![!_!`!c!p!|!}#e#o#p#q#r#s#t#u#v#w#x#y#z#}$T$W$`$a$e$g$h$q$r%_%b%i%j&U&Y&[&b&u'a'|(V(X(c(d(e(i(n(o(q({)h)o)p*e*h*j*k+Y+m,r,y-Q-S-f-l.h.|/]/a/c/f0d0f0l0|1i1s1}3`3b3g3i3l4Y4e4j4x5O5a5l5y6b6f7X7c7l7t8S8_8`8r9T9]9o9z9{9|9}:O:P:R:S:T:U:V:W:X:Y:Z:[:]:^:_:`:a:b:g:h:l:m:n:{:|;P$^%l!Q!n$O$u%n%o%p%q%y%{&P&p&r(p)r)w)x)}*O*Q*U*Z*]*d*m*v*w+T+U+g+n+|-h-u.T._.o.s.w.x/Y/Z/z/|0_0q0v0}1X1z2b2i2k2n2t2w3W3w3}4O4T4W4a4g4v5e5i5{6W6_6u6{6}7h7y8nQ&j!hQ&k!iQ&l!jQ&m!kQ&s!oQ)Z%QQ)[%RQ)]%SQ)^%TQ)a%WQ+_&oS,Q']1YQ.V)_S/q*t4VR4p0r+}TOY[bfgilmop!O!P!Q!T!Y!Z![!_!`!c!n!p!q!|!}#Q#U#Z#e#o#p#q#r#s#t#u#v#w#x#y#z#}$O$T$W$`$a$e$g$h$q$r$u$y%X%_%b%i%j%n%o%p%q%y%{&P&U&Y&[&b&o&p&r&u&z&|'P']'a'k'm'n'|(V(X(a(c(d(e(i(n(o(p(q({)R)T)h)o)p)r)w)x)}*O*Q*U*Y*Z*]*d*e*h*j*k*m*v*w+T+U+Y+g+m+n+y+|,k,o,r,y-Q-S-f-h-l-s-u.T._.h.o.s.w.x.|/Y/Z/]/a/c/f/z/|0_0d0f0l0q0v0|0}1O1X1Y1i1s1z1}2b2i2k2n2t2w3W3`3b3g3i3l3w3}4O4T4W4Y4a4e4g4j4v4x5O5X5a5e5i5l5y5{6W6_6b6f6u6{6}7X7c7h7l7t7y8S8_8`8n8r9T9]9o9z9{9|9}:O:P:R:S:T:U:V:W:X:Y:Z:[:]:^:_:`:a:b:g:h:l:m:n:{:|;PQ'[!xQ'g#PQ)k%gU)q%m*S*VR.e)jQ,S']R5R1Y#t%s!Q!n$O$u%p%q&P&p&r(p)w)x)}*Q*U*Z*]*d*m*v+U+g+n+|-h-u.T._.s.w.x/Y/Z/z/|0_0q0v0}1X1z2b2i2k2n2w3W3w3}4O4W4g4v5e5i5{6W6_6u6{6}7h7y8nQ)w%oQ+^&oQ,T']l,_'b'c'd,Y,e,f/l/m1`3p3s5V5W7pS.p)r2tQ.}*OQ/P*RQ/p*tS0P*w4TQ0`+T[0n+Y.i0f4j6O7cQ2w.oS4f0d2sQ4o0rQ5S1YQ6Y3SQ7P4RQ7T4VQ7^4aR9a8o&pVOfgilmop!O!P!T!Y!Z![!_!`!c!p#e#o#p#q#r#s#t#u#v#w#x#y#z#}$T$W$e$g$h$q$r%_%b&U&Y&[&b&u']'|(V(X(a(e(i(n(o(q({)h*e*h*j*k+Y+m,j,k,r,y-S-f-l.|/]/a/c/f0d0f0l0|1Y1i1s1}3`3b3g3i3l4Y4e4j4x5O5X5a5l5y6b6f7X7c7l7t8S8_8`8r9T9]9o9|9}:O:P:R:S:V:W:X:Y:Z:[:]:^:_:`:a:b:l:mU&g!g%P%[m,_'b'c'd,Y,e,f/l/m1`3p3s5V5W7p$nsOfgilm!O!P!T!Y!Z![!_!`#e#o#p#q#r#s#t#u#v#w#x#y#z#}$T$W%_%b&Y'|(V(X({)h*h*k+m,r,y-l.|/]0|1i1}3`3b3i3l4Y4x5O5l5y6b7X7l8S8_8`8r9T9]9o9|9}:R:S:V:W:X:Y:Z:[:]:^:_:`:a:lS$tp:PS&O!W#bS&Q!X#cQ&`!bQ*^&RQ*`&VS*c&[:mQ*g&^Q,S']Q-i(vQ/h*iQ0o+ZS2g.W0pQ3^/^Q3_/_Q3h/gQ3j/jQ5R1YU5g2S2h4nU7v5h5j5wQ8d6iS8|7w7xS9f8}9OR9p9gi{Ob!O!P!T$y%_%b)R)T)h-shxOb!O!P!T$y%_%b)R)T)h-sW/u*u/s3y6vQ/|*vW0Z+P0X4]7ZQ3}/zQ6}4OR8n6{!h$do!c!p$e$g$h$q$r&U&b&u(e(i(n(o(q*e*j+Y-S-f/a/c/f0d0f0l1s3g4e4j5a6f7c7tQ&d!dQ&f!fQ&n!mW&x!q%X&|1OQ'S!rQ)W$}Q)X%OQ)`%VU)c%Y'T'UQ*r&hS+r&z'PS-X(j1tQ-t)VQ-w)YS.`)d)eS0w+b/rQ1R+yQ1V+zS1w-^-_Q2l.aQ3u/oQ5b1yR5m2W${sOfgilmp!O!P!T!Y!Z![!_!`#e#o#p#q#r#s#t#u#v#w#x#y#z#}$T$W%_%b&Y&['|(V(X({)h*h*k+m,r,y-l.|/]0|1i1}3`3b3i3l4Y4x5O5l5y6b7X7l8S8_8`8r9T9]9o9|9}:O:P:R:S:V:W:X:Y:Z:[:]:^:_:`:a:b:l:m$zsOfgilmp!O!P!T!Y!Z![!_!`#e#o#p#q#r#s#t#u#v#w#x#y#z#}$T$W%_%b&Y&['|(V(X({)h*h*k+m,r,y-l.|/]0|1i1}3`3b3i3l4Y4x5O5l5y6b7X7l8S8_8`8r9T9]9o9|9}:O:P:R:S:V:W:X:Y:Z:[:]:^:_:`:a:b:l:mR3^/^V&T!Y!`*h!i$lo!c!p$e$g$h$q$r&U&b&u(e(i(n(o(q*e*j+Y-S-f/a/c/f0d0f0l1s3g4e4j5a6f7c7t!k$^o!c!p$e$g$h$q$r&U&b&u(a(e(i(n(o(q*e*j+Y-S-f/a/c/f0d0f0l1s3g4e4j5a6f7c7t!i$co!c!p$e$g$h$q$r&U&b&u(e(i(n(o(q*e*j+Y-S-f/a/c/f0d0f0l1s3g4e4j5a6f7c7t&e^Ofgilmop!O!P!T!Y!Z![!_!`!c!p#e#o#p#q#r#s#t#u#v#w#x#y#z#}$T$W$e$g$h$q$r%_%b&U&Y&[&b&u'|(V(X(e(i(n(o(q({)h*e*h*j*k+Y+m,r,y-S-f-l.|/]/a/c/f0d0f0l0|1i1s1}3`3b3g3i3l4Y4e4j4x5O5a5l5y6b6f7X7c7l7t8S8_8`8r9T9]9o9|9}:O:P:R:S:V:W:X:Y:Z:[:]:^:_:`:a:b:l:mR(k$fQ-Z(jR5_1tQ(R#|S(z$v-nS-Y(j1tQ-k(xW/t*u/s3y6vS1x-^-_Q3x/uR5c1yQ'e#Oh,b'b'c'd,Y,e,f/l/m1`3p3s5WQ,m'lQ,p'oQ.t)tR8f6kQ'f#Oh,b'b'c'd,Y,e,f/l/m1`3p3s5WQ,n'lQ,p'oQ.t)tR8f6ki,b'b'c'd,Y,e,f/l/m1`3p3s5WR*f&]X/b*e/c/f3g!}aOb!O!P!T#z$v$y%_%b'|(x)R)T)h)r*e*u*v+P+Y,r-n-s.i/a/c/f/s/z0X0f1i2t3g3y4O4]4j6O6f6v6{7Z7cQ3a/`Q6d3cQ8a6eR9^8b${rOfgilmp!O!P!T!Y!Z![!_!`#e#o#p#q#r#s#t#u#v#w#x#y#z#}$T$W%_%b&Y&['|(V(X({)h*h*k+m,r,y-l.|/]0|1i1}3`3b3i3l4Y4x5O5l5y6b7X7l8S8_8`8r9T9]9o9|9}:O:P:R:S:V:W:X:Y:Z:[:]:^:_:`:a:b:l:m#nfOfglmp!O!P!T!Z![#e#o#p#q#r#s#t#u#v#w#x#z#}$T$W%_%b&Y&['|(V(X({)h+m,r,y-l.|0|1i1}3`3b3l4Y4x5O5l5y6b7X7l8S8_8`8r9T9]9o!T9|!Y!_!`*h*k/]3i9|9}:P:R:S:V:W:X:Y:Z:[:]:^:_:`:a:l:m#rfOfgilmp!O!P!T!Z![#e#o#p#q#r#s#t#u#v#w#x#y#z#}$T$W%_%b&Y&['|(V(X({)h+m,r,y-l.|0|1i1}3`3b3l4Y4x5O5l5y6b7X7l8S8_8`8r9T9]9o!X9|!Y!_!`*h*k/]3i9|9}:O:P:R:S:V:W:X:Y:Z:[:]:^:_:`:a:b:l:m$srOfglmp!O!P!T!Y!Z![!_!`#e#o#p#q#r#s#t#u#v#w#x#z#}$T$W%_%b&Y&['|(V(X({)h*h*k+m,r,y-l.|/]0|1i1}3`3b3i3l4Y4x5O5l5y6b7X7l8S8_8`8r9T9]9o9|9}:P:R:S:V:W:X:Y:Z:[:]:^:_:`:a:l:m#U#oh#d$P$Q$V$s%^&W&X'p's't'u'v'w'x'y'z'{'}(T(Z(_*a*b,q,v,x-m0y1j1m2O3Q4y5[5f6c6j7W7j7m7z8Q8q8x9S9c9i}:W&S&]/j3]6i:c:d:j:k:o:q:r:s:t:u:v:w:x:y:};O;S#W#ph#d$P$Q$V$s%^&W&X'p'q's't'u'v'w'x'y'z'{'}(T(Z(_*a*b,q,v,x-m0y1j1m2O3Q4y5[5f6c6j7W7j7m7z8Q8q8x9S9c9i!P:X&S&]/j3]6i:c:d:j:k:o:p:q:r:s:t:u:v:w:x:y:};O;S#S#qh#d$P$Q$V$s%^&W&X'p't'u'v'w'x'y'z'{'}(T(Z(_*a*b,q,v,x-m0y1j1m2O3Q4y5[5f6c6j7W7j7m7z8Q8q8x9S9c9i{:Y&S&]/j3]6i:c:d:j:k:o:r:s:t:u:v:w:x:y:};O;S#Q#rh#d$P$Q$V$s%^&W&X'p'u'v'w'x'y'z'{'}(T(Z(_*a*b,q,v,x-m0y1j1m2O3Q4y5[5f6c6j7W7j7m7z8Q8q8x9S9c9iy:Z&S&]/j3]6i:c:d:j:k:o:s:t:u:v:w:x:y:};O;S#O#sh#d$P$Q$V$s%^&W&X'p'v'w'x'y'z'{'}(T(Z(_*a*b,q,v,x-m0y1j1m2O3Q4y5[5f6c6j7W7j7m7z8Q8q8x9S9c9iw:[&S&]/j3]6i:c:d:j:k:o:t:u:v:w:x:y:};O;S!|#th#d$P$Q$V$s%^&W&X'p'w'x'y'z'{'}(T(Z(_*a*b,q,v,x-m0y1j1m2O3Q4y5[5f6c6j7W7j7m7z8Q8q8x9S9c9iu:]&S&]/j3]6i:c:d:j:k:o:u:v:w:x:y:};O;S!x#vh#d$P$Q$V$s%^&W&X'p'y'z'{'}(T(Z(_*a*b,q,v,x-m0y1j1m2O3Q4y5[5f6c6j7W7j7m7z8Q8q8x9S9c9iq:_&S&]/j3]6i:c:d:j:k:o:w:x:y:};O;S!v#wh#d$P$Q$V$s%^&W&X'p'z'{'}(T(Z(_*a*b,q,v,x-m0y1j1m2O3Q4y5[5f6c6j7W7j7m7z8Q8q8x9S9c9io:`&S&]/j3]6i:c:d:j:k:o:x:y:};O;S$]#{h#`#d$P$Q$V$s%^&S&W&X&]'p'q'r's't'u'v'w'x'y'z'{'}(T(Z(_*a*b,q,v,x-m/j0y1j1m2O3Q3]4y5[5f6c6i6j7W7j7m7z8Q8q8x9S9c9i:c:d:j:k:o:p:q:r:s:t:u:v:w:x:y:};O;S${jOfgilmp!O!P!T!Y!Z![!_!`#e#o#p#q#r#s#t#u#v#w#x#y#z#}$T$W%_%b&Y&['|(V(X({)h*h*k+m,r,y-l.|/]0|1i1}3`3b3i3l4Y4x5O5l5y6b7X7l8S8_8`8r9T9]9o9|9}:O:P:R:S:V:W:X:Y:Z:[:]:^:_:`:a:b:l:m$v!aOfgilmp!O!P!T!Y!Z!_!`#e#o#p#q#r#s#t#u#v#w#x#y#z#}$T$W%_%b&Y&['|(V(X({)h*h*k+m,r,y-l.|/]0|1i1}3`3b3i3l4Y4x5O5l5y6b7X7l8S8_8`8r9T9]9o9|9}:O:P:R:V:W:X:Y:Z:[:]:^:_:`:a:b:l:mQ&Y![Q&Z!]R:l:S#rpOfgilmp!O!P!T!Z![#e#o#p#q#r#s#t#u#v#w#x#y#z#}$T$W%_%b&Y&['|(V(X({)h+m,r,y-l.|0|1i1}3`3b3l4Y4x5O5l5y6b7X7l8S8_8`8r9T9]9oQ&[!^!W:P!Y!_!`*h*k/]3i9|9}:O:P:R:S:V:W:X:Y:Z:[:]:^:_:`:a:b:l:mR:m;RR$moR-e(qR$wqT(|$v-nQ/e*eS3e/c/fR6h3gQ3o/lQ3r/mQ6n3pR6q3sQ$zwQ)U${Q*p&fQ+e&qQ+h&sQ-v)XW.Y)a+i+j+kS/W*[+fW2c.V.Z.[.]U3X/X/[0xU5t2d2e2fS6]3Y3[S8O5u5vS8X6[6^Q9Q8PS9U8Y8ZR9j9V^|O!O!P!T%_%b)hX)Q$y)R)T-sQ&r!nQ*]&PQ*{&jQ+O&kQ+S&lQ+V&mQ+[&nQ+k&sQ-|)ZQ.P)[Q.S)]Q.U)^Q.X)`Q.])aQ2T-tQ2f.VR4W0UU+`&o*t4VR4q0rQ+X&mQ+j&sS.[)a+k^0u+^+_/p/q4o4p7TS2e.V.]S4S0Q0RR5v2fS0Q*w4TQ0`+TR7^4aU+c&o*t4VR4r0rQ*y&jQ*}&kQ+R&lQ+f&qQ+i&sS-z)Z*{S.O)[+OS.R)]+SU.Z)a+j+kQ/X*[Q0W*zQ0p+ZQ2Y-{Q2Z-|Q2^.PQ2`.SU2d.V.[.]Q2h.WS3[/[0xS5h2S4nQ5o2[S5u2e2fQ6^3YS7x5j5wQ8P5vQ8Y6[Q8}7wQ9V8ZR9g9OQ0S*wR7R4TQ*x&jQ*|&kU-y)Z*y*{U-})[*}+OS2X-z-|S2].O.PQ4Z0YQ5n2ZQ5p2^R7Y4[Q/v*uQ3v/sQ6w3yR8k6vQ*z&jS-{)Z*{Q2[-|Q4Z0YR7Y4[Q+Q&lU.Q)]+R+SS2_.R.SR5q2`Q0[+PQ4X0XQ7[4]R8s7ZQ+Z&nS.W)`+[S2S-t.XR5j2TQ0h+YQ4h0fQ7e4jR8t7cQ.l)rQ0h+YQ2q.iQ4h0fQ6R2tQ7e4jQ8U6OR8t7cQ0h+YR4h0fX'O!q%X&|1OX&{!q%X&|1OW'O!q%X&|1OS+t&z'PR1T+y_|O!O!P!T%_%b)hQ%a!PS)g%_%bR.c)h$^%u!Q!n$O$u%o%p%q&P&o&p&r'](p)r)w)x)}*O*Q*U*Z*]*d*m*v*w+T+U+g+n+|-h-u.T._.o.s.w.x/Y/Z/z/|0_0q0v0}1X1Y1z2b2i2k2n2t2w3W3w3}4O4T4W4a4g4v5e5i5{6W6_6u6{6}7h7y8nQ*T%yR*W%{$c%n!Q!n$O$u%o%p%q%y%{&P&o&p&r'](p)r)w)x)}*O*Q*U*Z*]*d*m*v*w+T+U+g+n+|-h-u.T._.o.s.w.x/Y/Z/z/|0_0q0v0}1X1Y1z2b2i2k2n2t2w3W3w3}4O4T4W4a4g4v5e5i5{6W6_6u6{6}7h7y8nW)s%m%x*S*VQ.d)iR2|.uR.l)rR6R2tQ'W!sR+}'WQ!TOQ$TlQ$WmQ%b!P[%|!T$T$W%b)T/fQ)T$yR/f*e$b%i!Q!n$O$u%o%p%q%y%{&P&o&p&r'](p)r)w)x)}*O*Q*U*Z*]*d*m*v*w+T+U+g+n+|-h-u.T._.o.s.w.x/Y/Z/z/|0_0q0v0}1X1Y1z2b2i2k2n2t2w3W3w3}4O4T4W4a4g4v5e5i5{6W6_6u6{6}7h7y8n[)m%i)o.h:g:{;PQ)o%jQ.h)pQ:g%nQ:{:hR;P:|Q!vUR'Y!vS!OO!TU%]!O%_)hQ%_!PR)h%b#rYOfgilmp!O!P!T!Z![#e#o#p#q#r#s#t#u#v#w#x#y#z#}$T$W%_%b&Y&['|(V(X({)h+m,r,y-l.|0|1i1}3`3b3l4Y4x5O5l5y6b7X7l8S8_8`8r9T9]9oh!yY!|#U$`'a'm(c,o-Q9z:T:nQ!|[f#Ub#Q$y'k(a)R)T*Y,k-s5X!h$`o!c!p$e$g$h$q$r&U&b&u(e(i(n(o(q*e*j+Y-S-f/a/c/f0d0f0l1s3g4e4j5a6f7c7tQ'a!}Q'm#ZQ(c$aQ,o'nQ-Q(d!W9z!Y!_!`*h*k/]3i9|9}:O:P:R:S:V:W:X:Y:Z:[:]:^:_:`:a:b:l:mQ:T9{R:n:UQ-T(fR1q-TQ1u-ZR5`1uQ,Y'bQ,e'cQ,f'dW1^,Y,e,f5WR5W1`Q/c*eS3d/c3gR3g/ffbO!O!P!T$y%_%b)R)T)h-sp#Wb'|(x.i/a/s/z0X0f1i6O6f6v6{7Z7cQ'|#zS(x$v-nQ.i)rW/a*e/c/f3gQ/s*uQ/z*vQ0X+PQ0f+YQ1i,rQ6O2tQ6v3yQ6{4OQ7Z4]R7c4jQ,s'}Q1h,qT1k,s1hS(W$Q(ZQ(]$VU,w(W(],|R,|(_Q(r$mR-g(rQ-o(}R2Q-oQ3p/lQ3s/mT6o3p3sQ)R$yS-q)R-sR-s)TQ4b0`R7_4b`0s+]+^+_+`+c/p/q7TR4s0sQ8p7PR9b8pQ4U0SR7S4UQ3z/vQ6s3vT6x3z6sQ4P/{Q6y3|U7O4P6y8lR8l6zQ4^0[Q7V4XT7]4^7VhzOb!O!P!T$y%_%b)R)T)h-sQ$|xW%Zz$|%f)u$b%f!Q!n$O$u%o%p%q%y%{&P&o&p&r'](p)r)w)x)}*O*Q*U*Z*]*d*m*v*w+T+U+g+n+|-h-u.T._.o.s.w.x/Y/Z/z/|0_0q0v0}1X1Y1z2b2i2k2n2t2w3W3w3}4O4T4W4a4g4v5e5i5{6W6_6u6{6}7h7y8nR)u%nS4k0h0mS7b4h4iT7g4k7bW&z!q%X&|1OS+q&z+yR+y'PQ1P+vR4|1PU1Z,R,S,TR5T1ZS3T/P7TR6Z3TQ2u.lQ5}2qT6S2u5}Q.y)yR3P.yQ5Y1dR7r5Y^_O!O!P!T%_%b)hY#Xb$y)R)T-s$l#_fgilmp!Y!Z![!_!`#e#o#p#q#r#s#t#u#v#w#x#y#z#}$T$W&Y&['|(V(X({*h*k+m,r,y-l.|/]0|1i1}3`3b3i3l4Y4x5O5l5y6b7X7l8S8_8`8r9T9]9o9|9}:O:P:R:S:V:W:X:Y:Z:[:]:^:_:`:a:b:l:m!h$io!c!p$e$g$h$q$r&U&b&u(e(i(n(o(q*e*j+Y-S-f/a/c/f0d0f0l1s3g4e4j5a6f7c7tW'i#Q'k,k5XQ-O(aR/U*Y&z!RObfgilmop!O!P!T!Y!Z![!_!`!c!p#Q#e#o#p#q#r#s#t#u#v#w#x#y#z#}$T$W$e$g$h$q$r$y%_%b&U&Y&[&b&u'k'|(V(X(a(e(i(n(o(q({)R)T)h*Y*e*h*j*k+Y+m,k,r,y-S-f-l-s.|/]/a/c/f0d0f0l0|1i1s1}3`3b3g3i3l4Y4e4j4x5O5X5a5l5y6b6f7X7c7l7t8S8_8`8r9T9]9o9|9}:O:P:R:S:V:W:X:Y:Z:[:]:^:_:`:a:b:l:m[!{Y[#U#Z9z9{W&{!q%X&|1O['`!|!}'m'n:T:US(b$`$aS+s&z'PU,W'a,o:nS-P(c(dQ1S+yR1o-QS%t!Q&oQ&q!nQ(U$OQ(v$uS)v%o.oQ)y%pQ)|%qS*[&P&rQ+d&pQ,R']Q-c(pQ.k)rU.v)w)x2wS.})}*OQ/O*QQ/S*UQ/V*ZQ/[*]Q/_*dQ/k*mQ/{*vS0R*w4TQ0`+TQ0b+UQ0x+gQ0z+nQ1W+|Q1|-hQ2U-uQ2a.TQ2j._Q2{.sQ2}.wQ3O.xQ3Y/YQ3Z/ZS3|/z/|Q4`0_Q4n0qQ4u0vQ4z0}Q5P1XQ5Q1YQ5d1zQ5s2bQ5w2iQ5z2kQ5|2nQ6Q2tQ6[3WQ6t3wQ6z3}Q6|4OQ7U4WQ7^4aQ7a4gQ7i4vQ7u5eQ7w5iQ8T5{Q8W6WQ8Z6_Q8j6uS8m6{6}Q8v7hQ9O7yR9`8n$^%m!Q!n$O$u%o%p%q&P&o&p&r'](p)r)w)x)}*O*Q*U*Z*]*d*m*v*w+T+U+g+n+|-h-u.T._.o.s.w.x/Y/Z/z/|0_0q0v0}1X1Y1z2b2i2k2n2t2w3W3w3}4O4T4W4a4g4v5e5i5{6W6_6u6{6}7h7y8nQ)i%nQ*S%yR*V%{$y%h!Q!n$O$u%i%j%n%o%p%q%y%{&P&o&p&r'](p)o)p)r)w)x)}*O*Q*U*Z*]*d*m*v*w+T+U+g+n+|-h-u.T._.h.o.s.w.x/Y/Z/z/|0_0q0v0}1X1Y1z2b2i2k2n2t2w3W3w3}4O4T4W4a4g4v5e5i5{6W6_6u6{6}7h7y8n:g:h:{:|;P'tWOY[bfgilmop!O!P!T!Y!Z![!_!`!c!p!|!}#Q#U#Z#e#o#p#q#r#s#t#u#v#w#x#y#z#}$T$W$`$a$e$g$h$q$r$y%_%b&U&Y&[&b&u'a'k'm'n'|(V(X(a(c(d(e(i(n(o(q({)R)T)h*Y*e*h*j*k+Y+m,k,o,r,y-Q-S-f-l-s.|/]/a/c/f0d0f0l0|1i1s1}3`3b3g3i3l4Y4e4j4x5O5X5a5l5y6b6f7X7c7l7t8S8_8`8r9T9]9o9z9{9|9}:O:P:R:S:T:U:V:W:X:Y:Z:[:]:^:_:`:a:b:l:m:n$x%g!Q!n$O$u%i%j%n%o%p%q%y%{&P&o&p&r'](p)o)p)r)w)x)}*O*Q*U*Z*]*d*m*v*w+T+U+g+n+|-h-u.T._.h.o.s.w.x/Y/Z/z/|0_0q0v0}1X1Y1z2b2i2k2n2t2w3W3w3}4O4T4W4a4g4v5e5i5{6W6_6u6{6}7h7y8n:g:h:{:|;P_&y!q%X&z&|'P+y1OR,U']$zrOfgilmp!O!P!T!Y!Z![!_!`#e#o#p#q#r#s#t#u#v#w#x#y#z#}$T$W%_%b&Y&['|(V(X({)h*h*k+m,r,y-l.|/]0|1i1}3`3b3i3l4Y4x5O5l5y6b7X7l8S8_8`8r9T9]9o9|9}:O:P:R:S:V:W:X:Y:Z:[:]:^:_:`:a:b:l:m!j$]o!c!p$e$g$h$q$r&U&b&u(a(e(i(n(o(q*e*j+Y-S-f/a/c/f0d0f0l1s3g4e4j5a6f7c7tQ,S']Q1c,jQ1d,kQ5R1YR7q5X_}O!O!P!T%_%b)h^|O!O!P!T%_%b)hQ#YbX)Q$y)R)T-sbhO!O!T3`6b8_8`9]9oS#`f9|Q#dgQ$PiQ$QlQ$VmQ$spW%^!P%_%b)hU&S!Y!`*hQ&W!ZQ&X![Q&]!_Q'p#eQ'q#oS'r#p:XQ's#qQ't#rQ'u#sQ'v#tQ'w#uQ'x#vQ'y#wQ'z#xQ'{#yQ'}#zQ(T#}Q(Z$TQ(_$WQ*a&YQ*b&[Q,q'|Q,v(VQ,x(XQ-m({Q/j*kQ0y+mQ1j,rQ1m,yQ2O-lQ3Q.|Q3]/]Q4y0|Q5[1iQ5f1}Q6c3bQ6i3iQ6j3lQ7W4YQ7j4xQ7m5OQ7z5lQ8Q5yQ8q7XQ8x7lQ9S8SQ9c8rQ9i9TQ:c:OQ:d:PQ:j:RQ:k:SQ:o:VQ:p:WQ:q:YQ:r:ZQ:s:[Q:t:]Q:u:^Q:v:_Q:w:`Q:x:aQ:y:bQ:}:lQ;O:mR;S9}^tO!O!P!T%_%b)h$`#afgilmp!Y!Z![!_!`#e#o#p#q#r#s#t#u#v#w#x#y#z#}$T$W&Y&['|(V(X({*h*k+m,r,y-l.|/]0|1i1}3b3i3l4Y4x5O5l5y7X7l8S8r9T9|9}:O:P:R:S:V:W:X:Y:Z:[:]:^:_:`:a:b:l:mQ6a3`Q8^6bQ9Y8_Q9[8`Q9n9]R9t9oQ&V!YQ&^!`R/g*hQ$joQ&a!cQ&t!pU(f$e$g(iS(m$h0dQ(t$qQ(u$rQ*_&UQ*l&bQ+o&uQ-R(eS-a(n4eQ-b(oQ-d(qW/`*e/c/f3gQ/i*jW0e+Y0f4j7cQ1p-SQ1{-fQ3c/aQ4m0lQ5^1sQ7s5aQ8b6fR8{7t!h$_o!c!p$e$g$h$q$r&U&b&u(e(i(n(o(q*e*j+Y-S-f/a/c/f0d0f0l1s3g4e4j5a6f7c7tR-O(a'uXOY[bfgilmop!O!P!T!Y!Z![!_!`!c!p!|!}#Q#U#Z#e#o#p#q#r#s#t#u#v#w#x#y#z#}$T$W$`$a$e$g$h$q$r$y%_%b&U&Y&[&b&u'a'k'm'n'|(V(X(a(c(d(e(i(n(o(q({)R)T)h*Y*e*h*j*k+Y+m,k,o,r,y-Q-S-f-l-s.|/]/a/c/f0d0f0l0|1i1s1}3`3b3g3i3l4Y4e4j4x5O5X5a5l5y6b6f7X7c7l7t8S8_8`8r9T9]9o9z9{9|9}:O:P:R:S:T:U:V:W:X:Y:Z:[:]:^:_:`:a:b:l:m:n$zqOfgilmp!O!P!T!Y!Z![!_!`#e#o#p#q#r#s#t#u#v#w#x#y#z#}$T$W%_%b&Y&['|(V(X({)h*h*k+m,r,y-l.|/]0|1i1}3`3b3i3l4Y4x5O5l5y6b7X7l8S8_8`8r9T9]9o9|9}:O:P:R:S:V:W:X:Y:Z:[:]:^:_:`:a:b:l:m!i$fo!c!p$e$g$h$q$r&U&b&u(e(i(n(o(q*e*j+Y-S-f/a/c/f0d0f0l1s3g4e4j5a6f7c7t&d^Ofgilmop!O!P!T!Y!Z![!_!`!c!p#e#o#p#q#r#s#t#u#v#w#x#y#z#}$T$W$e$g$h$q$r%_%b&U&Y&[&b&u'|(V(X(e(i(n(o(q({)h*e*h*j*k+Y+m,r,y-S-f-l.|/]/a/c/f0d0f0l0|1i1s1}3`3b3g3i3l4Y4e4j4x5O5a5l5y6b6f7X7c7l7t8S8_8`8r9T9]9o9|9}:O:P:R:S:V:W:X:Y:Z:[:]:^:_:`:a:b:l:m[!zY[$`$a9z9{['_!|!}(c(d:T:UW)n%i%j:g:hU,V'a-Q:nW.g)o)p:{:|T2p.h;PQ(h$eQ(l$gR-V(iV(g$e$g(iR-](jR-[(j$znOfgilmp!O!P!T!Y!Z![!_!`#e#o#p#q#r#s#t#u#v#w#x#y#z#}$T$W%_%b&Y&['|(V(X({)h*h*k+m,r,y-l.|/]0|1i1}3`3b3i3l4Y4x5O5l5y6b7X7l8S8_8`8r9T9]9o9|9}:O:P:R:S:V:W:X:Y:Z:[:]:^:_:`:a:b:l:m!i$ko!c!p$e$g$h$q$r&U&b&u(e(i(n(o(q*e*j+Y-S-f/a/c/f0d0f0l1s3g4e4j5a6f7c7t`,c'b'c'd,Y,e,f1`5WX3n/l/m3p3sh,b'b'c'd,Y,e,f/l/m1`3p3s5WQ7o5VR8z7p^uO!O!P!T%_%b)h$`#afgilmp!Y!Z![!_!`#e#o#p#q#r#s#t#u#v#w#x#y#z#}$T$W&Y&['|(V(X({*h*k+m,r,y-l.|/]0|1i1}3b3i3l4Y4x5O5l5y7X7l8S8r9T9|9}:O:P:R:S:V:W:X:Y:Z:[:]:^:_:`:a:b:l:mQ6`3`Q8]6bQ9X8_Q9Z8`Q9m9]R9s9oR(P#zR(O#zQ$SlR([$TR$ooR$noR)P$vR)O$vQ(}$vR2P-nhwOb!O!P!T$y%_%b)R)T)h-s$l!lz!Q!n$O$u$|%f%n%o%p%q%y%{&P&o&p&r'](p)r)u)w)x)}*O*Q*U*Z*]*d*m*v*w+T+U+g+n+|-h-u.T._.o.s.w.x/Y/Z/z/|0_0q0v0}1X1Y1z2b2i2k2n2t2w3W3w3}4O4T4W4a4g4v5e5i5{6W6_6u6{6}7h7y8nR${xR0a+TR0V*wR0T*wR7Q4RR/x*uR/w*uR0O*vR/}*vR0^+PR0]+P%XyObxz!O!P!Q!T!n$O$u$y$|%_%b%f%n%o%p%q%y%{&P&o&p&r'](p)R)T)h)r)u)w)x)}*O*Q*U*Z*]*d*m*v*w+T+U+g+n+|-h-s-u.T._.o.s.w.x/Y/Z/z/|0_0q0v0}1X1Y1z2b2i2k2n2t2w3W3w3}4O4T4W4a4g4v5e5i5{6W6_6u6{6}7h7y8nR0j+YR0i+YQ'R!qQ)b%XQ+v&|R4{1OX'Q!q%X&|1OR+x&|R+w&|T/R*R4VT/Q*R4VR.n)rR.m)rR)z%pR1f,kR1e,k",
  nodeNames: "\u26A0 | < > RawString Float LineComment BlockComment SourceFile ] InnerAttribute ! [ MetaItem self Metavariable super crate Identifier ScopedIdentifier :: QualifiedScope AbstractType impl SelfType MetaType TypeIdentifier ScopedTypeIdentifier ScopeIdentifier TypeArgList TypeBinding = Lifetime String Escape Char Boolean Integer } { Block ; ConstItem Vis pub ( in ) const BoundIdentifier : UnsafeBlock unsafe AsyncBlock async move IfExpression if LetDeclaration let LiteralPattern ArithOp MetaPattern SelfPattern ScopedIdentifier TuplePattern ScopedTypeIdentifier , StructPattern FieldPatternList FieldPattern ref mut FieldIdentifier .. RefPattern SlicePattern CapturedPattern ReferencePattern & MutPattern RangePattern ... OrPattern MacroPattern ParenthesizedTokens BracketedTokens BracedTokens TokenBinding Identifier TokenRepetition ArithOp BitOp LogicOp UpdateOp CompareOp -> => ArithOp _ else MatchExpression match MatchBlock MatchArm Attribute Guard UnaryExpression ArithOp DerefOp LogicOp ReferenceExpression TryExpression BinaryExpression ArithOp ArithOp BitOp BitOp BitOp BitOp LogicOp LogicOp AssignmentExpression TypeCastExpression as ReturnExpression return RangeExpression CallExpression ArgList AwaitExpression await FieldExpression GenericFunction BreakExpression break LoopLabel ContinueExpression continue IndexExpression ArrayExpression TupleExpression MacroInvocation UnitExpression ClosureExpression ParamList Parameter Parameter ParenthesizedExpression StructExpression FieldInitializerList ShorthandFieldInitializer FieldInitializer BaseFieldInitializer MatchArm WhileExpression while LoopExpression loop ForExpression for MacroInvocation MacroDefinition macro_rules MacroRule EmptyStatement ModItem mod DeclarationList AttributeItem ForeignModItem extern StructItem struct TypeParamList ConstrainedTypeParameter TraitBounds HigherRankedTraitBound RemovedTraitBound OptionalTypeParameter ConstParameter WhereClause where LifetimeClause TypeBoundClause FieldDeclarationList FieldDeclaration OrderedFieldDeclarationList UnionItem union EnumItem enum EnumVariantList EnumVariant TypeItem type FunctionItem default fn ParamList Parameter SelfParameter VariadicParameter VariadicParameter ImplItem TraitItem trait AssociatedType LetDeclaration UseDeclaration use ScopedIdentifier UseAsClause ScopedIdentifier UseList ScopedUseList UseWildcard ExternCrateDeclaration StaticItem static ExpressionStatement ExpressionStatement GenericType FunctionType ForLifetimes ParamList VariadicParameter Parameter VariadicParameter Parameter ReferenceType PointerType TupleType UnitType ArrayType MacroInvocation EmptyType DynamicType dyn BoundedType",
  maxTerm: 361,
  nodeProps: [
    [NodeProp.group, -42, 4, 5, 14, 15, 16, 17, 18, 19, 33, 35, 36, 37, 40, 51, 53, 56, 101, 107, 111, 112, 113, 122, 123, 125, 127, 128, 130, 132, 133, 134, 137, 139, 140, 141, 142, 143, 144, 148, 149, 155, 157, 159, "Expression", -16, 22, 24, 25, 26, 27, 222, 223, 230, 231, 232, 233, 234, 235, 236, 237, 239, "Type", -20, 42, 161, 162, 165, 166, 169, 170, 172, 188, 190, 194, 196, 204, 205, 207, 208, 209, 217, 218, 220, "Statement", -17, 49, 60, 62, 63, 64, 65, 68, 74, 75, 76, 77, 78, 80, 81, 83, 84, 99, "Pattern"],
    [NodeProp.openedBy, 9, "[", 38, "{", 47, "("],
    [NodeProp.closedBy, 12, "]", 39, "}", 45, ")"]
  ],
  skippedNodes: [0, 6, 7, 240],
  repeatNodeCount: 33,
  tokenData: "#CO_R!VOX$hXY1_YZ2ZZ]$h]^1_^p$hpq1_qr2srs4qst5Ztu6Vuv9lvw;jwx=nxy!#yyz!$uz{!%q{|!'k|}!(m}!O!)i!O!P!+j!P!Q!/f!Q!R!7q!R![!9f![!]!La!]!^!N_!^!_# Z!_!`##b!`!a#%c!a!b#'j!b!c#(f!c!}#)b!}#O#+X#O#P#,T#P#Q#4d#Q#R#5`#R#S#)b#S#T$h#T#U#)b#U#V#6b#V#f#)b#f#g#9u#g#o#)b#o#p#?S#p#q#@O#q#r#BS#r${$h${$|#)b$|4w$h4w5b#)b5b5i$h5i6S#)b6S~$hU$oZ'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q#O$h#O#P%x#P~$hU%iT'`Q'PSOz%xz{&^{!P%x!P!Q'S!Q~%xS%}T'PSOz%xz{&^{!P%x!P!Q'S!Q~%xS&aTOz&pz{&^{!P&p!P!Q({!Q~&pS&sTOz%xz{&^{!P%x!P!Q'S!Q~%xS'VSOz&p{!P&p!P!Q'c!Q~&pS'fSOz'r{!P'r!P!Q'c!Q~'rS'uTOz(Uz{(l{!P(U!P!Q'c!Q~(US(]T'QS'PSOz(Uz{(l{!P(U!P!Q'c!Q~(US(oSOz'rz{(l{!P'r!Q~'rS)QO'QSU)VZ'`QOY)xYZ+hZr)xrs&psz)xz{)Q{!P)x!P!Q0w!Q#O)x#O#P&p#P~)xU)}Z'`QOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q#O$h#O#P%x#P~$hU*uZ'`QOY)xYZ+hZr)xrs&psz)xz{+|{!P)x!P!Q,g!Q#O)x#O#P&p#P~)xU+mT'`QOz%xz{&^{!P%x!P!Q'S!Q~%xQ,RT'`QOY+|YZ,bZr+|s#O+|#P~+|Q,gO'`QU,lZ'`QOY-_YZ0cZr-_rs'rsz-_z{+|{!P-_!P!Q,g!Q#O-_#O#P'r#P~-_U-dZ'`QOY.VYZ/RZr.Vrs(Usz.Vz{/k{!P.V!P!Q,g!Q#O.V#O#P(U#P~.VU.`Z'`Q'QS'PSOY.VYZ/RZr.Vrs(Usz.Vz{/k{!P.V!P!Q,g!Q#O.V#O#P(U#P~.VU/[T'`Q'QS'PSOz(Uz{(l{!P(U!P!Q'c!Q~(UU/pZ'`QOY-_YZ0cZr-_rs'rsz-_z{/k{!P-_!P!Q+|!Q#O-_#O#P'r#P~-_U0hT'`QOz(Uz{(l{!P(U!P!Q'c!Q~(UU1OT'`Q'QSOY+|YZ,bZr+|s#O+|#P~+|_1hZ'`Q&}X'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q#O$h#O#P%x#P~$h_2dT'`Q&}X'PSOz%xz{&^{!P%x!P!Q'S!Q~%x_2|]ZX'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q!_$h!_!`3u!`#O$h#O#P%x#P~$h_4OZ#RX'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q#O$h#O#P%x#P~$h_4zT'^Q'PS'_XOz%xz{&^{!P%x!P!Q'S!Q~%x_5dZ'RX'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q#O$h#O#P%x#P~$h_6`g'`Q'vW'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q!c$h!c!}7w!}#O$h#O#P%x#P#R$h#R#S7w#S#T$h#T#o7w#o${$h${$|7w$|4w$h4w5b7w5b5i$h5i6S7w6S~$h_8Qh'`Q_X'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q![7w![!c$h!c!}7w!}#O$h#O#P%x#P#R$h#R#S7w#S#T$h#T#o7w#o${$h${$|7w$|4w$h4w5b7w5b5i$h5i6S7w6S~$h_9u](TP'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q!_$h!_!`:n!`#O$h#O#P%x#P~$h_:wZ#QX'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q#O$h#O#P%x#P~$h_;s_!qX'`Q'PSOY$hYZ%bZr$hrs%xsv$hvw<rwz$hz{)Q{!P$h!P!Q*p!Q!_$h!_!`:n!`#O$h#O#P%x#P~$h_<{Z'}X'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q#O$h#O#P%x#P~$h_=ui'`Q'PSOY?dYZA`Zr?drsBdsw?dwx@dxz?dz{CO{!P?d!P!QDv!Q!c?d!c!}Et!}#O?d#O#PId#P#R?d#R#SEt#S#T?d#T#oEt#o${?d${$|Et$|4w?d4w5bEt5b5i?d5i6SEt6S~?d_?k]'`Q'PSOY$hYZ%bZr$hrs%xsw$hwx@dxz$hz{)Q{!P$h!P!Q*p!Q#O$h#O#P%x#P~$h_@mZ'`Q'PSsXOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q#O$h#O#P%x#P~$h_AgV'`Q'PSOw%xwxA|xz%xz{&^{!P%x!P!Q'S!Q~%x]BTT'PSsXOz%xz{&^{!P%x!P!Q'S!Q~%x]BiV'PSOw%xwxA|xz%xz{&^{!P%x!P!Q'S!Q~%x_CT]'`QOY)xYZ+hZr)xrs&psw)xwxC|xz)xz{)Q{!P)x!P!Q0w!Q#O)x#O#P&p#P~)x_DTZ'`QsXOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q#O$h#O#P%x#P~$h_D{]'`QOY)xYZ+hZr)xrs&psw)xwxC|xz)xz{+|{!P)x!P!Q,g!Q#O)x#O#P&p#P~)x_E}j'`Q'PS'[XOY$hYZ%bZr$hrs%xsw$hwx@dxz$hz{)Q{!P$h!P!Q*p!Q![Go![!c$h!c!}Go!}#O$h#O#P%x#P#R$h#R#SGo#S#T$h#T#oGo#o${$h${$|Go$|4w$h4w5bGo5b5i$h5i6SGo6S~$h_Gxh'`Q'PS'[XOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q![Go![!c$h!c!}Go!}#O$h#O#P%x#P#R$h#R#SGo#S#T$h#T#oGo#o${$h${$|Go$|4w$h4w5bGo5b5i$h5i6SGo6S~$h]IiX'PSOzBdz{JU{!PBd!P!QKS!Q#iBd#i#jKi#j#lBd#l#m!!a#m~Bd]JXVOw&pwxJnxz&pz{&^{!P&p!P!Q({!Q~&p]JsTsXOz%xz{&^{!P%x!P!Q'S!Q~%x]KVUOw&pwxJnxz&p{!P&p!P!Q'c!Q~&p]Kn['PSOz%xz{&^{!P%x!P!Q'S!Q![Ld![!c%x!c!iLd!i#T%x#T#ZLd#Z#o%x#o#pNq#p~%x]LiY'PSOz%xz{&^{!P%x!P!Q'S!Q![MX![!c%x!c!iMX!i#T%x#T#ZMX#Z~%x]M^Y'PSOz%xz{&^{!P%x!P!Q'S!Q![M|![!c%x!c!iM|!i#T%x#T#ZM|#Z~%x]NRY'PSOz%xz{&^{!P%x!P!Q'S!Q![Bd![!c%x!c!iBd!i#T%x#T#ZBd#Z~%x]NvY'PSOz%xz{&^{!P%x!P!Q'S!Q![! f![!c%x!c!i! f!i#T%x#T#Z! f#Z~%x]! k['PSOz%xz{&^{!P%x!P!Q'S!Q![! f![!c%x!c!i! f!i#T%x#T#Z! f#Z#q%x#q#rBd#r~%x]!!fY'PSOz%xz{&^{!P%x!P!Q'S!Q![!#U![!c%x!c!i!#U!i#T%x#T#Z!#U#Z~%x]!#ZY'PSOz%xz{&^{!P%x!P!Q'S!Q![Bd![!c%x!c!iBd!i#T%x#T#ZBd#Z~%x_!$SZ}X'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q#O$h#O#P%x#P~$h_!%OZ!PX'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q#O$h#O#P%x#P~$h_!%x](QX'`QOY)xYZ+hZr)xrs&psz)xz{)Q{!P)x!P!Q0w!Q!_)x!_!`!&q!`#O)x#O#P&p#P~)x_!&xZ#QX'`QOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q#O$h#O#P%x#P~$h_!'t](PX'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q!_$h!_!`:n!`#O$h#O#P%x#P~$h_!(vZ!eX'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q#O$h#O#P%x#P~$h_!)r^'hX'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q!_$h!_!`:n!`!a!*n!a#O$h#O#P%x#P~$h_!*wZ#SX'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q#O$h#O#P%x#P~$h_!+s[(OX'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!O$h!O!P!,i!P!Q*p!Q#O$h#O#P%x#P~$h_!,r^!lX'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!O$h!O!P!-n!P!Q*p!Q!_$h!_!`!.j!`#O$h#O#P%x#P~$h_!-wZ!tX'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q#O$h#O#P%x#P~$hV!.sZ'rP'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q#O$h#O#P%x#P~$h_!/m]'`Q'xXOY)xYZ+hZr)xrs&psz)xz{!0f{!P)x!P!Q!0|!Q!_)x!_!`!&q!`#O)x#O#P&p#P~)x_!0mT'O]'`QOY+|YZ,bZr+|s#O+|#P~+|_!1TZ'`QUXOY!1vYZ0cZr!1vrs!4xsz!1vz{!7T{!P!1v!P!Q!0|!Q#O!1v#O#P!4x#P~!1v_!1}Z'`QUXOY!2pYZ/RZr!2prs!3nsz!2pz{!6Z{!P!2p!P!Q!0|!Q#O!2p#O#P!3n#P~!2p_!2{Z'`QUX'QS'PSOY!2pYZ/RZr!2prs!3nsz!2pz{!6Z{!P!2p!P!Q!0|!Q#O!2p#O#P!3n#P~!2p]!3wVUX'QS'PSOY!3nYZ(UZz!3nz{!4^{!P!3n!P!Q!5d!Q~!3n]!4cVUXOY!4xYZ'rZz!4xz{!4^{!P!4x!P!Q!6O!Q~!4x]!4}VUXOY!3nYZ(UZz!3nz{!4^{!P!3n!P!Q!5d!Q~!3n]!5iVUXOY!4xYZ'rZz!4xz{!6O{!P!4x!P!Q!5d!Q~!4xX!6TQUXOY!6OZ~!6O_!6bZ'`QUXOY!1vYZ0cZr!1vrs!4xsz!1vz{!6Z{!P!1v!P!Q!7T!Q#O!1v#O#P!4x#P~!1vZ!7[V'`QUXOY!7TYZ,bZr!7Trs!6Os#O!7T#O#P!6O#P~!7T_!7zhuX'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q![!9f![#O$h#O#P%x#P#R$h#R#S!9f#S#U$h#U#V!Dc#V#]$h#]#^!:w#^#c$h#c#d!F}#d#i$h#i#j!:w#j#l$h#l#m!Ic#m~$h_!9obuX'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q![!9f![#O$h#O#P%x#P#R$h#R#S!9f#S#]$h#]#^!:w#^#i$h#i#j!:w#j~$h_!;Oe'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q!R$h!R!S!<a!S!T$h!T!U!?c!U!W$h!W!X!@c!X!Y$h!Y!Z!>g!Z#O$h#O#P%x#P#g$h#g#h!Ac#h~$h_!<h_'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q!S$h!S!T!=g!T!W$h!W!X!>g!X#O$h#O#P%x#P~$h_!=n]'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q!Y$h!Y!Z!>g!Z#O$h#O#P%x#P~$h_!>pZuX'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q#O$h#O#P%x#P~$h_!?j]'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q!S$h!S!T!>g!T#O$h#O#P%x#P~$h_!@j]'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q!U$h!U!V!>g!V#O$h#O#P%x#P~$h_!Aj]'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q#O$h#O#P%x#P#]$h#]#^!Bc#^~$h_!Bj]'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q#O$h#O#P%x#P#n$h#n#o!Cc#o~$h_!Cj]'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q#O$h#O#P%x#P#X$h#X#Y!>g#Y~$h_!Dj_'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q!R!Ei!R!S!Ei!S#O$h#O#P%x#P#R$h#R#S!Ei#S~$h_!ErcuX'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q!R!Ei!R!S!Ei!S#O$h#O#P%x#P#R$h#R#S!Ei#S#]$h#]#^!:w#^#i$h#i#j!:w#j~$h_!GU^'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q!Y!HQ!Y#O$h#O#P%x#P#R$h#R#S!HQ#S~$h_!HZbuX'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q!Y!HQ!Y#O$h#O#P%x#P#R$h#R#S!HQ#S#]$h#]#^!:w#^#i$h#i#j!:w#j~$h_!Ijb'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q![!Jr![!c$h!c!i!Jr!i#O$h#O#P%x#P#R$h#R#S!Jr#S#T$h#T#Z!Jr#Z~$h_!J{fuX'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q![!Jr![!c$h!c!i!Jr!i#O$h#O#P%x#P#R$h#R#S!Jr#S#T$h#T#Z!Jr#Z#]$h#]#^!:w#^#i$h#i#j!:w#j~$h_!Lj]!SX'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q![$h![!]!Mc!]#O$h#O#P%x#P~$h_!MlZdX'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q#O$h#O#P%x#P~$h_!NhZyX'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q#O$h#O#P%x#P~$h_# d^#RX'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q!^$h!^!_#!`!_!`3u!`#O$h#O#P%x#P~$h_#!i]'yX'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q!_$h!_!`:n!`#O$h#O#P%x#P~$h_##k^oX'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q!_$h!_!`3u!`!a#$g!a#O$h#O#P%x#P~$h_#$pZ#TX'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q#O$h#O#P%x#P~$h_#%l^#RX'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q!_$h!_!`3u!`!a#&h!a#O$h#O#P%x#P~$h_#&q]'zX'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q!_$h!_!`:n!`#O$h#O#P%x#P~$h_#'sZ(RX'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q#O$h#O#P%x#P~$hV#(oZ'qP'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q#O$h#O#P%x#P~$h_#)mh'`Q'PS!{W'UPOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q![#)b![!c$h!c!}#)b!}#O$h#O#P%x#P#R$h#R#S#)b#S#T$h#T#o#)b#o${$h${$|#)b$|4w$h4w5b#)b5b5i$h5i6S#)b6S~$h_#+bZ[X'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q#O$h#O#P%x#P~$hU#,YX'PSOz#,uz{#-]{!P#,u!P!Q#-q!Q#i#,u#i#j#.S#j#l#,u#l#m#2z#m~#,uU#,|TrQ'PSOz%xz{&^{!P%x!P!Q'S!Q~%xU#-bTrQOz&pz{&^{!P&p!P!Q({!Q~&pU#-vSrQOz&p{!P&p!P!Q'c!Q~&pU#.X['PSOz%xz{&^{!P%x!P!Q'S!Q![#.}![!c%x!c!i#.}!i#T%x#T#Z#.}#Z#o%x#o#p#1[#p~%xU#/SY'PSOz%xz{&^{!P%x!P!Q'S!Q![#/r![!c%x!c!i#/r!i#T%x#T#Z#/r#Z~%xU#/wY'PSOz%xz{&^{!P%x!P!Q'S!Q![#0g![!c%x!c!i#0g!i#T%x#T#Z#0g#Z~%xU#0lY'PSOz%xz{&^{!P%x!P!Q'S!Q![#,u![!c%x!c!i#,u!i#T%x#T#Z#,u#Z~%xU#1aY'PSOz%xz{&^{!P%x!P!Q'S!Q![#2P![!c%x!c!i#2P!i#T%x#T#Z#2P#Z~%xU#2U['PSOz%xz{&^{!P%x!P!Q'S!Q![#2P![!c%x!c!i#2P!i#T%x#T#Z#2P#Z#q%x#q#r#,u#r~%xU#3PY'PSOz%xz{&^{!P%x!P!Q'S!Q![#3o![!c%x!c!i#3o!i#T%x#T#Z#3o#Z~%xU#3tY'PSOz%xz{&^{!P%x!P!Q'S!Q![#,u![!c%x!c!i#,u!i#T%x#T#Z#,u#Z~%x_#4mZXX'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q#O$h#O#P%x#P~$h_#5i]'{X'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q!_$h!_!`:n!`#O$h#O#P%x#P~$h_#6mj'`Q'PS!{W'UPOY$hYZ%bZr$hrs#8_sw$hwx#8uxz$hz{)Q{!P$h!P!Q*p!Q![#)b![!c$h!c!}#)b!}#O$h#O#P%x#P#R$h#R#S#)b#S#T$h#T#o#)b#o${$h${$|#)b$|4w$h4w5b#)b5b5i$h5i6S#)b6S~$h]#8fT'PS'_XOz%xz{&^{!P%x!P!Q'S!Q~%x_#8|]'`Q'PSOY?dYZA`Zr?drsBdsw?dwx@dxz?dz{CO{!P?d!P!QDv!Q#O?d#O#PId#P~?d_#:Qi'`Q'PS!{W'UPOY$hYZ%bZr$hrs%xst#;otz$hz{)Q{!P$h!P!Q*p!Q![#)b![!c$h!c!}#)b!}#O$h#O#P%x#P#R$h#R#S#)b#S#T$h#T#o#)b#o${$h${$|#)b$|4w$h4w5b#)b5b5i$h5i6S#)b6S~$hV#;vg'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q!c$h!c!}#=_!}#O$h#O#P%x#P#R$h#R#S#=_#S#T$h#T#o#=_#o${$h${$|#=_$|4w$h4w5b#=_5b5i$h5i6S#=_6S~$hV#=hh'`Q'PS'UPOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q![#=_![!c$h!c!}#=_!}#O$h#O#P%x#P#R$h#R#S#=_#S#T$h#T#o#=_#o${$h${$|#=_$|4w$h4w5b#=_5b5i$h5i6S#=_6S~$h_#?]ZwX'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q#O$h#O#P%x#P~$h_#@X_'sX'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q!_$h!_!`:n!`#O$h#O#P%x#P#p$h#p#q#AW#q~$h_#AaZ'|X'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q#O$h#O#P%x#P~$h_#B]ZvX'`Q'PSOY$hYZ%bZr$hrs%xsz$hz{)Q{!P$h!P!Q*p!Q#O$h#O#P%x#P~$h",
  tokenizers: [closureParam, tpDelim, literalTokens, 0, 1, 2, 3],
  topRules: {SourceFile: [0, 8]},
  specialized: [{term: 282, get: (value) => spec_identifier4[value] || -1}],
  tokenPrec: 15890
});

// node_modules/@codemirror/lang-rust/dist/index.js
var rustLanguage = LezerLanguage.define({
  parser: parser4.configure({
    props: [
      indentNodeProp.add({
        IfExpression: continuedIndent({except: /^\s*({|else\b)/}),
        "String BlockComment": () => -1,
        "Statement MatchArm": continuedIndent()
      }),
      foldNodeProp.add((type2) => {
        if (/(Block|edTokens|List)$/.test(type2.name))
          return foldInside;
        if (type2.name == "BlockComment")
          return (tree) => ({from: tree.from + 2, to: tree.to - 2});
        return void 0;
      }),
      styleTags({
        "const macro_rules mod struct union enum type fn impl trait let use crate static": tags.definitionKeyword,
        "pub unsafe async mut extern default move": tags.modifier,
        "for if else loop while match continue break return await": tags.controlKeyword,
        "as in ref": tags.operatorKeyword,
        "where _ crate super dyn": tags.keyword,
        self: tags.self,
        String: tags.string,
        RawString: tags.special(tags.string),
        Boolean: tags.bool,
        Identifier: tags.variableName,
        "CallExpression/Identifier": tags.function(tags.variableName),
        BoundIdentifier: tags.definition(tags.variableName),
        LoopLabel: tags.labelName,
        FieldIdentifier: tags.propertyName,
        "CallExpression/FieldExpression/FieldIdentifier": tags.function(tags.propertyName),
        Lifetime: tags.special(tags.variableName),
        ScopeIdentifier: tags.namespace,
        TypeIdentifier: tags.typeName,
        "MacroInvocation/Identifier MacroInvocation/ScopedIdentifier/Identifier": tags.macroName,
        "MacroInvocation/TypeIdentifier MacroInvocation/ScopedIdentifier/TypeIdentifier": tags.macroName,
        '"!"': tags.macroName,
        UpdateOp: tags.updateOperator,
        LineComment: tags.lineComment,
        BlockComment: tags.blockComment,
        Integer: tags.integer,
        Float: tags.float,
        ArithOp: tags.arithmeticOperator,
        LogicOp: tags.logicOperator,
        BitOp: tags.bitwiseOperator,
        CompareOp: tags.compareOperator,
        "=": tags.definitionOperator,
        ".. ... => ->": tags.punctuation,
        "( )": tags.paren,
        "[ ]": tags.squareBracket,
        "{ }": tags.brace,
        ".": tags.derefOperator,
        "&": tags.operator,
        ", ; ::": tags.separator
      })
    ]
  }),
  languageData: {
    commentTokens: {line: "//", block: {open: "/*", close: "*/"}},
    indentOnInput: /^\s*(?:\{|\})$/
  }
});
function rust() {
  return new LanguageSupport(rustLanguage);
}

// node_modules/@codemirror/legacy-modes/mode/shell.js
var words3 = {};
function define(style, dict) {
  for (var i = 0; i < dict.length; i++) {
    words3[dict[i]] = style;
  }
}
var commonAtoms = ["true", "false"];
var commonKeywords = [
  "if",
  "then",
  "do",
  "else",
  "elif",
  "while",
  "until",
  "for",
  "in",
  "esac",
  "fi",
  "fin",
  "fil",
  "done",
  "exit",
  "set",
  "unset",
  "export",
  "function"
];
var commonCommands = [
  "ab",
  "awk",
  "bash",
  "beep",
  "cat",
  "cc",
  "cd",
  "chown",
  "chmod",
  "chroot",
  "clear",
  "cp",
  "curl",
  "cut",
  "diff",
  "echo",
  "find",
  "gawk",
  "gcc",
  "get",
  "git",
  "grep",
  "hg",
  "kill",
  "killall",
  "ln",
  "ls",
  "make",
  "mkdir",
  "openssl",
  "mv",
  "nc",
  "nl",
  "node",
  "npm",
  "ping",
  "ps",
  "restart",
  "rm",
  "rmdir",
  "sed",
  "service",
  "sh",
  "shopt",
  "shred",
  "source",
  "sort",
  "sleep",
  "ssh",
  "start",
  "stop",
  "su",
  "sudo",
  "svn",
  "tee",
  "telnet",
  "top",
  "touch",
  "vi",
  "vim",
  "wall",
  "wc",
  "wget",
  "who",
  "write",
  "yes",
  "zsh"
];
define("atom", commonAtoms);
define("keyword", commonKeywords);
define("builtin", commonCommands);
function tokenBase7(stream, state) {
  if (stream.eatSpace())
    return null;
  var sol = stream.sol();
  var ch = stream.next();
  if (ch === "\\") {
    stream.next();
    return null;
  }
  if (ch === "'" || ch === '"' || ch === "`") {
    state.tokens.unshift(tokenString3(ch, ch === "`" ? "quote" : "string"));
    return tokenize(stream, state);
  }
  if (ch === "#") {
    if (sol && stream.eat("!")) {
      stream.skipToEnd();
      return "meta";
    }
    stream.skipToEnd();
    return "comment";
  }
  if (ch === "$") {
    state.tokens.unshift(tokenDollar);
    return tokenize(stream, state);
  }
  if (ch === "+" || ch === "=") {
    return "operator";
  }
  if (ch === "-") {
    stream.eat("-");
    stream.eatWhile(/\w/);
    return "attribute";
  }
  if (ch == "<") {
    if (stream.match("<<"))
      return "operator";
    var heredoc = stream.match(/^<-?\s*['"]?([^'"]*)['"]?/);
    if (heredoc) {
      state.tokens.unshift(tokenHeredoc(heredoc[1]));
      return "string.special";
    }
  }
  if (/\d/.test(ch)) {
    stream.eatWhile(/\d/);
    if (stream.eol() || !/\w/.test(stream.peek())) {
      return "number";
    }
  }
  stream.eatWhile(/[\w-]/);
  var cur2 = stream.current();
  if (stream.peek() === "=" && /\w+/.test(cur2))
    return "def";
  return words3.hasOwnProperty(cur2) ? words3[cur2] : null;
}
function tokenString3(quote, style) {
  var close = quote == "(" ? ")" : quote == "{" ? "}" : quote;
  return function(stream, state) {
    var next, escaped = false;
    while ((next = stream.next()) != null) {
      if (next === close && !escaped) {
        state.tokens.shift();
        break;
      } else if (next === "$" && !escaped && quote !== "'" && stream.peek() != close) {
        escaped = true;
        stream.backUp(1);
        state.tokens.unshift(tokenDollar);
        break;
      } else if (!escaped && quote !== close && next === quote) {
        state.tokens.unshift(tokenString3(quote, style));
        return tokenize(stream, state);
      } else if (!escaped && /['"]/.test(next) && !/['"]/.test(quote)) {
        state.tokens.unshift(tokenStringStart(next, "string"));
        stream.backUp(1);
        break;
      }
      escaped = !escaped && next === "\\";
    }
    return style;
  };
}
function tokenStringStart(quote, style) {
  return function(stream, state) {
    state.tokens[0] = tokenString3(quote, style);
    stream.next();
    return tokenize(stream, state);
  };
}
var tokenDollar = function(stream, state) {
  if (state.tokens.length > 1)
    stream.eat("$");
  var ch = stream.next();
  if (/['"({]/.test(ch)) {
    state.tokens[0] = tokenString3(ch, ch == "(" ? "quote" : ch == "{" ? "def" : "string");
    return tokenize(stream, state);
  }
  if (!/\d/.test(ch))
    stream.eatWhile(/\w/);
  state.tokens.shift();
  return "def";
};
function tokenHeredoc(delim) {
  return function(stream, state) {
    if (stream.sol() && stream.string == delim)
      state.tokens.shift();
    stream.skipToEnd();
    return "string.special";
  };
}
function tokenize(stream, state) {
  return (state.tokens[0] || tokenBase7)(stream, state);
}
var shell = {
  startState: function() {
    return {tokens: []};
  },
  token: function(stream, state) {
    return tokenize(stream, state);
  },
  languageData: {
    autocomplete: commonAtoms.concat(commonKeywords, commonCommands),
    closeBrackets: {brackets: ["(", "[", "{", "'", '"', "`"]},
    commentTokens: {line: "#"}
  }
};

// node_modules/@codemirror/lang-sql/dist/index.js
var whitespace = 34;
var LineComment = 1;
var BlockComment = 2;
var String2 = 3;
var Number2 = 4;
var Bool = 5;
var Null = 6;
var ParenL = 7;
var ParenR = 8;
var BraceL = 9;
var BraceR = 10;
var BracketL = 11;
var BracketR = 12;
var Semi = 13;
var Dot2 = 14;
var Operator = 15;
var Punctuation = 16;
var SpecialVar = 17;
var Identifier = 18;
var QuotedIdentifier = 19;
var Keyword = 20;
var Type = 21;
var Builtin = 22;
function isAlpha(ch) {
  return ch >= 65 && ch <= 90 || ch >= 97 && ch <= 122 || ch >= 48 && ch <= 57;
}
function isHexDigit(ch) {
  return ch >= 48 && ch <= 57 || ch >= 97 && ch <= 102 || ch >= 65 && ch <= 70;
}
function readLiteral(input, pos, endQuote, backslashEscapes) {
  for (let escaped = false; ; ) {
    let next = input.get(pos++);
    if (next < 0)
      return pos - 1;
    if (next == endQuote && !escaped)
      return pos;
    escaped = backslashEscapes && !escaped && next == 92;
  }
}
function readWord(input, pos) {
  for (; ; pos++) {
    let next = input.get(pos);
    if (next != 95 && !isAlpha(next))
      break;
  }
  return pos;
}
function readWordOrQuoted(input, pos) {
  let next = input.get(pos);
  if (next == 39 || next == 34 || next == 96)
    return readLiteral(input, pos + 1, next, false);
  return readWord(input, pos);
}
function readNumber(input, pos, sawDot) {
  let next;
  for (; ; pos++) {
    next = input.get(pos);
    if (next == 46) {
      if (sawDot)
        break;
      sawDot = true;
    } else if (next < 48 || next > 57) {
      break;
    }
  }
  if (next == 69 || next == 101) {
    next = input.get(++pos);
    if (next == 43 || next == 45)
      pos++;
    for (; ; pos++) {
      next = input.get(pos);
      if (next < 48 || next > 57)
        break;
    }
  }
  return pos;
}
function eol(input, pos) {
  for (; ; pos++) {
    let next = input.get(pos);
    if (next < 0 || next == 10)
      return pos;
  }
}
function inString2(ch, str) {
  for (let i = 0; i < str.length; i++)
    if (str.charCodeAt(i) == ch)
      return true;
  return false;
}
var Space = " 	\r\n";
function keywords9(keywords11, types4, builtin) {
  let result = Object.create(null);
  result["true"] = result["false"] = Bool;
  result["null"] = result["unknown"] = Null;
  for (let kw of keywords11.split(" "))
    if (kw)
      result[kw] = Keyword;
  for (let tp of types4.split(" "))
    if (tp)
      result[tp] = Type;
  for (let kw of (builtin || "").split(" "))
    if (kw)
      result[kw] = Builtin;
  return result;
}
var SQLTypes = "array binary bit boolean char character clob date decimal double float int integer interval large national nchar nclob numeric object precision real smallint time timestamp varchar varying ";
var SQLKeywords = "absolute action add after all allocate alter and any are as asc assertion at authorization before begin between blob both breadth by call cascade cascaded case cast catalog check close collate collation column commit condition connect connection constraint constraints constructor continue corresponding count create cross cube current current_date current_default_transform_group current_transform_group_for_type current_path current_role current_time current_timestamp current_user cursor cycle data day deallocate dec declare default deferrable deferred delete depth deref desc describe descriptor deterministic diagnostics disconnect distinct do domain drop dynamic each else elseif end end-exec equals escape except exception exec execute exists exit external fetch first for foreign found from free full function general get global go goto grant group grouping handle having hold hour identity if immediate in indicator initially inner inout input insert intersect into is isolation join key language last lateral leading leave left level like limit local localtime localtimestamp locator loop map match method minute modifies module month names natural nesting new next no none not of old on only open option or order ordinality out outer output overlaps pad parameter partial path prepare preserve primary prior privileges procedure public read reads recursive redo ref references referencing relative release repeat resignal restrict result return returns revoke right role rollback rollup routine row rows savepoint schema scroll search second section select session session_user set sets signal similar size some space specific specifictype sql sqlexception sqlstate sqlwarning start state static system_user table temporary then timezone_hour timezone_minute to trailing transaction translation treat trigger under undo union unique unnest until update usage user using value values view when whenever where while with without work write year zone ";
var defaults3 = {
  backslashEscapes: false,
  hashComments: false,
  spaceAfterDashes: false,
  slashComments: false,
  doubleQuotedStrings: false,
  charSetCasts: false,
  operatorChars: "*+-%<>!=&|~^/",
  specialVar: "?",
  identifierQuotes: '"',
  words: keywords9(SQLKeywords, SQLTypes)
};
function dialect(spec, kws, types4, builtin) {
  let dialect2 = {};
  for (let prop in defaults3)
    dialect2[prop] = (spec.hasOwnProperty(prop) ? spec : defaults3)[prop];
  if (kws)
    dialect2.words = keywords9(kws, types4 || "", builtin);
  return dialect2;
}
function tokensFor(d) {
  return new ExternalTokenizer((input, token) => {
    var _a;
    let pos = token.start, next = input.get(pos++), next2 = input.get(pos);
    if (inString2(next, Space)) {
      while (inString2(input.get(pos), Space))
        pos++;
      token.accept(whitespace, pos);
    } else if (next == 39 || next == 34 && d.doubleQuotedStrings) {
      token.accept(String2, readLiteral(input, pos, next, d.backslashEscapes));
    } else if (next == 35 && d.hashComments || next == 47 && next2 == 47 && d.slashComments) {
      token.accept(LineComment, eol(input, pos));
    } else if (next == 45 && next2 == 45 && (!d.spaceAfterDashes || input.get(pos + 1) == 32)) {
      token.accept(LineComment, eol(input, pos + 1));
    } else if (next == 47 && next2 == 42) {
      pos++;
      for (let prev = -1, depth = 1; ; ) {
        let next3 = input.get(pos++);
        if (next3 < 0) {
          pos--;
          break;
        } else if (prev == 42 && next3 == 47) {
          depth--;
          if (!depth)
            break;
          next3 = -1;
        } else if (prev == 47 && next3 == 42) {
          depth++;
          next3 = -1;
        }
        prev = next3;
      }
      token.accept(BlockComment, pos);
    } else if ((next == 101 || next == 69) && next2 == 39) {
      token.accept(String2, readLiteral(input, pos + 1, 39, true));
    } else if ((next == 110 || next == 78) && next2 == 39 && d.charSetCasts) {
      token.accept(String2, readLiteral(input, pos + 1, 39, d.backslashEscapes));
    } else if (next == 95 && d.charSetCasts) {
      for (; ; ) {
        let next3 = input.get(pos++);
        if (next3 == 39 && pos > token.start + 2) {
          token.accept(String2, readLiteral(input, pos, 39, d.backslashEscapes));
          break;
        }
        if (!isAlpha(next3))
          break;
      }
    } else if (next == 40) {
      token.accept(ParenL, pos);
    } else if (next == 41) {
      token.accept(ParenR, pos);
    } else if (next == 123) {
      token.accept(BraceL, pos);
    } else if (next == 125) {
      token.accept(BraceR, pos);
    } else if (next == 91) {
      token.accept(BracketL, pos);
    } else if (next == 93) {
      token.accept(BracketR, pos);
    } else if (next == 59) {
      token.accept(Semi, pos);
    } else if (next == 48 && (next2 == 98 || next2 == 66) || (next == 98 || next == 66) && next2 == 39) {
      let quoted = next2 == 39;
      pos++;
      while ((next = input.get(pos)) == 48 || next == 49)
        pos++;
      if (quoted && next == 39)
        pos++;
      token.accept(Number2, pos);
    } else if (next == 48 && (next2 == 120 || next2 == 88) || (next == 120 || next == 88) && next2 == 39) {
      let quoted = next2 == 39;
      pos++;
      while (isHexDigit(next = input.get(pos)))
        pos++;
      if (quoted && next == 39)
        pos++;
      token.accept(Number2, pos);
    } else if (next == 46 && next2 >= 48 && next2 <= 57) {
      token.accept(Number2, readNumber(input, pos + 1, true));
    } else if (next == 46) {
      token.accept(Dot2, pos);
    } else if (next >= 48 && next <= 57) {
      token.accept(Number2, readNumber(input, pos, false));
    } else if (inString2(next, d.operatorChars)) {
      while (inString2(input.get(pos), d.operatorChars))
        pos++;
      token.accept(Operator, pos);
    } else if (inString2(next, d.specialVar)) {
      token.accept(SpecialVar, readWordOrQuoted(input, next2 == next ? pos + 1 : pos));
    } else if (inString2(next, d.identifierQuotes)) {
      token.accept(QuotedIdentifier, readLiteral(input, pos, next, false));
    } else if (next == 58 || next == 44) {
      token.accept(Punctuation, pos);
    } else if (isAlpha(next)) {
      pos = readWord(input, pos);
      token.accept((_a = d.words[input.read(token.start, pos).toLowerCase()]) !== null && _a !== void 0 ? _a : Identifier, pos);
    }
  });
}
var tokens = tokensFor(defaults3);
var parser5 = Parser.deserialize({
  version: 13,
  states: "%dQ]QQOOO#kQRO'#DQO#rQQO'#CuO%RQQO'#CvO%YQQO'#CwO%aQQO'#CxOOQQ'#DQ'#DQOOQQ'#C{'#C{O&lQRO'#CyOOQQ'#Ct'#CtOOQQ'#Cz'#CzQ]QQOOQOQQOOO&vQQO,59aO'RQQO,59aO'WQQO'#DQOOQQ,59b,59bO'eQQO,59bOOQQ,59c,59cO'lQQO,59cOOQQ,59d,59dO'sQQO,59dOOQQ-E6y-E6yOOQQ,59`,59`OOQQ-E6x-E6xOOQQ'#C|'#C|OOQQ1G.{1G.{O&vQQO1G.{OOQQ1G.|1G.|OOQQ1G.}1G.}OOQQ1G/O1G/OP'zQQO'#C{POQQ-E6z-E6zOOQQ7+$g7+$g",
  stateData: "(R~OrOSPOSQOS~ORUOSUOTUOUUOVROXSOZTO]XO^QO_UO`UOaPObPOcPOdUOeUOfUO~O^]ORtXStXTtXUtXVtXXtXZtX]tX_tX`tXatXbtXctXdtXetXftX~OqtX~P!dOa^Ob^Oc^O~ORUOSUOTUOUUOVROXSOZTO^QO_UO`UOa_Ob_Oc_OdUOeUOfUO~OW`O~P#}OYbO~P#}O[dO~P#}ORUOSUOTUOUUOVROXSOZTO^QO_UO`UOaPObPOcPOdUOeUOfUO~O]gOqmX~P%hOaiObiOciO~O^kO~OWtXYtX[tX~P!dOWlO~P#}OYmO~P#}O[nO~P#}O]gO~P#}O",
  goto: "#YuPPPPPPPPPPPPPPPPPPPPPPPPvzzzz!W![!b!vPPP!|TYOZeUORSTWZaceoT[OZQZORhZSWOZQaRQcSQeTZfWaceoQj]RqkeVORSTWZaceo",
  nodeNames: "\u26A0 LineComment BlockComment String Number Bool Null ( ) [ ] { } ; . Operator Punctuation SpecialVar Identifier QuotedIdentifier Keyword Type Builtin Script Statement CompositeIdentifier Parens Braces Brackets Statement",
  maxTerm: 36,
  skippedNodes: [0, 1, 2],
  repeatNodeCount: 3,
  tokenData: "RORO",
  tokenizers: [0, tokens],
  topRules: {Script: [0, 23]},
  tokenPrec: 0
});
function tokenBefore(tree) {
  let cursor = tree.cursor.moveTo(tree.from, -1);
  while (/Comment/.test(cursor.name))
    cursor.moveTo(cursor.from, -1);
  return cursor.node;
}
function stripQuotes(name2) {
  let quoted = /^[`'"](.*)[`'"]$/.exec(name2);
  return quoted ? quoted[1] : name2;
}
function sourceContext(state, startPos) {
  let pos = syntaxTree(state).resolve(startPos, -1);
  let empty = false;
  if (pos.name == "Identifier" || pos.name == "QuotedIdentifier") {
    empty = false;
    let parent = null;
    let dot2 = tokenBefore(pos);
    if (dot2 && dot2.name == ".") {
      let before = tokenBefore(dot2);
      if (before && before.name == "Identifier" || before.name == "QuotedIdentifier")
        parent = stripQuotes(state.sliceDoc(before.from, before.to).toLowerCase());
    }
    return {
      parent,
      from: pos.from,
      quoted: pos.name == "QuotedIdentifier" ? state.sliceDoc(pos.from, pos.from + 1) : null
    };
  } else if (pos.name == ".") {
    let before = tokenBefore(pos);
    if (before && before.name == "Identifier" || before.name == "QuotedIdentifier")
      return {
        parent: stripQuotes(state.sliceDoc(before.from, before.to).toLowerCase()),
        from: startPos,
        quoted: null
      };
  } else {
    empty = true;
  }
  return {parent: null, from: startPos, quoted: null, empty};
}
function maybeQuoteCompletions(quote, completions) {
  if (!quote)
    return completions;
  return completions.map((c2) => Object.assign(Object.assign({}, c2), {label: quote + c2.label + quote, apply: void 0}));
}
var Span = /^\w*$/;
var QuotedSpan = /^[`'"]?\w*[`'"]?$/;
function completeFromSchema(schema, tables, defaultTable) {
  let byTable = Object.create(null);
  for (let table in schema)
    byTable[table] = schema[table].map((val) => {
      return typeof val == "string" ? {label: val, type: "property"} : val;
    });
  let topOptions = (tables || Object.keys(byTable).map((name2) => ({label: name2, type: "type"}))).concat(defaultTable && byTable[defaultTable] || []);
  return (context) => {
    let {parent, from, quoted, empty} = sourceContext(context.state, context.pos);
    if (empty && !context.explicit)
      return null;
    let options = topOptions;
    if (parent) {
      let columns = byTable[parent];
      if (!columns)
        return null;
      options = columns;
    }
    let quoteAfter = quoted && context.state.sliceDoc(context.pos, context.pos + 1) == quoted;
    return {
      from,
      to: quoteAfter ? context.pos + 1 : void 0,
      options: maybeQuoteCompletions(quoted, options),
      span: quoted ? QuotedSpan : Span
    };
  };
}
function completeKeywords(keywords11, upperCase) {
  let completions = Object.keys(keywords11).map((keyword2) => ({
    label: upperCase ? keyword2.toUpperCase() : keyword2,
    type: keywords11[keyword2] == Type ? "type" : keywords11[keyword2] == Keyword ? "keyword" : "variable",
    boost: -1
  }));
  return ifNotIn(["QuotedIdentifier", "SpecialVar", "String", "LineComment", "BlockComment", "."], completeFromList(completions));
}
var parser$1 = parser5.configure({
  props: [
    indentNodeProp.add({
      Statement: continuedIndent()
    }),
    foldNodeProp.add({
      Statement(tree) {
        return {from: tree.firstChild.to, to: tree.to};
      },
      BlockComment(tree) {
        return {from: tree.from + 2, to: tree.to - 2};
      }
    }),
    styleTags({
      Keyword: tags.keyword,
      Type: tags.typeName,
      Builtin: tags.standard(tags.name),
      Bool: tags.bool,
      Null: tags.null,
      Number: tags.number,
      String: tags.string,
      Identifier: tags.name,
      QuotedIdentifier: tags.special(tags.string),
      SpecialVar: tags.special(tags.name),
      LineComment: tags.lineComment,
      BlockComment: tags.blockComment,
      Operator: tags.operator,
      "Semi Punctuation": tags.punctuation,
      "( )": tags.paren,
      "{ }": tags.brace,
      "[ ]": tags.squareBracket
    })
  ]
});
var SQLDialect = class {
  constructor(dialect2, language2) {
    this.dialect = dialect2;
    this.language = language2;
  }
  get extension() {
    return this.language.extension;
  }
  static define(spec) {
    let d = dialect(spec, spec.keywords, spec.types, spec.builtin);
    let language2 = LezerLanguage.define({
      parser: parser$1.configure({
        tokenizers: [{from: tokens, to: tokensFor(d)}]
      }),
      languageData: {
        commentTokens: {line: "--", block: {open: "/*", close: "*/"}},
        closeBrackets: {brackets: ["(", "[", "{", "'", '"', "`"]}
      }
    });
    return new SQLDialect(d, language2);
  }
};
function keywordCompletion(dialect2, upperCase = false) {
  return dialect2.language.data.of({
    autocomplete: completeKeywords(dialect2.dialect.words, upperCase)
  });
}
function schemaCompletion(config2) {
  return config2.schema ? (config2.dialect || StandardSQL).language.data.of({
    autocomplete: completeFromSchema(config2.schema, config2.tables, config2.defaultTable)
  }) : [];
}
function sql(config2 = {}) {
  let lang = config2.dialect || StandardSQL;
  return new LanguageSupport(lang.language, [schemaCompletion(config2), keywordCompletion(lang, !!config2.upperCaseKeywords)]);
}
var StandardSQL = SQLDialect.define({});
var PostgreSQL = SQLDialect.define({
  charSetCasts: true,
  operatorChars: "+-*/<>=~!@#%^&|`?",
  specialVar: "",
  keywords: SQLKeywords + "a abort abs absent access according ada admin aggregate alias also always analyse analyze array_agg array_max_cardinality asensitive assert assignment asymmetric atomic attach attribute attributes avg backward base64 begin_frame begin_partition bernoulli bit_length blocked bom c cache called cardinality catalog_name ceil ceiling chain char_length character_length character_set_catalog character_set_name character_set_schema characteristics characters checkpoint class class_origin cluster coalesce cobol collation_catalog collation_name collation_schema collect column_name columns command_function command_function_code comment comments committed concurrently condition_number configuration conflict connection_name constant constraint_catalog constraint_name constraint_schema contains content control conversion convert copy corr cost covar_pop covar_samp csv cume_dist current_catalog current_row current_schema cursor_name database datalink datatype datetime_interval_code datetime_interval_precision db debug defaults defined definer degree delimiter delimiters dense_rank depends derived detach detail dictionary disable discard dispatch dlnewcopy dlpreviouscopy dlurlcomplete dlurlcompleteonly dlurlcompletewrite dlurlpath dlurlpathonly dlurlpathwrite dlurlscheme dlurlserver dlvalue document dump dynamic_function dynamic_function_code element elsif empty enable encoding encrypted end_frame end_partition endexec enforced enum errcode error event every exclude excluding exclusive exp explain expression extension extract family file filter final first_value flag floor following force foreach fortran forward frame_row freeze fs functions fusion g generated granted greatest groups handler header hex hierarchy hint id ignore ilike immediately immutable implementation implicit import include including increment indent index indexes info inherit inherits inline insensitive instance instantiable instead integrity intersection invoker isnull k key_member key_type label lag last_value lead leakproof least length library like_regex link listen ln load location lock locked log logged lower m mapping matched materialized max max_cardinality maxvalue member merge message message_length message_octet_length message_text min minvalue mod mode more move multiset mumps name namespace nfc nfd nfkc nfkd nil normalize normalized nothing notice notify notnull nowait nth_value ntile nullable nullif nulls number occurrences_regex octet_length octets off offset oids operator options ordering others over overlay overriding owned owner p parallel parameter_mode parameter_name parameter_ordinal_position parameter_specific_catalog parameter_specific_name parameter_specific_schema parser partition pascal passing passthrough password percent percent_rank percentile_cont percentile_disc perform period permission pg_context pg_datatype_name pg_exception_context pg_exception_detail pg_exception_hint placing plans pli policy portion position position_regex power precedes preceding prepared print_strict_params procedural procedures program publication query quote raise range rank reassign recheck recovery refresh regr_avgx regr_avgy regr_count regr_intercept regr_r2 regr_slope regr_sxx regr_sxy regr_syy reindex rename repeatable replace replica requiring reset respect restart restore result_oid returned_cardinality returned_length returned_octet_length returned_sqlstate returning reverse routine_catalog routine_name routine_schema routines row_count row_number rowtype rule scale schema_name schemas scope scope_catalog scope_name scope_schema security selective self sensitive sequence sequences serializable server server_name setof share show simple skip slice snapshot source specific_name sqlcode sqlerror sqrt stable stacked standalone statement statistics stddev_pop stddev_samp stdin stdout storage strict strip structure style subclass_origin submultiset subscription substring substring_regex succeeds sum symmetric sysid system system_time t table_name tables tablesample tablespace temp template ties token top_level_count transaction_active transactions_committed transactions_rolled_back transform transforms translate translate_regex trigger_catalog trigger_name trigger_schema trim trim_array truncate trusted type types uescape unbounded uncommitted unencrypted unlink unlisten unlogged unnamed untyped upper uri use_column use_variable user_defined_type_catalog user_defined_type_code user_defined_type_name user_defined_type_schema vacuum valid validate validator value_of var_pop var_samp varbinary variable_conflict variadic verbose version versioning views volatile warning whitespace width_bucket window within wrapper xmlagg xmlattributes xmlbinary xmlcast xmlcomment xmlconcat xmldeclaration xmldocument xmlelement xmlexists xmlforest xmliterate xmlnamespaces xmlparse xmlpi xmlquery xmlroot xmlschema xmlserialize xmltable xmltext xmlvalidate yes",
  types: SQLTypes + "bigint int8 bigserial serial8 varbit bool box bytea cidr circle precision float8 inet int4 json jsonb line lseg macaddr macaddr8 money numeric path pg_lsn point polygon float4 int2 smallserial serial2 serial serial4 text without zone with timetz timestamptz tsquery tsvector txid_snapshot uuid xml"
});
var MySQLKeywords = "accessible algorithm analyze asensitive authors auto_increment autocommit avg avg_row_length binlog btree cache catalog_name chain change changed checkpoint checksum class_origin client_statistics coalesce code collations columns comment committed completion concurrent consistent contains contributors convert database databases day_hour day_microsecond day_minute day_second delay_key_write delayed delimiter des_key_file dev_pop dev_samp deviance directory disable discard distinctrow div dual dumpfile enable enclosed ends engine engines enum errors escaped even event events every explain extended fast field fields flush force found_rows fulltext grants handler hash high_priority hosts hour_microsecond hour_minute hour_second ignore ignore_server_ids import index index_statistics infile innodb insensitive insert_method install invoker iterate keys kill linear lines list load lock logs low_priority master master_heartbeat_period master_ssl_verify_server_cert masters max max_rows maxvalue message_text middleint migrate min min_rows minute_microsecond minute_second mod mode modify mutex mysql_errno no_write_to_binlog offline offset one online optimize optionally outfile pack_keys parser partition partitions password phase plugin plugins prev processlist profile profiles purge query quick range read_write rebuild recover regexp relaylog remove rename reorganize repair repeatable replace require resume rlike row_format rtree schedule schema_name schemas second_microsecond security sensitive separator serializable server share show slave slow snapshot soname spatial sql_big_result sql_buffer_result sql_cache sql_calc_found_rows sql_no_cache sql_small_result ssl starting starts std stddev stddev_pop stddev_samp storage straight_join subclass_origin sum suspend table_name table_statistics tables tablespace terminated triggers truncate uncommitted uninstall unlock upgrade use use_frm user_resources user_statistics utc_date utc_time utc_timestamp variables views warnings xa xor year_month zerofill";
var MySQLTypes = SQLTypes + "bool blob long longblob longtext medium mediumblob mediumint mediumtext tinyblob tinyint tinytext text bigint int1 int2 int3 int4 int8 float4 float8 varbinary varcharacter precision datetime year unsigned signed";
var MySQLBuiltin = "charset clear connect edit ego exit go help nopager notee nowarning pager print prompt quit rehash source status system tee";
var MySQL = SQLDialect.define({
  operatorChars: "*+-%<>!=&|^",
  charSetCasts: true,
  doubleQuotedStrings: true,
  hashComments: true,
  spaceAfterDashes: true,
  specialVar: "@?",
  identifierQuotes: "`",
  keywords: SQLKeywords + "group_concat " + MySQLKeywords,
  types: MySQLTypes,
  builtin: MySQLBuiltin
});
var MariaSQL = SQLDialect.define({
  operatorChars: "*+-%<>!=&|^",
  charSetCasts: true,
  doubleQuotedStrings: true,
  hashComments: true,
  spaceAfterDashes: true,
  specialVar: "@?",
  identifierQuotes: "`",
  keywords: SQLKeywords + "always generated groupby_concat hard persistent shutdown soft virtual " + MySQLKeywords,
  types: MySQLTypes,
  builtin: MySQLBuiltin
});
var MSSQL = SQLDialect.define({
  keywords: SQLKeywords + "trigger proc view index for add constraint key primary foreign collate clustered nonclustered declare exec go if use index holdlock nolock nowait paglock pivot readcommitted readcommittedlock readpast readuncommitted repeatableread rowlock serializable snapshot tablock tablockx unpivot updlock with",
  types: SQLTypes + "bigint smallint smallmoney tinyint money real text nvarchar ntext varbinary image cursor hierarchyid uniqueidentifier sql_variant xml table",
  builtin: "binary_checksum checksum connectionproperty context_info current_request_id error_line error_message error_number error_procedure error_severity error_state formatmessage get_filestream_transaction_context getansinull host_id host_name isnull isnumeric min_active_rowversion newid newsequentialid rowcount_big xact_state object_id",
  operatorChars: "*+-%<>!=^&|/",
  specialVar: "@"
});
var SQLite = SQLDialect.define({
  keywords: SQLKeywords + "abort analyze attach autoincrement conflict database detach exclusive fail glob ignore index indexed instead isnull notnull offset plan pragma query raise regexp reindex rename replace temp vacuum virtual",
  types: SQLTypes + "bool blob long longblob longtext medium mediumblob mediumint mediumtext tinyblob tinyint tinytext text bigint int2 int8 year unsigned signed real",
  builtin: "auth backup bail binary changes check clone databases dbinfo dump echo eqp exit explain fullschema headers help import imposter indexes iotrace limit lint load log mode nullvalue once open output print prompt quit read restore save scanstats schema separator session shell show stats system tables testcase timeout timer trace vfsinfo vfslist vfsname width",
  operatorChars: "*+-%<>!=&|/~",
  identifierQuotes: '`"',
  specialVar: "@:?$"
});
var Cassandra = SQLDialect.define({
  keywords: "add all allow alter and any apply as asc authorize batch begin by clustering columnfamily compact consistency count create custom delete desc distinct drop each_quorum exists filtering from grant if in index insert into key keyspace keyspaces level limit local_one local_quorum modify nan norecursive nosuperuser not of on one order password permission permissions primary quorum rename revoke schema select set storage superuser table three to token truncate ttl two type unlogged update use user users using values where with writetime infinity NaN",
  types: SQLTypes + "ascii bigint blob counter frozen inet list map static text timeuuid tuple uuid varint",
  slashComments: true
});
var PLSQL = SQLDialect.define({
  keywords: SQLKeywords + "abort accept access add all alter and any array arraylen as asc assert assign at attributes audit authorization avg base_table begin between binary_integer body boolean by case cast char char_base check close cluster clusters colauth column comment commit compress connect connected constant constraint crash create current currval cursor data_base database date dba deallocate debugoff debugon decimal declare default definition delay delete desc digits dispose distinct do drop else elseif elsif enable end entry escape exception exception_init exchange exclusive exists exit external fast fetch file for force form from function generic goto grant group having identified if immediate in increment index indexes indicator initial initrans insert interface intersect into is key level library like limited local lock log logging long loop master maxextents maxtrans member minextents minus mislabel mode modify multiset new next no noaudit nocompress nologging noparallel not nowait number_base object of off offline on online only open option or order out package parallel partition pctfree pctincrease pctused pls_integer positive positiven pragma primary prior private privileges procedure public raise range raw read rebuild record ref references refresh release rename replace resource restrict return returning returns reverse revoke rollback row rowid rowlabel rownum rows run savepoint schema segment select separate session set share snapshot some space split sql start statement storage subtype successful synonym tabauth table tables tablespace task terminate then to trigger truncate type union unique unlimited unrecoverable unusable update use using validate value values variable view views when whenever where while with work",
  builtin: "appinfo arraysize autocommit autoprint autorecovery autotrace blockterminator break btitle cmdsep colsep compatibility compute concat copycommit copytypecheck define describe echo editfile embedded escape exec execute feedback flagger flush heading headsep instance linesize lno loboffset logsource long longchunksize markup native newpage numformat numwidth pagesize pause pno recsep recsepchar release repfooter repheader serveroutput shiftinout show showmode size spool sqlblanklines sqlcase sqlcode sqlcontinue sqlnumber sqlpluscompatibility sqlprefix sqlprompt sqlterminator suffix tab term termout time timing trimout trimspool ttitle underline verify version wrap",
  types: SQLTypes + "ascii bfile bfilename bigserial bit blob dec number nvarchar nvarchar2 serial smallint string text uid varchar2 xml",
  operatorChars: "*/+-%<>!=~",
  doubleQuotedStrings: true,
  charSetCasts: true
});

// node_modules/@codemirror/legacy-modes/mode/swift.js
function wordSet(words4) {
  var set = {};
  for (var i = 0; i < words4.length; i++)
    set[words4[i]] = true;
  return set;
}
var keywords10 = wordSet([
  "_",
  "var",
  "let",
  "class",
  "enum",
  "extension",
  "import",
  "protocol",
  "struct",
  "func",
  "typealias",
  "associatedtype",
  "open",
  "public",
  "internal",
  "fileprivate",
  "private",
  "deinit",
  "init",
  "new",
  "override",
  "self",
  "subscript",
  "super",
  "convenience",
  "dynamic",
  "final",
  "indirect",
  "lazy",
  "required",
  "static",
  "unowned",
  "unowned(safe)",
  "unowned(unsafe)",
  "weak",
  "as",
  "is",
  "break",
  "case",
  "continue",
  "default",
  "else",
  "fallthrough",
  "for",
  "guard",
  "if",
  "in",
  "repeat",
  "switch",
  "where",
  "while",
  "defer",
  "return",
  "inout",
  "mutating",
  "nonmutating",
  "catch",
  "do",
  "rethrows",
  "throw",
  "throws",
  "try",
  "didSet",
  "get",
  "set",
  "willSet",
  "assignment",
  "associativity",
  "infix",
  "left",
  "none",
  "operator",
  "postfix",
  "precedence",
  "precedencegroup",
  "prefix",
  "right",
  "Any",
  "AnyObject",
  "Type",
  "dynamicType",
  "Self",
  "Protocol",
  "__COLUMN__",
  "__FILE__",
  "__FUNCTION__",
  "__LINE__"
]);
var definingKeywords = wordSet(["var", "let", "class", "enum", "extension", "import", "protocol", "struct", "func", "typealias", "associatedtype", "for"]);
var atoms3 = wordSet(["true", "false", "nil", "self", "super", "_"]);
var types3 = wordSet([
  "Array",
  "Bool",
  "Character",
  "Dictionary",
  "Double",
  "Float",
  "Int",
  "Int8",
  "Int16",
  "Int32",
  "Int64",
  "Never",
  "Optional",
  "Set",
  "String",
  "UInt8",
  "UInt16",
  "UInt32",
  "UInt64",
  "Void"
]);
var operators4 = "+-/*%=|&<>~^?!";
var punc = ":;,.(){}[]";
var binary = /^\-?0b[01][01_]*/;
var octal = /^\-?0o[0-7][0-7_]*/;
var hexadecimal = /^\-?0x[\dA-Fa-f][\dA-Fa-f_]*(?:(?:\.[\dA-Fa-f][\dA-Fa-f_]*)?[Pp]\-?\d[\d_]*)?/;
var decimal = /^\-?\d[\d_]*(?:\.\d[\d_]*)?(?:[Ee]\-?\d[\d_]*)?/;
var identifier = /^\$\d+|(`?)[_A-Za-z][_A-Za-z$0-9]*\1/;
var property = /^\.(?:\$\d+|(`?)[_A-Za-z][_A-Za-z$0-9]*\1)/;
var instruction = /^\#[A-Za-z]+/;
var attribute = /^@(?:\$\d+|(`?)[_A-Za-z][_A-Za-z$0-9]*\1)/;
function tokenBase8(stream, state, prev) {
  if (stream.sol())
    state.indented = stream.indentation();
  if (stream.eatSpace())
    return null;
  var ch = stream.peek();
  if (ch == "/") {
    if (stream.match("//")) {
      stream.skipToEnd();
      return "comment";
    }
    if (stream.match("/*")) {
      state.tokenize.push(tokenComment4);
      return tokenComment4(stream, state);
    }
  }
  if (stream.match(instruction))
    return "builtin";
  if (stream.match(attribute))
    return "attribute";
  if (stream.match(binary))
    return "number";
  if (stream.match(octal))
    return "number";
  if (stream.match(hexadecimal))
    return "number";
  if (stream.match(decimal))
    return "number";
  if (stream.match(property))
    return "property";
  if (operators4.indexOf(ch) > -1) {
    stream.next();
    return "operator";
  }
  if (punc.indexOf(ch) > -1) {
    stream.next();
    stream.match("..");
    return "punctuation";
  }
  var stringMatch;
  if (stringMatch = stream.match(/("""|"|')/)) {
    var tokenize2 = tokenString4.bind(null, stringMatch[0]);
    state.tokenize.push(tokenize2);
    return tokenize2(stream, state);
  }
  if (stream.match(identifier)) {
    var ident = stream.current();
    if (types3.hasOwnProperty(ident))
      return "type";
    if (atoms3.hasOwnProperty(ident))
      return "atom";
    if (keywords10.hasOwnProperty(ident)) {
      if (definingKeywords.hasOwnProperty(ident))
        state.prev = "define";
      return "keyword";
    }
    if (prev == "define")
      return "def";
    return "variable";
  }
  stream.next();
  return null;
}
function tokenUntilClosingParen() {
  var depth = 0;
  return function(stream, state, prev) {
    var inner = tokenBase8(stream, state, prev);
    if (inner == "punctuation") {
      if (stream.current() == "(")
        ++depth;
      else if (stream.current() == ")") {
        if (depth == 0) {
          stream.backUp(1);
          state.tokenize.pop();
          return state.tokenize[state.tokenize.length - 1](stream, state);
        } else
          --depth;
      }
    }
    return inner;
  };
}
function tokenString4(openQuote, stream, state) {
  var singleLine = openQuote.length == 1;
  var ch, escaped = false;
  while (ch = stream.peek()) {
    if (escaped) {
      stream.next();
      if (ch == "(") {
        state.tokenize.push(tokenUntilClosingParen());
        return "string";
      }
      escaped = false;
    } else if (stream.match(openQuote)) {
      state.tokenize.pop();
      return "string";
    } else {
      stream.next();
      escaped = ch == "\\";
    }
  }
  if (singleLine) {
    state.tokenize.pop();
  }
  return "string";
}
function tokenComment4(stream, state) {
  var ch;
  while (true) {
    stream.match(/^[^/*]+/, true);
    ch = stream.next();
    if (!ch)
      break;
    if (ch === "/" && stream.eat("*")) {
      state.tokenize.push(tokenComment4);
    } else if (ch === "*" && stream.eat("/")) {
      state.tokenize.pop();
    }
  }
  return "comment";
}
function Context3(prev, align, indented) {
  this.prev = prev;
  this.align = align;
  this.indented = indented;
}
function pushContext3(state, stream) {
  var align = stream.match(/^\s*($|\/[\/\*])/, false) ? null : stream.column() + 1;
  state.context = new Context3(state.context, align, state.indented);
}
function popContext3(state) {
  if (state.context) {
    state.indented = state.context.indented;
    state.context = state.context.prev;
  }
}
var swift = {
  startState: function() {
    return {
      prev: null,
      context: null,
      indented: 0,
      tokenize: []
    };
  },
  token: function(stream, state) {
    var prev = state.prev;
    state.prev = null;
    var tokenize2 = state.tokenize[state.tokenize.length - 1] || tokenBase8;
    var style = tokenize2(stream, state, prev);
    if (!style || style == "comment")
      state.prev = prev;
    else if (!state.prev)
      state.prev = style;
    if (style == "punctuation") {
      var bracket2 = /[\(\[\{]|([\]\)\}])/.exec(stream.current());
      if (bracket2)
        (bracket2[1] ? popContext3 : pushContext3)(state, stream);
    }
    return style;
  },
  indent: function(state, textAfter, iCx) {
    var cx = state.context;
    if (!cx)
      return 0;
    var closing3 = /^[\]\}\)]/.test(textAfter);
    if (cx.align != null)
      return cx.align - (closing3 ? 1 : 0);
    return cx.indented + (closing3 ? 0 : iCx.unit);
  },
  languageData: {
    indentOnInput: /^\s*[\)\}\]]$/,
    commentTokens: {line: "//", block: {open: "/*", close: "*/"}},
    closeBrackets: {brackets: ["(", "[", "{", "'", '"', "`"]}
  }
};

// editor.js
var extensions = [
  EditorView.lineWrapping,
  bracketMatching(),
  closeBrackets(),
  defaultHighlightStyle,
  indentOnInput(),
  keymap.of([...closeBracketsKeymap, defaultTabBinding, ...standardKeymap]),
  lineNumbers()
];
var languages = {
  bash: StreamLanguage.define(shell),
  brainfuck: StreamLanguage.define(brainfuck),
  c: StreamLanguage.define(c),
  "c-sharp": StreamLanguage.define(csharp),
  cobol: StreamLanguage.define(cobol),
  crystal: StreamLanguage.define(crystal),
  "f-sharp": StreamLanguage.define(fSharp),
  fortran: StreamLanguage.define(fortran),
  go: StreamLanguage.define(go),
  haskell: StreamLanguage.define(haskell),
  java: java2(),
  javascript: javascript(),
  julia: StreamLanguage.define(julia),
  lisp: StreamLanguage.define(commonLisp),
  lua: StreamLanguage.define(lua),
  perl: StreamLanguage.define(perl),
  powershell: StreamLanguage.define(powerShell),
  python: python(),
  ruby: StreamLanguage.define(ruby),
  rust: rust(),
  sql: sql({dialect: SQLite}),
  swift: StreamLanguage.define(swift)
};
export {
  EditorState,
  EditorView,
  extensions,
  languages
};
