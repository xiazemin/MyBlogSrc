---
title: Goimports explained
layout: post
category: golang
author: 夏泽民
---
https://nikodoko.com/posts/goimports_explained/
Go users out there are probably familiar with gofmt and it’s brother goimports (which actually uses gofmt under the hood). These are two little CLI tools (written in Go, of course), that have become core parts of many developers’ Go workflow.

Over time, I have personally grown very fond of these tools. They are undoubtedly great productivity boosters, but what I love about them is not so much their raw utility as their design. They are simple, elegant, get the job done without any fuss, and, most importantly, are usable from a shell. This means you can use them in combination with other Unix tools, like ls, cat, pipes, etc. In this regard, they are aligned with the Unix “software tools” philosophy:

using programs separately or in combination to get a job done, rather than doing it by hand, by monolithic self-sufficient subsystems, or by special-purpose, one-time programs1

Besides, making a tool available to the shell is much more “democratic” than baking it into an IDE. It is (almost) universally accessible, and really easy to integrate with any editor (including Vim, of course!).

While equivalents to gofmt can be found in other languages (like black for Python or google-java-format for Java), I cannot seem to find an equivalent to goimports. It’s has been bothering me quite a bit, and I have spent a lot of time wondering if I could replicate it for another language. But I first needed to understand how it works. And this is what today’s article is all about!
<!-- more -->
What does goimports do?
Before diving into the details of how goimports is implemented, let’s take a moment to try and think about what it does. When I began working as a software engineer, I started to read lots of source code, and realized that I had a bad tendency to dive directly into it, without taking time to think beforehand about what I was expecting to see. As a result I often got lost in small details, and after digging into 15 functions or so, there often came a time when I wondered: “what is this function trying to achieve already?“.

I have since devised a little exercise that I try to do every time before starting to dive in an unknown codebase: whenever I have a rough idea of a function’s contract (given this condition on the inputs, return such and such outputs), I take a few minutes to jolt down on paper what I feel are the main steps to fulfill it.

Of course, I am not talking about guessing each and every little step, but simply about roughly outlining the process. Something like “it should first parse the text, then do this, then…”. I have found that in spite of the very little time investment needed (I rarely spend more than 5 or 10 minutes on it), the benefits are great:

when my outline is correct (which is the case most of the time), it helps me avoid getting lost in the code, and
if I realize I got it wrong, I do not just move on but use my written notes to think about why I guessed it wrong, which helps me grow my software engineering “instinct”.
I would recommend everyone to do the same!

In the case of goimports, the contract is given a valid Go file, return the same valid Go file formatted by gofmt, with exactly the imports needed for it to compile. In order to fulfill it, a possible solution could be:

parse the file and find unresolved references (symbols not declared in that file),
list all the imports and find the unneeded ones using the list of unresolved references,
look for the packages containing the unresolved references and add those to the import list,
rectify the input’s import list, removing and adding according to the results of 2. and 3., and
pass the fixed file to gofmt, returning its result.
How does it do it?
The outline we came up with in the previous section of course omits lots of important details, like handling other files of the same package or the import identifiers (the m in import m "math"), etc. Still, it is fairly accurate.

Parsing a Go file to find unresolved references
The Go parser, like many other programming language parsers out there, generates an AST out of a source file. The difference is that Go’s parser is smart enough to handle scopes as it builds the AST, resulting in an ability to tell which identifier (variable, function…) refers to something declared in the file and which does not (an unresolved reference). This analysis in itself is not trivial, and may be a good topic for another article, but for now we will say that:

There is a parse function that takes a Go source file and returns all unresolved references, along with the list of packages imported by this file and its top level declarations.

The question then becomes: how do we find which symbols need importing, and what are the right packages to import?

Finding imports
While the goimports command is defined in golang.org/x/tools/cmd/goimports/goimports.go, it is in fact just a CLI wrapper around the functionality contained in golang.org/x/tools/internal/imports.

The keystone of goimports is a structure called pass, which looks like this:

type pass struct {
	// Inputs. These must be set before a call to load, and not modified after.
	f                    *ast.File      // the file being fixed.
	srcDir               string         // the directory containing f.

	// Intermediate state, generated by load.
	existingImports map[string]*ImportInfo
	allRefs         references
	missingRefs     references

	// Inputs to fix. These can be augmented between successive fix calls.
	candidates    []*ImportInfo           // candidate imports in priority order.

  // Other fields omitted for now ...
}
existingImports’s map key is the name under which a package is imported, with ImportInfo being pretty straightforward:

// An ImportInfo represents a single import statement.
type ImportInfo struct {
	ImportPath string // import path, e.g. "crypto/rand".
	Name       string // import name, e.g. "crand", or "" if none.
}
references is defined as map[string]map[string]bool, where the first map key is the package name, and the second map key is the imported symbol. This means that something like fmt.Println would be stored in {"fmt": {"Println": true}}

pass has two core methods, named load and fix:

load determines which references are missing, and collects import candidates from various sources,
fix scans the candidates and looks for the right packages to import in order to cover as many missing references as possible.
So far, given a file f to fix, we can:

