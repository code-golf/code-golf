import ts from 'typescript'
import * as path from 'path'

function checkTypes(wrappedProgram) {
	const options = {
		allowJs: false,
		strict: true,
		target: ts.ScriptTarget.ESNext,
		module: ts.ModuleKind.None,
		noEmitOnError: true,
		noErrorTruncation: true
	}
	
	const file = ts.createSourceFile('code.ts', wrappedProgram, {
		languageVersion: ts.ScriptTarget.ESNext
	})
	
	const libLocation = path.dirname(ts.sys.getExecutingFilePath())
	
	const compilerHost = {
		getSourceFile: (fileName, languageVersionOrOptions) => {
			if(fileName === 'code.ts') return file
			
			const text = ts.sys.readFile(fileName).toString()
			return ts.createSourceFile(fileName, text, languageVersionOrOptions, false)
		},
		getDefaultLibLocation: () => libLocation,
		getDefaultLibFileName: () => path.join(libLocation, 'lib.esnext.full.d.ts'),
		writeFile: () => {},
		getCurrentDirectory: ts.sys.getCurrentDirectory,
		getDirectories: ts.sys.getDirectories,
		getCanonicalFileName: fileName => fileName,
		getNewLine: () => ts.sys.newLine,
		fileExists: ts.sys.fileExists,
		readFile: ts.sys.readFile,
		resolveModuleNames: (moduleNames, containingFile) => moduleNames.map((moduleName) => ts.resolveModuleName(moduleName, containingFile, options, ts.sys).resolvedModule),
		useCaseSensitiveFileNames: () => false
	}
	
	const program = ts.createProgram({
		rootNames: ['code.ts'],
		options: options,
		host: compilerHost,
	})
	
	const diagnostics = ts.getPreEmitDiagnostics(program)
	if(diagnostics.length) {
		throw new Error(diagnostics
			.map(e => `T${e.code}: ${e.messageText}`)
			.join('\n'))
	}
	
	const checker = program.getTypeChecker()
	const blockScopes = []

	function searchForBlockScopes(node) {
		if(node.kind === ts.SyntaxKind.Block) {
			blockScopes.push(node)
		}
		ts.forEachChild(node, searchForBlockScopes)
	}
	ts.forEachChild(file, searchForBlockScopes)
	
	const lastBlockScope = blockScopes.at(-1)
	if(!lastBlockScope) {
		throw new Error('no last block scope')
	}

	const type = checker.getTypeAtLocation(lastBlockScope.statements[0])
	const typeString = checker.typeToString(type)

	let parsedType
	try {
		parsedType = JSON.parse(typeString)
	} catch(err) {
		throw new Error(typeString)
	}

	return parsedType
}

let input = ''
process.stdin.on('data', data => {
	input += data.toString()
})
await new Promise(r => {
	process.stdin.on('end', r)
})

const parsedInput = JSON.parse(input)

const wrappedProgram = `
type args = ${JSON.stringify(parsedInput.args)};
${parsedInput.code}
{
	type _ = output
}`

let output
try {
	output = {
		stderr: '',
		stdout: checkTypes(wrappedProgram)
	}
} catch(err) {
	output = {
		stderr: err.message ?? err.toString(),
		stdout: ''
	}
}

process.stdout.write(output.stdout)
process.stderr.write(output.stderr)