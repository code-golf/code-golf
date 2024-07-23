#!/usr/bin/env nim
######################################################
# Arturo
# Programming Language + Bytecode VM compiler
# Copyright (c) 2019-2024 Yanis Zafirópulos (aka Dr.Kameleon)
#  
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
# 
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
# 
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.
#
# @file: build.nims
######################################################
# initial conversion to NimScript thanks to:
# - Patrick (skydive241@gmx.de)

#=======================================
# Libraries
#=======================================

import std/json, os
import strformat, strutils

import ".config/utils/ui.nims"
import ".config/utils/cli.nims"

#=======================================
# Initialize globals
#=======================================

mode = ScriptMode.Silent
--hints:off

#=======================================
# Flag system
#=======================================

include ".config/utils/flags.nims"

include ".config/arch.nims"
include ".config/buildmode.nims"
include ".config/devtools.nims"
include ".config/who.nims"

#=======================================
# Constants
#=======================================

let
    targetDir = getHomeDir()/".arturo"

    paths: tuple = (
        targetBin:      targetDir/"bin",
        targetLib:      targetDir/"lib",
        targetStores:   targetDir/"stores",
        mainFile:       "src"/"arturo.nim",
    )

#=======================================
# Types
#=======================================

type BuildConfig = tuple
    binary, version, bundle: string
    shouldCompress, shouldInstall, shouldLog, generateBundle, isDeveloper: bool

func webVersion(config: BuildConfig): bool

func backend(config: BuildConfig): string =
    result = "c"
    if config.webVersion:
        return "js"

func silentCompilation(config: BuildConfig): bool =
    ## CI and User builds should actually be silent,
    ## the most important is the exit code.
    ## But for developers, it's useful to have a detailed log.
    not (config.isDeveloper or config.shouldLog)

func webVersion(config: BuildConfig): bool =
    config.version == "@web"

func buildConfig(): BuildConfig =
    (
        binary:             "bin/arturo".toExe,
        version:            "@full",
        bundle:             "",
        shouldCompress:     true,
        shouldInstall:      false,
        shouldLog:          false,
        generateBundle:     false,
        isDeveloper:        false,
    )

#=======================================
# Helpers
#=======================================

func toErrorCode(a: bool): int =
    if a:
        return QuitSuccess
    else:
        return QuitFailure

template unless(condition: bool, body: untyped) =
    if not condition:
        body

# TODO(build.nims) JavaScript compression not working correctly
#  labels: web,bug
proc recompressJS*(jsFile: string, config: BuildConfig) =
    var js: string
    "testsed.txt".writeFile("""
        s/Field([0-5])/F\1/g
        s/field [^\"]+ is not accessible [^\"]+//g
    """)

    let CompressionResult =
        gorgeEx fmt"""
            sed -E -f testsed.txt {jsFile}
        """

    if CompressionResult.exitCode != QuitSuccess:
        js = readFile(jsFile)
            .replaceWord("Field0", "F0")
            .replaceWord("Field1", "F1")
            .replaceWord("Field2", "F2")
            .replaceWord("Field3", "F3")
    else:
        js = CompressionResult.output

    jsFile.writeFile js

proc miniBuild*() =
    # all the necessary "modes" for mini builds
    miniBuildConfig()

    # plus, shrinking + the MINI flag
    if hostOS=="freebsd" or hostOS=="openbsd" or hostOS=="netbsd":
        --verbosity:3

proc compressBinary(config: BuildConfig) =

    if (not config.shouldCompress) or (not config.webVersion):
        return

    section "Post-processing..."

    log "compressing binary..."
    let minBin = config.binary.replace(".js",".min.js")

    let CompressionResult =
        gorgeEx fmt"uglifyjs {config.binary} -c -m ""toplevel,reserved=['A$']"" -c -o {minBin}"

    if CompressionResult.exitCode != QuitSuccess:
        warn "uglifyjs: 3rd-party tool not available"
        minBin.writeFile readFile(config.binary)
    
    recompressJS(minBin, config)

proc verifyDirectories*() =
    ## Create target dirs recursively, if they don't exist
    log "setting up directories..."
    for path in [paths.targetBin, paths.targetLib, paths.targetStores]:
        mkdir path