create a pass with f and it’s source directory,
use the parse function we talked about earlier to create the allRefs and the existingImports, and use these to decide which references belong to missingRefs.
However, we’ve only covered about half of the load function! The only thing we can do at this point is decide whether the file is complete or not (i.e. if existingImports covers allRefs so that there is no missingRefs). If it is, all is well and we can exit. But what if it’s not? How do we collect import candidates and generate the fixes?

Simply using the file that needs fixing does not give us enough information: if it is incomplete, we obviously need to find the packages to import (the candidates) somewhere else. Besides, at this point, missingRefs might contain false positives, as we only used the information contained in the file to decide if a reference is unresolved or not. We cannot know if it comes from the same package (and does not need to be imported), or from a different one (and needs to be imported).

And so this is where pass got it’s name from: goimports will do successive runs of load and fix, adding more and more information as it progresses.

First pass: one file only
This is the pass we talked about: we simply scan the file to fix, and determine if there is any missingRefs using only the imports at the top of the file. If nothing is found (and the file is complete) then goimports exits without doing anything. Otherwise, it moves on to the second pass.

As we said just before, this pass does not enable us to add any fixes, and also generates false positives (some missingRefs are in fact not missing). So why bother? You might think that it’s wasteful to call parse just for this, and you would not be wrong. But, as I will explain later, fetching additional information can be comparatively expensive, so this first pass is just a chance to exit as soon as possible with minimal effort.

Second pass: sibling files
Time to reveal a little more about pass:

type pass struct {
	// Inputs. These must be set before a call to load, and not modified after.
	fset                 *token.FileSet // fset used to parse f and its siblings.
	f                    *ast.File      // the file being fixed.
	srcDir               string         // the directory containing f.
	otherFiles           []*ast.File    // sibling files.

	// Intermediate state, generated by load.
	existingImports map[string]*ImportInfo
	allRefs         references
	missingRefs     references

	// Inputs to fix. These can be augmented between successive fix calls.
	candidates    []*ImportInfo           // candidate imports in priority order.

  // Other fields omitted for now ...
}
goimports will use srcDir to load all Go files in the directory of the file to fix, and parse all of them: if using the recommended Go architecture, this should yield all the package files. Then, when going over allRefs to determine the missingRefs, it will not only use existingImports (the import declarations on top of the file to fix), but also all top declarations of the sibling files. With this, we can say for sure that missingRefs contains only the references that need to be imported (as opposed to the first pass, where the only thing we could say was that it contained references that need to be imported and references found in the package).

It is now time to find some candidates! For this purpose, Go uses a clever trick: while parsing the files of the same package, it keeps track of all their imports and adds them to the candidates list. Then, if any missingRefs can be satisfied by one of these candidates (meaning that it has the same package name and contains all imported symbols), the candidate will be added to the list of fixes.

If all missingRefs are found, then fix returns the fixes (a list of ImportedInfo), else it moves on to the next pass.

Third pass: adding the standard library
goimports in fact contains a fixed map[string][]string with all the standard library packages and their exported symbols (like archive/zip, time, etc.), stored in a file called zstdlib.go. This map is generated by a script (mkstdlib.go), using all goX.Y.txt stored in GOROOT/api. These files are added at each release and contain all new packages and their exported symbols. For example, go1.13.txt starts by:

pkg bytes, func ToValidUTF8([]uint8, []uint8) []uint8
pkg crypto/ed25519, const PrivateKeySize = 64
pkg crypto/ed25519, const PrivateKeySize ideal-int
pkg crypto/ed25519, const PublicKeySize = 32

etc.
The third pass is identical to the second one, except all the contents from the standard library are added as candidates in addition to the imports of sibling files.

Fourth pass: external packages
If it is still not enough, goimports attempts a last pass, using the environment (GOROOT, GOPATH etc.) to parse all external packages that can be found. I will not go into details, as this step is fairly complex, and simply mention that it uses a “distance” system to sort packages with similar names (it essentially assumes that the sorter import path is the best), and manages to stay pretty fast by making heavy usage of goroutines (using all CPU cores available).

Once external candidates are found this way, they are added to the candidates list and the process is identical to the second and third pass. This fourth path is obviously the most expensive one, that we try to avoid.

Because there is no next pass, the fourth one returns all the fixes it can, even if that list is incomplete.

Conclusion
Once we’ve built a list of ImportInfo containing all fixes, the last step is to apply them. This is essentially text formatting, and while interesting in its own right, I will not cover it in this article.

That’s it! I hope this article gives you a good enough view of goimports’s internals. I’ve gone into a decent amount of details, but if you ever decide to read the source yourself you’ll quickly realize I’ve skipped some steps, simplified others and outright ignored several subtleties (especially in external package resolving). If you wish to know more, I hope this article can become your hitchhiker’s guide, helping you to not get lost in the source code!

As always, shoot me a message or tweet @nicol4s_c if you want to chat about any of this, if you spotted any mistakes or typos, or if you’d like me to cover anything else! Have a great day :)

Brian Kernighan and Rob Pike, “Program Design in the UNIX Environment,” in AT&T Bell Labs Technical Journal, October, 1984, p. 1596. [return]