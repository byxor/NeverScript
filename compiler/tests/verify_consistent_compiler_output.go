package main

import (
    "errors"
    "fmt"
    "github.com/byxor/NeverScript/compiler"
    "github.com/pmezard/go-difflib/difflib"
    "io/ioutil"
    "log"
    "os"
    "os/exec"
    "strings"
)

/*
 * Rather than writing a gigantic suite of tests like I did the past ~4 times I've started this project,
 * I've written a short program that warns me when the compiler is producing different output than before.
 *
 * Snapshot testing is easier than maintaining a bunch of tests that check sequences of bytes.
 */

func main() {
    if err := runTests(); err != nil {
        log.Fatal(err)
    }
}

func runTests() error {
    tempDir, err := ioutil.TempDir(os.TempDir(), "neverscript-temporary-testing-tempDir")
    if err != nil { return err }
    defer os.RemoveAll(tempDir)

    err = runTest("thug2", "./neverscript.ns", tempDir+"/thug2.qb", "./thug2.roq")
    if err != nil { return err }

    err = runTest("thps4", "./neverscript.ns", tempDir+"/thps4.qb", "./thps4.roq")
    if err != nil { return err }

    return nil
}

func runTest(targetGame, nsPath, qbPath, roqPath string) error {
    compileNeverscript := func() error {
        var lexer compiler.Lexer
        var parser compiler.Parser
        var bytecodeCompiler compiler.BytecodeCompiler
        bytecodeCompiler.TargetGame = targetGame
        compiler.Compile(nsPath, qbPath, &lexer, &parser, &bytecodeCompiler)
        return nil
    }

    cleanWhitespace := func(text string) string {
        text = strings.Replace(text, "\r", "", -1)
        lines := strings.Split(text,"\n")
        for i, line := range lines {
            lines[i] = strings.TrimRight(line, " \t")
        }
        text = strings.Join(lines, "\n")
        return strings.TrimSpace(text)
    }

    decompileQbToRoq := func() (string, error) {
        roqCmd := exec.Command("roq.exe", "-d", qbPath)
        decompiledRoqBytes, err := roqCmd.Output()
        if err != nil { return "", err }

        roq := cleanWhitespace(string(decompiledRoqBytes))
        return roq, nil
    }

    loadExpectedDecompiledRoq := func() (string, error) {
        expectedDecompiledRoqBytes, err := ioutil.ReadFile(roqPath)
        if err != nil { return "", err }
        expectedDecompiledRoq := cleanWhitespace(string(expectedDecompiledRoqBytes))
        return expectedDecompiledRoq, nil
    }

    compareResults := func(expected, __actual string) (string, error) {
        banner := "--------------------------------------------------------"
        if __actual == expected {
            report := "Success: ✓"
            return report, nil
        } else {
            diff := difflib.UnifiedDiff{
                A:        difflib.SplitLines(expected),
                B:        difflib.SplitLines(__actual),
                FromFile: "Expected",
                ToFile:   "Actual",
                Context:  3,
            }
            diffString, err := difflib.GetUnifiedDiffString(diff)
            if err != nil { return "", err }
            report := banner + "\n"
            report += "WARNING: COMPILER OUTPUT CHANGED (" + targetGame + ")\n"
            report += banner + "\n"
            report += "Expected: (" + targetGame + ")\n"
            report += banner + "\n"
            report += expected + "\n"
            report += banner + "\n"
            report += "But got: (" + targetGame + ")\n"
            report += banner + "\n"
            report += __actual + "\n"
            report += banner + "\n"
            report += "Diff (" + targetGame + "):\n"
            report += banner + "\n"
            report += diffString
            return report, errors.New(report)
        }
    }

    fmt.Println("Compiling ns... (" + targetGame + ")")
    err := compileNeverscript()
    if err != nil { return err }

    fmt.Println("Decompiling roq...")
    actualRoq, err := decompileQbToRoq()
    if err != nil { return err }

    expectedRoq, err := loadExpectedDecompiledRoq()
    if err != nil { return err }

    report, err := compareResults(expectedRoq, actualRoq)
    if err != nil { return err }
    fmt.Println(report + "\n")

    return nil
}
