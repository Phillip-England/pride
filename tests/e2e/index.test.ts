import { expect, test } from "bun:test"
import { $ } from "bun"
import { rm } from "node:fs/promises"


const PridePath = '/home/jacex/src/pride/main.go'
const TestSitePath = '/home/jacex/src/pride/tmp/test-site'

async function clean() {
    await rm(TestSitePath, {recursive: true, force: true})
}

async function cmdNew(flag: string, dest: string) {
    let result = await $`go run ${PridePath} new ${flag} ${dest}`
    return result
}

async function cmdHelp() {
    let result = await $`go run ${PridePath} help`
    return result
}

test('pride new', async () => {
    await clean()
    let result = await cmdNew('site', TestSitePath)
    result = await cmdNew('site', TestSitePath)
    let stderr = result.stderr.toString()
    console.log(stderr)
})