proc updateBuild*() =
    ## Increment the build version by one and perform a commit.
    
    proc commit(file: string): string =
        let cmd = fmt"git commit -m 'build update' {file}"
        cmd.gorgeEx().output

    proc increaseVersion(file: string) =
        let buildVersion: int = file.readFile()
                                    .strip()
                                    .parseInt()
                                    .succ()

        file.writeFile $buildVersion
    
    proc main() =
        let buildFile = "version/build"
        increaseVersion(buildFile)
        for line in commit(buildFile).splitLines:
            echo line.strip()

    main()

proc compile*(config: BuildConfig, showFooter: bool = false): int
    {. raises: [OSError, ValueError, Exception] .} =

    proc windowsHostSpecific() =
        if config.isDeveloper and not flags.contains("NOWEBVIEW"):
            discard gorgeEx "src\\extras\\webview\\deps\\build.bat"
            #discard gorgeEx "src\\extras\\webview\\deps\\build-new.bat"
        --passL:"\"-static-libstdc++ -static-libgcc -Wl,-Bstatic -lstdc++ -Wl,-Bdynamic\""
        --gcc.linkerexe:"g++"

    proc unixHostSpecific() =
        --passL:"\"-lm\""

    result = QuitSuccess
    let
        params = flags.join(" ")
        cmd = fmt"nim {config.backend} {params} -o:{config.binary} {paths.mainFile}"

    if "windows" == hostOS:
         windowsHostSpecific()
    else:
        unixHostSpecific()

    if config.silentCompilation:
        return cmd.gorgeEx().exitCode
    else:
        echo fmt"{colors.gray}"
        cmd.exec()

proc installAll*(config: BuildConfig, targetFile: string) =

    # Helper functions

    proc copy(file: string, source: string, target: string) =
        cpFile source.joinPath(file), target.joinPath(file)

    # Methods

    proc copyWebView() =
        let 
            sourcePath = "src\\extras\\webview\\deps\\dlls\\x64\\"
            targetPath = "bin"
        log "copying webview..."
        "webview.dll".copy(sourcePath, targetPath)
        "WebView2Loader.dll".copy(sourcePath, targetPath)

    proc copyArturo(config: BuildConfig, targetFile: string) =
        log "copying files..."
        cpFile(config.binary, targetFile)

    proc giveBinaryPermission(targetFile: string) =
        exec fmt"chmod +x {targetFile}"

    proc main(config: BuildConfig) =

        if not config.shouldInstall:
            return

        section "Installing..."

        if config.webVersion:
            panic "Web builds can't be installed, please don't use --install"

        verifyDirectories()
        config.copyArturo(targetFile)

        if hostOS != "windows":
            giveBinaryPermission(targetFile)
        else:
            copyWebView()

        log fmt"deployed to: {targetDir}"

    main(config)

proc showBuildInfo*(config: BuildConfig) =
    let
        params = flags.join(" ")
        version = "version/version".staticRead()

    if config.generateBundle:
        section "Bundling..."
    else:
        section "Building..."
    log fmt"version: {version}"
    log fmt"config: {config.version}"

    if not config.silentCompilation:
        log fmt"flags: {params}"

#=======================================
# Methods
#=======================================

proc buildArturo*(config: BuildConfig, targetFile: string) =
    
    # Methods 

    proc showInfo(config: BuildConfig) =
        showEnvironment()
        config.showBuildInfo()

    proc setDevmodeUp() =
        section "Updating build..."
        updateBuild()
        devConfig()

    proc setBundlemodeUp() =
        bundleConfig()
        putEnv "BUNDLE_CONFIG", config.bundle

    proc tryCompilation(config: BuildConfig) =
        ## Panics if can't compile.
        if (let cd = config.compile(showFooter=true); cd != 0):
            panic "Compilation failed. Please try again with --log and report it.", cd

    proc main() =
        showHeader "install"

        if config.isDeveloper:
            setDevmodeUp()

        if config.generateBundle:
            setBundlemodeUp()

        config.showInfo()
        config.tryCompilation()
        config.compressBinary()

        if config.shouldInstall:
            config.installAll(targetFile)

        showFooter()

    main()

