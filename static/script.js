let matches;

if (matches = /\/(javascript|perl6?|php|python|ruby)$/.exec(location.pathname))
    onload = function() {
        CodeMirror.fromTextArea(
            document.getElementsByTagName('textarea')[0],
            { lineNumbers: true, mode: matches[1] }
        );
    };
