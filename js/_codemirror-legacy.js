import CodeMirror from './codemirror-legacy/_codemirror';

export default CodeMirror;

// Addons.
import './codemirror-legacy/_vim';

// Modes.
import './codemirror-legacy/_bash';
import './codemirror-legacy/_brainfuck';
import './codemirror-legacy/_cobol';
import './codemirror-legacy/_crystal';
import './codemirror-legacy/_d';
import './codemirror-legacy/_fortran';
import './codemirror-legacy/_go';
import './codemirror-legacy/_haskell';
import './codemirror-legacy/_javascript';
import './codemirror-legacy/_julia';
import './codemirror-legacy/_lisp';
import './codemirror-legacy/_lua';
import './codemirror-legacy/_mllike';
import './codemirror-legacy/_pascal';
import './codemirror-legacy/_perl';
import './codemirror-legacy/_php';
import './codemirror-legacy/_powershell';
import './codemirror-legacy/_python';
import './codemirror-legacy/_raku';
import './codemirror-legacy/_ruby';
import './codemirror-legacy/_rust';
import './codemirror-legacy/_sql';
import './codemirror-legacy/_swift';

// Nim comes from NPM.
import { nim } from 'nim-codemirror-mode';
CodeMirror.defineMode('nim', nim);
CodeMirror.defineMIME('text/x-nim', 'nim');