proc buildPackage*(config: BuildConfig) =

    # Helper functions

    proc dataFile(package: string): string =
        return fmt"{package}.data.json"

    proc file(package: string): string =
        return fmt"{package}.art"

    proc info(package: string): string =
        staticExec fmt"arturo --package-info {package.file}"

    # Subroutines

    proc generateData(package: string) =
        section "Processing data..."
        (package.dataFile).writeFile(package.info)
        log fmt"written to: {package.dataFile}"

    proc setEnvUp(package: string) =
        section "Setting up options..."

        putEnv "PORTABLE_INPUT", package.file
        putEnv "PORTABLE_DATA", package.dataFile

        log fmt"done!"

    proc setFlagsUp() =
        --forceBuild:on
        --opt:size
        --define:NOERRORLINES
        --define:PORTABLE

    proc showFlags() =
        let params = flags.join(" ")
        log fmt"FLAGS: {params}"
        echo ""

    proc cleanUp(package: string) =
        rmFile package.dataFile
        echo fmt"{styles.clear}"

    proc main() =
        let package = config.binary

        showHeader "package"

        package.generateData()
        package.setEnvUp()
        showEnvironment()
        config.showBuildInfo()

        setFlagsUp()
        showFlags()

        if (let cd = compile(config, showFooter=false); cd != 0):
            panic "Package building failed. Please try again with --log and report it.", cd

        package.cleanUp()

    main()


proc buildDocs*() =
    let 
        params = flags.join(" ")
        genDocs = fmt"nim doc --project --index:on --outdir:dev-docs {params} src/arturo.nim"
        genIndex = "nim buildIndex -o:dev-docs/theindex.html dev-docs"

    showHeader "docs"

    section "Generating documentation..."
    genDocs.exec()
    genIndex.exec()

proc performTests*(binary: string): bool =
    result = true

    showHeader "test"
    try:
        exec fmt"{binary} ./tools/tester.art"
    except:
        return false

proc performBenchmarks*(binary: string): bool =
    result = true

    showHeader "benchmark"
    try:
        exec fmt"{binary} ./tools/benchmarker.art"
    except:
        return false

#=======================================
# Main
#=======================================

cliInstance.header = getLogo()
cliInstance.defaultCommand = "build"
let 
    args = cliInstance.args

