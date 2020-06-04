using System;
using System.IO;
using System.Linq;
using System.Management.Automation;
using System.Management.Automation.Runspaces;

class Program
{
	static int Main(string[] args)
	{
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

		var requireExplicitOutput = args.Length > 0 && args[0] == "--explicit";
		if (requireExplicitOutput)
		{
			args = args[1..];
		}

		return Run(Console.In.ReadToEnd(), args, requireExplicitOutput);
	}

	static int Run(string code, string[] args, bool requireExplicitOutput)
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

				var results = powerShell.Invoke();
				if (!requireExplicitOutput)
				{
					foreach (var item in results)
					{
						Console.WriteLine(item);
					}
				}

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
