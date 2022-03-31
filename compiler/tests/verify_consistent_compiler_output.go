package main

import (
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
 * I'm just going to write a short program that warns me when the compiler is producing different output than before.
 *
 * I suspect that the change-detection will be easier to work with than a bunch of ever-changing tests.
 * Even if a bit of human effort is required when reviewing the output,
 * it still saves vastly more time than updating a very granular set of "this code produces this sequence of bytes"-type tests.
 */

func main() {
    tempDir, err := ioutil.TempDir(os.TempDir(), "neverscript-temporary-testing-tempDir")
    if err != nil {
        log.Fatal(err)
    }
    defer os.RemoveAll(tempDir)

    nsPath := tempDir + "/code.ns"
    qbPath := tempDir + "/code.qb"

    // Create neverscript file
    ioutil.WriteFile(nsPath, []byte(code), 0644)

    // Compile
    fmt.Println("Compiling ns...")
    var lexer compiler.Lexer
    var parser compiler.Parser
    var bytecodeCompiler compiler.BytecodeCompiler
    bytecodeCompiler.TargetGame = "thug2"
    compiler.Compile(nsPath, qbPath, &lexer, &parser, &bytecodeCompiler)
    fmt.Println()

    // Decompile
    fmt.Println("Decompiling roq...")
    roqCmd := exec.Command("roq.exe", "-d", qbPath)
    decompiledRoq, err := roqCmd.Output()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println()

    // Compare results (the stupid variable name aligns the output in the debugger)
    __actual := strings.TrimSpace(strings.Replace(string(decompiledRoq), "\r", "", -1))
    expected := strings.TrimSpace(strings.Replace(expectedDecompiledRoq, "\r", "", -1))
    banner := "--------------------------------------------------------"
    if __actual != expected {
        fmt.Println(banner)
        fmt.Println("WARNING: COMPILER OUTPUT CHANGED")
        fmt.Println(banner)
        fmt.Println("Expected:")
        fmt.Println(banner)
        fmt.Println(expected)
        fmt.Println(banner)
        fmt.Println("But got:")
        fmt.Println(banner)
        fmt.Println(__actual)
        fmt.Println(banner)
        fmt.Println("Diff:")
        fmt.Println(banner)

        diff := difflib.UnifiedDiff{
            A:        difflib.SplitLines(expected),
            B:        difflib.SplitLines(__actual),
            FromFile: "Expected",
            ToFile:   "Actual",
            Context:  3,
        }
        text, _ := difflib.GetUnifiedDiffString(diff)
        fmt.Println(text)
    } else {
        fmt.Println(banner)
        fmt.Println("Hooray! Still working!")
    }

    // fmt.Printf("Hex dump:\n%s", qbPath, hex.Dump(bytecodeCompiler.Bytes))
}










// 'code' is on line 101 so the line numbers are easy to see when debugging
// (just add 100 to the line number shown by the debugger/compiler)
const code = `
my_int = 10
my_int = -10
my_float = 0.1
my_float = -0.1
my_string = "hey"
my_pair = (1.0, 2.0)
my_vector = (100.0, 200.0, 300.0)
my_array = [1, 2, 3]
my_struct = { x=1, y=2, z=3 }
my_checksum = #deadf00d
my_checksum = identifiers_are_checksums_too
my_checksum =`+" `checksums between backticks can have spaces` "+`

// my comment
x = 10 // comment after assignment

my_struct = {
    1 // one
    2 // two
    3 /* three */
}

/* my comment */
x = 10 /* comment
after assignment */

/*big multiline comment
script 7 2 1 1 2 3 4 
"will not compile
zzzz = - ~ ¬!!()&&&^^{}{ { { (<{ / / / *
*/

/* multiline
/* with */
/* /* nested */ */
/* comments */ */

script TestBasicExpressions {
    x = 1
    x = (1)
    
    description = "Positive ints:"
    x = (1 + 2)
    x = (1 - 2)
    x = (1 * 3)
    x = (1 / 2)
    
    description = "Negative ints:"
    x = (-1 + -2)
    x = (-1 - -2)
    x = (-1 * -3)
    x = (-1 / -2)
    
    description = "Positive floats:"
    x = (1.0 + 2.0)
    x = (1.0 - 2.0)
    x = (1.0 * 3.0)
    x = (1.0 / 2.0)

    description = "Negative floats:"
    x = (-1.0 + -2.0)
    x = (-1.0 - -2.0)
    x = (-1.0 * -3.0)
    x = (-1.0 / -2.0)
}

script TestStringEscaping {
    description = "String containing a single backslash:"
    s = "\\"

    description = "String containing two backslashes:"
    s = "\\\\"
    
    description = "String containing a quote:"
    s = "\""
}

script TestShorthandMath {
    description = "Global variables:"
    Change x += 5
    Change x -= 6
    Change x *= 7
    Change x *= 8

    description = "Local variables:"
    <x> += 5
    <x> -= 6
    <x> *= 7
    <x> /= 8
}

script TestInvocations {
    description = "Invocation with checksum parameters:"
    NameOfScript param1 param2 param3

    description = "Invocation with assigned parameters:"
    NameOfScript param1=1 param2=2.0 param3="3" param4=(4.0, 0.4) param5=[5 5.0 "5" five "five"] param6={six=6}

    description = "Invocation across multiple lines:"
    NameOfScript param1 = 1 \
                 param2 = 2.0 \
                 param3 = "3"

    description = "Invocation across multiple lines (1st param on next line):"
    NameOfScript \
        param1 = 1 \
        param2 = 2.0 \
        param3 = "3"
}

script TestIfStatements {
    description = "Basic if:"
    if something {}

    description = "Basic if/else:"
    if something {} else {}

    description = "Basic if/elseif/else:"
    if c1 {} else if c2 {} else {}

    description = "Condition with logical not:"
    if ! condition {}

    description = "Condition with logical and:"
    if c1 and c2 {}

    description = "Condition with invocation:"
    if GotParam Foo {}

    description = "Condition with invocation with struct parameter:"
    if IsOld {name="byxor", age=23} {
        MakeYounger
    }

    description = "Condition with logical not with invocation with struct parameter:"
    if ! IsFinished {progress=10, finish=100} {
        MakeProgress
    }

    description = "Condition with member function invocation:"
    if Object:GetCollision {
        PlayCollisionSound
    }

    description = "Condition with member function invocation with struct parameter:"
    if Object:GetCollision { length=20 } {
        PlayCollisionSound
    }

    // TODO(brandon): Add more variations of invocations with final struct params here
    // e.g. ! Object:MemberFunction {distance=15000000}
    // e.g. & struct.script {distance=15000000}

    description = "Comparisons:"
    if (c1 = c2) {}
    if (c1 < c2) {}
    if (c1 > c2) {}
    if (c1 != c2) {}
    if (c1 <= c2) {}
    if (c1 >= c2) {}
}

script TestEmptyReturn {
    return
}

script TestReturningMultipleParametersOnSingleLine {
    return x=1 y=2 z=3 w={what="the", heckIsHeDoingHere}
}

script TestReturningMultipleParametersOnMultipleLines {
    return \
        x = 11 \
        y = 22 \
        z = 33
}

script TestWhile {
    while {
        Tick
        Tock
    }
}

script TestNestedWhile {
    while {
        while {
            // should have different variable names when bypassing infinite loop checks.
        }
    }
}

script TestRandom {
    random {
        10 {
            print "this is gonna happen 10/15 times on average"
            print "yo yo"
        }
        5 {
            print "this is gonna happen 5/15 times on average"
            print "skrrrrrt"
        }
    }

    x = random {
        9 { "Hey" }
        4 { "Hello" }
        10 { "Yo" }
        2 { "What's up?" }
    }
}

script TestIdentifiersWithKeywordPrefixes {
    script_1 = ""
    or_1 = ""
    if_1 = ""
    and_1 = ""
    else_1 = ""
    while_1 = ""
    break_1 = ""
    random_1 = ""
    return_1 = ""
}

/*
TODO(brandon): Add tests for this, it won't compile
script TestFloatsToVector {
    return vector = (
        ((1.0 * <x>) * (1.0, 0.0, 0.0)) +
        ((1.0 * <y>) * (0.0, 1.0, 0.0)) +
        ((1.0 * <z>) * (0.0, 0.0, 1.0))
    )
}
*/
`

const expectedDecompiledRoq = `
:i $my_int$ = %i(10,0000000a)
:i $my_int$ = %i(4294967286,fffffff6)
:i $my_float$ = %f(0.100000)
:i $my_float$ = %f(-0.100000)
:i $my_string$ = %s(3,"hey")
:i $my_pair$ = %vec2(1.000000,2.000000)
:i $my_vector$ = %vec3(100.000000,200.000000,300.000000)
:i $my_array$ = :a{%i(1,00000001);%i(2,00000002);%i(3,00000003):a}
:i $my_struct$ = :s{$x$ = %i(1,00000001);$y$ = %i(2,00000002);$z$ = %i(3,00000003):s}
:i $my_checksum$ = $[0df0adde]$
:i $my_checksum$ = $identifiers_are_checksums_too$
:i $my_checksum$ = $checksums between backticks can have spaces$
:i $x$ = %i(10,0000000a)
:i $my_struct$ = :s{
	:i %i(1,00000001)
	:i %i(2,00000002)
	:i %i(3,00000003)
:i :s}
:i $x$ = %i(10,0000000a)
:i function $TestBasicExpressions$
	:i $x$ = %i(1,00000001)
	:i $x$ =  (%i(1,00000001)) 
	:i $description$ = %s(14,"Positive ints:")
	:i $x$ =  (%i(1,00000001) + %i(2,00000002)) 
	:i $x$ =  (%i(1,00000001) - %i(2,00000002)) 
	:i $x$ =  (%i(1,00000001) * %i(3,00000003)) 
	:i $x$ =  (%i(1,00000001) / %i(2,00000002)) 
	:i $description$ = %s(14,"Negative ints:")
	:i $x$ =  (%i(4294967295,ffffffff) + %i(4294967294,fffffffe)) 
	:i $x$ =  (%i(4294967295,ffffffff) - %i(4294967294,fffffffe)) 
	:i $x$ =  (%i(4294967295,ffffffff) * %i(4294967293,fffffffd)) 
	:i $x$ =  (%i(4294967295,ffffffff) / %i(4294967294,fffffffe)) 
	:i $description$ = %s(16,"Positive floats:")
	:i $x$ =  (%f(1.000000) + %f(2.000000)) 
	:i $x$ =  (%f(1.000000) - %f(2.000000)) 
	:i $x$ =  (%f(1.000000) * %f(3.000000)) 
	:i $x$ =  (%f(1.000000) / %f(2.000000)) 
	:i $description$ = %s(16,"Negative floats:")
	:i $x$ =  (%f(-1.000000) + %f(-2.000000)) 
	:i $x$ =  (%f(-1.000000) - %f(-2.000000)) 
	:i $x$ =  (%f(-1.000000) * %f(-3.000000)) 
	:i $x$ =  (%f(-1.000000) / %f(-2.000000)) 
:i endfunction
:i function $TestStringEscaping$
	:i $description$ = %s(37,"String containing a single backslash:")
	:i $s$ = %s(1,"\")
	:i $description$ = %s(34,"String containing two backslashes:")
	:i $s$ = %s(2,"\\")
	:i $description$ = %s(26,"String containing a quote:")
	:i $s$ = %s(1,""")
:i endfunction
:i function $TestShorthandMath$
	:i $description$ = %s(17,"Global variables:")
	:i $Change$$x$ =  ($x$ + %i(5,00000005)) 
	:i $Change$$x$ =  ($x$ - %i(6,00000006)) 
	:i $Change$$x$ =  ($x$ * %i(7,00000007)) 
	:i $Change$$x$ =  ($x$ * %i(8,00000008)) 
	:i $description$ = %s(16,"Local variables:")
	:i %GLOBAL%$x$ =  (%GLOBAL%$x$ + %i(5,00000005)) 
	:i %GLOBAL%$x$ =  (%GLOBAL%$x$ - %i(6,00000006)) 
	:i %GLOBAL%$x$ =  (%GLOBAL%$x$ * %i(7,00000007)) 
	:i %GLOBAL%$x$ =  (%GLOBAL%$x$ / %i(8,00000008)) 
:i endfunction
:i function $TestInvocations$
	:i $description$ = %s(36,"Invocation with checksum parameters:")
	:i $NameOfScript$$param1$$param2$$param3$
	:i $description$ = %s(36,"Invocation with assigned parameters:")
	:i $NameOfScript$$param1$ = %i(1,00000001)$param2$ = %f(2.000000)$param3$ = %s(1,"3")$param4$ = %vec2(4.000000,0.400000)$param5$ = :a{%i(5,00000005)%f(5.000000)%s(1,"5")$five$%s(4,"five"):a}$param6$ = :s{$six$ = %i(6,00000006):s}
	:i $description$ = %s(33,"Invocation across multiple lines:")
	:i $NameOfScript$$param1$ = %i(1,00000001)$param2$ = %f(2.000000)$param3$ = %s(1,"3")
	:i $description$ = %s(58,"Invocation across multiple lines (1st param on next line):")
	:i $NameOfScript$$param1$ = %i(1,00000001)$param2$ = %f(2.000000)$param3$ = %s(1,"3")
:i endfunction
:i function $TestIfStatements$
	:i $description$ = %s(9,"Basic if:")
	:i if $something$endif
	:i $description$ = %s(14,"Basic if/else:")
	:i if $something$else endif
	:i $description$ = %s(21,"Basic if/elseif/else:")
	:i if $c1$else 
		:i if $c2$else endif
	:i endif
	:i $description$ = %s(27,"Condition with logical not:")
	:i if NOT $condition$endif
	:i $description$ = %s(27,"Condition with logical and:")
	:i if $c1$ AND $c2$endif
	:i $description$ = %s(26,"Condition with invocation:")
	:i if $GotParam$$Foo$endif
	:i $description$ = %s(48,"Condition with invocation with struct parameter:")
	:i if $IsOld$:s{$name$ = %s(5,"byxor");$age$ = %i(23,00000017):s}
		:i $MakeYounger$
	:i endif
	:i $description$ = %s(65,"Condition with logical not with invocation with struct parameter:")
	:i if NOT $IsFinished$:s{$progress$ = %i(10,0000000a);$finish$ = %i(100,00000064):s}
		:i $MakeProgress$
	:i endif
	:i $description$ = %s(42,"Condition with member function invocation:")
	:i if $Object$.$GetCollision$
		:i $PlayCollisionSound$
	:i endif
	:i $description$ = %s(64,"Condition with member function invocation with struct parameter:")
	:i if $Object$.$GetCollision$:s{$length$ = %i(20,00000014):s}
		:i $PlayCollisionSound$
	:i endif
	:i $description$ = %s(12,"Comparisons:")
	:i if  ($c1$ = $c2$) endif
	:i if  ($c1$ < $c2$) endif
	:i if  ($c1$ > $c2$) endif
	:i if NOT  ($c1$ = $c2$) endif
	:i if NOT  ($c1$ > $c2$) endif
	:i if NOT  ($c1$ < $c2$) endif
:i endfunction
:i function $TestEmptyReturn$
	:i return
	
:i endfunction
:i function $TestReturningMultipleParametersOnSingleLine$
	:i return
	$x$ = %i(1,00000001)$y$ = %i(2,00000002)$z$ = %i(3,00000003)$w$ = :s{$what$ = %s(3,"the");$heckIsHeDoingHere$:s}
:i endfunction
:i function $TestReturningMultipleParametersOnMultipleLines$
	:i return
	$x$ = %i(11,0000000b)$y$ = %i(22,00000016)$z$ = %i(33,00000021)
:i endfunction
:i function $TestWhile$
	:i $__COMPILER__infinite_loop_bypasser_0$ = %i(0,00000000)
	:i while
		
		:i if  (%GLOBAL%$__COMPILER__infinite_loop_bypasser_0$ > %i(0,00000000)) 
			:i continue
			
		:i endif
		:i $Tick$
		:i $Tock$
	:i loop_to 
:i endfunction
:i function $TestNestedWhile$
	:i $__COMPILER__infinite_loop_bypasser_1$ = %i(0,00000000)
	:i while
		
		:i if  (%GLOBAL%$__COMPILER__infinite_loop_bypasser_1$ > %i(0,00000000)) 
			:i continue
			
		:i endif
		:i $__COMPILER__infinite_loop_bypasser_2$ = %i(0,00000000)
		:i while
			
			:i if  (%GLOBAL%$__COMPILER__infinite_loop_bypasser_2$ > %i(0,00000000)) 
				:i continue
				
			:i endif
		:i loop_to 
	:i loop_to 
:i endfunction
:i function $TestRandom$
	:i select(2f,2, 0a 00 05 00) :OFFSET(0):OFFSET(1)
		 :POS(0) 
		:i $print$%s(43,"this is gonna happen 10/15 times on average")
		:i $print$%s(5,"yo yo")
		:i 
	:BREAKTO(2)
		 :POS(1) 
		:i $print$%s(42,"this is gonna happen 5/15 times on average")
		:i $print$%s(8,"skrrrrrt")
		:i  :POS(2) 
	:i $x$ = select(2f,4, 09 00 04 00 0a 00 02 00) :OFFSET(3):OFFSET(4):OFFSET(5):OFFSET(6)
		 :POS(3) %s(3,"Hey")
	:BREAKTO(7)
		 :POS(4) %s(5,"Hello")
	:BREAKTO(7)
		 :POS(5) %s(2,"Yo")
	:BREAKTO(7)
		 :POS(6) %s(10,"What's up?") :POS(7) 
	:i endfunction
	:i function $TestIdentifiersWithKeywordPrefixes$
		:i $script_1$ = %s(0,"")
		:i $or_1$ = %s(0,"")
		:i $if_1$ = %s(0,"")
		:i $and_1$ = %s(0,"")
		:i $else_1$ = %s(0,"")
		:i $while_1$ = %s(0,"")
		:i $break_1$ = %s(0,"")
		:i $random_1$ = %s(0,"")
		:i $return_1$ = %s(0,"")
	:i endfunction
	:i :end
`