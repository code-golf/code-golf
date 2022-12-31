open System
open System.IO
open System.Reflection
open FSharp.Compiler.CodeAnalysis

let runCode (assemblyPath: string) (args: string[]) =
    let assembly = Assembly.LoadFile assemblyPath
    let entryPoint =
        assembly.GetTypes()
        |> Seq.map (fun t -> t.GetMethods())
        |> Seq.concat
        |> Seq.tryFind (fun m ->
            m.GetCustomAttributes(typeof<EntryPointAttribute>, true).Length > 0)

    match entryPoint with
    | Some method -> method.Invoke(null, [| args |]) |> unbox<int>
    | None ->
        let main =
            assembly.GetTypes()
            |> Seq.pick (fun t ->
                match t.GetMethod("main@", BindingFlags.Public ||| BindingFlags.Static) with
                | null -> None
                | m -> Some m)
        main.Invoke(null, [| |]) |> ignore
        0

let compile(checker: FSharpChecker) (codeFile: string) (assemblyPath: string) =
    let compileArgs =
        [|
            "fsi.exe";
            "-o"; assemblyPath;
            codeFile;
            "--targetprofile:netcore";
            "--target:exe";
            "--optimize";
            "-r:/usr/bin/FSharp.Core.dll";
            "-r:/usr/bin/System.Private.CoreLib.dll";
            "-r:/usr/bin/System.Runtime.dll";
            "-r:/usr/bin/netstandard.dll";
            "--nowin32manifest";
        |]
    let info, result =
        checker.Compile(compileArgs)
        |> Async.RunSynchronously
    info |> Seq.iter (fprintfn stderr "%O")
    if result <> 0 then
        Environment.Exit result

[<EntryPoint>]
let main args =
    if args.Length < 1 then
        fprintfn stderr "Arguments required."
        1
    elif args.[0] = "--version" then
        let assembly = Assembly.Load "FSharp.Compiler.Service"
        let versionType = assembly.GetType "FSharp.Compiler.Features+LanguageVersion"
        let field = versionType.GetField("latestVersion", BindingFlags.Static ||| BindingFlags.NonPublic)
        let version = field.GetValue null :?> decimal
        let framework = System.Environment.Version
        printfn "F# %s on .NET %O" (version.ToString "g") framework
        0
    else
        let codeFile = "/tmp/Code.fs"
        let assemblyPath = "/tmp/Code.exe"
        File.WriteAllText(codeFile, stdin.ReadToEnd())
        let checker = FSharpChecker.Create()
        compile checker codeFile assemblyPath
        File.Delete codeFile
        runCode assemblyPath args.[1..]
