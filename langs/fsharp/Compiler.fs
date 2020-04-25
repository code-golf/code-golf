open System
open System.IO
open System.Reflection
open FSharp.Compiler.SourceCodeServices

let runCode (assembly: Assembly) (args: string[]) =
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

let compileArgs =
    [|
        "-o"; "/tmp/Code.exe";
        "/tmp/Code.fs";
        "--targetprofile:netcore";
        "--target:exe";
        "--optimize";
        "-r:/compiler/FSharp.Core.dll";
        "-r:/compiler/System.Private.CoreLib.dll";
        "-r:/compiler/System.Runtime.dll";
        "-r:/compiler/netstandard.dll";
        "--nowin32manifest";
    |]

let compile(checker: FSharpChecker) =
    let info, result, assembly =
        checker.CompileToDynamicAssembly(compileArgs, None)
        |> Async.RunSynchronously
    info |> Seq.iter (fprintfn stderr "%O")
    if result <> 0 then
        Environment.Exit result
    assembly.Value

[<EntryPoint>]
let main args =
    if args.Length < 1 then
        fprintfn stderr "Arguments required."
        Environment.Exit 0

    if args.[0] = "--version" then
        let info = File.ReadAllText "/compiler/version.txt"
        let framework = System.Runtime.InteropServices.RuntimeInformation.FrameworkDescription
        printfn "%s on %s" (info.Trim()) framework
        Environment.Exit 0

    let codeFile = "/tmp/Code.fs"
    File.WriteAllText(codeFile, stdin.ReadToEnd())
    let checker = FSharpChecker.Create()
    let assembly = compile checker
    File.Delete codeFile
    runCode assembly args.[1..]
