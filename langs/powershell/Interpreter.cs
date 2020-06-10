using System;
using System.Collections.Generic;
using System.Collections.ObjectModel;
using System.Globalization;
using System.Management.Automation;
using System.Management.Automation.Host;
using System.Management.Automation.Runspaces;
using System.Security;

class Program
{
	static int Main(string[] args)
	{
		if (args.Length > 0 && args[0] == "--version")
		{
			var framework = System.Runtime.InteropServices.RuntimeInformation.FrameworkDescription;
			Console.WriteLine($"PowerShell {PSVersionInfo.PSVersion} on {framework}");
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

		// These aliases are from:
		// https://github.com/PowerShell/PowerShell/blob/master/src/System.Management.Automation/engine/InitialSessionState.cs
		// They are disabled by default on Unix-like systems to avoid conflicting with built-in commands, but
		// there are no conflicts within the code golf environment.
		state.Commands.Add(new SessionStateCommandEntry[]
		{
			new SessionStateAliasEntry("ac", "Add-Content", string.Empty, ScopedItemOptions.ReadOnly),
			new SessionStateAliasEntry("clear", "Clear-Host"),
			new SessionStateAliasEntry("compare", "Compare-Object", string.Empty, ScopedItemOptions.ReadOnly),
			new SessionStateAliasEntry("cpp", "Copy-ItemProperty", string.Empty, ScopedItemOptions.ReadOnly),
			new SessionStateAliasEntry("diff", "Compare-Object", string.Empty, ScopedItemOptions.ReadOnly),
			new SessionStateAliasEntry("gsv", "Get-Service", string.Empty, ScopedItemOptions.ReadOnly),
			new SessionStateAliasEntry("sleep", "Start-Sleep", string.Empty, ScopedItemOptions.ReadOnly),
			new SessionStateAliasEntry("sort", "Sort-Object", string.Empty, ScopedItemOptions.ReadOnly),
			new SessionStateAliasEntry("start", "Start-Process", string.Empty, ScopedItemOptions.ReadOnly),
			new SessionStateAliasEntry("sasv", "Start-Service", string.Empty, ScopedItemOptions.ReadOnly),
			new SessionStateAliasEntry("spsv", "Stop-Service", string.Empty, ScopedItemOptions.ReadOnly),
			new SessionStateAliasEntry("tee", "Tee-Object", string.Empty, ScopedItemOptions.ReadOnly),
			new SessionStateAliasEntry("write", "Write-Output", string.Empty, ScopedItemOptions.ReadOnly),
			// These were transferred from the "transferred from the profile" section
			new SessionStateAliasEntry("cat", "Get-Content"),
			new SessionStateAliasEntry("cp", "Copy-Item", string.Empty, ScopedItemOptions.AllScope),
			new SessionStateAliasEntry("ls", "Get-ChildItem"),
			new SessionStateAliasEntry("man", "help"),
			new SessionStateAliasEntry("mount", "New-PSDrive"),
			new SessionStateAliasEntry("mv", "Move-Item"),
			new SessionStateAliasEntry("ps", "Get-Process"),
			new SessionStateAliasEntry("rm", "Remove-Item"),
			new SessionStateAliasEntry("rmdir", "Remove-Item"),
			new SessionStateAliasEntry("cnsn", "Connect-PSSession", string.Empty, ScopedItemOptions.ReadOnly),
			new SessionStateAliasEntry("dnsn", "Disconnect-PSSession", string.Empty, ScopedItemOptions.ReadOnly),
			new SessionStateAliasEntry("ogv", "Out-GridView", string.Empty, ScopedItemOptions.ReadOnly),
			new SessionStateAliasEntry("shcm", "Show-Command", string.Empty, ScopedItemOptions.ReadOnly),
		});

		var host = new Host();
		using var runspace = RunspaceFactory.CreateRunspace(host, state);
		using var powerShell = PowerShell.Create();
		runspace.Open();
		powerShell.Runspace = runspace;
		try
		{
			powerShell.AddScript(code);
			foreach (var arg in args)
			{
				powerShell.AddArgument(arg);
			}

			if (!requireExplicitOutput)
			{
				powerShell.Commands.AddCommand("Out-Default");
			}

			powerShell.Invoke();

			foreach (var item in powerShell.Streams.Error)
			{
				Console.Error.WriteLine(item);
			}

			return host.ExitCode;
		}
		catch (Exception e)
		{
			Console.Error.WriteLine($"{e.GetType().FullName}: {e.Message}");
			return 1;
		}
	}
}

class Host : PSHost
{
	public int ExitCode { get; private set; }

