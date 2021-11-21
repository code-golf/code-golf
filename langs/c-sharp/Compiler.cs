using System.Reflection;
using System.Runtime.InteropServices;
using Microsoft.CodeAnalysis;
using Microsoft.CodeAnalysis.CSharp;

if (args.Length > 0 && args[0] == "--version")
{
	Console.WriteLine("C# {0} on .NET {1}",
		LanguageVersion.Latest.MapSpecifiedToEffectiveVersion().ToDisplayString(),
		Environment.Version);
	return 0;
}

var code = Console.In.ReadToEnd();

var options = new CSharpParseOptions(LanguageVersion.Latest);
var syntaxTree = CSharpSyntaxTree.ParseText(code, options);

// Explicitly enable default C# 10 using directives for console applications.
const string GlobalUsings = @"
global using global::System;
global using global::System.Collections.Generic;
global using global::System.IO;
global using global::System.Linq;
global using global::System.Net.Http;
global using global::System.Threading;
global using global::System.Threading.Tasks;";
var usingsSyntaxTree = CSharpSyntaxTree.ParseText(GlobalUsings, options);

var compilationOptions = new CSharpCompilationOptions(OutputKind.ConsoleApplication,
	optimizationLevel: OptimizationLevel.Release,
	allowUnsafe: true);

var references = new MetadataReference[]
{
	MetadataReference.CreateFromFile(Assembly.Load("netstandard").Location),
	MetadataReference.CreateFromFile(Assembly.Load("System.Collections").Location),
	MetadataReference.CreateFromFile(Assembly.Load("System.Runtime").Location),
	MetadataReference.CreateFromFile(typeof(object).Assembly.Location),
	MetadataReference.CreateFromFile(typeof(Console).Assembly.Location),
	MetadataReference.CreateFromFile(typeof(Microsoft.CSharp.RuntimeBinder.Binder).Assembly.Location),
	MetadataReference.CreateFromFile(typeof(HashSet<>).Assembly.Location),
	MetadataReference.CreateFromFile(typeof(List<>).Assembly.Location),
	MetadataReference.CreateFromFile(typeof(System.Dynamic.ExpandoObject).Assembly.Location),
	MetadataReference.CreateFromFile(typeof(Enumerable).Assembly.Location),
	MetadataReference.CreateFromFile(typeof(System.Numerics.BigInteger).Assembly.Location),
	MetadataReference.CreateFromFile(typeof(System.Text.RegularExpressions.Regex).Assembly.Location),
};

var compilation = CSharpCompilation.Create(
    "code", new[] { syntaxTree, usingsSyntaxTree }, references, compilationOptions);

var exePath = "/tmp/code.exe";

var result = compilation.Emit(exePath);
if (!result.Success)
{
	foreach (var d in result.Diagnostics)
	{
		Console.Error.WriteLine(d);
	}

	return 1;
}

return AppDomain.CurrentDomain.ExecuteAssembly(exePath, args[1..]);
