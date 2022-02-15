import std.process;
import std.stdio;

void main(string[] args) {

	args[0] = "/usr/bin/ldc2";

    if (args.length > 0 && args[1] != "--version") {

        auto f = File("/tmp/code.d", "w");

        foreach(ubyte[] buffer; stdin.byChunk(4096)) {

            f.rawWrite(buffer);

        }

        f.close();
    }

    environment["PATH"] = "/usr/bin/";
    execv("/usr/bin/ldc2", args);

    stderr.writeln("execv");

}