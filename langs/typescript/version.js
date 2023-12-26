import * as fs from 'fs/promises'

const packageJson = JSON.parse(await fs.readFile('package-lock.json'))

process.stdout.write(packageJson.packages['node_modules/typescript'].version)