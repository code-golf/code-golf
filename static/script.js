let matches;

if (matches = /\/(perl)$/.exec(location.pathname))
    onload = function() {
        CodeMirror.fromTextArea(
            document.getElementsByTagName('textarea')[0],
            { lineNumbers: true, mode: matches[1] }
        );
    };