	public override CultureInfo CurrentCulture => CultureInfo.CurrentCulture;

	public override CultureInfo CurrentUICulture => CultureInfo.CurrentUICulture;

	public override Guid InstanceId { get; } = Guid.NewGuid();

	public override string Name => "CodeGolfHost";

	public override PSHostUserInterface UI { get; } = new UserInterface();

	public override Version Version => PSVersionInfo.PSVersion;

	public override void EnterNestedPrompt()
	{
	}

	public override void ExitNestedPrompt()
	{
	}

	public override void NotifyBeginApplication()
	{
	}

	public override void NotifyEndApplication()
	{
	}

	public override void SetShouldExit(int exitCode)
	{
		ExitCode = exitCode;
	}
}

class UserInterface : PSHostUserInterface
{
	public override PSHostRawUserInterface RawUI { get; } = new RawUserInterface();

	public override Dictionary<string, PSObject> Prompt(string caption, string message, Collection<FieldDescription> descriptions)
	{
		throw new NotImplementedException();
	}

	public override int PromptForChoice(string caption, string message, Collection<ChoiceDescription> choices, int defaultChoice)
	{
		throw new NotImplementedException();
	}

	public override PSCredential PromptForCredential(string caption, string message, string userName, string targetName)
	{
		throw new NotImplementedException();
	}

	public override PSCredential PromptForCredential(string caption, string message, string userName, string targetName, PSCredentialTypes allowedCredentialTypes, PSCredentialUIOptions options)
	{
		throw new NotImplementedException();
	}

	public override string ReadLine()
	{
		throw new NotImplementedException();
	}

	public override SecureString ReadLineAsSecureString()
	{
		throw new NotImplementedException();
	}

	public override void Write(ConsoleColor foregroundColor, ConsoleColor backgroundColor, string value)
	{
		Console.Write(value);
	}

	public override void Write(string value)
	{
		Console.Write(value);
	}

	public override void WriteDebugLine(string message)
	{
		// pwsh.exe doesn't show this by default.
	}

	public override void WriteErrorLine(string value)
	{
		Console.Error.WriteLine(value);
	}

	public override void WriteLine(string value)
	{
		Console.WriteLine(value);
	}

	public override void WriteProgress(long sourceId, ProgressRecord record)
	{
	}

	public override void WriteVerboseLine(string message)
	{
		// pwsh.exe doesn't show this by default.
	}

	public override void WriteWarningLine(string message)
	{
		// Match behavior of pwsh.exe
		Console.WriteLine($"WARNING: {message}");
	}
}

class RawUserInterface : PSHostRawUserInterface
{
	public override void FlushInputBuffer()
	{
	}

	public override BufferCell[,] GetBufferContents(Rectangle rectangle)
	{
		throw new NotImplementedException();
	}

	public override KeyInfo ReadKey(ReadKeyOptions options)
	{
		throw new NotImplementedException();
	}

	public override void ScrollBufferContents(Rectangle source, Coordinates destination, Rectangle clip, BufferCell fill)
	{
		throw new NotImplementedException();
	}

	public override void SetBufferContents(Coordinates origin, BufferCell[,] contents)
	{
		throw new NotImplementedException();
	}

	public override void SetBufferContents(Rectangle rectangle, BufferCell fill)
	{
		throw new NotImplementedException();
	}

	public override ConsoleColor BackgroundColor { get; set; }

	public override Size BufferSize { get; set; } = new Size(1000, 1000);

	public override Coordinates CursorPosition { get; set; }

	public override int CursorSize { get; set; }

	public override ConsoleColor ForegroundColor { get; set; }

	public override bool KeyAvailable => false;

	public override Size MaxPhysicalWindowSize => new Size(1000, 1000);

	public override Size MaxWindowSize => new Size(1000, 1000);

	public override Coordinates WindowPosition { get; set; }

	public override Size WindowSize { get; set; } = new Size(1000, 1000);

	public override string WindowTitle { get; set; } = "";
}