cmd build, "[default] Build arturo and optionally install the executable":
    ## build:
    ##     Provides a cross-compilation for the Arturo's binary.
    ##
    ##     --arch -a: $hostCPU          chooses the target CPU
    ##          [amd64, arm, arm64, i386, x86]
    ##     --as: arturo                 changes the name of the binary
    ##     --mode -m: full              chooses the target Build Version
    ##          [full, mini, web]
    ##     --os: $hostOS                chooses the target OS
    ##          [freebsd, linux, openbsd, mac, macos, macosx, netbsd, win, windows]
    ##     --profiler -p: none          defines which profiler use
    ##          [default, mem, native, none, profile]
    ##     --who: none                  defines who is compiling the code
    ##          [dev, user]
    ##     --debug -d                   enables debugging
    ##     --install -i                 installs the final binary
    ##     --log -l                     shows compilation logs
    ##     --raw                        disables compression
    ##     --release                    enable release config mode
    ##     --help

    let
        availableCPUs = @["amd-64", "x64", "x86-64", "arm-64", "i386", "x86",
                          "x86-32", "arm", "arm-32"]
        availableOSes = @["freebsd", "openbsd", "netbsd", "linux", "mac",
                          "macos", "macosx", "win", "windows",]
        availableBuilds = @["full", "mini", "safe", "web"]
        availableProfilers = @["default", "mem", "native", "profile"]

    var config = buildConfig()

    config.binary = "bin"/args.getOptionValue("as", default="arturo").toExe

    match args.getOptionValue("arch", short="a",
                              default=hostCPU,
                              into=availableCPUs):
        let
            amd64 = availableCPUs[0..2]
            arm64 = [availableCPUs[3]]
            x86 = availableCPUs[4..6]
            arm32 = availableCPUs[7..8]

        >> amd64: amd64Config()
        >> arm64: arm64Config()
        >> x86:   arm64Config()
        >> arm32: arm32Config()

    match args.getOptionValue("mode", short="m", default="full", into=availableBuilds):
        >> ["full"]:
            fullBuildConfig()
        >> ["mini"]:
            miniBuildConfig()
            config.version = "@mini"
            miniBuild()
        >> ["safe"]:
            safeBuildConfig()
            miniBuild()
        >> ["web"]:
            config.binary     = config.binary.replace(".exe", "") & ".js"
            config.version    = "@web"
            webBuildConfig()
            miniBuild()

    match args.getOptionValue("os", default=hostOS, into=availableOSes):
        let
            bsd = availableOSes[0..2]
            linux = [availableOSes[3]]
            macos = availableOSes[4..6]
            windows = availableOSes[7..8]

        >> bsd:     discard
        >> linux:   discard
        >> macos:   discard
        >> windows: discard

    match args.getOptionValue("profiler", default="none", short="p",
                              into=availableProfilers):
        >> ["default"]: profilerConfig()
        >> ["mem"]:     memProfileConfig()
        >> ["native"]:  nativeProfileConfig()
        >> ["profile"]: profileConfig()

    match args.getOptionValue("who", default="", into= @["user", "dev"]):
        >> ["user"]:
            config.isDeveloper = false
            userConfig()
        >> ["dev"]:
            config.isDeveloper = true
            devConfig()

    if args.hasFlag("bundle", "b"):
        config.generateBundle = true
        config.bundle = args.getPositionalArg(2)

    if args.hasFlag("debug", "d"):
        config.shouldCompress = false
        debugConfig()

    if args.hasFlag("install", "i"):
        config.shouldInstall = true

    if args.hasFlag("log", "l"):
        config.shouldLog = true

    if args.hasFlag("raw"):
        config.shouldCompress = false

    if args.hasFlag("release"):
        releaseConfig()

    config.buildArturo(targetDir/config.binary)

cmd package, "Package arturo app and build executable":
    ## package <pkg-name>:
    ##     Compiles packages into executables.
    ##
    ##     --arch: $hostCPU             chooses the target CPU
    ##          [amd64, arm, arm64, i386, x86]
    ##     --debug -d                   enables debugging
    ##     --help

    const availableCPUs = @["amd-64", "x64", "x86-64", "arm-64", "i386", "x86",
                          "x86-32", "arm", "arm-32"]

    var config = buildConfig()
    config.binary = args.getPositionalArg(2)

    match args.getOptionValue("arch", short="a",
                              default=hostCPU,
                              into=availableCPUs):
        let
            amd64 = availableCPUs[0..2]
            arm64 = [availableCPUs[3]]
            x86 = availableCPUs[4..6]
            arm32 = availableCPUs[7..8]

        >> amd64: amd64Config()
        >> arm64: arm64Config()
        >> x86:   arm64Config()
        >> arm32: arm32Config()

    if args.hasFlag("debug", "d"):
        config.shouldCompress = false
        debugConfig()

    config.buildPackage()

cmd docs, "Build the documentation":
    ## docs:
    ##     Builds the developer documentation
    ##
    ##     --help

    --define:DOCGEN
    buildDocs()

cmd test, "Run test suite":
    ## test:
    ##     Runs test suite
    ##
    ##     --using -u: arturo           runs with the given binary
    ##     --help

    let
        binary = args.getOptionValue("using", default="arturo", short="u").toExe
        paths: tuple = (
            local: "bin"/binary,
            global: paths.targetBin/binary
        )

    unless paths.global.performTests():
        quit paths.local.performTests().toErrorCode

cmd benchmark, "Run benchmark suite":
    ## benchmark:
    ##     Runs benchmark suite
    ## 
    ##     --using -u: arturo           runs with the given binary
    ##     --help
    
    let
        binary = args.getOptionValue("using", default="arturo", short="u").toExe
        paths: tuple = (
            local: "bin"/binary,
            global: paths.targetBin/binary
        )

    unless paths.global.performBenchmarks():
        quit paths.local.performBenchmarks().toErrorCode

helpForMissingCommand()