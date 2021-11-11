using System;
using System.Reflection;
using Microsoft.CodeAnalysis;
using Microsoft.CodeAnalysis.CSharp;

namespace Compiler
{
	class Program
	{
		// Explicitly enable default C# 10 using directives for console applications.
		const string GlobalUsings = @"
			global using global::System;
			global using global::System.Collections.Generic;
			global using global::System.IO;
			global using global::System.Linq;
			global using global::System.Net.Http;
			global using global::System.Threading;
			global using global::System.Threading.Tasks;";

		static int Main(string[] args)
		{
			if (args.Length > 0 && args[0] == "--version")
			{
				var version = LanguageVersion.Latest.MapSpecifiedToEffectiveVersion().ToDisplayString();
				var framework = System.Runtime.InteropServices.RuntimeInformation.FrameworkDescription;
				Console.WriteLine($"C# {version} on {framework}");
				return 0;
			}

			var code = Console.In.ReadToEnd();

			var options = new CSharpParseOptions(LanguageVersion.Latest);
			var syntaxTree = CSharpSyntaxTree.ParseText(code, options);
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
				MetadataReference.CreateFromFile(typeof(System.Collections.Generic.HashSet<>).Assembly.Location),
				MetadataReference.CreateFromFile(typeof(System.Collections.Generic.List<>).Assembly.Location),
				MetadataReference.CreateFromFile(typeof(System.Dynamic.ExpandoObject).Assembly.Location),
				MetadataReference.CreateFromFile(typeof(System.Linq.Enumerable).Assembly.Location),
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
		}
	}
}
