using System;
using System.IO;
using System.Linq;
using System.Management.Automation;
using System.Management.Automation.Runspaces;

class Program
{
	static int Main(string[] args)
	{
		if (args.Length == 0)
		{
			Console.Error.WriteLine("Arguments required.");
			return 1;
		}
		if (args.Length > 0 && args[0] == "--version")
		{
			string version;
			using (var powerShell = PowerShell.Create())
			{
				version = powerShell.AddScript("(Get-Host).Version.ToString()").Invoke()[0].ToString();
			}
			var framework = System.Runtime.InteropServices.RuntimeInformation.FrameworkDescription;
			Console.WriteLine($"PowerShell {version} on {framework}");
			return 0;
		}

		var code = args[0] == "-" ? Console.In.ReadToEnd() : File.ReadAllText(args[0]);
		return Run(code, args[1..]);
	}

	static int Run(string code, string[] args)
	{
		var state = InitialSessionState.CreateDefault();
		using (var runspace = RunspaceFactory.CreateRunspace(state))
		using (var powerShell = PowerShell.Create())
		{
			runspace.Open();
			powerShell.Runspace = runspace;
			try
			{
				powerShell.AddScript(code);
				foreach (var arg in args)
				{
					powerShell.AddArgument(arg);
				}

				// Ignore the results of invoking the script. This requires explicit output to be
				// used and prevents a two character Quine solution.
				powerShell.Invoke();

				foreach (var item in powerShell.Streams.Information)
				{
					Console.WriteLine(item);
				}
				foreach (var item in powerShell.Streams.Error)
				{
					Console.Error.WriteLine(item);
				}

				return 0;
			}
			catch (Exception e)
			{
				Console.Error.WriteLine($"{e.GetType().FullName}: {e.Message}");
				return 1;
			}
		}
	}
}
