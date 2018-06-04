package routes

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func about(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	printHeader(w, r, 200)

	const html = "<main id=about><p>" +
		"Code Golf is a game designed to let you show off your code-fu " +
		"by solving problems in the least number of characters. " +
		"It is written in <a href=//golang.org>Go</a> and is " +
		"<a href=//github.com/JRaspass/code-golf>open source</a>, " +
		"patches welcome!</p>" +
		"<h2>Frequently Asked Questions</h2><dl>" +
		"<dt>Do I Need to Login to Play?" +
		"<dd>No. Submitted solutions will be executed and checked without " +
		"logging in, but nothing will be saved and you won't appear on the " +
		"<a href=scores>leaderboards</a>." +
		"<dt>Which Languages Are Supported?" +
		"<dd>Currently <a href=//www.gnu.org/software/bash/>Bash</a>, " +
		"<a href=//www.haskell.org/ghc/>Haskell</a>, " +
		"<a href=//developer.mozilla.org/en-US/docs/Mozilla/Projects/SpiderMonkey/Introduction_to_the_JavaScript_shell>JavaScript</a>, " +
		"<a href=//clisp.sourceforge.io>Lisp</a>, " +
		"<a href=//www.lua.org/>Lua</a>, " +
		"<a href=//www.perl.org>Perl</a>, " +
		"<a href=//perl6.org>Perl 6</a>, " +
		"<a href=http://www.php.net>PHP</a>, " +
		"<a href=//www.python.org>Python</a>, " +
		"and <a href=//www.ruby-lang.org>Ruby</a>. " +
		"If you'd like to see another language added then raise an " +
		"<a href=//github.com/JRaspass/code-golf/issues/new>issue</a>." +
		"<dt>What Version Is Each Language At?" +
		"<dd>Code Golf ties to track upstream closely, manually amending " +
		"solutions if they cease to work on newer versions, the current " +
		"version matrix is as follows:<br><br><table>" + versionTable + "</table>" +
		"<dt>Is Perl customized in any way?" +
		"<dd>Yes. Code Golf uses miniperl with no ability to load modules, " +
		`but modern features like "say" and "state" enabled by default. ` +
		"For full details read the <a href=//github.com/JRaspass/code-golf/blob/master/containers/perl/Dockerfile>dockerfile</a>." +
		"<dt>Are Warnings Ignored?" +
		"<dd>Yes. Only STDOUT is checked against the solution, STDERR is " +
		"however shown back to you to ease debugging." +
		"<dt>How Are Arguments Passed to My Program?" +
		"<dd>Some holes pass arguments, for those your program should read " +
		"them from ARGV." +
		"<dt>How Do I Access ARGV in JavaScript?" +
		"<dd>Code Golf uses Mozilla's <a href=//developer.mozilla.org/en-US/docs/Mozilla/Projects/SpiderMonkey/Introduction_to_the_JavaScript_shell>JavaScript Shell</a>, " +
		`they extend the language to have a global variable "scriptArgs" ` +
		"which will contain ARGV. For more details read the <a href=//developer.mozilla.org/en-US/docs/Mozilla/Projects/SpiderMonkey/Shell_global_objects>full docs</a>." +
		"<dt>How Are Solutions Scored?" +
		"<dd>The score of your solution is the sum of the UTF-8 characters " +
		"in your source code. This means both " +
		`"A" (U+0041 Latin Capital Letter A) and "😉" (U+1F609 Winking Face) ` +
		"cost the same despite the 1:4 ratio in byte count." +
		"<dt>How Is My Overall Score Computed?" +
		"<dd>For each hole, the shortest solution is awarded 100 points, " +
		"with the points descreasing in uniform decrements per rank. " +
		"Your overall score is simply the sum of your points in each hole." +
		"<dt>Are Submissions Resource Constrained?" +
		"<dd>Yes. Execution time is limited to 5 seconds, CPU & RAM usage " +
		"is unbounded but this will probably change soon." +
		"<dt>Can I See Other Peoples Solutions?" +
		"<dd>Currently no, but a feature is in development to allow you to " +
		"see any solution that you already beaten in score." +
		"</dt>"

	w.Write([]byte(html))
}
